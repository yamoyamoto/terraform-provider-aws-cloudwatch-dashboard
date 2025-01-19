package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

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
	Statistic      types.String                     `tfsdk:"statistic"`
	Timezone       types.String                     `tfsdk:"timezone"`
	Title          types.String                     `tfsdk:"title"`
	View           types.String                     `tfsdk:"view"`
	Width          types.Int32                      `tfsdk:"width"`
	Json           types.String                     `tfsdk:"json"`
}

func (d *graphWidgetDataSourceModel) Validate() error {
	// Period must be 60 or a multiple of 60
	if period := d.Period.ValueInt32(); period < 60 || period%60 != 0 {
		return fmt.Errorf("period must be 60 or a multiple of 60, got: %d", period)
	}

	// Validate LegendPosition
	if !d.LegendPosition.IsNull() {
		legendPos := d.LegendPosition.ValueString()
		validPositions := map[string]bool{
			"right":  true,
			"bottom": true,
			"hidden": true,
		}
		if !validPositions[legendPos] {
			return fmt.Errorf("legend_position must be one of 'right', 'bottom', or 'hidden', got: %s", legendPos)
		}
	}

	// Validate Statistic
	if !d.Statistic.IsNull() {
		stat := d.Statistic.ValueString()
		validStats := map[string]bool{
			"SampleCount": true,
			"Average":     true,
			"Sum":         true,
			"Minimum":     true,
			"Maximum":     true,
		}

		if !validStats[stat] {
			// Check if it's a percentile statistic (p??)
			if strings.HasPrefix(stat, "p") {
				percentile, err := strconv.ParseFloat(strings.TrimPrefix(stat, "p"), 64)
				if err != nil || percentile < 0 || percentile > 100 {
					return fmt.Errorf("invalid percentile statistic: %s, must be between p0 and p100", stat)
				}
			} else {
				return fmt.Errorf("statistic must be one of 'SampleCount', 'Average', 'Sum', 'Minimum', 'Maximum', or a percentile (p0-p100), got: %s", stat)
			}
		}
	}

	// Validate Timezone format if present
	if !d.Timezone.IsNull() {
		timezone := d.Timezone.ValueString()
		if timezone != "" {
			// Timezone format: [+-]HHMM
			timezonePattern := regexp.MustCompile(`^[+-][0-9]{4}$`)
			if !timezonePattern.MatchString(timezone) {
				return fmt.Errorf("invalid timezone format: %s. Must be in format +/-HHMM (e.g., +0130)", timezone)
			}

			// Validate hours (00-23) and minutes (00-59)
			hours, _ := strconv.Atoi(timezone[1:3])
			minutes, _ := strconv.Atoi(timezone[3:])
			if hours > 23 {
				return fmt.Errorf("invalid timezone hours: %02d. Must be between 00 and 23", hours)
			}
			if minutes > 59 {
				return fmt.Errorf("invalid timezone minutes: %02d. Must be between 00 and 59", minutes)
			}
		}
	}

	// Validate View
	if !d.View.IsNull() {
		view := d.View.ValueString()
		validViews := map[string]bool{
			"timeSeries":  true,
			"singleValue": true,
		}
		if !validViews[view] {
			return fmt.Errorf("view must be either 'timeSeries' or 'singleValue', got: %s", view)
		}
	}

	return nil
}

type graphWidgetYAxisDataSourceSettings struct {
	Label     string  `json:"label,omitempty"`
	Max       float64 `json:"max,omitempty"`
	Min       float64 `json:"min,omitempty"`
	ShowUnits bool    `json:"show_units,omitempty"`
}

const (
	typeGraphWidget = "graph"
)

type graphWidgetDataSourceSettings struct {
	Type           string                              `json:"type"`
	Height         int32                               `json:"height"`
	Left           []IMetricSettings                   `json:"left,omitempty"`
	LeftYAxis      *graphWidgetYAxisDataSourceSettings `json:"left_y_axis,omitempty"`
	LegendPosition string                              `json:"legend_position,omitempty"`
	LiveData       bool                                `json:"live_data,omitempty"`
	Period         int32                               `json:"period,omitempty"`
	Region         string                              `json:"region,omitempty"`
	Right          []IMetricSettings                   `json:"right,omitempty"`
	RightYAxis     *graphWidgetYAxisDataSourceSettings `json:"right_y_axis,omitempty"`
	Sparkline      bool                                `json:"sparkline,omitempty"`
	Stacked        bool                                `json:"stacked,omitempty"`
	Statistic      string                              `json:"statistic,omitempty"`
	Timezone       string                              `json:"timezone,omitempty"`
	Title          string                              `json:"title,omitempty"`
	View           string                              `json:"view,omitempty"`
	Width          int32                               `json:"width"`
}

func (s *graphWidgetDataSourceSettings) UnmarshalJSON(data []byte) error {
	var intermediate struct {
		Type           string                              `json:"type"`
		Height         int32                               `json:"height"`
		LeftYAxis      *graphWidgetYAxisDataSourceSettings `json:"left_y_axis,omitempty"`
		LegendPosition string                              `json:"legend_position,omitempty"`
		LiveData       bool                                `json:"live_data,omitempty"`
		Period         int32                               `json:"period,omitempty"`
		Region         string                              `json:"region,omitempty"`
		RightYAxis     *graphWidgetYAxisDataSourceSettings `json:"right_y_axis,omitempty"`
		Sparkline      bool                                `json:"sparkline,omitempty"`
		Stacked        bool                                `json:"stacked,omitempty"`
		Statistic      string                              `json:"statistic,omitempty"`
		Timezone       string                              `json:"timezone,omitempty"`
		Title          string                              `json:"title,omitempty"`
		View           string                              `json:"view,omitempty"`
		Width          int32                               `json:"width"`
		// Left/Right has multiple types, so we need to unmarshal them separately
		Left  []interface{} `json:"left"`
		Right []interface{} `json:"right"`
	}

	if err := json.Unmarshal(data, &intermediate); err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}

	s.Type = intermediate.Type
	s.Height = intermediate.Height
	s.LeftYAxis = intermediate.LeftYAxis
	s.LegendPosition = intermediate.LegendPosition
	s.LiveData = intermediate.LiveData
	s.Period = intermediate.Period
	s.Region = intermediate.Region
	s.RightYAxis = intermediate.RightYAxis
	s.Sparkline = intermediate.Sparkline
	s.Stacked = intermediate.Stacked
	s.Statistic = intermediate.Statistic
	s.Timezone = intermediate.Timezone
	s.Title = intermediate.Title
	s.View = intermediate.View
	s.Width = intermediate.Width

	// Process left and right metrics separately
	left, err := processMetrics(intermediate.Left)
	if err != nil {
		return err
	}
	right, err := processMetrics(intermediate.Right)
	if err != nil {
		return err
	}

	s.Left = left
	s.Right = right

	return nil
}

func processMetrics(metrics []interface{}) ([]IMetricSettings, error) {
	var result []IMetricSettings

	for _, m := range metrics {
		m2, ok := m.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid metric")
		}

		t, ok := m2["type"].(string)
		if !ok {
			return nil, fmt.Errorf("missing metric type")
		}

		switch t {
		case typeNameOfMetricDataSource:
			b, err := json.Marshal(m)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal metric settings: %w", err)
			}

			m := &metricDataSourceSettings{}
			if err := json.Unmarshal(b, m); err != nil {
				return nil, fmt.Errorf("failed to unmarshal metric settings: %w", err)
			}
			result = append(result, m)

		case typeNameOfMetricExpressionDataSource:
			b, err := json.Marshal(m)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal metric expression settings: %w", err)
			}

			m := &metricExpressionDataSourceSettings{}
			if err := json.Unmarshal(b, m); err != nil {
				return nil, fmt.Errorf("failed to unmarshal metric expression settings: %w", err)
			}
			result = append(result, m)

		default:
			return nil, fmt.Errorf("unsupported metric type")
		}
	}

	return result, nil
}

func (d *graphWidgetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state graphWidgetDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse left metrics from JSON
	leftMetrics := make([]IMetricSettings, len(state.Left))
	for i, metricJson := range state.Left {
		metric := &metricDataSourceSettings{}
		if err := json.Unmarshal([]byte(metricJson.ValueString()), &metric); err != nil {
			resp.Diagnostics.AddError("failed to unmarshal left metric", err.Error())
			return
		}

		if metric.Type != typeNameOfMetricDataSource {
			metric := &metricExpressionDataSourceSettings{}
			if err := json.Unmarshal([]byte(metricJson.ValueString()), &metric); err != nil {
				resp.Diagnostics.AddError("failed to unmarshal left metric", err.Error())
				return
			}
			leftMetrics[i] = metric
		} else {
			leftMetrics[i] = metric
		}
	}

	// Parse right metrics from JSON
	rightMetrics := make([]IMetricSettings, len(state.Right))
	for i, metricJson := range state.Right {
		var metric *metricDataSourceSettings
		if err := json.Unmarshal([]byte(metricJson.ValueString()), &metric); err != nil {
			resp.Diagnostics.AddError("failed to unmarshal right metric", err.Error())
			return
		}

		if metric.Type != typeNameOfMetricDataSource {
			metric := &metricExpressionDataSourceSettings{}
			if err := json.Unmarshal([]byte(metricJson.ValueString()), &metric); err != nil {
				resp.Diagnostics.AddError("failed to unmarshal right metric", err.Error())
				return
			}
			rightMetrics[i] = metric
		} else {
			rightMetrics[i] = metric
		}
	}

	settings := graphWidgetDataSourceSettings{
		Type:           typeGraphWidget,
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
		switch m := metric.(type) {
		case *metricDataSourceSettings:
			settings, err := m.buildMetricWidgetMetricsSettings(true, nil)
			if err != nil {
				return CWDashboardBodyWidget{}, fmt.Errorf("failed to build metric settings: %w", err)
			}
			metrics = append(metrics, settings)
		case *metricExpressionDataSourceSettings:
			settingsList, err := m.buildMetricWidgetMetricSettingsList(true)
			if err != nil {
				return CWDashboardBodyWidget{}, fmt.Errorf("failed to build metric settings: %w", err)
			}
			metrics = append(metrics, settingsList...)
		default:
			return CWDashboardBodyWidget{}, fmt.Errorf("unsupported metric type: %T", metric)

		}
	}

	for _, metric := range w.Right {
		switch m := metric.(type) {
		case *metricDataSourceSettings:
			settings, err := m.buildMetricWidgetMetricsSettings(false, nil)
			if err != nil {
				return CWDashboardBodyWidget{}, fmt.Errorf("failed to build metric settings: %w", err)
			}
			metrics = append(metrics, settings)
		case *metricExpressionDataSourceSettings:
			settingsList, err := m.buildMetricWidgetMetricSettingsList(false)
			if err != nil {
				return CWDashboardBodyWidget{}, fmt.Errorf("failed to build metric settings: %w", err)
			}
			metrics = append(metrics, settingsList...)
		default:
			return CWDashboardBodyWidget{}, fmt.Errorf("unsupported metric type: %T", metric)
		}
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
