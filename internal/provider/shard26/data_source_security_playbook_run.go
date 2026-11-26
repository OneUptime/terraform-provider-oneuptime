package shard26

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
var _ datasource.DataSource = &SecurityPlaybookRunDataSource{}

func NewSecurityPlaybookRunDataSource() datasource.DataSource {
    return &SecurityPlaybookRunDataSource{}
}

// SecurityPlaybookRunDataSource defines the data source implementation.
type SecurityPlaybookRunDataSource struct {
    client *Client
}

// SecurityPlaybookRunDataSourceModel describes the data source data model.
type SecurityPlaybookRunDataSourceModel struct {
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
    SecurityCaseId types.String `tfsdk:"security_case_id"`
    SecurityPlaybookId types.String `tfsdk:"security_playbook_id"`
    PlaybookName types.String `tfsdk:"playbook_name"`
    Trigger types.String `tfsdk:"trigger"`
    TriggeredByUserId types.String `tfsdk:"triggered_by_user_id"`
    SecurityAutomaticResponsePolicyId types.String `tfsdk:"security_automatic_response_policy_id"`
    Status types.String `tfsdk:"status"`
    Steps types.String `tfsdk:"steps"`
    StepResults types.String `tfsdk:"step_results"`
    CurrentStepIndex types.Number `tfsdk:"current_step_index"`
    ThreatIntelMatched types.Bool `tfsdk:"threat_intel_matched"`
    StartedAt types.String `tfsdk:"started_at"`
    CompletedAt types.String `tfsdk:"completed_at"`
    ErrorMessage types.String `tfsdk:"error_message"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *SecurityPlaybookRunDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_security_playbook_run"
}

func (d *SecurityPlaybookRunDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "One run of a security playbook on a case: the steps it took, what each did, and who approved the steps that asked. Look up an existing security playbook run by `id`, or by any of its other arguments (`created_by_actor_id`, `created_by_actor_type`, `created_by_user_id`, ...): each one set must match, and exactly one security playbook run may match them all.",

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
                MarkdownDescription: "ID of the project this run belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "security_case_id": schema.StringAttribute{
                MarkdownDescription: "ID of the case the playbook runs on. The ID of a `oneuptime_security_case`.",
                Optional: true,
                Computed: true,
                Validators: RecordIDValidators("oneuptime_security_case"),
            },
            "security_playbook_id": schema.StringAttribute{
                MarkdownDescription: "ID of the playbook being run. The ID of a `oneuptime_security_playbook`.",
                Optional: true,
                Computed: true,
                Validators: RecordIDValidators("oneuptime_security_playbook"),
            },
            "playbook_name": schema.StringAttribute{
                MarkdownDescription: "The playbook's name when the run started, kept if the playbook is renamed or deleted.",
                Optional: true,
                Computed: true,
            },
            "trigger": schema.StringAttribute{
                MarkdownDescription: "What started the run: Case Opened, Detection Added, Manual for a person, or Automatic Response for an automatic response policy's block.",
                Optional: true,
                Computed: true,
            },
            "triggered_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who ran the playbook. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
            "security_automatic_response_policy_id": schema.StringAttribute{
                MarkdownDescription: "ID of the automatic response policy that started the run. The ID of a `oneuptime_automatic_response_policy`.",
                Optional: true,
                Computed: true,
                Validators: RecordIDValidators("oneuptime_automatic_response_policy"),
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Queued, Running, Waiting for Approval, Succeeded, Failed or Cancelled. Set it to Cancelled to stop a run that has not finished.",
                Optional: true,
                Computed: true,
            },
            "steps": schema.StringAttribute{
                MarkdownDescription: "The playbook's steps as they were when the run started. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "step_results": schema.StringAttribute{
                MarkdownDescription: "What each step did: its status, when it started and finished, what it reported, and who decided on its approval. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "current_step_index": schema.NumberAttribute{
                MarkdownDescription: "The step the run is on, counting from 0.",
                Optional: true,
                Computed: true,
            },
            "threat_intel_matched": schema.BoolAttribute{
                MarkdownDescription: "Whether a Check Threat Intel step of this run found a known-bad entity. Steps set to run only on a match read it.",
                Optional: true,
                Computed: true,
            },
            "started_at": schema.StringAttribute{
                MarkdownDescription: "When the run started its first step.",
                Computed: true,
            },
            "completed_at": schema.StringAttribute{
                MarkdownDescription: "When the run finished, failed or was cancelled.",
                Computed: true,
            },
            "error_message": schema.StringAttribute{
                MarkdownDescription: "Why the run failed, when it did.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who started the run. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
        },
    }
}

func (d *SecurityPlaybookRunDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SecurityPlaybookRunDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data SecurityPlaybookRunDataSourceModel

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
    if !data.SecurityCaseId.IsNull() && !data.SecurityCaseId.IsUnknown() {
        filters["securityCaseId"] = data.SecurityCaseId.ValueString()
        filterNames = append(filterNames, "security_case_id = "+fmt.Sprintf("%q", data.SecurityCaseId.ValueString()))
    }
    if !data.SecurityPlaybookId.IsNull() && !data.SecurityPlaybookId.IsUnknown() {
        filters["securityPlaybookId"] = data.SecurityPlaybookId.ValueString()
        filterNames = append(filterNames, "security_playbook_id = "+fmt.Sprintf("%q", data.SecurityPlaybookId.ValueString()))
    }
    if !data.PlaybookName.IsNull() && !data.PlaybookName.IsUnknown() {
        filters["playbookName"] = data.PlaybookName.ValueString()
        filterNames = append(filterNames, "playbook_name = "+fmt.Sprintf("%q", data.PlaybookName.ValueString()))
    }
    if !data.Trigger.IsNull() && !data.Trigger.IsUnknown() {
        filters["trigger"] = data.Trigger.ValueString()
        filterNames = append(filterNames, "trigger = "+fmt.Sprintf("%q", data.Trigger.ValueString()))
    }
    if !data.TriggeredByUserId.IsNull() && !data.TriggeredByUserId.IsUnknown() {
        filters["triggeredByUserId"] = data.TriggeredByUserId.ValueString()
        filterNames = append(filterNames, "triggered_by_user_id = "+fmt.Sprintf("%q", data.TriggeredByUserId.ValueString()))
    }
    if !data.SecurityAutomaticResponsePolicyId.IsNull() && !data.SecurityAutomaticResponsePolicyId.IsUnknown() {
        filters["securityAutomaticResponsePolicyId"] = data.SecurityAutomaticResponsePolicyId.ValueString()
        filterNames = append(filterNames, "security_automatic_response_policy_id = "+fmt.Sprintf("%q", data.SecurityAutomaticResponsePolicyId.ValueString()))
    }
    if !data.Status.IsNull() && !data.Status.IsUnknown() {
        filters["status"] = data.Status.ValueString()
        filterNames = append(filterNames, "status = "+fmt.Sprintf("%q", data.Status.ValueString()))
    }
    if !data.CurrentStepIndex.IsNull() && !data.CurrentStepIndex.IsUnknown() {
        filters["currentStepIndex"] = LookupNumber(data.CurrentStepIndex)
        filterNames = append(filterNames, "current_step_index = "+data.CurrentStepIndex.ValueBigFloat().String())
    }
    if !data.ThreatIntelMatched.IsNull() && !data.ThreatIntelMatched.IsUnknown() {
        filters["threatIntelMatched"] = data.ThreatIntelMatched.ValueBool()
        filterNames = append(filterNames, "threat_intel_matched = "+fmt.Sprintf("%t", data.ThreatIntelMatched.ValueBool()))
    }
    if !data.ErrorMessage.IsNull() && !data.ErrorMessage.IsUnknown() {
        filters["errorMessage"] = data.ErrorMessage.ValueString()
        filterNames = append(filterNames, "error_message = "+fmt.Sprintf("%q", data.ErrorMessage.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the security playbook run up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the security playbook run up by.",
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
        "securityCaseId": true,
        "securityPlaybookId": true,
        "playbookName": true,
        "trigger": true,
        "triggeredByUserId": true,
        "securityAutomaticResponsePolicyId": true,
        "status": true,
        "steps": true,
        "stepResults": true,
        "currentStepIndex": true,
        "threatIntelMatched": true,
        "startedAt": true,
        "completedAt": true,
        "errorMessage": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/security-playbook-run/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read security_playbook_run, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No security playbook run found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read security_playbook_run: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/security-playbook-run/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list security_playbook_run, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list security_playbook_run: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No security playbook run matches %s.", DescribeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one security playbook run matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", DescribeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for security_playbook_run.")
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
    if obj, ok := item["securityPlaybookId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SecurityPlaybookId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SecurityPlaybookId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SecurityPlaybookId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SecurityPlaybookId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityPlaybookId = types.StringNull()
        }
    } else if arr, ok := item["securityPlaybookId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.SecurityPlaybookId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityPlaybookId = types.StringNull()
        }
    } else if val, ok := item["securityPlaybookId"].(string); ok {
        data.SecurityPlaybookId = types.StringValue(val)
    } else {
        data.SecurityPlaybookId = types.StringNull()
    }
    if obj, ok := item["playbookName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PlaybookName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PlaybookName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PlaybookName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PlaybookName = types.StringValue(string(jsonBytes))
        } else {
            data.PlaybookName = types.StringNull()
        }
    } else if arr, ok := item["playbookName"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.PlaybookName = types.StringValue(string(jsonBytes))
        } else {
            data.PlaybookName = types.StringNull()
        }
    } else if val, ok := item["playbookName"].(string); ok {
        data.PlaybookName = types.StringValue(val)
    } else {
        data.PlaybookName = types.StringNull()
    }
    if obj, ok := item["trigger"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Trigger = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Trigger = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Trigger = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Trigger = types.StringValue(string(jsonBytes))
        } else {
            data.Trigger = types.StringNull()
        }
    } else if arr, ok := item["trigger"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Trigger = types.StringValue(string(jsonBytes))
        } else {
            data.Trigger = types.StringNull()
        }
    } else if val, ok := item["trigger"].(string); ok {
        data.Trigger = types.StringValue(val)
    } else {
        data.Trigger = types.StringNull()
    }
    if obj, ok := item["triggeredByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TriggeredByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TriggeredByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TriggeredByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TriggeredByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.TriggeredByUserId = types.StringNull()
        }
    } else if arr, ok := item["triggeredByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.TriggeredByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.TriggeredByUserId = types.StringNull()
        }
    } else if val, ok := item["triggeredByUserId"].(string); ok {
        data.TriggeredByUserId = types.StringValue(val)
    } else {
        data.TriggeredByUserId = types.StringNull()
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
    if obj, ok := item["steps"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Steps = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Steps = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Steps = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Steps = types.StringValue(string(jsonBytes))
        } else {
            data.Steps = types.StringNull()
        }
    } else if arr, ok := item["steps"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Steps = types.StringValue(string(jsonBytes))
        } else {
            data.Steps = types.StringNull()
        }
    } else if val, ok := item["steps"].(string); ok {
        data.Steps = types.StringValue(val)
    } else {
        data.Steps = types.StringNull()
    }
    if obj, ok := item["stepResults"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StepResults = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StepResults = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StepResults = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StepResults = types.StringValue(string(jsonBytes))
        } else {
            data.StepResults = types.StringNull()
        }
    } else if arr, ok := item["stepResults"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.StepResults = types.StringValue(string(jsonBytes))
        } else {
            data.StepResults = types.StringNull()
        }
    } else if val, ok := item["stepResults"].(string); ok {
        data.StepResults = types.StringValue(val)
    } else {
        data.StepResults = types.StringNull()
    }
    data.CurrentStepIndex = NumberFromAPI(data.CurrentStepIndex, item["currentStepIndex"])
    if val, ok := item["threatIntelMatched"].(bool); ok {
        data.ThreatIntelMatched = types.BoolValue(val)
    } else {
        data.ThreatIntelMatched = types.BoolNull()
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
    if obj, ok := item["completedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CompletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CompletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CompletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CompletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CompletedAt = types.StringNull()
        }
    } else if arr, ok := item["completedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CompletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CompletedAt = types.StringNull()
        }
    } else if val, ok := item["completedAt"].(string); ok {
        data.CompletedAt = types.StringValue(val)
    } else {
        data.CompletedAt = types.StringNull()
    }
    if obj, ok := item["errorMessage"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ErrorMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ErrorMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ErrorMessage = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ErrorMessage = types.StringValue(string(jsonBytes))
        } else {
            data.ErrorMessage = types.StringNull()
        }
    } else if arr, ok := item["errorMessage"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ErrorMessage = types.StringValue(string(jsonBytes))
        } else {
            data.ErrorMessage = types.StringNull()
        }
    } else if val, ok := item["errorMessage"].(string); ok {
        data.ErrorMessage = types.StringValue(val)
    } else {
        data.ErrorMessage = types.StringNull()
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
