package provider

import (
	"context"

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
			"text": schema.StringAttribute{
				Description: "The text content of the widget",
				Required:    true,
			},
		},
	}
}

type textWidgetDataSourceModel struct {
	Text types.String `tfsdk:"text"`
}

func (d *textWidgetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state textWidgetDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
