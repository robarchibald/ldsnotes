package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type note struct {
	ID             string    `json:"id"`
	ContentVersion string    `json:"contentVersion"`
	DocumentID     string    `json:"docId"`
	Locale         string    `json:"locale"`
	URI            string    `json:"uri"`
	PersonID       string    `json:"personId"`
	Type           string    `json:"type"`
	Highlight      highlight `json:"highlight"`
	Source         string    `json:"source"`
	Device         string    `json:"device"`
	Tags           []string  `json:"tags"`
	Refs           []string  `json:"refs"`
	Note           noteValue `json:"note"`
	Folders        []itemRef `json:"folders"`
	LastUpdated    time.Time `json:"lastUpdated"`
	Created        time.Time `json:"created"`
}

type highlight struct {
	Content []highlightContent `json:"content"`
}

type highlightContent struct {
	StartOffset string `json:"startOffset"`
	Color       string `json:"color"`
	EndOffset   string `json:"endOffset"`
	URI         string `json:"uri"`
	PID         string `json:"pid"`
}

type noteValue struct {
	Content string `json:"content"`
}

type itemRef struct {
	URI string `json:"uri"`
	ID  string `json:"id"`
}

func main() {
	f, err := os.Open("lds notes.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	w, err := os.Create("out.html")
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	result := []note{}
	if err := json.NewDecoder(f).Decode(&result); err != nil {
		log.Fatal(err)
	}
	w.WriteString(`<!DOCTYPE html>
<html>
<head>
<title>Rob's notes</title>
</head>
<body>
<table>
<tr>
	<td>Tags</td>
	<td>Reference</td>
	<td>Note</td>
</tr>`)
	for _, note := range result {
		uri := note.URI
		if len(note.Highlight.Content) > 0 {
			uri = note.Highlight.Content[0].URI
		}
		w.WriteString(fmt.Sprintf(`
<tr>
	<td>%v</td>
	<td>%s</td>
	<td>%s</td>
</tr>`+"\n", note.Tags, getReference(uri), note.Note.Content))
	}
	w.WriteString(`
</table></body></html>`)
}

func getReference(ref string) string {
	refs := strings.Split(ref, "/")
	if len(refs) > 1 {
		switch refs[1] {
		case "scriptures":
			volume, book, chapter, verse := parseRefs(refs[2:])
			return fmt.Sprintf("%s %s:%s", getBook(volume, book), chapter, getVerse(verse))
		case "general-conference", "ensign":
			year, month, title, _ := parseRefs(refs[2:])
			return fmt.Sprintf("%s (%s %s-%s)", formatTitle(title), refs[1], formatMonth(month), year)
		case "manual":
			volume, page, _, paragraph := parseRefs(refs[2:])
			return fmt.Sprintf("%s (%s:%s)", formatTitle(volume), page, getVerse(paragraph))
		}
	}
	return ref
}

func parseRefs(refs []string) (string, string, string, string) {
	var major, minor, sectionParagraph, section, paragraph string
	getValues(refs, &major, &minor, &sectionParagraph)
	if sectionParagraph == "" && minor != "" {
		getValues(strings.Split(minor, "."), &minor, &paragraph)
		section = "1"
	} else {
		getValues(strings.Split(sectionParagraph, "."), &section, &paragraph)
	}
	return major, minor, section, paragraph
}

func formatTitle(title string) string {
	words := strings.Split(title, "-")
	var buf strings.Builder
	for i, word := range words {
		buf.WriteString(strings.Title(word))
		if i != len(words)-1 {
			buf.WriteByte(' ')
		}
	}
	return buf.String()
}

func formatMonth(month string) string {
	switch month {
	case "1", "01":
		return "Jan"
	case "2", "02":
		return "Feb"
	case "3", "03":
		return "Mar"
	case "4", "04":
		return "Apr"
	case "5", "05":
		return "May"
	case "6", "06":
		return "Jun"
	case "7", "07":
		return "Jul"
	case "8", "08":
		return "Aug"
	case "9", "09":
		return "Sep"
	case "10":
		return "Oct"
	case "11":
		return "Nov"
	case "12":
		return "Dec"
	default:
		return month
	}
}

func getValues(refs []string, value ...*string) {
	for i := 0; i < len(refs) && i < len(value); i++ {
		*value[i] = refs[i]
	}
}

func getBook(volume, book string) string {
	switch volume {
	case "bofm":
		return getBOMBook(book)
	case "dc-testament":
		return "D&C"
	}
	return strings.Title(book)
}

func getBOMBook(book string) string {
	switch book {
	case "moro":
		return "Moroni"
	case "morm":
		return "Mormon"
	case "4-ne":
		return "4 Nephi"
	case "3-ne":
		return "3 Nephi"
	case "hel":
		return "Helaman"
	case "w-of-m":
		return "Words of Mormon"
	case "2-ne":
		return "2 Nephi"
	case "1-ne":
		return "1 Nephi"
	case "three":
		return "Testimony of Three Witnesses"
	default:
		return strings.Title(book)
	}
}

func getVerse(paragraph string) string {
	if len(paragraph) == 0 || paragraph[0] != 'p' {
		return paragraph
	}
	return paragraph[1:]
}
