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
	// URL is the download URL
	URL string
	// Raspberry versions
	Raspberry []string
	// Filesize in bytes
	Filesize int64
	// ExtractedFilesize is the unpacted file size in bytes
	ExtractedFilesize int64
	// Checksum
	Checksum string
	// ImageChecksum
	ImageChecksum string
}

var ImageSources = []ImageSource{
	{
		Release:           "16.04",
		Codename:          "xenial",
		Architecture:      ARMHF,
		URL:               "https://cdimage.ubuntu.com/ubuntu/releases/16.04.7/release/ubuntu-16.04.4-preinstalled-server-armhf+raspi2.img.xz",
		Raspberry:         []string{RPI2},
		Filesize:          247934984,
		ExtractedFilesize: 4000000000,
		Checksum:          "sha256:a6c96eee607f80c25436f7b786d74fb38d39894c2ea034d2aabfcc8c4cecd101",
		ImageChecksum:     "sha256:69830f5f8d9360156ed869751474898b3de83d02eae0f427e0727f78f89633f1",
	},
	{
		Release:           "18.04",
		Codename:          "bionic",
		Architecture:      ARM64,
		URL:               "https://cdimage.ubuntu.com/ubuntu/releases/18.04.6/release/ubuntu-18.04.5-preinstalled-server-arm64+raspi3.img.xz",
		Filesize:          509610804,
		ExtractedFilesize: 2653289472,
		// this image is identical to the rpi4 version
		Raspberry:     []string{RPI3, RPI4},
		Checksum:      "sha256:69cbd2c0b70bc2cd1d8b6aa0a98bd64d59617b03f2681ff4ac56c85daa44cde5",
		ImageChecksum: "sha256:c581f72ec6d17d270216bb2b18aca7ed47b21990560985222c2259d766403262",
	},
	{
		Release:      "18.04",
		Codename:     "bionic",
		Architecture: ARMHF,
		URL:          "https://cdimage.ubuntu.com/ubuntu/releases/18.04.6/release/ubuntu-18.04.5-preinstalled-server-armhf+raspi2.img.xz",
		// this image is identical to the rpi3/rpi4 images
		Raspberry:         []string{RPI2, RPI3, RPI4},
		Filesize:          499504344,
		ExtractedFilesize: 2417323008,
		Checksum:          "sha256:343692137d74490dabbb8b20568827e5d53c5d01ff760244ac37cf559f3ee3b4",
		ImageChecksum:     "sha256:3edf405bd9f1418bc257ad3a448b02f1d0f943b4d7ecd9d5384f26193cb80a23",
	},
	{
		Release:      "20.04",
		Codename:     "focal",
		Architecture: ARM64,
		URL:          "https://cdimage.ubuntu.com/ubuntu/releases/20.04/release/ubuntu-20.04.4-preinstalled-server-arm64+raspi.img.xz",
		// newer versions of the RPI2 use the 64bit CPU. It may or may not work. 20.04 advertises support for RPI2
		Raspberry:         []string{RPI2, RPI3, RPI4},
		Filesize:          769906140,
		ExtractedFilesize: 3487742976,
		Checksum:          "sha256:6aeba20c00ef13ee7b48c57217ad0d6fc3b127b3734c113981d9477aceb4dad7",
		ImageChecksum:     "sha256:fc3e8fbeb7c3536aaeaf1a77ddbc95eb9bb98de17346d24ee68716856ce7f0fd",
	},
	{
		Release:           "20.04",
		Codename:          "focal",
		Architecture:      ARMHF,
		URL:               "https://cdimage.ubuntu.com/ubuntu/releases/20.04/release/ubuntu-20.04.4-preinstalled-server-armhf+raspi.img.xz",
		Raspberry:         []string{RPI2, RPI3, RPI4},
		Filesize:          727412296,
		ExtractedFilesize: 3160341504,
		Checksum:          "sha256:3b1704e8e4ff8e01dd89b9dd6adf9b99b48b2a7530d6f7676ce8c37772ff4178",
		ImageChecksum:     "sha256:4281cce92e532e8443253e88c0924bd55100bc9c5f17c393619fe37560e5a586",
	},
	{
		Release:           "21.04",
		Codename:          "hirsute",
		Architecture:      ARM64,
		URL:               "https://cdimage.ubuntu.com/releases/21.04/release/ubuntu-21.04-preinstalled-server-arm64+raspi.img.xz",
		Raspberry:         []string{RPI2, RPI3, RPI4},
		Filesize:          787262724,
		ExtractedFilesize: 3491662848,
		Checksum:          "sha256:3df85b93b66ccd2d370c844568b37888de66c362eebae5204bf017f6f5875207",
		ImageChecksum:     "sha256:0d1eb068e55879ed279a3c9ba79fe186919db5746606477a3ba4b76318fd10cd",
	},
	{
		Release:           "21.04",
		Codename:          "hirsute",
		Architecture:      ARMHF,
		URL:               "https://cdimage.ubuntu.com/releases/21.04/release/ubuntu-21.04-preinstalled-server-armhf+raspi.img.xz",
		Raspberry:         []string{RPI2, RPI3, RPI4},
		Filesize:          756021496,
		ExtractedFilesize: 3211852800,
		Checksum:          "sha256:c9a9a5177a03fcbb6203b38e5c3c4e5447fd9e8891515da4146f319f04eb3495",
		ImageChecksum:     "sha256:8eeeaa116b91f4f622bc3abaa80453387783de6bfce77f7a51a02c5337f5c7e1",
	},
	{
		Release:           "21.10",
		Codename:          "impish",
		Architecture:      ARM64,
		URL:               "https://cdimage.ubuntu.com/releases/21.10/release/ubuntu-21.10-preinstalled-server-arm64+raspi.img.xz",
		Raspberry:         []string{RPI2, RPI3, RPI4},
		Filesize:          921117368,
		ExtractedFilesize: 4068480000,
		Checksum:          "sha256:126f940d3b270a6c1fc5a183ac8a3d193805fead4f517296a7df9d3e7d691a03",
		ImageChecksum:     "sha256:4cf06429e0367f0a59b890819d1792b0d816bc531fcb5bd092e441d8d6a942b9",
	},
	{
		Release:           "21.10",
		Codename:          "impish",
		Architecture:      ARMHF,
		URL:               "https://cdimage.ubuntu.com/releases/21.10/release/ubuntu-21.10-preinstalled-server-armhf+raspi.img.xz",
		Raspberry:         []string{RPI2, RPI3, RPI4},
		Filesize:          886173524,
		ExtractedFilesize: 3732335616,
		Checksum:          "sha256:341593c9607ed20744cd86941d94d73e3ba4f74e8ef2633eec63ce9b0cfac5a1",
		ImageChecksum:     "sha256:56caec8fd34aa4aec01641aa3ac3993d21468b375835ca40a7ccf948947ca353",
	},
	{
		Release:           "22.04",
		Codename:          "jammy",
		Architecture:      ARM64,
		URL:               "https://cdimage.ubuntu.com/releases/22.04/release/ubuntu-22.04-preinstalled-server-arm64+raspi.img.xz",
		Raspberry:         []string{RPI2, RPI3, RPI4},
		Filesize:          953648752,
		ExtractedFilesize: 3937402880,
		Checksum:          "sha256:f41936779e66ee1942878aa7e4b16df74d4b1b90f644da4d7ce70993b4fb3c46",
		ImageChecksum:     "sha256:9cd6b5e6b4e4a7453cfde276927efe20380ab9ec0377318d5ce0bc8c8a56993b",
	},
	{
		Release:           "22.04",
		Codename:          "jammy",
		Architecture:      ARMHF,
		URL:               "https://cdimage.ubuntu.com/releases/22.04/release/ubuntu-22.04-preinstalled-server-armhf+raspi.img.xz",
		Raspberry:         []string{RPI2, RPI3, RPI4},
		Filesize:          918555612,
		ExtractedFilesize: 3674210304,
		Checksum:          "sha256:d93bc4b7f4040c04e17ca44e61f1e08aa8fb9086671b4065b8722e2b5dd6543e",
		ImageChecksum:     "sha256:569cd1a885cdfd839689d70794778a807737c30b6cf63fb982f62799ba54d6f6",
	},
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
