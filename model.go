package w2v

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"sort"
	"strconv"
	"strings"
)

// Vector is struct hold word and vector.
type Vector struct {
	word  string
	vec   []float64
	elems []string
}

// String return string that have word and slice of vector.
func (v *Vector) String() string {
	return fmt.Sprintf("%s%v", v.word, v.vec)
}

// Word return word.
func (v *Vector) Word() string {
	return v.word
}

// Vector return vector.
func (v *Vector) Vector() []float64 {
	return v.vec
}

// Normalize return normalized vector.
func (v *Vector) Normalize() *Vector {
	w := snrm2(len(v.vec), v.vec)
	sscal(len(v.vec), 1/w, v.vec)
	return v
}

// Model is slice of Vector.
type Model []*Vector

// LoadBinary load model binary-file generated word2vec.
func LoadBinary(r io.Reader, bitsize int) (Model, error) {
	br := bufio.NewReader(r)
	result := Model{}

	b, _, err := br.ReadLine()
	if err != nil {
		return nil, errors.New("bad format")
	}
	token := strings.Split(string(b), " ")
	if len(token) != 2 {
		return nil, errors.New("bad format")
	}
	nsize, err := strconv.ParseInt(token[1], 10, 64)
	if err != nil {
		return nil, errors.New("bad format")
	}

	for {
		s, err := br.ReadString(' ')
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		s = s[:len(s)-1]
		vec := []float64{}
		for i := 0; i < int(nsize); i++ {
			if bitsize == 32 {
				var val float32
				err = binary.Read(br, binary.LittleEndian, &val)
				if err != nil {
					return nil, err
				}
				vec = append(vec, float64(val))
			} else {
				var val float64
				err = binary.Read(br, binary.LittleEndian, &val)
				if err != nil {
					return nil, err
				}
				vec = append(vec, float64(val))
			}
		}
		b, err := br.ReadByte()
		if err != nil {
			return nil, err
		}
		if b != '\n' {
			return nil, errors.New("bad format")
		}
		result = append(result, &Vector{
			word:  s,
			vec:   vec,
			elems: []string{s},
		})

	}
	return result, nil
}

// LoadText load model text-file generated word2vec.
func LoadText(r io.Reader) (Model, error) {
	scanner := bufio.NewScanner(r)
	result := Model{}
	for scanner.Scan() {
		token := strings.Split(strings.TrimSpace(scanner.Text()), " ")
		if len(token) < 3 || token[0] == "</s>" {
			continue
		}
		vec := []float64{}
		for i := 1; i < len(token); i++ {
			val, err := strconv.ParseFloat(token[i], 64)
			if err != nil {
				return nil, err
			}
			vec = append(vec, val)
		}
		result = append(result, &Vector{
			word:  token[0],
			vec:   vec,
			elems: []string{token[0]},
		})

	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

// Find return Vector matched with word.
func (m Model) Find(word string) *Vector {
	for _, vector := range m {
		if vector.word == word {
			return vector
		}
	}
	return nil
}

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}

func max(i, j int) int {
	if i < j {
		return j
	}
	return i
}

// Add return new Vector added vector given.
func (v *Vector) Add(rhs *Vector) *Vector {
	if v == nil {
		return rhs
	}
	if rhs == nil {
		return v
	}
	l := min(len(v.vec), len(rhs.vec))
	vec := make([]float64, l)
	copy(vec, v.vec)
	saxpy(l, 1, rhs.vec, 1, vec, 1)
	elems := make([]string, len(v.elems)+len(rhs.elems))
	elems = append(elems, rhs.elems...)
	elems = append(elems, v.elems...)
	return &Vector{
		word:  v.word + " + " + rhs.word,
		vec:   vec,
		elems: elems,
	}
}

// Sub return new Vector subtracted vector given.
func (v *Vector) Sub(rhs *Vector) *Vector {
	if rhs == nil {
		return v
	}
	if v == nil {
		v = &Vector{
			word:  "",
			vec:   make([]float64, len(rhs.vec)),
			elems: nil,
		}
	}

	l := min(len(v.vec), len(rhs.vec))
	vec := make([]float64, l)
	copy(vec, v.vec)
	saxpy(l, -1, rhs.vec, 1, vec, 1)
	elems := make([]string, len(v.elems)+len(rhs.elems))
	elems = append(elems, rhs.elems...)
	elems = append(elems, v.elems...)
	return &Vector{
		word:  v.word + " - " + rhs.word,
		vec:   vec,
		elems: elems,
	}
}

// Distance return distance between v and rhs.
func (v *Vector) Distance(rhs *Vector) float64 {
	distance := float64(0)
	for i := 0; i < len(rhs.vec); i++ {
		distance += math.Pow(float64(rhs.vec[i]-v.vec[i]), float64(2))
	}
	return float64(math.Sqrt(distance))
}

// Cosine return how rhs is similar to v.
func (v *Vector) Cosine(rhs *Vector) float64 {
	n := 0
	llhs := len(v.vec)
	lrhs := len(rhs.vec)
	n = max(llhs, lrhs)
	svec, slhs, srhs := 0.0, 0.0, 0.0
	for i := 0; i < n; i++ {
		if i >= llhs {
			srhs += math.Pow(rhs.vec[i], 2.0)
		} else if i >= lrhs {
			slhs += math.Pow(v.vec[i], 2.0)
		} else {
			svec += v.vec[i] * rhs.vec[i]
			slhs += math.Pow(v.vec[i], 2.0)
			srhs += math.Pow(rhs.vec[i], 2.0)
		}
	}
	if slhs == 0 || srhs == 0 {
		return 0.0
	}
	return svec / (math.Sqrt(slhs) * math.Sqrt(srhs))
}

// Entry is struct used by CosineSimilars and Neighbourhood.
type Entry struct {
	Value  float64
	Vector *Vector
}

// String return string that have Vector and Value.
func (e *Entry) String() string {
	return fmt.Sprintf("%v(%v)", e.Vector, e.Value)
}

func (m Model) order(vector *Vector, calc func(*Vector) float64, less func(*Entry, *Entry) bool) []*Entry {
	if len(m) == 0 || vector == nil {
		return nil
	}

	entries := []*Entry{}
	for _, lhs := range m {
		if lhs.word == vector.word {
			continue
		}
		found := false
		for _, v := range vector.elems {
			if lhs.word == v {
				found = true
				break
			}
		}
		if found {
			continue
		}
		entries = append(entries, &Entry{
			Vector: lhs,
			Value:  calc(lhs),
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		return less(entries[i], entries[j])
	})
	return entries
}

// Neighbourhood return entries order by how near by vector.
func (m Model) Neighbourhood(vector *Vector) []*Entry {
	return m.order(vector, vector.Distance, func(lhs, rhs *Entry) bool {
		return lhs.Value < rhs.Value
	})
}

// CosineSimilars return entries order by how similar to vector.
func (m Model) CosineSimilars(vector *Vector) []*Entry {
	return m.order(vector, vector.Cosine, func(lhs, rhs *Entry) bool {
		return lhs.Value > rhs.Value
	})
}

func snrm2(n int, x []float64) float64 {
	var a, b, c, d float64
	var xi int
	for ; n >= 4; n -= 4 {
		a += x[xi] * x[xi]
		xi++
		b += x[xi] * x[xi]
		xi++
		c += x[xi] * x[xi]
		xi++
		d += x[xi] * x[xi]
		xi++
	}
	for ; n > 0; n-- {
		a += x[xi] * x[xi]
		xi++
	}
	return math.Sqrt(a + b + c + d)
}

func sscal(n int, a float64, x []float64) {
	var xi int
	for ; n >= 2; n -= 2 {
		x[xi] = a * x[xi]
		xi++
		x[xi] = a * x[xi]
		xi++
	}
	if n != 0 {
		x[xi] = a * x[xi]
	}
}

func saxpy(n int, a float64, x []float64, dx int, y []float64, dy int) {
	var xi, yi int
	if a > 0 {
		for ; n >= 2; n -= 2 {
			y[yi] += x[xi]
			xi += dx
			yi += dy

			y[yi] += x[xi]
			xi += dx
			yi += dy
		}
		if n != 0 {
			y[yi] += a * x[xi]
		}
	} else {
		for ; n >= 2; n -= 2 {
			y[yi] -= x[xi]
			xi += dx
			yi += dy

			y[yi] -= x[xi]
			xi += dx
			yi += dy
		}
		if n != 0 {
			y[yi] -= x[xi]
		}
	}
}
