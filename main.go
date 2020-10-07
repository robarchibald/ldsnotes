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
	"introduction": {"Introduction to the Book of Mormon", 101},
	"three":        {"Testimony of Three Witnesses", 102},
	"1-ne":         {"1 Nephi", 103},
	"2-ne":         {"2 Nephi", 104},
	"jacob":        {"Jacob", 105},
	"enos":         {"Enos", 106},
	"jarom":        {"Jarom", 107},
	"omni":         {"Omni", 108},
	"w-of-m":       {"Words of Mormon", 109},
	"mosiah":       {"Mosiah", 110},
	"alma":         {"Alma", 111},
	"hel":          {"Helaman", 112},
	"3-ne":         {"3 Nephi", 113},
	"4-ne":         {"4 Nephi", 114},
	"morm":         {"Mormon", 115},
	"ether":        {"Ether", 116},
	"moro":         {"Moroni", 117},
}

var pgp = map[string]*sortedBook{
	"introduction": {"Introduction to the Pearl of Great Price", 401},
	"moses":        {"Moses", 402},
	"abr":          {"Abraham", 403},
	"js-m":         {"Joseph Smith--Matthew", 404},
	"js-h":         {"Joseph Smith--History", 405},
	"a-of-f":       {"Articles of Faith", 406},
}

var nt = map[string]*sortedBook{
	"matt":   {"Matthew", 501},
	"mark":   {"Mark", 502},
	"luke":   {"Luke", 503},
	"john":   {"John", 504},
	"acts":   {"Acts", 505},
	"rom":    {"Romans", 506},
	"1-cor":  {"1 Corinthians", 507},
	"2-cor":  {"2 Corinthians", 508},
	"gal":    {"Galatians", 509},
	"eph":    {"Ephesians", 510},
	"philip": {"Philippians", 511},
	"col":    {"Colossians", 512},
	"1-thes": {"1 Thessalonians", 513},
	"2-thes": {"2 Thessalonians", 514},
	"1-tim":  {"1 Timothy", 515},
	"2-tim":  {"2 Timothy", 516},
	"titus":  {"Titus", 517},
	"philem": {"Philemon", 518},
	"heb":    {"Hebrews", 519},
	"james":  {"James", 520},
	"1-pet":  {"1 Peter", 521},
	"2-pet":  {"2 Peter", 522},
	"1-jn":   {"1 John", 523},
	"2-jn":   {"2 John", 524},
	"3-jn":   {"3 John", 525},
	"jude":   {"Jude", 526},
	"rev":    {"Revelation", 527},
}

var ot = map[string]*sortedBook{
	"gen":   {"Genesis", 601},
	"ex":    {"Exodus", 602},
	"lev":   {"Leviticus", 603},
	"num":   {"Numbers", 604},
	"deut":  {"Deuteronomy", 605},
	"josh":  {"Joshua", 606},
	"judg":  {"Judges", 607},
	"ruth":  {"Ruth", 608},
	"1-sam": {"1 Samuel", 609},
	"2-sam": {"2 Samuel", 610},
	"1-kgs": {"1 Kings", 611},
	"2-kgs": {"2 Kings", 612},
	"1-chr": {"1 Chronicles", 613},
	"2-chr": {"2 Chronicles", 614},
	"ezra":  {"Ezra", 615},
	"neh":   {"Nehemiah", 616},
	"esth":  {"Esther", 617},
	"job":   {"Job", 618},
	"ps":    {"Psalms", 619},
	"prov":  {"Proverbs", 620},
	"eccl":  {"Ecclesiastes", 621},
	"song":  {"Song of Solomon", 622},
	"isa":   {"Isaiah", 623},
	"jer":   {"Jeremiah", 624},
	"lam":   {"Lamentations", 625},
	"ezek":  {"Ezekiel", 626},
	"dan":   {"Daniel", 627},
	"hosea": {"Hosea", 628},
	"joel":  {"Joel", 629},
	"amos":  {"Amos", 630},
	"obad":  {"Obadiah", 631},
	"jonah": {"Jonah", 632},
	"micah": {"Micah", 633},
	"nahum": {"Nahum", 634},
	"hab":   {"Habakkuk", 635},
	"zeph":  {"Zephaniah", 636},
	"hag":   {"Haggai", 637},
	"zech":  {"Zechariah", 638},
	"mal":   {"Malachi", 639},
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
	//simpleNotes = filterSimpleNotes("Faith", simpleNotes)
	orderedNotes := orderSimpleNotes(simpleNotes)

	w.WriteString(`<!DOCTYPE html>
<html>
<head>
<title>Rob's notes</title>
</head>
<body>
<table border="1" cellPadding="5" cellSpacing="0">
<tr>
	<td width="50%">Reference</td>
	<td width="50%">Note</td>
</tr>`)
	for _, note := range orderedNotes {
		w.WriteString(fmt.Sprintf(`
<tr>
	<td>%s</td>
	<td>%s</td>
</tr>`+"\n", formatReferences(note.References), note.Notes))
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

func filterSimpleNotes(tag string, notes []simpleNote) []simpleNote {
	for i := 0; i < len(notes); i++ {
		note := &notes[i]
		if !hasTag(note.Tags, tag) {
			notes = append(notes[:i], notes[i+1:]...)
			i--
		}
	}
	return notes
}

func hasTag(tags []string, find string) bool {
	for _, tag := range tags {
		if tag == find {
			return true
		}
	}
	return false
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
			buf.WriteString("<br/><br/>")
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
		content := strings.ReplaceAll(strings.ReplaceAll(note.Note.Content, "\n", "<br/>"), "\\", "")
		notes = append(notes, *newSimpleNote(uris, content, note.Tags, db))
	}
	return notes
}

func newSimpleNote(uris []string, notes string, tags []string, db pgx.PGXer) *simpleNote {
	note := &simpleNote{Notes: notes, Tags: tags}
	for i, uri := range uris {
		reference, order := getReference(uri, db)
		note.References = append(note.References, *reference)
		if i == 0 {
			note.Order = order
		}
	}
	return note
}

func getReference(uri string, db pgx.PGXer) (*reference, int) {
	refs := strings.Split(uri, "/")
	if len(refs) > 1 {
		switch refs[1] {
		case "scriptures":
			return newScriptureReference(refs[2:], db)
		case "general-conference", "ensign":
			return newConferenceReference(refs[2:])
		case "manual":
			return newManualReference(refs[2:])
		}
	}
	return &reference{FullReference: uri, ShortReference: uri}, 0
}

func newScriptureReference(refs []string, db pgx.PGXer) (*reference, int) {
	order := 100 // BOM starts at 100
	volume, book, chapter, verse := parseRefs(refs)
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
	order = order*1000000 + chapterInt*1000 + verseInt // order starts at 100 million for BOM intro
	return &reference{fmt.Sprintf("%s %s:%s", book, chapter, verse), verse, text}, order
}

func newConferenceReference(refs []string) (*reference, int) {
	year, month, title, _ := parseRefs(refs)
	yearInt, _ := strconv.Atoi(year)
	monthInt, _ := strconv.Atoi(month)
	order := 1000000 + yearInt*100 + monthInt
	return &reference{fmt.Sprintf("%s (%s %s)", formatTitle(title), formatMonth(monthInt), year), "", ""}, order
}

func newManualReference(refs []string) (*reference, int) {
	volume, page, _, paragraph := parseRefs(refs)
	order := 2000000
	return &reference{fmt.Sprintf("%s (%s:%s)", formatTitle(volume), page, getVerse(paragraph)), "", ""}, order
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

func formatMonth(month int) string {
	switch month {
	case 1:
		return "Jan"
	case 2:
		return "Feb"
	case 3:
		return "Mar"
	case 4:
		return "Apr"
	case 5:
		return "May"
	case 6:
		return "Jun"
	case 7:
		return "Jul"
	case 8:
		return "Aug"
	case 9:
		return "Sep"
	case 10:
		return "Oct"
	case 11:
		return "Nov"
	case 12:
		return "Dec"
	default:
		return strconv.Itoa(month)
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
