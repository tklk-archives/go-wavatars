package wavatars // import "src.techknowlogick.com/wavatars"

import (
	"bytes"
	"crypto/sha512"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math/rand"
)

var (
	face  = 11
	bg_color  = 240
	fade = 4 // BACKGROUNDS
	wav_color  = 240
	brow  = 8
	eyes = 13
	pupil = 11
	mouth = 19
)

var bodyParts = []string{"fade", "mask", "shine", "brow", "eyes", "pupils", "mouth"}

func New(seed []byte) image.Image {
	buf := sha512.Sum512(seed)
	seed_int := binary.BigEndian.Uint64(buf[56:])
	rand.Seed(int64(seed_int))
	
	w_face := rand.Intn(face) + 1
	w_bg_color := rand.Intn(bg_color) + 1
	w_fade := rand.Intn(fade) + 1
	w_wav_color := rand.Intn(wav_color) + 1
	w_brow := rand.Intn(brow) + 1
	w_eyes := rand.Intn(eyes) + 1
	w_pupil := rand.Intn(pupil) + 1
	w_mouth := rand.Intn(mouth) + 1

	img := image.NewRGBA(image.Rect(0, 0, 80, 80))

	// Pick a random color for the background
	c := wavatar_hsl(w_bg_color, 240, 50)
	bg := color.RGBA{uint8(c[0]), uint8(c[1]), uint8(c[2]), 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{bg}, image.ZP, draw.Src)
	c1 := wavatar_hsl(w_wav_color, 240, 170)
	bg1 := color.RGBA{uint8(c1[0]), uint8(c1[1]), uint8(c1[2]), 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{bg1}, image.ZP, draw.Over)

	// Now add the various layers onto the image
	for _, part := range bodyParts {
		num := 0
		switch part {
		case "fade":
			num = w_fade
		case "mask":
			num = w_face
		case "shine":
			num = w_face
		case "brow":
			num = w_brow
		case "eyes":
			num = w_eyes
		case "pupils":
			num = w_pupil
		case "mouth":
			num = w_mouth
		}
		fileName := fmt.Sprintf("parts/%s%d.png", part, num)
		asset, err := Asset(fileName)
		if err != nil {
			log.Fatal(err)
		}
		assetFull, err := png.Decode(bytes.NewReader(asset))
		if err != nil {
			log.Fatal(err)
		}
		draw.Draw(img, img.Bounds(), assetFull, image.ZP, draw.Over)
	}
	return img
}

func wavatar_hsl(h, s, l int) ([]int) {
	var r,g,b int
	if h > 240 || h < 0 || s > 240 || s < 0 || l > 240 || l < 0 {
		return []int{0, 0, 0}
	}
	if h <= 40 {
		r = 255
		g = int(h/24*256)
		b = 0
	} else if h > 40 && h <= 80 {
		r = int(1-(h-40)/40)*256
		g = 255
		b = 0
	} else if h > 80 && h <= 120 {
		r = 0
		g = 255
		b = int((h-80)/40*256)
	} else if h > 120 && h <= 160 {
		r = 0
		g = int(1-(h-120)/40*256)
		b = 255
	} else if h > 160 && h <= 200 {
		r = int((h-160)/40*256)
		g = 0
		b = 255
	} else if h > 200 {
		r = 255
		g = 0
		b = int((h-200)/40*256)
	}
	r = r + (240-s)/240*(128-r)
	g = g + (240-s)/240*(128-g)
	b = b + (240-s)/240*(128-b)
	if l < 120 {
		r = (r/120) * l
		g = (g/120) * l
		b = (b/120) * l
	} else {
		r = l * ((256-r)/120)+2*r-256
		g = l * ((256-g)/120)+2*g-256
		b = l * ((256-b)/120)+2*b-256
	}
	return []int{wavatar_clamp(r), wavatar_clamp(g), wavatar_clamp(b)}
	
}

func wavatar_clamp(v int) int {
	if v < 0 {
		return 0
	} else if v > 255 {
		return 255
	}
	return v
}
