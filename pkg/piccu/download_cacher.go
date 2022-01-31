package piccu

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/schollz/progressbar/v3"
)

const cacheDirName = ".cache"

type flock int

func newFlock(lockname string) (lock flock, err error) {
	lockfd, fd_err := syscall.Open(lockname, syscall.O_CREAT, 0644)
	if fd_err != nil {
		return flock(lockfd), fd_err
	}
	flock_err := syscall.Flock(lockfd, syscall.LOCK_EX)
	if flock_err != nil {
		syscall.Close(lockfd)
		lockfd = 0
	}
	return flock(lockfd), flock_err
}

func (f flock) Unlock() error {
	return syscall.Close(int(f))
}

func verifyDownload(img ImageSource, file string, refresh time.Duration) (bool, error) {
	stat, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if stat.IsDir() {
		return false, os.ErrExist
	}
	if stat.Size() != int64(img.Filesize) && img.Filesize > 0 {
		// wrong size
		return false, nil
	}

	if img.Checksum == "" {
		if stat.ModTime().Before(time.Now().Add(-refresh)) {
			// file is too old
			return false, nil
		}
		// no checksum
		return true, nil
	}

	// we have a file of the right size, let's sha256 it
	hash := sha256.New()
	f, err := os.Open(file)
	if err != nil {
		// we failed to verify the file, this is non-fatal
		return false, nil
	}
	defer f.Close()

	bar := progressbar.DefaultBytes(
		stat.Size(),
		"verify "+filepath.Base(file),
	)

	_, err = io.Copy(io.MultiWriter(hash, bar), f)
	if err != nil {
		return false, err
	}

	ref := "sha256:" + hex.EncodeToString(hash.Sum(nil))

	return img.Checksum == ref, nil
}

func verifyImage(img ImageSource, file string, refresh time.Duration) (bool, error) {
	stat, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if stat.IsDir() {
		return false, os.ErrExist
	}
	if stat.Size() != int64(img.ExtractedFilesize) && img.ExtractedFilesize > 0 {
		// wrong size
		return false, nil
	}

	if img.ImageChecksum == "" {
		if stat.ModTime().Before(time.Now().Add(-refresh)) {
			// file is too old
			return false, nil
		}
		// no checksum
		return true, nil
	}

	// we have a file of the right size, let's sha256 it
	hash := sha256.New()
	f, err := os.Open(file)
	if err != nil {
		// we failed to verify the file, this is non-fatal
		return false, nil
	}
	defer f.Close()

	bar := progressbar.DefaultBytes(
		stat.Size(),
		"verify "+filepath.Base(file),
	)

	_, err = io.Copy(io.MultiWriter(hash, bar), f)
	if err != nil {
		return false, err
	}

	ref := "sha256:" + hex.EncodeToString(hash.Sum(nil))

	return img.ImageChecksum == ref, nil
}

func Download(url, target, checksum string, expectedSize int64) error {
	os.Remove(target)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	l := resp.ContentLength
	if l != 0 && expectedSize != 0 && l != expectedSize {
		return fmt.Errorf("expected download size of %d, got %d", expectedSize, l)
	}
	if l > expectedSize {
		expectedSize = l
	}
	if expectedSize == 0 {
		expectedSize = -1
	}

	hash := sha256.New()

	bar := progressbar.DefaultBytes(
		expectedSize,
		"download "+filepath.Base(target),
	)

	_, err = io.Copy(io.MultiWriter(f, hash, bar), resp.Body)
	if err != nil {
		return err
	}

	ref := "sha256:" + hex.EncodeToString(hash.Sum(nil))
	if checksum != "" && checksum != ref {
		return fmt.Errorf("broken download, expected %s, got %s", checksum, ref)
	}

	return nil
}

func LockfileName(img ImageSource, cachedir string) (string, error) {
	dir := cachedir
	if cachedir == "" {
		dir = cacheDirName
	}
	err := os.Mkdir(dir, os.FileMode(0755))
	if err != nil && !os.IsExist(err) {
		return "", err
	}

	return filepath.Join(dir, fmt.Sprintf("%s-%s-%s.lock", img.Release, img.Codename, img.Architecture)), nil
}

func DownloadName(img ImageSource, cachedir string) (string, error) {
	dir := cachedir
	if cachedir == "" {
		dir = cacheDirName
	}
	err := os.Mkdir(dir, os.FileMode(0755))
	if err != nil && !os.IsExist(err) {
		return "", err
	}

	return filepath.Join(dir, fmt.Sprintf("%s-%s-%s.img.xz", img.Release, img.Codename, img.Architecture)), nil
}

func ImageFilename(img ImageSource, cachedir string) (string, error) {
	dir := cachedir
	if cachedir == "" {
		dir = cacheDirName
	}
	err := os.Mkdir(dir, os.FileMode(0755))
	if err != nil && !os.IsExist(err) {
		return "", err
	}

	return filepath.Join(dir, fmt.Sprintf("%s-%s-%s.img", img.Release, img.Codename, img.Architecture)), nil
}

func fetchDownload(img ImageSource, cachedir string, refresh time.Duration) error {
	downloadName, err := DownloadName(img, cachedir)
	if err != nil {
		return err
	}

	// we are done if we can verify the download
	verified, err := verifyDownload(img, downloadName, refresh)
	if err != nil {
		return err
	}
	if verified {
		return nil
	}

	return Download(img.URL, downloadName, img.Checksum, img.Filesize)
}

// Fetch fetches the image and returns a raw image file path
func Fetch(img ImageSource, cachedir string, refresh time.Duration) (string, error) {
	downloadName, err := DownloadName(img, cachedir)
	if err != nil {
		return "", err
	}
	lockName, err := LockfileName(img, cachedir)
	if err != nil {
		return "", err
	}
	imageName, err := ImageFilename(img, cachedir)
	if err != nil {
		return "", err
	}

	lock, err := newFlock(lockName)
	if err != nil {
		return "", err
	}
	defer lock.Unlock()

	// check if we can verify the image
	verified, err := verifyImage(img, imageName, refresh)
	if err != nil {
		return "", err
	}
	if verified {
		return imageName, nil
	}

	// we need to at least download the file
	err = fetchDownload(img, cachedir, refresh)
	if err != nil {
		return "", err
	}

	// and extract it
	err = ExtractXz(downloadName, imageName, img.ImageChecksum, img.ExtractedFilesize)
	if err != nil {
		return "", err
	}

	return imageName, nil
}
