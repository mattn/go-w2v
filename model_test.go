package w2v

import (
	"math"
	"os"
	"reflect"
	"testing"
)

func TestVectorString(t *testing.T) {
	want := "foo[1 2 3]"
	got := (&Vector{
		word: "foo",
		vec:  []float64{1, 2, 3},
	}).String()
	if got != want {
		t.Fatalf("want %v, but %v:", want, got)
	}
}

func TestVectorWord(t *testing.T) {
	want := "foo"
	got := (&Vector{
		word: "foo",
		vec:  []float64{1, 2, 3},
	}).Word()
	if got != want {
		t.Fatalf("want %v, but %v:", want, got)
	}
}

func TestVectorNormalize(t *testing.T) {
	want := "foo[0.2672612419124244 0.5345224838248488 0.8017837257372732]"
	got := (&Vector{
		word: "foo",
		vec:  []float64{1, 2, 3},
	}).Normalize().String()
	if got != want {
		t.Fatalf("want %v, but %v:", want, got)
	}
}

func TestVectorVector(t *testing.T) {
	want := []float64{1, 2, 3}
	got := (&Vector{
		word: "foo",
		vec:  []float64{1, 2, 3},
	}).Vector()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("want %v, but %v:", want, got)
	}
}

func TestVectorAdd(t *testing.T) {
	f, err := os.Open("testdata/text.model")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	m, err := LoadText(f)
	if err != nil {
		t.Fatal(err)
	}
	v := m.Find("あ").Add(m.Find("い"))
	if v == nil {
		t.Fatal("must not be nil")
	}
	want := "あ + い"
	got := v.word
	if got != want {
		t.Fatalf("want %v, but %v:", want, got)
	}
}

func TestVectorSub(t *testing.T) {
	f, err := os.Open("testdata/text.model")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	m, err := LoadText(f)
	if err != nil {
		t.Fatal(err)
	}
	v := m.Find("あ").Sub(m.Find("い"))
	if v == nil {
		t.Fatal("must not be nil")
	}
	want := "あ - い"
	got := v.word
	if got != want {
		t.Fatalf("want %v, but %v:", want, got)
	}
}

func TestVectorCosine(t *testing.T) {
	f, err := os.Open("testdata/text.model")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	m, err := LoadText(f)
	if err != nil {
		t.Fatal(err)
	}
	v := m.Find("あ").Cosine(m.Find("い"))
	want := 0.19579290043116676
	got := v
	if got != want {
		t.Fatalf("want %v, but %v:", want, got)
	}
}

func TestLoadBinary(t *testing.T) {
	f, err := os.Open("testdata/binary.model")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	m, err := LoadBinary(f, 32)
	if err != nil {
		t.Fatal(err)
	}
	v := m.Find("あ")
	if v == nil {
		t.Fatal("must not be nil")
	}
	want := []float64{0.002279, -0.004997, 0.004357, 0.000486, -0.001852, -0.002103, -0.000743, -0.002969, -0.004407, -0.002839, 0.004320, -0.001797, -0.002695, 0.002324, -0.001714, 0.004798, -0.001596, -0.004614, 0.002748, -0.000066, -0.003513, 0.000998, 0.001459, -0.000467, 0.001829, -0.003965, 0.000903, 0.001824, -0.003046, 0.000161, -0.000058, -0.002422, -0.002854, -0.004729, 0.000010, 0.000151, 0.000008, -0.003979, -0.005007, -0.004784, -0.004356, -0.000543, 0.001403, 0.001191, -0.003193, 0.005038, -0.002024, -0.001268, -0.001418, 0.004757, -0.003789, 0.001157, -0.001356, 0.003208, 0.000057, 0.004293, -0.002839, -0.002638, -0.004208, -0.003697, -0.003210, -0.003144, 0.002457, -0.001880, 0.002713, 0.003715, 0.001379, 0.002967, 0.002498, 0.002354, -0.003588, -0.003438, -0.003679, -0.000145, 0.004170, -0.002772, -0.003095, -0.004212, 0.003314, -0.004218, -0.000570, 0.002206, -0.004646, -0.004491, 0.001470, 0.003562, 0.004297, 0.002143, 0.003191, -0.003167, -0.003670, 0.004006, -0.002779, 0.001760, 0.000589, 0.001786, -0.001196, 0.000195, -0.001702, -0.001123}
	got := v.vec
	for i := 0; i < len(want); i++ {
		if math.Abs(got[i]-want[i]) > 0.000001 {
			t.Fatalf("want[%d]=%v, got[%d]=%v:", i, float32(want[i]), i, float32(got[i]))
		}
	}
}

func TestLoadText(t *testing.T) {
	f, err := os.Open("testdata/text.model")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	m, err := LoadText(f)
	if err != nil {
		t.Fatal(err)
	}
	v := m.Find("あ")
	if v == nil {
		t.Fatal("must not be nil")
	}
	want := []float64{0.002279, -0.004997, 0.004357, 0.000486, -0.001852, -0.002103, -0.000743, -0.002969, -0.004407, -0.002839, 0.004320, -0.001797, -0.002695, 0.002324, -0.001714, 0.004798, -0.001596, -0.004614, 0.002748, -0.000066, -0.003513, 0.000998, 0.001459, -0.000467, 0.001829, -0.003965, 0.000903, 0.001824, -0.003046, 0.000161, -0.000058, -0.002422, -0.002854, -0.004729, 0.000010, 0.000151, 0.000008, -0.003979, -0.005007, -0.004784, -0.004356, -0.000543, 0.001403, 0.001191, -0.003193, 0.005038, -0.002024, -0.001268, -0.001418, 0.004757, -0.003789, 0.001157, -0.001356, 0.003208, 0.000057, 0.004293, -0.002839, -0.002638, -0.004208, -0.003697, -0.003210, -0.003144, 0.002457, -0.001880, 0.002713, 0.003715, 0.001379, 0.002967, 0.002498, 0.002354, -0.003588, -0.003438, -0.003679, -0.000145, 0.004170, -0.002772, -0.003095, -0.004212, 0.003314, -0.004218, -0.000570, 0.002206, -0.004646, -0.004491, 0.001470, 0.003562, 0.004297, 0.002143, 0.003191, -0.003167, -0.003670, 0.004006, -0.002779, 0.001760, 0.000589, 0.001786, -0.001196, 0.000195, -0.001702, -0.001123}
	got := v.vec
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("want %v, but %v:", want, got)
	}
}

func TestModelFind(t *testing.T) {
	f, err := os.Open("testdata/text.model")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	m, err := LoadText(f)
	if err != nil {
		t.Fatal(err)
	}
	v := m.Find("ん")
	if v != nil {
		t.Fatalf("must be nil: %v", v)
	}
}

func TestModelNeighbourhood(t *testing.T) {
	f, err := os.Open("testdata/text.model")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	m, err := LoadText(f)
	if err != nil {
		t.Fatal(err)
	}
	v := m.Neighbourhood(m.Find("い"))
	s := ""
	for _, vv := range v {
		s += vv.Vector.Word()
	}
	want := "あうえ"
	got := s
	if got != want {
		t.Fatalf("want %v, but %v:", want, got)
	}
}

func TestModelCosineSimilars(t *testing.T) {
	f, err := os.Open("testdata/text.model")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	m, err := LoadText(f)
	if err != nil {
		t.Fatal(err)
	}
	v := m.CosineSimilars(m.Find("う"))
	s := ""
	for _, vv := range v {
		s += vv.Vector.Word()
	}
	want := "えあい"
	got := s
	if got != want {
		t.Fatalf("want %v, but %v:", want, got)
	}
}
