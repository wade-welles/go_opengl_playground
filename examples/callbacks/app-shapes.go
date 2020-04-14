package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/akosgarai/opengl_playground/pkg/primitives/camera"
	"github.com/akosgarai/opengl_playground/pkg/primitives/cube"
	sp "github.com/akosgarai/opengl_playground/pkg/primitives/sphere"
	trans "github.com/akosgarai/opengl_playground/pkg/primitives/transformations"
	vec "github.com/akosgarai/opengl_playground/pkg/primitives/vector"
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
	DebugPrint = false
)

type Application struct {
	sphere               *sp.Sphere
	cube                 *cube.Cube
	cubePosition         vec.Vector
	camera               *camera.Camera
	cameraDirection      float64
	cameraDirectionSpeed float64
	cameraLastUpdate     int64
	worldUpDirection     *vec.Vector
	moveSpeed            float64
	epsilon              float64

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
	app.camera = camera.NewCamera(vec.Vector{-3, -5, 18.0}, vec.Vector{0, 1, 0}, -90.0, 0.0)
	app.cameraDirection = 0.1
	app.cameraDirectionSpeed = 5
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
	app.KeyDowns["dLeft"] = false
	app.KeyDowns["dRight"] = false
	app.KeyDowns["dUp"] = false
	app.KeyDowns["dDown"] = false

	return &app
}

// It generates a cube.
func (a *Application) GenerateCube() {
	a.cube = cube.NewCubeByVectorAndLength(vec.Vector{-0.5, -0.5, -0.5}, 1.0)
	a.cubePosition = vec.Vector{0, 0, 1}
}

// It generates a Sphere.
func (a *Application) GenerateSphere() {
	a.sphere = sp.NewSphere()
	a.sphere.SetCenter(vec.Vector{0, 0, 5})
	a.sphere.SetColor(vec.Vector{0, 0, 1})
	a.sphere.SetRadius(2.0)
}

// KeyCallback is responsible for the keyboard event handling.
func (a *Application) KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	switch key {
	case glfw.KeyH:
		if action != glfw.Release {
			fmt.Printf("app.camera: %s\n", a.camera.Log())
			fmt.Printf("app.cube: %v\n", a.cube)
			fmt.Printf("app.cubePosition: %v\n", a.cubePosition)
		}
		break
	case glfw.KeyW:
		if action != glfw.Release {
			a.KeyDowns["W"] = true
		} else {
			a.KeyDowns["W"] = false
		}
		break
	case glfw.KeyA:
		if action != glfw.Release {
			a.KeyDowns["A"] = true
		} else {
			a.KeyDowns["A"] = false
		}
		break
	case glfw.KeyS:
		if action != glfw.Release {
			a.KeyDowns["S"] = true
		} else {
			a.KeyDowns["S"] = false
		}
		break
	case glfw.KeyD:
		if action != glfw.Release {
			a.KeyDowns["D"] = true
		} else {
			a.KeyDowns["D"] = false
		}
		break
	case glfw.KeyQ:
		if action != glfw.Release {
			a.KeyDowns["Q"] = true
		} else {
			a.KeyDowns["Q"] = false
		}
		break
	case glfw.KeyE:
		if action != glfw.Release {
			a.KeyDowns["E"] = true
		} else {
			a.KeyDowns["E"] = false
		}
		break
	}
}
func (a *Application) Update() {
	//calculate delta
	nowUnix := time.Now().UnixNano()
	delta := nowUnix - a.cameraLastUpdate
	moveTime := float64(delta / int64(time.Millisecond))
	// rotate camera
	currX, currY := a.window.GetCursorPos()
	x, y := trans.MouseCoordinates(currX, currY, windowWidth, windowHeight)
	// dUp
	if y > 1.0-a.cameraDirection && y < 1.0 {
		a.KeyDowns["dUp"] = true
	} else {
		a.KeyDowns["dUp"] = false
	}
	// dDown
	if y < -1.0+a.cameraDirection && y > -1.0 {
		a.KeyDowns["dDown"] = true
	} else {
		a.KeyDowns["dDown"] = false
	}
	// dLeft
	if x < -1.0+a.cameraDirection && x > -1.0 {
		a.KeyDowns["dLeft"] = true
	} else {
		a.KeyDowns["dLeft"] = false
	}
	// dRight
	if x > 1.0-a.cameraDirection && x < 1.0 {
		a.KeyDowns["dRight"] = true
	} else {
		a.KeyDowns["dRight"] = false
	}

	dX := 0.0
	dY := 0.0
	if a.KeyDowns["dUp"] && !a.KeyDowns["dDown"] {
		dY = 0.01 * a.cameraDirectionSpeed
	} else if a.KeyDowns["dDown"] && !a.KeyDowns["dUp"] {
		dY = -0.01 * a.cameraDirectionSpeed
	}
	if a.KeyDowns["dLeft"] && !a.KeyDowns["dRight"] {
		dX = -0.01 * a.cameraDirectionSpeed
	} else if a.KeyDowns["dRight"] && !a.KeyDowns["dLeft"] {
		dX = 0.01 * a.cameraDirectionSpeed
	}
	a.camera.UpdateDirection(dX, dY)
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
func (a *Application) Draw() {
	a.Update()
	modelLocation := gl.GetUniformLocation(a.program, gl.Str("model\x00"))
	// view aka camera
	viewLocation := gl.GetUniformLocation(a.program, gl.Str("view\x00"))
	projectionLocation := gl.GetUniformLocation(a.program, gl.Str("projection\x00"))

	// mvp - modelview - projection matrix
	V := a.camera.GetViewMatrix().GetMatrix()
	gl.UniformMatrix4fv(viewLocation, 1, false, &V[0])
	// Should Be fine 'P'
	P := a.camera.GetProjectionMatrix().GetMatrix()
	gl.UniformMatrix4fv(projectionLocation, 1, false, &P[0])
	M := trans.TranslationMatrix(float32(a.cubePosition.X), float32(a.cubePosition.Y), float32(a.cubePosition.Z)).TransposeMatrix().GetMatrix()
	gl.UniformMatrix4fv(modelLocation, 1, false, &M[0])
	a.cube.Draw()
	// Update M for the sphere
	M = trans.TranslationMatrix(
		float32(a.sphere.GetCenter().X),
		float32(a.sphere.GetCenter().Y),
		float32(a.sphere.GetCenter().Z)).Dot(trans.ScaleMatrix(
		float32(a.sphere.GetRadius()),
		float32(a.sphere.GetRadius()),
		float32(a.sphere.GetRadius()))).TransposeMatrix().GetMatrix()
	gl.UniformMatrix4fv(modelLocation, 1, false, &M[0])
	a.sphere.Draw()
}

// MouseButtonCallback is responsible for the mouse button event handling.
func (a *Application) MouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	fmt.Printf("MouseButtonCallback has been called with the following options: button: '%d', action: '%d'!, mods: '%d'\n", button, action, mods)
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

	// register keyboard button callback
	app.window.SetKeyCallback(app.KeyCallback)
	// register mouse button callback
	app.window.SetMouseButtonCallback(app.MouseButtonCallback)

	for !app.window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		// Update the camera.
		app.Draw()
		glfw.PollEvents()
		app.window.SwapBuffers()
	}
}
