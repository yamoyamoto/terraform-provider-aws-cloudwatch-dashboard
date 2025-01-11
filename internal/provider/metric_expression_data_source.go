package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

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
				Optional:    true,
			},
			"label": schema.StringAttribute{
				Description: "The label of the metric",
				Optional:    true,
			},
			"period": schema.Int32Attribute{
				Description: "The period of the metric",
				Optional:    true,
			},
			"using_metrics": schema.MapAttribute{
				Description: "The metrics used in the expression",
				Optional:    true,
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
	Expression   types.String `tfsdk:"expression"`
	Color        types.String `tfsdk:"color"`
	Label        types.String `tfsdk:"label"`
	Period       types.Int32  `tfsdk:"period"`
	UsingMetrics types.Map    `tfsdk:"using_metrics"`
	Json         types.String `tfsdk:"json"`
}

type metricExpressionDataSourceSettings struct {
	Type         string            `json:"type"`
	Expression   string            `json:"expression"`
	Color        string            `json:"color"`
	Label        string            `json:"label"`
	Period       int32             `json:"period"`
	UsingMetrics map[string]string `json:"using_metrics"`
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

	usingMetrics := make(map[string]string)
	for k, v := range state.UsingMetrics.Elements() {
		vs, ok := v.(types.String)
		if !ok {
			resp.Diagnostics.AddError("failed to convert using_metrics", "failed to convert using_metrics")
			return
		}

		usingMetrics[k] = vs.ValueString()
	}

	settings := metricExpressionDataSourceSettings{
		Type:         typeNameOfMetricExpressionDataSource,
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

	// sort keys to make the order of metrics deterministic
	keys := make([]string, 0, len(s.UsingMetrics))
	for id := range s.UsingMetrics {
		keys = append(keys, id)
	}
	sort.Strings(keys)

	for _, id := range keys {
		usingMetric := s.UsingMetrics[id]
		var m metricDataSourceSettings
		if err := json.Unmarshal([]byte(usingMetric), &m); err != nil {
			return nil, err
		}

		ms, err := m.buildMetricWidgetMetricsSettings(left, map[string]interface{}{
			"id":      id,
			"visible": false,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to build metric widget metric settings: %w", err)
		}

		settings = append(settings, ms)
	}

	metricExpressionSettings := []interface{}{
		map[string]interface{}{
			"expression": s.Expression,
		},
	}

	if s.Color != "" {
		metricExpressionSettings[0].(map[string]interface{})["color"] = s.Color
	}

	if s.Label != "" {
		metricExpressionSettings[0].(map[string]interface{})["label"] = s.Label
	}

	settings = append(settings, metricExpressionSettings)

	return settings, nil
}
