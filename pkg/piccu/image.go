package piccu

import (
	"fmt"
	"os"

	blockfile "go.fuchsia.dev/fuchsia/src/lib/thinfs/block/file"
	"go.fuchsia.dev/fuchsia/src/lib/thinfs/fs"
	"go.fuchsia.dev/fuchsia/src/lib/thinfs/fs/msdosfs"

	"github.com/diskfs/go-diskfs"
	"github.com/diskfs/go-diskfs/filesystem"
)

type Image struct {
	path       string
	underlying *os.File
	offset     int64
	length     int64
	block      *blockfile.File
	fs         fs.FileSystem
}

func OpenImage(file string) (result *Image, err error) {
	result = &Image{
		path: file,
	}

	disk, err := diskfs.Open(file)
	if err != nil {
		return
	}

	blockSize := disk.PhysicalBlocksize

	partitionTable, err := disk.GetPartitionTable()
	if err != nil {
		disk.File.Close()
		return
	}

	partitions := partitionTable.GetPartitions()
	for i := 1; i <= len(partitions); i++ {
		fs, err := disk.GetFilesystem(i)
		if err != nil || fs.Type() != filesystem.TypeFat32 {
			continue
		}

		result.offset = partitions[i-1].GetStart()
		result.length = partitions[i-1].GetSize()
		break
	}

	disk.File.Close()

	result.underlying, err = os.OpenFile(file, os.O_RDWR|os.O_EXCL, os.FileMode(0644))
	if err != nil {
		return nil, err
	}

	result.block, err = blockfile.NewRange(result.underlying, blockSize, result.offset, result.length)
	if err != nil {
		result.underlying.Close()
		return nil, err
	}

	result.fs, err = msdosfs.New("/", result.block, fs.ReadWrite)
	if err != nil {
		result.underlying.Close()
		return nil, err
	}

	return result, nil
}

func (img *Image) Close() error {
	// close the filesystem, block backend, file
	img.fs.Close()
	img.block.Flush()
	img.block.Close()
	img.underlying.Sync()
	return img.underlying.Close()
}

func (img *Image) InjectFile(path string, payload []byte) error {
	img.fs.RootDirectory().Unlink(path)
	file, _, _, err := img.fs.RootDirectory().Open(path, fs.OpenFlagCreate|fs.OpenFlagWrite|fs.OpenFlagFile)
	if err != nil {
		return err
	}
	defer func() {
		file.Sync()
		file.Close()
	}()
	n, err := file.Write(payload, 0, fs.WhenceFromStart)
	if err != nil {
		return err
	}
	if n != len(payload) {
		return fmt.Errorf("expected to write %d, wrote %d", len(payload), n)
	}
	return nil
}
