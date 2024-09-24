// Package iniparser implements a parser for ini files
//
// The iniparser package should only be used for for ini files and doesn't support any other types of files
package iniparser

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
)

var ErrNoKey = errors.New("key not found")
var ErrNoSection = errors.New("section not found")
var ErrNotINI = errors.New("this is not an ini file")
var ErrGlobalKey = errors.New("global keys aren't supported")

// The Parser acts as the data structure storing all of the parsed sections
type Parser struct {
	sections map[string]map[string]string
}

// NewParser returns a Parser type object
// NewParser is an essential call to get a parser to be able to access its APIs
func NewParser() Parser {
	return Parser{
		make(map[string]map[string]string),
	}
}

func (i Parser) parse(lines []string) error {
	var title string
	var sec map[string]string
	inSection := false
	for _, line := range lines {

		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			inSection = true
			if title != "" {
				i.sections[title] = sec
			}
			title = strings.Trim(line, "[]")
			sec = make(map[string]string)

		} else if strings.Contains(line, "=") {
			if !inSection {
				return ErrGlobalKey
			}
			parts := strings.SplitN(line, "=", 2)
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			sec[key] = value
		}
	}
	if title != "" {
		i.sections[title] = sec
	}
	return nil
}

// LoadFromString is a method that takes a string ini configs and parses them into the caller Parser object
// LoadFromString assumes that:
// 1- There're no global keys, every keys need to be part of a section
// 2- The key value separator is just =
// 3- Comments are only valid at the beginning of the line
func (i Parser) LoadFromString(data string) (err error) {

	lines := strings.Split(data, "\n")
	err = i.parse(lines)
	return
}

// LoadFromFile is a method that takes a path to an ini file and parses it into the caller Parser object
// LoadFromFile assumes that:
// 1- There're no global keys, every keys need to be part of a section
// 2- The key value separator is just =
// 3- Comments are only valid at the beginning of the line
func (i Parser) LoadFromFile(path string) error {

	if !strings.HasSuffix(path, ".ini") {
		return ErrNotINI
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error in reading the file: %w", err)
	}
	i.LoadFromString(string(data))
	return nil
}

// GetSectionNames is a method that returns an array of strings having the names
// of the sections of the caller Parser object and an error in case of empty sections
func (i Parser) GetSectionNames() ([]string, error) {
	if len(i.sections) == 0 {
		return make([]string, 0), nil
	}
	names := make([]string, 0)
	for key := range i.sections {
		names = append(names, key)
	}
	sort.Strings(names)
	return names, nil
}

// GetSections is a method that returns a map[string]section representing
// the data structure of the caller Parser object and an error in case of empty sections
func (i Parser) GetSections() (map[string]map[string]string, error) {
	if len(i.sections) == 0 {
		return make(map[string]map[string]string, 0), nil
	}
	return i.sections, nil
}

// Get is a method that takes a string for the sectionName and a string for the key
// and returns the value of this key and a boolean to indicate if found or not
func (i Parser) Get(sectionName string, key string) (string, bool) {

	val, exists := i.sections[sectionName][key]
	return val, exists
}

// Set is a method that takes a string for the sectionName, a string for the key and a string for the value of this key
// It sets the passed key if found with the passed value and returns an error
// Error formats for different cases:
//
//	If the section name passed isn't found --> "section not found"
//	If the key passed isn't found in the passed section --> "key not found"
//	else --> nil
func (i Parser) Set(sectionName string, key string, value string) error {
	if _, ok := i.sections[sectionName]; !ok {
		return ErrNoSection
	}
	if _, ok := i.sections[sectionName][key]; !ok {
		return ErrNoKey
	}
	i.sections[sectionName][key] = value
	return nil
}

// String is a method that returns the parsed ini map of the caller object as one string
// The returned string won't include the comments
// Also,it tells fmt pkg how to print the object
func (i Parser) String() string {
	var result string
	sectionNames := make([]string, 0)
	for sectionName := range i.sections {
		sectionNames = append(sectionNames, sectionName)
	}
	sort.Strings(sectionNames)

	for _, sectionName := range sectionNames {
		keys := make([]string, 0)
		result += fmt.Sprintf("[%s]\n", sectionName)
		for key := range i.sections[sectionName] {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			result += fmt.Sprintf("%s = %s\n", key, i.sections[sectionName][key])
		}
	}
	return fmt.Sprint(result)
}

// SaveToFile is a method that takes a path to an output file and returns an error
// It saves the parsed ini map of the caller object into a file
// Error formats for different cases:
//
//	If the file couldn't be opened --> "error in opening the file:"
//	If writing to the file failed --> "error in writing to the file:"
//	else --> nil
func (i Parser) SaveToFile(path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error in opening the file: %w", err)
	}
	defer file.Close()

	stringFile := i.String()
	_, err = file.WriteString(stringFile)
	if err != nil {
		return fmt.Errorf("error in writing to the file: %w", err)
	}
	return nil
}
