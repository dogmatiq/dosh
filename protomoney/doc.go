// Package protomoney provides basic, low-level mathematical and comparison
// operations for the "well-known" protocol buffers Money type without
// conversion to another data type.
//
// The functions in this package treat a *money.Money value as immutable.
// Allocations are avoided as best possible, and minimal validation is
// performed.
package protomoney
