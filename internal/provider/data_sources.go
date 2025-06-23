package provider

import (
    "github.com/hashicorp/terraform-plugin-framework/datasource"
)

// GetDataSources returns all available data sources
func GetDataSources() []func() datasource.DataSource {
    return []func() datasource.DataSource{

    }
}
