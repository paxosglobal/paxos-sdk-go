package paxos

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	shopspring "github.com/shopspring/decimal"
)

const maxDecimalExponent = 1000

func checkDecimalMagnitude(d shopspring.Decimal) error {
	if d.Sign() == 0 {
		return nil
	}
	c := d.Coefficient()
	digits := int64(len(c.Abs(c).String()))
	adjusted := digits - 1 + int64(d.Exponent())
	if adjusted > maxDecimalExponent || adjusted < -maxDecimalExponent {
		return fmt.Errorf("decimal magnitude must be within 10^+-%d", maxDecimalExponent)
	}
	return nil
}

// Decimal is an arbitrary-precision decimal with JSON serialization that
// marshals to/from JSON strings. This matches the Paxos API convention where
// monetary amounts are transmitted as string values (e.g., "1.50000000")
// to preserve precision.
//
// Use NewDecimal or NewDecimalFromString to create Decimal values:
//
//	amount, err := paxos.NewDecimalFromString("1.50")
//	fee := paxos.NewDecimal(0) // zero
//
// Arithmetic such as Add, Sub, Mul and Div returns a Decimal, so results can be
// passed to the SDK directly:
//
//	total := amount.Add(fee)
type Decimal struct {
	d shopspring.Decimal
}

// NewDecimal creates a Decimal from an int64 value.
func NewDecimal(value int64) Decimal {
	return Decimal{d: shopspring.NewFromInt(value)}
}

// NewDecimalFromString creates a Decimal from a string representation.
// Returns an error if the string is not a valid decimal number.
func NewDecimalFromString(s string) (Decimal, error) {
	d, err := shopspring.NewFromString(s)
	if err != nil {
		return Decimal{}, err
	}
	if err := checkDecimalMagnitude(d); err != nil {
		return Decimal{}, err
	}
	return Decimal{d: d}, nil
}

// MustDecimal creates a Decimal from a string, panicking on error.
// Use only for compile-time constants or test fixtures.
func MustDecimal(s string) Decimal {
	d, err := NewDecimalFromString(s)
	if err != nil {
		panic("paxos.MustDecimal: " + err.Error())
	}
	return d
}

// String returns the plain-notation decimal with trailing zeros trimmed.
func (d Decimal) String() string { return d.d.String() }

// Cmp returns -1, 0 or 1 as d is less than, equal to or greater than o.
func (d Decimal) Cmp(o Decimal) int { return d.d.Cmp(o.d) }

// Equal reports whether d and o are numerically equal, ignoring scale.
func (d Decimal) Equal(o Decimal) bool { return d.d.Equal(o.d) }

// IsZero reports whether d is zero.
func (d Decimal) IsZero() bool { return d.d.IsZero() }

// Sign returns -1, 0 or 1 as d is negative, zero or positive.
func (d Decimal) Sign() int { return d.d.Sign() }

// Add returns d + o.
func (d Decimal) Add(o Decimal) Decimal { return Decimal{d: d.d.Add(o.d)} }

// Sub returns d - o.
func (d Decimal) Sub(o Decimal) Decimal { return Decimal{d: d.d.Sub(o.d)} }

// Mul returns d * o.
func (d Decimal) Mul(o Decimal) Decimal { return Decimal{d: d.d.Mul(o.d)} }

// Div returns d / o rounded to precision digits after the decimal point,
// rounding half away from zero. It panics if o is zero.
func (d Decimal) Div(o Decimal, precision int32) Decimal {
	return Decimal{d: d.d.DivRound(o.d, precision)}
}

// Neg returns -d.
func (d Decimal) Neg() Decimal { return Decimal{d: d.d.Neg()} }

// Abs returns the absolute value of d.
func (d Decimal) Abs() Decimal { return Decimal{d: d.d.Abs()} }

// Round returns d rounded half away from zero to places digits after the decimal point.
func (d Decimal) Round(places int32) Decimal { return Decimal{d: d.d.Round(places)} }

// Truncate returns d truncated toward zero to places digits after the decimal point.
func (d Decimal) Truncate(places int32) Decimal { return Decimal{d: d.d.Truncate(places)} }

// StringFixed returns d rounded to places digits after the decimal point and
// formatted with exactly that many digits.
func (d Decimal) StringFixed(places int32) string { return d.d.StringFixed(places) }

// IntPart returns the integer part of d, truncated toward zero.
func (d Decimal) IntPart() int64 { return d.d.IntPart() }

// InexactFloat64 returns the nearest float64 to d. The conversion can lose
// precision; do not use it for money arithmetic.
func (d Decimal) InexactFloat64() float64 { return d.d.InexactFloat64() }

// MarshalText returns the plain-notation string form of d, for use in text
// encoders and query parameters.
func (d Decimal) MarshalText() ([]byte, error) {
	if err := checkDecimalMagnitude(d.d); err != nil {
		return nil, err
	}
	return []byte(d.d.String()), nil
}

// UnmarshalText parses a decimal from its string form.
func (d *Decimal) UnmarshalText(text []byte) error {
	parsed, err := NewDecimalFromString(string(text))
	if err != nil {
		return err
	}
	*d = parsed
	return nil
}

// MarshalBinary implements encoding.BinaryMarshaler, so encoding/gob can encode
// a Decimal. The encoding is the plain-notation string form of d.
func (d Decimal) MarshalBinary() ([]byte, error) { return d.MarshalText() }

// UnmarshalBinary implements encoding.BinaryUnmarshaler and decodes the output
// of MarshalBinary.
func (d *Decimal) UnmarshalBinary(data []byte) error { return d.UnmarshalText(data) }

// Scan implements sql.Scanner for string, []byte, integer and float column values.
func (d *Decimal) Scan(value any) error {
	var parsed shopspring.Decimal
	if err := parsed.Scan(value); err != nil {
		return err
	}
	if err := checkDecimalMagnitude(parsed); err != nil {
		return err
	}
	d.d = parsed
	return nil
}

// Value implements driver.Valuer, storing d as its string form.
func (d Decimal) Value() (driver.Value, error) {
	if err := checkDecimalMagnitude(d.d); err != nil {
		return nil, err
	}
	return d.d.String(), nil
}

// MarshalJSON serializes the Decimal as a JSON string (e.g., "1.5"; trailing zeros are trimmed).
// This matches the Paxos API wire format.
func (d Decimal) MarshalJSON() ([]byte, error) {
	if d.d.Sign() == 0 {
		return []byte(`"0"`), nil
	}
	if err := checkDecimalMagnitude(d.d); err != nil {
		return nil, err
	}
	return json.Marshal(d.d.String())
}

// UnmarshalJSON deserializes a Decimal from either a JSON string ("1.50")
// or a JSON number (1.50). The string format is the Paxos API convention.
func (d *Decimal) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	// Try as JSON string first (Paxos API format).
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		parsed, err := shopspring.NewFromString(s)
		if err != nil {
			return err
		}
		if err := checkDecimalMagnitude(parsed); err != nil {
			return err
		}
		d.d = parsed
		return nil
	}

	var parsed shopspring.Decimal
	if err := parsed.UnmarshalJSON(data); err != nil {
		return err
	}
	if err := checkDecimalMagnitude(parsed); err != nil {
		return err
	}
	d.d = parsed
	return nil
}
