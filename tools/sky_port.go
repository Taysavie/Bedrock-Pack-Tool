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
	"github.com/swim-services/swim_porter/cubemap"
)

type SkyPorter struct {
	picker *utils.FilePicker
	button *widget.Button
}

func (f *SkyPorter) View(w fyne.Window) fyne.CanvasObject {
	shiftCubemapCheck := widget.NewCheck("Disable skybox offset", nil)

	if f.button == nil {
		f.button = widget.NewButton("Port sky", func() {
			if f.picker == nil {
				return
			}
			f.button.Disable()
			f.button.SetText("Porting sky...")
			defer f.button.SetText("Port sky")
			defer f.picker.Clear()
			data, err := io.ReadAll(f.picker.Reader())
			if err != nil {
				dialog.NewError(err, w).Show()
				return
			}
			name := utils.TrimMultiSpace(utils.RemoveExtension(f.picker.Filename()))
			img, err := utils.ReadImage(data)
			if err != nil {
				dialog.NewError(err, w).Show()
				return
			}

			cubeMap := cubemap.BuildCubemap(img)
			if !shiftCubemapCheck.Checked {
				totalWidth := 0
				for _, img := range cubeMap {
					totalWidth += img.Bounds().Dx()
				}
				multAmt := max(4.5, min(8, float64(totalWidth)/1024))
				equi := cubemap.CubemapToEquirectangular(cubeMap, multAmt)
				cubeMap = cubemap.CubemapFromImage(equi, cubemap.CubemapImageOpts{VertOffset: 0.41, DivAmt: multAmt})
			}
			zipped, err := cubemap.SkyPack(name, cubeMap, "Sky ported")
			if err != nil {
				dialog.NewError(err, w).Show()
				return
			}
			utils.SaveFile(zipped, name+"_cubemap.mcpack", w)
		})
		f.picker = utils.NewFilePicker(func(isReady bool) {
			if isReady {
				f.button.Enable()
			} else {
				f.button.Disable()
			}
		}, []string{".jpg", ".jpeg", ".png", ".webp"})
	}
	f.button.Disable()

	text := canvas.NewText("Swim Sky Porter", fyne.CurrentApp().Settings().Theme().Color(theme.ColorNameForeground, fyne.CurrentApp().Settings().ThemeVariant()))
	text.TextSize = 50
	return container.NewCenter(container.NewVBox(text, f.picker.Show(w), shiftCubemapCheck, f.button))
}

func (f *SkyPorter) OnDrop(uri fyne.URI) {
	if f.picker != nil {
		f.picker.OnDrop(uri)
	}
}
