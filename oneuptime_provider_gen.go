// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package oneuptime

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
)

func OneuptimeProviderSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				Description:         "API key for authentication",
				MarkdownDescription: "API key for authentication",
			},
			"api_url": schema.StringAttribute{
				Optional:            true,
				Description:         "The base URL for the API",
				MarkdownDescription: "The base URL for the API",
			},
		},
	}
}

type OneuptimeModel struct {
	ApiKey types.String `tfsdk:"api_key"`
	ApiUrl types.String `tfsdk:"api_url"`
}
