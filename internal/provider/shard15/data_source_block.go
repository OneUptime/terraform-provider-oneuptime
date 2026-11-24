package shard15

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
var _ datasource.DataSource = &BlockDataSource{}

func NewBlockDataSource() datasource.DataSource {
    return &BlockDataSource{}
}

// BlockDataSource defines the data source implementation.
type BlockDataSource struct {
    client *Client
}

// BlockDataSourceModel describes the data source data model.
type BlockDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    IndicatorType types.String `tfsdk:"indicator_type"`
    Value types.String `tfsdk:"value"`
    Action types.String `tfsdk:"action"`
    State types.String `tfsdk:"state"`
    ExpiresAt types.String `tfsdk:"expires_at"`
    Reason types.String `tfsdk:"reason"`
    SecurityCaseId types.String `tfsdk:"security_case_id"`
    Source types.String `tfsdk:"source"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    SecurityPlaybookRunId types.String `tfsdk:"security_playbook_run_id"`
    LiftedAt types.String `tfsdk:"lifted_at"`
    LiftedByUserId types.String `tfsdk:"lifted_by_user_id"`
    LiftReason types.String `tfsdk:"lift_reason"`
}

func (d *BlockDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_block"
}

func (d *BlockDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "The IP addresses and ranges the project blocks: why, until when, and where each block is enforced. Look up an existing block by `id`, or by any of its other arguments (`action`, `created_by_user_id`, `indicator_type`, ...): each one set must match, and exactly one block may match them all.",

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
                MarkdownDescription: "ID of the project this block belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "indicator_type": schema.StringAttribute{
                MarkdownDescription: "What the block holds: IpAddress for one address, IpRange for a range. Worked out from the value.",
                Optional: true,
                Computed: true,
            },
            "value": schema.StringAttribute{
                MarkdownDescription: "The address or range blocked: an IPv4 or IPv6 address, or a range such as 203.0.113.0/24. Kept as its first address and its length, an IPv6 address compressed in lowercase, a /32 or /128 as the bare address.",
                Optional: true,
                Computed: true,
            },
            "action": schema.StringAttribute{
                MarkdownDescription: "Block turns the address's requests away. Challenge asks for a challenge a person passes and a script or an API client does not.",
                Optional: true,
                Computed: true,
            },
            "state": schema.StringAttribute{
                MarkdownDescription: "Active until it expires or someone lifts it. Set it to Lifted to lift the block; Expired is set by OneUptime when its time is up.",
                Optional: true,
                Computed: true,
            },
            "expires_at": schema.StringAttribute{
                MarkdownDescription: "When the block ends by itself: 24 hours after it is made unless the request says otherwise. Null keeps it until someone lifts it, which takes a Security Admin.",
                Computed: true,
            },
            "reason": schema.StringAttribute{
                MarkdownDescription: "Why the address is blocked, for whoever reads the list later.",
                Optional: true,
                Computed: true,
            },
            "security_case_id": schema.StringAttribute{
                MarkdownDescription: "ID of the case the block was made for. The block and its lifting go on that case's timeline. The ID of a `oneuptime_security_case`.",
                Optional: true,
                Computed: true,
                Validators: RecordIDValidators("oneuptime_security_case"),
            },
            "source": schema.StringAttribute{
                MarkdownDescription: "Who or what made the block: a Person, a Playbook, an AutomaticResponse policy, OneUptimeAi or an Api key.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who blocked it. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
            "security_playbook_run_id": schema.StringAttribute{
                MarkdownDescription: "ID of the playbook run whose step made the block. The ID of a `oneuptime_security_playbook_run`.",
                Optional: true,
                Computed: true,
                Validators: RecordIDValidators("oneuptime_security_playbook_run"),
            },
            "lifted_at": schema.StringAttribute{
                MarkdownDescription: "When someone lifted the block.",
                Computed: true,
            },
            "lifted_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who lifted the block. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
            "lift_reason": schema.StringAttribute{
                MarkdownDescription: "Why the block was lifted, given with the lift.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *BlockDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *BlockDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data BlockDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.IndicatorType.IsNull() && !data.IndicatorType.IsUnknown() {
        filters["indicatorType"] = data.IndicatorType.ValueString()
        filterNames = append(filterNames, "indicator_type = "+fmt.Sprintf("%q", data.IndicatorType.ValueString()))
    }
    if !data.Value.IsNull() && !data.Value.IsUnknown() {
        filters["value"] = data.Value.ValueString()
        filterNames = append(filterNames, "value = "+fmt.Sprintf("%q", data.Value.ValueString()))
    }
    if !data.Action.IsNull() && !data.Action.IsUnknown() {
        filters["action"] = data.Action.ValueString()
        filterNames = append(filterNames, "action = "+fmt.Sprintf("%q", data.Action.ValueString()))
    }
    if !data.State.IsNull() && !data.State.IsUnknown() {
        filters["state"] = data.State.ValueString()
        filterNames = append(filterNames, "state = "+fmt.Sprintf("%q", data.State.ValueString()))
    }
    if !data.Reason.IsNull() && !data.Reason.IsUnknown() {
        filters["reason"] = data.Reason.ValueString()
        filterNames = append(filterNames, "reason = "+fmt.Sprintf("%q", data.Reason.ValueString()))
    }
    if !data.SecurityCaseId.IsNull() && !data.SecurityCaseId.IsUnknown() {
        filters["securityCaseId"] = data.SecurityCaseId.ValueString()
        filterNames = append(filterNames, "security_case_id = "+fmt.Sprintf("%q", data.SecurityCaseId.ValueString()))
    }
    if !data.Source.IsNull() && !data.Source.IsUnknown() {
        filters["source"] = data.Source.ValueString()
        filterNames = append(filterNames, "source = "+fmt.Sprintf("%q", data.Source.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.SecurityPlaybookRunId.IsNull() && !data.SecurityPlaybookRunId.IsUnknown() {
        filters["securityPlaybookRunId"] = data.SecurityPlaybookRunId.ValueString()
        filterNames = append(filterNames, "security_playbook_run_id = "+fmt.Sprintf("%q", data.SecurityPlaybookRunId.ValueString()))
    }
    if !data.LiftedByUserId.IsNull() && !data.LiftedByUserId.IsUnknown() {
        filters["liftedByUserId"] = data.LiftedByUserId.ValueString()
        filterNames = append(filterNames, "lifted_by_user_id = "+fmt.Sprintf("%q", data.LiftedByUserId.ValueString()))
    }
    if !data.LiftReason.IsNull() && !data.LiftReason.IsUnknown() {
        filters["liftReason"] = data.LiftReason.ValueString()
        filterNames = append(filterNames, "lift_reason = "+fmt.Sprintf("%q", data.LiftReason.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the block up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the block up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "indicatorType": true,
        "value": true,
        "action": true,
        "state": true,
        "expiresAt": true,
        "reason": true,
        "securityCaseId": true,
        "source": true,
        "createdByUserId": true,
        "securityPlaybookRunId": true,
        "liftedAt": true,
        "liftedByUserId": true,
        "liftReason": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/security-block/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read block, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No block found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read block: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/security-block/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list block, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list block: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No block matches %s.", DescribeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one block matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", DescribeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for block.")
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
    if obj, ok := item["indicatorType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IndicatorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IndicatorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IndicatorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IndicatorType = types.StringValue(string(jsonBytes))
        } else {
            data.IndicatorType = types.StringNull()
        }
    } else if arr, ok := item["indicatorType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.IndicatorType = types.StringValue(string(jsonBytes))
        } else {
            data.IndicatorType = types.StringNull()
        }
    } else if val, ok := item["indicatorType"].(string); ok {
        data.IndicatorType = types.StringValue(val)
    } else {
        data.IndicatorType = types.StringNull()
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
    if obj, ok := item["action"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Action = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Action = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Action = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Action = types.StringValue(string(jsonBytes))
        } else {
            data.Action = types.StringNull()
        }
    } else if arr, ok := item["action"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Action = types.StringValue(string(jsonBytes))
        } else {
            data.Action = types.StringNull()
        }
    } else if val, ok := item["action"].(string); ok {
        data.Action = types.StringValue(val)
    } else {
        data.Action = types.StringNull()
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
    if obj, ok := item["expiresAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ExpiresAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ExpiresAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ExpiresAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ExpiresAt = types.StringValue(string(jsonBytes))
        } else {
            data.ExpiresAt = types.StringNull()
        }
    } else if arr, ok := item["expiresAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ExpiresAt = types.StringValue(string(jsonBytes))
        } else {
            data.ExpiresAt = types.StringNull()
        }
    } else if val, ok := item["expiresAt"].(string); ok {
        data.ExpiresAt = types.StringValue(val)
    } else {
        data.ExpiresAt = types.StringNull()
    }
    if obj, ok := item["reason"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Reason = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Reason = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Reason = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Reason = types.StringValue(string(jsonBytes))
        } else {
            data.Reason = types.StringNull()
        }
    } else if arr, ok := item["reason"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Reason = types.StringValue(string(jsonBytes))
        } else {
            data.Reason = types.StringNull()
        }
    } else if val, ok := item["reason"].(string); ok {
        data.Reason = types.StringValue(val)
    } else {
        data.Reason = types.StringNull()
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
    if obj, ok := item["source"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Source = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Source = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Source = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Source = types.StringValue(string(jsonBytes))
        } else {
            data.Source = types.StringNull()
        }
    } else if arr, ok := item["source"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Source = types.StringValue(string(jsonBytes))
        } else {
            data.Source = types.StringNull()
        }
    } else if val, ok := item["source"].(string); ok {
        data.Source = types.StringValue(val)
    } else {
        data.Source = types.StringNull()
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
    if obj, ok := item["liftedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LiftedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LiftedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LiftedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LiftedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LiftedAt = types.StringNull()
        }
    } else if arr, ok := item["liftedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LiftedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LiftedAt = types.StringNull()
        }
    } else if val, ok := item["liftedAt"].(string); ok {
        data.LiftedAt = types.StringValue(val)
    } else {
        data.LiftedAt = types.StringNull()
    }
    if obj, ok := item["liftedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LiftedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LiftedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LiftedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LiftedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LiftedByUserId = types.StringNull()
        }
    } else if arr, ok := item["liftedByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LiftedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LiftedByUserId = types.StringNull()
        }
    } else if val, ok := item["liftedByUserId"].(string); ok {
        data.LiftedByUserId = types.StringValue(val)
    } else {
        data.LiftedByUserId = types.StringNull()
    }
    if obj, ok := item["liftReason"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LiftReason = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LiftReason = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LiftReason = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LiftReason = types.StringValue(string(jsonBytes))
        } else {
            data.LiftReason = types.StringNull()
        }
    } else if arr, ok := item["liftReason"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LiftReason = types.StringValue(string(jsonBytes))
        } else {
            data.LiftReason = types.StringNull()
        }
    } else if val, ok := item["liftReason"].(string); ok {
        data.LiftReason = types.StringValue(val)
    } else {
        data.LiftReason = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
