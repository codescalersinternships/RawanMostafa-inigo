# INI Config Parser

This repository implements an ini config parser

## In this README 👇

- [Features](#features)
- [Usage](#usage)

## Features
 - `InitParser()`  returns a parser object to call our APIs
 - `LoadFromString()`  takes a string ini configs and parses them into the caller iniParser object
 - `LoadFromFile()`  takes an ini file path and parses it into the caller iniParser object
 - `GetSectionNames()`  returns an array of strings having the names of the sections of the caller iniParser object
 - `GetSections()`  returns a representing the parsed data structure of the caller iniParser object
 - `Get()` takes a sectionName and a key and returns the value of this key and an error if found
 - `Set()` takes a sectionName, a key and a value, it sets the passed key with the passed value and returns an error if found
 - `ToString()` returns the parsed ini map of the caller object as one string
 - `SaveToFile()` takes a path to an output file, saves the parsed ini map of the caller object into a file and returns an error if found

## Usage

1. 
    ```go
        import github.com/codescalersinternships/RawanMostafa-inigo
    ```

2. Initialize the parser first
    ```go
        parser := InitParser()
    ```

3. Example usage:
    ```go
        parser.LoadFromString(s)
    ```
    ```go
    	parser.LoadFromFile(filepath)
    ```
    ```go
    	names := parser.GetSectionNames()
    ```
    ```go
    	sections := parser.GetSections()
    ```
    ```go
    	value, err := parser.Get(sectionName, key)
    ```
    ```go
    	err = parser.Set(sectionName, key, value)
    ```
    ```go
    	s := parser.ToString()
    ```
    ```go
    	err = parser.SaveToFile(outPath)
    ```



