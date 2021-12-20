package main

import (
    "testing"
    "regexp"
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
        t.Fatalf(`findPhrase("") = %q, want "", error`, word, )
    }
}

func TestErrorPhraseNoEnd(t *testing.T) {
    phrase := "abcdefghijklmnopqrstuvwxyz"
    // want := ""
    word, err := findPhrase(phrase, "l", "1")
    if word != "" || err == nil {
        t.Fatalf(`findPhrase("") = %q, want "", error`, word, )
    }
}

func TestGetWotd(t *testing.T) {
    url := "https://www.merriam-webster.com/word-of-the-day/"
    pageContent, err := getWotd(url)
    word, err2 := findPhrase(pageContent, "<h1>", "</h1>")
    if word == "" || err != nil || err2 != nil {
        t.Fatalf(`getWotd()  = %q, want match for %#q, nil`, word, "")
    }
}

func TestGetWotdWithDate(t *testing.T) {
    url := "https://www.merriam-webster.com/word-of-the-day/2021-01-01"
    want := "reprise"
    pageContent, err := getWotd(url)
    word, err2 := findPhrase(pageContent, "<h1>", "</h1>")
    if word != want || err != nil || err2 != nil {
        t.Fatalf(`getWotd()  = %q, want %#q, nil`, word, want)
    }
}

func TestErrorGetWotdBadURL(t *testing.T) {
    url := "https://uuu.merriam-webster.com/word-of-the-day/2000-01-01"
    pageContent, err := getWotd(url)
    if pageContent != "" || err == nil {
        t.Fatalf(`getWotd() (bad url) = %q, want match for %#q, nil`, pageContent, "")
    }
}

func TestErrorGetWotdBadFormat(t *testing.T) {
    url := "https://www.merriam-webster.com/word-of-the-day/2000-01-01"
    pageContent, err := getWotd(url)
    word, err := findPhrase(pageContent, "<h1>", "</h1>")
    if word != "" || err == nil {
        t.Fatalf(`getWotd() (bad web content) = %q, want match for %#q, nil`, word, "")
    }
}

func TestMain(t *testing.T) {
    // This is a bad test, it only checks if it can run without erroring, not if the Printf is accurate
    main()
}
