package main

import (
	"math/rand"
	"runtime"

	"github.com/akosgarai/opengl_playground/pkg/application"
	wrapper "github.com/akosgarai/opengl_playground/pkg/glwrapper"
	"github.com/akosgarai/opengl_playground/pkg/mesh"
	trans "github.com/akosgarai/opengl_playground/pkg/primitives/transformations"
	"github.com/akosgarai/opengl_playground/pkg/primitives/vertex"
	"github.com/akosgarai/opengl_playground/pkg/shader"
	"github.com/akosgarai/opengl_playground/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth  = 800
	WindowHeight = 800
	WindowTitle  = "Example - draw points from mouse inputs"

	LEFT_MOUSE_BUTTON = glfw.MouseButtonLeft
)

var (
	addPoint = false

	app *application.Application

	PointMesh *mesh.PointMesh
)

func Update() {
	if !app.GetMouseButtonState(LEFT_MOUSE_BUTTON) && addPoint {
		mX, mY := trans.MouseCoordinates(app.MousePosX, app.MousePosY, WindowWidth, WindowHeight)
		coords := mgl32.Vec3{float32(mX), float32(mY), 0.0}
		color := mgl32.Vec3{rand.Float32(), rand.Float32(), rand.Float32()}
		size := float32(3 + rand.Intn(17))
		vert := vertex.Vertex{
			Position:  coords,
			Color:     color,
			PointSize: size,
		}
		PointMesh.AddVertex(vert)

		addPoint = false
	} else if app.GetMouseButtonState(LEFT_MOUSE_BUTTON) {
		addPoint = true
	}
}

func main() {
	runtime.LockOSThread()
	app = application.New()
	app.SetWindow(window.InitGlfw(WindowWidth, WindowHeight, WindowTitle))
	defer glfw.Terminate()
	wrapper.InitOpenGL()

	shaderProgram := shader.NewShader("examples/06-draw-points-from-mouse-input/shaders/vertexshader.vert", "examples/06-draw-points-from-mouse-input/shaders/fragmentshader.frag")
	app.AddShader(shaderProgram)

	PointMesh = mesh.NewPointMesh()
	app.AddMeshToShader(PointMesh, shaderProgram)

	app.GetWindow().SetMouseButtonCallback(app.MouseButtonCallback)
	app.GetWindow().SetKeyCallback(app.KeyCallback)

	wrapper.Enable(wrapper.PROGRAM_POINT_SIZE)
	wrapper.Enable(wrapper.DEPTH_TEST)
	wrapper.DepthFunc(wrapper.LESS)
	wrapper.ClearColor(0.3, 0.3, 0.3, 1.0)

	for !app.GetWindow().ShouldClose() {
		wrapper.Clear(wrapper.COLOR_BUFFER_BIT | wrapper.DEPTH_BUFFER_BIT)
		glfw.PollEvents()
		Update()
		app.Draw()
		app.GetWindow().SwapBuffers()
	}
}
