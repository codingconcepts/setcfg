package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"regexp"

	"gopkg.in/yaml.v2"
)

// Initialise a default pattern that finds placeholders that are wrapped in tilde characters.
// For example: ~hello~
var placeholderPattern = regexp.MustCompile("~(.*?)~")

func main() {
	input := flag.String("i", "", "Absolute or relative path to input YAML file.")
	parts := flag.String("p", "", "Absolute or relative path to the parts YAML file.")
	pattern := flag.String("pattern", "~(.*?)~", "The regex pattern to use for extracting part keys.")
	flag.Parse()

	if *input == "" {
		flag.Usage()
		os.Exit(2)
	}
	if *parts == "" {
		flag.Usage()
		os.Exit(2)
	}
	placeholderPattern = regexp.MustCompile(*pattern)

	inputFile, err := os.Open(*input)
	if err != nil {
		log.Fatalf("error reading input file: %v", err)
	}

	partsFile, err := os.Open(*parts)
	if err != nil {
		log.Fatalf("error reading input parts file: %v", err)
	}

	output, err := set(inputFile, partsFile)
	if err != nil {
		log.Fatalf("error setting file: %v", err)
	}

	fmt.Println(output)
}

func set(inputFile, partsFile io.Reader) (string, error) {
	inputParsed, err := parse(inputFile)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling input file: %w", err)
	}
	partsParsed, err := parse(partsFile)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling parts file: %w", err)
	}

	if err = setParsed(inputParsed, partsParsed); err != nil {
		return "", fmt.Errorf("error setting yaml placeholders: %w", err)
	}

	output, err := yaml.Marshal(inputParsed)
	if err != nil {
		return "", err
	}

	return string(output), err
}

func parse(input io.Reader) (map[interface{}]interface{}, error) {
	var parsed map[interface{}]interface{}
	err := yaml.NewDecoder(input).Decode(&parsed)

	// If there aren't any parts, just return an empty collection.
	if err == io.EOF {
		return map[interface{}]interface{}{}, nil
	}

	return parsed, err
}

func setParsed(inputParsed, partsParsed map[interface{}]interface{}) error {
	for k, v := range inputParsed {
		// Walk the YAML file recursively.
		if reflect.ValueOf(v).Kind() == reflect.Map {
			if err := setParsed(v.(map[interface{}]interface{}), partsParsed); err != nil {
				return err
			}
		}

		// Look for the SYAML_ prefix for the element value and if we don't have one,
		// continue through the rest of the YAML file.
		key, ok := isPlaceholder(v)
		if !ok {
			continue
		}

		// We've got a SYAML_ prefix. Find the part from the parts file and apply. If
		// we don't have a part, return an error, as we'll never be able to complete
		// this file.
		part, ok := partsParsed[key]
		if !ok {
			return fmt.Errorf("missing part for %q", v)
		}

		inputParsed[k] = part
	}

	return nil
}

func isPlaceholder(value interface{}) (string, bool) {
	valueStr, ok := value.(string)
	if !ok {
		return "", false
	}

	matches := placeholderPattern.FindStringSubmatch(valueStr)
	if len(matches) < 2 {
		return "", false
	}

	return matches[1], true
}
