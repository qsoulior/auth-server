package uuid

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	_, err := New()
	if err != nil {
		t.Errorf("New() error = %v, wantErr %v", err, false)
		return
	}
}

func TestUUID_String(t *testing.T) {
	u := UUID{57, 146, 39, 158, 92, 80, 78, 141, 145, 22, 41, 155, 253, 255, 11, 95}
	want := "3992279e-5c50-4e8d-9116-299bfdff0b5f"

	if got := u.String(); got != want {
		t.Errorf("UUID.String() = %v, want %v", got, want)
	}
}

func TestUUID_MarshalJSON(t *testing.T) {
	u := UUID{57, 146, 39, 158, 92, 80, 78, 141, 145, 22, 41, 155, 253, 255, 11, 95}
	want := []byte{34, 51, 57, 57, 50, 50, 55, 57, 101, 45, 53, 99, 53, 48, 45, 52, 101, 56, 100, 45, 57, 49, 49, 54, 45, 50, 57, 57, 98, 102, 100, 102, 102, 48, 98, 53, 102, 34}
	wantErr := false

	got, err := u.MarshalJSON()
	if (err != nil) != wantErr {
		t.Errorf("UUID.MarshalJSON() error = %v, wantErr %v", err, wantErr)
		return
	}
	if !bytes.Equal(got, want) {
		t.Errorf("UUID.MarshalJSON() = %v, want %v", got, want)
	}
}

func TestUUID_Scan(t *testing.T) {
	type args struct {
		src any
	}
	tests := []struct {
		name    string
		u       *UUID
		args    args
		wantErr bool
	}{
		{"ValidSource", &UUID{}, args{"3992279e-5c50-4e8d-9116-299bfdff0b5f"}, false},
		{"InvalidSourceType", &UUID{}, args{0}, true},
		{"InvalidSourceString", &UUID{}, args{""}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.u.Scan(tt.args.src); (err != nil) != tt.wantErr {
				t.Errorf("UUID.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFromString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    UUID
		wantErr bool
	}{
		{"ValidString", args{"3992279e-5c50-4e8d-9116-299bfdff0b5f"}, UUID{57, 146, 39, 158, 92, 80, 78, 141, 145, 22, 41, 155, 253, 255, 11, 95}, false},
		{"InvalidStringType", args{""}, UUID{}, true},
		{"InvalidStringHex", args{"3992279g-5c50-4e8d-9116-299bfdff0b5f"}, UUID{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromString(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromBytes(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    UUID
		wantErr bool
	}{
		{"ValidBytes", args{[]byte{57, 146, 39, 158, 92, 80, 78, 141, 145, 22, 41, 155, 253, 255, 11, 95}}, UUID{57, 146, 39, 158, 92, 80, 78, 141, 145, 22, 41, 155, 253, 255, 11, 95}, false},
		{"InvalidBytes", args{[]byte{57, 146, 39, 158, 92, 80, 78, 141}}, UUID{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromBytes(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FromBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
