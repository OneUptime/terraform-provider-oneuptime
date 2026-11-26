package shard06

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
var _ datasource.DataSource = &AutomaticResponseActionDataSource{}

func NewAutomaticResponseActionDataSource() datasource.DataSource {
    return &AutomaticResponseActionDataSource{}
}

// AutomaticResponseActionDataSource defines the data source implementation.
type AutomaticResponseActionDataSource struct {
    client *Client
}

// AutomaticResponseActionDataSourceModel describes the data source data model.
type AutomaticResponseActionDataSourceModel struct {
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
    SecurityAutomaticResponsePolicyId types.String `tfsdk:"security_automatic_response_policy_id"`
    PolicyName types.String `tfsdk:"policy_name"`
    SecurityCaseId types.String `tfsdk:"security_case_id"`
    Value types.String `tfsdk:"value"`
    Outcome types.String `tfsdk:"outcome"`
    IsTrialRun types.Bool `tfsdk:"is_trial_run"`
    BlockAction types.String `tfsdk:"block_action"`
    Duration types.String `tfsdk:"duration"`
    PriorBlockCount types.Number `tfsdk:"prior_block_count"`
    DetectionTitle types.String `tfsdk:"detection_title"`
    FirstRequestAt types.String `tfsdk:"first_request_at"`
    BlockedAt types.String `tfsdk:"blocked_at"`
    AppliedAt types.String `tfsdk:"applied_at"`
    TimeToBlockInSeconds types.Number `tfsdk:"time_to_block_in_seconds"`
    EnforcedAt types.String `tfsdk:"enforced_at"`
    Detail types.String `tfsdk:"detail"`
    PagedAt types.String `tfsdk:"paged_at"`
    DoubtSeenAt types.String `tfsdk:"doubt_seen_at"`
    DoubtDetail types.String `tfsdk:"doubt_detail"`
    SecurityBlockId types.String `tfsdk:"security_block_id"`
    SecurityPlaybookRunId types.String `tfsdk:"security_playbook_run_id"`
}

func (d *AutomaticResponseActionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_automatic_response_action"
}

func (d *AutomaticResponseActionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "What each automatic response policy did about an address: challenged or blocked it (or, in its trial run, would have), or a guardrail spared it - and how long the block took. Look up an existing automatic response action by `id`, or by any of its other arguments (`block_action`, `created_by_actor_id`, `created_by_actor_type`, ...): each one set must match, and exactly one automatic response action may match them all.",

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
                MarkdownDescription: "ID of the project the policy acted for. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "security_automatic_response_policy_id": schema.StringAttribute{
                MarkdownDescription: "ID of the policy that acted. The ID of a `oneuptime_automatic_response_policy`.",
                Optional: true,
                Computed: true,
                Validators: RecordIDValidators("oneuptime_automatic_response_policy"),
            },
            "policy_name": schema.StringAttribute{
                MarkdownDescription: "The policy's name when it acted, kept if the policy is renamed or deleted.",
                Optional: true,
                Computed: true,
            },
            "security_case_id": schema.StringAttribute{
                MarkdownDescription: "ID of the case whose detection the policy acted on. The ID of a `oneuptime_security_case`.",
                Optional: true,
                Computed: true,
                Validators: RecordIDValidators("oneuptime_security_case"),
            },
            "value": schema.StringAttribute{
                MarkdownDescription: "The address the policy acted on: 35.205.11.131.",
                Optional: true,
                Computed: true,
            },
            "outcome": schema.StringAttribute{
                MarkdownDescription: "What the policy did: Challenged or Blocked (in a trial run, what it would have done), or why it left the address alone.",
                Optional: true,
                Computed: true,
            },
            "is_trial_run": schema.BoolAttribute{
                MarkdownDescription: "Decided in the policy's trial run: recorded, and nothing done.",
                Optional: true,
                Computed: true,
            },
            "block_action": schema.StringAttribute{
                MarkdownDescription: "What the policy did, or would have done: Block or Challenge.",
                Optional: true,
                Computed: true,
            },
            "duration": schema.StringAttribute{
                MarkdownDescription: "How long the block lasts: the policy's first length, longer for an address blocked before in the last 30 days.",
                Optional: true,
                Computed: true,
            },
            "prior_block_count": schema.NumberAttribute{
                MarkdownDescription: "How many times the address, or a range holding it, was blocked in the 30 days before.",
                Optional: true,
                Computed: true,
            },
            "detection_title": schema.StringAttribute{
                MarkdownDescription: "The detection the policy acted on, as its case names it.",
                Optional: true,
                Computed: true,
            },
            "first_request_at": schema.StringAttribute{
                MarkdownDescription: "When the address's first request of the detection happened: time to block is measured from here.",
                Computed: true,
            },
            "blocked_at": schema.StringAttribute{
                MarkdownDescription: "When the policy's block was made on the blocklist.",
                Computed: true,
            },
            "applied_at": schema.StringAttribute{
                MarkdownDescription: "When the block first took hold at a destination that enforces the blocklist: Cloudflare, AWS WAF or a firewall.",
                Computed: true,
            },
            "time_to_block_in_seconds": schema.NumberAttribute{
                MarkdownDescription: "Seconds from the address's first request to the block taking hold.",
                Optional: true,
                Computed: true,
            },
            "enforced_at": schema.StringAttribute{
                MarkdownDescription: "Where the block took hold, by name, and whether an edge that cannot challenge holds a challenge as a block there. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "detail": schema.StringAttribute{
                MarkdownDescription: "Why the policy left the address alone or could not block it, in a sentence: what protects it, who signed in from it, or what failed.",
                Optional: true,
                Computed: true,
            },
            "paged_at": schema.StringAttribute{
                MarkdownDescription: "For an action stopped at the hourly limit: when the project's security owners were paged about it.",
                Computed: true,
            },
            "doubt_seen_at": schema.StringAttribute{
                MarkdownDescription: "When a signed-in user was seen from the address after the policy acted on it: the policy lifted its block then (in a trial run, it would have).",
                Computed: true,
            },
            "doubt_detail": schema.StringAttribute{
                MarkdownDescription: "Who was seen signed in from the address after the policy acted on it.",
                Optional: true,
                Computed: true,
            },
            "security_block_id": schema.StringAttribute{
                MarkdownDescription: "ID of the block the policy made. The ID of a `oneuptime_block`.",
                Optional: true,
                Computed: true,
                Validators: RecordIDValidators("oneuptime_block"),
            },
            "security_playbook_run_id": schema.StringAttribute{
                MarkdownDescription: "ID of the run of the playbook engine that made the block. The ID of a `oneuptime_security_playbook_run`.",
                Optional: true,
                Computed: true,
                Validators: RecordIDValidators("oneuptime_security_playbook_run"),
            },
        },
    }
}

func (d *AutomaticResponseActionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AutomaticResponseActionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data AutomaticResponseActionDataSourceModel

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
    if !data.SecurityAutomaticResponsePolicyId.IsNull() && !data.SecurityAutomaticResponsePolicyId.IsUnknown() {
        filters["securityAutomaticResponsePolicyId"] = data.SecurityAutomaticResponsePolicyId.ValueString()
        filterNames = append(filterNames, "security_automatic_response_policy_id = "+fmt.Sprintf("%q", data.SecurityAutomaticResponsePolicyId.ValueString()))
    }
    if !data.PolicyName.IsNull() && !data.PolicyName.IsUnknown() {
        filters["policyName"] = data.PolicyName.ValueString()
        filterNames = append(filterNames, "policy_name = "+fmt.Sprintf("%q", data.PolicyName.ValueString()))
    }
    if !data.SecurityCaseId.IsNull() && !data.SecurityCaseId.IsUnknown() {
        filters["securityCaseId"] = data.SecurityCaseId.ValueString()
        filterNames = append(filterNames, "security_case_id = "+fmt.Sprintf("%q", data.SecurityCaseId.ValueString()))
    }
    if !data.Value.IsNull() && !data.Value.IsUnknown() {
        filters["value"] = data.Value.ValueString()
        filterNames = append(filterNames, "value = "+fmt.Sprintf("%q", data.Value.ValueString()))
    }
    if !data.Outcome.IsNull() && !data.Outcome.IsUnknown() {
        filters["outcome"] = data.Outcome.ValueString()
        filterNames = append(filterNames, "outcome = "+fmt.Sprintf("%q", data.Outcome.ValueString()))
    }
    if !data.IsTrialRun.IsNull() && !data.IsTrialRun.IsUnknown() {
        filters["isTrialRun"] = data.IsTrialRun.ValueBool()
        filterNames = append(filterNames, "is_trial_run = "+fmt.Sprintf("%t", data.IsTrialRun.ValueBool()))
    }
    if !data.BlockAction.IsNull() && !data.BlockAction.IsUnknown() {
        filters["blockAction"] = data.BlockAction.ValueString()
        filterNames = append(filterNames, "block_action = "+fmt.Sprintf("%q", data.BlockAction.ValueString()))
    }
    if !data.Duration.IsNull() && !data.Duration.IsUnknown() {
        filters["duration"] = data.Duration.ValueString()
        filterNames = append(filterNames, "duration = "+fmt.Sprintf("%q", data.Duration.ValueString()))
    }
    if !data.PriorBlockCount.IsNull() && !data.PriorBlockCount.IsUnknown() {
        filters["priorBlockCount"] = LookupNumber(data.PriorBlockCount)
        filterNames = append(filterNames, "prior_block_count = "+data.PriorBlockCount.ValueBigFloat().String())
    }
    if !data.DetectionTitle.IsNull() && !data.DetectionTitle.IsUnknown() {
        filters["detectionTitle"] = data.DetectionTitle.ValueString()
        filterNames = append(filterNames, "detection_title = "+fmt.Sprintf("%q", data.DetectionTitle.ValueString()))
    }
    if !data.TimeToBlockInSeconds.IsNull() && !data.TimeToBlockInSeconds.IsUnknown() {
        filters["timeToBlockInSeconds"] = LookupNumber(data.TimeToBlockInSeconds)
        filterNames = append(filterNames, "time_to_block_in_seconds = "+data.TimeToBlockInSeconds.ValueBigFloat().String())
    }
    if !data.Detail.IsNull() && !data.Detail.IsUnknown() {
        filters["detail"] = data.Detail.ValueString()
        filterNames = append(filterNames, "detail = "+fmt.Sprintf("%q", data.Detail.ValueString()))
    }
    if !data.DoubtDetail.IsNull() && !data.DoubtDetail.IsUnknown() {
        filters["doubtDetail"] = data.DoubtDetail.ValueString()
        filterNames = append(filterNames, "doubt_detail = "+fmt.Sprintf("%q", data.DoubtDetail.ValueString()))
    }
    if !data.SecurityBlockId.IsNull() && !data.SecurityBlockId.IsUnknown() {
        filters["securityBlockId"] = data.SecurityBlockId.ValueString()
        filterNames = append(filterNames, "security_block_id = "+fmt.Sprintf("%q", data.SecurityBlockId.ValueString()))
    }
    if !data.SecurityPlaybookRunId.IsNull() && !data.SecurityPlaybookRunId.IsUnknown() {
        filters["securityPlaybookRunId"] = data.SecurityPlaybookRunId.ValueString()
        filterNames = append(filterNames, "security_playbook_run_id = "+fmt.Sprintf("%q", data.SecurityPlaybookRunId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the automatic response action up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the automatic response action up by.",
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
        "securityAutomaticResponsePolicyId": true,
        "policyName": true,
        "securityCaseId": true,
        "value": true,
        "outcome": true,
        "isTrialRun": true,
        "blockAction": true,
        "duration": true,
        "priorBlockCount": true,
        "detectionTitle": true,
        "firstRequestAt": true,
        "blockedAt": true,
        "appliedAt": true,
        "timeToBlockInSeconds": true,
        "enforcedAt": true,
        "detail": true,
        "pagedAt": true,
        "doubtSeenAt": true,
        "doubtDetail": true,
        "securityBlockId": true,
        "securityPlaybookRunId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/security-automatic-response-action/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read automatic_response_action, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No automatic response action found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read automatic_response_action: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/security-automatic-response-action/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list automatic_response_action, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list automatic_response_action: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No automatic response action matches %s.", DescribeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one automatic response action matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", DescribeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for automatic_response_action.")
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
    if obj, ok := item["securityAutomaticResponsePolicyId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SecurityAutomaticResponsePolicyId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SecurityAutomaticResponsePolicyId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SecurityAutomaticResponsePolicyId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SecurityAutomaticResponsePolicyId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityAutomaticResponsePolicyId = types.StringNull()
        }
    } else if arr, ok := item["securityAutomaticResponsePolicyId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.SecurityAutomaticResponsePolicyId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityAutomaticResponsePolicyId = types.StringNull()
        }
    } else if val, ok := item["securityAutomaticResponsePolicyId"].(string); ok {
        data.SecurityAutomaticResponsePolicyId = types.StringValue(val)
    } else {
        data.SecurityAutomaticResponsePolicyId = types.StringNull()
    }
    if obj, ok := item["policyName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PolicyName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PolicyName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PolicyName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PolicyName = types.StringValue(string(jsonBytes))
        } else {
            data.PolicyName = types.StringNull()
        }
    } else if arr, ok := item["policyName"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.PolicyName = types.StringValue(string(jsonBytes))
        } else {
            data.PolicyName = types.StringNull()
        }
    } else if val, ok := item["policyName"].(string); ok {
        data.PolicyName = types.StringValue(val)
    } else {
        data.PolicyName = types.StringNull()
    }
    if obj, ok := item["securityCaseId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SecurityCaseId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SecurityCaseId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SecurityCaseId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SecurityCaseId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityCaseId = types.StringNull()
        }
    } else if arr, ok := item["securityCaseId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.SecurityCaseId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityCaseId = types.StringNull()
        }
    } else if val, ok := item["securityCaseId"].(string); ok {
        data.SecurityCaseId = types.StringValue(val)
    } else {
        data.SecurityCaseId = types.StringNull()
    }
    if obj, ok := item["value"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Value = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Value = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Value = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Value = types.StringValue(string(jsonBytes))
        } else {
            data.Value = types.StringNull()
        }
    } else if arr, ok := item["value"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Value = types.StringValue(string(jsonBytes))
        } else {
            data.Value = types.StringNull()
        }
    } else if val, ok := item["value"].(string); ok {
        data.Value = types.StringValue(val)
    } else {
        data.Value = types.StringNull()
    }
    if obj, ok := item["outcome"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Outcome = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Outcome = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Outcome = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Outcome = types.StringValue(string(jsonBytes))
        } else {
            data.Outcome = types.StringNull()
        }
    } else if arr, ok := item["outcome"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Outcome = types.StringValue(string(jsonBytes))
        } else {
            data.Outcome = types.StringNull()
        }
    } else if val, ok := item["outcome"].(string); ok {
        data.Outcome = types.StringValue(val)
    } else {
        data.Outcome = types.StringNull()
    }
    if val, ok := item["isTrialRun"].(bool); ok {
        data.IsTrialRun = types.BoolValue(val)
    } else {
        data.IsTrialRun = types.BoolNull()
    }
    if obj, ok := item["blockAction"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.BlockAction = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.BlockAction = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.BlockAction = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.BlockAction = types.StringValue(string(jsonBytes))
        } else {
            data.BlockAction = types.StringNull()
        }
    } else if arr, ok := item["blockAction"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.BlockAction = types.StringValue(string(jsonBytes))
        } else {
            data.BlockAction = types.StringNull()
        }
    } else if val, ok := item["blockAction"].(string); ok {
        data.BlockAction = types.StringValue(val)
    } else {
        data.BlockAction = types.StringNull()
    }
    if obj, ok := item["duration"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Duration = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Duration = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Duration = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Duration = types.StringValue(string(jsonBytes))
        } else {
            data.Duration = types.StringNull()
        }
    } else if arr, ok := item["duration"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Duration = types.StringValue(string(jsonBytes))
        } else {
            data.Duration = types.StringNull()
        }
    } else if val, ok := item["duration"].(string); ok {
        data.Duration = types.StringValue(val)
    } else {
        data.Duration = types.StringNull()
    }
    data.PriorBlockCount = NumberFromAPI(data.PriorBlockCount, item["priorBlockCount"])
    if obj, ok := item["detectionTitle"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DetectionTitle = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DetectionTitle = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DetectionTitle = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DetectionTitle = types.StringValue(string(jsonBytes))
        } else {
            data.DetectionTitle = types.StringNull()
        }
    } else if arr, ok := item["detectionTitle"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.DetectionTitle = types.StringValue(string(jsonBytes))
        } else {
            data.DetectionTitle = types.StringNull()
        }
    } else if val, ok := item["detectionTitle"].(string); ok {
        data.DetectionTitle = types.StringValue(val)
    } else {
        data.DetectionTitle = types.StringNull()
    }
    if obj, ok := item["firstRequestAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FirstRequestAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FirstRequestAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FirstRequestAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FirstRequestAt = types.StringValue(string(jsonBytes))
        } else {
            data.FirstRequestAt = types.StringNull()
        }
    } else if arr, ok := item["firstRequestAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.FirstRequestAt = types.StringValue(string(jsonBytes))
        } else {
            data.FirstRequestAt = types.StringNull()
        }
    } else if val, ok := item["firstRequestAt"].(string); ok {
        data.FirstRequestAt = types.StringValue(val)
    } else {
        data.FirstRequestAt = types.StringNull()
    }
    if obj, ok := item["blockedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.BlockedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.BlockedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.BlockedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.BlockedAt = types.StringValue(string(jsonBytes))
        } else {
            data.BlockedAt = types.StringNull()
        }
    } else if arr, ok := item["blockedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.BlockedAt = types.StringValue(string(jsonBytes))
        } else {
            data.BlockedAt = types.StringNull()
        }
    } else if val, ok := item["blockedAt"].(string); ok {
        data.BlockedAt = types.StringValue(val)
    } else {
        data.BlockedAt = types.StringNull()
    }
    if obj, ok := item["appliedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AppliedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AppliedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AppliedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AppliedAt = types.StringValue(string(jsonBytes))
        } else {
            data.AppliedAt = types.StringNull()
        }
    } else if arr, ok := item["appliedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AppliedAt = types.StringValue(string(jsonBytes))
        } else {
            data.AppliedAt = types.StringNull()
        }
    } else if val, ok := item["appliedAt"].(string); ok {
        data.AppliedAt = types.StringValue(val)
    } else {
        data.AppliedAt = types.StringNull()
    }
    data.TimeToBlockInSeconds = NumberFromAPI(data.TimeToBlockInSeconds, item["timeToBlockInSeconds"])
    if obj, ok := item["enforcedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EnforcedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EnforcedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EnforcedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EnforcedAt = types.StringValue(string(jsonBytes))
        } else {
            data.EnforcedAt = types.StringNull()
        }
    } else if arr, ok := item["enforcedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.EnforcedAt = types.StringValue(string(jsonBytes))
        } else {
            data.EnforcedAt = types.StringNull()
        }
    } else if val, ok := item["enforcedAt"].(string); ok {
        data.EnforcedAt = types.StringValue(val)
    } else {
        data.EnforcedAt = types.StringNull()
    }
    if obj, ok := item["detail"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Detail = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Detail = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Detail = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Detail = types.StringValue(string(jsonBytes))
        } else {
            data.Detail = types.StringNull()
        }
    } else if arr, ok := item["detail"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Detail = types.StringValue(string(jsonBytes))
        } else {
            data.Detail = types.StringNull()
        }
    } else if val, ok := item["detail"].(string); ok {
        data.Detail = types.StringValue(val)
    } else {
        data.Detail = types.StringNull()
    }
    if obj, ok := item["pagedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PagedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PagedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PagedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PagedAt = types.StringValue(string(jsonBytes))
        } else {
            data.PagedAt = types.StringNull()
        }
    } else if arr, ok := item["pagedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.PagedAt = types.StringValue(string(jsonBytes))
        } else {
            data.PagedAt = types.StringNull()
        }
    } else if val, ok := item["pagedAt"].(string); ok {
        data.PagedAt = types.StringValue(val)
    } else {
        data.PagedAt = types.StringNull()
    }
    if obj, ok := item["doubtSeenAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DoubtSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DoubtSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DoubtSeenAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DoubtSeenAt = types.StringValue(string(jsonBytes))
        } else {
            data.DoubtSeenAt = types.StringNull()
        }
    } else if arr, ok := item["doubtSeenAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.DoubtSeenAt = types.StringValue(string(jsonBytes))
        } else {
            data.DoubtSeenAt = types.StringNull()
        }
    } else if val, ok := item["doubtSeenAt"].(string); ok {
        data.DoubtSeenAt = types.StringValue(val)
    } else {
        data.DoubtSeenAt = types.StringNull()
    }
    if obj, ok := item["doubtDetail"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DoubtDetail = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DoubtDetail = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DoubtDetail = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DoubtDetail = types.StringValue(string(jsonBytes))
        } else {
            data.DoubtDetail = types.StringNull()
        }
    } else if arr, ok := item["doubtDetail"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.DoubtDetail = types.StringValue(string(jsonBytes))
        } else {
            data.DoubtDetail = types.StringNull()
        }
    } else if val, ok := item["doubtDetail"].(string); ok {
        data.DoubtDetail = types.StringValue(val)
    } else {
        data.DoubtDetail = types.StringNull()
    }
    if obj, ok := item["securityBlockId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SecurityBlockId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SecurityBlockId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SecurityBlockId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SecurityBlockId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityBlockId = types.StringNull()
        }
    } else if arr, ok := item["securityBlockId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.SecurityBlockId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityBlockId = types.StringNull()
        }
    } else if val, ok := item["securityBlockId"].(string); ok {
        data.SecurityBlockId = types.StringValue(val)
    } else {
        data.SecurityBlockId = types.StringNull()
    }
    if obj, ok := item["securityPlaybookRunId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SecurityPlaybookRunId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SecurityPlaybookRunId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SecurityPlaybookRunId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SecurityPlaybookRunId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityPlaybookRunId = types.StringNull()
        }
    } else if arr, ok := item["securityPlaybookRunId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.SecurityPlaybookRunId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityPlaybookRunId = types.StringNull()
        }
    } else if val, ok := item["securityPlaybookRunId"].(string); ok {
        data.SecurityPlaybookRunId = types.StringValue(val)
    } else {
        data.SecurityPlaybookRunId = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
