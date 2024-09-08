package iniparser

import (
	"log"
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

	for _, line := range lines {
		var title string
		var section section
		line := strings.TrimSpace(line)
		if strings.Contains(line, "[") && strings.Contains(line, "]") {
			title = strings.Trim(line, "[]")
		}

		log.Printf(title, section)
	}

}
