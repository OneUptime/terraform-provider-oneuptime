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
var _ datasource.DataSource = &WorkflowLogDataSource{}

func NewWorkflowLogDataSource() datasource.DataSource {
    return &WorkflowLogDataSource{}
}

// WorkflowLogDataSource defines the data source implementation.
type WorkflowLogDataSource struct {
    client *Client
}

// WorkflowLogDataSourceModel describes the data source data model.
type WorkflowLogDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    WorkflowId types.String `tfsdk:"workflow_id"`
    Logs types.String `tfsdk:"logs"`
    WorkflowStatus types.String `tfsdk:"workflow_status"`
    StartedAt types.String `tfsdk:"started_at"`
    CompletedAt types.String `tfsdk:"completed_at"`
    ResumeAt types.String `tfsdk:"resume_at"`
    StepTrace types.String `tfsdk:"step_trace"`
    IsTestRun types.Bool `tfsdk:"is_test_run"`
    CanRunAgain types.Bool `tfsdk:"can_run_again"`
    StartedByUserId types.String `tfsdk:"started_by_user_id"`
    RunAs types.String `tfsdk:"run_as"`
    RunAsCustomRoleName types.String `tfsdk:"run_as_custom_role_name"`
    RunAsGivenByUserId types.String `tfsdk:"run_as_given_by_user_id"`
    RunAsGivenByApiKeyId types.String `tfsdk:"run_as_given_by_api_key_id"`
    RunAsGivenByApiKeyName types.String `tfsdk:"run_as_given_by_api_key_name"`
}

func (d *WorkflowLogDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_workflow_log"
}

func (d *WorkflowLogDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Logs of the workflows executed Look up an existing workflow log by `id`, or by any of its other arguments (`can_run_again`, `is_test_run`, `logs`, ...): each one set must match, and exactly one workflow log may match them all.",

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
            "workflow_id": schema.StringAttribute{
                MarkdownDescription: "ID of Workflow this logs belong to. The ID of a `oneuptime_workflow`.",
                Optional: true,
                Computed: true,
                Validators: recordIDValidators("oneuptime_workflow"),
            },
            "logs": schema.StringAttribute{
                MarkdownDescription: "Logs.",
                Optional: true,
                Computed: true,
            },
            "workflow_status": schema.StringAttribute{
                MarkdownDescription: "Status of this workflow.",
                Optional: true,
                Computed: true,
            },
            "started_at": schema.StringAttribute{
                MarkdownDescription: "When did this workflow start.",
                Computed: true,
            },
            "completed_at": schema.StringAttribute{
                MarkdownDescription: "When did this workflow complete.",
                Computed: true,
            },
            "resume_at": schema.StringAttribute{
                MarkdownDescription: "When this workflow run is scheduled to resume after a Sleep step (only set while the run is Waiting).",
                Computed: true,
            },
            "step_trace": schema.StringAttribute{
                MarkdownDescription: "Structured per-step record of this run: arguments, return values, the port taken and timing. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "is_test_run": schema.BoolAttribute{
                MarkdownDescription: "True for a run a person started by hand from the Builder while the workflow was turned off. Its steps ran for real; the workflow's trigger stayed paused.",
                Optional: true,
                Computed: true,
            },
            "can_run_again": schema.BoolAttribute{
                MarkdownDescription: "True when OneUptime kept what this run started with, so a person can run the workflow again with the same values.",
                Optional: true,
                Computed: true,
            },
            "started_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who started this run by hand from the Builder. Empty for a run its trigger or another workflow started. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: dataSourceRecordIDValidators("oneuptime_user"),
            },
            "run_as": schema.StringAttribute{
                MarkdownDescription: "What the run's steps could do in the project: Default, ProjectOwner or CustomRole (runAsCustomRoleName names the role). Set by OneUptime as the run starts.",
                Optional: true,
                Computed: true,
            },
            "run_as_custom_role_name": schema.StringAttribute{
                MarkdownDescription: "The name of the custom role the run's steps acted as, as it was when the run started. Empty unless the run ran as a custom role.",
                Optional: true,
                Computed: true,
            },
            "run_as_given_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person whose permissions the run ran with, when it ran as more than Default. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: dataSourceRecordIDValidators("oneuptime_user"),
            },
            "run_as_given_by_api_key_id": schema.StringAttribute{
                MarkdownDescription: "ID of the project API key whose permissions the run ran with, when it ran as more than Default and a key gave them.",
                Optional: true,
                Computed: true,
            },
            "run_as_given_by_api_key_name": schema.StringAttribute{
                MarkdownDescription: "The name of the project API key whose permissions the run ran with, as it was when the run started. Empty unless a key gave them.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *WorkflowLogDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *WorkflowLogDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data WorkflowLogDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.WorkflowId.IsNull() && !data.WorkflowId.IsUnknown() {
        filters["workflowId"] = data.WorkflowId.ValueString()
        filterNames = append(filterNames, "workflow_id = "+fmt.Sprintf("%q", data.WorkflowId.ValueString()))
    }
    if !data.Logs.IsNull() && !data.Logs.IsUnknown() {
        filters["logs"] = data.Logs.ValueString()
        filterNames = append(filterNames, "logs = "+fmt.Sprintf("%q", data.Logs.ValueString()))
    }
    if !data.WorkflowStatus.IsNull() && !data.WorkflowStatus.IsUnknown() {
        filters["workflowStatus"] = data.WorkflowStatus.ValueString()
        filterNames = append(filterNames, "workflow_status = "+fmt.Sprintf("%q", data.WorkflowStatus.ValueString()))
    }
    if !data.IsTestRun.IsNull() && !data.IsTestRun.IsUnknown() {
        filters["isTestRun"] = data.IsTestRun.ValueBool()
        filterNames = append(filterNames, "is_test_run = "+fmt.Sprintf("%t", data.IsTestRun.ValueBool()))
    }
    if !data.CanRunAgain.IsNull() && !data.CanRunAgain.IsUnknown() {
        filters["canRunAgain"] = data.CanRunAgain.ValueBool()
        filterNames = append(filterNames, "can_run_again = "+fmt.Sprintf("%t", data.CanRunAgain.ValueBool()))
    }
    if !data.StartedByUserId.IsNull() && !data.StartedByUserId.IsUnknown() {
        filters["startedByUserId"] = data.StartedByUserId.ValueString()
        filterNames = append(filterNames, "started_by_user_id = "+fmt.Sprintf("%q", data.StartedByUserId.ValueString()))
    }
    if !data.RunAs.IsNull() && !data.RunAs.IsUnknown() {
        filters["runAs"] = data.RunAs.ValueString()
        filterNames = append(filterNames, "run_as = "+fmt.Sprintf("%q", data.RunAs.ValueString()))
    }
    if !data.RunAsCustomRoleName.IsNull() && !data.RunAsCustomRoleName.IsUnknown() {
        filters["runAsCustomRoleName"] = data.RunAsCustomRoleName.ValueString()
        filterNames = append(filterNames, "run_as_custom_role_name = "+fmt.Sprintf("%q", data.RunAsCustomRoleName.ValueString()))
    }
    if !data.RunAsGivenByUserId.IsNull() && !data.RunAsGivenByUserId.IsUnknown() {
        filters["runAsGivenByUserId"] = data.RunAsGivenByUserId.ValueString()
        filterNames = append(filterNames, "run_as_given_by_user_id = "+fmt.Sprintf("%q", data.RunAsGivenByUserId.ValueString()))
    }
    if !data.RunAsGivenByApiKeyId.IsNull() && !data.RunAsGivenByApiKeyId.IsUnknown() {
        filters["runAsGivenByApiKeyId"] = data.RunAsGivenByApiKeyId.ValueString()
        filterNames = append(filterNames, "run_as_given_by_api_key_id = "+fmt.Sprintf("%q", data.RunAsGivenByApiKeyId.ValueString()))
    }
    if !data.RunAsGivenByApiKeyName.IsNull() && !data.RunAsGivenByApiKeyName.IsUnknown() {
        filters["runAsGivenByApiKeyName"] = data.RunAsGivenByApiKeyName.ValueString()
        filterNames = append(filterNames, "run_as_given_by_api_key_name = "+fmt.Sprintf("%q", data.RunAsGivenByApiKeyName.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the workflow log up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the workflow log up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "workflowId": true,
        "logs": true,
        "workflowStatus": true,
        "startedAt": true,
        "completedAt": true,
        "resumeAt": true,
        "stepTrace": true,
        "isTestRun": true,
        "canRunAgain": true,
        "startedByUserId": true,
        "runAs": true,
        "runAsCustomRoleName": true,
        "runAsGivenByUserId": true,
        "runAsGivenByApiKeyId": true,
        "runAsGivenByApiKeyName": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/workflow-log/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read workflow_log, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No workflow log found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read workflow_log: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/workflow-log/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list workflow_log, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list workflow_log: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No workflow log matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one workflow log matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for workflow_log.")
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
    if obj, ok := item["workflowId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.WorkflowId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.WorkflowId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.WorkflowId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.WorkflowId = types.StringValue(string(jsonBytes))
        } else {
            data.WorkflowId = types.StringNull()
        }
    } else if arr, ok := item["workflowId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.WorkflowId = types.StringValue(string(jsonBytes))
        } else {
            data.WorkflowId = types.StringNull()
        }
    } else if val, ok := item["workflowId"].(string); ok {
        data.WorkflowId = types.StringValue(val)
    } else {
        data.WorkflowId = types.StringNull()
    }
    if obj, ok := item["logs"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Logs = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Logs = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Logs = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Logs = types.StringValue(string(jsonBytes))
        } else {
            data.Logs = types.StringNull()
        }
    } else if arr, ok := item["logs"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Logs = types.StringValue(string(jsonBytes))
        } else {
            data.Logs = types.StringNull()
        }
    } else if val, ok := item["logs"].(string); ok {
        data.Logs = types.StringValue(val)
    } else {
        data.Logs = types.StringNull()
    }
    if obj, ok := item["workflowStatus"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.WorkflowStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.WorkflowStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.WorkflowStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.WorkflowStatus = types.StringValue(string(jsonBytes))
        } else {
            data.WorkflowStatus = types.StringNull()
        }
    } else if arr, ok := item["workflowStatus"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.WorkflowStatus = types.StringValue(string(jsonBytes))
        } else {
            data.WorkflowStatus = types.StringNull()
        }
    } else if val, ok := item["workflowStatus"].(string); ok {
        data.WorkflowStatus = types.StringValue(val)
    } else {
        data.WorkflowStatus = types.StringNull()
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
    if obj, ok := item["resumeAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResumeAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ResumeAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ResumeAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ResumeAt = types.StringValue(string(jsonBytes))
        } else {
            data.ResumeAt = types.StringNull()
        }
    } else if arr, ok := item["resumeAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ResumeAt = types.StringValue(string(jsonBytes))
        } else {
            data.ResumeAt = types.StringNull()
        }
    } else if val, ok := item["resumeAt"].(string); ok {
        data.ResumeAt = types.StringValue(val)
    } else {
        data.ResumeAt = types.StringNull()
    }
    if obj, ok := item["stepTrace"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StepTrace = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StepTrace = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StepTrace = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StepTrace = types.StringValue(string(jsonBytes))
        } else {
            data.StepTrace = types.StringNull()
        }
    } else if arr, ok := item["stepTrace"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.StepTrace = types.StringValue(string(jsonBytes))
        } else {
            data.StepTrace = types.StringNull()
        }
    } else if val, ok := item["stepTrace"].(string); ok {
        data.StepTrace = types.StringValue(val)
    } else {
        data.StepTrace = types.StringNull()
    }
    if val, ok := item["isTestRun"].(bool); ok {
        data.IsTestRun = types.BoolValue(val)
    } else {
        data.IsTestRun = types.BoolNull()
    }
    if val, ok := item["canRunAgain"].(bool); ok {
        data.CanRunAgain = types.BoolValue(val)
    } else {
        data.CanRunAgain = types.BoolNull()
    }
    if obj, ok := item["startedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StartedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StartedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StartedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.StartedByUserId = types.StringNull()
        }
    } else if arr, ok := item["startedByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.StartedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.StartedByUserId = types.StringNull()
        }
    } else if val, ok := item["startedByUserId"].(string); ok {
        data.StartedByUserId = types.StringValue(val)
    } else {
        data.StartedByUserId = types.StringNull()
    }
    if obj, ok := item["runAs"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RunAs = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RunAs = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RunAs = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RunAs = types.StringValue(string(jsonBytes))
        } else {
            data.RunAs = types.StringNull()
        }
    } else if arr, ok := item["runAs"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.RunAs = types.StringValue(string(jsonBytes))
        } else {
            data.RunAs = types.StringNull()
        }
    } else if val, ok := item["runAs"].(string); ok {
        data.RunAs = types.StringValue(val)
    } else {
        data.RunAs = types.StringNull()
    }
    if obj, ok := item["runAsCustomRoleName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RunAsCustomRoleName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RunAsCustomRoleName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RunAsCustomRoleName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RunAsCustomRoleName = types.StringValue(string(jsonBytes))
        } else {
            data.RunAsCustomRoleName = types.StringNull()
        }
    } else if arr, ok := item["runAsCustomRoleName"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.RunAsCustomRoleName = types.StringValue(string(jsonBytes))
        } else {
            data.RunAsCustomRoleName = types.StringNull()
        }
    } else if val, ok := item["runAsCustomRoleName"].(string); ok {
        data.RunAsCustomRoleName = types.StringValue(val)
    } else {
        data.RunAsCustomRoleName = types.StringNull()
    }
    if obj, ok := item["runAsGivenByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RunAsGivenByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RunAsGivenByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RunAsGivenByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RunAsGivenByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.RunAsGivenByUserId = types.StringNull()
        }
    } else if arr, ok := item["runAsGivenByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.RunAsGivenByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.RunAsGivenByUserId = types.StringNull()
        }
    } else if val, ok := item["runAsGivenByUserId"].(string); ok {
        data.RunAsGivenByUserId = types.StringValue(val)
    } else {
        data.RunAsGivenByUserId = types.StringNull()
    }
    if obj, ok := item["runAsGivenByApiKeyId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RunAsGivenByApiKeyId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RunAsGivenByApiKeyId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RunAsGivenByApiKeyId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RunAsGivenByApiKeyId = types.StringValue(string(jsonBytes))
        } else {
            data.RunAsGivenByApiKeyId = types.StringNull()
        }
    } else if arr, ok := item["runAsGivenByApiKeyId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.RunAsGivenByApiKeyId = types.StringValue(string(jsonBytes))
        } else {
            data.RunAsGivenByApiKeyId = types.StringNull()
        }
    } else if val, ok := item["runAsGivenByApiKeyId"].(string); ok {
        data.RunAsGivenByApiKeyId = types.StringValue(val)
    } else {
        data.RunAsGivenByApiKeyId = types.StringNull()
    }
    if obj, ok := item["runAsGivenByApiKeyName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RunAsGivenByApiKeyName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RunAsGivenByApiKeyName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RunAsGivenByApiKeyName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RunAsGivenByApiKeyName = types.StringValue(string(jsonBytes))
        } else {
            data.RunAsGivenByApiKeyName = types.StringNull()
        }
    } else if arr, ok := item["runAsGivenByApiKeyName"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.RunAsGivenByApiKeyName = types.StringValue(string(jsonBytes))
        } else {
            data.RunAsGivenByApiKeyName = types.StringNull()
        }
    } else if val, ok := item["runAsGivenByApiKeyName"].(string); ok {
        data.RunAsGivenByApiKeyName = types.StringValue(val)
    } else {
        data.RunAsGivenByApiKeyName = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
