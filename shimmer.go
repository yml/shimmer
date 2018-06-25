// +build js,wasm

package shimmer

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"math"
	"strings"
	"syscall/js"
	"time"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

type Shimmer struct {
	buf                        bytes.Buffer
	onImgLoadCb, shutdownCb    js.Callback
	brightnessCb, contrastCb   js.Callback
	hueCb, satCb, fileUploadCb js.Callback
	sourceImg                  image.Image
	dstWidth                   int
	resizedImg                 image.Image

	done chan struct{}
}

// New returns a new instance of shimmer
func New() *Shimmer {
	return &Shimmer{
		dstWidth: 300,
		done:     make(chan struct{}),
	}
}

// Start sets up all the callbacks and waits for the close signal
// to be sent from the browser.
func (s *Shimmer) Start() {
	// Setup callbacks
	s.setupOnFileUploadCb()
	js.Global.Get("document").
		Call("getElementById", "uploader").
		Call("addEventListener", "change", s.fileUploadCb)

	s.setupBrightnessCb()
	js.Global.Get("document").
		Call("getElementById", "brightness").
		Call("addEventListener", "change", s.brightnessCb)

	s.setupContrastCb()
	js.Global.Get("document").
		Call("getElementById", "contrast").
		Call("addEventListener", "change", s.contrastCb)

	s.setupHueCb()
	js.Global.Get("document").
		Call("getElementById", "hue").
		Call("addEventListener", "change", s.hueCb)

	s.setupSatCb()
	js.Global.Get("document").
		Call("getElementById", "sat").
		Call("addEventListener", "change", s.satCb)

	s.setupShutdownCb()
	js.Global.Get("document").
		Call("getElementById", "close").
		Call("addEventListener", "click", s.shutdownCb)

	<-s.done
	s.log("Shutting down app")
	s.onImgLoadCb.Close()
	s.brightnessCb.Close()
	s.contrastCb.Close()
	s.hueCb.Close()
	s.satCb.Close()
	s.shutdownCb.Close()
}

func (s *Shimmer) setSourceImageFromString(simg string) {
	switch {
	case strings.HasPrefix(simg, jpegPrefix):
		simg = strings.TrimPrefix(simg, jpegPrefix)
	case strings.HasPrefix(simg, pngPrefix):
		simg = strings.TrimPrefix(simg, pngPrefix)
	default:
		s.log("unrecognized image format")
		return
	}

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(simg))
	var err error
	s.sourceImg, _, err = image.Decode(reader)
	if err != nil {
		s.log(err.Error())
		return
	}
	srcWidth, srcHeight := s.sourceImg.Bounds().Dx(), s.sourceImg.Bounds().Dy()
	dstWidth := s.dstWidth
	ratio := float64(srcHeight) / float64(srcWidth)
	dstHeight := int(math.Ceil(ratio * float64(dstWidth)))
	s.resizedImg = transform.Resize(s.sourceImg, dstWidth, dstHeight, transform.Linear)
	enc := imgio.JPEGEncoder(90)
	err = enc(&s.buf, s.resizedImg)
	if err != nil {
		s.log(err.Error())
		return
	}
	// setting the previewImg src property
	js.Global.Get("document").
		Call("getElementById", "previewImg").
		Set("src", jpegPrefix+base64.StdEncoding.EncodeToString(s.buf.Bytes()))
	js.Global.Get("document").
		Call("getElementById", "targetImg").
		Set("src", jpegPrefix+base64.StdEncoding.EncodeToString(s.buf.Bytes()))
	s.buf.Reset()
}

// updateImage writes the image to a byte buffer and then converts it to base64.
// Then it sets the value to the src attribute of the target image.
func (s *Shimmer) updateImage(img *image.RGBA, start time.Time) {
	enc := imgio.JPEGEncoder(90)
	err := enc(&s.buf, img)
	if err != nil {
		s.log(err.Error())
		return
	}
	// Setting the src property
	js.Global.Get("document").
		Call("getElementById", "targetImg").
		Set("src", jpegPrefix+base64.StdEncoding.EncodeToString(s.buf.Bytes()))
	fmt.Println("time taken:", time.Now().Sub(start))
	s.buf.Reset()
}

// utility function to log a msg to the UI from inside a callback
func (s *Shimmer) log(msg string) {
	js.Global.Get("document").
		Call("getElementById", "status").
		Set("innerText", msg)
}

func (s *Shimmer) setupShutdownCb() {
	s.shutdownCb = js.NewEventCallback(js.PreventDefault, func(ev js.Value) {
		s.done <- struct{}{}
	})
}
