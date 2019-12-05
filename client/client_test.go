package kclient

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/nick-uk/ksql/ksql"
)

func TestNewClient(t *testing.T) {

	if got := NewClient("fakehost"); got == nil {
		t.Errorf("NewClient() = got: %+v, want: %v", got, nil)
	}
}

func TestNewClientWithTimeout(t *testing.T) {
	type args struct {
		host    string
		timeout time.Duration
	}

	var cancel context.CancelFunc

	tests := []struct {
		name  string
		args  args
		want  reflect.Type
		want1 reflect.Type
	}{
		{"Test ksql client with timeout", args{host: "whatever", timeout: time.Second * 1}, reflect.TypeOf(&ksql.Client{}), reflect.TypeOf(cancel)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := NewClientWithTimeout(tt.args.host, tt.args.timeout)

			if reflect.TypeOf(got) != tt.want {
				t.Errorf("NewClientWithTimeout() got: %+v, want: %+v", reflect.TypeOf(got), tt.want)
			}
			if reflect.TypeOf(got1) != tt.want1 {
				t.Errorf("NewClientWithTimeout() got: %+v, want: %+v", reflect.TypeOf(got1), tt.want1)
			}
		})
	}
}

func TestNewClientWithCancel(t *testing.T) {
	type args struct {
		host string
	}
	var cancel context.CancelFunc

	tests := []struct {
		name  string
		args  args
		want  reflect.Type
		want1 reflect.Type
	}{
		{"Test ksql client with cancel", args{host: "whatever"}, reflect.TypeOf(&ksql.Client{}), reflect.TypeOf(cancel)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := NewClientWithCancel(tt.args.host)

			if reflect.TypeOf(got) != tt.want {
				t.Errorf("NewClientWithCancel() got: %+v, want: %+v", reflect.TypeOf(got), tt.want)
			}
			if reflect.TypeOf(got1) != tt.want1 {
				t.Errorf("NewClientWithCancel() got: %+v, want: %+v", reflect.TypeOf(got1), tt.want1)
			}
		})
	}
}
