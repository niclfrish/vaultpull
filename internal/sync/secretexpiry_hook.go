package sync

import (
	"fmt"
	"io"
	"os"
	"time"
)

// ExpiryCheckResult holds the outcome of an expiry check pass.
type ExpiryCheckResult struct {
	Infos   []ExpiryInfo
	Expired []ExpiryInfo
}

// CheckExpiry inspects secrets for expiry metadata and writes a report to w.
// It returns an error if any secrets are expired and failOnExpired is true.
func CheckExpiry(secrets map[string]string, failOnExpired bool, w io.Writer) (*ExpiryCheckResult, error) {
	if w == nil {
		w = os.Stdout
	}
	now := time.Now().UTC()
	infos := ClassifyExpiry(secrets, now)
	result := &ExpiryCheckResult{Infos: infos}

	for _, info := range infos {
		if info.Expired {
			result.Expired = append(result.Expired, info)
			fmt.Fprintf(w, "[expiry] EXPIRED: %s (expired at %s)\n",
				info.Key, info.ExpiresAt.Format(time.RFC3339))
		} else {
			fmt.Fprintf(w, "[expiry] OK: %s (expires in %s)\n",
				info.Key, info.TTL.Round(time.Second))
		}
	}

	fmt.Fprintln(w, "[expiry]", ExpirySummary(infos))

	if failOnExpired && len(result.Expired) > 0 {
		return result, fmt.Errorf("%d secret(s) have expired", len(result.Expired))
	}
	return result, nil
}

// ExpiryFilterStage returns a PipelineStage that removes expired secrets
// (those whose __expires_at__ companion key indicates past expiry).
func ExpiryFilterStage() PipelineStage {
	return PipelineStage{
		Name: "expiry-filter",
		Fn: func(secrets map[string]string) (map[string]string, error) {
			now := time.Now().UTC()
			infos := ClassifyExpiry(secrets, now)
			expiredKeys := make(map[string]bool, len(infos))
			for _, i := range infos {
				if i.Expired {
					expiredKeys[i.Key] = true
				}
			}
			out := make(map[string]string, len(secrets))
			for k, v := range secrets {
				if !expiredKeys[k] {
					out[k] = v
				}
			}
			return out, nil
		},
	}
}
