package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func findPhrase(pageContent string, start string, end string) (string, error) {
	wordIndex := strings.Index(pageContent, start)
	// Error if start not found
	if wordIndex == -1 {
		return "", fmt.Errorf("could not find phrase")
	}
	wordIndex += len(start)
	wordEndIndex := strings.Index(pageContent[wordIndex:], end) + wordIndex
	// Error if end not found
	if wordIndex > wordEndIndex {
		return "", fmt.Errorf("could not find phrase")
	}
	word := []byte(pageContent[wordIndex:wordEndIndex])
	strWord := fmt.Sprintf("%s", word)
	return strWord, nil
}

func getWotd(url string, date *string) (string, error) {
	// Validate date
	if *date != "" {
		parsedDate, err := time.Parse("2006-01-02", *date)
		if err != nil || parsedDate.After(time.Now()) {
			return "", fmt.Errorf("invalid date format")
		}
	}

	resp, err := http.Get(url + *date)
	if err != nil {
		return "", fmt.Errorf("error with getting Wotd from url")
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)

	// Simple parse HTML
	pageContent := string(data)
	return pageContent, nil
}

func color(word string, color string, noColor bool) string {
	if noColor {
		return word
	}

	// Bold ansi + color ansi + word + reset ansi
	return "\u001b[1m" + color + word + "\033[0m"
}

func main() {
	// Get specific date if specified
	date := flag.String("date", "", "Optional date of Word of the Day, YYYY-MM-DD")
	flag.StringVar(date, "d", "", "Optional date of Word of the Day, YYYY-MM-DD")
	noColor := flag.Bool("no-color", false, "Disables color output")
	flag.Parse()

	// Get Word of the Day from Merriam-Webster
	pageContent, err := getWotd("https://www.merriam-webster.com/word-of-the-day/", date)
	if err != nil {
		log.Fatal(err)
	}
	word, err := findPhrase(pageContent, "<h1>", "</h1>")
	def, err := findPhrase(pageContent, "<p>", "</p>")
	def = regexp.MustCompile("</?em>").ReplaceAllString(def, "")

	// Print result
	blue := "\u001b[36m"
	green := "\u001b[32m"
	fmt.Printf("Word of the Day: %s\nDefinition: %s\n", color(word, blue, *noColor), color(def, green, *noColor))
}
