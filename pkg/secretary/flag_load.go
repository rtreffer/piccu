package secretary

import (
	"flag"
	"strings"
)

type LoadType int

const (
	LoadUnknown LoadType = iota
	LoadPlainFile
	LoadPass
	Set
	Unset
)

type LoadFlag struct {
	Type LoadType
	Name string
}

type LoadFlags []LoadFlag

type LoadPlainFileFlag LoadFlags
type LoadPassFlag LoadFlags
type SetFlag LoadFlags
type UnsetFlag LoadFlags

var _ flag.Value = (*LoadPlainFileFlag)(nil)
var _ flag.Value = (*LoadPassFlag)(nil)
var _ flag.Value = (*SetFlag)(nil)
var _ flag.Value = (*UnsetFlag)(nil)

func NewMultiFlagset() (*LoadFlags, *LoadPlainFileFlag, *LoadPassFlag, *SetFlag, *UnsetFlag) {
	result := make(LoadFlags, 0)
	return &result, (*LoadPlainFileFlag)(&result), (*LoadPassFlag)(&result), (*SetFlag)(&result), (*UnsetFlag)(&result)
}

func (f *LoadPlainFileFlag) Set(value string) error {
	*f = append(*f, LoadFlag{
		Type: LoadPlainFile,
		Name: value,
	})
	return nil
}

func (f *LoadPlainFileFlag) String() string {
	files := make([]string, 0, len(*f))
	for _, e := range *f {
		if e.Type != LoadPlainFile {
			continue
		}
		files = append(files, e.Name)
	}
	return strings.Join(files, " ")
}

func (f *LoadPassFlag) Set(value string) error {
	*f = append(*f, LoadFlag{
		Type: LoadPass,
		Name: value,
	})
	return nil
}

func (f *LoadPassFlag) String() string {
	files := make([]string, 0, len(*f))
	for _, e := range *f {
		if e.Type != LoadPass {
			continue
		}
		files = append(files, e.Name)
	}
	return strings.Join(files, " ")
}

func (f *SetFlag) Set(value string) error {
	*f = append(*f, LoadFlag{
		Type: Set,
		Name: value,
	})
	return nil
}

func (f *SetFlag) String() string {
	files := make([]string, 0, len(*f))
	for _, e := range *f {
		if e.Type != Set {
			continue
		}
		files = append(files, e.Name)
	}
	return strings.Join(files, " ")
}

func (f *UnsetFlag) Set(value string) error {
	*f = append(*f, LoadFlag{
		Type: Unset,
		Name: value,
	})
	return nil
}

func (f *UnsetFlag) String() string {
	files := make([]string, 0, len(*f))
	for _, e := range *f {
		if e.Type != Unset {
			continue
		}
		files = append(files, e.Name)
	}
	return strings.Join(files, " ")
}
