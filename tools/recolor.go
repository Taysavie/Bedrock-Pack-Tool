package tools

import (
	"image/color"
	"io"
	"swim-pack-tool/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/swim-services/swim_porter/recolor"
)

type Recolor struct {
	picker *utils.FilePickerOrLink
	button *widget.Button
	ready  bool
}

func (r *Recolor) View(w fyne.Window) fyne.CanvasObject {
	var chosenColor color.Color
	colorRect := canvas.NewRectangle(chosenColor)
	colorPicker := container.NewBorder(nil, nil, nil, container.New(utils.NewFixedSizeLayout(fyne.NewSize(35, 30)), colorRect), widget.NewButton("Choose color", func() {
		p := dialog.NewColorPicker("Choose color", "", func(c color.Color) {
			chosenColor = c
			colorRect.FillColor = c
			colorRect.Refresh()
			newReady := chosenColor != nil && r.ready
			if newReady {
				r.button.Enable()
			} else {
				r.button.Disable()
			}
		}, w)
		p.Resize(fyne.Size{Width: 700, Height: 600})
		p.Advanced = true
		p.Show()
	}))

	alg := utils.NewObjectWithLabel("Algorithm:", widget.NewSelect([]string{"gray_tint", "tint", "hue"}, nil))
	alg.Obj().SetSelectedIndex(0)

	if r.button == nil {
		r.button = widget.NewButton("Recolor", func() {
			if r.picker == nil {
				return
			}
			r.button.Disable()
			r.button.SetText("Recoloring...")
			defer r.button.SetText("Recolor")
			defer r.picker.Clear()
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
			name := utils.TrimMultiSpace(utils.RemoveExtension(r.picker.Filename()))
			extracted, err := utils.LoadArchive(r.picker.Filename(), data)
			if err != nil {
				dialog.NewError(err, w).Show()
				return
			}
			r, g, b, _ := chosenColor.RGBA()
			err = recolor.RecolorRaw(extracted, recolor.RecolorOptions{
				NewColor:    color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: 255},
				ShowCredits: true,
				Alg:         alg.Obj().Selected,
			})
			if err != nil {
				dialog.NewError(err, w).Show()
				return
			}
			if err := utils.SaveMapFsAsZip(extracted, name+"_RECOLOR.mcpack", w); err != nil {
				dialog.NewError(err, w).Show()
				return
			}
		})
		r.picker = utils.NewFilePickerOrMediafire(func(isReady bool) {
			r.ready = isReady
			if isReady && chosenColor != nil {
				r.button.Enable()
			} else {
				r.button.Disable()
			}
		}, []string{".zip", ".rar", ".mcpack"})
	}
	r.button.Disable()

	text := canvas.NewText("Swim Pack Recolorer", fyne.CurrentApp().Settings().Theme().Color(theme.ColorNameForeground, fyne.CurrentApp().Settings().ThemeVariant()))
	text.TextSize = 50
	return container.NewCenter(container.NewVBox(text, r.picker.Show(w), colorPicker, alg.Container(), r.button))
}

func (r *Recolor) OnDrop(uri fyne.URI) {
	if r.picker != nil {
		r.picker.OnDrop(uri)
	}
}
