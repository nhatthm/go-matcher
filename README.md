# Just some random matchers 

[![GitHub Releases](https://img.shields.io/github/v/release/nhatthm/go-matcher)](https://github.com/nhatthm/go-matcher/releases/latest)
[![Build Status](https://github.com/nhatthm/go-matcher/actions/workflows/test.yaml/badge.svg)](https://github.com/nhatthm/go-matcher/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/nhatthm/go-matcher/branch/master/graph/badge.svg?token=eTdAgDE2vR)](https://codecov.io/gh/nhatthm/go-matcher)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/go.nhat.io/matcher/v3)
[![Donate](https://img.shields.io/badge/%20-Donate-%20?style=flat&logo=githubsponsors&color=E5E4E2)](http://donate.nhat.me)

The package provides a matcher interface to match a given value of any types. 

## Prerequisites

- `Go >= 1.23`

## Install

```bash
go get go.nhat.io/matcher/v3
```

## Usage

You could use it in a test or anywhere that needs a value matcher.

```go
package mypackage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.nhat.io/matcher/v3"
)

func TestValue(t *testing.T) {
	m := matcher.Equal("foobar")
	actual := "FOOBAR"

	result, err := m.Match(actual)

	assert.True(t, result, "got: %s, want: %s", actual, m.Expected())
	assert.NoError(t, err)
}

```

## Donation

If this project saved you some development time, buy me a cup of coffee :)

[![donate](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](http://donate.nhat.me)

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;or scan this

<img src="https://github.com/nhatthm/donate.nhat.me/blob/master/images/qr_sponsor.png" width="147px" />

[<sub><sup>[table of contents]</sup></sub>](#table-of-contents)
