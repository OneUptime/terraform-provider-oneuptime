package shard08

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
var _ resource.Resource = &BlocklistFeedResource{}
var _ resource.ResourceWithImportState = &BlocklistFeedResource{}
var _ resource.ResourceWithModifyPlan = &BlocklistFeedResource{}

func NewBlocklistFeedResource() resource.Resource {
    return &BlocklistFeedResource{}
}

// BlocklistFeedResource defines the resource implementation.
type BlocklistFeedResource struct {
    client *Client
}

// BlocklistFeedResourceModel describes the resource data model.
type BlocklistFeedResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Firewall types.String `tfsdk:"firewall"`
    Family types.String `tfsdk:"family"`
    AllowedSourceAddresses types.String `tfsdk:"allowed_source_addresses"`
    CreatedAt RFC3339Value `tfsdk:"created_at"`
    UpdatedAt RFC3339Value `tfsdk:"updated_at"`
    CreatedByActorType types.String `tfsdk:"created_by_actor_type"`
    CreatedByActorId types.String `tfsdk:"created_by_actor_id"`
    LastModifiedAt RFC3339Value `tfsdk:"last_modified_at"`
    LastModifiedByActorType types.String `tfsdk:"last_modified_by_actor_type"`
    LastModifiedByActorId types.String `tfsdk:"last_modified_by_actor_id"`
    LastModifiedByUserId types.String `tfsdk:"last_modified_by_user_id"`
    SecretHint types.String `tfsdk:"secret_hint"`
    SecretMadeAt RFC3339Value `tfsdk:"secret_made_at"`
    PreviousSecretExpiresAt RFC3339Value `tfsdk:"previous_secret_expires_at"`
    RevokedAt RFC3339Value `tfsdk:"revoked_at"`
    LastFetchedAt RFC3339Value `tfsdk:"last_fetched_at"`
    LastFetchSourceIp types.String `tfsdk:"last_fetch_source_ip"`
    LastFetchUserAgent types.String `tfsdk:"last_fetch_user_agent"`
    UsualFetchIntervalInSeconds types.Number `tfsdk:"usual_fetch_interval_in_seconds"`
    StoppedReadingToldAt RFC3339Value `tfsdk:"stopped_reading_told_at"`
    ServedCount types.Number `tfsdk:"served_count"`
    ServedConfirmedAt RFC3339Value `tfsdk:"served_confirmed_at"`
    Secret types.String `tfsdk:"secret"`
}

func (r *BlocklistFeedResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_blocklist_feed"
}

func (r *BlocklistFeedResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = r.schemaDefinition()
}

func (r *BlocklistFeedResource) schemaDefinition() schema.Schema {
    return schema.Schema{
        MarkdownDescription: "A secret URL a firewall or a proxy reads the project's blocklist from, on a schedule: Palo Alto, FortiGate, Check Point, Cisco, OPNsense, pfSense, nftables, nginx and others.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of the project this feed belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "What the feed is called where it is enforced: the firewall or proxy that reads it, such as Palo Alto edge firewall.",
                Required: true,
            },
            "firewall": schema.StringAttribute{
                MarkdownDescription: "The firewall or proxy that reads the feed, whose setup steps its page opens on. Any of them can read any feed.",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("Other"),
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    stringvalidator.OneOf("PaloAlto", "FortiGate", "CheckPoint", "Cisco", "OPNsense", "PfSense", "Nftables", "Ipset", "Nginx", "HAProxy", "WindowsFirewall", "Other"),
                },
            },
            "family": schema.StringAttribute{
                MarkdownDescription: "What the feed serves: IPv4 and IPv6 addresses (Both), or one of them, for a firewall that keeps them apart.",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("Both"),
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    stringvalidator.OneOf("Both", "IPv4", "IPv6"),
                },
            },
            "allowed_source_addresses": schema.StringAttribute{
                MarkdownDescription: "The only addresses and ranges the feed answers, one a line: the firewall's own. Empty to answer any address that has the secret.",
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
            "secret_hint": schema.StringAttribute{
                MarkdownDescription: "The last four characters of the feed's secret, shown in its place so a person can tell which secret a firewall holds. Never enough to use it.",
                Computed: true,
            },
            "secret_made_at": schema.StringAttribute{
                MarkdownDescription: "When the feed's current secret was made.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "previous_secret_expires_at": schema.StringAttribute{
                MarkdownDescription: "When the secret the last rotation replaced stops working: a day after the rotation, or at once when the feed's secret is revoked.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "revoked_at": schema.StringAttribute{
                MarkdownDescription: "When the feed's secret was revoked. A revoked feed answers no firewall until a new secret is made.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "last_fetched_at": schema.StringAttribute{
                MarkdownDescription: "When a firewall last fetched the feed.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "last_fetch_source_ip": schema.StringAttribute{
                MarkdownDescription: "The address the feed was last fetched from.",
                Computed: true,
            },
            "last_fetch_user_agent": schema.StringAttribute{
                MarkdownDescription: "The user agent the feed was last fetched with: the firewall or tool that read it.",
                Computed: true,
            },
            "usual_fetch_interval_in_seconds": schema.NumberAttribute{
                MarkdownDescription: "How often the firewall usually fetches the feed: the median gap between its latest fetches. Three times it without a fetch, and the firewall has stopped reading the feed.",
                Computed: true,
            },
            "stopped_reading_told_at": schema.StringAttribute{
                MarkdownDescription: "When the Project Owners were told the firewall stopped reading the feed. Cleared when it reads it again, so they are told once each time.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "served_count": schema.NumberAttribute{
                MarkdownDescription: "How many addresses and ranges the last list a firewall fetched held.",
                Computed: true,
            },
            "served_confirmed_at": schema.StringAttribute{
                MarkdownDescription: "When a firewall last fetched the feed and was found to hold its last list: the time its blocks were last confirmed there.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "secret": schema.StringAttribute{
                MarkdownDescription: "The feed's secret, which a firewall sends as a password, as a bearer token or in the URL. Shown once: in the answer to the create that makes the feed, or to a rotation of it. OneUptime keeps only its hash. Kept in state from the answer to the create: OneUptime never reads it back, so a resource imported into Terraform has none. Replace the resource for a new one.",
                Computed: true,
                Sensitive: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
        },
    }
}

// ModifyPlan keeps the state when nothing a configuration sets changes (see
// importsettle.go).
func (r *BlocklistFeedResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
    // Destroying: nothing to plan.
    if req.Plan.Raw.IsNull() {
        return
    }

    // Under Crossplane, what a manifest leaves out is the record's (see crossplane.go).
    KeepLeftOutUnderCrossplane(ctx, req, resp)

    KeepStateWhenNothingChanges(ctx, req, resp)
}

func (r *BlocklistFeedResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *BlocklistFeedResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data BlocklistFeedResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // What the configuration sets, and what Terraform planned (see keepPlannedValues).
    var config BlocklistFeedResourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
    if resp.Diagnostics.HasError() {
        return
    }
    plan := data



    // Create API request body. Unset (null/unknown) optional fields are
    // omitted so server-side defaults apply instead of being overwritten
    // with zero values.
    blocklistFeedRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := blocklistFeedRequest["data"].(map[string]interface{})

    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        requestDataMap["name"] = data.Name.ValueString()
    }
    if !data.Firewall.IsNull() && !data.Firewall.IsUnknown() {
        requestDataMap["firewall"] = data.Firewall.ValueString()
    }
    if !data.Family.IsNull() && !data.Family.IsUnknown() {
        requestDataMap["family"] = data.Family.ValueString()
    }
    if !data.AllowedSourceAddresses.IsNull() && !data.AllowedSourceAddresses.IsUnknown() {
        requestDataMap["allowedSourceAddresses"] = data.AllowedSourceAddresses.ValueString()
    }

    // Make API call
    httpResp, err := r.client.Post(ctx, "/security-blocklist-feed", blocklistFeedRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create blocklist_feed, got error: %s", err))
        return
    }

    var blocklistFeedResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &blocklistFeedResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to create blocklist_feed: %s", err))
        return
    }

    // What OneUptime left out because only it writes it, if anything (see entries.go).
    resp.Diagnostics.Append(FieldsLeftOutWarning(httpResp, "blocklist feed")...)

    // What the write does that someone should know, if anything (see entries.go).
    resp.Diagnostics.Append(WriteWarnings(httpResp, "blocklist feed")...)

    // Extract the new resource id from the create response.
    createdId := ""
    if wrapper, ok := blocklistFeedResponse["data"].(map[string]interface{}); ok {
        if val, ok := wrapper["_id"].(string); ok {
            createdId = val
        }
    } else if val, ok := blocklistFeedResponse["_id"].(string); ok {
        createdId = val
    }
    if createdId == "" {
        resp.Diagnostics.AddError("OneUptime API Error", "Create response for blocklist_feed did not contain an id. This is a bug in the provider or the API; please report it.")
        return
    }
    data.Id = types.StringValue(createdId)

    // What the create answers once, and nothing reads back: kept from its answer.
    var shownOnceMap map[string]interface{}
    if wrapper, ok := blocklistFeedResponse["data"].(map[string]interface{}); ok {
        shownOnceMap = wrapper
    } else {
        shownOnceMap = blocklistFeedResponse
    }
    if obj, ok := shownOnceMap["secret"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Secret = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Secret = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Secret = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Secret = types.StringValue(string(jsonBytes))
            } else {
                data.Secret = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Secret = types.StringValue(string(jsonBytes))
            } else {
                data.Secret = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Secret = types.StringValue(string(jsonBytes))
        } else {
            data.Secret = types.StringNull()
        }
    } else if arr, ok := shownOnceMap["secret"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Secret = types.StringValue(string(jsonBytes))
        } else {
            data.Secret = types.StringNull()
        }
    } else if val, ok := shownOnceMap["secret"].(string); ok {
        data.Secret = types.StringValue(val)
    } else {
        data.Secret = types.StringNull()
    }

    /*
     * The server has committed the row. Persist what we know to state BEFORE
     * the read-back: if the read-back fails and we return without setting
     * state, Terraform never learns the resource exists and the created
     * blocklist_feed is orphaned server-side — never refreshed, never
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
        "firewall": true,
        "family": true,
        "allowedSourceAddresses": true,
        "createdAt": true,
        "updatedAt": true,
        "createdByActorType": true,
        "createdByActorId": true,
        "lastModifiedAt": true,
        "lastModifiedByActorType": true,
        "lastModifiedByActorId": true,
        "lastModifiedByUserId": true,
        "secretHint": true,
        "secretMadeAt": true,
        "previousSecretExpiresAt": true,
        "revokedAt": true,
        "lastFetchedAt": true,
        "lastFetchSourceIp": true,
        "lastFetchUserAgent": true,
        "usualFetchIntervalInSeconds": true,
        "stoppedReadingToldAt": true,
        "servedCount": true,
        "servedConfirmedAt": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/security-blocklist-feed/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        /*
         * State already owns the id, so the resource is tracked and the next
         * refresh reconciles the remaining attributes. Warn rather than
         * error: erroring here would strand a real resource.
         */
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created blocklist_feed but could not read it back; state is incomplete until the next refresh: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created blocklist_feed but could not parse the read-back response; state is incomplete until the next refresh: %s", err))
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
    } else if arr, ok := dataMap["name"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["firewall"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Firewall = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Firewall = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Firewall = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Firewall = types.StringValue(string(jsonBytes))
            } else {
                data.Firewall = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Firewall = types.StringValue(string(jsonBytes))
            } else {
                data.Firewall = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Firewall = types.StringValue(string(jsonBytes))
        } else {
            data.Firewall = types.StringNull()
        }
    } else if arr, ok := dataMap["firewall"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Firewall = types.StringValue(string(jsonBytes))
        } else {
            data.Firewall = types.StringNull()
        }
    } else if val, ok := dataMap["firewall"].(string); ok {
        data.Firewall = types.StringValue(val)
    } else {
        data.Firewall = types.StringNull()
    }
    if obj, ok := dataMap["family"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Family = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Family = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Family = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Family = types.StringValue(string(jsonBytes))
            } else {
                data.Family = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Family = types.StringValue(string(jsonBytes))
            } else {
                data.Family = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Family = types.StringValue(string(jsonBytes))
        } else {
            data.Family = types.StringNull()
        }
    } else if arr, ok := dataMap["family"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Family = types.StringValue(string(jsonBytes))
        } else {
            data.Family = types.StringNull()
        }
    } else if val, ok := dataMap["family"].(string); ok {
        data.Family = types.StringValue(val)
    } else {
        data.Family = types.StringNull()
    }
    if obj, ok := dataMap["allowedSourceAddresses"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AllowedSourceAddresses = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AllowedSourceAddresses = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AllowedSourceAddresses = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AllowedSourceAddresses = types.StringValue(string(jsonBytes))
            } else {
                data.AllowedSourceAddresses = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AllowedSourceAddresses = types.StringValue(string(jsonBytes))
            } else {
                data.AllowedSourceAddresses = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AllowedSourceAddresses = types.StringValue(string(jsonBytes))
        } else {
            data.AllowedSourceAddresses = types.StringNull()
        }
    } else if arr, ok := dataMap["allowedSourceAddresses"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AllowedSourceAddresses = types.StringValue(string(jsonBytes))
        } else {
            data.AllowedSourceAddresses = types.StringNull()
        }
    } else if val, ok := dataMap["allowedSourceAddresses"].(string); ok {
        data.AllowedSourceAddresses = types.StringValue(val)
    } else {
        data.AllowedSourceAddresses = types.StringNull()
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
    if obj, ok := dataMap["secretHint"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SecretHint = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SecretHint = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SecretHint = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SecretHint = types.StringValue(string(jsonBytes))
            } else {
                data.SecretHint = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SecretHint = types.StringValue(string(jsonBytes))
            } else {
                data.SecretHint = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SecretHint = types.StringValue(string(jsonBytes))
        } else {
            data.SecretHint = types.StringNull()
        }
    } else if arr, ok := dataMap["secretHint"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.SecretHint = types.StringValue(string(jsonBytes))
        } else {
            data.SecretHint = types.StringNull()
        }
    } else if val, ok := dataMap["secretHint"].(string); ok {
        data.SecretHint = types.StringValue(val)
    } else {
        data.SecretHint = types.StringNull()
    }
    if obj, ok := dataMap["secretMadeAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.SecretMadeAt = NewRFC3339Value(val)
        } else {
            data.SecretMadeAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["secretMadeAt"].(string); ok && val != "" {
        data.SecretMadeAt = NewRFC3339Value(val)
    } else {
        data.SecretMadeAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["previousSecretExpiresAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.PreviousSecretExpiresAt = NewRFC3339Value(val)
        } else {
            data.PreviousSecretExpiresAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["previousSecretExpiresAt"].(string); ok && val != "" {
        data.PreviousSecretExpiresAt = NewRFC3339Value(val)
    } else {
        data.PreviousSecretExpiresAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["revokedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.RevokedAt = NewRFC3339Value(val)
        } else {
            data.RevokedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["revokedAt"].(string); ok && val != "" {
        data.RevokedAt = NewRFC3339Value(val)
    } else {
        data.RevokedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastFetchedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastFetchedAt = NewRFC3339Value(val)
        } else {
            data.LastFetchedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastFetchedAt"].(string); ok && val != "" {
        data.LastFetchedAt = NewRFC3339Value(val)
    } else {
        data.LastFetchedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastFetchSourceIp"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastFetchSourceIp = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastFetchSourceIp = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastFetchSourceIp = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastFetchSourceIp = types.StringValue(string(jsonBytes))
            } else {
                data.LastFetchSourceIp = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastFetchSourceIp = types.StringValue(string(jsonBytes))
            } else {
                data.LastFetchSourceIp = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastFetchSourceIp = types.StringValue(string(jsonBytes))
        } else {
            data.LastFetchSourceIp = types.StringNull()
        }
    } else if arr, ok := dataMap["lastFetchSourceIp"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LastFetchSourceIp = types.StringValue(string(jsonBytes))
        } else {
            data.LastFetchSourceIp = types.StringNull()
        }
    } else if val, ok := dataMap["lastFetchSourceIp"].(string); ok {
        data.LastFetchSourceIp = types.StringValue(val)
    } else {
        data.LastFetchSourceIp = types.StringNull()
    }
    if obj, ok := dataMap["lastFetchUserAgent"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastFetchUserAgent = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastFetchUserAgent = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastFetchUserAgent = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastFetchUserAgent = types.StringValue(string(jsonBytes))
            } else {
                data.LastFetchUserAgent = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastFetchUserAgent = types.StringValue(string(jsonBytes))
            } else {
                data.LastFetchUserAgent = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastFetchUserAgent = types.StringValue(string(jsonBytes))
        } else {
            data.LastFetchUserAgent = types.StringNull()
        }
    } else if arr, ok := dataMap["lastFetchUserAgent"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LastFetchUserAgent = types.StringValue(string(jsonBytes))
        } else {
            data.LastFetchUserAgent = types.StringNull()
        }
    } else if val, ok := dataMap["lastFetchUserAgent"].(string); ok {
        data.LastFetchUserAgent = types.StringValue(val)
    } else {
        data.LastFetchUserAgent = types.StringNull()
    }
    data.UsualFetchIntervalInSeconds = NumberFromAPI(data.UsualFetchIntervalInSeconds, dataMap["usualFetchIntervalInSeconds"])
    if obj, ok := dataMap["stoppedReadingToldAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.StoppedReadingToldAt = NewRFC3339Value(val)
        } else {
            data.StoppedReadingToldAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["stoppedReadingToldAt"].(string); ok && val != "" {
        data.StoppedReadingToldAt = NewRFC3339Value(val)
    } else {
        data.StoppedReadingToldAt = NewRFC3339Null()
    }
    data.ServedCount = NumberFromAPI(data.ServedCount, dataMap["servedCount"])
    if obj, ok := dataMap["servedConfirmedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ServedConfirmedAt = NewRFC3339Value(val)
        } else {
            data.ServedConfirmedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["servedConfirmedAt"].(string); ok && val != "" {
        data.ServedConfirmedAt = NewRFC3339Value(val)
    } else {
        data.ServedConfirmedAt = NewRFC3339Null()
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

func (r *BlocklistFeedResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data BlocklistFeedResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "firewall": true,
        "family": true,
        "allowedSourceAddresses": true,
        "createdAt": true,
        "updatedAt": true,
        "createdByActorType": true,
        "createdByActorId": true,
        "lastModifiedAt": true,
        "lastModifiedByActorType": true,
        "lastModifiedByActorId": true,
        "lastModifiedByUserId": true,
        "secretHint": true,
        "secretMadeAt": true,
        "previousSecretExpiresAt": true,
        "revokedAt": true,
        "lastFetchedAt": true,
        "lastFetchSourceIp": true,
        "lastFetchUserAgent": true,
        "usualFetchIntervalInSeconds": true,
        "stoppedReadingToldAt": true,
        "servedCount": true,
        "servedConfirmedAt": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect(ctx, "/security-blocklist-feed/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read blocklist_feed, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var blocklistFeedResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &blocklistFeedResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse blocklist_feed response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := blocklistFeedResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = blocklistFeedResponse
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
    } else if arr, ok := dataMap["name"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["firewall"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Firewall = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Firewall = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Firewall = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Firewall = types.StringValue(string(jsonBytes))
            } else {
                data.Firewall = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Firewall = types.StringValue(string(jsonBytes))
            } else {
                data.Firewall = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Firewall = types.StringValue(string(jsonBytes))
        } else {
            data.Firewall = types.StringNull()
        }
    } else if arr, ok := dataMap["firewall"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Firewall = types.StringValue(string(jsonBytes))
        } else {
            data.Firewall = types.StringNull()
        }
    } else if val, ok := dataMap["firewall"].(string); ok {
        data.Firewall = types.StringValue(val)
    } else {
        data.Firewall = types.StringNull()
    }
    if obj, ok := dataMap["family"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Family = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Family = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Family = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Family = types.StringValue(string(jsonBytes))
            } else {
                data.Family = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Family = types.StringValue(string(jsonBytes))
            } else {
                data.Family = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Family = types.StringValue(string(jsonBytes))
        } else {
            data.Family = types.StringNull()
        }
    } else if arr, ok := dataMap["family"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Family = types.StringValue(string(jsonBytes))
        } else {
            data.Family = types.StringNull()
        }
    } else if val, ok := dataMap["family"].(string); ok {
        data.Family = types.StringValue(val)
    } else {
        data.Family = types.StringNull()
    }
    if obj, ok := dataMap["allowedSourceAddresses"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AllowedSourceAddresses = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AllowedSourceAddresses = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AllowedSourceAddresses = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AllowedSourceAddresses = types.StringValue(string(jsonBytes))
            } else {
                data.AllowedSourceAddresses = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AllowedSourceAddresses = types.StringValue(string(jsonBytes))
            } else {
                data.AllowedSourceAddresses = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AllowedSourceAddresses = types.StringValue(string(jsonBytes))
        } else {
            data.AllowedSourceAddresses = types.StringNull()
        }
    } else if arr, ok := dataMap["allowedSourceAddresses"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AllowedSourceAddresses = types.StringValue(string(jsonBytes))
        } else {
            data.AllowedSourceAddresses = types.StringNull()
        }
    } else if val, ok := dataMap["allowedSourceAddresses"].(string); ok {
        data.AllowedSourceAddresses = types.StringValue(val)
    } else {
        data.AllowedSourceAddresses = types.StringNull()
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
    if obj, ok := dataMap["secretHint"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SecretHint = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SecretHint = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SecretHint = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SecretHint = types.StringValue(string(jsonBytes))
            } else {
                data.SecretHint = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SecretHint = types.StringValue(string(jsonBytes))
            } else {
                data.SecretHint = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SecretHint = types.StringValue(string(jsonBytes))
        } else {
            data.SecretHint = types.StringNull()
        }
    } else if arr, ok := dataMap["secretHint"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.SecretHint = types.StringValue(string(jsonBytes))
        } else {
            data.SecretHint = types.StringNull()
        }
    } else if val, ok := dataMap["secretHint"].(string); ok {
        data.SecretHint = types.StringValue(val)
    } else {
        data.SecretHint = types.StringNull()
    }
    if obj, ok := dataMap["secretMadeAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.SecretMadeAt = NewRFC3339Value(val)
        } else {
            data.SecretMadeAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["secretMadeAt"].(string); ok && val != "" {
        data.SecretMadeAt = NewRFC3339Value(val)
    } else {
        data.SecretMadeAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["previousSecretExpiresAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.PreviousSecretExpiresAt = NewRFC3339Value(val)
        } else {
            data.PreviousSecretExpiresAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["previousSecretExpiresAt"].(string); ok && val != "" {
        data.PreviousSecretExpiresAt = NewRFC3339Value(val)
    } else {
        data.PreviousSecretExpiresAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["revokedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.RevokedAt = NewRFC3339Value(val)
        } else {
            data.RevokedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["revokedAt"].(string); ok && val != "" {
        data.RevokedAt = NewRFC3339Value(val)
    } else {
        data.RevokedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastFetchedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastFetchedAt = NewRFC3339Value(val)
        } else {
            data.LastFetchedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastFetchedAt"].(string); ok && val != "" {
        data.LastFetchedAt = NewRFC3339Value(val)
    } else {
        data.LastFetchedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastFetchSourceIp"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastFetchSourceIp = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastFetchSourceIp = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastFetchSourceIp = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastFetchSourceIp = types.StringValue(string(jsonBytes))
            } else {
                data.LastFetchSourceIp = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastFetchSourceIp = types.StringValue(string(jsonBytes))
            } else {
                data.LastFetchSourceIp = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastFetchSourceIp = types.StringValue(string(jsonBytes))
        } else {
            data.LastFetchSourceIp = types.StringNull()
        }
    } else if arr, ok := dataMap["lastFetchSourceIp"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LastFetchSourceIp = types.StringValue(string(jsonBytes))
        } else {
            data.LastFetchSourceIp = types.StringNull()
        }
    } else if val, ok := dataMap["lastFetchSourceIp"].(string); ok {
        data.LastFetchSourceIp = types.StringValue(val)
    } else {
        data.LastFetchSourceIp = types.StringNull()
    }
    if obj, ok := dataMap["lastFetchUserAgent"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastFetchUserAgent = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastFetchUserAgent = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastFetchUserAgent = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastFetchUserAgent = types.StringValue(string(jsonBytes))
            } else {
                data.LastFetchUserAgent = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastFetchUserAgent = types.StringValue(string(jsonBytes))
            } else {
                data.LastFetchUserAgent = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastFetchUserAgent = types.StringValue(string(jsonBytes))
        } else {
            data.LastFetchUserAgent = types.StringNull()
        }
    } else if arr, ok := dataMap["lastFetchUserAgent"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LastFetchUserAgent = types.StringValue(string(jsonBytes))
        } else {
            data.LastFetchUserAgent = types.StringNull()
        }
    } else if val, ok := dataMap["lastFetchUserAgent"].(string); ok {
        data.LastFetchUserAgent = types.StringValue(val)
    } else {
        data.LastFetchUserAgent = types.StringNull()
    }
    data.UsualFetchIntervalInSeconds = NumberFromAPI(data.UsualFetchIntervalInSeconds, dataMap["usualFetchIntervalInSeconds"])
    if obj, ok := dataMap["stoppedReadingToldAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.StoppedReadingToldAt = NewRFC3339Value(val)
        } else {
            data.StoppedReadingToldAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["stoppedReadingToldAt"].(string); ok && val != "" {
        data.StoppedReadingToldAt = NewRFC3339Value(val)
    } else {
        data.StoppedReadingToldAt = NewRFC3339Null()
    }
    data.ServedCount = NumberFromAPI(data.ServedCount, dataMap["servedCount"])
    if obj, ok := dataMap["servedConfirmedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ServedConfirmedAt = NewRFC3339Value(val)
        } else {
            data.ServedConfirmedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["servedConfirmedAt"].(string); ok && val != "" {
        data.ServedConfirmedAt = NewRFC3339Value(val)
    } else {
        data.ServedConfirmedAt = NewRFC3339Null()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BlocklistFeedResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data BlocklistFeedResourceModel
    var state BlocklistFeedResourceModel

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
    // Answered once, by the create: kept as state holds it.
    data.Secret = state.Secret

    // What the configuration sets, and what Terraform planned (see keepPlannedValues).
    var config BlocklistFeedResourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
    if resp.Diagnostics.HasError() {
        return
    }
    plan := data

    // Create API request body
    blocklistFeedRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := blocklistFeedRequest["data"].(map[string]interface{})

    if !data.Name.IsUnknown() && !state.Name.IsUnknown() && !data.Name.Equal(state.Name) {
        requestDataMap["name"] = data.Name.ValueString()
    }
    if !data.Firewall.IsUnknown() && !state.Firewall.IsUnknown() && !data.Firewall.Equal(state.Firewall) {
        requestDataMap["firewall"] = data.Firewall.ValueString()
    }
    if !data.Family.IsUnknown() && !state.Family.IsUnknown() && !data.Family.Equal(state.Family) {
        requestDataMap["family"] = data.Family.ValueString()
    }
    if !data.AllowedSourceAddresses.IsUnknown() && !state.AllowedSourceAddresses.IsUnknown() && !data.AllowedSourceAddresses.Equal(state.AllowedSourceAddresses) {
        requestDataMap["allowedSourceAddresses"] = data.AllowedSourceAddresses.ValueString()
    }

    // Only call the API when there are changed fields to send. An empty
    // update body is rejected by the API; state is still refreshed below so
    // this method never writes unverified plan values into state.
    if len(blocklistFeedRequest["data"].(map[string]interface{})) > 0 {
        httpResp, err := r.client.Put(ctx, "/security-blocklist-feed/" + data.Id.ValueString() + "", blocklistFeedRequest)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update blocklist_feed, got error: %s", err))
            return
        }

        // Parse the update response
        var blocklistFeedResponse map[string]interface{}
        err = r.client.ParseResponse(httpResp, &blocklistFeedResponse)
        if err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to update blocklist_feed: %s", err))
            return
        }
        _ = blocklistFeedResponse

        // What the write does that someone should know, if anything (see entries.go).
        resp.Diagnostics.Append(WriteWarnings(httpResp, "blocklist feed")...)
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "firewall": true,
        "family": true,
        "allowedSourceAddresses": true,
        "createdAt": true,
        "updatedAt": true,
        "createdByActorType": true,
        "createdByActorId": true,
        "lastModifiedAt": true,
        "lastModifiedByActorType": true,
        "lastModifiedByActorId": true,
        "lastModifiedByUserId": true,
        "secretHint": true,
        "secretMadeAt": true,
        "previousSecretExpiresAt": true,
        "revokedAt": true,
        "lastFetchedAt": true,
        "lastFetchSourceIp": true,
        "lastFetchUserAgent": true,
        "usualFetchIntervalInSeconds": true,
        "stoppedReadingToldAt": true,
        "servedCount": true,
        "servedConfirmedAt": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/security-blocklist-feed/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read blocklist_feed after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read blocklist_feed after update: %s", err))
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
    } else if arr, ok := dataMap["name"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["firewall"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Firewall = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Firewall = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Firewall = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Firewall = types.StringValue(string(jsonBytes))
            } else {
                data.Firewall = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Firewall = types.StringValue(string(jsonBytes))
            } else {
                data.Firewall = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Firewall = types.StringValue(string(jsonBytes))
        } else {
            data.Firewall = types.StringNull()
        }
    } else if arr, ok := dataMap["firewall"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Firewall = types.StringValue(string(jsonBytes))
        } else {
            data.Firewall = types.StringNull()
        }
    } else if val, ok := dataMap["firewall"].(string); ok {
        data.Firewall = types.StringValue(val)
    } else {
        data.Firewall = types.StringNull()
    }
    if obj, ok := dataMap["family"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Family = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Family = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Family = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Family = types.StringValue(string(jsonBytes))
            } else {
                data.Family = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Family = types.StringValue(string(jsonBytes))
            } else {
                data.Family = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Family = types.StringValue(string(jsonBytes))
        } else {
            data.Family = types.StringNull()
        }
    } else if arr, ok := dataMap["family"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Family = types.StringValue(string(jsonBytes))
        } else {
            data.Family = types.StringNull()
        }
    } else if val, ok := dataMap["family"].(string); ok {
        data.Family = types.StringValue(val)
    } else {
        data.Family = types.StringNull()
    }
    if obj, ok := dataMap["allowedSourceAddresses"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AllowedSourceAddresses = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AllowedSourceAddresses = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AllowedSourceAddresses = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AllowedSourceAddresses = types.StringValue(string(jsonBytes))
            } else {
                data.AllowedSourceAddresses = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AllowedSourceAddresses = types.StringValue(string(jsonBytes))
            } else {
                data.AllowedSourceAddresses = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AllowedSourceAddresses = types.StringValue(string(jsonBytes))
        } else {
            data.AllowedSourceAddresses = types.StringNull()
        }
    } else if arr, ok := dataMap["allowedSourceAddresses"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AllowedSourceAddresses = types.StringValue(string(jsonBytes))
        } else {
            data.AllowedSourceAddresses = types.StringNull()
        }
    } else if val, ok := dataMap["allowedSourceAddresses"].(string); ok {
        data.AllowedSourceAddresses = types.StringValue(val)
    } else {
        data.AllowedSourceAddresses = types.StringNull()
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
    if obj, ok := dataMap["secretHint"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SecretHint = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SecretHint = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SecretHint = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SecretHint = types.StringValue(string(jsonBytes))
            } else {
                data.SecretHint = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SecretHint = types.StringValue(string(jsonBytes))
            } else {
                data.SecretHint = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SecretHint = types.StringValue(string(jsonBytes))
        } else {
            data.SecretHint = types.StringNull()
        }
    } else if arr, ok := dataMap["secretHint"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.SecretHint = types.StringValue(string(jsonBytes))
        } else {
            data.SecretHint = types.StringNull()
        }
    } else if val, ok := dataMap["secretHint"].(string); ok {
        data.SecretHint = types.StringValue(val)
    } else {
        data.SecretHint = types.StringNull()
    }
    if obj, ok := dataMap["secretMadeAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.SecretMadeAt = NewRFC3339Value(val)
        } else {
            data.SecretMadeAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["secretMadeAt"].(string); ok && val != "" {
        data.SecretMadeAt = NewRFC3339Value(val)
    } else {
        data.SecretMadeAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["previousSecretExpiresAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.PreviousSecretExpiresAt = NewRFC3339Value(val)
        } else {
            data.PreviousSecretExpiresAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["previousSecretExpiresAt"].(string); ok && val != "" {
        data.PreviousSecretExpiresAt = NewRFC3339Value(val)
    } else {
        data.PreviousSecretExpiresAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["revokedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.RevokedAt = NewRFC3339Value(val)
        } else {
            data.RevokedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["revokedAt"].(string); ok && val != "" {
        data.RevokedAt = NewRFC3339Value(val)
    } else {
        data.RevokedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastFetchedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastFetchedAt = NewRFC3339Value(val)
        } else {
            data.LastFetchedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastFetchedAt"].(string); ok && val != "" {
        data.LastFetchedAt = NewRFC3339Value(val)
    } else {
        data.LastFetchedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["lastFetchSourceIp"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastFetchSourceIp = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastFetchSourceIp = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastFetchSourceIp = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastFetchSourceIp = types.StringValue(string(jsonBytes))
            } else {
                data.LastFetchSourceIp = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastFetchSourceIp = types.StringValue(string(jsonBytes))
            } else {
                data.LastFetchSourceIp = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastFetchSourceIp = types.StringValue(string(jsonBytes))
        } else {
            data.LastFetchSourceIp = types.StringNull()
        }
    } else if arr, ok := dataMap["lastFetchSourceIp"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LastFetchSourceIp = types.StringValue(string(jsonBytes))
        } else {
            data.LastFetchSourceIp = types.StringNull()
        }
    } else if val, ok := dataMap["lastFetchSourceIp"].(string); ok {
        data.LastFetchSourceIp = types.StringValue(val)
    } else {
        data.LastFetchSourceIp = types.StringNull()
    }
    if obj, ok := dataMap["lastFetchUserAgent"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastFetchUserAgent = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastFetchUserAgent = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastFetchUserAgent = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastFetchUserAgent = types.StringValue(string(jsonBytes))
            } else {
                data.LastFetchUserAgent = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastFetchUserAgent = types.StringValue(string(jsonBytes))
            } else {
                data.LastFetchUserAgent = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastFetchUserAgent = types.StringValue(string(jsonBytes))
        } else {
            data.LastFetchUserAgent = types.StringNull()
        }
    } else if arr, ok := dataMap["lastFetchUserAgent"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LastFetchUserAgent = types.StringValue(string(jsonBytes))
        } else {
            data.LastFetchUserAgent = types.StringNull()
        }
    } else if val, ok := dataMap["lastFetchUserAgent"].(string); ok {
        data.LastFetchUserAgent = types.StringValue(val)
    } else {
        data.LastFetchUserAgent = types.StringNull()
    }
    data.UsualFetchIntervalInSeconds = NumberFromAPI(data.UsualFetchIntervalInSeconds, dataMap["usualFetchIntervalInSeconds"])
    if obj, ok := dataMap["stoppedReadingToldAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.StoppedReadingToldAt = NewRFC3339Value(val)
        } else {
            data.StoppedReadingToldAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["stoppedReadingToldAt"].(string); ok && val != "" {
        data.StoppedReadingToldAt = NewRFC3339Value(val)
    } else {
        data.StoppedReadingToldAt = NewRFC3339Null()
    }
    data.ServedCount = NumberFromAPI(data.ServedCount, dataMap["servedCount"])
    if obj, ok := dataMap["servedConfirmedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ServedConfirmedAt = NewRFC3339Value(val)
        } else {
            data.ServedConfirmedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["servedConfirmedAt"].(string); ok && val != "" {
        data.ServedConfirmedAt = NewRFC3339Value(val)
    } else {
        data.ServedConfirmedAt = NewRFC3339Null()
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

func (r *BlocklistFeedResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data BlocklistFeedResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    httpResp, err := r.client.Delete(ctx, "/security-blocklist-feed/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete blocklist_feed, got error: %s", err))
        return
    }

    // A failed delete must keep the resource in state — silently dropping it
    // orphans real infrastructure. 404 means it is already gone.
    if httpResp.StatusCode >= 400 && httpResp.StatusCode != http.StatusNotFound {
        err = r.client.ParseResponse(httpResp, nil)
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to delete blocklist_feed: %s", err))
        return
    }
    if httpResp.Body != nil {
        httpResp.Body.Close()
    }
}


func (r *BlocklistFeedResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}


// keepPlannedValues puts back, after a create or an update, the planned value
// of each optional attribute the configuration leaves out. The server keeps
// some of these up to date on its own (when it last checked a heartbeat, the
// status a probe last reported...), so the value read back after the write can
// already differ from the plan, and Terraform would fail the apply with
// "Provider produced inconsistent result after apply". The next refresh reads
// the server's value, which is never a diff for an attribute nobody configured.
func (r *BlocklistFeedResource) keepPlannedValues(data *BlocklistFeedResourceModel, plan *BlocklistFeedResourceModel, config *BlocklistFeedResourceModel) {
    if config.Firewall.IsNull() && !plan.Firewall.IsUnknown() {
        data.Firewall = plan.Firewall
    }
    if config.Family.IsNull() && !plan.Family.IsUnknown() {
        data.Family = plan.Family
    }
    if config.AllowedSourceAddresses.IsNull() && !plan.AllowedSourceAddresses.IsUnknown() {
        data.AllowedSourceAddresses = plan.AllowedSourceAddresses
    }
}

// Helper method to convert Terraform map to Go interface{}
func (r *BlocklistFeedResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *BlocklistFeedResource) convertTerraformListToInterface(terraformList types.List) interface{} {
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
func (r *BlocklistFeedResource) convertTerraformSetToInterface(terraformSet types.Set) interface{} {
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
func (r *BlocklistFeedResource) parseJSONField(terraformString basetypes.StringValuable) interface{} {
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
func (r *BlocklistFeedResource) normalizeURLWrappers(value interface{}) interface{} {
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

func (r *BlocklistFeedResource) normalizeURLString(value string) string {
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
func (r *BlocklistFeedResource) bigFloatToFloat64(bf *big.Float) interface{} {
    if bf == nil {
        return nil
    }
    f, _ := bf.Float64()
    return f
}

// Helper method to check if a type string is a valid OneUptime ObjectType.
// The registry itself lives in objecttypes.go, shared across the package.
func (r *BlocklistFeedResource) isValidOneUptimeObjectType(typeStr string) bool {
    return ValidOneUptimeObjectTypes[typeStr]
}
