package content

import "testing"

func Test_stripHTML(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "basic HTML tags",
			s:    "<p>Hello, <b>world</b>!</p>",
			want: "Hello, world!",
		},
		{
			name: "nested HTML tags",
			s:    "<p>Hello, <b>world <i>again</i></b>!</p>",
			want: "Hello, world again!",
		},
		{
			name: "link tag",
			s:    `<a href="http://example.com">ayo</a>`,
			want: "ayo",
		},
		{
			name: "Nested tags",
			s:    `<div><a href="#">Nested <span>n shi</span></a></div>`,
			want: "Nested n shi",
		},
		{
			name: "multiple spaces",
			s:    `<p>Multiple   spaces    here</p>`,
			want: "Multiple spaces here",
		},
		{
			name: "script tag",
			s:    `<script>alert('evil')</script><p>Content</p>`,
			want: "Content",
		},
		{
			name: "br tag",
			s: `Line 1
		<br>
		Line 2`,
			want: "Line 1 Line 2",
		},
		{
			name: "complex URL",
			s:    `<a href="http://site.com?param=1&other=2">Complex URL</a>`,
			want: "Complex URL",
		},
		{
			name: "",
			// test case from: https://grahamhelton.com
			s:    "Charm released an AI coding agent that works in your terminal called &lt;a href=\"https://github.com/charmbracelet/crush\"&gt;Crush.",
			want: "Charm released an AI coding agent that works in your terminal called Crush.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StripHTML(tt.s)
			if got != tt.want {
				t.Errorf("stripHTML() = %v, want %v", got, tt.want)
			}
		})
	}
}
