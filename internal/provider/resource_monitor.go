package provider

import (
    "context"
    "fmt"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
    "math/big"
    "net/http"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &MonitorResource{}
var _ resource.ResourceWithImportState = &MonitorResource{}

func NewMonitorResource() resource.Resource {
    return &MonitorResource{}
}

// MonitorResource defines the resource implementation.
type MonitorResource struct {
    client *Client
}

// MonitorResourceModel describes the resource data model.
type MonitorResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    Labels types.List `tfsdk:"labels"`
    MonitorType types.String `tfsdk:"monitor_type"`
    CurrentMonitorStatusId types.String `tfsdk:"current_monitor_status_id"`
    MonitorSteps types.Map `tfsdk:"monitor_steps"`
    MonitoringInterval types.String `tfsdk:"monitoring_interval"`
    CustomFields types.Map `tfsdk:"custom_fields"`
    IsOwnerNotifiedOfResourceCreation types.Bool `tfsdk:"is_owner_notified_of_resource_creation"`
    DisableActiveMonitoring types.Bool `tfsdk:"disable_active_monitoring"`
    IncomingRequestMonitorHeartbeatCheckedAt types.Map `tfsdk:"incoming_request_monitor_heartbeat_checked_at"`
    TelemetryMonitorNextMonitorAt types.Map `tfsdk:"telemetry_monitor_next_monitor_at"`
    TelemetryMonitorLastMonitorAt types.Map `tfsdk:"telemetry_monitor_last_monitor_at"`
    ServerMonitorRequestReceivedAt types.Map `tfsdk:"server_monitor_request_received_at"`
    ServerMonitorSecretKey types.String `tfsdk:"server_monitor_secret_key"`
    IncomingRequestSecretKey types.String `tfsdk:"incoming_request_secret_key"`
    IncomingMonitorRequest types.Map `tfsdk:"incoming_monitor_request"`
    ServerMonitorResponse types.Map `tfsdk:"server_monitor_response"`
    CreatedAt types.Map `tfsdk:"created_at"`
    UpdatedAt types.Map `tfsdk:"updated_at"`
    DeletedAt types.Map `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    Slug types.String `tfsdk:"slug"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DisableActiveMonitoringBecauseOfScheduledMaintenanceEvent types.Bool `tfsdk:"disable_active_monitoring_because_of_scheduled_maintenance_event"`
    DisableActiveMonitoringBecauseOfManualIncident types.Bool `tfsdk:"disable_active_monitoring_because_of_manual_incident"`
    IsAllProbesDisconnectedFromThisMonitor types.Bool `tfsdk:"is_all_probes_disconnected_from_this_monitor"`
    IsNoProbeEnabledOnThisMonitor types.Bool `tfsdk:"is_no_probe_enabled_on_this_monitor"`
}

func (r *MonitorResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_monitor"
}

func (r *MonitorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "monitor resource",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Required: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name",
                Required: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description",
                Optional: true,
            },
            "labels": schema.ListAttribute{
                MarkdownDescription: "Labels",
                Optional: true,
                ElementType: types.StringType,
            },
            "monitor_type": schema.StringAttribute{
                MarkdownDescription: "Monitor Type",
                Required: true,
            },
            "current_monitor_status_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Required: true,
            },
            "monitor_steps": schema.MapAttribute{
                MarkdownDescription: "MonitorSteps object",
                Optional: true,
                ElementType: types.StringType,
            },
            "monitoring_interval": schema.StringAttribute{
                MarkdownDescription: "Monitoring Interval",
                Optional: true,
            },
            "custom_fields": schema.MapAttribute{
                MarkdownDescription: "Custom Fields",
                Optional: true,
                ElementType: types.StringType,
            },
            "is_owner_notified_of_resource_creation": schema.BoolAttribute{
                MarkdownDescription: "Are Owners Notified Of Resource Creation?",
                Required: true,
            },
            "disable_active_monitoring": schema.BoolAttribute{
                MarkdownDescription: "Disable Monitoring",
                Required: true,
            },
            "incoming_request_monitor_heartbeat_checked_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Optional: true,
                ElementType: types.StringType,
            },
            "telemetry_monitor_next_monitor_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Optional: true,
                ElementType: types.StringType,
            },
            "telemetry_monitor_last_monitor_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Optional: true,
                ElementType: types.StringType,
            },
            "server_monitor_request_received_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Optional: true,
                ElementType: types.StringType,
            },
            "server_monitor_secret_key": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
            },
            "incoming_request_secret_key": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
            },
            "incoming_monitor_request": schema.MapAttribute{
                MarkdownDescription: "Incoming Monitor Request",
                Optional: true,
                ElementType: types.StringType,
            },
            "server_monitor_response": schema.MapAttribute{
                MarkdownDescription: "Server Monitor Response",
                Optional: true,
                ElementType: types.StringType,
            },
            "created_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
                ElementType: types.StringType,
            },
            "updated_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
                ElementType: types.StringType,
            },
            "deleted_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
                ElementType: types.StringType,
            },
            "version": schema.NumberAttribute{
                MarkdownDescription: "Version",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Monitor], Read: [Project Owner, Project Admin, Project Member, Read Monitor], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "disable_active_monitoring_because_of_scheduled_maintenance_event": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Monitor], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "disable_active_monitoring_because_of_manual_incident": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Monitor], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "is_all_probes_disconnected_from_this_monitor": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Monitor], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "is_no_probe_enabled_on_this_monitor": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Monitor], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
        },
    }
}

func (r *MonitorResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    // Prevent panic if the provider has not been configured.
    if req.ProviderData == nil {
        return
    }

    client, ok := req.ProviderData.(*Client)

    if !ok {
        resp.Diagnostics.AddError(
            "Unexpected Resource Configure Type",
            fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
        )

        return
    }

    r.client = client
}


func (r *MonitorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data MonitorResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create API request body
    monitorRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "projectId": data.ProjectId.ValueString(),
        "name": data.Name.ValueString(),
        "description": data.Description.ValueString(),
        "labels": r.convertTerraformListToInterface(data.Labels),
        "monitorType": data.MonitorType.ValueString(),
        "currentMonitorStatusId": data.CurrentMonitorStatusId.ValueString(),
        "monitorSteps": r.convertTerraformMapToInterface(data.MonitorSteps),
        "monitoringInterval": data.MonitoringInterval.ValueString(),
        "customFields": r.convertTerraformMapToInterface(data.CustomFields),
        "isOwnerNotifiedOfResourceCreation": data.IsOwnerNotifiedOfResourceCreation.ValueBool(),
        "disableActiveMonitoring": data.DisableActiveMonitoring.ValueBool(),
        "incomingRequestMonitorHeartbeatCheckedAt": r.convertTerraformMapToInterface(data.IncomingRequestMonitorHeartbeatCheckedAt),
        "telemetryMonitorNextMonitorAt": r.convertTerraformMapToInterface(data.TelemetryMonitorNextMonitorAt),
        "telemetryMonitorLastMonitorAt": r.convertTerraformMapToInterface(data.TelemetryMonitorLastMonitorAt),
        "serverMonitorRequestReceivedAt": r.convertTerraformMapToInterface(data.ServerMonitorRequestReceivedAt),
        "serverMonitorSecretKey": data.ServerMonitorSecretKey.ValueString(),
        "incomingRequestSecretKey": data.IncomingRequestSecretKey.ValueString(),
        "incomingMonitorRequest": r.convertTerraformMapToInterface(data.IncomingMonitorRequest),
        "serverMonitorResponse": r.convertTerraformMapToInterface(data.ServerMonitorResponse),
        },
    }

    // Make API call
    httpResp, err := r.client.Post("/monitor", monitorRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create monitor, got error: %s", err))
        return
    }

    var monitorResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &monitorResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse monitor response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := monitorResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = monitorResponse
    }

    if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if val, ok := dataMap["name"].(string); ok && val != "" {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if val, ok := dataMap["description"].(string); ok && val != "" {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.Labels = listValue
    } else if dataMap["labels"] == nil {
        data.Labels = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["monitorType"].(string); ok && val != "" {
        data.MonitorType = types.StringValue(val)
    } else {
        data.MonitorType = types.StringNull()
    }
    if val, ok := dataMap["currentMonitorStatusId"].(string); ok && val != "" {
        data.CurrentMonitorStatusId = types.StringValue(val)
    } else {
        data.CurrentMonitorStatusId = types.StringNull()
    }
    if val, ok := dataMap["monitorSteps"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.MonitorSteps = mapValue
    } else if dataMap["monitorSteps"] == nil {
        data.MonitorSteps = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["monitoringInterval"].(string); ok && val != "" {
        data.MonitoringInterval = types.StringValue(val)
    } else {
        data.MonitoringInterval = types.StringNull()
    }
    if val, ok := dataMap["customFields"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CustomFields = mapValue
    } else if dataMap["customFields"] == nil {
        data.CustomFields = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["isOwnerNotifiedOfResourceCreation"].(bool); ok {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolValue(val)
    } else if dataMap["isOwnerNotifiedOfResourceCreation"] == nil {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolNull()
    }
    if val, ok := dataMap["disableActiveMonitoring"].(bool); ok {
        data.DisableActiveMonitoring = types.BoolValue(val)
    } else if dataMap["disableActiveMonitoring"] == nil {
        data.DisableActiveMonitoring = types.BoolNull()
    }
    if val, ok := dataMap["incomingRequestMonitorHeartbeatCheckedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.IncomingRequestMonitorHeartbeatCheckedAt = mapValue
    } else if dataMap["incomingRequestMonitorHeartbeatCheckedAt"] == nil {
        data.IncomingRequestMonitorHeartbeatCheckedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["telemetryMonitorNextMonitorAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.TelemetryMonitorNextMonitorAt = mapValue
    } else if dataMap["telemetryMonitorNextMonitorAt"] == nil {
        data.TelemetryMonitorNextMonitorAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["telemetryMonitorLastMonitorAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.TelemetryMonitorLastMonitorAt = mapValue
    } else if dataMap["telemetryMonitorLastMonitorAt"] == nil {
        data.TelemetryMonitorLastMonitorAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["serverMonitorRequestReceivedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.ServerMonitorRequestReceivedAt = mapValue
    } else if dataMap["serverMonitorRequestReceivedAt"] == nil {
        data.ServerMonitorRequestReceivedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["serverMonitorSecretKey"].(string); ok && val != "" {
        data.ServerMonitorSecretKey = types.StringValue(val)
    } else {
        data.ServerMonitorSecretKey = types.StringNull()
    }
    if val, ok := dataMap["incomingRequestSecretKey"].(string); ok && val != "" {
        data.IncomingRequestSecretKey = types.StringValue(val)
    } else {
        data.IncomingRequestSecretKey = types.StringNull()
    }
    if val, ok := dataMap["incomingMonitorRequest"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.IncomingMonitorRequest = mapValue
    } else if dataMap["incomingMonitorRequest"] == nil {
        data.IncomingMonitorRequest = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["serverMonitorResponse"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.ServerMonitorResponse = mapValue
    } else if dataMap["serverMonitorResponse"] == nil {
        data.ServerMonitorResponse = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CreatedAt = mapValue
    } else if dataMap["createdAt"] == nil {
        data.CreatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.UpdatedAt = mapValue
    } else if dataMap["updatedAt"] == nil {
        data.UpdatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.DeletedAt = mapValue
    } else if dataMap["deletedAt"] == nil {
        data.DeletedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if dataMap["version"] == nil {
        data.Version = types.NumberNull()
    }
    if val, ok := dataMap["slug"].(string); ok && val != "" {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["disableActiveMonitoringBecauseOfScheduledMaintenanceEvent"].(bool); ok {
        data.DisableActiveMonitoringBecauseOfScheduledMaintenanceEvent = types.BoolValue(val)
    } else if dataMap["disableActiveMonitoringBecauseOfScheduledMaintenanceEvent"] == nil {
        data.DisableActiveMonitoringBecauseOfScheduledMaintenanceEvent = types.BoolNull()
    }
    if val, ok := dataMap["disableActiveMonitoringBecauseOfManualIncident"].(bool); ok {
        data.DisableActiveMonitoringBecauseOfManualIncident = types.BoolValue(val)
    } else if dataMap["disableActiveMonitoringBecauseOfManualIncident"] == nil {
        data.DisableActiveMonitoringBecauseOfManualIncident = types.BoolNull()
    }
    if val, ok := dataMap["isAllProbesDisconnectedFromThisMonitor"].(bool); ok {
        data.IsAllProbesDisconnectedFromThisMonitor = types.BoolValue(val)
    } else if dataMap["isAllProbesDisconnectedFromThisMonitor"] == nil {
        data.IsAllProbesDisconnectedFromThisMonitor = types.BoolNull()
    }
    if val, ok := dataMap["isNoProbeEnabledOnThisMonitor"].(bool); ok {
        data.IsNoProbeEnabledOnThisMonitor = types.BoolValue(val)
    } else if dataMap["isNoProbeEnabledOnThisMonitor"] == nil {
        data.IsNoProbeEnabledOnThisMonitor = types.BoolNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "created a resource")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MonitorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data MonitorResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "description": true,
        "labels": true,
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
        "serverMonitorRequestReceivedAt": true,
        "serverMonitorSecretKey": true,
        "incomingRequestSecretKey": true,
        "incomingMonitorRequest": true,
        "serverMonitorResponse": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "slug": true,
        "createdByUserId": true,
        "disableActiveMonitoringBecauseOfScheduledMaintenanceEvent": true,
        "disableActiveMonitoringBecauseOfManualIncident": true,
        "isAllProbesDisconnectedFromThisMonitor": true,
        "isNoProbeEnabledOnThisMonitor": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect("/monitor/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read monitor, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var monitorResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &monitorResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse monitor response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := monitorResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = monitorResponse
    }

    if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if val, ok := dataMap["name"].(string); ok && val != "" {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if val, ok := dataMap["description"].(string); ok && val != "" {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.Labels = listValue
    } else if dataMap["labels"] == nil {
        data.Labels = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["monitorType"].(string); ok && val != "" {
        data.MonitorType = types.StringValue(val)
    } else {
        data.MonitorType = types.StringNull()
    }
    if val, ok := dataMap["currentMonitorStatusId"].(string); ok && val != "" {
        data.CurrentMonitorStatusId = types.StringValue(val)
    } else {
        data.CurrentMonitorStatusId = types.StringNull()
    }
    if val, ok := dataMap["monitorSteps"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.MonitorSteps = mapValue
    } else if dataMap["monitorSteps"] == nil {
        data.MonitorSteps = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["monitoringInterval"].(string); ok && val != "" {
        data.MonitoringInterval = types.StringValue(val)
    } else {
        data.MonitoringInterval = types.StringNull()
    }
    if val, ok := dataMap["customFields"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CustomFields = mapValue
    } else if dataMap["customFields"] == nil {
        data.CustomFields = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["isOwnerNotifiedOfResourceCreation"].(bool); ok {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolValue(val)
    } else if dataMap["isOwnerNotifiedOfResourceCreation"] == nil {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolNull()
    }
    if val, ok := dataMap["disableActiveMonitoring"].(bool); ok {
        data.DisableActiveMonitoring = types.BoolValue(val)
    } else if dataMap["disableActiveMonitoring"] == nil {
        data.DisableActiveMonitoring = types.BoolNull()
    }
    if val, ok := dataMap["incomingRequestMonitorHeartbeatCheckedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.IncomingRequestMonitorHeartbeatCheckedAt = mapValue
    } else if dataMap["incomingRequestMonitorHeartbeatCheckedAt"] == nil {
        data.IncomingRequestMonitorHeartbeatCheckedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["telemetryMonitorNextMonitorAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.TelemetryMonitorNextMonitorAt = mapValue
    } else if dataMap["telemetryMonitorNextMonitorAt"] == nil {
        data.TelemetryMonitorNextMonitorAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["telemetryMonitorLastMonitorAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.TelemetryMonitorLastMonitorAt = mapValue
    } else if dataMap["telemetryMonitorLastMonitorAt"] == nil {
        data.TelemetryMonitorLastMonitorAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["serverMonitorRequestReceivedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.ServerMonitorRequestReceivedAt = mapValue
    } else if dataMap["serverMonitorRequestReceivedAt"] == nil {
        data.ServerMonitorRequestReceivedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["serverMonitorSecretKey"].(string); ok && val != "" {
        data.ServerMonitorSecretKey = types.StringValue(val)
    } else {
        data.ServerMonitorSecretKey = types.StringNull()
    }
    if val, ok := dataMap["incomingRequestSecretKey"].(string); ok && val != "" {
        data.IncomingRequestSecretKey = types.StringValue(val)
    } else {
        data.IncomingRequestSecretKey = types.StringNull()
    }
    if val, ok := dataMap["incomingMonitorRequest"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.IncomingMonitorRequest = mapValue
    } else if dataMap["incomingMonitorRequest"] == nil {
        data.IncomingMonitorRequest = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["serverMonitorResponse"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.ServerMonitorResponse = mapValue
    } else if dataMap["serverMonitorResponse"] == nil {
        data.ServerMonitorResponse = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CreatedAt = mapValue
    } else if dataMap["createdAt"] == nil {
        data.CreatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.UpdatedAt = mapValue
    } else if dataMap["updatedAt"] == nil {
        data.UpdatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.DeletedAt = mapValue
    } else if dataMap["deletedAt"] == nil {
        data.DeletedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if dataMap["version"] == nil {
        data.Version = types.NumberNull()
    }
    if val, ok := dataMap["slug"].(string); ok && val != "" {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["disableActiveMonitoringBecauseOfScheduledMaintenanceEvent"].(bool); ok {
        data.DisableActiveMonitoringBecauseOfScheduledMaintenanceEvent = types.BoolValue(val)
    } else if dataMap["disableActiveMonitoringBecauseOfScheduledMaintenanceEvent"] == nil {
        data.DisableActiveMonitoringBecauseOfScheduledMaintenanceEvent = types.BoolNull()
    }
    if val, ok := dataMap["disableActiveMonitoringBecauseOfManualIncident"].(bool); ok {
        data.DisableActiveMonitoringBecauseOfManualIncident = types.BoolValue(val)
    } else if dataMap["disableActiveMonitoringBecauseOfManualIncident"] == nil {
        data.DisableActiveMonitoringBecauseOfManualIncident = types.BoolNull()
    }
    if val, ok := dataMap["isAllProbesDisconnectedFromThisMonitor"].(bool); ok {
        data.IsAllProbesDisconnectedFromThisMonitor = types.BoolValue(val)
    } else if dataMap["isAllProbesDisconnectedFromThisMonitor"] == nil {
        data.IsAllProbesDisconnectedFromThisMonitor = types.BoolNull()
    }
    if val, ok := dataMap["isNoProbeEnabledOnThisMonitor"].(bool); ok {
        data.IsNoProbeEnabledOnThisMonitor = types.BoolValue(val)
    } else if dataMap["isNoProbeEnabledOnThisMonitor"] == nil {
        data.IsNoProbeEnabledOnThisMonitor = types.BoolNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MonitorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data MonitorResourceModel
    var state MonitorResourceModel

    // Read Terraform current state data to get the ID
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Read Terraform plan data to get the new values
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Use the ID from the current state
    data.Id = state.Id

    // Create API request body
    monitorRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "name": data.Name.ValueString(),
        "description": data.Description.ValueString(),
        "labels": r.convertTerraformListToInterface(data.Labels),
        "monitorType": data.MonitorType.ValueString(),
        "currentMonitorStatusId": data.CurrentMonitorStatusId.ValueString(),
        "monitorSteps": r.convertTerraformMapToInterface(data.MonitorSteps),
        "monitoringInterval": data.MonitoringInterval.ValueString(),
        "customFields": r.convertTerraformMapToInterface(data.CustomFields),
        "isOwnerNotifiedOfResourceCreation": data.IsOwnerNotifiedOfResourceCreation.ValueBool(),
        "disableActiveMonitoring": data.DisableActiveMonitoring.ValueBool(),
        "incomingRequestMonitorHeartbeatCheckedAt": r.convertTerraformMapToInterface(data.IncomingRequestMonitorHeartbeatCheckedAt),
        "telemetryMonitorNextMonitorAt": r.convertTerraformMapToInterface(data.TelemetryMonitorNextMonitorAt),
        "telemetryMonitorLastMonitorAt": r.convertTerraformMapToInterface(data.TelemetryMonitorLastMonitorAt),
        "serverMonitorRequestReceivedAt": r.convertTerraformMapToInterface(data.ServerMonitorRequestReceivedAt),
        "serverMonitorSecretKey": data.ServerMonitorSecretKey.ValueString(),
        "incomingRequestSecretKey": data.IncomingRequestSecretKey.ValueString(),
        "incomingMonitorRequest": r.convertTerraformMapToInterface(data.IncomingMonitorRequest),
        "serverMonitorResponse": r.convertTerraformMapToInterface(data.ServerMonitorResponse),
        },
    }

    // Make API call
    httpResp, err := r.client.Put("/monitor/" + data.Id.ValueString() + "", monitorRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update monitor, got error: %s", err))
        return
    }

    // Parse the update response
    var monitorResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &monitorResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse monitor response, got error: %s", err))
        return
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "description": true,
        "labels": true,
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
        "serverMonitorRequestReceivedAt": true,
        "serverMonitorSecretKey": true,
        "incomingRequestSecretKey": true,
        "incomingMonitorRequest": true,
        "serverMonitorResponse": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "slug": true,
        "createdByUserId": true,
        "disableActiveMonitoringBecauseOfScheduledMaintenanceEvent": true,
        "disableActiveMonitoringBecauseOfManualIncident": true,
        "isAllProbesDisconnectedFromThisMonitor": true,
        "isNoProbeEnabledOnThisMonitor": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect("/monitor/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read monitor after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse monitor read response, got error: %s", err))
        return
    }

    // Update the model with response data from the Read operation
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := readResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = readResponse
    }

    if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if val, ok := dataMap["name"].(string); ok && val != "" {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if val, ok := dataMap["description"].(string); ok && val != "" {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.Labels = listValue
    } else if dataMap["labels"] == nil {
        data.Labels = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["monitorType"].(string); ok && val != "" {
        data.MonitorType = types.StringValue(val)
    } else {
        data.MonitorType = types.StringNull()
    }
    if val, ok := dataMap["currentMonitorStatusId"].(string); ok && val != "" {
        data.CurrentMonitorStatusId = types.StringValue(val)
    } else {
        data.CurrentMonitorStatusId = types.StringNull()
    }
    if val, ok := dataMap["monitorSteps"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.MonitorSteps = mapValue
    } else if dataMap["monitorSteps"] == nil {
        data.MonitorSteps = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["monitoringInterval"].(string); ok && val != "" {
        data.MonitoringInterval = types.StringValue(val)
    } else {
        data.MonitoringInterval = types.StringNull()
    }
    if val, ok := dataMap["customFields"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CustomFields = mapValue
    } else if dataMap["customFields"] == nil {
        data.CustomFields = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["isOwnerNotifiedOfResourceCreation"].(bool); ok {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolValue(val)
    } else if dataMap["isOwnerNotifiedOfResourceCreation"] == nil {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolNull()
    }
    if val, ok := dataMap["disableActiveMonitoring"].(bool); ok {
        data.DisableActiveMonitoring = types.BoolValue(val)
    } else if dataMap["disableActiveMonitoring"] == nil {
        data.DisableActiveMonitoring = types.BoolNull()
    }
    if val, ok := dataMap["incomingRequestMonitorHeartbeatCheckedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.IncomingRequestMonitorHeartbeatCheckedAt = mapValue
    } else if dataMap["incomingRequestMonitorHeartbeatCheckedAt"] == nil {
        data.IncomingRequestMonitorHeartbeatCheckedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["telemetryMonitorNextMonitorAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.TelemetryMonitorNextMonitorAt = mapValue
    } else if dataMap["telemetryMonitorNextMonitorAt"] == nil {
        data.TelemetryMonitorNextMonitorAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["telemetryMonitorLastMonitorAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.TelemetryMonitorLastMonitorAt = mapValue
    } else if dataMap["telemetryMonitorLastMonitorAt"] == nil {
        data.TelemetryMonitorLastMonitorAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["serverMonitorRequestReceivedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.ServerMonitorRequestReceivedAt = mapValue
    } else if dataMap["serverMonitorRequestReceivedAt"] == nil {
        data.ServerMonitorRequestReceivedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["serverMonitorSecretKey"].(string); ok && val != "" {
        data.ServerMonitorSecretKey = types.StringValue(val)
    } else {
        data.ServerMonitorSecretKey = types.StringNull()
    }
    if val, ok := dataMap["incomingRequestSecretKey"].(string); ok && val != "" {
        data.IncomingRequestSecretKey = types.StringValue(val)
    } else {
        data.IncomingRequestSecretKey = types.StringNull()
    }
    if val, ok := dataMap["incomingMonitorRequest"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.IncomingMonitorRequest = mapValue
    } else if dataMap["incomingMonitorRequest"] == nil {
        data.IncomingMonitorRequest = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["serverMonitorResponse"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.ServerMonitorResponse = mapValue
    } else if dataMap["serverMonitorResponse"] == nil {
        data.ServerMonitorResponse = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CreatedAt = mapValue
    } else if dataMap["createdAt"] == nil {
        data.CreatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.UpdatedAt = mapValue
    } else if dataMap["updatedAt"] == nil {
        data.UpdatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.DeletedAt = mapValue
    } else if dataMap["deletedAt"] == nil {
        data.DeletedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if dataMap["version"] == nil {
        data.Version = types.NumberNull()
    }
    if val, ok := dataMap["slug"].(string); ok && val != "" {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["disableActiveMonitoringBecauseOfScheduledMaintenanceEvent"].(bool); ok {
        data.DisableActiveMonitoringBecauseOfScheduledMaintenanceEvent = types.BoolValue(val)
    } else if dataMap["disableActiveMonitoringBecauseOfScheduledMaintenanceEvent"] == nil {
        data.DisableActiveMonitoringBecauseOfScheduledMaintenanceEvent = types.BoolNull()
    }
    if val, ok := dataMap["disableActiveMonitoringBecauseOfManualIncident"].(bool); ok {
        data.DisableActiveMonitoringBecauseOfManualIncident = types.BoolValue(val)
    } else if dataMap["disableActiveMonitoringBecauseOfManualIncident"] == nil {
        data.DisableActiveMonitoringBecauseOfManualIncident = types.BoolNull()
    }
    if val, ok := dataMap["isAllProbesDisconnectedFromThisMonitor"].(bool); ok {
        data.IsAllProbesDisconnectedFromThisMonitor = types.BoolValue(val)
    } else if dataMap["isAllProbesDisconnectedFromThisMonitor"] == nil {
        data.IsAllProbesDisconnectedFromThisMonitor = types.BoolNull()
    }
    if val, ok := dataMap["isNoProbeEnabledOnThisMonitor"].(bool); ok {
        data.IsNoProbeEnabledOnThisMonitor = types.BoolValue(val)
    } else if dataMap["isNoProbeEnabledOnThisMonitor"] == nil {
        data.IsNoProbeEnabledOnThisMonitor = types.BoolNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MonitorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data MonitorResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    _, err := r.client.Delete("/monitor/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete monitor, got error: %s", err))
        return
    }
}


func (r *MonitorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *MonitorResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
    if terraformMap.IsNull() || terraformMap.IsUnknown() {
        return nil
    }
    
    result := make(map[string]string)
    terraformMap.ElementsAs(context.Background(), &result, false)
    
    // Convert map[string]string to map[string]interface{}
    interfaceResult := make(map[string]interface{})
    for key, value := range result {
        interfaceResult[key] = value
    }
    
    return interfaceResult
}

// Helper method to convert Terraform list to Go interface{}
func (r *MonitorResource) convertTerraformListToInterface(terraformList types.List) interface{} {
    if terraformList.IsNull() || terraformList.IsUnknown() {
        return nil
    }
    
    var result []interface{}
    terraformList.ElementsAs(context.Background(), &result, false)
    return result
}
