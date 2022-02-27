package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/klauspost/readahead"
	"github.com/schollz/progressbar/v3"
	"github.com/ulikunitz/xz"
)

func checkPanic(err error) {
	if err == nil {
		return
	}
	panic(err)
}

type CountWriter struct {
	Count int
}

func (c *CountWriter) Write(buf []byte) (int, error) {
	c.Count += len(buf)
	return len(buf), nil
}

func main() {
	for _, img := range os.Args[1:] {
		resp, err := http.Get(img)
		checkPanic(err)
		bar := progressbar.DefaultBytes(
			resp.ContentLength,
			"downloading "+filepath.Base(img),
		)

		origSize := &CountWriter{}
		unxzSize := &CountWriter{}
		hash := sha256.New()
		hashunxz := sha256.New()

		xzIn, xzOut := io.Pipe()

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			plain, err := xz.NewReader(xzIn)
			ra := readahead.NewReader(plain)
			checkPanic(err)
			_, err = io.Copy(io.MultiWriter(hashunxz, unxzSize), ra)
			checkPanic(err)
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			ra := readahead.NewReader(resp.Body)
			_, err := io.Copy(io.MultiWriter(xzOut, origSize, hash, bar), ra)
			checkPanic(err)
			xzOut.Close()
		}()
		wg.Wait()
		bar.Finish()

		fmt.Printf("%s\t%d\t%s\t%d\t%s\n",
			img,
			origSize.Count,
			"sha256:"+hex.EncodeToString(hash.Sum(nil)),
			unxzSize.Count,
			"sha256:"+hex.EncodeToString(hashunxz.Sum(nil)),
		)
	}
}
