package handler

import (
	"drk-url-shortener/internal/repository"
	"net/http"
	"testing"
)

func Test_redirect(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		repo repository.Storage
		want http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := redirect(tt.repo)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("redirect() = %v, want %v", got, tt.want)
			}
		})
	}
}
