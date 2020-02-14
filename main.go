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
	windowHeight = 600
	windowTitle  = "Example - dynamically subdivision"
)

var (
	square = primitives.NewSquare(
		primitives.Vector{-5, 0, -5},
		primitives.Vector{-5, 0, 5},
		primitives.Vector{5, 0, 5},
		primitives.Vector{5, 0, -5},
	)
	sub_division       = 1.0
	mouseButtonPressed = false
	mousePositionX     = 0
	mousePositionY     = 0
	rotationX          = 45.0
	rotationY          = 45.0
	distance           = 1.0
)

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

	vertexShader, err := shader.CompileShader(shader.VertexShaderDirectOutputSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	geometryShader, err := shader.CompileShader(shader.GeometryShaderQuadSubdivisionSource, gl.GEOMETRY_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := shader.CompileShader(shader.FragmentShaderConstantSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, geometryShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)
	return program
}

/*
* Mouse click handler logic:
* - if the button is pressed: state is on(1), else off(0)
* - if the state is on(1) - we can store the old x,y coordinates
* Mouse move handler logic:
* - if the state is on: we can calculate the dist. variable.
* - if the state is off: we can calculate rx, rY rotation values.
 */
func mouseHandler(window *glfw.Window) {
	x, y := window.GetCursorPos()

	if window.GetMouseButton(glfw.MouseButtonMiddle) == glfw.Press {
		if !mouseButtonPressed {
			mousePositionX = int(x)
			mousePositionY = int(y)
			mouseButtonPressed = true
			fmt.Print("Button Pressed x: ")
			fmt.Print(x)
			fmt.Print(", y: ")
			fmt.Print(y)
		}
	} else {
		mouseButtonPressed = false
	}
	if mousePositionX != int(x) || mousePositionY != int(y) {
		if mouseButtonPressed {
			distance *= float64((1 + (int(y)-mousePositionY)/60.0))
			fmt.Print("Distance : ")
			fmt.Print(distance)
		} else {
			rotationY += float64((int(x) - mousePositionX) / 5.0)
			rotationX += float64((int(y) - mousePositionY) / 5.0)
		}
		mousePositionX = int(x)
		mousePositionY = int(y)
	}
}
func keyHandler(window *glfw.Window) {
	if window.GetKey(glfw.KeyT) == glfw.Press {
		sub_division = sub_division + 1.0
	}
	if window.GetKey(glfw.KeyR) == glfw.Press {
		sub_division = sub_division - 1.0
	}
	if sub_division < 1 {
		sub_division = 1.0
	}
	if sub_division > 8 {
		sub_division = 8.0
	}
	if window.GetKey(glfw.KeyO) == glfw.Press {
		x, y := window.GetCursorPos()
		fmt.Print("x: ")
		fmt.Print(x)
		fmt.Print(", y: ")
		fmt.Print(y)
	}
}

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()
	program := initOpenGL()
	square.SetupVaoPoligonMode()

	// Configure global settings
	gl.UseProgram(program)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	// projection matrix
	angelOfView := float64(45)
	near := float64(1)
	far := float64(1000)
	// P = glm::perspective(45.0f, (GLfloat)w/h, 1.f, 1000.f);
	P := primitives.ProjectionMatrix4x4(angelOfView, near, far)

	sub_divisionLocation := gl.GetUniformLocation(program, gl.Str("sub_divisions\x00"))
	mvpLocation := gl.GetUniformLocation(program, gl.Str("MVP\x00"))

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		// define the matrixes:
		translationMatrix := primitives.TranslationMatrix4x4(0, 0, float32(distance))
		rotationXMatrix := primitives.RotationXMatrix4x4(rotationX)
		rotationYMatrix := primitives.RotationYMatrix4x4(rotationY)

		MV := translationMatrix.Dot(rotationXMatrix.Dot(rotationYMatrix))

		translationMatrix = primitives.TranslationMatrix4x4(-5, 0, -5)
		MV = MV.Dot(translationMatrix)
		// sub_division
		gl.Uniform1f(sub_divisionLocation, float32(sub_division))
		mvpPoints := (P.Dot(MV)).Points
		gl.UniformMatrix4fv(mvpLocation, 1, false, &mvpPoints[0])

		gl.DrawElements(gl.TRIANGLES, 6, gl.FLOAT, gl.PtrOffset(0))

		translationMatrix = primitives.TranslationMatrix4x4(10, 0, 0)
		MV = MV.Dot(translationMatrix)
		mvpPoints = (P.Dot(MV)).Points
		gl.UniformMatrix4fv(mvpLocation, 1, false, &mvpPoints[0])
		gl.DrawElements(gl.TRIANGLES, 6, gl.FLOAT, gl.PtrOffset(0))

		translationMatrix = primitives.TranslationMatrix4x4(0, 0, 10)
		MV = MV.Dot(translationMatrix)
		mvpPoints = (P.Dot(MV)).Points
		gl.UniformMatrix4fv(mvpLocation, 1, false, &mvpPoints[0])
		gl.DrawElements(gl.TRIANGLES, 6, gl.FLOAT, gl.PtrOffset(0))

		translationMatrix = primitives.TranslationMatrix4x4(-10, 0, 0)
		MV = MV.Dot(translationMatrix)
		mvpPoints = (P.Dot(MV)).Points
		gl.UniformMatrix4fv(mvpLocation, 1, false, &mvpPoints[0])
		gl.DrawElements(gl.TRIANGLES, 6, gl.FLOAT, gl.PtrOffset(0))

		keyHandler(window)
		mouseHandler(window)
		glfw.PollEvents()
		window.SwapBuffers()
	}
}
