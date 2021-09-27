package main

import (
	"fmt"
	"github.com/fivemoreminix/dos"
	"github.com/gdamore/tcell/v2"
	"os"
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create tcell screen: %v", err)
	}
	if err = screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize: %v", err)
	}

	widget := dos.Center{
		Child: dos.Label{
			Text: "Hello, world!",
		},
	}

	app := dos.App{
		MainWidget: &widget,
	}
	app.Run(screen)
}
