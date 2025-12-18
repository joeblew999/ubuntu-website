package youtube

import "testing"

func TestBuildFormatString(t *testing.T) {
	tests := []struct {
		name      string
		format    string
		quality   string
		audioOnly bool
		want      string
	}{
		{name: "explicit format wins", format: "bestvideo", quality: "720", want: "bestvideo"},
		{name: "audio only", audioOnly: true, want: "bestaudio/best"},
		{name: "quality limit", quality: "1080", want: "bestvideo[height<=?1080]+bestaudio/best[height<=?1080]/best"},
		{name: "default", want: "bestvideo+bestaudio/best"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildFormatString(tt.format, tt.quality, tt.audioOnly)
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestHumanBytes(t *testing.T) {
	tests := []struct {
		value int64
		want  string
	}{
		{value: 0, want: "0B"},
		{value: 1023, want: "1023B"},
		{value: 1024, want: "1.0KB"},
		{value: 1048576, want: "1.0MB"},
		{value: 1073741824, want: "1.0GB"},
	}

	for _, tt := range tests {
		got := humanBytes(tt.value)
		if got != tt.want {
			t.Fatalf("humanBytes(%d) = %q, want %q", tt.value, got, tt.want)
		}
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "Simple Title.mp4", want: "simple-title.mp4"},
		{name: "Gaussian Splatting reconstructions with Android & iPhone.mp4", want: "gaussian-splatting-reconstructions-with-android-iphone.mp4"},
		{name: "Hello   World!!!.webm", want: "hello-world.webm"},
		{name: "Test_Video_2024.MP4", want: "test-video-2024.mp4"},
		{name: "---Leading--Trailing---.mp4", want: "leading-trailing.mp4"},
		{name: "日本語タイトル.mp4", want: "video.mp4"}, // Non-ASCII falls back to "video"
		{name: "Mix of 日本語 and English.mp4", want: "mix-of-and-english.mp4"},
		{name: "file.with" + ".many.dots.mp4", want: "filewithmanydots.mp4"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeFilename(tt.name)
			if got != tt.want {
				t.Fatalf("sanitizeFilename(%q) = %q, want %q", tt.name, got, tt.want)
			}
		})
	}
}
