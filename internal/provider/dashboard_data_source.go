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
			"start": schema.StringAttribute{
				Description: "The start of the time range to use for each widget on the dashboard",
				Required:    true,
			},
			"end": schema.StringAttribute{
				Description: "The end of the time range to use for each widget on the dashboard when the dashboard loads",
				Required:    true,
			},
			"period_override": schema.StringAttribute{
				Description: "Use this field to specify the period for the graphs when the dashboard loads",
				Required:    true,
			},

			"widgets": schema.ListAttribute{
				Description: "The widgets of the dashboard",
				Required:    true,
				ElementType: types.StringType,
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
			widget, err := w.ToCWDashboardBodyWidget(ctx, w, currentPosition)
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
