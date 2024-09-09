package iniparser

import (
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

func (i *iniParser) LoadFromString(file string) {

	lines := strings.Split(file, "\n")

	var title string
	var sec section
	for _, line := range lines {
		line = strings.TrimSpace(line)
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
