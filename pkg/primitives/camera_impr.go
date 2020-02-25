package primitives

import (
	"math"
	"strconv"
)

type CameraImpr struct {
	// camera position
	position Vector
	// look at position (the camera is looking at lookat from it's position.)
	lookAt Vector
	// world up direction.
	UpDirection Vector

	projectionOptions struct {
		fov         float64
		aspectRatio float64

		far  float64
		near float64
	}
	// projection matrix
	projectionMatrix *Matrix4x4

	cameraUpDirection    Vector
	cameraRightDirection Vector
	cameraFrontDirection Vector
	// view matrix
	viewMatrix *Matrix4x4

	cameraOptions struct {
		yaw         float64
		pitch       float64
		roll        float64
		minRy       float64
		maxRy       float64
		distance    float64
		minDistance float64
		maxDistance float64
	}
}

// Log returns the string representation of this object.
func (c *CameraImpr) Log() string {
	logString := "Position: Vector{" + c.position.ToString() + "}\n"
	logString += "LookAt: Vector{" + c.lookAt.ToString() + "}\n"
	logString += "UpDirection: Vector{" + c.UpDirection.ToString() + "}\n"
	logString += "ProjectionOptions:\n"
	logString += " - fov : " + strconv.FormatFloat(c.projectionOptions.fov, 'f', 6, 64) + "\n"
	logString += " - aspectRatio : " + strconv.FormatFloat(c.projectionOptions.aspectRatio, 'f', 6, 64) + "\n"
	logString += " - far : " + strconv.FormatFloat(c.projectionOptions.far, 'f', 6, 64) + "\n"
	logString += " - near : " + strconv.FormatFloat(c.projectionOptions.near, 'f', 6, 64) + "\n"
	logString += "CameraUpDirection: Vector{" + c.cameraUpDirection.ToString() + "}\n"
	logString += "CameraRightDirection: Vector{" + c.cameraRightDirection.ToString() + "}\n"
	logString += "CameraForwardDirection: Vector{" + c.cameraFrontDirection.ToString() + "}\n"
	logString += "CameraOptions:\n"
	logString += " - yaw : " + strconv.FormatFloat(c.cameraOptions.yaw, 'f', 6, 64) + "\n"
	logString += " - pitch : " + strconv.FormatFloat(c.cameraOptions.pitch, 'f', 6, 64) + "\n"
	logString += " - roll : " + strconv.FormatFloat(c.cameraOptions.roll, 'f', 6, 64) + "\n"
	logString += " - minRy : " + strconv.FormatFloat(c.cameraOptions.minRy, 'f', 6, 64) + "\n"
	logString += " - maxRy : " + strconv.FormatFloat(c.cameraOptions.maxRy, 'f', 6, 64) + "\n"
	logString += " - distance : " + strconv.FormatFloat(c.cameraOptions.distance, 'f', 6, 64) + "\n"
	logString += " - minDistance : " + strconv.FormatFloat(c.cameraOptions.minDistance, 'f', 6, 64) + "\n"
	logString += " - maxDistance : " + strconv.FormatFloat(c.cameraOptions.maxDistance, 'f', 6, 64) + "\n"
	return logString
}

func NewCameraImpr() *CameraImpr {
	var cam CameraImpr
	cam.cameraOptions.minRy = -60
	cam.cameraOptions.maxRy = 60
	cam.cameraOptions.minDistance = 0.1
	cam.cameraOptions.maxDistance = 1000
	cam.projectionOptions.fov = 45
	cam.projectionOptions.aspectRatio = 1
	cam.projectionOptions.near = 0.1
	cam.projectionOptions.far = 1000
	return &cam
}

// GetViewMatrix returns the viewMatrix of the camera
func (c *CameraImpr) GetViewMatrix() *Matrix4x4 {
	return c.viewMatrix
}

// GetProjectionMatrix returns the viewMatrix of the camera
func (c *CameraImpr) GetProjectionMatrix() *Matrix4x4 {
	return c.projectionMatrix
}

// SetPosition updates the position of the camera
func (c *CameraImpr) SetPosition(newPos Vector) {
	c.position = newPos
}

// SetupProjection creates the projection matrix and setups the fow and the aspect_ration
func (c *CameraImpr) SetupProjection(fov, aspRatio float64) {
	c.projectionOptions.fov = fov
	c.projectionOptions.aspectRatio = aspRatio
	c.projectionMatrix = Perspective(float32(c.projectionOptions.fov), float32(c.projectionOptions.aspectRatio), float32(c.projectionOptions.near), float32(c.projectionOptions.far))
}

// setupCameraDirections is a helper function for updating the camera[Up|Right|Forward]Diraction variables
func (c *CameraImpr) setupCameraDirections() {
	c.cameraFrontDirection = (c.position.Add(c.lookAt.MultiplyScalar(-1))).Normalize()

	c.cameraRightDirection = (c.UpDirection.Cross(c.cameraFrontDirection)).Normalize()
	c.cameraUpDirection = c.cameraFrontDirection.Cross(c.cameraRightDirection)
}

// TargetCameraSetTarget updates the camera based on the new targetPoint
// It updates the lookAt vector (target - position normalized)
// It calculates the distance
// Sets the viewMatrix
// Sets the yaw, pitch values
func (c *CameraImpr) TargetCameraSetTarget(target Vector) {
	c.lookAt = (target.Subtract(c.position)).Normalize()
	c.cameraOptions.distance = (c.position.Subtract(target)).Length()
	c.cameraOptions.distance = math.Max(
		c.cameraOptions.minDistance, math.Min(
			c.cameraOptions.distance, c.cameraOptions.maxDistance))
	c.setupCameraDirections()
	c.viewMatrix = LookAt(c.position, c.lookAt, c.cameraUpDirection)
	c.cameraOptions.yaw = 0
	c.cameraOptions.pitch = 0
	if c.viewMatrix.Points[0] < 0 {
		c.cameraOptions.yaw = RadToDeg(math.Pi - math.Asin(-float64(c.viewMatrix.Points[8])))
	} else {
		c.cameraOptions.yaw = RadToDeg(math.Asin(-float64(c.viewMatrix.Points[8])))
	}
	c.cameraOptions.pitch = RadToDeg(math.Asin(-float64(c.viewMatrix.Points[6])))
}

// Update setup the viewMatrix based on the given cameraOptions
func (c *CameraImpr) Update() {
	rotation := YPR(
		DegToRad(c.cameraOptions.yaw),
		DegToRad(c.cameraOptions.pitch),
		DegToRad(c.cameraOptions.roll))
	T := rotation.MultiVector(Vector{0, 0, c.cameraOptions.distance})
	c.position = c.lookAt.Add(*T)
	c.lookAt = (c.lookAt.Subtract(c.position)).Normalize()
	c.cameraUpDirection = *(rotation.MultiVector(c.UpDirection))
	c.cameraRightDirection = c.lookAt.Cross(c.cameraUpDirection)
	c.viewMatrix = LookAt(c.position, c.lookAt, c.cameraUpDirection)
}

// SetViewMatrix setup the viewMatrix based on the given cameraOptions
func (c *CameraImpr) SetViewMatrix() {
	c.setupCameraDirections()
	c.viewMatrix = LookAt(c.position, c.lookAt, c.cameraUpDirection)
}
func (c *CameraImpr) TargetCameraMove(dX, dY float64) {
	x := c.cameraRightDirection.MultiplyScalar(dX)
	y := c.lookAt.MultiplyScalar(dY)
	c.position = c.position.Add(x.Add(y))
	c.lookAt = c.lookAt.Add(x.Add(y))
	c.Update()
}

// Walk updates the translation for the transformation (forward, back directions)
func (c *CameraImpr) Walk(amount float64) {
	c.position = c.position.Add(
		c.lookAt.MultiplyScalar(amount))
	c.SetViewMatrix()
}

// Strafe updates the translation for the transformation (left, right directions)
func (c *CameraImpr) Strafe(amount float64) {
	c.position = c.position.Add(
		c.cameraRightDirection.MultiplyScalar(amount))
	c.SetViewMatrix()
}
