package iniparser

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
)

type section struct {
	map_ map[string]string
}

type iniParser struct {
	sections map[string]section
}

func InitParser() iniParser {
	return iniParser{
		make(map[string]section),
	}
}

func (i *iniParser) populateINI(lines []string) {
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

func (i *iniParser) LoadFromString(file string) {

	lines := strings.Split(file, "\n")
	i.populateINI(lines)
}

func (i *iniParser) LoadFromFile(path string) error {

	if !strings.HasSuffix(path, ".ini") {
		return fmt.Errorf("this is not an ini file")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error in reading the file: %v", err)
	}
	i.LoadFromString(string(data))
	return nil
}

func (i *iniParser) GetSectionNames() []string {
	names := make([]string, 0)
	for key := range i.sections {
		names = append(names, key)
	}
	return names
}

func (i *iniParser) GetSections() map[string]section {
	return i.sections
}

func (i *iniParser) Get(sectionName string, key string) (string, error) {
	if reflect.DeepEqual(i.sections[sectionName], section{}) {
		return "", fmt.Errorf("section not found")
	}
	if i.sections[sectionName].map_[key] == "" {
		return "", fmt.Errorf("key not found")
	}
	return i.sections[sectionName].map_[key], nil

}

func (i *iniParser) Set(sectionName string, key string, value string) error {
	if reflect.DeepEqual(i.sections[sectionName], section{}) {
		return fmt.Errorf("section not found")
	}
	if i.sections[sectionName].map_[key] == "" {
		return fmt.Errorf("key not found")
	}
	i.sections[sectionName].map_[key] = value
	return nil
}

func (i *iniParser) ToString() string {
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
	return result
}

func (i *iniParser) SaveToFile(path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error in opening the file: %v", err)
	}
	defer file.Close()

	stringFile := i.ToString()
	_, err = file.WriteString(stringFile)
	if err != nil {
		return fmt.Errorf("error in writing to the file: %v", err)
	}
	return nil
}
