package main

import (
	"testing"

	"github.com/EndFirstCorp/onedb/pgx"
)

type fakeData struct {
	Value string
}

func TestGetReference(t *testing.T) {
	uri := "/general-conference/2015/10/faith-is-not-by-chance-but-by-choice.p13"
	if ref, _ := getReference(uri, nil); ref.FullReference != "Faith Is Not By Chance But By Choice (Oct 2015)" {
		t.Error("Expected valid reference", ref)
	}

	db := pgx.NewMock(nil, nil, []fakeData{{"hello"}}, []fakeData{{"hello"}}, []fakeData{{"hello"}})
	uri = "/scriptures/bofm/mosiah/9.p14"
	if ref, _ := getReference(uri, db); ref.FullReference != "Mosiah 9:14" || ref.Text != "hello" {
		t.Error("Expected valid reference", ref)
	}

	uri = "/scriptures/bofm/introduction.p8"
	if ref, _ := getReference(uri, db); ref.FullReference != "Introduction to the Book of Mormon 1:8" {
		t.Error("Expected valid reference", ref)
	}

	uri = "/manual/come-follow-me-for-individuals-and-families-new-testament-2019/18.p3"
	if ref, _ := getReference(uri, nil); ref.FullReference != "Come Follow Me For Individuals And Families New Testament 2019 (18:3)" {
		t.Error("Expected valid reference", ref)
	}
}

func TestFilterSimpleNotes(t *testing.T) {
	if notes := filterSimpleNotes("tag", []simpleNote{
		{Tags: []string{"Faith", "tag"}},
		{Tags: []string{"Other"}},
	}); len(notes) != 1 {
		t.Error("Expected to filter to only one note", notes)
	}
}
