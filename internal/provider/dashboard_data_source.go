package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource = &dashboardDataSource{}
)

type dashboardDataSource struct {
}

func NewDashboardDataSource() func() datasource.DataSource {
	return func() datasource.DataSource {
		return &dashboardDataSource{}
	}
}

func (d *dashboardDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName
}

func (d *dashboardDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"widgets": schema.ListAttribute{
				Description: `The list of widgets in the dashboard.`,
				Required:    true,
				ElementType: types.StringType,
			},
			"end": schema.StringAttribute{
				Description: `The end of the time range to use for each widget on the dashboard when the dashboard loads. If you specify a value for end, you must also specify a value for ` + "`start`" + `. For each of these values, specify an absolute time in the ISO 8601 format. For example, ` + "`2018-12-17T06:00:00.000Z`" + `.`,
				Optional:    true,
			},
			"start": schema.StringAttribute{
				Description: `The start of the time range to use for each widget on the dashboard.

You can specify ` + "`start`" + ` without specifying end to specify a relative time range that ends with the current time. In this case, the value of ` + "`start`" + ` must begin with ` + "`-PT`" + ` if you specify a time range in minutes or hours, and must begin with ` + "`-P`" + ` if you specify a time range in days, weeks, or months. You can then use M, H, D, W and M as abbreviations for minutes, hours, days, weeks and months. For example, ` + "`-PT5M`" + ` shows the last 5 minutes, ` + "`-PT8H`" + ` shows the last 8 hours, and ` + "`-P3M`" + ` shows the last three months.

You can also use ` + "`start`" + ` along with an end field, to specify an absolute time range. When specifying an absolute time range, use the ISO 8601 format. For example, ` + "`2018-12-17T06:00:00.000Z`" + `.

If you omit ` + "`start`" + `, the dashboard shows the default time range when it loads.`,
				Optional: true,
			},
			"period_override": schema.StringAttribute{
				Description: `Use this field to specify the period for the graphs when the dashboard loads. Specifying ` + "`auto`" + ` causes the period of all graphs on the dashboard to automatically adapt to the time range of the dashboard. Specifying ` + "`inherit`" + ` ensures that the period set for each graph is always obeyed.

Valid Values: auto | inherit`,
				Required: true,
			},
			"json": schema.StringAttribute{
				Description: "The json of the dashboard body",
				Computed:    true,
			},
		},
	}
}

type dashboardDataSourceModel struct {
	Start          types.String `tfsdk:"start"`
	End            types.String `tfsdk:"end"`
	PeriodOverride types.String `tfsdk:"period_override"`
	Widgets        types.List   `tfsdk:"widgets"`
	Json           types.String `tfsdk:"json"`
}

func (d *dashboardDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state dashboardDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	widgets, err := d.parseToWidgetSettings(ctx, state.Widgets.Elements())
	if err != nil {
		resp.Diagnostics.AddError("failed to parse widgets", err.Error())
		return
	}

	dashboardJson, err := buildDashboardBodyJson(ctx, state, widgets)
	if err != nil {
		resp.Diagnostics.AddError("failed to build dashboard json", err.Error())
		return
	}

	tflog.Info(ctx, "built dashboard json", map[string]interface{}{
		"dashboard_json": dashboardJson,
	})

	state.Json = types.StringValue(dashboardJson)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *dashboardDataSource) parseToWidgetSettings(ctx context.Context, elements []attr.Value) ([]interface{}, error) {
	widgets := make([]interface{}, 0)
	var currentPosition *widgetPosition
	for _, elem := range elements {
		// NOTE: Unmarshal twice because of double escaping by Terraform
		var escaped string
		if err := json.Unmarshal([]byte(elem.String()), &escaped); err != nil {
			return nil, fmt.Errorf("failed to unmarshal widget json: %w", err)
		}

		w := map[string]interface{}{}
		if err := json.Unmarshal([]byte(escaped), &w); err != nil {
			return nil, fmt.Errorf("failed to unmarshal widget json: %w", err)
		}

		widgetType, ok := w["type"].(string)
		if !ok {
			return nil, fmt.Errorf("missing widget type")
		}

		switch widgetType {
		case "text":
			var w textWidgetDataSourceSettings
			if err := json.Unmarshal([]byte(escaped), &w); err != nil {
				return nil, fmt.Errorf("failed to unmarshal text widget json: %w", err)
			}
			widget, err := w.ToCWDashboardBodyWidget(ctx, w, currentPosition)
			if err != nil {
				return nil, fmt.Errorf("failed to parse text widget: %w", err)
			}
			currentPosition = &widgetPosition{X: widget.X, Y: widget.Y}
			widgets = append(widgets, w)
		case "graph":
			var w graphWidgetDataSourceSettings
			if err := json.Unmarshal([]byte(escaped), &w); err != nil {
				return nil, fmt.Errorf("failed to unmarshal graph widget json: %w", err)
			}

			widget, err := w.ToCWDashboardBodyWidget(ctx, currentPosition)
			if err != nil {
				return nil, fmt.Errorf("failed to parse graph widget: %w", err)
			}
			currentPosition = &widgetPosition{X: widget.X, Y: widget.Y}
			widgets = append(widgets, w)
		default:
			return nil, fmt.Errorf("unsupported widget type")
		}
	}

	return widgets, nil
}
