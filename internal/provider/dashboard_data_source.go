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

	widgets := []map[string]interface{}{}
	for _, elem := range state.Widgets.Elements() {
		var intermediateStr string
		if err := json.Unmarshal([]byte(elem.String()), &intermediateStr); err != nil {
			resp.Diagnostics.AddError("failed to unmarshal widget json", err.Error())
			return
		}

		widget := map[string]interface{}{}
		if err := json.Unmarshal([]byte(intermediateStr), &widget); err != nil {
			resp.Diagnostics.AddError("failed to unmarshal widget json", err.Error())
			return
		}

		widgets = append(widgets, widget)
	}

	tflog.Info(ctx, "widgets!!!!", map[string]interface{}{
		"widgets": widgets,
	})

	dashboardJson, err := buildDashboardBodyJson(ctx, state, widgets)
	if err != nil {
		resp.Diagnostics.AddError("failed to build dashboard json", err.Error())
		return
	}

	tflog.Info(ctx, "dashboard json", map[string]interface{}{
		"dashboard_json": dashboardJson,
	})

	state.Json = types.StringValue(dashboardJson)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
