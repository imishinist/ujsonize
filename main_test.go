package main

import (
	"bytes"
	"testing"
)

func TestDecode(t *testing.T) {
	cases := []struct {
		name     string
		in       []byte
		expected string
	}{
		{
			name:     "empty",
			in:       []byte(`{}`),
			expected: "\n",
		},
		{
			name:     "single",
			in:       []byte(`{"foo":["bar"]}`),
			expected: "foo=bar\n",
		},
		{
			name:     "multiple",
			in:       []byte(`{"foo":["bar","baz","qux"]}`),
			expected: "foo=bar&foo=baz&foo=qux\n",
		},
		{
			name:     "multiple keys",
			in:       []byte(`{"foo":["bar","baz","qux"],"hoge":["fuga","piyo"]}`),
			expected: "foo=bar&foo=baz&foo=qux&hoge=fuga&hoge=piyo\n",
		},
		{
			name:     "primitive string not array",
			in:       []byte(`{"foo":"bar"}`),
			expected: "foo=bar\n",
		},
		{
			name:     "primitive number not array",
			in:       []byte(`{"foo":123}`),
			expected: "foo=123\n",
		},
		{
			name:     "primitive boolean not array",
			in:       []byte(`{"foo":true}`),
			expected: "foo=true\n",
		},
		{
			name:     "primitive null not array",
			in:       []byte(`{"foo":null}`),
			expected: "foo=null\n",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := new(bytes.Buffer)
			if err := decode(tc.in, out); err != nil {
				t.Errorf("decode() error = %v", err)
				return
			}
			got := out.String()
			if got != tc.expected {
				t.Errorf("expected = %q, but got = %q", tc.expected, got)
			}
		})
	}
}
