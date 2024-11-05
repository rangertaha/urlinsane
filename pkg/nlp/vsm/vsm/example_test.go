package vsm

import (
	"context"
	"fmt"
	"unicode"

	"golang.org/x/text/runes"
)

func Example() {
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
			Sentence: "Shipment of gold arrived in a truck.",
			Class:    "d3",
		},
	}

	vsm := New(nil)

	for _, doc := range docs {
		fmt.Println(vsm.StaticTraining(doc))
	}

	doc, err := vsm.Search("gold silver truck.")

	fmt.Println(doc.Class, err)
	// Output:
	// <nil>
	// <nil>
	// <nil>
	// d2 <nil>
}

func ExampleVSM_Search() {
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

	// transformer will be applied to every sentence
	// and query.
	transformer := runes.Map(func(r rune) rune {
		// Replaces hyphens by space.
		if unicode.Is(unicode.Hyphen, r) {
			return ' '
		}
		return r
	})

	vsm := New(transformer)

	for _, doc := range docs {
		fmt.Println(vsm.StaticTraining(doc))
	}

	doc, err := vsm.Search("shipment gold in a flying truck.")

	fmt.Println(doc.Class, err)
	// Output:
	// <nil>
	// <nil>
	// <nil>
	// d3 <nil>
}

func ExampleVSM_DynamicTraining() {
	docs := []Document{
		{
			Sentence: "Shipment of gold damaged in a fire.",
			Class:    "d1",
		},
		{
			Sentence: "Shipment of gold arrived in a truck.",
			Class:    "d3",
		},
	}

	vsm := New(nil)

	for _, doc := range docs {
		err := vsm.StaticTraining(doc)
		fmt.Println(err)
	}

	docCh := make(chan Document)

	go func() {
		defer close(docCh)

		// Loads document from some source dynamically
		// and sends it to the training channel.
		docCh <- Document{
			Sentence: "Delivery of silver arrived in a silver truck.",
			Class:    "d2",
		}
	}()

	trainCh := vsm.DynamicTraining(context.Background(), docCh)

	// Waits until all documents are consumed by the training process.
	for {
		_, ok := <-trainCh
		if !ok {
			break
		}
	}

	doc, err := vsm.Search("gold silver truck.")

	fmt.Println(doc.Class, err)
	// Output:
	// <nil>
	// <nil>
	// d2 <nil>
}
