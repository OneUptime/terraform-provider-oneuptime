package shard21

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
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
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
    "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &BlockResource{}
var _ resource.ResourceWithImportState = &BlockResource{}
var _ resource.ResourceWithModifyPlan = &BlockResource{}

func NewBlockResource() resource.Resource {
    return &BlockResource{}
}

// BlockResource defines the resource implementation.
type BlockResource struct {
    client *Client
}

// BlockResourceModel describes the resource data model.
type BlockResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    Value types.String `tfsdk:"value"`
    Action types.String `tfsdk:"action"`
    Reach types.String `tfsdk:"reach"`
    ExpiresAt RFC3339Value `tfsdk:"expires_at"`
    Reason types.String `tfsdk:"reason"`
    SecurityCaseId types.String `tfsdk:"security_case_id"`
    State types.String `tfsdk:"state"`
    LiftReason types.String `tfsdk:"lift_reason"`
    CreatedAt RFC3339Value `tfsdk:"created_at"`
    UpdatedAt RFC3339Value `tfsdk:"updated_at"`
    CreatedByActorType types.String `tfsdk:"created_by_actor_type"`
    CreatedByActorId types.String `tfsdk:"created_by_actor_id"`
    LastModifiedAt RFC3339Value `tfsdk:"last_modified_at"`
    LastModifiedByActorType types.String `tfsdk:"last_modified_by_actor_type"`
    LastModifiedByActorId types.String `tfsdk:"last_modified_by_actor_id"`
    LastModifiedByUserId types.String `tfsdk:"last_modified_by_user_id"`
    IndicatorType types.String `tfsdk:"indicator_type"`
    Source types.String `tfsdk:"source"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    SecurityPlaybookRunId types.String `tfsdk:"security_playbook_run_id"`
    SecurityAutomaticResponsePolicyId types.String `tfsdk:"security_automatic_response_policy_id"`
    ChallengeFallsBackToBlock types.Bool `tfsdk:"challenge_falls_back_to_block"`
    LiftedAt RFC3339Value `tfsdk:"lifted_at"`
    LiftedByUserId types.String `tfsdk:"lifted_by_user_id"`
}

func (r *BlockResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_block"
}

func (r *BlockResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = r.schemaDefinition()
}

func (r *BlockResource) schemaDefinition() schema.Schema {
    return schema.Schema{
        MarkdownDescription: "The IP addresses, ranges, domains, URLs and files the project blocks: why, until when, and where each block is enforced.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of the project this block belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "value": schema.StringAttribute{
                MarkdownDescription: "What is blocked: an IPv4 or IPv6 address, a range such as 203.0.113.0/24, a domain, an http or https URL, or a file's MD5, SHA-1 or SHA-256 hash. Kept one way: a range as its first address and its length, an IPv6 address compressed in lowercase, a domain lower case without a trailing dot, a URL as a browser sends it, a hash in lower-case hex.",
                Required: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "action": schema.StringAttribute{
                MarkdownDescription: "Block turns it away. Challenge, for an address at the edge, asks for a challenge a person passes and a script or an API client does not. BlockAndQuarantine, for a file, also quarantines the copies an endpoint product finds.",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("Block"),
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                    stringplanmodifier.RequiresReplace(),
                },
                Validators: []validator.String{
                    stringvalidator.OneOf("Block", "Challenge", "BlockAndQuarantine"),
                },
            },
            "reach": schema.StringAttribute{
                MarkdownDescription: "Where it is enforced. Edge, for an address: at the edge in front of the app. Everywhere: on every machine too, through Microsoft Defender or CrowdStrike, which a domain, a URL and a file always are.",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("Edge"),
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                    stringplanmodifier.RequiresReplace(),
                },
                Validators: []validator.String{
                    stringvalidator.OneOf("Edge", "Everywhere"),
                },
            },
            "expires_at": schema.StringAttribute{
                MarkdownDescription: "When the block ends by itself: 24 hours after it is made unless the request says otherwise. Null keeps it until someone lifts it, which takes a Security Admin.",
                CustomType: RFC3339Type{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    RFC3339KeepSameInstant(),
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "reason": schema.StringAttribute{
                MarkdownDescription: "Why the address is blocked, for whoever reads the list later.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "security_case_id": schema.StringAttribute{
                MarkdownDescription: "ID of the case the block was made for. The block and its lifting go on that case's timeline. The ID of a `oneuptime_security_case`.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                    stringplanmodifier.RequiresReplace(),
                },
                Validators: RecordIDValidators("oneuptime_security_case"),
            },
            "state": schema.StringAttribute{
                MarkdownDescription: "Active until it expires or someone lifts it. Set it to Lifted to lift the block; Expired is set by OneUptime when its time is up.",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("Active"),
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    stringvalidator.OneOf("Active", "Lifted", "Expired"),
                },
            },
            "lift_reason": schema.StringAttribute{
                MarkdownDescription: "Why the block was lifted, given with the lift.",
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
            "created_by_actor_type": schema.StringAttribute{
                MarkdownDescription: "What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "created_by_actor_id": schema.StringAttribute{
                MarkdownDescription: "ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "last_modified_at": schema.StringAttribute{
                MarkdownDescription: "When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "last_modified_by_actor_type": schema.StringAttribute{
                MarkdownDescription: "What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.",
                Computed: true,
            },
            "last_modified_by_actor_id": schema.StringAttribute{
                MarkdownDescription: "ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.",
                Computed: true,
            },
            "last_modified_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).",
                Computed: true,
            },
            "indicator_type": schema.StringAttribute{
                MarkdownDescription: "What the block holds: IpAddress, IpRange, Domain, Url, or a file's FileSha256, FileSha1 or FileMd5 hash. Worked out from the value.",
                Computed: true,
            },
            "source": schema.StringAttribute{
                MarkdownDescription: "Who or what made the block: a Person, a Playbook, an AutomaticResponse policy, OneUptimeAi or an Api key.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who blocked it. The ID of a `oneuptime_user` (see the data source).",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "security_playbook_run_id": schema.StringAttribute{
                MarkdownDescription: "ID of the playbook run whose step made the block. The ID of a `oneuptime_security_playbook_run`.",
                Computed: true,
            },
            "security_automatic_response_policy_id": schema.StringAttribute{
                MarkdownDescription: "ID of the automatic response policy that made the block. The ID of a `oneuptime_automatic_response_policy`.",
                Computed: true,
            },
            "challenge_falls_back_to_block": schema.BoolAttribute{
                MarkdownDescription: "For a challenge: whether an edge that cannot challenge holds it as a block for the same time. Set for an automatic response policy's challenges.",
                Computed: true,
            },
            "lifted_at": schema.StringAttribute{
                MarkdownDescription: "When someone lifted the block.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "lifted_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who lifted the block. The ID of a `oneuptime_user` (see the data source).",
                Computed: true,
            },
        },
    }
}

// ModifyPlan keeps the state when nothing a configuration sets changes (see
// importsettle.go).
func (r *BlockResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
    // Destroying: nothing to plan.
    if req.Plan.Raw.IsNull() {
        return
    }

    // Under Crossplane, what a manifest leaves out is the record's (see crossplane.go).
    KeepLeftOutUnderCrossplane(ctx, req, resp)

    KeepStateWhenNothingChanges(ctx, req, resp)
}

func (r *BlockResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *BlockResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data BlockResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // What the configuration sets, and what Terraform planned (see keepPlannedValues).
    var config BlockResourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
    if resp.Diagnostics.HasError() {
        return
    }
    plan := data



    // Create API request body. Unset (null/unknown) optional fields are
    // omitted so server-side defaults apply instead of being overwritten
    // with zero values.
    blockRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := blockRequest["data"].(map[string]interface{})

    if !data.Value.IsNull() && !data.Value.IsUnknown() {
        requestDataMap["value"] = data.Value.ValueString()
    }
    if !data.Action.IsNull() && !data.Action.IsUnknown() {
        requestDataMap["action"] = data.Action.ValueString()
    }
    if !data.Reach.IsNull() && !data.Reach.IsUnknown() {
        requestDataMap["reach"] = data.Reach.ValueString()
    }
    if !data.ExpiresAt.IsNull() && !data.ExpiresAt.IsUnknown() {
        requestDataMap["expiresAt"] = data.ExpiresAt.ValueString()
    }
    if !data.Reason.IsNull() && !data.Reason.IsUnknown() {
        requestDataMap["reason"] = data.Reason.ValueString()
    }
    if !data.SecurityCaseId.IsNull() && !data.SecurityCaseId.IsUnknown() {
        requestDataMap["securityCaseId"] = data.SecurityCaseId.ValueString()
    }

    // Make API call
    httpResp, err := r.client.Post(ctx, "/security-block", blockRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create block, got error: %s", err))
        return
    }

    var blockResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &blockResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to create block: %s", err))
        return
    }

    // What OneUptime left out because only it writes it, if anything (see entries.go).
    resp.Diagnostics.Append(FieldsLeftOutWarning(httpResp, "block")...)

    // What the write does that someone should know, if anything (see entries.go).
    resp.Diagnostics.Append(WriteWarnings(httpResp, "block")...)

    // Extract the new resource id from the create response.
    createdId := ""
    if wrapper, ok := blockResponse["data"].(map[string]interface{}); ok {
        if val, ok := wrapper["_id"].(string); ok {
            createdId = val
        }
    } else if val, ok := blockResponse["_id"].(string); ok {
        createdId = val
    }
    if createdId == "" {
        resp.Diagnostics.AddError("OneUptime API Error", "Create response for block did not contain an id. This is a bug in the provider or the API; please report it.")
        return
    }
    data.Id = types.StringValue(createdId)


    /*
     * The server has committed the row. Persist what we know to state BEFORE
     * the read-back: if the read-back fails and we return without setting
     * state, Terraform never learns the resource exists and the created
     * block is orphaned server-side — never refreshed, never
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
        "value": true,
        "action": true,
        "reach": true,
        "expiresAt": true,
        "reason": true,
        "securityCaseId": true,
        "state": true,
        "liftReason": true,
        "createdAt": true,
        "updatedAt": true,
        "createdByActorType": true,
        "createdByActorId": true,
        "lastModifiedAt": true,
        "lastModifiedByActorType": true,
        "lastModifiedByActorId": true,
        "lastModifiedByUserId": true,
        "indicatorType": true,
        "source": true,
        "createdByUserId": true,
        "securityPlaybookRunId": true,
        "securityAutomaticResponsePolicyId": true,
        "challengeFallsBackToBlock": true,
        "liftedAt": true,
        "liftedByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/security-block/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        /*
         * State already owns the id, so the resource is tracked and the next
         * refresh reconciles the remaining attributes. Warn rather than
         * error: erroring here would strand a real resource.
         */
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created block but could not read it back; state is incomplete until the next refresh: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created block but could not parse the read-back response; state is incomplete until the next refresh: %s", err))
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
    if obj, ok := dataMap["value"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Value = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Value = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Value = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Value = types.StringValue(string(jsonBytes))
            } else {
                data.Value = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Value = types.StringValue(string(jsonBytes))
            } else {
                data.Value = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Value = types.StringValue(string(jsonBytes))
        } else {
            data.Value = types.StringNull()
        }
    } else if arr, ok := dataMap["value"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Value = types.StringValue(string(jsonBytes))
        } else {
            data.Value = types.StringNull()
        }
    } else if val, ok := dataMap["value"].(string); ok {
        data.Value = types.StringValue(val)
    } else {
        data.Value = types.StringNull()
    }
    if obj, ok := dataMap["action"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Action = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Action = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Action = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Action = types.StringValue(string(jsonBytes))
            } else {
                data.Action = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Action = types.StringValue(string(jsonBytes))
            } else {
                data.Action = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Action = types.StringValue(string(jsonBytes))
        } else {
            data.Action = types.StringNull()
        }
    } else if arr, ok := dataMap["action"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Action = types.StringValue(string(jsonBytes))
        } else {
            data.Action = types.StringNull()
        }
    } else if val, ok := dataMap["action"].(string); ok {
        data.Action = types.StringValue(val)
    } else {
        data.Action = types.StringNull()
    }
    if obj, ok := dataMap["reach"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Reach = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Reach = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Reach = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Reach = types.StringValue(string(jsonBytes))
            } else {
                data.Reach = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Reach = types.StringValue(string(jsonBytes))
            } else {
                data.Reach = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Reach = types.StringValue(string(jsonBytes))
        } else {
            data.Reach = types.StringNull()
        }
    } else if arr, ok := dataMap["reach"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Reach = types.StringValue(string(jsonBytes))
        } else {
            data.Reach = types.StringNull()
        }
    } else if val, ok := dataMap["reach"].(string); ok {
        data.Reach = types.StringValue(val)
    } else {
        data.Reach = types.StringNull()
    }
    if obj, ok := dataMap["expiresAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ExpiresAt = NewRFC3339Value(val)
        } else {
            data.ExpiresAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["expiresAt"].(string); ok && val != "" {
        data.ExpiresAt = NewRFC3339Value(val)
    } else {
        data.ExpiresAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["reason"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Reason = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Reason = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Reason = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Reason = types.StringValue(string(jsonBytes))
            } else {
                data.Reason = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Reason = types.StringValue(string(jsonBytes))
            } else {
                data.Reason = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Reason = types.StringValue(string(jsonBytes))
        } else {
            data.Reason = types.StringNull()
        }
    } else if arr, ok := dataMap["reason"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Reason = types.StringValue(string(jsonBytes))
        } else {
            data.Reason = types.StringNull()
        }
    } else if val, ok := dataMap["reason"].(string); ok {
        data.Reason = types.StringValue(val)
    } else {
        data.Reason = types.StringNull()
    }
    if obj, ok := dataMap["securityCaseId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SecurityCaseId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SecurityCaseId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SecurityCaseId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SecurityCaseId = types.StringValue(string(jsonBytes))
            } else {
                data.SecurityCaseId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SecurityCaseId = types.StringValue(string(jsonBytes))
            } else {
                data.SecurityCaseId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SecurityCaseId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityCaseId = types.StringNull()
        }
    } else if arr, ok := dataMap["securityCaseId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.SecurityCaseId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityCaseId = types.StringNull()
        }
    } else if val, ok := dataMap["securityCaseId"].(string); ok {
        data.SecurityCaseId = types.StringValue(val)
    } else {
        data.SecurityCaseId = types.StringNull()
    }
    if obj, ok := dataMap["state"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.State = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.State = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.State = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.State = types.StringValue(string(jsonBytes))
            } else {
                data.State = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.State = types.StringValue(string(jsonBytes))
            } else {
                data.State = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.State = types.StringValue(string(jsonBytes))
        } else {
            data.State = types.StringNull()
        }
    } else if arr, ok := dataMap["state"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.State = types.StringValue(string(jsonBytes))
        } else {
            data.State = types.StringNull()
        }
    } else if val, ok := dataMap["state"].(string); ok {
        data.State = types.StringValue(val)
    } else {
        data.State = types.StringNull()
    }
    if obj, ok := dataMap["liftReason"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LiftReason = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LiftReason = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LiftReason = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LiftReason = types.StringValue(string(jsonBytes))
            } else {
                data.LiftReason = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LiftReason = types.StringValue(string(jsonBytes))
            } else {
                data.LiftReason = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LiftReason = types.StringValue(string(jsonBytes))
        } else {
            data.LiftReason = types.StringNull()
        }
    } else if arr, ok := dataMap["liftReason"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LiftReason = types.StringValue(string(jsonBytes))
        } else {
            data.LiftReason = types.StringNull()
        }
    } else if val, ok := dataMap["liftReason"].(string); ok {
        data.LiftReason = types.StringValue(val)
    } else {
        data.LiftReason = types.StringNull()
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
    if obj, ok := dataMap["createdByActorType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByActorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByActorType = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByActorType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByActorType = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByActorType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorType = types.StringNull()
        }
    } else if arr, ok := dataMap["createdByActorType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CreatedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorType = types.StringNull()
        }
    } else if val, ok := dataMap["createdByActorType"].(string); ok {
        data.CreatedByActorType = types.StringValue(val)
    } else {
        data.CreatedByActorType = types.StringNull()
    }
    if obj, ok := dataMap["createdByActorId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByActorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByActorId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByActorId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByActorId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByActorId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorId = types.StringNull()
        }
    } else if arr, ok := dataMap["createdByActorId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CreatedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByActorId"].(string); ok {
        data.CreatedByActorId = types.StringValue(val)
    } else {
        data.CreatedByActorId = types.StringNull()
    }
    if obj, ok := dataMap["lastModifiedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastModifiedAt = NewRFC3339Value(val)
        } else {
            data.LastModifiedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastModifiedAt"].(string); ok && val != "" {
        data.LastModifiedAt = NewRFC3339Value(val)
    } else {
        data.LastModifiedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastModifiedByActorType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastModifiedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastModifiedByActorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByActorType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByActorType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorType = types.StringNull()
        }
    } else if arr, ok := dataMap["lastModifiedByActorType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorType = types.StringNull()
        }
    } else if val, ok := dataMap["lastModifiedByActorType"].(string); ok {
        data.LastModifiedByActorType = types.StringValue(val)
    } else {
        data.LastModifiedByActorType = types.StringNull()
    }
    if obj, ok := dataMap["lastModifiedByActorId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastModifiedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastModifiedByActorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByActorId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByActorId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorId = types.StringNull()
        }
    } else if arr, ok := dataMap["lastModifiedByActorId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorId = types.StringNull()
        }
    } else if val, ok := dataMap["lastModifiedByActorId"].(string); ok {
        data.LastModifiedByActorId = types.StringValue(val)
    } else {
        data.LastModifiedByActorId = types.StringNull()
    }
    if obj, ok := dataMap["lastModifiedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastModifiedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastModifiedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["lastModifiedByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["lastModifiedByUserId"].(string); ok {
        data.LastModifiedByUserId = types.StringValue(val)
    } else {
        data.LastModifiedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["indicatorType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IndicatorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IndicatorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IndicatorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.IndicatorType = types.StringValue(string(jsonBytes))
            } else {
                data.IndicatorType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.IndicatorType = types.StringValue(string(jsonBytes))
            } else {
                data.IndicatorType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IndicatorType = types.StringValue(string(jsonBytes))
        } else {
            data.IndicatorType = types.StringNull()
        }
    } else if arr, ok := dataMap["indicatorType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.IndicatorType = types.StringValue(string(jsonBytes))
        } else {
            data.IndicatorType = types.StringNull()
        }
    } else if val, ok := dataMap["indicatorType"].(string); ok {
        data.IndicatorType = types.StringValue(val)
    } else {
        data.IndicatorType = types.StringNull()
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
    } else if arr, ok := dataMap["source"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Source = types.StringValue(string(jsonBytes))
        } else {
            data.Source = types.StringNull()
        }
    } else if val, ok := dataMap["source"].(string); ok {
        data.Source = types.StringValue(val)
    } else {
        data.Source = types.StringNull()
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
    if obj, ok := dataMap["securityPlaybookRunId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SecurityPlaybookRunId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SecurityPlaybookRunId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SecurityPlaybookRunId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SecurityPlaybookRunId = types.StringValue(string(jsonBytes))
            } else {
                data.SecurityPlaybookRunId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SecurityPlaybookRunId = types.StringValue(string(jsonBytes))
            } else {
                data.SecurityPlaybookRunId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SecurityPlaybookRunId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityPlaybookRunId = types.StringNull()
        }
    } else if arr, ok := dataMap["securityPlaybookRunId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.SecurityPlaybookRunId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityPlaybookRunId = types.StringNull()
        }
    } else if val, ok := dataMap["securityPlaybookRunId"].(string); ok {
        data.SecurityPlaybookRunId = types.StringValue(val)
    } else {
        data.SecurityPlaybookRunId = types.StringNull()
    }
    if obj, ok := dataMap["securityAutomaticResponsePolicyId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SecurityAutomaticResponsePolicyId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SecurityAutomaticResponsePolicyId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SecurityAutomaticResponsePolicyId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SecurityAutomaticResponsePolicyId = types.StringValue(string(jsonBytes))
            } else {
                data.SecurityAutomaticResponsePolicyId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SecurityAutomaticResponsePolicyId = types.StringValue(string(jsonBytes))
            } else {
                data.SecurityAutomaticResponsePolicyId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SecurityAutomaticResponsePolicyId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityAutomaticResponsePolicyId = types.StringNull()
        }
    } else if arr, ok := dataMap["securityAutomaticResponsePolicyId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.SecurityAutomaticResponsePolicyId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityAutomaticResponsePolicyId = types.StringNull()
        }
    } else if val, ok := dataMap["securityAutomaticResponsePolicyId"].(string); ok {
        data.SecurityAutomaticResponsePolicyId = types.StringValue(val)
    } else {
        data.SecurityAutomaticResponsePolicyId = types.StringNull()
    }
    if val, ok := dataMap["challengeFallsBackToBlock"].(bool); ok {
        data.ChallengeFallsBackToBlock = types.BoolValue(val)
    }
    if obj, ok := dataMap["liftedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LiftedAt = NewRFC3339Value(val)
        } else {
            data.LiftedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["liftedAt"].(string); ok && val != "" {
        data.LiftedAt = NewRFC3339Value(val)
    } else {
        data.LiftedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["liftedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LiftedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LiftedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LiftedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LiftedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.LiftedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LiftedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.LiftedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LiftedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LiftedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["liftedByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LiftedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LiftedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["liftedByUserId"].(string); ok {
        data.LiftedByUserId = types.StringValue(val)
    } else {
        data.LiftedByUserId = types.StringNull()
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

func (r *BlockResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data BlockResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "value": true,
        "action": true,
        "reach": true,
        "expiresAt": true,
        "reason": true,
        "securityCaseId": true,
        "state": true,
        "liftReason": true,
        "createdAt": true,
        "updatedAt": true,
        "createdByActorType": true,
        "createdByActorId": true,
        "lastModifiedAt": true,
        "lastModifiedByActorType": true,
        "lastModifiedByActorId": true,
        "lastModifiedByUserId": true,
        "indicatorType": true,
        "source": true,
        "createdByUserId": true,
        "securityPlaybookRunId": true,
        "securityAutomaticResponsePolicyId": true,
        "challengeFallsBackToBlock": true,
        "liftedAt": true,
        "liftedByUserId": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect(ctx, "/security-block/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read block, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var blockResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &blockResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse block response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := blockResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = blockResponse
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
    if obj, ok := dataMap["value"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Value = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Value = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Value = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Value = types.StringValue(string(jsonBytes))
            } else {
                data.Value = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Value = types.StringValue(string(jsonBytes))
            } else {
                data.Value = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Value = types.StringValue(string(jsonBytes))
        } else {
            data.Value = types.StringNull()
        }
    } else if arr, ok := dataMap["value"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Value = types.StringValue(string(jsonBytes))
        } else {
            data.Value = types.StringNull()
        }
    } else if val, ok := dataMap["value"].(string); ok {
        data.Value = types.StringValue(val)
    } else {
        data.Value = types.StringNull()
    }
    if obj, ok := dataMap["action"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Action = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Action = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Action = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Action = types.StringValue(string(jsonBytes))
            } else {
                data.Action = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Action = types.StringValue(string(jsonBytes))
            } else {
                data.Action = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Action = types.StringValue(string(jsonBytes))
        } else {
            data.Action = types.StringNull()
        }
    } else if arr, ok := dataMap["action"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Action = types.StringValue(string(jsonBytes))
        } else {
            data.Action = types.StringNull()
        }
    } else if val, ok := dataMap["action"].(string); ok {
        data.Action = types.StringValue(val)
    } else {
        data.Action = types.StringNull()
    }
    if obj, ok := dataMap["reach"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Reach = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Reach = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Reach = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Reach = types.StringValue(string(jsonBytes))
            } else {
                data.Reach = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Reach = types.StringValue(string(jsonBytes))
            } else {
                data.Reach = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Reach = types.StringValue(string(jsonBytes))
        } else {
            data.Reach = types.StringNull()
        }
    } else if arr, ok := dataMap["reach"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Reach = types.StringValue(string(jsonBytes))
        } else {
            data.Reach = types.StringNull()
        }
    } else if val, ok := dataMap["reach"].(string); ok {
        data.Reach = types.StringValue(val)
    } else {
        data.Reach = types.StringNull()
    }
    if obj, ok := dataMap["expiresAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ExpiresAt = NewRFC3339Value(val)
        } else {
            data.ExpiresAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["expiresAt"].(string); ok && val != "" {
        data.ExpiresAt = NewRFC3339Value(val)
    } else {
        data.ExpiresAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["reason"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Reason = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Reason = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Reason = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Reason = types.StringValue(string(jsonBytes))
            } else {
                data.Reason = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Reason = types.StringValue(string(jsonBytes))
            } else {
                data.Reason = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Reason = types.StringValue(string(jsonBytes))
        } else {
            data.Reason = types.StringNull()
        }
    } else if arr, ok := dataMap["reason"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Reason = types.StringValue(string(jsonBytes))
        } else {
            data.Reason = types.StringNull()
        }
    } else if val, ok := dataMap["reason"].(string); ok {
        data.Reason = types.StringValue(val)
    } else {
        data.Reason = types.StringNull()
    }
    if obj, ok := dataMap["securityCaseId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SecurityCaseId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SecurityCaseId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SecurityCaseId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SecurityCaseId = types.StringValue(string(jsonBytes))
            } else {
                data.SecurityCaseId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SecurityCaseId = types.StringValue(string(jsonBytes))
            } else {
                data.SecurityCaseId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SecurityCaseId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityCaseId = types.StringNull()
        }
    } else if arr, ok := dataMap["securityCaseId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.SecurityCaseId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityCaseId = types.StringNull()
        }
    } else if val, ok := dataMap["securityCaseId"].(string); ok {
        data.SecurityCaseId = types.StringValue(val)
    } else {
        data.SecurityCaseId = types.StringNull()
    }
    if obj, ok := dataMap["state"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.State = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.State = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.State = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.State = types.StringValue(string(jsonBytes))
            } else {
                data.State = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.State = types.StringValue(string(jsonBytes))
            } else {
                data.State = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.State = types.StringValue(string(jsonBytes))
        } else {
            data.State = types.StringNull()
        }
    } else if arr, ok := dataMap["state"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.State = types.StringValue(string(jsonBytes))
        } else {
            data.State = types.StringNull()
        }
    } else if val, ok := dataMap["state"].(string); ok {
        data.State = types.StringValue(val)
    } else {
        data.State = types.StringNull()
    }
    if obj, ok := dataMap["liftReason"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LiftReason = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LiftReason = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LiftReason = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LiftReason = types.StringValue(string(jsonBytes))
            } else {
                data.LiftReason = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LiftReason = types.StringValue(string(jsonBytes))
            } else {
                data.LiftReason = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LiftReason = types.StringValue(string(jsonBytes))
        } else {
            data.LiftReason = types.StringNull()
        }
    } else if arr, ok := dataMap["liftReason"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LiftReason = types.StringValue(string(jsonBytes))
        } else {
            data.LiftReason = types.StringNull()
        }
    } else if val, ok := dataMap["liftReason"].(string); ok {
        data.LiftReason = types.StringValue(val)
    } else {
        data.LiftReason = types.StringNull()
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
    if obj, ok := dataMap["createdByActorType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByActorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByActorType = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByActorType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByActorType = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByActorType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorType = types.StringNull()
        }
    } else if arr, ok := dataMap["createdByActorType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CreatedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorType = types.StringNull()
        }
    } else if val, ok := dataMap["createdByActorType"].(string); ok {
        data.CreatedByActorType = types.StringValue(val)
    } else {
        data.CreatedByActorType = types.StringNull()
    }
    if obj, ok := dataMap["createdByActorId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByActorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByActorId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByActorId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByActorId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByActorId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorId = types.StringNull()
        }
    } else if arr, ok := dataMap["createdByActorId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CreatedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByActorId"].(string); ok {
        data.CreatedByActorId = types.StringValue(val)
    } else {
        data.CreatedByActorId = types.StringNull()
    }
    if obj, ok := dataMap["lastModifiedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastModifiedAt = NewRFC3339Value(val)
        } else {
            data.LastModifiedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastModifiedAt"].(string); ok && val != "" {
        data.LastModifiedAt = NewRFC3339Value(val)
    } else {
        data.LastModifiedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastModifiedByActorType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastModifiedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastModifiedByActorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByActorType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByActorType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorType = types.StringNull()
        }
    } else if arr, ok := dataMap["lastModifiedByActorType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorType = types.StringNull()
        }
    } else if val, ok := dataMap["lastModifiedByActorType"].(string); ok {
        data.LastModifiedByActorType = types.StringValue(val)
    } else {
        data.LastModifiedByActorType = types.StringNull()
    }
    if obj, ok := dataMap["lastModifiedByActorId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastModifiedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastModifiedByActorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByActorId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByActorId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorId = types.StringNull()
        }
    } else if arr, ok := dataMap["lastModifiedByActorId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorId = types.StringNull()
        }
    } else if val, ok := dataMap["lastModifiedByActorId"].(string); ok {
        data.LastModifiedByActorId = types.StringValue(val)
    } else {
        data.LastModifiedByActorId = types.StringNull()
    }
    if obj, ok := dataMap["lastModifiedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastModifiedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastModifiedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["lastModifiedByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["lastModifiedByUserId"].(string); ok {
        data.LastModifiedByUserId = types.StringValue(val)
    } else {
        data.LastModifiedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["indicatorType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IndicatorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IndicatorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IndicatorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.IndicatorType = types.StringValue(string(jsonBytes))
            } else {
                data.IndicatorType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.IndicatorType = types.StringValue(string(jsonBytes))
            } else {
                data.IndicatorType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IndicatorType = types.StringValue(string(jsonBytes))
        } else {
            data.IndicatorType = types.StringNull()
        }
    } else if arr, ok := dataMap["indicatorType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.IndicatorType = types.StringValue(string(jsonBytes))
        } else {
            data.IndicatorType = types.StringNull()
        }
    } else if val, ok := dataMap["indicatorType"].(string); ok {
        data.IndicatorType = types.StringValue(val)
    } else {
        data.IndicatorType = types.StringNull()
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
    } else if arr, ok := dataMap["source"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Source = types.StringValue(string(jsonBytes))
        } else {
            data.Source = types.StringNull()
        }
    } else if val, ok := dataMap["source"].(string); ok {
        data.Source = types.StringValue(val)
    } else {
        data.Source = types.StringNull()
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
    if obj, ok := dataMap["securityPlaybookRunId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SecurityPlaybookRunId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SecurityPlaybookRunId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SecurityPlaybookRunId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SecurityPlaybookRunId = types.StringValue(string(jsonBytes))
            } else {
                data.SecurityPlaybookRunId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SecurityPlaybookRunId = types.StringValue(string(jsonBytes))
            } else {
                data.SecurityPlaybookRunId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SecurityPlaybookRunId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityPlaybookRunId = types.StringNull()
        }
    } else if arr, ok := dataMap["securityPlaybookRunId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.SecurityPlaybookRunId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityPlaybookRunId = types.StringNull()
        }
    } else if val, ok := dataMap["securityPlaybookRunId"].(string); ok {
        data.SecurityPlaybookRunId = types.StringValue(val)
    } else {
        data.SecurityPlaybookRunId = types.StringNull()
    }
    if obj, ok := dataMap["securityAutomaticResponsePolicyId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SecurityAutomaticResponsePolicyId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SecurityAutomaticResponsePolicyId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SecurityAutomaticResponsePolicyId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SecurityAutomaticResponsePolicyId = types.StringValue(string(jsonBytes))
            } else {
                data.SecurityAutomaticResponsePolicyId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SecurityAutomaticResponsePolicyId = types.StringValue(string(jsonBytes))
            } else {
                data.SecurityAutomaticResponsePolicyId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SecurityAutomaticResponsePolicyId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityAutomaticResponsePolicyId = types.StringNull()
        }
    } else if arr, ok := dataMap["securityAutomaticResponsePolicyId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.SecurityAutomaticResponsePolicyId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityAutomaticResponsePolicyId = types.StringNull()
        }
    } else if val, ok := dataMap["securityAutomaticResponsePolicyId"].(string); ok {
        data.SecurityAutomaticResponsePolicyId = types.StringValue(val)
    } else {
        data.SecurityAutomaticResponsePolicyId = types.StringNull()
    }
    if val, ok := dataMap["challengeFallsBackToBlock"].(bool); ok {
        data.ChallengeFallsBackToBlock = types.BoolValue(val)
    }
    if obj, ok := dataMap["liftedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LiftedAt = NewRFC3339Value(val)
        } else {
            data.LiftedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["liftedAt"].(string); ok && val != "" {
        data.LiftedAt = NewRFC3339Value(val)
    } else {
        data.LiftedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["liftedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LiftedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LiftedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LiftedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LiftedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.LiftedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LiftedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.LiftedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LiftedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LiftedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["liftedByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LiftedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LiftedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["liftedByUserId"].(string); ok {
        data.LiftedByUserId = types.StringValue(val)
    } else {
        data.LiftedByUserId = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BlockResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data BlockResourceModel
    var state BlockResourceModel

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
    var config BlockResourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
    if resp.Diagnostics.HasError() {
        return
    }
    plan := data

    // Create API request body
    blockRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := blockRequest["data"].(map[string]interface{})

    if !data.State.IsUnknown() && !state.State.IsUnknown() && !data.State.Equal(state.State) {
        requestDataMap["state"] = data.State.ValueString()
    }
    if !data.ExpiresAt.IsUnknown() && !state.ExpiresAt.IsUnknown() && !data.ExpiresAt.Equal(state.ExpiresAt) {
        requestDataMap["expiresAt"] = data.ExpiresAt.ValueString()
    }
    if !data.Reason.IsUnknown() && !state.Reason.IsUnknown() && !data.Reason.Equal(state.Reason) {
        requestDataMap["reason"] = data.Reason.ValueString()
    }
    if !data.LiftReason.IsUnknown() && !state.LiftReason.IsUnknown() && !data.LiftReason.Equal(state.LiftReason) {
        requestDataMap["liftReason"] = data.LiftReason.ValueString()
    }

    // Only call the API when there are changed fields to send. An empty
    // update body is rejected by the API; state is still refreshed below so
    // this method never writes unverified plan values into state.
    if len(blockRequest["data"].(map[string]interface{})) > 0 {
        httpResp, err := r.client.Put(ctx, "/security-block/" + data.Id.ValueString() + "", blockRequest)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update block, got error: %s", err))
            return
        }

        // Parse the update response
        var blockResponse map[string]interface{}
        err = r.client.ParseResponse(httpResp, &blockResponse)
        if err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to update block: %s", err))
            return
        }
        _ = blockResponse

        // What the write does that someone should know, if anything (see entries.go).
        resp.Diagnostics.Append(WriteWarnings(httpResp, "block")...)
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "value": true,
        "action": true,
        "reach": true,
        "expiresAt": true,
        "reason": true,
        "securityCaseId": true,
        "state": true,
        "liftReason": true,
        "createdAt": true,
        "updatedAt": true,
        "createdByActorType": true,
        "createdByActorId": true,
        "lastModifiedAt": true,
        "lastModifiedByActorType": true,
        "lastModifiedByActorId": true,
        "lastModifiedByUserId": true,
        "indicatorType": true,
        "source": true,
        "createdByUserId": true,
        "securityPlaybookRunId": true,
        "securityAutomaticResponsePolicyId": true,
        "challengeFallsBackToBlock": true,
        "liftedAt": true,
        "liftedByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/security-block/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read block after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read block after update: %s", err))
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
    if obj, ok := dataMap["value"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Value = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Value = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Value = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Value = types.StringValue(string(jsonBytes))
            } else {
                data.Value = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Value = types.StringValue(string(jsonBytes))
            } else {
                data.Value = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Value = types.StringValue(string(jsonBytes))
        } else {
            data.Value = types.StringNull()
        }
    } else if arr, ok := dataMap["value"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Value = types.StringValue(string(jsonBytes))
        } else {
            data.Value = types.StringNull()
        }
    } else if val, ok := dataMap["value"].(string); ok {
        data.Value = types.StringValue(val)
    } else {
        data.Value = types.StringNull()
    }
    if obj, ok := dataMap["action"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Action = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Action = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Action = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Action = types.StringValue(string(jsonBytes))
            } else {
                data.Action = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Action = types.StringValue(string(jsonBytes))
            } else {
                data.Action = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Action = types.StringValue(string(jsonBytes))
        } else {
            data.Action = types.StringNull()
        }
    } else if arr, ok := dataMap["action"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Action = types.StringValue(string(jsonBytes))
        } else {
            data.Action = types.StringNull()
        }
    } else if val, ok := dataMap["action"].(string); ok {
        data.Action = types.StringValue(val)
    } else {
        data.Action = types.StringNull()
    }
    if obj, ok := dataMap["reach"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Reach = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Reach = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Reach = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Reach = types.StringValue(string(jsonBytes))
            } else {
                data.Reach = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Reach = types.StringValue(string(jsonBytes))
            } else {
                data.Reach = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Reach = types.StringValue(string(jsonBytes))
        } else {
            data.Reach = types.StringNull()
        }
    } else if arr, ok := dataMap["reach"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Reach = types.StringValue(string(jsonBytes))
        } else {
            data.Reach = types.StringNull()
        }
    } else if val, ok := dataMap["reach"].(string); ok {
        data.Reach = types.StringValue(val)
    } else {
        data.Reach = types.StringNull()
    }
    if obj, ok := dataMap["expiresAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ExpiresAt = NewRFC3339Value(val)
        } else {
            data.ExpiresAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["expiresAt"].(string); ok && val != "" {
        data.ExpiresAt = NewRFC3339Value(val)
    } else {
        data.ExpiresAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["reason"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Reason = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Reason = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Reason = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Reason = types.StringValue(string(jsonBytes))
            } else {
                data.Reason = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Reason = types.StringValue(string(jsonBytes))
            } else {
                data.Reason = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Reason = types.StringValue(string(jsonBytes))
        } else {
            data.Reason = types.StringNull()
        }
    } else if arr, ok := dataMap["reason"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Reason = types.StringValue(string(jsonBytes))
        } else {
            data.Reason = types.StringNull()
        }
    } else if val, ok := dataMap["reason"].(string); ok {
        data.Reason = types.StringValue(val)
    } else {
        data.Reason = types.StringNull()
    }
    if obj, ok := dataMap["securityCaseId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SecurityCaseId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SecurityCaseId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SecurityCaseId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SecurityCaseId = types.StringValue(string(jsonBytes))
            } else {
                data.SecurityCaseId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SecurityCaseId = types.StringValue(string(jsonBytes))
            } else {
                data.SecurityCaseId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SecurityCaseId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityCaseId = types.StringNull()
        }
    } else if arr, ok := dataMap["securityCaseId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.SecurityCaseId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityCaseId = types.StringNull()
        }
    } else if val, ok := dataMap["securityCaseId"].(string); ok {
        data.SecurityCaseId = types.StringValue(val)
    } else {
        data.SecurityCaseId = types.StringNull()
    }
    if obj, ok := dataMap["state"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.State = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.State = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.State = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.State = types.StringValue(string(jsonBytes))
            } else {
                data.State = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.State = types.StringValue(string(jsonBytes))
            } else {
                data.State = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.State = types.StringValue(string(jsonBytes))
        } else {
            data.State = types.StringNull()
        }
    } else if arr, ok := dataMap["state"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.State = types.StringValue(string(jsonBytes))
        } else {
            data.State = types.StringNull()
        }
    } else if val, ok := dataMap["state"].(string); ok {
        data.State = types.StringValue(val)
    } else {
        data.State = types.StringNull()
    }
    if obj, ok := dataMap["liftReason"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LiftReason = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LiftReason = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LiftReason = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LiftReason = types.StringValue(string(jsonBytes))
            } else {
                data.LiftReason = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LiftReason = types.StringValue(string(jsonBytes))
            } else {
                data.LiftReason = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LiftReason = types.StringValue(string(jsonBytes))
        } else {
            data.LiftReason = types.StringNull()
        }
    } else if arr, ok := dataMap["liftReason"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LiftReason = types.StringValue(string(jsonBytes))
        } else {
            data.LiftReason = types.StringNull()
        }
    } else if val, ok := dataMap["liftReason"].(string); ok {
        data.LiftReason = types.StringValue(val)
    } else {
        data.LiftReason = types.StringNull()
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
    if obj, ok := dataMap["createdByActorType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByActorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByActorType = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByActorType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByActorType = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByActorType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorType = types.StringNull()
        }
    } else if arr, ok := dataMap["createdByActorType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CreatedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorType = types.StringNull()
        }
    } else if val, ok := dataMap["createdByActorType"].(string); ok {
        data.CreatedByActorType = types.StringValue(val)
    } else {
        data.CreatedByActorType = types.StringNull()
    }
    if obj, ok := dataMap["createdByActorId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByActorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByActorId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByActorId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByActorId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByActorId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorId = types.StringNull()
        }
    } else if arr, ok := dataMap["createdByActorId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.CreatedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByActorId"].(string); ok {
        data.CreatedByActorId = types.StringValue(val)
    } else {
        data.CreatedByActorId = types.StringNull()
    }
    if obj, ok := dataMap["lastModifiedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastModifiedAt = NewRFC3339Value(val)
        } else {
            data.LastModifiedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastModifiedAt"].(string); ok && val != "" {
        data.LastModifiedAt = NewRFC3339Value(val)
    } else {
        data.LastModifiedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastModifiedByActorType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastModifiedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastModifiedByActorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByActorType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByActorType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorType = types.StringNull()
        }
    } else if arr, ok := dataMap["lastModifiedByActorType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorType = types.StringNull()
        }
    } else if val, ok := dataMap["lastModifiedByActorType"].(string); ok {
        data.LastModifiedByActorType = types.StringValue(val)
    } else {
        data.LastModifiedByActorType = types.StringNull()
    }
    if obj, ok := dataMap["lastModifiedByActorId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastModifiedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastModifiedByActorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByActorId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByActorId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorId = types.StringNull()
        }
    } else if arr, ok := dataMap["lastModifiedByActorId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorId = types.StringNull()
        }
    } else if val, ok := dataMap["lastModifiedByActorId"].(string); ok {
        data.LastModifiedByActorId = types.StringValue(val)
    } else {
        data.LastModifiedByActorId = types.StringNull()
    }
    if obj, ok := dataMap["lastModifiedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastModifiedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastModifiedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.LastModifiedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["lastModifiedByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["lastModifiedByUserId"].(string); ok {
        data.LastModifiedByUserId = types.StringValue(val)
    } else {
        data.LastModifiedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["indicatorType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IndicatorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IndicatorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IndicatorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.IndicatorType = types.StringValue(string(jsonBytes))
            } else {
                data.IndicatorType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.IndicatorType = types.StringValue(string(jsonBytes))
            } else {
                data.IndicatorType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IndicatorType = types.StringValue(string(jsonBytes))
        } else {
            data.IndicatorType = types.StringNull()
        }
    } else if arr, ok := dataMap["indicatorType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.IndicatorType = types.StringValue(string(jsonBytes))
        } else {
            data.IndicatorType = types.StringNull()
        }
    } else if val, ok := dataMap["indicatorType"].(string); ok {
        data.IndicatorType = types.StringValue(val)
    } else {
        data.IndicatorType = types.StringNull()
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
    } else if arr, ok := dataMap["source"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Source = types.StringValue(string(jsonBytes))
        } else {
            data.Source = types.StringNull()
        }
    } else if val, ok := dataMap["source"].(string); ok {
        data.Source = types.StringValue(val)
    } else {
        data.Source = types.StringNull()
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
    if obj, ok := dataMap["securityPlaybookRunId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SecurityPlaybookRunId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SecurityPlaybookRunId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SecurityPlaybookRunId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SecurityPlaybookRunId = types.StringValue(string(jsonBytes))
            } else {
                data.SecurityPlaybookRunId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SecurityPlaybookRunId = types.StringValue(string(jsonBytes))
            } else {
                data.SecurityPlaybookRunId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SecurityPlaybookRunId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityPlaybookRunId = types.StringNull()
        }
    } else if arr, ok := dataMap["securityPlaybookRunId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.SecurityPlaybookRunId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityPlaybookRunId = types.StringNull()
        }
    } else if val, ok := dataMap["securityPlaybookRunId"].(string); ok {
        data.SecurityPlaybookRunId = types.StringValue(val)
    } else {
        data.SecurityPlaybookRunId = types.StringNull()
    }
    if obj, ok := dataMap["securityAutomaticResponsePolicyId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SecurityAutomaticResponsePolicyId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SecurityAutomaticResponsePolicyId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SecurityAutomaticResponsePolicyId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SecurityAutomaticResponsePolicyId = types.StringValue(string(jsonBytes))
            } else {
                data.SecurityAutomaticResponsePolicyId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SecurityAutomaticResponsePolicyId = types.StringValue(string(jsonBytes))
            } else {
                data.SecurityAutomaticResponsePolicyId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SecurityAutomaticResponsePolicyId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityAutomaticResponsePolicyId = types.StringNull()
        }
    } else if arr, ok := dataMap["securityAutomaticResponsePolicyId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.SecurityAutomaticResponsePolicyId = types.StringValue(string(jsonBytes))
        } else {
            data.SecurityAutomaticResponsePolicyId = types.StringNull()
        }
    } else if val, ok := dataMap["securityAutomaticResponsePolicyId"].(string); ok {
        data.SecurityAutomaticResponsePolicyId = types.StringValue(val)
    } else {
        data.SecurityAutomaticResponsePolicyId = types.StringNull()
    }
    if val, ok := dataMap["challengeFallsBackToBlock"].(bool); ok {
        data.ChallengeFallsBackToBlock = types.BoolValue(val)
    }
    if obj, ok := dataMap["liftedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LiftedAt = NewRFC3339Value(val)
        } else {
            data.LiftedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["liftedAt"].(string); ok && val != "" {
        data.LiftedAt = NewRFC3339Value(val)
    } else {
        data.LiftedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["liftedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LiftedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LiftedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LiftedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LiftedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.LiftedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LiftedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.LiftedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LiftedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LiftedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["liftedByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LiftedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LiftedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["liftedByUserId"].(string); ok {
        data.LiftedByUserId = types.StringValue(val)
    } else {
        data.LiftedByUserId = types.StringNull()
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

func (r *BlockResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    // This resource does not have a delete API endpoint.
    // Simply remove the resource from Terraform state.
    tflog.Trace(ctx, "deleted a resource (no-op: removed from state only)")
}


func (r *BlockResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}


// keepPlannedValues puts back, after a create or an update, the planned value
// of each optional attribute the configuration leaves out. The server keeps
// some of these up to date on its own (when it last checked a heartbeat, the
// status a probe last reported...), so the value read back after the write can
// already differ from the plan, and Terraform would fail the apply with
// "Provider produced inconsistent result after apply". The next refresh reads
// the server's value, which is never a diff for an attribute nobody configured.
func (r *BlockResource) keepPlannedValues(data *BlockResourceModel, plan *BlockResourceModel, config *BlockResourceModel) {
    if config.Action.IsNull() && !plan.Action.IsUnknown() {
        data.Action = plan.Action
    }
    if config.Reach.IsNull() && !plan.Reach.IsUnknown() {
        data.Reach = plan.Reach
    }
    if config.ExpiresAt.IsNull() && !plan.ExpiresAt.IsUnknown() {
        data.ExpiresAt = plan.ExpiresAt
    }
    if config.Reason.IsNull() && !plan.Reason.IsUnknown() {
        data.Reason = plan.Reason
    }
    if config.SecurityCaseId.IsNull() && !plan.SecurityCaseId.IsUnknown() {
        data.SecurityCaseId = plan.SecurityCaseId
    }
    if config.State.IsNull() && !plan.State.IsUnknown() {
        data.State = plan.State
    }
    if config.LiftReason.IsNull() && !plan.LiftReason.IsUnknown() {
        data.LiftReason = plan.LiftReason
    }
}

// Helper method to convert Terraform map to Go interface{}
func (r *BlockResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *BlockResource) convertTerraformListToInterface(terraformList types.List) interface{} {
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
func (r *BlockResource) convertTerraformSetToInterface(terraformSet types.Set) interface{} {
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
func (r *BlockResource) parseJSONField(terraformString basetypes.StringValuable) interface{} {
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
func (r *BlockResource) normalizeURLWrappers(value interface{}) interface{} {
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

func (r *BlockResource) normalizeURLString(value string) string {
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
func (r *BlockResource) bigFloatToFloat64(bf *big.Float) interface{} {
    if bf == nil {
        return nil
    }
    f, _ := bf.Float64()
    return f
}

// Helper method to check if a type string is a valid OneUptime ObjectType.
// The registry itself lives in objecttypes.go, shared across the package.
func (r *BlockResource) isValidOneUptimeObjectType(typeStr string) bool {
    return ValidOneUptimeObjectTypes[typeStr]
}
