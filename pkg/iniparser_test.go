package iniparser

import (
	"reflect"
	"testing"
)

func TestLoadFromString(t *testing.T) {
	s := `; last modified 1 April 2001 by John Doe
			[owner]
			name = John Doe
			organization = Acme Widgets Inc.`

	section1 := section{make(map[string]string)}
	section1.map_["name"] = "John Doe"
	section1.map_["organization"] = "Acme Widgets Inc."
	expected := make(map[string]section)
	expected["owner"] = section1

	parser := InitParser()
	parser.LoadFromString(s)

	if !reflect.DeepEqual(parser.sections, expected) {
		t.Errorf("Error! wanted:%v \n , got : %v \n", expected, parser.sections)
	}
}
