package flags

import "strings"

type StringMap map[string]string
type UnsetMap StringMap

func NewStringMap() (*StringMap, *UnsetMap) {
	stringMap := make(StringMap)
	result := &stringMap
	return result, (*UnsetMap)(result)
}

func (f *StringMap) Set(value string) error {
	parts := strings.SplitN(value, "=", 2)
	if len(parts) == 1 {
		(*f)[value] = ""
		return nil
	}
	(*f)[parts[0]] = parts[1]
	return nil
}

func (f *StringMap) String() string {
	entries := make([]string, 0, len(*f))
	for k, v := range *f {
		entries = append(entries, k+"="+v)
	}
	return strings.Join(entries, " ")
}

func (f *UnsetMap) Set(value string) error {
	delete(*f, value)
	return nil
}

func (f *UnsetMap) String() string {
	return ""
}
