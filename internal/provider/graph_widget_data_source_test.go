package provider

import (
	"context"
	"testing"

	"github.com/tj/assert"
)

func TestGraphWidgetDatasourceSettings_ToCWDashboardBodyWidget(t *testing.T) {
	t.Run("should successfully parse when all fields are specified", func(t *testing.T) {
		input := graphWidgetDataSourceSettings{
			Height: 6,
			Left: []metricDataSourceSettings{
				{
					MetricName: "CPUUtilization",
					Namespace:  "AWS/EC2",
					Account:    "123456789012",
					Color:      "#ff0000",
					DimensionsMap: map[string]string{
						"InstanceId": "i-1234567890abcdef0",
					},
					Label:     "CPU Utilization",
					Period:    300,
					Region:    "us-east-1",
					Statistic: "Average",
					Unit:      "Percent",
				},
			},
			LeftYAxis: &graphWidgetYAxisDataSourceSettings{
				Label:     "percentage",
				Min:       0,
				Max:       100,
				ShowUnits: true,
			},
			LegendPosition: "bottom",
			LiveData:       true,
			Period:         300,
			Region:         "us-east-1",
			Right: []metricDataSourceSettings{
				{
					MetricName: "NetworkIn",
					Namespace:  "AWS/EC2",
					Account:    "123456789012",
					Color:      "#0000ff",
					DimensionsMap: map[string]string{
						"InstanceId": "i-1234567890abcdef0",
					},
					Label:     "Network In",
					Period:    300,
					Region:    "us-east-1",
					Statistic: "Average",
					Unit:      "Bytes",
				},
			},
			RightYAxis: &graphWidgetYAxisDataSourceSettings{
				Label:     "bytes",
				Min:       0,
				ShowUnits: true,
			},
			Sparkline: true,
			Stacked:   true,
			Statistic: "Average",
			Timezone:  "+0000",
			Title:     "EC2 Instance Monitoring",
			View:      "timeSeries",
			Width:     12,
		}
		beforeWidgetPosition := &widgetPosition{X: 6, Y: 10}

		cwWidget, err := input.ToCWDashboardBodyWidget(context.TODO(), beforeWidgetPosition)

		assert.NoError(t, err)

		assert.Equal(t, "metric", cwWidget.Type)
		assert.Equal(t, int32(6), cwWidget.X)
		assert.Equal(t, int32(10), cwWidget.Y)
		assert.Equal(t, int32(12), cwWidget.Width)
		assert.Equal(t, int32(6), cwWidget.Height)

		cwWidgetProperties, ok := cwWidget.Properties.(CWDashboardBodyWidgetPropertyMetric)
		assert.True(t, ok)
		assert.Equal(t, true, cwWidgetProperties.LiveData)
		assert.Equal(t, "bottom", cwWidgetProperties.Legend.Position)
		assert.Equal(t, int32(300), cwWidgetProperties.Period)
		assert.Equal(t, "us-east-1", cwWidgetProperties.Region)
		assert.Equal(t, "Average", cwWidgetProperties.Stat)
		assert.Equal(t, "EC2 Instance Monitoring", cwWidgetProperties.Title)
		assert.Equal(t, "timeSeries", cwWidgetProperties.View)
		assert.Equal(t, true, cwWidgetProperties.Stacked)
		assert.Equal(t, true, cwWidgetProperties.Sparkline)
		assert.Equal(t, "+0000", cwWidgetProperties.Timezone)
		assert.Equal(t, &CWDashboardBodyWidgetPropertyMetricYAxisSide{
			Label:     "percentage",
			Min:       0,
			Max:       100,
			ShowUnits: true,
		}, cwWidgetProperties.YAxis.Left)
		assert.Equal(t, &CWDashboardBodyWidgetPropertyMetricYAxisSide{
			Label:     "bytes",
			Min:       0,
			ShowUnits: true,
		}, cwWidgetProperties.YAxis.Right)
		assert.Nil(t, cwWidgetProperties.Table)

		// Metrics assertions
		assert.Len(t, cwWidgetProperties.Metrics, 2)
		assert.Equal(t, []interface{}{
			"AWS/EC2",
			"CPUUtilization",
			"InstanceId",
			"i-1234567890abcdef0",
			map[string]interface{}{
				"color":  "#ff0000",
				"label":  "CPU Utilization",
				"period": int32(300),
				"stat":   "Average",
				"yAxis":  "left",
			},
		}, cwWidgetProperties.Metrics[0])
		assert.Equal(t, []interface{}{
			"AWS/EC2",
			"NetworkIn",
			"InstanceId",
			"i-1234567890abcdef0",
			map[string]interface{}{
				"color":  "#0000ff",
				"label":  "Network In",
				"period": int32(300),
				"stat":   "Average",
				"yAxis":  "right",
			},
		}, cwWidgetProperties.Metrics[1])

	})
}
