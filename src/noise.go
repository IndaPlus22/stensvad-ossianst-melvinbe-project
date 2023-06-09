package main

import (
	"math"
)

var perm = [512]uint8{
	151, 160, 137, 91, 90, 15,
	131, 13, 201, 95, 96, 53, 194, 233, 7, 225, 140, 36, 103, 30, 69, 142, 8, 99, 37, 240, 21, 10, 23,
	190, 6, 148, 247, 120, 234, 75, 0, 26, 197, 62, 94, 252, 219, 203, 117, 35, 11, 32, 57, 177, 33,
	88, 237, 149, 56, 87, 174, 20, 125, 136, 171, 168, 68, 175, 74, 165, 71, 134, 139, 48, 27, 166,
	77, 146, 158, 231, 83, 111, 229, 122, 60, 211, 133, 230, 220, 105, 92, 41, 55, 46, 245, 40, 244,
	102, 143, 54, 65, 25, 63, 161, 1, 216, 80, 73, 209, 76, 132, 187, 208, 89, 18, 169, 200, 196,
	135, 130, 116, 188, 159, 86, 164, 100, 109, 198, 173, 186, 3, 64, 52, 217, 226, 250, 124, 123,
	5, 202, 38, 147, 118, 126, 255, 82, 85, 212, 207, 206, 59, 227, 47, 16, 58, 17, 182, 189, 28, 42,
	223, 183, 170, 213, 119, 248, 152, 2, 44, 154, 163, 70, 221, 153, 101, 155, 167, 43, 172, 9,
	129, 22, 39, 253, 19, 98, 108, 110, 79, 113, 224, 232, 178, 185, 112, 104, 218, 246, 97, 228,
	251, 34, 242, 193, 238, 210, 144, 12, 191, 179, 162, 241, 81, 51, 145, 235, 249, 14, 239, 107,
	49, 192, 214, 31, 181, 199, 106, 157, 184, 84, 204, 176, 115, 121, 50, 45, 127, 4, 150, 254,
	138, 236, 205, 93, 222, 114, 67, 29, 24, 72, 243, 141, 128, 195, 78, 66, 215, 61, 156, 180,
	151, 160, 137, 91, 90, 15,
}

func grad3(hash uint8, x, y, z float64) float64 {
	grad := [16][3]float64{
		{1, 1, 0}, {-1, 1, 0}, {1, -1, 0}, {-1, -1, 0},
		{1, 0, 1}, {-1, 0, 1}, {1, 0, -1}, {-1, 0, -1},
		{0, 1, 1}, {0, -1, 1}, {0, 1, -1}, {0, -1, -1},
		{1, 1, 0}, {-1, 1, 0}, {0, -1, 1}, {0, -1, -1},
	}

	h := hash & 15
	return grad[h][0]*x + grad[h][1]*y + grad[h][2]*z
}

/*
Snoise samples 3D simplex noise and returns a value between -1.0 and 1.0.
Close samples will return similar values.

Parameters:
- x, y, z: Sample coordinates

Returns:
- value: A smooth value between -1.0, 1.0

Example usage:

	sample := Snoise(1.0, 2.0, 3.0)
*/
func Snoise(x, y, z float32) float32 {
	// Offset sampling by seed
	x += seed

	const F3 = 1.0 / 3.0
	const G3 = 1.0 / 6.0

	var n0, n1, n2, n3 float64

	// Calculate fractional and integer part of input coordinates
	s := (x + y + z) * F3
	xs := x + s
	ys := y + s
	zs := z + s
	i := int(math.Floor(float64(xs)))
	j := int(math.Floor(float64(ys)))
	k := int(math.Floor(float64(zs)))

	// Calculate offsets and fractional parts within unit cube
	t := float64(i+j+k) * G3
	X0 := float64(i) - t
	Y0 := float64(j) - t
	Z0 := float64(k) - t
	x0 := float64(x) - X0
	y0 := float64(y) - Y0
	z0 := float64(z) - Z0

	var i1, j1, k1 int
	var i2, j2, k2 int

	// Determine which simplex we are in and set up corresponding indices
	if x0 >= y0 {
		if y0 >= z0 {
			i1, j1, k1, i2, j2, k2 = 1, 0, 0, 1, 1, 0
		} else if x0 >= z0 {
			i1, j1, k1, i2, j2, k2 = 1, 0, 0, 1, 0, 1
		} else {
			i1, j1, k1, i2, j2, k2 = 0, 0, 1, 1, 0, 1
		}
	} else { // x0 < y0
		if y0 < z0 {
			i1, j1, k1, i2, j2, k2 = 0, 0, 1, 0, 1, 1
		} else if x0 < z0 {
			i1, j1, k1, i2, j2, k2 = 0, 1, 0, 0, 1, 1
		} else {
			i1, j1, k1, i2, j2, k2 = 0, 1, 0, 1, 1, 0
		}
	}

	x1 := x0 - float64(i1) + G3
	y1 := y0 - float64(j1) + G3
	z1 := z0 - float64(k1) + G3
	x2 := x0 - float64(i2) + 2*G3
	y2 := y0 - float64(j2) + 2*G3
	z2 := z0 - float64(k2) + 2*G3
	x3 := x0 - 1 + 3*G3
	y3 := y0 - 1 + 3*G3
	z3 := z0 - 1 + 3*G3

	// Wrap the integer indices at 256, to avoid indexing perm[] out of bounds
	i = (i + 256) % 256
	j = (j + 256) % 256
	k = (k + 256) % 256

	// Calculate the contribution from the four corners
	t0 := 0.6 - x0*x0 - y0*y0 - z0*z0
	if t0 < 0 {
		n0 = 0
	} else {
		t0 *= t0
		// Calculate the contribution using the grad3 function and permuted indices
		n0 = t0 * t0 * grad3(perm[(i+int(perm[(j+int(perm[k]))%256]))%256], x0, y0, z0)
	}

	t1 := 0.6 - x1*x1 - y1*y1 - z1*z1
	if t1 < 0 {
		n1 = 0
	} else {
		t1 *= t1
		// Calculate the contribution using the grad3 function and permuted indices
		n1 = t1 * t1 * grad3(perm[(i+i1+int(perm[(j+j1+int(perm[(k+k1)%256]))%256]))%256], x1, y1, z1)
	}

	t2 := 0.6 - x2*x2 - y2*y2 - z2*z2
	if t2 < 0 {
		n2 = 0
	} else {
		t2 *= t2
		// Calculate the contribution using the grad3 function and permuted indices
		n2 = t2 * t2 * grad3(perm[(i+i2+int(perm[(j+j2+int(perm[(k+k2)%256]))%256]))%256], x2, y2, z2)
	}

	t3 := 0.6 - x3*x3 - y3*y3 - z3*z3
	if t3 < 0 {
		n3 = 0
	} else {
		t3 *= t3
		// Calculate the contribution using the grad3 function and permuted indices
		n3 = t3 * t3 * grad3(perm[(i+1+int(perm[(j+1+int(perm[(k+1)%256]))%256]))%256], x3, y3, z3)
	}

	// Scale the result to be within -1.0 and 1.0
	return float32((n0 + n1 + n2 + n3) / 0.030555466710745972)
}
