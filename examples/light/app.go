package main

import (
	"runtime"
	"time"

	"github.com/akosgarai/opengl_playground/examples/light/application"
	"github.com/akosgarai/opengl_playground/examples/light/primitives"
	"github.com/akosgarai/opengl_playground/pkg/shader"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	windowWidth  = 800
	windowHeight = 800
	windowTitle  = "Example - the house"
	moveSpeed    = 1.0 / 100.0
	epsilon      = 1000.0
	// buttons
	FORWARD  = glfw.KeyW // Go forward
	BACKWARD = glfw.KeyS // Go backward
	LEFT     = glfw.KeyA // Turn 90 deg. left
	RIGHT    = glfw.KeyD // Turn 90 deg. right
)

var (
	cameraLastUpdate int64
	app              *application.Application
)

// It creates a new camera with the necessary setup
func CreateCamera() *primitives.Camera {
	camera := primitives.NewCamera(mgl32.Vec3{75, 30, 0.0}, mgl32.Vec3{0, -1, 0}, 90.0, 0.0)
	camera.SetupProjection(45, float32(windowWidth)/float32(windowHeight), 0.1, 1000.0)
	return camera
}

// Create the keymap
func SetupKeyMap() map[glfw.Key]bool {
	keyDowns := make(map[glfw.Key]bool)
	keyDowns[FORWARD] = false
	keyDowns[LEFT] = false
	keyDowns[RIGHT] = false
	keyDowns[BACKWARD] = false

	return keyDowns
}

// the ground
func Ground(shaderProgram *shader.Shader) {
	coordinates := [4]mgl32.Vec3{
		mgl32.Vec3{-1000, -0.02, 1000},
		mgl32.Vec3{-1000, -0.02, -1000},
		mgl32.Vec3{1000, -0.02, -1000},
		mgl32.Vec3{1000, -0.02, 1000},
	}
	material := primitives.TestMaterialGreen
	rect := primitives.NewRectangle(coordinates, material, 20, shaderProgram)
	rect.SetInvertNormal(true)
	app.AddItem(rect)
}

// the Floor of the building
func Floor(shaderProgram *shader.Shader) {
	coordinates := [4]mgl32.Vec3{
		mgl32.Vec3{-40, 0, -20},
		mgl32.Vec3{-40, 0, 80},
		mgl32.Vec3{60, 0, 80},
		mgl32.Vec3{60, 0, -20},
	}
	material := primitives.TestMaterialRed
	rect := primitives.NewRectangle(coordinates, material, 20, shaderProgram)
	rect.SetInvertNormal(false)
	app.AddItem(rect)
}

// the left wall of the building. 1 width. 100 lenght, 50 height. I decided to draw it from 10 cuboids.
func LeftWall(shaderProgram *shader.Shader) {
	material := primitives.TestMaterialRed
	farest := -20
	step := 10
	for i := 0; i < 10; i++ {
		// cuboid is based on the bottom rect.
		z0 := float32(farest + (step * i))
		z1 := float32(farest + (step * (i + 1)))
		bottomSide := [4]mgl32.Vec3{
			mgl32.Vec3{-40, 0, z0},
			mgl32.Vec3{-39, 0, z0},
			mgl32.Vec3{-39, 0, z1},
			mgl32.Vec3{-40, 0, z1},
		}
		bottomRect := primitives.NewRectangle(bottomSide, material, 10, shaderProgram)
		wall := primitives.NewCuboid(bottomRect, 50.0, material, 10, shaderProgram)
		app.AddItem(wall)
	}
}

// the right wall of the building. 1 width.
func RightWall(shaderProgram *shader.Shader) {
	material := primitives.TestMaterialRed
	farest := -20
	step := 10
	for i := 0; i < 10; i++ {
		// cuboid is based on the bottom rect.
		z0 := float32(farest + (step * i))
		z1 := float32(farest + (step * (i + 1)))
		bottomSide := [4]mgl32.Vec3{
			mgl32.Vec3{59, 0, z0},
			mgl32.Vec3{60, 0, z0},
			mgl32.Vec3{60, 0, z1},
			mgl32.Vec3{59, 0, z1},
		}
		bottomRect := primitives.NewRectangle(bottomSide, material, 10, shaderProgram)
		wall := primitives.NewCuboid(bottomRect, 50.0, material, 10, shaderProgram)
		app.AddItem(wall)
	}
}

func Update() {
	//calculate delta
	nowUnix := time.Now().UnixNano()
	delta := nowUnix - cameraLastUpdate
	moveTime := float64(delta / int64(time.Millisecond))

	if epsilon > moveTime {
		return
	}
	cameraLastUpdate = nowUnix

	forward := 0.0
	if app.GetKeyState(FORWARD) && !app.GetKeyState(BACKWARD) {
		forward = moveSpeed * moveTime
	} else if app.GetKeyState(BACKWARD) && !app.GetKeyState(FORWARD) {
		forward = -moveSpeed * moveTime
	}
	if forward != 0 {
		app.GetCamera().Walk(float32(forward))
	}
	dX := float32(0.0)
	dY := float32(0.0)
	if app.GetKeyState(LEFT) && !app.GetKeyState(RIGHT) {
		dX = -90
	} else if app.GetKeyState(RIGHT) && !app.GetKeyState(LEFT) {
		dX = 90
	}
	if dX != 0.0 {
		app.GetCamera().UpdateDirection(dX, dY)
	}
}
func main() {
	runtime.LockOSThread()

	app = application.New()

	app.SetWindow(application.InitGlfw(windowWidth, windowHeight, windowTitle))
	defer glfw.Terminate()
	application.InitOpenGL()

	shaderProgram := shader.NewShader("examples/light/vertexshader.vert", "examples/light/fragmentshader.frag")

	app.SetCamera(CreateCamera())
	cameraLastUpdate = time.Now().UnixNano()

	app.SetKeys(SetupKeyMap())
	Ground(shaderProgram)
	Floor(shaderProgram)
	LeftWall(shaderProgram)
	RightWall(shaderProgram)

	gl.ClearColor(0.3, 0.3, 0.3, 1.0)
	gl.Viewport(0, 0, windowWidth, windowHeight)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	// register keyboard button callback
	app.GetWindow().SetKeyCallback(app.KeyCallback)
	// register mouse button callback
	app.GetWindow().SetMouseButtonCallback(app.DummyMouseButtonCallback)

	for !app.GetWindow().ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		glfw.PollEvents()
		Update()
		app.Draw()
		app.GetWindow().SwapBuffers()
	}
}