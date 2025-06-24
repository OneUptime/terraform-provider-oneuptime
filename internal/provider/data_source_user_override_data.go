package provider

import (
    "context"
    "fmt"
    "math/big"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &UserOverrideDataDataSource{}

func NewUserOverrideDataDataSource() datasource.DataSource {
    return &UserOverrideDataDataSource{}
}

// UserOverrideDataDataSource defines the data source implementation.
type UserOverrideDataDataSource struct {
    client *Client
}

// UserOverrideDataDataSourceModel describes the data source data model.
type UserOverrideDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    OnCallDutyPolicyId types.String `tfsdk:"on_call_duty_policy_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    OverrideUserId types.String `tfsdk:"override_user_id"`
    RouteAlertsToUserId types.String `tfsdk:"route_alerts_to_user_id"`
    StartsAt types.String `tfsdk:"starts_at"`
    EndsAt types.String `tfsdk:"ends_at"`
}

func (d *UserOverrideDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_user_override_data"
}

func (d *UserOverrideDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "user_override_data data source",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Identifier to filter by",
                Optional: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name to filter by",
                Optional: true,
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "deleted_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "version": schema.NumberAttribute{
                MarkdownDescription: "Version",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "on_call_duty_policy_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "override_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "route_alerts_to_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "starts_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "ends_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
        },
    }
}

func (d *UserOverrideDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *UserOverrideDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data UserOverrideDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "on-call-duty-policy-user-override" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read user_override_data, got error: %s", err))
        return
    }

    var userOverrideDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &userOverrideDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse user_override_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := userOverrideDataResponse["data"].(map[string]interface{}); ok {
        userOverrideDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := userOverrideDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := userOverrideDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := userOverrideDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := userOverrideDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := userOverrideDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := userOverrideDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := userOverrideDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := userOverrideDataResponse["on_call_duty_policy_id"].(string); ok {
        data.OnCallDutyPolicyId = types.StringValue(val)
    }
    if val, ok := userOverrideDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := userOverrideDataResponse["override_user_id"].(string); ok {
        data.OverrideUserId = types.StringValue(val)
    }
    if val, ok := userOverrideDataResponse["route_alerts_to_user_id"].(string); ok {
        data.RouteAlertsToUserId = types.StringValue(val)
    }
    if val, ok := userOverrideDataResponse["starts_at"].(string); ok {
        data.StartsAt = types.StringValue(val)
    }
    if val, ok := userOverrideDataResponse["ends_at"].(string); ok {
        data.EndsAt = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
