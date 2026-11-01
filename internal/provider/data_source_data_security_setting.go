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
var _ datasource.DataSource = &DataSecuritySettingDataSource{}

func NewDataSecuritySettingDataSource() datasource.DataSource {
    return &DataSecuritySettingDataSource{}
}

// DataSecuritySettingDataSource defines the data source implementation.
type DataSecuritySettingDataSource struct {
    client *Client
}

// DataSecuritySettingDataSourceModel describes the data source data model.
type DataSecuritySettingDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
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
    AllowedValues types.String `tfsdk:"allowed_values"`
    DisabledDataTypeIds types.String `tfsdk:"disabled_data_type_ids"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *DataSecuritySettingDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_data_security_setting"
}

func (d *DataSecuritySettingDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "How Data Security runs in a project: what it scans, how much it reads, and what it does when it finds a Critical issue. Look up an existing data security setting by `id`, or by any of its other arguments (`alert_severity_id`, `content_scan_interval_in_hours`, `created_by_user_id`, ...): each one set must match, and exactly one data security setting may match them all.",

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
                MarkdownDescription: "ID of the project this settings belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "is_configuration_scan_enabled": schema.BoolAttribute{
                MarkdownDescription: "Read stores' configuration and access, and run the checks.",
                Optional: true,
                Computed: true,
            },
            "is_content_scan_enabled": schema.BoolAttribute{
                MarkdownDescription: "Sample stores' content for sensitive data, where the content permission was granted.",
                Optional: true,
                Computed: true,
            },
            "is_telemetry_scan_enabled": schema.BoolAttribute{
                MarkdownDescription: "Sample the project's logs, traces, exceptions, replays and AI conversations for sensitive data.",
                Optional: true,
                Computed: true,
            },
            "keep_masked_samples": schema.BoolAttribute{
                MarkdownDescription: "Keep up to three masked samples with each finding. Off keeps counts only.",
                Optional: true,
                Computed: true,
            },
            "open_cases_for_critical_issues": schema.BoolAttribute{
                MarkdownDescription: "Open a security case, grouped by store, when a Critical issue opens.",
                Optional: true,
                Computed: true,
            },
            "raise_alerts_for_critical_issues": schema.BoolAttribute{
                MarkdownDescription: "Raise an alert, which pages whoever is on call, when a Critical issue opens.",
                Optional: true,
                Computed: true,
            },
            "alert_severity_id": schema.StringAttribute{
                MarkdownDescription: "The severity of the alerts raised. The project's most severe when empty. The ID of a `oneuptime_alert_severity`.",
                Optional: true,
                Computed: true,
            },
            "declare_incidents_for_critical_issues": schema.BoolAttribute{
                MarkdownDescription: "Declare an incident, which runs your incident process, when a Critical issue opens.",
                Optional: true,
                Computed: true,
            },
            "incident_severity_id": schema.StringAttribute{
                MarkdownDescription: "The severity of the incidents declared. The project's most severe when empty. The ID of a `oneuptime_incident_severity`.",
                Optional: true,
                Computed: true,
            },
            "notify_owners": schema.BoolAttribute{
                MarkdownDescription: "Tell a store's owners when an issue opens on it.",
                Optional: true,
                Computed: true,
            },
            "max_objects_per_scan": schema.NumberAttribute{
                MarkdownDescription: "The most objects one content scan of a bucket reads.",
                Optional: true,
                Computed: true,
            },
            "max_bytes_per_object": schema.NumberAttribute{
                MarkdownDescription: "The most bytes read from one object.",
                Optional: true,
                Computed: true,
            },
            "max_bytes_per_scan": schema.NumberAttribute{
                MarkdownDescription: "The most bytes one content scan reads in all.",
                Optional: true,
                Computed: true,
            },
            "max_rows_per_table": schema.NumberAttribute{
                MarkdownDescription: "The most rows read from one table.",
                Optional: true,
                Computed: true,
            },
            "max_tables_per_scan": schema.NumberAttribute{
                MarkdownDescription: "The most tables one database scan reads.",
                Optional: true,
                Computed: true,
            },
            "telemetry_sample_size": schema.NumberAttribute{
                MarkdownDescription: "How many recent records of each service and signal a telemetry scan reads.",
                Optional: true,
                Computed: true,
            },
            "content_scan_interval_in_hours": schema.NumberAttribute{
                MarkdownDescription: "How often a store's content is sampled again.",
                Optional: true,
                Computed: true,
            },
            "allowed_values": schema.StringAttribute{
                MarkdownDescription: "Values that are never findings, such as your own support address. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "disabled_data_type_ids": schema.StringAttribute{
                MarkdownDescription: "Built-in data types the project's scans do not look for. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the user who created this settings. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *DataSecuritySettingDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataSecuritySettingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data DataSecuritySettingDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.IsConfigurationScanEnabled.IsNull() && !data.IsConfigurationScanEnabled.IsUnknown() {
        filters["isConfigurationScanEnabled"] = data.IsConfigurationScanEnabled.ValueBool()
        filterNames = append(filterNames, "is_configuration_scan_enabled = "+fmt.Sprintf("%t", data.IsConfigurationScanEnabled.ValueBool()))
    }
    if !data.IsContentScanEnabled.IsNull() && !data.IsContentScanEnabled.IsUnknown() {
        filters["isContentScanEnabled"] = data.IsContentScanEnabled.ValueBool()
        filterNames = append(filterNames, "is_content_scan_enabled = "+fmt.Sprintf("%t", data.IsContentScanEnabled.ValueBool()))
    }
    if !data.IsTelemetryScanEnabled.IsNull() && !data.IsTelemetryScanEnabled.IsUnknown() {
        filters["isTelemetryScanEnabled"] = data.IsTelemetryScanEnabled.ValueBool()
        filterNames = append(filterNames, "is_telemetry_scan_enabled = "+fmt.Sprintf("%t", data.IsTelemetryScanEnabled.ValueBool()))
    }
    if !data.KeepMaskedSamples.IsNull() && !data.KeepMaskedSamples.IsUnknown() {
        filters["keepMaskedSamples"] = data.KeepMaskedSamples.ValueBool()
        filterNames = append(filterNames, "keep_masked_samples = "+fmt.Sprintf("%t", data.KeepMaskedSamples.ValueBool()))
    }
    if !data.OpenCasesForCriticalIssues.IsNull() && !data.OpenCasesForCriticalIssues.IsUnknown() {
        filters["openCasesForCriticalIssues"] = data.OpenCasesForCriticalIssues.ValueBool()
        filterNames = append(filterNames, "open_cases_for_critical_issues = "+fmt.Sprintf("%t", data.OpenCasesForCriticalIssues.ValueBool()))
    }
    if !data.RaiseAlertsForCriticalIssues.IsNull() && !data.RaiseAlertsForCriticalIssues.IsUnknown() {
        filters["raiseAlertsForCriticalIssues"] = data.RaiseAlertsForCriticalIssues.ValueBool()
        filterNames = append(filterNames, "raise_alerts_for_critical_issues = "+fmt.Sprintf("%t", data.RaiseAlertsForCriticalIssues.ValueBool()))
    }
    if !data.AlertSeverityId.IsNull() && !data.AlertSeverityId.IsUnknown() {
        filters["alertSeverityId"] = data.AlertSeverityId.ValueString()
        filterNames = append(filterNames, "alert_severity_id = "+fmt.Sprintf("%q", data.AlertSeverityId.ValueString()))
    }
    if !data.DeclareIncidentsForCriticalIssues.IsNull() && !data.DeclareIncidentsForCriticalIssues.IsUnknown() {
        filters["declareIncidentsForCriticalIssues"] = data.DeclareIncidentsForCriticalIssues.ValueBool()
        filterNames = append(filterNames, "declare_incidents_for_critical_issues = "+fmt.Sprintf("%t", data.DeclareIncidentsForCriticalIssues.ValueBool()))
    }
    if !data.IncidentSeverityId.IsNull() && !data.IncidentSeverityId.IsUnknown() {
        filters["incidentSeverityId"] = data.IncidentSeverityId.ValueString()
        filterNames = append(filterNames, "incident_severity_id = "+fmt.Sprintf("%q", data.IncidentSeverityId.ValueString()))
    }
    if !data.NotifyOwners.IsNull() && !data.NotifyOwners.IsUnknown() {
        filters["notifyOwners"] = data.NotifyOwners.ValueBool()
        filterNames = append(filterNames, "notify_owners = "+fmt.Sprintf("%t", data.NotifyOwners.ValueBool()))
    }
    if !data.MaxObjectsPerScan.IsNull() && !data.MaxObjectsPerScan.IsUnknown() {
        filters["maxObjectsPerScan"] = lookupNumber(data.MaxObjectsPerScan)
        filterNames = append(filterNames, "max_objects_per_scan = "+data.MaxObjectsPerScan.ValueBigFloat().String())
    }
    if !data.MaxBytesPerObject.IsNull() && !data.MaxBytesPerObject.IsUnknown() {
        filters["maxBytesPerObject"] = lookupNumber(data.MaxBytesPerObject)
        filterNames = append(filterNames, "max_bytes_per_object = "+data.MaxBytesPerObject.ValueBigFloat().String())
    }
    if !data.MaxBytesPerScan.IsNull() && !data.MaxBytesPerScan.IsUnknown() {
        filters["maxBytesPerScan"] = lookupNumber(data.MaxBytesPerScan)
        filterNames = append(filterNames, "max_bytes_per_scan = "+data.MaxBytesPerScan.ValueBigFloat().String())
    }
    if !data.MaxRowsPerTable.IsNull() && !data.MaxRowsPerTable.IsUnknown() {
        filters["maxRowsPerTable"] = lookupNumber(data.MaxRowsPerTable)
        filterNames = append(filterNames, "max_rows_per_table = "+data.MaxRowsPerTable.ValueBigFloat().String())
    }
    if !data.MaxTablesPerScan.IsNull() && !data.MaxTablesPerScan.IsUnknown() {
        filters["maxTablesPerScan"] = lookupNumber(data.MaxTablesPerScan)
        filterNames = append(filterNames, "max_tables_per_scan = "+data.MaxTablesPerScan.ValueBigFloat().String())
    }
    if !data.TelemetrySampleSize.IsNull() && !data.TelemetrySampleSize.IsUnknown() {
        filters["telemetrySampleSize"] = lookupNumber(data.TelemetrySampleSize)
        filterNames = append(filterNames, "telemetry_sample_size = "+data.TelemetrySampleSize.ValueBigFloat().String())
    }
    if !data.ContentScanIntervalInHours.IsNull() && !data.ContentScanIntervalInHours.IsUnknown() {
        filters["contentScanIntervalInHours"] = lookupNumber(data.ContentScanIntervalInHours)
        filterNames = append(filterNames, "content_scan_interval_in_hours = "+data.ContentScanIntervalInHours.ValueBigFloat().String())
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the data security setting up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the data security setting up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
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
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/data-security-setting/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read data_security_setting, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data security setting found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read data_security_setting: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/data-security-setting/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list data_security_setting, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list data_security_setting: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data security setting matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one data security setting matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for data_security_setting.")
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
    if val, ok := item["isConfigurationScanEnabled"].(bool); ok {
        data.IsConfigurationScanEnabled = types.BoolValue(val)
    } else {
        data.IsConfigurationScanEnabled = types.BoolNull()
    }
    if val, ok := item["isContentScanEnabled"].(bool); ok {
        data.IsContentScanEnabled = types.BoolValue(val)
    } else {
        data.IsContentScanEnabled = types.BoolNull()
    }
    if val, ok := item["isTelemetryScanEnabled"].(bool); ok {
        data.IsTelemetryScanEnabled = types.BoolValue(val)
    } else {
        data.IsTelemetryScanEnabled = types.BoolNull()
    }
    if val, ok := item["keepMaskedSamples"].(bool); ok {
        data.KeepMaskedSamples = types.BoolValue(val)
    } else {
        data.KeepMaskedSamples = types.BoolNull()
    }
    if val, ok := item["openCasesForCriticalIssues"].(bool); ok {
        data.OpenCasesForCriticalIssues = types.BoolValue(val)
    } else {
        data.OpenCasesForCriticalIssues = types.BoolNull()
    }
    if val, ok := item["raiseAlertsForCriticalIssues"].(bool); ok {
        data.RaiseAlertsForCriticalIssues = types.BoolValue(val)
    } else {
        data.RaiseAlertsForCriticalIssues = types.BoolNull()
    }
    if obj, ok := item["alertSeverityId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AlertSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AlertSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AlertSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.AlertSeverityId = types.StringNull()
        }
    } else if val, ok := item["alertSeverityId"].(string); ok {
        data.AlertSeverityId = types.StringValue(val)
    } else {
        data.AlertSeverityId = types.StringNull()
    }
    if val, ok := item["declareIncidentsForCriticalIssues"].(bool); ok {
        data.DeclareIncidentsForCriticalIssues = types.BoolValue(val)
    } else {
        data.DeclareIncidentsForCriticalIssues = types.BoolNull()
    }
    if obj, ok := item["incidentSeverityId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncidentSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncidentSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentSeverityId = types.StringNull()
        }
    } else if val, ok := item["incidentSeverityId"].(string); ok {
        data.IncidentSeverityId = types.StringValue(val)
    } else {
        data.IncidentSeverityId = types.StringNull()
    }
    if val, ok := item["notifyOwners"].(bool); ok {
        data.NotifyOwners = types.BoolValue(val)
    } else {
        data.NotifyOwners = types.BoolNull()
    }
    data.MaxObjectsPerScan = numberFromAPI(data.MaxObjectsPerScan, item["maxObjectsPerScan"])
    data.MaxBytesPerObject = numberFromAPI(data.MaxBytesPerObject, item["maxBytesPerObject"])
    data.MaxBytesPerScan = numberFromAPI(data.MaxBytesPerScan, item["maxBytesPerScan"])
    data.MaxRowsPerTable = numberFromAPI(data.MaxRowsPerTable, item["maxRowsPerTable"])
    data.MaxTablesPerScan = numberFromAPI(data.MaxTablesPerScan, item["maxTablesPerScan"])
    data.TelemetrySampleSize = numberFromAPI(data.TelemetrySampleSize, item["telemetrySampleSize"])
    data.ContentScanIntervalInHours = numberFromAPI(data.ContentScanIntervalInHours, item["contentScanIntervalInHours"])
    if obj, ok := item["allowedValues"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AllowedValues = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AllowedValues = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AllowedValues = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AllowedValues = types.StringValue(string(jsonBytes))
        } else {
            data.AllowedValues = types.StringNull()
        }
    } else if val, ok := item["allowedValues"].(string); ok {
        data.AllowedValues = types.StringValue(val)
    } else {
        data.AllowedValues = types.StringNull()
    }
    if obj, ok := item["disabledDataTypeIds"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DisabledDataTypeIds = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DisabledDataTypeIds = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DisabledDataTypeIds = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DisabledDataTypeIds = types.StringValue(string(jsonBytes))
        } else {
            data.DisabledDataTypeIds = types.StringNull()
        }
    } else if val, ok := item["disabledDataTypeIds"].(string); ok {
        data.DisabledDataTypeIds = types.StringValue(val)
    } else {
        data.DisabledDataTypeIds = types.StringNull()
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
