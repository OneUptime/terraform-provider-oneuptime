package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "math/big"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &MonitorDataSource{}

func NewMonitorDataSource() datasource.DataSource {
    return &MonitorDataSource{}
}

// MonitorDataSource defines the data source implementation.
type MonitorDataSource struct {
    client *Client
}

// MonitorDataSourceModel describes the data source data model.
type MonitorDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Description types.String `tfsdk:"description"`
    Slug types.String `tfsdk:"slug"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    Labels types.Set `tfsdk:"labels"`
    DependsOnMonitors types.Set `tfsdk:"depends_on_monitors"`
    SuppressAlertsWhenParentMonitorStatuses types.Set `tfsdk:"suppress_alerts_when_parent_monitor_statuses"`
    MonitorTemplateId types.String `tfsdk:"monitor_template_id"`
    MonitorType types.String `tfsdk:"monitor_type"`
    CurrentMonitorStatusId types.String `tfsdk:"current_monitor_status_id"`
    MonitorSteps types.String `tfsdk:"monitor_steps"`
    MonitoringInterval types.String `tfsdk:"monitoring_interval"`
    CustomFields types.String `tfsdk:"custom_fields"`
    IsOwnerNotifiedOfResourceCreation types.Bool `tfsdk:"is_owner_notified_of_resource_creation"`
    DisableActiveMonitoring types.Bool `tfsdk:"disable_active_monitoring"`
    IncomingRequestMonitorHeartbeatCheckedAt types.String `tfsdk:"incoming_request_monitor_heartbeat_checked_at"`
    TelemetryMonitorNextMonitorAt types.String `tfsdk:"telemetry_monitor_next_monitor_at"`
    TelemetryMonitorLastMonitorAt types.String `tfsdk:"telemetry_monitor_last_monitor_at"`
    DisableActiveMonitoringBecauseOfScheduledMaintenanceEvent types.Bool `tfsdk:"disable_active_monitoring_because_of_scheduled_maintenance_event"`
    DisableActiveMonitoringBecauseOfManualIncident types.Bool `tfsdk:"disable_active_monitoring_because_of_manual_incident"`
    ServerMonitorRequestReceivedAt types.String `tfsdk:"server_monitor_request_received_at"`
    ServerMonitorSecretKey types.String `tfsdk:"server_monitor_secret_key"`
    IncomingRequestSecretKey types.String `tfsdk:"incoming_request_secret_key"`
    IncomingMonitorRequest types.String `tfsdk:"incoming_monitor_request"`
    IncomingEmailSecretKey types.String `tfsdk:"incoming_email_secret_key"`
    IncomingEmailMonitorLastEmailReceivedAt types.String `tfsdk:"incoming_email_monitor_last_email_received_at"`
    IncomingEmailMonitorRequest types.String `tfsdk:"incoming_email_monitor_request"`
    IncomingEmailMonitorHeartbeatCheckedAt types.String `tfsdk:"incoming_email_monitor_heartbeat_checked_at"`
    ServerMonitorResponse types.String `tfsdk:"server_monitor_response"`
    IsAllProbesDisconnectedFromThisMonitor types.Bool `tfsdk:"is_all_probes_disconnected_from_this_monitor"`
    IsNoProbeEnabledOnThisMonitor types.Bool `tfsdk:"is_no_probe_enabled_on_this_monitor"`
    MinimumProbeAgreement types.Number `tfsdk:"minimum_probe_agreement"`
}

func (d *MonitorDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_monitor"
}

func (d *MonitorDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Monitor is anything that monitors your API, Websites, IP, Network or more. You can also create static monitor that does not monitor anything. Look up an existing monitor by `id` or by `name`.",

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
                MarkdownDescription: "Friendly description that will help you remember.",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in..",
                Computed: true,
                ElementType: types.StringType,
            },
            "depends_on_monitors": schema.SetAttribute{
                MarkdownDescription: "Parent monitors this monitor depends on. When a parent is offline (or in one of the configured suppression statuses), alerts and incidents from this monitor are suppressed at creation time — the monitor keeps evaluating and its status timeline still updates..",
                Computed: true,
                ElementType: types.StringType,
            },
            "suppress_alerts_when_parent_monitor_statuses": schema.SetAttribute{
                MarkdownDescription: "Parent monitor statuses that suppress this monitor's alerts and incidents. When empty, statuses flagged as offline suppress (the default). Only used when Depends On Monitors is set..",
                Computed: true,
                ElementType: types.StringType,
            },
            "monitor_template_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "monitor_type": schema.StringAttribute{
                MarkdownDescription: "What is the type of this monitor? Website? API? etc..",
                Computed: true,
            },
            "current_monitor_status_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "monitor_steps": schema.StringAttribute{
                MarkdownDescription: "MonitorSteps object",
                Computed: true,
            },
            "monitoring_interval": schema.StringAttribute{
                MarkdownDescription: "How often would you like OneUptime to monitor this resource?.",
                Computed: true,
            },
            "custom_fields": schema.StringAttribute{
                MarkdownDescription: "Custom Fields on this resource..",
                Computed: true,
            },
            "is_owner_notified_of_resource_creation": schema.BoolAttribute{
                MarkdownDescription: "Are owners notified of when this resource is created?.",
                Computed: true,
            },
            "disable_active_monitoring": schema.BoolAttribute{
                MarkdownDescription: "Disable active monitoring for this resource?.",
                Computed: true,
            },
            "incoming_request_monitor_heartbeat_checked_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "telemetry_monitor_next_monitor_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "telemetry_monitor_last_monitor_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "disable_active_monitoring_because_of_scheduled_maintenance_event": schema.BoolAttribute{
                MarkdownDescription: "Disable Monitoring because of Ongoing Scheduled Maintenance Event.",
                Computed: true,
            },
            "disable_active_monitoring_because_of_manual_incident": schema.BoolAttribute{
                MarkdownDescription: "Disable Monitoring because of Incident which is creeated manually by user..",
                Computed: true,
            },
            "server_monitor_request_received_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "server_monitor_secret_key": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "incoming_request_secret_key": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "incoming_monitor_request": schema.StringAttribute{
                MarkdownDescription: "Incoming Monitor Request for Incoming Request Monitor.",
                Computed: true,
            },
            "incoming_email_secret_key": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "incoming_email_monitor_last_email_received_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "incoming_email_monitor_request": schema.StringAttribute{
                MarkdownDescription: "This field is for Incoming Email Monitor only. Last email data received..",
                Computed: true,
            },
            "incoming_email_monitor_heartbeat_checked_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "server_monitor_response": schema.StringAttribute{
                MarkdownDescription: "Server Monitor Response for Server Monitor.",
                Computed: true,
            },
            "is_all_probes_disconnected_from_this_monitor": schema.BoolAttribute{
                MarkdownDescription: "All Probes Disconnected From This Monitor. Is this monitor not being monitored?.",
                Computed: true,
            },
            "is_no_probe_enabled_on_this_monitor": schema.BoolAttribute{
                MarkdownDescription: "No Probe Enabled On This Monitor. Is this monitor not being monitored?.",
                Computed: true,
            },
            "minimum_probe_agreement": schema.NumberAttribute{
                MarkdownDescription: "Minimum number of probes that must agree on a status before the monitor status changes. If null, all enabled and connected probes must agree..",
                Computed: true,
            },
        },
    }
}

func (d *MonitorDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MonitorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data MonitorDataSourceModel

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
            "Exactly one of `id` or `name` must be set to look up a monitor.",
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
        "slug": true,
        "createdByUserId": true,
        "labels": true,
        "dependsOnMonitors": true,
        "suppressAlertsWhenParentMonitorStatuses": true,
        "monitorTemplateId": true,
        "monitorType": true,
        "currentMonitorStatusId": true,
        "monitorSteps": true,
        "monitoringInterval": true,
        "customFields": true,
        "isOwnerNotifiedOfResourceCreation": true,
        "disableActiveMonitoring": true,
        "incomingRequestMonitorHeartbeatCheckedAt": true,
        "telemetryMonitorNextMonitorAt": true,
        "telemetryMonitorLastMonitorAt": true,
        "disableActiveMonitoringBecauseOfScheduledMaintenanceEvent": true,
        "disableActiveMonitoringBecauseOfManualIncident": true,
        "serverMonitorRequestReceivedAt": true,
        "serverMonitorSecretKey": true,
        "incomingRequestSecretKey": true,
        "incomingMonitorRequest": true,
        "incomingEmailSecretKey": true,
        "incomingEmailMonitorLastEmailReceivedAt": true,
        "incomingEmailMonitorRequest": true,
        "incomingEmailMonitorHeartbeatCheckedAt": true,
        "serverMonitorResponse": true,
        "isAllProbesDisconnectedFromThisMonitor": true,
        "isNoProbeEnabledOnThisMonitor": true,
        "minimumProbeAgreement": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/monitor/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read monitor, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No monitor found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read monitor: %s", err))
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
        httpResp, err := d.client.Post(ctx, "/monitor/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list monitor, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list monitor: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No monitor found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one monitor matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for monitor.")
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
    if val, ok := item["dependsOnMonitors"].([]interface{}); ok {
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
        data.DependsOnMonitors = types.SetValueMust(types.StringType, setItems)
    } else {
        data.DependsOnMonitors = types.SetNull(types.StringType)
    }
    if val, ok := item["suppressAlertsWhenParentMonitorStatuses"].([]interface{}); ok {
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
        data.SuppressAlertsWhenParentMonitorStatuses = types.SetValueMust(types.StringType, setItems)
    } else {
        data.SuppressAlertsWhenParentMonitorStatuses = types.SetNull(types.StringType)
    }
    if obj, ok := item["monitorTemplateId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorTemplateId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MonitorTemplateId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MonitorTemplateId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MonitorTemplateId = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorTemplateId = types.StringNull()
        }
    } else if val, ok := item["monitorTemplateId"].(string); ok {
        data.MonitorTemplateId = types.StringValue(val)
    } else {
        data.MonitorTemplateId = types.StringNull()
    }
    if obj, ok := item["monitorType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MonitorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MonitorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MonitorType = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorType = types.StringNull()
        }
    } else if val, ok := item["monitorType"].(string); ok {
        data.MonitorType = types.StringValue(val)
    } else {
        data.MonitorType = types.StringNull()
    }
    if obj, ok := item["currentMonitorStatusId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CurrentMonitorStatusId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CurrentMonitorStatusId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CurrentMonitorStatusId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CurrentMonitorStatusId = types.StringValue(string(jsonBytes))
        } else {
            data.CurrentMonitorStatusId = types.StringNull()
        }
    } else if val, ok := item["currentMonitorStatusId"].(string); ok {
        data.CurrentMonitorStatusId = types.StringValue(val)
    } else {
        data.CurrentMonitorStatusId = types.StringNull()
    }
    if obj, ok := item["monitorSteps"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorSteps = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MonitorSteps = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MonitorSteps = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MonitorSteps = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorSteps = types.StringNull()
        }
    } else if val, ok := item["monitorSteps"].(string); ok {
        data.MonitorSteps = types.StringValue(val)
    } else {
        data.MonitorSteps = types.StringNull()
    }
    if obj, ok := item["monitoringInterval"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitoringInterval = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MonitoringInterval = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MonitoringInterval = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MonitoringInterval = types.StringValue(string(jsonBytes))
        } else {
            data.MonitoringInterval = types.StringNull()
        }
    } else if val, ok := item["monitoringInterval"].(string); ok {
        data.MonitoringInterval = types.StringValue(val)
    } else {
        data.MonitoringInterval = types.StringNull()
    }
    if obj, ok := item["customFields"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CustomFields = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CustomFields = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CustomFields = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CustomFields = types.StringValue(string(jsonBytes))
        } else {
            data.CustomFields = types.StringNull()
        }
    } else if val, ok := item["customFields"].(string); ok {
        data.CustomFields = types.StringValue(val)
    } else {
        data.CustomFields = types.StringNull()
    }
    if val, ok := item["isOwnerNotifiedOfResourceCreation"].(bool); ok {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolValue(val)
    } else {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolNull()
    }
    if val, ok := item["disableActiveMonitoring"].(bool); ok {
        data.DisableActiveMonitoring = types.BoolValue(val)
    } else {
        data.DisableActiveMonitoring = types.BoolNull()
    }
    if obj, ok := item["incomingRequestMonitorHeartbeatCheckedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncomingRequestMonitorHeartbeatCheckedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncomingRequestMonitorHeartbeatCheckedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncomingRequestMonitorHeartbeatCheckedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncomingRequestMonitorHeartbeatCheckedAt = types.StringValue(string(jsonBytes))
        } else {
            data.IncomingRequestMonitorHeartbeatCheckedAt = types.StringNull()
        }
    } else if val, ok := item["incomingRequestMonitorHeartbeatCheckedAt"].(string); ok {
        data.IncomingRequestMonitorHeartbeatCheckedAt = types.StringValue(val)
    } else {
        data.IncomingRequestMonitorHeartbeatCheckedAt = types.StringNull()
    }
    if obj, ok := item["telemetryMonitorNextMonitorAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TelemetryMonitorNextMonitorAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TelemetryMonitorNextMonitorAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TelemetryMonitorNextMonitorAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TelemetryMonitorNextMonitorAt = types.StringValue(string(jsonBytes))
        } else {
            data.TelemetryMonitorNextMonitorAt = types.StringNull()
        }
    } else if val, ok := item["telemetryMonitorNextMonitorAt"].(string); ok {
        data.TelemetryMonitorNextMonitorAt = types.StringValue(val)
    } else {
        data.TelemetryMonitorNextMonitorAt = types.StringNull()
    }
    if obj, ok := item["telemetryMonitorLastMonitorAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TelemetryMonitorLastMonitorAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TelemetryMonitorLastMonitorAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TelemetryMonitorLastMonitorAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TelemetryMonitorLastMonitorAt = types.StringValue(string(jsonBytes))
        } else {
            data.TelemetryMonitorLastMonitorAt = types.StringNull()
        }
    } else if val, ok := item["telemetryMonitorLastMonitorAt"].(string); ok {
        data.TelemetryMonitorLastMonitorAt = types.StringValue(val)
    } else {
        data.TelemetryMonitorLastMonitorAt = types.StringNull()
    }
    if val, ok := item["disableActiveMonitoringBecauseOfScheduledMaintenanceEvent"].(bool); ok {
        data.DisableActiveMonitoringBecauseOfScheduledMaintenanceEvent = types.BoolValue(val)
    } else {
        data.DisableActiveMonitoringBecauseOfScheduledMaintenanceEvent = types.BoolNull()
    }
    if val, ok := item["disableActiveMonitoringBecauseOfManualIncident"].(bool); ok {
        data.DisableActiveMonitoringBecauseOfManualIncident = types.BoolValue(val)
    } else {
        data.DisableActiveMonitoringBecauseOfManualIncident = types.BoolNull()
    }
    if obj, ok := item["serverMonitorRequestReceivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServerMonitorRequestReceivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ServerMonitorRequestReceivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ServerMonitorRequestReceivedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ServerMonitorRequestReceivedAt = types.StringValue(string(jsonBytes))
        } else {
            data.ServerMonitorRequestReceivedAt = types.StringNull()
        }
    } else if val, ok := item["serverMonitorRequestReceivedAt"].(string); ok {
        data.ServerMonitorRequestReceivedAt = types.StringValue(val)
    } else {
        data.ServerMonitorRequestReceivedAt = types.StringNull()
    }
    if obj, ok := item["serverMonitorSecretKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServerMonitorSecretKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ServerMonitorSecretKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ServerMonitorSecretKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ServerMonitorSecretKey = types.StringValue(string(jsonBytes))
        } else {
            data.ServerMonitorSecretKey = types.StringNull()
        }
    } else if val, ok := item["serverMonitorSecretKey"].(string); ok {
        data.ServerMonitorSecretKey = types.StringValue(val)
    } else {
        data.ServerMonitorSecretKey = types.StringNull()
    }
    if obj, ok := item["incomingRequestSecretKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncomingRequestSecretKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncomingRequestSecretKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncomingRequestSecretKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncomingRequestSecretKey = types.StringValue(string(jsonBytes))
        } else {
            data.IncomingRequestSecretKey = types.StringNull()
        }
    } else if val, ok := item["incomingRequestSecretKey"].(string); ok {
        data.IncomingRequestSecretKey = types.StringValue(val)
    } else {
        data.IncomingRequestSecretKey = types.StringNull()
    }
    if obj, ok := item["incomingMonitorRequest"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncomingMonitorRequest = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncomingMonitorRequest = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncomingMonitorRequest = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncomingMonitorRequest = types.StringValue(string(jsonBytes))
        } else {
            data.IncomingMonitorRequest = types.StringNull()
        }
    } else if val, ok := item["incomingMonitorRequest"].(string); ok {
        data.IncomingMonitorRequest = types.StringValue(val)
    } else {
        data.IncomingMonitorRequest = types.StringNull()
    }
    if obj, ok := item["incomingEmailSecretKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncomingEmailSecretKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncomingEmailSecretKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncomingEmailSecretKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncomingEmailSecretKey = types.StringValue(string(jsonBytes))
        } else {
            data.IncomingEmailSecretKey = types.StringNull()
        }
    } else if val, ok := item["incomingEmailSecretKey"].(string); ok {
        data.IncomingEmailSecretKey = types.StringValue(val)
    } else {
        data.IncomingEmailSecretKey = types.StringNull()
    }
    if obj, ok := item["incomingEmailMonitorLastEmailReceivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncomingEmailMonitorLastEmailReceivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncomingEmailMonitorLastEmailReceivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncomingEmailMonitorLastEmailReceivedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncomingEmailMonitorLastEmailReceivedAt = types.StringValue(string(jsonBytes))
        } else {
            data.IncomingEmailMonitorLastEmailReceivedAt = types.StringNull()
        }
    } else if val, ok := item["incomingEmailMonitorLastEmailReceivedAt"].(string); ok {
        data.IncomingEmailMonitorLastEmailReceivedAt = types.StringValue(val)
    } else {
        data.IncomingEmailMonitorLastEmailReceivedAt = types.StringNull()
    }
    if obj, ok := item["incomingEmailMonitorRequest"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncomingEmailMonitorRequest = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncomingEmailMonitorRequest = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncomingEmailMonitorRequest = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncomingEmailMonitorRequest = types.StringValue(string(jsonBytes))
        } else {
            data.IncomingEmailMonitorRequest = types.StringNull()
        }
    } else if val, ok := item["incomingEmailMonitorRequest"].(string); ok {
        data.IncomingEmailMonitorRequest = types.StringValue(val)
    } else {
        data.IncomingEmailMonitorRequest = types.StringNull()
    }
    if obj, ok := item["incomingEmailMonitorHeartbeatCheckedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncomingEmailMonitorHeartbeatCheckedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncomingEmailMonitorHeartbeatCheckedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncomingEmailMonitorHeartbeatCheckedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncomingEmailMonitorHeartbeatCheckedAt = types.StringValue(string(jsonBytes))
        } else {
            data.IncomingEmailMonitorHeartbeatCheckedAt = types.StringNull()
        }
    } else if val, ok := item["incomingEmailMonitorHeartbeatCheckedAt"].(string); ok {
        data.IncomingEmailMonitorHeartbeatCheckedAt = types.StringValue(val)
    } else {
        data.IncomingEmailMonitorHeartbeatCheckedAt = types.StringNull()
    }
    if obj, ok := item["serverMonitorResponse"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServerMonitorResponse = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ServerMonitorResponse = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ServerMonitorResponse = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ServerMonitorResponse = types.StringValue(string(jsonBytes))
        } else {
            data.ServerMonitorResponse = types.StringNull()
        }
    } else if val, ok := item["serverMonitorResponse"].(string); ok {
        data.ServerMonitorResponse = types.StringValue(val)
    } else {
        data.ServerMonitorResponse = types.StringNull()
    }
    if val, ok := item["isAllProbesDisconnectedFromThisMonitor"].(bool); ok {
        data.IsAllProbesDisconnectedFromThisMonitor = types.BoolValue(val)
    } else {
        data.IsAllProbesDisconnectedFromThisMonitor = types.BoolNull()
    }
    if val, ok := item["isNoProbeEnabledOnThisMonitor"].(bool); ok {
        data.IsNoProbeEnabledOnThisMonitor = types.BoolValue(val)
    } else {
        data.IsNoProbeEnabledOnThisMonitor = types.BoolNull()
    }
    if val, ok := item["minimumProbeAgreement"].(float64); ok {
        data.MinimumProbeAgreement = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["minimumProbeAgreement"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.MinimumProbeAgreement = types.NumberValue(big.NewFloat(val))
        } else {
            data.MinimumProbeAgreement = types.NumberNull()
        }
    } else {
        data.MinimumProbeAgreement = types.NumberNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
