package data

import (
	"fmt"
	"strconv"
)

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {

	jsonValue := fmt.Sprintf("%d mins", r)

	// need to be double quote for valid json string
	quotedJSONVValue := strconv.Quote(jsonValue)
	return []byte(quotedJSONVValue), nil
}
