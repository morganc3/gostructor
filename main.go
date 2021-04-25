package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"reflect"
	"strings"
)

type data struct {
	name     string
	contents map[string]interface{}
}

func main() {
	buff := bytes.NewBuffer([]byte(`{"name":"colby", "age": 26.4, "education":{"school":"UMD", "asd":"bsd"}}`))
	output := generateFromJSON(buff, "Person")
	fmt.Println(output)
}

func generateFromJSON(r io.Reader, name string) string {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatalln("Failed to read JSON bytes")
	}
	var jsonData map[string]interface{}
	err = json.Unmarshal(bytes, &jsonData)
	if err != nil {
		log.Fatal("failed to unmarshal")
	}

	d := data{name: name, contents: jsonData}
	queue := []*data{&d}
	return structify(queue)
}

func structify(queue []*data) string {
	result := ""
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		var currStruct string
		currStruct = fmt.Sprintf("type %s struct {\n", curr.name)

		for jsonKey, v := range curr.contents {
			t := getType(v)
			structKey := getCamelCase(jsonKey)
			typeString := t.String()
			if typeString == "map[string]interface {}" {
				key := strings.ToUpper(string(jsonKey[0]))
				if len(jsonKey) > 1 {
					key += jsonKey[1:]
				}
				queue = append(queue, &data{name: key, contents: v.(map[string]interface{})})
				typeString = key
			}
			if typeString == "[]interface {}" {
				// array
				if len(v.([]interface{})) > 0 {
					sliceElemType := getType(v.([]interface{})[0])
					if sliceElemType.String() == "map[string]interface {}" {
						key := strings.ToUpper(string(jsonKey[0]))
						if len(jsonKey) > 1 {
							key += jsonKey[1:]
						}
						queue = append(queue, &data{name: key, contents: v.(map[string]interface{})})
						typeString = key
					}
					typeString = "[]" + sliceElemType.String()
				}
			}

			currStruct += fmt.Sprintf("\t%s\t%s\t`json:\"%s\"`\n", structKey, typeString, jsonKey)
		}
		currStruct += "}\n\n"
		result += currStruct
	}
	return result
}

func processNext(queue []map[string]interface{}, name string) ([]map[string]interface{}, string) {
	// pop first item off the queue
	curr := queue[0]
	queue = queue[1:]
	out := fmt.Sprintf("type %s struct {\n", name)

	for jsonKey, v := range curr {
		t := getType(v)
		structKey := getCamelCase(jsonKey)
		typeString := t.String()
		if typeString == "map[string]interface {}" {
			queue = append(queue, v.(map[string]interface{}))
		}
		if typeString == "[]map[string]interface {}" {

		}
		if typeString == "[]interface {}" {
			// array
			if len(v.([]interface{})) > 0 {
				sliceElemType := getType(v.([]interface{})[0])
				typeString = "[]" + sliceElemType.String()
			}
		}
		out += fmt.Sprintf("\t%s\t%s\t`json:\"%s\"`\n", structKey, typeString, jsonKey)
	}

	out += "}\n\n"
	return queue, out
}

func getType(in interface{}) reflect.Type {
	t := reflect.TypeOf(in)
	ts := t.String()
	if ts == "float64" {
		numAsString := fmt.Sprint(in.(float64))
		if strings.Contains(numAsString, ".") {
			return t
		} else {
			return reflect.TypeOf(1)
		}
	}

	return t
}

func getCamelCase(str string) string {
	str = strings.ReplaceAll(str, "-", "_")
	ret := ""
	spl := strings.Split(str, "_")
	for _, s := range spl {
		switch len(s) {
		case 0:
			continue
		case 1:
			ret += strings.ToUpper(string(s[0]))
		default:
			ret += strings.ToUpper(string(s[0])) + s[1:]
		}
	}
	return ret
}
