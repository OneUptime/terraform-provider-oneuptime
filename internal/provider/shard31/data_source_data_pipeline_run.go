package shard31

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
var _ datasource.DataSource = &DataPipelineRunDataSource{}

func NewDataPipelineRunDataSource() datasource.DataSource {
    return &DataPipelineRunDataSource{}
}

// DataPipelineRunDataSource defines the data source implementation.
type DataPipelineRunDataSource struct {
    client *Client
}

// DataPipelineRunDataSourceModel describes the data source data model.
type DataPipelineRunDataSourceModel struct {
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
    DataPipelineId types.String `tfsdk:"data_pipeline_id"`
    RunId types.String `tfsdk:"run_id"`
    State types.String `tfsdk:"state"`
    StartedAt types.String `tfsdk:"started_at"`
    ScheduledAt types.String `tfsdk:"scheduled_at"`
    EndedAt types.String `tfsdk:"ended_at"`
    DurationInSeconds types.Number `tfsdk:"duration_in_seconds"`
    ErrorMessage types.String `tfsdk:"error_message"`
    Inputs types.String `tfsdk:"inputs"`
    Outputs types.String `tfsdk:"outputs"`
    OutputDataAssetIds types.String `tfsdk:"output_data_asset_ids"`
    AttemptNumber types.Number `tfsdk:"attempt_number"`
    TraceId types.String `tfsdk:"trace_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *DataPipelineRunDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_data_pipeline_run"
}

func (d *DataPipelineRunDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "One run of a data pipeline: when it started and ended, how it went, and the datasets it read and wrote. Look up an existing data pipeline run by `id`, or by any of its other arguments (`attempt_number`, `created_by_actor_id`, `created_by_actor_type`, ...): each one set must match, and exactly one data pipeline run may match them all.",

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
            "data_pipeline_id": schema.StringAttribute{
                MarkdownDescription: "ID of the pipeline this is a run of. The ID of a `oneuptime_data_pipeline` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_data_pipeline"),
            },
            "run_id": schema.StringAttribute{
                MarkdownDescription: "The run's OpenLineage id.",
                Optional: true,
                Computed: true,
            },
            "state": schema.StringAttribute{
                MarkdownDescription: "How the run stands: Running, Complete, Fail or Abort.",
                Optional: true,
                Computed: true,
            },
            "started_at": schema.StringAttribute{
                MarkdownDescription: "When the run started.",
                Computed: true,
            },
            "scheduled_at": schema.StringAttribute{
                MarkdownDescription: "When the run was scheduled for, when its producer says (OpenLineage's nominalTime facet): a backfilled run's scheduled time, not when it ran.",
                Computed: true,
            },
            "ended_at": schema.StringAttribute{
                MarkdownDescription: "When the run ended.",
                Computed: true,
            },
            "duration_in_seconds": schema.NumberAttribute{
                MarkdownDescription: "How long the run took, in seconds.",
                Optional: true,
                Computed: true,
            },
            "error_message": schema.StringAttribute{
                MarkdownDescription: "The error a failed run reported.",
                Optional: true,
                Computed: true,
            },
            "inputs": schema.StringAttribute{
                MarkdownDescription: "The datasets the run read, as namespace and name. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "outputs": schema.StringAttribute{
                MarkdownDescription: "The datasets the run wrote, as namespace and name. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "output_data_asset_ids": schema.StringAttribute{
                MarkdownDescription: "The watched tables among the datasets the run wrote. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "attempt_number": schema.NumberAttribute{
                MarkdownDescription: "Which try of its task this run is, when its scheduler says (Airflow's try number).",
                Optional: true,
                Computed: true,
            },
            "trace_id": schema.StringAttribute{
                MarkdownDescription: "The trace of a run found in traces: a CronJob or a scheduled function.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
        },
    }
}

func (d *DataPipelineRunDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataPipelineRunDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data DataPipelineRunDataSourceModel

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
    if !data.DataPipelineId.IsNull() && !data.DataPipelineId.IsUnknown() {
        filters["dataPipelineId"] = data.DataPipelineId.ValueString()
        filterNames = append(filterNames, "data_pipeline_id = "+fmt.Sprintf("%q", data.DataPipelineId.ValueString()))
    }
    if !data.RunId.IsNull() && !data.RunId.IsUnknown() {
        filters["runId"] = data.RunId.ValueString()
        filterNames = append(filterNames, "run_id = "+fmt.Sprintf("%q", data.RunId.ValueString()))
    }
    if !data.State.IsNull() && !data.State.IsUnknown() {
        filters["state"] = data.State.ValueString()
        filterNames = append(filterNames, "state = "+fmt.Sprintf("%q", data.State.ValueString()))
    }
    if !data.DurationInSeconds.IsNull() && !data.DurationInSeconds.IsUnknown() {
        filters["durationInSeconds"] = LookupNumber(data.DurationInSeconds)
        filterNames = append(filterNames, "duration_in_seconds = "+data.DurationInSeconds.ValueBigFloat().String())
    }
    if !data.ErrorMessage.IsNull() && !data.ErrorMessage.IsUnknown() {
        filters["errorMessage"] = data.ErrorMessage.ValueString()
        filterNames = append(filterNames, "error_message = "+fmt.Sprintf("%q", data.ErrorMessage.ValueString()))
    }
    if !data.AttemptNumber.IsNull() && !data.AttemptNumber.IsUnknown() {
        filters["attemptNumber"] = LookupNumber(data.AttemptNumber)
        filterNames = append(filterNames, "attempt_number = "+data.AttemptNumber.ValueBigFloat().String())
    }
    if !data.TraceId.IsNull() && !data.TraceId.IsUnknown() {
        filters["traceId"] = data.TraceId.ValueString()
        filterNames = append(filterNames, "trace_id = "+fmt.Sprintf("%q", data.TraceId.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the data pipeline run up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the data pipeline run up by.",
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
        "dataPipelineId": true,
        "runId": true,
        "state": true,
        "startedAt": true,
        "scheduledAt": true,
        "endedAt": true,
        "durationInSeconds": true,
        "errorMessage": true,
        "inputs": true,
        "outputs": true,
        "outputDataAssetIds": true,
        "attemptNumber": true,
        "traceId": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/data-pipeline-run/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read data_pipeline_run, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data pipeline run found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read data_pipeline_run: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/data-pipeline-run/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list data_pipeline_run, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list data_pipeline_run: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data pipeline run matches %s.", DescribeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one data pipeline run matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", DescribeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for data_pipeline_run.")
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
    if obj, ok := item["dataPipelineId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataPipelineId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DataPipelineId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DataPipelineId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DataPipelineId = types.StringValue(string(jsonBytes))
        } else {
            data.DataPipelineId = types.StringNull()
        }
    } else if arr, ok := item["dataPipelineId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.DataPipelineId = types.StringValue(string(jsonBytes))
        } else {
            data.DataPipelineId = types.StringNull()
        }
    } else if val, ok := item["dataPipelineId"].(string); ok {
        data.DataPipelineId = types.StringValue(val)
    } else {
        data.DataPipelineId = types.StringNull()
    }
    if obj, ok := item["runId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RunId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RunId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RunId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RunId = types.StringValue(string(jsonBytes))
        } else {
            data.RunId = types.StringNull()
        }
    } else if arr, ok := item["runId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.RunId = types.StringValue(string(jsonBytes))
        } else {
            data.RunId = types.StringNull()
        }
    } else if val, ok := item["runId"].(string); ok {
        data.RunId = types.StringValue(val)
    } else {
        data.RunId = types.StringNull()
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
    if obj, ok := item["scheduledAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScheduledAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ScheduledAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ScheduledAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ScheduledAt = types.StringValue(string(jsonBytes))
        } else {
            data.ScheduledAt = types.StringNull()
        }
    } else if arr, ok := item["scheduledAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ScheduledAt = types.StringValue(string(jsonBytes))
        } else {
            data.ScheduledAt = types.StringNull()
        }
    } else if val, ok := item["scheduledAt"].(string); ok {
        data.ScheduledAt = types.StringValue(val)
    } else {
        data.ScheduledAt = types.StringNull()
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
    data.DurationInSeconds = NumberFromAPI(data.DurationInSeconds, item["durationInSeconds"])
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
    if obj, ok := item["outputs"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Outputs = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Outputs = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Outputs = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Outputs = types.StringValue(string(jsonBytes))
        } else {
            data.Outputs = types.StringNull()
        }
    } else if arr, ok := item["outputs"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Outputs = types.StringValue(string(jsonBytes))
        } else {
            data.Outputs = types.StringNull()
        }
    } else if val, ok := item["outputs"].(string); ok {
        data.Outputs = types.StringValue(val)
    } else {
        data.Outputs = types.StringNull()
    }
    if obj, ok := item["outputDataAssetIds"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OutputDataAssetIds = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OutputDataAssetIds = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OutputDataAssetIds = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OutputDataAssetIds = types.StringValue(string(jsonBytes))
        } else {
            data.OutputDataAssetIds = types.StringNull()
        }
    } else if arr, ok := item["outputDataAssetIds"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.OutputDataAssetIds = types.StringValue(string(jsonBytes))
        } else {
            data.OutputDataAssetIds = types.StringNull()
        }
    } else if val, ok := item["outputDataAssetIds"].(string); ok {
        data.OutputDataAssetIds = types.StringValue(val)
    } else {
        data.OutputDataAssetIds = types.StringNull()
    }
    data.AttemptNumber = NumberFromAPI(data.AttemptNumber, item["attemptNumber"])
    if obj, ok := item["traceId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TraceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TraceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TraceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TraceId = types.StringValue(string(jsonBytes))
        } else {
            data.TraceId = types.StringNull()
        }
    } else if arr, ok := item["traceId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.TraceId = types.StringValue(string(jsonBytes))
        } else {
            data.TraceId = types.StringNull()
        }
    } else if val, ok := item["traceId"].(string); ok {
        data.TraceId = types.StringValue(val)
    } else {
        data.TraceId = types.StringNull()
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
