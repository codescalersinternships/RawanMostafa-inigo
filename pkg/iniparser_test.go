package iniparser

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

const path = "testdata/file.ini"

const stringINI = `; last modified 1 April 2001 by John Doe
[owner]
name = John Doe
organization = Acme Widgets Inc.

[database]
; use IP address in case network name resolution is not working
server = 192.0.2.62
port = 143
[section]
key0 = val0    
key1 = val1`

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
		"section": {
			map_: map[string]string{
				"key0": "val0",
				"key1": "val1",
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

func assertFile(t *testing.T, filePath string, expectedData string) {
	t.Helper()
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Errorf("Error! : %v\n", err)
	}
	assertEquality(t, expectedData, string(data))

}

func TestLoadFromString(t *testing.T) {
	t.Run("test normal ini file", func(t *testing.T) {

		parser := InitParser()
		parser.LoadFromString(stringINI)

		expected := populateExpectedNormal(t)

		assertEquality(t, expected, parser.sections)

	})

	t.Run("test empty section ini file", func(t *testing.T) {
		const emptySectionINI = `; last modified 1 April 2001 by John Doe
				[owner]
				name = John Doe
				organization = Acme Widgets Inc.
	
				[database]`
		parser := InitParser()
		parser.LoadFromString(emptySectionINI)

		expected := populateExpectedEmptySection(t)

		assertEquality(t, expected, parser.sections)
	})
}

func TestLoadFromFile(t *testing.T) {

	expected := populateExpectedNormal(t)

	parser := InitParser()

	err := parser.LoadFromFile(path)
	if err != nil {
		t.Errorf("Error! %v", err)
	}

	assertEquality(t, expected, parser.sections)
}

func TestGetSectionNames(t *testing.T) {
	parser := InitParser()
	err := parser.LoadFromFile(path)
	if err != nil {
		t.Errorf("Error! %v", err)
	}
	names := parser.GetSectionNames()

	expected := []string{"owner", "database", "section"}

	assertEquality(t, expected, names)
}

func TestGetSections(t *testing.T) {
	parser := InitParser()
	err := parser.LoadFromFile(path)
	if err != nil {
		t.Errorf("Error! %v", err)
	}

	got := parser.GetSections()

	expected := populateExpectedNormal(t)
	assertEquality(t, expected, got)
}

func TestGet(t *testing.T) {

	testcases := []struct {
		testcaseName string
		sectionName  string
		key          string
		expected     string
		err          error
	}{
		{
			testcaseName: "Normal case: section and key are present",
			sectionName:  "database",
			key:          "server",
			expected:     "192.0.2.62",
			err:          nil,
		},
		{
			testcaseName: "corner case: section not found",
			sectionName:  "user",
			key:          "server",
			err:          fmt.Errorf("section not found"),
		},
		{
			testcaseName: "corner case: key not found",
			sectionName:  "database",
			key:          "size",
			err:          fmt.Errorf("key not found"),
		},
	}
	for _, testcase := range testcases {

		t.Run(testcase.testcaseName, func(t *testing.T) {
			parser := InitParser()
			err := parser.LoadFromFile(path)
			if err != nil {
				t.Errorf("Error! %v", err)
			}
			got, err := parser.Get(testcase.sectionName, testcase.key)
			if testcase.err == nil {
				assertEquality(t, testcase.expected, got)
			} else {
				assertEquality(t, testcase.err, err)
			}
		})
	}

}

func TestSet(t *testing.T) {

	testcases := []struct {
		testcaseName string
		sectionName  string
		key          string
		value        string
		err          error
	}{
		{
			testcaseName: "Normal case: section and key are present",
			sectionName:  "database",
			key:          "server",
			value:        "127.0.0.1",
			err:          nil,
		},
		{
			testcaseName: "corner case: section not found",
			sectionName:  "user",
			key:          "server",
			err:          fmt.Errorf("section not found"),
		},
		{
			testcaseName: "corner case: key not found",
			sectionName:  "database",
			key:          "size",
			err:          fmt.Errorf("key not found"),
		},
	}
	for _, testcase := range testcases {

		t.Run(testcase.testcaseName, func(t *testing.T) {
			parser := InitParser()
			err := parser.LoadFromFile(path)
			if err != nil {
				t.Errorf("Error! %v", err)
			}
			err = parser.Set(testcase.sectionName, testcase.key, testcase.value)
			if testcase.err == nil {
				value := parser.sections[testcase.sectionName].map_[testcase.key]
				assertEquality(t, testcase.value, value)
			} else {
				assertEquality(t, testcase.err, err)
			}
		})
	}

}

func TestToString(t *testing.T) {
	parser1 := InitParser()
	parser2 := InitParser()

	parser1.LoadFromString(stringINI)
	got := parser1.ToString()

	parser2.LoadFromString(got)

	assertEquality(t, parser1.sections, parser2.sections)
}

func TestSaveToFile(t *testing.T) {
	const outPath = "testdata/out.ini"
	parser := InitParser()
	err := parser.LoadFromFile(path)
	if err != nil {
		t.Errorf("Error! %v", err)
	}

	err = parser.SaveToFile(outPath)
	if err != nil {
		t.Errorf("Error! %v", err)
	}

	stringFile := parser.ToString()
	assertFile(t, outPath, stringFile)
}
