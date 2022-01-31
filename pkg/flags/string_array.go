package flags

import "strings"

type StringArray []string

func (f *StringArray) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func (f *StringArray) String() string {
	return strings.Join(*f, " ")
}
