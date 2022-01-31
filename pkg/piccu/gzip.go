package piccu

import (
	"bytes"
	"compress/gzip"
)

func GzipString(input string) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, len(input)))
	w, err := gzip.NewWriterLevel(buf, gzip.BestCompression)
	if err != nil {
		return nil, err
	}
	_, err = w.Write([]byte(input))
	if err != nil {
		return nil, err
	}
	err = w.Close()
	return buf.Bytes(), err
}
