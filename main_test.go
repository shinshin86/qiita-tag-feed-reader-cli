package main

import "testing"

func TestRemoveHTMLTag(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "<div>test</div>",
			want:  "test",
		},
		{
			input: "<div class='test'>test</div>",
			want:  "test",
		},
		{
			input: "<div class='test'><!-- test -->test</div>",
			want:  "test",
		},
		{
			input: "<div class='test'><a href='test' alt='test'><img src='./test' /></a></div>",
			want:  "",
		},
		{
			input: "<script src='test'></script>",
			want:  "",
		},
	}

	for _, test := range tests {
		got := removeHTMLTag(test.input)
		if got != test.want {
			t.Fatalf("want %v, but %v", test.want, got)
		}
	}
}

func TestRemoveNewlineTag(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "t\r\nest",
			want:  "test",
		},
		{
			input: "t\nest",
			want:  "test",
		},
		{
			input: "t\rest",
			want:  "test",
		},
	}

	for _, test := range tests {
		got := removeNewlineTag(test.input)
		if got != test.want {
			t.Fatalf("want %v, but %v", test.want, got)
		}
	}
}
