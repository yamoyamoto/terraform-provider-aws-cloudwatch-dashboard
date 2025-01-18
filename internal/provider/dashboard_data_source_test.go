package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestDashboardDataSourceModel_Validate(t *testing.T) {
	tests := []struct {
		name    string
		model   dashboardDataSourceModel
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid model with all fields",
			model: dashboardDataSourceModel{
				Start:          types.StringValue("2024-01-01T00:00:00Z"),
				End:            types.StringValue("2024-01-02T00:00:00Z"),
				PeriodOverride: types.StringValue("auto"),
				Widgets:        types.ListValueMust(types.StringType, []attr.Value{}),
				Json:           types.StringValue("{}"),
			},
			wantErr: false,
		},
		{
			name: "valid model with minutes relative time for start",
			model: dashboardDataSourceModel{
				Start:          types.StringValue("-PT15M"),
				End:            types.StringValue("2024-01-02T00:00:00Z"),
				PeriodOverride: types.StringValue("auto"),
				Widgets:        types.ListValueMust(types.StringType, []attr.Value{}),
			},
			wantErr: false,
		},
		{
			name: "valid model with hours relative time for start",
			model: dashboardDataSourceModel{
				Start:          types.StringValue("-PT2H"),
				End:            types.StringValue("2024-01-02T00:00:00Z"),
				PeriodOverride: types.StringValue("auto"),
			},
			wantErr: false,
		},
		{
			name: "valid model with days relative time for start",
			model: dashboardDataSourceModel{
				Start:          types.StringValue("-P7D"),
				End:            types.StringValue("2024-01-02T00:00:00Z"),
				PeriodOverride: types.StringValue("auto"),
			},
			wantErr: false,
		},
		{
			name: "valid model with weeks relative time for start",
			model: dashboardDataSourceModel{
				Start:          types.StringValue("-P2W"),
				End:            types.StringValue("2024-01-02T00:00:00Z"),
				PeriodOverride: types.StringValue("auto"),
			},
			wantErr: false,
		},
		{
			name: "valid model with months relative time for start",
			model: dashboardDataSourceModel{
				Start:          types.StringValue("-P3M"),
				End:            types.StringValue("2024-01-02T00:00:00Z"),
				PeriodOverride: types.StringValue("auto"),
			},
			wantErr: false,
		},
		{
			name: "invalid relative time format for start (invalid minutes)",
			model: dashboardDataSourceModel{
				Start: types.StringValue("-PT15X"),
			},
			wantErr: true,
			errMsg:  "start must be a valid ISO8601 date or a valid relative time",
		},
		{
			name: "invalid relative time format for start (invalid days)",
			model: dashboardDataSourceModel{
				Start: types.StringValue("-P7X"),
			},
			wantErr: true,
			errMsg:  "start must be a valid ISO8601 date or a valid relative time",
		},
		{
			name: "invalid start date format",
			model: dashboardDataSourceModel{
				Start: types.StringValue("invalid-date"),
			},
			wantErr: true,
			errMsg:  "start must be a valid ISO8601 date or a valid relative time",
		},
		{
			name: "invalid end date format",
			model: dashboardDataSourceModel{
				End: types.StringValue("invalid-date"),
			},
			wantErr: true,
			errMsg:  "end must be a valid ISO8601 date",
		},
		{
			name: "invalid period override value",
			model: dashboardDataSourceModel{
				PeriodOverride: types.StringValue("invalid"),
			},
			wantErr: true,
			errMsg:  "period_override must be either 'auto' or 'inherit'",
		},
		{
			name: "too many widgets",
			model: dashboardDataSourceModel{
				Widgets: createWidgetsList(501),
			},
			wantErr: true,
			errMsg:  "maximum number of widgets is 500",
		},
		{
			name: "valid model with inherit period override",
			model: dashboardDataSourceModel{
				PeriodOverride: types.StringValue("inherit"),
			},
			wantErr: false,
		},
		{
			name:    "empty model",
			model:   dashboardDataSourceModel{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.model.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Helper function to create a list with specified number of widgets
func createWidgetsList(count int) types.List {
	elements := make([]attr.Value, count)
	for i := 0; i < count; i++ {
		elements[i] = types.StringValue("")
	}
	return types.ListValueMust(types.StringType, elements)
}
