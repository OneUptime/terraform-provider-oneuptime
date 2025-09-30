package provider

import (
    "context"
    "fmt"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &TotpAuthDataDataSource{}

func NewTotpAuthDataDataSource() datasource.DataSource {
    return &TotpAuthDataDataSource{}
}

// TotpAuthDataDataSource defines the data source implementation.
type TotpAuthDataDataSource struct {
    client *Client
}

// TotpAuthDataDataSourceModel describes the data source data model.
type TotpAuthDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
}

func (d *TotpAuthDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_totp_auth_data"
}

func (d *TotpAuthDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "totp_auth_data data source",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Identifier to filter by",
                Optional: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name to filter by",
                Optional: true,
            },
        },
    }
}

func (d *TotpAuthDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
    // Prevent panic if the provider has not been configured.
    if req.ProviderData == nil {
        return
    }

    client, ok := req.ProviderData.(*Client)

    if !ok {
        resp.Diagnostics.AddError(
            "Unexpected Data Source Configure Type",
            fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
        )

        return
    }

    d.client = client
}

func (d *TotpAuthDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data TotpAuthDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build request body with query parameters
    requestBody := map[string]interface{}{
        "query": map[string]interface{}{},
        "select": map[string]interface{}{},
    }
    
    // Add filters based on data source inputs
    queryFilters := map[string]interface{}{}
    if !data.Id.IsNull() {
        queryFilters["_id"] = data.Id.ValueString()
    }
    if !data.Name.IsNull() {
        queryFilters["name"] = data.Name.ValueString()
    }
    if len(queryFilters) > 0 {
        requestBody["query"] = queryFilters
    }
    
    // Make API call
    httpResp, err := d.client.Post("/user-totp-auth/get-list", requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read totp_auth_data, got error: %s", err))
        return
    }

    var totpAuthDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &totpAuthDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse totp_auth_data response, got error: %s", err))
        return
    }

    // For list operations, take the first matching item
    if items, ok := totpAuthDataResponse["data"].([]interface{}); ok && len(items) > 0 {
        if firstItem, ok := items[0].(map[string]interface{}); ok {
            totpAuthDataResponse = firstItem
        }
    }

    // Update the model with response data
    if val, ok := totpAuthDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := totpAuthDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
