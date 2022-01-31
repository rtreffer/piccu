package secretary

import (
	"fmt"
	"io"
	"os"
	"strings"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/syntax"
)

type Environment map[string]string

func (e Environment) Envv() []string {
	result := make([]string, 0, len(e))
	for k, v := range e {
		result = append(result, k+"="+v)
	}
	return result
}

func ParseEnvironmentFile(name, content string) (Environment, error) {
	// an environment file can be seen as a POSIX shell script
	// let's try to parse it
	parser := syntax.NewParser(syntax.Variant(syntax.LangBash))
	parsed, err := parser.Parse(strings.NewReader(content), name)
	if err != nil {
		return nil, err
	}

	// needed for e.g. ${HOME}
	parent, err := ParseEnvironment(os.Environ())
	if err != nil {
		return nil, err
	}

	shenv := shellEnv{
		parent:  parent,
		complex: make(map[string]expand.Variable),
	}
	cfg := &expand.Config{
		Env: shenv,
		CmdSubst: func(w io.Writer, cs *syntax.CmdSubst) error {
			return fmt.Errorf("command substitution not supported for security reasons: %v", cs)
		},
		ProcSubst: func(ps *syntax.ProcSubst) (string, error) {
			return "", fmt.Errorf("process substitution not supported for security reasons: %v", ps)
		},
		ReadDir: func(s string) ([]os.FileInfo, error) {
			return nil, fmt.Errorf("globbing not implemented")
		},
	}

	for _, stmt := range parsed.Stmts {
		call, ok := stmt.Cmd.(*syntax.CallExpr)
		if !ok {
			continue
		}
		if len(call.Args) > 0 {
			continue
		}
		for _, assign := range call.Assigns {
			name := assign.Name.Value
			value, err := expand.Literal(cfg, assign.Value)
			if err != nil {
				return nil, err
			}
			shenv.Set(name, expand.Variable{
				Exported: true,
				Kind:     expand.String,
				Str:      value,
			})
		}
	}

	return shenv.Environment(), nil
}

func ParseEnvironment(envv []string) (Environment, error) {
	out := make(Environment)
	for _, env := range envv {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid environment entry: %s", env)
		}
		out[parts[0]] = parts[1]
	}
	return out, nil
}
