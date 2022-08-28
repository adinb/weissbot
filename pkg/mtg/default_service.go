package mtg

import (
	"image"
	"image/draw"

	//"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
)

type DefaultService struct {
	Repo RepositoryContract
}

func (s *DefaultService) SearchCardByName(name string) ([]*MagicCard, error) {
	cards, err := s.Repo.Find(name)
	if err != nil {
		return nil, err
	}

	return cards, nil
}

func TileImagesHorizontally(imgA image.Image, imgB image.Image) image.Image {
	//w := imgA.Bounds().Size().X + imgB.Bounds().Size().X
	//oldh := imgA.Bounds().Size().Y
	h := imgA.Bounds().Size().Y + imgB.Bounds().Size().Y
	//h := 0
	//if imgA.Bounds().Size().Y >= imgB.Bounds().Size().Y {
	//	h = imgA.Bounds().Size().Y
	//} else {
	//	h = imgB.Bounds().Size().Y
	//}

	//tiledImage := image.NewRGBA(image.Rect(0, 0, w, h))
	//rectA := image.Rect(0, 0, imgA.Bounds().Size().X, imgA.Bounds().Size().Y)
	//rectB := image.Rect(0, 0, w, h)

	//draw.Draw(tiledImage, rectA, imgA, image.Pt(0, 0), draw.Over)

	extra := make([]uint8, imgA.Bounds().Size().X*imgA.Bounds().Size().Y*4)
	imgA.(*image.NRGBA).Pix = append(imgA.(*image.NRGBA).Pix, extra...)
	imgA.(*image.NRGBA).Rect = image.Rect(0, 0, imgA.Bounds().Size().X, h)
	draw.Draw(imgA.(*image.NRGBA), imgA.(*image.NRGBA).Rect, imgB, image.Pt(0, imgB.Bounds().Size().Y*-1), draw.Over)

	resultimage := imgA
	return resultimage
}

func CreateImageReader(img image.Image) io.Reader {
	reader, writer := io.Pipe()
	go func() {
		jpeg.Encode(writer, img, &jpeg.Options{Quality: 80})
		writer.Close()
	}()

	return reader
}

func RetrievePNG(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	image, err := png.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	return image, nil
}
