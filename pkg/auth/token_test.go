package auth

import (
	"reflect"
	"testing"
	"time"
)

func TestCreateToken(t *testing.T) {
	type args struct {
		ttl        time.Duration
		payload    interface{}
		privateKey string
	}
	var tests []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateToken(tt.args.ttl, tt.args.payload, tt.args.privateKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CreateToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	type args struct {
		token     string
		publicKey string
	}
	var tests []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateToken(tt.args.token, tt.args.publicKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}
