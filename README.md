# w2v

word2vec library for japanese written in Go.

## Usage

```go
m, err := w2v.LoadText("data.model")
if err != nil {
	log.Fatal(err)
}
m.Find("Go").Add("Language").CosineSimilars()
```

## Installation

```
$ go get github.com/mattn/go-w2v
```

## License

MIT

## Author

Yasuhiro Matsumoto (a.k.a. mattn)
