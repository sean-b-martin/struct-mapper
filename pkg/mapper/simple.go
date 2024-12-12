package mapper

import (
	"errors"
	"fmt"
	"reflect"
)

type mappingData struct {
	attributes map[string]string
	embedded   map[reflect.Type]mappingData
}

func newMappingData() mappingData {
	return mappingData{attributes: make(map[string]string), embedded: make(map[reflect.Type]mappingData)}
}

func (data *mappingData) addAttribute(key, value string) error {
	if _, ok := data.attributes[key]; ok {
		return ErrDuplicateID
	}
	data.attributes[key] = value

	return nil
}

type SimpleMapper struct {
	tagIDFunc ExtractTagIDFunc
	cache     map[reflect.Type]mappingData
}

func NewSimpleMapper(defaultTag string) SimpleMapper {
	return SimpleMapper{tagIDFunc: getIDFromTag(defaultTag), cache: make(map[reflect.Type]mappingData)}
}

func (mapper *SimpleMapper) RegisterStructF(s any, tagIDFunc ExtractTagIDFunc) error {
	if s == nil {
		return ErrNotStruct
	}

	if reflect.TypeOf(s).Kind() != reflect.Struct {
		return fmt.Errorf("s-type %T: %w", s, ErrNotStruct)
	}

	sType := reflect.TypeOf(s)
	data, err := mapper.processEmbeddedStruct(s, tagIDFunc)
	if err != nil {
		return err
	}

	mapper.cache[sType] = data
	return nil
}

func (mapper *SimpleMapper) processEmbeddedStruct(s any, tagIDFunc ExtractTagIDFunc) (mappingData, error) {
	if tagIDFunc == nil {
		tagIDFunc = mapper.tagIDFunc
	}

	sData := reflect.ValueOf(s)
	sType := reflect.TypeOf(s)
	data := newMappingData()

	for i := 0; i < sType.NumField(); i++ {
		field := sType.Field(i)

		if field.Anonymous {
			embeddedData, err := mapper.processEmbeddedStruct(sData.Field(i).Interface(), tagIDFunc)
			if err != nil {
				return data, err
			}

			data.embedded[field.Type] = embeddedData
			continue
		}

		tagID, err := tagIDFunc(field)
		if err != nil {
			if errors.Is(err, ErrIgnoreField) {
				continue
			}
			return data, err
		}

		if err := data.addAttribute(tagID, field.Name); err != nil {
			return data, err
		}
	}

	return data, nil
}

func (mapper *SimpleMapper) RegisterStruct(s any) error {
	return mapper.RegisterStructF(s, nil)
}
