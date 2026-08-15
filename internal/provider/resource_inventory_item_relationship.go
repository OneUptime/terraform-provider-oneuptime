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
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/numberplanmodifier"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &InventoryItemRelationshipResource{}
var _ resource.ResourceWithImportState = &InventoryItemRelationshipResource{}

func NewInventoryItemRelationshipResource() resource.Resource {
    return &InventoryItemRelationshipResource{}
}

// InventoryItemRelationshipResource defines the resource implementation.
type InventoryItemRelationshipResource struct {
    client *Client
}

// InventoryItemRelationshipResourceModel describes the resource data model.
type InventoryItemRelationshipResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    FromEntityKey types.String `tfsdk:"from_entity_key"`
    ToEntityKey types.String `tfsdk:"to_entity_key"`
    RelationshipType types.String `tfsdk:"relationship_type"`
    Source types.String `tfsdk:"source"`
    FirstSeenAt RFC3339Value `tfsdk:"first_seen_at"`
    LastSeenAt RFC3339Value `tfsdk:"last_seen_at"`
    CallCount types.Number `tfsdk:"call_count"`
    ErrorCount types.Number `tfsdk:"error_count"`
    AvgDurationMs types.Number `tfsdk:"avg_duration_ms"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
    CreatedAt RFC3339Value `tfsdk:"created_at"`
    UpdatedAt RFC3339Value `tfsdk:"updated_at"`
    DeletedAt RFC3339Value `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
}

func (r *InventoryItemRelationshipResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_inventory_item_relationship"
}

func (r *InventoryItemRelationshipResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Directed relationships between telemetry entities (runs-on, member-of, hosted-on, part-of, instance-of), inferred from resource co-occurrence.",

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
            "from_entity_key": schema.StringAttribute{
                MarkdownDescription: "Stable identity key of the source entity of this edge..",
                Required: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "to_entity_key": schema.StringAttribute{
                MarkdownDescription: "Stable identity key of the target entity of this edge..",
                Required: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "relationship_type": schema.StringAttribute{
                MarkdownDescription: "The inferred relationship (runs-on, member-of, hosted-on, part-of, instance-of)..",
                Required: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "source": schema.StringAttribute{
                MarkdownDescription: "Whether this edge was derived from telemetry or drawn manually by a user. Determines whether stale-edge pruning applies..",
                Required: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "first_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: RFC3339Type{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "call_count": schema.NumberAttribute{
                MarkdownDescription: "Calls observed over this edge in the most recent computation window (depends-on edges only)..",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "error_count": schema.NumberAttribute{
                MarkdownDescription: "Errored calls observed over this edge in the most recent computation window (depends-on edges only)..",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "avg_duration_ms": schema.NumberAttribute{
                MarkdownDescription: "Average call duration in milliseconds over this edge in the most recent computation window (depends-on edges only)..",
                Optional: true,
                Computed: true,
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
            "deleted_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
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
        },
    }
}

func (r *InventoryItemRelationshipResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *InventoryItemRelationshipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data InventoryItemRelationshipResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }



    // Create API request body. Unset (null/unknown) optional fields are
    // omitted so server-side defaults apply instead of being overwritten
    // with zero values.
    inventoryItemRelationshipRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := inventoryItemRelationshipRequest["data"].(map[string]interface{})

    if !data.FromEntityKey.IsNull() && !data.FromEntityKey.IsUnknown() {
        requestDataMap["fromEntityKey"] = data.FromEntityKey.ValueString()
    }
    if !data.ToEntityKey.IsNull() && !data.ToEntityKey.IsUnknown() {
        requestDataMap["toEntityKey"] = data.ToEntityKey.ValueString()
    }
    if !data.RelationshipType.IsNull() && !data.RelationshipType.IsUnknown() {
        requestDataMap["relationshipType"] = data.RelationshipType.ValueString()
    }
    if !data.Source.IsNull() && !data.Source.IsUnknown() {
        requestDataMap["source"] = data.Source.ValueString()
    }
    if !data.FirstSeenAt.IsNull() && !data.FirstSeenAt.IsUnknown() {
        requestDataMap["firstSeenAt"] = data.FirstSeenAt.ValueString()
    }
    if !data.LastSeenAt.IsNull() && !data.LastSeenAt.IsUnknown() {
        requestDataMap["lastSeenAt"] = data.LastSeenAt.ValueString()
    }
    if !data.CallCount.IsNull() && !data.CallCount.IsUnknown() {
        requestDataMap["callCount"] = r.bigFloatToFloat64(data.CallCount.ValueBigFloat())
    }
    if !data.ErrorCount.IsNull() && !data.ErrorCount.IsUnknown() {
        requestDataMap["errorCount"] = r.bigFloatToFloat64(data.ErrorCount.ValueBigFloat())
    }
    if !data.AvgDurationMs.IsNull() && !data.AvgDurationMs.IsUnknown() {
        requestDataMap["avgDurationMs"] = r.bigFloatToFloat64(data.AvgDurationMs.ValueBigFloat())
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        requestDataMap["createdByUserId"] = data.CreatedByUserId.ValueString()
    }
    if !data.DeletedByUserId.IsNull() && !data.DeletedByUserId.IsUnknown() {
        requestDataMap["deletedByUserId"] = data.DeletedByUserId.ValueString()
    }

    // Make API call
    httpResp, err := r.client.Post(ctx, "/inventory-item-relationship", inventoryItemRelationshipRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create inventory_item_relationship, got error: %s", err))
        return
    }

    var inventoryItemRelationshipResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &inventoryItemRelationshipResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to create inventory_item_relationship: %s", err))
        return
    }

    // Extract the new resource id from the create response.
    createdId := ""
    if wrapper, ok := inventoryItemRelationshipResponse["data"].(map[string]interface{}); ok {
        if val, ok := wrapper["_id"].(string); ok {
            createdId = val
        }
    } else if val, ok := inventoryItemRelationshipResponse["_id"].(string); ok {
        createdId = val
    }
    if createdId == "" {
        resp.Diagnostics.AddError("OneUptime API Error", "Create response for inventory_item_relationship did not contain an id. This is a bug in the provider or the API; please report it.")
        return
    }
    data.Id = types.StringValue(createdId)

    /*
     * The server has committed the row. Persist what we know to state BEFORE
     * the read-back: if the read-back fails and we return without setting
     * state, Terraform never learns the resource exists and the created
     * inventory_item_relationship is orphaned server-side — never refreshed, never
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
        "fromEntityKey": true,
        "toEntityKey": true,
        "relationshipType": true,
        "source": true,
        "firstSeenAt": true,
        "lastSeenAt": true,
        "callCount": true,
        "errorCount": true,
        "avgDurationMs": true,
        "createdByUserId": true,
        "deletedByUserId": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/inventory-item-relationship/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        /*
         * State already owns the id, so the resource is tracked and the next
         * refresh reconciles the remaining attributes. Warn rather than
         * error: erroring here would strand a real resource.
         */
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created inventory_item_relationship but could not read it back; state is incomplete until the next refresh: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created inventory_item_relationship but could not parse the read-back response; state is incomplete until the next refresh: %s", err))
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
    if obj, ok := dataMap["fromEntityKey"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FromEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.FromEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.FromEntityKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.FromEntityKey = types.StringValue(string(jsonBytes))
            } else {
                data.FromEntityKey = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.FromEntityKey = types.StringValue(string(jsonBytes))
            } else {
                data.FromEntityKey = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.FromEntityKey = types.StringValue(string(jsonBytes))
        } else {
            data.FromEntityKey = types.StringNull()
        }
    } else if val, ok := dataMap["fromEntityKey"].(string); ok {
        data.FromEntityKey = types.StringValue(val)
    } else {
        data.FromEntityKey = types.StringNull()
    }
    if obj, ok := dataMap["toEntityKey"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ToEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ToEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ToEntityKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ToEntityKey = types.StringValue(string(jsonBytes))
            } else {
                data.ToEntityKey = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ToEntityKey = types.StringValue(string(jsonBytes))
            } else {
                data.ToEntityKey = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ToEntityKey = types.StringValue(string(jsonBytes))
        } else {
            data.ToEntityKey = types.StringNull()
        }
    } else if val, ok := dataMap["toEntityKey"].(string); ok {
        data.ToEntityKey = types.StringValue(val)
    } else {
        data.ToEntityKey = types.StringNull()
    }
    if obj, ok := dataMap["relationshipType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RelationshipType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.RelationshipType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.RelationshipType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.RelationshipType = types.StringValue(string(jsonBytes))
            } else {
                data.RelationshipType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.RelationshipType = types.StringValue(string(jsonBytes))
            } else {
                data.RelationshipType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.RelationshipType = types.StringValue(string(jsonBytes))
        } else {
            data.RelationshipType = types.StringNull()
        }
    } else if val, ok := dataMap["relationshipType"].(string); ok {
        data.RelationshipType = types.StringValue(val)
    } else {
        data.RelationshipType = types.StringNull()
    }
    if obj, ok := dataMap["source"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Source = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Source = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Source = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Source = types.StringValue(string(jsonBytes))
            } else {
                data.Source = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Source = types.StringValue(string(jsonBytes))
            } else {
                data.Source = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Source = types.StringValue(string(jsonBytes))
        } else {
            data.Source = types.StringNull()
        }
    } else if val, ok := dataMap["source"].(string); ok {
        data.Source = types.StringValue(val)
    } else {
        data.Source = types.StringNull()
    }
    if obj, ok := dataMap["firstSeenAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.FirstSeenAt = NewRFC3339Value(val)
        } else {
            data.FirstSeenAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["firstSeenAt"].(string); ok && val != "" {
        data.FirstSeenAt = NewRFC3339Value(val)
    } else {
        data.FirstSeenAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastSeenAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastSeenAt = NewRFC3339Value(val)
        } else {
            data.LastSeenAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastSeenAt"].(string); ok && val != "" {
        data.LastSeenAt = NewRFC3339Value(val)
    } else {
        data.LastSeenAt = NewRFC3339Null()
    }
    if val, ok := dataMap["callCount"].(float64); ok {
        data.CallCount = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["callCount"].(int); ok {
        data.CallCount = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["callCount"].(int64); ok {
        data.CallCount = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["callCount"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.CallCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.CallCount = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.CallCount = types.NumberNull()
    }
    if val, ok := dataMap["errorCount"].(float64); ok {
        data.ErrorCount = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["errorCount"].(int); ok {
        data.ErrorCount = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["errorCount"].(int64); ok {
        data.ErrorCount = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["errorCount"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ErrorCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.ErrorCount = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ErrorCount = types.NumberNull()
    }
    if val, ok := dataMap["avgDurationMs"].(float64); ok {
        data.AvgDurationMs = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["avgDurationMs"].(int); ok {
        data.AvgDurationMs = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["avgDurationMs"].(int64); ok {
        data.AvgDurationMs = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["avgDurationMs"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AvgDurationMs = types.NumberValue(big.NewFloat(val))
        } else {
            data.AvgDurationMs = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AvgDurationMs = types.NumberNull()
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
    if obj, ok := dataMap["deletedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DeletedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DeletedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DeletedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["deletedByUserId"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    } else {
        data.DeletedByUserId = types.StringNull()
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

func (r *InventoryItemRelationshipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data InventoryItemRelationshipResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "fromEntityKey": true,
        "toEntityKey": true,
        "relationshipType": true,
        "source": true,
        "firstSeenAt": true,
        "lastSeenAt": true,
        "callCount": true,
        "errorCount": true,
        "avgDurationMs": true,
        "createdByUserId": true,
        "deletedByUserId": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect(ctx, "/inventory-item-relationship/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read inventory_item_relationship, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var inventoryItemRelationshipResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &inventoryItemRelationshipResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse inventory_item_relationship response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := inventoryItemRelationshipResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = inventoryItemRelationshipResponse
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
    if obj, ok := dataMap["fromEntityKey"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FromEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.FromEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.FromEntityKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.FromEntityKey = types.StringValue(string(jsonBytes))
            } else {
                data.FromEntityKey = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.FromEntityKey = types.StringValue(string(jsonBytes))
            } else {
                data.FromEntityKey = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.FromEntityKey = types.StringValue(string(jsonBytes))
        } else {
            data.FromEntityKey = types.StringNull()
        }
    } else if val, ok := dataMap["fromEntityKey"].(string); ok {
        data.FromEntityKey = types.StringValue(val)
    } else {
        data.FromEntityKey = types.StringNull()
    }
    if obj, ok := dataMap["toEntityKey"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ToEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ToEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ToEntityKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ToEntityKey = types.StringValue(string(jsonBytes))
            } else {
                data.ToEntityKey = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ToEntityKey = types.StringValue(string(jsonBytes))
            } else {
                data.ToEntityKey = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ToEntityKey = types.StringValue(string(jsonBytes))
        } else {
            data.ToEntityKey = types.StringNull()
        }
    } else if val, ok := dataMap["toEntityKey"].(string); ok {
        data.ToEntityKey = types.StringValue(val)
    } else {
        data.ToEntityKey = types.StringNull()
    }
    if obj, ok := dataMap["relationshipType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RelationshipType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.RelationshipType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.RelationshipType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.RelationshipType = types.StringValue(string(jsonBytes))
            } else {
                data.RelationshipType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.RelationshipType = types.StringValue(string(jsonBytes))
            } else {
                data.RelationshipType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.RelationshipType = types.StringValue(string(jsonBytes))
        } else {
            data.RelationshipType = types.StringNull()
        }
    } else if val, ok := dataMap["relationshipType"].(string); ok {
        data.RelationshipType = types.StringValue(val)
    } else {
        data.RelationshipType = types.StringNull()
    }
    if obj, ok := dataMap["source"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Source = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Source = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Source = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Source = types.StringValue(string(jsonBytes))
            } else {
                data.Source = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Source = types.StringValue(string(jsonBytes))
            } else {
                data.Source = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Source = types.StringValue(string(jsonBytes))
        } else {
            data.Source = types.StringNull()
        }
    } else if val, ok := dataMap["source"].(string); ok {
        data.Source = types.StringValue(val)
    } else {
        data.Source = types.StringNull()
    }
    if obj, ok := dataMap["firstSeenAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.FirstSeenAt = NewRFC3339Value(val)
        } else {
            data.FirstSeenAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["firstSeenAt"].(string); ok && val != "" {
        data.FirstSeenAt = NewRFC3339Value(val)
    } else {
        data.FirstSeenAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastSeenAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastSeenAt = NewRFC3339Value(val)
        } else {
            data.LastSeenAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastSeenAt"].(string); ok && val != "" {
        data.LastSeenAt = NewRFC3339Value(val)
    } else {
        data.LastSeenAt = NewRFC3339Null()
    }
    if val, ok := dataMap["callCount"].(float64); ok {
        data.CallCount = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["callCount"].(int); ok {
        data.CallCount = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["callCount"].(int64); ok {
        data.CallCount = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["callCount"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.CallCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.CallCount = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.CallCount = types.NumberNull()
    }
    if val, ok := dataMap["errorCount"].(float64); ok {
        data.ErrorCount = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["errorCount"].(int); ok {
        data.ErrorCount = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["errorCount"].(int64); ok {
        data.ErrorCount = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["errorCount"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ErrorCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.ErrorCount = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ErrorCount = types.NumberNull()
    }
    if val, ok := dataMap["avgDurationMs"].(float64); ok {
        data.AvgDurationMs = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["avgDurationMs"].(int); ok {
        data.AvgDurationMs = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["avgDurationMs"].(int64); ok {
        data.AvgDurationMs = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["avgDurationMs"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AvgDurationMs = types.NumberValue(big.NewFloat(val))
        } else {
            data.AvgDurationMs = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AvgDurationMs = types.NumberNull()
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
    if obj, ok := dataMap["deletedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DeletedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DeletedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DeletedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["deletedByUserId"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    } else {
        data.DeletedByUserId = types.StringNull()
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
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InventoryItemRelationshipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data InventoryItemRelationshipResourceModel
    var state InventoryItemRelationshipResourceModel

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
    inventoryItemRelationshipRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := inventoryItemRelationshipRequest["data"].(map[string]interface{})

    if !data.FirstSeenAt.IsUnknown() && !state.FirstSeenAt.IsUnknown() && !data.FirstSeenAt.Equal(state.FirstSeenAt) {
        requestDataMap["firstSeenAt"] = data.FirstSeenAt.ValueString()
    }
    if !data.LastSeenAt.IsUnknown() && !state.LastSeenAt.IsUnknown() && !data.LastSeenAt.Equal(state.LastSeenAt) {
        requestDataMap["lastSeenAt"] = data.LastSeenAt.ValueString()
    }
    if !data.CallCount.IsUnknown() && !state.CallCount.IsUnknown() && !data.CallCount.Equal(state.CallCount) {
        requestDataMap["callCount"] = r.bigFloatToFloat64(data.CallCount.ValueBigFloat())
    }
    if !data.ErrorCount.IsUnknown() && !state.ErrorCount.IsUnknown() && !data.ErrorCount.Equal(state.ErrorCount) {
        requestDataMap["errorCount"] = r.bigFloatToFloat64(data.ErrorCount.ValueBigFloat())
    }
    if !data.AvgDurationMs.IsUnknown() && !state.AvgDurationMs.IsUnknown() && !data.AvgDurationMs.Equal(state.AvgDurationMs) {
        requestDataMap["avgDurationMs"] = r.bigFloatToFloat64(data.AvgDurationMs.ValueBigFloat())
    }
    if !data.DeletedByUserId.IsUnknown() && !state.DeletedByUserId.IsUnknown() && !data.DeletedByUserId.Equal(state.DeletedByUserId) {
        requestDataMap["deletedByUserId"] = data.DeletedByUserId.ValueString()
    }

    // Only call the API when there are changed fields to send. An empty
    // update body is rejected by the API; state is still refreshed below so
    // this method never writes unverified plan values into state.
    if len(inventoryItemRelationshipRequest["data"].(map[string]interface{})) > 0 {
        httpResp, err := r.client.Put(ctx, "/inventory-item-relationship/" + data.Id.ValueString() + "", inventoryItemRelationshipRequest)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update inventory_item_relationship, got error: %s", err))
            return
        }

        // Parse the update response
        var inventoryItemRelationshipResponse map[string]interface{}
        err = r.client.ParseResponse(httpResp, &inventoryItemRelationshipResponse)
        if err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to update inventory_item_relationship: %s", err))
            return
        }
        _ = inventoryItemRelationshipResponse
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "fromEntityKey": true,
        "toEntityKey": true,
        "relationshipType": true,
        "source": true,
        "firstSeenAt": true,
        "lastSeenAt": true,
        "callCount": true,
        "errorCount": true,
        "avgDurationMs": true,
        "createdByUserId": true,
        "deletedByUserId": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/inventory-item-relationship/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read inventory_item_relationship after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read inventory_item_relationship after update: %s", err))
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
    if obj, ok := dataMap["fromEntityKey"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FromEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.FromEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.FromEntityKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.FromEntityKey = types.StringValue(string(jsonBytes))
            } else {
                data.FromEntityKey = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.FromEntityKey = types.StringValue(string(jsonBytes))
            } else {
                data.FromEntityKey = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.FromEntityKey = types.StringValue(string(jsonBytes))
        } else {
            data.FromEntityKey = types.StringNull()
        }
    } else if val, ok := dataMap["fromEntityKey"].(string); ok {
        data.FromEntityKey = types.StringValue(val)
    } else {
        data.FromEntityKey = types.StringNull()
    }
    if obj, ok := dataMap["toEntityKey"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ToEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ToEntityKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ToEntityKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ToEntityKey = types.StringValue(string(jsonBytes))
            } else {
                data.ToEntityKey = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ToEntityKey = types.StringValue(string(jsonBytes))
            } else {
                data.ToEntityKey = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ToEntityKey = types.StringValue(string(jsonBytes))
        } else {
            data.ToEntityKey = types.StringNull()
        }
    } else if val, ok := dataMap["toEntityKey"].(string); ok {
        data.ToEntityKey = types.StringValue(val)
    } else {
        data.ToEntityKey = types.StringNull()
    }
    if obj, ok := dataMap["relationshipType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RelationshipType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.RelationshipType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.RelationshipType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.RelationshipType = types.StringValue(string(jsonBytes))
            } else {
                data.RelationshipType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.RelationshipType = types.StringValue(string(jsonBytes))
            } else {
                data.RelationshipType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.RelationshipType = types.StringValue(string(jsonBytes))
        } else {
            data.RelationshipType = types.StringNull()
        }
    } else if val, ok := dataMap["relationshipType"].(string); ok {
        data.RelationshipType = types.StringValue(val)
    } else {
        data.RelationshipType = types.StringNull()
    }
    if obj, ok := dataMap["source"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Source = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Source = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Source = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Source = types.StringValue(string(jsonBytes))
            } else {
                data.Source = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Source = types.StringValue(string(jsonBytes))
            } else {
                data.Source = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Source = types.StringValue(string(jsonBytes))
        } else {
            data.Source = types.StringNull()
        }
    } else if val, ok := dataMap["source"].(string); ok {
        data.Source = types.StringValue(val)
    } else {
        data.Source = types.StringNull()
    }
    if obj, ok := dataMap["firstSeenAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.FirstSeenAt = NewRFC3339Value(val)
        } else {
            data.FirstSeenAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["firstSeenAt"].(string); ok && val != "" {
        data.FirstSeenAt = NewRFC3339Value(val)
    } else {
        data.FirstSeenAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastSeenAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastSeenAt = NewRFC3339Value(val)
        } else {
            data.LastSeenAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastSeenAt"].(string); ok && val != "" {
        data.LastSeenAt = NewRFC3339Value(val)
    } else {
        data.LastSeenAt = NewRFC3339Null()
    }
    if val, ok := dataMap["callCount"].(float64); ok {
        data.CallCount = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["callCount"].(int); ok {
        data.CallCount = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["callCount"].(int64); ok {
        data.CallCount = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["callCount"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.CallCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.CallCount = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.CallCount = types.NumberNull()
    }
    if val, ok := dataMap["errorCount"].(float64); ok {
        data.ErrorCount = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["errorCount"].(int); ok {
        data.ErrorCount = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["errorCount"].(int64); ok {
        data.ErrorCount = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["errorCount"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.ErrorCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.ErrorCount = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.ErrorCount = types.NumberNull()
    }
    if val, ok := dataMap["avgDurationMs"].(float64); ok {
        data.AvgDurationMs = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["avgDurationMs"].(int); ok {
        data.AvgDurationMs = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["avgDurationMs"].(int64); ok {
        data.AvgDurationMs = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["avgDurationMs"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AvgDurationMs = types.NumberValue(big.NewFloat(val))
        } else {
            data.AvgDurationMs = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AvgDurationMs = types.NumberNull()
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
    if obj, ok := dataMap["deletedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DeletedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DeletedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DeletedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["deletedByUserId"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    } else {
        data.DeletedByUserId = types.StringNull()
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
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    data.Id = state.Id

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InventoryItemRelationshipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data InventoryItemRelationshipResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    httpResp, err := r.client.Delete(ctx, "/inventory-item-relationship/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete inventory_item_relationship, got error: %s", err))
        return
    }

    // A failed delete must keep the resource in state — silently dropping it
    // orphans real infrastructure. 404 means it is already gone.
    if httpResp.StatusCode >= 400 && httpResp.StatusCode != http.StatusNotFound {
        err = r.client.ParseResponse(httpResp, nil)
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to delete inventory_item_relationship: %s", err))
        return
    }
    if httpResp.Body != nil {
        httpResp.Body.Close()
    }
}


func (r *InventoryItemRelationshipResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *InventoryItemRelationshipResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *InventoryItemRelationshipResource) convertTerraformListToInterface(terraformList types.List) interface{} {
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
func (r *InventoryItemRelationshipResource) convertTerraformSetToInterface(terraformSet types.Set) interface{} {
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
func (r *InventoryItemRelationshipResource) parseJSONField(terraformString basetypes.StringValuable) interface{} {
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
func (r *InventoryItemRelationshipResource) normalizeURLWrappers(value interface{}) interface{} {
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

func (r *InventoryItemRelationshipResource) normalizeURLString(value string) string {
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
func (r *InventoryItemRelationshipResource) bigFloatToFloat64(bf *big.Float) interface{} {
    if bf == nil {
        return nil
    }
    f, _ := bf.Float64()
    return f
}

// Helper method to check if a type string is a valid OneUptime ObjectType.
// The registry itself lives in objecttypes.go, shared across the package.
func (r *InventoryItemRelationshipResource) isValidOneUptimeObjectType(typeStr string) bool {
    return validOneUptimeObjectTypes[typeStr]
}
