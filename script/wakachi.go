package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ikawaha/kagome/tokenizer"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	t := tokenizer.New()
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}
		tokens := t.Tokenize(text)
		n := 0
		for _, token := range tokens {
			features := token.Features()
			if len(features) == 0 || token.Surface == "BOS" || token.Surface == "EOS" {
				continue
			}
			if n > 0 {
				fmt.Print(" ")
			}
			if features[6] == "*" {
				fmt.Print(token.Surface)
			} else {
				fmt.Print(features[6])
			}
			n++
		}
		fmt.Println()
	}
	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}
}
