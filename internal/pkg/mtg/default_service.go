package mtg

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"

	"github.com/fogleman/gg"
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

func CombineImage(imgA image.Image, imgB image.Image) io.Reader {
	w := imgA.Bounds().Size().X * 2
	h := imgA.Bounds().Size().Y
	dc := gg.NewContext(w, h)

	dc.DrawImage(imgA, 0, 0)
	dc.DrawImage(imgB, imgA.Bounds().Size().X, 0)

	reader, writer := io.Pipe()
	go func() {
		jpeg.Encode(writer, dc.Image(), &jpeg.Options{Quality: 80})
		writer.Close()
	}()

	return reader
}

func RetrievePNG(url string) (image.Image, error){
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