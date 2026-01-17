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
var _ datasource.DataSource = &IncomingCallLogDataDataSource{}

func NewIncomingCallLogDataDataSource() datasource.DataSource {
    return &IncomingCallLogDataDataSource{}
}

// IncomingCallLogDataDataSource defines the data source implementation.
type IncomingCallLogDataDataSource struct {
    client *Client
}

// IncomingCallLogDataDataSourceModel describes the data source data model.
type IncomingCallLogDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    IncomingCallPolicyId types.String `tfsdk:"incoming_call_policy_id"`
    CallerPhoneNumber types.String `tfsdk:"caller_phone_number"`
    RoutingPhoneNumber types.String `tfsdk:"routing_phone_number"`
    CallProviderCallId types.String `tfsdk:"call_provider_call_id"`
    Status types.String `tfsdk:"status"`
    StatusMessage types.String `tfsdk:"status_message"`
    CallDurationInSeconds types.Number `tfsdk:"call_duration_in_seconds"`
    CallCostInUsdCents types.Number `tfsdk:"call_cost_in_usd_cents"`
    IncomingCallCostInUsdCents types.Number `tfsdk:"incoming_call_cost_in_usd_cents"`
    OutgoingCallCostInUsdCents types.Number `tfsdk:"outgoing_call_cost_in_usd_cents"`
    StartedAt types.String `tfsdk:"started_at"`
    EndedAt types.String `tfsdk:"ended_at"`
    AnsweredByUserId types.String `tfsdk:"answered_by_user_id"`
    CurrentEscalationRuleOrder types.Number `tfsdk:"current_escalation_rule_order"`
    RepeatCount types.Number `tfsdk:"repeat_count"`
}

func (d *IncomingCallLogDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incoming_call_log_data"
}

func (d *IncomingCallLogDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "incoming_call_log_data data source",

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
            "incoming_call_policy_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "caller_phone_number": schema.StringAttribute{
                MarkdownDescription: "Phone object",
                Computed: true,
            },
            "routing_phone_number": schema.StringAttribute{
                MarkdownDescription: "Phone object",
                Computed: true,
            },
            "call_provider_call_id": schema.StringAttribute{
                MarkdownDescription: "Call provider's call identifier. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Current status of the incoming call. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Additional status information. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "call_duration_in_seconds": schema.NumberAttribute{
                MarkdownDescription: "Total call duration in seconds. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "call_cost_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Total cost for this call in USD cents. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "incoming_call_cost_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Cost for incoming leg in USD cents. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "outgoing_call_cost_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Cost for all forwarding attempts in USD cents. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Log], Update: [No access - you don't have permission for this operation]",
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
            "answered_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "current_escalation_rule_order": schema.NumberAttribute{
                MarkdownDescription: "The current escalation rule order being processed. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "repeat_count": schema.NumberAttribute{
                MarkdownDescription: "Number of times the policy has been repeated. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Incoming Call Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
        },
    }
}

func (d *IncomingCallLogDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IncomingCallLogDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IncomingCallLogDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "incoming-call-log" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incoming_call_log_data, got error: %s", err))
        return
    }

    var incomingCallLogDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &incomingCallLogDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incoming_call_log_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := incomingCallLogDataResponse["data"].(map[string]interface{}); ok {
        incomingCallLogDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := incomingCallLogDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := incomingCallLogDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := incomingCallLogDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := incomingCallLogDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := incomingCallLogDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := incomingCallLogDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incomingCallLogDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := incomingCallLogDataResponse["incoming_call_policy_id"].(string); ok {
        data.IncomingCallPolicyId = types.StringValue(val)
    }
    if val, ok := incomingCallLogDataResponse["caller_phone_number"].(string); ok {
        data.CallerPhoneNumber = types.StringValue(val)
    }
    if val, ok := incomingCallLogDataResponse["routing_phone_number"].(string); ok {
        data.RoutingPhoneNumber = types.StringValue(val)
    }
    if val, ok := incomingCallLogDataResponse["call_provider_call_id"].(string); ok {
        data.CallProviderCallId = types.StringValue(val)
    }
    if val, ok := incomingCallLogDataResponse["status"].(string); ok {
        data.Status = types.StringValue(val)
    }
    if val, ok := incomingCallLogDataResponse["status_message"].(string); ok {
        data.StatusMessage = types.StringValue(val)
    }
    if val, ok := incomingCallLogDataResponse["call_duration_in_seconds"].(float64); ok {
        data.CallDurationInSeconds = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incomingCallLogDataResponse["call_cost_in_usd_cents"].(float64); ok {
        data.CallCostInUsdCents = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incomingCallLogDataResponse["incoming_call_cost_in_usd_cents"].(float64); ok {
        data.IncomingCallCostInUsdCents = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incomingCallLogDataResponse["outgoing_call_cost_in_usd_cents"].(float64); ok {
        data.OutgoingCallCostInUsdCents = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incomingCallLogDataResponse["started_at"].(string); ok {
        data.StartedAt = types.StringValue(val)
    }
    if val, ok := incomingCallLogDataResponse["ended_at"].(string); ok {
        data.EndedAt = types.StringValue(val)
    }
    if val, ok := incomingCallLogDataResponse["answered_by_user_id"].(string); ok {
        data.AnsweredByUserId = types.StringValue(val)
    }
    if val, ok := incomingCallLogDataResponse["current_escalation_rule_order"].(float64); ok {
        data.CurrentEscalationRuleOrder = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incomingCallLogDataResponse["repeat_count"].(float64); ok {
        data.RepeatCount = types.NumberValue(big.NewFloat(val))
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
