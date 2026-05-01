package handler

import (
	"net/http"
	"reflect"
	"testing"
)

func Test_redirect(t *testing.T) {
	type args struct {
		repo map[string]string
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := redirect(tt.args.repo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("redirect() = %v, want %v", got, tt.want)
			}
		})
	}
}
