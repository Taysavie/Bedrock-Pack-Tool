package utils

import (
	"io"
	"os"
	"path"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type FilePicker struct {
	reader           io.ReadCloser
	button           *widget.Button
	isReady          bool
	readyStateChange func(isReady bool)

	allowedExtensions []string

	filename string
}

func NewFilePicker(readyStateChange func(isReady bool), allowedExtensions []string) *FilePicker {
	f := &FilePicker{readyStateChange: readyStateChange, allowedExtensions: allowedExtensions}
	return f
}

func (f *FilePicker) Show(w fyne.Window) fyne.CanvasObject {
	if f.reader != nil {
		f.reader.Close()
		f.reader = nil
		f.button.SetText("Choose file")
	}
	f.button = widget.NewButton("Choose file", func() {
		if f.reader != nil {
			f.reader.Close()
			f.reader = nil
			f.button.SetText("Choose file")
		}
		d := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if reader == nil {
				return
			}
			f.reader = reader
			f.filename = reader.URI().Name()
			f.button.SetText("Choose file (" + f.filename + ")")
			f.isReady = true
			f.readyStateChange(f.isReady)
		}, w)
		d.Resize(fyne.Size{Width: 1000, Height: 700})
		d.SetFilter(storage.NewExtensionFileFilter(f.allowedExtensions))
		d.Show()
	})
	return f.button
}

func (f *FilePicker) Clear() {
	if f.reader != nil {
		f.reader.Close()
		f.reader = nil
	}
	if f.isReady {
		f.readyStateChange(false)
		f.isReady = false
	}
	f.button.SetText("Choose file")
}

func (f *FilePicker) Reader() io.ReadCloser {
	return f.reader
}

func (f *FilePicker) Filename() string {
	return f.filename
}

func (f *FilePicker) OnDrop(uri fyne.URI) {
	ext := path.Ext(uri.Name())
	if slices.Contains(f.allowedExtensions, ext) {
		var err error
		f.reader, err = os.Open(uri.Path())
		if err == nil {
			f.filename = uri.Name()
			f.button.SetText("Choose file (" + f.filename + ")")
			if !f.isReady {
				f.isReady = true
				f.readyStateChange(true)
			}
		}
	}
}
