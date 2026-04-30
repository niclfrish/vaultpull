package sync

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// WatchRule defines a condition to monitor on a secret key.
type WatchRule struct {
	// Key is the secret key to watch.
	Key string
	// Pattern is an optional substring the value must contain.
	Pattern string
	// AbsentOK means no alert is raised if the key is missing.
	AbsentOK bool
}

// WatchAlert describes a triggered watch rule.
type WatchAlert struct {
	Key       string
	Rule      WatchRule
	Message   string
	Triggered time.Time
}

// DefaultSecretWatchConfig returns a WatchConfig with sensible defaults.
func DefaultSecretWatchConfig() SecretWatchConfig {
	return SecretWatchConfig{
		AlertOnMissing:  true,
		AlertOnEmpty:    true,
		AlertOnPattern:  true,
	}
}

// SecretWatchConfig controls which conditions trigger an alert.
type SecretWatchConfig struct {
	// AlertOnMissing raises an alert when a watched key is absent from secrets.
	AlertOnMissing bool
	// AlertOnEmpty raises an alert when a watched key has an empty value.
	AlertOnEmpty bool
	// AlertOnPattern raises an alert when a value does not match the rule pattern.
	AlertOnPattern bool
}

// WatchSecrets evaluates rules against secrets and returns any triggered alerts.
// It never modifies the secrets map.
func WatchSecrets(secrets map[string]string, rules []WatchRule, cfg SecretWatchConfig) ([]WatchAlert, error) {
	if secrets == nil {
		return nil, fmt.Errorf("watchsecrets: secrets map is nil")
	}
	if len(rules) == 0 {
		return nil, nil
	}

	now := time.Now().UTC()
	var alerts []WatchAlert

	for _, rule := range rules {
		if strings.TrimSpace(rule.Key) == "" {
			continue
		}

		val, exists := secrets[rule.Key]

		if !exists {
			if !rule.AbsentOK && cfg.AlertOnMissing {
				alerts = append(alerts, WatchAlert{
					Key:       rule.Key,
					Rule:      rule,
					Message:   fmt.Sprintf("key %q is missing from secrets", rule.Key),
					Triggered: now,
				})
			}
			continue
		}

		if cfg.AlertOnEmpty && val == "" {
			alerts = append(alerts, WatchAlert{
				Key:       rule.Key,
				Rule:      rule,
				Message:   fmt.Sprintf("key %q has an empty value", rule.Key),
				Triggered: now,
			})
			continue
		}

		if cfg.AlertOnPattern && rule.Pattern != "" && !strings.Contains(val, rule.Pattern) {
			alerts = append(alerts, WatchAlert{
				Key:       rule.Key,
				Rule:      rule,
				Message:   fmt.Sprintf("key %q value does not contain expected pattern %q", rule.Key, rule.Pattern),
				Triggered: now,
			})
		}
	}

	// Sort alerts deterministically by key then message.
	sort.Slice(alerts, func(i, j int) bool {
		if alerts[i].Key != alerts[j].Key {
			return alerts[i].Key < alerts[j].Key
		}
		return alerts[i].Message < alerts[j].Message
	})

	return alerts, nil
}

// WatchSummary returns a human-readable summary of the alerts.
func WatchSummary(alerts []WatchAlert) string {
	if len(alerts) == 0 {
		return "watch: no alerts triggered"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "watch: %d alert(s) triggered\n", len(alerts))
	for _, a := range alerts {
		fmt.Fprintf(&sb, "  [%s] %s\n", a.Triggered.Format(time.RFC3339), a.Message)
	}
	return strings.TrimRight(sb.String(), "\n")
}
