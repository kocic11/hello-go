package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"reflect"
	"strings"
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
			want:    "200 OK",
			wantErr: false,
		},
		{
			name: "Test2",
			args: args{
				input: &Input{
					XStorageUser: "Storage-gse00013735:cloud.admin",
					XStoragePass: "WrongPassword",
					Content:      "{\"name\":\"Test123\"}"},
			},
			want:    "401 Unauthorized",
			wantErr: true,
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
		content string
		token   string
	}
	tests := []struct {
		name  string
		args  args
		want  *http.Response
		want1 string
	}{
		{
			name: "Test1",
			args: args{
				content: "{\"name\":\"Test123\"}",
				token:   "AUTH_tk98be995c2ce3eba6189ffa828be0bfb8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := putObject(tt.args.content, tt.args.token)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("putObject() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("putObject() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_getObject(t *testing.T) {
	type args struct {
		url   string
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Response
		wantErr bool
	}{
		{
			name: "getObject_Test1",
			args: args{
				url:   "https://gse00013735.storage.oraclecloud.com/v1/Storage-gse00013735/fn_container/log1522341159078238800.json",
				token: "AUTH_tkbab23e5e15dc1fefd4be57def2bfe076",
			},
			want:    &http.Response{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getObject(tt.args.url, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("getObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_myHandler(t *testing.T) {
	type args struct {
		ctx context.Context
		in  io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
	}{
		{
			name: "myHandler_test1",
			args: args{
				ctx: nil,
				in: strings.NewReader(`
					{
						"XStorageUser":"Storage-gse00013735:cloud.admin",
						"XStoragePass":"ImmaNeNt@0SumP",
						"Content":{"name":"Test123"}
					}
				`),
			},
			wantOut: "",
		},
		{
			name: "myHandler_test2",
			args: args{
				ctx: nil,
				in: strings.NewReader(`
					{
						"XStorageUser": "Storage-gse00013735:cloud.admin",
						"XStoragePass": "ImmaNeNt@0SumP",
						"Content": 
							{
								"Exception": {
									"ApplicationID": "str1234",
									"ExceptionMessage": "str1234",
									"ExceptionNumber": "str1234",
									"ExceptionSeverity": "str1234",
									"ExceptionType": "str1234",
									"ServiceName": "str1234",
									"ServiceOperationName": "str1234",
									"ExceptionContextList": {
										"ExceptionContext": {
											"Key": "str1234",
											"Value": "str1234"
										}
									}
								}
							}
					}
				`),
			},
			wantOut: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			myHandler(tt.args.ctx, tt.args.in, out)
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("myHandler() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
