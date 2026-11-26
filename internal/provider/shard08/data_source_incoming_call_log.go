package shard08

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
    "context"
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &IncomingCallLogDataSource{}

func NewIncomingCallLogDataSource() datasource.DataSource {
    return &IncomingCallLogDataSource{}
}

// IncomingCallLogDataSource defines the data source implementation.
type IncomingCallLogDataSource struct {
    client *Client
}

// IncomingCallLogDataSourceModel describes the data source data model.
type IncomingCallLogDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    CreatedByActorType types.String `tfsdk:"created_by_actor_type"`
    CreatedByActorId types.String `tfsdk:"created_by_actor_id"`
    LastModifiedAt types.String `tfsdk:"last_modified_at"`
    LastModifiedByActorType types.String `tfsdk:"last_modified_by_actor_type"`
    LastModifiedByActorId types.String `tfsdk:"last_modified_by_actor_id"`
    LastModifiedByUserId types.String `tfsdk:"last_modified_by_user_id"`
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

func (d *IncomingCallLogDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incoming_call_log"
}

func (d *IncomingCallLogDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Parent log for each incoming call instance. Groups all escalation attempts together. Look up an existing incoming call log by `id`, or by any of its other arguments (`answered_by_user_id`, `call_cost_in_usd_cents`, `call_duration_in_seconds`, ...): each one set must match, and exactly one incoming call log may match them all.",

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
            "created_by_actor_type": schema.StringAttribute{
                MarkdownDescription: "What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.",
                Optional: true,
                Computed: true,
            },
            "created_by_actor_id": schema.StringAttribute{
                MarkdownDescription: "ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.",
                Optional: true,
                Computed: true,
            },
            "last_modified_at": schema.StringAttribute{
                MarkdownDescription: "When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.",
                Computed: true,
            },
            "last_modified_by_actor_type": schema.StringAttribute{
                MarkdownDescription: "What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.",
                Optional: true,
                Computed: true,
            },
            "last_modified_by_actor_id": schema.StringAttribute{
                MarkdownDescription: "ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.",
                Optional: true,
                Computed: true,
            },
            "last_modified_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "incoming_call_policy_id": schema.StringAttribute{
                MarkdownDescription: "ID of the Incoming Call Policy. The ID of a `oneuptime_incoming_call_policy`.",
                Optional: true,
                Computed: true,
                Validators: RecordIDValidators("oneuptime_incoming_call_policy"),
            },
            "caller_phone_number": schema.StringAttribute{
                MarkdownDescription: "Incoming caller's phone number.",
                Computed: true,
            },
            "routing_phone_number": schema.StringAttribute{
                MarkdownDescription: "The routing number that was called.",
                Computed: true,
            },
            "call_provider_call_id": schema.StringAttribute{
                MarkdownDescription: "Call provider's call identifier.",
                Optional: true,
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Current status of the incoming call.",
                Optional: true,
                Computed: true,
            },
            "status_message": schema.StringAttribute{
                MarkdownDescription: "Additional status information.",
                Optional: true,
                Computed: true,
            },
            "call_duration_in_seconds": schema.NumberAttribute{
                MarkdownDescription: "Total call duration in seconds.",
                Optional: true,
                Computed: true,
            },
            "call_cost_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Total cost for this call in USD cents.",
                Optional: true,
                Computed: true,
            },
            "incoming_call_cost_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Cost for incoming leg in USD cents.",
                Optional: true,
                Computed: true,
            },
            "outgoing_call_cost_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Cost for all forwarding attempts in USD cents.",
                Optional: true,
                Computed: true,
            },
            "started_at": schema.StringAttribute{
                MarkdownDescription: "When the call started.",
                Computed: true,
            },
            "ended_at": schema.StringAttribute{
                MarkdownDescription: "When the call ended.",
                Computed: true,
            },
            "answered_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who answered the call. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
            "current_escalation_rule_order": schema.NumberAttribute{
                MarkdownDescription: "The current escalation rule order being processed.",
                Optional: true,
                Computed: true,
            },
            "repeat_count": schema.NumberAttribute{
                MarkdownDescription: "Number of times the policy has been repeated.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *IncomingCallLogDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IncomingCallLogDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IncomingCallLogDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.CreatedByActorType.IsNull() && !data.CreatedByActorType.IsUnknown() {
        filters["createdByActorType"] = data.CreatedByActorType.ValueString()
        filterNames = append(filterNames, "created_by_actor_type = "+fmt.Sprintf("%q", data.CreatedByActorType.ValueString()))
    }
    if !data.CreatedByActorId.IsNull() && !data.CreatedByActorId.IsUnknown() {
        filters["createdByActorId"] = data.CreatedByActorId.ValueString()
        filterNames = append(filterNames, "created_by_actor_id = "+fmt.Sprintf("%q", data.CreatedByActorId.ValueString()))
    }
    if !data.LastModifiedByActorType.IsNull() && !data.LastModifiedByActorType.IsUnknown() {
        filters["lastModifiedByActorType"] = data.LastModifiedByActorType.ValueString()
        filterNames = append(filterNames, "last_modified_by_actor_type = "+fmt.Sprintf("%q", data.LastModifiedByActorType.ValueString()))
    }
    if !data.LastModifiedByActorId.IsNull() && !data.LastModifiedByActorId.IsUnknown() {
        filters["lastModifiedByActorId"] = data.LastModifiedByActorId.ValueString()
        filterNames = append(filterNames, "last_modified_by_actor_id = "+fmt.Sprintf("%q", data.LastModifiedByActorId.ValueString()))
    }
    if !data.LastModifiedByUserId.IsNull() && !data.LastModifiedByUserId.IsUnknown() {
        filters["lastModifiedByUserId"] = data.LastModifiedByUserId.ValueString()
        filterNames = append(filterNames, "last_modified_by_user_id = "+fmt.Sprintf("%q", data.LastModifiedByUserId.ValueString()))
    }
    if !data.IncomingCallPolicyId.IsNull() && !data.IncomingCallPolicyId.IsUnknown() {
        filters["incomingCallPolicyId"] = data.IncomingCallPolicyId.ValueString()
        filterNames = append(filterNames, "incoming_call_policy_id = "+fmt.Sprintf("%q", data.IncomingCallPolicyId.ValueString()))
    }
    if !data.CallProviderCallId.IsNull() && !data.CallProviderCallId.IsUnknown() {
        filters["callProviderCallId"] = data.CallProviderCallId.ValueString()
        filterNames = append(filterNames, "call_provider_call_id = "+fmt.Sprintf("%q", data.CallProviderCallId.ValueString()))
    }
    if !data.Status.IsNull() && !data.Status.IsUnknown() {
        filters["status"] = data.Status.ValueString()
        filterNames = append(filterNames, "status = "+fmt.Sprintf("%q", data.Status.ValueString()))
    }
    if !data.StatusMessage.IsNull() && !data.StatusMessage.IsUnknown() {
        filters["statusMessage"] = data.StatusMessage.ValueString()
        filterNames = append(filterNames, "status_message = "+fmt.Sprintf("%q", data.StatusMessage.ValueString()))
    }
    if !data.CallDurationInSeconds.IsNull() && !data.CallDurationInSeconds.IsUnknown() {
        filters["callDurationInSeconds"] = LookupNumber(data.CallDurationInSeconds)
        filterNames = append(filterNames, "call_duration_in_seconds = "+data.CallDurationInSeconds.ValueBigFloat().String())
    }
    if !data.CallCostInUsdCents.IsNull() && !data.CallCostInUsdCents.IsUnknown() {
        filters["callCostInUSDCents"] = LookupNumber(data.CallCostInUsdCents)
        filterNames = append(filterNames, "call_cost_in_usd_cents = "+data.CallCostInUsdCents.ValueBigFloat().String())
    }
    if !data.IncomingCallCostInUsdCents.IsNull() && !data.IncomingCallCostInUsdCents.IsUnknown() {
        filters["incomingCallCostInUSDCents"] = LookupNumber(data.IncomingCallCostInUsdCents)
        filterNames = append(filterNames, "incoming_call_cost_in_usd_cents = "+data.IncomingCallCostInUsdCents.ValueBigFloat().String())
    }
    if !data.OutgoingCallCostInUsdCents.IsNull() && !data.OutgoingCallCostInUsdCents.IsUnknown() {
        filters["outgoingCallCostInUSDCents"] = LookupNumber(data.OutgoingCallCostInUsdCents)
        filterNames = append(filterNames, "outgoing_call_cost_in_usd_cents = "+data.OutgoingCallCostInUsdCents.ValueBigFloat().String())
    }
    if !data.AnsweredByUserId.IsNull() && !data.AnsweredByUserId.IsUnknown() {
        filters["answeredByUserId"] = data.AnsweredByUserId.ValueString()
        filterNames = append(filterNames, "answered_by_user_id = "+fmt.Sprintf("%q", data.AnsweredByUserId.ValueString()))
    }
    if !data.CurrentEscalationRuleOrder.IsNull() && !data.CurrentEscalationRuleOrder.IsUnknown() {
        filters["currentEscalationRuleOrder"] = LookupNumber(data.CurrentEscalationRuleOrder)
        filterNames = append(filterNames, "current_escalation_rule_order = "+data.CurrentEscalationRuleOrder.ValueBigFloat().String())
    }
    if !data.RepeatCount.IsNull() && !data.RepeatCount.IsUnknown() {
        filters["repeatCount"] = LookupNumber(data.RepeatCount)
        filterNames = append(filterNames, "repeat_count = "+data.RepeatCount.ValueBigFloat().String())
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the incoming call log up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the incoming call log up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "createdByActorType": true,
        "createdByActorId": true,
        "lastModifiedAt": true,
        "lastModifiedByActorType": true,
        "lastModifiedByActorId": true,
        "lastModifiedByUserId": true,
        "projectId": true,
        "incomingCallPolicyId": true,
        "callerPhoneNumber": true,
        "routingPhoneNumber": true,
        "callProviderCallId": true,
        "status": true,
        "statusMessage": true,
        "callDurationInSeconds": true,
        "callCostInUSDCents": true,
        "incomingCallCostInUSDCents": true,
        "outgoingCallCostInUSDCents": true,
        "startedAt": true,
        "endedAt": true,
        "answeredByUserId": true,
        "currentEscalationRuleOrder": true,
        "repeatCount": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/incoming-call-log/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incoming_call_log, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No incoming call log found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read incoming_call_log: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/incoming-call-log/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list incoming_call_log, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list incoming_call_log: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No incoming call log matches %s.", DescribeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one incoming call log matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", DescribeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for incoming_call_log.")
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
    } else if arr, ok := item["_id"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["createdAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["updatedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := item["updatedAt"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if obj, ok := item["createdByActorType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByActorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorType = types.StringNull()
        }
    } else if arr, ok := item["createdByActorType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CreatedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorType = types.StringNull()
        }
    } else if val, ok := item["createdByActorType"].(string); ok {
        data.CreatedByActorType = types.StringValue(val)
    } else {
        data.CreatedByActorType = types.StringNull()
    }
    if obj, ok := item["createdByActorId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByActorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorId = types.StringNull()
        }
    } else if arr, ok := item["createdByActorId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CreatedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorId = types.StringNull()
        }
    } else if val, ok := item["createdByActorId"].(string); ok {
        data.CreatedByActorId = types.StringValue(val)
    } else {
        data.CreatedByActorId = types.StringNull()
    }
    if obj, ok := item["lastModifiedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedAt = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedAt = types.StringNull()
        }
    } else if val, ok := item["lastModifiedAt"].(string); ok {
        data.LastModifiedAt = types.StringValue(val)
    } else {
        data.LastModifiedAt = types.StringNull()
    }
    if obj, ok := item["lastModifiedByActorType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedByActorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorType = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedByActorType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorType = types.StringNull()
        }
    } else if val, ok := item["lastModifiedByActorType"].(string); ok {
        data.LastModifiedByActorType = types.StringValue(val)
    } else {
        data.LastModifiedByActorType = types.StringNull()
    }
    if obj, ok := item["lastModifiedByActorId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedByActorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorId = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedByActorId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorId = types.StringNull()
        }
    } else if val, ok := item["lastModifiedByActorId"].(string); ok {
        data.LastModifiedByActorId = types.StringValue(val)
    } else {
        data.LastModifiedByActorId = types.StringNull()
    }
    if obj, ok := item["lastModifiedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByUserId = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByUserId = types.StringNull()
        }
    } else if val, ok := item["lastModifiedByUserId"].(string); ok {
        data.LastModifiedByUserId = types.StringValue(val)
    } else {
        data.LastModifiedByUserId = types.StringNull()
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
    } else if arr, ok := item["projectId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ProjectId = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := item["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := item["incomingCallPolicyId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncomingCallPolicyId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncomingCallPolicyId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncomingCallPolicyId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncomingCallPolicyId = types.StringValue(string(jsonBytes))
        } else {
            data.IncomingCallPolicyId = types.StringNull()
        }
    } else if arr, ok := item["incomingCallPolicyId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.IncomingCallPolicyId = types.StringValue(string(jsonBytes))
        } else {
            data.IncomingCallPolicyId = types.StringNull()
        }
    } else if val, ok := item["incomingCallPolicyId"].(string); ok {
        data.IncomingCallPolicyId = types.StringValue(val)
    } else {
        data.IncomingCallPolicyId = types.StringNull()
    }
    if obj, ok := item["callerPhoneNumber"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CallerPhoneNumber = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CallerPhoneNumber = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CallerPhoneNumber = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CallerPhoneNumber = types.StringValue(string(jsonBytes))
        } else {
            data.CallerPhoneNumber = types.StringNull()
        }
    } else if arr, ok := item["callerPhoneNumber"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CallerPhoneNumber = types.StringValue(string(jsonBytes))
        } else {
            data.CallerPhoneNumber = types.StringNull()
        }
    } else if val, ok := item["callerPhoneNumber"].(string); ok {
        data.CallerPhoneNumber = types.StringValue(val)
    } else {
        data.CallerPhoneNumber = types.StringNull()
    }
    if obj, ok := item["routingPhoneNumber"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RoutingPhoneNumber = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RoutingPhoneNumber = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RoutingPhoneNumber = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RoutingPhoneNumber = types.StringValue(string(jsonBytes))
        } else {
            data.RoutingPhoneNumber = types.StringNull()
        }
    } else if arr, ok := item["routingPhoneNumber"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.RoutingPhoneNumber = types.StringValue(string(jsonBytes))
        } else {
            data.RoutingPhoneNumber = types.StringNull()
        }
    } else if val, ok := item["routingPhoneNumber"].(string); ok {
        data.RoutingPhoneNumber = types.StringValue(val)
    } else {
        data.RoutingPhoneNumber = types.StringNull()
    }
    if obj, ok := item["callProviderCallId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CallProviderCallId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CallProviderCallId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CallProviderCallId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CallProviderCallId = types.StringValue(string(jsonBytes))
        } else {
            data.CallProviderCallId = types.StringNull()
        }
    } else if arr, ok := item["callProviderCallId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CallProviderCallId = types.StringValue(string(jsonBytes))
        } else {
            data.CallProviderCallId = types.StringNull()
        }
    } else if val, ok := item["callProviderCallId"].(string); ok {
        data.CallProviderCallId = types.StringValue(val)
    } else {
        data.CallProviderCallId = types.StringNull()
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
    } else if arr, ok := item["status"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["statusMessage"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.StatusMessage = types.StringValue(string(jsonBytes))
        } else {
            data.StatusMessage = types.StringNull()
        }
    } else if val, ok := item["statusMessage"].(string); ok {
        data.StatusMessage = types.StringValue(val)
    } else {
        data.StatusMessage = types.StringNull()
    }
    data.CallDurationInSeconds = NumberFromAPI(data.CallDurationInSeconds, item["callDurationInSeconds"])
    data.CallCostInUsdCents = NumberFromAPI(data.CallCostInUsdCents, item["callCostInUSDCents"])
    data.IncomingCallCostInUsdCents = NumberFromAPI(data.IncomingCallCostInUsdCents, item["incomingCallCostInUSDCents"])
    data.OutgoingCallCostInUsdCents = NumberFromAPI(data.OutgoingCallCostInUsdCents, item["outgoingCallCostInUSDCents"])
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
    } else if arr, ok := item["startedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["endedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.EndedAt = types.StringValue(string(jsonBytes))
        } else {
            data.EndedAt = types.StringNull()
        }
    } else if val, ok := item["endedAt"].(string); ok {
        data.EndedAt = types.StringValue(val)
    } else {
        data.EndedAt = types.StringNull()
    }
    if obj, ok := item["answeredByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AnsweredByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AnsweredByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AnsweredByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AnsweredByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.AnsweredByUserId = types.StringNull()
        }
    } else if arr, ok := item["answeredByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AnsweredByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.AnsweredByUserId = types.StringNull()
        }
    } else if val, ok := item["answeredByUserId"].(string); ok {
        data.AnsweredByUserId = types.StringValue(val)
    } else {
        data.AnsweredByUserId = types.StringNull()
    }
    data.CurrentEscalationRuleOrder = NumberFromAPI(data.CurrentEscalationRuleOrder, item["currentEscalationRuleOrder"])
    data.RepeatCount = NumberFromAPI(data.RepeatCount, item["repeatCount"])

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
