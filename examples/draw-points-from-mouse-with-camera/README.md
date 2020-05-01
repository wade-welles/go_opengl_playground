# Mouse input example with keyboard color management, and camera

This application demonstrates the mouse handler options with some keyboard inputs. The middle mouse button selects a new point. Every selected point is displayed on the screen. The button `r` sets the red, the `g` the green the `b` the blue part of the color of the point. The camera position is updatable with the `q`, `w`, `e`, `a`, `s`, `d` keys. (The coordinates of the point is calculated with some magic, but i can't remember why.)

## Application

This stucture represents our world. It contains the selected points. Every point contains the coordinate ([-1,1]) and the color ([0-1]). The AddPoint function is calle from the mouse event handler.

### MouseHandler

This logic is responsible for handling the mouse events. It maintains 3 global variable. The `mousePositionX` is for storing the x coordinate of the mouse click, `mousePositionY` is for storing the y coordinate of the mouse click, `mouseButtonPressed` is for handling the button state. If the button is released, then the selected point (coordinates of the click event) is converted then inserted to the app.Points.

### KeyHandler

This is a basic function for supporting the debug. In case of the `h` button is clicked, it prints out the app.Points. It's other responsibility is the color management. in case of the button `r` is clicked, the red part of the color is updated to 1, else it fallbacks to 0. The same logic is implemented for the `g` button end green color and for the `b` button and blue color. The camera position is updatable with the `q`, `w`, `e`, `a`, `s`, `d` keys.

### buildVAO

This function builds the `[]float32`, that can be used as vertex data object. Currently it inserts the coordinates and the colors to the vao.

### Draw

This function is responsible for drawing the points to the screen. It creates the vao, sets the buffer data, enables and sets the attribute arrays and then draw the points.

## Functions

- `initGlfw`

Basic function for glfw initialization.

- `initOpenGL`

It is responsible for openGL initialization. It uses the `shader.FragmentShaderBasicSource` fragment shader and the `shader.VertexShaderPointWithColorModelViewProjectionSource` vertex shader.