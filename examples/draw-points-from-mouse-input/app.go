package main

import (
	"fmt"
	"runtime"

	"github.com/akosgarai/opengl_playground/pkg/primitives"
	"github.com/akosgarai/opengl_playground/pkg/shader"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	windowWidth  = 800
	windowHeight = 800
	windowTitle  = "Example - draw points from mouse inputs"
)

var (
	DebugPrint         = false
	mouseButtonPressed = false
	mousePositionX     = 0.0
	mousePositionY     = 0.0
)

type Application struct {
	Points []primitives.Point
}

func (a *Application) AddPoint(point primitives.Point) {
	a.Points = append(a.Points, point)
}

var app Application

func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(fmt.Errorf("could not initialize glfw: %v", err))
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(windowWidth, windowHeight, windowTitle, nil, nil)

	if err != nil {
		panic(fmt.Errorf("could not create opengl renderer: %v", err))
	}

	window.MakeContextCurrent()

	return window
}

func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	vertexShader, err := shader.CompileShader(shader.VertexShaderPointSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := shader.CompileShader(shader.FragmentShaderConstantSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)
	return program
}

/*
* Mouse click handler logic:
* - if the left mouse button is not pressed, and the button is just released, App.AddPoint(), clean up the temp.point.
* - if the button is just pressed, set the point that needs to be added.
 */
func mouseHandler(window *glfw.Window) {
	x, y := window.GetCursorPos()

	if window.GetMouseButton(glfw.MouseButtonMiddle) == glfw.Press {
		if !mouseButtonPressed {
			mousePositionX = x
			mousePositionY = y
			mouseButtonPressed = true
		}
	} else {
		if mouseButtonPressed {
			mouseButtonPressed = false
			x, y := convertMouseCoordinates()
			app.AddPoint(
				primitives.Point{
					primitives.Vector{x, y, 0.0},
					primitives.Vector{1, 1, 1},
				})
		}
	}
}

// Key handler function. it supports the debug option. (print out the points of the app)
func keyHandler(window *glfw.Window) {
	if window.GetKey(glfw.KeyD) == glfw.Press {
		if !DebugPrint {
			DebugPrint = true
			fmt.Println(app.Points)
		}
	} else {
		DebugPrint = false
	}
}
func convertMouseCoordinates() (float64, float64) {
	halfWidth := windowWidth / 2.0
	halfHeight := windowHeight / 2.0
	x := (mousePositionX - halfWidth) / (halfWidth)
	y := (halfHeight - mousePositionY) / (halfHeight)
	return x, y
}
func buildVAO() []float32 {
	var vao []float32
	for _, item := range app.Points {
		vao = append(vao, float32(item.Coordinate.X))
		vao = append(vao, float32(item.Coordinate.Y))
		vao = append(vao, float32(item.Coordinate.Z))
	}
	return vao
}
func Draw() {
	if len(app.Points) < 1 {
		return
	}
	points := buildVAO()
	var vertexBufferObject uint32
	gl.GenBuffers(1, &vertexBufferObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	// a 32-bit float has 4 bytes, so we are saying the size of the buffer,
	// in bytes, is 4 times the number of points
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vertexArrayObject uint32
	gl.GenVertexArrays(1, &vertexArrayObject)
	gl.BindVertexArray(vertexArrayObject)
	// setup points
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 4*3, gl.PtrOffset(0))

	gl.BindVertexArray(vertexArrayObject)
	gl.DrawArrays(gl.POINTS, 0, int32(len(app.Points)))
}

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()
	program := initOpenGL()

	gl.UseProgram(program)

	gl.Enable(gl.PROGRAM_POINT_SIZE)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)
		mouseHandler(window)
		keyHandler(window)
		Draw()
		glfw.PollEvents()
		window.SwapBuffers()
	}
}
