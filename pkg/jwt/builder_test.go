package jwt

import (
	"reflect"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestNewBuilder(t *testing.T) {
	type args struct {
		params Params
	}
	tests := []struct {
		name    string
		args    args
		want    *builder
		wantErr bool
	}{
		{"ValidAlg", args{Params{"", "HS256", nil}}, &builder{Params{"", "HS256", nil}, jwt.SigningMethodHS256}, false},
		{"InvalidAlg", args{Params{"", "HS255", nil}}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBuilder(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBuilder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBuilder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_builder_Build(t *testing.T) {
	b := &builder{Params{"", "HS256", []byte("secret")}, jwt.SigningMethodHS256}

	_, err := b.Build("", 10, "", nil)
	if err != nil {
		t.Errorf("builder.Build() error = %v, wantErr %v", err, false)
		return
	}
}
