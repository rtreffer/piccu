package secretary

import "mvdan.cc/sh/v3/expand"

// shellEnv provides a WriteEnvironment suitable for shell expansion
// in environment files
type shellEnv struct {
	parent  Environment
	complex map[string]expand.Variable
}

var _ expand.WriteEnviron = &shellEnv{}

func (e shellEnv) Each(f func(name string, vr expand.Variable) bool) {
	for name, value := range e.parent {
		if _, hasComplex := e.complex[name]; hasComplex {
			continue
		}
		if !f(name, expand.Variable{
			Exported: true,
			Kind:     expand.String,
			Str:      value,
		}) {
			return
		}
	}
	for name, value := range e.complex {
		if !f(name, value) {
			return
		}
	}
}

func (e shellEnv) Get(name string) expand.Variable {
	value, found := e.complex[name]
	if found {
		return value
	}
	simple, found := e.parent[name]
	if !found {
		return expand.Variable{
			Kind: expand.Unset,
		}
	}
	return expand.Variable{
		Kind:     expand.String,
		Exported: true,
		Str:      simple,
	}
}

func (e shellEnv) Set(name string, value expand.Variable) error {
	e.complex[name] = value
	return nil
}

func (e shellEnv) Environment() Environment {
	out := make(Environment)
	for name, value := range e.complex {
		if value.Kind == expand.Unset {
			continue
		}
		out[name] = value.String()
	}
	return out
}
