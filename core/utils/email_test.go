package utils

import (
	"testing"
)

func TestSendEmailGun(t *testing.T) {
	err := SendEmailGun(
		"locci.carlo.85@gmail.com",
		&Email{
			Title:   "Test Mail Gun!",
			Message: "Hello Email Gun!",
		},
		"registration",
	)
	if err != nil {
		t.Fatal(err)
	}
}
