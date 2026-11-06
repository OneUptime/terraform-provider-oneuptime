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
var _ datasource.DataSource = &DataPipelineDataSource{}

func NewDataPipelineDataSource() datasource.DataSource {
    return &DataPipelineDataSource{}
}

// DataPipelineDataSource defines the data source implementation.
type DataPipelineDataSource struct {
    client *Client
}

// DataPipelineDataSourceModel describes the data source data model.
type DataPipelineDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    JobNamespace types.String `tfsdk:"job_namespace"`
    JobName types.String `tfsdk:"job_name"`
    PipelineKey types.String `tfsdk:"pipeline_key"`
    Description types.String `tfsdk:"description"`
    LastRunId types.String `tfsdk:"last_run_id"`
    LastRunState types.String `tfsdk:"last_run_state"`
    LastRunStartedAt types.String `tfsdk:"last_run_started_at"`
    LastRunEndedAt types.String `tfsdk:"last_run_ended_at"`
    LastRunDurationInSeconds types.Number `tfsdk:"last_run_duration_in_seconds"`
    LastSuccessAt types.String `tfsdk:"last_success_at"`
    LastFailureAt types.String `tfsdk:"last_failure_at"`
    LastErrorMessage types.String `tfsdk:"last_error_message"`
    RunCount types.Number `tfsdk:"run_count"`
    TypicalRunDurationInSeconds types.Number `tfsdk:"typical_run_duration_in_seconds"`
    ExpectedRunIntervalInSeconds types.Number `tfsdk:"expected_run_interval_in_seconds"`
    Inputs types.String `tfsdk:"inputs"`
    Outputs types.String `tfsdk:"outputs"`
    FirstSeenAt types.String `tfsdk:"first_seen_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *DataPipelineDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_data_pipeline"
}

func (d *DataPipelineDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "A job that moves data - an Airflow task, a Spark job, a dbt model - as its OpenLineage events describe it, with how its last run went. Look up an existing data pipeline by `id`, or by any of its other arguments (`created_by_user_id`, `description`, `expected_run_interval_in_seconds`, ...): each one set must match, and exactly one data pipeline may match them all.",

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
            "job_namespace": schema.StringAttribute{
                MarkdownDescription: "The job's OpenLineage namespace: where it runs, such as airflow://prod.",
                Optional: true,
                Computed: true,
            },
            "job_name": schema.StringAttribute{
                MarkdownDescription: "The job's OpenLineage name, such as orders_daily.load_orders.",
                Optional: true,
                Computed: true,
            },
            "pipeline_key": schema.StringAttribute{
                MarkdownDescription: "The job's namespace and name together: what its events are matched by.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "What the pipeline does and who relies on it, in your own words.",
                Optional: true,
                Computed: true,
            },
            "last_run_id": schema.StringAttribute{
                MarkdownDescription: "The OpenLineage run id of the pipeline's latest run.",
                Optional: true,
                Computed: true,
            },
            "last_run_state": schema.StringAttribute{
                MarkdownDescription: "How the latest run stands: Running, Complete, Fail or Abort.",
                Optional: true,
                Computed: true,
            },
            "last_run_started_at": schema.StringAttribute{
                MarkdownDescription: "When the latest run started.",
                Computed: true,
            },
            "last_run_ended_at": schema.StringAttribute{
                MarkdownDescription: "When the latest run ended.",
                Computed: true,
            },
            "last_run_duration_in_seconds": schema.NumberAttribute{
                MarkdownDescription: "How long the latest finished run took, in seconds.",
                Optional: true,
                Computed: true,
            },
            "last_success_at": schema.StringAttribute{
                MarkdownDescription: "When a run last completed.",
                Computed: true,
            },
            "last_failure_at": schema.StringAttribute{
                MarkdownDescription: "When a run last failed or was aborted.",
                Computed: true,
            },
            "last_error_message": schema.StringAttribute{
                MarkdownDescription: "The error the latest failed run reported.",
                Optional: true,
                Computed: true,
            },
            "run_count": schema.NumberAttribute{
                MarkdownDescription: "How many runs OneUptime has seen.",
                Optional: true,
                Computed: true,
            },
            "typical_run_duration_in_seconds": schema.NumberAttribute{
                MarkdownDescription: "How long a run usually takes, in seconds: the median of the last 20 completed runs, once there are five.",
                Optional: true,
                Computed: true,
            },
            "expected_run_interval_in_seconds": schema.NumberAttribute{
                MarkdownDescription: "The longest the pipeline goes between runs on its schedule, in seconds: the longest gap between the starts of its last 50 runs, once there are five.",
                Optional: true,
                Computed: true,
            },
            "inputs": schema.StringAttribute{
                MarkdownDescription: "The datasets the pipeline reads, as namespace and name. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "outputs": schema.StringAttribute{
                MarkdownDescription: "The datasets the pipeline writes, as namespace and name. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "first_seen_at": schema.StringAttribute{
                MarkdownDescription: "When the pipeline's first event arrived.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: dataSourceRecordIDValidators("oneuptime_user"),
            },
        },
    }
}

func (d *DataPipelineDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataPipelineDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data DataPipelineDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.JobNamespace.IsNull() && !data.JobNamespace.IsUnknown() {
        filters["jobNamespace"] = data.JobNamespace.ValueString()
        filterNames = append(filterNames, "job_namespace = "+fmt.Sprintf("%q", data.JobNamespace.ValueString()))
    }
    if !data.JobName.IsNull() && !data.JobName.IsUnknown() {
        filters["jobName"] = data.JobName.ValueString()
        filterNames = append(filterNames, "job_name = "+fmt.Sprintf("%q", data.JobName.ValueString()))
    }
    if !data.PipelineKey.IsNull() && !data.PipelineKey.IsUnknown() {
        filters["pipelineKey"] = data.PipelineKey.ValueString()
        filterNames = append(filterNames, "pipeline_key = "+fmt.Sprintf("%q", data.PipelineKey.ValueString()))
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.LastRunId.IsNull() && !data.LastRunId.IsUnknown() {
        filters["lastRunId"] = data.LastRunId.ValueString()
        filterNames = append(filterNames, "last_run_id = "+fmt.Sprintf("%q", data.LastRunId.ValueString()))
    }
    if !data.LastRunState.IsNull() && !data.LastRunState.IsUnknown() {
        filters["lastRunState"] = data.LastRunState.ValueString()
        filterNames = append(filterNames, "last_run_state = "+fmt.Sprintf("%q", data.LastRunState.ValueString()))
    }
    if !data.LastRunDurationInSeconds.IsNull() && !data.LastRunDurationInSeconds.IsUnknown() {
        filters["lastRunDurationInSeconds"] = lookupNumber(data.LastRunDurationInSeconds)
        filterNames = append(filterNames, "last_run_duration_in_seconds = "+data.LastRunDurationInSeconds.ValueBigFloat().String())
    }
    if !data.LastErrorMessage.IsNull() && !data.LastErrorMessage.IsUnknown() {
        filters["lastErrorMessage"] = data.LastErrorMessage.ValueString()
        filterNames = append(filterNames, "last_error_message = "+fmt.Sprintf("%q", data.LastErrorMessage.ValueString()))
    }
    if !data.RunCount.IsNull() && !data.RunCount.IsUnknown() {
        filters["runCount"] = lookupNumber(data.RunCount)
        filterNames = append(filterNames, "run_count = "+data.RunCount.ValueBigFloat().String())
    }
    if !data.TypicalRunDurationInSeconds.IsNull() && !data.TypicalRunDurationInSeconds.IsUnknown() {
        filters["typicalRunDurationInSeconds"] = lookupNumber(data.TypicalRunDurationInSeconds)
        filterNames = append(filterNames, "typical_run_duration_in_seconds = "+data.TypicalRunDurationInSeconds.ValueBigFloat().String())
    }
    if !data.ExpectedRunIntervalInSeconds.IsNull() && !data.ExpectedRunIntervalInSeconds.IsUnknown() {
        filters["expectedRunIntervalInSeconds"] = lookupNumber(data.ExpectedRunIntervalInSeconds)
        filterNames = append(filterNames, "expected_run_interval_in_seconds = "+data.ExpectedRunIntervalInSeconds.ValueBigFloat().String())
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the data pipeline up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the data pipeline up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "jobNamespace": true,
        "jobName": true,
        "pipelineKey": true,
        "description": true,
        "lastRunId": true,
        "lastRunState": true,
        "lastRunStartedAt": true,
        "lastRunEndedAt": true,
        "lastRunDurationInSeconds": true,
        "lastSuccessAt": true,
        "lastFailureAt": true,
        "lastErrorMessage": true,
        "runCount": true,
        "typicalRunDurationInSeconds": true,
        "expectedRunIntervalInSeconds": true,
        "inputs": true,
        "outputs": true,
        "firstSeenAt": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/data-pipeline/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read data_pipeline, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data pipeline found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read data_pipeline: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/data-pipeline/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list data_pipeline, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list data_pipeline: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data pipeline matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one data pipeline matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for data_pipeline.")
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
    if obj, ok := item["jobNamespace"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.JobNamespace = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.JobNamespace = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.JobNamespace = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.JobNamespace = types.StringValue(string(jsonBytes))
        } else {
            data.JobNamespace = types.StringNull()
        }
    } else if val, ok := item["jobNamespace"].(string); ok {
        data.JobNamespace = types.StringValue(val)
    } else {
        data.JobNamespace = types.StringNull()
    }
    if obj, ok := item["jobName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.JobName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.JobName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.JobName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.JobName = types.StringValue(string(jsonBytes))
        } else {
            data.JobName = types.StringNull()
        }
    } else if val, ok := item["jobName"].(string); ok {
        data.JobName = types.StringValue(val)
    } else {
        data.JobName = types.StringNull()
    }
    if obj, ok := item["pipelineKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PipelineKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PipelineKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PipelineKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PipelineKey = types.StringValue(string(jsonBytes))
        } else {
            data.PipelineKey = types.StringNull()
        }
    } else if val, ok := item["pipelineKey"].(string); ok {
        data.PipelineKey = types.StringValue(val)
    } else {
        data.PipelineKey = types.StringNull()
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
    } else if val, ok := item["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := item["lastRunId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastRunId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastRunId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastRunId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastRunId = types.StringValue(string(jsonBytes))
        } else {
            data.LastRunId = types.StringNull()
        }
    } else if val, ok := item["lastRunId"].(string); ok {
        data.LastRunId = types.StringValue(val)
    } else {
        data.LastRunId = types.StringNull()
    }
    if obj, ok := item["lastRunState"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastRunState = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastRunState = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastRunState = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastRunState = types.StringValue(string(jsonBytes))
        } else {
            data.LastRunState = types.StringNull()
        }
    } else if val, ok := item["lastRunState"].(string); ok {
        data.LastRunState = types.StringValue(val)
    } else {
        data.LastRunState = types.StringNull()
    }
    if obj, ok := item["lastRunStartedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastRunStartedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastRunStartedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastRunStartedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastRunStartedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastRunStartedAt = types.StringNull()
        }
    } else if val, ok := item["lastRunStartedAt"].(string); ok {
        data.LastRunStartedAt = types.StringValue(val)
    } else {
        data.LastRunStartedAt = types.StringNull()
    }
    if obj, ok := item["lastRunEndedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastRunEndedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastRunEndedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastRunEndedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastRunEndedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastRunEndedAt = types.StringNull()
        }
    } else if val, ok := item["lastRunEndedAt"].(string); ok {
        data.LastRunEndedAt = types.StringValue(val)
    } else {
        data.LastRunEndedAt = types.StringNull()
    }
    data.LastRunDurationInSeconds = numberFromAPI(data.LastRunDurationInSeconds, item["lastRunDurationInSeconds"])
    if obj, ok := item["lastSuccessAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastSuccessAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastSuccessAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastSuccessAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastSuccessAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastSuccessAt = types.StringNull()
        }
    } else if val, ok := item["lastSuccessAt"].(string); ok {
        data.LastSuccessAt = types.StringValue(val)
    } else {
        data.LastSuccessAt = types.StringNull()
    }
    if obj, ok := item["lastFailureAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastFailureAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastFailureAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastFailureAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastFailureAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastFailureAt = types.StringNull()
        }
    } else if val, ok := item["lastFailureAt"].(string); ok {
        data.LastFailureAt = types.StringValue(val)
    } else {
        data.LastFailureAt = types.StringNull()
    }
    if obj, ok := item["lastErrorMessage"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastErrorMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastErrorMessage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastErrorMessage = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastErrorMessage = types.StringValue(string(jsonBytes))
        } else {
            data.LastErrorMessage = types.StringNull()
        }
    } else if val, ok := item["lastErrorMessage"].(string); ok {
        data.LastErrorMessage = types.StringValue(val)
    } else {
        data.LastErrorMessage = types.StringNull()
    }
    data.RunCount = numberFromAPI(data.RunCount, item["runCount"])
    data.TypicalRunDurationInSeconds = numberFromAPI(data.TypicalRunDurationInSeconds, item["typicalRunDurationInSeconds"])
    data.ExpectedRunIntervalInSeconds = numberFromAPI(data.ExpectedRunIntervalInSeconds, item["expectedRunIntervalInSeconds"])
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
    } else if val, ok := item["outputs"].(string); ok {
        data.Outputs = types.StringValue(val)
    } else {
        data.Outputs = types.StringNull()
    }
    if obj, ok := item["firstSeenAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FirstSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FirstSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FirstSeenAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FirstSeenAt = types.StringValue(string(jsonBytes))
        } else {
            data.FirstSeenAt = types.StringNull()
        }
    } else if val, ok := item["firstSeenAt"].(string); ok {
        data.FirstSeenAt = types.StringValue(val)
    } else {
        data.FirstSeenAt = types.StringNull()
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
