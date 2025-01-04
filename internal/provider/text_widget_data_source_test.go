package provider

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func TestTextWidgetDataSourceSettingsToCWDashboardBodyWidget(t *testing.T) {
	type testCase struct {
		name                 string
		widget               textWidgetDataSourceSettings
		beforeWidgetPosition *widgetPosition
		expected             CWDashboardBodyWidget
	}

	tests := []testCase{
		{
			name: "should successfully parse when all fields are specified",
			widget: textWidgetDataSourceSettings{
				Width:      8,
				Height:     6,
				Markdown:   "# Test Header",
				Background: "#ffffff",
			},
			beforeWidgetPosition: &widgetPosition{X: 0, Y: 0},
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
			widget: textWidgetDataSourceSettings{
				Width:    8,
				Height:   6,
				Markdown: "# Test Header",
			},
			beforeWidgetPosition: &widgetPosition{X: 8, Y: 0},
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
			name: "should start from origin when beforeWidgetPosition is nil",
			widget: textWidgetDataSourceSettings{
				Width:    8,
				Height:   6,
				Markdown: "# Test Header",
			},
			beforeWidgetPosition: nil,
			expected: CWDashboardBodyWidget{
				Type:   "text",
				X:      0,
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
			widget: textWidgetDataSourceSettings{
				Width:    8,
				Height:   6,
				Markdown: "# Test Header",
			},
			beforeWidgetPosition: &widgetPosition{X: 20, Y: 0},
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
			name: "should parse with empty background",
			widget: textWidgetDataSourceSettings{
				Width:    8,
				Height:   6,
				Markdown: "# Test Header",
			},
			beforeWidgetPosition: &widgetPosition{X: 0, Y: 6},
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
	}

	ctx := context.Background()

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actual, err := tc.widget.ToCWDashboardBodyWidget(ctx, tc.widget, tc.beforeWidgetPosition)
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
