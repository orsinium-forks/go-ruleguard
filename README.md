# go-ruleguard

[![Build Status](https://travis-ci.com/quasilyte/go-ruleguard.svg?branch=master)](https://travis-ci.com/quasilyte/go-ruleguard)
[![GoDoc](https://godoc.org/github.com/quasilyte/go-ruleguard/dsl?status.svg)](https://godoc.org/github.com/quasilyte/go-ruleguard)
[![Go Report Card](https://goreportcard.com/badge/github.com/quasilyte/go-ruleguard)](https://goreportcard.com/report/github.com/quasilyte/go-ruleguard)

![Logo](docs/logo2.png)

## Overview

[analysis](https://godoc.org/golang.org/x/tools/go/analysis)-based Go linter that runs dynamically loaded rules.

You write the rules, `ruleguard` checks whether they are satisfied.

`ruleguard` has some similarities with [GitHub CodeQL](https://securitylab.github.com/tools/codeql), but only focuses on Go code queries.

**Features:**

* Custom linting rules without re-compilation and Go plugins.
* Diagnostics are written in a declarative way.
* [Quickfix](docs/gorules.md#suggestions-quickfix-support) actions support.
* Powerful match filtering features, like expression [type pattern matching](docs/gorules.md#type-pattern-matching).
* Remote configuration. Reuse rules between projects, use ruleguard as zero-configuration linter.

`ruleguard` comes with [rules.go](rules.go) file that can be used as linter or as an example.

It can also be easily embedded into other static analyzers. [go-critic](https://github.com/go-critic/go-critic) can be used as an example.

## Quick start

To install `ruleguard` binary under your `$(go env GOPATH)/bin`:

```bash
$ go get -v -u github.com/quasilyte/go-ruleguard/...
```

If `$GOPATH/bin` is under your system `$PATH`, `ruleguard` command should be available after that.<br>

```bash
$ ruleguard -help
ruleguard: execute dynamic gogrep-based rules

Usage: ruleguard [-flag] [package]

Flags:
  -rules string
        comma-separated list of gorule file paths or URLs
  -e string
        execute a single rule from a given string
  -fix
        apply all suggested fixes
  -c int
        display offending line with this many lines of context (default -1)
  -json
        emit JSON output
```


Create a test `example.go` target file:

```go
package main

func main() {
	var v1, v2 int
	println(!(v1 != v2))
	println(!(v1 == v2))
	if v1 == 0 && v1 == 0 {
		println("hello, world!")
	}
}
```

Run `ruleguard` on that target file:

```bash
$ ruleguard \
    -rules https://raw.githubusercontent.com/quasilyte/go-ruleguard/master/rules.go \
    -fix example.go
example.go:5:10: suggestion: v1 == v2
example.go:6:10: suggestion: v1 != v2
example.go:7:5: suspicious identical LHS and RHS
```

The `-rules` option can be a path to a local file with rules, an URL, or comma-separated list of paths and URLs in any combination.

Since we ran `ruleguard` with `-fix` argument, both **suggested** changes are applied to `example.go`.

## Using custom rules

Create a test `example.rules.go` file:

```go
// +build ignore

package gorules

import "github.com/quasilyte/go-ruleguard/dsl/fluent"

func _(m fluent.Matcher) {
	m.Match(`$x || $x`,
		`$x && $x`).
		Where(m["x"].Pure).
		Report(`suspicious identical LHS and RHS`)

	m.Match(`!($x != $y)`).Suggest(`$x == $y`)
	m.Match(`!($x == $y)`).Suggest(`$x != $y`)
}
```

Run `ruleguard` on that target file:

```bash
$ ruleguard -rules example.rules.go -fix example.go
example.go:5:10: hint: suggested: v1 == v2
example.go:6:10: hint: suggested: v1 != v2
example.go:7:5: error: suspicious identical LHS and RHS
```

There is also a `-e` mode that is useful during pattern debugging:

```bash
$ ruleguard -e 'm.Match(`!($x != $y)`)' example.go
example.go:5:10: !(v1 != v2)
```

It automatically inserts `Report("$$")` into the specified pattern.

## How it works

`ruleguard` parses [gorules](docs/gorules.md) (e.g. `rules.go`) during the start to load the rule set.
Loaded rules are then used to check the specified targets (Go files, packages).
The `rules.go` file itself is never compiled, nor executed.

A `rules.go` file, as interpreted by a [`dsl/fluent`](https://godoc.org/github.com/quasilyte/go-ruleguard/dsl/fluent) API, is a set of functions that serve as a rule groups. Every function accepts a single [`fluent.Matcher`](https://godoc.org/github.com/quasilyte/go-ruleguard/dsl/fluent#Matcher) argument that is then used to define and configure rules inside the group.

A rule definition always starts from a [`Match(patterns...)`](https://godoc.org/github.com/quasilyte/go-ruleguard/dsl/fluent#Matcher.Match) method call and ends with a [`Report(message)`](https://godoc.org/github.com/quasilyte/go-ruleguard/dsl/fluent#Matcher.Report) method call.

There can be additional calls in between these two. For example, a [`Where(cond)`](https://godoc.org/github.com/quasilyte/go-ruleguard/dsl/fluent#Matcher.Where) call applies constraints to a match to decide whether its accepted or rejected. So even if there is a match for a pattern, it won't produce a report message unless it satisfies a `Where()` condition.

To learn more, check out the documentation and/or the source code.

## Documentation

* Example rule files: [rules.go](rules.go)
* Another great example: [github.com/dgryski/semgrep-go/ruleguard.rules.go](https://github.com/dgryski/semgrep-go/blob/master/ruleguard.rules.go)
* [gorules](docs/gorules.md) format documentation
* [dsl/fluent package](https://godoc.org/github.com/quasilyte/go-ruleguard/dsl/fluent) reference
* [ruleguard package](https://godoc.org/github.com/quasilyte/go-ruleguard/ruleguard) reference
* Introduction article: [EN](https://quasilyte.dev/blog/post/ruleguard/), [RU](https://habr.com/ru/post/481696/)

## Ready to use rules

There is the list of open-source rules you can use in your project:

* [go-ruleguard/rules.go](https://github.com/quasilyte/go-ruleguard/blob/master/rules.go)
* [semgrep-go/ruleguard.rules.go](https://github.com/dgryski/semgrep-go/blob/master/ruleguard.rules.go)
* [go-critic/rules.go](https://github.com/go-critic/go-critic/blob/master/rules.go)

## Extra references

* Online ruleguard playground: [go-ruleguard.github.io/play](https://go-ruleguard.github.io/play)
* [gogrep](https://github.com/mvdan/gogrep) - underlying AST matching engine
* [NoVerify: Dynamic Rules for Static Analysis](https://medium.com/@vktech/noverify-dynamic-rules-for-static-analysis-8f42859e9253)
* [Ruleguard comparison with Semgrep and CodeQL](https://speakerdeck.com/quasilyte/ruleguard-vs-semgrep-vs-codeql)
