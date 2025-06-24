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
var _ datasource.DataSource = &SmsLogDataDataSource{}

func NewSmsLogDataDataSource() datasource.DataSource {
    return &SmsLogDataDataSource{}
}

// SmsLogDataDataSource defines the data source implementation.
type SmsLogDataDataSource struct {
    client *Client
}

// SmsLogDataDataSourceModel describes the data source data model.
type SmsLogDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    ToNumber types.String `tfsdk:"to_number"`
    FromNumber types.String `tfsdk:"from_number"`
    SmsText types.String `tfsdk:"sms_text"`
    StatusMessage types.String `tfsdk:"status_message"`
    Status types.String `tfsdk:"status"`
    SmsCostInUsdCents types.Number `tfsdk:"sms_cost_in_usd_cents"`
}

func (d *SmsLogDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_sms_log_data"
}

func (d *SmsLogDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "sms_log_data data source",

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
            "to_number": schema.StringAttribute{
                MarkdownDescription: "Phone object",
                Computed: true,
            },
            "from_number": schema.StringAttribute{
                MarkdownDescription: "Phone object",
                Computed: true,
            },
            "sms_text": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read SMS Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read SMS Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read SMS Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "sms_cost_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read SMS Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
        },
    }
}

func (d *SmsLogDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SmsLogDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data SmsLogDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "sms-log" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read sms_log_data, got error: %s", err))
        return
    }

    var smsLogDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &smsLogDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse sms_log_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := smsLogDataResponse["data"].(map[string]interface{}); ok {
        smsLogDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := smsLogDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := smsLogDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := smsLogDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := smsLogDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := smsLogDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := smsLogDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := smsLogDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := smsLogDataResponse["to_number"].(string); ok {
        data.ToNumber = types.StringValue(val)
    }
    if val, ok := smsLogDataResponse["from_number"].(string); ok {
        data.FromNumber = types.StringValue(val)
    }
    if val, ok := smsLogDataResponse["sms_text"].(string); ok {
        data.SmsText = types.StringValue(val)
    }
    if val, ok := smsLogDataResponse["status_message"].(string); ok {
        data.StatusMessage = types.StringValue(val)
    }
    if val, ok := smsLogDataResponse["status"].(string); ok {
        data.Status = types.StringValue(val)
    }
    if val, ok := smsLogDataResponse["sms_cost_in_usd_cents"].(float64); ok {
        data.SmsCostInUsdCents = types.NumberValue(big.NewFloat(val))
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
