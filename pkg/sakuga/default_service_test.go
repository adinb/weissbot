package sakuga

import (
	"errors"
	"testing"
)

type MockSakugaRepository struct {
	IsFail bool
}

func (r MockSakugaRepository) Get() (Sakuga, error){
	if r.IsFail {
		return Sakuga{}, errors.New("Failed to get sakuga") 
	}

	return Sakuga{URL: "https://sakugabooru.com/data/ed900fb57e3031f19f95077e0ebdfead.mp4"}, nil
}

func TestDefaultServiceGetSakugaSuccess(t *testing.T) {
	r := MockSakugaRepository{IsFail: false}
	s := DefaultService{Repo: r}
	sakuga, err := s.GetSakuga()
	if err != nil {
		t.Error("Failed to get sakuga")
	}

	if sakuga.URL != "https://sakugabooru.com/data/ed900fb57e3031f19f95077e0ebdfead.mp4" {
		t.Error("Failed to get sakuga")
	}
}

func TestDefaultServiceGetSakugaFail(t *testing.T) {
	r := MockSakugaRepository{IsFail: true}
	s := DefaultService{Repo: r}
	_, err := s.GetSakuga()
	if err == nil {
		t.Error("Failed to return error")
	}
}