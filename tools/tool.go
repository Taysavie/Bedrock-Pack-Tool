package tools

import "fyne.io/fyne/v2"

type Tool interface {
	View(w fyne.Window) fyne.CanvasObject
	OnDrop(uri fyne.URI)
}

func RegisterTools(r *Registry) {
	r.RegisterTool("Port", &Port{})
	r.RegisterTool("Recolor", &Recolor{})
	r.RegisterTool("Rescale", &Rescale{})
	r.RegisterTool("Fix font", &FixFont{})
	r.RegisterTool("Fix sky", &FixSky{})
	r.RegisterTool("Fix particles", &FixParticles{})
	r.RegisterTool("Sky overlay", &SkyOverlay{})
	r.RegisterTool("Java sky porter", &SkyPorter{})
	r.RegisterTool("Inventory maker", &InventoryMaker{})
	r.RegisterTool("Crosshair maker", &CrosshairMaker{})
	r.RegisterTool("Compressor", &Compressor{})
}
