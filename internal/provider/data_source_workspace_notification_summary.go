package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "math/big"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &WorkspaceNotificationSummaryDataSource{}

func NewWorkspaceNotificationSummaryDataSource() datasource.DataSource {
    return &WorkspaceNotificationSummaryDataSource{}
}

// WorkspaceNotificationSummaryDataSource defines the data source implementation.
type WorkspaceNotificationSummaryDataSource struct {
    client *Client
}

// WorkspaceNotificationSummaryDataSourceModel describes the data source data model.
type WorkspaceNotificationSummaryDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Description types.String `tfsdk:"description"`
    WorkspaceType types.String `tfsdk:"workspace_type"`
    SummaryType types.String `tfsdk:"summary_type"`
    RecurringInterval types.String `tfsdk:"recurring_interval"`
    NumberOfDaysOfData types.Number `tfsdk:"number_of_days_of_data"`
    SendFirstReportAt types.String `tfsdk:"send_first_report_at"`
    ChannelNames types.String `tfsdk:"channel_names"`
    TeamName types.String `tfsdk:"team_name"`
    SummaryItems types.String `tfsdk:"summary_items"`
    Filters types.String `tfsdk:"filters"`
    FilterCondition types.String `tfsdk:"filter_condition"`
    NextSendAt types.String `tfsdk:"next_send_at"`
    LastSentAt types.String `tfsdk:"last_sent_at"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
}

func (d *WorkspaceNotificationSummaryDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_workspace_notification_summary"
}

func (d *WorkspaceNotificationSummaryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Recurring summary reports for incidents and alerts sent to Slack or Microsoft Teams Look up an existing workspace_notification_summary by `id` or by `name`.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Exactly one of `id` or `name` must be set.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.",
                Optional: true,
                Computed: true,
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "deleted_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "version": schema.NumberAttribute{
                MarkdownDescription: "Object version",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of the Summary Rule.",
                Computed: true,
            },
            "workspace_type": schema.StringAttribute{
                MarkdownDescription: "Type of Workspace - Slack, Microsoft Teams, etc..",
                Computed: true,
            },
            "summary_type": schema.StringAttribute{
                MarkdownDescription: "Type of summary - Incident, Alert, Incident Episode, or Alert Episode.",
                Computed: true,
            },
            "recurring_interval": schema.StringAttribute{
                MarkdownDescription: "How often should the summary be sent?.",
                Computed: true,
            },
            "number_of_days_of_data": schema.NumberAttribute{
                MarkdownDescription: "How many days of data to include in the summary.",
                Computed: true,
            },
            "send_first_report_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "channel_names": schema.StringAttribute{
                MarkdownDescription: "List of channel names to post the summary to.",
                Computed: true,
            },
            "team_name": schema.StringAttribute{
                MarkdownDescription: "Microsoft Teams team name (only for Microsoft Teams).",
                Computed: true,
            },
            "summary_items": schema.StringAttribute{
                MarkdownDescription: "Checklist of items to include in the summary.",
                Computed: true,
            },
            "filters": schema.StringAttribute{
                MarkdownDescription: "Filter conditions for which items to include in the summary.",
                Computed: true,
            },
            "filter_condition": schema.StringAttribute{
                MarkdownDescription: "How to combine filters - Any or All.",
                Computed: true,
            },
            "next_send_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "last_sent_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Is this summary rule enabled?.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "deleted_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (d *WorkspaceNotificationSummaryDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *WorkspaceNotificationSummaryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data WorkspaceNotificationSummaryDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && data.Id.ValueString() != ""
    hasName := !data.Name.IsNull() && data.Name.ValueString() != ""
    if hasId == hasName {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Exactly one of `id` or `name` must be set to look up a workspace_notification_summary.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "name": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "projectId": true,
        "description": true,
        "workspaceType": true,
        "summaryType": true,
        "recurringInterval": true,
        "numberOfDaysOfData": true,
        "sendFirstReportAt": true,
        "channelNames": true,
        "teamName": true,
        "summaryItems": true,
        "filters": true,
        "filterCondition": true,
        "nextSendAt": true,
        "lastSentAt": true,
        "isEnabled": true,
        "createdByUserId": true,
        "deletedByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/workspace-notification-summary/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read workspace_notification_summary, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No workspace_notification_summary found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read workspace_notification_summary: %s", err))
            return
        }
        if wrapper, ok := itemResponse["data"].(map[string]interface{}); ok {
            item = wrapper
        } else {
            item = itemResponse
        }
    } else {
        listBody := map[string]interface{}{
            "query": map[string]interface{}{
                "name": data.Name.ValueString(),
            },
            "select": selectParam,
            // limit 2 is enough to detect ambiguity without paging.
            "limit": 2,
        }
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/workspace-notification-summary/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list workspace_notification_summary, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list workspace_notification_summary: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No workspace_notification_summary found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one workspace_notification_summary matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for workspace_notification_summary.")
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
    if obj, ok := item["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
        }
    } else if val, ok := item["deletedAt"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    } else {
        data.DeletedAt = types.StringNull()
    }
    if val, ok := item["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["version"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        data.Version = types.NumberNull()
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
    if obj, ok := item["workspaceType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.WorkspaceType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.WorkspaceType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.WorkspaceType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.WorkspaceType = types.StringValue(string(jsonBytes))
        } else {
            data.WorkspaceType = types.StringNull()
        }
    } else if val, ok := item["workspaceType"].(string); ok {
        data.WorkspaceType = types.StringValue(val)
    } else {
        data.WorkspaceType = types.StringNull()
    }
    if obj, ok := item["summaryType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SummaryType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SummaryType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SummaryType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SummaryType = types.StringValue(string(jsonBytes))
        } else {
            data.SummaryType = types.StringNull()
        }
    } else if val, ok := item["summaryType"].(string); ok {
        data.SummaryType = types.StringValue(val)
    } else {
        data.SummaryType = types.StringNull()
    }
    if obj, ok := item["recurringInterval"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RecurringInterval = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RecurringInterval = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RecurringInterval = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RecurringInterval = types.StringValue(string(jsonBytes))
        } else {
            data.RecurringInterval = types.StringNull()
        }
    } else if val, ok := item["recurringInterval"].(string); ok {
        data.RecurringInterval = types.StringValue(val)
    } else {
        data.RecurringInterval = types.StringNull()
    }
    if val, ok := item["numberOfDaysOfData"].(float64); ok {
        data.NumberOfDaysOfData = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["numberOfDaysOfData"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.NumberOfDaysOfData = types.NumberValue(big.NewFloat(val))
        } else {
            data.NumberOfDaysOfData = types.NumberNull()
        }
    } else {
        data.NumberOfDaysOfData = types.NumberNull()
    }
    if obj, ok := item["sendFirstReportAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SendFirstReportAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SendFirstReportAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SendFirstReportAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SendFirstReportAt = types.StringValue(string(jsonBytes))
        } else {
            data.SendFirstReportAt = types.StringNull()
        }
    } else if val, ok := item["sendFirstReportAt"].(string); ok {
        data.SendFirstReportAt = types.StringValue(val)
    } else {
        data.SendFirstReportAt = types.StringNull()
    }
    if obj, ok := item["channelNames"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ChannelNames = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ChannelNames = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ChannelNames = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ChannelNames = types.StringValue(string(jsonBytes))
        } else {
            data.ChannelNames = types.StringNull()
        }
    } else if val, ok := item["channelNames"].(string); ok {
        data.ChannelNames = types.StringValue(val)
    } else {
        data.ChannelNames = types.StringNull()
    }
    if obj, ok := item["teamName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TeamName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TeamName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TeamName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TeamName = types.StringValue(string(jsonBytes))
        } else {
            data.TeamName = types.StringNull()
        }
    } else if val, ok := item["teamName"].(string); ok {
        data.TeamName = types.StringValue(val)
    } else {
        data.TeamName = types.StringNull()
    }
    if obj, ok := item["summaryItems"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SummaryItems = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SummaryItems = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SummaryItems = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SummaryItems = types.StringValue(string(jsonBytes))
        } else {
            data.SummaryItems = types.StringNull()
        }
    } else if val, ok := item["summaryItems"].(string); ok {
        data.SummaryItems = types.StringValue(val)
    } else {
        data.SummaryItems = types.StringNull()
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
    if obj, ok := item["filterCondition"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FilterCondition = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FilterCondition = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FilterCondition = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FilterCondition = types.StringValue(string(jsonBytes))
        } else {
            data.FilterCondition = types.StringNull()
        }
    } else if val, ok := item["filterCondition"].(string); ok {
        data.FilterCondition = types.StringValue(val)
    } else {
        data.FilterCondition = types.StringNull()
    }
    if obj, ok := item["nextSendAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NextSendAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.NextSendAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.NextSendAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.NextSendAt = types.StringValue(string(jsonBytes))
        } else {
            data.NextSendAt = types.StringNull()
        }
    } else if val, ok := item["nextSendAt"].(string); ok {
        data.NextSendAt = types.StringValue(val)
    } else {
        data.NextSendAt = types.StringNull()
    }
    if obj, ok := item["lastSentAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastSentAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastSentAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastSentAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastSentAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastSentAt = types.StringNull()
        }
    } else if val, ok := item["lastSentAt"].(string); ok {
        data.LastSentAt = types.StringValue(val)
    } else {
        data.LastSentAt = types.StringNull()
    }
    if val, ok := item["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    } else {
        data.IsEnabled = types.BoolNull()
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
    if obj, ok := item["deletedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeletedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedByUserId = types.StringNull()
        }
    } else if val, ok := item["deletedByUserId"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    } else {
        data.DeletedByUserId = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
