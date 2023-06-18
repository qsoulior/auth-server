package jwt

import (
	"reflect"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestGetSigningMethod(t *testing.T) {
	type args struct {
		alg string
	}
	tests := []struct {
		name    string
		args    args
		want    jwt.SigningMethod
		wantErr bool
	}{
		{"ValidAlg", args{"ES256"}, jwt.SigningMethodES256, false},
		{"InvalidAlg", args{"ES255"}, nil, true},
		{"NoneAlg", args{"none"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSigningMethod(tt.args.alg)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSigningMethod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSigningMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}
