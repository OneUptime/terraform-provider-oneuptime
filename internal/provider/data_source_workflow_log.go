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
}

func (d *WorkflowLogDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_workflow_log"
}

func (d *WorkflowLogDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Logs of the workflows executed Look up an existing workflow log by `id`, or by any of its other arguments (`logs`, `workflow_id`, `workflow_status`): each one set must match, and exactly one workflow log may match them all.",

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
                MarkdownDescription: "Structured per-step record of this run: arguments, return values, the port taken and timing. A JSON value: write it with `jsonencode()`.",
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
    } else if val, ok := item["stepTrace"].(string); ok {
        data.StepTrace = types.StringValue(val)
    } else {
        data.StepTrace = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
