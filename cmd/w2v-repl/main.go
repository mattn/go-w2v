package main

//go:generate golex -o tokenizer.go tokenizer.l
//go:generate goyacc -o expr.go expr.y

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/fatih/color"
	"github.com/ikawaha/kagome/tokenizer"
	"github.com/mattn/go-w2v"
)

var re = regexp.MustCompile(`[\[\]「」『』()（）、。*:]`)

var ignoreType = []string{
	"フィラー",
	"記号",
	"動詞",
	"感動詞",
	"助詞",
	"助動詞",
	"副詞",
	"接続詞",
	"連体詞",
}

func ignore(features []string) bool {
	if len(features) < 7 {
		return true
	}
	for _, v := range ignoreType {
		if features[0] == v {
			return true
		}
	}
	if features[1] == "接尾" || features[1] == "非自立" {
		return true
	}
	if features[0] == "名詞" && features[1] == "数" {
		return true
	}
	return false
}

func filter(model *w2v.Model) {
	var n w2v.Model
	t := tokenizer.New()
	for _, v := range *model {
		word := v.Word()
		if re.MatchString(word) {
			continue
		}
		tokens := t.Tokenize(word)
		bad := false
		for _, token := range tokens {
			features := token.Features()
			if len(features) == 0 || features[0] == "BOS" || features[0] == "EOS" {
				continue
			}
			if ignore(features) {
				bad = true
				break
			}
		}
		if bad {
			continue
		}
		n = append(n, v)
	}
	*model = n
}

func candidates(model w2v.Model, vector *w2v.Vector, count int) {
	var elems []*w2v.Entry

	elems = model.CosineSimilars(vector)
	n := 0
	for _, elem := range elems {
		word := elem.Vector.Word()
		line := color.YellowString("%v", word) + " " + color.GreenString("%f", elem.Value)
		fmt.Fprintln(color.Output, line)
		if n++; n == count {
			break
		}
	}
}

func run() int {
	var filename, query string
	var binfmt, count int
	flag.IntVar(&count, "n", 4, "number of results")
	flag.StringVar(&query, "q", "", "query")
	flag.StringVar(&filename, "f", "data.model", "word2vec model file")
	flag.IntVar(&binfmt, "b", 0, "binary format (0:text, 32/64: binary")
	flag.Parse()

	switch binfmt {
	case 0, 32, 64:
	default:
		flag.Usage()
		return 1
	}

	f, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: %v\n", os.Args[0], err)
		return 1
	}
	defer f.Close()

	var model w2v.Model
	switch binfmt {
	case 0:
		model, err = w2v.LoadText(f)
	case 32, 64:
		model, err = w2v.LoadBinary(f, binfmt)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: %v\n", os.Args[0], err)
		return 1
	}

	filter(&model)

	eval := newEval(model)
	if query != "" {
		vector, err := eval.Do(query)
		if err != nil {
			color.Red("%v\n", err)
			return 1
		}
		vector.Normalize()
		n := 0
		for _, elem := range model.CosineSimilars(vector) {
			fmt.Println(elem.Vector.Word())
			if n++; n > count {
				break
			}
		}
		return 0
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Fprint(color.Output, color.GreenString("> "))
		if !scanner.Scan() {
			break
		}
		vector, err := eval.Do(scanner.Text())
		if err != nil {
			color.Red("%v\n", err)
			continue
		}
		vector.Normalize()
		candidates(model, vector, count)
	}
	return 0
}

func main() {
	os.Exit(run())
}
