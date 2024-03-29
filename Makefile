.PHONY: default
default: help

.PHONY: help
help:
	@echo
	@echo "targets:"
	@echo "     schema     build the json schema from cloud-init sources"
	@echo "     build      build the piccu binary"
	@echo

pkg/cicci/cloud-config.schema.json:
	cd _external/github.com/canonical/cloud-init/ && \
	python3 -c 'import json; from cloudinit.config.schema import *; print(json.dumps(get_schema()))' | \
	  jq . > "../../../../$@.tmp"
	mv "$@.tmp" "$@"

schema: pkg/cicci/cloud-config.schema.json

.PHONY: build
build: piccu

.PHONY: piccu
piccu:
	go build github.com/rtreffer/piccu/cmd/piccu
