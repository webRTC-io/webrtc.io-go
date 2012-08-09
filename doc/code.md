[Documents](/doc/) [References](/ref/) [Packages](/pkg/) [The
Project](/project/) [Help](/help/)

[The Go Programming Language](/)

How to Write Go Code
====================

Introduction
------------

This document demonstrates the development of a simple Go package and
introduces the [go command](/cmd/go/), the standard way to fetch, build,
and install Go packages and commands.

Code organization
-----------------

### `GOPATH` and workspaces

One of Go's design goals is to make writing software easier. To that
end, the `go` command doesn't use Makefiles or other configuration files
to guide program construction. Instead, it uses the source code to find
dependencies and determine build conditions. This means your source code
and build scripts are always in sync; they are one and the same.

The one thing you must do is set a `GOPATH` environment variable.
`GOPATH` tells the `go` command (and other related tools) where to find
and install the Go packages on your system.

`GOPATH` is a list of paths. It shares the syntax of your system's
`PATH` environment variable. A typical `GOPATH` on a Unix system might
look like this:

    GOPATH=/home/user/ext:/home/user/mygo

(On a Windows system use semicolons as the path separator instead of
colons.)

Each path in the list (in this case `/home/user/ext` or
`/home/user/mygo`) specifies the location of a *workspace*. A workspace
contains Go source files and their associated package objects, and
command executables. It has a prescribed structure of three
subdirectories:

-   `src` contains Go source files,
-   `pkg` contains compiled package objects, and
-   `bin` contains executable commands.

Subdirectories of the `src` directory hold independent packages, and all
source files (`.go`, `.c`, `.h`, and `.s`) in each subdirectory are
elements of that subdirectory's package.

When building a program that imports the package "`widget`" the `go`
command looks for `src/pkg/widget` inside the Go root, and then—if the
package source isn't found there—it searches for `src/widget` inside
each workspace in order.

Multiple workspaces can offer some flexibility and convenience, but for
now we'll concern ourselves with only a single workspace.

Let's work through a simple example. First, create a `$HOME/mygo`
directory and its `src` subdirectory:

    $ mkdir -p $HOME/mygo/src # create a place to put source code

Next, set it as the `GOPATH`. You should also add the `bin` subdirectory
to your `PATH` environment variable so that you can run the commands
therein without specifying their full path. To do this, add the
following lines to `$HOME/.profile` (or equivalent):

    export GOPATH=$HOME/mygo
    export PATH=$PATH:$HOME/mygo/bin

### Import paths

The standard packages are given short import paths such as `"fmt"` and
`"net/http"` for convenience. For your own projects, it is important to
choose a base import path that is unlikely to collide with future
additions to the standard library or other external libraries.

The best way to choose an import path is to use the location of your
version control repository. For instance, if your source repository is
at `example.com` or `code.google.com/p/example`, you should begin your
package paths with that URL, as in "`example.com/foo/bar`" or
"`code.google.com/p/example/foo/bar`". Using this convention, the `go`
command can automatically check out and build the source code by its
import path alone.

If you don't intend to install your code in this way, you should at
least use a unique prefix like "`widgets/`", as in "`widgets/foo/bar`".
A good rule is to use a prefix such as your company or project name,
since it is unlikely to be used by another group.

We'll use `example/` as our base import path:

    $ mkdir -p $GOPATH/src/example

### Package names

The first statement in a Go source file should be

    package name

where `name` is the package's default name for imports. (All files in a
package must use the same `name`.)

Go's convention is that the package name is the last element of the
import path: the package imported as "`crypto/rot13`" should be named
`rot13`. There is no requirement that package names be unique across all
packages linked into a single binary, only that the import paths (their
full file names) be unique.

Create a new package under `example` called `newmath`:

    $ cd $GOPATH/src/example
    $ mkdir newmath

Then create a file named `$GOPATH/src/example/newmath/sqrt.go`
containing the following Go code:

    // Package newmath is a trivial example package.
    package newmath

    // Sqrt returns an approximation to the square root of x.
    func Sqrt(x float64) float64 {
            // This is a terrible implementation.
            // Real code should import "math" and use math.Sqrt.
            z := 0.0
            for i := 0; i < 1000; i++ {
                    z -= (z*z - x) / (2 * x)
            }
            return z
    }

This package is imported by the path name of the directory it's in,
starting after the `src` component:

    import "example/newmath"

See [Effective Go](/doc/effective_go.html#names) to learn more about
Go's naming conventions.

Building and installing
-----------------------

The `go` command comprises several subcommands, the most central being
`install`. Running `go install importpath` builds and installs a package
and its dependencies.

To "install a package" means to write the package object or executable
command to the `pkg` or `bin` subdirectory of the workspace in which the
source resides.

### Building a package

To build and install the `newmath` package, type

    $ go install example/newmath

This command will produce no output if the package and its dependencies
are built and installed correctly.

As a convenience, the `go` command will assume the current directory if
no import path is specified on the command line. This sequence of
commands has the same effect as the one above:

    $ cd $GOPATH/src/example/newmath
    $ go install

The resulting workspace directory tree (assuming we're running Linux on
a 64-bit system) looks like this:

    pkg/
        linux_amd64/
            example/
                newmath.a  # package object
    src/
        example/
            newmath/
                sqrt.go    # package source

### Building a command

The `go` command treats code belonging to `package main` as an
executable command and installs the package binary to the `GOPATH`'s
`bin` subdirectory.

Add a command named `hello` to the source tree. First create the
`example/hello` directory:

    $ cd $GOPATH/src/example
    $ mkdir hello

Then create the file `$GOPATH/src/example/hello/hello.go` containing the
following Go code.

    // Hello is a trivial example of a main package.
    package main

    import (
            "example/newmath"
            "fmt"
    )

    func main() {
            fmt.Printf("Hello, world.  Sqrt(2) = %v\n", newmath.Sqrt(2))
    }

Next, run `go install`, which builds and installs the binary to
`$GOPATH/bin` (or `$GOBIN`, if set; to simplify presentation, this
document assumes `GOBIN` is unset):

    $ go install example/hello

To run the program, invoke it by name as you would any other command:

    $ $GOPATH/bin/hello
    Hello, world.  Sqrt(2) = 1.414213562373095

If you added `$HOME/mygo/bin` to your `PATH`, you may omit the path to
the executable:

    $ hello
    Hello, world.  Sqrt(2) = 1.414213562373095

The workspace directory tree now looks like this:

    bin/
        hello              # command executable
    pkg/
        linux_amd64/ 
            example/
                newmath.a  # package object
    src/
        example/
            hello/
                hello.go   # command source
            newmath/
                sqrt.go    # package source

The `go` command also provides a `build` command, which is like
`install` except it builds all objects in a temporary directory and does
not install them under `pkg` or `bin`. When building a command an
executable named after the last element of the import path is written to
the current directory. When building a package, `go build` serves merely
to test that the package and its dependencies can be built. (The
resulting package object is thrown away.)

Testing
-------

Go has a lightweight test framework composed of the `go test` command
and the `testing` package.

You write a test by creating a file with a name ending in `_test.go`
that contains functions named `TestXXX` with signature
`func (t *testing.T)`. The test framework runs each such function; if
the function calls a failure function such as `t.Error` or `t.Fail`, the
test is considered to have failed.

Add a test to the `newmath` package by creating the file
`$GOPATH/src/example/newmath/sqrt_test.go` containing the following Go
code.

    package newmath

    import "testing"

    func TestSqrt(t *testing.T) {
        const in, out = 4, 2
        if x := Sqrt(in); x != out {
            t.Errorf("Sqrt(%v) = %v, want %v", in, x, out)
            }
    }

Now run the test with `go test`:

    $ go test example/newmath
    ok      example/newmath 0.165s

Run `go help test` and see the [testing package
documentation](/pkg/testing/) for more detail.

Remote packages
---------------

An import path can describe how to obtain the package source code using
a revision control system such as Git or Mercurial. The `go` command
uses this property to automatically fetch packages from remote
repositories. For instance, the examples described in this document are
also kept in a Mercurial repository hosted at Google Code,
`code.google.com/p/go.example`. If you include the repository URL in the
package's import path, `go get` will fetch, build, and install it
automatically:

    $ go get code.google.com/p/go.example/hello
    $ $GOPATH/bin/hello
    Hello, world.  Sqrt(2) = 1.414213562373095

If the specified package is not present in a workspace, `go get` will
place it inside the first workspace specified by `GOPATH`. (If the
package does already exist, `go get` skips the remote fetch and behaves
the same as `go install`.)

After issuing the above `go get` command, the workspace directory tree
should now now look like this:

    bin/
        hello                 # command executable
    pkg/
        linux_amd64/ 
            code.google.com/p/go.example/
                newmath.a     # package object
            example/
                newmath.a     # package object
    src/
        code.google.com/p/go.example/
            hello/
                hello.go      # command source
            newmath/
                sqrt.go       # package source
                sqrt_test.go  # test source
        example/
            hello/
                hello.go      # command source
            newmath/
                sqrt.go       # package source
                sqrt_test.go  # test source

The `hello` command hosted at Google Code depends on the `newmath`
package within the same repository. The imports in `hello.go` file use
the same import path convention, so the `go get` command is able to
locate and install the dependent package, too.

    import "code.google.com/p/go.example/newmath"

This convention is the easiest way to make your Go packages available
for others to use. The [Go Project
Dashboard](http://godashboard.appspot.com) is a list of external Go
projects including programs and libraries.

For more information on using remote repositories with the `go` command,
see `go help remote`.

Further reading
---------------

See [Effective Go](/doc/effective_go.html) for tips on writing clear,
idiomatic Go code.

Take [A Tour of Go](http://tour.golang.org/) to learn the language
proper.

Visit the [documentation page](/doc/#articles) for a set of in-depth
articles about the Go language and its libraries and tools.

Build version go1.0.2.\
 Except as [noted](http://code.google.com/policies.html#restrictions),
the content of this page is licensed under the Creative Commons
Attribution 3.0 License, and code is licensed under a [BSD
license](/LICENSE).\
 [Terms of Service](/doc/tos.html) | [Privacy
Policy](http://www.google.com/intl/en/privacy/privacy-policy.html)
