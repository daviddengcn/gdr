gdr
===

go-debug-run tool

Go is designed not to run when some useless declaration exists, e.g.some imported but not used packages.
Gdr appends some placeholders to the original source to go throught this kind of errors from the compiler,
so that you can focus on more important problems at an early stage of the implementation.

Installation and Run
--------------------
```bash
$ go get -u github.com/daviddengcn/gdr
```
Use `gdr` to replace the `go run`:
```bash
$ gdr yourapp.go
```

Features
--------

* Currently, only unused package errors can be avoided.
* ToDO: To avoid unused variables.
