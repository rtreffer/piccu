module github.com/rtreffer/piccu

go 1.17

require (
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/diskfs/go-diskfs v1.2.0
	github.com/golang/glog v1.0.0
	github.com/gopasspw/gopass v1.13.1
	github.com/klauspost/readahead v1.4.0
	github.com/santhosh-tekuri/jsonschema/v5 v5.0.0
	github.com/schollz/progressbar/v3 v3.8.5
	github.com/ulikunitz/xz v0.5.10
	go.fuchsia.dev/fuchsia/src v0.0.0-20210227002123-220857068aaf
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	mvdan.cc/sh/v3 v3.4.2
)

// use our special version of thinfs
replace go.fuchsia.dev/fuchsia/src/lib/thinfs/ => ./pkg/thinfs

require (
	filippo.io/age v1.0.0 // indirect
	filippo.io/edwards25519 v1.0.0-rc.1 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/caspr-io/yamlpath v0.0.0-20200722075116-502e8d113a9b // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.1 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/google/go-github v17.0.0+incompatible // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gopasspw/pinentry v0.0.2 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/pierrec/lz4 v2.3.0+incompatible // indirect
	github.com/pkg/xattr v0.4.1 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/urfave/cli/v2 v2.3.0 // indirect
	golang.org/x/crypto v0.0.0-20220128200615-198e4374d7ed // indirect
	golang.org/x/sys v0.0.0-20220128215802-99c3d69c2c27 // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	gopkg.in/djherbis/times.v1 v1.2.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
