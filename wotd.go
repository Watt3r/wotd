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
	source := flag.String("source", "", "Optional source to get Word of the Day, Allowed values: merriam, dictionary")
	flag.StringVar(source, "s", "", "Optional source to get Word of the Day, Allowed values: merriam, dictionary")
	noColor := flag.Bool("no-color", false, "Disables color output")
	flag.Parse()

	// Get Word of the Day from Merriam-Webster
	url := ""
        wordStart := ""
        defStart := ""
        if (*source == "dictionary") {
          url = "https://www.dictionary.com/e/word-of-the-day/"
          wordStart =  "<h1 class=\"js-fit-text\" style=\"color: #00248B\">"
          // Dictionary.com doesn't use any class or id on the definition <p> element, this is the best way to isolate
          defStart = "</p>\n\n                \n                <p>"
          // Dictionary.com doesn't have a searchable archive, manually ensure date is set to today
          *date = ""
        } else {
          url = "https://www.merriam-webster.com/word-of-the-day/"
          wordStart = "<h1>"
          defStart = "<p>"
        }

	pageContent, err := getWotd(url, date)
	if err != nil {
		log.Fatal(err)
	}
	word, err := findPhrase(pageContent, wordStart, "</h1>")
	def, err := findPhrase(pageContent, defStart, "</p>")
	def = regexp.MustCompile("</?em>").ReplaceAllString(def, "")

	// Print result
	blue := "\u001b[36m"
	green := "\u001b[32m"
	fmt.Printf("Word of the Day: %s\nDefinition: %s\n", color(word, blue, *noColor), color(def, green, *noColor))
}
