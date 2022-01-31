package cicci

import "fmt"

type ExpandedFiles []ExpandedFile

type ExpandedFile struct {
	OriginalFilename string
	Filename         string
	Content          string
	IsScript         bool
	IsTemplate       bool
}

func (e *ExpandedFile) Validate() error {
	if !e.IsScript {
		// validate YAML/JSON
		if err := ValidateCloudConfig(e.Content); err != nil {
			return fmt.Errorf("error in %s: %s", e.OriginalFilename, err)
		}
		return nil
	}
	// validate shell scripts
	return ValidateScript(e.Filename, e.Content)
}

func (s ExpandedFiles) Validate() (result []error) {
	for _, e := range s {
		if err := e.Validate(); err != nil {
			result = append(result,
				fmt.Errorf("%s: %s", e.OriginalFilename, err))
		}
	}
	return
}
