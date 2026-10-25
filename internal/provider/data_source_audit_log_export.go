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
var _ datasource.DataSource = &AuditLogExportDataSource{}

func NewAuditLogExportDataSource() datasource.DataSource {
    return &AuditLogExportDataSource{}
}

// AuditLogExportDataSource defines the data source implementation.
type AuditLogExportDataSource struct {
    client *Client
}

// AuditLogExportDataSourceModel describes the data source data model.
type AuditLogExportDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    DestinationType types.String `tfsdk:"destination_type"`
    Config types.String `tfsdk:"config"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    StartsAt types.String `tfsdk:"starts_at"`
    CursorCreatedAt types.String `tfsdk:"cursor_created_at"`
    CursorEntryId types.String `tfsdk:"cursor_entry_id"`
    ExportedEntryCount types.Number `tfsdk:"exported_entry_count"`
    LastExportedAt types.String `tfsdk:"last_exported_at"`
    LastAttemptAt types.String `tfsdk:"last_attempt_at"`
    LastError types.String `tfsdk:"last_error"`
    FailingSince types.String `tfsdk:"failing_since"`
    ConsecutiveFailureCount types.Number `tfsdk:"consecutive_failure_count"`
    NextAttemptAt types.String `tfsdk:"next_attempt_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *AuditLogExportDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_audit_log_export"
}

func (d *AuditLogExportDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Where the audit log is exported: S3 or S3-compatible storage, Google Cloud Storage, Azure Blob Storage, a webhook or a syslog receiver. Look up an existing audit log export by `id`, or by any of its other arguments (`name`, `consecutive_failure_count`, `created_by_user_id`, ...): each one set must match, and exactly one audit log export may match them all.",

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
                MarkdownDescription: "ID of the project whose audit log is exported. Empty for an export of the whole instance. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "A name for this export, such as Compliance archive.",
                Optional: true,
                Computed: true,
            },
            "destination_type": schema.StringAttribute{
                MarkdownDescription: "Where entries are sent: s3, s3-compatible, google-cloud-storage, azure-blob-storage, webhook or syslog. Fixed once created.",
                Optional: true,
                Computed: true,
            },
            "config": schema.StringAttribute{
                MarkdownDescription: "The destination's settings that are not secret, such as the bucket and the region. Keys are defined by the destination. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether entries are exported. A paused export keeps its place.",
                Optional: true,
                Computed: true,
            },
            "starts_at": schema.StringAttribute{
                MarkdownDescription: "Entries recorded from this time on are exported. Empty on create means from now; an earlier date also sends the entries kept since then.",
                Computed: true,
            },
            "cursor_created_at": schema.StringAttribute{
                MarkdownDescription: "When the last entry delivered was recorded: where the next batch starts.",
                Computed: true,
            },
            "cursor_entry_id": schema.StringAttribute{
                MarkdownDescription: "The ID of the last entry delivered, which orders the entries recorded in the same second.",
                Optional: true,
                Computed: true,
            },
            "exported_entry_count": schema.NumberAttribute{
                MarkdownDescription: "How many entries have been delivered so far.",
                Optional: true,
                Computed: true,
            },
            "last_exported_at": schema.StringAttribute{
                MarkdownDescription: "When entries were last delivered.",
                Computed: true,
            },
            "last_attempt_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime last tried to deliver entries.",
                Computed: true,
            },
            "last_error": schema.StringAttribute{
                MarkdownDescription: "Why the last delivery failed, with credentials redacted. Cleared by the next delivery that succeeds.",
                Optional: true,
                Computed: true,
            },
            "failing_since": schema.StringAttribute{
                MarkdownDescription: "When deliveries started failing. Cleared by the next delivery that succeeds.",
                Computed: true,
            },
            "consecutive_failure_count": schema.NumberAttribute{
                MarkdownDescription: "How many deliveries in a row have failed. It sets how long the next retry waits.",
                Optional: true,
                Computed: true,
            },
            "next_attempt_at": schema.StringAttribute{
                MarkdownDescription: "While deliveries fail: when OneUptime tries again. Saving the export's settings tries again at once.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the user who created this export. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *AuditLogExportDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AuditLogExportDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data AuditLogExportDataSourceModel

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
    if !data.DestinationType.IsNull() && !data.DestinationType.IsUnknown() {
        filters["destinationType"] = data.DestinationType.ValueString()
        filterNames = append(filterNames, "destination_type = "+fmt.Sprintf("%q", data.DestinationType.ValueString()))
    }
    if !data.IsEnabled.IsNull() && !data.IsEnabled.IsUnknown() {
        filters["isEnabled"] = data.IsEnabled.ValueBool()
        filterNames = append(filterNames, "is_enabled = "+fmt.Sprintf("%t", data.IsEnabled.ValueBool()))
    }
    if !data.CursorEntryId.IsNull() && !data.CursorEntryId.IsUnknown() {
        filters["cursorEntryId"] = data.CursorEntryId.ValueString()
        filterNames = append(filterNames, "cursor_entry_id = "+fmt.Sprintf("%q", data.CursorEntryId.ValueString()))
    }
    if !data.ExportedEntryCount.IsNull() && !data.ExportedEntryCount.IsUnknown() {
        filters["exportedEntryCount"] = lookupNumber(data.ExportedEntryCount)
        filterNames = append(filterNames, "exported_entry_count = "+data.ExportedEntryCount.ValueBigFloat().String())
    }
    if !data.LastError.IsNull() && !data.LastError.IsUnknown() {
        filters["lastError"] = data.LastError.ValueString()
        filterNames = append(filterNames, "last_error = "+fmt.Sprintf("%q", data.LastError.ValueString()))
    }
    if !data.ConsecutiveFailureCount.IsNull() && !data.ConsecutiveFailureCount.IsUnknown() {
        filters["consecutiveFailureCount"] = lookupNumber(data.ConsecutiveFailureCount)
        filterNames = append(filterNames, "consecutive_failure_count = "+data.ConsecutiveFailureCount.ValueBigFloat().String())
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the audit log export up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the audit log export up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "name": true,
        "destinationType": true,
        "config": true,
        "isEnabled": true,
        "startsAt": true,
        "cursorCreatedAt": true,
        "cursorEntryId": true,
        "exportedEntryCount": true,
        "lastExportedAt": true,
        "lastAttemptAt": true,
        "lastError": true,
        "failingSince": true,
        "consecutiveFailureCount": true,
        "nextAttemptAt": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/audit-log-export/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read audit_log_export, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No audit log export found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read audit_log_export: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/audit-log-export/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list audit_log_export, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list audit_log_export: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No audit log export matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one audit log export matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for audit_log_export.")
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
    if obj, ok := item["destinationType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DestinationType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DestinationType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DestinationType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DestinationType = types.StringValue(string(jsonBytes))
        } else {
            data.DestinationType = types.StringNull()
        }
    } else if val, ok := item["destinationType"].(string); ok {
        data.DestinationType = types.StringValue(val)
    } else {
        data.DestinationType = types.StringNull()
    }
    if obj, ok := item["config"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Config = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Config = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Config = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Config = types.StringValue(string(jsonBytes))
        } else {
            data.Config = types.StringNull()
        }
    } else if val, ok := item["config"].(string); ok {
        data.Config = types.StringValue(val)
    } else {
        data.Config = types.StringNull()
    }
    if val, ok := item["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    } else {
        data.IsEnabled = types.BoolNull()
    }
    if obj, ok := item["startsAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartsAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StartsAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StartsAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StartsAt = types.StringValue(string(jsonBytes))
        } else {
            data.StartsAt = types.StringNull()
        }
    } else if val, ok := item["startsAt"].(string); ok {
        data.StartsAt = types.StringValue(val)
    } else {
        data.StartsAt = types.StringNull()
    }
    if obj, ok := item["cursorCreatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CursorCreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CursorCreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CursorCreatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CursorCreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CursorCreatedAt = types.StringNull()
        }
    } else if val, ok := item["cursorCreatedAt"].(string); ok {
        data.CursorCreatedAt = types.StringValue(val)
    } else {
        data.CursorCreatedAt = types.StringNull()
    }
    if obj, ok := item["cursorEntryId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CursorEntryId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CursorEntryId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CursorEntryId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CursorEntryId = types.StringValue(string(jsonBytes))
        } else {
            data.CursorEntryId = types.StringNull()
        }
    } else if val, ok := item["cursorEntryId"].(string); ok {
        data.CursorEntryId = types.StringValue(val)
    } else {
        data.CursorEntryId = types.StringNull()
    }
    data.ExportedEntryCount = numberFromAPI(data.ExportedEntryCount, item["exportedEntryCount"])
    if obj, ok := item["lastExportedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastExportedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastExportedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastExportedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastExportedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastExportedAt = types.StringNull()
        }
    } else if val, ok := item["lastExportedAt"].(string); ok {
        data.LastExportedAt = types.StringValue(val)
    } else {
        data.LastExportedAt = types.StringNull()
    }
    if obj, ok := item["lastAttemptAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastAttemptAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastAttemptAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastAttemptAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastAttemptAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastAttemptAt = types.StringNull()
        }
    } else if val, ok := item["lastAttemptAt"].(string); ok {
        data.LastAttemptAt = types.StringValue(val)
    } else {
        data.LastAttemptAt = types.StringNull()
    }
    if obj, ok := item["lastError"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastError = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastError = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastError = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastError = types.StringValue(string(jsonBytes))
        } else {
            data.LastError = types.StringNull()
        }
    } else if val, ok := item["lastError"].(string); ok {
        data.LastError = types.StringValue(val)
    } else {
        data.LastError = types.StringNull()
    }
    if obj, ok := item["failingSince"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FailingSince = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FailingSince = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FailingSince = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FailingSince = types.StringValue(string(jsonBytes))
        } else {
            data.FailingSince = types.StringNull()
        }
    } else if val, ok := item["failingSince"].(string); ok {
        data.FailingSince = types.StringValue(val)
    } else {
        data.FailingSince = types.StringNull()
    }
    data.ConsecutiveFailureCount = numberFromAPI(data.ConsecutiveFailureCount, item["consecutiveFailureCount"])
    if obj, ok := item["nextAttemptAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NextAttemptAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.NextAttemptAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.NextAttemptAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.NextAttemptAt = types.StringValue(string(jsonBytes))
        } else {
            data.NextAttemptAt = types.StringNull()
        }
    } else if val, ok := item["nextAttemptAt"].(string); ok {
        data.NextAttemptAt = types.StringValue(val)
    } else {
        data.NextAttemptAt = types.StringNull()
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
