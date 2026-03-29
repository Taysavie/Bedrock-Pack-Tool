package utils

import (
	"reflect"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type SliderWithValue struct {
	container    *fyne.Container
	slider       *widget.Slider
	valueDisplay *widget.Label
}

func NewSliderWithValue[T int | float64](min, max, value T) *SliderWithValue {
	isInt := reflect.TypeOf(min).Kind() == reflect.Int
	slider := widget.NewSlider(float64(min), float64(max))
	if !isInt {
		slider.Step = 0.01
	}

	label := widget.NewLabel("")

	slider.OnChanged = func(f float64) {
		if isInt {
			label.SetText(strconv.Itoa(int(f)))
		} else {
			label.SetText(strconv.FormatFloat(f, 'f', 1, 64))
		}
	}

	slider.SetValue(float64(value))

	return &SliderWithValue{slider: slider, valueDisplay: label, container: container.NewBorder(nil, nil, nil, label, slider)}
}

func (s *SliderWithValue) Slider() *widget.Slider {
	return s.slider
}

func (s *SliderWithValue) Label() *widget.Label {
	return s.valueDisplay
}

func (s *SliderWithValue) Container() *fyne.Container {
	return s.container
}
