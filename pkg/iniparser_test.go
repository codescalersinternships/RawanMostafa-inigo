package iniparser

import (
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

	expectedFile := populateExpectedNormal(t)

	testcases := []struct {
		testcaseName string
		filePath     string
		expected     map[string]section
		err          string
	}{
		{
			testcaseName: "Normal case: ini file is present",
			filePath:     "testdata/file.ini",
			expected:     expectedFile,
			err:          "",
		},
		{
			testcaseName: "corner case: not an ini file",
			filePath:     "testdata/file.txt",
			err:          "this is not an ini file",
		},
		{
			testcaseName: "corner case: file not found",
			filePath:     "testdata/filex.ini",
			err:          "error in reading the file: open testdata/filex.ini: no such file or directory",
		},
	}
	for _, testcase := range testcases {

		t.Run(testcase.testcaseName, func(t *testing.T) {
			parser := InitParser()

			err := parser.LoadFromFile(testcase.filePath)
			if testcase.err == "" {
				assertEquality(t, expectedFile, parser.sections)
			} else {
				assertEquality(t, testcase.err, err.Error())
			}

		})
	}
}

func TestGetSectionNames(t *testing.T) {
	t.Run("Normal case: sections are not empty", func(t *testing.T) {

		parser := InitParser()
		err := parser.LoadFromFile(path)
		if err != nil {
			t.Errorf("Error! %v", err)
		}
		names, err := parser.GetSectionNames()

		expected := []string{"database", "owner", "section"}

		assertEquality(t, expected, names)
		assertEquality(t, nil, err)
	})

	t.Run("Corner case: sections are empty", func(t *testing.T) {
		parser := InitParser()

		_, err := parser.GetSectionNames()

		assertEquality(t, "this parser has no sections", err.Error())
	})
}

func TestGetSections(t *testing.T) {
	t.Run("Normal case: sections are not empty", func(t *testing.T) {

		parser := InitParser()
		err := parser.LoadFromFile(path)
		if err != nil {
			t.Errorf("Error! %v", err)
		}

		got, err := parser.GetSections()

		expected := populateExpectedNormal(t)
		assertEquality(t, expected, got)
		assertEquality(t, nil, err)
	})

	t.Run("Corner case: sections are empty", func(t *testing.T) {

		parser := InitParser()

		_, err := parser.GetSections()
		assertEquality(t, "this parser has no sections", err.Error())
	})

}

func TestGet(t *testing.T) {

	testcases := []struct {
		testcaseName string
		sectionName  string
		key          string
		expected     string
		err          string
	}{
		{
			testcaseName: "Normal case: section and key are present",
			sectionName:  "database",
			key:          "server",
			expected:     "192.0.2.62",
			err:          "",
		},
		{
			testcaseName: "corner case: section not found",
			sectionName:  "user",
			key:          "server",
			err:          "section not found",
		},
		{
			testcaseName: "corner case: key not found",
			sectionName:  "database",
			key:          "size",
			err:          "key not found",
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
			if testcase.err == "" {
				assertEquality(t, testcase.expected, got)
			} else {
				assertEquality(t, testcase.err, err.Error())
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
		err          string
	}{
		{
			testcaseName: "Normal case: section and key are present",
			sectionName:  "database",
			key:          "server",
			value:        "127.0.0.1",
			err:          "",
		},
		{
			testcaseName: "corner case: section not found",
			sectionName:  "user",
			key:          "server",
			err:          "section not found",
		},
		{
			testcaseName: "corner case: key not found",
			sectionName:  "database",
			key:          "size",
			err:          "key not found",
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
			if testcase.err == "" {
				value := parser.sections[testcase.sectionName].map_[testcase.key]
				assertEquality(t, testcase.value, value)
			} else {
				assertEquality(t, testcase.err, err.Error())
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
