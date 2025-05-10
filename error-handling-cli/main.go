package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"error-handling-cli/cmd"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "goerrors",
		Short: "A CLI tool for learning Go error handling patterns",
		Long: `
Go Error Handling CLI Tutorial
-------------------------------
This CLI tool demonstrates various Go error handling patterns and best practices.
It provides interactive examples with step-by-step explanations of different error
handling techniques in Go.

Use the subcommands to explore different error handling patterns.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Welcome to the Go Error Handling Tutorial")
			fmt.Println("Use --help to see available commands")
		},
	}

	// Add subcommands
	cmd.AddBasicErrorHandlingCmd(rootCmd)
	cmd.AddCustomErrorsCmd(rootCmd)
	cmd.AddErrorWrappingCmd(rootCmd)
	cmd.AddPanicRecoveryCmd(rootCmd)
	cmd.AddContextErrorsCmd(rootCmd)
	cmd.AddErrorGroupsCmd(rootCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
