package cicci

import (
	"os"
	"path/filepath"
)

type CCFiles []CCFile

func (files CCFiles) LoadAndExpand(extraKeys map[string]string) (result ExpandedFiles, err error) {
	result = make(ExpandedFiles, len(files))
	for i, file := range files {
		result[i], err = file.LoadAndExpand(extraKeys)
		if err != nil {
			return
		}
	}
	return
}

func CollectFiles(inputs []string) (CCFiles, error) {
	output := make([]CCFile, 0, len(inputs))
	for _, input := range inputs {
		stat, err := os.Stat(input)

		if os.IsNotExist(err) {
			// try glob
			matches, err := filepath.Glob(input)
			if err != nil {
				// this may happen if we tried a glob on something that wasn't
				continue
			}
			if len(matches) == 0 {
				continue
			}
			files, err := CollectFiles(matches)
			if err != nil {
				return nil, err
			}
			output = append(output, files...)
			continue
		}

		if err != nil {
			return nil, err
		}

		if stat.IsDir() {
			entries, err := os.ReadDir(input)
			if err != nil {
				return nil, err
			}
			names := make([]string, 0, len(entries))
			for _, e := range entries {
				if e.Name() == "." || e.Name() == ".." {
					continue
				}
				names = append(names, filepath.Join(input, e.Name()))
			}
			recursion, err := CollectFiles(names)
			if err != nil {
				return nil, err
			}
			output = append(output, recursion...)
			continue
		}

		if filenamePattern.MatchString(input) {
			output = append(output, CCFile(input))
		}
	}
	return output, nil
}
