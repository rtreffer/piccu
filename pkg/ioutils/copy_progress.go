package ioutils

import (
	"io"
	"os"
	"path/filepath"

	"github.com/klauspost/readahead"
	"github.com/schollz/progressbar/v3"
)

func Copy(src, dst string) error {
	srcStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	bar := progressbar.DefaultBytes(
		srcStat.Size(),
		"copy "+filepath.Base(dst),
	)

	in, err := os.OpenFile(src, os.O_RDONLY, os.FileMode(0644))
	if err != nil {
		return err
	}
	defer in.Close()

	ra := readahead.NewReader(in)
	defer ra.Close()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, os.FileMode(0644))
	if err != nil {
		return err
	}

	_, err = io.Copy(io.MultiWriter(out, bar), ra)
	out.Close()
	return err
}
