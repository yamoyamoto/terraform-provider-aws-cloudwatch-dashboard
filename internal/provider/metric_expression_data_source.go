package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource = &metricExpressionDataSource{}
)

type metricExpressionDataSource struct {
}

func NewMetricExpressionDataSource() func() datasource.DataSource {
	return func() datasource.DataSource {
		return &metricExpressionDataSource{}
	}
}

func (d *metricExpressionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_metric_expression"
}

func (d *metricExpressionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"expression": schema.StringAttribute{
				Description: "The expression of the metric",
				Required:    true,
			},
			"color": schema.StringAttribute{
				Description: "The color of the metric",
				Required:    true,
			},
			"label": schema.StringAttribute{
				Description: "The label of the metric",
				Required:    true,
			},
			"period": schema.Int32Attribute{
				Description: "The period of the metric",
				Required:    true,
			},
			"using_metrics": schema.ListAttribute{
				Description: "The metrics used in the expression",
				Required:    true,
				ElementType: types.StringType,
			},
			"json": schema.StringAttribute{
				Description: "The settings of the metric",
				Computed:    true,
			},
		},
	}
}

type metricExpressionDataSourceModel struct {
	Expression   types.String   `tfsdk:"expression"`
	Color        types.String   `tfsdk:"color"`
	Label        types.String   `tfsdk:"label"`
	Period       types.Int32    `tfsdk:"period"`
	UsingMetrics []types.String `tfsdk:"using_metrics"`
	Json         types.String   `tfsdk:"json"`
}

type metricExpressionDataSourceSettings struct {
	Expression   string   `json:"expression"`
	Color        string   `json:"color"`
	Label        string   `json:"label"`
	Period       int32    `json:"period"`
	UsingMetrics []string `json:"using_metrics"`
}

const (
	typeNameOfMetricExpressionDataSource = "metric_expression"
)

func (s *metricExpressionDataSourceSettings) GetType() string {
	return typeNameOfMetricExpressionDataSource
}

func (d *metricExpressionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state metricExpressionDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: metrics validation

	usingMetrics := make([]string, 0, len(state.UsingMetrics))
	for _, m := range state.UsingMetrics {
		usingMetrics = append(usingMetrics, m.ValueString())
	}

	settings := metricExpressionDataSourceSettings{
		Expression:   state.Expression.ValueString(),
		Color:        state.Color.ValueString(),
		Label:        state.Label.ValueString(),
		Period:       state.Period.ValueInt32(),
		UsingMetrics: usingMetrics,
	}

	b, err := json.Marshal(settings)
	if err != nil {
		resp.Diagnostics.AddError("failed to marshal settings", err.Error())
		return
	}

	state.Json = types.StringValue(string(b))

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (s *metricExpressionDataSourceSettings) buildMetricWidgetMetricSettingsList(left bool) ([][]interface{}, error) {
	settings := make([][]interface{}, 0)

	for _, usingMetric := range s.UsingMetrics {
		var m metricDataSourceSettings
		if err := json.Unmarshal([]byte(usingMetric), &m); err != nil {
			return nil, err
		}

		ms, err := m.buildMetricWidgetMetricsSettings(left)
		if err != nil {
			return nil, fmt.Errorf("failed to build metric widget metric settings: %w", err)
		}

		settings = append(settings, ms)
	}

	metricExpressionSettings := []interface{}{
		map[string]interface{}{
			"expression": s.Expression,
			"label":      s.Label,
			"color":      s.Color,
		},
	}

	settings = append(settings, metricExpressionSettings)

	return settings, nil
}
