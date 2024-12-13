# gitignore

[![Go Documentation](https://godocs.io/git.sr.ht/~jamesponddotco/gitignore-go?status.svg)](https://godocs.io/git.sr.ht/~jamesponddotco/gitignore-go)
[![Go Report Card](https://goreportcard.com/badge/git.sr.ht/~jamesponddotco/gitignore-go)](https://goreportcard.com/report/git.sr.ht/~jamesponddotco/gitignore-go)
[![Coverage Report](https://img.shields.io/badge/coverage-92%25-brightgreen)](https://git.sr.ht/~jamesponddotco/gitignore-go/tree/trunk/item/cover.out)
[![builds.sr.ht status](https://builds.sr.ht/~jamesponddotco/gitignore-go.svg)](https://builds.sr.ht/~jamesponddotco/gitignore-go?)

Package `gitignore` provides a simple way to parse `.gitignore` files
and match paths against the rules defined in those files. The package
tries to implement [the gitignore
specification](https://git-scm.com/docs/gitignore) as defined in the git
documentation, but it is not guaranteed to be 100% accurate.

`gitignore` is kind of a fork of
[`sabhiram/go-gitignore`](https://github.com/sabhiram/go-gitignore), as
a lot of the logic comes from that package, but offers a different
public API, usage of modern Go features, and is actively maintained.

As the package uses the standard `regexp` package for the bulk of its
logic, performance for big files could be improved. [Patches are
welcome](https://lists.sr.ht/~jamesponddotco/gitignore-devel)!

## Installation

To install `gitignore` and use it in your project, run:

```console
go get git.sr.ht/~jamesponddotco/gitignore-go@latest
```

## Usage

To parse a `.gitignore` file and use it to match a path, call
`gitignore.New` with the path to the `.gitignore` file as an argument,
and then call `.Match` with the path to be matched as an argument.

```go
package main

import (
	"log"

	"git.sr.ht/~jamesponddotco/gitignore-go"
)

func main() {
	matcher, err := gitignore.New("/path/to/.gitignore")
	if err != nil {
		log.Fatal(err)
	}

	// Match the given path against the rules defined in the .gitignore file and
	// do something with it.
	if matcher.Match("/path/to/file.txt") {
		log.Println("File is ignored.")
	} else {
		log.Println("File is not ignored.")
	}
}
```

You can also use `gitignore.NewFromLines` to parse a slice of strings
representing the lines of a `.gitignore` file, if you prefer.

## Contributing

Anyone can help make `gitignore` better. Send patches on the [mailing
list](https://lists.sr.ht/~jamesponddotco/gitignore-devel) and report
bugs on the [issue
tracker](https://todo.sr.ht/~jamesponddotco/gitignore).

You must sign-off your work using `git commit --signoff`. Follow the
[Linux kernel developer's certificate of
origin](https://www.kernel.org/doc/html/latest/process/submitting-patches.html#sign-your-work-the-developer-s-certificate-of-origin)
for more details.

All contributions are made under [the MIT License](LICENSE.md).

## Resources

The following resources are available:

- [Package documentation](https://godocs.io/git.sr.ht/~jamesponddotco/gitignore-go).
- [Support and general discussions](https://lists.sr.ht/~jamesponddotco/gitignore-discuss).
- [Patches and development related questions](https://lists.sr.ht/~jamesponddotco/gitignore-devel).
- [Instructions on how to prepare patches](https://git-send-email.io/).
- [Feature requests and bug reports](https://todo.sr.ht/~jamesponddotco/gitignore).

---

Released under the [MIT License](LICENSE.md).
