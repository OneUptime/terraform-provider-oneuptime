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
var _ datasource.DataSource = &ExceptionActivityDataSource{}

func NewExceptionActivityDataSource() datasource.DataSource {
    return &ExceptionActivityDataSource{}
}

// ExceptionActivityDataSource defines the data source implementation.
type ExceptionActivityDataSource struct {
    client *Client
}

// ExceptionActivityDataSourceModel describes the data source data model.
type ExceptionActivityDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    TelemetryExceptionId types.String `tfsdk:"telemetry_exception_id"`
    ActivityType types.String `tfsdk:"activity_type"`
    UserId types.String `tfsdk:"user_id"`
    ApiKeyName types.String `tfsdk:"api_key_name"`
    AssigneeUserId types.String `tfsdk:"assignee_user_id"`
    AssigneeTeamId types.String `tfsdk:"assignee_team_id"`
    RelatedTelemetryExceptionId types.String `tfsdk:"related_telemetry_exception_id"`
    Release types.String `tfsdk:"release"`
    IsAssigneeNotified types.Bool `tfsdk:"is_assignee_notified"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *ExceptionActivityDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_exception_activity"
}

func (d *ExceptionActivityDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "The activity of an exception: who resolved, reopened, archived, assigned, merged or unmerged it and when, and when it came back after it was resolved. Read one exception's activity at a time, by its telemetryExceptionId. Look up an existing exception activity by `id`, or by any of its other arguments (`activity_type`, `api_key_name`, `assignee_team_id`, ...): each one set must match, and exactly one exception activity may match them all.",

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
                MarkdownDescription: "ID of the project this entry belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "telemetry_exception_id": schema.StringAttribute{
                MarkdownDescription: "ID of the exception this entry is about. The ID of a `oneuptime_exception`.",
                Optional: true,
                Computed: true,
                Validators: recordIDValidators("oneuptime_exception"),
            },
            "activity_type": schema.StringAttribute{
                MarkdownDescription: "What happened: Resolved, ResolvedInNextRelease, Reopened, Regressed, Archived, Unarchived, Assigned, Unassigned, Merged, MergedInto, Unmerged or UnmergedFrom.",
                Optional: true,
                Computed: true,
            },
            "user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who did it, when a person did. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: dataSourceRecordIDValidators("oneuptime_user"),
            },
            "api_key_name": schema.StringAttribute{
                MarkdownDescription: "The name the API key that made the change had at the time, when an API key did.",
                Optional: true,
                Computed: true,
            },
            "assignee_user_id": schema.StringAttribute{
                MarkdownDescription: "For Assigned, ID of the person it was assigned to; for Unassigned, of the person it was taken from. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: dataSourceRecordIDValidators("oneuptime_user"),
            },
            "assignee_team_id": schema.StringAttribute{
                MarkdownDescription: "For Assigned, ID of the team it was assigned to; for Unassigned, of the team it was taken from. The ID of a `oneuptime_team`.",
                Optional: true,
                Computed: true,
                Validators: recordIDValidators("oneuptime_team"),
            },
            "related_telemetry_exception_id": schema.StringAttribute{
                MarkdownDescription: "For a merge or an unmerge, ID of the other exception. The ID of a `oneuptime_exception`.",
                Optional: true,
                Computed: true,
                Validators: recordIDValidators("oneuptime_exception"),
            },
            "release": schema.StringAttribute{
                MarkdownDescription: "For ResolvedInNextRelease, the release it waits out; for Regressed, the release of the occurrence that brought it back, when it names one.",
                Optional: true,
                Computed: true,
            },
            "is_assignee_notified": schema.BoolAttribute{
                MarkdownDescription: "For Assigned, whether the person, or the members of the team, it was assigned to were notified. True for every other entry.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person whose action wrote the entry. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: dataSourceRecordIDValidators("oneuptime_user"),
            },
        },
    }
}

func (d *ExceptionActivityDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ExceptionActivityDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ExceptionActivityDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.TelemetryExceptionId.IsNull() && !data.TelemetryExceptionId.IsUnknown() {
        filters["telemetryExceptionId"] = data.TelemetryExceptionId.ValueString()
        filterNames = append(filterNames, "telemetry_exception_id = "+fmt.Sprintf("%q", data.TelemetryExceptionId.ValueString()))
    }
    if !data.ActivityType.IsNull() && !data.ActivityType.IsUnknown() {
        filters["activityType"] = data.ActivityType.ValueString()
        filterNames = append(filterNames, "activity_type = "+fmt.Sprintf("%q", data.ActivityType.ValueString()))
    }
    if !data.UserId.IsNull() && !data.UserId.IsUnknown() {
        filters["userId"] = data.UserId.ValueString()
        filterNames = append(filterNames, "user_id = "+fmt.Sprintf("%q", data.UserId.ValueString()))
    }
    if !data.ApiKeyName.IsNull() && !data.ApiKeyName.IsUnknown() {
        filters["apiKeyName"] = data.ApiKeyName.ValueString()
        filterNames = append(filterNames, "api_key_name = "+fmt.Sprintf("%q", data.ApiKeyName.ValueString()))
    }
    if !data.AssigneeUserId.IsNull() && !data.AssigneeUserId.IsUnknown() {
        filters["assigneeUserId"] = data.AssigneeUserId.ValueString()
        filterNames = append(filterNames, "assignee_user_id = "+fmt.Sprintf("%q", data.AssigneeUserId.ValueString()))
    }
    if !data.AssigneeTeamId.IsNull() && !data.AssigneeTeamId.IsUnknown() {
        filters["assigneeTeamId"] = data.AssigneeTeamId.ValueString()
        filterNames = append(filterNames, "assignee_team_id = "+fmt.Sprintf("%q", data.AssigneeTeamId.ValueString()))
    }
    if !data.RelatedTelemetryExceptionId.IsNull() && !data.RelatedTelemetryExceptionId.IsUnknown() {
        filters["relatedTelemetryExceptionId"] = data.RelatedTelemetryExceptionId.ValueString()
        filterNames = append(filterNames, "related_telemetry_exception_id = "+fmt.Sprintf("%q", data.RelatedTelemetryExceptionId.ValueString()))
    }
    if !data.Release.IsNull() && !data.Release.IsUnknown() {
        filters["release"] = data.Release.ValueString()
        filterNames = append(filterNames, "release = "+fmt.Sprintf("%q", data.Release.ValueString()))
    }
    if !data.IsAssigneeNotified.IsNull() && !data.IsAssigneeNotified.IsUnknown() {
        filters["isAssigneeNotified"] = data.IsAssigneeNotified.ValueBool()
        filterNames = append(filterNames, "is_assignee_notified = "+fmt.Sprintf("%t", data.IsAssigneeNotified.ValueBool()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the exception activity up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the exception activity up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "telemetryExceptionId": true,
        "activityType": true,
        "userId": true,
        "apiKeyName": true,
        "assigneeUserId": true,
        "assigneeTeamId": true,
        "relatedTelemetryExceptionId": true,
        "release": true,
        "isAssigneeNotified": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/telemetry-exception-activity/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read exception_activity, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No exception activity found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read exception_activity: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/telemetry-exception-activity/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list exception_activity, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list exception_activity: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No exception activity matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one exception activity matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for exception_activity.")
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
    if obj, ok := item["telemetryExceptionId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TelemetryExceptionId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TelemetryExceptionId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TelemetryExceptionId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TelemetryExceptionId = types.StringValue(string(jsonBytes))
        } else {
            data.TelemetryExceptionId = types.StringNull()
        }
    } else if arr, ok := item["telemetryExceptionId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.TelemetryExceptionId = types.StringValue(string(jsonBytes))
        } else {
            data.TelemetryExceptionId = types.StringNull()
        }
    } else if val, ok := item["telemetryExceptionId"].(string); ok {
        data.TelemetryExceptionId = types.StringValue(val)
    } else {
        data.TelemetryExceptionId = types.StringNull()
    }
    if obj, ok := item["activityType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ActivityType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ActivityType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ActivityType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ActivityType = types.StringValue(string(jsonBytes))
        } else {
            data.ActivityType = types.StringNull()
        }
    } else if arr, ok := item["activityType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ActivityType = types.StringValue(string(jsonBytes))
        } else {
            data.ActivityType = types.StringNull()
        }
    } else if val, ok := item["activityType"].(string); ok {
        data.ActivityType = types.StringValue(val)
    } else {
        data.ActivityType = types.StringNull()
    }
    if obj, ok := item["userId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UserId = types.StringValue(string(jsonBytes))
        } else {
            data.UserId = types.StringNull()
        }
    } else if arr, ok := item["userId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.UserId = types.StringValue(string(jsonBytes))
        } else {
            data.UserId = types.StringNull()
        }
    } else if val, ok := item["userId"].(string); ok {
        data.UserId = types.StringValue(val)
    } else {
        data.UserId = types.StringNull()
    }
    if obj, ok := item["apiKeyName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ApiKeyName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ApiKeyName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ApiKeyName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ApiKeyName = types.StringValue(string(jsonBytes))
        } else {
            data.ApiKeyName = types.StringNull()
        }
    } else if arr, ok := item["apiKeyName"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ApiKeyName = types.StringValue(string(jsonBytes))
        } else {
            data.ApiKeyName = types.StringNull()
        }
    } else if val, ok := item["apiKeyName"].(string); ok {
        data.ApiKeyName = types.StringValue(val)
    } else {
        data.ApiKeyName = types.StringNull()
    }
    if obj, ok := item["assigneeUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AssigneeUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AssigneeUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AssigneeUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AssigneeUserId = types.StringValue(string(jsonBytes))
        } else {
            data.AssigneeUserId = types.StringNull()
        }
    } else if arr, ok := item["assigneeUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AssigneeUserId = types.StringValue(string(jsonBytes))
        } else {
            data.AssigneeUserId = types.StringNull()
        }
    } else if val, ok := item["assigneeUserId"].(string); ok {
        data.AssigneeUserId = types.StringValue(val)
    } else {
        data.AssigneeUserId = types.StringNull()
    }
    if obj, ok := item["assigneeTeamId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AssigneeTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AssigneeTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AssigneeTeamId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AssigneeTeamId = types.StringValue(string(jsonBytes))
        } else {
            data.AssigneeTeamId = types.StringNull()
        }
    } else if arr, ok := item["assigneeTeamId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AssigneeTeamId = types.StringValue(string(jsonBytes))
        } else {
            data.AssigneeTeamId = types.StringNull()
        }
    } else if val, ok := item["assigneeTeamId"].(string); ok {
        data.AssigneeTeamId = types.StringValue(val)
    } else {
        data.AssigneeTeamId = types.StringNull()
    }
    if obj, ok := item["relatedTelemetryExceptionId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RelatedTelemetryExceptionId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RelatedTelemetryExceptionId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RelatedTelemetryExceptionId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RelatedTelemetryExceptionId = types.StringValue(string(jsonBytes))
        } else {
            data.RelatedTelemetryExceptionId = types.StringNull()
        }
    } else if arr, ok := item["relatedTelemetryExceptionId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.RelatedTelemetryExceptionId = types.StringValue(string(jsonBytes))
        } else {
            data.RelatedTelemetryExceptionId = types.StringNull()
        }
    } else if val, ok := item["relatedTelemetryExceptionId"].(string); ok {
        data.RelatedTelemetryExceptionId = types.StringValue(val)
    } else {
        data.RelatedTelemetryExceptionId = types.StringNull()
    }
    if obj, ok := item["release"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Release = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Release = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Release = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Release = types.StringValue(string(jsonBytes))
        } else {
            data.Release = types.StringNull()
        }
    } else if arr, ok := item["release"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Release = types.StringValue(string(jsonBytes))
        } else {
            data.Release = types.StringNull()
        }
    } else if val, ok := item["release"].(string); ok {
        data.Release = types.StringValue(val)
    } else {
        data.Release = types.StringNull()
    }
    if val, ok := item["isAssigneeNotified"].(bool); ok {
        data.IsAssigneeNotified = types.BoolValue(val)
    } else {
        data.IsAssigneeNotified = types.BoolNull()
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
