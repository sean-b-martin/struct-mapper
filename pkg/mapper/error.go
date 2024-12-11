package mapper

import "errors"

var (
	ErrNotStruct   = errors.New("attribute is not a struct")
	ErrMissingTag  = errors.New("missing tag")
	ErrIgnoreField = errors.New("ignore field")
)
