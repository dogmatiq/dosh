package dosh

// Ceil returns an amount with the the nearest integer value less than or equal
// to a.
func (a Amount) Floor() Amount {
	a.magnitude = a.magnitude.Floor()
	return a
}

// Ceil returns an amount with the nearest integer value greater than or equal
// to a.
func (a Amount) Ceil() Amount {
	a.magnitude = a.magnitude.Ceil()
	return a
}

// Round returns the amount rounded to n decimal places.
//
// If n is negative the result is rounded to the -n'th integer place. For
// example, an amount with a magnitude of 543 rounded to -1 places results in an
// amount with a magnitude of 540.
func (a Amount) Round(n int32) Amount {
	a.magnitude = a.magnitude.Round(n)
	return a
}

// RoundBank returns the amount rounded to n decimal places, rounded towards
// even digits, also known as "banker's rounding".
//
// See https://wiki.c2.com/?BankersRounding.
//
// If n is negative the result is rounded to the -n'th integer place. For
// example, an amount with a magnitude of 543 rounded to -1 places results in an
// amount with a magnitude of 540.
func (a Amount) RoundBank(n int32) Amount {
	a.magnitude = a.magnitude.RoundBank(n)
	return a
}

// Truncate returns the amount truncated to n decimal places without performing
// any rounding.
//
// n must be positive.
func (a Amount) Truncate(n int32) Amount {
	a.magnitude = a.magnitude.Truncate(n)
	return a
}
