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
var _ datasource.DataSource = &TeamComplianceSettingDataSource{}

func NewTeamComplianceSettingDataSource() datasource.DataSource {
    return &TeamComplianceSettingDataSource{}
}

// TeamComplianceSettingDataSource defines the data source implementation.
type TeamComplianceSettingDataSource struct {
    client *Client
}

// TeamComplianceSettingDataSourceModel describes the data source data model.
type TeamComplianceSettingDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    TeamId types.String `tfsdk:"team_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    RuleType types.String `tfsdk:"rule_type"`
    Enabled types.Bool `tfsdk:"enabled"`
    Options types.String `tfsdk:"options"`
    NotificationChannel types.String `tfsdk:"notification_channel"`
    NotificationChannels types.String `tfsdk:"notification_channels"`
    IncidentSeverities types.Set `tfsdk:"incident_severities"`
    AlertSeverities types.Set `tfsdk:"alert_severities"`
}

func (d *TeamComplianceSettingDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_team_compliance_setting"
}

func (d *TeamComplianceSettingDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Compliance settings for your OneUptime team Look up an existing team compliance setting by `id`, or by any of its other arguments (`created_by_user_id`, `enabled`, `notification_channel`, ...): each one set must match, and exactly one team compliance setting may match them all.",

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
                MarkdownDescription: "ID of Team this compliance setting belongs to. The ID of a `oneuptime_team`.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "rule_type": schema.StringAttribute{
                MarkdownDescription: "Type of compliance rule.",
                Optional: true,
                Computed: true,
            },
            "enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this compliance rule is enabled.",
                Optional: true,
                Computed: true,
            },
            "options": schema.StringAttribute{
                MarkdownDescription: "Additional options for this compliance rule. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "notification_channel": schema.StringAttribute{
                MarkdownDescription: "Deprecated: use notificationChannels. The first of the rule's notification channels, or empty when it accepts any channel. Sending this field without notificationChannels sets the rule to that one channel (Call, SMS, Push, Email, WhatsApp, Telegram, Slack, MicrosoftTeams or Webhook).",
                Optional: true,
                Computed: true,
            },
            "notification_channels": schema.StringAttribute{
                MarkdownDescription: "On-call rules only: the channels members must be notified on, as a list - each member needs a rule on every one of them (Call, SMS, Push, Email, WhatsApp, Telegram, Slack, MicrosoftTeams or Webhook). Leave empty to accept any channel. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "incident_severities": schema.SetAttribute{
                MarkdownDescription: "Incident and incident episode on-call rules only: the severities members must have a rule for. Leave empty to require every incident severity. IDs of `oneuptime_incident_severity` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "alert_severities": schema.SetAttribute{
                MarkdownDescription: "Alert and alert episode on-call rules only: the severities members must have a rule for. Leave empty to require every alert severity. IDs of `oneuptime_alert_severity` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
        },
    }
}

func (d *TeamComplianceSettingDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TeamComplianceSettingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data TeamComplianceSettingDataSourceModel

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
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.RuleType.IsNull() && !data.RuleType.IsUnknown() {
        filters["ruleType"] = data.RuleType.ValueString()
        filterNames = append(filterNames, "rule_type = "+fmt.Sprintf("%q", data.RuleType.ValueString()))
    }
    if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
        filters["enabled"] = data.Enabled.ValueBool()
        filterNames = append(filterNames, "enabled = "+fmt.Sprintf("%t", data.Enabled.ValueBool()))
    }
    if !data.NotificationChannel.IsNull() && !data.NotificationChannel.IsUnknown() {
        filters["notificationChannel"] = data.NotificationChannel.ValueString()
        filterNames = append(filterNames, "notification_channel = "+fmt.Sprintf("%q", data.NotificationChannel.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the team compliance setting up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the team compliance setting up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "teamId": true,
        "createdByUserId": true,
        "ruleType": true,
        "enabled": true,
        "options": true,
        "notificationChannel": true,
        "notificationChannels": true,
        "incidentSeverities": true,
        "alertSeverities": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/team-compliance-setting/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read team_compliance_setting, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No team compliance setting found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read team_compliance_setting: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/team-compliance-setting/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list team_compliance_setting, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list team_compliance_setting: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No team compliance setting matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one team compliance setting matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for team_compliance_setting.")
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
    } else if val, ok := item["teamId"].(string); ok {
        data.TeamId = types.StringValue(val)
    } else {
        data.TeamId = types.StringNull()
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
    if obj, ok := item["ruleType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RuleType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RuleType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RuleType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RuleType = types.StringValue(string(jsonBytes))
        } else {
            data.RuleType = types.StringNull()
        }
    } else if val, ok := item["ruleType"].(string); ok {
        data.RuleType = types.StringValue(val)
    } else {
        data.RuleType = types.StringNull()
    }
    if val, ok := item["enabled"].(bool); ok {
        data.Enabled = types.BoolValue(val)
    } else {
        data.Enabled = types.BoolNull()
    }
    if obj, ok := item["options"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Options = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Options = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Options = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Options = types.StringValue(string(jsonBytes))
        } else {
            data.Options = types.StringNull()
        }
    } else if val, ok := item["options"].(string); ok {
        data.Options = types.StringValue(val)
    } else {
        data.Options = types.StringNull()
    }
    if obj, ok := item["notificationChannel"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NotificationChannel = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.NotificationChannel = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.NotificationChannel = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.NotificationChannel = types.StringValue(string(jsonBytes))
        } else {
            data.NotificationChannel = types.StringNull()
        }
    } else if val, ok := item["notificationChannel"].(string); ok {
        data.NotificationChannel = types.StringValue(val)
    } else {
        data.NotificationChannel = types.StringNull()
    }
    if obj, ok := item["notificationChannels"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NotificationChannels = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.NotificationChannels = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.NotificationChannels = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.NotificationChannels = types.StringValue(string(jsonBytes))
        } else {
            data.NotificationChannels = types.StringNull()
        }
    } else if val, ok := item["notificationChannels"].(string); ok {
        data.NotificationChannels = types.StringValue(val)
    } else {
        data.NotificationChannels = types.StringNull()
    }
    if val, ok := item["incidentSeverities"].([]interface{}); ok {
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
        data.IncidentSeverities = types.SetValueMust(types.StringType, setItems)
    } else {
        data.IncidentSeverities = types.SetNull(types.StringType)
    }
    if val, ok := item["alertSeverities"].([]interface{}); ok {
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
        data.AlertSeverities = types.SetValueMust(types.StringType, setItems)
    } else {
        data.AlertSeverities = types.SetNull(types.StringType)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
