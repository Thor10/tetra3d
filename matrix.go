package tetra3d

import (
	"math"
	"strconv"
)

// Matrix4 represents a 4x4 matrix for translation, scale, and rotation. A Matrix4 in Tetra3D is row-major (i.e. the X axis is matrix[0]).
type Matrix4 [4][4]float64

// NewMatrix4 returns a new identity Matrix4. A Matrix4 in Tetra3D is row-major (i.e. the X axis for a rotation Matrix4 is matrix[0][0], matrix[0][1], matrix[0][2]).
func NewMatrix4() Matrix4 {

	mat := Matrix4{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
	return mat

}

func NewEmptyMatrix4() Matrix4 {

	mat := Matrix4{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}
	return mat

}

func (matrix *Matrix4) Clear() {
	for y := 0; y < len(matrix); y++ {
		for x := 0; x < len(matrix[y]); x++ {
			matrix[y][x] = 0
		}
	}
}

// Clone clones the Matrix4, returning a new copy.
func (matrix Matrix4) Clone() Matrix4 {
	newMat := NewMatrix4()
	for y := 0; y < len(matrix); y++ {
		for x := 0; x < len(matrix[y]); x++ {
			newMat[y][x] = matrix[y][x]
		}
	}
	return newMat
}

// Set allows you to set the Matrix4 to the same values as another Matrix4.
func (matrix *Matrix4) Set(other Matrix4) {
	for y := 0; y < len(matrix); y++ {
		for x := 0; x < len(matrix[y]); x++ {
			matrix[y][x] = other[y][x]
		}
	}
}

// BlenderToTetra returns a Matrix with the rows altered such that Blender's +Z is now Tetra's +Y and Blender's +Y is now Tetra's -Z.
func (matrix Matrix4) BlenderToTetra() Matrix4 {
	newMat := matrix.Clone()
	newMat.SetRow(1, matrix.Row(2).Invert())
	newMat.SetRow(2, matrix.Row(1))
	return newMat
}

// NewMatrix4Scale returns a new identity Matrix4, but with the x, y, and z translation components set as provided.
func NewMatrix4Translate(x, y, z float64) Matrix4 {
	mat := NewMatrix4()
	mat[3][0] = x
	mat[3][1] = y
	mat[3][2] = z
	return mat
}

// NewMatrix4Scale returns a new identity Matrix4, but with the scale components set as provided. 1, 1, 1 is the default.
func NewMatrix4Scale(x, y, z float64) Matrix4 {
	mat := NewMatrix4()
	mat[0][0] = x
	mat[1][1] = y
	mat[2][2] = z
	return mat
}

// NewMatrix4Rotate returns a new Matrix4 designed to rotate by the angle given (in radians) along the axis given [x, y, z].
// This rotation works as though you pierced the object utilizing the matrix through by the axis, and then rotated it
// counter-clockwise by the angle in radians.
func NewMatrix4Rotate(x, y, z, angle float64) Matrix4 {

	// Default to spinning on +Y axis if there is no valid axis
	if x == 0 && y == 0 && z == 0 {
		y = 1
	}

	mat := NewMatrix4()
	vector := Vector{x, y, z, 1}.Unit()
	s := math.Sin(angle)
	c := math.Cos(angle)
	m := 1 - c

	mat[0][0] = m*vector.X*vector.X + c
	mat[0][1] = m*vector.X*vector.Y + vector.Z*s
	mat[0][2] = m*vector.Z*vector.X - vector.Y*s

	mat[1][0] = m*vector.X*vector.Y - vector.Z*s
	mat[1][1] = m*vector.Y*vector.Y + c
	mat[1][2] = m*vector.Y*vector.Z + vector.X*s

	mat[2][0] = m*vector.Z*vector.X + vector.Y*s
	mat[2][1] = m*vector.Y*vector.Z - vector.X*s
	mat[2][2] = m*vector.Z*vector.Z + c

	return mat

}

// ToQuaternion returns a Quaternion representative of the Matrix4's rotation (assuming it is just a purely rotational Matrix4).
func (matrix Matrix4) ToQuaternion() Quaternion {

	sqrt := math.Sqrt(1 + matrix[0][0] + matrix[1][1] + matrix[2][2])

	if sqrt != 0 {
		qw := sqrt / 2

		return NewQuaternion(
			(matrix[1][2]-matrix[2][1])/(4*qw),
			(matrix[2][0]-matrix[0][2])/(4*qw),
			(matrix[0][1]-matrix[1][0])/(4*qw),
			qw,
		)

	}

	return NewQuaternion(0, 1, 0, 0)

}

// Right returns the right-facing rotational component of the Matrix4. For an identity matrix, this would be [1, 0, 0], or +X.
func (matrix Matrix4) Right() Vector {
	return Vector{
		matrix[0][0],
		matrix[0][1],
		matrix[0][2],
		1,
	}.Unit()
}

// Up returns the upward rotational component of the Matrix4. For an identity matrix, this would be [0, 1, 0], or +Y.
func (matrix Matrix4) Up() Vector {
	return Vector{
		matrix[1][0],
		matrix[1][1],
		matrix[1][2],
		1,
	}.Unit()
}

// Forward returns the forward rotational component of the Matrix4. For an identity matrix, this would be [0, 0, 1], or +Z (towards camera).
func (matrix Matrix4) Forward() Vector {
	return Vector{
		matrix[2][0],
		matrix[2][1],
		matrix[2][2],
		1,
	}.Unit()
}

// Decompose decomposes the Matrix4 and returns three components - the position (a 3D Vector), scale (another 3D Vector), and rotation (an AxisAngle)
// indicated by the Matrix4. Note that this is mainly used when loading a mesh from a 3D modeler - this being the case, it may not be the most precise, and negative
// scales are not supported.
func (matrix Matrix4) Decompose() (Vector, Vector, Matrix4) {

	position := Vector{matrix[3][0], matrix[3][1], matrix[3][2], 1}

	rotation := NewMatrix4()
	rotation.SetRow(0, matrix.Row(0).Unit())
	rotation.SetRow(1, matrix.Row(1).Unit())
	rotation.SetRow(2, matrix.Row(2).Unit())

	in := matrix.Mult(rotation.Transposed())

	scale := Vector{in.Row(0).Magnitude(), in.Row(1).Magnitude(), in.Row(2).Magnitude(), 1}

	return position, scale, rotation

}

// Transposed transposes a Matrix4, switching the Matrix from being Row Major to being Column Major. For orthonormalized Matrices (matrices
// that have rows that are normalized (having a length of 1), like rotation matrices), this is equivalent to inverting it.
func (matrix Matrix4) Transposed() Matrix4 {

	new := NewMatrix4()

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			new[i][j] = matrix[j][i]
		}
	}

	return new

}

// Inverted returns an inverted (reversed) clone of a Matrix4. An inverted matrix is defined here as a matrix
// composed of a decomposed matrix's transposed rotation, negated position, and the same scale (since scale is multiplicative).
// func (matrix Matrix4) InvertedOld() Matrix4 {

// 	p, s, r := matrix.Decompose()

// 	newMat := NewMatrix4()
// 	newMat = newMat.SetRow(0, r.Row(0))
// 	newMat = newMat.SetRow(1, r.Row(1))
// 	newMat = newMat.SetRow(2, r.Row(2))
// 	newMat = newMat.Transposed()

// 	newMat[0][0] *= 1 / s[0]
// 	newMat[1][1] *= 1 / s[1]
// 	newMat[2][2] *= 1 / s[2]

// 	newMat = newMat.SetRow(3, Vector{-p[0], -p[1], -p[2], 1})

// 	return newMat

// }

// The ultimate sin; I'm just going to copy this code for inverting a 4x4 Matrix and call it a day.
// This code was obtained from https://stackoverflow.com/questions/1148309/inverting-a-4x4-matrix,
// and is like 200x faster than my old inversion code, whaaaaa
func (matrix Matrix4) Inverted() Matrix4 {

	var A2323 = matrix[2][2]*matrix[3][3] - matrix[2][3]*matrix[3][2]
	var A1323 = matrix[2][1]*matrix[3][3] - matrix[2][3]*matrix[3][1]
	var A1223 = matrix[2][1]*matrix[3][2] - matrix[2][2]*matrix[3][1]
	var A0323 = matrix[2][0]*matrix[3][3] - matrix[2][3]*matrix[3][0]
	var A0223 = matrix[2][0]*matrix[3][2] - matrix[2][2]*matrix[3][0]
	var A0123 = matrix[2][0]*matrix[3][1] - matrix[2][1]*matrix[3][0]
	var A2313 = matrix[1][2]*matrix[3][3] - matrix[1][3]*matrix[3][2]
	var A1313 = matrix[1][1]*matrix[3][3] - matrix[1][3]*matrix[3][1]
	var A1213 = matrix[1][1]*matrix[3][2] - matrix[1][2]*matrix[3][1]
	var A2312 = matrix[1][2]*matrix[2][3] - matrix[1][3]*matrix[2][2]
	var A1312 = matrix[1][1]*matrix[2][3] - matrix[1][3]*matrix[2][1]
	var A1212 = matrix[1][1]*matrix[2][2] - matrix[1][2]*matrix[2][1]
	var A0313 = matrix[1][0]*matrix[3][3] - matrix[1][3]*matrix[3][0]
	var A0213 = matrix[1][0]*matrix[3][2] - matrix[1][2]*matrix[3][0]
	var A0312 = matrix[1][0]*matrix[2][3] - matrix[1][3]*matrix[2][0]
	var A0212 = matrix[1][0]*matrix[2][2] - matrix[1][2]*matrix[2][0]
	var A0113 = matrix[1][0]*matrix[3][1] - matrix[1][1]*matrix[3][0]
	var A0112 = matrix[1][0]*matrix[2][1] - matrix[1][1]*matrix[2][0]

	var det = matrix[0][0]*(matrix[1][1]*A2323-matrix[1][2]*A1323+matrix[1][3]*A1223) -
		matrix[0][1]*(matrix[1][0]*A2323-matrix[1][2]*A0323+matrix[1][3]*A0223) +
		matrix[0][2]*(matrix[1][0]*A1323-matrix[1][1]*A0323+matrix[1][3]*A0123) -
		matrix[0][3]*(matrix[1][0]*A1223-matrix[1][1]*A0223+matrix[1][2]*A0123)

	det = 1 / det

	m := NewMatrix4()

	m[0][0] = det * (matrix[1][1]*A2323 - matrix[1][2]*A1323 + matrix[1][3]*A1223)
	m[0][1] = det * -(matrix[0][1]*A2323 - matrix[0][2]*A1323 + matrix[0][3]*A1223)
	m[0][2] = det * (matrix[0][1]*A2313 - matrix[0][2]*A1313 + matrix[0][3]*A1213)
	m[0][3] = det * -(matrix[0][1]*A2312 - matrix[0][2]*A1312 + matrix[0][3]*A1212)
	m[1][0] = det * -(matrix[1][0]*A2323 - matrix[1][2]*A0323 + matrix[1][3]*A0223)
	m[1][1] = det * (matrix[0][0]*A2323 - matrix[0][2]*A0323 + matrix[0][3]*A0223)
	m[1][2] = det * -(matrix[0][0]*A2313 - matrix[0][2]*A0313 + matrix[0][3]*A0213)
	m[1][3] = det * (matrix[0][0]*A2312 - matrix[0][2]*A0312 + matrix[0][3]*A0212)
	m[2][0] = det * (matrix[1][0]*A1323 - matrix[1][1]*A0323 + matrix[1][3]*A0123)
	m[2][1] = det * -(matrix[0][0]*A1323 - matrix[0][1]*A0323 + matrix[0][3]*A0123)
	m[2][2] = det * (matrix[0][0]*A1313 - matrix[0][1]*A0313 + matrix[0][3]*A0113)
	m[2][3] = det * -(matrix[0][0]*A1312 - matrix[0][1]*A0312 + matrix[0][3]*A0112)
	m[3][0] = det * -(matrix[1][0]*A1223 - matrix[1][1]*A0223 + matrix[1][2]*A0123)
	m[3][1] = det * (matrix[0][0]*A1223 - matrix[0][1]*A0223 + matrix[0][2]*A0123)
	m[3][2] = det * -(matrix[0][0]*A1213 - matrix[0][1]*A0213 + matrix[0][2]*A0113)
	m[3][3] = det * (matrix[0][0]*A1212 - matrix[0][1]*A0212 + matrix[0][2]*A0112)

	return m

}

// Inverted returns an inverted (reversed) clone of a Matrix4. See the above StackOverflow link.
// func (matrix Matrix4) oldInverted() Matrix4 {

// 	inv := NewMatrix4()

// 	m := matrix.Index

// 	inv.setIndex(0, m(5)*m(10)*m(15)-
// 		m(5)*m(11)*m(14)-
// 		m(9)*m(6)*m(15)+
// 		m(9)*m(7)*m(14)+
// 		m(13)*m(6)*m(11)-
// 		m(13)*m(7)*m(10))

// 	inv.setIndex(4, -m(4)*m(10)*m(15)+
// 		m(4)*m(11)*m(14)+
// 		m(8)*m(6)*m(15)-
// 		m(8)*m(7)*m(14)-
// 		m(12)*m(6)*m(11)+
// 		m(12)*m(7)*m(10))

// 	inv.setIndex(8, m(4)*m(9)*m(15)-
// 		m(4)*m(11)*m(13)-
// 		m(8)*m(5)*m(15)+
// 		m(8)*m(7)*m(13)+
// 		m(12)*m(5)*m(11)-
// 		m(12)*m(7)*m(9))

// 	inv.setIndex(12, -m(4)*m(9)*m(14)+
// 		m(4)*m(10)*m(13)+
// 		m(8)*m(5)*m(14)-
// 		m(8)*m(6)*m(13)-
// 		m(12)*m(5)*m(10)+
// 		m(12)*m(6)*m(9))

// 	inv.setIndex(1, -m(1)*m(10)*m(15)+
// 		m(1)*m(11)*m(14)+
// 		m(9)*m(2)*m(15)-
// 		m(9)*m(3)*m(14)-
// 		m(13)*m(2)*m(11)+
// 		m(13)*m(3)*m(10))

// 	inv.setIndex(5, m(0)*m(10)*m(15)-
// 		m(0)*m(11)*m(14)-
// 		m(8)*m(2)*m(15)+
// 		m(8)*m(3)*m(14)+
// 		m(12)*m(2)*m(11)-
// 		m(12)*m(3)*m(10))

// 	inv.setIndex(9, -m(0)*m(9)*m(15)+
// 		m(0)*m(11)*m(13)+
// 		m(8)*m(1)*m(15)-
// 		m(8)*m(3)*m(13)-
// 		m(12)*m(1)*m(11)+
// 		m(12)*m(3)*m(9))

// 	inv.setIndex(13, m(0)*m(9)*m(14)-
// 		m(0)*m(10)*m(13)-
// 		m(8)*m(1)*m(14)+
// 		m(8)*m(2)*m(13)+
// 		m(12)*m(1)*m(10)-
// 		m(12)*m(2)*m(9))

// 	inv.setIndex(2, m(1)*m(6)*m(15)-
// 		m(1)*m(7)*m(14)-
// 		m(5)*m(2)*m(15)+
// 		m(5)*m(3)*m(14)+
// 		m(13)*m(2)*m(7)-
// 		m(13)*m(3)*m(6))

// 	inv.setIndex(6, -m(0)*m(6)*m(15)+
// 		m(0)*m(7)*m(14)+
// 		m(4)*m(2)*m(15)-
// 		m(4)*m(3)*m(14)-
// 		m(12)*m(2)*m(7)+
// 		m(12)*m(3)*m(6))

// 	inv.setIndex(10, m(0)*m(5)*m(15)-
// 		m(0)*m(7)*m(13)-
// 		m(4)*m(1)*m(15)+
// 		m(4)*m(3)*m(13)+
// 		m(12)*m(1)*m(7)-
// 		m(12)*m(3)*m(5))

// 	inv.setIndex(14, -m(0)*m(5)*m(14)+
// 		m(0)*m(6)*m(13)+
// 		m(4)*m(1)*m(14)-
// 		m(4)*m(2)*m(13)-
// 		m(12)*m(1)*m(6)+
// 		m(12)*m(2)*m(5))

// 	inv.setIndex(3, -m(1)*m(6)*m(11)+
// 		m(1)*m(7)*m(10)+
// 		m(5)*m(2)*m(11)-
// 		m(5)*m(3)*m(10)-
// 		m(9)*m(2)*m(7)+
// 		m(9)*m(3)*m(6))

// 	inv.setIndex(7, m(0)*m(6)*m(11)-
// 		m(0)*m(7)*m(10)-
// 		m(4)*m(2)*m(11)+
// 		m(4)*m(3)*m(10)+
// 		m(8)*m(2)*m(7)-
// 		m(8)*m(3)*m(6))

// 	inv.setIndex(11, -m(0)*m(5)*m(11)+
// 		m(0)*m(7)*m(9)+
// 		m(4)*m(1)*m(11)-
// 		m(4)*m(3)*m(9)-
// 		m(8)*m(1)*m(7)+
// 		m(8)*m(3)*m(5))

// 	inv.setIndex(15, m(0)*m(5)*m(10)-
// 		m(0)*m(6)*m(9)-
// 		m(4)*m(1)*m(10)+
// 		m(4)*m(2)*m(9)+
// 		m(8)*m(1)*m(6)-
// 		m(8)*m(2)*m(5))

// 	det := m(0)*inv.Index(0) + m(1)*inv.Index(4) + m(2)*inv.Index(8) + m(3)*inv.Index(12)

// 	if det == 0 {
// 		return NewMatrix4()
// 	}

// 	det = 1.0 / det

// 	for i := 0; i < 16; i++ {
// 		inv.setIndex(i, inv.Index(i)*det)
// 	}

// 	return inv

// }

func (matrix *Matrix4) setIndex(index int, value float64) {
	matrix[index/4][index%4] = value
}

func (matrix *Matrix4) Index(index int) float64 {
	y := index / 4
	x := index % 4
	return matrix[y][x]
}

// Equals returns true if the matrix equals the same values in the provided Other Matrix4.
func (matrix Matrix4) Equals(other Matrix4) bool {

	eps := 0.0001 // epsilon floating point error value
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[i]); j++ {
			if math.Abs(matrix[i][j]-other[i][j]) > eps {
				return false
			}
		}
	}
	return true
}

var identityMatrix = NewMatrix4()

// IsIdentity returns true if the matrix is an unmodified identity matrix.
func (matrix Matrix4) IsIdentity() bool {
	return matrix.Equals(identityMatrix)
}

// Row returns the indiced row from the Matrix4 as a Vector.
func (matrix Matrix4) Row(rowIndex int) Vector {
	vec := Vector{
		matrix[rowIndex][0],
		matrix[rowIndex][1],
		matrix[rowIndex][2],
		matrix[rowIndex][3],
	}
	return vec
}

// Column returns the indiced column from the Matrix4 as a Vector.
func (matrix Matrix4) Column(columnIndex int) Vector {
	vec := Vector{
		matrix[0][columnIndex],
		matrix[1][columnIndex],
		matrix[2][columnIndex],
		matrix[3][columnIndex],
	}
	return vec
}

// SetRow returns a clone of the Matrix4 with the row in rowIndex set to the 4D vector passed.
func (matrix *Matrix4) SetRow(rowIndex int, vec Vector) {
	matrix[rowIndex][0] = vec.X
	matrix[rowIndex][1] = vec.Y
	matrix[rowIndex][2] = vec.Z
	matrix[rowIndex][3] = vec.W
}

// SetColumn returns a clone of the Matrix4 with the column in columnIndex set to the 4D vector passed.
func (matrix *Matrix4) SetColumn(columnIndex int, vec Vector) {
	matrix[0][columnIndex] = vec.X
	matrix[1][columnIndex] = vec.Y
	matrix[2][columnIndex] = vec.Z
	matrix[3][columnIndex] = vec.W
}

// Rotated returns a clone of the Matrix4 rotated along the local axis by the angle given (in radians). This rotation works as though
// you pierced the object through by the axis, and then rotated it counter-clockwise by the angle
// in radians. The axis is relative to any existing rotation contained in the matrix.
func (matrix Matrix4) Rotated(x, y, z, angle float64) Matrix4 {
	return NewMatrix4Rotate(x, y, z, angle).Mult(matrix)
}

// NewProjectionPerspective generates a perspective frustum Matrix4. fovy is the vertical field of view in degrees, near and far are the near and far clipping plane,
// while viewWidth and viewHeight is the width and height of the backing texture / camera. Generally, you won't need to use this directly.
func NewProjectionPerspective(fovy, near, far, viewWidth, viewHeight float64) Matrix4 {

	aspect := viewWidth / viewHeight

	t := math.Tan(fovy * math.Pi / 360)
	b := -t
	r := t * aspect
	l := -r

	return Matrix4{
		{(2 * near) / (r - l), 0, (r + l) / (r - l), 0},
		{0, (2 * near) / (t - b), (t + b) / (t - b), 0},
		{0, 0, -((far + near) / (far - near)), -((2 * far * near) / (far - near))},
		{0, 0, -1, 0},
	}

}

// NewProjectionOrthographic generates an orthographic frustum Matrix4. near and far are the near and far clipping plane. right, left, top, and bottom
// are the right, left, top, and bottom planes (usually 1 and -1 for right and left, and the aspect ratio of the window and negative for top and bottom).
// Generally, you won't need to use this directly.
func NewProjectionOrthographic(near, far, right, left, top, bottom float64) Matrix4 {
	return Matrix4{
		{2 / (right - left), 0, 0, 0},
		{0, 2 / (top - bottom), 0, 0},
		{0, 0, -2 / (far - near), 0},
		{0, 0, 0, 1},
	}
}

// MultVec multiplies the vector provided by the Matrix4, giving a vector that has been rotated, scaled, or translated as desired.
func (matrix Matrix4) MultVec(vect Vector) Vector {

	return Vector{
		matrix[0][0]*vect.X + matrix[1][0]*vect.Y + matrix[2][0]*vect.Z + matrix[3][0],
		matrix[0][1]*vect.X + matrix[1][1]*vect.Y + matrix[2][1]*vect.Z + matrix[3][1],
		matrix[0][2]*vect.X + matrix[1][2]*vect.Y + matrix[2][2]*vect.Z + matrix[3][2],
		1,
	}

}

// MultVecW multiplies the vector provided by the Matrix4, including the fourth (W) component, giving a vector that has been rotated, scaled, or translated as desired.
func (matrix Matrix4) MultVecW(vect Vector) Vector {

	return Vector{
		matrix[0][0]*vect.X + matrix[1][0]*vect.Y + matrix[2][0]*vect.Z + matrix[3][0],
		matrix[0][1]*vect.X + matrix[1][1]*vect.Y + matrix[2][1]*vect.Z + matrix[3][1],
		matrix[0][2]*vect.X + matrix[1][2]*vect.Y + matrix[2][2]*vect.Z + matrix[3][2],
		matrix[0][3]*vect.X + matrix[1][3]*vect.Y + matrix[2][3]*vect.Z + matrix[3][3],
	}

}

// Mult multiplies a Matrix4 by another provided Matrix4 - this effectively combines them.
func (matrix Matrix4) Mult(other Matrix4) Matrix4 {

	newMat := NewMatrix4()

	newMat[0][0] = matrix[0][0]*other[0][0] + matrix[0][1]*other[1][0] + matrix[0][2]*other[2][0] + matrix[0][3]*other[3][0]
	newMat[1][0] = matrix[1][0]*other[0][0] + matrix[1][1]*other[1][0] + matrix[1][2]*other[2][0] + matrix[1][3]*other[3][0]
	newMat[2][0] = matrix[2][0]*other[0][0] + matrix[2][1]*other[1][0] + matrix[2][2]*other[2][0] + matrix[2][3]*other[3][0]
	newMat[3][0] = matrix[3][0]*other[0][0] + matrix[3][1]*other[1][0] + matrix[3][2]*other[2][0] + matrix[3][3]*other[3][0]

	newMat[0][1] = matrix[0][0]*other[0][1] + matrix[0][1]*other[1][1] + matrix[0][2]*other[2][1] + matrix[0][3]*other[3][1]
	newMat[1][1] = matrix[1][0]*other[0][1] + matrix[1][1]*other[1][1] + matrix[1][2]*other[2][1] + matrix[1][3]*other[3][1]
	newMat[2][1] = matrix[2][0]*other[0][1] + matrix[2][1]*other[1][1] + matrix[2][2]*other[2][1] + matrix[2][3]*other[3][1]
	newMat[3][1] = matrix[3][0]*other[0][1] + matrix[3][1]*other[1][1] + matrix[3][2]*other[2][1] + matrix[3][3]*other[3][1]

	newMat[0][2] = matrix[0][0]*other[0][2] + matrix[0][1]*other[1][2] + matrix[0][2]*other[2][2] + matrix[0][3]*other[3][2]
	newMat[1][2] = matrix[1][0]*other[0][2] + matrix[1][1]*other[1][2] + matrix[1][2]*other[2][2] + matrix[1][3]*other[3][2]
	newMat[2][2] = matrix[2][0]*other[0][2] + matrix[2][1]*other[1][2] + matrix[2][2]*other[2][2] + matrix[2][3]*other[3][2]
	newMat[3][2] = matrix[3][0]*other[0][2] + matrix[3][1]*other[1][2] + matrix[3][2]*other[2][2] + matrix[3][3]*other[3][2]

	newMat[0][3] = matrix[0][0]*other[0][3] + matrix[0][1]*other[1][3] + matrix[0][2]*other[2][3] + matrix[0][3]*other[3][3]
	newMat[1][3] = matrix[1][0]*other[0][3] + matrix[1][1]*other[1][3] + matrix[1][2]*other[2][3] + matrix[1][3]*other[3][3]
	newMat[2][3] = matrix[2][0]*other[0][3] + matrix[2][1]*other[1][3] + matrix[2][2]*other[2][3] + matrix[2][3]*other[3][3]
	newMat[3][3] = matrix[3][0]*other[0][3] + matrix[3][1]*other[1][3] + matrix[3][2]*other[2][3] + matrix[3][3]*other[3][3]

	return newMat

}

func (matrix Matrix4) Add(other Matrix4) Matrix4 {

	newMat := matrix.Clone()

	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[i]); j++ {
			newMat[i][j] += other[i][j]
		}
	}

	return newMat

}

func (matrix Matrix4) ScaleByScalar(scalar float64) Matrix4 {

	newMat := matrix.Clone()

	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[i]); j++ {
			newMat[i][j] *= scalar
		}
	}

	return newMat

}

func (matrix Matrix4) Sub(other Matrix4) Matrix4 {

	newMat := matrix.Clone()

	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[i]); j++ {
			newMat[i][j] -= other[i][j]
		}
	}

	return newMat

}

// Lerping matrices is not very useful, so I'm hesitant to add it in
// func (matrix Matrix4) Lerp(other Matrix4, perc float64) Matrix4 {
// 	return matrix.Add(other.Sub(matrix).ScaleByScalar(perc))
// }

// Columns returns the Matrix4 as a slice of []float64, in column-major order (so it's transposed from the row-major default).
func (matrix Matrix4) Columns() [][]float64 {

	columns := [][]float64{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}

	for r := range matrix {
		for c := range matrix[r] {
			columns[c][r] = matrix[r][c]
		}
	}

	return columns
}

// IsZero returns true if the Matrix is zero'd out (all values are 0).
func (matrix Matrix4) IsZero() bool {
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[i]); j++ {
			if matrix[i][j] != 0 {
				return false
			}
		}
	}
	return true
}

// HasValidRotation returns if the first three vectors in the Matrix are non-zero
func (matrix Matrix4) HasValidRotation() bool {

	return matrix[0][0] != 0 || matrix[0][1] != 0 || matrix[0][2] != 0 &&
		matrix[1][0] != 0 || matrix[1][1] != 0 || matrix[1][2] != 0 &&
		matrix[2][0] != 0 || matrix[2][1] != 0 || matrix[2][2] != 0

}

func (matrix Matrix4) String() string {
	s := "{"
	for i, y := range matrix {
		for _, x := range y {
			s += strconv.FormatFloat(x, 'f', -1, 64) + ", "
		}
		if i < len(matrix)-1 {
			s += "\n"
		}
	}
	s += "}"
	return s
}

// ToEuler converts the rotation portion of the Matrix to a euler values,
// returned in the form of a Vector.
// func (matrix Matrix4) ToEuler() Vector {
// 	// This function doesn't work properly yet.
// 	eps := 0.998
// 	vec := NewVectorZero()

// 	if matrix[1][0] > eps {
// 		vec.Y = math.Atan2(matrix[0][2], matrix[2][2])
// 		vec.Z = math.Pi / 2
// 		vec.X = 0
// 	} else if matrix[1][0] < -eps {
// 		vec.Y = math.Atan2(matrix[0][2], matrix[2][2])
// 		vec.Z = -math.Pi / 2
// 		vec.X = 0
// 	} else {
// 		vec.Y = math.Atan2(-matrix[2][0], matrix[0][0])
// 		vec.X = math.Atan2(-matrix[1][2], matrix[1][1])
// 		vec.Z = math.Asin(matrix[1][0])
// 	}
// 	return vec
// }

// NewMatrix4RotateFromEuler creates a rotation Matrix4 from the euler values contained within the
// Vector.
func NewMatrix4RotateFromEuler(euler Vector) Matrix4 {
	mat := NewMatrix4()

	ch := math.Cos(euler.Y)
	sh := math.Sin(euler.Y)

	ca := math.Cos(euler.Z)
	sb := math.Sin(euler.Z)

	caa := math.Cos(euler.X)
	saa := math.Sin(euler.X)

	mat[0][0] = ch * ca
	mat[0][1] = sh*saa - ch*sb*caa
	mat[0][2] = ch*sb*saa + sh*caa

	mat[1][0] = sb
	mat[1][1] = ca * caa
	mat[1][2] = -ca * saa

	mat[2][0] = -sh * ca
	mat[2][1] = sh*sb*caa + ch*saa
	mat[2][2] = -sh*sb*saa + ch*caa

	return mat
}

// NewLookAtMatrix generates a new Matrix4 to rotate an object to point towards another object. to is the target's world position,
// from is the world position of the object looking towards the target, and up is the upward vector ( usually +Y, or [0, 1, 0] ).
func NewLookAtMatrix(from, to, up Vector) Matrix4 {

	// If from and to are the same, then an identity Matrix4 should be a sensible default
	if from.Equals(to) {
		return NewMatrix4()
	}
	z := to.Sub(from).Unit()

	up = up.Unit()

	// If z == up, then the matrix will be unusable, so we sub up out with another angle
	if z.Equals(up) || z.Equals(up.Invert()) {
		if !up.Equals(WorldRight) {
			up = WorldRight
		} else {
			up = WorldBack
		}
	}

	x := up.Cross(z).Unit()
	y := z.Cross(x)
	return Matrix4{
		{x.X, x.Y, x.Z, 0},
		{y.X, y.Y, y.Z, 0},
		{z.X, z.Y, z.Z, 0},
		{0, 0, 0, 1},
	}
}

// Lerp lerps a matrix to another, destination Matrix by the percent given. It does this by converting both
// Matrices to Quaternions, lerping them, then converting the result back to a Matrix4.
func (mat Matrix4) Lerp(other Matrix4, percent float64) Matrix4 {
	q1 := mat.ToQuaternion()
	q2 := other.ToQuaternion()
	return q1.Lerp(q2, percent).ToMatrix4()
}
