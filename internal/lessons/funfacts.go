package lessons

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"math/rand"
	"slices"

	"github.com/bobparsons/rootcamp/internal/types"
)

//go:embed data/funfacts.json
var embeddedFunFacts []byte

var cachedFacts *types.FunFactsData

func LoadFunFacts() (*types.FunFactsData, error) {
	if cachedFacts != nil {
		return cachedFacts, nil
	}

	var data types.FunFactsData
	if err := json.Unmarshal(embeddedFunFacts, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal fun facts: %w", err)
	}

	cachedFacts = &data
	return cachedFacts, nil
}

func GetRandomFact() (*types.FunFact, error) {
	data, err := LoadFunFacts()
	if err != nil {
		return nil, err
	}

	if len(data.Facts) == 0 {
		return nil, fmt.Errorf("no facts available")
	}

	idx := rand.Intn(len(data.Facts))
	return &data.Facts[idx], nil
}

// TODO: Use or remove
func GetFactsByTag(tag string) ([]types.FunFact, error) {
	data, err := LoadFunFacts()
	if err != nil {
		return nil, err
	}

	var filtered []types.FunFact
	for _, fact := range data.Facts {
		if slices.Contains(fact.Tags, tag) {
			filtered = append(filtered, fact)
		}
	}

	return filtered, nil
}
