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
var _ datasource.DataSource = &TeamOrLabelBudgetDataSource{}

func NewTeamOrLabelBudgetDataSource() datasource.DataSource {
    return &TeamOrLabelBudgetDataSource{}
}

// TeamOrLabelBudgetDataSource defines the data source implementation.
type TeamOrLabelBudgetDataSource struct {
    client *Client
}

// TeamOrLabelBudgetDataSourceModel describes the data source data model.
type TeamOrLabelBudgetDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    TeamId types.String `tfsdk:"team_id"`
    LabelId types.String `tfsdk:"label_id"`
    LogsDailyIngestionBudgetInGb types.Number `tfsdk:"logs_daily_ingestion_budget_in_gb"`
    TracesDailyIngestionBudgetInGb types.Number `tfsdk:"traces_daily_ingestion_budget_in_gb"`
    MetricsDailyIngestionBudgetInGb types.Number `tfsdk:"metrics_daily_ingestion_budget_in_gb"`
    ProfilesDailyIngestionBudgetInGb types.Number `tfsdk:"profiles_daily_ingestion_budget_in_gb"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *TeamOrLabelBudgetDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_team_or_label_budget"
}

func (d *TeamOrLabelBudgetDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Daily ingestion budgets for the telemetry of the services a team owns, or that carry a label. Look up an existing team or label budget by `id`, or by any of its other arguments (`created_by_user_id`, `label_id`, `logs_daily_ingestion_budget_in_gb`, ...): each one set must match, and exactly one team or label budget may match them all.",

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
            "team_id": schema.StringAttribute{
                MarkdownDescription: "ID of the team whose services the budget is for. Empty for a label's budget. The ID of a `oneuptime_team`.",
                Optional: true,
                Computed: true,
                Validators: RecordIDValidators("oneuptime_team"),
            },
            "label_id": schema.StringAttribute{
                MarkdownDescription: "ID of the label whose services the budget is for. Empty for a team's budget. The ID of a `oneuptime_label`.",
                Optional: true,
                Computed: true,
                Validators: RecordIDValidators("oneuptime_label"),
            },
            "logs_daily_ingestion_budget_in_gb": schema.NumberAttribute{
                MarkdownDescription: "The most logs the team's or label's services may send in a budget day, in GB, counted decompressed as they arrive; once used up, their logs are dropped until the budget resets. At most the project's own budget. Empty: no logs budget of their own (the project's still applies).",
                Optional: true,
                Computed: true,
            },
            "traces_daily_ingestion_budget_in_gb": schema.NumberAttribute{
                MarkdownDescription: "The most traces the team's or label's services may send in a budget day, in GB, counted decompressed as they arrive; once used up, their traces are dropped until the budget resets. At most the project's own budget. Empty: no traces budget of their own (the project's still applies).",
                Optional: true,
                Computed: true,
            },
            "metrics_daily_ingestion_budget_in_gb": schema.NumberAttribute{
                MarkdownDescription: "The most metrics the team's or label's services may send in a budget day, in GB, counted decompressed as they arrive; once used up, their metrics are dropped until the budget resets. At most the project's own budget. Empty: no metrics budget of their own (the project's still applies).",
                Optional: true,
                Computed: true,
            },
            "profiles_daily_ingestion_budget_in_gb": schema.NumberAttribute{
                MarkdownDescription: "The most profiles the team's or label's services may send in a budget day, in GB, counted decompressed as they arrive; once used up, their profiles are dropped until the budget resets. At most the project's own budget. Empty: no profiles budget of their own (the project's still applies).",
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

func (d *TeamOrLabelBudgetDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TeamOrLabelBudgetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data TeamOrLabelBudgetDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.TeamId.IsNull() && !data.TeamId.IsUnknown() {
        filters["teamId"] = data.TeamId.ValueString()
        filterNames = append(filterNames, "team_id = "+fmt.Sprintf("%q", data.TeamId.ValueString()))
    }
    if !data.LabelId.IsNull() && !data.LabelId.IsUnknown() {
        filters["labelId"] = data.LabelId.ValueString()
        filterNames = append(filterNames, "label_id = "+fmt.Sprintf("%q", data.LabelId.ValueString()))
    }
    if !data.LogsDailyIngestionBudgetInGb.IsNull() && !data.LogsDailyIngestionBudgetInGb.IsUnknown() {
        filters["logsDailyIngestionBudgetInGB"] = LookupNumber(data.LogsDailyIngestionBudgetInGb)
        filterNames = append(filterNames, "logs_daily_ingestion_budget_in_gb = "+data.LogsDailyIngestionBudgetInGb.ValueBigFloat().String())
    }
    if !data.TracesDailyIngestionBudgetInGb.IsNull() && !data.TracesDailyIngestionBudgetInGb.IsUnknown() {
        filters["tracesDailyIngestionBudgetInGB"] = LookupNumber(data.TracesDailyIngestionBudgetInGb)
        filterNames = append(filterNames, "traces_daily_ingestion_budget_in_gb = "+data.TracesDailyIngestionBudgetInGb.ValueBigFloat().String())
    }
    if !data.MetricsDailyIngestionBudgetInGb.IsNull() && !data.MetricsDailyIngestionBudgetInGb.IsUnknown() {
        filters["metricsDailyIngestionBudgetInGB"] = LookupNumber(data.MetricsDailyIngestionBudgetInGb)
        filterNames = append(filterNames, "metrics_daily_ingestion_budget_in_gb = "+data.MetricsDailyIngestionBudgetInGb.ValueBigFloat().String())
    }
    if !data.ProfilesDailyIngestionBudgetInGb.IsNull() && !data.ProfilesDailyIngestionBudgetInGb.IsUnknown() {
        filters["profilesDailyIngestionBudgetInGB"] = LookupNumber(data.ProfilesDailyIngestionBudgetInGb)
        filterNames = append(filterNames, "profiles_daily_ingestion_budget_in_gb = "+data.ProfilesDailyIngestionBudgetInGb.ValueBigFloat().String())
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the team or label budget up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the team or label budget up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "teamId": true,
        "labelId": true,
        "logsDailyIngestionBudgetInGB": true,
        "tracesDailyIngestionBudgetInGB": true,
        "metricsDailyIngestionBudgetInGB": true,
        "profilesDailyIngestionBudgetInGB": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/telemetry-ingestion-team-label-budget/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read team_or_label_budget, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No team or label budget found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read team_or_label_budget: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/telemetry-ingestion-team-label-budget/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list team_or_label_budget, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list team_or_label_budget: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No team or label budget matches %s.", DescribeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one team or label budget matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", DescribeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for team_or_label_budget.")
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
    if obj, ok := item["teamId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TeamId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TeamId = types.StringValue(string(jsonBytes))
        } else {
            data.TeamId = types.StringNull()
        }
    } else if arr, ok := item["teamId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.TeamId = types.StringValue(string(jsonBytes))
        } else {
            data.TeamId = types.StringNull()
        }
    } else if val, ok := item["teamId"].(string); ok {
        data.TeamId = types.StringValue(val)
    } else {
        data.TeamId = types.StringNull()
    }
    if obj, ok := item["labelId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LabelId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LabelId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LabelId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LabelId = types.StringValue(string(jsonBytes))
        } else {
            data.LabelId = types.StringNull()
        }
    } else if arr, ok := item["labelId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LabelId = types.StringValue(string(jsonBytes))
        } else {
            data.LabelId = types.StringNull()
        }
    } else if val, ok := item["labelId"].(string); ok {
        data.LabelId = types.StringValue(val)
    } else {
        data.LabelId = types.StringNull()
    }
    data.LogsDailyIngestionBudgetInGb = NumberFromAPI(data.LogsDailyIngestionBudgetInGb, item["logsDailyIngestionBudgetInGB"])
    data.TracesDailyIngestionBudgetInGb = NumberFromAPI(data.TracesDailyIngestionBudgetInGb, item["tracesDailyIngestionBudgetInGB"])
    data.MetricsDailyIngestionBudgetInGb = NumberFromAPI(data.MetricsDailyIngestionBudgetInGb, item["metricsDailyIngestionBudgetInGB"])
    data.ProfilesDailyIngestionBudgetInGb = NumberFromAPI(data.ProfilesDailyIngestionBudgetInGb, item["profilesDailyIngestionBudgetInGB"])
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
