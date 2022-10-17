package wpscheduler

import "testing"

func TestGetTimeInterval(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			name: "TestGetTimeInterval",
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTimeInterval(); got != tt.want {
				t.Errorf("GetTimeInterval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetWallpaperPath(t *testing.T) {
	var tests = []struct {
		name string
		want string
	}{
		{
			name: "TestGetWallpaperPath",
			want: "/System/Library/Desktop Pictures/Monterey Graphic.heic",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetWallpaperPath(); got != tt.want {
				t.Errorf("GetWallpaperPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
