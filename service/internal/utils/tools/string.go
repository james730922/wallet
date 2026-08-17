package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

//ConcatStrings concatenate many strings to one string
func ConcatStrings(ss ...string) string {
	var buf bytes.Buffer
	for _, s := range ss {
		buf.WriteString(s)
	}
	return buf.String()
}

func ArrayInt64ToString(slice []int64, delim string, removeBracket bool) string {
	if removeBracket {
		return strings.Trim(strings.Replace(fmt.Sprint(slice), " ", delim, -1), "[]")
	}
	return strings.Trim(strings.Replace(fmt.Sprint(slice), " ", delim, -1), "")
}

func StringToArrayInt64(a string) ([]int64, error) {
	var res []int64
	err := json.Unmarshal([]byte(a), &res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func ArrayStringToString(slice []string, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(slice), " ", delim, -1), "")
}
