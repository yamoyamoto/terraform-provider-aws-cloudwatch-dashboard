package provider

import (
	"testing"

	"github.com/tj/assert"
)

func TestCalculatePosition(t *testing.T) {
	tests := []struct {
		name                 string
		size                 widgetSize
		beforeWidgetPosition *widgetPosition
		expected             widgetPosition
	}{
		{
			name: "should start from origin when beforeWidgetPosition is nil",
			size: widgetSize{
				Width:  8,
				Height: 6,
			},
			beforeWidgetPosition: nil,
			expected: widgetPosition{
				X: 0,
				Y: 0,
			},
		},
		{
			name: "should calculate next position in same row",
			size: widgetSize{
				Width:  8,
				Height: 6,
			},
			beforeWidgetPosition: &widgetPosition{X: 0, Y: 0},
			expected: widgetPosition{
				X: 8,
				Y: 0,
			},
		},
		{
			name: "should move to next row when exceeding max width",
			size: widgetSize{
				Width:  8,
				Height: 6,
			},
			beforeWidgetPosition: &widgetPosition{X: 20, Y: 0},
			expected: widgetPosition{
				X: 0,
				Y: 6,
			},
		},
		{
			name: "should calculate next position in current row when exactly at max width",
			size: widgetSize{
				Width:  8,
				Height: 6,
			},
			beforeWidgetPosition: &widgetPosition{X: 24, Y: 6},
			expected: widgetPosition{
				X: 0,
				Y: 12,
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actual := calculatePosition(tc.size, tc.beforeWidgetPosition)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
