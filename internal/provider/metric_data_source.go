package provider

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type metricDataSource struct {
}

func NewMetricDataSource() func() datasource.DataSource {
	return func() datasource.DataSource {
		return &metricDataSource{}
	}
}

func (d *metricDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_metric"
}

func (d *metricDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"metric_name": schema.StringAttribute{
				Description: "Name of the metric",
				Required:    true,
			},
			"namespace": schema.StringAttribute{
				Description: "Namespace of the metric",
				Required:    true,
			},
			"account": schema.StringAttribute{
				Description: "Account which this metric comes from",
				Optional:    true,
			},
			"color": schema.StringAttribute{
				Description: "The hex color code, prefixed with '#' (e.g. '#00ff00'), to use when this metric is rendered on a graph",
				Optional:    true,
			},
			"dimensions_map": schema.MapAttribute{
				Description: "Dimensions of the metric",
				Optional:    true,
				ElementType: types.StringType,
			},
			"label": schema.StringAttribute{
				Description: "Label for this metric when added to a Graph in a Dashboard",
				Optional:    true,
			},
			"period": schema.Int32Attribute{
				Description: "The period over which the specified statistic is applied",
				Optional:    true,
			},
			"region": schema.StringAttribute{
				Description: "Region which this metric comes from",
				Optional:    true,
			},
			"statistic": schema.StringAttribute{
				Description: "What function to use for aggregating",
				Optional:    true,
			},
			"unit": schema.StringAttribute{
				Description: "Unit used to filter the metric stream",
				Optional:    true,
			},

			"json": schema.StringAttribute{
				Description: "The settings of the metric",
				Computed:    true,
			},
		},
	}
}

type metricDataSourceModel struct {
	MetricName    types.String `tfsdk:"metric_name"`
	Namespace     types.String `tfsdk:"namespace"`
	Account       types.String `tfsdk:"account"`
	Color         types.String `tfsdk:"color"`
	DimensionsMap types.Map    `tfsdk:"dimensions_map"`
	Label         types.String `tfsdk:"label"`
	Period        types.Int32  `tfsdk:"period"`
	Region        types.String `tfsdk:"region"`
	Statistic     types.String `tfsdk:"statistic"`
	Unit          types.String `tfsdk:"unit"`

	Json types.String `tfsdk:"json"`
}

type metricDataSourceSettings struct {
	Type          string            `json:"type"`
	MetricName    string            `json:"metricName"`
	Namespace     string            `json:"namespace"`
	Account       string            `json:"account,omitempty"`
	Color         string            `json:"color,omitempty"`
	DimensionsMap map[string]string `json:"dimensionsMap,omitempty"`
	Label         string            `json:"label,omitempty"`
	Period        int32             `json:"period,omitempty"`
	Region        string            `json:"region,omitempty"`
	Statistic     string            `json:"statistic,omitempty"`
	Unit          string            `json:"unit,omitempty"`
}

const (
	typeNameOfMetricDataSource = "metric"
)

func (s *metricDataSourceSettings) GetType() string {
	return typeNameOfMetricDataSource
}

func (d *metricDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state metricDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dimensionsMap := make(map[string]string)
	for i, v := range state.DimensionsMap.Elements() {
		switch typedV := v.(type) {
		case types.String:
			dimensionsMap[i] = typedV.ValueString()
		default:
			resp.Diagnostics.AddError("invalid type for dimensions map", "expected string")
			return
		}
	}

	settings := metricDataSourceSettings{
		Type:          typeNameOfMetricDataSource,
		MetricName:    state.MetricName.ValueString(),
		Namespace:     state.Namespace.ValueString(),
		Account:       state.Account.ValueString(),
		Color:         state.Color.ValueString(),
		DimensionsMap: dimensionsMap,
		Label:         state.Label.ValueString(),
		Period:        state.Period.ValueInt32(),
		Region:        state.Region.ValueString(),
		Statistic:     state.Statistic.ValueString(),
		Unit:          state.Unit.ValueString(),
	}

	b, err := json.Marshal(settings)
	if err != nil {
		resp.Diagnostics.AddError("failed to marshal metric settings", err.Error())
		return
	}
	tflog.Info(ctx, "metric settings", map[string]interface{}{
		"settings": string(b),
	})

	state.Json = types.StringValue(string(b))

	resp.State.Set(ctx, state)
}

// https://docs.aws.amazon.com/AmazonCloudWatch/latest/APIReference/CloudWatch-Dashboard-Body-Structure.html#CloudWatch-Dashboard-Properties-Metrics-Array-Format
func (s *metricDataSourceSettings) buildMetricWidgetMetricsSettings(left bool, extra map[string]interface{}) ([]interface{}, error) {
	settings := make([]interface{}, 0)

	settings = append(settings, s.Namespace)
	settings = append(settings, s.MetricName)

	for dimKey, dimVal := range s.DimensionsMap {
		settings = append(settings, dimKey)
		settings = append(settings, dimVal)
	}

	renderingProperties := map[string]interface{}{}

	if s.Color != "" {
		renderingProperties["color"] = s.Color
	}
	if s.Label != "" {
		renderingProperties["label"] = s.Label
	}
	if s.Period != 0 {
		renderingProperties["period"] = s.Period
	}
	if s.Statistic != "" {
		renderingProperties["stat"] = s.Statistic
	}

	if left {
		renderingProperties["yAxis"] = "left"
	} else {
		renderingProperties["yAxis"] = "right"
	}

	for k, v := range extra {
		renderingProperties[k] = v
	}

	settings = append(settings, renderingProperties)

	return settings, nil
}
