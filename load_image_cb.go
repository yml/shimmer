// +build js,wasm

package shimmer

import (
	"syscall/js"
)

const (
	jpegPrefix = "data:image/jpeg;base64,"
	pngPrefix  = "data:image/png;base64,"
)

func (s *Shimmer) setupOnFileUploadCb() {
	s.fileUploadCb = js.NewEventCallback(js.PreventDefault, func(ev js.Value) {
		file := ev.Get("srcElement").Get("files").Get("0")
		freader := js.Global.Get("FileReader").New()
		freader.Set("onload", js.NewEventCallback(js.PreventDefault, func(ev js.Value) {
			s.setSourceImageFromString(ev.Get("target").Get("result").String())
		}))
		freader.Call("readAsDataURL", file)
	})
}
