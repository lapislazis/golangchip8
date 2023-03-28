/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package main
import ( 
	"alex/chip8/cmd"
	"github.com/faiface/pixel/pixelgl"
)

func main() {
	pixelgl.Run(runCHIP8)
}

func runCHIP8() {
	cmd.Execute()
}