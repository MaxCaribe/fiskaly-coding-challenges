package domain

import (
	"reflect"
	"testing"
)

func TestAlgorithm_String(t *testing.T) {
	tests := []struct {
		name      string
		algorithm Algorithm
		want      string
	}{
		{
			"1 to ECC",
			Algorithm(1),
			"ECC",
		},
		{
			"2 to RSA",
			Algorithm(2),
			"RSA",
		},
		{
			"0 to Empty string",
			Algorithm(0),
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.algorithm.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseAlgorithm(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    Algorithm
		wantErr bool
	}{
		{
			"ECC to 1",
			args{s: "ECC"},
			Algorithm(1),
			false,
		},
		{
			"RSA to 2",
			args{s: "RSA"},
			Algorithm(2),
			false,
		},
		{
			"invalid value to 0",
			args{s: "SHA"},
			Algorithm(0),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAlgorithm(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAlgorithm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseAlgorithm() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlgorithm_MarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		algorithm Algorithm
		want      []byte
		wantErr   bool
	}{
		{
			"1 to ECC",
			Algorithm(1),
			[]byte(`"ECC"`),
			false,
		},
		{
			"2 to RSA",
			Algorithm(2),
			[]byte(`"RSA"`),
			false,
		},
		{
			"0 to Empty string",
			Algorithm(0),
			[]byte(`""`),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.algorithm.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlgorithm_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name      string
		algorithm Algorithm
		args      args
		wantErr   bool
	}{
		{
			"ECC to 1",
			Algorithm(1),
			args{data: []byte(`"ECC"`)},
			false,
		},
		{
			"RSA to 2",
			Algorithm(2),
			args{data: []byte(`"RSA"`)},
			false,
		},
		{
			"invalid value to 0",
			Algorithm(0),
			args{data: []byte(`"SHA"`)},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.algorithm.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAlgorithm_GenerateKeyPairsInBytesECC(t *testing.T) {
	keyPair, err := Algorithm(1).GenerateKeyPairsInBytes()
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(keyPair.PrivateKey) == 0 {
		t.Errorf("ECC private key is absent")
	}
	if len(keyPair.PublicKey) == 0 {
		t.Errorf("ECC public key is absent")
	}
}

func TestAlgorithm_GenerateKeyPairsInBytesRSA(t *testing.T) {
	keyPair, err := Algorithm(2).GenerateKeyPairsInBytes()
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(keyPair.PrivateKey) == 0 {
		t.Errorf("RSA private key is absent")
	}
	if len(keyPair.PublicKey) == 0 {
		t.Errorf("RSA public key is absent")
	}
}

func TestAlgorithm_GenerateKeyPairsInBytesInvalid(t *testing.T) {
	_, err := Algorithm(0).GenerateKeyPairsInBytes()
	if err == nil {
		t.Errorf("error should present for invalid algorithm")
	}
}
