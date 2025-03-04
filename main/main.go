// Copyright 2018 Brian Noyama. Subject to the the Apache License, Version 2.0.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	bvh "github.com/briannoyama/bvh/bvh"
	"github.com/briannoyama/bvh/math32"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

type Number = math32.Number

func main() {

	config := flag.String("config", "test.json",
		"JSON configuration for the test.")
	compare := flag.Bool("compare", false,
		"Compare with Top Down method? Default False.")
	flag.Parse()
	configFile, err := os.Open(*config)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	defer configFile.Close()

	configBytes, _ := ioutil.ReadAll(configFile)

	test := &bvhTest[float32]{}
	json.Unmarshal([]byte(configBytes), test)
	if *compare {
		test.comparisonTest()
	} else {
		test.runTest()
	}
}

type operation[T math32.Number] struct {
	orth   *math32.Orthotope[T]
	opcode int
}

type bvhTest[T math32.Number] struct {
	MaxBounds *math32.Orthotope[T]
	MinVol    *[math32.DIMENSIONS]T
	MaxVol    *[math32.DIMENSIONS]T
	Additions int
	Removals  int
	Queries   int
	RandSeed  int64
}

func (b *bvhTest[T]) comparisonTest() {
	orths := make([]*math32.Orthotope[T], 0, b.Additions)
	r := rand.New(rand.NewSource(b.RandSeed))
	bvol := &bvh.BVol[*math32.Orthotope[T], T]{}
	iter := bvol.Iterator()
	for a := 0; a < b.Additions; a += 1 {
		orth := b.makeOrth(r)
		orths = append(orths, orth)

		iter.Add(orth)
		bvol2 := bvh.TopDownBVH(orths)

		fmt.Printf("%d, %d, %d, %d, %d\n", a, bvol.GetDepth(), iter.Score(),
			bvol2.GetDepth(), bvol2.Score())
	}
}

func (b *bvhTest[T]) runTest() {
	orths := make([]*math32.Orthotope[T], 0, b.Additions)
	removed := make(map[int]bool, b.Additions)
	bvol := &bvh.BVol[*math32.Orthotope[T], T]{}
	iter := bvol.Iterator()
	r := rand.New(rand.NewSource(b.RandSeed))

	if b.Removals > b.Additions {
		fmt.Printf("Incorrect config, removals larger than additions.\n")
		return
	}

	removals := *distribute(r, b.Removals, b.Additions)
	queries := *distribute(r, b.Queries, b.Additions)
	total := 0

	for a := 0; a < b.Additions; a += 1 {
		orth := b.makeOrth(r)
		orths = append(orths, orth)

		// Test the addition operation.
		t := time.Now()
		iter.Add(orth)
		duration := time.Now().Sub(t).Nanoseconds()
		total += 1
		fmt.Printf("add, %d, %d, %d\n", total, bvol.GetDepth(), duration)

		for removal := 0; removal < removals[a]; removal += 1 {
			toRemove := r.Intn(a + 1)
			for ; removed[toRemove] && toRemove <= a; toRemove += 1 {
			}
			if toRemove <= a {
				removed[toRemove] = true

				// Test the removal operation.
				t = time.Now()
				iter.Remove(orths[toRemove])
				duration := time.Now().Sub(t).Nanoseconds()
				total -= 1
				fmt.Printf("sub, %d, %d, %d\n", total, bvol.GetDepth(), duration)
			} else if a+1 < len(removals) {
				removals[a+1] += 1
			}
		}
		for query := 0; query < queries[a]; query += 1 {
			q := b.makeOrth(r)
			iter.Reset()
			count := 0

			// Test the query operation.
			t = time.Now()
			for r := iter.Query(q); r != nil; r = iter.Query(q) {
				count += 1
			}
			duration := time.Now().Sub(t).Nanoseconds()
			fmt.Printf("que, %d, %d, %d, %d\n", total, bvol.GetDepth(),
				duration, count)
		}
	}
}

func distribute(r *rand.Rand, totalEvents int, steps int) *[]int {
	events := make([]int, steps)
	for e := 0; e < totalEvents; e += 1 {
		events[r.Intn(steps)] += 1
	}

	return &events
}

func (b *bvhTest[T]) makeOrth(r *rand.Rand) *math32.Orthotope[T] {
	orth := &math32.Orthotope[T]{}
	for d := 0; d < math32.DIMENSIONS; d++ {
		orth.Delta[d] = randomValue[T](r, b.MinVol[d], b.MaxVol[d])
		maxPos := b.MaxBounds.Point[d] + b.MaxBounds.Delta[d]
		minPos := b.MaxBounds.Point[d]
		orth.Point[d] = randomValue[T](r, minPos, maxPos-orth.Delta[d])
	}
	return orth
}

func randomValue[T math32.Number](r *rand.Rand, min, max T) T {
	var zero T
	switch any(zero).(type) {
	case int32:
		return min + T(r.Int31n(int32(max-min)))
	case int64:
		return min + T(r.Int63n(int64(max-min)))
	case float32:
		return min + T(r.Float32()*float32(max-min))
	case float64:
		return min + T(r.Float64()*float64(max-min))
	default:
		panic("unsupported number type")
	}
}
