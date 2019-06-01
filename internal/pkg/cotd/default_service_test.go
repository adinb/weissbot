package cotd

import (
	"errors"
	"testing"
)

type mockCOTDRepository struct {
	isFailed bool
	cotdList    []COTD
}

func (m *mockCOTDRepository) Get() ([]COTD, error) {
	if m.isFailed {
		return nil, errors.New("Failed to fetch the cards")
	}

	return m.cotdList, nil
}

func TestRetrieveCOTDSuccess(t *testing.T) {
	m := new(mockCOTDRepository)
	m.isFailed = false
	m.cotdList = append(m.cotdList, COTD{ImageURL: "https://s3-ap-northeast-1.amazonaws.com/cf-vanguard.com/wordpress/wp-content/images/todays-card/vgd_today0523_696468926482.png"})

	service := DefaultService{Repo: m}
	cotdList, err := service.GetCOTD()
	if err != nil {
		t.Error("Failed to get COTD", err)
	}

	if len(m.cotdList) != len(cotdList) {
		t.Error("Card entry mismatch")
	}
}

func TestRetrieveVanguardCOTDFail(t *testing.T) {
	m := new(mockCOTDRepository)
	m.isFailed = true
	service := DefaultService{Repo: m}
	_, err := service.GetCOTD()
	if err == nil {
		t.Error("Service didn't return error")
	}
}

func TestRetrieveVanguardCOTDNoCard(t *testing.T) {
	m := new(mockCOTDRepository)
	service := DefaultService{Repo: m}
	_, err := service.GetCOTD()
	if err == nil {
		t.Error("Service didn't return error")
	}
}
