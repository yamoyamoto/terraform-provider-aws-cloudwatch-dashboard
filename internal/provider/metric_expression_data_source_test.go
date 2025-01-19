package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestMetricExpressionDataSourceModel_Validate(t *testing.T) {
	tests := []struct {
		name    string
		model   metricExpressionDataSourceModel
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid complete model",
			model: metricExpressionDataSourceModel{
				Period:     types.Int32Value(60),
				Color:      types.StringValue("#FF0000"),
				Expression: types.StringValue("m1 + m2"),
				UsingMetrics: createMapFromElements(map[string]string{
					"m1": "metric1",
					"m2": "metric2",
				}),
			},
			wantErr: false,
		},
		{
			name: "valid model without optional color",
			model: metricExpressionDataSourceModel{
				Period:     types.Int32Value(300),
				Expression: types.StringValue("m1 * 2"),
				UsingMetrics: createMapFromElements(map[string]string{
					"m1": "metric1",
				}),
			},
			wantErr: false,
		},
		{
			name: "invalid period - less than 60",
			model: metricExpressionDataSourceModel{
				Period:     types.Int32Value(30),
				Expression: types.StringValue("m1"),
				UsingMetrics: createMapFromElements(map[string]string{
					"m1": "metric1",
				}),
			},
			wantErr: true,
			errMsg:  "period must be 60 or a multiple of 60, got: 30",
		},
		{
			name: "invalid period - not multiple of 60",
			model: metricExpressionDataSourceModel{
				Period:     types.Int32Value(90),
				Expression: types.StringValue("m1"),
				UsingMetrics: createMapFromElements(map[string]string{
					"m1": "metric1",
				}),
			},
			wantErr: true,
			errMsg:  "period must be 60 or a multiple of 60, got: 90",
		},
		{
			name: "invalid color - missing #",
			model: metricExpressionDataSourceModel{
				Period:     types.Int32Value(60),
				Color:      types.StringValue("FF0000"),
				Expression: types.StringValue("m1"),
				UsingMetrics: createMapFromElements(map[string]string{
					"m1": "metric1",
				}),
			},
			wantErr: true,
			errMsg:  "invalid color format: FF0000, must be a six-digit hex color code (e.g., #FF0000)",
		},
		{
			name: "invalid color - wrong length",
			model: metricExpressionDataSourceModel{
				Period:     types.Int32Value(60),
				Color:      types.StringValue("#FF00"),
				Expression: types.StringValue("m1"),
				UsingMetrics: createMapFromElements(map[string]string{
					"m1": "metric1",
				}),
			},
			wantErr: true,
			errMsg:  "invalid color format: #FF00, must be a six-digit hex color code (e.g., #FF0000)",
		},
		{
			name: "empty expression",
			model: metricExpressionDataSourceModel{
				Period:       types.Int32Value(60),
				Expression:   types.StringValue(""),
				UsingMetrics: createMapFromElements(map[string]string{}),
			},
			wantErr: true,
			errMsg:  "expression cannot be empty",
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

func TestMetricExpressionDataSourceModel_ValidateExpressions(t *testing.T) {
	tests := []struct {
		name    string
		model   metricExpressionDataSourceModel
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid expression with matching metrics",
			model: metricExpressionDataSourceModel{
				Period:     types.Int32Value(60),
				Expression: types.StringValue("m1 + m2"),
				UsingMetrics: createMapFromElements(map[string]string{
					"m1": "metric1",
					"m2": "metric2",
				}),
			},
			wantErr: false,
		},
		{
			name: "valid SEARCH expression",
			model: metricExpressionDataSourceModel{
				Period:       types.Int32Value(60),
				Expression:   types.StringValue("SEARCH('{AWS/EC2,InstanceId} MetricName=\"CPUUtilization\"')"),
				UsingMetrics: createMapFromElements(map[string]string{}),
			},
			wantErr: false,
		},
		{
			name: "valid SELECT expression",
			model: metricExpressionDataSourceModel{
				Period:       types.Int32Value(60),
				Expression:   types.StringValue("SELECT SUM(CPUUtilization) FROM AWS/EC2"),
				UsingMetrics: createMapFromElements(map[string]string{}),
			},
			wantErr: false,
		},
		{
			name: "invalid variable name in using_metrics",
			model: metricExpressionDataSourceModel{
				Period:     types.Int32Value(60),
				Expression: types.StringValue("m1 + m2"),
				UsingMetrics: createMapFromElements(map[string]string{
					"M1": "metric1",
					"m2": "metric2",
				}),
			},
			wantErr: true,
			errMsg:  "invalid variable names in expression: [M1]. Must start with lowercase letter and only contain alphanumerics",
		},
		{
			name: "missing metrics in using_metrics",
			model: metricExpressionDataSourceModel{
				Period:     types.Int32Value(60),
				Expression: types.StringValue("m1 + m2 + m3"),
				UsingMetrics: createMapFromElements(map[string]string{
					"m1": "metric1",
					"m2": "metric2",
				}),
			},
			wantErr: true,
			errMsg:  "missing metrics in using_metrics: [m3]",
		},
		{
			name: "complex expression with functions",
			model: metricExpressionDataSourceModel{
				Period:     types.Int32Value(60),
				Expression: types.StringValue("AVG(m1) + SUM(m2) * 2"),
				UsingMetrics: createMapFromElements(map[string]string{
					"m1": "metric1",
					"m2": "metric2",
				}),
			},
			wantErr: false,
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

// Helper function to create types.Map from a map of strings
func createMapFromElements(elements map[string]string) types.Map {
	elemMap := make(map[string]attr.Value)
	for k, v := range elements {
		elemMap[k] = types.StringValue(v)
	}
	m, _ := types.MapValueFrom(context.Background(), types.StringType, elemMap)
	return m
}
