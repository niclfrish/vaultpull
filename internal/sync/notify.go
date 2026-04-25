package sync

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// NotifyLevel represents the severity of a notification.
type NotifyLevel string

const (
	NotifyInfo  NotifyLevel = "INFO"
	NotifyWarn  NotifyLevel = "WARN"
	NotifyError NotifyLevel = "ERROR"
)

// NotifyEvent holds the data for a single notification.
type NotifyEvent struct {
	Level     NotifyLevel
	Message   string
	Namespace string
	Timestamp time.Time
	Meta      map[string]string
}

// Notifier dispatches sync lifecycle events to one or more sinks.
type Notifier struct {
	sinks []NotifySink
}

// NotifySink is implemented by anything that can receive a NotifyEvent.
type NotifySink interface {
	Send(event NotifyEvent) error
}

// NewNotifier creates a Notifier with the provided sinks.
func NewNotifier(sinks ...NotifySink) *Notifier {
	return &Notifier{sinks: sinks}
}

// Emit sends an event to all registered sinks, collecting errors.
func (n *Notifier) Emit(level NotifyLevel, namespace, message string, meta map[string]string) []error {
	event := NotifyEvent{
		Level:     level,
		Message:   message,
		Namespace: namespace,
		Timestamp: time.Now().UTC(),
		Meta:      meta,
	}
	var errs []error
	for _, s := range n.sinks {
		if err := s.Send(event); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// WriterSink writes human-readable notifications to an io.Writer.
type WriterSink struct {
	w io.Writer
}

// NewWriterSink creates a WriterSink; falls back to os.Stdout when w is nil.
func NewWriterSink(w io.Writer) *WriterSink {
	if w == nil {
		w = os.Stdout
	}
	return &WriterSink{w: w}
}

// Send formats and writes the event.
func (ws *WriterSink) Send(event NotifyEvent) error {
	ns := event.Namespace
	if ns == "" {
		ns = "default"
	}
	var extras []string
	for k, v := range event.Meta {
		extras = append(extras, fmt.Sprintf("%s=%s", k, v))
	}
	extra := ""
	if len(extras) > 0 {
		extra = " " + strings.Join(extras, " ")
	}
	_, err := fmt.Fprintf(ws.w, "[%s] %s ns=%s%s %s\n",
		event.Timestamp.Format(time.RFC3339),
		event.Level, ns, extra, event.Message)
	return err
}
