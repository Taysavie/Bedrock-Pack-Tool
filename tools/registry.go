package tools

import (
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type Registry struct {
	tools *orderedmap.OrderedMap[string, Tool]
}

func NewToolRegistry() *Registry {
	return &Registry{tools: orderedmap.New[string, Tool]()}
}

func (r *Registry) RegisterTool(name string, tool Tool) {
	r.tools.Set(name, tool)
}

func (r *Registry) Names() []string {
	names := make([]string, r.tools.Len())
	i := 0
	for n := range r.tools.FromOldest() {
		names[i] = n
		i++
	}
	return names
}

func (r *Registry) ToolByName(name string) Tool {
	t, _ := r.tools.Get(name)
	return t
}
