package primitives

import (
	"math"
)

type Matrix4x4 struct {
	Points [16]float32
}

func NullMatrix4x4() *Matrix4x4 {
	return &Matrix4x4{
		[16]float32{
			0.0, 0.0, 0.0, 0.0,
			0.0, 0.0, 0.0, 0.0,
			0.0, 0.0, 0.0, 0.0,
			0.0, 0.0, 0.0, 0.0,
		},
	}
}
func UnitMatrix4x4() *Matrix4x4 {
	return &Matrix4x4{
		[16]float32{
			1.0, 0.0, 0.0, 0.0,
			0.0, 1.0, 0.0, 0.0,
			0.0, 0.0, 1.0, 0.0,
			0.0, 0.0, 0.0, 1.0,
		},
	}
}
func ScaleMatrix4x4(scaleX, scaleY, scaleZ float32) *Matrix4x4 {
	return &Matrix4x4{
		[16]float32{
			scaleX, 0.0, 0.0, 0.0,
			0.0, scaleY, 0.0, 0.0,
			0.0, 0.0, scaleZ, 0.0,
			0.0, 0.0, 0.0, 1.0,
		},
	}
}
func TranslationMatrix4x4(translationX, translationY, translationZ float32) *Matrix4x4 {
	return &Matrix4x4{
		[16]float32{
			1.0, 0.0, 0.0, translationX,
			0.0, 1.0, 0.0, translationY,
			0.0, 0.0, 1.0, translationZ,
			0.0, 0.0, 0.0, 1.0,
		},
	}
}

//func ProjectionMatrix4x4(const float &angleOfView, const float &near, const float &far, Matrix44f &M) *Matrix4x4 {
func ProjectionMatrix4x4(angleOfView, near, far float64) *Matrix4x4 {
	scale := float32(1 / math.Tan(angleOfView*0.5*math.Pi/180))
	projection := NullMatrix4x4()
	projection.Points[0] = scale
	projection.Points[5] = scale
	projection.Points[10] = float32(-far / (far - near))
	projection.Points[14] = float32(-far * near / (far - near))
	projection.Points[11] = -1
	projection.Points[15] = 0
	return projection
}

func (m *Matrix4x4) Dot(m2 *Matrix4x4) *Matrix4x4 {
	result := NullMatrix4x4()
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			result.Points[4*i+j] = m.Points[4*i+0]*m2.Points[4*0+j] +
				m.Points[4*i+1]*m2.Points[4*1+j] +
				m.Points[4*i+2]*m2.Points[4*2+j] +
				m.Points[4*i+3]*m2.Points[4*3+j]
		}
	}
	return result
}
