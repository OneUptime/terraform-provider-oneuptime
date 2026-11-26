package shard00

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
var _ datasource.DataSource = &BlockEnforcementDataSource{}

func NewBlockEnforcementDataSource() datasource.DataSource {
    return &BlockEnforcementDataSource{}
}

// BlockEnforcementDataSource defines the data source implementation.
type BlockEnforcementDataSource struct {
    client *Client
}

// BlockEnforcementDataSourceModel describes the data source data model.
type BlockEnforcementDataSourceModel struct {
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
    SecurityBlockId types.String `tfsdk:"security_block_id"`
    DestinationKey types.String `tfsdk:"destination_key"`
    ProviderValue types.String `tfsdk:"provider_value"`
    DestinationName types.String `tfsdk:"destination_name"`
    ConnectionId types.String `tfsdk:"connection_id"`
    State types.String `tfsdk:"state"`
    AppliedAt types.String `tfsdk:"applied_at"`
    LastConfirmedAt types.String `tfsdk:"last_confirmed_at"`
    ProviderRef types.String `tfsdk:"provider_ref"`
    Error types.String `tfsdk:"error"`
    ErrorReason types.String `tfsdk:"error_reason"`
    FailedAttempts types.Number `tfsdk:"failed_attempts"`
    RetryAt types.String `tfsdk:"retry_at"`
}

func (d *BlockEnforcementDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_block_enforcement"
}

func (d *BlockEnforcementDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Where each block of the blocklist is enforced, and whether that holds: one row per block and destination, written by OneUptime as it applies the block and reads it back. Look up an existing block enforcement by `id`, or by any of its other arguments (`connection_id`, `created_by_actor_id`, `created_by_actor_type`, ...): each one set must match, and exactly one block enforcement may match them all.",

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
                MarkdownDescription: "ID of the project this row belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "security_block_id": schema.StringAttribute{
                MarkdownDescription: "ID of the block this row says the enforcement of. The ID of a `oneuptime_block`.",
                Optional: true,
                Computed: true,
                Validators: RecordIDValidators("oneuptime_block"),
            },
            "destination_key": schema.StringAttribute{
                MarkdownDescription: "Which destination this is, as OneUptime keys it: its provider and its connection, and what in it (a zone, a web ACL) when it holds more than one.",
                Optional: true,
                Computed: true,
            },
            "provider_value": schema.StringAttribute{
                MarkdownDescription: "The kind of destination: Cloudflare, AWS WAF, Azure WAF, Cloud Armor, a feed, Microsoft Defender or CrowdStrike.",
                Optional: true,
                Computed: true,
            },
            "destination_name": schema.StringAttribute{
                MarkdownDescription: "The destination as people read it: Cloudflare, AWS WAF (prod-alb).",
                Optional: true,
                Computed: true,
            },
            "connection_id": schema.StringAttribute{
                MarkdownDescription: "ID of the connection the destination is enforced through: a security connection, a cloud connection or a feed.",
                Optional: true,
                Computed: true,
            },
            "state": schema.StringAttribute{
                MarkdownDescription: "Pending until the destination is seen holding the block, Applied once it is, Failed with why it was not applied, then Removing and Removed when the block ends.",
                Optional: true,
                Computed: true,
            },
            "applied_at": schema.StringAttribute{
                MarkdownDescription: "When the destination was first seen holding the block.",
                Computed: true,
            },
            "last_confirmed_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime last read the destination back and found the block there.",
                Computed: true,
            },
            "provider_ref": schema.StringAttribute{
                MarkdownDescription: "The provider's own ID of what holds the block: a list item, a rule, an indicator.",
                Optional: true,
                Computed: true,
            },
            "error": schema.StringAttribute{
                MarkdownDescription: "Why the block is not applied there: the destination cannot challenge, matches IPv4 only, is full, or the provider refused.",
                Optional: true,
                Computed: true,
            },
            "error_reason": schema.StringAttribute{
                MarkdownDescription: "Why the block is not applied there, as a sentence and its values that the Dashboard words in the reader's language. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "failed_attempts": schema.NumberAttribute{
                MarkdownDescription: "How many tries in a row failed at the destination; the next waits longer each time.",
                Optional: true,
                Computed: true,
            },
            "retry_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime tries the destination again after a failure; empty when nothing waits.",
                Computed: true,
            },
        },
    }
}

func (d *BlockEnforcementDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *BlockEnforcementDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data BlockEnforcementDataSourceModel

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
    if !data.SecurityBlockId.IsNull() && !data.SecurityBlockId.IsUnknown() {
        filters["securityBlockId"] = data.SecurityBlockId.ValueString()
        filterNames = append(filterNames, "security_block_id = "+fmt.Sprintf("%q", data.SecurityBlockId.ValueString()))
    }
    if !data.DestinationKey.IsNull() && !data.DestinationKey.IsUnknown() {
        filters["destinationKey"] = data.DestinationKey.ValueString()
        filterNames = append(filterNames, "destination_key = "+fmt.Sprintf("%q", data.DestinationKey.ValueString()))
    }
    if !data.ProviderValue.IsNull() && !data.ProviderValue.IsUnknown() {
        filters["provider"] = data.ProviderValue.ValueString()
        filterNames = append(filterNames, "provider_value = "+fmt.Sprintf("%q", data.ProviderValue.ValueString()))
    }
    if !data.DestinationName.IsNull() && !data.DestinationName.IsUnknown() {
        filters["destinationName"] = data.DestinationName.ValueString()
        filterNames = append(filterNames, "destination_name = "+fmt.Sprintf("%q", data.DestinationName.ValueString()))
    }
    if !data.ConnectionId.IsNull() && !data.ConnectionId.IsUnknown() {
        filters["connectionId"] = data.ConnectionId.ValueString()
        filterNames = append(filterNames, "connection_id = "+fmt.Sprintf("%q", data.ConnectionId.ValueString()))
    }
    if !data.State.IsNull() && !data.State.IsUnknown() {
        filters["state"] = data.State.ValueString()
        filterNames = append(filterNames, "state = "+fmt.Sprintf("%q", data.State.ValueString()))
    }
    if !data.ProviderRef.IsNull() && !data.ProviderRef.IsUnknown() {
        filters["providerRef"] = data.ProviderRef.ValueString()
        filterNames = append(filterNames, "provider_ref = "+fmt.Sprintf("%q", data.ProviderRef.ValueString()))
    }
    if !data.Error.IsNull() && !data.Error.IsUnknown() {
        filters["error"] = data.Error.ValueString()
        filterNames = append(filterNames, "error = "+fmt.Sprintf("%q", data.Error.ValueString()))
    }
    if !data.FailedAttempts.IsNull() && !data.FailedAttempts.IsUnknown() {
        filters["failedAttempts"] = LookupNumber(data.FailedAttempts)
        filterNames = append(filterNames, "failed_attempts = "+data.FailedAttempts.ValueBigFloat().String())
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the block enforcement up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the block enforcement up by.",
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
        "securityBlockId": true,
        "destinationKey": true,
        "provider": true,
        "destinationName": true,
        "connectionId": true,
        "state": true,
        "appliedAt": true,
        "lastConfirmedAt": true,
        "providerRef": true,
        "error": true,
        "errorReason": true,
        "failedAttempts": true,
        "retryAt": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/security-block-enforcement/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read block_enforcement, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No block enforcement found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read block_enforcement: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/security-block-enforcement/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list block_enforcement, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list block_enforcement: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No block enforcement matches %s.", DescribeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one block enforcement matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", DescribeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for block_enforcement.")
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
    if obj, ok := item["destinationKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DestinationKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DestinationKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DestinationKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DestinationKey = types.StringValue(string(jsonBytes))
        } else {
            data.DestinationKey = types.StringNull()
        }
    } else if arr, ok := item["destinationKey"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.DestinationKey = types.StringValue(string(jsonBytes))
        } else {
            data.DestinationKey = types.StringNull()
        }
    } else if val, ok := item["destinationKey"].(string); ok {
        data.DestinationKey = types.StringValue(val)
    } else {
        data.DestinationKey = types.StringNull()
    }
    if obj, ok := item["provider"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProviderValue = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProviderValue = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProviderValue = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProviderValue = types.StringValue(string(jsonBytes))
        } else {
            data.ProviderValue = types.StringNull()
        }
    } else if arr, ok := item["provider"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ProviderValue = types.StringValue(string(jsonBytes))
        } else {
            data.ProviderValue = types.StringNull()
        }
    } else if val, ok := item["provider"].(string); ok {
        data.ProviderValue = types.StringValue(val)
    } else {
        data.ProviderValue = types.StringNull()
    }
    if obj, ok := item["destinationName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DestinationName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DestinationName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DestinationName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DestinationName = types.StringValue(string(jsonBytes))
        } else {
            data.DestinationName = types.StringNull()
        }
    } else if arr, ok := item["destinationName"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.DestinationName = types.StringValue(string(jsonBytes))
        } else {
            data.DestinationName = types.StringNull()
        }
    } else if val, ok := item["destinationName"].(string); ok {
        data.DestinationName = types.StringValue(val)
    } else {
        data.DestinationName = types.StringNull()
    }
    if obj, ok := item["connectionId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ConnectionId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ConnectionId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ConnectionId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ConnectionId = types.StringValue(string(jsonBytes))
        } else {
            data.ConnectionId = types.StringNull()
        }
    } else if arr, ok := item["connectionId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ConnectionId = types.StringValue(string(jsonBytes))
        } else {
            data.ConnectionId = types.StringNull()
        }
    } else if val, ok := item["connectionId"].(string); ok {
        data.ConnectionId = types.StringValue(val)
    } else {
        data.ConnectionId = types.StringNull()
    }
    if obj, ok := item["state"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.State = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.State = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.State = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.State = types.StringValue(string(jsonBytes))
        } else {
            data.State = types.StringNull()
        }
    } else if arr, ok := item["state"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.State = types.StringValue(string(jsonBytes))
        } else {
            data.State = types.StringNull()
        }
    } else if val, ok := item["state"].(string); ok {
        data.State = types.StringValue(val)
    } else {
        data.State = types.StringNull()
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
    if obj, ok := item["lastConfirmedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastConfirmedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastConfirmedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastConfirmedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastConfirmedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastConfirmedAt = types.StringNull()
        }
    } else if arr, ok := item["lastConfirmedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastConfirmedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastConfirmedAt = types.StringNull()
        }
    } else if val, ok := item["lastConfirmedAt"].(string); ok {
        data.LastConfirmedAt = types.StringValue(val)
    } else {
        data.LastConfirmedAt = types.StringNull()
    }
    if obj, ok := item["providerRef"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProviderRef = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProviderRef = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProviderRef = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProviderRef = types.StringValue(string(jsonBytes))
        } else {
            data.ProviderRef = types.StringNull()
        }
    } else if arr, ok := item["providerRef"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ProviderRef = types.StringValue(string(jsonBytes))
        } else {
            data.ProviderRef = types.StringNull()
        }
    } else if val, ok := item["providerRef"].(string); ok {
        data.ProviderRef = types.StringValue(val)
    } else {
        data.ProviderRef = types.StringNull()
    }
    if obj, ok := item["error"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Error = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Error = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Error = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Error = types.StringValue(string(jsonBytes))
        } else {
            data.Error = types.StringNull()
        }
    } else if arr, ok := item["error"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Error = types.StringValue(string(jsonBytes))
        } else {
            data.Error = types.StringNull()
        }
    } else if val, ok := item["error"].(string); ok {
        data.Error = types.StringValue(val)
    } else {
        data.Error = types.StringNull()
    }
    if obj, ok := item["errorReason"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ErrorReason = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ErrorReason = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ErrorReason = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ErrorReason = types.StringValue(string(jsonBytes))
        } else {
            data.ErrorReason = types.StringNull()
        }
    } else if arr, ok := item["errorReason"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ErrorReason = types.StringValue(string(jsonBytes))
        } else {
            data.ErrorReason = types.StringNull()
        }
    } else if val, ok := item["errorReason"].(string); ok {
        data.ErrorReason = types.StringValue(val)
    } else {
        data.ErrorReason = types.StringNull()
    }
    data.FailedAttempts = NumberFromAPI(data.FailedAttempts, item["failedAttempts"])
    if obj, ok := item["retryAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RetryAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RetryAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RetryAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RetryAt = types.StringValue(string(jsonBytes))
        } else {
            data.RetryAt = types.StringNull()
        }
    } else if arr, ok := item["retryAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.RetryAt = types.StringValue(string(jsonBytes))
        } else {
            data.RetryAt = types.StringNull()
        }
    } else if val, ok := item["retryAt"].(string); ok {
        data.RetryAt = types.StringValue(val)
    } else {
        data.RetryAt = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
