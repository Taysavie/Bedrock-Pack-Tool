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
	"github.com/swim-services/swim_porter/skyfix"
)

type FixSky struct {
	picker *utils.FilePickerOrLink
	button *widget.Button
}

func (f *FixSky) View(w fyne.Window) fyne.CanvasObject {
	if f.button == nil {
		f.button = widget.NewButton("Fix sky", func() {
			if f.picker == nil {
				return
			}
			f.button.Disable()
			f.button.SetText("Fixing sky...")
			defer f.button.SetText("Fix sky")
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
			err = skyfix.FixSkyRaw(extracted)
			if err != nil {
				dialog.NewError(err, w).Show()
				return
			}
			if err := utils.SaveMapFsAsZip(extracted, name+"_SKYFIX.mcpack", w); err != nil {
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
	text := canvas.NewText("Swim Sky Fixer", fyne.CurrentApp().Settings().Theme().Color(theme.ColorNameForeground, fyne.CurrentApp().Settings().ThemeVariant()))
	text.TextSize = 50
	return container.NewCenter(container.NewVBox(text, f.picker.Show(w), f.button))
}

func (f *FixSky) OnDrop(uri fyne.URI) {
	if f.picker != nil {
		f.picker.OnDrop(uri)
	}
}
