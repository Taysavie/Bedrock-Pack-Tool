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
	"github.com/klauspost/compress/flate"
	"github.com/swim-services/swim_porter/compressor"
)

type Compressor struct {
	picker *utils.FilePickerOrLink
	button *widget.Button
}

func (c *Compressor) View(w fyne.Window) fyne.CanvasObject {
	if c.button == nil {
		c.button = widget.NewButton("Compress", func() {
			if c.picker == nil {
				return
			}
			c.button.Disable()
			c.button.SetText("Compressing... (may take a while)")
			defer c.button.SetText("Compress")
			defer c.picker.Clear()
			reader, err := c.picker.Reader()
			if err != nil {
				dialog.NewError(err, w).Show()
				return
			}
			data, err := io.ReadAll(reader)
			if err != nil {
				dialog.NewError(err, w).Show()
				return
			}
			name := utils.TrimMultiSpace(utils.RemoveExtension(c.picker.Filename()))
			extracted, err := utils.LoadArchive(c.picker.Filename(), data)
			if err != nil {
				dialog.NewError(err, w).Show()
				return
			}
			compressor.CompressRaw(extracted)
			if err := utils.SaveMapFsAsZipCompressionLevel(extracted, name+"_compressed.mcpack", w, flate.BestCompression); err != nil {
				dialog.NewError(err, w).Show()
				return
			}
		})
		c.picker = utils.NewFilePickerOrMediafire(func(isReady bool) {
			if isReady {
				c.button.Enable()
			} else {
				c.button.Disable()
			}
		}, []string{".zip", ".rar", ".mcpack"})
	}
	c.button.Disable()
	text := canvas.NewText("Swim Pack Compressor", fyne.CurrentApp().Settings().Theme().Color(theme.ColorNameForeground, fyne.CurrentApp().Settings().ThemeVariant()))
	text.TextSize = 50
	return container.NewCenter(container.NewVBox(text, c.picker.Show(w), c.button))
}

func (c *Compressor) OnDrop(uri fyne.URI) {
	if c.picker != nil {
		c.picker.OnDrop(uri)
	}
}
