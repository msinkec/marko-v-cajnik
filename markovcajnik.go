package main


import (
    "fmt"
    "os"
    "strings"
    "unicode"
    "io/ioutil"

    "github.com/jmcvetta/randutil"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func readCorpus(files []string) string {
    var res strings.Builder

    for i := range files {
        dat, err := ioutil.ReadFile(files[i])
        check(err)
        res.Write(dat)
    }

    return res.String()
}

func getTrailingWords(target string, words []string) []string {
    res := make([]string, 0, 5)

    for i := range words {
        word := words[i]
        if word == target && i < len(words) - 1 {
           res = append(res, words[i+1])
        }
    }

    return res
}

func tokenizeText(text string) []string {
    var sb strings.Builder

    for _, r := range text {
        if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
            sb.WriteRune(' ')
        }
        sb.WriteRune(r)
    }

    return strings.Fields(sb.String())
}

func buildDict(text string) map[string]map[string]int {
    res := make(map[string]map[string]int)

    words := tokenizeText(text)
    for i := range words {
        word := words[i]
        res[word] = make(map[string]int)
        trailingWords := getTrailingWords(word, words)
        for j := range trailingWords {
            trailingWord := trailingWords[j]
            res[word][trailingWord]++
        }
    }

    return res
}

func generateSentence(dict map[string]map[string]int) string {
    var sb strings.Builder

    // Start with a random word
    var word string
    for word = range dict {
        if word != "." {
            break
        }
    }
    sb.WriteString(word)

    for true {
        all := dict[word]

        // Choose subsequent word with weighted probabilities
        choices := make([]randutil.Choice, len(all))
        for key, count := range all {
            choices = append(choices, randutil.Choice{count, key})
        }

        choice, err := randutil.WeightedChoice(choices)
        check(err)

        word = choice.Item.(string)

        if word == "." {
            sb.WriteRune('.')
            break
        }

        sb.WriteRune(' ')
        sb.WriteString(word)
    }

    return sb.String()
}

func main() {
    argsWithoutProg := os.Args[1:]

    if len(argsWithoutProg) == 0 {
        fmt.Println("Too few input arguments. Expected 1 or more.")
        os.Exit(1)
    }

    text := readCorpus(argsWithoutProg)
    dict := buildDict(text)


    for i := 0; i < 15; i++ {
        gen := generateSentence(dict)
        fmt.Println(gen)
    }
}

