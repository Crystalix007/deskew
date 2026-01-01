# Deskew

Small utility to deskew images of textual documents.

## Building

This links with the [`leptonica` library][1] via CGO, so make sure to install
that first.

Then install the latest version with

```shell
$ go install github.com/Crystalix007/deskew@latest
```

## Running

```shell
$ deskew [opts...] <input-image>
```

[1]: https://github.com/DanBloomberg/leptonica
