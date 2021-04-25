package main

import (
	"encoding/json"
	"log"
	"testing"
)

func TestStringUtils(t *testing.T) {
	in := "my-test-string"
	out := "MyTestString"
	if getCamelCase(in) != out {
		t.Fatalf("Got invalid camel case for %s: %s\n", in, getCamelCase(in))
	}

	in = "My_test_String"
	if getCamelCase(in) != out {
		t.Fatalf("Got invalid camel case for %s: %s\n", in, getCamelCase(in))
	}

	in = "My_test_String-"
	if getCamelCase(in) != out {
		t.Fatalf("Got invalid camel case for %s: %s\n", in, getCamelCase(in))
	}
	in = "_My_test_String-"
	if getCamelCase(in) != out {
		t.Fatalf("Got invalid camel case for %s: %s\n", in, getCamelCase(in))
	}
}

func TestGenerateStruct(t *testing.T) {
	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(`{"name":"colby", "age": 71}`), &jsonData)
	if err != nil {
		log.Fatal("failed to unmarshal")
	}
	convert("person", jsonData)
}
