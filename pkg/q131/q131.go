// Package q131 implements Q1.31 fixed-point arithmetic for deterministic financial calculations.
// Q1.31 format: 1 sign bit + 31 fractional bits
// Precision: ~0.0000000005 (2^-31)
// Range: [-1.0, +1.0)
// Guarantees bit-exact reproducibility across all hardware platforms.
package q131

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

// Q131 represents a Q1.31 fixed-point number.
// The underlying int32 stores the fractional value with 31 bits of precision.
type Q131 int32

const (
	// FracBits is the number of fractional bits in Q1.31 format
	FracBits = 31

	// Scale is 2^31 used for conversion
	Scale = 1 << FracBits

	// MaxValue is the maximum representable value (approaching 1.0)
	MaxValue = Q131(0x7FFFFFFF)

	// MinValue is the minimum representable value (-1.0)
	MinValue = Q131(-0x80000000)

	// One represents 1.0 in Q1.31 format
	One = Q131(0x7FFFFFFF)

	// Zero represents 0.0 in Q1.31 format
	Zero = Q131(0)

	// Epsilon is the smallest representable positive value
	Epsilon = Q131(1)
)

// FromFloat64 converts a float64 to Q1.31 format.
// Input must be in range [-1.0, 1.0) or it will be clamped.
func FromFloat64(f float64) Q131 {
	// Clamp to valid range
	if f >= 1.0 {
		return MaxValue
	}
	if f <= -1.0 {
		return MinValue
	}

	// Convert to Q1.31
	scaled := f * float64(Scale)
	return Q131(int32(scaled))
}

// ToFloat64 converts Q1.31 to float64 for display purposes only.
// WARNING: This conversion is lossy and should NOT be used for calculations.
func (q Q131) ToFloat64() float64 {
	return float64(q) / float64(Scale)
}

// Add performs deterministic addition with overflow protection.
func (q Q131) Add(other Q131) Q131 {
	result := int64(q) + int64(other)

	// Clamp on overflow
	if result > int64(MaxValue) {
		return MaxValue
	}
	if result < int64(MinValue) {
		return MinValue
	}

	return Q131(result)
}

// Sub performs deterministic subtraction with underflow protection.
func (q Q131) Sub(other Q131) Q131 {
	result := int64(q) - int64(other)

	// Clamp on underflow
	if result > int64(MaxValue) {
		return MaxValue
	}
	if result < int64(MinValue) {
		return MinValue
	}

	return Q131(result)
}

// Mul performs deterministic multiplication.
// Result = (q * other) / 2^31
func (q Q131) Mul(other Q131) Q131 {
	// Use 64-bit intermediate to prevent overflow
	result := (int64(q) * int64(other)) >> FracBits

	// Clamp to valid range
	if result > int64(MaxValue) {
		return MaxValue
	}
	if result < int64(MinValue) {
		return MinValue
	}

	return Q131(result)
}

// Div performs deterministic division.
// Result = (q * 2^31) / other
func (q Q131) Div(other Q131) Q131 {
	if other == Zero {
		// Division by zero returns max/min based on sign
		if q >= Zero {
			return MaxValue
		}
		return MinValue
	}

	// Use 64-bit intermediate to maintain precision
	result := (int64(q) << FracBits) / int64(other)

	// Clamp to valid range
	if result > int64(MaxValue) {
		return MaxValue
	}
	if result < int64(MinValue) {
		return MinValue
	}

	return Q131(result)
}

// Neg returns the negation of the value.
func (q Q131) Neg() Q131 {
	if q == MinValue {
		return MaxValue // -(-1.0) = 1.0 (clamped)
	}
	return -q
}

// Abs returns the absolute value.
func (q Q131) Abs() Q131 {
	if q < Zero {
		return q.Neg()
	}
	return q
}

// Sqrt computes the square root using integer Newton-Raphson method.
// For Q1.31 format: q represents the value q/2^31.
// sqrt(q/2^31) in Q1.31 = sqrt(q/2^31) * 2^31 = sqrt(q) * 2^15.5 = sqrt(q << 31).
// So we compute isqrt(val << 31) where val = uint64(q).
func (q Q131) Sqrt() Q131 {
	if q <= Zero {
		return Zero
	}
	if q == MaxValue || q == One {
		return MaxValue // sqrt(1) = 1
	}

	// val is the raw Q1.31 integer representation
	val := uint64(q)

	// We want sqrt(q/2^31) expressed as Q1.31.
	// In Q1.31: result = sqrt(val/2^31) * 2^31 = sqrt(val * 2^31) = isqrt(val << 31).
	// Compute val << 31 carefully to avoid overflow (val < 2^31, so val<<31 < 2^62, fits uint64).
	scaled := val << 31

	// Use a good initial estimate: start near sqrt(val) using bit-length.
	// Find the most significant bit of val to get a close initial estimate.
	x := uint64(1)
	for x*x < val {
		x <<= 1
	}
	// x is now >= sqrt(val), refine with Newton-Raphson on scaled
	// Re-initialize x as an estimate for sqrt(scaled)
	x = x << 16 // x ~ sqrt(val) * 2^16, which approximates sqrt(val << 31) since 2^16 ~ 2^15.5
	for {
		xNext := (x + scaled/x) >> 1
		if xNext >= x {
			break
		}
		x = xNext
	}

	// Clamp to valid Q1.31 range
	if x > uint64(MaxValue) {
		return MaxValue
	}

	return Q131(x)
}

// Exp computes e^q using Taylor series expansion.
// Deterministic implementation for range [-1, 1].
func (q Q131) Exp() Q131 {
	// e^x = 1 + x + x^2/2! + x^3/3! + ...
	// Compute first 15 terms for sufficient precision

	result := One
	term := One

	for n := 1; n <= 15; n++ {
		term = term.Mul(q).Div(FromFloat64(float64(n)))
		result = result.Add(term)

		// Early termination if term becomes negligible
		if term.Abs() <= Epsilon {
			break
		}
	}

	return result
}

// Ln computes natural logarithm using series expansion.
// Valid for q in (0, 1]. Deterministic implementation.
func (q Q131) Ln() Q131 {
	if q <= Zero {
		return MinValue // ln(0) = -inf
	}

	// Use ln(x) = 2 * atanh((x-1)/(x+1)) for better convergence
	numerator := q.Sub(One)
	denominator := q.Add(One)

	if denominator == Zero {
		return Zero
	}

	z := numerator.Div(denominator)

	// atanh(z) = z + z^3/3 + z^5/5 + z^7/7 + ...
	result := Zero
	zSquared := z.Mul(z)
	term := z

	for n := 1; n <= 15; n += 2 {
		result = result.Add(term.Div(FromFloat64(float64(n))))
		term = term.Mul(zSquared)

		if term.Abs() <= Epsilon {
			break
		}
	}

	// Multiply by 2
	return result.Add(result)
}

// Compare returns -1 if q < other, 0 if q == other, 1 if q > other.
func (q Q131) Compare(other Q131) int {
	if q < other {
		return -1
	}
	if q > other {
		return 1
	}
	return 0
}

// Hash returns a SHA-256 hash of the Q1.31 value for verification.
func (q Q131) Hash() [32]byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(q))
	return sha256.Sum256(buf)
}

// String returns a string representation for debugging.
func (q Q131) String() string {
	return fmt.Sprintf("Q1.31(%d = %.10f)", int32(q), q.ToFloat64())
}

// Vector represents a vector of Q1.31 values for batch operations.
type Vector []Q131

// NewVector creates a vector from float64 slice.
func NewVector(values []float64) Vector {
	v := make(Vector, len(values))
	for i, val := range values {
		v[i] = FromFloat64(val)
	}
	return v
}

// Add performs element-wise addition.
func (v Vector) Add(other Vector) Vector {
	if len(v) != len(other) {
		panic("vector dimension mismatch")
	}

	result := make(Vector, len(v))
	for i := range v {
		result[i] = v[i].Add(other[i])
	}
	return result
}

// Mul performs element-wise multiplication.
func (v Vector) Mul(other Vector) Vector {
	if len(v) != len(other) {
		panic("vector dimension mismatch")
	}

	result := make(Vector, len(v))
	for i := range v {
		result[i] = v[i].Mul(other[i])
	}
	return result
}

// Dot computes the dot product of two vectors.
func (v Vector) Dot(other Vector) Q131 {
	if len(v) != len(other) {
		panic("vector dimension mismatch")
	}

	result := Zero
	for i := range v {
		result = result.Add(v[i].Mul(other[i]))
	}
	return result
}

// Scale multiplies all elements by a scalar.
func (v Vector) Scale(scalar Q131) Vector {
	result := make(Vector, len(v))
	for i := range v {
		result[i] = v[i].Mul(scalar)
	}
	return result
}

// Hash returns a SHA-256 hash of the entire vector for verification.
func (v Vector) Hash() [32]byte {
	h := sha256.New()
	for _, val := range v {
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(val))
		h.Write(buf)
	}
	var result [32]byte
	copy(result[:], h.Sum(nil))
	return result
}

// Matrix represents a matrix of Q1.31 values.
type Matrix struct {
	Rows int
	Cols int
	Data []Q131
}

// NewMatrix creates a matrix from float64 2D slice.
func NewMatrix(values [][]float64) *Matrix {
	if len(values) == 0 {
		return &Matrix{Rows: 0, Cols: 0, Data: []Q131{}}
	}

	rows := len(values)
	cols := len(values[0])
	data := make([]Q131, rows*cols)

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			data[i*cols+j] = FromFloat64(values[i][j])
		}
	}

	return &Matrix{
		Rows: rows,
		Cols: cols,
		Data: data,
	}
}

// Get returns the element at (row, col).
func (m *Matrix) Get(row, col int) Q131 {
	return m.Data[row*m.Cols+col]
}

// Set sets the element at (row, col).
func (m *Matrix) Set(row, col int, value Q131) {
	m.Data[row*m.Cols+col] = value
}

// Mul performs matrix multiplication.
func (m *Matrix) Mul(other *Matrix) *Matrix {
	if m.Cols != other.Rows {
		panic("matrix dimension mismatch")
	}

	result := &Matrix{
		Rows: m.Rows,
		Cols: other.Cols,
		Data: make([]Q131, m.Rows*other.Cols),
	}

	for i := 0; i < m.Rows; i++ {
		for j := 0; j < other.Cols; j++ {
			sum := Zero
			for k := 0; k < m.Cols; k++ {
				sum = sum.Add(m.Get(i, k).Mul(other.Get(k, j)))
			}
			result.Set(i, j, sum)
		}
	}

	return result
}

// Hash returns a SHA-256 hash of the entire matrix for verification.
func (m *Matrix) Hash() [32]byte {
	h := sha256.New()
	for _, val := range m.Data {
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(val))
		h.Write(buf)
	}
	var result [32]byte
	copy(result[:], h.Sum(nil))
	return result
}

// VerifyDeterminism checks that two calculations produce identical results.
func VerifyDeterminism(q1, q2 Q131) bool {
	return q1 == q2 && q1.Hash() == q2.Hash()
}
