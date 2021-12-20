package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "log"
    "strings"
    "regexp"
    "github.com/fatih/color"
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

func getWotd(url string) (string, error) {
    resp, err := http.Get(url)
    if err != nil {
      return "", fmt.Errorf("error with getting Wotd from url")
    }
    defer resp.Body.Close()
    data, err := ioutil.ReadAll(resp.Body)

    // Simple parse HTML 
    pageContent := string(data)
    return pageContent, nil
}

func main() {
    // Color output
    blue := color.New(color.FgBlue).Add(color.Bold).SprintFunc()
    green := color.New(color.FgGreen).Add(color.Bold).SprintFunc()

    // Get Word of the Day from Merriam-Webster
    pageContent, err := getWotd("https://www.merriam-webster.com/word-of-the-day/")
    if err != nil {
      log.Fatal("Wotd could not load URL, are you connected to the internet?")
    }
    word, err := findPhrase(pageContent, "<h1>", "</h1>")
    def, err := findPhrase(pageContent, "<p>", "</p>")

    // Print result
    fmt.Printf("Word of the Day: %s\nDefinition: %s\n",  blue(word), green(regexp.MustCompile("</?em>").ReplaceAllString(def, "")))
}
