package provider

type WidgetDataSourceModel struct {
	Type     string                 `tfsdk:"type"`
	Settings map[string]interface{} `tfsdk:"settings"`
}
