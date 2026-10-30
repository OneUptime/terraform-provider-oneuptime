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
var _ datasource.DataSource = &HuntDataSource{}

func NewHuntDataSource() datasource.DataSource {
    return &HuntDataSource{}
}

// HuntDataSource defines the data source implementation.
type HuntDataSource struct {
    client *Client
}

// HuntDataSourceModel describes the data source data model.
type HuntDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    SearchQuery types.String `tfsdk:"search_query"`
    Filters types.String `tfsdk:"filters"`
    TimeRange types.String `tfsdk:"time_range"`
    Schedule types.String `tfsdk:"schedule"`
    NotifyUsers types.Set `tfsdk:"notify_users"`
    NotifyTeams types.Set `tfsdk:"notify_teams"`
    ScheduledByUserId types.String `tfsdk:"scheduled_by_user_id"`
    SearchedUpToAt types.String `tfsdk:"searched_up_to_at"`
    NextRunAt types.String `tfsdk:"next_run_at"`
    LastRunAt types.String `tfsdk:"last_run_at"`
    LastRunEventCount types.Number `tfsdk:"last_run_event_count"`
    LastRunError types.String `tfsdk:"last_run_error"`
    LastFoundAt types.String `tfsdk:"last_found_at"`
    LastFoundSinceAt types.String `tfsdk:"last_found_since_at"`
    LastFoundEventCount types.Number `tfsdk:"last_found_event_count"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *HuntDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_hunt"
}

func (d *HuntDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Searches over security events saved with a name, to open again or to run on a schedule that tells people when it finds something. Look up an existing hunt by `id`, or by any of its other arguments (`name`, `created_by_user_id`, `description`, ...): each one set must match, and exactly one hunt may match them all.",

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
                MarkdownDescription: "ID of the project this hunt belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "What the hunt looks for, in a few words. At most 100 characters.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Why the hunt matters and what to do when it finds something, for whoever opens it later.",
                Optional: true,
                Computed: true,
            },
            "search_query": schema.StringAttribute{
                MarkdownDescription: "What the hunt searches for, in the Events search bar's language: words a message contains, fields such as severity:High or user:alice, and source attributes such as @threat.matched:true. At most 2,000 characters.",
                Optional: true,
                Computed: true,
            },
            "filters": schema.StringAttribute{
                MarkdownDescription: "The Events sidebar's filters kept beside the search: each field's values, or with ! before the field the values it leaves out. At most 50 values. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "time_range": schema.StringAttribute{
                MarkdownDescription: "How far back the hunt searches when it is opened, as the Events time picker names it: Past 1 Hour, Past 1 Day, Past 1 Week, Past 1 Month and the rest.",
                Optional: true,
                Computed: true,
            },
            "schedule": schema.StringAttribute{
                MarkdownDescription: "How often the hunt runs on its own and tells people what it found: Not Scheduled, Every Hour, Every 6 Hours, Every Day or Every Week. Each run searches the events that arrived since the run before.",
                Optional: true,
                Computed: true,
            },
            "notify_users": schema.SetAttribute{
                MarkdownDescription: "The people told when a scheduled run finds events, each as their notification settings say. Only people who can read security events are told. IDs of `oneuptime_user` records.",
                Computed: true,
                ElementType: types.StringType,
            },
            "notify_teams": schema.SetAttribute{
                MarkdownDescription: "The teams whose members are told when a scheduled run finds events, each as their notification settings say. Only members who can read security events are told. IDs of `oneuptime_team` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "scheduled_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person the hunt's scheduled runs run as. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "searched_up_to_at": schema.StringAttribute{
                MarkdownDescription: "How far the scheduled runs have searched: the next run reads the events from here on. Set to the moment the schedule is set or the search changes.",
                Computed: true,
            },
            "next_run_at": schema.StringAttribute{
                MarkdownDescription: "When the hunt next runs on its schedule. Empty for a hunt that is not scheduled.",
                Computed: true,
            },
            "last_run_at": schema.StringAttribute{
                MarkdownDescription: "When the hunt last ran on its schedule, whether it found anything or not.",
                Computed: true,
            },
            "last_run_event_count": schema.NumberAttribute{
                MarkdownDescription: "How many events the hunt's last scheduled run found.",
                Optional: true,
                Computed: true,
            },
            "last_run_error": schema.StringAttribute{
                MarkdownDescription: "Why the hunt's last scheduled run could not be made. Empty when it ran.",
                Optional: true,
                Computed: true,
            },
            "last_found_at": schema.StringAttribute{
                MarkdownDescription: "The end of the events the last run that found something searched: what it found arrived before this.",
                Computed: true,
            },
            "last_found_since_at": schema.StringAttribute{
                MarkdownDescription: "The start of the events the last run that found something searched: what it found arrived after this.",
                Computed: true,
            },
            "last_found_event_count": schema.NumberAttribute{
                MarkdownDescription: "How many events the last run that found something found.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who saved the hunt. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *HuntDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *HuntDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data HuntDataSourceModel

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
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.SearchQuery.IsNull() && !data.SearchQuery.IsUnknown() {
        filters["searchQuery"] = data.SearchQuery.ValueString()
        filterNames = append(filterNames, "search_query = "+fmt.Sprintf("%q", data.SearchQuery.ValueString()))
    }
    if !data.TimeRange.IsNull() && !data.TimeRange.IsUnknown() {
        filters["timeRange"] = data.TimeRange.ValueString()
        filterNames = append(filterNames, "time_range = "+fmt.Sprintf("%q", data.TimeRange.ValueString()))
    }
    if !data.Schedule.IsNull() && !data.Schedule.IsUnknown() {
        filters["schedule"] = data.Schedule.ValueString()
        filterNames = append(filterNames, "schedule = "+fmt.Sprintf("%q", data.Schedule.ValueString()))
    }
    if !data.ScheduledByUserId.IsNull() && !data.ScheduledByUserId.IsUnknown() {
        filters["scheduledByUserId"] = data.ScheduledByUserId.ValueString()
        filterNames = append(filterNames, "scheduled_by_user_id = "+fmt.Sprintf("%q", data.ScheduledByUserId.ValueString()))
    }
    if !data.LastRunEventCount.IsNull() && !data.LastRunEventCount.IsUnknown() {
        filters["lastRunEventCount"] = lookupNumber(data.LastRunEventCount)
        filterNames = append(filterNames, "last_run_event_count = "+data.LastRunEventCount.ValueBigFloat().String())
    }
    if !data.LastRunError.IsNull() && !data.LastRunError.IsUnknown() {
        filters["lastRunError"] = data.LastRunError.ValueString()
        filterNames = append(filterNames, "last_run_error = "+fmt.Sprintf("%q", data.LastRunError.ValueString()))
    }
    if !data.LastFoundEventCount.IsNull() && !data.LastFoundEventCount.IsUnknown() {
        filters["lastFoundEventCount"] = lookupNumber(data.LastFoundEventCount)
        filterNames = append(filterNames, "last_found_event_count = "+data.LastFoundEventCount.ValueBigFloat().String())
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the hunt up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the hunt up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "name": true,
        "description": true,
        "searchQuery": true,
        "filters": true,
        "timeRange": true,
        "schedule": true,
        "notifyUsers": true,
        "notifyTeams": true,
        "scheduledByUserId": true,
        "searchedUpToAt": true,
        "nextRunAt": true,
        "lastRunAt": true,
        "lastRunEventCount": true,
        "lastRunError": true,
        "lastFoundAt": true,
        "lastFoundSinceAt": true,
        "lastFoundEventCount": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/security-hunt/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read hunt, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No hunt found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read hunt: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/security-hunt/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list hunt, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list hunt: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No hunt matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one hunt matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for hunt.")
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
    if obj, ok := item["searchQuery"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SearchQuery = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SearchQuery = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SearchQuery = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SearchQuery = types.StringValue(string(jsonBytes))
        } else {
            data.SearchQuery = types.StringNull()
        }
    } else if val, ok := item["searchQuery"].(string); ok {
        data.SearchQuery = types.StringValue(val)
    } else {
        data.SearchQuery = types.StringNull()
    }
    if obj, ok := item["filters"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Filters = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Filters = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Filters = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Filters = types.StringValue(string(jsonBytes))
        } else {
            data.Filters = types.StringNull()
        }
    } else if val, ok := item["filters"].(string); ok {
        data.Filters = types.StringValue(val)
    } else {
        data.Filters = types.StringNull()
    }
    if obj, ok := item["timeRange"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TimeRange = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TimeRange = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TimeRange = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TimeRange = types.StringValue(string(jsonBytes))
        } else {
            data.TimeRange = types.StringNull()
        }
    } else if val, ok := item["timeRange"].(string); ok {
        data.TimeRange = types.StringValue(val)
    } else {
        data.TimeRange = types.StringNull()
    }
    if obj, ok := item["schedule"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Schedule = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Schedule = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Schedule = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Schedule = types.StringValue(string(jsonBytes))
        } else {
            data.Schedule = types.StringNull()
        }
    } else if val, ok := item["schedule"].(string); ok {
        data.Schedule = types.StringValue(val)
    } else {
        data.Schedule = types.StringNull()
    }
    if val, ok := item["notifyUsers"].([]interface{}); ok {
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
        data.NotifyUsers = types.SetValueMust(types.StringType, setItems)
    } else {
        data.NotifyUsers = types.SetNull(types.StringType)
    }
    if val, ok := item["notifyTeams"].([]interface{}); ok {
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
        data.NotifyTeams = types.SetValueMust(types.StringType, setItems)
    } else {
        data.NotifyTeams = types.SetNull(types.StringType)
    }
    if obj, ok := item["scheduledByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScheduledByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ScheduledByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ScheduledByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ScheduledByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ScheduledByUserId = types.StringNull()
        }
    } else if val, ok := item["scheduledByUserId"].(string); ok {
        data.ScheduledByUserId = types.StringValue(val)
    } else {
        data.ScheduledByUserId = types.StringNull()
    }
    if obj, ok := item["searchedUpToAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SearchedUpToAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SearchedUpToAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SearchedUpToAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SearchedUpToAt = types.StringValue(string(jsonBytes))
        } else {
            data.SearchedUpToAt = types.StringNull()
        }
    } else if val, ok := item["searchedUpToAt"].(string); ok {
        data.SearchedUpToAt = types.StringValue(val)
    } else {
        data.SearchedUpToAt = types.StringNull()
    }
    if obj, ok := item["nextRunAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NextRunAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.NextRunAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.NextRunAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.NextRunAt = types.StringValue(string(jsonBytes))
        } else {
            data.NextRunAt = types.StringNull()
        }
    } else if val, ok := item["nextRunAt"].(string); ok {
        data.NextRunAt = types.StringValue(val)
    } else {
        data.NextRunAt = types.StringNull()
    }
    if obj, ok := item["lastRunAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastRunAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastRunAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastRunAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastRunAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastRunAt = types.StringNull()
        }
    } else if val, ok := item["lastRunAt"].(string); ok {
        data.LastRunAt = types.StringValue(val)
    } else {
        data.LastRunAt = types.StringNull()
    }
    data.LastRunEventCount = numberFromAPI(data.LastRunEventCount, item["lastRunEventCount"])
    if obj, ok := item["lastRunError"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastRunError = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastRunError = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastRunError = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastRunError = types.StringValue(string(jsonBytes))
        } else {
            data.LastRunError = types.StringNull()
        }
    } else if val, ok := item["lastRunError"].(string); ok {
        data.LastRunError = types.StringValue(val)
    } else {
        data.LastRunError = types.StringNull()
    }
    if obj, ok := item["lastFoundAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastFoundAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastFoundAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastFoundAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastFoundAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastFoundAt = types.StringNull()
        }
    } else if val, ok := item["lastFoundAt"].(string); ok {
        data.LastFoundAt = types.StringValue(val)
    } else {
        data.LastFoundAt = types.StringNull()
    }
    if obj, ok := item["lastFoundSinceAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastFoundSinceAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastFoundSinceAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastFoundSinceAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastFoundSinceAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastFoundSinceAt = types.StringNull()
        }
    } else if val, ok := item["lastFoundSinceAt"].(string); ok {
        data.LastFoundSinceAt = types.StringValue(val)
    } else {
        data.LastFoundSinceAt = types.StringNull()
    }
    data.LastFoundEventCount = numberFromAPI(data.LastFoundEventCount, item["lastFoundEventCount"])
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
