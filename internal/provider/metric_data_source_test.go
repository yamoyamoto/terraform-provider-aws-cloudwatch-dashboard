package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestMetricDataSourceModel_Validate(t *testing.T) {
	tests := []struct {
		name    string
		model   metricDataSourceModel
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid metric with all fields",
			model: metricDataSourceModel{
				Period:    types.Int32Value(60),
				Statistic: types.StringValue("Average"),
				Color:     types.StringValue("#FF0000"),
			},
			wantErr: false,
		},
		{
			name: "valid metric without optional color",
			model: metricDataSourceModel{
				Period:    types.Int32Value(300),
				Statistic: types.StringValue("Maximum"),
			},
			wantErr: false,
		},
		{
			name: "valid percentile statistic",
			model: metricDataSourceModel{
				Period:    types.Int32Value(60),
				Statistic: types.StringValue("p95"),
			},
			wantErr: false,
		},
		{
			name: "invalid period - less than 60",
			model: metricDataSourceModel{
				Period:    types.Int32Value(30),
				Statistic: types.StringValue("Average"),
			},
			wantErr: true,
			errMsg:  "period must be 60 or a multiple of 60, got: 30",
		},
		{
			name: "invalid period - not multiple of 60",
			model: metricDataSourceModel{
				Period:    types.Int32Value(90),
				Statistic: types.StringValue("Average"),
			},
			wantErr: true,
			errMsg:  "period must be 60 or a multiple of 60, got: 90",
		},
		{
			name: "invalid statistic",
			model: metricDataSourceModel{
				Period:    types.Int32Value(60),
				Statistic: types.StringValue("InvalidStat"),
			},
			wantErr: true,
			errMsg:  "invalid statistic: InvalidStat",
		},
		{
			name: "invalid percentile - negative",
			model: metricDataSourceModel{
				Period:    types.Int32Value(60),
				Statistic: types.StringValue("p-1"),
			},
			wantErr: true,
			errMsg:  "invalid percentile statistic: p-1, must be between p0 and p100",
		},
		{
			name: "invalid percentile - over 100",
			model: metricDataSourceModel{
				Period:    types.Int32Value(60),
				Statistic: types.StringValue("p101"),
			},
			wantErr: true,
			errMsg:  "invalid percentile statistic: p101, must be between p0 and p100",
		},
		{
			name: "invalid color format - missing #",
			model: metricDataSourceModel{
				Period:    types.Int32Value(60),
				Statistic: types.StringValue("Average"),
				Color:     types.StringValue("FF0000"),
			},
			wantErr: true,
			errMsg:  "invalid color format: FF0000, must be a six-digit hex color code (e.g., #FF0000)",
		},
		{
			name: "invalid color format - wrong length",
			model: metricDataSourceModel{
				Period:    types.Int32Value(60),
				Statistic: types.StringValue("Average"),
				Color:     types.StringValue("#FF00"),
			},
			wantErr: true,
			errMsg:  "invalid color format: #FF00, must be a six-digit hex color code (e.g., #FF0000)",
		},
		{
			name: "invalid color format - invalid characters",
			model: metricDataSourceModel{
				Period:    types.Int32Value(60),
				Statistic: types.StringValue("Average"),
				Color:     types.StringValue("#GG0000"),
			},
			wantErr: true,
			errMsg:  "invalid color format: #GG0000, must be a six-digit hex color code (e.g., #FF0000)",
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
