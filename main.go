package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/EndFirstCorp/onedb"
	"github.com/EndFirstCorp/onedb/pgx"
)

var bom = map[string]string{
	"introduction": "Introduction to the Book of Mormon",
	"three":        "Testimony of Three Witnesses",
	"1-ne":         "1 Nephi",
	"2-ne":         "2 Nephi",
	"jacob":        "Jacob",
	"enos":         "Enos",
	"jarom":        "Jarom",
	"omni":         "Omni",
	"w-of-m":       "Words of Mormon",
	"mosiah":       "Mosiah",
	"alma":         "Alma",
	"hel":          "Helaman",
	"3-ne":         "3 Nephi",
	"4-ne":         "4 Nephi",
	"morm":         "Mormon",
	"ether":        "Ether",
	"moro":         "Moroni",
}

var nt = map[string]string{
	"matt":   "Matthew",
	"mark":   "Mark",
	"luke":   "Luke",
	"john":   "John",
	"acts":   "Acts",
	"rom":    "Romans",
	"1-cor":  "1 Corinthians",
	"2-cor":  "2 Corinthians",
	"gal":    "Galatians",
	"eph":    "Ephesians",
	"philip": "Philippians",
	"col":    "Colossians",
	"1-thes": "1 Thessalonians",
	"2-thes": "2 Thessalonians",
	"1-tim":  "1 Timothy",
	"2-tim":  "2 Timothy",
	"titus":  "Titus",
	"philem": "Philemon",
	"heb":    "Hebrews",
	"james":  "James",
	"1-pet":  "1 Peter",
	"2-pet":  "2 Peter",
	"1-jn":   "1 John",
	"2-jn":   "2 John",
	"3-jn":   "3 John",
	"jude":   "Jude",
	"rev":    "Revelation",
}

var ot = map[string]string{
	"gen":   "Genesis",
	"ex":    "Exodus",
	"lev":   "Leviticus",
	"num":   "Numbers",
	"deut":  "Deuteronomy",
	"josh":  "Joshua",
	"judg":  "Judges",
	"ruth":  "Ruth",
	"1-sam": "1 Samuel",
	"2-sam": "2 Samuel",
	"1-kgs": "1 Kings",
	"2-kgs": "2 Kings",
	"1-chr": "1 Chronicles",
	"2-chr": "2 Chronicles",
	"ezra":  "Ezra",
	"neh":   "Nehemiah",
	"esth":  "Esther",
	"job":   "Job",
	"ps":    "Psalms",
	"prov":  "Proverbs",
	"eccl":  "Ecclesiastes",
	"song":  "Song of Solomon",
	"isa":   "Isaiah",
	"jer":   "Jeremiah",
	"lam":   "Lamentations",
	"ezek":  "Ezekial",
	"dan":   "Daniel",
	"hosea": "Hosea",
	"joel":  "Joel",
	"amos":  "Amos",
	"obad":  "Obadiah",
	"jonah": "Jonah",
	"micah": "Micah",
	"nahum": "Nahum",
	"hab":   "Habakkuk",
	"zeph":  "Zephaniah",
	"hag":   "Haggai",
	"zech":  "Zechariah",
	"mal":   "Malachi",
}

var pgp = map[string]string{
	"introduction": "Introduction to the Pearl of Great Price",
	"moses":        "Moses",
	"abr":          "Abraham",
	"js-m":         "Joseph Smith--Matthew",
	"js-h":         "Joseph Smith--History",
	"a-of-f":       "Articles of Faith",
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
	w.WriteString(`<!DOCTYPE html>
<html>
<head>
<title>Rob's notes</title>
</head>
<body>
<table>
<tr>
	<td width="100">Tags</td>
	<td width="50%">Reference</td>
	<td width="50%">Note</td>
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
</tr>`+"\n", note.Tags, getReference(uri, db), note.Note.Content))
	}
	w.WriteString(`
</table></body></html>`)
}

const query = `SELECT v.scripture_text 
FROM verses v 
JOIN chapters c ON v.chapter_id = c.id 
JOIN books b ON c.book_id = b.id 
JOIN volumes vo ON b.volume_id = vo.id 
WHERE b.book_title = $1 and c.chapter_number = $2 and v.verse_number = $3;`

func getText(db pgx.PGXer, book string, chapter, verse int) (string, error) {
	var text string
	return text, db.QueryValues(onedb.NewQuery(query, book, chapter, verse), &text)
}

func getReference(ref string, db pgx.PGXer) string {
	refs := strings.Split(ref, "/")
	if len(refs) > 1 {
		switch refs[1] {
		case "scriptures":
			volume, book, chapter, verse := parseRefs(refs[2:])
			if newBook := getBook(volume, book); newBook != "" {
				book = newBook
			}
			verse = getVerse(verse)
			verseInt, _ := strconv.Atoi(verse)
			chapterInt, _ := strconv.Atoi(chapter)
			text, err := getText(db, book, chapterInt, verseInt)
			if err != nil {
				fmt.Println(book, chapterInt, verseInt, err)
			}
			return fmt.Sprintf("%s %s:%s\n%s", book, chapter, verse, text)
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
		return bom[book]
	case "dc-testament":
		return "Doctrine and Covenants"
	case "nt":
		return nt[book]
	case "ot":
		return ot[book]
	case "pgp":
		return pgp[book]
	}
	return strings.Title(book)
}

func getVerse(paragraph string) string {
	if len(paragraph) == 0 || paragraph[0] != 'p' {
		return paragraph
	}
	return paragraph[1:]
}
