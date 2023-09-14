package json

import (
	"encoding/json"
	"strconv"
)

func NumberType(value any) (Number, bool) {
	number, ok := value.(json.Number)
	return Number{Number: number}, ok
}

type Number struct {
	json.Number
}

func (number Number) Int() int {
	value, _ := number.Int64()
	return int(value)
}

func (number Number) Normalize() any {
	str := number.String()

	_int64, err := strconv.ParseInt(str, 10, 64)
	if err == nil {
		return _int64
	}

	_float64, err := strconv.ParseFloat(str, 64)
	if err == nil {
		return _float64
	}

	return 0
}
