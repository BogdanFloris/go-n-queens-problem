package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	width  = 500 // width of the game board
	height = 500 // height of the game board

	// source code for the vertex shader
	vertexShaderSource = `
    #version 410
	layout(location = 0) in vec3 vp;
	layout(location = 1) in vec3 vc;
	out vec3 fragmentColor;
    void main() {
		gl_Position = vec4(vp, 1.0);
		fragmentColor = vc;
    }
` + "\x00"

	// source code for the fragment shader
	fragmentShaderSource = `
	#version 410
	precision highp float;
	in vec3 fragmentColor;
    out vec3 frag_colour;
    void main() {
        frag_colour = fragmentColor;
    }
` + "\x00"
)

// struct that represents a cell
type cell struct {
	// drawable openGL VAO
	drawable uint32

	// whether this cell has a queen on it or not
	hasQueen bool

	// coordinates
	x int
	y int
}

// draws the cell
func (c *cell) draw() {
	gl.BindVertexArray(c.drawable)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square)/3))
}

// function that makes the cells of the board
func makeCells() [][]*cell {
	cells := make([][]*cell, *N, *N)
	// first cell in the bottom left corner is black
	isWhite := false
	for x := 0; x < *N; x++ {
		for y := 0; y < *N; y++ {
			// make a new cell
			c := newCell(x, y, isWhite)
			cells[x] = append(cells[x], c)
			// next cell is white
			isWhite = !isWhite
		}
		// only switch the color again if the board has an even N
		if *N%2 == 0 {
			isWhite = !isWhite
		}
	}
	return cells
}

func newCell(x, y int, isWhite bool) *cell {
	// copy the square in a new array such that we can modify the points
	points := make([]float32, len(square), len(square))
	copy(points, square)

	for i := 0; i < len(points); i++ {
		var position float32
		var size float32
		switch i % 3 {
		// if x coordinate
		case 0:
			size = 1.0 / float32(*N)
			position = float32(x) * size
		// if y coordinate
		case 1:
			size = 1.0 / float32(*N)
			position = float32(y) * size
		default:
			continue
		}

		if points[i] < 0 {
			points[i] = (position * 2) - 1
		} else {
			points[i] = ((position + size) * 2) - 1
		}
	}
	// if the cell is white, use the white color buffer
	if isWhite {
		return &cell{
			// TODO: replace squareColor with new colors
			drawable: makeVao(points, squareColorWhite),
			x:        x,
			y:        y,
		}
	}
	// else use the black color buffer
	return &cell{
		// TODO: replace squareColor with new colors
		drawable: makeVao(points, squareColorBlack),
		x:        x,
		y:        y,
	}
}

func draw(cells [][]*cell, window *glfw.Window, program uint32) {
	// clear everything
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	// grey background RGB(0.25, 0.25, 0.25)
	gl.ClearColor(0.25, 0.25, 0.25, 1)
	gl.UseProgram(program)

	// draw each cell
	for x := range cells {
		for _, c := range cells[x] {
			c.draw()
		}
	}

	// poll events and swap the buffers
	glfw.PollEvents()
	window.SwapBuffers()
}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	// window settings
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "N-Queens Problem", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	// compile the vertex shader
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	// compile the fragment shader
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	// create the program and attach the shaders
	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

// function that makes a VAO based on an array of points and an array of colors
func makeVao(points []float32, colors []float32) uint32 {
	// buffer for the vertices
	var vertexBuffer uint32
	gl.GenBuffers(1, &vertexBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	// buffer for the colors
	var colorBuffer uint32
	gl.GenBuffers(1, &colorBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, colorBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(colors), gl.Ptr(colors), gl.STATIC_DRAW)

	// initialize the VAO
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// bind the vertex buffer
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	// bind the color buffer
	gl.EnableVertexAttribArray(1)
	gl.BindBuffer(gl.ARRAY_BUFFER, colorBuffer)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 0, nil)

	return vao
}

// function that compiles the source code of a shader based on the given shader type
func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
