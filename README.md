# Just some random matchers 

[![GitHub Releases](https://img.shields.io/github/v/release/nhatthm/go-matcher)](https://github.com/nhatthm/go-matcher/releases/latest)
[![Build Status](https://github.com/nhatthm/go-matcher/actions/workflows/test.yaml/badge.svg)](https://github.com/nhatthm/go-matcher/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/nhatthm/go-matcher/branch/master/graph/badge.svg?token=eTdAgDE2vR)](https://codecov.io/gh/nhatthm/go-matcher)
[![Go Report Card](https://goreportcard.com/badge/github.com/nhatthm/go-matcher)](https://goreportcard.com/report/github.com/nhatthm/go-matcher)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/github.com/nhatthm/go-matcher)
[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

The package provides a matcher interface to match a given value of any types. 

## Prerequisites

- `Go >= 1.16`

## Install

```bash
go get github.com/nhatthm/go-matcher
```

## Usage

You could use it in a test or anywhere that needs a value matcher.

```go
package mypackage

import (
	"testing"

	"github.com/nhatthm/go-matcher"
	"github.com/stretchr/testify/assert"
)

func TestValue(t *testing.T) {
	m := matcher.Exact("foobar")
	actual := "FOOBAR"

	assert.True(t, m.Match(actual), "got: %s, want: %s", actual, m.Expected())
}

```

## Donation

If this project help you reduce time to develop, you can give me a cup of coffee :)

### Paypal donation

[![paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;or scan this

<img src="https://user-images.githubusercontent.com/1154587/113494222-ad8cb200-94e6-11eb-9ef3-eb883ada222a.png" width="147px" />
