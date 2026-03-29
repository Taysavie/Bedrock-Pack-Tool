package utils

import (
	"io"
	"net/http"
	"os"
	"path"
	"slices"
	"swim-pack-tool/mediafire"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/gameparrot/gifdl"
)

type FilePickerOrLink struct {
	reader    io.ReadCloser
	button    *widget.Button
	linkEntry *widget.Entry

	isReady          bool
	readyStateChange func(isReady bool)

	allowedExtensions []string

	filename string

	links    []string
	linkName string

	linkFetchFunc func(url string) (data io.ReadCloser, name string, err error)
}

func NewFilePickerOrGifLink(readyStateChange func(isReady bool), allowedExtensions []string) *FilePickerOrLink {
	return NewFilePickerOrLink(readyStateChange, allowedExtensions, []string{"tenor.com", "giphy.com"}, "or enter GIF link (tenor/giphy):", func(url string) (data io.ReadCloser, name string, err error) {
		downloadUrl, title, err := gifdl.GetGIFDownloadUrl(url)
		if err != nil {
			return nil, "", err
		}
		resp, err := http.Get(downloadUrl)
		if err != nil {
			return nil, "", err
		}
		return resp.Body, title + ".gif", nil
	})
}

func NewFilePickerOrMediafire(readyStateChange func(isReady bool), allowedExtensions []string) *FilePickerOrLink {
	return NewFilePickerOrLink(readyStateChange, allowedExtensions, []string{"mediafire.com"}, "or enter Mediafire link:", func(url string) (data io.ReadCloser, name string, err error) {
		data, _, name, err = mediafire.MediafireDownload(url)
		return data, name, err
	})
}

func NewFilePickerOrLink(readyStateChange func(isReady bool), allowedExtensions []string, links []string, linkName string, linkFetchFunc func(url string) (data io.ReadCloser, name string, err error)) *FilePickerOrLink {
	entry := widget.NewEntry()
	f := &FilePickerOrLink{linkEntry: entry, readyStateChange: readyStateChange, allowedExtensions: allowedExtensions, links: links, linkName: linkName, linkFetchFunc: linkFetchFunc}
	entry.OnChanged = func(s string) {
		valid := slices.Contains(links, GetURLHost(s)) || f.reader != nil
		if valid != f.isReady {
			f.isReady = valid
			f.readyStateChange(f.isReady)
		}
	}
	return f
}

func (f *FilePickerOrLink) Show(w fyne.Window) fyne.CanvasObject {
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
	return container.NewVBox(
		f.button,
		widget.NewLabelWithStyle(f.linkName, fyne.TextAlignCenter, fyne.TextStyle{}),
		f.linkEntry,
	)
}

func (f *FilePickerOrLink) Clear() {
	if f.reader != nil {
		f.reader.Close()
		f.reader = nil
	}
	f.linkEntry.SetText("")
	if f.isReady {
		f.readyStateChange(false)
		f.isReady = false
	}
	f.button.SetText("Choose file")
}

func (f *FilePickerOrLink) Reader() (io.ReadCloser, error) {
	if f.reader != nil {
		return f.reader, nil
	}
	data, name, err := f.linkFetchFunc(f.linkEntry.Text)
	if err != nil {
		return nil, err
	}
	f.filename = name
	return data, nil
}

func (f *FilePickerOrLink) Filename() string {
	return f.filename
}

func (f *FilePickerOrLink) OnDrop(uri fyne.URI) {
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
