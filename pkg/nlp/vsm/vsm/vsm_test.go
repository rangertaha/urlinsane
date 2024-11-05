package vsm

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
)

var fromFile = flag.Bool("fromfile", false, `test from files inside "testdata" dir.`)
var fileName = flag.String("filename", "training.json", "name of the file that contains the tests.")

type fileTest struct {
	Docs      []Document `json:"documents"`
	Transform *struct {
		Map *struct {
			Runes string `json:"runes"`
			To    string `json:"to"`
		} `json:"map"`
	} `json:"transform"`
	Tests []struct {
		Query string `json:"query"`
		Want  string `json:"wantClass"`
	} `json:"tests"`
}

// openTestFile reads a file from `name` and returns a file descriptor.
// The file should be in `testdata` dir, as the `name` will be prefixed
// with `testdata/`. It the caller's responsability to close the file when needed.
func openTestFile(name string) (*os.File, error) {
	file, err := os.Open(fmt.Sprintf("testdata%s%s", string(os.PathSeparator), name))
	if err != nil {
		return nil, err
	}

	return file, nil
}

func TestVSMDynamicTraining(t *testing.T) {
	assertClosedChan := func(t *testing.T, c <-chan TrainingResult) {
		select {
		case _, ok := <-c:
			if ok {
				t.Error("Got channel open, want it to be closed")
			}
		case <-time.Tick(500 * time.Millisecond):
			t.Error("trainCh timed out.")
		}

	}

	t.Run("ClosedInputChannel", func(t *testing.T) {
		docCh := make(chan Document)
		close(docCh)

		vsm := New(nil)

		assertClosedChan(t, vsm.DynamicTraining(context.Background(), docCh))
	})

	t.Run("ContextDone", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		vsm := New(nil)

		docCh := make(chan Document)
		defer close(docCh)

		assertClosedChan(t, vsm.DynamicTraining(ctx, docCh))
	})

}

// testingTransformer is used for stubbing the
// sentences filtering.
type testingTransformer struct {
	nDst, nSrc int
	err        error
}

func (t *testingTransformer) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	return t.nDst, t.nSrc, t.err
}

func (t *testingTransformer) Reset() {
}

func TestVSMStaticTrainingError(t *testing.T) {
	vsm := New(&testingTransformer{err: errors.New("Testing Error")})

	if err := vsm.StaticTraining(Document{}); err == nil {
		t.Error("Got error nil; want not nil.")
	}
}

func TestVSMDynamicTrainingError(t *testing.T) {
	vsm := New(&testingTransformer{err: errors.New("Testing Error")})

	trainingCh := make(chan Document, 1)
	defer close(trainingCh)

	trainingCh <- Document{}

	resCh := vsm.DynamicTraining(context.Background(), trainingCh)

	select {
	case res := <-resCh:
		if res.Err == nil {
			t.Error("Got error nil, want not nil.")
		}
	case <-time.Tick(500 * time.Millisecond):
		t.Fatal("error channel timed out")
	}
}

func TestVSMSearch(t *testing.T) {
	docs := []Document{
		{
			Sentence: "Shipment of gold damaged in a fire.",
			Class:    "d1",
		},
		{
			Sentence: "Delivery of silver arrived in a silver truck.",
			Class:    "d2",
		},
		{
			Sentence: "Shipment-of-gold-arrived in a truck.",
			Class:    "d3",
		},
	}

	testCases := []struct {
		transformer transform.Transformer
		query       string
		want        *Document
	}{
		{
			transformer: nil,
			query:       "gold silver truck.",
			want:        &Document{"Delivery of silver arrived in a silver truck.", "d2"},
		},
		{
			transformer: nil,
			query:       "shipment gold fire.",
			want:        &Document{"Shipment of gold damaged in a fire.", "d1"},
		},
		{
			transformer: runes.Map(func(r rune) rune {
				if unicode.Is(unicode.Hyphen, r) {
					return ' '
				}
				return r
			}),
			query: "shipment gold in a flying truck.",
			want:  &Document{"Shipment-of-gold-arrived in a truck.", "d3"},
		},
		{
			transformer: nil,
			query:       "this query should result an empty document.",
			want:        nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.query, func(t *testing.T) {
			vsm := New(tc.transformer)

			for _, doc := range docs {
				if err := vsm.StaticTraining(doc); err != nil {
					t.Fatalf("Got error while training: %q; want nil.", err)
				}
			}

			got, err := vsm.Search(tc.query)
			if err != nil {
				t.Fatalf("Got error while searching for %q: %q.", tc.query, err)
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Got %+v classifier; want %+v.", got, tc.want)
			}
		})
	}
}

// setupTransformer get a fileTest struct and returns a
// transform.Transformer with all transformations configured
// in the fileTest.Transform.
func setupTransformer(training fileTest) transform.Transformer {
	var transf transform.Transformer
	if training.Transform != nil {
		var transformers []transform.Transformer

		if m := training.Transform.Map; m != nil {
			newRune := []rune(m.To)
			transformers = append(transformers, runes.Map(func(r rune) rune {
				if strings.ContainsRune(m.Runes, r) && len(newRune) > 0 {
					return newRune[0]
				}
				return r
			}))
		}

		if len(transformers) > 0 {
			transf = transform.Chain(transformers...)
		}
	}

	return transf
}

func TestVSMSearchFromFile(t *testing.T) {
	if !*fromFile {
		t.Skip("Skipping tests loaded from file")
	}

	f, err := openTestFile(*fileName)
	if err != nil {
		t.Fatalf("got error loading test file: 'testdata/%s'.", *fileName)
	}
	defer f.Close()

	var training fileTest

	if err := json.NewDecoder(f).Decode(&training); err != nil {
		t.Fatalf("got error parsing test file 'testdata/%s'.", *fileName)
	}

	for _, tc := range training.Tests {
		t.Run(tc.Query, func(t *testing.T) {
			vsm := New(setupTransformer(training))

			for _, doc := range training.Docs {
				if err := vsm.StaticTraining(doc); err != nil {
					t.Fatalf("Got error while training: %q; want nil.", err)
				}
			}

			doc, err := vsm.Search(tc.Query)
			if err != nil {
				t.Fatalf("Got error while searching for %q: %q.", tc.Query, err)
			}

			if doc == nil {
				t.Fatalf("Got no document found for query: %q.", tc.Query)
			}

			if got := doc.Class; got != tc.Want {
				t.Errorf("Got %q class; want %q.", got, tc.Want)
			}
		})
	}
}

func TestVSMSearchError(t *testing.T) {
	vsm := New(&testingTransformer{err: errors.New("Testing Error")})

	if _, err := vsm.Search("testing"); err == nil {
		t.Error("Got error nil while searching, want not nil.")
	}
}
