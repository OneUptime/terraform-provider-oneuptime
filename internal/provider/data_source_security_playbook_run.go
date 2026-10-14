package provider

import (
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
    ProjectId types.String `tfsdk:"project_id"`
    SecurityCaseId types.String `tfsdk:"security_case_id"`
    SecurityPlaybookId types.String `tfsdk:"security_playbook_id"`
    PlaybookName types.String `tfsdk:"playbook_name"`
    Trigger types.String `tfsdk:"trigger"`
    TriggeredByUserId types.String `tfsdk:"triggered_by_user_id"`
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
        MarkdownDescription: "One run of a security playbook on a case: the steps it took, what each did, and who approved the steps that asked. Look up an existing security playbook run by `id`, or by any of its other arguments (`created_by_user_id`, `current_step_index`, `error_message`, ...): each one set must match, and exactly one security playbook run may match them all.",

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
                MarkdownDescription: "ID of the project this run belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "security_case_id": schema.StringAttribute{
                MarkdownDescription: "ID of the case the playbook runs on. The ID of a `oneuptime_security_case`.",
                Optional: true,
                Computed: true,
            },
            "security_playbook_id": schema.StringAttribute{
                MarkdownDescription: "ID of the playbook being run. The ID of a `oneuptime_security_playbook`.",
                Optional: true,
                Computed: true,
            },
            "playbook_name": schema.StringAttribute{
                MarkdownDescription: "The playbook's name when the run started, kept if the playbook is renamed or deleted.",
                Computed: true,
            },
            "trigger": schema.StringAttribute{
                MarkdownDescription: "What started the run: Case Opened, Detection Added, or Manual for a person.",
                Optional: true,
                Computed: true,
            },
            "triggered_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who ran the playbook. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Queued, Running, Waiting for Approval, Succeeded, Failed or Cancelled. Set it to Cancelled to stop a run that has not finished.",
                Optional: true,
                Computed: true,
            },
            "steps": schema.StringAttribute{
                MarkdownDescription: "The playbook's steps as they were when the run started. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "step_results": schema.StringAttribute{
                MarkdownDescription: "What each step did: its status, when it started and finished, what it reported, and who decided on its approval. A JSON value: write it with `jsonencode()`.",
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
    if !data.SecurityCaseId.IsNull() && !data.SecurityCaseId.IsUnknown() {
        filters["securityCaseId"] = data.SecurityCaseId.ValueString()
        filterNames = append(filterNames, "security_case_id = "+fmt.Sprintf("%q", data.SecurityCaseId.ValueString()))
    }
    if !data.SecurityPlaybookId.IsNull() && !data.SecurityPlaybookId.IsUnknown() {
        filters["securityPlaybookId"] = data.SecurityPlaybookId.ValueString()
        filterNames = append(filterNames, "security_playbook_id = "+fmt.Sprintf("%q", data.SecurityPlaybookId.ValueString()))
    }
    if !data.Trigger.IsNull() && !data.Trigger.IsUnknown() {
        filters["trigger"] = data.Trigger.ValueString()
        filterNames = append(filterNames, "trigger = "+fmt.Sprintf("%q", data.Trigger.ValueString()))
    }
    if !data.TriggeredByUserId.IsNull() && !data.TriggeredByUserId.IsUnknown() {
        filters["triggeredByUserId"] = data.TriggeredByUserId.ValueString()
        filterNames = append(filterNames, "triggered_by_user_id = "+fmt.Sprintf("%q", data.TriggeredByUserId.ValueString()))
    }
    if !data.Status.IsNull() && !data.Status.IsUnknown() {
        filters["status"] = data.Status.ValueString()
        filterNames = append(filterNames, "status = "+fmt.Sprintf("%q", data.Status.ValueString()))
    }
    if !data.CurrentStepIndex.IsNull() && !data.CurrentStepIndex.IsUnknown() {
        filters["currentStepIndex"] = lookupNumber(data.CurrentStepIndex)
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
        "projectId": true,
        "securityCaseId": true,
        "securityPlaybookId": true,
        "playbookName": true,
        "trigger": true,
        "triggeredByUserId": true,
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No security playbook run matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one security playbook run matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
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
    } else if val, ok := item["triggeredByUserId"].(string); ok {
        data.TriggeredByUserId = types.StringValue(val)
    } else {
        data.TriggeredByUserId = types.StringNull()
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
    } else if val, ok := item["stepResults"].(string); ok {
        data.StepResults = types.StringValue(val)
    } else {
        data.StepResults = types.StringNull()
    }
    data.CurrentStepIndex = numberFromAPI(data.CurrentStepIndex, item["currentStepIndex"])
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
