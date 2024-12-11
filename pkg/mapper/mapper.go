package mapper

import (
	"errors"
	"fmt"
	"reflect"
)

type StructMapper interface {
}

type ExtractTagIDFunc func(field reflect.StructField) (string, error)

type SimpleMapper struct {
	tagIDFunc ExtractTagIDFunc
	cache     map[reflect.Type]map[string]string
}

func NewSimpleMapper(defaultTag string) SimpleMapper {
	return SimpleMapper{tagIDFunc: getIDFromTag(defaultTag), cache: make(map[reflect.Type]map[string]string)}
}

func (mapper *SimpleMapper) RegisterStructF(s any, tagIDFunc ExtractTagIDFunc) error {
	if reflect.TypeOf(s).Kind() != reflect.Struct {
		return fmt.Errorf("s-type %T: %w", s, ErrNotStruct)
	}

	if tagIDFunc == nil {
		tagIDFunc = mapper.tagIDFunc
	}

	sType := reflect.TypeOf(s)
	cache := make(map[string]string)

	for i := 0; i < sType.NumField(); i++ {
		field := sType.Field(i)

		if field.Anonymous {
			// TODO
			continue
		}

		tagID, err := tagIDFunc(field)
		if err != nil {
			if errors.Is(err, ErrIgnoreField) {
				continue
			}
			return err
		}
		cache[tagID] = field.Name
	}
	mapper.cache[sType] = cache
	return nil
}

func (mapper *SimpleMapper) RegisterStruct(s any) error {
	return mapper.RegisterStructF(s, nil)
}
