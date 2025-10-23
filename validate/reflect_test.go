package validate

import "testing"

func Test_reflectOne(t *testing.T) {
	type User struct {
		Name string `valid:"min=1"`
		Age  int    `valid:"min=18"`
	}

	type Member struct {
		ID  string `json:"id" valid:"required"`
		Who User
	}

	tests := []struct {
		name    string
		tag     string
		object  any
		wantErr bool
	}{
		{
			name:    "nestedStruct",
			tag:     "",
			object:  &Member{},
			wantErr: false,
		},
		{
			name:    "simpleStruct",
			tag:     "",
			object:  &User{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := reflectOne(tt.name, tt.tag, tt.object)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("reflectOne() failed: %v", gotErr)
				}

				return
			}

			if tt.wantErr {
				t.Fatal("reflectOne() succeeded unexpectedly")
			}
		})
	}
}
