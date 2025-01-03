package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ provider.Provider = &cwDashboardProvider{}
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &cwDashboardProvider{
			version: version,
		}
	}
}

type cwDashboardProvider struct {
	version string
}

func (p *cwDashboardProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "dashboard"
	resp.Version = p.version
}

func (p *cwDashboardProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
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

func (p *cwDashboardProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

func (p *cwDashboardProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDashboardDataSource(),
	}
}

func (p *cwDashboardProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}