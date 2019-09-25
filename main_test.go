package main

import (
	"testing"
)

func TestGetReference(t *testing.T) {
	uri := "/general-conference/2015/10/faith-is-not-by-chance-but-by-choice"
	if ref := getReference(uri); ref != "Faith Is Not By Chance But By Choice (Oct 2015)" {
		t.Error("Expected valid reference", ref)
	}

	uri = "/scriptures/bofm/mosiah/9"
	if ref := getReference(uri); ref != "Mosiah 9:" {
		t.Error("Expected valid reference", ref)
	}

}
