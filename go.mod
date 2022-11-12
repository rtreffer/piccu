module github.com/rtreffer/piccu

go 1.19

require (
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/diskfs/go-diskfs v1.2.0
	github.com/golang/glog v1.0.0
	github.com/gopasspw/gopass v1.14.10
	github.com/klauspost/readahead v1.4.0
	github.com/santhosh-tekuri/jsonschema/v5 v5.1.0
	github.com/schollz/progressbar/v3 v3.12.1
	github.com/ulikunitz/xz v0.5.10
	go.fuchsia.dev/fuchsia/src v0.0.0-20210227002123-220857068aaf
	gopkg.in/yaml.v3 v3.0.1
	mvdan.cc/sh/v3 v3.5.1
)

// use our special version of thinfs
replace go.fuchsia.dev/fuchsia/src/lib/thinfs/ => ./pkg/thinfs

require (
	bitbucket.org/creachadair/stringset v0.0.10 // indirect
	filippo.io/age v1.0.0 // indirect
	filippo.io/edwards25519 v1.0.0 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20221026131551-cf6655e29de4 // indirect
	github.com/alessio/shellescape v1.4.1 // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/caspr-io/yamlpath v0.0.0-20200722075116-502e8d113a9b // indirect
	github.com/cloudflare/circl v1.3.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/danieljoos/wincred v1.1.2 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/go-github v17.0.0+incompatible // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/nbutton23/zxcvbn-go v0.0.0-20210217022336-fa2cb2858354 // indirect
	github.com/pierrec/lz4 v2.3.0+incompatible // indirect
	github.com/pkg/xattr v0.4.1 // indirect
	github.com/rivo/uniseg v0.4.2 // indirect
	github.com/rs/zerolog v1.28.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/twpayne/go-pinentry v0.2.0 // indirect
	github.com/urfave/cli/v2 v2.23.5 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	github.com/zalando/go-keyring v0.2.1 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/crypto v0.2.0 // indirect
	golang.org/x/exp v0.0.0-20221111204811-129d8d6c17ab // indirect
	golang.org/x/sys v0.2.0 // indirect
	golang.org/x/term v0.2.0 // indirect
	gopkg.in/djherbis/times.v1 v1.2.0 // indirect
)
