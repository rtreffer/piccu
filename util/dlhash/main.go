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

	"github.com/schollz/progressbar/v3"
	"github.com/ulikunitz/xz"
)

func checkPanic(err error) {
	if err == nil {
		return
	}
	panic(err)
}

func main() {
	for _, img := range os.Args[1:] {
		resp, err := http.Get(img)
		checkPanic(err)
		bar := progressbar.DefaultBytes(
			resp.ContentLength,
			"download "+filepath.Base(img),
		)

		hash := sha256.New()
		hashunxz := sha256.New()

		xzIn, xzOut := io.Pipe()

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			plain, err := xz.NewReader(xzIn)
			checkPanic(err)
			_, err = io.Copy(hashunxz, plain)
			checkPanic(err)
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := io.Copy(io.MultiWriter(xzOut, hash, bar), resp.Body)
			checkPanic(err)
			xzOut.Close()
		}()
		wg.Wait()
		bar.Finish()

		fmt.Printf("%s\t%s\t%s\n",
			img,
			"sha256:"+hex.EncodeToString(hash.Sum(nil)),
			"sha256:"+hex.EncodeToString(hashunxz.Sum(nil)),
		)
	}
}
