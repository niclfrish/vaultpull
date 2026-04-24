package sync

import (
	"fmt"
	"io"
	"os"
)

// RenderTemplate applies a TemplateRenderer to the provided secrets and writes
// the rendered result to w under outputKey. If w is nil, os.Stdout is used.
// The original secrets map is returned unchanged so the hook is non-destructive.
func RenderTemplate(r *TemplateRenderer, outputKey string, w io.Writer) func(map[string]string) (map[string]string, error) {
	if w == nil {
		w = os.Stdout
	}
	return func(secrets map[string]string) (map[string]string, error) {
		if r == nil {
			return secrets, nil
		}
		out, err := r.Render(secrets)
		if err != nil {
			return nil, fmt.Errorf("template hook: %w", err)
		}
		_, werr := fmt.Fprintln(w, out)
		if werr != nil {
			return nil, fmt.Errorf("template hook write: %w", werr)
		}
		return secrets, nil
	}
}

// TemplateStage returns a pipeline Stage that renders secrets through the
// given TemplateRenderer and stores the result under outputKey in the map.
func TemplateStage(r *TemplateRenderer, outputKey string) Stage {
	return Stage{
		Name: "template",
		Fn: func(secrets map[string]string) (map[string]string, error) {
			if r == nil {
				return secrets, nil
			}
			result, err := r.RenderToMap(secrets, outputKey)
			if err != nil {
				return nil, err
			}
			// Merge rendered key back into the secrets map.
			merged := make(map[string]string, len(secrets)+1)
			for k, v := range secrets {
				merged[k] = v
			}
			for k, v := range result {
				merged[k] = v
			}
			return merged, nil
		},
	}
}
