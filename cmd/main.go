package main

import (
	"fmt"
	"os"

	"github.com/Dowdow/gosible/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Usage: gosible <config.json>")
		return
	}

	p := tea.NewProgram(ui.NewMainModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprint(os.Stderr, ui.PrintError(err))
		os.Exit(1)
	}
}
