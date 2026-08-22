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
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/numberplanmodifier"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ScheduledMaintenanceMeasurementResource{}
var _ resource.ResourceWithImportState = &ScheduledMaintenanceMeasurementResource{}

func NewScheduledMaintenanceMeasurementResource() resource.Resource {
    return &ScheduledMaintenanceMeasurementResource{}
}

// ScheduledMaintenanceMeasurementResource defines the resource implementation.
type ScheduledMaintenanceMeasurementResource struct {
    client *Client
}

// ScheduledMaintenanceMeasurementResourceModel describes the resource data model.
type ScheduledMaintenanceMeasurementResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Key types.String `tfsdk:"key"`
    Description types.String `tfsdk:"description"`
    StartAnchorType types.String `tfsdk:"start_anchor_type"`
    EndAnchorType types.String `tfsdk:"end_anchor_type"`
    StartScheduledMaintenanceStateId types.String `tfsdk:"start_scheduled_maintenance_state_id"`
    EndScheduledMaintenanceStateId types.String `tfsdk:"end_scheduled_maintenance_state_id"`
    StartScheduledMaintenanceStateRole types.String `tfsdk:"start_scheduled_maintenance_state_role"`
    EndScheduledMaintenanceStateRole types.String `tfsdk:"end_scheduled_maintenance_state_role"`
    StartStateOccurrence types.String `tfsdk:"start_state_occurrence"`
    EndStateOccurrence types.String `tfsdk:"end_state_occurrence"`
    Unit types.String `tfsdk:"unit"`
    AggregationType types.String `tfsdk:"aggregation_type"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    ShowOnScheduledMaintenanceView types.Bool `tfsdk:"show_on_scheduled_maintenance_view"`
    Order types.Number `tfsdk:"order"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    CreatedAt RFC3339Value `tfsdk:"created_at"`
    UpdatedAt RFC3339Value `tfsdk:"updated_at"`
    DeletedAt RFC3339Value `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    MetricName types.String `tfsdk:"metric_name"`
    IsSystemDefined types.Bool `tfsdk:"is_system_defined"`
    BackfillRequestedAt RFC3339Value `tfsdk:"backfill_requested_at"`
    BackfillCursorCreatedAt RFC3339Value `tfsdk:"backfill_cursor_created_at"`
    BackfillCompletedAt RFC3339Value `tfsdk:"backfill_completed_at"`
}

func (r *ScheduledMaintenanceMeasurementResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_scheduled_maintenance_measurement"
}

func (r *ScheduledMaintenanceMeasurementResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "A named duration between two points in a scheduled maintenance event's life, computed automatically for every event",

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
                MarkdownDescription: "Name of this measurement, shown on charts and on the scheduled maintenance event page.",
                Required: true,
            },
            "key": schema.StringAttribute{
                MarkdownDescription: "Stable machine-readable key for this measurement. It is part of the metric name, so it cannot be changed once the measurement is created..",
                Required: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of what this measurement means to your team.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "start_anchor_type": schema.StringAttribute{
                MarkdownDescription: "Where the measurement starts - the moment the event was created, either end of the planned window, the start of its timeline, a specific state, or a state role..",
                Required: true,
            },
            "end_anchor_type": schema.StringAttribute{
                MarkdownDescription: "Where the measurement ends - the moment the event was created, either end of the planned window, the start of its timeline, a specific state, or a state role..",
                Required: true,
            },
            "start_scheduled_maintenance_state_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "end_scheduled_maintenance_state_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "start_scheduled_maintenance_state_role": schema.StringAttribute{
                MarkdownDescription: "The role of the state that starts this measurement (Scheduled, Ongoing, Ended or Resolved), when the start anchor is a state role. Resolving by role keeps working when a project renames or replaces the state..",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "end_scheduled_maintenance_state_role": schema.StringAttribute{
                MarkdownDescription: "The role of the state that ends this measurement (Scheduled, Ongoing, Ended or Resolved), when the end anchor is a state role. Resolving by role keeps working when a project renames or replaces the state..",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "start_state_occurrence": schema.StringAttribute{
                MarkdownDescription: "Which entry to use when the start state is entered more than once - the first time it was entered, or the last..",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("First"),
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "end_state_occurrence": schema.StringAttribute{
                MarkdownDescription: "Which entry to use when the end state is entered more than once - the first time it was entered, or the last..",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("First"),
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "unit": schema.StringAttribute{
                MarkdownDescription: "The unit this measurement is displayed in. Values are always stored in seconds..",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("seconds"),
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "aggregation_type": schema.StringAttribute{
                MarkdownDescription: "The aggregation this measurement's chart defaults to. Summing durations across events produces a number with no meaning, so Sum is not offered..",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("Avg"),
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this measurement is computed for scheduled maintenance events.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "show_on_scheduled_maintenance_view": schema.BoolAttribute{
                MarkdownDescription: "Whether this measurement is shown on the scheduled maintenance event page.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "order": schema.NumberAttribute{
                MarkdownDescription: "Order in which this measurement is displayed. Lowest first..",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(1)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
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
            "metric_name": schema.StringAttribute{
                MarkdownDescription: "Name of the metric this measurement writes to, derived from the key as oneuptime.scheduled-maintenance.measurement.<key>.",
                Computed: true,
            },
            "is_system_defined": schema.BoolAttribute{
                MarkdownDescription: "Whether this measurement was created by OneUptime rather than by your team.",
                Computed: true,
            },
            "backfill_requested_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "backfill_cursor_created_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "backfill_completed_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
        },
    }
}

func (r *ScheduledMaintenanceMeasurementResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *ScheduledMaintenanceMeasurementResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data ScheduledMaintenanceMeasurementResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }



    // Create API request body. Unset (null/unknown) optional fields are
    // omitted so server-side defaults apply instead of being overwritten
    // with zero values.
    scheduledMaintenanceMeasurementRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := scheduledMaintenanceMeasurementRequest["data"].(map[string]interface{})

    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        requestDataMap["name"] = data.Name.ValueString()
    }
    if !data.Key.IsNull() && !data.Key.IsUnknown() {
        requestDataMap["key"] = data.Key.ValueString()
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        requestDataMap["description"] = data.Description.ValueString()
    }
    if !data.StartAnchorType.IsNull() && !data.StartAnchorType.IsUnknown() {
        requestDataMap["startAnchorType"] = data.StartAnchorType.ValueString()
    }
    if !data.EndAnchorType.IsNull() && !data.EndAnchorType.IsUnknown() {
        requestDataMap["endAnchorType"] = data.EndAnchorType.ValueString()
    }
    if !data.StartScheduledMaintenanceStateId.IsNull() && !data.StartScheduledMaintenanceStateId.IsUnknown() {
        requestDataMap["startScheduledMaintenanceStateId"] = data.StartScheduledMaintenanceStateId.ValueString()
    }
    if !data.EndScheduledMaintenanceStateId.IsNull() && !data.EndScheduledMaintenanceStateId.IsUnknown() {
        requestDataMap["endScheduledMaintenanceStateId"] = data.EndScheduledMaintenanceStateId.ValueString()
    }
    if !data.StartScheduledMaintenanceStateRole.IsNull() && !data.StartScheduledMaintenanceStateRole.IsUnknown() {
        requestDataMap["startScheduledMaintenanceStateRole"] = data.StartScheduledMaintenanceStateRole.ValueString()
    }
    if !data.EndScheduledMaintenanceStateRole.IsNull() && !data.EndScheduledMaintenanceStateRole.IsUnknown() {
        requestDataMap["endScheduledMaintenanceStateRole"] = data.EndScheduledMaintenanceStateRole.ValueString()
    }
    if !data.StartStateOccurrence.IsNull() && !data.StartStateOccurrence.IsUnknown() {
        requestDataMap["startStateOccurrence"] = data.StartStateOccurrence.ValueString()
    }
    if !data.EndStateOccurrence.IsNull() && !data.EndStateOccurrence.IsUnknown() {
        requestDataMap["endStateOccurrence"] = data.EndStateOccurrence.ValueString()
    }
    if !data.Unit.IsNull() && !data.Unit.IsUnknown() {
        requestDataMap["unit"] = data.Unit.ValueString()
    }
    if !data.AggregationType.IsNull() && !data.AggregationType.IsUnknown() {
        requestDataMap["aggregationType"] = data.AggregationType.ValueString()
    }
    if !data.IsEnabled.IsNull() && !data.IsEnabled.IsUnknown() {
        requestDataMap["isEnabled"] = data.IsEnabled.ValueBool()
    }
    if !data.ShowOnScheduledMaintenanceView.IsNull() && !data.ShowOnScheduledMaintenanceView.IsUnknown() {
        requestDataMap["showOnScheduledMaintenanceView"] = data.ShowOnScheduledMaintenanceView.ValueBool()
    }
    if !data.Order.IsNull() && !data.Order.IsUnknown() {
        requestDataMap["order"] = r.bigFloatToFloat64(data.Order.ValueBigFloat())
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        requestDataMap["createdByUserId"] = data.CreatedByUserId.ValueString()
    }

    // Make API call
    httpResp, err := r.client.Post(ctx, "/scheduled-maintenance-measurement", scheduledMaintenanceMeasurementRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create scheduled_maintenance_measurement, got error: %s", err))
        return
    }

    var scheduledMaintenanceMeasurementResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &scheduledMaintenanceMeasurementResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to create scheduled_maintenance_measurement: %s", err))
        return
    }

    // Extract the new resource id from the create response.
    createdId := ""
    if wrapper, ok := scheduledMaintenanceMeasurementResponse["data"].(map[string]interface{}); ok {
        if val, ok := wrapper["_id"].(string); ok {
            createdId = val
        }
    } else if val, ok := scheduledMaintenanceMeasurementResponse["_id"].(string); ok {
        createdId = val
    }
    if createdId == "" {
        resp.Diagnostics.AddError("OneUptime API Error", "Create response for scheduled_maintenance_measurement did not contain an id. This is a bug in the provider or the API; please report it.")
        return
    }
    data.Id = types.StringValue(createdId)

    /*
     * The server has committed the row. Persist what we know to state BEFORE
     * the read-back: if the read-back fails and we return without setting
     * state, Terraform never learns the resource exists and the created
     * scheduled_maintenance_measurement is orphaned server-side — never refreshed, never
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
        "key": true,
        "description": true,
        "startAnchorType": true,
        "endAnchorType": true,
        "startScheduledMaintenanceStateId": true,
        "endScheduledMaintenanceStateId": true,
        "startScheduledMaintenanceStateRole": true,
        "endScheduledMaintenanceStateRole": true,
        "startStateOccurrence": true,
        "endStateOccurrence": true,
        "unit": true,
        "aggregationType": true,
        "isEnabled": true,
        "showOnScheduledMaintenanceView": true,
        "order": true,
        "createdByUserId": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "metricName": true,
        "isSystemDefined": true,
        "backfillRequestedAt": true,
        "backfillCursorCreatedAt": true,
        "backfillCompletedAt": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/scheduled-maintenance-measurement/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        /*
         * State already owns the id, so the resource is tracked and the next
         * refresh reconciles the remaining attributes. Warn rather than
         * error: erroring here would strand a real resource.
         */
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created scheduled_maintenance_measurement but could not read it back; state is incomplete until the next refresh: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created scheduled_maintenance_measurement but could not parse the read-back response; state is incomplete until the next refresh: %s", err))
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
    if obj, ok := dataMap["key"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Key = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Key = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Key = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Key = types.StringValue(string(jsonBytes))
            } else {
                data.Key = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Key = types.StringValue(string(jsonBytes))
            } else {
                data.Key = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Key = types.StringValue(string(jsonBytes))
        } else {
            data.Key = types.StringNull()
        }
    } else if val, ok := dataMap["key"].(string); ok {
        data.Key = types.StringValue(val)
    } else {
        data.Key = types.StringNull()
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
    if obj, ok := dataMap["startAnchorType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartAnchorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StartAnchorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StartAnchorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.StartAnchorType = types.StringValue(string(jsonBytes))
            } else {
                data.StartAnchorType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.StartAnchorType = types.StringValue(string(jsonBytes))
            } else {
                data.StartAnchorType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StartAnchorType = types.StringValue(string(jsonBytes))
        } else {
            data.StartAnchorType = types.StringNull()
        }
    } else if val, ok := dataMap["startAnchorType"].(string); ok {
        data.StartAnchorType = types.StringValue(val)
    } else {
        data.StartAnchorType = types.StringNull()
    }
    if obj, ok := dataMap["endAnchorType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndAnchorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EndAnchorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EndAnchorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EndAnchorType = types.StringValue(string(jsonBytes))
            } else {
                data.EndAnchorType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EndAnchorType = types.StringValue(string(jsonBytes))
            } else {
                data.EndAnchorType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EndAnchorType = types.StringValue(string(jsonBytes))
        } else {
            data.EndAnchorType = types.StringNull()
        }
    } else if val, ok := dataMap["endAnchorType"].(string); ok {
        data.EndAnchorType = types.StringValue(val)
    } else {
        data.EndAnchorType = types.StringNull()
    }
    if obj, ok := dataMap["startScheduledMaintenanceStateId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartScheduledMaintenanceStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StartScheduledMaintenanceStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StartScheduledMaintenanceStateId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.StartScheduledMaintenanceStateId = types.StringValue(string(jsonBytes))
            } else {
                data.StartScheduledMaintenanceStateId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.StartScheduledMaintenanceStateId = types.StringValue(string(jsonBytes))
            } else {
                data.StartScheduledMaintenanceStateId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StartScheduledMaintenanceStateId = types.StringValue(string(jsonBytes))
        } else {
            data.StartScheduledMaintenanceStateId = types.StringNull()
        }
    } else if val, ok := dataMap["startScheduledMaintenanceStateId"].(string); ok {
        data.StartScheduledMaintenanceStateId = types.StringValue(val)
    } else {
        data.StartScheduledMaintenanceStateId = types.StringNull()
    }
    if obj, ok := dataMap["endScheduledMaintenanceStateId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndScheduledMaintenanceStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EndScheduledMaintenanceStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EndScheduledMaintenanceStateId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EndScheduledMaintenanceStateId = types.StringValue(string(jsonBytes))
            } else {
                data.EndScheduledMaintenanceStateId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EndScheduledMaintenanceStateId = types.StringValue(string(jsonBytes))
            } else {
                data.EndScheduledMaintenanceStateId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EndScheduledMaintenanceStateId = types.StringValue(string(jsonBytes))
        } else {
            data.EndScheduledMaintenanceStateId = types.StringNull()
        }
    } else if val, ok := dataMap["endScheduledMaintenanceStateId"].(string); ok {
        data.EndScheduledMaintenanceStateId = types.StringValue(val)
    } else {
        data.EndScheduledMaintenanceStateId = types.StringNull()
    }
    if obj, ok := dataMap["startScheduledMaintenanceStateRole"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartScheduledMaintenanceStateRole = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StartScheduledMaintenanceStateRole = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StartScheduledMaintenanceStateRole = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.StartScheduledMaintenanceStateRole = types.StringValue(string(jsonBytes))
            } else {
                data.StartScheduledMaintenanceStateRole = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.StartScheduledMaintenanceStateRole = types.StringValue(string(jsonBytes))
            } else {
                data.StartScheduledMaintenanceStateRole = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StartScheduledMaintenanceStateRole = types.StringValue(string(jsonBytes))
        } else {
            data.StartScheduledMaintenanceStateRole = types.StringNull()
        }
    } else if val, ok := dataMap["startScheduledMaintenanceStateRole"].(string); ok {
        data.StartScheduledMaintenanceStateRole = types.StringValue(val)
    } else {
        data.StartScheduledMaintenanceStateRole = types.StringNull()
    }
    if obj, ok := dataMap["endScheduledMaintenanceStateRole"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndScheduledMaintenanceStateRole = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EndScheduledMaintenanceStateRole = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EndScheduledMaintenanceStateRole = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EndScheduledMaintenanceStateRole = types.StringValue(string(jsonBytes))
            } else {
                data.EndScheduledMaintenanceStateRole = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EndScheduledMaintenanceStateRole = types.StringValue(string(jsonBytes))
            } else {
                data.EndScheduledMaintenanceStateRole = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EndScheduledMaintenanceStateRole = types.StringValue(string(jsonBytes))
        } else {
            data.EndScheduledMaintenanceStateRole = types.StringNull()
        }
    } else if val, ok := dataMap["endScheduledMaintenanceStateRole"].(string); ok {
        data.EndScheduledMaintenanceStateRole = types.StringValue(val)
    } else {
        data.EndScheduledMaintenanceStateRole = types.StringNull()
    }
    if obj, ok := dataMap["startStateOccurrence"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartStateOccurrence = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StartStateOccurrence = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StartStateOccurrence = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.StartStateOccurrence = types.StringValue(string(jsonBytes))
            } else {
                data.StartStateOccurrence = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.StartStateOccurrence = types.StringValue(string(jsonBytes))
            } else {
                data.StartStateOccurrence = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StartStateOccurrence = types.StringValue(string(jsonBytes))
        } else {
            data.StartStateOccurrence = types.StringNull()
        }
    } else if val, ok := dataMap["startStateOccurrence"].(string); ok {
        data.StartStateOccurrence = types.StringValue(val)
    } else {
        data.StartStateOccurrence = types.StringNull()
    }
    if obj, ok := dataMap["endStateOccurrence"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndStateOccurrence = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EndStateOccurrence = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EndStateOccurrence = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EndStateOccurrence = types.StringValue(string(jsonBytes))
            } else {
                data.EndStateOccurrence = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EndStateOccurrence = types.StringValue(string(jsonBytes))
            } else {
                data.EndStateOccurrence = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EndStateOccurrence = types.StringValue(string(jsonBytes))
        } else {
            data.EndStateOccurrence = types.StringNull()
        }
    } else if val, ok := dataMap["endStateOccurrence"].(string); ok {
        data.EndStateOccurrence = types.StringValue(val)
    } else {
        data.EndStateOccurrence = types.StringNull()
    }
    if obj, ok := dataMap["unit"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Unit = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Unit = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Unit = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Unit = types.StringValue(string(jsonBytes))
            } else {
                data.Unit = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Unit = types.StringValue(string(jsonBytes))
            } else {
                data.Unit = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Unit = types.StringValue(string(jsonBytes))
        } else {
            data.Unit = types.StringNull()
        }
    } else if val, ok := dataMap["unit"].(string); ok {
        data.Unit = types.StringValue(val)
    } else {
        data.Unit = types.StringNull()
    }
    if obj, ok := dataMap["aggregationType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AggregationType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AggregationType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AggregationType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AggregationType = types.StringValue(string(jsonBytes))
            } else {
                data.AggregationType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AggregationType = types.StringValue(string(jsonBytes))
            } else {
                data.AggregationType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AggregationType = types.StringValue(string(jsonBytes))
        } else {
            data.AggregationType = types.StringNull()
        }
    } else if val, ok := dataMap["aggregationType"].(string); ok {
        data.AggregationType = types.StringValue(val)
    } else {
        data.AggregationType = types.StringNull()
    }
    if val, ok := dataMap["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := dataMap["showOnScheduledMaintenanceView"].(bool); ok {
        data.ShowOnScheduledMaintenanceView = types.BoolValue(val)
    }
    if val, ok := dataMap["order"].(float64); ok {
        data.Order = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["order"].(int); ok {
        data.Order = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["order"].(int64); ok {
        data.Order = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["order"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.Order = types.NumberValue(big.NewFloat(val))
        } else {
            data.Order = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.Order = types.NumberNull()
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
    if obj, ok := dataMap["metricName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MetricName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MetricName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MetricName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MetricName = types.StringValue(string(jsonBytes))
            } else {
                data.MetricName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MetricName = types.StringValue(string(jsonBytes))
            } else {
                data.MetricName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MetricName = types.StringValue(string(jsonBytes))
        } else {
            data.MetricName = types.StringNull()
        }
    } else if val, ok := dataMap["metricName"].(string); ok {
        data.MetricName = types.StringValue(val)
    } else {
        data.MetricName = types.StringNull()
    }
    if val, ok := dataMap["isSystemDefined"].(bool); ok {
        data.IsSystemDefined = types.BoolValue(val)
    }
    if obj, ok := dataMap["backfillRequestedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.BackfillRequestedAt = NewRFC3339Value(val)
        } else {
            data.BackfillRequestedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["backfillRequestedAt"].(string); ok && val != "" {
        data.BackfillRequestedAt = NewRFC3339Value(val)
    } else {
        data.BackfillRequestedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["backfillCursorCreatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.BackfillCursorCreatedAt = NewRFC3339Value(val)
        } else {
            data.BackfillCursorCreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["backfillCursorCreatedAt"].(string); ok && val != "" {
        data.BackfillCursorCreatedAt = NewRFC3339Value(val)
    } else {
        data.BackfillCursorCreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["backfillCompletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.BackfillCompletedAt = NewRFC3339Value(val)
        } else {
            data.BackfillCompletedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["backfillCompletedAt"].(string); ok && val != "" {
        data.BackfillCompletedAt = NewRFC3339Value(val)
    } else {
        data.BackfillCompletedAt = NewRFC3339Null()
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

func (r *ScheduledMaintenanceMeasurementResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data ScheduledMaintenanceMeasurementResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "key": true,
        "description": true,
        "startAnchorType": true,
        "endAnchorType": true,
        "startScheduledMaintenanceStateId": true,
        "endScheduledMaintenanceStateId": true,
        "startScheduledMaintenanceStateRole": true,
        "endScheduledMaintenanceStateRole": true,
        "startStateOccurrence": true,
        "endStateOccurrence": true,
        "unit": true,
        "aggregationType": true,
        "isEnabled": true,
        "showOnScheduledMaintenanceView": true,
        "order": true,
        "createdByUserId": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "metricName": true,
        "isSystemDefined": true,
        "backfillRequestedAt": true,
        "backfillCursorCreatedAt": true,
        "backfillCompletedAt": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect(ctx, "/scheduled-maintenance-measurement/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read scheduled_maintenance_measurement, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var scheduledMaintenanceMeasurementResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &scheduledMaintenanceMeasurementResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse scheduled_maintenance_measurement response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := scheduledMaintenanceMeasurementResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = scheduledMaintenanceMeasurementResponse
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
    if obj, ok := dataMap["key"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Key = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Key = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Key = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Key = types.StringValue(string(jsonBytes))
            } else {
                data.Key = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Key = types.StringValue(string(jsonBytes))
            } else {
                data.Key = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Key = types.StringValue(string(jsonBytes))
        } else {
            data.Key = types.StringNull()
        }
    } else if val, ok := dataMap["key"].(string); ok {
        data.Key = types.StringValue(val)
    } else {
        data.Key = types.StringNull()
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
    if obj, ok := dataMap["startAnchorType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartAnchorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StartAnchorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StartAnchorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.StartAnchorType = types.StringValue(string(jsonBytes))
            } else {
                data.StartAnchorType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.StartAnchorType = types.StringValue(string(jsonBytes))
            } else {
                data.StartAnchorType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StartAnchorType = types.StringValue(string(jsonBytes))
        } else {
            data.StartAnchorType = types.StringNull()
        }
    } else if val, ok := dataMap["startAnchorType"].(string); ok {
        data.StartAnchorType = types.StringValue(val)
    } else {
        data.StartAnchorType = types.StringNull()
    }
    if obj, ok := dataMap["endAnchorType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndAnchorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EndAnchorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EndAnchorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EndAnchorType = types.StringValue(string(jsonBytes))
            } else {
                data.EndAnchorType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EndAnchorType = types.StringValue(string(jsonBytes))
            } else {
                data.EndAnchorType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EndAnchorType = types.StringValue(string(jsonBytes))
        } else {
            data.EndAnchorType = types.StringNull()
        }
    } else if val, ok := dataMap["endAnchorType"].(string); ok {
        data.EndAnchorType = types.StringValue(val)
    } else {
        data.EndAnchorType = types.StringNull()
    }
    if obj, ok := dataMap["startScheduledMaintenanceStateId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartScheduledMaintenanceStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StartScheduledMaintenanceStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StartScheduledMaintenanceStateId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.StartScheduledMaintenanceStateId = types.StringValue(string(jsonBytes))
            } else {
                data.StartScheduledMaintenanceStateId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.StartScheduledMaintenanceStateId = types.StringValue(string(jsonBytes))
            } else {
                data.StartScheduledMaintenanceStateId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StartScheduledMaintenanceStateId = types.StringValue(string(jsonBytes))
        } else {
            data.StartScheduledMaintenanceStateId = types.StringNull()
        }
    } else if val, ok := dataMap["startScheduledMaintenanceStateId"].(string); ok {
        data.StartScheduledMaintenanceStateId = types.StringValue(val)
    } else {
        data.StartScheduledMaintenanceStateId = types.StringNull()
    }
    if obj, ok := dataMap["endScheduledMaintenanceStateId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndScheduledMaintenanceStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EndScheduledMaintenanceStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EndScheduledMaintenanceStateId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EndScheduledMaintenanceStateId = types.StringValue(string(jsonBytes))
            } else {
                data.EndScheduledMaintenanceStateId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EndScheduledMaintenanceStateId = types.StringValue(string(jsonBytes))
            } else {
                data.EndScheduledMaintenanceStateId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EndScheduledMaintenanceStateId = types.StringValue(string(jsonBytes))
        } else {
            data.EndScheduledMaintenanceStateId = types.StringNull()
        }
    } else if val, ok := dataMap["endScheduledMaintenanceStateId"].(string); ok {
        data.EndScheduledMaintenanceStateId = types.StringValue(val)
    } else {
        data.EndScheduledMaintenanceStateId = types.StringNull()
    }
    if obj, ok := dataMap["startScheduledMaintenanceStateRole"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartScheduledMaintenanceStateRole = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StartScheduledMaintenanceStateRole = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StartScheduledMaintenanceStateRole = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.StartScheduledMaintenanceStateRole = types.StringValue(string(jsonBytes))
            } else {
                data.StartScheduledMaintenanceStateRole = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.StartScheduledMaintenanceStateRole = types.StringValue(string(jsonBytes))
            } else {
                data.StartScheduledMaintenanceStateRole = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StartScheduledMaintenanceStateRole = types.StringValue(string(jsonBytes))
        } else {
            data.StartScheduledMaintenanceStateRole = types.StringNull()
        }
    } else if val, ok := dataMap["startScheduledMaintenanceStateRole"].(string); ok {
        data.StartScheduledMaintenanceStateRole = types.StringValue(val)
    } else {
        data.StartScheduledMaintenanceStateRole = types.StringNull()
    }
    if obj, ok := dataMap["endScheduledMaintenanceStateRole"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndScheduledMaintenanceStateRole = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EndScheduledMaintenanceStateRole = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EndScheduledMaintenanceStateRole = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EndScheduledMaintenanceStateRole = types.StringValue(string(jsonBytes))
            } else {
                data.EndScheduledMaintenanceStateRole = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EndScheduledMaintenanceStateRole = types.StringValue(string(jsonBytes))
            } else {
                data.EndScheduledMaintenanceStateRole = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EndScheduledMaintenanceStateRole = types.StringValue(string(jsonBytes))
        } else {
            data.EndScheduledMaintenanceStateRole = types.StringNull()
        }
    } else if val, ok := dataMap["endScheduledMaintenanceStateRole"].(string); ok {
        data.EndScheduledMaintenanceStateRole = types.StringValue(val)
    } else {
        data.EndScheduledMaintenanceStateRole = types.StringNull()
    }
    if obj, ok := dataMap["startStateOccurrence"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartStateOccurrence = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StartStateOccurrence = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StartStateOccurrence = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.StartStateOccurrence = types.StringValue(string(jsonBytes))
            } else {
                data.StartStateOccurrence = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.StartStateOccurrence = types.StringValue(string(jsonBytes))
            } else {
                data.StartStateOccurrence = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StartStateOccurrence = types.StringValue(string(jsonBytes))
        } else {
            data.StartStateOccurrence = types.StringNull()
        }
    } else if val, ok := dataMap["startStateOccurrence"].(string); ok {
        data.StartStateOccurrence = types.StringValue(val)
    } else {
        data.StartStateOccurrence = types.StringNull()
    }
    if obj, ok := dataMap["endStateOccurrence"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndStateOccurrence = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EndStateOccurrence = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EndStateOccurrence = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EndStateOccurrence = types.StringValue(string(jsonBytes))
            } else {
                data.EndStateOccurrence = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EndStateOccurrence = types.StringValue(string(jsonBytes))
            } else {
                data.EndStateOccurrence = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EndStateOccurrence = types.StringValue(string(jsonBytes))
        } else {
            data.EndStateOccurrence = types.StringNull()
        }
    } else if val, ok := dataMap["endStateOccurrence"].(string); ok {
        data.EndStateOccurrence = types.StringValue(val)
    } else {
        data.EndStateOccurrence = types.StringNull()
    }
    if obj, ok := dataMap["unit"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Unit = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Unit = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Unit = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Unit = types.StringValue(string(jsonBytes))
            } else {
                data.Unit = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Unit = types.StringValue(string(jsonBytes))
            } else {
                data.Unit = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Unit = types.StringValue(string(jsonBytes))
        } else {
            data.Unit = types.StringNull()
        }
    } else if val, ok := dataMap["unit"].(string); ok {
        data.Unit = types.StringValue(val)
    } else {
        data.Unit = types.StringNull()
    }
    if obj, ok := dataMap["aggregationType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AggregationType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AggregationType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AggregationType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AggregationType = types.StringValue(string(jsonBytes))
            } else {
                data.AggregationType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AggregationType = types.StringValue(string(jsonBytes))
            } else {
                data.AggregationType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AggregationType = types.StringValue(string(jsonBytes))
        } else {
            data.AggregationType = types.StringNull()
        }
    } else if val, ok := dataMap["aggregationType"].(string); ok {
        data.AggregationType = types.StringValue(val)
    } else {
        data.AggregationType = types.StringNull()
    }
    if val, ok := dataMap["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := dataMap["showOnScheduledMaintenanceView"].(bool); ok {
        data.ShowOnScheduledMaintenanceView = types.BoolValue(val)
    }
    if val, ok := dataMap["order"].(float64); ok {
        data.Order = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["order"].(int); ok {
        data.Order = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["order"].(int64); ok {
        data.Order = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["order"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.Order = types.NumberValue(big.NewFloat(val))
        } else {
            data.Order = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.Order = types.NumberNull()
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
    if obj, ok := dataMap["metricName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MetricName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MetricName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MetricName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MetricName = types.StringValue(string(jsonBytes))
            } else {
                data.MetricName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MetricName = types.StringValue(string(jsonBytes))
            } else {
                data.MetricName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MetricName = types.StringValue(string(jsonBytes))
        } else {
            data.MetricName = types.StringNull()
        }
    } else if val, ok := dataMap["metricName"].(string); ok {
        data.MetricName = types.StringValue(val)
    } else {
        data.MetricName = types.StringNull()
    }
    if val, ok := dataMap["isSystemDefined"].(bool); ok {
        data.IsSystemDefined = types.BoolValue(val)
    }
    if obj, ok := dataMap["backfillRequestedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.BackfillRequestedAt = NewRFC3339Value(val)
        } else {
            data.BackfillRequestedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["backfillRequestedAt"].(string); ok && val != "" {
        data.BackfillRequestedAt = NewRFC3339Value(val)
    } else {
        data.BackfillRequestedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["backfillCursorCreatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.BackfillCursorCreatedAt = NewRFC3339Value(val)
        } else {
            data.BackfillCursorCreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["backfillCursorCreatedAt"].(string); ok && val != "" {
        data.BackfillCursorCreatedAt = NewRFC3339Value(val)
    } else {
        data.BackfillCursorCreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["backfillCompletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.BackfillCompletedAt = NewRFC3339Value(val)
        } else {
            data.BackfillCompletedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["backfillCompletedAt"].(string); ok && val != "" {
        data.BackfillCompletedAt = NewRFC3339Value(val)
    } else {
        data.BackfillCompletedAt = NewRFC3339Null()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ScheduledMaintenanceMeasurementResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data ScheduledMaintenanceMeasurementResourceModel
    var state ScheduledMaintenanceMeasurementResourceModel

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
    scheduledMaintenanceMeasurementRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := scheduledMaintenanceMeasurementRequest["data"].(map[string]interface{})

    if !data.Name.IsUnknown() && !state.Name.IsUnknown() && !data.Name.Equal(state.Name) {
        requestDataMap["name"] = data.Name.ValueString()
    }
    if !data.Description.IsUnknown() && !state.Description.IsUnknown() && !data.Description.Equal(state.Description) {
        requestDataMap["description"] = data.Description.ValueString()
    }
    if !data.StartAnchorType.IsUnknown() && !state.StartAnchorType.IsUnknown() && !data.StartAnchorType.Equal(state.StartAnchorType) {
        requestDataMap["startAnchorType"] = data.StartAnchorType.ValueString()
    }
    if !data.EndAnchorType.IsUnknown() && !state.EndAnchorType.IsUnknown() && !data.EndAnchorType.Equal(state.EndAnchorType) {
        requestDataMap["endAnchorType"] = data.EndAnchorType.ValueString()
    }
    if !data.StartScheduledMaintenanceStateId.IsUnknown() && !state.StartScheduledMaintenanceStateId.IsUnknown() && !data.StartScheduledMaintenanceStateId.Equal(state.StartScheduledMaintenanceStateId) {
        requestDataMap["startScheduledMaintenanceStateId"] = data.StartScheduledMaintenanceStateId.ValueString()
    }
    if !data.EndScheduledMaintenanceStateId.IsUnknown() && !state.EndScheduledMaintenanceStateId.IsUnknown() && !data.EndScheduledMaintenanceStateId.Equal(state.EndScheduledMaintenanceStateId) {
        requestDataMap["endScheduledMaintenanceStateId"] = data.EndScheduledMaintenanceStateId.ValueString()
    }
    if !data.StartScheduledMaintenanceStateRole.IsUnknown() && !state.StartScheduledMaintenanceStateRole.IsUnknown() && !data.StartScheduledMaintenanceStateRole.Equal(state.StartScheduledMaintenanceStateRole) {
        requestDataMap["startScheduledMaintenanceStateRole"] = data.StartScheduledMaintenanceStateRole.ValueString()
    }
    if !data.EndScheduledMaintenanceStateRole.IsUnknown() && !state.EndScheduledMaintenanceStateRole.IsUnknown() && !data.EndScheduledMaintenanceStateRole.Equal(state.EndScheduledMaintenanceStateRole) {
        requestDataMap["endScheduledMaintenanceStateRole"] = data.EndScheduledMaintenanceStateRole.ValueString()
    }
    if !data.StartStateOccurrence.IsUnknown() && !state.StartStateOccurrence.IsUnknown() && !data.StartStateOccurrence.Equal(state.StartStateOccurrence) {
        requestDataMap["startStateOccurrence"] = data.StartStateOccurrence.ValueString()
    }
    if !data.EndStateOccurrence.IsUnknown() && !state.EndStateOccurrence.IsUnknown() && !data.EndStateOccurrence.Equal(state.EndStateOccurrence) {
        requestDataMap["endStateOccurrence"] = data.EndStateOccurrence.ValueString()
    }
    if !data.Unit.IsUnknown() && !state.Unit.IsUnknown() && !data.Unit.Equal(state.Unit) {
        requestDataMap["unit"] = data.Unit.ValueString()
    }
    if !data.AggregationType.IsUnknown() && !state.AggregationType.IsUnknown() && !data.AggregationType.Equal(state.AggregationType) {
        requestDataMap["aggregationType"] = data.AggregationType.ValueString()
    }
    if !data.IsEnabled.IsUnknown() && !state.IsEnabled.IsUnknown() && !data.IsEnabled.Equal(state.IsEnabled) {
        requestDataMap["isEnabled"] = data.IsEnabled.ValueBool()
    }
    if !data.ShowOnScheduledMaintenanceView.IsUnknown() && !state.ShowOnScheduledMaintenanceView.IsUnknown() && !data.ShowOnScheduledMaintenanceView.Equal(state.ShowOnScheduledMaintenanceView) {
        requestDataMap["showOnScheduledMaintenanceView"] = data.ShowOnScheduledMaintenanceView.ValueBool()
    }
    if !data.Order.IsUnknown() && !state.Order.IsUnknown() && !data.Order.Equal(state.Order) {
        requestDataMap["order"] = r.bigFloatToFloat64(data.Order.ValueBigFloat())
    }

    // Only call the API when there are changed fields to send. An empty
    // update body is rejected by the API; state is still refreshed below so
    // this method never writes unverified plan values into state.
    if len(scheduledMaintenanceMeasurementRequest["data"].(map[string]interface{})) > 0 {
        httpResp, err := r.client.Put(ctx, "/scheduled-maintenance-measurement/" + data.Id.ValueString() + "", scheduledMaintenanceMeasurementRequest)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update scheduled_maintenance_measurement, got error: %s", err))
            return
        }

        // Parse the update response
        var scheduledMaintenanceMeasurementResponse map[string]interface{}
        err = r.client.ParseResponse(httpResp, &scheduledMaintenanceMeasurementResponse)
        if err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to update scheduled_maintenance_measurement: %s", err))
            return
        }
        _ = scheduledMaintenanceMeasurementResponse
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "key": true,
        "description": true,
        "startAnchorType": true,
        "endAnchorType": true,
        "startScheduledMaintenanceStateId": true,
        "endScheduledMaintenanceStateId": true,
        "startScheduledMaintenanceStateRole": true,
        "endScheduledMaintenanceStateRole": true,
        "startStateOccurrence": true,
        "endStateOccurrence": true,
        "unit": true,
        "aggregationType": true,
        "isEnabled": true,
        "showOnScheduledMaintenanceView": true,
        "order": true,
        "createdByUserId": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "metricName": true,
        "isSystemDefined": true,
        "backfillRequestedAt": true,
        "backfillCursorCreatedAt": true,
        "backfillCompletedAt": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/scheduled-maintenance-measurement/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read scheduled_maintenance_measurement after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read scheduled_maintenance_measurement after update: %s", err))
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
    if obj, ok := dataMap["key"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Key = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Key = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Key = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Key = types.StringValue(string(jsonBytes))
            } else {
                data.Key = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Key = types.StringValue(string(jsonBytes))
            } else {
                data.Key = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Key = types.StringValue(string(jsonBytes))
        } else {
            data.Key = types.StringNull()
        }
    } else if val, ok := dataMap["key"].(string); ok {
        data.Key = types.StringValue(val)
    } else {
        data.Key = types.StringNull()
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
    if obj, ok := dataMap["startAnchorType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartAnchorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StartAnchorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StartAnchorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.StartAnchorType = types.StringValue(string(jsonBytes))
            } else {
                data.StartAnchorType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.StartAnchorType = types.StringValue(string(jsonBytes))
            } else {
                data.StartAnchorType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StartAnchorType = types.StringValue(string(jsonBytes))
        } else {
            data.StartAnchorType = types.StringNull()
        }
    } else if val, ok := dataMap["startAnchorType"].(string); ok {
        data.StartAnchorType = types.StringValue(val)
    } else {
        data.StartAnchorType = types.StringNull()
    }
    if obj, ok := dataMap["endAnchorType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndAnchorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EndAnchorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EndAnchorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EndAnchorType = types.StringValue(string(jsonBytes))
            } else {
                data.EndAnchorType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EndAnchorType = types.StringValue(string(jsonBytes))
            } else {
                data.EndAnchorType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EndAnchorType = types.StringValue(string(jsonBytes))
        } else {
            data.EndAnchorType = types.StringNull()
        }
    } else if val, ok := dataMap["endAnchorType"].(string); ok {
        data.EndAnchorType = types.StringValue(val)
    } else {
        data.EndAnchorType = types.StringNull()
    }
    if obj, ok := dataMap["startScheduledMaintenanceStateId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartScheduledMaintenanceStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StartScheduledMaintenanceStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StartScheduledMaintenanceStateId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.StartScheduledMaintenanceStateId = types.StringValue(string(jsonBytes))
            } else {
                data.StartScheduledMaintenanceStateId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.StartScheduledMaintenanceStateId = types.StringValue(string(jsonBytes))
            } else {
                data.StartScheduledMaintenanceStateId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StartScheduledMaintenanceStateId = types.StringValue(string(jsonBytes))
        } else {
            data.StartScheduledMaintenanceStateId = types.StringNull()
        }
    } else if val, ok := dataMap["startScheduledMaintenanceStateId"].(string); ok {
        data.StartScheduledMaintenanceStateId = types.StringValue(val)
    } else {
        data.StartScheduledMaintenanceStateId = types.StringNull()
    }
    if obj, ok := dataMap["endScheduledMaintenanceStateId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndScheduledMaintenanceStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EndScheduledMaintenanceStateId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EndScheduledMaintenanceStateId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EndScheduledMaintenanceStateId = types.StringValue(string(jsonBytes))
            } else {
                data.EndScheduledMaintenanceStateId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EndScheduledMaintenanceStateId = types.StringValue(string(jsonBytes))
            } else {
                data.EndScheduledMaintenanceStateId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EndScheduledMaintenanceStateId = types.StringValue(string(jsonBytes))
        } else {
            data.EndScheduledMaintenanceStateId = types.StringNull()
        }
    } else if val, ok := dataMap["endScheduledMaintenanceStateId"].(string); ok {
        data.EndScheduledMaintenanceStateId = types.StringValue(val)
    } else {
        data.EndScheduledMaintenanceStateId = types.StringNull()
    }
    if obj, ok := dataMap["startScheduledMaintenanceStateRole"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartScheduledMaintenanceStateRole = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StartScheduledMaintenanceStateRole = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StartScheduledMaintenanceStateRole = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.StartScheduledMaintenanceStateRole = types.StringValue(string(jsonBytes))
            } else {
                data.StartScheduledMaintenanceStateRole = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.StartScheduledMaintenanceStateRole = types.StringValue(string(jsonBytes))
            } else {
                data.StartScheduledMaintenanceStateRole = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StartScheduledMaintenanceStateRole = types.StringValue(string(jsonBytes))
        } else {
            data.StartScheduledMaintenanceStateRole = types.StringNull()
        }
    } else if val, ok := dataMap["startScheduledMaintenanceStateRole"].(string); ok {
        data.StartScheduledMaintenanceStateRole = types.StringValue(val)
    } else {
        data.StartScheduledMaintenanceStateRole = types.StringNull()
    }
    if obj, ok := dataMap["endScheduledMaintenanceStateRole"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndScheduledMaintenanceStateRole = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EndScheduledMaintenanceStateRole = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EndScheduledMaintenanceStateRole = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EndScheduledMaintenanceStateRole = types.StringValue(string(jsonBytes))
            } else {
                data.EndScheduledMaintenanceStateRole = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EndScheduledMaintenanceStateRole = types.StringValue(string(jsonBytes))
            } else {
                data.EndScheduledMaintenanceStateRole = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EndScheduledMaintenanceStateRole = types.StringValue(string(jsonBytes))
        } else {
            data.EndScheduledMaintenanceStateRole = types.StringNull()
        }
    } else if val, ok := dataMap["endScheduledMaintenanceStateRole"].(string); ok {
        data.EndScheduledMaintenanceStateRole = types.StringValue(val)
    } else {
        data.EndScheduledMaintenanceStateRole = types.StringNull()
    }
    if obj, ok := dataMap["startStateOccurrence"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StartStateOccurrence = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StartStateOccurrence = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StartStateOccurrence = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.StartStateOccurrence = types.StringValue(string(jsonBytes))
            } else {
                data.StartStateOccurrence = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.StartStateOccurrence = types.StringValue(string(jsonBytes))
            } else {
                data.StartStateOccurrence = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StartStateOccurrence = types.StringValue(string(jsonBytes))
        } else {
            data.StartStateOccurrence = types.StringNull()
        }
    } else if val, ok := dataMap["startStateOccurrence"].(string); ok {
        data.StartStateOccurrence = types.StringValue(val)
    } else {
        data.StartStateOccurrence = types.StringNull()
    }
    if obj, ok := dataMap["endStateOccurrence"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EndStateOccurrence = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.EndStateOccurrence = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.EndStateOccurrence = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.EndStateOccurrence = types.StringValue(string(jsonBytes))
            } else {
                data.EndStateOccurrence = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.EndStateOccurrence = types.StringValue(string(jsonBytes))
            } else {
                data.EndStateOccurrence = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.EndStateOccurrence = types.StringValue(string(jsonBytes))
        } else {
            data.EndStateOccurrence = types.StringNull()
        }
    } else if val, ok := dataMap["endStateOccurrence"].(string); ok {
        data.EndStateOccurrence = types.StringValue(val)
    } else {
        data.EndStateOccurrence = types.StringNull()
    }
    if obj, ok := dataMap["unit"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Unit = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Unit = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Unit = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Unit = types.StringValue(string(jsonBytes))
            } else {
                data.Unit = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Unit = types.StringValue(string(jsonBytes))
            } else {
                data.Unit = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Unit = types.StringValue(string(jsonBytes))
        } else {
            data.Unit = types.StringNull()
        }
    } else if val, ok := dataMap["unit"].(string); ok {
        data.Unit = types.StringValue(val)
    } else {
        data.Unit = types.StringNull()
    }
    if obj, ok := dataMap["aggregationType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AggregationType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AggregationType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AggregationType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AggregationType = types.StringValue(string(jsonBytes))
            } else {
                data.AggregationType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AggregationType = types.StringValue(string(jsonBytes))
            } else {
                data.AggregationType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AggregationType = types.StringValue(string(jsonBytes))
        } else {
            data.AggregationType = types.StringNull()
        }
    } else if val, ok := dataMap["aggregationType"].(string); ok {
        data.AggregationType = types.StringValue(val)
    } else {
        data.AggregationType = types.StringNull()
    }
    if val, ok := dataMap["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := dataMap["showOnScheduledMaintenanceView"].(bool); ok {
        data.ShowOnScheduledMaintenanceView = types.BoolValue(val)
    }
    if val, ok := dataMap["order"].(float64); ok {
        data.Order = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["order"].(int); ok {
        data.Order = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["order"].(int64); ok {
        data.Order = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["order"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.Order = types.NumberValue(big.NewFloat(val))
        } else {
            data.Order = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.Order = types.NumberNull()
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
    if obj, ok := dataMap["metricName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MetricName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MetricName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MetricName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MetricName = types.StringValue(string(jsonBytes))
            } else {
                data.MetricName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MetricName = types.StringValue(string(jsonBytes))
            } else {
                data.MetricName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MetricName = types.StringValue(string(jsonBytes))
        } else {
            data.MetricName = types.StringNull()
        }
    } else if val, ok := dataMap["metricName"].(string); ok {
        data.MetricName = types.StringValue(val)
    } else {
        data.MetricName = types.StringNull()
    }
    if val, ok := dataMap["isSystemDefined"].(bool); ok {
        data.IsSystemDefined = types.BoolValue(val)
    }
    if obj, ok := dataMap["backfillRequestedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.BackfillRequestedAt = NewRFC3339Value(val)
        } else {
            data.BackfillRequestedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["backfillRequestedAt"].(string); ok && val != "" {
        data.BackfillRequestedAt = NewRFC3339Value(val)
    } else {
        data.BackfillRequestedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["backfillCursorCreatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.BackfillCursorCreatedAt = NewRFC3339Value(val)
        } else {
            data.BackfillCursorCreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["backfillCursorCreatedAt"].(string); ok && val != "" {
        data.BackfillCursorCreatedAt = NewRFC3339Value(val)
    } else {
        data.BackfillCursorCreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["backfillCompletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.BackfillCompletedAt = NewRFC3339Value(val)
        } else {
            data.BackfillCompletedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["backfillCompletedAt"].(string); ok && val != "" {
        data.BackfillCompletedAt = NewRFC3339Value(val)
    } else {
        data.BackfillCompletedAt = NewRFC3339Null()
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

func (r *ScheduledMaintenanceMeasurementResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data ScheduledMaintenanceMeasurementResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    httpResp, err := r.client.Delete(ctx, "/scheduled-maintenance-measurement/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete scheduled_maintenance_measurement, got error: %s", err))
        return
    }

    // A failed delete must keep the resource in state — silently dropping it
    // orphans real infrastructure. 404 means it is already gone.
    if httpResp.StatusCode >= 400 && httpResp.StatusCode != http.StatusNotFound {
        err = r.client.ParseResponse(httpResp, nil)
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to delete scheduled_maintenance_measurement: %s", err))
        return
    }
    if httpResp.Body != nil {
        httpResp.Body.Close()
    }
}


func (r *ScheduledMaintenanceMeasurementResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *ScheduledMaintenanceMeasurementResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *ScheduledMaintenanceMeasurementResource) convertTerraformListToInterface(terraformList types.List) interface{} {
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
func (r *ScheduledMaintenanceMeasurementResource) convertTerraformSetToInterface(terraformSet types.Set) interface{} {
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
func (r *ScheduledMaintenanceMeasurementResource) parseJSONField(terraformString basetypes.StringValuable) interface{} {
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
func (r *ScheduledMaintenanceMeasurementResource) normalizeURLWrappers(value interface{}) interface{} {
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

func (r *ScheduledMaintenanceMeasurementResource) normalizeURLString(value string) string {
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
func (r *ScheduledMaintenanceMeasurementResource) bigFloatToFloat64(bf *big.Float) interface{} {
    if bf == nil {
        return nil
    }
    f, _ := bf.Float64()
    return f
}

// Helper method to check if a type string is a valid OneUptime ObjectType.
// The registry itself lives in objecttypes.go, shared across the package.
func (r *ScheduledMaintenanceMeasurementResource) isValidOneUptimeObjectType(typeStr string) bool {
    return validOneUptimeObjectTypes[typeStr]
}
