package cicci

import (
	"os"
	"regexp"
	"strings"
)

var filenamePattern = regexp.MustCompile(".*([.]tpl)?.(yaml|yml|json|sh|bash)$")
var templatePattern = regexp.MustCompile(".*[.]tpl.(yaml|yml|json|sh|bash)$")

type CCFile string

func (f *CCFile) IsTemplate() bool {
	return templatePattern.MatchString(string(*f))
}

func (f *CCFile) IsScript() bool {
	return strings.HasSuffix(string(*f), ".sh") || strings.HasSuffix(string(*f), ".bash")
}

func (f *CCFile) NonTemplateName() string {
	parts := strings.Split(string(*f), ".")
	p := len(parts) - 2
	if p < 0 || parts[p] != "tpl" {
		return string(*f)
	}
	return strings.Join(append(parts[0:p], parts[p+1]), ".")
}

func (f *CCFile) Load() (string, error) {
	buf, err := os.ReadFile(string(*f))
	return string(buf), err
}

func (f *CCFile) LoadAndExpand(extraKeys map[string]string) (ExpandedFile, error) {
	var err error

	result := ExpandedFile{
		OriginalFilename: string(*f),
		Filename:         f.NonTemplateName(),
		IsScript:         f.IsScript(),
		IsTemplate:       f.IsTemplate(),
	}

	result.Content, err = f.Load()
	if err != nil {
		return result, err
	}

	if !result.IsTemplate {
		return result, nil
	}

	expanded, err := ExpandTemplate(result.Filename, result.Content, extraKeys)
	if err != nil {
		return result, err
	}

	result.Content = expanded

	return result, nil
}
