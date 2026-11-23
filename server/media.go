package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"image"
	"image/color"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
)

const (
	maxImageBytes  = 30 << 20
	maxVideoBytes  = 800 << 20
	maxImagePixels = 80_000_000 // refuse decompression bombs before decoding
	jpegQuality    = 84
)

// Image sizes (longest side, px) per use on the site.
const (
	sizePerson  = 1000
	sizePoster  = 1000
	sizeGallery = 1920
	sizeThumb   = 640
	sizeWide    = 1600
)

var errBadImage = errors.New("файл не похож на фото: загрузите JPEG, PNG или WebP")

// Media stores uploads in DATA_DIR/media under random names, served at /media/.
type Media struct{ dir string }

func NewMedia(dir string) (*Media, error) {
	return &Media{dir: dir}, os.MkdirAll(dir, 0o755)
}

// SaveImage normalises an uploaded photo: fixes phone rotation, flattens transparency
// onto black (the site background), shrinks to maxSide and stores it as JPEG.
func (m *Media) SaveImage(r io.Reader, maxSide int) (string, error) {
	img, err := decodeUpload(r)
	if err != nil {
		return "", err
	}
	return m.writeJPEG(fit(img, maxSide))
}

// SaveImageWithThumb stores a full-size photo plus a small preview for grids.
func (m *Media) SaveImageWithThumb(r io.Reader, maxSide, thumbSide int) (full, thumb string, err error) {
	img, err := decodeUpload(r)
	if err != nil {
		return "", "", err
	}
	if full, err = m.writeJPEG(fit(img, maxSide)); err != nil {
		return "", "", err
	}
	if thumb, err = m.writeJPEG(fit(img, thumbSide)); err != nil {
		return "", "", err
	}
	return full, thumb, nil
}

// SaveVideo stores an MP4/WebM upload as-is.
func (m *Media) SaveVideo(r io.Reader) (string, error) {
	head := make([]byte, 512)
	n, err := io.ReadFull(r, head)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return "", errors.New("не удалось прочитать видео")
	}
	ext := map[string]string{"video/mp4": ".mp4", "video/webm": ".webm"}[http.DetectContentType(head[:n])]
	if ext == "" {
		return "", errors.New("видео должно быть в формате MP4 или WebM")
	}
	return m.write(ext, io.MultiReader(bytes.NewReader(head[:n]), r))
}

func (m *Media) writeJPEG(img image.Image) (string, error) {
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: jpegQuality}); err != nil {
		return "", err
	}
	return m.write(".jpg", &buf)
}

func (m *Media) write(ext string, r io.Reader) (string, error) {
	name := randHex(12) + ext
	tmp, err := os.CreateTemp(m.dir, ".up-*")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmp.Name())
	if _, err := io.Copy(tmp, r); err != nil {
		tmp.Close()
		return "", err
	}
	if err := tmp.Close(); err != nil {
		return "", err
	}
	if err := os.Rename(tmp.Name(), filepath.Join(m.dir, name)); err != nil {
		return "", err
	}
	return "/media/" + name, nil
}

func decodeUpload(r io.Reader) (image.Image, error) {
	data, err := io.ReadAll(io.LimitReader(r, maxImageBytes+1))
	if err != nil {
		return nil, err
	}
	if len(data) > maxImageBytes {
		return nil, errors.New("фото больше 30 МБ")
	}
	cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, errBadImage
	}
	if cfg.Width*cfg.Height > maxImagePixels {
		return nil, errors.New("фото слишком большое по размеру в пикселях")
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, errBadImage
	}
	if format == "jpeg" {
		img = orient(img, exifOrientation(data))
	}
	return flatten(img), nil
}

// fit scales img down so its longest side is at most max; it never upscales.
func fit(img image.Image, max int) image.Image {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	if w <= max && h <= max {
		return img
	}
	if w >= h {
		h, w = h*max/w, max
	} else {
		w, h = w*max/h, max
	}
	dst := image.NewRGBA(image.Rect(0, 0, max1(w), max1(h)))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, b, draw.Src, nil)
	return dst
}

// flatten draws img over black so transparent PNG areas match the dark site.
func flatten(img image.Image) *image.RGBA {
	b := img.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(dst, dst.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)
	draw.Draw(dst, dst.Bounds(), img, b.Min, draw.Over)
	return dst
}

// orient applies EXIF orientation (2..8) so phone photos aren't sideways.
func orient(img image.Image, o int) image.Image {
	if o < 2 || o > 8 {
		return img
	}
	src := flatten(img)
	w, h := src.Bounds().Dx(), src.Bounds().Dy()
	dw, dh := w, h
	if o >= 5 {
		dw, dh = h, w
	}
	dst := image.NewRGBA(image.Rect(0, 0, dw, dh))
	for y := 0; y < dh; y++ {
		for x := 0; x < dw; x++ {
			var sx, sy int
			switch o {
			case 2:
				sx, sy = w-1-x, y
			case 3:
				sx, sy = w-1-x, h-1-y
			case 4:
				sx, sy = x, h-1-y
			case 5:
				sx, sy = y, x
			case 6:
				sx, sy = y, h-1-x
			case 7:
				sx, sy = w-1-y, h-1-x
			case 8:
				sx, sy = w-1-y, x
			}
			dst.SetRGBA(x, y, src.RGBAAt(sx, sy))
		}
	}
	return dst
}

// exifOrientation reads tag 0x0112 from a JPEG's APP1/Exif segment; 1 if absent.
func exifOrientation(b []byte) int {
	if len(b) < 4 || b[0] != 0xFF || b[1] != 0xD8 {
		return 1
	}
	for i := 2; i+4 <= len(b); {
		if b[i] != 0xFF {
			return 1
		}
		marker := b[i+1]
		if marker == 0xDA || marker == 0xD9 { // start of scan / end: no more metadata
			return 1
		}
		size := int(binary.BigEndian.Uint16(b[i+2:]))
		if size < 2 || i+2+size > len(b) {
			return 1
		}
		seg := b[i+4 : i+2+size]
		if marker == 0xE1 && len(seg) > 14 && string(seg[:6]) == "Exif\x00\x00" {
			return tiffOrientation(seg[6:])
		}
		i += 2 + size
	}
	return 1
}

func tiffOrientation(t []byte) int {
	if len(t) < 8 {
		return 1
	}
	var bo binary.ByteOrder
	switch string(t[:2]) {
	case "II":
		bo = binary.LittleEndian
	case "MM":
		bo = binary.BigEndian
	default:
		return 1
	}
	ifd := int(bo.Uint32(t[4:]))
	if ifd+2 > len(t) {
		return 1
	}
	n := int(bo.Uint16(t[ifd:]))
	for k := 0; k < n; k++ {
		e := ifd + 2 + k*12
		if e+12 > len(t) {
			return 1
		}
		if bo.Uint16(t[e:]) == 0x0112 {
			if v := int(bo.Uint16(t[e+8:])); v >= 1 && v <= 8 {
				return v
			}
			return 1
		}
	}
	return 1
}

func max1(v int) int {
	if v < 1 {
		return 1
	}
	return v
}
