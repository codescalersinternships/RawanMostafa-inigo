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

var errEmptyParser = errors.New("this parser has no sections")
var errNoKey = errors.New("key not found")
var errNoSection = errors.New("section not found")
var errNotINI = errors.New("this is not an ini file")

// The section acts as a type representing the value of the iniparser map
type section struct {
	map_ map[string]string
}

// The iniParser acts as the data structure storing all of the parsed sections
// It has the following format:
//
//			"sectionName": {
//				map_: map[string]string{
//					"key0":   "value0",
//					"key1":   "value1",
//				},
//	},
type iniParser struct {
	sections map[string]section
}

// InitParser returns an iniParser type object
// InitParser is an essential call to get a parser to be able to access its APIs
func InitParser() iniParser {
	return iniParser{
		make(map[string]section),
	}
}

func (i iniParser) populateINI(lines []string) {
	var title string
	var sec section
	for _, line := range lines {

		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, ";") {
			continue
		}
		if strings.Contains(line, "[") && strings.Contains(line, "]") {
			if title != "" {
				i.sections[title] = sec
			}
			title = strings.Trim(line, "[]")
			title = strings.TrimSpace(title)
			sec = section{map_: make(map[string]string)}

		} else if strings.Contains(line, "=") {

			parts := strings.Split(line, "=")
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			sec.map_[key] = value
		}
	}
	if title != "" {
		i.sections[title] = sec
	}
}

// LoadFromString is a method that takes a string ini configs and parses them into the caller iniParser object
// LoadFromString assumes that:
// 1- There're no global keys, every keys need to be part of a section
// 2- The key value separator is just =
// 3- Comments are only valid at the beginning of the line
func (i iniParser) LoadFromString(file string) {

	lines := strings.Split(file, "\n")
	i.populateINI(lines)
}

// LoadFromFile is a method that takes a path to an ini file and parses it into the caller iniParser object
// LoadFromFile assumes that:
// 1- There're no global keys, every keys need to be part of a section
// 2- The key value separator is just =
// 3- Comments are only valid at the beginning of the line
func (i iniParser) LoadFromFile(path string) error {

	if !strings.HasSuffix(path, ".ini") {
		return errNotINI
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error in reading the file: %v", err)
	}
	i.LoadFromString(string(data))
	return nil
}

// GetSectionNames is a method that returns an array of strings having the names
// of the sections of the caller iniParser object and an error in case of empty sections
func (i iniParser) GetSectionNames() ([]string, error) {
	if len(i.sections) == 0 {
		return make([]string, 0), errEmptyParser
	}
	names := make([]string, 0)
	for key := range i.sections {
		names = append(names, key)
	}
	sort.Strings(names)
	return names, nil
}

// GetSections is a method that returns a map[string]section representing
// the data structure of the caller iniParser object and an error in case of empty sections
func (i iniParser) GetSections() (map[string]section, error) {
	if len(i.sections) == 0 {
		return make(map[string]section, 0), errEmptyParser
	}
	return i.sections, nil
}

// Get is a method that takes a string for the sectionName and a string for the key
// and returns the value of this key and an error
// Error formats for different cases:
//
//	If the section name passed isn't found --> "section not found"
//	If the key passed isn't found in the passed section --> "key not found"
//	else --> nil
func (i iniParser) Get(sectionName string, key string) (string, error) {

	if _, ok := i.sections[sectionName]; !ok {
		return "", errNoSection
	}
	if _, ok := i.sections[sectionName].map_[key]; !ok {
		return "", errNoKey
	}
	return i.sections[sectionName].map_[key], nil

}

// Set is a method that takes a string for the sectionName, a string for the key and a string for the value of this key
// It sets the passed key if found with the passed value and returns an error
// Error formats for different cases:
//
//	If the section name passed isn't found --> "section not found"
//	If the key passed isn't found in the passed section --> "key not found"
//	else --> nil
func (i iniParser) Set(sectionName string, key string, value string) error {
	if _, ok := i.sections[sectionName]; !ok {
		return errNoSection
	}
	if _, ok := i.sections[sectionName].map_[key]; !ok {
		return errNoKey
	}
	i.sections[sectionName].map_[key] = value
	return nil
}


// ToString is a method that returns the parsed ini map of the caller object as one string
// The returned string won't include the comments
// Also,it tells fmt pkg how to print the object
func (i iniParser) String() string {
	var result string
	sectionNames := make([]string, 0)
	for sectionName := range i.sections {
		sectionNames = append(sectionNames, sectionName)
	}
	sort.Strings(sectionNames)

	for _, sectionName := range sectionNames {
		keys := make([]string, 0)
		result += "[" + sectionName + "]\n"
		for key := range i.sections[sectionName].map_ {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			result += key + " = " + i.sections[sectionName].map_[key] + "\n"
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
func (i iniParser) SaveToFile(path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error in opening the file: %v", err)
	}
	defer file.Close()

	stringFile := i.String()
	_, err = file.WriteString(stringFile)
	if err != nil {
		return fmt.Errorf("error in writing to the file: %v", err)
	}
	return nil
}
