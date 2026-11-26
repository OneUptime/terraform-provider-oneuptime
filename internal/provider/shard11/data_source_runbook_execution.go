package shard11

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
var _ datasource.DataSource = &RunbookExecutionDataSource{}

func NewRunbookExecutionDataSource() datasource.DataSource {
    return &RunbookExecutionDataSource{}
}

// RunbookExecutionDataSource defines the data source implementation.
type RunbookExecutionDataSource struct {
    client *Client
}

// RunbookExecutionDataSourceModel describes the data source data model.
type RunbookExecutionDataSourceModel struct {
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
    RunbookId types.String `tfsdk:"runbook_id"`
    RunbookNameSnapshot types.String `tfsdk:"runbook_name_snapshot"`
    Status types.String `tfsdk:"status"`
    StepExecutions types.String `tfsdk:"step_executions"`
    IncidentId types.String `tfsdk:"incident_id"`
    AlertId types.String `tfsdk:"alert_id"`
    ScheduledMaintenanceId types.String `tfsdk:"scheduled_maintenance_id"`
    TriggeredByUserId types.String `tfsdk:"triggered_by_user_id"`
    TriggerType types.String `tfsdk:"trigger_type"`
    TriggerDetails types.String `tfsdk:"trigger_details"`
    Inputs types.String `tfsdk:"inputs"`
    StartedAt types.String `tfsdk:"started_at"`
    CompletedAt types.String `tfsdk:"completed_at"`
    FailureReason types.String `tfsdk:"failure_reason"`
}

func (d *RunbookExecutionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_runbook_execution"
}

func (d *RunbookExecutionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "A single run of a Runbook. Look up an existing runbook execution by `id`, or by any of its other arguments (`alert_id`, `created_by_actor_id`, `created_by_actor_type`, ...): each one set must match, and exactly one runbook execution may match them all.",

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
            "runbook_id": schema.StringAttribute{
                MarkdownDescription: "ID of the Runbook this execution belongs to. The ID of a `oneuptime_runbook`.",
                Optional: true,
                Computed: true,
                Validators: RecordIDValidators("oneuptime_runbook"),
            },
            "runbook_name_snapshot": schema.StringAttribute{
                MarkdownDescription: "Name of the runbook at the time this execution was created (preserved even if the runbook is later renamed or deleted).",
                Optional: true,
                Computed: true,
            },
            "status": schema.StringAttribute{
                MarkdownDescription: "Current status of this runbook execution.",
                Optional: true,
                Computed: true,
            },
            "step_executions": schema.StringAttribute{
                MarkdownDescription: "Per-step execution state. Each entry mirrors a step from the runbook with status, output, and timestamps. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "incident_id": schema.StringAttribute{
                MarkdownDescription: "ID of the incident that triggered or hosts this runbook execution. The ID of a `oneuptime_incident`.",
                Optional: true,
                Computed: true,
                Validators: RecordIDValidators("oneuptime_incident"),
            },
            "alert_id": schema.StringAttribute{
                MarkdownDescription: "ID of the alert that triggered or hosts this runbook execution. The ID of a `oneuptime_alert`.",
                Optional: true,
                Computed: true,
                Validators: RecordIDValidators("oneuptime_alert"),
            },
            "scheduled_maintenance_id": schema.StringAttribute{
                MarkdownDescription: "ID of the maintenance event that triggered this runbook execution. The ID of a `oneuptime_scheduled_maintenance_event`.",
                Optional: true,
                Computed: true,
                Validators: RecordIDValidators("oneuptime_scheduled_maintenance_event"),
            },
            "triggered_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the User who triggered this runbook execution. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
            "trigger_type": schema.StringAttribute{
                MarkdownDescription: "What started this run: \"Manual\" (a person, or the API), \"Rule\" (a runbook rule), \"AutoRemediation\", \"AI\" (OneUptime AI), \"Schedule\" (the runbook's schedule) or \"Webhook\" (a call to the runbook's webhook URL). Empty for runs started before it was recorded.",
                Optional: true,
                Computed: true,
            },
            "trigger_details": schema.StringAttribute{
                MarkdownDescription: "More about what started this run: for a scheduled run, the time it was due (scheduledFor) and whether it made up for times OneUptime missed (isCatchUp) or waited for the run before it (waitedForPreviousRun); for a webhook call, the address it came from (callerIp). A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "inputs": schema.StringAttribute{
                MarkdownDescription: "The values this run was started with: the JSON body of the webhook call that started it. Steps read them as data - Bash as ONEUPTIME_INPUT_<NAME> environment variables, JavaScript as args.inputs - never as script text. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "started_at": schema.StringAttribute{
                MarkdownDescription: "Time at which this runbook execution started.",
                Computed: true,
            },
            "completed_at": schema.StringAttribute{
                MarkdownDescription: "Time at which this runbook execution completed.",
                Computed: true,
            },
            "failure_reason": schema.StringAttribute{
                MarkdownDescription: "Reason this runbook execution failed (if it did).",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *RunbookExecutionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RunbookExecutionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data RunbookExecutionDataSourceModel

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
    if !data.RunbookId.IsNull() && !data.RunbookId.IsUnknown() {
        filters["runbookId"] = data.RunbookId.ValueString()
        filterNames = append(filterNames, "runbook_id = "+fmt.Sprintf("%q", data.RunbookId.ValueString()))
    }
    if !data.RunbookNameSnapshot.IsNull() && !data.RunbookNameSnapshot.IsUnknown() {
        filters["runbookNameSnapshot"] = data.RunbookNameSnapshot.ValueString()
        filterNames = append(filterNames, "runbook_name_snapshot = "+fmt.Sprintf("%q", data.RunbookNameSnapshot.ValueString()))
    }
    if !data.Status.IsNull() && !data.Status.IsUnknown() {
        filters["status"] = data.Status.ValueString()
        filterNames = append(filterNames, "status = "+fmt.Sprintf("%q", data.Status.ValueString()))
    }
    if !data.IncidentId.IsNull() && !data.IncidentId.IsUnknown() {
        filters["incidentId"] = data.IncidentId.ValueString()
        filterNames = append(filterNames, "incident_id = "+fmt.Sprintf("%q", data.IncidentId.ValueString()))
    }
    if !data.AlertId.IsNull() && !data.AlertId.IsUnknown() {
        filters["alertId"] = data.AlertId.ValueString()
        filterNames = append(filterNames, "alert_id = "+fmt.Sprintf("%q", data.AlertId.ValueString()))
    }
    if !data.ScheduledMaintenanceId.IsNull() && !data.ScheduledMaintenanceId.IsUnknown() {
        filters["scheduledMaintenanceId"] = data.ScheduledMaintenanceId.ValueString()
        filterNames = append(filterNames, "scheduled_maintenance_id = "+fmt.Sprintf("%q", data.ScheduledMaintenanceId.ValueString()))
    }
    if !data.TriggeredByUserId.IsNull() && !data.TriggeredByUserId.IsUnknown() {
        filters["triggeredByUserId"] = data.TriggeredByUserId.ValueString()
        filterNames = append(filterNames, "triggered_by_user_id = "+fmt.Sprintf("%q", data.TriggeredByUserId.ValueString()))
    }
    if !data.TriggerType.IsNull() && !data.TriggerType.IsUnknown() {
        filters["triggerType"] = data.TriggerType.ValueString()
        filterNames = append(filterNames, "trigger_type = "+fmt.Sprintf("%q", data.TriggerType.ValueString()))
    }
    if !data.FailureReason.IsNull() && !data.FailureReason.IsUnknown() {
        filters["failureReason"] = data.FailureReason.ValueString()
        filterNames = append(filterNames, "failure_reason = "+fmt.Sprintf("%q", data.FailureReason.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the runbook execution up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the runbook execution up by.",
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
        "runbookId": true,
        "runbookNameSnapshot": true,
        "status": true,
        "stepExecutions": true,
        "incidentId": true,
        "alertId": true,
        "scheduledMaintenanceId": true,
        "triggeredByUserId": true,
        "triggerType": true,
        "triggerDetails": true,
        "inputs": true,
        "startedAt": true,
        "completedAt": true,
        "failureReason": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/runbook-execution/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read runbook_execution, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No runbook execution found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read runbook_execution: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/runbook-execution/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list runbook_execution, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list runbook_execution: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No runbook execution matches %s.", DescribeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one runbook execution matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", DescribeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for runbook_execution.")
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
    if obj, ok := item["runbookId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RunbookId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RunbookId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RunbookId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RunbookId = types.StringValue(string(jsonBytes))
        } else {
            data.RunbookId = types.StringNull()
        }
    } else if arr, ok := item["runbookId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.RunbookId = types.StringValue(string(jsonBytes))
        } else {
            data.RunbookId = types.StringNull()
        }
    } else if val, ok := item["runbookId"].(string); ok {
        data.RunbookId = types.StringValue(val)
    } else {
        data.RunbookId = types.StringNull()
    }
    if obj, ok := item["runbookNameSnapshot"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RunbookNameSnapshot = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RunbookNameSnapshot = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RunbookNameSnapshot = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RunbookNameSnapshot = types.StringValue(string(jsonBytes))
        } else {
            data.RunbookNameSnapshot = types.StringNull()
        }
    } else if arr, ok := item["runbookNameSnapshot"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.RunbookNameSnapshot = types.StringValue(string(jsonBytes))
        } else {
            data.RunbookNameSnapshot = types.StringNull()
        }
    } else if val, ok := item["runbookNameSnapshot"].(string); ok {
        data.RunbookNameSnapshot = types.StringValue(val)
    } else {
        data.RunbookNameSnapshot = types.StringNull()
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
    if obj, ok := item["stepExecutions"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StepExecutions = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StepExecutions = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StepExecutions = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StepExecutions = types.StringValue(string(jsonBytes))
        } else {
            data.StepExecutions = types.StringNull()
        }
    } else if arr, ok := item["stepExecutions"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.StepExecutions = types.StringValue(string(jsonBytes))
        } else {
            data.StepExecutions = types.StringNull()
        }
    } else if val, ok := item["stepExecutions"].(string); ok {
        data.StepExecutions = types.StringValue(val)
    } else {
        data.StepExecutions = types.StringNull()
    }
    if obj, ok := item["incidentId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncidentId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncidentId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncidentId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentId = types.StringNull()
        }
    } else if arr, ok := item["incidentId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.IncidentId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentId = types.StringNull()
        }
    } else if val, ok := item["incidentId"].(string); ok {
        data.IncidentId = types.StringValue(val)
    } else {
        data.IncidentId = types.StringNull()
    }
    if obj, ok := item["alertId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AlertId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AlertId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AlertId = types.StringValue(string(jsonBytes))
        } else {
            data.AlertId = types.StringNull()
        }
    } else if arr, ok := item["alertId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AlertId = types.StringValue(string(jsonBytes))
        } else {
            data.AlertId = types.StringNull()
        }
    } else if val, ok := item["alertId"].(string); ok {
        data.AlertId = types.StringValue(val)
    } else {
        data.AlertId = types.StringNull()
    }
    if obj, ok := item["scheduledMaintenanceId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScheduledMaintenanceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ScheduledMaintenanceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ScheduledMaintenanceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ScheduledMaintenanceId = types.StringValue(string(jsonBytes))
        } else {
            data.ScheduledMaintenanceId = types.StringNull()
        }
    } else if arr, ok := item["scheduledMaintenanceId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ScheduledMaintenanceId = types.StringValue(string(jsonBytes))
        } else {
            data.ScheduledMaintenanceId = types.StringNull()
        }
    } else if val, ok := item["scheduledMaintenanceId"].(string); ok {
        data.ScheduledMaintenanceId = types.StringValue(val)
    } else {
        data.ScheduledMaintenanceId = types.StringNull()
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
    if obj, ok := item["triggerType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TriggerType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TriggerType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TriggerType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TriggerType = types.StringValue(string(jsonBytes))
        } else {
            data.TriggerType = types.StringNull()
        }
    } else if arr, ok := item["triggerType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.TriggerType = types.StringValue(string(jsonBytes))
        } else {
            data.TriggerType = types.StringNull()
        }
    } else if val, ok := item["triggerType"].(string); ok {
        data.TriggerType = types.StringValue(val)
    } else {
        data.TriggerType = types.StringNull()
    }
    if obj, ok := item["triggerDetails"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TriggerDetails = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TriggerDetails = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TriggerDetails = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TriggerDetails = types.StringValue(string(jsonBytes))
        } else {
            data.TriggerDetails = types.StringNull()
        }
    } else if arr, ok := item["triggerDetails"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.TriggerDetails = types.StringValue(string(jsonBytes))
        } else {
            data.TriggerDetails = types.StringNull()
        }
    } else if val, ok := item["triggerDetails"].(string); ok {
        data.TriggerDetails = types.StringValue(val)
    } else {
        data.TriggerDetails = types.StringNull()
    }
    if obj, ok := item["inputs"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Inputs = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Inputs = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Inputs = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Inputs = types.StringValue(string(jsonBytes))
        } else {
            data.Inputs = types.StringNull()
        }
    } else if arr, ok := item["inputs"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Inputs = types.StringValue(string(jsonBytes))
        } else {
            data.Inputs = types.StringNull()
        }
    } else if val, ok := item["inputs"].(string); ok {
        data.Inputs = types.StringValue(val)
    } else {
        data.Inputs = types.StringNull()
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
    if obj, ok := item["failureReason"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FailureReason = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FailureReason = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FailureReason = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FailureReason = types.StringValue(string(jsonBytes))
        } else {
            data.FailureReason = types.StringNull()
        }
    } else if arr, ok := item["failureReason"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.FailureReason = types.StringValue(string(jsonBytes))
        } else {
            data.FailureReason = types.StringNull()
        }
    } else if val, ok := item["failureReason"].(string); ok {
        data.FailureReason = types.StringValue(val)
    } else {
        data.FailureReason = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
