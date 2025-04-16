package utils

import (
	"errors"
	"reflect"
	"slices"
	"strings"
)

const (
	omitemptyTag = "omitempty"
)

var (
	ErrorMustBeValue = errors.New("not null value but empty input")
)

func StructToMap(s any, full bool) (map[string]any, error) {
	res := map[string]any{}
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()

	for i := range v.NumField() {
		field := v.Field(i)
		if !field.CanInterface() {
			continue
		}
		tags := strings.Split(t.Field(i).Tag.Get("db"), ",")
		if tags[0] == "" {
			continue
		}

		if field.IsZero() {
			if full && !isOmitempty(tags) {
				return nil, ErrorMustBeValue
			}
			if !full || isOmitempty(tags) {
				continue
			}
		}
		value := field.Interface()
		res[tags[0]] = value
	}
	return res, nil
}

func isOmitempty(tags []string) bool {
	return slices.Contains(tags, omitemptyTag)
}
