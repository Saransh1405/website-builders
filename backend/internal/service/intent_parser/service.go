package intent_parser

import (
	"context"
	"log"
	"strings"
	"unicode"
	"website-builder/internal/domain"
)

type IntentParserService struct {
}

func NewIntentParserService() *IntentParserService {
	return &IntentParserService{}
}

func (s *IntentParserService) ParseIntent(ctx context.Context, prompt string) (*domain.Intent, error) {
	NormalizePrompt := NormalizeInput(prompt)
	log.Println("normalized the inerted prompt")

	var (
		componentType   string
		componentName   string
		features        []string
		confidence      float64
		confidenceArray []float64
		avgConfidence   float64
		options         domain.RequestOptions
	)

	componentType, confidence = DetectComponentType(NormalizePrompt)
	log.Println("got the componentType & confidence")

	confidenceArray = append(confidenceArray, confidence)

	componentName, confidence = ComponentNamePatternMatching(NormalizePrompt, prompt)
	log.Println("got the componentName & confidence")

	confidenceArray = append(confidenceArray, confidence)

	features, confidence = DetectFeatures(NormalizePrompt, options)
	log.Println("got the features & confidence")

	confidenceArray = append(confidenceArray, confidence)

	avgConfidence = Avg(confidenceArray)

	return &domain.Intent{
		Type:          componentType,
		ComponentName: componentName,
		Features:      features,
		Confidence:    avgConfidence,
	}, nil
}

func NormalizeInput(input string) string {
	input = strings.ToLower(input)
	input = strings.TrimSpace(input)

	result := make([]rune, 0, len(input))
	for _, r := range input {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			result = append(result, r)
		}
	}

	return string(result)
}

func DetectComponentType(normalizeInput string) (string, float64) {
	for _, check := range componentTypeChecks {
		for _, keyword := range check.keywords {
			if strings.Contains(normalizeInput, keyword) {
				return check.typeName, check.confidence
			}
		}
	}

	return "react-component", 0.3
}

func ComponentNamePatternMatching(normalizeInput, prompt string) (string, float64) {
	normalized := normalizeInput

	// Pattern 1: Explicit naming - "create a [Name] component"
	searchZone := prompt
	foundVerb := false

	for _, verb := range actionVerbs_90 {
		verbLower := strings.ToLower(verb)
		idx := strings.Index(normalized, verbLower)
		if idx != -1 {
			searchZone = strings.TrimSpace(prompt[idx+len(verb):])
			foundVerb = true
			break
		}
	}

	if foundVerb && searchZone != "" {
		tokens := strings.Fields(searchZone)

		ignoreWords := map[string]bool{
			"a": true, "an": true, "the": true,
			"component": true, "app": true, "page": true,
			"hook": true, "form": true, "with": true,
		}

		for _, token := range tokens {
			clean := strings.Trim(token, ",.!?")
			lower := strings.ToLower(clean)

			if ignoreWords[lower] {
				continue
			}

			if len(clean) > 0 && unicode.IsUpper(rune(clean[0])) {
				return clean, 0.9
			}
		}
	}

	// Pattern 2: Direct mention - "[Name] component" or "[Name] hook"
	for _, componentType := range componentTypeWords_70 {
		idx := strings.Index(normalized, componentType)
		if idx != -1 {
			beforeText := prompt[:idx]
			tokens := strings.Fields(beforeText)

			startIdx := len(tokens) - 3
			if startIdx < 0 {
				startIdx = 0
			}

			for i := len(tokens) - 1; i >= startIdx; i-- {
				if i < 0 {
					break
				}
				clean := strings.Trim(tokens[i], ",.!?")
				if len(clean) > 0 && unicode.IsUpper(rune(clean[0])) {
					return clean, 0.7
				}
			}
		}
	}

	// Pattern 3: PascalCase naming - "TodoApp", "LoginForm", etc.
	tokens := strings.Fields(prompt)
	var pascalCaseWords []string

	for _, token := range tokens {
		clean := strings.Trim(token, ",.!?")
		if len(clean) == 0 {
			continue
		}

		if unicode.IsUpper(rune(clean[0])) {
			upperCount := 0
			for _, r := range clean {
				if unicode.IsUpper(r) {
					upperCount++
				}
			}

			if upperCount >= 1 && len(clean) >= 2 {
				lower := strings.ToLower(clean)
				ignoreWords := map[string]bool{
					"the": true, "a": true, "an": true, "i": true,
					"create": true, "build": true, "make": true,
				}
				if !ignoreWords[lower] {
					pascalCaseWords = append(pascalCaseWords, clean)
				}
			}
		}
	}

	if len(pascalCaseWords) > 0 {
		return pascalCaseWords[0], 0.6
	}

	// Pattern 4: Common app patterns (lowercase → PascalCase)
	for pattern, pascalCase := range commonPatterns_50 {
		if strings.Contains(normalized, pattern) {
			return pascalCase, 0.5
		}
	}

	// Pattern 5: Two-word combinations (lowercase → PascalCase)
	words := strings.Fields(normalized)
	for i := 0; i < len(words)-1; i++ {
		twoWordPattern := words[i] + " " + words[i+1]
		if pascalCase, exists := commonPatterns_50[twoWordPattern]; exists {
			return pascalCase, 0.5
		}

		if len(words[i]) > 0 && len(words[i+1]) > 0 {
			word1 := words[i]
			word2 := words[i+1]

			ignoreWords := map[string]bool{
				"a": true, "an": true, "the": true, "with": true,
				"for": true, "and": true, "or": true, "to": true,
			}
			if ignoreWords[word1] || ignoreWords[word2] {
				continue
			}

			pascalCase := strings.Title(word1) + strings.Title(word2)
			return pascalCase, 0.4
		}
	}

	// Fallback: return empty string with 0 confidence
	return "", 0.0
}

func DetectFeatures(normalizeInput string, options domain.RequestOptions) ([]string, float64) {
	featuresMap := make(map[string]bool)
	confidence := 0.0

	for keyword, featureName := range featureKeywords {
		keywordLower := strings.ToLower(keyword)
		if strings.Contains(normalizeInput, keywordLower) {
			featuresMap[featureName] = true
			confidence += 0.1
		}
	}

	if options.UseTypeScript {
		featuresMap["typescript"] = true
		confidence += 0.1
	}

	if options.StyleLibrary != "" {
		styleLib := strings.ToLower(options.StyleLibrary)
		switch styleLib {
		case "tailwind", "tailwindcss":
			featuresMap["tailwind"] = true
		case "styled-components", "styledcomponents":
			featuresMap["styled-components"] = true
		case "css", "css-modules", "cssmodules":
			featuresMap["css"] = true
		case "scss", "sass":
			featuresMap["scss"] = true
		}
		confidence += 0.1
	}

	if options.IncludeTests {
		featuresMap["testing"] = true
		confidence += 0.1
	}

	features := make([]string, 0, len(featuresMap))
	for feature := range featuresMap {
		features = append(features, feature)
	}

	if confidence > 1.0 {
		confidence = 1.0
	}

	if len(features) == 0 {
		return features, 0.0
	}

	return features, confidence
}

func Avg(values []float64) float64 {
	var sum float64
	for _, v := range values {
		sum += v
	}

	avg := sum / float64(len(values))
	return avg
}
