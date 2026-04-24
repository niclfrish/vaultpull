package sync

import (
	"bytes"
	"fmt"
	"text/template"
)

// TemplateRenderer renders secrets into a user-defined template string.
type TemplateRenderer struct {
	tmpl *template.Template
}

// NewTemplateRenderer parses the given template text and returns a TemplateRenderer.
// Returns an error if the template text is invalid.
func NewTemplateRenderer(text string) (*TemplateRenderer, error) {
	if text == "" {
		return nil, fmt.Errorf("template text must not be empty")
	}
	tmpl, err := template.New("secrets").Option("missingkey=error").Parse(text)
	if err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}
	return &TemplateRenderer{tmpl: tmpl}, nil
}

// Render executes the template with the provided secrets map and returns the
// rendered output as a string.
func (r *TemplateRenderer) Render(secrets map[string]string) (string, error) {
	if secrets == nil {
		secrets = map[string]string{}
	}
	var buf bytes.Buffer
	if err := r.tmpl.Execute(&buf, secrets); err != nil {
		return "", fmt.Errorf("render template: %w", err)
	}
	return buf.String(), nil
}

// RenderToMap renders the template and returns a single-key map under the
// provided output key, suitable for passing to downstream pipeline stages.
func (r *TemplateRenderer) RenderToMap(secrets map[string]string, outputKey string) (map[string]string, error) {
	if outputKey == "" {
		return nil, fmt.Errorf("outputKey must not be empty")
	}
	out, err := r.Render(secrets)
	if err != nil {
		return nil, err
	}
	return map[string]string{outputKey: out}, nil
}
