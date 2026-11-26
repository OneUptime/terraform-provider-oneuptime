package shard14

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
var _ datasource.DataSource = &AutomaticResponsePolicyDataSource{}

func NewAutomaticResponsePolicyDataSource() datasource.DataSource {
    return &AutomaticResponsePolicyDataSource{}
}

// AutomaticResponsePolicyDataSource defines the data source implementation.
type AutomaticResponsePolicyDataSource struct {
    client *Client
}

// AutomaticResponsePolicyDataSourceModel describes the data source data model.
type AutomaticResponsePolicyDataSourceModel struct {
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
    Kind types.String `tfsdk:"kind"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    Mode types.String `tfsdk:"mode"`
    BlockAction types.String `tfsdk:"block_action"`
    Duration types.String `tfsdk:"duration"`
    MatchConditions types.String `tfsdk:"match_conditions"`
    TrialStartedAt types.String `tfsdk:"trial_started_at"`
    TrialEndsAt types.String `tfsdk:"trial_ends_at"`
    TrialFinishedAt types.String `tfsdk:"trial_finished_at"`
    TrialSkippedAt types.String `tfsdk:"trial_skipped_at"`
    TrialSkippedByUserId types.String `tfsdk:"trial_skipped_by_user_id"`
    TurnedOnAt types.String `tfsdk:"turned_on_at"`
    TurnedOnByUserId types.String `tfsdk:"turned_on_by_user_id"`
    ModeChangedAt types.String `tfsdk:"mode_changed_at"`
    ModeChangedByUserId types.String `tfsdk:"mode_changed_by_user_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *AutomaticResponsePolicyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_automatic_response_policy"
}

func (d *AutomaticResponsePolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Policies that challenge or block the address a security detection is about as the detection is made, under guardrails: never a protected or own address, never one a signed-in user used, longer blocks for repeat offenders, at most 50 an hour, and a trial run first. Look up an existing automatic response policy by `id`, or by any of its other arguments (`name`, `block_action`, `created_by_actor_id`, ...): each one set must match, and exactly one automatic response policy may match them all.",

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
                MarkdownDescription: "ID of the project the policy acts for. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "kind": schema.StringAttribute{
                MarkdownDescription: "Which policy it is: one of the ready ones (Scanners, Password Guessing, Known-Bad Addresses), or Custom for one the project made.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "The policy's name: Scanners.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "What the policy does, in a sentence.",
                Optional: true,
                Computed: true,
            },
            "mode": schema.StringAttribute{
                MarkdownDescription: "Off (does nothing), Trial Run (records what it would have done, for a week, and does nothing else) or On (acts).",
                Optional: true,
                Computed: true,
            },
            "block_action": schema.StringAttribute{
                MarkdownDescription: "What the policy does to the address: Block, or Challenge (people pass a challenge, scripts do not). Where an edge cannot challenge, a challenge is enforced there as a block for the same time.",
                Optional: true,
                Computed: true,
            },
            "duration": schema.StringAttribute{
                MarkdownDescription: "How long the policy's first block of an address lasts: 1 hour, 24 hours or 7 days. An address blocked before in the last 30 days is blocked for longer.",
                Optional: true,
                Computed: true,
            },
            "match_conditions": schema.StringAttribute{
                MarkdownDescription: "What a detection must be for a project's own policy to act on it: its ATT&CK techniques or detection rules, whether threat intel must list the address, and whether every request so far must have been turned away. A ready policy's are its own. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "trial_started_at": schema.StringAttribute{
                MarkdownDescription: "When the policy's trial run started.",
                Computed: true,
            },
            "trial_ends_at": schema.StringAttribute{
                MarkdownDescription: "When the policy's trial run has run its week, and it asks to be turned on.",
                Computed: true,
            },
            "trial_finished_at": schema.StringAttribute{
                MarkdownDescription: "When the policy last finished a whole week of trial run. Empty while it has never finished one.",
                Computed: true,
            },
            "trial_skipped_at": schema.StringAttribute{
                MarkdownDescription: "When the policy was turned on before a week of trial run, skipping the rest of it.",
                Computed: true,
            },
            "trial_skipped_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who turned the policy on before its trial run was up. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
            "turned_on_at": schema.StringAttribute{
                MarkdownDescription: "When the policy was last turned on: its blocks run without approval from here.",
                Computed: true,
            },
            "turned_on_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who last turned the policy on. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
            "mode_changed_at": schema.StringAttribute{
                MarkdownDescription: "When the policy was last turned on, off or into its trial run.",
                Computed: true,
            },
            "mode_changed_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who last turned the policy on, off or into its trial run. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who made the policy. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
        },
    }
}

func (d *AutomaticResponsePolicyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AutomaticResponsePolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data AutomaticResponsePolicyDataSourceModel

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
    if !data.Kind.IsNull() && !data.Kind.IsUnknown() {
        filters["kind"] = data.Kind.ValueString()
        filterNames = append(filterNames, "kind = "+fmt.Sprintf("%q", data.Kind.ValueString()))
    }
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.Mode.IsNull() && !data.Mode.IsUnknown() {
        filters["mode"] = data.Mode.ValueString()
        filterNames = append(filterNames, "mode = "+fmt.Sprintf("%q", data.Mode.ValueString()))
    }
    if !data.BlockAction.IsNull() && !data.BlockAction.IsUnknown() {
        filters["blockAction"] = data.BlockAction.ValueString()
        filterNames = append(filterNames, "block_action = "+fmt.Sprintf("%q", data.BlockAction.ValueString()))
    }
    if !data.Duration.IsNull() && !data.Duration.IsUnknown() {
        filters["duration"] = data.Duration.ValueString()
        filterNames = append(filterNames, "duration = "+fmt.Sprintf("%q", data.Duration.ValueString()))
    }
    if !data.TrialSkippedByUserId.IsNull() && !data.TrialSkippedByUserId.IsUnknown() {
        filters["trialSkippedByUserId"] = data.TrialSkippedByUserId.ValueString()
        filterNames = append(filterNames, "trial_skipped_by_user_id = "+fmt.Sprintf("%q", data.TrialSkippedByUserId.ValueString()))
    }
    if !data.TurnedOnByUserId.IsNull() && !data.TurnedOnByUserId.IsUnknown() {
        filters["turnedOnByUserId"] = data.TurnedOnByUserId.ValueString()
        filterNames = append(filterNames, "turned_on_by_user_id = "+fmt.Sprintf("%q", data.TurnedOnByUserId.ValueString()))
    }
    if !data.ModeChangedByUserId.IsNull() && !data.ModeChangedByUserId.IsUnknown() {
        filters["modeChangedByUserId"] = data.ModeChangedByUserId.ValueString()
        filterNames = append(filterNames, "mode_changed_by_user_id = "+fmt.Sprintf("%q", data.ModeChangedByUserId.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the automatic response policy up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the automatic response policy up by.",
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
        "kind": true,
        "name": true,
        "description": true,
        "mode": true,
        "blockAction": true,
        "duration": true,
        "matchConditions": true,
        "trialStartedAt": true,
        "trialEndsAt": true,
        "trialFinishedAt": true,
        "trialSkippedAt": true,
        "trialSkippedByUserId": true,
        "turnedOnAt": true,
        "turnedOnByUserId": true,
        "modeChangedAt": true,
        "modeChangedByUserId": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/security-automatic-response-policy/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read automatic_response_policy, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No automatic response policy found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read automatic_response_policy: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/security-automatic-response-policy/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list automatic_response_policy, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list automatic_response_policy: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No automatic response policy matches %s.", DescribeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one automatic response policy matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", DescribeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for automatic_response_policy.")
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
    if obj, ok := item["kind"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Kind = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Kind = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Kind = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Kind = types.StringValue(string(jsonBytes))
        } else {
            data.Kind = types.StringNull()
        }
    } else if arr, ok := item["kind"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Kind = types.StringValue(string(jsonBytes))
        } else {
            data.Kind = types.StringNull()
        }
    } else if val, ok := item["kind"].(string); ok {
        data.Kind = types.StringValue(val)
    } else {
        data.Kind = types.StringNull()
    }
    if obj, ok := item["name"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if arr, ok := item["name"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := item["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := item["description"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if arr, ok := item["description"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := item["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := item["mode"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Mode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Mode = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Mode = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Mode = types.StringValue(string(jsonBytes))
        } else {
            data.Mode = types.StringNull()
        }
    } else if arr, ok := item["mode"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Mode = types.StringValue(string(jsonBytes))
        } else {
            data.Mode = types.StringNull()
        }
    } else if val, ok := item["mode"].(string); ok {
        data.Mode = types.StringValue(val)
    } else {
        data.Mode = types.StringNull()
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
    if obj, ok := item["matchConditions"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MatchConditions = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MatchConditions = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MatchConditions = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MatchConditions = types.StringValue(string(jsonBytes))
        } else {
            data.MatchConditions = types.StringNull()
        }
    } else if arr, ok := item["matchConditions"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.MatchConditions = types.StringValue(string(jsonBytes))
        } else {
            data.MatchConditions = types.StringNull()
        }
    } else if val, ok := item["matchConditions"].(string); ok {
        data.MatchConditions = types.StringValue(val)
    } else {
        data.MatchConditions = types.StringNull()
    }
    if obj, ok := item["trialStartedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TrialStartedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TrialStartedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TrialStartedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TrialStartedAt = types.StringValue(string(jsonBytes))
        } else {
            data.TrialStartedAt = types.StringNull()
        }
    } else if arr, ok := item["trialStartedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.TrialStartedAt = types.StringValue(string(jsonBytes))
        } else {
            data.TrialStartedAt = types.StringNull()
        }
    } else if val, ok := item["trialStartedAt"].(string); ok {
        data.TrialStartedAt = types.StringValue(val)
    } else {
        data.TrialStartedAt = types.StringNull()
    }
    if obj, ok := item["trialEndsAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TrialEndsAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TrialEndsAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TrialEndsAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TrialEndsAt = types.StringValue(string(jsonBytes))
        } else {
            data.TrialEndsAt = types.StringNull()
        }
    } else if arr, ok := item["trialEndsAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.TrialEndsAt = types.StringValue(string(jsonBytes))
        } else {
            data.TrialEndsAt = types.StringNull()
        }
    } else if val, ok := item["trialEndsAt"].(string); ok {
        data.TrialEndsAt = types.StringValue(val)
    } else {
        data.TrialEndsAt = types.StringNull()
    }
    if obj, ok := item["trialFinishedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TrialFinishedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TrialFinishedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TrialFinishedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TrialFinishedAt = types.StringValue(string(jsonBytes))
        } else {
            data.TrialFinishedAt = types.StringNull()
        }
    } else if arr, ok := item["trialFinishedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.TrialFinishedAt = types.StringValue(string(jsonBytes))
        } else {
            data.TrialFinishedAt = types.StringNull()
        }
    } else if val, ok := item["trialFinishedAt"].(string); ok {
        data.TrialFinishedAt = types.StringValue(val)
    } else {
        data.TrialFinishedAt = types.StringNull()
    }
    if obj, ok := item["trialSkippedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TrialSkippedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TrialSkippedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TrialSkippedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TrialSkippedAt = types.StringValue(string(jsonBytes))
        } else {
            data.TrialSkippedAt = types.StringNull()
        }
    } else if arr, ok := item["trialSkippedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.TrialSkippedAt = types.StringValue(string(jsonBytes))
        } else {
            data.TrialSkippedAt = types.StringNull()
        }
    } else if val, ok := item["trialSkippedAt"].(string); ok {
        data.TrialSkippedAt = types.StringValue(val)
    } else {
        data.TrialSkippedAt = types.StringNull()
    }
    if obj, ok := item["trialSkippedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TrialSkippedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TrialSkippedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TrialSkippedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TrialSkippedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.TrialSkippedByUserId = types.StringNull()
        }
    } else if arr, ok := item["trialSkippedByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.TrialSkippedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.TrialSkippedByUserId = types.StringNull()
        }
    } else if val, ok := item["trialSkippedByUserId"].(string); ok {
        data.TrialSkippedByUserId = types.StringValue(val)
    } else {
        data.TrialSkippedByUserId = types.StringNull()
    }
    if obj, ok := item["turnedOnAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TurnedOnAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TurnedOnAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TurnedOnAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TurnedOnAt = types.StringValue(string(jsonBytes))
        } else {
            data.TurnedOnAt = types.StringNull()
        }
    } else if arr, ok := item["turnedOnAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.TurnedOnAt = types.StringValue(string(jsonBytes))
        } else {
            data.TurnedOnAt = types.StringNull()
        }
    } else if val, ok := item["turnedOnAt"].(string); ok {
        data.TurnedOnAt = types.StringValue(val)
    } else {
        data.TurnedOnAt = types.StringNull()
    }
    if obj, ok := item["turnedOnByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TurnedOnByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TurnedOnByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TurnedOnByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TurnedOnByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.TurnedOnByUserId = types.StringNull()
        }
    } else if arr, ok := item["turnedOnByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.TurnedOnByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.TurnedOnByUserId = types.StringNull()
        }
    } else if val, ok := item["turnedOnByUserId"].(string); ok {
        data.TurnedOnByUserId = types.StringValue(val)
    } else {
        data.TurnedOnByUserId = types.StringNull()
    }
    if obj, ok := item["modeChangedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ModeChangedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ModeChangedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ModeChangedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ModeChangedAt = types.StringValue(string(jsonBytes))
        } else {
            data.ModeChangedAt = types.StringNull()
        }
    } else if arr, ok := item["modeChangedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ModeChangedAt = types.StringValue(string(jsonBytes))
        } else {
            data.ModeChangedAt = types.StringNull()
        }
    } else if val, ok := item["modeChangedAt"].(string); ok {
        data.ModeChangedAt = types.StringValue(val)
    } else {
        data.ModeChangedAt = types.StringNull()
    }
    if obj, ok := item["modeChangedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ModeChangedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ModeChangedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ModeChangedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ModeChangedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ModeChangedByUserId = types.StringNull()
        }
    } else if arr, ok := item["modeChangedByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ModeChangedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ModeChangedByUserId = types.StringNull()
        }
    } else if val, ok := item["modeChangedByUserId"].(string); ok {
        data.ModeChangedByUserId = types.StringValue(val)
    } else {
        data.ModeChangedByUserId = types.StringNull()
    }
    if obj, ok := item["createdByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if arr, ok := item["createdByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := item["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
