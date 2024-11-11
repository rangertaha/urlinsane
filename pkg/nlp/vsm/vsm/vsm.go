// Package vsm implements algebraic
// vector space model for comparison between documents
// and queries.
//
// Vector space model represents documents and
// queries as objects (points or vector) in an
// N-dimensional space where each term is a dimension.
// So a document d with N number of terms can be
// represented as point or vector with coordinates
// d(w₁, w₂, ... wN), where w = term weight.
//
// The term weighting scheme used in this package is the TFIDF:
//
//	w = term.Count * log(|Docs| / |{ d ∈ Docs | term ∈ d}|)
//
// Where Docs are the the known documents in the corpus.
//
// Vector space model uses the deviation of angles
// between each document vector and the query vector
// to check their similarities by calculating
// the cosine of the angle between the vectors.
// So for each document vector dᵢ and a query vector q:
//
//	cos0 = dᵢ•q / ||dᵢ|| * ||q||
//
// where ||dᵢ|| = the magnitude of the document vector
// and ||q|| = the magnitude of the query vector.
//
// See: http://www.minerazzi.com/tutorials/term-vector-1.pdf
package vsm

import (
	"context"
	"math"
	"strings"
	"sync"
	"sync/atomic"

	"golang.org/x/text/transform"
)

// VSM holds the corpus for the algebraic
// vector space model calculation.
type VSM struct {
	terms *terms

	mu   sync.RWMutex
	docs []document

	docsCount uint64

	// transformer is used for filtering
	// the documents and query sentences.
	transformer transform.Transformer
}

// Document  holds a sentence, which is tokenized and
// filtered in order to be represented as a vector,
// and a class for keywording/tagging the sentence.
type Document struct {
	Sentence string
	Class    string
}

type term struct {
	// DocsSeen holds the number of times
	// the term has been seen in document.
	// But it is only counted once even
	// if it's in a document multiple times.
	docsSeen uint64
}

// terms holds a read/write mutual exclusion lock for
// synchronizing access to the map of terms.
type terms struct {
	mu    sync.RWMutex
	terms map[string]term
}

func (t *terms) Get(term string) (term, bool) {
	t.mu.RLock()
	res, ok := t.terms[term]
	t.mu.RUnlock()
	return res, ok
}

func (t *terms) Set(k string, v term) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.terms[k] = v
}

type document struct {
	Document

	termFreq map[string]uint64
}

// New returns a VSM structure used
// for searching documents in the corpus.
// The Transformer is used for filtering
// the documents sentences and queries.
func New(t transform.Transformer) *VSM {
	vsm := &VSM{
		terms:       &terms{sync.RWMutex{}, make(map[string]term)},
		mu:          sync.RWMutex{},
		docs:        []document{},
		transformer: t,
	}

	return vsm
}

// StaticTraining receives Document used for augmenting the corpus.
// It returns error if the training was unsucessful.
func (v *VSM) StaticTraining(dc Document) error {
	doc := document{Document: dc, termFreq: make(map[string]uint64)}

	// Uses map of string of pmpty struct used
	// for it occupies no space in memory.
	seenTerms := make(map[string]struct{})

	sentence, err := v.sanitize(dc.Sentence)
	if err != nil {
		return err
	}

	for _, trm := range strings.Split(sentence, " ") {
		t := strings.ToLower(strings.TrimSpace(trm))

		if _, ok := v.terms.Get(t); !ok {
			v.terms.Set(t, term{})
		}

		doc.termFreq[t]++

		seenTerms[t] = struct{}{}
	}

	for trm := range seenTerms {
		t, _ := v.terms.Get(trm)
		t.docsSeen++
		v.terms.Set(trm, t)
	}

	v.mu.Lock()
	v.docs = append(v.docs, doc)
	v.mu.Unlock()

	atomic.StoreUint64(&v.docsCount, uint64(len(v.docs)))

	return nil
}

// TrainingResult holds the Document used for the
// training process and an error which will be nil
// if no error occurred.
type TrainingResult struct {
	Doc Document
	Err error
}

// DynamicTraining receives a producer channel of Document used
// for dynamically augmenting the corpus. It returns a TrainingResult
// channel which can be used to check if an error occurred
// during the training process.
func (v *VSM) DynamicTraining(ctx context.Context, docCh <-chan Document) <-chan TrainingResult {
	trainCh := make(chan TrainingResult)

	go func() {
		// Releases resources of trainCh.
		defer close(trainCh)

		for {
			select {
			case doc, ok := <-docCh:
				// Checks if sender is still open.
				// If not, exists.
				if !ok {
					return
				}

				// Try to send the result of the training.
				// If no one is interested, discards this value.
				select {
				default:
				case trainCh <- TrainingResult{Doc: doc, Err: v.StaticTraining(doc)}:
				}

			case <-ctx.Done():
				return
			}
		}
	}()

	return trainCh
}

// Search returns the most similar document from the corpus
// with the query based on vector space model, or an error.
// A nil Document means there's no similarity between any
// document in the corpus and the query.
func (v *VSM) Search(query string) (*Document, error) {
	queryDoc := document{termFreq: make(map[string]uint64)}

	query, err := v.sanitize(query)
	if err != nil {
		return nil, err
	}

	for _, term := range strings.Split(query, " ") {
		t := strings.ToLower(strings.TrimSpace(term))

		queryDoc.termFreq[t]++
	}

	totalDocs := atomic.LoadUint64(&v.docsCount)

	var querySum float64
	for trm, freq := range queryDoc.termFreq {
		t, ok := v.terms.Get(trm)
		if !ok {
			continue
		}

		idf := math.Log(float64(totalDocs) / float64(t.docsSeen))

		weight := float64(freq) * idf

		querySum += math.Pow(weight, 2)
	}

	queryMag := math.Sqrt(querySum)

	var foundDoc *Document

	var maxSim float64
	v.mu.RLock()
	for _, doc := range v.docs {
		var docSum float64
		var coeff float64

		for trm, freq := range doc.termFreq {
			t, _ := v.terms.Get(trm)

			idf := math.Log(float64(totalDocs) / float64(t.docsSeen))

			weight := float64(freq) * idf

			queryTermWeight := float64(queryDoc.termFreq[trm]) * idf

			coeff += weight * queryTermWeight

			docSum += math.Pow(float64(weight), 2)
		}

		docMag := math.Sqrt(docSum)

		sim := coeff / (docMag * queryMag)
		if sim > maxSim {
			foundDoc = &Document{Sentence: doc.Sentence, Class: doc.Class}
			maxSim = sim
		}
	}
	v.mu.RUnlock()

	return foundDoc, nil
}

// sanatize applies the transformations to the sentence.
func (v *VSM) sanitize(sentence string) (string, error) {
	if v.transformer != nil {
		sanitized, _, err := transform.String(v.transformer, sentence)
		if err != nil {
			return sentence, err
		}

		return sanitized, nil
	}

	// If no transformer is set, returns the plain sentence.
	return sentence, nil
}
