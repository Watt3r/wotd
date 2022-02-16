package main

import (
	"regexp"
	"testing"
	"time"
)

func TestFindPhrase(t *testing.T) {
	phrase := "abcdefghijklmnopqrstuvwxyz"
	want := regexp.MustCompile("lmno")
	word, err := findPhrase(phrase, "k", "p")
	if !want.MatchString(word) || err != nil {
		t.Fatalf(`findPhrase("lmno") = %q, want match for %#q, nil`, word, want)
	}
}

func TestErrorPhraseNoStart(t *testing.T) {
	phrase := "abcdefghijklmnopqrstuvwxyz"
	// want := ""
	word, err := findPhrase(phrase, "1", "a")
	if word != "" || err == nil {
		t.Fatalf(`findPhrase("") = %q, want "", error`, word)
	}
}

func TestErrorPhraseNoEnd(t *testing.T) {
	phrase := "abcdefghijklmnopqrstuvwxyz"
	// want := ""
	word, err := findPhrase(phrase, "l", "1")
	if word != "" || err == nil {
		t.Fatalf(`findPhrase("") = %q, want "", error`, word)
	}
}

func TestColor(t *testing.T) {
	phrase := "abcdefghijklmnopqrstuvwxyz"
	want := "\u001b[1m\u001b[36mabcdefghijklmnopqrstuvwxyz\033[0m"
	colorOutput := color(phrase, "\u001b[36m", false)
	t.Logf("\n\n\n" + want + "\n" + colorOutput + "\n")
	if colorOutput != want {
		t.Fatalf(`color("abc", "ansi", false) = %q, want match for %#q, nil`, colorOutput, want)
	}
}

func TestColorNoColor(t *testing.T) {
	phrase := "abcdefghijklmnopqrstuvwxyz"
	word := color(phrase, "legitanythingthisdoesntmatter", true)
	if word != phrase {
		t.Fatalf(`color("abc", "doesntmatter", true) = %q, want "", error`, word)
	}
}

func TestGetWotd(t *testing.T) {
	url := "https://www.merriam-webster.com/word-of-the-day/"
	date := ""
	pageContent, err := getWotd(url, &date)
	word, err2 := findPhrase(pageContent, "<h1>", "</h1>")
	if word == "" || err != nil || err2 != nil {
		t.Fatalf(`getWotd()  = %q, want match for %#q, nil`, word, "")
	}
}

func TestGetWotdWithDate(t *testing.T) {
	url := "https://www.merriam-webster.com/word-of-the-day/"
	wantWord := "benign"
	wantDef := "<em>Benign</em> means \"not causing harm or injury.\" In medicine, it refers to tumors that are not cancerous."
	date := "2021-12-22"
	pageContent, err := getWotd(url, &date)
	word, err2 := findPhrase(pageContent, "<h1>", "</h1>")
	def, err3 := findPhrase(pageContent, "<p>", "</p>")
	if word != wantWord || def != wantDef || err != nil || err2 != nil || err3 != nil {
		t.Fatalf(`getWotd()  = %q, want %#q, nil`, word, wantWord)
	}
}

func TestGetWotdWithBadFormatDate(t *testing.T) {
	url := "https://www.merriam-webster.com/word-of-the-day/"
	date := "2021-12-32"
	pageContent, err := getWotd(url, &date)
	if err == nil || pageContent != "" {
		t.Fatalf(`getWotd()  = %q, want %#q, nil`, pageContent, "")
	}
}

func TestGetWotdWithBadFormatDate2(t *testing.T) {
	url := "https://www.merriam-webster.com/word-of-the-day/"
	date := "1-1-2021"
	pageContent, err := getWotd(url, &date)
	if err == nil || pageContent != "" {
		t.Fatalf(`getWotd()  = %q, want %#q, nil`, pageContent, "")
	}
}

func TestGetWotdWithBadFutureDate(t *testing.T) {
	url := "https://www.merriam-webster.com/word-of-the-day/"
	date := time.Now().AddDate(0, 1, 0).Format("2006-01-02")
	pageContent, err := getWotd(url, &date)
	if err == nil || pageContent != "" {
		t.Fatalf(`getWotd()  = %q, want %#q, nil`, pageContent, "")
	}
}

func TestErrorGetWotdBadURL(t *testing.T) {
	url := "https://uuu.merriam-webster.com/word-of-the-day/"
	date := ""
	pageContent, err := getWotd(url, &date)
	if pageContent != "" || err == nil {
		t.Fatalf(`getWotd() (bad url) = %q, want match for %#q, nil`, pageContent, "")
	}
}

func TestErrorGetWotdBadFormat(t *testing.T) {
	url := "https://www.merriam-webster.com/word-of-the-day/2000-01-01"
	date := ""
	pageContent, err := getWotd(url, &date)
	word, err := findPhrase(pageContent, "<h1>", "</h1>")
	if word != "" || err == nil {
		t.Fatalf(`getWotd() (bad web content) = %q, want match for %#q, nil`, word, "")
	}
}

func TestGetSourceWebster(t *testing.T) {
	source := "merriam"
	date := ""
	url, wordStart, defStart := getSource(&source, &date)
	if url != "https://www.merriam-webster.com/word-of-the-day/" || wordStart != "<h1>" || defStart != "<p>" {
		t.Fatalf(`getSource("merriam") = %q, %q, %q`, url, defStart, defStart)
	}
}

func TestGetSourceDictionary(t *testing.T) {
	source := "dictionary"
	date := ""
	url, wordStart, defStart := getSource(&source, &date)
	if url != "https://www.dictionary.com/e/word-of-the-day/" || wordStart != "<h1 class=\"js-fit-text\" style=\"color: #00248B\">" || defStart != "</p>\n\n                \n                <p>" {
		t.Fatalf(`getSource("dictionary") = %q, %q, %q`, url, defStart, defStart)
	}
}

func TestMain(t *testing.T) {
	// This is a bad test, it only checks if it can run without erroring, not if the Printf is accurate
	main()
}
