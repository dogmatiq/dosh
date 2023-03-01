<div align="center">

# Dosh

A Go module for representing, manipulating and serializing monetary values.

[![Documentation](https://img.shields.io/badge/go.dev-documentation-007d9c?&style=for-the-badge)](https://pkg.go.dev/github.com/dogmatiq/dosh)
[![Latest Version](https://img.shields.io/github/tag/dogmatiq/dosh.svg?&style=for-the-badge&label=semver)](https://github.com/dogmatiq/dosh/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/dogmatiq/dosh/ci.yml?style=for-the-badge&branch=main)](https://github.com/dogmatiq/dosh/actions/workflows/ci.yml)
[![Code Coverage](https://img.shields.io/codecov/c/github/dogmatiq/dosh/main.svg?style=for-the-badge)](https://codecov.io/github/dogmatiq/dosh)

</div>

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

## Caveats

Google's `money` package does _not_ include the source `.proto` file used to
generate the [`money.Money`](https://pkg.go.dev/google.golang.org/genproto/googleapis/type/money#Money)
type. This makes it difficult to use the `Money` type in user-defined protocol
buffers messages. For this reason, [the original `.proto` file is included in
Dosh's `protomoney` package](protomoney/money.proto).
