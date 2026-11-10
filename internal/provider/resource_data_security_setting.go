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
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/numberdefault"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/numberplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DataSecuritySettingResource{}
var _ resource.ResourceWithImportState = &DataSecuritySettingResource{}
var _ resource.ResourceWithModifyPlan = &DataSecuritySettingResource{}

func NewDataSecuritySettingResource() resource.Resource {
    return &DataSecuritySettingResource{}
}

// DataSecuritySettingResource defines the resource implementation.
type DataSecuritySettingResource struct {
    client *Client
}

// DataSecuritySettingResourceModel describes the resource data model.
type DataSecuritySettingResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    IsConfigurationScanEnabled types.Bool `tfsdk:"is_configuration_scan_enabled"`
    IsContentScanEnabled types.Bool `tfsdk:"is_content_scan_enabled"`
    IsTelemetryScanEnabled types.Bool `tfsdk:"is_telemetry_scan_enabled"`
    KeepMaskedSamples types.Bool `tfsdk:"keep_masked_samples"`
    OpenCasesForCriticalIssues types.Bool `tfsdk:"open_cases_for_critical_issues"`
    RaiseAlertsForCriticalIssues types.Bool `tfsdk:"raise_alerts_for_critical_issues"`
    AlertSeverityId types.String `tfsdk:"alert_severity_id"`
    DeclareIncidentsForCriticalIssues types.Bool `tfsdk:"declare_incidents_for_critical_issues"`
    IncidentSeverityId types.String `tfsdk:"incident_severity_id"`
    NotifyOwners types.Bool `tfsdk:"notify_owners"`
    MaxObjectsPerScan types.Number `tfsdk:"max_objects_per_scan"`
    MaxBytesPerObject types.Number `tfsdk:"max_bytes_per_object"`
    MaxBytesPerScan types.Number `tfsdk:"max_bytes_per_scan"`
    MaxRowsPerTable types.Number `tfsdk:"max_rows_per_table"`
    MaxTablesPerScan types.Number `tfsdk:"max_tables_per_scan"`
    TelemetrySampleSize types.Number `tfsdk:"telemetry_sample_size"`
    ContentScanIntervalInHours types.Number `tfsdk:"content_scan_interval_in_hours"`
    AllowedValues JSONSubsetValue `tfsdk:"allowed_values"`
    DisabledDataTypeIds JSONSubsetValue `tfsdk:"disabled_data_type_ids"`
    CreatedAt RFC3339Value `tfsdk:"created_at"`
    UpdatedAt RFC3339Value `tfsdk:"updated_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (r *DataSecuritySettingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_data_security_setting"
}

func (r *DataSecuritySettingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = r.schemaDefinition()
}

func (r *DataSecuritySettingResource) schemaDefinition() schema.Schema {
    return schema.Schema{
        MarkdownDescription: "How Data Security runs in a project: what it scans, how much it reads, and what it does when it finds a Critical issue.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of the project this settings belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "is_configuration_scan_enabled": schema.BoolAttribute{
                MarkdownDescription: "Read stores' configuration and access, and run the checks.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "is_content_scan_enabled": schema.BoolAttribute{
                MarkdownDescription: "Sample stores' content for sensitive data, where the content permission was granted.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "is_telemetry_scan_enabled": schema.BoolAttribute{
                MarkdownDescription: "Sample the project's logs, traces, exceptions, replays and AI conversations for sensitive data.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "keep_masked_samples": schema.BoolAttribute{
                MarkdownDescription: "Keep up to three masked samples with each finding. Off keeps counts only.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "open_cases_for_critical_issues": schema.BoolAttribute{
                MarkdownDescription: "Open a security case, grouped by store, when a Critical issue opens.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "raise_alerts_for_critical_issues": schema.BoolAttribute{
                MarkdownDescription: "Raise an alert, which pages whoever is on call, when a Critical issue opens.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "alert_severity_id": schema.StringAttribute{
                MarkdownDescription: "The severity of the alerts raised. The project's most severe when empty. The ID of a `oneuptime_alert_severity`.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: recordIDValidators("oneuptime_alert_severity"),
            },
            "declare_incidents_for_critical_issues": schema.BoolAttribute{
                MarkdownDescription: "Declare an incident, which runs your incident process, when a Critical issue opens.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "incident_severity_id": schema.StringAttribute{
                MarkdownDescription: "The severity of the incidents declared. The project's most severe when empty. The ID of a `oneuptime_incident_severity`.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: recordIDValidators("oneuptime_incident_severity"),
            },
            "notify_owners": schema.BoolAttribute{
                MarkdownDescription: "Tell a store's owners when an issue opens on it.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "max_objects_per_scan": schema.NumberAttribute{
                MarkdownDescription: "The most objects one content scan of a bucket reads.",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(2000)),
                PlanModifiers: []planmodifier.Number{
                    NumberKeepSameAPINumber(),
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "max_bytes_per_object": schema.NumberAttribute{
                MarkdownDescription: "The most bytes read from one object.",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(65536)),
                PlanModifiers: []planmodifier.Number{
                    NumberKeepSameAPINumber(),
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "max_bytes_per_scan": schema.NumberAttribute{
                MarkdownDescription: "The most bytes one content scan reads in all.",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(52428800)),
                PlanModifiers: []planmodifier.Number{
                    NumberKeepSameAPINumber(),
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "max_rows_per_table": schema.NumberAttribute{
                MarkdownDescription: "The most rows read from one table.",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(1000)),
                PlanModifiers: []planmodifier.Number{
                    NumberKeepSameAPINumber(),
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "max_tables_per_scan": schema.NumberAttribute{
                MarkdownDescription: "The most tables one database scan reads.",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(200)),
                PlanModifiers: []planmodifier.Number{
                    NumberKeepSameAPINumber(),
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "telemetry_sample_size": schema.NumberAttribute{
                MarkdownDescription: "How many recent records of each service and signal a telemetry scan reads.",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(500)),
                PlanModifiers: []planmodifier.Number{
                    NumberKeepSameAPINumber(),
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "content_scan_interval_in_hours": schema.NumberAttribute{
                MarkdownDescription: "How often a store's content is sampled again.",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(24)),
                PlanModifiers: []planmodifier.Number{
                    NumberKeepSameAPINumber(),
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "allowed_values": schema.StringAttribute{
                MarkdownDescription: "Values that are never findings, such as your own support address. A JSON value: write it with `jsonencode()`.",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    JSONSubsetKeepImported(),
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "disabled_data_type_ids": schema.StringAttribute{
                MarkdownDescription: "Built-in data types the project's scans do not look for. A JSON value: write it with `jsonencode()`.",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    JSONSubsetKeepImported(),
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was created.",
                CustomType: RFC3339Type{},
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was updated.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the user who created this settings. The ID of a `oneuptime_user` (see the data source).",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
        },
    }
}

// ModifyPlan keeps the state when nothing a configuration sets changes (see
// importsettle.go).
func (r *DataSecuritySettingResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
    // Destroying: nothing to plan.
    if req.Plan.Raw.IsNull() {
        return
    }

    keepStateWhenNothingChanges(ctx, req, resp)
}

func (r *DataSecuritySettingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *DataSecuritySettingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data DataSecuritySettingResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // What the configuration sets, and what Terraform planned (see keepPlannedValues).
    var config DataSecuritySettingResourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
    if resp.Diagnostics.HasError() {
        return
    }
    plan := data



    // Create API request body. Unset (null/unknown) optional fields are
    // omitted so server-side defaults apply instead of being overwritten
    // with zero values.
    dataSecuritySettingRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := dataSecuritySettingRequest["data"].(map[string]interface{})

    if !data.IsConfigurationScanEnabled.IsNull() && !data.IsConfigurationScanEnabled.IsUnknown() {
        requestDataMap["isConfigurationScanEnabled"] = data.IsConfigurationScanEnabled.ValueBool()
    }
    if !data.IsContentScanEnabled.IsNull() && !data.IsContentScanEnabled.IsUnknown() {
        requestDataMap["isContentScanEnabled"] = data.IsContentScanEnabled.ValueBool()
    }
    if !data.IsTelemetryScanEnabled.IsNull() && !data.IsTelemetryScanEnabled.IsUnknown() {
        requestDataMap["isTelemetryScanEnabled"] = data.IsTelemetryScanEnabled.ValueBool()
    }
    if !data.KeepMaskedSamples.IsNull() && !data.KeepMaskedSamples.IsUnknown() {
        requestDataMap["keepMaskedSamples"] = data.KeepMaskedSamples.ValueBool()
    }
    if !data.OpenCasesForCriticalIssues.IsNull() && !data.OpenCasesForCriticalIssues.IsUnknown() {
        requestDataMap["openCasesForCriticalIssues"] = data.OpenCasesForCriticalIssues.ValueBool()
    }
    if !data.RaiseAlertsForCriticalIssues.IsNull() && !data.RaiseAlertsForCriticalIssues.IsUnknown() {
        requestDataMap["raiseAlertsForCriticalIssues"] = data.RaiseAlertsForCriticalIssues.ValueBool()
    }
    if !data.AlertSeverityId.IsNull() && !data.AlertSeverityId.IsUnknown() {
        requestDataMap["alertSeverityId"] = data.AlertSeverityId.ValueString()
    }
    if !data.DeclareIncidentsForCriticalIssues.IsNull() && !data.DeclareIncidentsForCriticalIssues.IsUnknown() {
        requestDataMap["declareIncidentsForCriticalIssues"] = data.DeclareIncidentsForCriticalIssues.ValueBool()
    }
    if !data.IncidentSeverityId.IsNull() && !data.IncidentSeverityId.IsUnknown() {
        requestDataMap["incidentSeverityId"] = data.IncidentSeverityId.ValueString()
    }
    if !data.NotifyOwners.IsNull() && !data.NotifyOwners.IsUnknown() {
        requestDataMap["notifyOwners"] = data.NotifyOwners.ValueBool()
    }
    if !data.MaxObjectsPerScan.IsNull() && !data.MaxObjectsPerScan.IsUnknown() {
        requestDataMap["maxObjectsPerScan"] = r.bigFloatToFloat64(data.MaxObjectsPerScan.ValueBigFloat())
    }
    if !data.MaxBytesPerObject.IsNull() && !data.MaxBytesPerObject.IsUnknown() {
        requestDataMap["maxBytesPerObject"] = r.bigFloatToFloat64(data.MaxBytesPerObject.ValueBigFloat())
    }
    if !data.MaxBytesPerScan.IsNull() && !data.MaxBytesPerScan.IsUnknown() {
        requestDataMap["maxBytesPerScan"] = r.bigFloatToFloat64(data.MaxBytesPerScan.ValueBigFloat())
    }
    if !data.MaxRowsPerTable.IsNull() && !data.MaxRowsPerTable.IsUnknown() {
        requestDataMap["maxRowsPerTable"] = r.bigFloatToFloat64(data.MaxRowsPerTable.ValueBigFloat())
    }
    if !data.MaxTablesPerScan.IsNull() && !data.MaxTablesPerScan.IsUnknown() {
        requestDataMap["maxTablesPerScan"] = r.bigFloatToFloat64(data.MaxTablesPerScan.ValueBigFloat())
    }
    if !data.TelemetrySampleSize.IsNull() && !data.TelemetrySampleSize.IsUnknown() {
        requestDataMap["telemetrySampleSize"] = r.bigFloatToFloat64(data.TelemetrySampleSize.ValueBigFloat())
    }
    if !data.ContentScanIntervalInHours.IsNull() && !data.ContentScanIntervalInHours.IsUnknown() {
        requestDataMap["contentScanIntervalInHours"] = r.bigFloatToFloat64(data.ContentScanIntervalInHours.ValueBigFloat())
    }
    if parsedAllowedValues := r.parseJSONField(data.AllowedValues); parsedAllowedValues != nil {
        requestDataMap["allowedValues"] = parsedAllowedValues
    }
    if parsedDisabledDataTypeIds := r.parseJSONField(data.DisabledDataTypeIds); parsedDisabledDataTypeIds != nil {
        requestDataMap["disabledDataTypeIds"] = parsedDisabledDataTypeIds
    }

    // Make API call
    httpResp, err := r.client.Post(ctx, "/data-security-setting", dataSecuritySettingRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create data_security_setting, got error: %s", err))
        return
    }

    var dataSecuritySettingResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &dataSecuritySettingResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to create data_security_setting: %s", err))
        return
    }

    // Extract the new resource id from the create response.
    createdId := ""
    if wrapper, ok := dataSecuritySettingResponse["data"].(map[string]interface{}); ok {
        if val, ok := wrapper["_id"].(string); ok {
            createdId = val
        }
    } else if val, ok := dataSecuritySettingResponse["_id"].(string); ok {
        createdId = val
    }
    if createdId == "" {
        resp.Diagnostics.AddError("OneUptime API Error", "Create response for data_security_setting did not contain an id. This is a bug in the provider or the API; please report it.")
        return
    }
    data.Id = types.StringValue(createdId)

    /*
     * The server has committed the row. Persist what we know to state BEFORE
     * the read-back: if the read-back fails and we return without setting
     * state, Terraform never learns the resource exists and the created
     * data_security_setting is orphaned server-side — never refreshed, never
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
        "isConfigurationScanEnabled": true,
        "isContentScanEnabled": true,
        "isTelemetryScanEnabled": true,
        "keepMaskedSamples": true,
        "openCasesForCriticalIssues": true,
        "raiseAlertsForCriticalIssues": true,
        "alertSeverityId": true,
        "declareIncidentsForCriticalIssues": true,
        "incidentSeverityId": true,
        "notifyOwners": true,
        "maxObjectsPerScan": true,
        "maxBytesPerObject": true,
        "maxBytesPerScan": true,
        "maxRowsPerTable": true,
        "maxTablesPerScan": true,
        "telemetrySampleSize": true,
        "contentScanIntervalInHours": true,
        "allowedValues": true,
        "disabledDataTypeIds": true,
        "createdAt": true,
        "updatedAt": true,
        "createdByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/data-security-setting/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        /*
         * State already owns the id, so the resource is tracked and the next
         * refresh reconciles the remaining attributes. Warn rather than
         * error: erroring here would strand a real resource.
         */
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created data_security_setting but could not read it back; state is incomplete until the next refresh: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created data_security_setting but could not parse the read-back response; state is incomplete until the next refresh: %s", err))
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
    if val, ok := dataMap["isConfigurationScanEnabled"].(bool); ok {
        data.IsConfigurationScanEnabled = types.BoolValue(val)
    }
    if val, ok := dataMap["isContentScanEnabled"].(bool); ok {
        data.IsContentScanEnabled = types.BoolValue(val)
    }
    if val, ok := dataMap["isTelemetryScanEnabled"].(bool); ok {
        data.IsTelemetryScanEnabled = types.BoolValue(val)
    }
    if val, ok := dataMap["keepMaskedSamples"].(bool); ok {
        data.KeepMaskedSamples = types.BoolValue(val)
    }
    if val, ok := dataMap["openCasesForCriticalIssues"].(bool); ok {
        data.OpenCasesForCriticalIssues = types.BoolValue(val)
    }
    if val, ok := dataMap["raiseAlertsForCriticalIssues"].(bool); ok {
        data.RaiseAlertsForCriticalIssues = types.BoolValue(val)
    }
    if obj, ok := dataMap["alertSeverityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AlertSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AlertSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AlertSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.AlertSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AlertSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.AlertSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AlertSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.AlertSeverityId = types.StringNull()
        }
    } else if arr, ok := dataMap["alertSeverityId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AlertSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.AlertSeverityId = types.StringNull()
        }
    } else if val, ok := dataMap["alertSeverityId"].(string); ok {
        data.AlertSeverityId = types.StringValue(val)
    } else {
        data.AlertSeverityId = types.StringNull()
    }
    if val, ok := dataMap["declareIncidentsForCriticalIssues"].(bool); ok {
        data.DeclareIncidentsForCriticalIssues = types.BoolValue(val)
    }
    if obj, ok := dataMap["incidentSeverityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IncidentSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.IncidentSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.IncidentSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IncidentSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentSeverityId = types.StringNull()
        }
    } else if arr, ok := dataMap["incidentSeverityId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.IncidentSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentSeverityId = types.StringNull()
        }
    } else if val, ok := dataMap["incidentSeverityId"].(string); ok {
        data.IncidentSeverityId = types.StringValue(val)
    } else {
        data.IncidentSeverityId = types.StringNull()
    }
    if val, ok := dataMap["notifyOwners"].(bool); ok {
        data.NotifyOwners = types.BoolValue(val)
    }
    data.MaxObjectsPerScan = numberFromAPI(data.MaxObjectsPerScan, dataMap["maxObjectsPerScan"])
    data.MaxBytesPerObject = numberFromAPI(data.MaxBytesPerObject, dataMap["maxBytesPerObject"])
    data.MaxBytesPerScan = numberFromAPI(data.MaxBytesPerScan, dataMap["maxBytesPerScan"])
    data.MaxRowsPerTable = numberFromAPI(data.MaxRowsPerTable, dataMap["maxRowsPerTable"])
    data.MaxTablesPerScan = numberFromAPI(data.MaxTablesPerScan, dataMap["maxTablesPerScan"])
    data.TelemetrySampleSize = numberFromAPI(data.TelemetrySampleSize, dataMap["telemetrySampleSize"])
    data.ContentScanIntervalInHours = numberFromAPI(data.ContentScanIntervalInHours, dataMap["contentScanIntervalInHours"])
    if obj, ok := dataMap["allowedValues"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AllowedValues = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AllowedValues = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AllowedValues = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AllowedValues = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.AllowedValues = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AllowedValues = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.AllowedValues = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AllowedValues = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.AllowedValues = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["allowedValues"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AllowedValues = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.AllowedValues = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["allowedValues"].(string); ok {
        data.AllowedValues = NewJSONSubsetValue(val)
    } else {
        data.AllowedValues = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["disabledDataTypeIds"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DisabledDataTypeIds = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DisabledDataTypeIds = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DisabledDataTypeIds = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DisabledDataTypeIds = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DisabledDataTypeIds = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DisabledDataTypeIds = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DisabledDataTypeIds = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DisabledDataTypeIds = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DisabledDataTypeIds = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["disabledDataTypeIds"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DisabledDataTypeIds = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DisabledDataTypeIds = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["disabledDataTypeIds"].(string); ok {
        data.DisabledDataTypeIds = NewJSONSubsetValue(val)
    } else {
        data.DisabledDataTypeIds = NewJSONSubsetNull()
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
    } else if arr, ok := dataMap["createdByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    // The read response is authoritative, but never let it clobber the id we just received.
    data.Id = types.StringValue(createdId)

    // Unconfigured attributes keep their planned value; see keepPlannedValues.
    r.keepPlannedValues(&data, &plan, &config)

    // Write logs using the tflog package
    tflog.Trace(ctx, "created a resource")


    // Nothing a create made was imported (see importsettle.go).
    resp.Diagnostics.Append(forgetImportedFromServer(ctx, resp.Private)...)
    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DataSecuritySettingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data DataSecuritySettingResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "isConfigurationScanEnabled": true,
        "isContentScanEnabled": true,
        "isTelemetryScanEnabled": true,
        "keepMaskedSamples": true,
        "openCasesForCriticalIssues": true,
        "raiseAlertsForCriticalIssues": true,
        "alertSeverityId": true,
        "declareIncidentsForCriticalIssues": true,
        "incidentSeverityId": true,
        "notifyOwners": true,
        "maxObjectsPerScan": true,
        "maxBytesPerObject": true,
        "maxBytesPerScan": true,
        "maxRowsPerTable": true,
        "maxTablesPerScan": true,
        "telemetrySampleSize": true,
        "contentScanIntervalInHours": true,
        "allowedValues": true,
        "disabledDataTypeIds": true,
        "createdAt": true,
        "updatedAt": true,
        "createdByUserId": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect(ctx, "/data-security-setting/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read data_security_setting, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var dataSecuritySettingResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &dataSecuritySettingResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse data_security_setting response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := dataSecuritySettingResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = dataSecuritySettingResponse
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
    if val, ok := dataMap["isConfigurationScanEnabled"].(bool); ok {
        data.IsConfigurationScanEnabled = types.BoolValue(val)
    }
    if val, ok := dataMap["isContentScanEnabled"].(bool); ok {
        data.IsContentScanEnabled = types.BoolValue(val)
    }
    if val, ok := dataMap["isTelemetryScanEnabled"].(bool); ok {
        data.IsTelemetryScanEnabled = types.BoolValue(val)
    }
    if val, ok := dataMap["keepMaskedSamples"].(bool); ok {
        data.KeepMaskedSamples = types.BoolValue(val)
    }
    if val, ok := dataMap["openCasesForCriticalIssues"].(bool); ok {
        data.OpenCasesForCriticalIssues = types.BoolValue(val)
    }
    if val, ok := dataMap["raiseAlertsForCriticalIssues"].(bool); ok {
        data.RaiseAlertsForCriticalIssues = types.BoolValue(val)
    }
    if obj, ok := dataMap["alertSeverityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AlertSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AlertSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AlertSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.AlertSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AlertSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.AlertSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AlertSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.AlertSeverityId = types.StringNull()
        }
    } else if arr, ok := dataMap["alertSeverityId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AlertSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.AlertSeverityId = types.StringNull()
        }
    } else if val, ok := dataMap["alertSeverityId"].(string); ok {
        data.AlertSeverityId = types.StringValue(val)
    } else {
        data.AlertSeverityId = types.StringNull()
    }
    if val, ok := dataMap["declareIncidentsForCriticalIssues"].(bool); ok {
        data.DeclareIncidentsForCriticalIssues = types.BoolValue(val)
    }
    if obj, ok := dataMap["incidentSeverityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IncidentSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.IncidentSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.IncidentSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IncidentSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentSeverityId = types.StringNull()
        }
    } else if arr, ok := dataMap["incidentSeverityId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.IncidentSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentSeverityId = types.StringNull()
        }
    } else if val, ok := dataMap["incidentSeverityId"].(string); ok {
        data.IncidentSeverityId = types.StringValue(val)
    } else {
        data.IncidentSeverityId = types.StringNull()
    }
    if val, ok := dataMap["notifyOwners"].(bool); ok {
        data.NotifyOwners = types.BoolValue(val)
    }
    data.MaxObjectsPerScan = numberFromAPI(data.MaxObjectsPerScan, dataMap["maxObjectsPerScan"])
    data.MaxBytesPerObject = numberFromAPI(data.MaxBytesPerObject, dataMap["maxBytesPerObject"])
    data.MaxBytesPerScan = numberFromAPI(data.MaxBytesPerScan, dataMap["maxBytesPerScan"])
    data.MaxRowsPerTable = numberFromAPI(data.MaxRowsPerTable, dataMap["maxRowsPerTable"])
    data.MaxTablesPerScan = numberFromAPI(data.MaxTablesPerScan, dataMap["maxTablesPerScan"])
    data.TelemetrySampleSize = numberFromAPI(data.TelemetrySampleSize, dataMap["telemetrySampleSize"])
    data.ContentScanIntervalInHours = numberFromAPI(data.ContentScanIntervalInHours, dataMap["contentScanIntervalInHours"])
    if obj, ok := dataMap["allowedValues"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AllowedValues = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AllowedValues = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AllowedValues = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AllowedValues = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.AllowedValues = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AllowedValues = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.AllowedValues = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AllowedValues = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.AllowedValues = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["allowedValues"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AllowedValues = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.AllowedValues = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["allowedValues"].(string); ok {
        data.AllowedValues = NewJSONSubsetValue(val)
    } else {
        data.AllowedValues = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["disabledDataTypeIds"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DisabledDataTypeIds = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DisabledDataTypeIds = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DisabledDataTypeIds = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DisabledDataTypeIds = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DisabledDataTypeIds = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DisabledDataTypeIds = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DisabledDataTypeIds = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DisabledDataTypeIds = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DisabledDataTypeIds = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["disabledDataTypeIds"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DisabledDataTypeIds = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DisabledDataTypeIds = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["disabledDataTypeIds"].(string); ok {
        data.DisabledDataTypeIds = NewJSONSubsetValue(val)
    } else {
        data.DisabledDataTypeIds = NewJSONSubsetNull()
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
    } else if arr, ok := dataMap["createdByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DataSecuritySettingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data DataSecuritySettingResourceModel
    var state DataSecuritySettingResourceModel

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

    // The attributes this update writes from the configuration (see importsettle.go).
    writtenFromConfig := map[string]bool{
        "allowed_values": !data.AllowedValues.Equal(state.AllowedValues),
        "disabled_data_type_ids": !data.DisabledDataTypeIds.Equal(state.DisabledDataTypeIds),
    }

    // What the configuration sets, and what Terraform planned (see keepPlannedValues).
    var config DataSecuritySettingResourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
    if resp.Diagnostics.HasError() {
        return
    }
    plan := data

    // Create API request body
    dataSecuritySettingRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := dataSecuritySettingRequest["data"].(map[string]interface{})

    if !data.IsConfigurationScanEnabled.IsUnknown() && !state.IsConfigurationScanEnabled.IsUnknown() && !data.IsConfigurationScanEnabled.Equal(state.IsConfigurationScanEnabled) {
        requestDataMap["isConfigurationScanEnabled"] = data.IsConfigurationScanEnabled.ValueBool()
    }
    if !data.IsContentScanEnabled.IsUnknown() && !state.IsContentScanEnabled.IsUnknown() && !data.IsContentScanEnabled.Equal(state.IsContentScanEnabled) {
        requestDataMap["isContentScanEnabled"] = data.IsContentScanEnabled.ValueBool()
    }
    if !data.IsTelemetryScanEnabled.IsUnknown() && !state.IsTelemetryScanEnabled.IsUnknown() && !data.IsTelemetryScanEnabled.Equal(state.IsTelemetryScanEnabled) {
        requestDataMap["isTelemetryScanEnabled"] = data.IsTelemetryScanEnabled.ValueBool()
    }
    if !data.KeepMaskedSamples.IsUnknown() && !state.KeepMaskedSamples.IsUnknown() && !data.KeepMaskedSamples.Equal(state.KeepMaskedSamples) {
        requestDataMap["keepMaskedSamples"] = data.KeepMaskedSamples.ValueBool()
    }
    if !data.OpenCasesForCriticalIssues.IsUnknown() && !state.OpenCasesForCriticalIssues.IsUnknown() && !data.OpenCasesForCriticalIssues.Equal(state.OpenCasesForCriticalIssues) {
        requestDataMap["openCasesForCriticalIssues"] = data.OpenCasesForCriticalIssues.ValueBool()
    }
    if !data.RaiseAlertsForCriticalIssues.IsUnknown() && !state.RaiseAlertsForCriticalIssues.IsUnknown() && !data.RaiseAlertsForCriticalIssues.Equal(state.RaiseAlertsForCriticalIssues) {
        requestDataMap["raiseAlertsForCriticalIssues"] = data.RaiseAlertsForCriticalIssues.ValueBool()
    }
    if !data.AlertSeverityId.IsUnknown() && !state.AlertSeverityId.IsUnknown() && !data.AlertSeverityId.Equal(state.AlertSeverityId) {
        requestDataMap["alertSeverityId"] = data.AlertSeverityId.ValueString()
    }
    if !data.DeclareIncidentsForCriticalIssues.IsUnknown() && !state.DeclareIncidentsForCriticalIssues.IsUnknown() && !data.DeclareIncidentsForCriticalIssues.Equal(state.DeclareIncidentsForCriticalIssues) {
        requestDataMap["declareIncidentsForCriticalIssues"] = data.DeclareIncidentsForCriticalIssues.ValueBool()
    }
    if !data.IncidentSeverityId.IsUnknown() && !state.IncidentSeverityId.IsUnknown() && !data.IncidentSeverityId.Equal(state.IncidentSeverityId) {
        requestDataMap["incidentSeverityId"] = data.IncidentSeverityId.ValueString()
    }
    if !data.NotifyOwners.IsUnknown() && !state.NotifyOwners.IsUnknown() && !data.NotifyOwners.Equal(state.NotifyOwners) {
        requestDataMap["notifyOwners"] = data.NotifyOwners.ValueBool()
    }
    if !data.MaxObjectsPerScan.IsUnknown() && !state.MaxObjectsPerScan.IsUnknown() && !data.MaxObjectsPerScan.Equal(state.MaxObjectsPerScan) {
        requestDataMap["maxObjectsPerScan"] = r.bigFloatToFloat64(data.MaxObjectsPerScan.ValueBigFloat())
    }
    if !data.MaxBytesPerObject.IsUnknown() && !state.MaxBytesPerObject.IsUnknown() && !data.MaxBytesPerObject.Equal(state.MaxBytesPerObject) {
        requestDataMap["maxBytesPerObject"] = r.bigFloatToFloat64(data.MaxBytesPerObject.ValueBigFloat())
    }
    if !data.MaxBytesPerScan.IsUnknown() && !state.MaxBytesPerScan.IsUnknown() && !data.MaxBytesPerScan.Equal(state.MaxBytesPerScan) {
        requestDataMap["maxBytesPerScan"] = r.bigFloatToFloat64(data.MaxBytesPerScan.ValueBigFloat())
    }
    if !data.MaxRowsPerTable.IsUnknown() && !state.MaxRowsPerTable.IsUnknown() && !data.MaxRowsPerTable.Equal(state.MaxRowsPerTable) {
        requestDataMap["maxRowsPerTable"] = r.bigFloatToFloat64(data.MaxRowsPerTable.ValueBigFloat())
    }
    if !data.MaxTablesPerScan.IsUnknown() && !state.MaxTablesPerScan.IsUnknown() && !data.MaxTablesPerScan.Equal(state.MaxTablesPerScan) {
        requestDataMap["maxTablesPerScan"] = r.bigFloatToFloat64(data.MaxTablesPerScan.ValueBigFloat())
    }
    if !data.TelemetrySampleSize.IsUnknown() && !state.TelemetrySampleSize.IsUnknown() && !data.TelemetrySampleSize.Equal(state.TelemetrySampleSize) {
        requestDataMap["telemetrySampleSize"] = r.bigFloatToFloat64(data.TelemetrySampleSize.ValueBigFloat())
    }
    if !data.ContentScanIntervalInHours.IsUnknown() && !state.ContentScanIntervalInHours.IsUnknown() && !data.ContentScanIntervalInHours.Equal(state.ContentScanIntervalInHours) {
        requestDataMap["contentScanIntervalInHours"] = r.bigFloatToFloat64(data.ContentScanIntervalInHours.ValueBigFloat())
    }
    if !data.AllowedValues.IsUnknown() && !state.AllowedValues.IsUnknown() && !data.AllowedValues.Equal(state.AllowedValues) {
        var allowedvaluesData interface{}
        if err := json.Unmarshal([]byte(data.AllowedValues.ValueString()), &allowedvaluesData); err == nil {
            requestDataMap["allowedValues"] = allowedvaluesData
        } else {
            requestDataMap["allowedValues"] = data.AllowedValues.ValueString()
        }
    }
    if !data.DisabledDataTypeIds.IsUnknown() && !state.DisabledDataTypeIds.IsUnknown() && !data.DisabledDataTypeIds.Equal(state.DisabledDataTypeIds) {
        var disableddatatypeidsData interface{}
        if err := json.Unmarshal([]byte(data.DisabledDataTypeIds.ValueString()), &disableddatatypeidsData); err == nil {
            requestDataMap["disabledDataTypeIds"] = disableddatatypeidsData
        } else {
            requestDataMap["disabledDataTypeIds"] = data.DisabledDataTypeIds.ValueString()
        }
    }

    // Only call the API when there are changed fields to send. An empty
    // update body is rejected by the API; state is still refreshed below so
    // this method never writes unverified plan values into state.
    if len(dataSecuritySettingRequest["data"].(map[string]interface{})) > 0 {
        httpResp, err := r.client.Put(ctx, "/data-security-setting/" + data.Id.ValueString() + "", dataSecuritySettingRequest)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update data_security_setting, got error: %s", err))
            return
        }

        // Parse the update response
        var dataSecuritySettingResponse map[string]interface{}
        err = r.client.ParseResponse(httpResp, &dataSecuritySettingResponse)
        if err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to update data_security_setting: %s", err))
            return
        }
        _ = dataSecuritySettingResponse
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "isConfigurationScanEnabled": true,
        "isContentScanEnabled": true,
        "isTelemetryScanEnabled": true,
        "keepMaskedSamples": true,
        "openCasesForCriticalIssues": true,
        "raiseAlertsForCriticalIssues": true,
        "alertSeverityId": true,
        "declareIncidentsForCriticalIssues": true,
        "incidentSeverityId": true,
        "notifyOwners": true,
        "maxObjectsPerScan": true,
        "maxBytesPerObject": true,
        "maxBytesPerScan": true,
        "maxRowsPerTable": true,
        "maxTablesPerScan": true,
        "telemetrySampleSize": true,
        "contentScanIntervalInHours": true,
        "allowedValues": true,
        "disabledDataTypeIds": true,
        "createdAt": true,
        "updatedAt": true,
        "createdByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/data-security-setting/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read data_security_setting after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read data_security_setting after update: %s", err))
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
    if val, ok := dataMap["isConfigurationScanEnabled"].(bool); ok {
        data.IsConfigurationScanEnabled = types.BoolValue(val)
    }
    if val, ok := dataMap["isContentScanEnabled"].(bool); ok {
        data.IsContentScanEnabled = types.BoolValue(val)
    }
    if val, ok := dataMap["isTelemetryScanEnabled"].(bool); ok {
        data.IsTelemetryScanEnabled = types.BoolValue(val)
    }
    if val, ok := dataMap["keepMaskedSamples"].(bool); ok {
        data.KeepMaskedSamples = types.BoolValue(val)
    }
    if val, ok := dataMap["openCasesForCriticalIssues"].(bool); ok {
        data.OpenCasesForCriticalIssues = types.BoolValue(val)
    }
    if val, ok := dataMap["raiseAlertsForCriticalIssues"].(bool); ok {
        data.RaiseAlertsForCriticalIssues = types.BoolValue(val)
    }
    if obj, ok := dataMap["alertSeverityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AlertSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AlertSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AlertSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.AlertSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AlertSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.AlertSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AlertSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.AlertSeverityId = types.StringNull()
        }
    } else if arr, ok := dataMap["alertSeverityId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AlertSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.AlertSeverityId = types.StringNull()
        }
    } else if val, ok := dataMap["alertSeverityId"].(string); ok {
        data.AlertSeverityId = types.StringValue(val)
    } else {
        data.AlertSeverityId = types.StringNull()
    }
    if val, ok := dataMap["declareIncidentsForCriticalIssues"].(bool); ok {
        data.DeclareIncidentsForCriticalIssues = types.BoolValue(val)
    }
    if obj, ok := dataMap["incidentSeverityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IncidentSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.IncidentSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.IncidentSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IncidentSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentSeverityId = types.StringNull()
        }
    } else if arr, ok := dataMap["incidentSeverityId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.IncidentSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentSeverityId = types.StringNull()
        }
    } else if val, ok := dataMap["incidentSeverityId"].(string); ok {
        data.IncidentSeverityId = types.StringValue(val)
    } else {
        data.IncidentSeverityId = types.StringNull()
    }
    if val, ok := dataMap["notifyOwners"].(bool); ok {
        data.NotifyOwners = types.BoolValue(val)
    }
    data.MaxObjectsPerScan = numberFromAPI(data.MaxObjectsPerScan, dataMap["maxObjectsPerScan"])
    data.MaxBytesPerObject = numberFromAPI(data.MaxBytesPerObject, dataMap["maxBytesPerObject"])
    data.MaxBytesPerScan = numberFromAPI(data.MaxBytesPerScan, dataMap["maxBytesPerScan"])
    data.MaxRowsPerTable = numberFromAPI(data.MaxRowsPerTable, dataMap["maxRowsPerTable"])
    data.MaxTablesPerScan = numberFromAPI(data.MaxTablesPerScan, dataMap["maxTablesPerScan"])
    data.TelemetrySampleSize = numberFromAPI(data.TelemetrySampleSize, dataMap["telemetrySampleSize"])
    data.ContentScanIntervalInHours = numberFromAPI(data.ContentScanIntervalInHours, dataMap["contentScanIntervalInHours"])
    if obj, ok := dataMap["allowedValues"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AllowedValues = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AllowedValues = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AllowedValues = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AllowedValues = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.AllowedValues = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AllowedValues = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.AllowedValues = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AllowedValues = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.AllowedValues = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["allowedValues"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AllowedValues = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.AllowedValues = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["allowedValues"].(string); ok {
        data.AllowedValues = NewJSONSubsetValue(val)
    } else {
        data.AllowedValues = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["disabledDataTypeIds"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DisabledDataTypeIds = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DisabledDataTypeIds = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DisabledDataTypeIds = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DisabledDataTypeIds = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DisabledDataTypeIds = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DisabledDataTypeIds = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DisabledDataTypeIds = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DisabledDataTypeIds = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DisabledDataTypeIds = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["disabledDataTypeIds"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.DisabledDataTypeIds = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DisabledDataTypeIds = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["disabledDataTypeIds"].(string); ok {
        data.DisabledDataTypeIds = NewJSONSubsetValue(val)
    } else {
        data.DisabledDataTypeIds = NewJSONSubsetNull()
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
    } else if arr, ok := dataMap["createdByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    data.Id = state.Id

    // Unconfigured attributes keep their planned value; see keepPlannedValues.
    r.keepPlannedValues(&data, &plan, &config)

    // What this update wrote is what the configuration says now (see importsettle.go).
    resp.Diagnostics.Append(keepImportedFromServer(ctx, resp.Private, writtenFromConfig)...)

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DataSecuritySettingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data DataSecuritySettingResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    httpResp, err := r.client.Delete(ctx, "/data-security-setting/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete data_security_setting, got error: %s", err))
        return
    }

    // A failed delete must keep the resource in state — silently dropping it
    // orphans real infrastructure. 404 means it is already gone.
    if httpResp.StatusCode >= 400 && httpResp.StatusCode != http.StatusNotFound {
        err = r.client.ParseResponse(httpResp, nil)
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to delete data_security_setting: %s", err))
        return
    }
    if httpResp.Body != nil {
        httpResp.Body.Close()
    }
}


func (r *DataSecuritySettingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
    // These hold the record as OneUptime keeps it, keys it fills in itself
    // included, until an update writes them (see importsettle.go).
    resp.Diagnostics.Append(markImportedFromServer(ctx, resp.Private, []string{"allowed_values", "disabled_data_type_ids"})...)
}


// keepPlannedValues puts back, after a create or an update, the planned value
// of each optional attribute the configuration leaves out. The server keeps
// some of these up to date on its own (when it last checked a heartbeat, the
// status a probe last reported...), so the value read back after the write can
// already differ from the plan, and Terraform would fail the apply with
// "Provider produced inconsistent result after apply". The next refresh reads
// the server's value, which is never a diff for an attribute nobody configured.
func (r *DataSecuritySettingResource) keepPlannedValues(data *DataSecuritySettingResourceModel, plan *DataSecuritySettingResourceModel, config *DataSecuritySettingResourceModel) {
    if config.IsConfigurationScanEnabled.IsNull() && !plan.IsConfigurationScanEnabled.IsUnknown() {
        data.IsConfigurationScanEnabled = plan.IsConfigurationScanEnabled
    }
    if config.IsContentScanEnabled.IsNull() && !plan.IsContentScanEnabled.IsUnknown() {
        data.IsContentScanEnabled = plan.IsContentScanEnabled
    }
    if config.IsTelemetryScanEnabled.IsNull() && !plan.IsTelemetryScanEnabled.IsUnknown() {
        data.IsTelemetryScanEnabled = plan.IsTelemetryScanEnabled
    }
    if config.KeepMaskedSamples.IsNull() && !plan.KeepMaskedSamples.IsUnknown() {
        data.KeepMaskedSamples = plan.KeepMaskedSamples
    }
    if config.OpenCasesForCriticalIssues.IsNull() && !plan.OpenCasesForCriticalIssues.IsUnknown() {
        data.OpenCasesForCriticalIssues = plan.OpenCasesForCriticalIssues
    }
    if config.RaiseAlertsForCriticalIssues.IsNull() && !plan.RaiseAlertsForCriticalIssues.IsUnknown() {
        data.RaiseAlertsForCriticalIssues = plan.RaiseAlertsForCriticalIssues
    }
    if config.AlertSeverityId.IsNull() && !plan.AlertSeverityId.IsUnknown() {
        data.AlertSeverityId = plan.AlertSeverityId
    }
    if config.DeclareIncidentsForCriticalIssues.IsNull() && !plan.DeclareIncidentsForCriticalIssues.IsUnknown() {
        data.DeclareIncidentsForCriticalIssues = plan.DeclareIncidentsForCriticalIssues
    }
    if config.IncidentSeverityId.IsNull() && !plan.IncidentSeverityId.IsUnknown() {
        data.IncidentSeverityId = plan.IncidentSeverityId
    }
    if config.NotifyOwners.IsNull() && !plan.NotifyOwners.IsUnknown() {
        data.NotifyOwners = plan.NotifyOwners
    }
    if config.MaxObjectsPerScan.IsNull() && !plan.MaxObjectsPerScan.IsUnknown() {
        data.MaxObjectsPerScan = plan.MaxObjectsPerScan
    }
    if config.MaxBytesPerObject.IsNull() && !plan.MaxBytesPerObject.IsUnknown() {
        data.MaxBytesPerObject = plan.MaxBytesPerObject
    }
    if config.MaxBytesPerScan.IsNull() && !plan.MaxBytesPerScan.IsUnknown() {
        data.MaxBytesPerScan = plan.MaxBytesPerScan
    }
    if config.MaxRowsPerTable.IsNull() && !plan.MaxRowsPerTable.IsUnknown() {
        data.MaxRowsPerTable = plan.MaxRowsPerTable
    }
    if config.MaxTablesPerScan.IsNull() && !plan.MaxTablesPerScan.IsUnknown() {
        data.MaxTablesPerScan = plan.MaxTablesPerScan
    }
    if config.TelemetrySampleSize.IsNull() && !plan.TelemetrySampleSize.IsUnknown() {
        data.TelemetrySampleSize = plan.TelemetrySampleSize
    }
    if config.ContentScanIntervalInHours.IsNull() && !plan.ContentScanIntervalInHours.IsUnknown() {
        data.ContentScanIntervalInHours = plan.ContentScanIntervalInHours
    }
    if config.AllowedValues.IsNull() && !plan.AllowedValues.IsUnknown() {
        data.AllowedValues = plan.AllowedValues
    }
    if config.DisabledDataTypeIds.IsNull() && !plan.DisabledDataTypeIds.IsUnknown() {
        data.DisabledDataTypeIds = plan.DisabledDataTypeIds
    }
}

// Helper method to convert Terraform map to Go interface{}
func (r *DataSecuritySettingResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *DataSecuritySettingResource) convertTerraformListToInterface(terraformList types.List) interface{} {
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
func (r *DataSecuritySettingResource) convertTerraformSetToInterface(terraformSet types.Set) interface{} {
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
func (r *DataSecuritySettingResource) parseJSONField(terraformString basetypes.StringValuable) interface{} {
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
func (r *DataSecuritySettingResource) normalizeURLWrappers(value interface{}) interface{} {
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

func (r *DataSecuritySettingResource) normalizeURLString(value string) string {
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
func (r *DataSecuritySettingResource) bigFloatToFloat64(bf *big.Float) interface{} {
    if bf == nil {
        return nil
    }
    f, _ := bf.Float64()
    return f
}

// Helper method to check if a type string is a valid OneUptime ObjectType.
// The registry itself lives in objecttypes.go, shared across the package.
func (r *DataSecuritySettingResource) isValidOneUptimeObjectType(typeStr string) bool {
    return validOneUptimeObjectTypes[typeStr]
}
