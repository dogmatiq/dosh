package dosh

// Floor returns an amount with a magnitude equal to the nearest integer less
// than or equal to a.Magnitude().
func (a Amount) Floor() Amount {
	a.mag = a.mag.Floor()
	return a
}

// Ceil returns an amount with a magnitude equal to the the nearest integer
// greater than or equal to a.Magnitude().
func (a Amount) Ceil() Amount {
	a.mag = a.mag.Ceil()
	return a
}

// Truncate returns the amount truncated to n decimal places without performing
// any rounding.
//
// n must be positive.
func (a Amount) Truncate(n int32) Amount {
	a.mag = a.mag.Truncate(n)
	return a
}

// Round returns the amount rounded to n decimal places, rounded away from zero,
// also known as "commercial rounding".
//
// See https://en.wikipedia.org/wiki/Rounding#Round_half_away_from_zero.
//
// If n is negative the result is rounded to the -n'th integer place. For
// example, an amount with a magnitude of 543 rounded to -1 places results in an
// amount with a magnitude of 540.
func (a Amount) Round(n int32) Amount {
	a.mag = a.mag.Round(n)
	return a
}

// RoundBank returns the amount rounded to n decimal places, rounded towards
// even digits, also known as "banker's rounding".
//
// Banker's rounding is commonly used for rounding sales tax amounts.
// See https://wiki.c2.com/?BankersRounding.
//
// If n is negative the result is rounded to the -n'th integer place. For
// example, an amount with a magnitude of 543 rounded to -1 places results in an
// amount with a magnitude of 540.
func (a Amount) RoundBank(n int32) Amount {
	a.mag = a.mag.RoundBank(n)
	return a
}
