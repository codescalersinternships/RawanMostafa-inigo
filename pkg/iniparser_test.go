package iniparser

import (
	"reflect"
	"testing"
)

const path = "testdata/file.ini"

func populateExpected(t *testing.T) map[string]section {
	t.Helper()
	expected := map[string]section{
		"owner": {
			map_: map[string]string{
				"name":         "John Doe",
				"organization": "Acme Widgets Inc.",
			},
		},
		"database": {
			map_: map[string]string{
				"server": "192.0.2.62",
				"port":   "143",
			},
		},
	}
	return expected
}

func TestLoadFromString(t *testing.T) {
	s := `; last modified 1 April 2001 by John Doe
			[owner]
			name = John Doe
			organization = Acme Widgets Inc.

			[database]
			; use IP address in case network name resolution is not working
			server = 192.0.2.62
			port = 143`
	parser := InitParser()
	parser.LoadFromString(s)

	expected := populateExpected(t)

	if !reflect.DeepEqual(parser.sections, expected) {
		t.Errorf("Error! wanted:%v \n , got : %v \n", expected, parser.sections)
	}
}

func TestLoadFromFile(t *testing.T) {

	expected := populateExpected(t)

	parser := InitParser()
	parser.LoadFromFile(path)

	if !reflect.DeepEqual(parser.sections, expected) {
		t.Errorf("Error! wanted:%v \n , got : %v \n", expected, parser.sections)
	}
}
