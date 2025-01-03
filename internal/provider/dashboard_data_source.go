package provider

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/yamoyamoto/terraform-provider-aws-cloudwatch-dashboard/internal/provider/widget"
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
	// TODO: add more fields
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
						"type": schema.StringAttribute{
							Description: "The type of the widget",
							Required:    true,
						},
						"properties": schema.MapAttribute{
							Description: "The properties of the widget",
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
	Widgets []widget.Widget `tfsdk:"widgets"`
}

func (d *dashboardDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state dashboardDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	widgetsJSON := []string{}
	for _, widget := range state.Widgets {
		jsonData, err := widget.ToJSON()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error converting widget to JSON",
				err.Error(),
			)
			return
		}
		widgetsJSON = append(widgetsJSON, jsonData)
	}

	dashboardJSON := map[string]interface{}{
		"name":    state.Name.Value,
		"widgets": widgetsJSON,
	}

	jsonData, err := json.Marshal(dashboardJSON)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting dashboard to JSON",
			err.Error(),
		)
		return
	}

	stateJSON := types.String{Value: string(jsonData)}

	diags := resp.State.Set(ctx, &stateJSON)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
