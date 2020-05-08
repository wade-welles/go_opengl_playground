package main

import (
	"runtime"
	"time"

	"github.com/akosgarai/opengl_playground/pkg/application"
	"github.com/akosgarai/opengl_playground/pkg/primitives/camera"
	"github.com/akosgarai/opengl_playground/pkg/primitives/cuboid"
	"github.com/akosgarai/opengl_playground/pkg/primitives/light"
	"github.com/akosgarai/opengl_playground/pkg/primitives/material"
	"github.com/akosgarai/opengl_playground/pkg/primitives/rectangle"
	trans "github.com/akosgarai/opengl_playground/pkg/primitives/transformations"
	"github.com/akosgarai/opengl_playground/pkg/shader"
	"github.com/akosgarai/opengl_playground/pkg/window"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth  = 800
	WindowHeight = 800
	WindowTitle  = "Example - multiple light source"

	FORWARD  = glfw.KeyW
	BACKWARD = glfw.KeyS
	LEFT     = glfw.KeyA
	RIGHT    = glfw.KeyD
	UP       = glfw.KeyQ
	DOWN     = glfw.KeyE

	moveSpeed            = 0.005
	rotationSpeed        = float32(2.0)
	cameraDirectionSpeed = float32(0.00500)
	CameraMoveSpeed      = 0.005
	cameraDistance       = 0.1
)

var (
	app        *application.Application
	lastUpdate int64

	DirectionalLightDirection = (mgl32.Vec3{0.7, 0.7, 0.7}).Normalize()
	DirectionalLightAmbient   = mgl32.Vec3{0.2, 0.2, 0.2}
	DirectionalLightDiffuse   = mgl32.Vec3{0.5, 0.5, 0.5}
	DirectionalLightSpecular  = mgl32.Vec3{1.0, 1.0, 1.0}
	PointLightAmbient         = mgl32.Vec3{1, 1, 1}
	PointLightDiffuse         = mgl32.Vec3{1, 1, 1}
	PointLightSpecular        = mgl32.Vec3{1, 1, 1}
	PointLightDirection_1     = mgl32.Vec3{1, 1, 1}
	PointLightDirection_2     = mgl32.Vec3{2, 2, 2}
	LightConstantTerm         = float32(1.0)
	LightLinearTerm           = float32(1.0)
	LightQuadraticTerm        = float32(1.0)
	SpotLightAmbient          = mgl32.Vec3{1, 1, 1}
	SpotLightDiffuse          = mgl32.Vec3{1, 1, 1}
	SpotLightSpecular         = mgl32.Vec3{1, 1, 1}
	SpotLightDirection_1      = (mgl32.Vec3{1, 1, 1}).Normalize()
	SpotLightDirection_2      = (mgl32.Vec3{2, 2, 2}).Normalize()
	SpotLightPosition_1       = mgl32.Vec3{3, 3, 3}
	SpotLightPosition_2       = mgl32.Vec3{4, 4, 4}
	SpotLightCutoff           = float32(20)
	SpotLightOuterCutoff      = float32(22)
)

// It generates a square.
func Grass(shaderProgram *shader.Shader) {
	square := rectangle.NewSquare(mgl32.Vec3{-1000, 0, -1000}, mgl32.Vec3{1000, 0, 1000}, mgl32.Vec3{0, 1, 0}, mgl32.Vec3{0, 1, 0}, shaderProgram)
	square.SetPrecision(1)
	square.DrawMode(cuboid.DRAW_MODE_TEXTURED_LIGHT)
	app.AddItem(square)
}

// It generates the box.
func Box(shaderProgram *shader.Shader) {
	bottomRect := rectangle.NewSquare(mgl32.Vec3{5, 0, -5}, mgl32.Vec3{-5, 0, 5}, mgl32.Vec3{0, 1, 0}, mgl32.Vec3{0.0, 1.0, 1.0}, shaderProgram)
	box := cuboid.New(bottomRect, 10.0, shaderProgram)
	mat := material.New(mgl32.Vec3{1.0, 1.0, 1.0}, mgl32.Vec3{1.0, 1.0, 1.0}, mgl32.Vec3{1, 1, 1}, 64.0)
	box.SetMaterial(mat)
	box.DrawMode(cuboid.DRAW_MODE_TEXTURED_LIGHT)
	app.AddItem(box)
}

// It creates a new camera with the necessary setup
func CreateCamera() *camera.Camera {
	camera := camera.NewCamera(mgl32.Vec3{3.3, -10, 14.0}, mgl32.Vec3{0, 1, 0}, -101.0, 21.5)
	camera.SetupProjection(45, float32(WindowWidth)/float32(WindowHeight), 0.1, 1000.0)
	return camera
}
func Update() {
	nowNano := time.Now().UnixNano()
	moveTime := float64(nowNano-lastUpdate) / float64(time.Millisecond)
	lastUpdate = nowNano

	forward := 0.0
	if app.GetKeyState(FORWARD) && !app.GetKeyState(BACKWARD) {
		forward = moveSpeed * moveTime
	} else if app.GetKeyState(BACKWARD) && !app.GetKeyState(FORWARD) {
		forward = -moveSpeed * moveTime
	}
	if forward != 0 {
		app.GetCamera().Walk(float32(forward))
	}
	horisontal := 0.0
	if app.GetKeyState(LEFT) && !app.GetKeyState(RIGHT) {
		horisontal = -moveSpeed * moveTime
	} else if app.GetKeyState(RIGHT) && !app.GetKeyState(LEFT) {
		horisontal = moveSpeed * moveTime
	}
	if horisontal != 0 {
		app.GetCamera().Strafe(float32(horisontal))
	}
	vertical := 0.0
	if app.GetKeyState(UP) && !app.GetKeyState(DOWN) {
		vertical = -moveSpeed * moveTime
	} else if app.GetKeyState(DOWN) && !app.GetKeyState(UP) {
		vertical = moveSpeed * moveTime
	}
	if vertical != 0 {
		app.GetCamera().Lift(float32(vertical))
	}
	currX, currY := app.GetWindow().GetCursorPos()
	x, y := trans.MouseCoordinates(currX, currY, WindowWidth, WindowHeight)
	KeyDowns := make(map[string]bool)
	// dUp
	if y > 1.0-cameraDistance && y < 1.0 {
		KeyDowns["dUp"] = true
	} else {
		KeyDowns["dUp"] = false
	}
	// dDown
	if y < -1.0+cameraDistance && y > -1.0 {
		KeyDowns["dDown"] = true
	} else {
		KeyDowns["dDown"] = false
	}
	// dLeft
	if x < -1.0+cameraDistance && x > -1.0 {
		KeyDowns["dLeft"] = true
	} else {
		KeyDowns["dLeft"] = false
	}
	// dRight
	if x > 1.0-cameraDistance && x < 1.0 {
		KeyDowns["dRight"] = true
	} else {
		KeyDowns["dRight"] = false
	}

	dX := float32(0.0)
	dY := float32(0.0)
	if KeyDowns["dUp"] && !KeyDowns["dDown"] {
		dY = cameraDirectionSpeed
	} else if KeyDowns["dDown"] && !KeyDowns["dUp"] {
		dY = -cameraDirectionSpeed
	}
	if KeyDowns["dLeft"] && !KeyDowns["dRight"] {
		dX = -cameraDirectionSpeed
	} else if KeyDowns["dRight"] && !KeyDowns["dLeft"] {
		dX = cameraDirectionSpeed
	}
	app.GetCamera().UpdateDirection(dX, dY)
}

func main() {
	runtime.LockOSThread()
	app = application.New()
	app.SetWindow(window.InitGlfw(WindowWidth, WindowHeight, WindowTitle))
	defer glfw.Terminate()
	shader.InitOpenGL()

	app.SetCamera(CreateCamera())

	// directional light is coming from the up direction but not from too up.
	DirectionalLightSource := light.NewDirectionalLight([4]mgl32.Vec3{
		DirectionalLightDirection,
		DirectionalLightAmbient,
		DirectionalLightDiffuse,
		DirectionalLightSpecular,
	})
	PointLightSource_1 := light.NewPointLight([4]mgl32.Vec3{
		PointLightDirection_1,
		PointLightAmbient,
		PointLightDiffuse,
		PointLightSpecular},
		[3]float32{LightConstantTerm, LightLinearTerm, LightQuadraticTerm})
	PointLightSource_2 := light.NewPointLight([4]mgl32.Vec3{
		PointLightDirection_2,
		PointLightAmbient,
		PointLightDiffuse,
		PointLightSpecular},
		[3]float32{LightConstantTerm, LightLinearTerm, LightQuadraticTerm})
	SpotLightSource_1 := light.NewSpotLight([5]mgl32.Vec3{
		SpotLightDirection_1,
		SpotLightAmbient,
		SpotLightDiffuse,
		SpotLightSpecular},
		[5]float32{LightConstantTerm, LightLinearTerm, LightQuadraticTerm, SpotLightCutoff, SpotLightOuterCutoff})
	SpotLightSource_2 := light.NewSpotLight([5]mgl32.Vec3{
		SpotLightDirection_2,
		SpotLightAmbient,
		SpotLightDiffuse,
		SpotLightSpecular},
		[5]float32{LightConstantTerm, LightLinearTerm, LightQuadraticTerm, SpotLightCutoff, SpotLightOuterCutoff})
	//Define the shader application
	shaderProgramGrass := shader.NewShader("examples/08-multiple-light/shaders/texture.vert", "examples/08-multiple-light/shaders/texture.frag")
	// Add textures
	shaderProgramGrass.AddTexture("examples/08-multiple-light/assets/grass.jpg", gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE, gl.LINEAR, gl.LINEAR, "material.diffuse")
	shaderProgramGrass.AddTexture("examples/08-multiple-light/assets/grass.jpg", gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE, gl.LINEAR, gl.LINEAR, "material.specular")
	// Add light sources
	shaderProgramGrass.AddDirectionalLightSource(DirectionalLightSource, [4]string{"dirLight[0].direction", "dirLight[0].ambient", "dirLight[0].diffuse", "dirLight[0].specular"})
	shaderProgramGrass.AddPointLightSource(PointLightSource_1, [7]string{"pointLight[0].position", "pointLight[0].ambient", "pointLight[0].diffuse", "pointLight[0].specular", "pointLight[0].constant", "pointLight[0].linear", "pointLight[0].quadratic"})
	shaderProgramGrass.AddPointLightSource(PointLightSource_2, [7]string{"pointLight[1].position", "pointLight[1].ambient", "pointLight[1].diffuse", "pointLight[1].specular", "pointLight[1].constant", "pointLight[1].linear", "pointLight[1].quadratic"})
	shaderProgramGrass.AddSpotLightSource(SpotLightSource_1, [10]string{"spotLight[0].position", "spotLight[0].direction", "spotLight[0].ambient", "spotLight[0].diffuse", "spotLight[0].specular", "spotLight[0].constant", "spotLight[0].linear", "spotLight[0].quadratic", "spotLight[0].cutOff", "spotLight[0].outerCutOff"})
	shaderProgramGrass.AddSpotLightSource(SpotLightSource_2, [10]string{"spotLight[1].position", "spotLight[1].direction", "spotLight[1].ambient", "spotLight[1].diffuse", "spotLight[1].specular", "spotLight[1].constant", "spotLight[1].linear", "spotLight[1].quadratic", "spotLight[1].cutOff", "spotLight[1].outerCutOff"})
	shaderProgramGrass.SetViewPosition(app.GetCamera().GetPosition(), "viewPosition")
	Grass(shaderProgramGrass)
	shaderProgramBox := shader.NewShader("examples/08-multiple-light/shaders/texture.vert", "examples/08-multiple-light/shaders/texture.frag")
	shaderProgramBox.AddTexture("examples/08-multiple-light/assets/box-diffuse.png", gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE, gl.LINEAR, gl.LINEAR, "material.diffuse")
	shaderProgramBox.AddTexture("examples/08-multiple-light/assets/box-specular.png", gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE, gl.LINEAR, gl.LINEAR, "material.specular")
	shaderProgramBox.AddDirectionalLightSource(DirectionalLightSource, [4]string{"dirLight[0].direction", "dirLight[0].ambient", "dirLight[0].diffuse", "dirLight[0].specular"})
	shaderProgramBox.AddPointLightSource(PointLightSource_1, [7]string{"pointLight[0].position", "pointLight[0].ambient", "pointLight[0].diffuse", "pointLight[0].specular", "pointLight[0].constant", "pointLight[0].linear", "pointLight[0].quadratic"})
	shaderProgramBox.AddPointLightSource(PointLightSource_2, [7]string{"pointLight[1].position", "pointLight[1].ambient", "pointLight[1].diffuse", "pointLight[1].specular", "pointLight[1].constant", "pointLight[1].linear", "pointLight[1].quadratic"})
	shaderProgramBox.AddSpotLightSource(SpotLightSource_1, [10]string{"spotLight[0].position", "spotLight[0].direction", "spotLight[0].ambient", "spotLight[0].diffuse", "spotLight[0].specular", "spotLight[0].constant", "spotLight[0].linear", "spotLight[0].quadratic", "spotLight[0].cutOff", "spotLight[0].outerCutOff"})
	shaderProgramBox.AddSpotLightSource(SpotLightSource_2, [10]string{"spotLight[1].position", "spotLight[1].direction", "spotLight[1].ambient", "spotLight[1].diffuse", "spotLight[1].specular", "spotLight[1].constant", "spotLight[1].linear", "spotLight[1].quadratic", "spotLight[1].cutOff", "spotLight[1].outerCutOff"})
	shaderProgramBox.SetViewPosition(app.GetCamera().GetPosition(), "viewPosition")
	Box(shaderProgramBox)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.3, 0.1, 0.5, 1.0)

	lastUpdate = time.Now().UnixNano()
	// register keyboard button callback
	app.GetWindow().SetKeyCallback(app.KeyCallback)

	for !app.GetWindow().ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		Update()
		app.DrawWithUniforms()
		glfw.PollEvents()
		app.GetWindow().SwapBuffers()
	}
}
