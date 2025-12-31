package lessons

import (
	"embed"
	"encoding/json"
	"fmt"
	"math/rand"
	"path/filepath"
	"slices"
	"strings"

	"github.com/bobparsons/rootcamp/internal/types"
)

//go:embed data/funfacts/*.json
var embeddedFunFactsFS embed.FS

var cachedFacts *types.FunFactsData

func LoadFunFacts() (*types.FunFactsData, error) {
	if cachedFacts != nil {
		return cachedFacts, nil
	}

	data := &types.FunFactsData{
		Version: "1.0",
		Facts:   []types.FunFact{},
	}

	entries, err := embeddedFunFactsFS.ReadDir("data/funfacts")
	if err != nil {
		return nil, fmt.Errorf("failed to read funfacts directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			if err := loadFactsFromFile(entry.Name(), data); err != nil {
				return nil, fmt.Errorf("failed to load %s: %w", entry.Name(), err)
			}
		}
	}

	cachedFacts = data
	return cachedFacts, nil
}

func loadFactsFromFile(filename string, accumulator *types.FunFactsData) error {
	filePath := filepath.Join("data/funfacts", filename)
	content, err := embeddedFunFactsFS.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read file error: %w", err)
	}

	var fileData types.FunFactsData
	if err := json.Unmarshal(content, &fileData); err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}

	accumulator.Facts = append(accumulator.Facts, fileData.Facts...)
	return nil
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
