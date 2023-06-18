package jwt

import (
	"reflect"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestNewParser(t *testing.T) {
	type args struct {
		params Params
	}
	tests := []struct {
		name    string
		args    args
		want    *parser
		wantErr bool
	}{
		{"ValidAlg", args{Params{"", "HS256", nil}}, &parser{Params{"", "HS256", nil}, jwt.SigningMethodHS256}, false},
		{"InvalidAlg", args{Params{"", "HS255", nil}}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewParser(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewParser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewParser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parser_Parse(t *testing.T) {
	type args struct {
		tokenString string
	}
	tests := []struct {
		name    string
		p       *parser
		args    args
		want    *Claims
		wantErr bool
	}{
		{"ValidString", &parser{Params{"", "HS256", []byte("secret")}, jwt.SigningMethodHS256}, args{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.Rq8IxqeX7eA6GgYxlcHdPFVRNFFZc5rEI3MQTZZbK3I"}, &Claims{"", nil, jwt.RegisteredClaims{"", "1234567890", nil, nil, nil, nil, ""}}, false},
		{"InvalidString", &parser{Params{"", "HS256", []byte("secret")}, jwt.SigningMethodHS256}, args{""}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.Parse(tt.args.tokenString)
			if (err != nil) != tt.wantErr {
				t.Errorf("parser.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parser.Parse() = %#v, want %v", got, tt.want)
			}
		})
	}
}
