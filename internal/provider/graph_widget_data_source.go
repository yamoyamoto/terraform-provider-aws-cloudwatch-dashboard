package provider

import (
	"context"
	"encoding/json"

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
			"end": schema.StringAttribute{
				Description: "The end of the time range to use for each widget independently from those of the dashboard",
				Optional:    true,
			},
			"height": schema.Int32Attribute{
				Description: "Height of the widget",
				Required:    true,
			},
			"left": schema.ListAttribute{
				Description: "Metrics to display on left Y axis",
				Optional:    true,
				ElementType: types.StringType,
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
			"sparkline": schema.BoolAttribute{
				Description: "Whether the graph should be shown as a sparkline",
				Optional:    true,
			},
			"stacked": schema.BoolAttribute{
				Description: "Whether the graph should be shown as stacked lines",
				Optional:    true,
			},
			"start": schema.StringAttribute{
				Description: "The start of the time range to use for each widget independently from those of the dashboard",
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

const (
	typeGraphWidget = "graph"
)

type graphWidgetDataSourceModel struct {
	End            types.String   `tfsdk:"end"`
	Height         types.Int32    `tfsdk:"height"`
	Left           []types.String `tfsdk:"left"` // JSON string containing array of metrics
	LegendPosition types.String   `tfsdk:"legend_position"`
	LiveData       types.Bool     `tfsdk:"live_data"`
	Period         types.Int32    `tfsdk:"period"`
	Region         types.String   `tfsdk:"region"`
	Right          []types.String `tfsdk:"right"` // JSON string containing array of metrics
	Sparkline      types.Bool     `tfsdk:"sparkline"`
	Stacked        types.Bool     `tfsdk:"stacked"`
	Start          types.String   `tfsdk:"start"`
	Statistic      types.String   `tfsdk:"statistic"`
	Timezone       types.String   `tfsdk:"timezone"`
	Title          types.String   `tfsdk:"title"`
	View           types.String   `tfsdk:"view"`
	Width          types.Int32    `tfsdk:"width"`
	Json           types.String   `tfsdk:"json"`
}

type graphWidgetDataSourceSettings struct {
	Type           string                     `json:"type"`
	End            string                     `json:"end,omitempty"`
	Height         int32                      `json:"height"`
	Left           []metricDataSourceSettings `json:"left,omitempty"`
	LegendPosition string                     `json:"legend_position,omitempty"`
	LiveData       bool                       `json:"live_data,omitempty"`
	Period         int32                      `json:"period,omitempty"`
	Region         string                     `json:"region,omitempty"`
	Right          []metricDataSourceSettings `json:"right,omitempty"`
	Sparkline      bool                       `json:"sparkline,omitempty"`
	Stacked        bool                       `json:"stacked,omitempty"`
	Start          string                     `json:"start,omitempty"`
	Statistic      string                     `json:"statistic,omitempty"`
	Timezone       string                     `json:"timezone,omitempty"`
	Title          string                     `json:"title,omitempty"`
	View           string                     `json:"view,omitempty"`
	Width          int32                      `json:"width"`
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
		Type:           typeGraphWidget,
		End:            state.End.ValueString(),
		Height:         state.Height.ValueInt32(),
		Left:           leftMetrics,
		LegendPosition: state.LegendPosition.ValueString(),
		LiveData:       state.LiveData.ValueBool(),
		Period:         state.Period.ValueInt32(),
		Region:         state.Region.ValueString(),
		Right:          rightMetrics,
		Sparkline:      state.Sparkline.ValueBool(),
		Stacked:        state.Stacked.ValueBool(),
		Start:          state.Start.ValueString(),
		Statistic:      state.Statistic.ValueString(),
		Timezone:       state.Timezone.ValueString(),
		Title:          state.Title.ValueString(),
		View:           state.View.ValueString(),
		Width:          state.Width.ValueInt32(),
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

func (w graphWidgetDataSourceSettings) ToCWDashboardBodyWidget(ctx context.Context, widget graphWidgetDataSourceSettings, beforeWidgetPosition *widgetPosition) (CWDashboardBodyWidget, error) {
	cwWidget := CWDashboardBodyWidget{
		Type:   "graph",
		Width:  widget.Width,
		Height: widget.Height,
		Properties: CWDashboardBodyWidgetPropertyMetric{
			AccountId:   "",
			Annotations: nil,
			LiveData:    widget.LiveData,
			Legend: &CWDashboardBodyWidgetPropertyMetricLegend{
				Position: widget.LegendPosition,
			},
			Metrics:   nil,
			Period:    widget.Period,
			Region:    widget.Region,
			Stat:      widget.Statistic,
			Title:     widget.Title,
			View:      widget.View,
			Stacked:   widget.Stacked,
			Sparkline: widget.Sparkline,
			Timezone:  widget.Timezone,
			YAxis:     nil,
			Table:     nil,
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
