package provider

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func TestParseCWDashboardBodyWidgetText(t *testing.T) {
	type testCase struct {
		name                 string
		input                map[string]interface{}
		beforeWidgetPosition widgetPosition
		expected             CWDashboardBodyWidget
	}

	tests := []testCase{
		{
			name: "should successfully parse when all fields are specified",
			input: map[string]interface{}{
				"width":      float64(8),
				"height":     float64(6),
				"markdown":   "# Test Header",
				"background": "#ffffff",
			},
			beforeWidgetPosition: widgetPosition{X: 0, Y: 0},
			expected: CWDashboardBodyWidget{
				Type:   "text",
				X:      8,
				Y:      0,
				Width:  8,
				Height: 6,
				Properties: CWDashboardBodyWidgetPropertyText{
					Markdown:   "# Test Header",
					Background: "#ffffff",
				},
			},
		},
		{
			name: "should parse with default values when optional fields are omitted",
			input: map[string]interface{}{
				"width":    float64(8),
				"height":   float64(6),
				"markdown": "# Test Header",
			},
			beforeWidgetPosition: widgetPosition{X: 8, Y: 0},
			expected: CWDashboardBodyWidget{
				Type:   "text",
				X:      16,
				Y:      0,
				Width:  8,
				Height: 6,
				Properties: CWDashboardBodyWidgetPropertyText{
					Markdown: "# Test Header",
				},
			},
		},
		{
			name: "should move to next row when exceeding max width",
			input: map[string]interface{}{
				"width":    float64(8),
				"height":   float64(6),
				"markdown": "# Test Header",
			},
			beforeWidgetPosition: widgetPosition{X: 20, Y: 0},
			expected: CWDashboardBodyWidget{
				Type:   "text",
				X:      0,
				Y:      6,
				Width:  8,
				Height: 6,
				Properties: CWDashboardBodyWidgetPropertyText{
					Markdown: "# Test Header",
				},
			},
		},
		{
			name: "should ignore fields with incorrect types during parsing",
			input: map[string]interface{}{
				"width":      float64(8),
				"height":     float64(6),
				"markdown":   "# Test Header",
				"background": 123,
			},
			beforeWidgetPosition: widgetPosition{X: 0, Y: 6},
			expected: CWDashboardBodyWidget{
				Type:   "text",
				X:      8,
				Y:      6,
				Width:  8,
				Height: 6,
				Properties: CWDashboardBodyWidgetPropertyText{
					Markdown: "# Test Header",
				},
			},
		},
		{
			name: "should successfully parse while ignoring unknown fields",
			input: map[string]interface{}{
				"width":         float64(8),
				"height":        float64(6),
				"markdown":      "# Test Header",
				"background":    "#ffffff",
				"unknown_field": "value",
				"another_field": 123,
			},
			beforeWidgetPosition: widgetPosition{X: 8, Y: 6},
			expected: CWDashboardBodyWidget{
				Type:   "text",
				X:      16,
				Y:      6,
				Width:  8,
				Height: 6,
				Properties: CWDashboardBodyWidgetPropertyText{
					Markdown:   "# Test Header",
					Background: "#ffffff",
				},
			},
		},
	}

	ctx := context.Background()

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actual, err := parseCWDashboardBodyWidgetText(ctx, tc.input, tc.beforeWidgetPosition)
			require.NoError(t, err)

			assert.Equal(t, "text", actual.Type)

			assert.Equal(t, tc.expected.X, actual.X)
			assert.Equal(t, tc.expected.Y, actual.Y)
			assert.Equal(t, tc.expected.Width, actual.Width)
			assert.Equal(t, tc.expected.Height, actual.Height)

			actualProps, ok := actual.Properties.(CWDashboardBodyWidgetPropertyText)
			require.True(t, ok)
			expectedProps := tc.expected.Properties.(CWDashboardBodyWidgetPropertyText)

			assert.Equal(t, expectedProps.Markdown, actualProps.Markdown)
			assert.Equal(t, expectedProps.Background, actualProps.Background)
		})
	}
}
