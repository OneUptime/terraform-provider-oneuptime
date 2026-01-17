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
var _ datasource.DataSource = &IncomingCallLogItemDataDataSource{}

func NewIncomingCallLogItemDataDataSource() datasource.DataSource {
    return &IncomingCallLogItemDataDataSource{}
}

// IncomingCallLogItemDataDataSource defines the data source implementation.
type IncomingCallLogItemDataDataSource struct {
    client *Client
}

// IncomingCallLogItemDataDataSourceModel describes the data source data model.
type IncomingCallLogItemDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    IncomingCallLogId types.String `tfsdk:"incoming_call_log_id"`
    IncomingCallPolicyEscalationRuleId types.String `tfsdk:"incoming_call_policy_escalation_rule_id"`
    UserId types.String `tfsdk:"user_id"`
    UserPhoneNumber types.String `tfsdk:"user_phone_number"`
    Status types.String `tfsdk:"status"`
    StatusMessage types.String `tfsdk:"status_message"`
    DialDurationInSeconds types.Number `tfsdk:"dial_duration_in_seconds"`
    CallCostInUsdCents types.Number `tfsdk:"call_cost_in_usd_cents"`
    StartedAt types.String `tfsdk:"started_at"`
    EndedAt types.String `tfsdk:"ended_at"`
    IsAnswered types.Bool `tfsdk:"is_answered"`
}

func (d *IncomingCallLogItemDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incoming_call_log_item_data"
}

func (d *IncomingCallLogItemDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "incoming_call_log_item_data data source",

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
                MarkdownDescription: "Object version",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "incoming_call_log_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "incoming_call_policy_escalation_rule_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "user_phone_number": schema.StringAttribute{
                MarkdownDescription: "Phone object",
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Status of this dial attempt. Permissions - Create: [Project Owner, Project Admin], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Log Item], Update: [Project Owner, Project Admin]",
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Additional status information. Permissions - Create: [Project Owner, Project Admin], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Log Item], Update: [Project Owner, Project Admin]",
                Computed: true,
            },
            "dial_duration_in_seconds": schema.NumberAttribute{
                MarkdownDescription: "How long this dial lasted in seconds. Permissions - Create: [Project Owner, Project Admin], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Log Item], Update: [Project Owner, Project Admin]",
                Computed: true,
            },
            "call_cost_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Cost for this dial attempt in USD cents. Permissions - Create: [Project Owner, Project Admin], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Log Item], Update: [Project Owner, Project Admin]",
                Computed: true,
            },
            "started_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "ended_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "is_answered": schema.BoolAttribute{
                MarkdownDescription: "Whether this user answered the call. Permissions - Create: [Project Owner, Project Admin], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Log Item], Update: [Project Owner, Project Admin]",
                Computed: true,
            },
        },
    }
}

func (d *IncomingCallLogItemDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IncomingCallLogItemDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IncomingCallLogItemDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "incoming-call-log-item" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incoming_call_log_item_data, got error: %s", err))
        return
    }

    var incomingCallLogItemDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &incomingCallLogItemDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incoming_call_log_item_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := incomingCallLogItemDataResponse["data"].(map[string]interface{}); ok {
        incomingCallLogItemDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := incomingCallLogItemDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := incomingCallLogItemDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := incomingCallLogItemDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := incomingCallLogItemDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := incomingCallLogItemDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := incomingCallLogItemDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incomingCallLogItemDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := incomingCallLogItemDataResponse["incoming_call_log_id"].(string); ok {
        data.IncomingCallLogId = types.StringValue(val)
    }
    if val, ok := incomingCallLogItemDataResponse["incoming_call_policy_escalation_rule_id"].(string); ok {
        data.IncomingCallPolicyEscalationRuleId = types.StringValue(val)
    }
    if val, ok := incomingCallLogItemDataResponse["user_id"].(string); ok {
        data.UserId = types.StringValue(val)
    }
    if val, ok := incomingCallLogItemDataResponse["user_phone_number"].(string); ok {
        data.UserPhoneNumber = types.StringValue(val)
    }
    if val, ok := incomingCallLogItemDataResponse["status"].(string); ok {
        data.Status = types.StringValue(val)
    }
    if val, ok := incomingCallLogItemDataResponse["status_message"].(string); ok {
        data.StatusMessage = types.StringValue(val)
    }
    if val, ok := incomingCallLogItemDataResponse["dial_duration_in_seconds"].(float64); ok {
        data.DialDurationInSeconds = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incomingCallLogItemDataResponse["call_cost_in_usd_cents"].(float64); ok {
        data.CallCostInUsdCents = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incomingCallLogItemDataResponse["started_at"].(string); ok {
        data.StartedAt = types.StringValue(val)
    }
    if val, ok := incomingCallLogItemDataResponse["ended_at"].(string); ok {
        data.EndedAt = types.StringValue(val)
    }
    if val, ok := incomingCallLogItemDataResponse["is_answered"].(bool); ok {
        data.IsAnswered = types.BoolValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
