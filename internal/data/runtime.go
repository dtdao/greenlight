package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrorInvalidRuntimeFormat = errors.New("invalid runtime format")

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {

	jsonValue := fmt.Sprintf("%d mins", r)

	// need to be double quote for valid json string
	quotedJSONVValue := strconv.Quote(jsonValue)
	return []byte(quotedJSONVValue), nil
}

func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {

	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrorInvalidRuntimeFormat
	}

	parts := strings.Split(unquotedJSONValue, " ")

	if len(parts) != 2 || parts[1] != "mins" {
		return ErrorInvalidRuntimeFormat
	}

	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrorInvalidRuntimeFormat
	}

	*r = Runtime(i)
	return nil
}
