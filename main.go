// Main program of the go-n-queens-problem
package main

import (
	"flag"
	"runtime"

	"github.com/go-gl/glfw/v3.2/glfw"
)

// N flag: the number of queens and the size of the board
var N = flag.Int("N", 4, "the number of queens and the size of the board")

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()

	program := initOpenGL()

	for !window.ShouldClose() {
		// TODO
		draw(window, program)
	}
}
