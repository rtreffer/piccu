package piccu

import "sort"

const RPI2 = "raspberry pi 2"
const RPI3 = "raspberry pi 3"
const RPI4 = "raspberry pi 4"

const ARMHF = "armhf"
const ARM64 = "arm64"

type ImageSource struct {
	// Release is the LSB release string, e.g. "22.04"
	Release string
	// Codename is the LSB codename, e.g. "jammy"
	Codename string
	// Architecture is the architecture, e.g. "armhf" or "arm64"
	Architecture string
	// Raspberry versions
	Raspberry []string
	// URL is the download URL
	URL string
	// Filesize in bytes
	Filesize int64
	// ExtractedFilesize is the unpacted file size in bytes
	ExtractedFilesize int64
	// Checksum
	Checksum string
	// ImageChecksum
	ImageChecksum string
}

func (img *ImageSource) ArchitectureCode() (output string) {
	output = `""`
	if img == nil {
		return
	}
	if img.Architecture == ARM64 {
		return "ARM64"
	}
	if img.Architecture == ARMHF {
		return "ARMHF"
	}
	return
}

func (img *ImageSource) RpiCode() (output string) {
	output = "[]string{}"
	if img == nil || len(img.Raspberry) == 0 {
		return
	}
	if len(img.Raspberry) == 1 {
		if img.Raspberry[0] != RPI2 {
			return
		}
		return "[]string{RPI2}"
	}
	if len(img.Raspberry) == 2 {
		if img.Raspberry[0] != RPI3 || img.Raspberry[1] != RPI4 {
			return
		}
		return "[]string{RPI3, RPI4}"
	}
	if len(img.Raspberry) == 3 {
		if img.Raspberry[0] != RPI2 || img.Raspberry[1] != RPI3 || img.Raspberry[2] != RPI4 {
			return
		}
		return "[]string{RPI2, RPI3, RPI4}"
	}
	return
}

func ImagesByKey() map[string]ImageSource {
	names := make(map[string]ImageSource)
	for _, img := range ImageSources {
		if img.Architecture == ARMHF {
			names[img.Release+":armhf"] = img
			names[img.Codename+":armhf"] = img
			if _, found := names[img.Release]; !found {
				// ony add the armhf image if an amr64 one does not exist
				names[img.Release] = img
				names[img.Codename] = img
			}
		} else {
			names[img.Release+":arm64"] = img
			names[img.Codename+":arm64"] = img
			names[img.Release] = img
			names[img.Codename] = img
		}
	}
	return names
}

func GetImageNames() []string {
	names := ImagesByKey()
	sortedNames := make([]string, 0, len(names))
	for k := range names {
		sortedNames = append(sortedNames, k)
	}
	sort.Strings(sortedNames)
	return sortedNames
}
