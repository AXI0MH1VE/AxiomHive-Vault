package q131

import (
	"math"
	"testing"
)

func TestFromFloat64(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected Q131
	}{
		{"Zero", 0.0, Zero},
		{"One", 1.0, MaxValue},
		{"MinusOne", -1.0, MinValue},
		{"Half", 0.5, Q131(0x40000000)},
		{"Quarter", 0.25, Q131(0x20000000)},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FromFloat64(tt.input)
			if result != tt.expected {
				t.Errorf("FromFloat64(%f) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestAddition(t *testing.T) {
	a := FromFloat64(0.5)
	b := FromFloat64(0.3)
	result := a.Add(b)
	expected := FromFloat64(0.8)
	
	// Allow small tolerance due to fixed-point precision
	diff := result.Sub(expected).Abs()
	if diff > Epsilon*10 {
		t.Errorf("0.5 + 0.3 = %f, want ~0.8", result.ToFloat64())
	}
}

func TestSubtraction(t *testing.T) {
	a := FromFloat64(0.7)
	b := FromFloat64(0.3)
	result := a.Sub(b)
	expected := FromFloat64(0.4)
	
	diff := result.Sub(expected).Abs()
	if diff > Epsilon*10 {
		t.Errorf("0.7 - 0.3 = %f, want ~0.4", result.ToFloat64())
	}
}

func TestMultiplication(t *testing.T) {
	a := FromFloat64(0.5)
	b := FromFloat64(0.6)
	result := a.Mul(b)
	expected := FromFloat64(0.3)
	
	diff := result.Sub(expected).Abs()
	if diff > Epsilon*10 {
		t.Errorf("0.5 * 0.6 = %f, want ~0.3", result.ToFloat64())
	}
}

func TestDivision(t *testing.T) {
	a := FromFloat64(0.8)
	b := FromFloat64(0.4)
	result := a.Div(b)
	expected := FromFloat64(2.0) // Will clamp to MaxValue
	
	// Division result exceeds 1.0, should clamp
	if result != MaxValue {
		t.Errorf("0.8 / 0.4 should clamp to MaxValue")
	}
}

func TestSqrt(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected float64
		tolerance float64
	}{
		{"Sqrt(0.25)", 0.25, 0.5, 0.001},
		{"Sqrt(0.5)", 0.5, 0.707, 0.001},
		{"Sqrt(0.81)", 0.81, 0.9, 0.001},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := FromFloat64(tt.input)
			result := input.Sqrt()
			resultFloat := result.ToFloat64()
			
			if math.Abs(resultFloat-tt.expected) > tt.tolerance {
				t.Errorf("Sqrt(%f) = %f, want ~%f", tt.input, resultFloat, tt.expected)
			}
		})
	}
}

func TestDeterminism(t *testing.T) {
	// Test that same calculation produces identical results
	a := FromFloat64(0.123456789)
	b := FromFloat64(0.987654321)
	
	result1 := a.Add(b).Mul(a).Div(b)
	result2 := a.Add(b).Mul(a).Div(b)
	
	if !VerifyDeterminism(result1, result2) {
		t.Errorf("Determinism check failed: results differ")
	}
	
	// Verify hashes match
	hash1 := result1.Hash()
	hash2 := result2.Hash()
	
	if hash1 != hash2 {
		t.Errorf("Hash mismatch: %x != %x", hash1, hash2)
	}
}

func TestVectorOperations(t *testing.T) {
	v1 := NewVector([]float64{0.1, 0.2, 0.3})
	v2 := NewVector([]float64{0.4, 0.5, 0.6})
	
	// Test addition
	sum := v1.Add(v2)
	if len(sum) != 3 {
		t.Errorf("Vector addition length mismatch")
	}
	
	// Test dot product
	dot := v1.Dot(v2)
	expected := FromFloat64(0.1*0.4 + 0.2*0.5 + 0.3*0.6)
	diff := dot.Sub(expected).Abs()
	
	if diff > Epsilon*100 {
		t.Errorf("Dot product = %f, want ~%f", dot.ToFloat64(), expected.ToFloat64())
	}
}

func TestMatrixMultiplication(t *testing.T) {
	m1 := NewMatrix([][]float64{
		{0.1, 0.2},
		{0.3, 0.4},
	})
	
	m2 := NewMatrix([][]float64{
		{0.5, 0.6},
		{0.7, 0.8},
	})
	
	result := m1.Mul(m2)
	
	if result.Rows != 2 || result.Cols != 2 {
		t.Errorf("Matrix multiplication dimension mismatch")
	}
	
	// Verify result[0,0] = 0.1*0.5 + 0.2*0.7 = 0.19
	expected := FromFloat64(0.19)
	actual := result.Get(0, 0)
	diff := actual.Sub(expected).Abs()
	
	if diff > Epsilon*100 {
		t.Errorf("Matrix[0,0] = %f, want ~0.19", actual.ToFloat64())
	}
}

func BenchmarkAddition(b *testing.B) {
	a := FromFloat64(0.5)
	c := FromFloat64(0.3)
	
	for i := 0; i < b.N; i++ {
		_ = a.Add(c)
	}
}

func BenchmarkMultiplication(b *testing.B) {
	a := FromFloat64(0.5)
	c := FromFloat64(0.6)
	
	for i := 0; i < b.N; i++ {
		_ = a.Mul(c)
	}
}

func BenchmarkSqrt(b *testing.B) {
	a := FromFloat64(0.5)
	
	for i := 0; i < b.N; i++ {
		_ = a.Sqrt()
	}
}

func BenchmarkVectorDot(b *testing.B) {
	v1 := NewVector([]float64{0.1, 0.2, 0.3, 0.4, 0.5})
	v2 := NewVector([]float64{0.6, 0.7, 0.8, 0.9, 0.95})
	
	for i := 0; i < b.N; i++ {
		_ = v1.Dot(v2)
	}
}
