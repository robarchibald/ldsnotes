package main

import (
	"testing"
)

func TestGetReference(t *testing.T) {
	uri := "/general-conference/2015/10/faith-is-not-by-chance-but-by-choice.p13"
	if ref := getReference(uri); ref != "Faith Is Not By Chance But By Choice (Oct 2015)" {
		t.Error("Expected valid reference", ref)
	}

	uri = "/scriptures/bofm/mosiah/9.p14"
	if ref := getReference(uri); ref != "Mosiah 9:14" {
		t.Error("Expected valid reference", ref)
	}

	uri = "/scriptures/bofm/introduction.p8"
	if ref := getReference(uri); ref != "Introduction 1:8" {
		t.Error("Expected valid reference", ref)
	}

	uri = "/manual/come-follow-me-for-individuals-and-families-new-testament-2019/18.p3"
	if ref := getReference(uri); ref != "Come Follow Me For Individuals And Families New Testament 2019 (18:3)" {
		t.Error("Expected valid reference", ref)
	}
}
