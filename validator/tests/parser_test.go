package tests

import (
	"reflect"
	"testing"

	"github.com/tucats/validator"
)

func TestItem_Parser(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		tag     string
		kind    validator.Type
		want    validator.Item
		wantErr error
	}{
		{
			name:    "Invalid enum list",
			tag:     "enum=,required",
			wantErr: validator.ErrEmptyTagValue.Context("enum"),
		},
		{
			name: "Valid enum list",
			tag:  "enum=red|blue|green",
			kind: validator.TypeString,
			want: validator.Item{
				ValueType: validator.TypeString,
				Enums:     []string{"red", "blue", "green"},
			},
		},
		{
			name:    "Invalid min length",
			tag:     "required,minlength=bogus",
			wantErr: validator.ErrInvalidInteger.Context("minlength").Value("bogus"),
		},
		{
			name: "Valid max length",
			tag:  "required,maxlength=20",
			want: validator.Item{
				Required:     true,
				HasMaxLength: true,
				MaxLength:    20,
			},
			wantErr: nil,
		},
		{
			name: "Valid min and max length",
			tag:  "required,minlength=1,maxlength=20",
			want: validator.Item{
				Required:     true,
				HasMaxLength: true,
				MaxLength:    20,
				HasMinLength: true,
				MinLength:    1,
			},
			wantErr: nil,
		},
		{
			name:    "Invalid keyword",
			tag:     "invalid=keyword",
			wantErr: validator.ErrInvalidKeyword.Value("invalid"),
		},
		{
			name:    "Empty tag",
			tag:     "",
			wantErr: validator.ErrEmptyTag,
		},
		{
			name:    "Missing keyword value",
			tag:     "invalid=",
			wantErr: validator.ErrEmptyTagValue.Context("invalid"),
		},
		{
			name: "Valid keyword",
			tag:  "required",
			want: validator.Item{
				Required: true,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := validator.Item{ValueType: tt.kind}

			gotErr := item.ParseTag(tt.tag)

			var msg1, msg2 string

			if tt.wantErr != nil {
				msg1 = tt.wantErr.Error()
			}

			if gotErr != nil {
				msg2 = gotErr.Error()
			}

			if msg1 != msg2 {
				t.Errorf("Parser() unexpected result: %v", gotErr)
			} else {
				if gotErr == nil && !reflect.DeepEqual(item, tt.want) {
					t.Errorf("Parser() = %+v, want %+v", item, tt.want)
				}
			}
		})
	}
}
