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
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &HuntressConnectionResource{}
var _ resource.ResourceWithImportState = &HuntressConnectionResource{}

func NewHuntressConnectionResource() resource.Resource {
    return &HuntressConnectionResource{}
}

// HuntressConnectionResource defines the resource implementation.
type HuntressConnectionResource struct {
    client *Client
}

// HuntressConnectionResourceModel describes the resource data model.
type HuntressConnectionResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    Name JSONSubsetValue `tfsdk:"name"`
    SigningSecret types.String `tfsdk:"signing_secret"`
    IsSigningSecretSet types.Bool `tfsdk:"is_signing_secret_set"`
    OnCallDutyPolicies types.Set `tfsdk:"on_call_duty_policies"`
    PageOnCallFor types.String `tfsdk:"page_on_call_for"`
    CriticalIncidentSeverityId types.String `tfsdk:"critical_incident_severity_id"`
    HighIncidentSeverityId types.String `tfsdk:"high_incident_severity_id"`
    LowIncidentSeverityId types.String `tfsdk:"low_incident_severity_id"`
    WatchedOrganizations types.String `tfsdk:"watched_organizations"`
    Labels types.Set `tfsdk:"labels"`
    ResolveIncidentWhenReportCloses types.Bool `tfsdk:"resolve_incident_when_report_closes"`
    CreatedAt RFC3339Value `tfsdk:"created_at"`
    UpdatedAt RFC3339Value `tfsdk:"updated_at"`
    LastEventReceivedAt RFC3339Value `tfsdk:"last_event_received_at"`
    LastEventType types.String `tfsdk:"last_event_type"`
    LastError types.String `tfsdk:"last_error"`
    LastErrorAt RFC3339Value `tfsdk:"last_error_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (r *HuntressConnectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_huntress_connection"
}

func (r *HuntressConnectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = r.schemaDefinition()
}

func (r *HuntressConnectionResource) schemaDefinition() schema.Schema {
    return schema.Schema{
        MarkdownDescription: "Huntress webhook endpoints. Every Huntress incident report opens an incident that pages on-call, and closing the report in Huntress resolves it.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of the project this connection belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "A name for this connection, such as the Huntress account it receives from.",
                CustomType: JSONSubsetType{},
                Required: true,
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "signing_secret": schema.StringAttribute{
                MarkdownDescription: "The endpoint's signing secret from Huntress (whsec_...), used to verify that a request really came from Huntress. Encrypted at rest and never returned by the API.",
                Optional: true,
            },
            "is_signing_secret_set": schema.BoolAttribute{
                MarkdownDescription: "Whether a signing secret is saved. Requests are only accepted once it is. Set from the signing secret itself; a value sent for it is ignored.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "on_call_duty_policies": schema.SetAttribute{
                MarkdownDescription: "The on-call policies paged for an incident report at or above the Page On-Call For severity. IDs of `oneuptime_on_call_policy` resources.",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
            },
            "page_on_call_for": schema.StringAttribute{
                MarkdownDescription: "The lowest Huntress severity that pages the on-call policies: critical (critical reports only), high (high and critical reports) or low (every report). Every report opens an incident either way.",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("high"),
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "critical_incident_severity_id": schema.StringAttribute{
                MarkdownDescription: "ID of the incident severity critical reports open at. The ID of a `oneuptime_incident_severity`.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "high_incident_severity_id": schema.StringAttribute{
                MarkdownDescription: "ID of the incident severity high reports open at. The ID of a `oneuptime_incident_severity`.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "low_incident_severity_id": schema.StringAttribute{
                MarkdownDescription: "ID of the incident severity low reports open at. The ID of a `oneuptime_incident_severity`.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "watched_organizations": schema.StringAttribute{
                MarkdownDescription: "The Huntress organizations this connection opens incidents for, one per line, by name or by id. Empty watches every organization.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Labels added to every incident this connection opens, next to the label named after the report's organization. IDs of `oneuptime_label` resources.",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
            },
            "resolve_incident_when_report_closes": schema.BoolAttribute{
                MarkdownDescription: "Resolve the incident when the report is closed or dismissed in Huntress. When off, a note on the incident says so instead.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
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
            "last_event_received_at": schema.StringAttribute{
                MarkdownDescription: "When a request signed with the signing secret last arrived. Empty until Huntress sends the first one.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "last_event_type": schema.StringAttribute{
                MarkdownDescription: "The type of the last event received, such as incident_report.created.",
                Computed: true,
            },
            "last_error": schema.StringAttribute{
                MarkdownDescription: "Why the last request was refused or could not be handled, if it was. Cleared by the next request that is handled.",
                Computed: true,
            },
            "last_error_at": schema.StringAttribute{
                MarkdownDescription: "When the last error happened.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the user who created this connection. The ID of a `oneuptime_user` (see the data source).",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
        },
    }
}

func (r *HuntressConnectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *HuntressConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data HuntressConnectionResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // What the configuration sets, and what Terraform planned (see keepPlannedValues).
    var config HuntressConnectionResourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
    if resp.Diagnostics.HasError() {
        return
    }
    plan := data



    // Create API request body. Unset (null/unknown) optional fields are
    // omitted so server-side defaults apply instead of being overwritten
    // with zero values.
    huntressConnectionRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := huntressConnectionRequest["data"].(map[string]interface{})

    if parsedName := r.parseJSONField(data.Name); parsedName != nil {
        requestDataMap["name"] = parsedName
    }
    if !data.SigningSecret.IsNull() && !data.SigningSecret.IsUnknown() {
        requestDataMap["signingSecret"] = data.SigningSecret.ValueString()
    }
    if !data.IsSigningSecretSet.IsNull() && !data.IsSigningSecretSet.IsUnknown() {
        requestDataMap["isSigningSecretSet"] = data.IsSigningSecretSet.ValueBool()
    }
    if !data.OnCallDutyPolicies.IsNull() && !data.OnCallDutyPolicies.IsUnknown() {
        requestDataMap["onCallDutyPolicies"] = r.convertTerraformSetToInterface(data.OnCallDutyPolicies)
    }
    if !data.PageOnCallFor.IsNull() && !data.PageOnCallFor.IsUnknown() {
        requestDataMap["pageOnCallFor"] = data.PageOnCallFor.ValueString()
    }
    if !data.CriticalIncidentSeverityId.IsNull() && !data.CriticalIncidentSeverityId.IsUnknown() {
        requestDataMap["criticalIncidentSeverityId"] = data.CriticalIncidentSeverityId.ValueString()
    }
    if !data.HighIncidentSeverityId.IsNull() && !data.HighIncidentSeverityId.IsUnknown() {
        requestDataMap["highIncidentSeverityId"] = data.HighIncidentSeverityId.ValueString()
    }
    if !data.LowIncidentSeverityId.IsNull() && !data.LowIncidentSeverityId.IsUnknown() {
        requestDataMap["lowIncidentSeverityId"] = data.LowIncidentSeverityId.ValueString()
    }
    if !data.WatchedOrganizations.IsNull() && !data.WatchedOrganizations.IsUnknown() {
        requestDataMap["watchedOrganizations"] = data.WatchedOrganizations.ValueString()
    }
    if !data.Labels.IsNull() && !data.Labels.IsUnknown() {
        requestDataMap["labels"] = r.convertTerraformSetToInterface(data.Labels)
    }
    if !data.ResolveIncidentWhenReportCloses.IsNull() && !data.ResolveIncidentWhenReportCloses.IsUnknown() {
        requestDataMap["resolveIncidentWhenReportCloses"] = data.ResolveIncidentWhenReportCloses.ValueBool()
    }

    // Make API call
    httpResp, err := r.client.Post(ctx, "/huntress-connection", huntressConnectionRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create huntress_connection, got error: %s", err))
        return
    }

    var huntressConnectionResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &huntressConnectionResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to create huntress_connection: %s", err))
        return
    }

    // Extract the new resource id from the create response.
    createdId := ""
    if wrapper, ok := huntressConnectionResponse["data"].(map[string]interface{}); ok {
        if val, ok := wrapper["_id"].(string); ok {
            createdId = val
        }
    } else if val, ok := huntressConnectionResponse["_id"].(string); ok {
        createdId = val
    }
    if createdId == "" {
        resp.Diagnostics.AddError("OneUptime API Error", "Create response for huntress_connection did not contain an id. This is a bug in the provider or the API; please report it.")
        return
    }
    data.Id = types.StringValue(createdId)

    /*
     * The server has committed the row. Persist what we know to state BEFORE
     * the read-back: if the read-back fails and we return without setting
     * state, Terraform never learns the resource exists and the created
     * huntress_connection is orphaned server-side — never refreshed, never
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
        "isSigningSecretSet": true,
        "onCallDutyPolicies": true,
        "pageOnCallFor": true,
        "criticalIncidentSeverityId": true,
        "highIncidentSeverityId": true,
        "lowIncidentSeverityId": true,
        "watchedOrganizations": true,
        "labels": true,
        "resolveIncidentWhenReportCloses": true,
        "createdAt": true,
        "updatedAt": true,
        "lastEventReceivedAt": true,
        "lastEventType": true,
        "lastError": true,
        "lastErrorAt": true,
        "createdByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/huntress-connection/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        /*
         * State already owns the id, so the resource is tracked and the next
         * refresh reconciles the remaining attributes. Warn rather than
         * error: erroring here would strand a real resource.
         */
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created huntress_connection but could not read it back; state is incomplete until the next refresh: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created huntress_connection but could not parse the read-back response; state is incomplete until the next refresh: %s", err))
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
    if val, ok := dataMap["isSigningSecretSet"].(bool); ok {
        data.IsSigningSecretSet = types.BoolValue(val)
    }
    if val, ok := dataMap["onCallDutyPolicies"].([]interface{}); ok {
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
        data.OnCallDutyPolicies = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.OnCallDutyPolicies = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["pageOnCallFor"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PageOnCallFor = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PageOnCallFor = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PageOnCallFor = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PageOnCallFor = types.StringValue(string(jsonBytes))
            } else {
                data.PageOnCallFor = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PageOnCallFor = types.StringValue(string(jsonBytes))
            } else {
                data.PageOnCallFor = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PageOnCallFor = types.StringValue(string(jsonBytes))
        } else {
            data.PageOnCallFor = types.StringNull()
        }
    } else if val, ok := dataMap["pageOnCallFor"].(string); ok {
        data.PageOnCallFor = types.StringValue(val)
    } else {
        data.PageOnCallFor = types.StringNull()
    }
    if obj, ok := dataMap["criticalIncidentSeverityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CriticalIncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CriticalIncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CriticalIncidentSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CriticalIncidentSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.CriticalIncidentSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CriticalIncidentSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.CriticalIncidentSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CriticalIncidentSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.CriticalIncidentSeverityId = types.StringNull()
        }
    } else if val, ok := dataMap["criticalIncidentSeverityId"].(string); ok {
        data.CriticalIncidentSeverityId = types.StringValue(val)
    } else {
        data.CriticalIncidentSeverityId = types.StringNull()
    }
    if obj, ok := dataMap["highIncidentSeverityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HighIncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.HighIncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.HighIncidentSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.HighIncidentSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.HighIncidentSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.HighIncidentSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.HighIncidentSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.HighIncidentSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.HighIncidentSeverityId = types.StringNull()
        }
    } else if val, ok := dataMap["highIncidentSeverityId"].(string); ok {
        data.HighIncidentSeverityId = types.StringValue(val)
    } else {
        data.HighIncidentSeverityId = types.StringNull()
    }
    if obj, ok := dataMap["lowIncidentSeverityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LowIncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LowIncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LowIncidentSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LowIncidentSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.LowIncidentSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LowIncidentSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.LowIncidentSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LowIncidentSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.LowIncidentSeverityId = types.StringNull()
        }
    } else if val, ok := dataMap["lowIncidentSeverityId"].(string); ok {
        data.LowIncidentSeverityId = types.StringValue(val)
    } else {
        data.LowIncidentSeverityId = types.StringNull()
    }
    if obj, ok := dataMap["watchedOrganizations"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.WatchedOrganizations = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.WatchedOrganizations = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.WatchedOrganizations = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.WatchedOrganizations = types.StringValue(string(jsonBytes))
            } else {
                data.WatchedOrganizations = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.WatchedOrganizations = types.StringValue(string(jsonBytes))
            } else {
                data.WatchedOrganizations = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.WatchedOrganizations = types.StringValue(string(jsonBytes))
        } else {
            data.WatchedOrganizations = types.StringNull()
        }
    } else if val, ok := dataMap["watchedOrganizations"].(string); ok {
        data.WatchedOrganizations = types.StringValue(val)
    } else {
        data.WatchedOrganizations = types.StringNull()
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
    if val, ok := dataMap["resolveIncidentWhenReportCloses"].(bool); ok {
        data.ResolveIncidentWhenReportCloses = types.BoolValue(val)
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
    if obj, ok := dataMap["lastEventReceivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastEventReceivedAt = NewRFC3339Value(val)
        } else {
            data.LastEventReceivedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastEventReceivedAt"].(string); ok && val != "" {
        data.LastEventReceivedAt = NewRFC3339Value(val)
    } else {
        data.LastEventReceivedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastEventType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastEventType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastEventType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastEventType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastEventType = types.StringValue(string(jsonBytes))
            } else {
                data.LastEventType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastEventType = types.StringValue(string(jsonBytes))
            } else {
                data.LastEventType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastEventType = types.StringValue(string(jsonBytes))
        } else {
            data.LastEventType = types.StringNull()
        }
    } else if val, ok := dataMap["lastEventType"].(string); ok {
        data.LastEventType = types.StringValue(val)
    } else {
        data.LastEventType = types.StringNull()
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
    if obj, ok := dataMap["lastErrorAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastErrorAt = NewRFC3339Value(val)
        } else {
            data.LastErrorAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastErrorAt"].(string); ok && val != "" {
        data.LastErrorAt = NewRFC3339Value(val)
    } else {
        data.LastErrorAt = NewRFC3339Null()
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

func (r *HuntressConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data HuntressConnectionResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "isSigningSecretSet": true,
        "onCallDutyPolicies": true,
        "pageOnCallFor": true,
        "criticalIncidentSeverityId": true,
        "highIncidentSeverityId": true,
        "lowIncidentSeverityId": true,
        "watchedOrganizations": true,
        "labels": true,
        "resolveIncidentWhenReportCloses": true,
        "createdAt": true,
        "updatedAt": true,
        "lastEventReceivedAt": true,
        "lastEventType": true,
        "lastError": true,
        "lastErrorAt": true,
        "createdByUserId": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect(ctx, "/huntress-connection/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read huntress_connection, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var huntressConnectionResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &huntressConnectionResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse huntress_connection response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := huntressConnectionResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = huntressConnectionResponse
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
    if val, ok := dataMap["isSigningSecretSet"].(bool); ok {
        data.IsSigningSecretSet = types.BoolValue(val)
    }
    if val, ok := dataMap["onCallDutyPolicies"].([]interface{}); ok {
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
        data.OnCallDutyPolicies = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.OnCallDutyPolicies = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["pageOnCallFor"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PageOnCallFor = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PageOnCallFor = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PageOnCallFor = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PageOnCallFor = types.StringValue(string(jsonBytes))
            } else {
                data.PageOnCallFor = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PageOnCallFor = types.StringValue(string(jsonBytes))
            } else {
                data.PageOnCallFor = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PageOnCallFor = types.StringValue(string(jsonBytes))
        } else {
            data.PageOnCallFor = types.StringNull()
        }
    } else if val, ok := dataMap["pageOnCallFor"].(string); ok {
        data.PageOnCallFor = types.StringValue(val)
    } else {
        data.PageOnCallFor = types.StringNull()
    }
    if obj, ok := dataMap["criticalIncidentSeverityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CriticalIncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CriticalIncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CriticalIncidentSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CriticalIncidentSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.CriticalIncidentSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CriticalIncidentSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.CriticalIncidentSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CriticalIncidentSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.CriticalIncidentSeverityId = types.StringNull()
        }
    } else if val, ok := dataMap["criticalIncidentSeverityId"].(string); ok {
        data.CriticalIncidentSeverityId = types.StringValue(val)
    } else {
        data.CriticalIncidentSeverityId = types.StringNull()
    }
    if obj, ok := dataMap["highIncidentSeverityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HighIncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.HighIncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.HighIncidentSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.HighIncidentSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.HighIncidentSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.HighIncidentSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.HighIncidentSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.HighIncidentSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.HighIncidentSeverityId = types.StringNull()
        }
    } else if val, ok := dataMap["highIncidentSeverityId"].(string); ok {
        data.HighIncidentSeverityId = types.StringValue(val)
    } else {
        data.HighIncidentSeverityId = types.StringNull()
    }
    if obj, ok := dataMap["lowIncidentSeverityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LowIncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LowIncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LowIncidentSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LowIncidentSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.LowIncidentSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LowIncidentSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.LowIncidentSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LowIncidentSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.LowIncidentSeverityId = types.StringNull()
        }
    } else if val, ok := dataMap["lowIncidentSeverityId"].(string); ok {
        data.LowIncidentSeverityId = types.StringValue(val)
    } else {
        data.LowIncidentSeverityId = types.StringNull()
    }
    if obj, ok := dataMap["watchedOrganizations"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.WatchedOrganizations = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.WatchedOrganizations = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.WatchedOrganizations = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.WatchedOrganizations = types.StringValue(string(jsonBytes))
            } else {
                data.WatchedOrganizations = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.WatchedOrganizations = types.StringValue(string(jsonBytes))
            } else {
                data.WatchedOrganizations = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.WatchedOrganizations = types.StringValue(string(jsonBytes))
        } else {
            data.WatchedOrganizations = types.StringNull()
        }
    } else if val, ok := dataMap["watchedOrganizations"].(string); ok {
        data.WatchedOrganizations = types.StringValue(val)
    } else {
        data.WatchedOrganizations = types.StringNull()
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
    if val, ok := dataMap["resolveIncidentWhenReportCloses"].(bool); ok {
        data.ResolveIncidentWhenReportCloses = types.BoolValue(val)
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
    if obj, ok := dataMap["lastEventReceivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastEventReceivedAt = NewRFC3339Value(val)
        } else {
            data.LastEventReceivedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastEventReceivedAt"].(string); ok && val != "" {
        data.LastEventReceivedAt = NewRFC3339Value(val)
    } else {
        data.LastEventReceivedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastEventType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastEventType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastEventType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastEventType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastEventType = types.StringValue(string(jsonBytes))
            } else {
                data.LastEventType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastEventType = types.StringValue(string(jsonBytes))
            } else {
                data.LastEventType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastEventType = types.StringValue(string(jsonBytes))
        } else {
            data.LastEventType = types.StringNull()
        }
    } else if val, ok := dataMap["lastEventType"].(string); ok {
        data.LastEventType = types.StringValue(val)
    } else {
        data.LastEventType = types.StringNull()
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
    if obj, ok := dataMap["lastErrorAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastErrorAt = NewRFC3339Value(val)
        } else {
            data.LastErrorAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastErrorAt"].(string); ok && val != "" {
        data.LastErrorAt = NewRFC3339Value(val)
    } else {
        data.LastErrorAt = NewRFC3339Null()
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

func (r *HuntressConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data HuntressConnectionResourceModel
    var state HuntressConnectionResourceModel

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
    var config HuntressConnectionResourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
    if resp.Diagnostics.HasError() {
        return
    }
    plan := data

    // Create API request body
    huntressConnectionRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := huntressConnectionRequest["data"].(map[string]interface{})

    if !data.Name.IsUnknown() && !state.Name.IsUnknown() && !data.Name.Equal(state.Name) {
        var nameData interface{}
        if err := json.Unmarshal([]byte(data.Name.ValueString()), &nameData); err == nil {
            requestDataMap["name"] = nameData
        } else {
            requestDataMap["name"] = data.Name.ValueString()
        }
    }
    if !data.SigningSecret.IsUnknown() && !state.SigningSecret.IsUnknown() && !data.SigningSecret.Equal(state.SigningSecret) {
        requestDataMap["signingSecret"] = data.SigningSecret.ValueString()
    }
    if !data.IsSigningSecretSet.IsUnknown() && !state.IsSigningSecretSet.IsUnknown() && !data.IsSigningSecretSet.Equal(state.IsSigningSecretSet) {
        requestDataMap["isSigningSecretSet"] = data.IsSigningSecretSet.ValueBool()
    }
    if !data.OnCallDutyPolicies.IsUnknown() && !state.OnCallDutyPolicies.IsUnknown() && !data.OnCallDutyPolicies.Equal(state.OnCallDutyPolicies) {
        requestDataMap["onCallDutyPolicies"] = r.convertTerraformSetToInterface(data.OnCallDutyPolicies)
    }
    if !data.PageOnCallFor.IsUnknown() && !state.PageOnCallFor.IsUnknown() && !data.PageOnCallFor.Equal(state.PageOnCallFor) {
        requestDataMap["pageOnCallFor"] = data.PageOnCallFor.ValueString()
    }
    if !data.CriticalIncidentSeverityId.IsUnknown() && !state.CriticalIncidentSeverityId.IsUnknown() && !data.CriticalIncidentSeverityId.Equal(state.CriticalIncidentSeverityId) {
        requestDataMap["criticalIncidentSeverityId"] = data.CriticalIncidentSeverityId.ValueString()
    }
    if !data.HighIncidentSeverityId.IsUnknown() && !state.HighIncidentSeverityId.IsUnknown() && !data.HighIncidentSeverityId.Equal(state.HighIncidentSeverityId) {
        requestDataMap["highIncidentSeverityId"] = data.HighIncidentSeverityId.ValueString()
    }
    if !data.LowIncidentSeverityId.IsUnknown() && !state.LowIncidentSeverityId.IsUnknown() && !data.LowIncidentSeverityId.Equal(state.LowIncidentSeverityId) {
        requestDataMap["lowIncidentSeverityId"] = data.LowIncidentSeverityId.ValueString()
    }
    if !data.WatchedOrganizations.IsUnknown() && !state.WatchedOrganizations.IsUnknown() && !data.WatchedOrganizations.Equal(state.WatchedOrganizations) {
        requestDataMap["watchedOrganizations"] = data.WatchedOrganizations.ValueString()
    }
    if !data.Labels.IsUnknown() && !state.Labels.IsUnknown() && !data.Labels.Equal(state.Labels) {
        requestDataMap["labels"] = r.convertTerraformSetToInterface(data.Labels)
    }
    if !data.ResolveIncidentWhenReportCloses.IsUnknown() && !state.ResolveIncidentWhenReportCloses.IsUnknown() && !data.ResolveIncidentWhenReportCloses.Equal(state.ResolveIncidentWhenReportCloses) {
        requestDataMap["resolveIncidentWhenReportCloses"] = data.ResolveIncidentWhenReportCloses.ValueBool()
    }

    // Only call the API when there are changed fields to send. An empty
    // update body is rejected by the API; state is still refreshed below so
    // this method never writes unverified plan values into state.
    if len(huntressConnectionRequest["data"].(map[string]interface{})) > 0 {
        httpResp, err := r.client.Put(ctx, "/huntress-connection/" + data.Id.ValueString() + "", huntressConnectionRequest)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update huntress_connection, got error: %s", err))
            return
        }

        // Parse the update response
        var huntressConnectionResponse map[string]interface{}
        err = r.client.ParseResponse(httpResp, &huntressConnectionResponse)
        if err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to update huntress_connection: %s", err))
            return
        }
        _ = huntressConnectionResponse
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "isSigningSecretSet": true,
        "onCallDutyPolicies": true,
        "pageOnCallFor": true,
        "criticalIncidentSeverityId": true,
        "highIncidentSeverityId": true,
        "lowIncidentSeverityId": true,
        "watchedOrganizations": true,
        "labels": true,
        "resolveIncidentWhenReportCloses": true,
        "createdAt": true,
        "updatedAt": true,
        "lastEventReceivedAt": true,
        "lastEventType": true,
        "lastError": true,
        "lastErrorAt": true,
        "createdByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/huntress-connection/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read huntress_connection after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read huntress_connection after update: %s", err))
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
    if val, ok := dataMap["isSigningSecretSet"].(bool); ok {
        data.IsSigningSecretSet = types.BoolValue(val)
    }
    if val, ok := dataMap["onCallDutyPolicies"].([]interface{}); ok {
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
        data.OnCallDutyPolicies = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.OnCallDutyPolicies = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["pageOnCallFor"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PageOnCallFor = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PageOnCallFor = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PageOnCallFor = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PageOnCallFor = types.StringValue(string(jsonBytes))
            } else {
                data.PageOnCallFor = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PageOnCallFor = types.StringValue(string(jsonBytes))
            } else {
                data.PageOnCallFor = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PageOnCallFor = types.StringValue(string(jsonBytes))
        } else {
            data.PageOnCallFor = types.StringNull()
        }
    } else if val, ok := dataMap["pageOnCallFor"].(string); ok {
        data.PageOnCallFor = types.StringValue(val)
    } else {
        data.PageOnCallFor = types.StringNull()
    }
    if obj, ok := dataMap["criticalIncidentSeverityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CriticalIncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CriticalIncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CriticalIncidentSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CriticalIncidentSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.CriticalIncidentSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CriticalIncidentSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.CriticalIncidentSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CriticalIncidentSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.CriticalIncidentSeverityId = types.StringNull()
        }
    } else if val, ok := dataMap["criticalIncidentSeverityId"].(string); ok {
        data.CriticalIncidentSeverityId = types.StringValue(val)
    } else {
        data.CriticalIncidentSeverityId = types.StringNull()
    }
    if obj, ok := dataMap["highIncidentSeverityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HighIncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.HighIncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.HighIncidentSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.HighIncidentSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.HighIncidentSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.HighIncidentSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.HighIncidentSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.HighIncidentSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.HighIncidentSeverityId = types.StringNull()
        }
    } else if val, ok := dataMap["highIncidentSeverityId"].(string); ok {
        data.HighIncidentSeverityId = types.StringValue(val)
    } else {
        data.HighIncidentSeverityId = types.StringNull()
    }
    if obj, ok := dataMap["lowIncidentSeverityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LowIncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LowIncidentSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LowIncidentSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LowIncidentSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.LowIncidentSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LowIncidentSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.LowIncidentSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LowIncidentSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.LowIncidentSeverityId = types.StringNull()
        }
    } else if val, ok := dataMap["lowIncidentSeverityId"].(string); ok {
        data.LowIncidentSeverityId = types.StringValue(val)
    } else {
        data.LowIncidentSeverityId = types.StringNull()
    }
    if obj, ok := dataMap["watchedOrganizations"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.WatchedOrganizations = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.WatchedOrganizations = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.WatchedOrganizations = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.WatchedOrganizations = types.StringValue(string(jsonBytes))
            } else {
                data.WatchedOrganizations = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.WatchedOrganizations = types.StringValue(string(jsonBytes))
            } else {
                data.WatchedOrganizations = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.WatchedOrganizations = types.StringValue(string(jsonBytes))
        } else {
            data.WatchedOrganizations = types.StringNull()
        }
    } else if val, ok := dataMap["watchedOrganizations"].(string); ok {
        data.WatchedOrganizations = types.StringValue(val)
    } else {
        data.WatchedOrganizations = types.StringNull()
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
    if val, ok := dataMap["resolveIncidentWhenReportCloses"].(bool); ok {
        data.ResolveIncidentWhenReportCloses = types.BoolValue(val)
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
    if obj, ok := dataMap["lastEventReceivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastEventReceivedAt = NewRFC3339Value(val)
        } else {
            data.LastEventReceivedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastEventReceivedAt"].(string); ok && val != "" {
        data.LastEventReceivedAt = NewRFC3339Value(val)
    } else {
        data.LastEventReceivedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastEventType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastEventType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastEventType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastEventType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastEventType = types.StringValue(string(jsonBytes))
            } else {
                data.LastEventType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastEventType = types.StringValue(string(jsonBytes))
            } else {
                data.LastEventType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastEventType = types.StringValue(string(jsonBytes))
        } else {
            data.LastEventType = types.StringNull()
        }
    } else if val, ok := dataMap["lastEventType"].(string); ok {
        data.LastEventType = types.StringValue(val)
    } else {
        data.LastEventType = types.StringNull()
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
    if obj, ok := dataMap["lastErrorAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastErrorAt = NewRFC3339Value(val)
        } else {
            data.LastErrorAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastErrorAt"].(string); ok && val != "" {
        data.LastErrorAt = NewRFC3339Value(val)
    } else {
        data.LastErrorAt = NewRFC3339Null()
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

func (r *HuntressConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data HuntressConnectionResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    httpResp, err := r.client.Delete(ctx, "/huntress-connection/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete huntress_connection, got error: %s", err))
        return
    }

    // A failed delete must keep the resource in state — silently dropping it
    // orphans real infrastructure. 404 means it is already gone.
    if httpResp.StatusCode >= 400 && httpResp.StatusCode != http.StatusNotFound {
        err = r.client.ParseResponse(httpResp, nil)
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to delete huntress_connection: %s", err))
        return
    }
    if httpResp.Body != nil {
        httpResp.Body.Close()
    }
}


func (r *HuntressConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}


// keepPlannedValues puts back, after a create or an update, the planned value
// of each optional attribute the configuration leaves out. The server keeps
// some of these up to date on its own (when it last checked a heartbeat, the
// status a probe last reported...), so the value read back after the write can
// already differ from the plan, and Terraform would fail the apply with
// "Provider produced inconsistent result after apply". The next refresh reads
// the server's value, which is never a diff for an attribute nobody configured.
func (r *HuntressConnectionResource) keepPlannedValues(data *HuntressConnectionResourceModel, plan *HuntressConnectionResourceModel, config *HuntressConnectionResourceModel) {
    if config.IsSigningSecretSet.IsNull() && !plan.IsSigningSecretSet.IsUnknown() {
        data.IsSigningSecretSet = plan.IsSigningSecretSet
    }
    if config.OnCallDutyPolicies.IsNull() && !plan.OnCallDutyPolicies.IsUnknown() {
        data.OnCallDutyPolicies = plan.OnCallDutyPolicies
    }
    if config.PageOnCallFor.IsNull() && !plan.PageOnCallFor.IsUnknown() {
        data.PageOnCallFor = plan.PageOnCallFor
    }
    if config.CriticalIncidentSeverityId.IsNull() && !plan.CriticalIncidentSeverityId.IsUnknown() {
        data.CriticalIncidentSeverityId = plan.CriticalIncidentSeverityId
    }
    if config.HighIncidentSeverityId.IsNull() && !plan.HighIncidentSeverityId.IsUnknown() {
        data.HighIncidentSeverityId = plan.HighIncidentSeverityId
    }
    if config.LowIncidentSeverityId.IsNull() && !plan.LowIncidentSeverityId.IsUnknown() {
        data.LowIncidentSeverityId = plan.LowIncidentSeverityId
    }
    if config.WatchedOrganizations.IsNull() && !plan.WatchedOrganizations.IsUnknown() {
        data.WatchedOrganizations = plan.WatchedOrganizations
    }
    if config.Labels.IsNull() && !plan.Labels.IsUnknown() {
        data.Labels = plan.Labels
    }
    if config.ResolveIncidentWhenReportCloses.IsNull() && !plan.ResolveIncidentWhenReportCloses.IsUnknown() {
        data.ResolveIncidentWhenReportCloses = plan.ResolveIncidentWhenReportCloses
    }
}

// Helper method to convert Terraform map to Go interface{}
func (r *HuntressConnectionResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *HuntressConnectionResource) convertTerraformListToInterface(terraformList types.List) interface{} {
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
func (r *HuntressConnectionResource) convertTerraformSetToInterface(terraformSet types.Set) interface{} {
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
func (r *HuntressConnectionResource) parseJSONField(terraformString basetypes.StringValuable) interface{} {
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
func (r *HuntressConnectionResource) normalizeURLWrappers(value interface{}) interface{} {
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

func (r *HuntressConnectionResource) normalizeURLString(value string) string {
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
func (r *HuntressConnectionResource) bigFloatToFloat64(bf *big.Float) interface{} {
    if bf == nil {
        return nil
    }
    f, _ := bf.Float64()
    return f
}

// Helper method to check if a type string is a valid OneUptime ObjectType.
// The registry itself lives in objecttypes.go, shared across the package.
func (r *HuntressConnectionResource) isValidOneUptimeObjectType(typeStr string) bool {
    return validOneUptimeObjectTypes[typeStr]
}
