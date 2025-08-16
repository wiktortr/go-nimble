package seda

import (
	"github.com/wiktortr/go-nimble/nimble"
	"net/url"
	"testing"
)

func TestSeda_Key(t *testing.T) {
	seda := Seda{}
	key := seda.Key()
	if key != "seda" {
		t.Errorf("Seda.Key() did not return 'seda', got %s", key)
	}
}

func TestSeda_Instantiate_no_params_values(t *testing.T) {
	seda := Seda{}
	_, err := seda.Instantiate(&nimble.ComponentParams{
		Uri:    "seda:test",
		Key:    "seda",
		Name:   "test",
		Values: nil,
	})
	if err != nil {
		t.Errorf("Seda.Instantiate() error = %v", err)
	}
}

func TestSeda_Instantiate_valid_buff_size(t *testing.T) {
	seda := Seda{}
	_, err := seda.Instantiate(&nimble.ComponentParams{
		Uri:  "seda:test",
		Key:  "seda",
		Name: "test",
		Values: url.Values{
			"buffSize": []string{"100"},
		},
	})
	if err != nil {
		t.Errorf("Seda.Instantiate() error = %v", err)
	}
}

func TestSeda_Instantiate_invalid_buff_size(t *testing.T) {
	seda := Seda{}
	_, err := seda.Instantiate(&nimble.ComponentParams{
		Uri:  "seda:test",
		Key:  "seda",
		Name: "test",
		Values: url.Values{
			"buffSize": []string{"abc"},
		},
	})
	if err == nil {
		t.Errorf("Seda.Instantiate() should return err")
	}
}
