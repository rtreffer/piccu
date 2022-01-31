# thinfs shim to hide fuchsia

fuchsia is not designed to be built as part of normal go projects. However nothing breaks thinfs in such a context. It is a very self-contained package.

This subpath can only be used via a go.mod substitute:
```
replace go.fuchsia.dev/fuchsia/src/lib/thinfs/ => ./pkg/thinfs
```