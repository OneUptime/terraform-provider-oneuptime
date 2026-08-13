package provider

import (
    "context"
    "fmt"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
    "github.com/hashicorp/terraform-plugin-log/tflog"
    "math/big"
    "net/http"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "encoding/json"
    "net/url"
    "strings"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/numberplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
    "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    Labels types.Set `tfsdk:"labels"`
    DependsOnMonitors types.Set `tfsdk:"depends_on_monitors"`
    SuppressAlertsWhenParentMonitorStatuses types.Set `tfsdk:"suppress_alerts_when_parent_monitor_statuses"`
    MonitorTemplateId types.String `tfsdk:"monitor_template_id"`
    MonitorType types.String `tfsdk:"monitor_type"`
    CurrentMonitorStatusId types.String `tfsdk:"current_monitor_status_id"`
    MonitorSteps MonitorStepsValue `tfsdk:"monitor_steps"`
    MonitoringInterval types.String `tfsdk:"monitoring_interval"`
    CustomFields JSONSubsetValue `tfsdk:"custom_fields"`
    DisableActiveMonitoring types.Bool `tfsdk:"disable_active_monitoring"`
    IncomingRequestMonitorHeartbeatCheckedAt RFC3339Value `tfsdk:"incoming_request_monitor_heartbeat_checked_at"`
    TelemetryMonitorNextMonitorAt RFC3339Value `tfsdk:"telemetry_monitor_next_monitor_at"`
    TelemetryMonitorLastMonitorAt RFC3339Value `tfsdk:"telemetry_monitor_last_monitor_at"`
    ServerMonitorRequestReceivedAt RFC3339Value `tfsdk:"server_monitor_request_received_at"`
    IncomingMonitorRequest JSONSubsetValue `tfsdk:"incoming_monitor_request"`
    ServerMonitorResponse JSONSubsetValue `tfsdk:"server_monitor_response"`
    MinimumProbeAgreement types.Number `tfsdk:"minimum_probe_agreement"`
    CreatedAt RFC3339Value `tfsdk:"created_at"`
    UpdatedAt RFC3339Value `tfsdk:"updated_at"`
    DeletedAt RFC3339Value `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    Slug types.String `tfsdk:"slug"`
    IsOwnerNotifiedOfResourceCreation types.Bool `tfsdk:"is_owner_notified_of_resource_creation"`
    DisableActiveMonitoringBecauseOfScheduledMaintenanceEvent types.Bool `tfsdk:"disable_active_monitoring_because_of_scheduled_maintenance_event"`
    DisableActiveMonitoringBecauseOfManualIncident types.Bool `tfsdk:"disable_active_monitoring_because_of_manual_incident"`
    ServerMonitorSecretKey types.String `tfsdk:"server_monitor_secret_key"`
    IncomingRequestSecretKey types.String `tfsdk:"incoming_request_secret_key"`
    IncomingEmailSecretKey types.String `tfsdk:"incoming_email_secret_key"`
    IncomingEmailMonitorLastEmailReceivedAt RFC3339Value `tfsdk:"incoming_email_monitor_last_email_received_at"`
    IncomingEmailMonitorRequest JSONSubsetValue `tfsdk:"incoming_email_monitor_request"`
    IncomingEmailMonitorHeartbeatCheckedAt RFC3339Value `tfsdk:"incoming_email_monitor_heartbeat_checked_at"`
    IsAllProbesDisconnectedFromThisMonitor types.Bool `tfsdk:"is_all_probes_disconnected_from_this_monitor"`
    IsNoProbeEnabledOnThisMonitor types.Bool `tfsdk:"is_no_probe_enabled_on_this_monitor"`
}

func (r *MonitorResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_monitor"
}

func (r *MonitorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Monitor is anything that monitors your API, Websites, IP, Network or more. You can also create static monitor that does not monitor anything.",

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
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Any friendly name for this monitor.",
                Required: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description that will help you remember.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in..",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
            },
            "depends_on_monitors": schema.SetAttribute{
                MarkdownDescription: "Parent monitors this monitor depends on. When a parent is offline (or in one of the configured suppression statuses), alerts and incidents from this monitor are suppressed at creation time — the monitor keeps evaluating and its status timeline still updates..",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
            },
            "suppress_alerts_when_parent_monitor_statuses": schema.SetAttribute{
                MarkdownDescription: "Parent monitor statuses that suppress this monitor's alerts and incidents. When empty, statuses flagged as offline suppress (the default). Only used when Depends On Monitors is set..",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
            },
            "monitor_template_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "monitor_type": schema.StringAttribute{
                MarkdownDescription: "What is the type of this monitor? Website? API? etc..",
                Required: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
                Validators: []validator.String{
                    stringvalidator.OneOf("Manual", "Website", "API", "Ping", "Kubernetes", "Docker", "Host", "Podman", "Docker Swarm", "Proxmox", "Ceph", "IoT Device", "IP", "Incoming Request", "Incoming Email", "Port", "Server", "SSL Certificate", "SQL Query", "Synthetic Monitor", "Custom JavaScript Code", "Logs", "Metrics", "Traces", "Exceptions", "Profiles", "Network Device", "DNS", "DNSSEC", "Domain", "External Status Page"),
                },
            },
            "current_monitor_status_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "monitor_steps": MonitorStepsSchemaAttribute("MonitorSteps object"),
            "monitoring_interval": schema.StringAttribute{
                MarkdownDescription: "How often would you like OneUptime to monitor this resource?.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "custom_fields": schema.StringAttribute{
                MarkdownDescription: "Custom Fields on this resource..",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "disable_active_monitoring": schema.BoolAttribute{
                MarkdownDescription: "Disable active monitoring for this resource?.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "incoming_request_monitor_heartbeat_checked_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "telemetry_monitor_next_monitor_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "telemetry_monitor_last_monitor_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "server_monitor_request_received_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "incoming_monitor_request": schema.StringAttribute{
                MarkdownDescription: "Incoming Monitor Request for Incoming Request Monitor.",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                    stringplanmodifier.RequiresReplace(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "server_monitor_response": schema.StringAttribute{
                MarkdownDescription: "Server Monitor Response for Server Monitor.",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                    stringplanmodifier.RequiresReplace(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "minimum_probe_agreement": schema.NumberAttribute{
                MarkdownDescription: "Minimum number of probes that must agree on a status before the monitor status changes. If null, all enabled and connected probes must agree..",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "deleted_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "version": schema.NumberAttribute{
                MarkdownDescription: "Object version",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Computed: true,
            },
            "is_owner_notified_of_resource_creation": schema.BoolAttribute{
                MarkdownDescription: "Are owners notified of when this resource is created?.",
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
            "server_monitor_secret_key": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "incoming_request_secret_key": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "incoming_email_secret_key": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "incoming_email_monitor_last_email_received_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "incoming_email_monitor_request": schema.StringAttribute{
                MarkdownDescription: "This field is for Incoming Email Monitor only. Last email data received..",
                CustomType: JSONSubsetType{},
                Computed: true,
            },
            "incoming_email_monitor_heartbeat_checked_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
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



    // Create API request body. Unset (null/unknown) optional fields are
    // omitted so server-side defaults apply instead of being overwritten
    // with zero values.
    monitorRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := monitorRequest["data"].(map[string]interface{})

    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        requestDataMap["name"] = data.Name.ValueString()
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        requestDataMap["description"] = data.Description.ValueString()
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        requestDataMap["createdByUserId"] = data.CreatedByUserId.ValueString()
    }
    if !data.Labels.IsNull() && !data.Labels.IsUnknown() {
        requestDataMap["labels"] = r.convertTerraformSetToInterface(data.Labels)
    }
    if !data.DependsOnMonitors.IsNull() && !data.DependsOnMonitors.IsUnknown() {
        requestDataMap["dependsOnMonitors"] = r.convertTerraformSetToInterface(data.DependsOnMonitors)
    }
    if !data.SuppressAlertsWhenParentMonitorStatuses.IsNull() && !data.SuppressAlertsWhenParentMonitorStatuses.IsUnknown() {
        requestDataMap["suppressAlertsWhenParentMonitorStatuses"] = r.convertTerraformSetToInterface(data.SuppressAlertsWhenParentMonitorStatuses)
    }
    if !data.MonitorTemplateId.IsNull() && !data.MonitorTemplateId.IsUnknown() {
        requestDataMap["monitorTemplateId"] = data.MonitorTemplateId.ValueString()
    }
    if !data.MonitorType.IsNull() && !data.MonitorType.IsUnknown() {
        requestDataMap["monitorType"] = data.MonitorType.ValueString()
    }
    if !data.CurrentMonitorStatusId.IsNull() && !data.CurrentMonitorStatusId.IsUnknown() {
        requestDataMap["currentMonitorStatusId"] = data.CurrentMonitorStatusId.ValueString()
    }
    if !data.MonitorSteps.IsNull() && !data.MonitorSteps.IsUnknown() {
        monitorStepsValue, monitorStepsDiags := MonitorStepsToAPI(ctx, data.MonitorSteps.ListValue)
        resp.Diagnostics.Append(monitorStepsDiags...)
        if resp.Diagnostics.HasError() {
            return
        }
        requestDataMap["monitorSteps"] = monitorStepsValue
    }
    if !data.MonitoringInterval.IsNull() && !data.MonitoringInterval.IsUnknown() {
        requestDataMap["monitoringInterval"] = data.MonitoringInterval.ValueString()
    }
    if parsedCustomFields := r.parseJSONField(data.CustomFields); parsedCustomFields != nil {
        requestDataMap["customFields"] = parsedCustomFields
    }
    if !data.DisableActiveMonitoring.IsNull() && !data.DisableActiveMonitoring.IsUnknown() {
        requestDataMap["disableActiveMonitoring"] = data.DisableActiveMonitoring.ValueBool()
    }
    if !data.IncomingRequestMonitorHeartbeatCheckedAt.IsNull() && !data.IncomingRequestMonitorHeartbeatCheckedAt.IsUnknown() {
        requestDataMap["incomingRequestMonitorHeartbeatCheckedAt"] = data.IncomingRequestMonitorHeartbeatCheckedAt.ValueString()
    }
    if !data.TelemetryMonitorNextMonitorAt.IsNull() && !data.TelemetryMonitorNextMonitorAt.IsUnknown() {
        requestDataMap["telemetryMonitorNextMonitorAt"] = data.TelemetryMonitorNextMonitorAt.ValueString()
    }
    if !data.TelemetryMonitorLastMonitorAt.IsNull() && !data.TelemetryMonitorLastMonitorAt.IsUnknown() {
        requestDataMap["telemetryMonitorLastMonitorAt"] = data.TelemetryMonitorLastMonitorAt.ValueString()
    }
    if !data.ServerMonitorRequestReceivedAt.IsNull() && !data.ServerMonitorRequestReceivedAt.IsUnknown() {
        requestDataMap["serverMonitorRequestReceivedAt"] = data.ServerMonitorRequestReceivedAt.ValueString()
    }
    if parsedIncomingMonitorRequest := r.parseJSONField(data.IncomingMonitorRequest); parsedIncomingMonitorRequest != nil {
        requestDataMap["incomingMonitorRequest"] = parsedIncomingMonitorRequest
    }
    if parsedServerMonitorResponse := r.parseJSONField(data.ServerMonitorResponse); parsedServerMonitorResponse != nil {
        requestDataMap["serverMonitorResponse"] = parsedServerMonitorResponse
    }
    if !data.MinimumProbeAgreement.IsNull() && !data.MinimumProbeAgreement.IsUnknown() {
        requestDataMap["minimumProbeAgreement"] = r.bigFloatToFloat64(data.MinimumProbeAgreement.ValueBigFloat())
    }

    // Make API call
    httpResp, err := r.client.Post(ctx, "/monitor", monitorRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create monitor, got error: %s", err))
        return
    }

    var monitorResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &monitorResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to create monitor: %s", err))
        return
    }

    // Extract the new resource id from the create response.
    createdId := ""
    if wrapper, ok := monitorResponse["data"].(map[string]interface{}); ok {
        if val, ok := wrapper["_id"].(string); ok {
            createdId = val
        }
    } else if val, ok := monitorResponse["_id"].(string); ok {
        createdId = val
    }
    if createdId == "" {
        resp.Diagnostics.AddError("OneUptime API Error", "Create response for monitor did not contain an id. This is a bug in the provider or the API; please report it.")
        return
    }
    data.Id = types.StringValue(createdId)

    /*
     * The server has committed the row. Persist what we know to state BEFORE
     * the read-back: if the read-back fails and we return without setting
     * state, Terraform never learns the resource exists and the created
     * monitor is orphaned server-side — never refreshed, never
     * destroyed. Delete already refuses to drop state on failure for the
     * same reason; Create must not either.
     */
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Re-read the resource so state reflects server-normalized values.
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "description": true,
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
        "disableActiveMonitoring": true,
        "incomingRequestMonitorHeartbeatCheckedAt": true,
        "telemetryMonitorNextMonitorAt": true,
        "telemetryMonitorLastMonitorAt": true,
        "serverMonitorRequestReceivedAt": true,
        "incomingMonitorRequest": true,
        "serverMonitorResponse": true,
        "minimumProbeAgreement": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "slug": true,
        "isOwnerNotifiedOfResourceCreation": true,
        "disableActiveMonitoringBecauseOfScheduledMaintenanceEvent": true,
        "disableActiveMonitoringBecauseOfManualIncident": true,
        "serverMonitorSecretKey": true,
        "incomingRequestSecretKey": true,
        "incomingEmailSecretKey": true,
        "incomingEmailMonitorLastEmailReceivedAt": true,
        "incomingEmailMonitorRequest": true,
        "incomingEmailMonitorHeartbeatCheckedAt": true,
        "isAllProbesDisconnectedFromThisMonitor": true,
        "isNoProbeEnabledOnThisMonitor": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/monitor/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        /*
         * State already owns the id, so the resource is tracked and the next
         * refresh reconciles the remaining attributes. Warn rather than
         * error: erroring here would strand a real resource.
         */
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created monitor but could not read it back; state is incomplete until the next refresh: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created monitor but could not parse the read-back response; state is incomplete until the next refresh: %s", err))
        return
    }

    // Update the model with the authoritative read response
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := readResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = readResponse
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
    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Labels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Labels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["dependsOnMonitors"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.DependsOnMonitors = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.DependsOnMonitors = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["suppressAlertsWhenParentMonitorStatuses"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.SuppressAlertsWhenParentMonitorStatuses = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.SuppressAlertsWhenParentMonitorStatuses = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["monitorTemplateId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorTemplateId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MonitorTemplateId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MonitorTemplateId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MonitorTemplateId = types.StringValue(string(jsonBytes))
            } else {
                data.MonitorTemplateId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MonitorTemplateId = types.StringValue(string(jsonBytes))
            } else {
                data.MonitorTemplateId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MonitorTemplateId = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorTemplateId = types.StringNull()
        }
    } else if val, ok := dataMap["monitorTemplateId"].(string); ok {
        data.MonitorTemplateId = types.StringValue(val)
    } else {
        data.MonitorTemplateId = types.StringNull()
    }
    if obj, ok := dataMap["monitorType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MonitorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MonitorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MonitorType = types.StringValue(string(jsonBytes))
            } else {
                data.MonitorType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MonitorType = types.StringValue(string(jsonBytes))
            } else {
                data.MonitorType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MonitorType = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorType = types.StringNull()
        }
    } else if val, ok := dataMap["monitorType"].(string); ok {
        data.MonitorType = types.StringValue(val)
    } else {
        data.MonitorType = types.StringNull()
    }
    if obj, ok := dataMap["currentMonitorStatusId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CurrentMonitorStatusId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CurrentMonitorStatusId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CurrentMonitorStatusId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CurrentMonitorStatusId = types.StringValue(string(jsonBytes))
            } else {
                data.CurrentMonitorStatusId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CurrentMonitorStatusId = types.StringValue(string(jsonBytes))
            } else {
                data.CurrentMonitorStatusId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CurrentMonitorStatusId = types.StringValue(string(jsonBytes))
        } else {
            data.CurrentMonitorStatusId = types.StringNull()
        }
    } else if val, ok := dataMap["currentMonitorStatusId"].(string); ok {
        data.CurrentMonitorStatusId = types.StringValue(val)
    } else {
        data.CurrentMonitorStatusId = types.StringNull()
    }
    {
        mappedSteps, stepsDiags := MonitorStepsFromAPI(ctx, dataMap["monitorSteps"])
        resp.Diagnostics.Append(stepsDiags...)
        data.MonitorSteps = MonitorStepsValue{ListValue: mappedSteps}
    }
    if obj, ok := dataMap["monitoringInterval"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitoringInterval = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MonitoringInterval = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MonitoringInterval = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MonitoringInterval = types.StringValue(string(jsonBytes))
            } else {
                data.MonitoringInterval = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MonitoringInterval = types.StringValue(string(jsonBytes))
            } else {
                data.MonitoringInterval = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MonitoringInterval = types.StringValue(string(jsonBytes))
        } else {
            data.MonitoringInterval = types.StringNull()
        }
    } else if val, ok := dataMap["monitoringInterval"].(string); ok {
        data.MonitoringInterval = types.StringValue(val)
    } else {
        data.MonitoringInterval = types.StringNull()
    }
    if obj, ok := dataMap["customFields"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CustomFields = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CustomFields = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CustomFields = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CustomFields = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CustomFields = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CustomFields = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CustomFields = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CustomFields = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CustomFields = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["customFields"].(string); ok {
        data.CustomFields = NewJSONSubsetValue(val)
    } else {
        data.CustomFields = NewJSONSubsetNull()
    }
    if val, ok := dataMap["disableActiveMonitoring"].(bool); ok {
        data.DisableActiveMonitoring = types.BoolValue(val)
    }
    if obj, ok := dataMap["incomingRequestMonitorHeartbeatCheckedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.IncomingRequestMonitorHeartbeatCheckedAt = NewRFC3339Value(val)
        } else {
            data.IncomingRequestMonitorHeartbeatCheckedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["incomingRequestMonitorHeartbeatCheckedAt"].(string); ok && val != "" {
        data.IncomingRequestMonitorHeartbeatCheckedAt = NewRFC3339Value(val)
    } else {
        data.IncomingRequestMonitorHeartbeatCheckedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["telemetryMonitorNextMonitorAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.TelemetryMonitorNextMonitorAt = NewRFC3339Value(val)
        } else {
            data.TelemetryMonitorNextMonitorAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["telemetryMonitorNextMonitorAt"].(string); ok && val != "" {
        data.TelemetryMonitorNextMonitorAt = NewRFC3339Value(val)
    } else {
        data.TelemetryMonitorNextMonitorAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["telemetryMonitorLastMonitorAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.TelemetryMonitorLastMonitorAt = NewRFC3339Value(val)
        } else {
            data.TelemetryMonitorLastMonitorAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["telemetryMonitorLastMonitorAt"].(string); ok && val != "" {
        data.TelemetryMonitorLastMonitorAt = NewRFC3339Value(val)
    } else {
        data.TelemetryMonitorLastMonitorAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["serverMonitorRequestReceivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ServerMonitorRequestReceivedAt = NewRFC3339Value(val)
        } else {
            data.ServerMonitorRequestReceivedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["serverMonitorRequestReceivedAt"].(string); ok && val != "" {
        data.ServerMonitorRequestReceivedAt = NewRFC3339Value(val)
    } else {
        data.ServerMonitorRequestReceivedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["incomingMonitorRequest"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncomingMonitorRequest = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IncomingMonitorRequest = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IncomingMonitorRequest = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.IncomingMonitorRequest = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.IncomingMonitorRequest = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.IncomingMonitorRequest = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.IncomingMonitorRequest = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IncomingMonitorRequest = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.IncomingMonitorRequest = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["incomingMonitorRequest"].(string); ok {
        data.IncomingMonitorRequest = NewJSONSubsetValue(val)
    } else {
        data.IncomingMonitorRequest = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["serverMonitorResponse"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServerMonitorResponse = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ServerMonitorResponse = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ServerMonitorResponse = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ServerMonitorResponse = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.ServerMonitorResponse = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ServerMonitorResponse = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.ServerMonitorResponse = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ServerMonitorResponse = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.ServerMonitorResponse = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["serverMonitorResponse"].(string); ok {
        data.ServerMonitorResponse = NewJSONSubsetValue(val)
    } else {
        data.ServerMonitorResponse = NewJSONSubsetNull()
    }
    if val, ok := dataMap["minimumProbeAgreement"].(float64); ok {
        data.MinimumProbeAgreement = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["minimumProbeAgreement"].(int); ok {
        data.MinimumProbeAgreement = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["minimumProbeAgreement"].(int64); ok {
        data.MinimumProbeAgreement = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["minimumProbeAgreement"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.MinimumProbeAgreement = types.NumberValue(big.NewFloat(val))
        } else {
            data.MinimumProbeAgreement = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.MinimumProbeAgreement = types.NumberNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedAt = NewRFC3339Value(val)
        } else {
            data.CreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = NewRFC3339Value(val)
    } else {
        data.CreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.UpdatedAt = NewRFC3339Value(val)
        } else {
            data.UpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = NewRFC3339Value(val)
    } else {
        data.UpdatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.DeletedAt = NewRFC3339Value(val)
        } else {
            data.DeletedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = NewRFC3339Value(val)
    } else {
        data.DeletedAt = NewRFC3339Null()
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["version"].(int); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["version"].(int64); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["version"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.Version = types.NumberNull()
    }
    if obj, ok := dataMap["slug"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := dataMap["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if val, ok := dataMap["isOwnerNotifiedOfResourceCreation"].(bool); ok {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolValue(val)
    }
    if val, ok := dataMap["disableActiveMonitoringBecauseOfScheduledMaintenanceEvent"].(bool); ok {
        data.DisableActiveMonitoringBecauseOfScheduledMaintenanceEvent = types.BoolValue(val)
    }
    if val, ok := dataMap["disableActiveMonitoringBecauseOfManualIncident"].(bool); ok {
        data.DisableActiveMonitoringBecauseOfManualIncident = types.BoolValue(val)
    }
    if obj, ok := dataMap["serverMonitorSecretKey"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServerMonitorSecretKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ServerMonitorSecretKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ServerMonitorSecretKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ServerMonitorSecretKey = types.StringValue(string(jsonBytes))
            } else {
                data.ServerMonitorSecretKey = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ServerMonitorSecretKey = types.StringValue(string(jsonBytes))
            } else {
                data.ServerMonitorSecretKey = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ServerMonitorSecretKey = types.StringValue(string(jsonBytes))
        } else {
            data.ServerMonitorSecretKey = types.StringNull()
        }
    } else if val, ok := dataMap["serverMonitorSecretKey"].(string); ok {
        data.ServerMonitorSecretKey = types.StringValue(val)
    } else {
        data.ServerMonitorSecretKey = types.StringNull()
    }
    if obj, ok := dataMap["incomingRequestSecretKey"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncomingRequestSecretKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IncomingRequestSecretKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IncomingRequestSecretKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.IncomingRequestSecretKey = types.StringValue(string(jsonBytes))
            } else {
                data.IncomingRequestSecretKey = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.IncomingRequestSecretKey = types.StringValue(string(jsonBytes))
            } else {
                data.IncomingRequestSecretKey = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IncomingRequestSecretKey = types.StringValue(string(jsonBytes))
        } else {
            data.IncomingRequestSecretKey = types.StringNull()
        }
    } else if val, ok := dataMap["incomingRequestSecretKey"].(string); ok {
        data.IncomingRequestSecretKey = types.StringValue(val)
    } else {
        data.IncomingRequestSecretKey = types.StringNull()
    }
    if obj, ok := dataMap["incomingEmailSecretKey"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncomingEmailSecretKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IncomingEmailSecretKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IncomingEmailSecretKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.IncomingEmailSecretKey = types.StringValue(string(jsonBytes))
            } else {
                data.IncomingEmailSecretKey = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.IncomingEmailSecretKey = types.StringValue(string(jsonBytes))
            } else {
                data.IncomingEmailSecretKey = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IncomingEmailSecretKey = types.StringValue(string(jsonBytes))
        } else {
            data.IncomingEmailSecretKey = types.StringNull()
        }
    } else if val, ok := dataMap["incomingEmailSecretKey"].(string); ok {
        data.IncomingEmailSecretKey = types.StringValue(val)
    } else {
        data.IncomingEmailSecretKey = types.StringNull()
    }
    if obj, ok := dataMap["incomingEmailMonitorLastEmailReceivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.IncomingEmailMonitorLastEmailReceivedAt = NewRFC3339Value(val)
        } else {
            data.IncomingEmailMonitorLastEmailReceivedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["incomingEmailMonitorLastEmailReceivedAt"].(string); ok && val != "" {
        data.IncomingEmailMonitorLastEmailReceivedAt = NewRFC3339Value(val)
    } else {
        data.IncomingEmailMonitorLastEmailReceivedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["incomingEmailMonitorRequest"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncomingEmailMonitorRequest = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IncomingEmailMonitorRequest = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IncomingEmailMonitorRequest = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.IncomingEmailMonitorRequest = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.IncomingEmailMonitorRequest = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.IncomingEmailMonitorRequest = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.IncomingEmailMonitorRequest = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IncomingEmailMonitorRequest = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.IncomingEmailMonitorRequest = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["incomingEmailMonitorRequest"].(string); ok {
        data.IncomingEmailMonitorRequest = NewJSONSubsetValue(val)
    } else {
        data.IncomingEmailMonitorRequest = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["incomingEmailMonitorHeartbeatCheckedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.IncomingEmailMonitorHeartbeatCheckedAt = NewRFC3339Value(val)
        } else {
            data.IncomingEmailMonitorHeartbeatCheckedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["incomingEmailMonitorHeartbeatCheckedAt"].(string); ok && val != "" {
        data.IncomingEmailMonitorHeartbeatCheckedAt = NewRFC3339Value(val)
    } else {
        data.IncomingEmailMonitorHeartbeatCheckedAt = NewRFC3339Null()
    }
    if val, ok := dataMap["isAllProbesDisconnectedFromThisMonitor"].(bool); ok {
        data.IsAllProbesDisconnectedFromThisMonitor = types.BoolValue(val)
    }
    if val, ok := dataMap["isNoProbeEnabledOnThisMonitor"].(bool); ok {
        data.IsNoProbeEnabledOnThisMonitor = types.BoolValue(val)
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    // The read response is authoritative, but never let it clobber the id we just received.
    data.Id = types.StringValue(createdId)

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
        "disableActiveMonitoring": true,
        "incomingRequestMonitorHeartbeatCheckedAt": true,
        "telemetryMonitorNextMonitorAt": true,
        "telemetryMonitorLastMonitorAt": true,
        "serverMonitorRequestReceivedAt": true,
        "incomingMonitorRequest": true,
        "serverMonitorResponse": true,
        "minimumProbeAgreement": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "slug": true,
        "isOwnerNotifiedOfResourceCreation": true,
        "disableActiveMonitoringBecauseOfScheduledMaintenanceEvent": true,
        "disableActiveMonitoringBecauseOfManualIncident": true,
        "serverMonitorSecretKey": true,
        "incomingRequestSecretKey": true,
        "incomingEmailSecretKey": true,
        "incomingEmailMonitorLastEmailReceivedAt": true,
        "incomingEmailMonitorRequest": true,
        "incomingEmailMonitorHeartbeatCheckedAt": true,
        "isAllProbesDisconnectedFromThisMonitor": true,
        "isNoProbeEnabledOnThisMonitor": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect(ctx, "/monitor/" + data.Id.ValueString() + "/get-item", selectParam)
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
    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Labels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Labels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["dependsOnMonitors"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.DependsOnMonitors = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.DependsOnMonitors = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["suppressAlertsWhenParentMonitorStatuses"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.SuppressAlertsWhenParentMonitorStatuses = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.SuppressAlertsWhenParentMonitorStatuses = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["monitorTemplateId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorTemplateId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MonitorTemplateId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MonitorTemplateId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MonitorTemplateId = types.StringValue(string(jsonBytes))
            } else {
                data.MonitorTemplateId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MonitorTemplateId = types.StringValue(string(jsonBytes))
            } else {
                data.MonitorTemplateId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MonitorTemplateId = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorTemplateId = types.StringNull()
        }
    } else if val, ok := dataMap["monitorTemplateId"].(string); ok {
        data.MonitorTemplateId = types.StringValue(val)
    } else {
        data.MonitorTemplateId = types.StringNull()
    }
    if obj, ok := dataMap["monitorType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MonitorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MonitorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MonitorType = types.StringValue(string(jsonBytes))
            } else {
                data.MonitorType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MonitorType = types.StringValue(string(jsonBytes))
            } else {
                data.MonitorType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MonitorType = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorType = types.StringNull()
        }
    } else if val, ok := dataMap["monitorType"].(string); ok {
        data.MonitorType = types.StringValue(val)
    } else {
        data.MonitorType = types.StringNull()
    }
    if obj, ok := dataMap["currentMonitorStatusId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CurrentMonitorStatusId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CurrentMonitorStatusId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CurrentMonitorStatusId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CurrentMonitorStatusId = types.StringValue(string(jsonBytes))
            } else {
                data.CurrentMonitorStatusId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CurrentMonitorStatusId = types.StringValue(string(jsonBytes))
            } else {
                data.CurrentMonitorStatusId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CurrentMonitorStatusId = types.StringValue(string(jsonBytes))
        } else {
            data.CurrentMonitorStatusId = types.StringNull()
        }
    } else if val, ok := dataMap["currentMonitorStatusId"].(string); ok {
        data.CurrentMonitorStatusId = types.StringValue(val)
    } else {
        data.CurrentMonitorStatusId = types.StringNull()
    }
    {
        mappedSteps, stepsDiags := MonitorStepsFromAPI(ctx, dataMap["monitorSteps"])
        resp.Diagnostics.Append(stepsDiags...)
        data.MonitorSteps = MonitorStepsValue{ListValue: mappedSteps}
    }
    if obj, ok := dataMap["monitoringInterval"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitoringInterval = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MonitoringInterval = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MonitoringInterval = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MonitoringInterval = types.StringValue(string(jsonBytes))
            } else {
                data.MonitoringInterval = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MonitoringInterval = types.StringValue(string(jsonBytes))
            } else {
                data.MonitoringInterval = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MonitoringInterval = types.StringValue(string(jsonBytes))
        } else {
            data.MonitoringInterval = types.StringNull()
        }
    } else if val, ok := dataMap["monitoringInterval"].(string); ok {
        data.MonitoringInterval = types.StringValue(val)
    } else {
        data.MonitoringInterval = types.StringNull()
    }
    if obj, ok := dataMap["customFields"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CustomFields = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CustomFields = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CustomFields = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CustomFields = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CustomFields = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CustomFields = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CustomFields = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CustomFields = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CustomFields = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["customFields"].(string); ok {
        data.CustomFields = NewJSONSubsetValue(val)
    } else {
        data.CustomFields = NewJSONSubsetNull()
    }
    if val, ok := dataMap["disableActiveMonitoring"].(bool); ok {
        data.DisableActiveMonitoring = types.BoolValue(val)
    }
    if obj, ok := dataMap["incomingRequestMonitorHeartbeatCheckedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.IncomingRequestMonitorHeartbeatCheckedAt = NewRFC3339Value(val)
        } else {
            data.IncomingRequestMonitorHeartbeatCheckedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["incomingRequestMonitorHeartbeatCheckedAt"].(string); ok && val != "" {
        data.IncomingRequestMonitorHeartbeatCheckedAt = NewRFC3339Value(val)
    } else {
        data.IncomingRequestMonitorHeartbeatCheckedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["telemetryMonitorNextMonitorAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.TelemetryMonitorNextMonitorAt = NewRFC3339Value(val)
        } else {
            data.TelemetryMonitorNextMonitorAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["telemetryMonitorNextMonitorAt"].(string); ok && val != "" {
        data.TelemetryMonitorNextMonitorAt = NewRFC3339Value(val)
    } else {
        data.TelemetryMonitorNextMonitorAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["telemetryMonitorLastMonitorAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.TelemetryMonitorLastMonitorAt = NewRFC3339Value(val)
        } else {
            data.TelemetryMonitorLastMonitorAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["telemetryMonitorLastMonitorAt"].(string); ok && val != "" {
        data.TelemetryMonitorLastMonitorAt = NewRFC3339Value(val)
    } else {
        data.TelemetryMonitorLastMonitorAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["serverMonitorRequestReceivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ServerMonitorRequestReceivedAt = NewRFC3339Value(val)
        } else {
            data.ServerMonitorRequestReceivedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["serverMonitorRequestReceivedAt"].(string); ok && val != "" {
        data.ServerMonitorRequestReceivedAt = NewRFC3339Value(val)
    } else {
        data.ServerMonitorRequestReceivedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["incomingMonitorRequest"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncomingMonitorRequest = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IncomingMonitorRequest = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IncomingMonitorRequest = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.IncomingMonitorRequest = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.IncomingMonitorRequest = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.IncomingMonitorRequest = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.IncomingMonitorRequest = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IncomingMonitorRequest = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.IncomingMonitorRequest = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["incomingMonitorRequest"].(string); ok {
        data.IncomingMonitorRequest = NewJSONSubsetValue(val)
    } else {
        data.IncomingMonitorRequest = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["serverMonitorResponse"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServerMonitorResponse = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ServerMonitorResponse = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ServerMonitorResponse = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ServerMonitorResponse = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.ServerMonitorResponse = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ServerMonitorResponse = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.ServerMonitorResponse = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ServerMonitorResponse = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.ServerMonitorResponse = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["serverMonitorResponse"].(string); ok {
        data.ServerMonitorResponse = NewJSONSubsetValue(val)
    } else {
        data.ServerMonitorResponse = NewJSONSubsetNull()
    }
    if val, ok := dataMap["minimumProbeAgreement"].(float64); ok {
        data.MinimumProbeAgreement = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["minimumProbeAgreement"].(int); ok {
        data.MinimumProbeAgreement = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["minimumProbeAgreement"].(int64); ok {
        data.MinimumProbeAgreement = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["minimumProbeAgreement"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.MinimumProbeAgreement = types.NumberValue(big.NewFloat(val))
        } else {
            data.MinimumProbeAgreement = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.MinimumProbeAgreement = types.NumberNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedAt = NewRFC3339Value(val)
        } else {
            data.CreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = NewRFC3339Value(val)
    } else {
        data.CreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.UpdatedAt = NewRFC3339Value(val)
        } else {
            data.UpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = NewRFC3339Value(val)
    } else {
        data.UpdatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.DeletedAt = NewRFC3339Value(val)
        } else {
            data.DeletedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = NewRFC3339Value(val)
    } else {
        data.DeletedAt = NewRFC3339Null()
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["version"].(int); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["version"].(int64); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["version"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.Version = types.NumberNull()
    }
    if obj, ok := dataMap["slug"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := dataMap["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if val, ok := dataMap["isOwnerNotifiedOfResourceCreation"].(bool); ok {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolValue(val)
    }
    if val, ok := dataMap["disableActiveMonitoringBecauseOfScheduledMaintenanceEvent"].(bool); ok {
        data.DisableActiveMonitoringBecauseOfScheduledMaintenanceEvent = types.BoolValue(val)
    }
    if val, ok := dataMap["disableActiveMonitoringBecauseOfManualIncident"].(bool); ok {
        data.DisableActiveMonitoringBecauseOfManualIncident = types.BoolValue(val)
    }
    if obj, ok := dataMap["serverMonitorSecretKey"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServerMonitorSecretKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ServerMonitorSecretKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ServerMonitorSecretKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ServerMonitorSecretKey = types.StringValue(string(jsonBytes))
            } else {
                data.ServerMonitorSecretKey = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ServerMonitorSecretKey = types.StringValue(string(jsonBytes))
            } else {
                data.ServerMonitorSecretKey = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ServerMonitorSecretKey = types.StringValue(string(jsonBytes))
        } else {
            data.ServerMonitorSecretKey = types.StringNull()
        }
    } else if val, ok := dataMap["serverMonitorSecretKey"].(string); ok {
        data.ServerMonitorSecretKey = types.StringValue(val)
    } else {
        data.ServerMonitorSecretKey = types.StringNull()
    }
    if obj, ok := dataMap["incomingRequestSecretKey"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncomingRequestSecretKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IncomingRequestSecretKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IncomingRequestSecretKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.IncomingRequestSecretKey = types.StringValue(string(jsonBytes))
            } else {
                data.IncomingRequestSecretKey = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.IncomingRequestSecretKey = types.StringValue(string(jsonBytes))
            } else {
                data.IncomingRequestSecretKey = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IncomingRequestSecretKey = types.StringValue(string(jsonBytes))
        } else {
            data.IncomingRequestSecretKey = types.StringNull()
        }
    } else if val, ok := dataMap["incomingRequestSecretKey"].(string); ok {
        data.IncomingRequestSecretKey = types.StringValue(val)
    } else {
        data.IncomingRequestSecretKey = types.StringNull()
    }
    if obj, ok := dataMap["incomingEmailSecretKey"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncomingEmailSecretKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IncomingEmailSecretKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IncomingEmailSecretKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.IncomingEmailSecretKey = types.StringValue(string(jsonBytes))
            } else {
                data.IncomingEmailSecretKey = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.IncomingEmailSecretKey = types.StringValue(string(jsonBytes))
            } else {
                data.IncomingEmailSecretKey = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IncomingEmailSecretKey = types.StringValue(string(jsonBytes))
        } else {
            data.IncomingEmailSecretKey = types.StringNull()
        }
    } else if val, ok := dataMap["incomingEmailSecretKey"].(string); ok {
        data.IncomingEmailSecretKey = types.StringValue(val)
    } else {
        data.IncomingEmailSecretKey = types.StringNull()
    }
    if obj, ok := dataMap["incomingEmailMonitorLastEmailReceivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.IncomingEmailMonitorLastEmailReceivedAt = NewRFC3339Value(val)
        } else {
            data.IncomingEmailMonitorLastEmailReceivedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["incomingEmailMonitorLastEmailReceivedAt"].(string); ok && val != "" {
        data.IncomingEmailMonitorLastEmailReceivedAt = NewRFC3339Value(val)
    } else {
        data.IncomingEmailMonitorLastEmailReceivedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["incomingEmailMonitorRequest"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncomingEmailMonitorRequest = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IncomingEmailMonitorRequest = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IncomingEmailMonitorRequest = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.IncomingEmailMonitorRequest = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.IncomingEmailMonitorRequest = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.IncomingEmailMonitorRequest = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.IncomingEmailMonitorRequest = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IncomingEmailMonitorRequest = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.IncomingEmailMonitorRequest = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["incomingEmailMonitorRequest"].(string); ok {
        data.IncomingEmailMonitorRequest = NewJSONSubsetValue(val)
    } else {
        data.IncomingEmailMonitorRequest = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["incomingEmailMonitorHeartbeatCheckedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.IncomingEmailMonitorHeartbeatCheckedAt = NewRFC3339Value(val)
        } else {
            data.IncomingEmailMonitorHeartbeatCheckedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["incomingEmailMonitorHeartbeatCheckedAt"].(string); ok && val != "" {
        data.IncomingEmailMonitorHeartbeatCheckedAt = NewRFC3339Value(val)
    } else {
        data.IncomingEmailMonitorHeartbeatCheckedAt = NewRFC3339Null()
    }
    if val, ok := dataMap["isAllProbesDisconnectedFromThisMonitor"].(bool); ok {
        data.IsAllProbesDisconnectedFromThisMonitor = types.BoolValue(val)
    }
    if val, ok := dataMap["isNoProbeEnabledOnThisMonitor"].(bool); ok {
        data.IsNoProbeEnabledOnThisMonitor = types.BoolValue(val)
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
        "data": map[string]interface{}{},
    }
    requestDataMap := monitorRequest["data"].(map[string]interface{})

    if !data.Name.IsUnknown() && !state.Name.IsUnknown() && !data.Name.Equal(state.Name) {
        requestDataMap["name"] = data.Name.ValueString()
    }
    if !data.Description.IsUnknown() && !state.Description.IsUnknown() && !data.Description.Equal(state.Description) {
        requestDataMap["description"] = data.Description.ValueString()
    }
    if !data.Labels.IsUnknown() && !state.Labels.IsUnknown() && !data.Labels.Equal(state.Labels) {
        requestDataMap["labels"] = r.convertTerraformSetToInterface(data.Labels)
    }
    if !data.DependsOnMonitors.IsUnknown() && !state.DependsOnMonitors.IsUnknown() && !data.DependsOnMonitors.Equal(state.DependsOnMonitors) {
        requestDataMap["dependsOnMonitors"] = r.convertTerraformSetToInterface(data.DependsOnMonitors)
    }
    if !data.SuppressAlertsWhenParentMonitorStatuses.IsUnknown() && !state.SuppressAlertsWhenParentMonitorStatuses.IsUnknown() && !data.SuppressAlertsWhenParentMonitorStatuses.Equal(state.SuppressAlertsWhenParentMonitorStatuses) {
        requestDataMap["suppressAlertsWhenParentMonitorStatuses"] = r.convertTerraformSetToInterface(data.SuppressAlertsWhenParentMonitorStatuses)
    }
    if !data.MonitorTemplateId.IsUnknown() && !state.MonitorTemplateId.IsUnknown() && !data.MonitorTemplateId.Equal(state.MonitorTemplateId) {
        requestDataMap["monitorTemplateId"] = data.MonitorTemplateId.ValueString()
    }
    if !data.CurrentMonitorStatusId.IsUnknown() && !state.CurrentMonitorStatusId.IsUnknown() && !data.CurrentMonitorStatusId.Equal(state.CurrentMonitorStatusId) {
        requestDataMap["currentMonitorStatusId"] = data.CurrentMonitorStatusId.ValueString()
    }
    if !data.MonitorSteps.IsUnknown() && !state.MonitorSteps.IsUnknown() && !data.MonitorSteps.Equal(state.MonitorSteps) {
        monitorStepsValue, monitorStepsDiags := MonitorStepsToAPI(ctx, data.MonitorSteps.ListValue)
        resp.Diagnostics.Append(monitorStepsDiags...)
        if resp.Diagnostics.HasError() {
            return
        }
        requestDataMap["monitorSteps"] = monitorStepsValue
    }
    if !data.MonitoringInterval.IsUnknown() && !state.MonitoringInterval.IsUnknown() && !data.MonitoringInterval.Equal(state.MonitoringInterval) {
        requestDataMap["monitoringInterval"] = data.MonitoringInterval.ValueString()
    }
    if !data.CustomFields.IsUnknown() && !state.CustomFields.IsUnknown() && !data.CustomFields.Equal(state.CustomFields) {
        var customfieldsData interface{}
        if err := json.Unmarshal([]byte(data.CustomFields.ValueString()), &customfieldsData); err == nil {
            requestDataMap["customFields"] = customfieldsData
        } else {
            requestDataMap["customFields"] = data.CustomFields.ValueString()
        }
    }
    if !data.DisableActiveMonitoring.IsUnknown() && !state.DisableActiveMonitoring.IsUnknown() && !data.DisableActiveMonitoring.Equal(state.DisableActiveMonitoring) {
        requestDataMap["disableActiveMonitoring"] = data.DisableActiveMonitoring.ValueBool()
    }
    if !data.MinimumProbeAgreement.IsUnknown() && !state.MinimumProbeAgreement.IsUnknown() && !data.MinimumProbeAgreement.Equal(state.MinimumProbeAgreement) {
        requestDataMap["minimumProbeAgreement"] = r.bigFloatToFloat64(data.MinimumProbeAgreement.ValueBigFloat())
    }

    // Only call the API when there are changed fields to send. An empty
    // update body is rejected by the API; state is still refreshed below so
    // this method never writes unverified plan values into state.
    if len(monitorRequest["data"].(map[string]interface{})) > 0 {
        httpResp, err := r.client.Put(ctx, "/monitor/" + data.Id.ValueString() + "", monitorRequest)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update monitor, got error: %s", err))
            return
        }

        // Parse the update response
        var monitorResponse map[string]interface{}
        err = r.client.ParseResponse(httpResp, &monitorResponse)
        if err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to update monitor: %s", err))
            return
        }
        _ = monitorResponse
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "description": true,
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
        "disableActiveMonitoring": true,
        "incomingRequestMonitorHeartbeatCheckedAt": true,
        "telemetryMonitorNextMonitorAt": true,
        "telemetryMonitorLastMonitorAt": true,
        "serverMonitorRequestReceivedAt": true,
        "incomingMonitorRequest": true,
        "serverMonitorResponse": true,
        "minimumProbeAgreement": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "slug": true,
        "isOwnerNotifiedOfResourceCreation": true,
        "disableActiveMonitoringBecauseOfScheduledMaintenanceEvent": true,
        "disableActiveMonitoringBecauseOfManualIncident": true,
        "serverMonitorSecretKey": true,
        "incomingRequestSecretKey": true,
        "incomingEmailSecretKey": true,
        "incomingEmailMonitorLastEmailReceivedAt": true,
        "incomingEmailMonitorRequest": true,
        "incomingEmailMonitorHeartbeatCheckedAt": true,
        "isAllProbesDisconnectedFromThisMonitor": true,
        "isNoProbeEnabledOnThisMonitor": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/monitor/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read monitor after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read monitor after update: %s", err))
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
    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["labels"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.Labels = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.Labels = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["dependsOnMonitors"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.DependsOnMonitors = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.DependsOnMonitors = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if val, ok := dataMap["suppressAlertsWhenParentMonitorStatuses"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.SuppressAlertsWhenParentMonitorStatuses = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.SuppressAlertsWhenParentMonitorStatuses = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["monitorTemplateId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorTemplateId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MonitorTemplateId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MonitorTemplateId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MonitorTemplateId = types.StringValue(string(jsonBytes))
            } else {
                data.MonitorTemplateId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MonitorTemplateId = types.StringValue(string(jsonBytes))
            } else {
                data.MonitorTemplateId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MonitorTemplateId = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorTemplateId = types.StringNull()
        }
    } else if val, ok := dataMap["monitorTemplateId"].(string); ok {
        data.MonitorTemplateId = types.StringValue(val)
    } else {
        data.MonitorTemplateId = types.StringNull()
    }
    if obj, ok := dataMap["monitorType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MonitorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MonitorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MonitorType = types.StringValue(string(jsonBytes))
            } else {
                data.MonitorType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MonitorType = types.StringValue(string(jsonBytes))
            } else {
                data.MonitorType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MonitorType = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorType = types.StringNull()
        }
    } else if val, ok := dataMap["monitorType"].(string); ok {
        data.MonitorType = types.StringValue(val)
    } else {
        data.MonitorType = types.StringNull()
    }
    if obj, ok := dataMap["currentMonitorStatusId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CurrentMonitorStatusId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CurrentMonitorStatusId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CurrentMonitorStatusId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CurrentMonitorStatusId = types.StringValue(string(jsonBytes))
            } else {
                data.CurrentMonitorStatusId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CurrentMonitorStatusId = types.StringValue(string(jsonBytes))
            } else {
                data.CurrentMonitorStatusId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CurrentMonitorStatusId = types.StringValue(string(jsonBytes))
        } else {
            data.CurrentMonitorStatusId = types.StringNull()
        }
    } else if val, ok := dataMap["currentMonitorStatusId"].(string); ok {
        data.CurrentMonitorStatusId = types.StringValue(val)
    } else {
        data.CurrentMonitorStatusId = types.StringNull()
    }
    {
        mappedSteps, stepsDiags := MonitorStepsFromAPI(ctx, dataMap["monitorSteps"])
        resp.Diagnostics.Append(stepsDiags...)
        data.MonitorSteps = MonitorStepsValue{ListValue: mappedSteps}
    }
    if obj, ok := dataMap["monitoringInterval"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitoringInterval = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MonitoringInterval = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MonitoringInterval = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MonitoringInterval = types.StringValue(string(jsonBytes))
            } else {
                data.MonitoringInterval = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MonitoringInterval = types.StringValue(string(jsonBytes))
            } else {
                data.MonitoringInterval = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MonitoringInterval = types.StringValue(string(jsonBytes))
        } else {
            data.MonitoringInterval = types.StringNull()
        }
    } else if val, ok := dataMap["monitoringInterval"].(string); ok {
        data.MonitoringInterval = types.StringValue(val)
    } else {
        data.MonitoringInterval = types.StringNull()
    }
    if obj, ok := dataMap["customFields"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CustomFields = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CustomFields = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CustomFields = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CustomFields = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CustomFields = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CustomFields = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CustomFields = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CustomFields = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CustomFields = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["customFields"].(string); ok {
        data.CustomFields = NewJSONSubsetValue(val)
    } else {
        data.CustomFields = NewJSONSubsetNull()
    }
    if val, ok := dataMap["disableActiveMonitoring"].(bool); ok {
        data.DisableActiveMonitoring = types.BoolValue(val)
    }
    if obj, ok := dataMap["incomingRequestMonitorHeartbeatCheckedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.IncomingRequestMonitorHeartbeatCheckedAt = NewRFC3339Value(val)
        } else {
            data.IncomingRequestMonitorHeartbeatCheckedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["incomingRequestMonitorHeartbeatCheckedAt"].(string); ok && val != "" {
        data.IncomingRequestMonitorHeartbeatCheckedAt = NewRFC3339Value(val)
    } else {
        data.IncomingRequestMonitorHeartbeatCheckedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["telemetryMonitorNextMonitorAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.TelemetryMonitorNextMonitorAt = NewRFC3339Value(val)
        } else {
            data.TelemetryMonitorNextMonitorAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["telemetryMonitorNextMonitorAt"].(string); ok && val != "" {
        data.TelemetryMonitorNextMonitorAt = NewRFC3339Value(val)
    } else {
        data.TelemetryMonitorNextMonitorAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["telemetryMonitorLastMonitorAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.TelemetryMonitorLastMonitorAt = NewRFC3339Value(val)
        } else {
            data.TelemetryMonitorLastMonitorAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["telemetryMonitorLastMonitorAt"].(string); ok && val != "" {
        data.TelemetryMonitorLastMonitorAt = NewRFC3339Value(val)
    } else {
        data.TelemetryMonitorLastMonitorAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["serverMonitorRequestReceivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ServerMonitorRequestReceivedAt = NewRFC3339Value(val)
        } else {
            data.ServerMonitorRequestReceivedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["serverMonitorRequestReceivedAt"].(string); ok && val != "" {
        data.ServerMonitorRequestReceivedAt = NewRFC3339Value(val)
    } else {
        data.ServerMonitorRequestReceivedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["incomingMonitorRequest"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncomingMonitorRequest = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IncomingMonitorRequest = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IncomingMonitorRequest = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.IncomingMonitorRequest = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.IncomingMonitorRequest = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.IncomingMonitorRequest = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.IncomingMonitorRequest = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IncomingMonitorRequest = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.IncomingMonitorRequest = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["incomingMonitorRequest"].(string); ok {
        data.IncomingMonitorRequest = NewJSONSubsetValue(val)
    } else {
        data.IncomingMonitorRequest = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["serverMonitorResponse"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServerMonitorResponse = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ServerMonitorResponse = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ServerMonitorResponse = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ServerMonitorResponse = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.ServerMonitorResponse = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ServerMonitorResponse = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.ServerMonitorResponse = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ServerMonitorResponse = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.ServerMonitorResponse = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["serverMonitorResponse"].(string); ok {
        data.ServerMonitorResponse = NewJSONSubsetValue(val)
    } else {
        data.ServerMonitorResponse = NewJSONSubsetNull()
    }
    if val, ok := dataMap["minimumProbeAgreement"].(float64); ok {
        data.MinimumProbeAgreement = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["minimumProbeAgreement"].(int); ok {
        data.MinimumProbeAgreement = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["minimumProbeAgreement"].(int64); ok {
        data.MinimumProbeAgreement = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["minimumProbeAgreement"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.MinimumProbeAgreement = types.NumberValue(big.NewFloat(val))
        } else {
            data.MinimumProbeAgreement = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.MinimumProbeAgreement = types.NumberNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedAt = NewRFC3339Value(val)
        } else {
            data.CreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = NewRFC3339Value(val)
    } else {
        data.CreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.UpdatedAt = NewRFC3339Value(val)
        } else {
            data.UpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = NewRFC3339Value(val)
    } else {
        data.UpdatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.DeletedAt = NewRFC3339Value(val)
        } else {
            data.DeletedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = NewRFC3339Value(val)
    } else {
        data.DeletedAt = NewRFC3339Null()
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["version"].(int); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["version"].(int64); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["version"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.Version = types.NumberNull()
    }
    if obj, ok := dataMap["slug"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := dataMap["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if val, ok := dataMap["isOwnerNotifiedOfResourceCreation"].(bool); ok {
        data.IsOwnerNotifiedOfResourceCreation = types.BoolValue(val)
    }
    if val, ok := dataMap["disableActiveMonitoringBecauseOfScheduledMaintenanceEvent"].(bool); ok {
        data.DisableActiveMonitoringBecauseOfScheduledMaintenanceEvent = types.BoolValue(val)
    }
    if val, ok := dataMap["disableActiveMonitoringBecauseOfManualIncident"].(bool); ok {
        data.DisableActiveMonitoringBecauseOfManualIncident = types.BoolValue(val)
    }
    if obj, ok := dataMap["serverMonitorSecretKey"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServerMonitorSecretKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ServerMonitorSecretKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ServerMonitorSecretKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ServerMonitorSecretKey = types.StringValue(string(jsonBytes))
            } else {
                data.ServerMonitorSecretKey = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ServerMonitorSecretKey = types.StringValue(string(jsonBytes))
            } else {
                data.ServerMonitorSecretKey = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ServerMonitorSecretKey = types.StringValue(string(jsonBytes))
        } else {
            data.ServerMonitorSecretKey = types.StringNull()
        }
    } else if val, ok := dataMap["serverMonitorSecretKey"].(string); ok {
        data.ServerMonitorSecretKey = types.StringValue(val)
    } else {
        data.ServerMonitorSecretKey = types.StringNull()
    }
    if obj, ok := dataMap["incomingRequestSecretKey"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncomingRequestSecretKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IncomingRequestSecretKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IncomingRequestSecretKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.IncomingRequestSecretKey = types.StringValue(string(jsonBytes))
            } else {
                data.IncomingRequestSecretKey = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.IncomingRequestSecretKey = types.StringValue(string(jsonBytes))
            } else {
                data.IncomingRequestSecretKey = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IncomingRequestSecretKey = types.StringValue(string(jsonBytes))
        } else {
            data.IncomingRequestSecretKey = types.StringNull()
        }
    } else if val, ok := dataMap["incomingRequestSecretKey"].(string); ok {
        data.IncomingRequestSecretKey = types.StringValue(val)
    } else {
        data.IncomingRequestSecretKey = types.StringNull()
    }
    if obj, ok := dataMap["incomingEmailSecretKey"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncomingEmailSecretKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IncomingEmailSecretKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IncomingEmailSecretKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.IncomingEmailSecretKey = types.StringValue(string(jsonBytes))
            } else {
                data.IncomingEmailSecretKey = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.IncomingEmailSecretKey = types.StringValue(string(jsonBytes))
            } else {
                data.IncomingEmailSecretKey = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IncomingEmailSecretKey = types.StringValue(string(jsonBytes))
        } else {
            data.IncomingEmailSecretKey = types.StringNull()
        }
    } else if val, ok := dataMap["incomingEmailSecretKey"].(string); ok {
        data.IncomingEmailSecretKey = types.StringValue(val)
    } else {
        data.IncomingEmailSecretKey = types.StringNull()
    }
    if obj, ok := dataMap["incomingEmailMonitorLastEmailReceivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.IncomingEmailMonitorLastEmailReceivedAt = NewRFC3339Value(val)
        } else {
            data.IncomingEmailMonitorLastEmailReceivedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["incomingEmailMonitorLastEmailReceivedAt"].(string); ok && val != "" {
        data.IncomingEmailMonitorLastEmailReceivedAt = NewRFC3339Value(val)
    } else {
        data.IncomingEmailMonitorLastEmailReceivedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["incomingEmailMonitorRequest"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncomingEmailMonitorRequest = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IncomingEmailMonitorRequest = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IncomingEmailMonitorRequest = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.IncomingEmailMonitorRequest = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.IncomingEmailMonitorRequest = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.IncomingEmailMonitorRequest = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.IncomingEmailMonitorRequest = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IncomingEmailMonitorRequest = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.IncomingEmailMonitorRequest = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["incomingEmailMonitorRequest"].(string); ok {
        data.IncomingEmailMonitorRequest = NewJSONSubsetValue(val)
    } else {
        data.IncomingEmailMonitorRequest = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["incomingEmailMonitorHeartbeatCheckedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.IncomingEmailMonitorHeartbeatCheckedAt = NewRFC3339Value(val)
        } else {
            data.IncomingEmailMonitorHeartbeatCheckedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["incomingEmailMonitorHeartbeatCheckedAt"].(string); ok && val != "" {
        data.IncomingEmailMonitorHeartbeatCheckedAt = NewRFC3339Value(val)
    } else {
        data.IncomingEmailMonitorHeartbeatCheckedAt = NewRFC3339Null()
    }
    if val, ok := dataMap["isAllProbesDisconnectedFromThisMonitor"].(bool); ok {
        data.IsAllProbesDisconnectedFromThisMonitor = types.BoolValue(val)
    }
    if val, ok := dataMap["isNoProbeEnabledOnThisMonitor"].(bool); ok {
        data.IsNoProbeEnabledOnThisMonitor = types.BoolValue(val)
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    data.Id = state.Id

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
    httpResp, err := r.client.Delete(ctx, "/monitor/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete monitor, got error: %s", err))
        return
    }

    // A failed delete must keep the resource in state — silently dropping it
    // orphans real infrastructure. 404 means it is already gone.
    if httpResp.StatusCode >= 400 && httpResp.StatusCode != http.StatusNotFound {
        err = r.client.ParseResponse(httpResp, nil)
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to delete monitor: %s", err))
        return
    }
    if httpResp.Body != nil {
        httpResp.Body.Close()
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

    var stringList []string
    terraformList.ElementsAs(context.Background(), &stringList, false)

    // Convert string array to OneUptime format with _id fields
    var result []interface{}
    for _, str := range stringList {
        if str != "" {
            result = append(result, map[string]interface{}{
                "_id": str,
            })
        }
    }
    return result
}

// Helper method to convert Terraform set to Go interface{}
func (r *MonitorResource) convertTerraformSetToInterface(terraformSet types.Set) interface{} {
    if terraformSet.IsNull() || terraformSet.IsUnknown() {
        return nil
    }

    var stringList []string
    terraformSet.ElementsAs(context.Background(), &stringList, false)

    // Convert string array to OneUptime format with _id fields
    var result []interface{}
    for _, str := range stringList {
        if str != "" {
            result = append(result, map[string]interface{}{
                "_id": str,
            })
        }
    }
    return result
}


// Helper method to parse JSON field for complex objects
func (r *MonitorResource) parseJSONField(terraformString basetypes.StringValuable) interface{} {
    sv, _ := terraformString.ToStringValue(context.Background())
    if sv.IsNull() || sv.IsUnknown() || sv.ValueString() == "" {
        return nil
    }

    var result interface{}
    if err := json.Unmarshal([]byte(sv.ValueString()), &result); err != nil {
        // If JSON parsing fails, return the raw string
        return sv.ValueString()
    }

    return result
}

// Normalize URL wrapper objects to avoid drift (e.g., trailing slash differences).
func (r *MonitorResource) normalizeURLWrappers(value interface{}) interface{} {
    switch v := value.(type) {
    case map[string]interface{}:
        if typeStr, ok := v["_type"].(string); ok && typeStr == "URL" {
            if val, ok := v["value"].(string); ok {
                v["value"] = r.normalizeURLString(val)
            }
        }
        for key, child := range v {
            v[key] = r.normalizeURLWrappers(child)
        }
        return v
    case []interface{}:
        for i, child := range v {
            v[i] = r.normalizeURLWrappers(child)
        }
        return v
    default:
        return v
    }
}

func (r *MonitorResource) normalizeURLString(value string) string {
    parsed, err := url.Parse(value)
    if err != nil {
        return value
    }
    if parsed.Path == "/" && parsed.RawQuery == "" && parsed.Fragment == "" {
        return strings.TrimSuffix(value, "/")
    }
    return value
}

// Helper method to convert *big.Float to float64 for JSON serialization
func (r *MonitorResource) bigFloatToFloat64(bf *big.Float) interface{} {
    if bf == nil {
        return nil
    }
    f, _ := bf.Float64()
    return f
}

// Helper method to check if a type string is a valid OneUptime ObjectType.
// The registry itself lives in objecttypes.go, shared across the package.
func (r *MonitorResource) isValidOneUptimeObjectType(typeStr string) bool {
    return validOneUptimeObjectTypes[typeStr]
}
