package handler

import (
	"net/http"
	"testing"
)

func Test_server_shortenText(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		w http.ResponseWriter
		r *http.Request
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var s server
			s.shortenText(tt.w, tt.r)
		})
	}
}
