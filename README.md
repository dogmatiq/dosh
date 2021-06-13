# Dosh

[![Build Status](https://github.com/dogmatiq/dosh/workflows/CI/badge.svg)](https://github.com/dogmatiq/dosh/actions?workflow=CI)
[![Code Coverage](https://img.shields.io/codecov/c/github/dogmatiq/dosh/main.svg)](https://codecov.io/github/dogmatiq/dosh)
[![Latest Version](https://img.shields.io/github/tag/dogmatiq/dosh.svg?label=semver)](https://semver.org)
[![Documentation](https://img.shields.io/badge/go.dev-reference-007d9c)](https://pkg.go.dev/github.com/dogmatiq/dosh)
[![Go Report Card](https://goreportcard.com/badge/github.com/dogmatiq/dosh)](https://goreportcard.com/report/github.com/dogmatiq/dosh)

Dosh is a Go library for representing, manipulating and serializing monetary
values.

The API is centered around the [`Amount`](https://pkg.go.dev/github.com/dogmatiq/dosh#Amount)
type, which is an immutable monetary value in a specific currency. It is based
on Shopspring's arbitrary-precision [`Decimal`](https://pkg.go.dev/github.com/shopspring/decimal#Decimal)
type, making it suitable for mathematical operations.

## Serialization

`Amount` provides serialization logic for text, binary, JSON and protocol
buffers formats.

Both the text and binary formats are "lossless" insofar as they can encode any
amount that can be represented by the internal arbitrary precision decimal.

The JSON and protocol buffers formats are based on the [`google.type.Money`](https://github.com/googleapis/googleapis/blob/master/google/type/money)
"well-known" protocol buffers type, which has a fixed precision of 9 decimal
places.

The JSON representation of an amount is obtained by applying Protocol Buffers'
canonical [JSON mapping rules](https://developers.google.com/protocol-buffers/docs/proto3#json) to the
`Money` type. This encoding is widely used throughout Google's APIs, an example
of which can be [seen here](https://cloud.google.com/channel/docs/reference/rest/Shared.Types/Money).

Dosh also provides the [`protomoney`](https://pkg.go.dev/github.com/dogmatiq/dosh@main/protomoney)
package, which can be used to perform comparisons and basic mathematical
operations on `google.type.Money` values directly, without unmarshaling them to
an `Amount`.
