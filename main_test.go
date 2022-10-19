package main

import "testing"

func Test_getGreeting(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "greeting", want: "Hello, world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getGreeting(); got != tt.want {
				t.Errorf("getGreeting() = %v, want %v", got, tt.want)
			}
		})
	}
}
