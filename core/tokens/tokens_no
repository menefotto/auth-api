package tokens

import "testing"

func TestPut(t *testing.T) {
	err := BlackList.Put("Carlo", "Locci")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(BlackList.Valid("Carlo"))
}
