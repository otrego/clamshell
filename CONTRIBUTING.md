# Contributing to Clamshell

This guide will walk you through contributing to Clamshell, and assumes you
have some experience with the game Go, but might be relatively new to
programming or to the tools/technologies.

## Bugs, Features, Questions, Comments

If you have any bugs, feature requests, questions, or comments, please file
them on our Issue on our [github tracker](https://github.com/otrego/clamshell/issues)

## Getting Started on Go

First, if you're new to Go, the language, check out:

1. [Installing Go](https://golang.org/doc/install)
2. [The Go Tour](https://tour.golang.org/welcome/1)
3. [How to Write Go Code](https://golang.org/doc/code.html)

## Setting Up Your Dev Environment

We expect all developers, both external contributors and maintainers to
interact with Clamshell via a [fork-and-pull-request
model](https://help.github.com/en/github/getting-started-with-github/fork-a-repo).
So generally, developers will fork the Clamshell repository and then submit
pull requests to the primary `otrego/clamshell` repository.

In general, once you install the Go programming language, you should set
`GOPATH` to something that suits your tastes. I (kashomon) set `GOPATH` to be
`$HOME/inprogress/go`, but feel free to set it as you wish. For more about
setting GOPATH, check out the [golang
wiki](https://github.com/golang/go/wiki/SettingGOPATH).

Assuming you have created a fork of Clamshell (see Development Model), then
create the relevant directories & clone the repo:

```
mkdir -p $GOPATH/src/github.com/otrego
git clone github.com/<USERNAME>/clamshell $GOPATH/src/github.com/otrego/clamshell
```

Then, make sure it builds!

```
$GOPATH/src/github.com/otrego/clamshell
go test ./...
```
