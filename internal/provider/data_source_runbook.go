package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &RunbookDataSource{}

func NewRunbookDataSource() datasource.DataSource {
    return &RunbookDataSource{}
}

// RunbookDataSource defines the data source implementation.
type RunbookDataSource struct {
    client *Client
}

// RunbookDataSourceModel describes the data source data model.
type RunbookDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Slug types.String `tfsdk:"slug"`
    Description types.String `tfsdk:"description"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    Steps types.String `tfsdk:"steps"`
    Labels types.Set `tfsdk:"labels"`
    IsScheduleEnabled types.Bool `tfsdk:"is_schedule_enabled"`
    ScheduleCron types.String `tfsdk:"schedule_cron"`
    ScheduleTimezone types.String `tfsdk:"schedule_timezone"`
    ScheduleOverlapPolicy types.String `tfsdk:"schedule_overlap_policy"`
    ScheduleMissedRunPolicy types.String `tfsdk:"schedule_missed_run_policy"`
    ScheduleNextRunAt types.String `tfsdk:"schedule_next_run_at"`
    ScheduleQueuedRunAt types.String `tfsdk:"schedule_queued_run_at"`
    IsWebhookEnabled types.Bool `tfsdk:"is_webhook_enabled"`
}

func (d *RunbookDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_runbook"
}

func (d *RunbookDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Reusable response procedures (manual checklists or scripts) that can be attached to incidents, alerts, or scheduled maintenance. Look up an existing runbook by `id`, or by any of its other arguments (`name`, `created_by_user_id`, `description`, ...): each one set must match, and exactly one runbook may match them all.",

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
            "name": schema.StringAttribute{
                MarkdownDescription: "Any friendly name of this object.",
                Optional: true,
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description that will help you remember.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Is this runbook enabled?",
                Optional: true,
                Computed: true,
            },
            "steps": schema.StringAttribute{
                MarkdownDescription: "Ordered list of steps to run for this runbook. Each step is one of Manual, JavaScript, HTTP request, Bash or AI. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "is_schedule_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether the runbook runs on its schedule (Schedule Cron, read in Schedule Time Zone). Turning it on takes permission to run the runbook as well as to edit it.",
                Optional: true,
                Computed: true,
            },
            "schedule_cron": schema.StringAttribute{
                MarkdownDescription: "When the runbook runs on its schedule: a five-field cron expression (minute, hour, day of month, month, day of week), read on the wall clock of Schedule Time Zone. Runbooks run at most once a minute.",
                Optional: true,
                Computed: true,
            },
            "schedule_timezone": schema.StringAttribute{
                MarkdownDescription: "The IANA time zone the schedule is read in, such as Europe/London or America/New_York. A 09:00 schedule runs at 09:00 there all year, also after the clocks change for daylight saving time. Left out, it is UTC.",
                Optional: true,
                Computed: true,
            },
            "schedule_overlap_policy": schema.StringAttribute{
                MarkdownDescription: "What a scheduled time does while a run of the runbook is still queued, running or waiting for someone: \"Skip\" (the default) starts nothing for that time, \"RunAfter\" starts one run as soon as the current run ends.",
                Optional: true,
                Computed: true,
            },
            "schedule_missed_run_policy": schema.StringAttribute{
                MarkdownDescription: "What happens to scheduled times OneUptime could not start a run at, because it was down: \"RunOnce\" (the default) starts one run when it is back, \"Skip\" starts none.",
                Optional: true,
                Computed: true,
            },
            "schedule_next_run_at": schema.StringAttribute{
                MarkdownDescription: "When the runbook's schedule starts its next run. Worked out by OneUptime; empty while the schedule is off.",
                Computed: true,
            },
            "schedule_queued_run_at": schema.StringAttribute{
                MarkdownDescription: "The scheduled time a run is waiting for the current run to end for, when \"If a Run Is Still Going\" is \"RunAfter\". Set and cleared by OneUptime.",
                Computed: true,
            },
            "is_webhook_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether a POST request to the runbook's webhook URL starts a run. The URL holds a secret key OneUptime generates when this is first turned on, shown on the runbook's Triggers page to people who may edit and run the runbook. Turning it on takes permission to run the runbook as well as to edit it.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *RunbookDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RunbookDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data RunbookDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.Slug.IsNull() && !data.Slug.IsUnknown() {
        filters["slug"] = data.Slug.ValueString()
        filterNames = append(filterNames, "slug = "+fmt.Sprintf("%q", data.Slug.ValueString()))
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.IsEnabled.IsNull() && !data.IsEnabled.IsUnknown() {
        filters["isEnabled"] = data.IsEnabled.ValueBool()
        filterNames = append(filterNames, "is_enabled = "+fmt.Sprintf("%t", data.IsEnabled.ValueBool()))
    }
    if !data.IsScheduleEnabled.IsNull() && !data.IsScheduleEnabled.IsUnknown() {
        filters["isScheduleEnabled"] = data.IsScheduleEnabled.ValueBool()
        filterNames = append(filterNames, "is_schedule_enabled = "+fmt.Sprintf("%t", data.IsScheduleEnabled.ValueBool()))
    }
    if !data.ScheduleCron.IsNull() && !data.ScheduleCron.IsUnknown() {
        filters["scheduleCron"] = data.ScheduleCron.ValueString()
        filterNames = append(filterNames, "schedule_cron = "+fmt.Sprintf("%q", data.ScheduleCron.ValueString()))
    }
    if !data.ScheduleTimezone.IsNull() && !data.ScheduleTimezone.IsUnknown() {
        filters["scheduleTimezone"] = data.ScheduleTimezone.ValueString()
        filterNames = append(filterNames, "schedule_timezone = "+fmt.Sprintf("%q", data.ScheduleTimezone.ValueString()))
    }
    if !data.ScheduleOverlapPolicy.IsNull() && !data.ScheduleOverlapPolicy.IsUnknown() {
        filters["scheduleOverlapPolicy"] = data.ScheduleOverlapPolicy.ValueString()
        filterNames = append(filterNames, "schedule_overlap_policy = "+fmt.Sprintf("%q", data.ScheduleOverlapPolicy.ValueString()))
    }
    if !data.ScheduleMissedRunPolicy.IsNull() && !data.ScheduleMissedRunPolicy.IsUnknown() {
        filters["scheduleMissedRunPolicy"] = data.ScheduleMissedRunPolicy.ValueString()
        filterNames = append(filterNames, "schedule_missed_run_policy = "+fmt.Sprintf("%q", data.ScheduleMissedRunPolicy.ValueString()))
    }
    if !data.IsWebhookEnabled.IsNull() && !data.IsWebhookEnabled.IsUnknown() {
        filters["isWebhookEnabled"] = data.IsWebhookEnabled.ValueBool()
        filterNames = append(filterNames, "is_webhook_enabled = "+fmt.Sprintf("%t", data.IsWebhookEnabled.ValueBool()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the runbook up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the runbook up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "name": true,
        "slug": true,
        "description": true,
        "createdByUserId": true,
        "isEnabled": true,
        "steps": true,
        "labels": true,
        "isScheduleEnabled": true,
        "scheduleCron": true,
        "scheduleTimezone": true,
        "scheduleOverlapPolicy": true,
        "scheduleMissedRunPolicy": true,
        "scheduleNextRunAt": true,
        "scheduleQueuedRunAt": true,
        "isWebhookEnabled": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/runbook/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read runbook, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No runbook found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read runbook: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/runbook/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list runbook, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list runbook: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No runbook matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one runbook matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for runbook.")
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
    if obj, ok := item["name"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := item["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := item["slug"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := item["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
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
    if val, ok := item["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    } else {
        data.IsEnabled = types.BoolNull()
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
    if val, ok := item["labels"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.Labels = types.SetValueMust(types.StringType, setItems)
    } else {
        data.Labels = types.SetNull(types.StringType)
    }
    if val, ok := item["isScheduleEnabled"].(bool); ok {
        data.IsScheduleEnabled = types.BoolValue(val)
    } else {
        data.IsScheduleEnabled = types.BoolNull()
    }
    if obj, ok := item["scheduleCron"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScheduleCron = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ScheduleCron = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ScheduleCron = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ScheduleCron = types.StringValue(string(jsonBytes))
        } else {
            data.ScheduleCron = types.StringNull()
        }
    } else if val, ok := item["scheduleCron"].(string); ok {
        data.ScheduleCron = types.StringValue(val)
    } else {
        data.ScheduleCron = types.StringNull()
    }
    if obj, ok := item["scheduleTimezone"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScheduleTimezone = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ScheduleTimezone = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ScheduleTimezone = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ScheduleTimezone = types.StringValue(string(jsonBytes))
        } else {
            data.ScheduleTimezone = types.StringNull()
        }
    } else if val, ok := item["scheduleTimezone"].(string); ok {
        data.ScheduleTimezone = types.StringValue(val)
    } else {
        data.ScheduleTimezone = types.StringNull()
    }
    if obj, ok := item["scheduleOverlapPolicy"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScheduleOverlapPolicy = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ScheduleOverlapPolicy = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ScheduleOverlapPolicy = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ScheduleOverlapPolicy = types.StringValue(string(jsonBytes))
        } else {
            data.ScheduleOverlapPolicy = types.StringNull()
        }
    } else if val, ok := item["scheduleOverlapPolicy"].(string); ok {
        data.ScheduleOverlapPolicy = types.StringValue(val)
    } else {
        data.ScheduleOverlapPolicy = types.StringNull()
    }
    if obj, ok := item["scheduleMissedRunPolicy"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScheduleMissedRunPolicy = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ScheduleMissedRunPolicy = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ScheduleMissedRunPolicy = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ScheduleMissedRunPolicy = types.StringValue(string(jsonBytes))
        } else {
            data.ScheduleMissedRunPolicy = types.StringNull()
        }
    } else if val, ok := item["scheduleMissedRunPolicy"].(string); ok {
        data.ScheduleMissedRunPolicy = types.StringValue(val)
    } else {
        data.ScheduleMissedRunPolicy = types.StringNull()
    }
    if obj, ok := item["scheduleNextRunAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScheduleNextRunAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ScheduleNextRunAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ScheduleNextRunAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ScheduleNextRunAt = types.StringValue(string(jsonBytes))
        } else {
            data.ScheduleNextRunAt = types.StringNull()
        }
    } else if val, ok := item["scheduleNextRunAt"].(string); ok {
        data.ScheduleNextRunAt = types.StringValue(val)
    } else {
        data.ScheduleNextRunAt = types.StringNull()
    }
    if obj, ok := item["scheduleQueuedRunAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScheduleQueuedRunAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ScheduleQueuedRunAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ScheduleQueuedRunAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ScheduleQueuedRunAt = types.StringValue(string(jsonBytes))
        } else {
            data.ScheduleQueuedRunAt = types.StringNull()
        }
    } else if val, ok := item["scheduleQueuedRunAt"].(string); ok {
        data.ScheduleQueuedRunAt = types.StringValue(val)
    } else {
        data.ScheduleQueuedRunAt = types.StringNull()
    }
    if val, ok := item["isWebhookEnabled"].(bool); ok {
        data.IsWebhookEnabled = types.BoolValue(val)
    } else {
        data.IsWebhookEnabled = types.BoolNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
