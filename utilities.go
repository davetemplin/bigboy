package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const dbTimeLayout = "2006-01-02 15:04:05.000"

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func warn(message string) {
	fmt.Fprintf(os.Stderr, message+"\n")
}

func stop(message string, code int) {
	fmt.Fprintf(os.Stderr, message+"\n")
	os.Exit(code)
}

func jsonWriteln(file *os.File, obj map[string]interface{}) error {
	line, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	line = bytes.Replace(line, []byte("\\u003c"), []byte("<"), -1)
	line = bytes.Replace(line, []byte("\\u003e"), []byte(">"), -1)
	line = bytes.Replace(line, []byte("\\u0026"), []byte("&"), -1)

	_, err = file.Write(line)
	if err != nil {
		return err
	}
	_, err = file.WriteString("\n")
	if err != nil {
		return err
	}
	return nil
}

func to_uint64(value interface{}) uint64 {
	switch value.(type) {
	case int:
		return uint64(value.(int))
	case int8:
		return uint64(value.(int8))
	case int16:
		return uint64(value.(int16))
	case int32:
		return uint64(value.(int32))
	case int64:
		return uint64(value.(int64))
	case uint:
		return uint64(value.(uint))
	case uint8:
		return uint64(value.(uint8))
	case uint16:
		return uint64(value.(uint16))
	case uint32:
		return uint64(value.(uint32))
	case uint64:
		return uint64(value.(uint64))
	default:
		panic("cannot convert value to uint64")
	}
}

func take(list []map[string]interface{}, key string) []uint64 {
	result := make([]uint64, 0)
	for _, obj := range list {
		result = append(result, to_uint64(obj[key]))
	}
	return result
}

func distinct(list []uint64) []uint64 {
	set := make(map[uint64]struct{})
	for _, val := range list {
		set[val] = struct{}{}
	}
	keys := make([]uint64, 0, len(set))
	for k := range set {
		keys = append(keys, k)
	}
	return keys
}

func csv(list []uint64) string {
	return strings.Replace(strings.Trim(fmt.Sprint(list), "[]"), " ", ",", -1)
}

func formatTimeDb(location *time.Location, t time.Time) string {
	return t.In(location).Format(dbTimeLayout)
}

func applyTimezone(location *time.Location, t time.Time) time.Time {
	if location != nil {
		value := t.Format(dbTimeLayout)
		result, _ := time.ParseInLocation(dbTimeLayout, value, location)
		return result
	}
	return t
}

func applyTimezoneAll(location *time.Location, list []map[string]interface{}) {
	for _, m := range list {
		for key := range m {
			if value, ok := m[key].(time.Time); ok {
				m[key] = applyTimezone(location, value)
			}
		}
	}
}

func parseTime(text string) (time.Time, error) {
	if text == "now" {
		return time.Now(), nil
	} else if text == "today" {
		return today(), nil
	} else if text == "yesterday" {
		return yesterday(), nil
	}

	t, err := time.Parse(time.RFC3339, text)
	if err == nil {
		return t, nil
	}

	t, err = time.Parse("2006-01-02 15:04:05", text)
	if err == nil {
		return t, nil
	}

	t, err = time.Parse("2006-01-02", text)
	if err == nil {
		return t, err
	}

	t, err = parseRelativeTime(text, time.Now())
	return t, err
}

func parseRelativeTime(text string, t time.Time) (time.Time, error) {
	var (
		u int // indicates the sign: +1 or -1
		s string
	)
	if strings.HasPrefix(text, "+") {
		u = 1
		s = text[1:]
	} else if strings.HasPrefix(text, "-") {
		u = -1
		s = text[1:]
	} else {
		u = -1 // if no +/- prefix then assume duration specifies a time in the past
		s = text
		text = "-" + text
	}

	if strings.HasSuffix(s, "d") {
		n, err := strconv.Atoi(strings.TrimSuffix(s, "d"))
		return t.AddDate(0, 0, u*n), err
	} else if strings.HasSuffix(s, "m") {
		n, err := strconv.Atoi(strings.TrimSuffix(s, "m"))
		return t.AddDate(0, u*n, 0), err
	} else if strings.HasSuffix(s, "y") {
		n, err := strconv.Atoi(strings.TrimSuffix(s, "y"))
		return t.AddDate(u*n, 0, 0), err
	}

	duration, err := time.ParseDuration(text)
	return t.Add(duration), err
}

func midnight(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func today() time.Time {
	return midnight(time.Now())
}

func yesterday() time.Time {
	return midnight(time.Now().AddDate(0, 0, -1))
}

func day(t time.Time) string {
	year, month, day := t.Date()
	result := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return result.String()[:10]
}

func contains(list []string, value string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}

func fileExists(name string) bool {
	path, err := os.Stat(name)
	if err == nil {
		return !path.IsDir()
	} else if os.IsNotExist(err) {
		return false
	}
	check(err)
	return false
}

func checkFileExists(path string) {
	if !fileExists(path) {
		stop(fmt.Sprintf("Cannot find file \"%s\"\n", path), 1)
	}
}
