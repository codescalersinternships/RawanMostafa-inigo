package iniparser

import (
	"reflect"
	"testing"
)

const path = "testdata/file.ini"

func populateExpectedNormal(t *testing.T) map[string]section {
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

func populateExpectedEmptySection(t *testing.T) map[string]section {
	t.Helper()
	expected := map[string]section{
		"owner": {
			map_: map[string]string{
				"name":         "John Doe",
				"organization": "Acme Widgets Inc.",
			},
		},
		"database": {
			map_: map[string]string{},
		},
	}
	return expected
}

func assertEquality(t *testing.T, obj1 any, obj2 any) {
	t.Helper()
	if reflect.TypeOf(obj1) != reflect.TypeOf(obj2) {
		t.Errorf("Error! type mismatch, wanted %t got %t", reflect.TypeOf(obj1), reflect.TypeOf(obj2))
	}
	if !reflect.DeepEqual(obj1, obj2) {
		t.Errorf("Error! values mismatch, wanted %v got %v", obj1, obj2)
	}
}

func TestLoadFromString(t *testing.T) {
	t.Run("test normal ini file", func(t *testing.T) {
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

		expected := populateExpectedNormal(t)

		assertEquality(t, expected, parser.sections)

	})

	t.Run("test empty section ini file", func(t *testing.T) {
		s := `; last modified 1 April 2001 by John Doe
				[owner]
				name = John Doe
				organization = Acme Widgets Inc.
	
				[database]`
		parser := InitParser()
		parser.LoadFromString(s)

		expected := populateExpectedEmptySection(t)

		assertEquality(t, expected, parser.sections)
	})
}

func TestLoadFromFile(t *testing.T) {

	expected := populateExpectedNormal(t)

	parser := InitParser()
	parser.LoadFromFile(path)

	assertEquality(t, expected, parser.sections)
}

func TestGetSectionNames(t *testing.T) {
	parser := InitParser()
	parser.LoadFromFile(path)
	names := parser.GetSectionNames()

	expected := []string{"owner", "database"}

	assertEquality(t, names, expected)
}
