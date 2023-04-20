/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	chip8 "alex/chip8/emulator"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run 'path/to/rom'",
	Short: "Run the chip8 emulator",
	Run: runChip8,
	}

func runChip8(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("The run command takes one argument: a `path/to/rom`")
		os.Exit(1)
	}
	filePath := os.Args[2]

	//Starts new vm
	vm, err := chip8.NewVM(filePath, clockSpeed, debug)
	if err != nil {
		fmt.Printf("\nError creating a new CHIP-8 VM: %v\n", err)
		os.Exit(1)
	}

	go vm.Audio()
	go vm.Run()

	<-vm.ShutdownChan
}
