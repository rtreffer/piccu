package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/klauspost/readahead"
	"github.com/rtreffer/piccu/pkg/piccu"
	"github.com/schollz/progressbar/v3"
	"github.com/ulikunitz/xz"
)

const urlPrefix = "https://cdimage.ubuntu.com"

const metaURL = "ubuntu/releases/streams/v1/com.ubuntu.cdimage:ubuntu.json"

var goCodeTemplate = `package piccu

//go:generate go run github.com/rtreffer/piccu/pkg/piccu/gen/

var ImageSources = []ImageSource{
    {{- range $index, $src := .sources }}
	{
		Release:           "{{ $src.Release }}",
		Codename:          "{{ $src.Codename }}",
		Architecture:      {{ $src.ArchitectureCode }},
		URL:               "{{ $src.URL }}",
		Raspberry:         {{ $src.RpiCode }},
		Filesize:          {{ $src.Filesize }},
		ExtractedFilesize: {{ $src.ExtractedFilesize }},
		Checksum:          "{{ $src.Checksum }}",
		ImageChecksum:     "{{ $src.ImageChecksum }}",
	},
	{{- end }}
}
`

// ImageSourcesSource is the image source without size & hash information
var ImageSourcesSource = []piccu.ImageSource{
	{
		"16.04", "xenial", piccu.ARMHF, []string{piccu.RPI2},
		"ubuntu/releases/xenial/release/ubuntu-16.04.6-preinstalled-server-armhf+raspi2.img.xz",
		0, 0, "", "",
	},
	{
		"18.04", "bionic", piccu.ARM64, []string{piccu.RPI3, piccu.RPI4},
		"ubuntu/releases/bionic/release/ubuntu-18.04.5-preinstalled-server-arm64+raspi3.img.xz",
		0, 0, "", "",
	},
	{
		"18.04", "bionic", piccu.ARMHF, []string{piccu.RPI2, piccu.RPI3, piccu.RPI4},
		"ubuntu/releases/bionic/release/ubuntu-18.04.5-preinstalled-server-armhf+raspi2.img.xz",
		0, 0, "", "",
	},
	{
		"20.04", "focal", piccu.ARM64, []string{piccu.RPI2, piccu.RPI3, piccu.RPI4},
		"ubuntu/releases/focal/release/ubuntu-20.04.5-preinstalled-server-arm64+raspi.img.xz",
		0, 0, "", "",
	},
	{
		"20.04", "focal", piccu.ARMHF, []string{piccu.RPI2, piccu.RPI3, piccu.RPI4},
		"ubuntu/releases/focal/release/ubuntu-20.04.5-preinstalled-server-armhf+raspi.img.xz",
		0, 0, "", "",
	},
	{
		"22.04", "jammy", piccu.ARM64, []string{piccu.RPI2, piccu.RPI3, piccu.RPI4},
		"ubuntu/releases/jammy/release/ubuntu-22.04.1-preinstalled-server-arm64+raspi.img.xz",
		0, 0, "", "",
	},
	{
		"22.04", "jammy", piccu.ARMHF, []string{piccu.RPI2, piccu.RPI3, piccu.RPI4},
		"ubuntu/releases/jammy/release/ubuntu-22.04.1-preinstalled-server-armhf+raspi.img.xz",
		0, 0, "", "",
	},
}

type UbuntuProducts struct {
	Products map[string]UbuntuProduct `json:"products"`
}

type UbuntuProduct struct {
	Architecture    string                          `json:"arch"`
	Release         string                          `json:"release"`
	ReleaseCodename string                          `json:"release_codename"`
	ReleaseTitle    string                          `json:"release_title"`
	Version         string                          `json:"version"`
	ImageType       string                          `json:"image_type"`
	OS              string                          `json:"os"`
	Versions        map[string]UbuntuProductVersion `json:"versions"`
}

type UbuntuProductVersion struct {
	Items map[string]UbuntuProductVersionItems `json:"items"`
}

type UbuntuProductVersionItems struct {
	FileType string `json:"ftype"`
	Path     string `json:"path"`
	Size     uint64 `json:"size"`
	Sha256   string `json:"sha256"`
}

func (p *UbuntuProducts) FindFile(path string) (size uint64, hash string, found bool) {
	for _, p := range p.Products {
		for _, v := range p.Versions {
			for _, i := range v.Items {
				if i.Path == path {
					return i.Size, i.Sha256, true
				}
			}
		}
	}
	return 0, "", false
}

type CountWriter int

func (c *CountWriter) Write(buf []byte) (int, error) {
	*c = CountWriter(int(*c) + len(buf))
	return len(buf), nil
}

func checkXZDownload(url string, expectedSize int64) (size, extractedSize int64, checksum, extractedChecksum string, err error) {
	resp, httperr := http.Get(url)
	defer resp.Body.Close()
	if httperr != nil {
		err = httperr
		return
	}

	bar := progressbar.DefaultBytes(
		expectedSize,
		filepath.Base(url),
	)
	if expectedSize == 0 {
		bar = progressbar.DefaultBytes(
			-1,
			filepath.Base(url),
		)
	}

	sizeWriter := CountWriter(0)
	extractedSizeWriter := CountWriter(0)
	hash := sha256.New()
	extractedHash := sha256.New()

	extractedWriter := io.MultiWriter(&extractedSizeWriter, extractedHash)
	xzIn, xzOut := io.Pipe()
	xzWriter := io.MultiWriter(hash, xzOut, bar, &sizeWriter)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		plain, xzerr := xz.NewReader(xzIn)
		if xzerr != nil {
			err = xzerr
			return
		}
		ra := readahead.NewReader(plain)
		_, xzerr = io.Copy(extractedWriter, ra)
		if xzerr != nil {
			err = xzerr
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer xzOut.Close()
		ra := readahead.NewReader(resp.Body)
		_, httperr := io.Copy(xzWriter, ra)
		if httperr != nil {
			err = httperr
		}
	}()
	wg.Wait()
	bar.Finish()

	size = int64(sizeWriter)
	extractedSize = int64(extractedSizeWriter)
	checksum = "sha256:" + hex.EncodeToString(hash.Sum(nil))
	extractedChecksum = "sha256:" + hex.EncodeToString(extractedHash.Sum(nil))

	return
}

func main() {
	tmpl := template.New("image_source_table.go")
	tmpl, err := tmpl.Parse(goCodeTemplate)
	if err != nil {
		panic(err)
	}

	// download and parse the ubuntu product list
	resp, err := http.Get(urlPrefix + "/" + metaURL)
	if err != nil {
		panic(err)
	}
	var products UbuntuProducts
	if err = json.NewDecoder(resp.Body).Decode(&products); err != nil {
		panic(err)
	}
	for i, src := range ImageSourcesSource {
		fmt.Println("::", i, "::", src.Codename, src.Architecture, src.URL)
		size, hash, found := products.FindFile(src.URL)
		if found {
			src.Filesize = int64(size)
			src.Checksum = "sha256:" + hash
			fmt.Println("        size:", size, "sha256:", hash)
		}
		src.URL = urlPrefix + "/" + src.URL
		ImageSourcesSource[i] = src
		for _, other := range piccu.ImageSources {
			if other.Architecture != src.Architecture {
				continue
			}
			if other.Release != src.Release || other.Codename != src.Codename {
				continue
			}
			if other.URL == src.URL || other.Checksum == src.Checksum {
				// everything should be the same
				src.Filesize = other.Filesize
				src.Checksum = other.Checksum
				src.ExtractedFilesize = other.ExtractedFilesize
				src.ImageChecksum = other.ImageChecksum
				continue
			}
		}
		if src.ExtractedFilesize != 0 && src.ImageChecksum != "" {
			fmt.Println("        image size:", src.ExtractedFilesize, "image hash:", src.ImageChecksum)
			ImageSourcesSource[i] = src
			continue
		}
		dlsize, dlimgsize, hash, imghash, err := checkXZDownload(src.URL, src.Filesize)
		if err != nil {
			panic(err)
		}
		if dlsize != src.Filesize || hash != src.Checksum {
			panic(fmt.Errorf("mismatch on download size/hash, expected %d/%s, got %d/%s",
				src.Filesize, src.Checksum,
				dlsize, hash))
		}
		src.ExtractedFilesize = dlimgsize
		src.ImageChecksum = imghash
		fmt.Println("        image size:", src.ExtractedFilesize, "image hash:", src.ImageChecksum)
		ImageSourcesSource[i] = src
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"sources": ImageSourcesSource,
	})
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("image_sources_table.go", buf.Bytes(), os.FileMode(0644))
	if err != nil {
		panic(err)
	}
}
