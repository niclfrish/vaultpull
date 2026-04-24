package sync

import "fmt"

// Stage represents a named processing step in the pipeline.
type Stage struct {
	Name string
	Fn   func(secrets map[string]string) (map[string]string, error)
}

// Pipeline chains multiple transformation/validation stages together.
type Pipeline struct {
	stages []Stage
}

// NewPipeline creates an empty Pipeline.
func NewPipeline() *Pipeline {
	return &Pipeline{}
}

// AddStage appends a named stage to the pipeline.
func (p *Pipeline) AddStage(name string, fn func(map[string]string) (map[string]string, error)) *Pipeline {
	p.stages = append(p.stages, Stage{Name: name, Fn: fn})
	return p
}

// Run executes all stages in order, threading the secrets map through each one.
// If any stage returns an error, execution stops and the error is wrapped with
// the stage name for easy diagnosis.
func (p *Pipeline) Run(secrets map[string]string) (map[string]string, error) {
	current := secrets
	for _, stage := range p.stages {
		result, err := stage.Fn(current)
		if err != nil {
			return nil, fmt.Errorf("pipeline stage %q: %w", stage.Name, err)
		}
		current = result
	}
	return current, nil
}

// StageCount returns the number of registered stages.
func (p *Pipeline) StageCount() int {
	return len(p.stages)
}

// StageNames returns the names of all registered stages in order.
func (p *Pipeline) StageNames() []string {
	names := make([]string, len(p.stages))
	for i, s := range p.stages {
		names[i] = s.Name
	}
	return names
}
