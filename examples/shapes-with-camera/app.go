package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/akosgarai/opengl_playground/pkg/primitives"
	"github.com/akosgarai/opengl_playground/pkg/shader"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	windowWidth  = 800
	windowHeight = 800
	windowTitle  = "Example - shapes with camera"
)

var (
	cameraRotate = false
	DebugPrint   = false
)

type Application struct {
	sphere           *primitives.Sphere
	cube             *primitives.Cube
	cubePosition     primitives.Vector
	camera           *primitives.Camera
	cameraLastUpdate int64
	worldUpDirection *primitives.Vector
	moveSpeed        float64
	epsilon          float64

	window  *glfw.Window
	program uint32

	KeyDowns map[string]bool
}

func NewApplication() *Application {
	var app Application
	app.GenerateCube()
	app.GenerateSphere()
	app.moveSpeed = 1.0 / 1000.0
	app.epsilon = 50.0
	app.camera = primitives.NewCamera(primitives.Vector{-3, -5, 18.0}, primitives.Vector{0, 1, 0}, -90.0, 0.0)
	fmt.Println("Camera state after new function")
	fmt.Println(app.camera.Log())
	// Rotation related code comes here.
	app.camera.SetupProjection(45, float64(windowWidth/windowHeight), 0.1, 100.0)
	fmt.Println("Camera state after setupProjection function")
	fmt.Println(app.camera.Log())
	app.cameraLastUpdate = time.Now().UnixNano()
	app.KeyDowns = make(map[string]bool)
	app.KeyDowns["W"] = false
	app.KeyDowns["A"] = false
	app.KeyDowns["S"] = false
	app.KeyDowns["D"] = false
	app.KeyDowns["Q"] = false
	app.KeyDowns["E"] = false
	return &app
}

// It generates a cube.
func (a *Application) GenerateCube() {
	a.cube = primitives.NewCubeByVectorAndLength(primitives.Vector{-0.5, -0.5, -0.5}, 1.0)
	a.cubePosition = primitives.Vector{0, 0, 1}
}

// It generates a Sphere.
func (a *Application) GenerateSphere() {
	a.sphere = primitives.NewSphere()
	a.sphere.SetCenter(primitives.Vector{0, 0, 5})
	a.sphere.SetColor(primitives.Vector{0, 0, 1})
	a.sphere.SetRadius(2.0)
}

// Key handler function. it supports the debug option. (print out the points of the app)
func (a *Application) KeyHandler() {
	if a.window.GetKey(glfw.KeyH) == glfw.Press {
		if !DebugPrint {
			DebugPrint = true
			fmt.Printf("app.camera: %s\n", a.camera.Log())
			fmt.Printf("app.cube: %v\n", a.cube)
			fmt.Printf("app.cubePosition: %v\n", a.cubePosition)
		}
	} else {
		DebugPrint = false
	}
	if a.window.GetKey(glfw.KeyW) == glfw.Press {
		a.KeyDowns["W"] = true
	} else {
		a.KeyDowns["W"] = false
	}
	if a.window.GetKey(glfw.KeyA) == glfw.Press {
		a.KeyDowns["A"] = true
	} else {
		a.KeyDowns["A"] = false
	}
	if a.window.GetKey(glfw.KeyS) == glfw.Press {
		a.KeyDowns["S"] = true
	} else {
		a.KeyDowns["S"] = false
	}
	if a.window.GetKey(glfw.KeyD) == glfw.Press {
		a.KeyDowns["D"] = true
	} else {
		a.KeyDowns["D"] = false
	}
	if a.window.GetKey(glfw.KeyQ) == glfw.Press {
		a.KeyDowns["Q"] = true
	} else {
		a.KeyDowns["Q"] = false
	}
	if a.window.GetKey(glfw.KeyE) == glfw.Press {
		a.KeyDowns["E"] = true
	} else {
		a.KeyDowns["E"] = false
	}
	//calculate delta
	nowUnix := time.Now().UnixNano()
	delta := nowUnix - a.cameraLastUpdate
	moveTime := float64(delta / int64(time.Millisecond))
	// if the camera has been updated recently, we can skip now
	if a.epsilon > moveTime {
		return
	}
	a.cameraLastUpdate = nowUnix
	// Move camera
	forward := 0.0
	if a.KeyDowns["W"] && !a.KeyDowns["S"] {
		forward = a.moveSpeed * moveTime
	} else if a.KeyDowns["S"] && !a.KeyDowns["W"] {
		forward = -a.moveSpeed * moveTime
	}
	if forward != 0 {
		a.camera.Walk(forward)
	}
	horisontal := 0.0
	if a.KeyDowns["A"] && !a.KeyDowns["D"] {
		horisontal = -a.moveSpeed * moveTime
	} else if a.KeyDowns["D"] && !a.KeyDowns["A"] {
		horisontal = a.moveSpeed * moveTime
	}
	if horisontal != 0 {
		a.camera.Strafe(horisontal)
	}
	vertical := 0.0
	if a.KeyDowns["Q"] && !a.KeyDowns["E"] {
		vertical = -a.moveSpeed * moveTime
	} else if a.KeyDowns["E"] && !a.KeyDowns["Q"] {
		vertical = a.moveSpeed * moveTime
	}
	if vertical != 0 {
		a.camera.Lift(vertical)
	}
}

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

	vertexShader, err := shader.CompileShader(shader.VertexShaderModelViewProjectionSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := shader.CompileShader(shader.FragmentShaderBasicSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)
	gl.ClearColor(0.3, 0.3, 0.3, 1.0)
	gl.Viewport(0, 0, windowWidth, windowHeight)
	return program
}

func main() {
	runtime.LockOSThread()

	app := NewApplication()

	app.window = initGlfw()
	defer glfw.Terminate()
	app.program = initOpenGL()

	// Configure global settings
	gl.UseProgram(app.program)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	modelLocation := gl.GetUniformLocation(app.program, gl.Str("model\x00"))
	// view aka camera
	viewLocation := gl.GetUniformLocation(app.program, gl.Str("view\x00"))
	projectionLocation := gl.GetUniformLocation(app.program, gl.Str("projection\x00"))
	//CubePosition := app.cubePosition

	for !app.window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		// mvp - modelview - projection matrix
		V := app.camera.GetViewMatrix().GetMatrix()
		gl.UniformMatrix4fv(viewLocation, 1, false, &V[0])
		// Should Be fine 'P'
		P := app.camera.GetProjectionMatrix().GetMatrix()
		gl.UniformMatrix4fv(projectionLocation, 1, false, &P[0])
		M := primitives.TranslationMatrix4x4(float32(app.cubePosition.X), float32(app.cubePosition.Y), float32(app.cubePosition.Z)).TransposeMatrix().GetMatrix()
		gl.UniformMatrix4fv(modelLocation, 1, false, &M[0])
		app.cube.Draw()
		// Update M for the sphere
		M = primitives.TranslationMatrix4x4(
			float32(app.sphere.GetCenter().X),
			float32(app.sphere.GetCenter().Y),
			float32(app.sphere.GetCenter().Z)).Dot(primitives.ScaleMatrix4x4(
			float32(app.sphere.GetRadius()),
			float32(app.sphere.GetRadius()),
			float32(app.sphere.GetRadius()))).TransposeMatrix().GetMatrix()
		gl.UniformMatrix4fv(modelLocation, 1, false, &M[0])
		app.sphere.Draw()
		app.KeyHandler()
		glfw.PollEvents()
		app.window.SwapBuffers()
	}
}
