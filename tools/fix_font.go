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
	"github.com/swim-services/swim_porter/fontfix"
)

type FixFont struct {
	picker *utils.FilePickerOrLink
	button *widget.Button
}

func (f *FixFont) View(w fyne.Window) fyne.CanvasObject {
	if f.button == nil {
		f.button = widget.NewButton("Fix font", func() {
			if f.picker == nil {
				return
			}
			f.button.Disable()
			f.button.SetText("Fixing font...")
			defer f.button.SetText("Fix font")
			defer f.picker.Clear()
			reader, err := f.picker.Reader()
			if err != nil {
				dialog.NewError(err, w).Show()
				return
			}
			data, err := io.ReadAll(reader)
			if err != nil {
				dialog.NewError(err, w).Show()
				return
			}
			name := utils.TrimMultiSpace(utils.RemoveExtension(f.picker.Filename()))
			extracted, err := utils.LoadArchive(f.picker.Filename(), data)
			if err != nil {
				dialog.NewError(err, w).Show()
				return
			}
			err = fontfix.FixFontRaw(extracted)
			if err != nil {
				dialog.NewError(err, w).Show()
				return
			}
			if err := utils.SaveMapFsAsZip(extracted, name+"_FONTFIX.mcpack", w); err != nil {
				dialog.NewError(err, w).Show()
				return
			}
		})
		f.picker = utils.NewFilePickerOrMediafire(func(isReady bool) {
			if isReady {
				f.button.Enable()
			} else {
				f.button.Disable()
			}
		}, []string{".zip", ".rar", ".mcpack"})
	}
	f.button.Disable()

	text := canvas.NewText("Swim Font Fixer", fyne.CurrentApp().Settings().Theme().Color(theme.ColorNameForeground, fyne.CurrentApp().Settings().ThemeVariant()))
	text.TextSize = 50
	return container.NewCenter(container.NewVBox(text, f.picker.Show(w), f.button))
}

func (f *FixFont) OnDrop(uri fyne.URI) {
	if f.picker != nil {
		f.picker.OnDrop(uri)
	}
}
