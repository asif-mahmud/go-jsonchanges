package gojsonchanges

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"reflect"
)

func FindChanges(a, b []byte) (any, error) {
	if a == nil && b == nil {
		return nil, nil
	}

	if a == nil && b != nil {
		return unmarshal(bytes.NewReader(b))
	}

	if a != nil && b == nil {
		return unmarshal(bytes.NewReader(a))
	}

	// lets unmarshal, any of these ops can
	// fail so early exit
	left, err := unmarshal(bytes.NewReader(a))
	if err != nil {
		return nil, err
	}

	right, err := unmarshal(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	return findChanges(left, right)
}

func unmarshal(a io.Reader) (any, error) {
	var v any
	err := json.NewDecoder(a).Decode(&v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func findChanges(a, b any) (any, error) {
	// if any of the arguments are nil, then we have simple situation
	// left = nil, right = nil, return is nil
	// left = nil, right != nil, return is right
	// left != nil, right = nil, return is left
	if a == nil && b == nil {
		return nil, nil
	}

	if a == nil && b != nil {
		return b, nil
	}

	if a != nil && b == nil {
		return a, nil
	}

	// now we have both non-nil json values.
	ak := reflect.ValueOf(a).Kind()
	bk := reflect.ValueOf(b).Kind()

	// now we check if a & b are of different kind if so,
	// return b
	if ak != bk {
		return b, nil
	}

	// at this point a & b are of same type but,
	// we need to check their values
	switch ak {
	case reflect.Bool:
		av, bv := a.(bool), b.(bool)
		if av != bv {
			return b, nil
		}
		return nil, nil

	case reflect.String:
		av, bv := a.(string), b.(string)
		if av != bv {
			return b, nil
		}
		return nil, nil

	case reflect.Float64:
		av, bv := a.(float64), b.(float64)
		if av != bv {
			return b, nil
		}
		return nil, nil

	case reflect.Slice:
		av, bv := a.([]any), b.([]any)
		maxLen := math.Max(float64(len(av)), float64(len(bv)))
		rv := make([]any, int(maxLen))

		var ix int
		var avv any
		var nilCnt int

		for ix, avv = range av {
			if ix > len(bv)-1 {
				break
			}
			ch, er := findChanges(avv, bv[ix])
			if er != nil {
				return nil, er
			}
			if ch == nil {
				nilCnt++
			}
			rv[ix] = ch
		}

		if ix == len(bv) {
			for ; ix < len(av); ix++ {
				if av[ix] == nil {
					nilCnt++
				}
				rv[ix] = av[ix]
			}
		}

		ix++
		for ; ix < len(bv); ix++ {
			if bv[ix] == nil {
				nilCnt++
			}
			rv[ix] = bv[ix]
		}

		if nilCnt == int(maxLen) {
			return nil, nil
		}

		return rv, nil

	case reflect.Map:
		av, bv := a.(map[string]any), b.(map[string]any)
		rv := map[string]any{}

		for ak, avv := range av {
			bvv, ok := bv[ak]
			if !ok {
				rv[ak] = avv
				continue
			}

			ch, er := findChanges(avv, bvv)
			if er != nil {
				return nil, er
			}

			if ch == nil {
				continue
			}

			rv[ak] = ch
		}

		for bk, bvv := range bv {
			if _, ok := rv[bk]; ok {
				continue
			}

			if _, ok := av[bk]; ok {
				continue
			}

			rv[bk] = bvv
		}

		if len(rv) == 0 {
			return nil, nil
		}

		return rv, nil

	default:
		return nil, fmt.Errorf("findChanges: unknown type: %v", ak)
	}
}
