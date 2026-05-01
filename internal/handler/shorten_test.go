package handler

import (
	"drk-url-shortener/internal/config"
	"net/http"
	"reflect"
	"testing"
)

func Test_shortenText(t *testing.T) {
	type args struct {
		repo map[string]string
		cfg  config.Config
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
			if got := shortenText(tt.args.repo, tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("shortenText() = %v, want %v", got, tt.want)
			}
		})
	}
}
