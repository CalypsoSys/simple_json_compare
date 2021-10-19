// Package simple_json_compare implements a very simplicstic json compare function
//
package simple_json_compare

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// CompareJSONFiles - Compare JSON from 2 file sources
func CompareJSONFiles(leftPath string, rightPath string, ignorePaths []string) (bool, []string, error) {
	leftJSON, err := loadFlatJSONFile(leftPath)
	if err != nil {
		return false, nil, err
	}
	rightJSON, err := loadFlatJSONFile(rightPath)
	if err != nil {
		return false, nil, err
	}

	return compareJSON(leftJSON, rightJSON, ignorePaths)
}

// CompareJSONStrings - Compare JSON from 2 strings
func CompareJSONStrings(leftString string, rightString string, ignorePaths []string) (bool, []string, error) {
	leftJSON, err := loadFlatJSONString(leftString)
	if err != nil {
		return false, nil, err
	}
	rightJSON, err := loadFlatJSONString(rightString)
	if err != nil {
		return false, nil, err
	}

	return compareJSON(leftJSON, rightJSON, ignorePaths)
}

// CompareJSONBytes - Compare JSON from 2 byte arrays
func CompareJSONBytes(leftBytes []byte, rightBytes []byte, ignorePaths []string) (bool, []string, error) {
	leftJSON, err := loadFlatJSONBytes(leftBytes)
	if err != nil {
		return false, nil, err
	}
	rightJSON, err := loadFlatJSONBytes(rightBytes)
	if err != nil {
		return false, nil, err
	}

	return compareJSON(leftJSON, rightJSON, ignorePaths)
}

func compareJSON(leftJSON []string, rightJSON []string, ignorePaths []string) (bool, []string, error) {
	foundInLeft := make([]bool, len(leftJSON))
	foundInRight := make([]bool, len(rightJSON))
	for leftIndex, leftElm := range leftJSON {
		for rightIndex, rightElm := range rightJSON {
			if foundInRight[rightIndex] == false && leftElm == rightElm {
				foundInLeft[leftIndex] = true
				foundInRight[rightIndex] = true
				break
			}
		}
	}

	var differences []string
	diff := false
	for index, element := range leftJSON {
		if foundInLeft[index] == false && isIgnored(element, ignorePaths) == false {
			differences = append(differences, fmt.Sprintf("left: %s", element))
			diff = true
		}
	}
	for index, element := range rightJSON {
		if foundInRight[index] == false && isIgnored(element, ignorePaths) == false {
			differences = append(differences, fmt.Sprintf("right: %s", element))
			diff = true
		}
	}

	return diff, differences, nil
}

func isIgnored(element string, ignorePaths []string) bool {
	for _, ignore := range ignorePaths {
		if element == ignore {
			return true
		} else if strings.HasSuffix(ignore, "->*") {
			le := strings.LastIndex(element, "->")

			if element[:le] == ignore[:len(ignore)-3] {
				return true
			}
		} else if strings.Contains(ignore, "->*R[") && strings.HasSuffix(ignore, "]") {
			li := strings.LastIndex(ignore, "->*R[")
			if len(element) > li && element[:li] == ignore[:li] {
				if ok, _ := regexp.MatchString(fmt.Sprintf("^%s$", ignore[li+5:len(ignore)-1]), element[li:]); ok {
					return true
				}
			}
		}
	}

	return false
}

func loadFlatJSONFile(file string) ([]string, error) {
	jsonFile, err := os.Open(file)
	// if we os.Open returns an error then handle it
	if err != nil {
		return nil, fmt.Errorf("cannot open payload file %s\n%v+", file, err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read payload file %s\n%v+", file, err)
	} else {
		var flatJSON []string
		err = processNode(byteValue, "root", &flatJSON)
		return flatJSON, err
	}
}

func loadFlatJSONString(json string) ([]string, error) {
	var flatJSON []string
	err := processNode([]byte(json), "root", &flatJSON)
	return flatJSON, err
}

func loadFlatJSONBytes(json []byte) ([]string, error) {
	var flatJSON []string
	err := processNode(json, "root", &flatJSON)
	return flatJSON, err
}

func processNode(jsonBytes []byte, path string, flatJSON *[]string) error {
	var mapNode map[string]interface{}
	if json.Unmarshal(jsonBytes, &mapNode) == nil {
		for key, element := range mapNode {
			currentPath := fmt.Sprintf("%s->%s", path, key)
			*flatJSON = append(*flatJSON, currentPath)
			if val, isValue := checkType(element); !isValue {
				jsonData, err := json.Marshal(element)
				if err == nil {
					err = processNode(jsonData, currentPath, flatJSON)
				}
				if err != nil {
					return fmt.Errorf("json map error path %s\n%v+", currentPath, err)
				}
			} else {
				*flatJSON = append(*flatJSON, fmt.Sprintf("%s->%s", currentPath, val))
			}
		}

		return nil
	}
	var arrayNode []interface{}
	if json.Unmarshal(jsonBytes, &arrayNode) == nil {
		for index, element := range arrayNode {
			currentPath := fmt.Sprintf("%s->%d", path, index)
			*flatJSON = append(*flatJSON, currentPath)
			if val, isValue := checkType(element); !isValue {
				jsonData, err := json.Marshal(element)
				if err == nil {
					err = processNode(jsonData, fmt.Sprintf("%s->*", path), flatJSON)
				}
				if err != nil {
					return fmt.Errorf("json array error path %s\n%v+", currentPath, err)
				}
			} else {
				*flatJSON = append(*flatJSON, fmt.Sprintf("%s->%s", currentPath, val))
			}
		}
		return nil
	}

	return fmt.Errorf("json critical error path %s", path)
}

func checkType(element interface{}) (string, bool) {
	switch val := element.(type) {
	case int8, int16, int32, int64, int:
		return strconv.FormatInt(reflect.ValueOf(val).Int(), 10), true
	case uint8, uint16, uint32, uint64:
		return strconv.FormatUint(reflect.ValueOf(val).Uint(), 10), true
	case string:
		return reflect.ValueOf(val).String(), true
	case bool:
		return strconv.FormatBool(bool(val)), true
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 64), true
	case float64:
		return strconv.FormatFloat(float64(val), 'f', -1, 64), true
	}

	return "", false
}
