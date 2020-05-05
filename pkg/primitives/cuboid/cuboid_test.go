package cuboid

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/akosgarai/opengl_playground/pkg/primitives/material"
	"github.com/akosgarai/opengl_playground/pkg/primitives/rectangle"
)

var (
	DefaultCoordinates = [4]mgl32.Vec3{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 1, 0},
		mgl32.Vec3{0, 1, 0},
	}
	DefaultColors = [4]mgl32.Vec3{
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
	}
)

type testShader struct {
	HasTextureValue bool
}

func (t testShader) Use() {
}
func (t testShader) SetUniformMat4(s string, m mgl32.Mat4) {
}
func (t testShader) SetUniform1f(s string, f1 float32) {
}
func (t testShader) SetUniform3f(s string, f1, f2, f3 float32) {
}
func (t testShader) DrawTriangles(i int32) {
}
func (t testShader) Close(i int) {
}
func (t testShader) VertexAttribPointer(i uint32, c int32, s int32, o int) {
}
func (t testShader) BindVertexArray() {
}
func (t testShader) BindBufferData(d []float32) {
}
func (t testShader) HasTexture() bool {
	return t.HasTextureValue
}

var shader testShader

func TestNew(t *testing.T) {
	shader.HasTextureValue = false
	bottom := rectangle.New(DefaultCoordinates, DefaultColors, shader)
	cube := New(bottom, 1, shader)
	cubeBottom := cube.sides[0].Coordinates()
	cubeTop := cube.sides[1].Coordinates()
	var cubeTopExpected [4]mgl32.Vec3
	for i := 0; i < 4; i++ {
		cubeTopExpected[i] = cubeBottom[i].Add(mgl32.Vec3{0, 0, -1})
	}
	for i := 0; i < 4; i++ {
		if cubeBottom[i] != DefaultCoordinates[i] {
			t.Error("Mismatch in the bottom coordinates")
		}
		if cubeTop[i] != cubeTopExpected[i] {
			t.Error("Mismatch in the top coordinates")
		}
	}
}
func TestLog(t *testing.T) {
	shader.HasTextureValue = false
	bottom := rectangle.New(DefaultCoordinates, DefaultColors, shader)
	cube := New(bottom, 1, shader)
	log := cube.Log()
	if len(log) < 10 {
		t.Error("Log too short")
	}
}
func TestSetMaterial(t *testing.T) {
	shader.HasTextureValue = false
	bottom := rectangle.New(DefaultCoordinates, DefaultColors, shader)
	cube := New(bottom, 1, shader)
	cube.SetMaterial(material.Gold)
	if cube.material != material.Gold {
		t.Error("Invalid material")
	}
}
func TestSetColor(t *testing.T) {
	shader.HasTextureValue = false
	bottom := rectangle.New(DefaultCoordinates, DefaultColors, shader)
	cube := New(bottom, 1, shader)
	newColor := mgl32.Vec3{0, 1, 0}
	newColors := [4]mgl32.Vec3{newColor, newColor, newColor, newColor}
	cube.SetColor(newColor)
	for i := 0; i < 6; i++ {
		if cube.sides[i].Colors() != newColors {
			t.Error("Invalid color update")
		}
	}
}
func TestSetIndexColor(t *testing.T) {
	shader.HasTextureValue = false
	bottom := rectangle.New(DefaultCoordinates, DefaultColors, shader)
	cube := New(bottom, 1, shader)
	newColor := mgl32.Vec3{1, 1, 0}
	newColors := [4]mgl32.Vec3{DefaultColors[0], newColor, DefaultColors[2], newColor}
	cube.SetIndexColor(1, newColor)
	cube.SetIndexColor(3, newColor)
	for i := 0; i < 6; i++ {
		if cube.sides[i].Colors() != newColors {
			t.Error("Invalid index color update")
		}
	}
}
func TestSetSideColor(t *testing.T) {
	shader.HasTextureValue = false
	bottom := rectangle.New(DefaultCoordinates, DefaultColors, shader)
	cube := New(bottom, 1, shader)
	newColor := mgl32.Vec3{1, 1, 0}
	newColors := [4]mgl32.Vec3{newColor, newColor, newColor, newColor}
	cube.SetSideColor(5, newColor)
	for i := 0; i < 5; i++ {
		if cube.sides[i].Colors() != DefaultColors {
			t.Error("Invalid side color update")
		}
	}
	if cube.sides[5].Colors() != newColors {
		t.Error("Invalid side color update")
	}
}
func TestSetDirection(t *testing.T) {
	shader.HasTextureValue = false
	bottom := rectangle.New(DefaultCoordinates, DefaultColors, shader)
	cube := New(bottom, 1, shader)
	direction := mgl32.Vec3{0, 0, 1}
	speed := float32(10.0)
	cube.SetDirection(direction)
	cube.SetSpeed(speed)
	cube.Update(10)
	expectedCoordinates := [4]mgl32.Vec3{
		mgl32.Vec3{0, 0, 100},
		mgl32.Vec3{1, 0, 100},
		mgl32.Vec3{1, 1, 100},
		mgl32.Vec3{0, 1, 100},
	}
	if cube.sides[0].Coordinates() != expectedCoordinates {
		t.Error("Invalid update after direction & speed setup")
	}
}
func TestSetIndexDirection(t *testing.T) {
	shader.HasTextureValue = false
	bottom := rectangle.New(DefaultCoordinates, DefaultColors, shader)
	cube := New(bottom, 1, shader)
	speed := float32(10.0)
	cube.SetIndexDirection(2, 1)
	cube.SetSpeed(speed)
	cube.Update(10)
	expectedCoordinates := [4]mgl32.Vec3{
		mgl32.Vec3{0, 0, 100},
		mgl32.Vec3{1, 0, 100},
		mgl32.Vec3{1, 1, 100},
		mgl32.Vec3{0, 1, 100},
	}
	if cube.sides[0].Coordinates() != expectedCoordinates {
		t.Error("Invalid update after direction & speed setup")
	}
}
func TestSetSpeed(t *testing.T) {
	shader.HasTextureValue = false
	bottom := rectangle.New(DefaultCoordinates, DefaultColors, shader)
	cube := New(bottom, 1, shader)
	direction := mgl32.Vec3{0, 0, 1}
	speed := float32(10.0)
	cube.SetDirection(direction)
	cube.SetSpeed(speed)
	cube.Update(10)
	expectedCoordinates := [4]mgl32.Vec3{
		mgl32.Vec3{0, 0, 100},
		mgl32.Vec3{1, 0, 100},
		mgl32.Vec3{1, 1, 100},
		mgl32.Vec3{0, 1, 100},
	}
	if cube.sides[0].Coordinates() != expectedCoordinates {
		t.Error("Invalid update after direction & speed setup")
	}
}
func TestSetupVao(t *testing.T) {
	shader.HasTextureValue = false
	bottom := rectangle.New(DefaultCoordinates, DefaultColors, shader)
	cube := New(bottom, 1, shader)
	if len(cube.vao.Get()) != 0 {
		t.Error("Vao is not empty before the first setup.")
	}
	cube.setupVao()
	if len(cube.vao.Get()) == 0 {
		t.Error("Vao is empty after the first setup.")
	}
}
func TestBuildVaoWithoutTexture(t *testing.T) {
	shader.HasTextureValue = false
	bottom := rectangle.New(DefaultCoordinates, DefaultColors, shader)
	cube := New(bottom, 1, shader)
	if len(cube.vao.Get()) != 0 {
		t.Error("Vao is not empty before the first setup.")
	}
	cube.buildVaoWithoutTexture()
	if len(cube.vao.Get()) == 0 {
		t.Error("Vao is empty after the first setup.")
	}
}
func TestDraw(t *testing.T) {
	shader.HasTextureValue = false
	bottom := rectangle.New(DefaultCoordinates, DefaultColors, shader)
	cube := New(bottom, 1, shader)
	if len(cube.vao.Get()) != 0 {
		t.Error("Vao is not empty before the first setup.")
	}
	cube.Draw()
	if len(cube.vao.Get()) == 0 {
		t.Error("Vao is empty after the first setup.")
	}
}
func TestDrawTexture(t *testing.T) {
	shader.HasTextureValue = true
	bottom := rectangle.New(DefaultCoordinates, DefaultColors, shader)
	cube := New(bottom, 1, shader)
	if len(cube.vao.Get()) != 0 {
		t.Error("Vao is not empty before the first setup.")
	}
	cube.Draw()
	if len(cube.vao.Get()) == 0 {
		t.Error("Vao is empty after the first setup.")
	}
}
func TestDrawWithUniformsTexture(t *testing.T) {
	shader.HasTextureValue = true
	bottom := rectangle.New(DefaultCoordinates, DefaultColors, shader)
	cube := New(bottom, 1, shader)
	if len(cube.vao.Get()) != 0 {
		t.Error("Vao is not empty before the first setup.")
	}
	cube.DrawWithUniforms(mgl32.Ident4(), mgl32.Ident4())
	if len(cube.vao.Get()) == 0 {
		t.Error("Vao is empty after the first setup.")
	}
}
func TestDrawWithUniforms(t *testing.T) {
	shader.HasTextureValue = false
	bottom := rectangle.New(DefaultCoordinates, DefaultColors, shader)
	cube := New(bottom, 1, shader)
	if len(cube.vao.Get()) != 0 {
		t.Error("Vao is not empty before the first setup.")
	}
	cube.DrawWithUniforms(mgl32.Ident4(), mgl32.Ident4())
	if len(cube.vao.Get()) == 0 {
		t.Error("Vao is empty after the first setup.")
	}
}
func TestUpdate(t *testing.T) {
	shader.HasTextureValue = false
	bottom := rectangle.New(DefaultCoordinates, DefaultColors, shader)
	cube := New(bottom, 1, shader)
	direction := mgl32.Vec3{0, 0, 1}
	speed := float32(10.0)
	cube.SetDirection(direction)
	cube.SetSpeed(speed)
	cube.Update(10)
	expectedCoordinates := [4]mgl32.Vec3{
		mgl32.Vec3{0, 0, 100},
		mgl32.Vec3{1, 0, 100},
		mgl32.Vec3{1, 1, 100},
		mgl32.Vec3{0, 1, 100},
	}
	if cube.sides[0].Coordinates() != expectedCoordinates {
		t.Error("Invalid update after direction & speed setup")
	}
}
func TestSetPrecision(t *testing.T) {
	shader.HasTextureValue = false
	bottom := rectangle.New(DefaultCoordinates, DefaultColors, shader)
	cube := New(bottom, 1, shader)
	cube.SetPrecision(1)
	cube.setupVao()
	if len(cube.vao.Get()) != 216 {
		t.Error("Invalid number of elements in the vao.")
	}
	cube.SetPrecision(2)
	cube.setupVao()
	if len(cube.vao.Get()) != 864 {
		t.Error("Invalid number of elements in the vao.")
	}
}

func TestSetRotationAngle(t *testing.T) {
	shader.HasTextureValue = false
	bottom := rectangle.New(DefaultCoordinates, DefaultColors, shader)
	cube := New(bottom, 1, shader)
	cube.SetAngle(float32(1.0))
	if cube.angle != float32(1.0) {
		t.Error("Invalid angle")
	}
}
func TestSetRotationAxis(t *testing.T) {
	shader.HasTextureValue = false
	bottom := rectangle.New(DefaultCoordinates, DefaultColors, shader)
	cube := New(bottom, 1, shader)
	axis := mgl32.Vec3{0, 1, 0}
	cube.SetAxis(axis)
	if cube.axis != axis {
		t.Error("Invalid axis")
	}
}
func TestDrawMode(t *testing.T) {
	shader.HasTextureValue = false
	bottom := rectangle.New(DefaultCoordinates, DefaultColors, shader)
	cube := New(bottom, 1, shader)

	if cube.drawMode != 0 {
		t.Errorf("Invalid default draw mode. Instead of '0', we got '%d'", cube.drawMode)
	}
	cube.DrawMode(2) // should keep the original value
	if cube.drawMode != 0 {
		t.Errorf("Invalid default draw mode. Instead of '0', we got '%d'", cube.drawMode)
	}
	cube.DrawMode(1) // should update the original value
	if cube.drawMode != 1 {
		t.Errorf("Invalid default draw mode. Instead of '1', we got '%d'", cube.drawMode)
	}
}
func TestDrawWithUniformsLight(t *testing.T) {
	shader.HasTextureValue = false
	bottom := rectangle.New(DefaultCoordinates, DefaultColors, shader)
	cube := New(bottom, 1, shader)
	cube.DrawMode(1) // set light mode
	if len(cube.vao.Get()) != 0 {
		t.Error("Vao is not empty before the first setup.")
	}
	cube.DrawWithUniforms(mgl32.Ident4(), mgl32.Ident4())
	if len(cube.vao.Get()) != 36*6 {
		t.Errorf("Invalid vao length. Instead of '216', we got '%d'", len(cube.vao.Get()))
	}
}
