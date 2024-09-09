package iniparser

import (
	"fmt"
	"log"
	"os"
	"reflect"
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

func (i *iniParser) LoadFromFile(path string) {

	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("Error in reading the file")
	}
	i.LoadFromString(string(data))
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
