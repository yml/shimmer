// +build js,wasm

package shimmer

import (
	"syscall/js"
	"time"

	"github.com/anthonynsimon/bild/adjust"
)

func (s *Shimmer) setupBrightnessCb() {
	s.brightnessCb = js.NewEventCallback(js.PreventDefault, func(ev js.Value) {
		// quick return if no source image is yet uploaded
		if s.resizedImg == nil {
			return
		}
		delta := ev.Get("target").Get("value").Float()
		start := time.Now()
		res := adjust.Brightness(s.resizedImg, delta)
		s.updateImage(res, start)
	})
}

func (s *Shimmer) setupContrastCb() {
	s.contrastCb = js.NewEventCallback(js.PreventDefault, func(ev js.Value) {
		// quick return if no source image is yet uploaded
		if s.resizedImg == nil {
			return
		}
		delta := ev.Get("target").Get("value").Float()
		start := time.Now()
		res := adjust.Contrast(s.resizedImg, delta)
		s.updateImage(res, start)
	})
}

func (s *Shimmer) setupHueCb() {
	s.hueCb = js.NewEventCallback(js.PreventDefault, func(ev js.Value) {
		// quick return if no source image is yet uploaded
		if s.resizedImg == nil {
			return
		}
		delta := ev.Get("target").Get("value").Int()
		start := time.Now()
		res := adjust.Hue(s.resizedImg, delta)
		s.updateImage(res, start)
	})
}

func (s *Shimmer) setupSatCb() {
	s.satCb = js.NewEventCallback(js.PreventDefault, func(ev js.Value) {
		// quick return if no source image is yet uploaded
		if s.resizedImg == nil {
			return
		}
		delta := ev.Get("target").Get("value").Float()
		start := time.Now()
		res := adjust.Saturation(s.resizedImg, delta)
		s.updateImage(res, start)
	})
}
