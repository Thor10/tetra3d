package tetra3d

import (
	"math"
)

// // The goal of fastmath.go is to provide vector operations that don't clone the vector to use. This means the main usage is not to use the results
// // directly, but rather as intermediary steps (i.e. use fastVectorSub to compare distances, or fastMatrixMult to multiply a vector by that final matrix).
// // Be careful with it, me!

// var standinVector = NewVectorZero()
// var standinMatrix = NewEmptyMatrix4()

// func fastVectorSub(a, b Vector) Vector {

// 	standinVector[0] = a[0] - b[0]
// 	standinVector[1] = a[1] - b[1]
// 	standinVector[2] = a[2] - b[2]
// 	return standinVector

// }

// func fastVectorDistanceSquared(a, b Vector) float64 {
// 	subX := a[0] - b[0]
// 	subY := a[1] - b[1]
// 	subZ := a[2] - b[2]
// 	return subX*subX + subY*subY + subZ*subZ
// }

// func fastVectorMagnitudeSquared(vec Vector) float64 {
// 	return vec[0]*vec[0] + vec[1]*vec[1] + vec[2]*vec[2]
// }

// func fastMatrixMult(matrix, other Matrix4) Matrix4 {

// 	standinMatrix[0][0] = matrix[0][0]*other[0][0] + matrix[0][1]*other[1][0] + matrix[0][2]*other[2][0] + matrix[0][3]*other[3][0]
// 	standinMatrix[1][0] = matrix[1][0]*other[0][0] + matrix[1][1]*other[1][0] + matrix[1][2]*other[2][0] + matrix[1][3]*other[3][0]
// 	standinMatrix[2][0] = matrix[2][0]*other[0][0] + matrix[2][1]*other[1][0] + matrix[2][2]*other[2][0] + matrix[2][3]*other[3][0]
// 	standinMatrix[3][0] = matrix[3][0]*other[0][0] + matrix[3][1]*other[1][0] + matrix[3][2]*other[2][0] + matrix[3][3]*other[3][0]

// 	standinMatrix[0][1] = matrix[0][0]*other[0][1] + matrix[0][1]*other[1][1] + matrix[0][2]*other[2][1] + matrix[0][3]*other[3][1]
// 	standinMatrix[1][1] = matrix[1][0]*other[0][1] + matrix[1][1]*other[1][1] + matrix[1][2]*other[2][1] + matrix[1][3]*other[3][1]
// 	standinMatrix[2][1] = matrix[2][0]*other[0][1] + matrix[2][1]*other[1][1] + matrix[2][2]*other[2][1] + matrix[2][3]*other[3][1]
// 	standinMatrix[3][1] = matrix[3][0]*other[0][1] + matrix[3][1]*other[1][1] + matrix[3][2]*other[2][1] + matrix[3][3]*other[3][1]

// 	standinMatrix[0][2] = matrix[0][0]*other[0][2] + matrix[0][1]*other[1][2] + matrix[0][2]*other[2][2] + matrix[0][3]*other[3][2]
// 	standinMatrix[1][2] = matrix[1][0]*other[0][2] + matrix[1][1]*other[1][2] + matrix[1][2]*other[2][2] + matrix[1][3]*other[3][2]
// 	standinMatrix[2][2] = matrix[2][0]*other[0][2] + matrix[2][1]*other[1][2] + matrix[2][2]*other[2][2] + matrix[2][3]*other[3][2]
// 	standinMatrix[3][2] = matrix[3][0]*other[0][2] + matrix[3][1]*other[1][2] + matrix[3][2]*other[2][2] + matrix[3][3]*other[3][2]

// 	standinMatrix[0][3] = matrix[0][0]*other[0][3] + matrix[0][1]*other[1][3] + matrix[0][2]*other[2][3] + matrix[0][3]*other[3][3]
// 	standinMatrix[1][3] = matrix[1][0]*other[0][3] + matrix[1][1]*other[1][3] + matrix[1][2]*other[2][3] + matrix[1][3]*other[3][3]
// 	standinMatrix[2][3] = matrix[2][0]*other[0][3] + matrix[2][1]*other[1][3] + matrix[2][2]*other[2][3] + matrix[2][3]*other[3][3]
// 	standinMatrix[3][3] = matrix[3][0]*other[0][3] + matrix[3][1]*other[1][3] + matrix[3][2]*other[2][3] + matrix[3][3]*other[3][3]

// 	return standinMatrix

// }

// func fastMatrixMultVec(matrix Matrix4, vect Vector) (x, y, z float64) {

// 	x = matrix[0][0]*vect[0] + matrix[1][0]*vect[1] + matrix[2][0]*vect[2] + matrix[3][0]
// 	y = matrix[0][1]*vect[0] + matrix[1][1]*vect[1] + matrix[2][1]*vect[2] + matrix[3][1]
// 	z = matrix[0][2]*vect[0] + matrix[1][2]*vect[1] + matrix[2][2]*vect[2] + matrix[3][2]

// 	return

// }

// func fastMatrixMultVecW(matrix Matrix4, vect Vector) (x, y, z, w float64) {

// 	x = matrix[0][0]*vect[0] + matrix[1][0]*vect[1] + matrix[2][0]*vect[2] + matrix[3][0]
// 	y = matrix[0][1]*vect[0] + matrix[1][1]*vect[1] + matrix[2][1]*vect[2] + matrix[3][1]
// 	z = matrix[0][2]*vect[0] + matrix[1][2]*vect[1] + matrix[2][2]*vect[2] + matrix[3][2]
// 	w = matrix[0][3]*vect[0] + matrix[1][3]*vect[1] + matrix[2][3]*vect[2] + matrix[3][3]

// 	return

// }

// func vectorCross(vecA, vecB, failsafeVec Vector) Vector {
// 	cross, _ := vecA.Cross(vecB)

// 	if cross.Magnitude() < 0.0001 {
// 		cross, _ = vecA.Cross(failsafeVec)

// 		// If it's still < 0, then it's not a separating axis
// 		if cross.Magnitude() < 0.0001 {
// 			return nil
// 		}
// 	}

// 	return cross
// }

// func vectorCrossUnsafe(v0, v1, out Vector) Vector {
// 	out[0] = v0[1]*v1[2] - v1[1]*v0[2]
// 	out[1] = v0[2]*v1[0] - v1[2]*v0[0]
// 	out[2] = v0[0]*v1[1] - v1[0]*v0[1]
// 	return out
// }

// // vectorsEqual returns if two >=3D vectors are basically equal in position.
// func vectorsEqual(a, b Vector) bool {
// 	m := 0.0001
// 	return math.Abs(a[0]-b[0]) < m && math.Abs(a[1]-b[1]) < m && math.Abs(a[2]-b[2]) < m
// }

// type VectorPool struct {
// 	Vectors        []Vector
// 	RetrievalIndex int
// 	Vectors4D      bool
// }

// func NewVectorPool(vectorCount int, vectors4D bool) *VectorPool {
// 	pool := &VectorPool{
// 		Vectors:   make([]Vector, vectorCount),
// 		Vectors4D: vectors4D,
// 	}
// 	if vectors4D {
// 		for i := 0; i < vectorCount; i++ {
// 			pool.Vectors[i] = Vector{0, 0, 0, 0}
// 		}
// 	} else {
// 		for i := 0; i < vectorCount; i++ {
// 			pool.Vectors[i] = Vector{0, 0, 0}
// 		}
// 	}
// 	return pool
// }

// func (pool *VectorPool) Reset() {
// 	pool.RetrievalIndex = 0
// }

// func (pool *VectorPool) Get() Vector {
// 	v := pool.Vectors[pool.RetrievalIndex]
// 	pool.RetrievalIndex++
// 	return v
// }

// func (pool *VectorPool) MultVec(matrix Matrix4, vect Vector) Vector {

// 	v := pool.Get()

// 	v[0] = matrix[0][0]*vect[0] + matrix[1][0]*vect[1] + matrix[2][0]*vect[2] + matrix[3][0]
// 	v[1] = matrix[0][1]*vect[0] + matrix[1][1]*vect[1] + matrix[2][1]*vect[2] + matrix[3][1]
// 	v[2] = matrix[0][2]*vect[0] + matrix[1][2]*vect[1] + matrix[2][2]*vect[2] + matrix[3][2]

// 	return v[:3]

// }

// func (pool *VectorPool) MultVecW(matrix Matrix4, vect Vector) Vector {

// 	v := pool.Get()

// 	v[0] = matrix[0][0]*vect[0] + matrix[1][0]*vect[1] + matrix[2][0]*vect[2] + matrix[3][0]
// 	v[1] = matrix[0][1]*vect[0] + matrix[1][1]*vect[1] + matrix[2][1]*vect[2] + matrix[3][1]
// 	v[2] = matrix[0][2]*vect[0] + matrix[1][2]*vect[1] + matrix[2][2]*vect[2] + matrix[3][2]
// 	v[3] = matrix[0][3]*vect[0] + matrix[1][3]*vect[1] + matrix[2][3]*vect[2] + matrix[3][3]

// 	return v

// }

// func (pool *VectorPool) Sub(v0, v1 Vector) Vector {
// 	out := pool.Get()

// 	if pool.Vectors4D {
// 		out[0] = v0[0] - v1[0]
// 		out[1] = v0[1] - v1[1]
// 		out[2] = v0[2] - v1[2]
// 		out[3] = v0[3] - v1[3]
// 	} else {
// 		out[0] = v0[0] - v1[0]
// 		out[1] = v0[1] - v1[1]
// 		out[2] = v0[2] - v1[2]
// 	}
// 	return out
// }

// func (pool *VectorPool) Add(v0, v1 Vector) Vector {
// 	out := pool.Get()
// 	if pool.Vectors4D {
// 		out[0] = v0[0] + v1[0]
// 		out[1] = v0[1] + v1[1]
// 		out[2] = v0[2] + v1[2]
// 		out[3] = v0[3] + v1[3]
// 	} else {
// 		out[0] = v0[0] + v1[0]
// 		out[1] = v0[1] + v1[1]
// 		out[2] = v0[2] + v1[2]
// 	}
// 	return out
// }

// func (pool *VectorPool) Cross(v0, v1 Vector) Vector {
// 	out := pool.Get()

// 	out[0] = v0[1]*v1[2] - v1[1]*v0[2]
// 	out[1] = v0[2]*v1[0] - v1[2]*v0[0]
// 	out[2] = v0[0]*v1[1] - v1[0]*v0[1]

// 	return out[:3]
// }

// Fast dot that should never call append() on the input Vectors, regardless of dimensions
// func dot(a, b Vector) float64 {
// 	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2]
// }

// ToRadians is a helper function to easily convert degrees to radians (which is what the rotation-oriented functions in Tetra3D use).
func ToRadians(degrees float64) float64 {
	return math.Pi * degrees / 180
}

// ToDegrees is a helper function to easily convert radians to degrees for human readability.
func ToDegrees(radians float64) float64 {
	return radians / math.Pi * 180
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}

func pow(value float64, power int) float64 {
	x := value
	for i := 0; i < power; i++ {
		x *= x
	}
	return x
}

func round(value float64) float64 {

	iv := float64(int(value))

	if value > iv+0.5 {
		return iv + 1
	} else if value < iv-0.5 {
		return iv - 1
	}

	return iv

}
