package jwt

import (
	"reflect"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func Test_publicKey_HMAC(t *testing.T) {
	p := publicKey{[]byte{1}}
	want := p.data
	got, err := p.HMAC()

	if err != nil {
		t.Errorf("publicKey.HMAC() error = %v, wantErr %v", err, false)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("publicKey.HMAC() = %v, want %v", got, want)
	}
}

func Test_publicKey_RSA(t *testing.T) {
	validData := []byte("-----BEGIN PUBLIC KEY-----\nMFswDQYJKoZIhvcNAQEBBQADSgAwRwJAbDPtajJjN+ZwPot6DR0HimeO0A/j0Ozqkp0pjq99OcbUtMFa78//ileYf1Kllracdvrv9/Alv3j+s7o+HkTY3wIDAQAB\n-----END PUBLIC KEY-----")
	invalidData := []byte{0}
	wantKey, _ := jwt.ParseRSAPublicKeyFromPEM(validData)
	wantNil, _ := jwt.ParseRSAPublicKeyFromPEM(invalidData)

	tests := []struct {
		name    string
		p       publicKey
		want    any
		wantErr bool
	}{
		{"ValidData", publicKey{validData}, wantKey, false},
		{"InvalidData", publicKey{invalidData}, wantNil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.RSA()
			if (err != nil) != tt.wantErr {
				t.Errorf("publicKey.RSA() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("publicKey.RSA() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_publicKey_ECDSA(t *testing.T) {
	validData := []byte("-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE9pyVw0cc766PaQXgPxUvzw4gFp0oFlciHIP3uUFspJl4vJDEOndIysnIgx1ox4GVzLLASmtyJlLdgk6Xd3XxdQ==\n-----END PUBLIC KEY-----")
	invalidData := []byte{0}
	wantKey, _ := jwt.ParseECPublicKeyFromPEM(validData)

	type args struct {
		bitSize int
	}

	tests := []struct {
		name    string
		p       publicKey
		args    args
		want    any
		wantErr bool
	}{
		{"ValidData", publicKey{validData}, args{256}, wantKey, false},
		{"InvalidData", publicKey{invalidData}, args{256}, nil, true},
		{"InvalidSize", publicKey{validData}, args{384}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.ECDSA(tt.args.bitSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("publicKey.ECDSA() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("publicKey.ECDSA() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_publicKey_Ed25519(t *testing.T) {
	p := publicKey{[]byte{0}}
	got, err := p.Ed25519()

	if err == nil {
		t.Errorf("publicKey.Ed25519() error = %v, wantErr %v", err, true)
		return
	}
	if got != nil {
		t.Errorf("publicKey.Ed25519() = %v, want %v", got, nil)
	}
}

func Test_privateKey_HMAC(t *testing.T) {
	p := privateKey{[]byte{1}}
	want := p.data
	got, err := p.HMAC()

	if err != nil {
		t.Errorf("privateKey.HMAC() error = %v, wantErr %v", err, false)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("privateKey.HMAC() = %v, want %v", got, want)
	}
}

func Test_privateKey_RSA(t *testing.T) {
	validData := []byte("-----BEGIN RSA PRIVATE KEY-----\nMIIBOAIBAAJAelRocsy1yHdWVo5l8R31edg7oRdtPdJtkNGKO+CxYnbq1s55OlS5Wos6eyfW8tcNn3u/khfFvBrDIhzQczXKSwIDAQABAkAmQ06kUdmuSX2M9191isxkfzkvixdVVgOFX7VgQ0jYpjBy+4XsvN1Wg8CFSVUZww9IBR58ONmQ3oqmIVzYV+iJAiEAwNosYfgQUqY8lX+C22ot3GlekzcpzZNbzayieSd56pUCIQCiYrBXXwFzKbwRSDNFTT0xH3reFgbchpuF0Qr1BzaJXwIgQPe+x+pHpXA3LK3eKYilloEwyStmO8kOvkUQHvx7h9kCIF/EqlFtA5Liyzq6BRrbGbqt4S23eeZ3MKO0DK1DysMrAiBtkZXXYQnBmd6LkQSM3IRj7W3JG/NDWD3U/4m7Xj4awA==\n-----END RSA PRIVATE KEY-----")
	invalidData := []byte{0}
	wantKey, _ := jwt.ParseRSAPrivateKeyFromPEM(validData)
	wantNil, _ := jwt.ParseRSAPrivateKeyFromPEM(invalidData)

	tests := []struct {
		name    string
		p       privateKey
		want    any
		wantErr bool
	}{
		{"ValidData", privateKey{validData}, wantKey, false},
		{"InvalidData", privateKey{invalidData}, wantNil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.RSA()
			if (err != nil) != tt.wantErr {
				t.Errorf("privateKey.RSA() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("privateKey.RSA() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_privateKey_ECDSA(t *testing.T) {
	validData := []byte("-----BEGIN EC PRIVATE KEY-----\nMHcCAQEEID7b3BCBKIOd9s5XUFQV3uGLKMpELITxNN+7JJbe2JifoAoGCCqGSM49AwEHoUQDQgAEFg4ATkXslQmgWKijnsmfPTxBMeA0yhhgYqIncASN5OFZXYxUKgZkozECV7Uuk5izTlR9GaaQJdEVlM7D0zpsEw==\n-----END EC PRIVATE KEY-----")
	invalidData := []byte{0}
	wantKey, _ := jwt.ParseECPrivateKeyFromPEM(validData)

	type args struct {
		bitSize int
	}
	tests := []struct {
		name    string
		p       privateKey
		args    args
		want    any
		wantErr bool
	}{
		{"ValidData", privateKey{validData}, args{256}, wantKey, false},
		{"InvalidData", privateKey{invalidData}, args{256}, nil, true},
		{"InvalidSize", privateKey{validData}, args{384}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.ECDSA(tt.args.bitSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("privateKey.ECDSA() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("privateKey.ECDSA() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_privateKey_Ed25519(t *testing.T) {
	p := privateKey{[]byte{0}}
	got, err := p.Ed25519()

	if err == nil {
		t.Errorf("privateKey.Ed25519() error = %v, wantErr %v", err, true)
		return
	}
	if got != nil {
		t.Errorf("privateKey.Ed25519() = %v, want %v", got, nil)
	}
}

func Test_keyParser_Parse(t *testing.T) {
	nilRSA, _ := jwt.ParseRSAPublicKeyFromPEM(nil)

	type args struct {
		alg string
	}
	tests := []struct {
		name    string
		p       keyParser
		args    args
		want    any
		wantErr bool
	}{
		{"ValidAlg", keyParser{publicKey{[]byte{0}}}, args{"HS256"}, []byte{0}, false},
		{"InvalidAlg", keyParser{publicKey{}}, args{"HS255"}, nil, true},
		{"InvalidKeyRSA", keyParser{publicKey{}}, args{"RS256"}, nilRSA, true},
		{"InvalidKeyECDSA", keyParser{publicKey{}}, args{"ES256"}, nil, true},
		{"InvalidKeyEd25519", keyParser{publicKey{}}, args{"EdDSA"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.Parse(tt.args.alg)
			if (err != nil) != tt.wantErr {
				t.Errorf("keyParser.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("keyParser.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParsePublicKey(t *testing.T) {
	want := []byte{0}
	got, err := ParsePublicKey(want, "HS256")
	if err != nil {
		t.Errorf("ParsePublicKey() error = %v, wantErr %v", err, false)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ParsePublicKey() = %v, want %v", got, want)
	}
}

func TestParsePrivateKey(t *testing.T) {
	want := []byte{0}
	got, err := ParsePrivateKey(want, "HS256")
	if err != nil {
		t.Errorf("ParsePrivateKey() error = %v, wantErr %v", err, true)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ParsePrivateKey() = %v, want %v", got, nil)
	}
}
