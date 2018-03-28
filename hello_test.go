package main

import (
	"net/http"
	"reflect"
	"testing"
)

func Test_getToken(t *testing.T) {
	type args struct {
		input *Input
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test1",
			args: args{
				input: &Input{
					XStorageUser: "Storage-gse00013735:cloud.admin",
					XStoragePass: "ImmaNeNt@0SumP",
					Content:      "{\"name\":\"Test123\"}"},
			},
			want:    "AUTH_tk3fd53f5799fc9d273cb559a98c523536",
			wantErr: false,
		},
		{
			name: "Test2",
			args: args{
				input: &Input{
					XStorageUser: "Storage-gse00013735:cloud.admin",
					XStoragePass: "wrongpass",
					Content:      "{\"name\":\"Test123\"}"},
			},
			want:    "AUTH_tkec85e46b6f6e28bcdedabccf583126d7",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getToken(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("getToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_putObject(t *testing.T) {
	type args struct {
		o string
		t *Token
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Response
		wantErr bool
	}{
		{
			name: "Test1",
			args: args{
				o: "{\"name\":\"Test123\"}",
				t: &Token{
					Token: "AUTH_tk98be995c2ce3eba6189ffa828be0bfb8",
					Error: "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := putObject(tt.args.o, tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("putObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("putObject() = %v, want %v", got, tt.want)
			}
		})
	}
}
