package mapper

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type a struct {
	SomeAttribute  string `json:"some-attribute"`
	SomeAttribute2 int    `json:"some-attribute2"`
}

func TestSimpleMapper_RegisterStruct(t *testing.T) {
	type missingTags struct {
		SomeAttribute string
	}

	type args struct {
		s any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		cache   map[reflect.Type]map[string]string
	}{
		{name: "invalid type string", args: args{s: "123"}, wantErr: true},
		{name: "invalid type number", args: args{s: 123}, wantErr: true},
		{name: "missing attribute", args: args{s: missingTags{}}, wantErr: true},
		{name: "valid type", args: args{s: a{}}, wantErr: false, cache: map[reflect.Type]map[string]string{
			reflect.TypeOf(a{}): {"some-attribute": "SomeAttribute", "some-attribute2": "SomeAttribute2"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapper := NewSimpleMapper("json")
			if err := mapper.RegisterStruct(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("RegisterStruct() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				assert.Equal(t, tt.cache, mapper.cache)
			}
		})
	}
}
