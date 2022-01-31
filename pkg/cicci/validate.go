package cicci

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"gopkg.in/yaml.v3"
	"mvdan.cc/sh/v3/syntax"
)

//go:embed cloud-config.schema.json
var jsonSchema string
var compiledJsonSchema = jsonschema.MustCompileString("cloud-config.schema.json", jsonSchema)

func ValidateCloudConfig(payload string) error {
	var value interface{}
	err := yaml.Unmarshal([]byte(payload), &value)
	if err != nil {
		return err
	}
	return compiledJsonSchema.Validate(value)
}

func ValidateScript(filename, payload string) error {
	lang := syntax.LangBash
	if !strings.HasSuffix(filename, ".bash") {
		lang = syntax.LangPOSIX
	}
	parser := syntax.NewParser(syntax.Variant(lang))
	_, err := parser.Parse(strings.NewReader(payload), filename)
	if err != nil {
		if lang == syntax.LangBash {
			return err
		}
		lang = syntax.LangBash
	}
	parser = syntax.NewParser(syntax.Variant(lang))
	_, basherr := parser.Parse(strings.NewReader(payload), filename)
	if basherr != nil {
		// if it wasn't a bash script the let's return the POSIX shell error
		return err
	}
	return fmt.Errorf("%s is not a posix shell, please use .bash", filename)
}
