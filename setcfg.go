package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

// Initialise a default pattern that finds placeholders that are wrapped in tilde characters.
// For example: ~hello~
var placeholderPattern = regexp.MustCompile("~(.*?)~")

// A type to hold any number of input arguments to support ad-hoc
type flagStrings []string

func (fs *flagStrings) String() string {
	return strings.Join(*fs, ",")
}

func (fs *flagStrings) Set(value string) error {
	*fs = append(*fs, value)
	return nil
}

func main() {
	input := flag.String("i", "", "Absolute or relative path to input YAML file.")
	env := flag.String("e", "", "Absolute or relative path to the environment YAML file.")
	pattern := flag.String("pattern", "~(.*?)~", "The regex pattern to use for extracting keys.")

	adhoc := &flagStrings{}
	flag.Var(adhoc, "f", "A list of 'key=value' fields to substitute (useful as an alternative to -e if all you're substituting are simple fields).")

	flag.Parse()

	if *input == "" {
		flag.Usage()
		os.Exit(2)
	}

	placeholderPattern = regexp.MustCompile(*pattern)

	inputFile, err := os.Open(*input)
	if err != nil {
		log.Fatalf("error reading input file: %v", err)
	}
	inputParsed, err := parse(inputFile)
	if err != nil {
		log.Fatalf("error unmarshalling input file: %v", err)
	}

	envParsed := map[interface{}]interface{}{}
	if *env != "" {
		envFile, err := os.Open(*env)
		if err != nil {
			log.Fatalf("error reading input environment file: %v", err)
		}
		if envParsed, err = parse(envFile); err != nil {
			log.Fatalf("error unmarshalling environment file: %v", err)
		}
	}

	if err := addAdhocFields(envParsed, adhoc); err != nil {
		log.Fatalf("error setting adhoc fields: %v", err)
	}

	if err := setParsed(inputParsed, envParsed); err != nil {
		log.Fatalf("error setting file: %v", err)
	}

	output, err := yaml.Marshal(inputParsed)
	if err != nil {
		log.Fatalf("error marshalling output: %v", err)
	}

	fmt.Println(string(output))
}

func parse(input io.Reader) (map[interface{}]interface{}, error) {
	var parsed map[interface{}]interface{}
	err := yaml.NewDecoder(input).Decode(&parsed)

	// If there aren't any environment fields, just return an empty collection.
	if err == io.EOF {
		return map[interface{}]interface{}{}, nil
	}

	return parsed, err
}

func addAdhocFields(envParsed map[interface{}]interface{}, adhoc *flagStrings) error {
	if adhoc == nil {
		return nil
	}

	for _, field := range *adhoc {
		k, v, err := parseAdhocKeyValue(field)
		if err != nil {
			return err
		}

		envParsed[k] = v
	}

	return nil
}

func parseAdhocKeyValue(field string) (string, string, error) {
	parts := strings.Split(field, "=")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("adhoc fields must be in the format of key=value")
	}

	return parts[0], parts[1], nil
}

func setParsed(inputParsed, envParsed map[interface{}]interface{}) error {
	for k, v := range inputParsed {
		if v == nil {
			continue
		}

		switch reflect.TypeOf(v).Kind() {
		// Recurse into complex fields.
		case reflect.Map:
			if err := setParsed(v.(map[interface{}]interface{}), envParsed); err != nil {
				return err
			}
		// Set complex and scalar fields.
		case reflect.Slice:
			x := reflect.ValueOf(v)
			for i := 0; i < x.Len(); i++ {
				switch reflect.TypeOf(x.Index(i).Interface()).Kind() {
				case reflect.Slice:
					log.Println("implement support for slice of slices")
				case reflect.Map:
					if err := setParsed(x.Index(i).Interface().(map[interface{}]interface{}), envParsed); err != nil {
						return err
					}
				case reflect.String:
					if err := setValue(v, v, inputParsed, envParsed); err != nil {
						return err
					}
				}
			}
		// Set scalar field.
		default:
			if err := setValue(k, v, inputParsed, envParsed); err != nil {
				return err
			}
		}
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

func setValue(k, v interface{}, inputParsed, envParsed map[interface{}]interface{}) error {
	key, ok := isPlaceholder(v)
	if !ok {
		return nil
	}

	part, ok := envParsed[key]
	if !ok {
		return fmt.Errorf("missing part for %q", v)
	}

	inputParsed[k] = part
	return nil
}
