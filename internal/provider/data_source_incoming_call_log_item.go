package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "math/big"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &IncomingCallLogItemDataSource{}

func NewIncomingCallLogItemDataSource() datasource.DataSource {
    return &IncomingCallLogItemDataSource{}
}

// IncomingCallLogItemDataSource defines the data source implementation.
type IncomingCallLogItemDataSource struct {
    client *Client
}

// IncomingCallLogItemDataSourceModel describes the data source data model.
type IncomingCallLogItemDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
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

func (d *IncomingCallLogItemDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incoming_call_log_item"
}

func (d *IncomingCallLogItemDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Child log for each escalation attempt / user ring within a call. Look up an existing incoming call log item by `id`, or by any of its other arguments (`call_cost_in_usd_cents`, `dial_duration_in_seconds`, `incoming_call_log_id`, ...): each one set must match, and exactly one incoming call log item may match them all.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Leave unset to look up by the other arguments instead.",
                Optional: true,
                Computed: true,
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was created.",
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was updated.",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "incoming_call_log_id": schema.StringAttribute{
                MarkdownDescription: "ID of the parent Incoming Call Log. The ID of a `oneuptime_incoming_call_log` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "incoming_call_policy_escalation_rule_id": schema.StringAttribute{
                MarkdownDescription: "ID of the escalation rule used. The ID of a `oneuptime_incoming_call_policy_escalation_rule`.",
                Optional: true,
                Computed: true,
            },
            "user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who was called. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "user_phone_number": schema.StringAttribute{
                MarkdownDescription: "Phone number that was dialed.",
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Status of this dial attempt.",
                Optional: true,
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Additional status information.",
                Optional: true,
                Computed: true,
            },
            "dial_duration_in_seconds": schema.NumberAttribute{
                MarkdownDescription: "How long this dial lasted in seconds.",
                Optional: true,
                Computed: true,
            },
            "call_cost_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Cost for this dial attempt in USD cents.",
                Optional: true,
                Computed: true,
            },
            "started_at": schema.StringAttribute{
                MarkdownDescription: "When dial started.",
                Computed: true,
            },
            "ended_at": schema.StringAttribute{
                MarkdownDescription: "When dial ended.",
                Computed: true,
            },
            "is_answered": schema.BoolAttribute{
                MarkdownDescription: "Whether this user answered the call.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *IncomingCallLogItemDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IncomingCallLogItemDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IncomingCallLogItemDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.IncomingCallLogId.IsNull() && !data.IncomingCallLogId.IsUnknown() {
        filters["incomingCallLogId"] = data.IncomingCallLogId.ValueString()
        filterNames = append(filterNames, "incoming_call_log_id = "+fmt.Sprintf("%q", data.IncomingCallLogId.ValueString()))
    }
    if !data.IncomingCallPolicyEscalationRuleId.IsNull() && !data.IncomingCallPolicyEscalationRuleId.IsUnknown() {
        filters["incomingCallPolicyEscalationRuleId"] = data.IncomingCallPolicyEscalationRuleId.ValueString()
        filterNames = append(filterNames, "incoming_call_policy_escalation_rule_id = "+fmt.Sprintf("%q", data.IncomingCallPolicyEscalationRuleId.ValueString()))
    }
    if !data.UserId.IsNull() && !data.UserId.IsUnknown() {
        filters["userId"] = data.UserId.ValueString()
        filterNames = append(filterNames, "user_id = "+fmt.Sprintf("%q", data.UserId.ValueString()))
    }
    if !data.Status.IsNull() && !data.Status.IsUnknown() {
        filters["status"] = data.Status.ValueString()
        filterNames = append(filterNames, "status = "+fmt.Sprintf("%q", data.Status.ValueString()))
    }
    if !data.StatusMessage.IsNull() && !data.StatusMessage.IsUnknown() {
        filters["statusMessage"] = data.StatusMessage.ValueString()
        filterNames = append(filterNames, "status_message = "+fmt.Sprintf("%q", data.StatusMessage.ValueString()))
    }
    if !data.DialDurationInSeconds.IsNull() && !data.DialDurationInSeconds.IsUnknown() {
        filters["dialDurationInSeconds"] = lookupNumber(data.DialDurationInSeconds)
        filterNames = append(filterNames, "dial_duration_in_seconds = "+data.DialDurationInSeconds.ValueBigFloat().String())
    }
    if !data.CallCostInUsdCents.IsNull() && !data.CallCostInUsdCents.IsUnknown() {
        filters["callCostInUSDCents"] = lookupNumber(data.CallCostInUsdCents)
        filterNames = append(filterNames, "call_cost_in_usd_cents = "+data.CallCostInUsdCents.ValueBigFloat().String())
    }
    if !data.IsAnswered.IsNull() && !data.IsAnswered.IsUnknown() {
        filters["isAnswered"] = data.IsAnswered.ValueBool()
        filterNames = append(filterNames, "is_answered = "+fmt.Sprintf("%t", data.IsAnswered.ValueBool()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the incoming call log item up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the incoming call log item up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "incomingCallLogId": true,
        "incomingCallPolicyEscalationRuleId": true,
        "userId": true,
        "userPhoneNumber": true,
        "status": true,
        "statusMessage": true,
        "dialDurationInSeconds": true,
        "callCostInUSDCents": true,
        "startedAt": true,
        "endedAt": true,
        "isAnswered": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/incoming-call-log-item/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incoming_call_log_item, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No incoming call log item found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read incoming_call_log_item: %s", err))
            return
        }
        if wrapper, ok := itemResponse["data"].(map[string]interface{}); ok {
            item = wrapper
        } else {
            item = itemResponse
        }
    }
    if !hasId {
        listBody := map[string]interface{}{
            "query":  filters,
            "select": selectParam,
            // limit 2 is enough to detect ambiguity without paging.
            "limit": 2,
        }
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/incoming-call-log-item/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list incoming_call_log_item, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list incoming_call_log_item: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No incoming call log item matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one incoming call log item matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for incoming_call_log_item.")
            return
        }
        item = first
    }

    // Update the model with response data
    if obj, ok := item["_id"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Id = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Id = types.StringValue(string(jsonBytes))
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := item["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if obj, ok := item["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedAt = types.StringNull()
        }
    } else if val, ok := item["createdAt"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    } else {
        data.CreatedAt = types.StringNull()
    }
    if obj, ok := item["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := item["updatedAt"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if obj, ok := item["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProjectId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProjectId = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := item["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := item["incomingCallLogId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncomingCallLogId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncomingCallLogId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncomingCallLogId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncomingCallLogId = types.StringValue(string(jsonBytes))
        } else {
            data.IncomingCallLogId = types.StringNull()
        }
    } else if val, ok := item["incomingCallLogId"].(string); ok {
        data.IncomingCallLogId = types.StringValue(val)
    } else {
        data.IncomingCallLogId = types.StringNull()
    }
    if obj, ok := item["incomingCallPolicyEscalationRuleId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncomingCallPolicyEscalationRuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncomingCallPolicyEscalationRuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncomingCallPolicyEscalationRuleId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncomingCallPolicyEscalationRuleId = types.StringValue(string(jsonBytes))
        } else {
            data.IncomingCallPolicyEscalationRuleId = types.StringNull()
        }
    } else if val, ok := item["incomingCallPolicyEscalationRuleId"].(string); ok {
        data.IncomingCallPolicyEscalationRuleId = types.StringValue(val)
    } else {
        data.IncomingCallPolicyEscalationRuleId = types.StringNull()
    }
    if obj, ok := item["userId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UserId = types.StringValue(string(jsonBytes))
        } else {
            data.UserId = types.StringNull()
        }
    } else if val, ok := item["userId"].(string); ok {
        data.UserId = types.StringValue(val)
    } else {
        data.UserId = types.StringNull()
    }
    if obj, ok := item["userPhoneNumber"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UserPhoneNumber = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UserPhoneNumber = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UserPhoneNumber = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UserPhoneNumber = types.StringValue(string(jsonBytes))
        } else {
            data.UserPhoneNumber = types.StringNull()
        }
    } else if val, ok := item["userPhoneNumber"].(string); ok {
        data.UserPhoneNumber = types.StringValue(val)
    } else {
        data.UserPhoneNumber = types.StringNull()
    }
    if obj, ok := item["status"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Status = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Status = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Status = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Status = types.StringValue(string(jsonBytes))
        } else {
            data.Status = types.StringNull()
        }
    } else if val, ok := item["status"].(string); ok {
        data.Status = types.StringValue(val)
    } else {
        data.Status = types.StringNull()
    }
    if obj, ok := item["statusMessage"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StatusMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StatusMessage = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StatusMessage = types.StringValue(string(jsonBytes))
        } else {
            data.StatusMessage = types.StringNull()
        }
    } else if val, ok := item["statusMessage"].(string); ok {
        data.StatusMessage = types.StringValue(val)
    } else {
        data.StatusMessage = types.StringNull()
    }
    if val, ok := item["dialDurationInSeconds"].(float64); ok {
        data.DialDurationInSeconds = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["dialDurationInSeconds"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.DialDurationInSeconds = types.NumberValue(big.NewFloat(val))
        } else {
            data.DialDurationInSeconds = types.NumberNull()
        }
    } else {
        data.DialDurationInSeconds = types.NumberNull()
    }
    if val, ok := item["callCostInUSDCents"].(float64); ok {
        data.CallCostInUsdCents = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["callCostInUSDCents"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.CallCostInUsdCents = types.NumberValue(big.NewFloat(val))
        } else {
            data.CallCostInUsdCents = types.NumberNull()
        }
    } else {
        data.CallCostInUsdCents = types.NumberNull()
    }
    if obj, ok := item["startedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StartedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StartedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StartedAt = types.StringValue(string(jsonBytes))
        } else {
            data.StartedAt = types.StringNull()
        }
    } else if val, ok := item["startedAt"].(string); ok {
        data.StartedAt = types.StringValue(val)
    } else {
        data.StartedAt = types.StringNull()
    }
    if obj, ok := item["endedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EndedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EndedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EndedAt = types.StringValue(string(jsonBytes))
        } else {
            data.EndedAt = types.StringNull()
        }
    } else if val, ok := item["endedAt"].(string); ok {
        data.EndedAt = types.StringValue(val)
    } else {
        data.EndedAt = types.StringNull()
    }
    if val, ok := item["isAnswered"].(bool); ok {
        data.IsAnswered = types.BoolValue(val)
    } else {
        data.IsAnswered = types.BoolNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
