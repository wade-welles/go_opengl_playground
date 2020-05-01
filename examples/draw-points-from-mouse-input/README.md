# Mouse input example

This application demonstrates the mouse handler options. The middle mouse button selects a new point. Every selected point is displayed on the screen.

## Application

This stucture represents our world. It contains the selected points. Every point contains the coordinate ([-1,1]) and the color ([0-1]). The AddPoint function is calle from the mouse event handler.

### NewApplication

It returns an application with the minimal setup.

### MouseHandler

This logic is responsible for handling the mouse events. It maintains 3 global variable. The `mousePositionX` is for storing the x coordinate of the mouse click, `mousePositionY` is for storing the y coordinate of the mouse click, `mouseButtonPressed` is for handling the button state. If the button is released, then the selected point (coordinates of the click event) is converted then inserted to the app.Points.

### KeyHandler

This is a basic function for supporting the debug. In case of the `D` button is clicked, it prints out the app.Points.

### buildVAO

This function builds the `[]float32`, that can be used as vertex data object. Currently it only inserts the coordinates to the vao.

### Draw

This function is responsible for drawing the points to the screen. It creates the vao, sets the buffer data, enables and sets the attribute arrays and then draw the points.

## Functions

- `initGlfw`

Basic function for glfw initialization.

- `initOpenGL`

It is responsible for openGL initialization. It uses the `shader.FragmentShaderConstantSource` fragment shader and the `shader.VertexShaderPointSource` vertex shader.