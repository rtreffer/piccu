package cicci

import (
	"bytes"
	"os"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
)

func environmentContext() map[string]string {
	env := make(map[string]string)
	for _, envKeyValue := range os.Environ() {
		parts := strings.SplitN(envKeyValue, "=", 2)
		env[parts[0]] = parts[1]
	}
	return env
}

func ExpandTemplate(name, tpl string, extraKeys map[string]string) (string, error) {
	tmpl, err := template.New(name).Parse(string(tpl))
	if err != nil {
		return "", err
	}
	tmpl.Funcs(template.FuncMap(sprig.TxtFuncMap()))
	buf := new(bytes.Buffer)
	self := environmentContext()
	for k, v := range extraKeys {
		self[k] = v
	}

	if err := tmpl.Execute(buf, self); err != nil {
		return "", err
	}
	return buf.String(), nil
}
