package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tj/assert"
)

func TestGraphWidgetDataSourceModel_Validate(t *testing.T) {
	tests := []struct {
		name    string
		model   graphWidgetDataSourceModel
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid complete model",
			model: graphWidgetDataSourceModel{
				Period:         types.Int32Value(60),
				LegendPosition: types.StringValue("right"),
				Statistic:      types.StringValue("Average"),
				Timezone:       types.StringValue("+0900"),
				View:           types.StringValue("timeSeries"),
			},
			wantErr: false,
		},
		{
			name: "valid model with minimum fields",
			model: graphWidgetDataSourceModel{
				Period: types.Int32Value(60),
			},
			wantErr: false,
		},
		{
			name: "invalid period - not allowed value",
			model: graphWidgetDataSourceModel{
				Period: types.Int32Value(45),
			},
			wantErr: true,
			errMsg:  "period must be 60 or a multiple of 60, got: 45",
		},
		{
			name: "invalid legend position",
			model: graphWidgetDataSourceModel{
				Period:         types.Int32Value(60),
				LegendPosition: types.StringValue("top"),
			},
			wantErr: true,
			errMsg:  "legend_position must be one of 'right', 'bottom', or 'hidden', got: top",
		},
		{
			name: "valid percentile statistic",
			model: graphWidgetDataSourceModel{
				Period:    types.Int32Value(60),
				Statistic: types.StringValue("p95"),
			},
			wantErr: false,
		},
		{
			name: "invalid percentile statistic",
			model: graphWidgetDataSourceModel{
				Period:    types.Int32Value(60),
				Statistic: types.StringValue("p101"),
			},
			wantErr: true,
			errMsg:  "invalid percentile statistic: p101, must be between p0 and p100",
		},
		{
			name: "invalid statistic",
			model: graphWidgetDataSourceModel{
				Period:    types.Int32Value(60),
				Statistic: types.StringValue("InvalidStat"),
			},
			wantErr: true,
			errMsg:  "statistic must be one of 'SampleCount', 'Average', 'Sum', 'Minimum', 'Maximum', or a percentile (p0-p100), got: InvalidStat",
		},
		{
			name: "invalid timezone format",
			model: graphWidgetDataSourceModel{
				Period:   types.Int32Value(60),
				Timezone: types.StringValue("0900"),
			},
			wantErr: true,
			errMsg:  "invalid timezone format: 0900. Must be in format +/-HHMM (e.g., +0130)",
		},
		{
			name: "invalid timezone hours",
			model: graphWidgetDataSourceModel{
				Period:   types.Int32Value(60),
				Timezone: types.StringValue("+2500"),
			},
			wantErr: true,
			errMsg:  "invalid timezone hours: 25. Must be between 00 and 23",
		},
		{
			name: "invalid timezone minutes",
			model: graphWidgetDataSourceModel{
				Period:   types.Int32Value(60),
				Timezone: types.StringValue("+0060"),
			},
			wantErr: true,
			errMsg:  "invalid timezone minutes: 60. Must be between 00 and 59",
		},
		{
			name: "invalid view",
			model: graphWidgetDataSourceModel{
				Period: types.Int32Value(60),
				View:   types.StringValue("invalid"),
			},
			wantErr: true,
			errMsg:  "view must be either 'timeSeries' or 'singleValue', got: invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.model.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() error = nil, want error")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("Validate() error = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}
			if err != nil {
				t.Errorf("Validate() error = %v, want nil", err)
			}
		})
	}
}

func TestGraphWidgetDatasourceSettings_ToCWDashboardBodyWidget(t *testing.T) {
	t.Run("should successfully parse when all fields are specified", func(t *testing.T) {
		input := graphWidgetDataSourceSettings{
			Height: 6,
			Left: []IMetricSettings{
				&metricDataSourceSettings{
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
			Right: []IMetricSettings{
				&metricDataSourceSettings{
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
