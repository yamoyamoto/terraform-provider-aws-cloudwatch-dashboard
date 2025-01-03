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
	_ datasource.DataSource = &textWidgetDataSource{}
)

type textWidgetDataSource struct {
}

func NewTextWidgetDataSource() func() datasource.DataSource {
	return func() datasource.DataSource {
		return &textWidgetDataSource{}
	}
}

func (d *textWidgetDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_text_widget"
}

func (d *textWidgetDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"title": schema.StringAttribute{
				Description: "The name of the dashboard",
				Required:    true,
			},
			"x": schema.Int32Attribute{
				Description: "The x position of the widget",
				Optional:    true,
			},
			"y": schema.Int32Attribute{
				Description: "The y position of the widget",
				Optional:    true,
			},
			"markdown": schema.StringAttribute{
				Description: "The markdown content of the widget",
				Required:    true,
			},
			"background": schema.StringAttribute{
				Description: "The background color of the widget",
				Optional:    true,
			},
			"width": schema.Int32Attribute{
				Description: "The width of the widget",
				Required:    true,
			},
			"height": schema.Int32Attribute{
				Description: "The height of the widget",
				Required:    true,
			},

			"json": schema.StringAttribute{
				Description: "The settings of the widget",
				Computed:    true,
			},
		},
	}
}

type textWidgetDataSourceModel struct {
	Title      types.String `tfsdk:"title"`
	X          types.Int32  `tfsdk:"x"`
	Y          types.Int32  `tfsdk:"y"`
	Markdown   types.String `tfsdk:"markdown"`
	Background types.String `tfsdk:"background"`
	Width      types.Int32  `tfsdk:"width"`
	Height     types.Int32  `tfsdk:"height"`

	Json types.String `tfsdk:"json"`
}

type textWidgetDataSourceSettings struct {
	Type       string `json:"type"`
	Title      string `json:"title"`
	X          int32  `json:"x"`
	Y          int32  `json:"y"`
	Markdown   string `json:"markdown"`
	Background string `json:"background"`
	Width      int32  `json:"width"`
	Height     int32  `json:"height"`
}

const (
	typeTextWidget = "text"
)

func (d *textWidgetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state textWidgetDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	settings := textWidgetDataSourceSettings{
		Type:       typeTextWidget,
		Title:      state.Title.ValueString(),
		X:          state.X.ValueInt32(),
		Y:          state.Y.ValueInt32(),
		Markdown:   state.Markdown.ValueString(),
		Background: state.Background.ValueString(),
		Width:      state.Width.ValueInt32(),
		Height:     state.Height.ValueInt32(),
	}

	b, err := json.Marshal(settings)
	if err != nil {
		resp.Diagnostics.AddError("failed to marshal widget settings", err.Error())
		return
	}

	tflog.Info(ctx, "text widget settings", map[string]interface{}{
		"settings": string(b),
	})

	state.Json = types.StringValue(string(b))

	stateDiags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(stateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func parseCWDashboardBodyWidgetText(ctx context.Context, widget map[string]interface{}) (CWDashboardBodyWidget, error) {
	cwWidget := CWDashboardBodyWidget{
		Type: "text",
	}

	title, ok := widget["title"].(string)
	if ok {
		cwWidget.Properties = title
	}

	x, ok := widget["x"].(float64)

	if ok {
		cwWidget.X = int32(x)
	}

	y, ok := widget["y"].(float64)
	if ok {
		cwWidget.Y = int32(y)
	}

	width, ok := widget["width"].(float64)
	if ok {
		cwWidget.Width = int32(width)
	}

	height, ok := widget["height"].(float64)
	if ok {
		cwWidget.Height = int32(height)
	}

	property := CWDashboardBodyWidgetPropertyText{}
	mkDown, ok := widget["markdown"].(string)
	if ok {
		property.Markdown = mkDown
	}

	background, ok := widget["background"].(string)
	if ok {
		property.Background = background
	}

	cwWidget.Properties = property

	tflog.Debug(ctx, "built text widget", map[string]interface{}{
		"widget": cwWidget,
	})

	return cwWidget, nil
}
