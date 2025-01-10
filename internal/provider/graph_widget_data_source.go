package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource = &graphWidgetDataSource{}
)

type graphWidgetDataSource struct {
}

func NewGraphWidgetDataSource() func() datasource.DataSource {
	return func() datasource.DataSource {
		return &graphWidgetDataSource{}
	}
}

func (d *graphWidgetDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_graph_widget"
}

// Schema defines the schema for the data source.
func (d *graphWidgetDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"height": schema.Int32Attribute{
				Description: "Height of the widget",
				Required:    true,
			},
			"left": schema.ListAttribute{
				Description: "Metrics to display on left Y axis",
				Optional:    true,
				ElementType: types.StringType,
			},
			"left_y_axis": schema.SingleNestedAttribute{
				Description: "Settings for the left Y axis",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"label": schema.StringAttribute{
						Description: "The label",
						Optional:    true,
					},
					"max": schema.Float64Attribute{
						Description: "The maximum value",
						Optional:    true,
					},
					"min": schema.Float64Attribute{
						Description: "The minimum value",
						Optional:    true,
					},
					"show_units": schema.BoolAttribute{
						Description: "Whether to show units",
						Optional:    true,
					},
				},
			},
			"legend_position": schema.StringAttribute{
				Description: "Position of the legend",
				Optional:    true,
			},
			"live_data": schema.BoolAttribute{
				Description: "Whether the graph should show live data",
				Optional:    true,
			},
			"period": schema.Int32Attribute{
				Description: "The default period for all metrics in this widget",
				Optional:    true,
			},
			"region": schema.StringAttribute{
				Description: "The region the metrics of this graph should be taken from",
				Optional:    true,
			},
			"right": schema.ListAttribute{
				Description: "Metrics to display on right Y axis",
				Optional:    true,
				ElementType: types.StringType,
			},
			"right_y_axis": schema.SingleNestedAttribute{
				Description: "Settings for the right Y axis",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"label": schema.StringAttribute{
						Description: "The label",
						Optional:    true,
					},
					"max": schema.Float64Attribute{
						Description: "The maximum value",
						Optional:    true,
					},
					"min": schema.Float64Attribute{
						Description: "The minimum value",
						Optional:    true,
					},
					"show_units": schema.BoolAttribute{
						Description: "Whether to show units",
						Optional:    true,
					},
				},
			},
			"sparkline": schema.BoolAttribute{
				Description: "Whether the graph should be shown as a sparkline",
				Optional:    true,
			},
			"stacked": schema.BoolAttribute{
				Description: "Whether the graph should be shown as stacked lines",
				Optional:    true,
			},
			"statistic": schema.StringAttribute{
				Description: "The default statistic to be displayed for each metric",
				Optional:    true,
			},
			"timezone": schema.StringAttribute{
				Description: "The timezone to use for the widget",
				Optional:    true,
			},
			"title": schema.StringAttribute{
				Description: "Title for the graph",
				Optional:    true,
			},
			"view": schema.StringAttribute{
				Description: "Display this metric",
				Optional:    true,
			},
			"width": schema.Int32Attribute{
				Description: "Width of the widget, in a grid of 24 units wide",
				Required:    true,
			},
			"json": schema.StringAttribute{
				Description: "The settings of the widget",
				Computed:    true,
			},
		},
	}
}

type graphWidgetYAxisDataSourceModel struct {
	Label     types.String  `tfsdk:"label"`
	Max       types.Float64 `tfsdk:"max"`
	Min       types.Float64 `tfsdk:"min"`
	ShowUnits types.Bool    `tfsdk:"show_units"`
}

type graphWidgetDataSourceModel struct {
	Height         types.Int32                      `tfsdk:"height"`
	Left           []types.String                   `tfsdk:"left"` // JSON string containing array of metrics
	LeftYAxis      *graphWidgetYAxisDataSourceModel `tfsdk:"left_y_axis"`
	LegendPosition types.String                     `tfsdk:"legend_position"`
	LiveData       types.Bool                       `tfsdk:"live_data"`
	Period         types.Int32                      `tfsdk:"period"`
	Region         types.String                     `tfsdk:"region"`
	Right          []types.String                   `tfsdk:"right"` // JSON string containing array of metrics
	RightYAxis     *graphWidgetYAxisDataSourceModel `tfsdk:"right_y_axis"`
	Sparkline      types.Bool                       `tfsdk:"sparkline"`
	Stacked        types.Bool                       `tfsdk:"stacked"`
	Start          types.String                     `tfsdk:"start"`
	Statistic      types.String                     `tfsdk:"statistic"`
	Timezone       types.String                     `tfsdk:"timezone"`
	Title          types.String                     `tfsdk:"title"`
	View           types.String                     `tfsdk:"view"`
	Width          types.Int32                      `tfsdk:"width"`
	Json           types.String                     `tfsdk:"json"`
}

type graphWidgetYAxisDataSourceSettings struct {
	Label     string  `json:"label,omitempty"`
	Max       float64 `json:"max,omitempty"`
	Min       float64 `json:"min,omitempty"`
	ShowUnits bool    `json:"show_units,omitempty"`
}

type graphWidgetDataSourceSettings struct {
	Height         int32                               `json:"height"`
	Left           []metricDataSourceSettings          `json:"left,omitempty"`
	LeftYAxis      *graphWidgetYAxisDataSourceSettings `json:"left_y_axis,omitempty"`
	LegendPosition string                              `json:"legend_position,omitempty"`
	LiveData       bool                                `json:"live_data,omitempty"`
	Period         int32                               `json:"period,omitempty"`
	Region         string                              `json:"region,omitempty"`
	Right          []metricDataSourceSettings          `json:"right,omitempty"`
	RightYAxis     *graphWidgetYAxisDataSourceSettings `json:"right_y_axis,omitempty"`
	Sparkline      bool                                `json:"sparkline,omitempty"`
	Stacked        bool                                `json:"stacked,omitempty"`
	Statistic      string                              `json:"statistic,omitempty"`
	Timezone       string                              `json:"timezone,omitempty"`
	Title          string                              `json:"title,omitempty"`
	View           string                              `json:"view,omitempty"`
	Width          int32                               `json:"width"`
}

func (d *graphWidgetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state graphWidgetDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse left metrics from JSON
	leftMetrics := make([]metricDataSourceSettings, len(state.Left))
	for i, metricJson := range state.Left {
		var metric metricDataSourceSettings
		if err := json.Unmarshal([]byte(metricJson.ValueString()), &metric); err != nil {
			resp.Diagnostics.AddError("failed to unmarshal left metric", err.Error())
			return
		}
		leftMetrics[i] = metric
	}

	// Parse right metrics from JSON
	rightMetrics := make([]metricDataSourceSettings, len(state.Right))
	for i, metricJson := range state.Right {
		var metric metricDataSourceSettings
		if err := json.Unmarshal([]byte(metricJson.ValueString()), &metric); err != nil {
			resp.Diagnostics.AddError("failed to unmarshal right metric", err.Error())
			return
		}
		rightMetrics[i] = metric
	}

	settings := graphWidgetDataSourceSettings{
		Height:         state.Height.ValueInt32(),
		Left:           leftMetrics,
		LegendPosition: state.LegendPosition.ValueString(),
		LiveData:       state.LiveData.ValueBool(),
		Period:         state.Period.ValueInt32(),
		Region:         state.Region.ValueString(),
		Right:          rightMetrics,
		Sparkline:      state.Sparkline.ValueBool(),
		Stacked:        state.Stacked.ValueBool(),
		Statistic:      state.Statistic.ValueString(),
		Timezone:       state.Timezone.ValueString(),
		Title:          state.Title.ValueString(),
		View:           state.View.ValueString(),
		Width:          state.Width.ValueInt32(),
	}

	if state.LeftYAxis != nil {
		settings.LeftYAxis = &graphWidgetYAxisDataSourceSettings{
			Label:     state.LeftYAxis.Label.ValueString(),
			Max:       state.LeftYAxis.Max.ValueFloat64(),
			Min:       state.LeftYAxis.Min.ValueFloat64(),
			ShowUnits: state.LeftYAxis.ShowUnits.ValueBool(),
		}
	}
	if state.RightYAxis != nil {
		settings.RightYAxis = &graphWidgetYAxisDataSourceSettings{
			Label:     state.RightYAxis.Label.ValueString(),
			Max:       state.RightYAxis.Max.ValueFloat64(),
			Min:       state.RightYAxis.Min.ValueFloat64(),
			ShowUnits: state.RightYAxis.ShowUnits.ValueBool(),
		}
	}

	b, err := json.Marshal(settings)
	if err != nil {
		resp.Diagnostics.AddError("failed to marshal widget settings", err.Error())
		return
	}

	tflog.Info(ctx, "graph widget settings", map[string]interface{}{
		"settings": string(b),
	})

	state.Json = types.StringValue(string(b))

	stateDiags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(stateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (w graphWidgetDataSourceSettings) ToCWDashboardBodyWidget(ctx context.Context, beforeWidgetPosition *widgetPosition) (CWDashboardBodyWidget, error) {
	var leftYAxis *CWDashboardBodyWidgetPropertyMetricYAxisSide
	if w.LeftYAxis != nil {
		leftYAxis = &CWDashboardBodyWidgetPropertyMetricYAxisSide{
			Label:     w.LeftYAxis.Label,
			Max:       w.LeftYAxis.Max,
			Min:       w.LeftYAxis.Min,
			ShowUnits: w.LeftYAxis.ShowUnits,
		}
	}

	var rightYAxis *CWDashboardBodyWidgetPropertyMetricYAxisSide
	if w.RightYAxis != nil {
		rightYAxis = &CWDashboardBodyWidgetPropertyMetricYAxisSide{
			Label:     w.RightYAxis.Label,
			Max:       w.RightYAxis.Max,
			Min:       w.RightYAxis.Min,
			ShowUnits: w.RightYAxis.ShowUnits,
		}
	}

	var yAxis *CWDashboardBodyWidgetPropertyMetricYAxis
	if leftYAxis != nil || rightYAxis != nil {
		yAxis = &CWDashboardBodyWidgetPropertyMetricYAxis{
			Left:  leftYAxis,
			Right: rightYAxis,
		}
	}

	metrics := make([][]interface{}, 0)
	for _, metric := range w.Left {
		settings, err := buildMetricWidgetMetricsSettings(true, metric)
		if err != nil {
			return CWDashboardBodyWidget{}, fmt.Errorf("failed to build metric settings: %w", err)
		}
		metrics = append(metrics, settings)
	}

	for _, metric := range w.Right {
		settings, err := buildMetricWidgetMetricsSettings(false, metric)
		if err != nil {
			return CWDashboardBodyWidget{}, fmt.Errorf("failed to build metric settings: %w", err)
		}
		metrics = append(metrics, settings)
	}

	cwWidget := CWDashboardBodyWidget{
		Type:   "metric",
		Width:  w.Width,
		Height: w.Height,
		Properties: CWDashboardBodyWidgetPropertyMetric{
			LiveData: w.LiveData,
			Legend: &CWDashboardBodyWidgetPropertyMetricLegend{
				Position: w.LegendPosition,
			},
			Metrics:   metrics,
			Period:    w.Period,
			Region:    w.Region,
			Stat:      w.Statistic,
			Title:     w.Title,
			View:      w.View,
			Stacked:   w.Stacked,
			Sparkline: w.Sparkline,
			Timezone:  w.Timezone,
			YAxis:     yAxis,
			// NOTE: unnecessary to set because it's not used in the graph widget
			Table: nil,
		},
	}

	position := calculatePosition(widgetSize{Width: cwWidget.Width, Height: cwWidget.Height}, beforeWidgetPosition)
	cwWidget.X = position.X
	cwWidget.Y = position.Y

	tflog.Debug(ctx, "built graph widget", map[string]interface{}{
		"widget": cwWidget,
	})

	return cwWidget, nil
}

// https://docs.aws.amazon.com/AmazonCloudWatch/latest/APIReference/CloudWatch-Dashboard-Body-Structure.html#CloudWatch-Dashboard-Properties-Metrics-Array-Format
func buildMetricWidgetMetricsSettings(left bool, metricSettings metricDataSourceSettings) ([]interface{}, error) {
	settings := make([]interface{}, 0)

	settings = append(settings, metricSettings.Namespace)
	settings = append(settings, metricSettings.MetricName)

	for dimKey, dimVal := range metricSettings.DimensionsMap {
		settings = append(settings, dimKey)
		settings = append(settings, dimVal)
	}

	renderingProperties := map[string]interface{}{}

	if metricSettings.Color != "" {
		renderingProperties["color"] = metricSettings.Color
	}
	if metricSettings.Label != "" {
		renderingProperties["label"] = metricSettings.Label
	}
	if metricSettings.Period != 0 {
		renderingProperties["period"] = metricSettings.Period
	}
	if metricSettings.Statistic != "" {
		renderingProperties["stat"] = metricSettings.Statistic
	}

	if left {
		renderingProperties["yAxis"] = "left"
	} else {
		renderingProperties["yAxis"] = "right"
	}

	settings = append(settings, renderingProperties)

	return settings, nil
}
