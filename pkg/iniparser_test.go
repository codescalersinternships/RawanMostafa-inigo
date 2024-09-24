package iniparser

import (
	"errors"
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

func populateExpectedNormal(t *testing.T) map[string]map[string]string {
	t.Helper()
	expected := map[string]map[string]string{
		"owner": {
			"name":         "John Doe",
			"organization": "Acme Widgets Inc.",
		},
		"database": {
			"server": "192.0.2.62",
			"port":   "143",
		},
		"section": {
			"key0": "val0",
			"key1": "val1",
		},
	}
	return expected
}

func populateExpectedEmptySection(t *testing.T) map[string]map[string]string {
	t.Helper()
	expected := map[string]map[string]string{
		"owner": {
			"name":         "John Doe",
			"organization": "Acme Widgets Inc.",
		},
		"database": {},
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

		parser := NewParser()
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
		parser := NewParser()
		parser.LoadFromString(emptySectionINI)

		expected := populateExpectedEmptySection(t)

		assertEquality(t, expected, parser.sections)
	})

	t.Run("test global key found", func(t *testing.T) {
		const IniFile = `; last modified 1 April 2001 by John Doe
				name = John Doe
				[owner]
				organization = Acme Widgets Inc.
	
				[database]`

		parser := NewParser()
		err := parser.LoadFromString(IniFile)
		assertEquality(t, err, ErrGlobalKey)
	})
}

func TestLoadFromFile(t *testing.T) {

	expectedFile := populateExpectedNormal(t)

	testcases := []struct {
		testcaseName string
		filePath     string
		expected     map[string]map[string]string
		err          error
	}{
		{
			testcaseName: "Normal case: ini file is present",
			filePath:     "testdata/file.ini",
			expected:     expectedFile,
			err:          nil,
		},
		{
			testcaseName: "corner case: not an ini file",
			filePath:     "testdata/file.txt",
			err:          ErrNotINI,
		},
		{
			testcaseName: "corner case: file not found",
			filePath:     "testdata/filex.ini",
			err:          ErrFileNotExist,
		},
	}
	for _, testcase := range testcases {

		t.Run(testcase.testcaseName, func(t *testing.T) {
			parser := NewParser()

			err := parser.LoadFromFile(testcase.filePath)
			if testcase.err == nil {
				assertEquality(t, expectedFile, parser.sections)
			}
			assertEquality(t, errors.Is(err, testcase.err), true)

		})
	}
}

func TestGetSectionNames(t *testing.T) {
	t.Run("Normal case: sections are not empty", func(t *testing.T) {

		parser := NewParser()
		err := parser.LoadFromFile(path)
		if err != nil {
			t.Errorf("Error! %v", err)
		}
		names, err := parser.GetSectionNames()

		expected := []string{"database", "owner", "section"}

		assertEquality(t, expected, names)
		assertEquality(t, nil, err)
	})

	t.Run("Corner case: no sections", func(t *testing.T) {

		parser := NewParser()
		err := parser.LoadFromString("")
		if err != nil {
			t.Errorf("Error! %v", err)
		}
		names, err := parser.GetSectionNames()

		expected := make([]string, 0)

		assertEquality(t, expected, names)
		assertEquality(t, nil, err)
	})

}

func TestGetSections(t *testing.T) {
	t.Run("Normal case: sections are not empty", func(t *testing.T) {

		parser := NewParser()
		err := parser.LoadFromFile(path)
		if err != nil {
			t.Errorf("Error! %v", err)
		}

		got, err := parser.GetSections()

		expected := populateExpectedNormal(t)
		assertEquality(t, expected, got)
		assertEquality(t, nil, err)
	})

	t.Run("Corner case: no sections", func(t *testing.T) {

		parser := NewParser()
		err := parser.LoadFromString("")
		if err != nil {
			t.Errorf("Error! %v", err)
		}
		names, err := parser.GetSections()

		expected := make(map[string]map[string]string, 0)

		assertEquality(t, expected, names)
		assertEquality(t, nil, err)
	})

}

func TestGet(t *testing.T) {

	testcases := []struct {
		testcaseName string
		sectionName  string
		key          string
		expected     string
		found        bool
	}{
		{
			testcaseName: "Normal case: section and key are present",
			sectionName:  "database",
			key:          "server",
			expected:     "192.0.2.62",
			found:        true,
		},
		{
			testcaseName: "corner case: section not found",
			sectionName:  "user",
			key:          "server",
			found:        false,
		},
		{
			testcaseName: "corner case: key not found",
			sectionName:  "database",
			key:          "size",
			found:        false,
		},
	}
	for _, testcase := range testcases {

		t.Run(testcase.testcaseName, func(t *testing.T) {
			parser := NewParser()
			err := parser.LoadFromFile(path)
			if err != nil {
				t.Errorf("Error! %v", err)
			}
			got, found := parser.Get(testcase.sectionName, testcase.key)
			assertEquality(t, testcase.expected, got)
			assertEquality(t, testcase.found, found)

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
			err:          ErrNoSection,
		},
		{
			testcaseName: "corner case: key not found",
			sectionName:  "database",
			key:          "size",
			err:          ErrNoKey,
		},
	}
	for _, testcase := range testcases {

		t.Run(testcase.testcaseName, func(t *testing.T) {
			parser := NewParser()
			err := parser.LoadFromFile(path)
			if err != nil {
				t.Errorf("Error! %v", err)
			}
			err = parser.Set(testcase.sectionName, testcase.key, testcase.value)
			if testcase.err == nil {
				value := parser.sections[testcase.sectionName][testcase.key]
				assertEquality(t, testcase.value, value)
			}
			assertEquality(t, errors.Is(err, testcase.err), true)
		})
	}

}

func TestString(t *testing.T) {
	parser1 := NewParser()
	parser2 := NewParser()

	parser1.LoadFromString(stringINI)
	got := parser1.String()

	parser2.LoadFromString(got)

	assertEquality(t, parser1.sections, parser2.sections)
	assertEquality(t, fmt.Sprint(parser1), fmt.Sprint(parser2))
}

func TestSaveToFile(t *testing.T) {

	parser := NewParser()
	err := parser.LoadFromFile(path)
	if err != nil {
		t.Errorf("Error! %v", err)
	}

	file, err := os.CreateTemp("", "out.ini")
	if err != nil {
		t.Errorf("Error! %v", err)
	}
	defer os.Remove(file.Name())

	err = parser.SaveToFile(file.Name())
	if err != nil {
		t.Errorf("Error! %v", err)
	}

	stringFile := parser.String()
	assertFile(t, file.Name(), stringFile)
}
