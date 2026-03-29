package tools

import (
	"errors"
	"io"
	"math"
	"swim-pack-tool/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/swim-services/swim_porter/rescale"
)

type Rescale struct {
	picker *utils.FilePickerOrLink
	button *widget.Button
}

func (r *Rescale) View(w fyne.Window) fyne.CanvasObject {
	alg := utils.NewObjectWithLabel("Algorithm:", widget.NewSelect(rescale.GetAlgorithms(), nil))
	alg.Obj().SetSelected("box")

	scale := utils.NewObjectWithLabel("Scale:", widget.NewSelect([]string{"8x", "16x", "32x", "64x", "128x"}, nil))
	scale.Obj().SetSelectedIndex(2)

	if r.button == nil {
		r.button = widget.NewButton("Rescale", func() {
			if r.picker == nil {
				return
			}
			r.button.Disable()
			r.button.SetText("Rescaling...")
			defer r.button.SetText("Rescale")
			defer r.picker.Clear()
			alg, ok := rescale.ParseAlgorithm(alg.Obj().Selected)
			if !ok {
				dialog.NewError(errors.New("unknown algorithm"), w).Show()
				return
			}
			reader, err := r.picker.Reader()
			if err != nil {
				dialog.NewError(err, w).Show()
				return
			}
			data, err := io.ReadAll(reader)
			if err != nil {
				dialog.NewError(err, w).Show()
				return
			}
			scale := int(math.Pow(2, (float64(scale.Obj().SelectedIndex() + 3))))
			name := utils.TrimMultiSpace(utils.RemoveExtension(r.picker.Filename()))
			extracted, err := utils.LoadArchive(r.picker.Filename(), data)
			if err != nil {
				dialog.NewError(err, w).Show()
				return
			}
			err = rescale.RescaleRaw(extracted, scale, rescale.RescaleOptions{
				ShowCredits: true,
				Algorithm:   alg,
			})
			if err != nil {
				dialog.NewError(err, w).Show()
				return
			}
			if err := utils.SaveMapFsAsZip(extracted, name+"_RESCALE.mcpack", w); err != nil {
				dialog.NewError(err, w).Show()
				return
			}
		})
		r.picker = utils.NewFilePickerOrMediafire(func(isReady bool) {
			if isReady {
				r.button.Enable()
			} else {
				r.button.Disable()
			}
		}, []string{".zip", ".rar", ".mcpack"})
	}
	r.button.Disable()

	text := canvas.NewText("Swim Pack Rescaler", fyne.CurrentApp().Settings().Theme().Color(theme.ColorNameForeground, fyne.CurrentApp().Settings().ThemeVariant()))
	text.TextSize = 50
	return container.NewCenter(container.NewVBox(text, r.picker.Show(w), scale.Container(), alg.Container(), r.button))
}

func (r *Rescale) OnDrop(uri fyne.URI) {
	if r.picker != nil {
		r.picker.OnDrop(uri)
	}
}
