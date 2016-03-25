package regedit

import (
	"fmt"
	"testing"
)

func TestCurl(t *testing.T) {
	testSearch()
}

func testSearch() {
	reg := New(Query, HKCU, "\\Environment")
	regcmd := reg.Search(Reg{Key: "path"})
	if arr, err := regcmd.Exec(); err == nil {
		for _, v := range arr {
			fmt.Println(v)
		}
		fmt.Println("Complete!")
	}
}
