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
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &AuditLogExportResource{}
var _ resource.ResourceWithImportState = &AuditLogExportResource{}

func NewAuditLogExportResource() resource.Resource {
    return &AuditLogExportResource{}
}

// AuditLogExportResource defines the resource implementation.
type AuditLogExportResource struct {
    client *Client
}

// AuditLogExportResourceModel describes the resource data model.
type AuditLogExportResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    Name JSONSubsetValue `tfsdk:"name"`
    DestinationType types.String `tfsdk:"destination_type"`
    Config JSONSubsetValue `tfsdk:"config"`
    Secrets types.String `tfsdk:"secrets"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    StartsAt RFC3339Value `tfsdk:"starts_at"`
    CreatedAt RFC3339Value `tfsdk:"created_at"`
    UpdatedAt RFC3339Value `tfsdk:"updated_at"`
    CursorCreatedAt RFC3339Value `tfsdk:"cursor_created_at"`
    CursorEntryId types.String `tfsdk:"cursor_entry_id"`
    ExportedEntryCount types.Number `tfsdk:"exported_entry_count"`
    LastExportedAt RFC3339Value `tfsdk:"last_exported_at"`
    LastAttemptAt RFC3339Value `tfsdk:"last_attempt_at"`
    LastError types.String `tfsdk:"last_error"`
    FailingSince RFC3339Value `tfsdk:"failing_since"`
    ConsecutiveFailureCount types.Number `tfsdk:"consecutive_failure_count"`
    NextAttemptAt RFC3339Value `tfsdk:"next_attempt_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (r *AuditLogExportResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_audit_log_export"
}

func (r *AuditLogExportResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = r.schemaDefinition()
}

func (r *AuditLogExportResource) schemaDefinition() schema.Schema {
    return schema.Schema{
        MarkdownDescription: "Where the audit log is exported: S3 or S3-compatible storage, Google Cloud Storage, Azure Blob Storage, a webhook or a syslog receiver.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of the project whose audit log is exported. Empty for an export of the whole instance. The ID of a `oneuptime_project`.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "A name for this export, such as Compliance archive.",
                CustomType: JSONSubsetType{},
                Required: true,
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "destination_type": schema.StringAttribute{
                MarkdownDescription: "Where entries are sent: s3, s3-compatible, google-cloud-storage, azure-blob-storage, webhook or syslog. Fixed once created.",
                Required: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "config": schema.StringAttribute{
                MarkdownDescription: "The destination's settings that are not secret, such as the bucket and the region. Keys are defined by the destination. A JSON value: write it with `jsonencode()`.",
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
            "secrets": schema.StringAttribute{
                MarkdownDescription: "The destination's secrets, such as an access key's secret, as a JSON object. Encrypted at rest and never returned by the API.",
                Optional: true,
                Sensitive: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether entries are exported. A paused export keeps its place.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "starts_at": schema.StringAttribute{
                MarkdownDescription: "Entries recorded from this time on are exported. Empty on create means from now; an earlier date also sends the entries kept since then.",
                CustomType: RFC3339Type{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
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
            "cursor_created_at": schema.StringAttribute{
                MarkdownDescription: "When the last entry delivered was recorded: where the next batch starts.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "cursor_entry_id": schema.StringAttribute{
                MarkdownDescription: "The ID of the last entry delivered, which orders the entries recorded in the same second.",
                Computed: true,
            },
            "exported_entry_count": schema.NumberAttribute{
                MarkdownDescription: "How many entries have been delivered so far.",
                Computed: true,
            },
            "last_exported_at": schema.StringAttribute{
                MarkdownDescription: "When entries were last delivered.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "last_attempt_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime last tried to deliver entries.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "last_error": schema.StringAttribute{
                MarkdownDescription: "Why the last delivery failed, with credentials redacted. Cleared by the next delivery that succeeds.",
                Computed: true,
            },
            "failing_since": schema.StringAttribute{
                MarkdownDescription: "When deliveries started failing. Cleared by the next delivery that succeeds.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "consecutive_failure_count": schema.NumberAttribute{
                MarkdownDescription: "How many deliveries in a row have failed. It sets how long the next retry waits.",
                Computed: true,
            },
            "next_attempt_at": schema.StringAttribute{
                MarkdownDescription: "While deliveries fail: when OneUptime tries again. Saving the export's settings tries again at once.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the user who created this export. The ID of a `oneuptime_user` (see the data source).",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
        },
    }
}

func (r *AuditLogExportResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *AuditLogExportResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data AuditLogExportResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // What the configuration sets, and what Terraform planned (see keepPlannedValues).
    var config AuditLogExportResourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
    if resp.Diagnostics.HasError() {
        return
    }
    plan := data



    // Create API request body. Unset (null/unknown) optional fields are
    // omitted so server-side defaults apply instead of being overwritten
    // with zero values.
    auditLogExportRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := auditLogExportRequest["data"].(map[string]interface{})

    if parsedName := r.parseJSONField(data.Name); parsedName != nil {
        requestDataMap["name"] = parsedName
    }
    if !data.DestinationType.IsNull() && !data.DestinationType.IsUnknown() {
        requestDataMap["destinationType"] = data.DestinationType.ValueString()
    }
    if parsedConfig := r.parseJSONField(data.Config); parsedConfig != nil {
        requestDataMap["config"] = parsedConfig
    }
    if !data.Secrets.IsNull() && !data.Secrets.IsUnknown() {
        requestDataMap["secrets"] = data.Secrets.ValueString()
    }
    if !data.IsEnabled.IsNull() && !data.IsEnabled.IsUnknown() {
        requestDataMap["isEnabled"] = data.IsEnabled.ValueBool()
    }
    if !data.StartsAt.IsNull() && !data.StartsAt.IsUnknown() {
        requestDataMap["startsAt"] = data.StartsAt.ValueString()
    }

    // Make API call
    httpResp, err := r.client.Post(ctx, "/audit-log-export", auditLogExportRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create audit_log_export, got error: %s", err))
        return
    }

    var auditLogExportResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &auditLogExportResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to create audit_log_export: %s", err))
        return
    }

    // Extract the new resource id from the create response.
    createdId := ""
    if wrapper, ok := auditLogExportResponse["data"].(map[string]interface{}); ok {
        if val, ok := wrapper["_id"].(string); ok {
            createdId = val
        }
    } else if val, ok := auditLogExportResponse["_id"].(string); ok {
        createdId = val
    }
    if createdId == "" {
        resp.Diagnostics.AddError("OneUptime API Error", "Create response for audit_log_export did not contain an id. This is a bug in the provider or the API; please report it.")
        return
    }
    data.Id = types.StringValue(createdId)

    /*
     * The server has committed the row. Persist what we know to state BEFORE
     * the read-back: if the read-back fails and we return without setting
     * state, Terraform never learns the resource exists and the created
     * audit_log_export is orphaned server-side — never refreshed, never
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
        "destinationType": true,
        "config": true,
        "isEnabled": true,
        "startsAt": true,
        "createdAt": true,
        "updatedAt": true,
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

    readResp, err := r.client.PostWithSelect(ctx, "/audit-log-export/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        /*
         * State already owns the id, so the resource is tracked and the next
         * refresh reconciles the remaining attributes. Warn rather than
         * error: erroring here would strand a real resource.
         */
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created audit_log_export but could not read it back; state is incomplete until the next refresh: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created audit_log_export but could not parse the read-back response; state is incomplete until the next refresh: %s", err))
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
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Name = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Name = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Name = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.Name = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Name = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.Name = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Name = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.Name = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = NewJSONSubsetValue(val)
    } else {
        data.Name = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["destinationType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DestinationType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DestinationType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DestinationType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DestinationType = types.StringValue(string(jsonBytes))
            } else {
                data.DestinationType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DestinationType = types.StringValue(string(jsonBytes))
            } else {
                data.DestinationType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DestinationType = types.StringValue(string(jsonBytes))
        } else {
            data.DestinationType = types.StringNull()
        }
    } else if val, ok := dataMap["destinationType"].(string); ok {
        data.DestinationType = types.StringValue(val)
    } else {
        data.DestinationType = types.StringNull()
    }
    if obj, ok := dataMap["config"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Config = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Config = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Config = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Config = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.Config = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Config = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.Config = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Config = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.Config = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["config"].(string); ok {
        data.Config = NewJSONSubsetValue(val)
    } else {
        data.Config = NewJSONSubsetNull()
    }
    if val, ok := dataMap["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if obj, ok := dataMap["startsAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.StartsAt = NewRFC3339Value(val)
        } else {
            data.StartsAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["startsAt"].(string); ok && val != "" {
        data.StartsAt = NewRFC3339Value(val)
    } else {
        data.StartsAt = NewRFC3339Null()
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
    if obj, ok := dataMap["cursorCreatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CursorCreatedAt = NewRFC3339Value(val)
        } else {
            data.CursorCreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["cursorCreatedAt"].(string); ok && val != "" {
        data.CursorCreatedAt = NewRFC3339Value(val)
    } else {
        data.CursorCreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["cursorEntryId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CursorEntryId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CursorEntryId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CursorEntryId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CursorEntryId = types.StringValue(string(jsonBytes))
            } else {
                data.CursorEntryId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CursorEntryId = types.StringValue(string(jsonBytes))
            } else {
                data.CursorEntryId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CursorEntryId = types.StringValue(string(jsonBytes))
        } else {
            data.CursorEntryId = types.StringNull()
        }
    } else if val, ok := dataMap["cursorEntryId"].(string); ok {
        data.CursorEntryId = types.StringValue(val)
    } else {
        data.CursorEntryId = types.StringNull()
    }
    data.ExportedEntryCount = numberFromAPI(data.ExportedEntryCount, dataMap["exportedEntryCount"])
    if obj, ok := dataMap["lastExportedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastExportedAt = NewRFC3339Value(val)
        } else {
            data.LastExportedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastExportedAt"].(string); ok && val != "" {
        data.LastExportedAt = NewRFC3339Value(val)
    } else {
        data.LastExportedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastAttemptAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastAttemptAt = NewRFC3339Value(val)
        } else {
            data.LastAttemptAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastAttemptAt"].(string); ok && val != "" {
        data.LastAttemptAt = NewRFC3339Value(val)
    } else {
        data.LastAttemptAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastError"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastError = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastError = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastError = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastError = types.StringValue(string(jsonBytes))
            } else {
                data.LastError = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastError = types.StringValue(string(jsonBytes))
            } else {
                data.LastError = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastError = types.StringValue(string(jsonBytes))
        } else {
            data.LastError = types.StringNull()
        }
    } else if val, ok := dataMap["lastError"].(string); ok {
        data.LastError = types.StringValue(val)
    } else {
        data.LastError = types.StringNull()
    }
    if obj, ok := dataMap["failingSince"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.FailingSince = NewRFC3339Value(val)
        } else {
            data.FailingSince = NewRFC3339Null()
        }
    } else if val, ok := dataMap["failingSince"].(string); ok && val != "" {
        data.FailingSince = NewRFC3339Value(val)
    } else {
        data.FailingSince = NewRFC3339Null()
    }
    data.ConsecutiveFailureCount = numberFromAPI(data.ConsecutiveFailureCount, dataMap["consecutiveFailureCount"])
    if obj, ok := dataMap["nextAttemptAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.NextAttemptAt = NewRFC3339Value(val)
        } else {
            data.NextAttemptAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["nextAttemptAt"].(string); ok && val != "" {
        data.NextAttemptAt = NewRFC3339Value(val)
    } else {
        data.NextAttemptAt = NewRFC3339Null()
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

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AuditLogExportResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data AuditLogExportResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "destinationType": true,
        "config": true,
        "isEnabled": true,
        "startsAt": true,
        "createdAt": true,
        "updatedAt": true,
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

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect(ctx, "/audit-log-export/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read audit_log_export, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var auditLogExportResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &auditLogExportResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse audit_log_export response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := auditLogExportResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = auditLogExportResponse
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
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Name = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Name = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Name = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.Name = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Name = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.Name = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Name = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.Name = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = NewJSONSubsetValue(val)
    } else {
        data.Name = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["destinationType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DestinationType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DestinationType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DestinationType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DestinationType = types.StringValue(string(jsonBytes))
            } else {
                data.DestinationType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DestinationType = types.StringValue(string(jsonBytes))
            } else {
                data.DestinationType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DestinationType = types.StringValue(string(jsonBytes))
        } else {
            data.DestinationType = types.StringNull()
        }
    } else if val, ok := dataMap["destinationType"].(string); ok {
        data.DestinationType = types.StringValue(val)
    } else {
        data.DestinationType = types.StringNull()
    }
    if obj, ok := dataMap["config"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Config = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Config = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Config = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Config = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.Config = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Config = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.Config = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Config = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.Config = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["config"].(string); ok {
        data.Config = NewJSONSubsetValue(val)
    } else {
        data.Config = NewJSONSubsetNull()
    }
    if val, ok := dataMap["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if obj, ok := dataMap["startsAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.StartsAt = NewRFC3339Value(val)
        } else {
            data.StartsAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["startsAt"].(string); ok && val != "" {
        data.StartsAt = NewRFC3339Value(val)
    } else {
        data.StartsAt = NewRFC3339Null()
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
    if obj, ok := dataMap["cursorCreatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CursorCreatedAt = NewRFC3339Value(val)
        } else {
            data.CursorCreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["cursorCreatedAt"].(string); ok && val != "" {
        data.CursorCreatedAt = NewRFC3339Value(val)
    } else {
        data.CursorCreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["cursorEntryId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CursorEntryId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CursorEntryId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CursorEntryId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CursorEntryId = types.StringValue(string(jsonBytes))
            } else {
                data.CursorEntryId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CursorEntryId = types.StringValue(string(jsonBytes))
            } else {
                data.CursorEntryId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CursorEntryId = types.StringValue(string(jsonBytes))
        } else {
            data.CursorEntryId = types.StringNull()
        }
    } else if val, ok := dataMap["cursorEntryId"].(string); ok {
        data.CursorEntryId = types.StringValue(val)
    } else {
        data.CursorEntryId = types.StringNull()
    }
    data.ExportedEntryCount = numberFromAPI(data.ExportedEntryCount, dataMap["exportedEntryCount"])
    if obj, ok := dataMap["lastExportedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastExportedAt = NewRFC3339Value(val)
        } else {
            data.LastExportedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastExportedAt"].(string); ok && val != "" {
        data.LastExportedAt = NewRFC3339Value(val)
    } else {
        data.LastExportedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastAttemptAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastAttemptAt = NewRFC3339Value(val)
        } else {
            data.LastAttemptAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastAttemptAt"].(string); ok && val != "" {
        data.LastAttemptAt = NewRFC3339Value(val)
    } else {
        data.LastAttemptAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastError"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastError = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastError = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastError = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastError = types.StringValue(string(jsonBytes))
            } else {
                data.LastError = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastError = types.StringValue(string(jsonBytes))
            } else {
                data.LastError = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastError = types.StringValue(string(jsonBytes))
        } else {
            data.LastError = types.StringNull()
        }
    } else if val, ok := dataMap["lastError"].(string); ok {
        data.LastError = types.StringValue(val)
    } else {
        data.LastError = types.StringNull()
    }
    if obj, ok := dataMap["failingSince"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.FailingSince = NewRFC3339Value(val)
        } else {
            data.FailingSince = NewRFC3339Null()
        }
    } else if val, ok := dataMap["failingSince"].(string); ok && val != "" {
        data.FailingSince = NewRFC3339Value(val)
    } else {
        data.FailingSince = NewRFC3339Null()
    }
    data.ConsecutiveFailureCount = numberFromAPI(data.ConsecutiveFailureCount, dataMap["consecutiveFailureCount"])
    if obj, ok := dataMap["nextAttemptAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.NextAttemptAt = NewRFC3339Value(val)
        } else {
            data.NextAttemptAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["nextAttemptAt"].(string); ok && val != "" {
        data.NextAttemptAt = NewRFC3339Value(val)
    } else {
        data.NextAttemptAt = NewRFC3339Null()
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
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AuditLogExportResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data AuditLogExportResourceModel
    var state AuditLogExportResourceModel

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

    // What the configuration sets, and what Terraform planned (see keepPlannedValues).
    var config AuditLogExportResourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
    if resp.Diagnostics.HasError() {
        return
    }
    plan := data

    // Create API request body
    auditLogExportRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := auditLogExportRequest["data"].(map[string]interface{})

    if !data.Name.IsUnknown() && !state.Name.IsUnknown() && !data.Name.Equal(state.Name) {
        var nameData interface{}
        if err := json.Unmarshal([]byte(data.Name.ValueString()), &nameData); err == nil {
            requestDataMap["name"] = nameData
        } else {
            requestDataMap["name"] = data.Name.ValueString()
        }
    }
    if !data.Config.IsUnknown() && !state.Config.IsUnknown() && !data.Config.Equal(state.Config) {
        var configData interface{}
        if err := json.Unmarshal([]byte(data.Config.ValueString()), &configData); err == nil {
            requestDataMap["config"] = configData
        } else {
            requestDataMap["config"] = data.Config.ValueString()
        }
    }
    if !data.Secrets.IsUnknown() && !state.Secrets.IsUnknown() && !data.Secrets.Equal(state.Secrets) {
        requestDataMap["secrets"] = data.Secrets.ValueString()
    }
    if !data.IsEnabled.IsUnknown() && !state.IsEnabled.IsUnknown() && !data.IsEnabled.Equal(state.IsEnabled) {
        requestDataMap["isEnabled"] = data.IsEnabled.ValueBool()
    }
    if !data.StartsAt.IsUnknown() && !state.StartsAt.IsUnknown() && !data.StartsAt.Equal(state.StartsAt) {
        requestDataMap["startsAt"] = data.StartsAt.ValueString()
    }

    // Only call the API when there are changed fields to send. An empty
    // update body is rejected by the API; state is still refreshed below so
    // this method never writes unverified plan values into state.
    if len(auditLogExportRequest["data"].(map[string]interface{})) > 0 {
        httpResp, err := r.client.Put(ctx, "/audit-log-export/" + data.Id.ValueString() + "", auditLogExportRequest)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update audit_log_export, got error: %s", err))
            return
        }

        // Parse the update response
        var auditLogExportResponse map[string]interface{}
        err = r.client.ParseResponse(httpResp, &auditLogExportResponse)
        if err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to update audit_log_export: %s", err))
            return
        }
        _ = auditLogExportResponse
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "destinationType": true,
        "config": true,
        "isEnabled": true,
        "startsAt": true,
        "createdAt": true,
        "updatedAt": true,
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

    readResp, err := r.client.PostWithSelect(ctx, "/audit-log-export/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read audit_log_export after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read audit_log_export after update: %s", err))
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
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Name = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Name = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Name = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.Name = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Name = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.Name = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Name = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.Name = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = NewJSONSubsetValue(val)
    } else {
        data.Name = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["destinationType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DestinationType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DestinationType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DestinationType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DestinationType = types.StringValue(string(jsonBytes))
            } else {
                data.DestinationType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DestinationType = types.StringValue(string(jsonBytes))
            } else {
                data.DestinationType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DestinationType = types.StringValue(string(jsonBytes))
        } else {
            data.DestinationType = types.StringNull()
        }
    } else if val, ok := dataMap["destinationType"].(string); ok {
        data.DestinationType = types.StringValue(val)
    } else {
        data.DestinationType = types.StringNull()
    }
    if obj, ok := dataMap["config"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Config = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Config = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Config = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Config = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.Config = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Config = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.Config = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Config = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.Config = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["config"].(string); ok {
        data.Config = NewJSONSubsetValue(val)
    } else {
        data.Config = NewJSONSubsetNull()
    }
    if val, ok := dataMap["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if obj, ok := dataMap["startsAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.StartsAt = NewRFC3339Value(val)
        } else {
            data.StartsAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["startsAt"].(string); ok && val != "" {
        data.StartsAt = NewRFC3339Value(val)
    } else {
        data.StartsAt = NewRFC3339Null()
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
    if obj, ok := dataMap["cursorCreatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CursorCreatedAt = NewRFC3339Value(val)
        } else {
            data.CursorCreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["cursorCreatedAt"].(string); ok && val != "" {
        data.CursorCreatedAt = NewRFC3339Value(val)
    } else {
        data.CursorCreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["cursorEntryId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CursorEntryId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CursorEntryId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CursorEntryId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CursorEntryId = types.StringValue(string(jsonBytes))
            } else {
                data.CursorEntryId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CursorEntryId = types.StringValue(string(jsonBytes))
            } else {
                data.CursorEntryId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CursorEntryId = types.StringValue(string(jsonBytes))
        } else {
            data.CursorEntryId = types.StringNull()
        }
    } else if val, ok := dataMap["cursorEntryId"].(string); ok {
        data.CursorEntryId = types.StringValue(val)
    } else {
        data.CursorEntryId = types.StringNull()
    }
    data.ExportedEntryCount = numberFromAPI(data.ExportedEntryCount, dataMap["exportedEntryCount"])
    if obj, ok := dataMap["lastExportedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastExportedAt = NewRFC3339Value(val)
        } else {
            data.LastExportedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastExportedAt"].(string); ok && val != "" {
        data.LastExportedAt = NewRFC3339Value(val)
    } else {
        data.LastExportedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastAttemptAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastAttemptAt = NewRFC3339Value(val)
        } else {
            data.LastAttemptAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastAttemptAt"].(string); ok && val != "" {
        data.LastAttemptAt = NewRFC3339Value(val)
    } else {
        data.LastAttemptAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastError"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastError = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastError = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastError = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastError = types.StringValue(string(jsonBytes))
            } else {
                data.LastError = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastError = types.StringValue(string(jsonBytes))
            } else {
                data.LastError = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastError = types.StringValue(string(jsonBytes))
        } else {
            data.LastError = types.StringNull()
        }
    } else if val, ok := dataMap["lastError"].(string); ok {
        data.LastError = types.StringValue(val)
    } else {
        data.LastError = types.StringNull()
    }
    if obj, ok := dataMap["failingSince"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.FailingSince = NewRFC3339Value(val)
        } else {
            data.FailingSince = NewRFC3339Null()
        }
    } else if val, ok := dataMap["failingSince"].(string); ok && val != "" {
        data.FailingSince = NewRFC3339Value(val)
    } else {
        data.FailingSince = NewRFC3339Null()
    }
    data.ConsecutiveFailureCount = numberFromAPI(data.ConsecutiveFailureCount, dataMap["consecutiveFailureCount"])
    if obj, ok := dataMap["nextAttemptAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.NextAttemptAt = NewRFC3339Value(val)
        } else {
            data.NextAttemptAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["nextAttemptAt"].(string); ok && val != "" {
        data.NextAttemptAt = NewRFC3339Value(val)
    } else {
        data.NextAttemptAt = NewRFC3339Null()
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
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    data.Id = state.Id

    // Unconfigured attributes keep their planned value; see keepPlannedValues.
    r.keepPlannedValues(&data, &plan, &config)

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AuditLogExportResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data AuditLogExportResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    httpResp, err := r.client.Delete(ctx, "/audit-log-export/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete audit_log_export, got error: %s", err))
        return
    }

    // A failed delete must keep the resource in state — silently dropping it
    // orphans real infrastructure. 404 means it is already gone.
    if httpResp.StatusCode >= 400 && httpResp.StatusCode != http.StatusNotFound {
        err = r.client.ParseResponse(httpResp, nil)
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to delete audit_log_export: %s", err))
        return
    }
    if httpResp.Body != nil {
        httpResp.Body.Close()
    }
}


func (r *AuditLogExportResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}


// keepPlannedValues puts back, after a create or an update, the planned value
// of each optional attribute the configuration leaves out. The server keeps
// some of these up to date on its own (when it last checked a heartbeat, the
// status a probe last reported...), so the value read back after the write can
// already differ from the plan, and Terraform would fail the apply with
// "Provider produced inconsistent result after apply". The next refresh reads
// the server's value, which is never a diff for an attribute nobody configured.
func (r *AuditLogExportResource) keepPlannedValues(data *AuditLogExportResourceModel, plan *AuditLogExportResourceModel, config *AuditLogExportResourceModel) {
    if config.Config.IsNull() && !plan.Config.IsUnknown() {
        data.Config = plan.Config
    }
    if config.IsEnabled.IsNull() && !plan.IsEnabled.IsUnknown() {
        data.IsEnabled = plan.IsEnabled
    }
    if config.StartsAt.IsNull() && !plan.StartsAt.IsUnknown() {
        data.StartsAt = plan.StartsAt
    }
}

// Helper method to convert Terraform map to Go interface{}
func (r *AuditLogExportResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *AuditLogExportResource) convertTerraformListToInterface(terraformList types.List) interface{} {
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
func (r *AuditLogExportResource) convertTerraformSetToInterface(terraformSet types.Set) interface{} {
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
func (r *AuditLogExportResource) parseJSONField(terraformString basetypes.StringValuable) interface{} {
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
func (r *AuditLogExportResource) normalizeURLWrappers(value interface{}) interface{} {
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

func (r *AuditLogExportResource) normalizeURLString(value string) string {
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
func (r *AuditLogExportResource) bigFloatToFloat64(bf *big.Float) interface{} {
    if bf == nil {
        return nil
    }
    f, _ := bf.Float64()
    return f
}

// Helper method to check if a type string is a valid OneUptime ObjectType.
// The registry itself lives in objecttypes.go, shared across the package.
func (r *AuditLogExportResource) isValidOneUptimeObjectType(typeStr string) bool {
    return validOneUptimeObjectTypes[typeStr]
}
