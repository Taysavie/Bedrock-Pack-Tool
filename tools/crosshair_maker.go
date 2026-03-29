package tools

import (
	"io"
	"swim-pack-tool/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/swim-services/swim_porter/crosshairmaker"
)

type CrosshairMaker struct {
	picker *utils.FilePicker
	button *widget.Button
}

func (c *CrosshairMaker) View(w fyne.Window) fyne.CanvasObject {
	size := utils.NewObjectWithLabel("Size:", widget.NewSelect([]string{"Very small", "Small", "Medium", "Large", "Very large"}, nil))
	size.Obj().SetSelectedIndex(2)

	colorCheck := widget.NewCheck("Color", nil)
	if c.button == nil {
		c.button = widget.NewButton("Make crosshair", func() {
			if c.picker == nil {
				return
			}
			c.button.Disable()
			c.button.SetText("Making crosshair...")
			defer c.button.SetText("Make crosshair")
			defer c.picker.Clear()
			data, err := io.ReadAll(c.picker.Reader())
			if err != nil {
				dialog.NewError(err, w).Show()
				return
			}
			name := utils.TrimMultiSpace(utils.RemoveExtension(c.picker.Filename()))
			img, err := utils.ReadImage(data)
			if err != nil {
				dialog.NewError(err, w).Show()
				return
			}

			zipped, err := crosshairmaker.CrosshairPack(name, img, float64(size.Obj().SelectedIndex()+1)*0.2, colorCheck.Checked)

			if err != nil {
				dialog.NewError(err, w).Show()
				return
			}
			utils.SaveFile(zipped, name+"_crosshair.mcpack", w)
		})
		c.picker = utils.NewFilePicker(func(isReady bool) {
			if isReady {
				c.button.Enable()
			} else {
				c.button.Disable()
			}
		}, []string{".png", ".webp"})
	}
	c.button.Disable()

	text := canvas.NewText("Swim Crosshair Maker", fyne.CurrentApp().Settings().Theme().Color(theme.ColorNameForeground, fyne.CurrentApp().Settings().ThemeVariant()))
	text.TextSize = 50
	return container.NewCenter(container.NewVBox(text, c.picker.Show(w), size.Container(), colorCheck, c.button))
}

func (c *CrosshairMaker) OnDrop(uri fyne.URI) {
	if c.picker != nil {
		c.picker.OnDrop(uri)
	}
}
