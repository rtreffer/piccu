package piccu

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/klauspost/readahead"
	"github.com/schollz/progressbar/v3"
	"github.com/ulikunitz/xz"
)

func ExtractXz(file, target, checksum string, expectedSize int64) error {
	in, err := os.Open(file)
	if err != nil {
		return err
	}
	defer in.Close()
	r, err := xz.NewReader(in)
	if err != nil {
		return err
	}

	// readahead reader pushes the decompression to a dedicated goroutine
	// this will usually improve the performance of the copy
	ra := readahead.NewReader(r)
	defer ra.Close()

	out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, os.FileMode(0644))
	if err != nil {
		return err
	}

	if expectedSize == 0 {
		expectedSize = -1
	}

	hash := sha256.New()

	bar := progressbar.DefaultBytes(
		expectedSize,
		"extract "+filepath.Base(target),
	)
	_, err = io.Copy(io.MultiWriter(out, hash, bar), ra)
	if err != nil {
		return err
	}

	ref := "sha256:" + hex.EncodeToString(hash.Sum(nil))
	if checksum != "" && checksum != ref {
		return fmt.Errorf("broken extraction, expected %s, got %s", checksum, ref)
	}

	return nil
}
