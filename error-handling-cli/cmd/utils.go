package cmd

import (
        "fmt"
        "strings"

        "github.com/fatih/color"
)

// Helper functions for the tutorial UI

// ClearScreen clears the terminal screen
func ClearScreen() {
        fmt.Print("\033[H\033[2J")
}

// PrintTitle prints a formatted title
func PrintTitle(title string) {
        color.New(color.FgHiBlue, color.Bold).Println("\n" + title)
        color.New(color.FgHiBlue, color.Bold).Println(strings.Repeat("=", len(title)))
        fmt.Println()
}

// PrintSection prints a formatted section heading
func PrintSection(title string) {
        color.New(color.FgYellow, color.Bold).Println("\n" + title)
        color.New(color.FgYellow, color.Bold).Println(strings.Repeat("-", len(title)))
}

// PressEnterToContinue waits for the user to press Enter
func PressEnterToContinue() {
        fmt.Print("\nPress Enter to continue...")
        fmt.Scanln()
}
