package iniparser

import (
	"reflect"
	"testing"
)

func TestLoadFromString(t *testing.T) {
	s := `; last modified 1 April 2001 by John Doe
			[owner]
			name = John Doe
			organization = Acme Widgets Inc.

			[database]
			; use IP address in case network name resolution is not working
			server = 192.0.2.62
			port = 143`

	var expected map[string]section
	t.Run("populate expected map", func(t *testing.T) {
		t.Helper()
		expected = map[string]section{
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
	})

	parser := InitParser()
	parser.LoadFromString(s)

	if !reflect.DeepEqual(parser.sections, expected) {
		t.Errorf("Error! wanted:%v \n , got : %v \n", expected, parser.sections)
	}
}
