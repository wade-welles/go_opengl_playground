# Rectangle

It represents a rectangle, so it contains 4 `coordinate vectors`, 4 `color vectors`, a direction vector + speed (for moving objects).
It has a VAO and a Shader also. It implements the Drawable interface.

## Functions

The stuff that we can do with a rectangle.

## New

It creates a new rectangle. The inputs of this functions are the coordinates, the colors and the shader. The direction & speed is initialized as null vector & 0 speed.

### SetColor

It updates the colors for the given one. It has one input, the new color Vector.

### SetIndexColor

It updates the color with the new color for the given index. It has 2 inputs, the index and the new color Vector.

### SetDirection

It updates the direction of the rectangle with the given new one.

### SetIndexDirection

It updates the indexed direction of the rectangle to the given value.

### SetSpeed

It updates the speed of the rectangle.

### SetAngle

It updates the rotation angle (radian) of the rectangle.

### SetAxis

It updates the rotation axis of the rectangle.

### Log

The string representation of the current state of the object.

### Draw

It draws the rectangle. The MVP uniform matrix is set to ident. matrix for the shader program.

### DrawWithUniforms

It draws the rectangle. It gets the V & P matrices as inputs. It sets the model, view, projection uniforms for the shader program.

### Update

It updates the state of the rectangle. It gets the delta time as input and it calculates the movement of the rectangle.

### Coordinates

It returns the coordinates of the rectangle.

### Color

It returns the colors of the rectangle.
