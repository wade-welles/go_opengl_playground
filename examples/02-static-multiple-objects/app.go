package main

import (
	"runtime"

	"github.com/akosgarai/opengl_playground/pkg/application"
	wrapper "github.com/akosgarai/opengl_playground/pkg/glwrapper"
	"github.com/akosgarai/opengl_playground/pkg/primitives/rectangle"
	"github.com/akosgarai/opengl_playground/pkg/primitives/triangle"
	"github.com/akosgarai/opengl_playground/pkg/shader"
	"github.com/akosgarai/opengl_playground/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth  = 800
	WindowHeight = 600
	WindowTitle  = "Example - static triangle and square"
)

var (
	triangleCoordinates = [3]mgl32.Vec3{
		mgl32.Vec3{-0.75, 0.75, 0}, // top
		mgl32.Vec3{-0.75, 0.25, 0}, // left
		mgl32.Vec3{-0.25, 0.25, 0}, // right
	}
	triangleColors = [3]mgl32.Vec3{
		mgl32.Vec3{0, 1, 0}, // top
		mgl32.Vec3{0, 1, 0}, // left
		mgl32.Vec3{0, 1, 0}, // right
	}
	squareCoordinates = [4]mgl32.Vec3{
		mgl32.Vec3{0.25, -0.25, 0}, // top-left
		mgl32.Vec3{0.25, -0.75, 0}, // bottom-left
		mgl32.Vec3{0.75, -0.75, 0}, // bottom-right
		mgl32.Vec3{0.75, -0.25, 0}, // top-right
	}
	squareColors = [4]mgl32.Vec3{
		mgl32.Vec3{0, 1, 0},
		mgl32.Vec3{0, 1, 0},
		mgl32.Vec3{0, 1, 0},
		mgl32.Vec3{0, 1, 0},
	}

	app *application.Application
)

func main() {
	runtime.LockOSThread()

	app = application.New()
	app.SetWindow(window.InitGlfw(WindowWidth, WindowHeight, WindowTitle))
	defer glfw.Terminate()
	wrapper.InitOpenGL()

	shaderProgram := shader.NewShader("examples/02-static-multiple-objects/vertexshader.vert", "examples/02-static-multiple-objects/fragmentshader.frag")

	triang := triangle.New(triangleCoordinates, triangleColors, shaderProgram)
	app.AddItem(triang)
	square := rectangle.New(squareCoordinates, squareColors, shaderProgram)
	app.AddItem(square)

	wrapper.Enable(wrapper.DEPTH_TEST)
	wrapper.DepthFunc(wrapper.LESS)

	for !app.GetWindow().ShouldClose() {
		wrapper.Clear(wrapper.COLOR_BUFFER_BIT | wrapper.DEPTH_BUFFER_BIT)
		app.Draw()
		glfw.PollEvents()
		app.GetWindow().SwapBuffers()
	}
}
