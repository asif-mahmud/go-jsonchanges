package gojsonchanges_test

import (
	"encoding/json"
	"fmt"

	gojsonchanges "github.com/asif-mahmud/go-jsonchanges"
)

func toJson(a any) string {
	d, _ := json.Marshal(a)
	return string(d)
}

func ExampleFindChanges_bothNull() {
	A := `null`
	B := `null`

	output, _ := gojsonchanges.FindChanges([]byte(A), []byte(B))

	fmt.Println(toJson(output))
	// Output: null
}

func ExampleFindChanges_aNull() {
	A := `null`
	B := `1.0`

	output, _ := gojsonchanges.FindChanges([]byte(A), []byte(B))

	fmt.Println(toJson(output))
	// Output: 1
}

func ExampleFindChanges_bNull() {
	A := `1.0`
	B := `null`

	output, _ := gojsonchanges.FindChanges([]byte(A), []byte(B))

	fmt.Println(toJson(output))
	// Output: 1
}

func ExampleFindChanges_nested() {
	A := `{"a":1,"b":{"c":1},"c":[1]}`
	B := `{"a":1,"b":{"c":1,"d":2},"c":[1,2]}`

	output, _ := gojsonchanges.FindChanges([]byte(A), []byte(B))

	fmt.Println(toJson(output))
	// Output: {"b":{"d":2},"c":[null,2]}
}
