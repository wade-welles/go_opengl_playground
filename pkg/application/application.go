package application

import (
	"fmt"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	DEBUG = glfw.KeyH
)

type Drawable interface {
	Draw()
	DrawWithUniforms(mgl32.Mat4, mgl32.Mat4)
	Update(float64)
	Log() string
}

type Camera interface {
	Log() string
	GetViewMatrix() mgl32.Mat4
	GetProjectionMatrix() mgl32.Mat4
	Walk(float32)
	Strafe(float32)
	Lift(float32)
	UpdateDirection(float32, float32)
	GetPosition() mgl32.Vec3
}

type Application struct {
	window     Window
	camera     Camera
	cameraSet  bool
	keyDowns   map[glfw.Key]bool
	mouseDowns map[glfw.MouseButton]bool
	MousePosX  float64
	MousePosY  float64

	items []Drawable
}

type Window interface {
	GetCursorPos() (float64, float64)
	SetKeyCallback(glfw.KeyCallback) glfw.KeyCallback
	SetMouseButtonCallback(glfw.MouseButtonCallback) glfw.MouseButtonCallback
	ShouldClose() bool
	SwapBuffers()
}

// New returns an application instance
func New() *Application {
	return &Application{
		keyDowns:   make(map[glfw.Key]bool),
		mouseDowns: make(map[glfw.MouseButton]bool),
		items:      []Drawable{},
		cameraSet:  false,
	}
}

// Log returns the string representation of this object.
func (a *Application) Log() string {
	logString := "Application:\n"
	if a.cameraSet {
		logString += " - camera : " + a.camera.Log() + "\n"
	}
	logString += " - items :\n"
	for _, item := range a.items {
		logString += item.Log()
	}
	return logString
}

// SetWindow updates the window with the new one.
func (a *Application) SetWindow(w Window) {
	a.window = w
}

// GetWindow returns the current window of the application.
func (a *Application) GetWindow() Window {
	return a.window
}

// SetCamera updates the camera with the new one.
func (a *Application) SetCamera(c Camera) {
	a.cameraSet = true
	a.camera = c
}

// GetCamera returns the current camera of the application.
func (a *Application) GetCamera() Camera {
	return a.camera
}

// SetMouseButtons updates the mouseDowns with the new one.
func (a *Application) SetMouseButtons(m map[glfw.MouseButton]bool) {
	a.mouseDowns = m
}

// GetMouseButtons returns the current mouseDowns of the application.
func (a *Application) GetMouseButtons() map[glfw.MouseButton]bool {
	return a.mouseDowns
}

// SetKeys updates the keyDowns with the new one.
func (a *Application) SetKeys(m map[glfw.Key]bool) {
	a.keyDowns = m
}

// GetKeys returns the current keyDowns of the application.
func (a *Application) GetKeys() map[glfw.Key]bool {
	return a.keyDowns
}

// AddItem inserts a new drawable item
func (a *Application) AddItem(d Drawable) {
	a.items = append(a.items, d)
}

// Draw calls Draw function in every drawable item.
func (a *Application) Draw() {
	for index, _ := range a.items {
		a.items[index].Draw()
	}
}

// Update calls the Update function in every drawable item.
func (a *Application) Update(dt float64) {
	for index, _ := range a.items {
		a.items[index].Update(dt)
	}
}

// DrawWithUniforms calls DrawWithUniforms function in every drawable item with the calculated V & P.
func (a *Application) DrawWithUniforms() {
	V := mgl32.Ident4()
	P := mgl32.Ident4()
	if a.cameraSet {
		V = a.camera.GetViewMatrix()
		P = a.camera.GetProjectionMatrix()
	}

	for _, item := range a.items {
		item.DrawWithUniforms(V, P)
	}
}

// KeyCallback is responsible for the keyboard event handling.
func (a *Application) KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	switch key {
	case DEBUG:
		if action != glfw.Release {
			fmt.Printf("%s\n", a.Log())
		}
		break
	default:
		a.SetKeyState(key, action)
		break
	}
}

// MouseButtonCallback is responsible for the mouse button event handling.
func (a *Application) MouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	a.MousePosX, a.MousePosY = w.GetCursorPos()
	switch button {
	default:
		a.SetButtonState(button, action)
		break
	}
}

// SetKeyState setups the keyDowns based on the key and action
func (a *Application) SetKeyState(key glfw.Key, action glfw.Action) {
	var isButtonPressed bool
	if action != glfw.Release {
		isButtonPressed = true
	} else {
		isButtonPressed = false
	}
	a.keyDowns[key] = isButtonPressed
}

// SetKeyState setups the keyDowns based on the key and action
func (a *Application) SetButtonState(button glfw.MouseButton, action glfw.Action) {
	var isButtonPressed bool
	if action != glfw.Release {
		isButtonPressed = true
	} else {
		isButtonPressed = false
	}
	a.mouseDowns[button] = isButtonPressed
}

// GetMouseButtonState returns the state of the given button
func (a *Application) GetMouseButtonState(button glfw.MouseButton) bool {
	return a.mouseDowns[button]
}

// GetKeyState returns the state of the given key
func (a *Application) GetKeyState(key glfw.Key) bool {
	return a.keyDowns[key]
}
