package main

import (
	"io/ioutil"
	"testing"
)

func TestLoadScryfallObjectResource(t *testing.T) {
	bolasDat, err := ioutil.ReadFile("testdata/bolas_ravager.json")
	if err != nil {
		panic(err)
	}
	_, err = LoadScryfallObjectResourceFromJSON(bolasDat)
	if err != nil {
		t.Error("Unmarshal failed", err.Error())
	}

	ghaltaDat, err := ioutil.ReadFile("testdata/ghalta.json")
	if err != nil {
		panic(err)
	}
	_, err = LoadScryfallObjectResourceFromJSON(ghaltaDat)
	if err != nil {
		t.Error("Unmarshal failed", err.Error())
	}
}
