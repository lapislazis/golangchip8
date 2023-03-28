/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)



// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "chip8",
	Short: "A CHIP-8 emulator",
	Long: `This is a CHIP-8 emulator made as part of a college Computer Science project.`,
	Args: cobra.ExactArgs(1),
	Run: runRoot,
}

//Runs if an unknown command is used
func runRoot(cmd *cobra.Command, args []string) {
	fmt.Println("Unknown command. Try `chip8 help` for more information")
}

var clockSpeed int

func init() {
	rootCmd.AddCommand(runCmd)

	//Defines an optional flag to set the clock speed
	runCmd.Flags().IntVarP(&clockSpeed, "clockspeed", "c", 700, "Set clock speed")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}



