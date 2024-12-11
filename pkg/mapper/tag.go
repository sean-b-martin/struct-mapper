package mapper

import (
	"fmt"
	"reflect"
	"strings"
)

func getIDFromTag(tag string) ExtractTagIDFunc {
	return func(field reflect.StructField) (string, error) {
		tags, ok := field.Tag.Lookup(tag)
		if !ok {
			return "", fmt.Errorf("field %s: %w", field.Name, ErrMissingTag)
		}

		tagID := strings.SplitN(tags, ",", 2)[0]
		if tagID == "" || tagID == "-" {
			return "", fmt.Errorf("field %s: %w", field.Name, ErrIgnoreField)
		}

		return tagID, nil
	}
}
