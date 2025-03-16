package gojsonchanges_test

import (
	"fmt"
	"testing"

	gojsonchanges "github.com/asif-mahmud/go-jsonchanges"
	"github.com/stretchr/testify/assert"
)

func TestFindChanges(t *testing.T) {
	type testCase struct {
		a, b        []byte
		expected    any
		errExpected bool
	}

	cases := []testCase{
		{nil, nil, nil, false},
		{nil, []byte(`1`), 1.0, false},
		{[]byte(`1`), nil, 1.0, false},

		{[]byte(``), nil, nil, true},
		{nil, []byte(``), nil, true},

		{[]byte(`1`), []byte(`true`), true, false},
		{[]byte(`1`), []byte(`"value"`), "value", false},
		{[]byte(`true`), []byte(`false`), false, false},
		{[]byte(`1.4`), []byte(`1.2`), 1.2, false},

		{[]byte(`[]`), []byte(`[]`), nil, false},
		{[]byte(`[1,2]`), []byte(`[]`), []any{1.0, 2.0}, false},
		{[]byte(`[1,2]`), []byte(`[1,2]`), nil, false},
		{[]byte(`[1,2]`), []byte(`[1]`), []any{nil, 2.0}, false},
		{[]byte(`[1,2]`), []byte(`[1.2,2]`), []any{1.2, nil}, false},
		{[]byte(`[1,2]`), []byte(`[1,2,3.4]`), []any{nil, nil, 3.4}, false},
		{[]byte(`[1,2]`), []byte(`[1.2,2,3.4]`), []any{1.2, nil, 3.4}, false},
		{[]byte(`[1,2]`), []byte(`[1.2,true]`), []any{1.2, true}, false},

		{[]byte(`{}`), []byte(`{}`), nil, false},
		{[]byte(`{"a":1}`), []byte(`{}`), map[string]any{"a": 1.0}, false},
		{[]byte(`{"a":1}`), []byte(`{"a":1}`), nil, false},
		{[]byte(`{"a":1}`), []byte(`{"a":null}`), map[string]any{"a": 1.0}, false},
		{
			[]byte(`{"a":1,"b":{"c":1}}`),
			[]byte(`{}`),
			map[string]any{"a": 1.0, "b": map[string]any{"c": 1.0}},
			false,
		},
		{
			[]byte(`{"a":1,"b":{"c":1}}`),
			[]byte(`{"a":1,"b":{"c":1,"d":2}}`),
			map[string]any{"b": map[string]any{"d": 2.0}},
			false,
		},
		{
			[]byte(`{"a":[1,2]}`),
			[]byte(`{"a":[1.2,2]}`),
			map[string]any{"a": []any{1.2, nil}},
			false,
		},
	}

	for ix, tc := range cases {
		t.Run(fmt.Sprintf("%d", ix), func(t *testing.T) {
			actual, er := gojsonchanges.FindChanges(tc.a, tc.b)
			if tc.errExpected {
				assert.NotNil(t, er)
				return
			}

			assert.Equal(t, tc.expected, actual)
		})
	}
}
