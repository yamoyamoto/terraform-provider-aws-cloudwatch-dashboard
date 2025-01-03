package provider

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
			"content": schema.StringAttribute{
				Description: "The content of the text widget",
				Required:    true,
			},
		},
	}
}

type textWidgetDataSourceModel struct {
	Content types.String `tfsdk:"content"`
}

func (d *textWidgetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state textWidgetDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	widget := map[string]interface{}{
		"type": "text",
		"properties": map[string]interface{}{
			"content": state.Content.Value,
		},
	}

	jsonData, err := json.Marshal(widget)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting text widget to JSON",
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
