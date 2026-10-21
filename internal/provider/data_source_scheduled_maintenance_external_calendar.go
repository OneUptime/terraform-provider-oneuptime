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
var _ datasource.DataSource = &ScheduledMaintenanceExternalCalendarDataSource{}

func NewScheduledMaintenanceExternalCalendarDataSource() datasource.DataSource {
    return &ScheduledMaintenanceExternalCalendarDataSource{}
}

// ScheduledMaintenanceExternalCalendarDataSource defines the data source implementation.
type ScheduledMaintenanceExternalCalendarDataSource struct {
    client *Client
}

// ScheduledMaintenanceExternalCalendarDataSourceModel describes the data source data model.
type ScheduledMaintenanceExternalCalendarDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    CalendarHost types.String `tfsdk:"calendar_host"`
    Color types.String `tfsdk:"color"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    RefreshIntervalInMinutes types.Number `tfsdk:"refresh_interval_in_minutes"`
    LastReadAt types.String `tfsdk:"last_read_at"`
    LastReadAttemptAt types.String `tfsdk:"last_read_attempt_at"`
    LastReadError types.String `tfsdk:"last_read_error"`
    LastReadHttpStatus types.Number `tfsdk:"last_read_http_status"`
    EventCount types.Number `tfsdk:"event_count"`
    HasMoreEvents types.Bool `tfsdk:"has_more_events"`
    NextReadAt types.String `tfsdk:"next_read_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *ScheduledMaintenanceExternalCalendarDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_scheduled_maintenance_external_calendar"
}

func (d *ScheduledMaintenanceExternalCalendarDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Calendars from outside the project - public holidays, change freezes, another project's maintenance - read from their iCalendar (.ics) link and shown on the scheduled maintenance calendar. Read-only, and never part of the project's own calendar links. Look up an existing scheduled maintenance external calendar by `id`, or by any of its other arguments (`name`, `calendar_host`, `created_by_user_id`, ...): each one set must match, and exactly one scheduled maintenance external calendar may match them all.",

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
                MarkdownDescription: "What the calendar holds, as the maintenance calendar names it: Public holidays, Change freezes.",
                Optional: true,
                Computed: true,
            },
            "calendar_host": schema.StringAttribute{
                MarkdownDescription: "The host the calendar link points at, such as calendar.google.com - what the page shows in place of the secret link.",
                Optional: true,
                Computed: true,
            },
            "color": schema.StringAttribute{
                MarkdownDescription: "The colour the calendar's events are drawn in on the maintenance calendar, in hex (#0d9488 for example). Picked for you when left out.",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "When off, the calendar is not read and its events are not shown, until it is turned on again.",
                Optional: true,
                Computed: true,
            },
            "refresh_interval_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "How often OneUptime reads the calendar: 15, 30, 60, 360 or 1440 minutes.",
                Optional: true,
                Computed: true,
            },
            "last_read_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime last read the calendar without a problem. Empty until the first read.",
                Computed: true,
            },
            "last_read_attempt_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime last tried to read the calendar.",
                Computed: true,
            },
            "last_read_error": schema.StringAttribute{
                MarkdownDescription: "Why the last read failed - NotFound, AccessDenied, ServerError, Unreachable, NotAllowed, Timeout, TooLarge, TooManyRedirects, NotACalendar or TooComplex. Empty when it worked.",
                Optional: true,
                Computed: true,
            },
            "last_read_http_status": schema.NumberAttribute{
                MarkdownDescription: "The HTTP status the calendar's server answered the last failed read with, when it answered.",
                Optional: true,
                Computed: true,
            },
            "event_count": schema.NumberAttribute{
                MarkdownDescription: "How many of the calendar's events OneUptime keeps: those from 90 days back to a year ahead, up to 2,000.",
                Optional: true,
                Computed: true,
            },
            "has_more_events": schema.BoolAttribute{
                MarkdownDescription: "True when the calendar holds more events in the window than OneUptime keeps: the 2,000 nearest today are kept.",
                Optional: true,
                Computed: true,
            },
            "next_read_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime next reads the calendar.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *ScheduledMaintenanceExternalCalendarDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ScheduledMaintenanceExternalCalendarDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ScheduledMaintenanceExternalCalendarDataSourceModel

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
    if !data.CalendarHost.IsNull() && !data.CalendarHost.IsUnknown() {
        filters["calendarHost"] = data.CalendarHost.ValueString()
        filterNames = append(filterNames, "calendar_host = "+fmt.Sprintf("%q", data.CalendarHost.ValueString()))
    }
    if !data.IsEnabled.IsNull() && !data.IsEnabled.IsUnknown() {
        filters["isEnabled"] = data.IsEnabled.ValueBool()
        filterNames = append(filterNames, "is_enabled = "+fmt.Sprintf("%t", data.IsEnabled.ValueBool()))
    }
    if !data.RefreshIntervalInMinutes.IsNull() && !data.RefreshIntervalInMinutes.IsUnknown() {
        filters["refreshIntervalInMinutes"] = lookupNumber(data.RefreshIntervalInMinutes)
        filterNames = append(filterNames, "refresh_interval_in_minutes = "+data.RefreshIntervalInMinutes.ValueBigFloat().String())
    }
    if !data.LastReadError.IsNull() && !data.LastReadError.IsUnknown() {
        filters["lastReadError"] = data.LastReadError.ValueString()
        filterNames = append(filterNames, "last_read_error = "+fmt.Sprintf("%q", data.LastReadError.ValueString()))
    }
    if !data.LastReadHttpStatus.IsNull() && !data.LastReadHttpStatus.IsUnknown() {
        filters["lastReadHttpStatus"] = lookupNumber(data.LastReadHttpStatus)
        filterNames = append(filterNames, "last_read_http_status = "+data.LastReadHttpStatus.ValueBigFloat().String())
    }
    if !data.EventCount.IsNull() && !data.EventCount.IsUnknown() {
        filters["eventCount"] = lookupNumber(data.EventCount)
        filterNames = append(filterNames, "event_count = "+data.EventCount.ValueBigFloat().String())
    }
    if !data.HasMoreEvents.IsNull() && !data.HasMoreEvents.IsUnknown() {
        filters["hasMoreEvents"] = data.HasMoreEvents.ValueBool()
        filterNames = append(filterNames, "has_more_events = "+fmt.Sprintf("%t", data.HasMoreEvents.ValueBool()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the scheduled maintenance external calendar up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the scheduled maintenance external calendar up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "name": true,
        "calendarHost": true,
        "color": true,
        "isEnabled": true,
        "refreshIntervalInMinutes": true,
        "lastReadAt": true,
        "lastReadAttemptAt": true,
        "lastReadError": true,
        "lastReadHttpStatus": true,
        "eventCount": true,
        "hasMoreEvents": true,
        "nextReadAt": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/scheduled-maintenance-external-calendar/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read scheduled_maintenance_external_calendar, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No scheduled maintenance external calendar found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read scheduled_maintenance_external_calendar: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/scheduled-maintenance-external-calendar/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list scheduled_maintenance_external_calendar, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list scheduled_maintenance_external_calendar: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No scheduled maintenance external calendar matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one scheduled maintenance external calendar matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for scheduled_maintenance_external_calendar.")
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
    if obj, ok := item["calendarHost"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CalendarHost = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CalendarHost = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CalendarHost = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CalendarHost = types.StringValue(string(jsonBytes))
        } else {
            data.CalendarHost = types.StringNull()
        }
    } else if val, ok := item["calendarHost"].(string); ok {
        data.CalendarHost = types.StringValue(val)
    } else {
        data.CalendarHost = types.StringNull()
    }
    if obj, ok := item["color"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Color = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Color = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Color = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Color = types.StringValue(string(jsonBytes))
        } else {
            data.Color = types.StringNull()
        }
    } else if val, ok := item["color"].(string); ok {
        data.Color = types.StringValue(val)
    } else {
        data.Color = types.StringNull()
    }
    if val, ok := item["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    } else {
        data.IsEnabled = types.BoolNull()
    }
    data.RefreshIntervalInMinutes = numberFromAPI(data.RefreshIntervalInMinutes, item["refreshIntervalInMinutes"])
    if obj, ok := item["lastReadAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastReadAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastReadAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastReadAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastReadAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastReadAt = types.StringNull()
        }
    } else if val, ok := item["lastReadAt"].(string); ok {
        data.LastReadAt = types.StringValue(val)
    } else {
        data.LastReadAt = types.StringNull()
    }
    if obj, ok := item["lastReadAttemptAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastReadAttemptAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastReadAttemptAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastReadAttemptAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastReadAttemptAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastReadAttemptAt = types.StringNull()
        }
    } else if val, ok := item["lastReadAttemptAt"].(string); ok {
        data.LastReadAttemptAt = types.StringValue(val)
    } else {
        data.LastReadAttemptAt = types.StringNull()
    }
    if obj, ok := item["lastReadError"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastReadError = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastReadError = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastReadError = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastReadError = types.StringValue(string(jsonBytes))
        } else {
            data.LastReadError = types.StringNull()
        }
    } else if val, ok := item["lastReadError"].(string); ok {
        data.LastReadError = types.StringValue(val)
    } else {
        data.LastReadError = types.StringNull()
    }
    data.LastReadHttpStatus = numberFromAPI(data.LastReadHttpStatus, item["lastReadHttpStatus"])
    data.EventCount = numberFromAPI(data.EventCount, item["eventCount"])
    if val, ok := item["hasMoreEvents"].(bool); ok {
        data.HasMoreEvents = types.BoolValue(val)
    } else {
        data.HasMoreEvents = types.BoolNull()
    }
    if obj, ok := item["nextReadAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NextReadAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.NextReadAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.NextReadAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.NextReadAt = types.StringValue(string(jsonBytes))
        } else {
            data.NextReadAt = types.StringNull()
        }
    } else if val, ok := item["nextReadAt"].(string); ok {
        data.NextReadAt = types.StringValue(val)
    } else {
        data.NextReadAt = types.StringNull()
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
