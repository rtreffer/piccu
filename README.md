# `piccu`: `pi + cloud-config + ubuntu` image builder

**This software is currently alpha software - it can be tried if you are ready to fix it**

**Getting started:**
```
piccu -o grafana.img --ubuntu jammy --set hostname=grafana examples/piusers.yaml examples/grafana.yaml examples/hostname.tpl.yaml exampels/avahi.yaml
```
Wait a few minutes until grafana.img is built and put it onto an sd-card. boot the raspberry pi and enjoy [`http://grafana.local`](http://grafana.local)

# Limitations

`piccu` can be used for _initial_ provisioning. It is not a substitute for long-term configuration management. However many raspberry pi projects do not need extensive configuration management.

Merging cloud-configs can have side-effects. E.g. the first user will be the "default" user.

# Tools

`piccu` is the `pi cloud-config ubuntu` image creator.
It can embed cloud-config and instance information into an ubuntu raspberry pi image without a need for root privileges. `piccu` inherits all capabilities of `cicci` and `secretary`.
```
piccu -c jammy -o pi-node1.img --pass nodes/1 node1/
```

`cicci /ˈzɪsi/` is a cloud-config merging tool with the following features:

1. Create multi-part cloud-config
1. Environment based templating

`secretary` can be used to load local secrets into the enviroment

1. load pass stored secrets
1. load plain environment files

`secretary` and `cissi` can be used to construct powerful cloud-configs
```
secretary --pass automation/node1 --plan node1/env -- cicci users.tpl.yaml node1/base.yaml > node1.cloud-config
```

`piccu` can be used to put these cloud-config files directly into an ubuntu image.

**WARNING** cloud-config is not a safe way to stay secrets. As such it might be preferable to write the image to memory or to disk directly.
piccu injected cloud-config files are gzip encoded which obfuscate the payload.

# Thanks

This software is mostly glue between a for large and complex projects

- [pass](https://www.passwordstore.org/) and [gopass](https://github.com/gopasspw/gopass) offer great password management
- [go-diskfs](https://github.com/diskfs/go-diskfs) and [fuchsia/thinfs](https://pkg.go.dev/go.fuchsia.dev/fuchsia/src/lib/thinfs) offer a way to manipulate disk images and fat32 filesystems without super user privileges or other dependencies
- a pure go [xz](https://github.com/ulikunitz/xz) to extract the downloaded images
- [mvdan.cc/sh](https://github.com/mvdan/sh/) offers shell parsing and execution - used for (encrypted) environment files and shell script checking
- [yaml.v3](https://github.com/go-yaml/yaml/tree/v3) and [jsonschema](github.com/santhosh-tekuri/jsonschema) provide parsing and validation for cloud-config files
- [sprig](https://github.com/Masterminds/sprig) for a comprehensive set of templating functions
- [progressbar](https://github.com/schollz/progressbar) and [readahead](https://github.com/klauspost/readahead) to make the wait times nicer and shorter

And of course all the dependencies of these tools and the golang libraries.