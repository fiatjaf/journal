package main

import (
	"os"

	"github.com/therecipe/qt/widgets"
)

func main() {
	app := widgets.NewQApplication(len(os.Args), os.Args)

	widget := widgets.NewQWidget(nil, 0)
	widget.SetLayout(widgets.NewQVBoxLayout())

	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(800, 600)
	window.SetWindowTitle("journal")
	window.SetCentralWidget(widget)
	window.Show()

	app.Exec()
}
