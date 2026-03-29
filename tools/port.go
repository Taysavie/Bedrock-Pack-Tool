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
	"github.com/swim-services/swim_porter/port"
)

type Port struct {
	picker *utils.FilePickerOrLink
	button *widget.Button
}

func (p *Port) View(w fyne.Window) fyne.CanvasObject {
	shiftCubemapCheck := widget.NewCheck("Disable skybox offset", nil)
	cubemapOverrideTextbox := utils.NewObjectWithLabel("Cubemap override:", widget.NewEntry())
	cubemapOverrideTextbox.Obj().PlaceHolder = "Auto"
	if p.button == nil {
		p.button = widget.NewButton("Port", func() {
			p.button.Disable()
			p.button.SetText("Porting...")
			defer p.button.SetText("Port")
			if p.picker == nil {
				return
			}
			defer p.picker.Clear()
			reader, err := p.picker.Reader()
			if err != nil {
				dialog.NewError(err, w).Show()
				return
			}
			data, err := io.ReadAll(reader)
			if err != nil {
				dialog.NewError(err, w).Show()
				return
			}
			name := utils.TrimMultiSpace(utils.RemoveExtension(p.picker.Filename()))
			extracted, err := utils.LoadArchive(p.picker.Filename(), data)
			if err != nil {
				dialog.NewError(err, w).Show()
				return
			}
			ported, err := port.PortRaw(extracted, name, port.PortOptions{
				ShowCredits:    true,
				OffsetSky:      !shiftCubemapCheck.Checked,
				SkyboxOverride: cubemapOverrideTextbox.Obj().Text,
			})
			if err != nil {
				dialog.NewError(err, w).Show()
				return
			}
			if err := utils.SaveMapFsAsZip(ported, name+"_PORT.mcpack", w); err != nil {
				dialog.NewError(err, w).Show()
				return
			}
		})
		p.picker = utils.NewFilePickerOrMediafire(func(isReady bool) {
			if isReady {
				p.button.Enable()
			} else {
				p.button.Disable()
			}
		}, []string{".zip", ".rar"})
	}
	p.button.Disable()

	text := canvas.NewText("Swim Pack Porter", fyne.CurrentApp().Settings().Theme().Color(theme.ColorNameForeground, fyne.CurrentApp().Settings().ThemeVariant()))
	text.TextSize = 50
	return container.NewCenter(container.NewVBox(text, p.picker.Show(w), shiftCubemapCheck, cubemapOverrideTextbox.Container(), p.button))
}

func (p *Port) OnDrop(uri fyne.URI) {
	if p.picker != nil {
		p.picker.OnDrop(uri)
	}
}
