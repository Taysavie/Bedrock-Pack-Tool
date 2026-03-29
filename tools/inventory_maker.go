package tools

import (
	"bytes"
	"image/gif"
	"io"
	"path"
	"swim-pack-tool/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/swim-services/swim_porter/animatedinv"
)

type InventoryMaker struct {
	picker *utils.FilePickerOrLink
	button *widget.Button
}

func (i *InventoryMaker) View(w fyne.Window) fyne.CanvasObject {
	addSlotsCheck := widget.NewCheck("Add slots", nil)
	if i.button == nil {
		i.button = widget.NewButton("Make inventory", func() {
			if i.picker == nil {
				return
			}
			i.button.Disable()
			i.button.SetText("Making inventory...")
			defer i.button.SetText("Make inventory")
			defer i.picker.Clear()
			reader, err := i.picker.Reader()
			if err != nil {
				dialog.NewError(err, w).Show()
				return
			}
			data, err := io.ReadAll(reader)
			if err != nil {
				dialog.NewError(err, w).Show()
				return
			}
			name := utils.TrimMultiSpace(utils.RemoveExtension(i.picker.Filename()))
			ext := path.Ext(i.picker.Filename())
			var zipped []byte
			if ext == ".gif" {
				img, err := gif.DecodeAll(bytes.NewReader(data))
				if err != nil {
					dialog.NewError(err, w).Show()
					return
				}
				zipped, err = animatedinv.MakeAnimated(img, name, addSlotsCheck.Checked)
				if err != nil {
					dialog.NewError(err, w).Show()
					return
				}
			} else {
				img, err := utils.ReadImage(data)
				if err != nil {
					dialog.NewError(err, w).Show()
					return
				}
				zipped, err = animatedinv.MakeOverlay(img, name, addSlotsCheck.Checked)
				if err != nil {
					dialog.NewError(err, w).Show()
					return
				}
			}
			utils.SaveFile(zipped, name+"_inventory.mcpack", w)
		})
		i.button.Disable()
		i.picker = utils.NewFilePickerOrGifLink(func(isReady bool) {
			if isReady {
				i.button.Enable()
			} else {
				i.button.Disable()
			}
		}, []string{".jpg", ".jpeg", ".png", ".webp", ".gif"})
	}

	text := canvas.NewText("Swim Inventory Maker", fyne.CurrentApp().Settings().Theme().Color(theme.ColorNameForeground, fyne.CurrentApp().Settings().ThemeVariant()))
	text.TextSize = 50
	return container.NewCenter(container.NewVBox(text, i.picker.Show(w), addSlotsCheck, i.button))
}

func (i *InventoryMaker) OnDrop(uri fyne.URI) {
	if i.picker != nil {
		i.picker.OnDrop(uri)
	}
}
