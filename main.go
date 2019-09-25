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
			return getScripture(refs[2:])
		case "general-conference":
			return getConference(refs[2:])
		}
	}
	return ref
}

func getScripture(refs []string) string {
	var volume, book, chapterVerse, chapter, verse string
	getValues(refs, &volume, &book, &chapterVerse)
	getValues(strings.Split(chapterVerse, "."), &chapter, &verse)
	return fmt.Sprintf("%s %s:%s", getBook(volume, book), chapter, verse)
}

func getConference(refs []string) string {
	var year, month, title string
	getValues(refs, &year, &month, &title)
	return fmt.Sprintf("%s (%s %s)", formatTitle(title), formatMonth(month), year)
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
	case "10":
		return "Oct"
	case "4":
		return "Apr"
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
