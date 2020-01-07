# symwalk [![Build Status](https://travis-ci.com/xyproto/symwalk.svg?branch=master)](https://travis-ci.com/xyproto/symwalk) [![GoDoc](https://godoc.org/github.com/xyproto/symwalk?status.svg)](http://godoc.org/github.com/xyproto/symwalk)

Concurrently search directories while also following symlinks.

## Fork and license info

* `walker.go` and `walker_test.go` are based on [powerwalk](https://github.com/stretchr/powerwalk) (MIT license), but with added support for traversing symlinks that points to directories too.
* `modwalk.go` is based on `path/filepath` from the Go standard library (BSD license).
* This project is licensed under the MIT license. See [LICENSE](LICENSE) for more info.

## General info

* Version: 1.0.0

