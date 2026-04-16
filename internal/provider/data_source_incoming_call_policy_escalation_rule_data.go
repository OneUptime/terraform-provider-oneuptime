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
var _ datasource.DataSource = &IncomingCallPolicyEscalationRuleDataDataSource{}

func NewIncomingCallPolicyEscalationRuleDataDataSource() datasource.DataSource {
    return &IncomingCallPolicyEscalationRuleDataDataSource{}
}

// IncomingCallPolicyEscalationRuleDataDataSource defines the data source implementation.
type IncomingCallPolicyEscalationRuleDataDataSource struct {
    client *Client
}

// IncomingCallPolicyEscalationRuleDataDataSourceModel describes the data source data model.
type IncomingCallPolicyEscalationRuleDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    IncomingCallPolicyId types.String `tfsdk:"incoming_call_policy_id"`
    Description types.String `tfsdk:"description"`
    Order types.Number `tfsdk:"order"`
    EscalateAfterSeconds types.Number `tfsdk:"escalate_after_seconds"`
    OnCallDutyPolicyScheduleId types.String `tfsdk:"on_call_duty_policy_schedule_id"`
    UserId types.String `tfsdk:"user_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *IncomingCallPolicyEscalationRuleDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incoming_call_policy_escalation_rule_data"
}

func (d *IncomingCallPolicyEscalationRuleDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "incoming_call_policy_escalation_rule_data data source",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Identifier to filter by",
                Optional: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name to filter by",
                Computed: true,
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
            "description": schema.StringAttribute{
                MarkdownDescription: "Optional description of this escalation rule. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Manager, Create Incoming Call Policy Escalation Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Manager, Read Incoming Call Policy Escalation Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Settings Manager, Edit Incoming Call Policy Escalation Rule]",
                Computed: true,
            },
            "order": schema.NumberAttribute{
                MarkdownDescription: "Execution order (1, 2, 3...). Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Manager, Create Incoming Call Policy Escalation Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Manager, Read Incoming Call Policy Escalation Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Settings Manager, Edit Incoming Call Policy Escalation Rule]",
                Computed: true,
            },
            "escalate_after_seconds": schema.NumberAttribute{
                MarkdownDescription: "Seconds before escalating to next rule. Permissions - Create: [Project Owner, Project Admin, Project Member, Settings Manager, Create Incoming Call Policy Escalation Rule], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Manager, Read Incoming Call Policy Escalation Rule, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Settings Manager, Edit Incoming Call Policy Escalation Rule]",
                Computed: true,
            },
            "on_call_duty_policy_schedule_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (d *IncomingCallPolicyEscalationRuleDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IncomingCallPolicyEscalationRuleDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IncomingCallPolicyEscalationRuleDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "incoming-call-policy-escalation-rule" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incoming_call_policy_escalation_rule_data, got error: %s", err))
        return
    }

    var incomingCallPolicyEscalationRuleDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &incomingCallPolicyEscalationRuleDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse incoming_call_policy_escalation_rule_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := incomingCallPolicyEscalationRuleDataResponse["data"].(map[string]interface{}); ok {
        incomingCallPolicyEscalationRuleDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := incomingCallPolicyEscalationRuleDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := incomingCallPolicyEscalationRuleDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := incomingCallPolicyEscalationRuleDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := incomingCallPolicyEscalationRuleDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := incomingCallPolicyEscalationRuleDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := incomingCallPolicyEscalationRuleDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incomingCallPolicyEscalationRuleDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := incomingCallPolicyEscalationRuleDataResponse["incoming_call_policy_id"].(string); ok {
        data.IncomingCallPolicyId = types.StringValue(val)
    }
    if val, ok := incomingCallPolicyEscalationRuleDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := incomingCallPolicyEscalationRuleDataResponse["order"].(float64); ok {
        data.Order = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incomingCallPolicyEscalationRuleDataResponse["escalate_after_seconds"].(float64); ok {
        data.EscalateAfterSeconds = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := incomingCallPolicyEscalationRuleDataResponse["on_call_duty_policy_schedule_id"].(string); ok {
        data.OnCallDutyPolicyScheduleId = types.StringValue(val)
    }
    if val, ok := incomingCallPolicyEscalationRuleDataResponse["user_id"].(string); ok {
        data.UserId = types.StringValue(val)
    }
    if val, ok := incomingCallPolicyEscalationRuleDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
