package loader

import (
	"testing"
)

func TestLoadProductCatelog(t *testing.T) {
	catelog := LoadProductCatelog("../products.json")
	if catelog["Vintage Typewriter"].ID != "OLJCESPC7Z" {
		t.Fail()
	}

}
