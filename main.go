package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/EndFirstCorp/onedb"
	"github.com/EndFirstCorp/onedb/pgx"
)

type sortedBook struct {
	Book  string
	Order int
}

var bom = map[string]*sortedBook{
	"introduction": &sortedBook{"Introduction to the Book of Mormon", 101},
	"three":        &sortedBook{"Testimony of Three Witnesses", 102},
	"1-ne":         &sortedBook{"1 Nephi", 103},
	"2-ne":         &sortedBook{"2 Nephi", 104},
	"jacob":        &sortedBook{"Jacob", 105},
	"enos":         &sortedBook{"Enos", 106},
	"jarom":        &sortedBook{"Jarom", 107},
	"omni":         &sortedBook{"Omni", 108},
	"w-of-m":       &sortedBook{"Words of Mormon", 109},
	"mosiah":       &sortedBook{"Mosiah", 110},
	"alma":         &sortedBook{"Alma", 111},
	"hel":          &sortedBook{"Helaman", 112},
	"3-ne":         &sortedBook{"3 Nephi", 113},
	"4-ne":         &sortedBook{"4 Nephi", 114},
	"morm":         &sortedBook{"Mormon", 115},
	"ether":        &sortedBook{"Ether", 116},
	"moro":         &sortedBook{"Moroni", 117},
}

var pgp = map[string]*sortedBook{
	"introduction": &sortedBook{"Introduction to the Pearl of Great Price", 401},
	"moses":        &sortedBook{"Moses", 402},
	"abr":          &sortedBook{"Abraham", 403},
	"js-m":         &sortedBook{"Joseph Smith--Matthew", 404},
	"js-h":         &sortedBook{"Joseph Smith--History", 405},
	"a-of-f":       &sortedBook{"Articles of Faith", 406},
}

var nt = map[string]*sortedBook{
	"matt":   &sortedBook{"Matthew", 501},
	"mark":   &sortedBook{"Mark", 502},
	"luke":   &sortedBook{"Luke", 503},
	"john":   &sortedBook{"John", 504},
	"acts":   &sortedBook{"Acts", 505},
	"rom":    &sortedBook{"Romans", 506},
	"1-cor":  &sortedBook{"1 Corinthians", 507},
	"2-cor":  &sortedBook{"2 Corinthians", 508},
	"gal":    &sortedBook{"Galatians", 509},
	"eph":    &sortedBook{"Ephesians", 510},
	"philip": &sortedBook{"Philippians", 511},
	"col":    &sortedBook{"Colossians", 512},
	"1-thes": &sortedBook{"1 Thessalonians", 513},
	"2-thes": &sortedBook{"2 Thessalonians", 514},
	"1-tim":  &sortedBook{"1 Timothy", 515},
	"2-tim":  &sortedBook{"2 Timothy", 516},
	"titus":  &sortedBook{"Titus", 517},
	"philem": &sortedBook{"Philemon", 518},
	"heb":    &sortedBook{"Hebrews", 519},
	"james":  &sortedBook{"James", 520},
	"1-pet":  &sortedBook{"1 Peter", 521},
	"2-pet":  &sortedBook{"2 Peter", 522},
	"1-jn":   &sortedBook{"1 John", 523},
	"2-jn":   &sortedBook{"2 John", 524},
	"3-jn":   &sortedBook{"3 John", 525},
	"jude":   &sortedBook{"Jude", 526},
	"rev":    &sortedBook{"Revelation", 527},
}

var ot = map[string]*sortedBook{
	"gen":   &sortedBook{"Genesis", 601},
	"ex":    &sortedBook{"Exodus", 602},
	"lev":   &sortedBook{"Leviticus", 603},
	"num":   &sortedBook{"Numbers", 604},
	"deut":  &sortedBook{"Deuteronomy", 605},
	"josh":  &sortedBook{"Joshua", 606},
	"judg":  &sortedBook{"Judges", 607},
	"ruth":  &sortedBook{"Ruth", 608},
	"1-sam": &sortedBook{"1 Samuel", 609},
	"2-sam": &sortedBook{"2 Samuel", 610},
	"1-kgs": &sortedBook{"1 Kings", 611},
	"2-kgs": &sortedBook{"2 Kings", 612},
	"1-chr": &sortedBook{"1 Chronicles", 613},
	"2-chr": &sortedBook{"2 Chronicles", 614},
	"ezra":  &sortedBook{"Ezra", 615},
	"neh":   &sortedBook{"Nehemiah", 616},
	"esth":  &sortedBook{"Esther", 617},
	"job":   &sortedBook{"Job", 618},
	"ps":    &sortedBook{"Psalms", 619},
	"prov":  &sortedBook{"Proverbs", 620},
	"eccl":  &sortedBook{"Ecclesiastes", 621},
	"song":  &sortedBook{"Song of Solomon", 622},
	"isa":   &sortedBook{"Isaiah", 623},
	"jer":   &sortedBook{"Jeremiah", 624},
	"lam":   &sortedBook{"Lamentations", 625},
	"ezek":  &sortedBook{"Ezekial", 626},
	"dan":   &sortedBook{"Daniel", 627},
	"hosea": &sortedBook{"Hosea", 628},
	"joel":  &sortedBook{"Joel", 629},
	"amos":  &sortedBook{"Amos", 630},
	"obad":  &sortedBook{"Obadiah", 631},
	"jonah": &sortedBook{"Jonah", 632},
	"micah": &sortedBook{"Micah", 633},
	"nahum": &sortedBook{"Nahum", 634},
	"hab":   &sortedBook{"Habakkuk", 635},
	"zeph":  &sortedBook{"Zephaniah", 636},
	"hag":   &sortedBook{"Haggai", 637},
	"zech":  &sortedBook{"Zechariah", 638},
	"mal":   &sortedBook{"Malachi", 639},
}

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

type simpleNote struct {
	Notes      string
	Tags       []string
	Order      int
	References []reference
}

type reference struct {
	FullReference  string
	ShortReference string
	Text           string
}

const query = `SELECT v.scripture_text 
FROM verses v 
JOIN chapters c ON v.chapter_id = c.id 
JOIN books b ON c.book_id = b.id 
JOIN volumes vo ON b.volume_id = vo.id 
WHERE b.book_title = $1 and c.chapter_number = $2 and v.verse_number = $3;`

func main() {
	f, err := os.Open("lds notes.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	db, err := pgx.NewPgx("localhost", 5432, "postgres", "", "scriptures")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	w, err := os.Create("out.html")
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	result := []note{}
	if err := json.NewDecoder(f).Decode(&result); err != nil {
		log.Fatal(err)
	}

	simpleNotes := getSimpleNotes(result, db)
	orderedNotes := orderSimpleNotes(simpleNotes)

	w.WriteString(`<!DOCTYPE html>
<html>
<head>
<title>Rob's notes</title>
</head>
<body>
<table border="1" cellPadding="5" cellSpacing="0">
<tr>
	<td width="100">Tags</td>
	<td width="50%">Reference</td>
	<td width="50%">Note</td>
</tr>`)
	for _, note := range orderedNotes {
		w.WriteString(fmt.Sprintf(`
<tr>
	<td>%v</td>
	<td>%s</td>
	<td>%s</td>
</tr>`+"\n", note.Tags, formatReferences(note.References), note.Notes))
	}
	w.WriteString(`
</table></body></html>`)
}

func orderSimpleNotes(notes []simpleNote) []simpleNote {
	sort.Slice(notes, func(i, j int) bool {
		return notes[i].Order < notes[j].Order
	})
	return notes
}

func getText(db pgx.PGXer, book string, chapter, verse int) (string, error) {
	var text string
	return text, db.QueryValues(onedb.NewQuery(query, book, chapter, verse), &text)
}

func formatReferences(refs []reference) string {
	var buf strings.Builder
	for i, ref := range refs {
		if i == 0 {
			buf.WriteString(fmt.Sprintf("<b>%s</b> &nbsp;%s", ref.FullReference, ref.Text))
		} else {
			buf.WriteString(fmt.Sprintf("<b>%s</b> &nbsp;%s", ref.ShortReference, ref.Text))
		}
		if i != len(refs)-1 {
			buf.WriteString("<br/>")
		}
	}
	return buf.String()
}

func getSimpleNotes(fullNotes []note, db pgx.PGXer) []simpleNote {
	notes := []simpleNote{}
	for _, note := range fullNotes {
		uris := []string{note.URI}
		if len(note.Highlight.Content) > 0 {
			uris = []string{}
			for _, content := range note.Highlight.Content {
				uris = append(uris, content.URI)
			}
		}
		notes = append(notes, *newSimpleNote(uris, note.Note.Content, note.Tags, db))
	}
	return notes
}

func newSimpleNote(uris []string, notes string, tags []string, db pgx.PGXer) *simpleNote {
	note := &simpleNote{Notes: notes, Tags: tags}
	for i, uri := range uris {
		var order int
		refs := strings.Split(uri, "/")
		if len(refs) > 1 {
			switch refs[1] {
			case "scriptures":
				volume, book, chapter, verse := parseRefs(refs[2:])
				if newBook := getBook(volume, book); newBook != nil {
					book = newBook.Book
					order = newBook.Order
				}
				verse = getVerse(verse)
				verseInt, _ := strconv.Atoi(verse)
				chapterInt, _ := strconv.Atoi(chapter)
				text, err := getText(db, book, chapterInt, verseInt)
				if err != nil {
					fmt.Println(book, chapterInt, verseInt, err)
				}
				note.References = append(note.References, reference{fmt.Sprintf("%s %s:%s", book, chapter, verse), verse, text})
				if i == 0 {
					note.Order = order*1000000 + chapterInt*1000 + verseInt
				}
			case "general-conference", "ensign":
				year, month, title, _ := parseRefs(refs[2:])
				yearInt, _ := strconv.Atoi(year)
				monthInt, _ := strconv.Atoi(month)
				note.Order = 2000000 + yearInt*100 + monthInt
				note.References = append(note.References, reference{fmt.Sprintf("%s (%s %s-%s)", formatTitle(title), refs[1], formatMonth(month), year), "", ""})
			case "manual":
				volume, page, _, paragraph := parseRefs(refs[2:])
				note.Order = 3000000
				note.References = append(note.References, reference{fmt.Sprintf("%s (%s:%s)", formatTitle(volume), page, getVerse(paragraph)), "", ""})
			}
		}
	}
	return note
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

func getBook(volume, book string) *sortedBook {
	switch volume {
	case "bofm":
		return bom[book]
	case "dc-testament":
		return &sortedBook{"Doctrine and Covenants", 200}
	case "nt":
		return nt[book]
	case "ot":
		return ot[book]
	case "pgp":
		return pgp[book]
	}
	return &sortedBook{strings.Title(book), 0}
}

func getVerse(paragraph string) string {
	if len(paragraph) == 0 || paragraph[0] != 'p' {
		return paragraph
	}
	return paragraph[1:]
}
