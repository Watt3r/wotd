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

func findPhrase(pageContent string, start string, end string) []byte {
    wordIndex := strings.Index(pageContent, start)
    if wordIndex == -1 {
      log.Fatal("Word not found")
    }
    wordIndex += len(start)
    wordEndIndex := strings.Index(pageContent[wordIndex:], end) + wordIndex
    word := []byte(pageContent[wordIndex:wordEndIndex])
    return word
}
func main() {
    // Color output
    blue := color.New(color.FgBlue).Add(color.Bold).SprintFunc()
    green := color.New(color.FgGreen).Add(color.Bold).SprintFunc()

    // Get Word of the Day from Merriam-Webster
    resp, err := http.Get("https://www.merriam-webster.com/word-of-the-day/")
    if err != nil {
      log.Fatal(err)
    }
    defer resp.Body.Close()
    data, err := ioutil.ReadAll(resp.Body)

    // Simple parse HTML 
    pageContent := string(data)
    word := findPhrase(pageContent, "<h1>", "</h1>")
    def := findPhrase(pageContent, "<p>", "</p>")

    // Print result
    fmt.Printf("Word of the Day: %s\nDefinition: %s\n",  blue(fmt.Sprintf("%s", word)), green(regexp.MustCompile("</?em>").ReplaceAllString(fmt.Sprintf("%s", def), "")))
}
