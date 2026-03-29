package tools

import (
	"errors"
	"io"
	"swim-pack-tool/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/swim-services/swim_porter/cubemap"
	"github.com/swim-services/swim_porter/cubemap/blend"
)

type SkyOverlay struct {
	picker *utils.FilePicker
	button *widget.Button
}

func (f *SkyOverlay) View(w fyne.Window) fyne.CanvasObject {
	slider := utils.NewSliderWithValue(0, 49, 10)
	sliderContainer := utils.NewObjectWithLabel("Blend amount: ", slider.Container())
	shiftCubemapCheck := widget.NewCheck("Disable skybox offset", nil)

	if f.button == nil {
		f.button = widget.NewButton("Make sky overlay", func() {
			if f.picker == nil {
				return
			}
			f.button.Disable()
			f.button.SetText("Making sky overlay...")
			defer f.button.SetText("Make sky overlay")
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
			if img.Bounds().Dx() < 128 || img.Bounds().Dy() < 128 {
				dialog.NewError(errors.New("image must be at least 128x128"), w).Show()
				return
			}
			if slider.Slider().Value > 0 {
				img, err = blend.Blend(img, int(slider.Slider().Value))
				if err != nil {
					dialog.NewError(err, w).Show()
					return
				}
			}
			vertOffset := 0.41
			if shiftCubemapCheck.Checked {
				vertOffset = 0
			}
			zipped, err := cubemap.SkyPack(name, cubemap.CubemapFromImage(img, cubemap.CubemapImageOpts{VertOffset: vertOffset}), "Sky made")
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

	text := canvas.NewText("Swim Sky Maker", fyne.CurrentApp().Settings().Theme().Color(theme.ColorNameForeground, fyne.CurrentApp().Settings().ThemeVariant()))
	text.TextSize = 50
	return container.NewCenter(container.NewVBox(text, f.picker.Show(w), sliderContainer.Container(), shiftCubemapCheck, f.button))
}

func (f *SkyOverlay) OnDrop(uri fyne.URI) {
	if f.picker != nil {
		f.picker.OnDrop(uri)
	}
}
