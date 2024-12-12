package mapper

import (
	"reflect"
)

type StructMapper interface {
	RegisterStructF(s any, tagIDFunc ExtractTagIDFunc) error
	RegisterStruct(s any) error
}

type ExtractTagIDFunc func(field reflect.StructField) (string, error)
