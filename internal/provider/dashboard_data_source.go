package provider

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
			"name": schema.StringAttribute{
				Description: "The name of the dashboard",
				Required:    true,
			},
			"widgets": schema.ListNestedAttribute{
				Description: "List of widgets in the dashboard",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"text": schema.StringAttribute{
							Description: "The text content of the widget",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

type dashboardDataSourceModel struct {
	Name    types.String `tfsdk:"name"`
	Widgets []TextWidget `tfsdk:"widgets"`
}

func (d *dashboardDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state dashboardDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	widgetsJSON := make([]string, len(state.Widgets))
	for i, widget := range state.Widgets {
		widgetJSON, err := widget.ToJSON()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error converting widget to JSON",
				err.Error(),
			)
			return
		}
		widgetsJSON[i] = widgetJSON
	}

	dashboardJSON, err := json.Marshal(map[string]interface{}{
		"name":    state.Name.ValueString(),
		"widgets": widgetsJSON,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating dashboard JSON",
			err.Error(),
		)
		return
	}

	stateJSON := types.StringValue(string(dashboardJSON))
	diags := resp.State.Set(ctx, &stateJSON)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
