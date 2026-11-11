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
    "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ExceptionResource{}
var _ resource.ResourceWithImportState = &ExceptionResource{}
var _ resource.ResourceWithModifyPlan = &ExceptionResource{}

func NewExceptionResource() resource.Resource {
    return &ExceptionResource{}
}

// ExceptionResource defines the resource implementation.
type ExceptionResource struct {
    client *Client
}

// ExceptionResourceModel describes the resource data model.
type ExceptionResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    PrimaryEntityId types.String `tfsdk:"primary_entity_id"`
    PrimaryEntityType types.String `tfsdk:"primary_entity_type"`
    Message types.String `tfsdk:"message"`
    StackTrace types.String `tfsdk:"stack_trace"`
    ExceptionType types.String `tfsdk:"exception_type"`
    Fingerprint types.String `tfsdk:"fingerprint"`
    FirstSeenAt RFC3339Value `tfsdk:"first_seen_at"`
    LastSeenAt RFC3339Value `tfsdk:"last_seen_at"`
    AssignToUserId types.String `tfsdk:"assign_to_user_id"`
    AssignToTeamId types.String `tfsdk:"assign_to_team_id"`
    IsResolved types.Bool `tfsdk:"is_resolved"`
    IsArchived types.Bool `tfsdk:"is_archived"`
    OccuranceCount types.Number `tfsdk:"occurance_count"`
    FirstSeenInRelease types.String `tfsdk:"first_seen_in_release"`
    LastSeenInRelease types.String `tfsdk:"last_seen_in_release"`
    Environment types.String `tfsdk:"environment"`
    RegressedAt RFC3339Value `tfsdk:"regressed_at"`
    ResolvedInReleaseAfter types.String `tfsdk:"resolved_in_release_after"`
    CreatedAt RFC3339Value `tfsdk:"created_at"`
    UpdatedAt RFC3339Value `tfsdk:"updated_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    MarkedAsResolvedAt RFC3339Value `tfsdk:"marked_as_resolved_at"`
    MarkedAsArchivedAt RFC3339Value `tfsdk:"marked_as_archived_at"`
    MarkedAsResolvedByUserId types.String `tfsdk:"marked_as_resolved_by_user_id"`
    MarkedAsArchivedByUserId types.String `tfsdk:"marked_as_archived_by_user_id"`
    Unhandled types.Bool `tfsdk:"unhandled"`
    AiClassification types.String `tfsdk:"ai_classification"`
    ErrorClass types.String `tfsdk:"error_class"`
    ErrorClassSource types.String `tfsdk:"error_class_source"`
    AiFixDeclinedAt RFC3339Value `tfsdk:"ai_fix_declined_at"`
    AiExplanation JSONSubsetValue `tfsdk:"ai_explanation"`
    AiExplanationGeneratedAt RFC3339Value `tfsdk:"ai_explanation_generated_at"`
    MergedIntoTelemetryExceptionId types.String `tfsdk:"merged_into_telemetry_exception_id"`
    MergedAt RFC3339Value `tfsdk:"merged_at"`
    MergedByUserId types.String `tfsdk:"merged_by_user_id"`
    AssignToUserEmail types.String `tfsdk:"assign_to_user_email"`
}

func (r *ExceptionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_exception"
}

func (r *ExceptionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = r.schemaDefinition()
}

func (r *ExceptionResource) schemaDefinition() schema.Schema {
    return schema.Schema{
        MarkdownDescription: "List of all Telemetry Exceptions created for the telemetry service for this OneUptime project and it's status.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "primary_entity_id": schema.StringAttribute{
                MarkdownDescription: "ID of the resource this exception belongs to (Service / Host / DockerHost / KubernetesCluster, or the projectId for unattributed telemetry — disambiguated by primaryEntityType).",
                Required: true,
            },
            "primary_entity_type": schema.StringAttribute{
                MarkdownDescription: "Resource type that produced this exception (e.g. OpenTelemetry service, Host, DockerHost, KubernetesCluster, or Unknown for unattributed telemetry).",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                    stringplanmodifier.RequiresReplace(),
                },
                Validators: []validator.String{
                    stringvalidator.OneOf("OpenTelemetry", "Monitor", "NetworkDevice", "Alert", "Incident", "ScheduledMaintenance", "ServiceLevelObjective", "Host", "DockerHost", "PodmanHost", "KubernetesCluster", "ProxmoxCluster", "OpenStackCloud", "CephCluster", "StorageArray", "DockerSwarmCluster", "VMwareVCenter", "IoTDevice", "ServerlessFunction", "CloudResource", "RealUserMonitor", "DatabaseServer", "DataStore", "DataPipeline", "DataProduct", "DataReconciliation", "DataBiConnection", "Unknown"),
                },
            },
            "message": schema.StringAttribute{
                MarkdownDescription: "Exception message that was thrown by the telemetry service.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "stack_trace": schema.StringAttribute{
                MarkdownDescription: "Stack trace of the exception that was thrown by the telemetry service.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "exception_type": schema.StringAttribute{
                MarkdownDescription: "Type of the exception that was thrown by the telemetry service.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "fingerprint": schema.StringAttribute{
                MarkdownDescription: "Finger print of the exception that was thrown by the telemetry service.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "first_seen_at": schema.StringAttribute{
                MarkdownDescription: "When the earliest occurrence of this exception happened.",
                CustomType: RFC3339Type{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    RFC3339KeepSameInstant(),
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "When the latest occurrence of this exception happened.",
                CustomType: RFC3339Type{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    RFC3339KeepSameInstant(),
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "assign_to_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who this exception is assigned to. The ID of a `oneuptime_user` (see the data source). Or name them by email address, in `assign_to_user_email`.",
                Optional: true,
                Computed: true,
                Validators: personIDValidators("assign_to_user_email", false),
            },
            "assign_to_team_id": schema.StringAttribute{
                MarkdownDescription: "Team ID who this exception is assigned to. The ID of a `oneuptime_team`.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: recordIDValidators("oneuptime_team"),
            },
            "is_resolved": schema.BoolAttribute{
                MarkdownDescription: "Is this exception resolved?",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Is this exception archived?",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "occurance_count": schema.NumberAttribute{
                MarkdownDescription: "Number of times this exception has occurred.",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(1)),
                PlanModifiers: []planmodifier.Number{
                    NumberKeepSameAPINumber(),
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "first_seen_in_release": schema.StringAttribute{
                MarkdownDescription: "The service version / release in which this exception was first observed.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "last_seen_in_release": schema.StringAttribute{
                MarkdownDescription: "The most recent service version / release in which this exception was observed.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "environment": schema.StringAttribute{
                MarkdownDescription: "Deployment environment from deployment.environment resource attribute.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "regressed_at": schema.StringAttribute{
                MarkdownDescription: "When this exception last came back after it was resolved. Empty when it has not come back since someone last resolved or reopened it.",
                CustomType: RFC3339Type{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    RFC3339KeepSameInstant(),
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "resolved_in_release_after": schema.StringAttribute{
                MarkdownDescription: "Set when the exception was resolved in the next release: the release it was resolved after. Occurrences from this release do not reopen the exception; an occurrence from any other release does.",
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
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "marked_as_resolved_at": schema.StringAttribute{
                MarkdownDescription: "When this exception was last resolved. Empty while it is unresolved.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "marked_as_archived_at": schema.StringAttribute{
                MarkdownDescription: "When this exception was last archived. Empty while it is not archived.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "marked_as_resolved_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who marked this exception as resolved. The ID of a `oneuptime_user` (see the data source).",
                Computed: true,
            },
            "marked_as_archived_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who marked this exception as archived. The ID of a `oneuptime_user` (see the data source).",
                Computed: true,
            },
            "unhandled": schema.BoolAttribute{
                MarkdownDescription: "True when at least one occurrence of this exception escaped its span scope (was unhandled, per OTel exception.escaped).",
                Computed: true,
            },
            "ai_classification": schema.StringAttribute{
                MarkdownDescription: "AI triage verdict for this exception group (code-fault, user-error, expected-denial, infrastructure).",
                Computed: true,
            },
            "error_class": schema.StringAttribute{
                MarkdownDescription: "Fault domain of this exception group (code-fault, user-error, expected-denial, infrastructure, unknown). Non-actionable classes are excluded from the Issues list.",
                Computed: true,
            },
            "error_class_source": schema.StringAttribute{
                MarkdownDescription: "Where the error class came from: default (unclassified), declared (by the emitting code), ai (triage verdict) or manual (a human).",
                Computed: true,
            },
            "ai_fix_declined_at": schema.StringAttribute{
                MarkdownDescription: "Set when an AI-authored fix pull request for this exception was closed without merging; suppresses further automatic fix attempts.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "ai_explanation": schema.StringAttribute{
                MarkdownDescription: "OneUptime AI's explanation of this exception: a summary, the likely root cause, a suggested fix, how confident it is, the suspect stack frames and next steps. A JSON value: read it with `jsondecode()`.",
                CustomType: JSONSubsetType{},
                Computed: true,
            },
            "ai_explanation_generated_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime AI last wrote its explanation of this exception.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "merged_into_telemetry_exception_id": schema.StringAttribute{
                MarkdownDescription: "ID of the exception this one was merged into as the same bug. Empty for an exception that was not merged. The ID of a `oneuptime_exception`.",
                Computed: true,
            },
            "merged_at": schema.StringAttribute{
                MarkdownDescription: "When this exception was merged into another.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "merged_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the user who merged this exception into another. The ID of a `oneuptime_user` (see the data source).",
                Computed: true,
            },
            "assign_to_user_email": schema.StringAttribute{
                MarkdownDescription: "The email address of the person in `assign_to_user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.",
                Optional: true,
                Computed: true,
                PlanModifiers: personEmailPlanModifiers("assign_to_user_id"),
                Validators: personEmailValidators(),
            },
        },
    }
}

// ModifyPlan finds, by email address, the people this exception names (see people.go),
// and keeps the state when nothing a configuration sets changes (see
// importsettle.go).
func (r *ExceptionResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
    // Destroying: nothing to plan.
    if req.Plan.Raw.IsNull() {
        return
    }

    planPersonID(ctx, r.client, req, resp, personAttribute{
        id:       "assign_to_user_id",
        email:    "assign_to_user_email",
        invites:  false,
        replaces: false,
    })

    keepStateWhenNothingChanges(ctx, req, resp)
}

func (r *ExceptionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *ExceptionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data ExceptionResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // What the configuration sets, and what Terraform planned (see keepPlannedValues).
    var config ExceptionResourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
    if resp.Diagnostics.HasError() {
        return
    }
    plan := data


    // People named by email address: the user IDs to send (see people.go).
    if !personIDForWrite(ctx, r.client, &data.AssignToUserId, data.AssignToUserEmail, "assign_to_user_email", &resp.Diagnostics) {
        return
    }


    // Create API request body. Unset (null/unknown) optional fields are
    // omitted so server-side defaults apply instead of being overwritten
    // with zero values.
    exceptionRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := exceptionRequest["data"].(map[string]interface{})

    if !data.PrimaryEntityId.IsNull() && !data.PrimaryEntityId.IsUnknown() {
        requestDataMap["primaryEntityId"] = data.PrimaryEntityId.ValueString()
    }
    if !data.PrimaryEntityType.IsNull() && !data.PrimaryEntityType.IsUnknown() {
        requestDataMap["primaryEntityType"] = data.PrimaryEntityType.ValueString()
    }
    if !data.Message.IsNull() && !data.Message.IsUnknown() {
        requestDataMap["message"] = data.Message.ValueString()
    }
    if !data.StackTrace.IsNull() && !data.StackTrace.IsUnknown() {
        requestDataMap["stackTrace"] = data.StackTrace.ValueString()
    }
    if !data.ExceptionType.IsNull() && !data.ExceptionType.IsUnknown() {
        requestDataMap["exceptionType"] = data.ExceptionType.ValueString()
    }
    if !data.Fingerprint.IsNull() && !data.Fingerprint.IsUnknown() {
        requestDataMap["fingerprint"] = data.Fingerprint.ValueString()
    }
    if !data.FirstSeenAt.IsNull() && !data.FirstSeenAt.IsUnknown() {
        requestDataMap["firstSeenAt"] = data.FirstSeenAt.ValueString()
    }
    if !data.LastSeenAt.IsNull() && !data.LastSeenAt.IsUnknown() {
        requestDataMap["lastSeenAt"] = data.LastSeenAt.ValueString()
    }
    if !data.AssignToUserId.IsNull() && !data.AssignToUserId.IsUnknown() {
        requestDataMap["assignToUserId"] = data.AssignToUserId.ValueString()
    }
    if !data.AssignToTeamId.IsNull() && !data.AssignToTeamId.IsUnknown() {
        requestDataMap["assignToTeamId"] = data.AssignToTeamId.ValueString()
    }
    if !data.IsResolved.IsNull() && !data.IsResolved.IsUnknown() {
        requestDataMap["isResolved"] = data.IsResolved.ValueBool()
    }
    if !data.IsArchived.IsNull() && !data.IsArchived.IsUnknown() {
        requestDataMap["isArchived"] = data.IsArchived.ValueBool()
    }
    if !data.OccuranceCount.IsNull() && !data.OccuranceCount.IsUnknown() {
        requestDataMap["occuranceCount"] = r.bigFloatToFloat64(data.OccuranceCount.ValueBigFloat())
    }
    if !data.FirstSeenInRelease.IsNull() && !data.FirstSeenInRelease.IsUnknown() {
        requestDataMap["firstSeenInRelease"] = data.FirstSeenInRelease.ValueString()
    }
    if !data.LastSeenInRelease.IsNull() && !data.LastSeenInRelease.IsUnknown() {
        requestDataMap["lastSeenInRelease"] = data.LastSeenInRelease.ValueString()
    }
    if !data.Environment.IsNull() && !data.Environment.IsUnknown() {
        requestDataMap["environment"] = data.Environment.ValueString()
    }

    // Make API call
    httpResp, err := r.client.Post(ctx, "/telemetry-exception", exceptionRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create exception, got error: %s", err))
        return
    }

    var exceptionResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &exceptionResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to create exception: %s", err))
        return
    }

    // Extract the new resource id from the create response.
    createdId := ""
    if wrapper, ok := exceptionResponse["data"].(map[string]interface{}); ok {
        if val, ok := wrapper["_id"].(string); ok {
            createdId = val
        }
    } else if val, ok := exceptionResponse["_id"].(string); ok {
        createdId = val
    }
    if createdId == "" {
        resp.Diagnostics.AddError("OneUptime API Error", "Create response for exception did not contain an id. This is a bug in the provider or the API; please report it.")
        return
    }
    data.Id = types.StringValue(createdId)

    // The person a create named by address, as the API has them.
    if data.AssignToUserId.IsUnknown() {
        created := exceptionResponse
        if wrapper, ok := exceptionResponse["data"].(map[string]interface{}); ok {
            created = wrapper
        }
        if id := apiStringValue(created["assignToUserId"]); id != "" {
            data.AssignToUserId = types.StringValue(id)
        } else {
            data.AssignToUserId = types.StringNull()
        }
    }

    /*
     * The server has committed the row. Persist what we know to state BEFORE
     * the read-back: if the read-back fails and we return without setting
     * state, Terraform never learns the resource exists and the created
     * exception is orphaned server-side — never refreshed, never
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
        "primaryEntityId": true,
        "primaryEntityType": true,
        "message": true,
        "stackTrace": true,
        "exceptionType": true,
        "fingerprint": true,
        "firstSeenAt": true,
        "lastSeenAt": true,
        "assignToUserId": true,
        "assignToTeamId": true,
        "isResolved": true,
        "isArchived": true,
        "occuranceCount": true,
        "firstSeenInRelease": true,
        "lastSeenInRelease": true,
        "environment": true,
        "regressedAt": true,
        "resolvedInReleaseAfter": true,
        "createdAt": true,
        "updatedAt": true,
        "createdByUserId": true,
        "markedAsResolvedAt": true,
        "markedAsArchivedAt": true,
        "markedAsResolvedByUserId": true,
        "markedAsArchivedByUserId": true,
        "unhandled": true,
        "aiClassification": true,
        "errorClass": true,
        "errorClassSource": true,
        "aiFixDeclinedAt": true,
        "aiExplanation": true,
        "aiExplanationGeneratedAt": true,
        "mergedIntoTelemetryExceptionId": true,
        "mergedAt": true,
        "mergedByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/telemetry-exception/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        /*
         * State already owns the id, so the resource is tracked and the next
         * refresh reconciles the remaining attributes. Warn rather than
         * error: erroring here would strand a real resource.
         */
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created exception but could not read it back; state is incomplete until the next refresh: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created exception but could not parse the read-back response; state is incomplete until the next refresh: %s", err))
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
    if obj, ok := dataMap["primaryEntityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PrimaryEntityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PrimaryEntityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PrimaryEntityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PrimaryEntityId = types.StringValue(string(jsonBytes))
            } else {
                data.PrimaryEntityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PrimaryEntityId = types.StringValue(string(jsonBytes))
            } else {
                data.PrimaryEntityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PrimaryEntityId = types.StringValue(string(jsonBytes))
        } else {
            data.PrimaryEntityId = types.StringNull()
        }
    } else if arr, ok := dataMap["primaryEntityId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.PrimaryEntityId = types.StringValue(string(jsonBytes))
        } else {
            data.PrimaryEntityId = types.StringNull()
        }
    } else if val, ok := dataMap["primaryEntityId"].(string); ok {
        data.PrimaryEntityId = types.StringValue(val)
    } else {
        data.PrimaryEntityId = types.StringNull()
    }
    if obj, ok := dataMap["primaryEntityType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PrimaryEntityType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PrimaryEntityType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PrimaryEntityType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PrimaryEntityType = types.StringValue(string(jsonBytes))
            } else {
                data.PrimaryEntityType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PrimaryEntityType = types.StringValue(string(jsonBytes))
            } else {
                data.PrimaryEntityType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PrimaryEntityType = types.StringValue(string(jsonBytes))
        } else {
            data.PrimaryEntityType = types.StringNull()
        }
    } else if arr, ok := dataMap["primaryEntityType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.PrimaryEntityType = types.StringValue(string(jsonBytes))
        } else {
            data.PrimaryEntityType = types.StringNull()
        }
    } else if val, ok := dataMap["primaryEntityType"].(string); ok {
        data.PrimaryEntityType = types.StringValue(val)
    } else {
        data.PrimaryEntityType = types.StringNull()
    }
    if obj, ok := dataMap["message"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Message = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Message = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Message = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Message = types.StringValue(string(jsonBytes))
            } else {
                data.Message = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Message = types.StringValue(string(jsonBytes))
            } else {
                data.Message = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Message = types.StringValue(string(jsonBytes))
        } else {
            data.Message = types.StringNull()
        }
    } else if arr, ok := dataMap["message"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Message = types.StringValue(string(jsonBytes))
        } else {
            data.Message = types.StringNull()
        }
    } else if val, ok := dataMap["message"].(string); ok {
        data.Message = types.StringValue(val)
    } else {
        data.Message = types.StringNull()
    }
    if obj, ok := dataMap["stackTrace"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StackTrace = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StackTrace = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StackTrace = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.StackTrace = types.StringValue(string(jsonBytes))
            } else {
                data.StackTrace = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.StackTrace = types.StringValue(string(jsonBytes))
            } else {
                data.StackTrace = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StackTrace = types.StringValue(string(jsonBytes))
        } else {
            data.StackTrace = types.StringNull()
        }
    } else if arr, ok := dataMap["stackTrace"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.StackTrace = types.StringValue(string(jsonBytes))
        } else {
            data.StackTrace = types.StringNull()
        }
    } else if val, ok := dataMap["stackTrace"].(string); ok {
        data.StackTrace = types.StringValue(val)
    } else {
        data.StackTrace = types.StringNull()
    }
    if obj, ok := dataMap["exceptionType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ExceptionType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ExceptionType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ExceptionType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ExceptionType = types.StringValue(string(jsonBytes))
            } else {
                data.ExceptionType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ExceptionType = types.StringValue(string(jsonBytes))
            } else {
                data.ExceptionType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ExceptionType = types.StringValue(string(jsonBytes))
        } else {
            data.ExceptionType = types.StringNull()
        }
    } else if arr, ok := dataMap["exceptionType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ExceptionType = types.StringValue(string(jsonBytes))
        } else {
            data.ExceptionType = types.StringNull()
        }
    } else if val, ok := dataMap["exceptionType"].(string); ok {
        data.ExceptionType = types.StringValue(val)
    } else {
        data.ExceptionType = types.StringNull()
    }
    if obj, ok := dataMap["fingerprint"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Fingerprint = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Fingerprint = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Fingerprint = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Fingerprint = types.StringValue(string(jsonBytes))
            } else {
                data.Fingerprint = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Fingerprint = types.StringValue(string(jsonBytes))
            } else {
                data.Fingerprint = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Fingerprint = types.StringValue(string(jsonBytes))
        } else {
            data.Fingerprint = types.StringNull()
        }
    } else if arr, ok := dataMap["fingerprint"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Fingerprint = types.StringValue(string(jsonBytes))
        } else {
            data.Fingerprint = types.StringNull()
        }
    } else if val, ok := dataMap["fingerprint"].(string); ok {
        data.Fingerprint = types.StringValue(val)
    } else {
        data.Fingerprint = types.StringNull()
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
    if obj, ok := dataMap["assignToUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AssignToUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AssignToUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AssignToUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AssignToUserId = types.StringValue(string(jsonBytes))
            } else {
                data.AssignToUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AssignToUserId = types.StringValue(string(jsonBytes))
            } else {
                data.AssignToUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AssignToUserId = types.StringValue(string(jsonBytes))
        } else {
            data.AssignToUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["assignToUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AssignToUserId = types.StringValue(string(jsonBytes))
        } else {
            data.AssignToUserId = types.StringNull()
        }
    } else if val, ok := dataMap["assignToUserId"].(string); ok {
        data.AssignToUserId = types.StringValue(val)
    } else {
        data.AssignToUserId = types.StringNull()
    }
    if obj, ok := dataMap["assignToTeamId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AssignToTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AssignToTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AssignToTeamId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AssignToTeamId = types.StringValue(string(jsonBytes))
            } else {
                data.AssignToTeamId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AssignToTeamId = types.StringValue(string(jsonBytes))
            } else {
                data.AssignToTeamId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AssignToTeamId = types.StringValue(string(jsonBytes))
        } else {
            data.AssignToTeamId = types.StringNull()
        }
    } else if arr, ok := dataMap["assignToTeamId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AssignToTeamId = types.StringValue(string(jsonBytes))
        } else {
            data.AssignToTeamId = types.StringNull()
        }
    } else if val, ok := dataMap["assignToTeamId"].(string); ok {
        data.AssignToTeamId = types.StringValue(val)
    } else {
        data.AssignToTeamId = types.StringNull()
    }
    if val, ok := dataMap["isResolved"].(bool); ok {
        data.IsResolved = types.BoolValue(val)
    }
    if val, ok := dataMap["isArchived"].(bool); ok {
        data.IsArchived = types.BoolValue(val)
    }
    data.OccuranceCount = numberFromAPI(data.OccuranceCount, dataMap["occuranceCount"])
    if obj, ok := dataMap["firstSeenInRelease"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FirstSeenInRelease = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.FirstSeenInRelease = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.FirstSeenInRelease = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.FirstSeenInRelease = types.StringValue(string(jsonBytes))
            } else {
                data.FirstSeenInRelease = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.FirstSeenInRelease = types.StringValue(string(jsonBytes))
            } else {
                data.FirstSeenInRelease = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.FirstSeenInRelease = types.StringValue(string(jsonBytes))
        } else {
            data.FirstSeenInRelease = types.StringNull()
        }
    } else if arr, ok := dataMap["firstSeenInRelease"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.FirstSeenInRelease = types.StringValue(string(jsonBytes))
        } else {
            data.FirstSeenInRelease = types.StringNull()
        }
    } else if val, ok := dataMap["firstSeenInRelease"].(string); ok {
        data.FirstSeenInRelease = types.StringValue(val)
    } else {
        data.FirstSeenInRelease = types.StringNull()
    }
    if obj, ok := dataMap["lastSeenInRelease"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastSeenInRelease = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastSeenInRelease = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastSeenInRelease = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastSeenInRelease = types.StringValue(string(jsonBytes))
            } else {
                data.LastSeenInRelease = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastSeenInRelease = types.StringValue(string(jsonBytes))
            } else {
                data.LastSeenInRelease = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastSeenInRelease = types.StringValue(string(jsonBytes))
        } else {
            data.LastSeenInRelease = types.StringNull()
        }
    } else if arr, ok := dataMap["lastSeenInRelease"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LastSeenInRelease = types.StringValue(string(jsonBytes))
        } else {
            data.LastSeenInRelease = types.StringNull()
        }
    } else if val, ok := dataMap["lastSeenInRelease"].(string); ok {
        data.LastSeenInRelease = types.StringValue(val)
    } else {
        data.LastSeenInRelease = types.StringNull()
    }
    if obj, ok := dataMap["environment"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Environment = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Environment = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Environment = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Environment = types.StringValue(string(jsonBytes))
            } else {
                data.Environment = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Environment = types.StringValue(string(jsonBytes))
            } else {
                data.Environment = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Environment = types.StringValue(string(jsonBytes))
        } else {
            data.Environment = types.StringNull()
        }
    } else if arr, ok := dataMap["environment"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Environment = types.StringValue(string(jsonBytes))
        } else {
            data.Environment = types.StringNull()
        }
    } else if val, ok := dataMap["environment"].(string); ok {
        data.Environment = types.StringValue(val)
    } else {
        data.Environment = types.StringNull()
    }
    if obj, ok := dataMap["regressedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.RegressedAt = NewRFC3339Value(val)
        } else {
            data.RegressedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["regressedAt"].(string); ok && val != "" {
        data.RegressedAt = NewRFC3339Value(val)
    } else {
        data.RegressedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["resolvedInReleaseAfter"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResolvedInReleaseAfter = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ResolvedInReleaseAfter = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ResolvedInReleaseAfter = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ResolvedInReleaseAfter = types.StringValue(string(jsonBytes))
            } else {
                data.ResolvedInReleaseAfter = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ResolvedInReleaseAfter = types.StringValue(string(jsonBytes))
            } else {
                data.ResolvedInReleaseAfter = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ResolvedInReleaseAfter = types.StringValue(string(jsonBytes))
        } else {
            data.ResolvedInReleaseAfter = types.StringNull()
        }
    } else if arr, ok := dataMap["resolvedInReleaseAfter"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ResolvedInReleaseAfter = types.StringValue(string(jsonBytes))
        } else {
            data.ResolvedInReleaseAfter = types.StringNull()
        }
    } else if val, ok := dataMap["resolvedInReleaseAfter"].(string); ok {
        data.ResolvedInReleaseAfter = types.StringValue(val)
    } else {
        data.ResolvedInReleaseAfter = types.StringNull()
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
    if obj, ok := dataMap["markedAsResolvedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.MarkedAsResolvedAt = NewRFC3339Value(val)
        } else {
            data.MarkedAsResolvedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["markedAsResolvedAt"].(string); ok && val != "" {
        data.MarkedAsResolvedAt = NewRFC3339Value(val)
    } else {
        data.MarkedAsResolvedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["markedAsArchivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.MarkedAsArchivedAt = NewRFC3339Value(val)
        } else {
            data.MarkedAsArchivedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["markedAsArchivedAt"].(string); ok && val != "" {
        data.MarkedAsArchivedAt = NewRFC3339Value(val)
    } else {
        data.MarkedAsArchivedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["markedAsResolvedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MarkedAsResolvedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MarkedAsResolvedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MarkedAsResolvedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MarkedAsResolvedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.MarkedAsResolvedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MarkedAsResolvedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.MarkedAsResolvedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MarkedAsResolvedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.MarkedAsResolvedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["markedAsResolvedByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.MarkedAsResolvedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.MarkedAsResolvedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["markedAsResolvedByUserId"].(string); ok {
        data.MarkedAsResolvedByUserId = types.StringValue(val)
    } else {
        data.MarkedAsResolvedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["markedAsArchivedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MarkedAsArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MarkedAsArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MarkedAsArchivedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MarkedAsArchivedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.MarkedAsArchivedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MarkedAsArchivedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.MarkedAsArchivedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MarkedAsArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.MarkedAsArchivedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["markedAsArchivedByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.MarkedAsArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.MarkedAsArchivedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["markedAsArchivedByUserId"].(string); ok {
        data.MarkedAsArchivedByUserId = types.StringValue(val)
    } else {
        data.MarkedAsArchivedByUserId = types.StringNull()
    }
    if val, ok := dataMap["unhandled"].(bool); ok {
        data.Unhandled = types.BoolValue(val)
    }
    if obj, ok := dataMap["aiClassification"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiClassification = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AiClassification = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AiClassification = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AiClassification = types.StringValue(string(jsonBytes))
            } else {
                data.AiClassification = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AiClassification = types.StringValue(string(jsonBytes))
            } else {
                data.AiClassification = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AiClassification = types.StringValue(string(jsonBytes))
        } else {
            data.AiClassification = types.StringNull()
        }
    } else if arr, ok := dataMap["aiClassification"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AiClassification = types.StringValue(string(jsonBytes))
        } else {
            data.AiClassification = types.StringNull()
        }
    } else if val, ok := dataMap["aiClassification"].(string); ok {
        data.AiClassification = types.StringValue(val)
    } else {
        data.AiClassification = types.StringNull()
    }
    if obj, ok := dataMap["errorClass"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ErrorClass = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ErrorClass = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ErrorClass = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ErrorClass = types.StringValue(string(jsonBytes))
            } else {
                data.ErrorClass = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ErrorClass = types.StringValue(string(jsonBytes))
            } else {
                data.ErrorClass = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ErrorClass = types.StringValue(string(jsonBytes))
        } else {
            data.ErrorClass = types.StringNull()
        }
    } else if arr, ok := dataMap["errorClass"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ErrorClass = types.StringValue(string(jsonBytes))
        } else {
            data.ErrorClass = types.StringNull()
        }
    } else if val, ok := dataMap["errorClass"].(string); ok {
        data.ErrorClass = types.StringValue(val)
    } else {
        data.ErrorClass = types.StringNull()
    }
    if obj, ok := dataMap["errorClassSource"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ErrorClassSource = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ErrorClassSource = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ErrorClassSource = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ErrorClassSource = types.StringValue(string(jsonBytes))
            } else {
                data.ErrorClassSource = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ErrorClassSource = types.StringValue(string(jsonBytes))
            } else {
                data.ErrorClassSource = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ErrorClassSource = types.StringValue(string(jsonBytes))
        } else {
            data.ErrorClassSource = types.StringNull()
        }
    } else if arr, ok := dataMap["errorClassSource"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ErrorClassSource = types.StringValue(string(jsonBytes))
        } else {
            data.ErrorClassSource = types.StringNull()
        }
    } else if val, ok := dataMap["errorClassSource"].(string); ok {
        data.ErrorClassSource = types.StringValue(val)
    } else {
        data.ErrorClassSource = types.StringNull()
    }
    if obj, ok := dataMap["aiFixDeclinedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.AiFixDeclinedAt = NewRFC3339Value(val)
        } else {
            data.AiFixDeclinedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["aiFixDeclinedAt"].(string); ok && val != "" {
        data.AiFixDeclinedAt = NewRFC3339Value(val)
    } else {
        data.AiFixDeclinedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["aiExplanation"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiExplanation = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AiExplanation = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AiExplanation = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AiExplanation = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.AiExplanation = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AiExplanation = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.AiExplanation = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AiExplanation = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.AiExplanation = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["aiExplanation"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AiExplanation = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.AiExplanation = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["aiExplanation"].(string); ok {
        data.AiExplanation = NewJSONSubsetValue(val)
    } else {
        data.AiExplanation = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["aiExplanationGeneratedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.AiExplanationGeneratedAt = NewRFC3339Value(val)
        } else {
            data.AiExplanationGeneratedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["aiExplanationGeneratedAt"].(string); ok && val != "" {
        data.AiExplanationGeneratedAt = NewRFC3339Value(val)
    } else {
        data.AiExplanationGeneratedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["mergedIntoTelemetryExceptionId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MergedIntoTelemetryExceptionId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MergedIntoTelemetryExceptionId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MergedIntoTelemetryExceptionId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MergedIntoTelemetryExceptionId = types.StringValue(string(jsonBytes))
            } else {
                data.MergedIntoTelemetryExceptionId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MergedIntoTelemetryExceptionId = types.StringValue(string(jsonBytes))
            } else {
                data.MergedIntoTelemetryExceptionId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MergedIntoTelemetryExceptionId = types.StringValue(string(jsonBytes))
        } else {
            data.MergedIntoTelemetryExceptionId = types.StringNull()
        }
    } else if arr, ok := dataMap["mergedIntoTelemetryExceptionId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.MergedIntoTelemetryExceptionId = types.StringValue(string(jsonBytes))
        } else {
            data.MergedIntoTelemetryExceptionId = types.StringNull()
        }
    } else if val, ok := dataMap["mergedIntoTelemetryExceptionId"].(string); ok {
        data.MergedIntoTelemetryExceptionId = types.StringValue(val)
    } else {
        data.MergedIntoTelemetryExceptionId = types.StringNull()
    }
    if obj, ok := dataMap["mergedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.MergedAt = NewRFC3339Value(val)
        } else {
            data.MergedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["mergedAt"].(string); ok && val != "" {
        data.MergedAt = NewRFC3339Value(val)
    } else {
        data.MergedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["mergedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MergedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MergedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MergedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MergedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.MergedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MergedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.MergedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MergedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.MergedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["mergedByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.MergedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.MergedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["mergedByUserId"].(string); ok {
        data.MergedByUserId = types.StringValue(val)
    } else {
        data.MergedByUserId = types.StringNull()
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
    // The address of a person named by ID (see people.go).
    data.AssignToUserEmail = personEmailAfterWrite(ctx, r.client, data.AssignToUserId, data.AssignToUserEmail)

    // Write logs using the tflog package
    tflog.Trace(ctx, "created a resource")


    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExceptionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data ExceptionResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "primaryEntityId": true,
        "primaryEntityType": true,
        "message": true,
        "stackTrace": true,
        "exceptionType": true,
        "fingerprint": true,
        "firstSeenAt": true,
        "lastSeenAt": true,
        "assignToUserId": true,
        "assignToTeamId": true,
        "isResolved": true,
        "isArchived": true,
        "occuranceCount": true,
        "firstSeenInRelease": true,
        "lastSeenInRelease": true,
        "environment": true,
        "regressedAt": true,
        "resolvedInReleaseAfter": true,
        "createdAt": true,
        "updatedAt": true,
        "createdByUserId": true,
        "markedAsResolvedAt": true,
        "markedAsArchivedAt": true,
        "markedAsResolvedByUserId": true,
        "markedAsArchivedByUserId": true,
        "unhandled": true,
        "aiClassification": true,
        "errorClass": true,
        "errorClassSource": true,
        "aiFixDeclinedAt": true,
        "aiExplanation": true,
        "aiExplanationGeneratedAt": true,
        "mergedIntoTelemetryExceptionId": true,
        "mergedAt": true,
        "mergedByUserId": true,
        "_id": true,
    }

    // Who the state named, to keep their address (see people.go).
    priorAssignToUserId := data.AssignToUserId
    priorAssignToUserEmail := data.AssignToUserEmail

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect(ctx, "/telemetry-exception/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read exception, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var exceptionResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &exceptionResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse exception response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := exceptionResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = exceptionResponse
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
    if obj, ok := dataMap["primaryEntityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PrimaryEntityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PrimaryEntityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PrimaryEntityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PrimaryEntityId = types.StringValue(string(jsonBytes))
            } else {
                data.PrimaryEntityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PrimaryEntityId = types.StringValue(string(jsonBytes))
            } else {
                data.PrimaryEntityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PrimaryEntityId = types.StringValue(string(jsonBytes))
        } else {
            data.PrimaryEntityId = types.StringNull()
        }
    } else if arr, ok := dataMap["primaryEntityId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.PrimaryEntityId = types.StringValue(string(jsonBytes))
        } else {
            data.PrimaryEntityId = types.StringNull()
        }
    } else if val, ok := dataMap["primaryEntityId"].(string); ok {
        data.PrimaryEntityId = types.StringValue(val)
    } else {
        data.PrimaryEntityId = types.StringNull()
    }
    if obj, ok := dataMap["primaryEntityType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PrimaryEntityType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PrimaryEntityType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PrimaryEntityType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PrimaryEntityType = types.StringValue(string(jsonBytes))
            } else {
                data.PrimaryEntityType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PrimaryEntityType = types.StringValue(string(jsonBytes))
            } else {
                data.PrimaryEntityType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PrimaryEntityType = types.StringValue(string(jsonBytes))
        } else {
            data.PrimaryEntityType = types.StringNull()
        }
    } else if arr, ok := dataMap["primaryEntityType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.PrimaryEntityType = types.StringValue(string(jsonBytes))
        } else {
            data.PrimaryEntityType = types.StringNull()
        }
    } else if val, ok := dataMap["primaryEntityType"].(string); ok {
        data.PrimaryEntityType = types.StringValue(val)
    } else {
        data.PrimaryEntityType = types.StringNull()
    }
    if obj, ok := dataMap["message"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Message = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Message = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Message = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Message = types.StringValue(string(jsonBytes))
            } else {
                data.Message = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Message = types.StringValue(string(jsonBytes))
            } else {
                data.Message = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Message = types.StringValue(string(jsonBytes))
        } else {
            data.Message = types.StringNull()
        }
    } else if arr, ok := dataMap["message"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Message = types.StringValue(string(jsonBytes))
        } else {
            data.Message = types.StringNull()
        }
    } else if val, ok := dataMap["message"].(string); ok {
        data.Message = types.StringValue(val)
    } else {
        data.Message = types.StringNull()
    }
    if obj, ok := dataMap["stackTrace"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StackTrace = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StackTrace = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StackTrace = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.StackTrace = types.StringValue(string(jsonBytes))
            } else {
                data.StackTrace = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.StackTrace = types.StringValue(string(jsonBytes))
            } else {
                data.StackTrace = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StackTrace = types.StringValue(string(jsonBytes))
        } else {
            data.StackTrace = types.StringNull()
        }
    } else if arr, ok := dataMap["stackTrace"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.StackTrace = types.StringValue(string(jsonBytes))
        } else {
            data.StackTrace = types.StringNull()
        }
    } else if val, ok := dataMap["stackTrace"].(string); ok {
        data.StackTrace = types.StringValue(val)
    } else {
        data.StackTrace = types.StringNull()
    }
    if obj, ok := dataMap["exceptionType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ExceptionType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ExceptionType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ExceptionType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ExceptionType = types.StringValue(string(jsonBytes))
            } else {
                data.ExceptionType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ExceptionType = types.StringValue(string(jsonBytes))
            } else {
                data.ExceptionType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ExceptionType = types.StringValue(string(jsonBytes))
        } else {
            data.ExceptionType = types.StringNull()
        }
    } else if arr, ok := dataMap["exceptionType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ExceptionType = types.StringValue(string(jsonBytes))
        } else {
            data.ExceptionType = types.StringNull()
        }
    } else if val, ok := dataMap["exceptionType"].(string); ok {
        data.ExceptionType = types.StringValue(val)
    } else {
        data.ExceptionType = types.StringNull()
    }
    if obj, ok := dataMap["fingerprint"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Fingerprint = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Fingerprint = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Fingerprint = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Fingerprint = types.StringValue(string(jsonBytes))
            } else {
                data.Fingerprint = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Fingerprint = types.StringValue(string(jsonBytes))
            } else {
                data.Fingerprint = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Fingerprint = types.StringValue(string(jsonBytes))
        } else {
            data.Fingerprint = types.StringNull()
        }
    } else if arr, ok := dataMap["fingerprint"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Fingerprint = types.StringValue(string(jsonBytes))
        } else {
            data.Fingerprint = types.StringNull()
        }
    } else if val, ok := dataMap["fingerprint"].(string); ok {
        data.Fingerprint = types.StringValue(val)
    } else {
        data.Fingerprint = types.StringNull()
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
    if obj, ok := dataMap["assignToUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AssignToUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AssignToUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AssignToUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AssignToUserId = types.StringValue(string(jsonBytes))
            } else {
                data.AssignToUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AssignToUserId = types.StringValue(string(jsonBytes))
            } else {
                data.AssignToUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AssignToUserId = types.StringValue(string(jsonBytes))
        } else {
            data.AssignToUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["assignToUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AssignToUserId = types.StringValue(string(jsonBytes))
        } else {
            data.AssignToUserId = types.StringNull()
        }
    } else if val, ok := dataMap["assignToUserId"].(string); ok {
        data.AssignToUserId = types.StringValue(val)
    } else {
        data.AssignToUserId = types.StringNull()
    }
    if obj, ok := dataMap["assignToTeamId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AssignToTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AssignToTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AssignToTeamId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AssignToTeamId = types.StringValue(string(jsonBytes))
            } else {
                data.AssignToTeamId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AssignToTeamId = types.StringValue(string(jsonBytes))
            } else {
                data.AssignToTeamId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AssignToTeamId = types.StringValue(string(jsonBytes))
        } else {
            data.AssignToTeamId = types.StringNull()
        }
    } else if arr, ok := dataMap["assignToTeamId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AssignToTeamId = types.StringValue(string(jsonBytes))
        } else {
            data.AssignToTeamId = types.StringNull()
        }
    } else if val, ok := dataMap["assignToTeamId"].(string); ok {
        data.AssignToTeamId = types.StringValue(val)
    } else {
        data.AssignToTeamId = types.StringNull()
    }
    if val, ok := dataMap["isResolved"].(bool); ok {
        data.IsResolved = types.BoolValue(val)
    }
    if val, ok := dataMap["isArchived"].(bool); ok {
        data.IsArchived = types.BoolValue(val)
    }
    data.OccuranceCount = numberFromAPI(data.OccuranceCount, dataMap["occuranceCount"])
    if obj, ok := dataMap["firstSeenInRelease"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FirstSeenInRelease = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.FirstSeenInRelease = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.FirstSeenInRelease = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.FirstSeenInRelease = types.StringValue(string(jsonBytes))
            } else {
                data.FirstSeenInRelease = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.FirstSeenInRelease = types.StringValue(string(jsonBytes))
            } else {
                data.FirstSeenInRelease = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.FirstSeenInRelease = types.StringValue(string(jsonBytes))
        } else {
            data.FirstSeenInRelease = types.StringNull()
        }
    } else if arr, ok := dataMap["firstSeenInRelease"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.FirstSeenInRelease = types.StringValue(string(jsonBytes))
        } else {
            data.FirstSeenInRelease = types.StringNull()
        }
    } else if val, ok := dataMap["firstSeenInRelease"].(string); ok {
        data.FirstSeenInRelease = types.StringValue(val)
    } else {
        data.FirstSeenInRelease = types.StringNull()
    }
    if obj, ok := dataMap["lastSeenInRelease"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastSeenInRelease = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastSeenInRelease = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastSeenInRelease = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastSeenInRelease = types.StringValue(string(jsonBytes))
            } else {
                data.LastSeenInRelease = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastSeenInRelease = types.StringValue(string(jsonBytes))
            } else {
                data.LastSeenInRelease = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastSeenInRelease = types.StringValue(string(jsonBytes))
        } else {
            data.LastSeenInRelease = types.StringNull()
        }
    } else if arr, ok := dataMap["lastSeenInRelease"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LastSeenInRelease = types.StringValue(string(jsonBytes))
        } else {
            data.LastSeenInRelease = types.StringNull()
        }
    } else if val, ok := dataMap["lastSeenInRelease"].(string); ok {
        data.LastSeenInRelease = types.StringValue(val)
    } else {
        data.LastSeenInRelease = types.StringNull()
    }
    if obj, ok := dataMap["environment"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Environment = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Environment = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Environment = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Environment = types.StringValue(string(jsonBytes))
            } else {
                data.Environment = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Environment = types.StringValue(string(jsonBytes))
            } else {
                data.Environment = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Environment = types.StringValue(string(jsonBytes))
        } else {
            data.Environment = types.StringNull()
        }
    } else if arr, ok := dataMap["environment"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Environment = types.StringValue(string(jsonBytes))
        } else {
            data.Environment = types.StringNull()
        }
    } else if val, ok := dataMap["environment"].(string); ok {
        data.Environment = types.StringValue(val)
    } else {
        data.Environment = types.StringNull()
    }
    if obj, ok := dataMap["regressedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.RegressedAt = NewRFC3339Value(val)
        } else {
            data.RegressedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["regressedAt"].(string); ok && val != "" {
        data.RegressedAt = NewRFC3339Value(val)
    } else {
        data.RegressedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["resolvedInReleaseAfter"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResolvedInReleaseAfter = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ResolvedInReleaseAfter = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ResolvedInReleaseAfter = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ResolvedInReleaseAfter = types.StringValue(string(jsonBytes))
            } else {
                data.ResolvedInReleaseAfter = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ResolvedInReleaseAfter = types.StringValue(string(jsonBytes))
            } else {
                data.ResolvedInReleaseAfter = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ResolvedInReleaseAfter = types.StringValue(string(jsonBytes))
        } else {
            data.ResolvedInReleaseAfter = types.StringNull()
        }
    } else if arr, ok := dataMap["resolvedInReleaseAfter"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ResolvedInReleaseAfter = types.StringValue(string(jsonBytes))
        } else {
            data.ResolvedInReleaseAfter = types.StringNull()
        }
    } else if val, ok := dataMap["resolvedInReleaseAfter"].(string); ok {
        data.ResolvedInReleaseAfter = types.StringValue(val)
    } else {
        data.ResolvedInReleaseAfter = types.StringNull()
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
    if obj, ok := dataMap["markedAsResolvedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.MarkedAsResolvedAt = NewRFC3339Value(val)
        } else {
            data.MarkedAsResolvedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["markedAsResolvedAt"].(string); ok && val != "" {
        data.MarkedAsResolvedAt = NewRFC3339Value(val)
    } else {
        data.MarkedAsResolvedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["markedAsArchivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.MarkedAsArchivedAt = NewRFC3339Value(val)
        } else {
            data.MarkedAsArchivedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["markedAsArchivedAt"].(string); ok && val != "" {
        data.MarkedAsArchivedAt = NewRFC3339Value(val)
    } else {
        data.MarkedAsArchivedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["markedAsResolvedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MarkedAsResolvedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MarkedAsResolvedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MarkedAsResolvedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MarkedAsResolvedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.MarkedAsResolvedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MarkedAsResolvedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.MarkedAsResolvedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MarkedAsResolvedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.MarkedAsResolvedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["markedAsResolvedByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.MarkedAsResolvedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.MarkedAsResolvedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["markedAsResolvedByUserId"].(string); ok {
        data.MarkedAsResolvedByUserId = types.StringValue(val)
    } else {
        data.MarkedAsResolvedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["markedAsArchivedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MarkedAsArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MarkedAsArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MarkedAsArchivedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MarkedAsArchivedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.MarkedAsArchivedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MarkedAsArchivedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.MarkedAsArchivedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MarkedAsArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.MarkedAsArchivedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["markedAsArchivedByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.MarkedAsArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.MarkedAsArchivedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["markedAsArchivedByUserId"].(string); ok {
        data.MarkedAsArchivedByUserId = types.StringValue(val)
    } else {
        data.MarkedAsArchivedByUserId = types.StringNull()
    }
    if val, ok := dataMap["unhandled"].(bool); ok {
        data.Unhandled = types.BoolValue(val)
    }
    if obj, ok := dataMap["aiClassification"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiClassification = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AiClassification = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AiClassification = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AiClassification = types.StringValue(string(jsonBytes))
            } else {
                data.AiClassification = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AiClassification = types.StringValue(string(jsonBytes))
            } else {
                data.AiClassification = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AiClassification = types.StringValue(string(jsonBytes))
        } else {
            data.AiClassification = types.StringNull()
        }
    } else if arr, ok := dataMap["aiClassification"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AiClassification = types.StringValue(string(jsonBytes))
        } else {
            data.AiClassification = types.StringNull()
        }
    } else if val, ok := dataMap["aiClassification"].(string); ok {
        data.AiClassification = types.StringValue(val)
    } else {
        data.AiClassification = types.StringNull()
    }
    if obj, ok := dataMap["errorClass"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ErrorClass = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ErrorClass = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ErrorClass = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ErrorClass = types.StringValue(string(jsonBytes))
            } else {
                data.ErrorClass = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ErrorClass = types.StringValue(string(jsonBytes))
            } else {
                data.ErrorClass = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ErrorClass = types.StringValue(string(jsonBytes))
        } else {
            data.ErrorClass = types.StringNull()
        }
    } else if arr, ok := dataMap["errorClass"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ErrorClass = types.StringValue(string(jsonBytes))
        } else {
            data.ErrorClass = types.StringNull()
        }
    } else if val, ok := dataMap["errorClass"].(string); ok {
        data.ErrorClass = types.StringValue(val)
    } else {
        data.ErrorClass = types.StringNull()
    }
    if obj, ok := dataMap["errorClassSource"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ErrorClassSource = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ErrorClassSource = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ErrorClassSource = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ErrorClassSource = types.StringValue(string(jsonBytes))
            } else {
                data.ErrorClassSource = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ErrorClassSource = types.StringValue(string(jsonBytes))
            } else {
                data.ErrorClassSource = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ErrorClassSource = types.StringValue(string(jsonBytes))
        } else {
            data.ErrorClassSource = types.StringNull()
        }
    } else if arr, ok := dataMap["errorClassSource"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ErrorClassSource = types.StringValue(string(jsonBytes))
        } else {
            data.ErrorClassSource = types.StringNull()
        }
    } else if val, ok := dataMap["errorClassSource"].(string); ok {
        data.ErrorClassSource = types.StringValue(val)
    } else {
        data.ErrorClassSource = types.StringNull()
    }
    if obj, ok := dataMap["aiFixDeclinedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.AiFixDeclinedAt = NewRFC3339Value(val)
        } else {
            data.AiFixDeclinedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["aiFixDeclinedAt"].(string); ok && val != "" {
        data.AiFixDeclinedAt = NewRFC3339Value(val)
    } else {
        data.AiFixDeclinedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["aiExplanation"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiExplanation = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AiExplanation = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AiExplanation = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AiExplanation = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.AiExplanation = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AiExplanation = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.AiExplanation = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AiExplanation = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.AiExplanation = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["aiExplanation"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AiExplanation = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.AiExplanation = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["aiExplanation"].(string); ok {
        data.AiExplanation = NewJSONSubsetValue(val)
    } else {
        data.AiExplanation = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["aiExplanationGeneratedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.AiExplanationGeneratedAt = NewRFC3339Value(val)
        } else {
            data.AiExplanationGeneratedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["aiExplanationGeneratedAt"].(string); ok && val != "" {
        data.AiExplanationGeneratedAt = NewRFC3339Value(val)
    } else {
        data.AiExplanationGeneratedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["mergedIntoTelemetryExceptionId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MergedIntoTelemetryExceptionId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MergedIntoTelemetryExceptionId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MergedIntoTelemetryExceptionId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MergedIntoTelemetryExceptionId = types.StringValue(string(jsonBytes))
            } else {
                data.MergedIntoTelemetryExceptionId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MergedIntoTelemetryExceptionId = types.StringValue(string(jsonBytes))
            } else {
                data.MergedIntoTelemetryExceptionId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MergedIntoTelemetryExceptionId = types.StringValue(string(jsonBytes))
        } else {
            data.MergedIntoTelemetryExceptionId = types.StringNull()
        }
    } else if arr, ok := dataMap["mergedIntoTelemetryExceptionId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.MergedIntoTelemetryExceptionId = types.StringValue(string(jsonBytes))
        } else {
            data.MergedIntoTelemetryExceptionId = types.StringNull()
        }
    } else if val, ok := dataMap["mergedIntoTelemetryExceptionId"].(string); ok {
        data.MergedIntoTelemetryExceptionId = types.StringValue(val)
    } else {
        data.MergedIntoTelemetryExceptionId = types.StringNull()
    }
    if obj, ok := dataMap["mergedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.MergedAt = NewRFC3339Value(val)
        } else {
            data.MergedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["mergedAt"].(string); ok && val != "" {
        data.MergedAt = NewRFC3339Value(val)
    } else {
        data.MergedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["mergedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MergedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MergedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MergedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MergedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.MergedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MergedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.MergedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MergedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.MergedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["mergedByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.MergedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.MergedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["mergedByUserId"].(string); ok {
        data.MergedByUserId = types.StringValue(val)
    } else {
        data.MergedByUserId = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    // The person's address (see people.go).
    data.AssignToUserEmail = personEmailAfterRead(ctx, r.client, priorAssignToUserId, data.AssignToUserId, priorAssignToUserEmail)

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExceptionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data ExceptionResourceModel
    var state ExceptionResourceModel

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
    var config ExceptionResourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
    if resp.Diagnostics.HasError() {
        return
    }
    plan := data

    // People named by email address: the user IDs to send (see people.go).
    if !personIDForWrite(ctx, r.client, &data.AssignToUserId, data.AssignToUserEmail, "assign_to_user_email", &resp.Diagnostics) {
        return
    }

    // Create API request body
    exceptionRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := exceptionRequest["data"].(map[string]interface{})

    if !data.PrimaryEntityId.IsUnknown() && !state.PrimaryEntityId.IsUnknown() && !data.PrimaryEntityId.Equal(state.PrimaryEntityId) {
        requestDataMap["primaryEntityId"] = data.PrimaryEntityId.ValueString()
    }
    if !data.Message.IsUnknown() && !state.Message.IsUnknown() && !data.Message.Equal(state.Message) {
        requestDataMap["message"] = data.Message.ValueString()
    }
    if !data.StackTrace.IsUnknown() && !state.StackTrace.IsUnknown() && !data.StackTrace.Equal(state.StackTrace) {
        requestDataMap["stackTrace"] = data.StackTrace.ValueString()
    }
    if !data.ExceptionType.IsUnknown() && !state.ExceptionType.IsUnknown() && !data.ExceptionType.Equal(state.ExceptionType) {
        requestDataMap["exceptionType"] = data.ExceptionType.ValueString()
    }
    if !data.Fingerprint.IsUnknown() && !state.Fingerprint.IsUnknown() && !data.Fingerprint.Equal(state.Fingerprint) {
        requestDataMap["fingerprint"] = data.Fingerprint.ValueString()
    }
    if !data.FirstSeenAt.IsUnknown() && !state.FirstSeenAt.IsUnknown() && !data.FirstSeenAt.Equal(state.FirstSeenAt) {
        requestDataMap["firstSeenAt"] = data.FirstSeenAt.ValueString()
    }
    if !data.LastSeenAt.IsUnknown() && !state.LastSeenAt.IsUnknown() && !data.LastSeenAt.Equal(state.LastSeenAt) {
        requestDataMap["lastSeenAt"] = data.LastSeenAt.ValueString()
    }
    if !data.AssignToUserId.IsUnknown() && !state.AssignToUserId.IsUnknown() && !data.AssignToUserId.Equal(state.AssignToUserId) {
        requestDataMap["assignToUserId"] = data.AssignToUserId.ValueString()
    }
    if !data.AssignToTeamId.IsUnknown() && !state.AssignToTeamId.IsUnknown() && !data.AssignToTeamId.Equal(state.AssignToTeamId) {
        requestDataMap["assignToTeamId"] = data.AssignToTeamId.ValueString()
    }
    if !data.IsResolved.IsUnknown() && !state.IsResolved.IsUnknown() && !data.IsResolved.Equal(state.IsResolved) {
        requestDataMap["isResolved"] = data.IsResolved.ValueBool()
    }
    if !data.IsArchived.IsUnknown() && !state.IsArchived.IsUnknown() && !data.IsArchived.Equal(state.IsArchived) {
        requestDataMap["isArchived"] = data.IsArchived.ValueBool()
    }
    if !data.OccuranceCount.IsUnknown() && !state.OccuranceCount.IsUnknown() && !data.OccuranceCount.Equal(state.OccuranceCount) {
        requestDataMap["occuranceCount"] = r.bigFloatToFloat64(data.OccuranceCount.ValueBigFloat())
    }
    if !data.FirstSeenInRelease.IsUnknown() && !state.FirstSeenInRelease.IsUnknown() && !data.FirstSeenInRelease.Equal(state.FirstSeenInRelease) {
        requestDataMap["firstSeenInRelease"] = data.FirstSeenInRelease.ValueString()
    }
    if !data.LastSeenInRelease.IsUnknown() && !state.LastSeenInRelease.IsUnknown() && !data.LastSeenInRelease.Equal(state.LastSeenInRelease) {
        requestDataMap["lastSeenInRelease"] = data.LastSeenInRelease.ValueString()
    }
    if !data.Environment.IsUnknown() && !state.Environment.IsUnknown() && !data.Environment.Equal(state.Environment) {
        requestDataMap["environment"] = data.Environment.ValueString()
    }
    if !data.RegressedAt.IsUnknown() && !state.RegressedAt.IsUnknown() && !data.RegressedAt.Equal(state.RegressedAt) {
        requestDataMap["regressedAt"] = data.RegressedAt.ValueString()
    }
    if !data.ResolvedInReleaseAfter.IsUnknown() && !state.ResolvedInReleaseAfter.IsUnknown() && !data.ResolvedInReleaseAfter.Equal(state.ResolvedInReleaseAfter) {
        requestDataMap["resolvedInReleaseAfter"] = data.ResolvedInReleaseAfter.ValueString()
    }

    // Only call the API when there are changed fields to send. An empty
    // update body is rejected by the API; state is still refreshed below so
    // this method never writes unverified plan values into state.
    if len(exceptionRequest["data"].(map[string]interface{})) > 0 {
        httpResp, err := r.client.Put(ctx, "/telemetry-exception/" + data.Id.ValueString() + "", exceptionRequest)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update exception, got error: %s", err))
            return
        }

        // Parse the update response
        var exceptionResponse map[string]interface{}
        err = r.client.ParseResponse(httpResp, &exceptionResponse)
        if err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to update exception: %s", err))
            return
        }
        _ = exceptionResponse
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "primaryEntityId": true,
        "primaryEntityType": true,
        "message": true,
        "stackTrace": true,
        "exceptionType": true,
        "fingerprint": true,
        "firstSeenAt": true,
        "lastSeenAt": true,
        "assignToUserId": true,
        "assignToTeamId": true,
        "isResolved": true,
        "isArchived": true,
        "occuranceCount": true,
        "firstSeenInRelease": true,
        "lastSeenInRelease": true,
        "environment": true,
        "regressedAt": true,
        "resolvedInReleaseAfter": true,
        "createdAt": true,
        "updatedAt": true,
        "createdByUserId": true,
        "markedAsResolvedAt": true,
        "markedAsArchivedAt": true,
        "markedAsResolvedByUserId": true,
        "markedAsArchivedByUserId": true,
        "unhandled": true,
        "aiClassification": true,
        "errorClass": true,
        "errorClassSource": true,
        "aiFixDeclinedAt": true,
        "aiExplanation": true,
        "aiExplanationGeneratedAt": true,
        "mergedIntoTelemetryExceptionId": true,
        "mergedAt": true,
        "mergedByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/telemetry-exception/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read exception after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read exception after update: %s", err))
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
    if obj, ok := dataMap["primaryEntityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PrimaryEntityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PrimaryEntityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PrimaryEntityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PrimaryEntityId = types.StringValue(string(jsonBytes))
            } else {
                data.PrimaryEntityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PrimaryEntityId = types.StringValue(string(jsonBytes))
            } else {
                data.PrimaryEntityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PrimaryEntityId = types.StringValue(string(jsonBytes))
        } else {
            data.PrimaryEntityId = types.StringNull()
        }
    } else if arr, ok := dataMap["primaryEntityId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.PrimaryEntityId = types.StringValue(string(jsonBytes))
        } else {
            data.PrimaryEntityId = types.StringNull()
        }
    } else if val, ok := dataMap["primaryEntityId"].(string); ok {
        data.PrimaryEntityId = types.StringValue(val)
    } else {
        data.PrimaryEntityId = types.StringNull()
    }
    if obj, ok := dataMap["primaryEntityType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PrimaryEntityType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PrimaryEntityType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PrimaryEntityType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PrimaryEntityType = types.StringValue(string(jsonBytes))
            } else {
                data.PrimaryEntityType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PrimaryEntityType = types.StringValue(string(jsonBytes))
            } else {
                data.PrimaryEntityType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PrimaryEntityType = types.StringValue(string(jsonBytes))
        } else {
            data.PrimaryEntityType = types.StringNull()
        }
    } else if arr, ok := dataMap["primaryEntityType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.PrimaryEntityType = types.StringValue(string(jsonBytes))
        } else {
            data.PrimaryEntityType = types.StringNull()
        }
    } else if val, ok := dataMap["primaryEntityType"].(string); ok {
        data.PrimaryEntityType = types.StringValue(val)
    } else {
        data.PrimaryEntityType = types.StringNull()
    }
    if obj, ok := dataMap["message"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Message = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Message = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Message = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Message = types.StringValue(string(jsonBytes))
            } else {
                data.Message = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Message = types.StringValue(string(jsonBytes))
            } else {
                data.Message = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Message = types.StringValue(string(jsonBytes))
        } else {
            data.Message = types.StringNull()
        }
    } else if arr, ok := dataMap["message"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Message = types.StringValue(string(jsonBytes))
        } else {
            data.Message = types.StringNull()
        }
    } else if val, ok := dataMap["message"].(string); ok {
        data.Message = types.StringValue(val)
    } else {
        data.Message = types.StringNull()
    }
    if obj, ok := dataMap["stackTrace"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StackTrace = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StackTrace = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StackTrace = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.StackTrace = types.StringValue(string(jsonBytes))
            } else {
                data.StackTrace = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.StackTrace = types.StringValue(string(jsonBytes))
            } else {
                data.StackTrace = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StackTrace = types.StringValue(string(jsonBytes))
        } else {
            data.StackTrace = types.StringNull()
        }
    } else if arr, ok := dataMap["stackTrace"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.StackTrace = types.StringValue(string(jsonBytes))
        } else {
            data.StackTrace = types.StringNull()
        }
    } else if val, ok := dataMap["stackTrace"].(string); ok {
        data.StackTrace = types.StringValue(val)
    } else {
        data.StackTrace = types.StringNull()
    }
    if obj, ok := dataMap["exceptionType"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ExceptionType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ExceptionType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ExceptionType = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ExceptionType = types.StringValue(string(jsonBytes))
            } else {
                data.ExceptionType = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ExceptionType = types.StringValue(string(jsonBytes))
            } else {
                data.ExceptionType = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ExceptionType = types.StringValue(string(jsonBytes))
        } else {
            data.ExceptionType = types.StringNull()
        }
    } else if arr, ok := dataMap["exceptionType"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ExceptionType = types.StringValue(string(jsonBytes))
        } else {
            data.ExceptionType = types.StringNull()
        }
    } else if val, ok := dataMap["exceptionType"].(string); ok {
        data.ExceptionType = types.StringValue(val)
    } else {
        data.ExceptionType = types.StringNull()
    }
    if obj, ok := dataMap["fingerprint"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Fingerprint = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Fingerprint = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Fingerprint = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Fingerprint = types.StringValue(string(jsonBytes))
            } else {
                data.Fingerprint = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Fingerprint = types.StringValue(string(jsonBytes))
            } else {
                data.Fingerprint = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Fingerprint = types.StringValue(string(jsonBytes))
        } else {
            data.Fingerprint = types.StringNull()
        }
    } else if arr, ok := dataMap["fingerprint"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Fingerprint = types.StringValue(string(jsonBytes))
        } else {
            data.Fingerprint = types.StringNull()
        }
    } else if val, ok := dataMap["fingerprint"].(string); ok {
        data.Fingerprint = types.StringValue(val)
    } else {
        data.Fingerprint = types.StringNull()
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
    if obj, ok := dataMap["assignToUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AssignToUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AssignToUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AssignToUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AssignToUserId = types.StringValue(string(jsonBytes))
            } else {
                data.AssignToUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AssignToUserId = types.StringValue(string(jsonBytes))
            } else {
                data.AssignToUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AssignToUserId = types.StringValue(string(jsonBytes))
        } else {
            data.AssignToUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["assignToUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AssignToUserId = types.StringValue(string(jsonBytes))
        } else {
            data.AssignToUserId = types.StringNull()
        }
    } else if val, ok := dataMap["assignToUserId"].(string); ok {
        data.AssignToUserId = types.StringValue(val)
    } else {
        data.AssignToUserId = types.StringNull()
    }
    if obj, ok := dataMap["assignToTeamId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AssignToTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AssignToTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AssignToTeamId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AssignToTeamId = types.StringValue(string(jsonBytes))
            } else {
                data.AssignToTeamId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AssignToTeamId = types.StringValue(string(jsonBytes))
            } else {
                data.AssignToTeamId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AssignToTeamId = types.StringValue(string(jsonBytes))
        } else {
            data.AssignToTeamId = types.StringNull()
        }
    } else if arr, ok := dataMap["assignToTeamId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AssignToTeamId = types.StringValue(string(jsonBytes))
        } else {
            data.AssignToTeamId = types.StringNull()
        }
    } else if val, ok := dataMap["assignToTeamId"].(string); ok {
        data.AssignToTeamId = types.StringValue(val)
    } else {
        data.AssignToTeamId = types.StringNull()
    }
    if val, ok := dataMap["isResolved"].(bool); ok {
        data.IsResolved = types.BoolValue(val)
    }
    if val, ok := dataMap["isArchived"].(bool); ok {
        data.IsArchived = types.BoolValue(val)
    }
    data.OccuranceCount = numberFromAPI(data.OccuranceCount, dataMap["occuranceCount"])
    if obj, ok := dataMap["firstSeenInRelease"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FirstSeenInRelease = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.FirstSeenInRelease = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.FirstSeenInRelease = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.FirstSeenInRelease = types.StringValue(string(jsonBytes))
            } else {
                data.FirstSeenInRelease = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.FirstSeenInRelease = types.StringValue(string(jsonBytes))
            } else {
                data.FirstSeenInRelease = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.FirstSeenInRelease = types.StringValue(string(jsonBytes))
        } else {
            data.FirstSeenInRelease = types.StringNull()
        }
    } else if arr, ok := dataMap["firstSeenInRelease"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.FirstSeenInRelease = types.StringValue(string(jsonBytes))
        } else {
            data.FirstSeenInRelease = types.StringNull()
        }
    } else if val, ok := dataMap["firstSeenInRelease"].(string); ok {
        data.FirstSeenInRelease = types.StringValue(val)
    } else {
        data.FirstSeenInRelease = types.StringNull()
    }
    if obj, ok := dataMap["lastSeenInRelease"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastSeenInRelease = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.LastSeenInRelease = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.LastSeenInRelease = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.LastSeenInRelease = types.StringValue(string(jsonBytes))
            } else {
                data.LastSeenInRelease = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.LastSeenInRelease = types.StringValue(string(jsonBytes))
            } else {
                data.LastSeenInRelease = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.LastSeenInRelease = types.StringValue(string(jsonBytes))
        } else {
            data.LastSeenInRelease = types.StringNull()
        }
    } else if arr, ok := dataMap["lastSeenInRelease"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.LastSeenInRelease = types.StringValue(string(jsonBytes))
        } else {
            data.LastSeenInRelease = types.StringNull()
        }
    } else if val, ok := dataMap["lastSeenInRelease"].(string); ok {
        data.LastSeenInRelease = types.StringValue(val)
    } else {
        data.LastSeenInRelease = types.StringNull()
    }
    if obj, ok := dataMap["environment"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Environment = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Environment = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Environment = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Environment = types.StringValue(string(jsonBytes))
            } else {
                data.Environment = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Environment = types.StringValue(string(jsonBytes))
            } else {
                data.Environment = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Environment = types.StringValue(string(jsonBytes))
        } else {
            data.Environment = types.StringNull()
        }
    } else if arr, ok := dataMap["environment"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.Environment = types.StringValue(string(jsonBytes))
        } else {
            data.Environment = types.StringNull()
        }
    } else if val, ok := dataMap["environment"].(string); ok {
        data.Environment = types.StringValue(val)
    } else {
        data.Environment = types.StringNull()
    }
    if obj, ok := dataMap["regressedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.RegressedAt = NewRFC3339Value(val)
        } else {
            data.RegressedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["regressedAt"].(string); ok && val != "" {
        data.RegressedAt = NewRFC3339Value(val)
    } else {
        data.RegressedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["resolvedInReleaseAfter"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResolvedInReleaseAfter = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ResolvedInReleaseAfter = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ResolvedInReleaseAfter = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ResolvedInReleaseAfter = types.StringValue(string(jsonBytes))
            } else {
                data.ResolvedInReleaseAfter = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ResolvedInReleaseAfter = types.StringValue(string(jsonBytes))
            } else {
                data.ResolvedInReleaseAfter = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ResolvedInReleaseAfter = types.StringValue(string(jsonBytes))
        } else {
            data.ResolvedInReleaseAfter = types.StringNull()
        }
    } else if arr, ok := dataMap["resolvedInReleaseAfter"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ResolvedInReleaseAfter = types.StringValue(string(jsonBytes))
        } else {
            data.ResolvedInReleaseAfter = types.StringNull()
        }
    } else if val, ok := dataMap["resolvedInReleaseAfter"].(string); ok {
        data.ResolvedInReleaseAfter = types.StringValue(val)
    } else {
        data.ResolvedInReleaseAfter = types.StringNull()
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
    if obj, ok := dataMap["markedAsResolvedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.MarkedAsResolvedAt = NewRFC3339Value(val)
        } else {
            data.MarkedAsResolvedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["markedAsResolvedAt"].(string); ok && val != "" {
        data.MarkedAsResolvedAt = NewRFC3339Value(val)
    } else {
        data.MarkedAsResolvedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["markedAsArchivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.MarkedAsArchivedAt = NewRFC3339Value(val)
        } else {
            data.MarkedAsArchivedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["markedAsArchivedAt"].(string); ok && val != "" {
        data.MarkedAsArchivedAt = NewRFC3339Value(val)
    } else {
        data.MarkedAsArchivedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["markedAsResolvedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MarkedAsResolvedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MarkedAsResolvedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MarkedAsResolvedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MarkedAsResolvedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.MarkedAsResolvedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MarkedAsResolvedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.MarkedAsResolvedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MarkedAsResolvedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.MarkedAsResolvedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["markedAsResolvedByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.MarkedAsResolvedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.MarkedAsResolvedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["markedAsResolvedByUserId"].(string); ok {
        data.MarkedAsResolvedByUserId = types.StringValue(val)
    } else {
        data.MarkedAsResolvedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["markedAsArchivedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MarkedAsArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MarkedAsArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MarkedAsArchivedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MarkedAsArchivedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.MarkedAsArchivedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MarkedAsArchivedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.MarkedAsArchivedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MarkedAsArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.MarkedAsArchivedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["markedAsArchivedByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.MarkedAsArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.MarkedAsArchivedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["markedAsArchivedByUserId"].(string); ok {
        data.MarkedAsArchivedByUserId = types.StringValue(val)
    } else {
        data.MarkedAsArchivedByUserId = types.StringNull()
    }
    if val, ok := dataMap["unhandled"].(bool); ok {
        data.Unhandled = types.BoolValue(val)
    }
    if obj, ok := dataMap["aiClassification"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiClassification = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AiClassification = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AiClassification = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AiClassification = types.StringValue(string(jsonBytes))
            } else {
                data.AiClassification = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AiClassification = types.StringValue(string(jsonBytes))
            } else {
                data.AiClassification = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AiClassification = types.StringValue(string(jsonBytes))
        } else {
            data.AiClassification = types.StringNull()
        }
    } else if arr, ok := dataMap["aiClassification"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AiClassification = types.StringValue(string(jsonBytes))
        } else {
            data.AiClassification = types.StringNull()
        }
    } else if val, ok := dataMap["aiClassification"].(string); ok {
        data.AiClassification = types.StringValue(val)
    } else {
        data.AiClassification = types.StringNull()
    }
    if obj, ok := dataMap["errorClass"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ErrorClass = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ErrorClass = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ErrorClass = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ErrorClass = types.StringValue(string(jsonBytes))
            } else {
                data.ErrorClass = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ErrorClass = types.StringValue(string(jsonBytes))
            } else {
                data.ErrorClass = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ErrorClass = types.StringValue(string(jsonBytes))
        } else {
            data.ErrorClass = types.StringNull()
        }
    } else if arr, ok := dataMap["errorClass"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ErrorClass = types.StringValue(string(jsonBytes))
        } else {
            data.ErrorClass = types.StringNull()
        }
    } else if val, ok := dataMap["errorClass"].(string); ok {
        data.ErrorClass = types.StringValue(val)
    } else {
        data.ErrorClass = types.StringNull()
    }
    if obj, ok := dataMap["errorClassSource"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ErrorClassSource = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ErrorClassSource = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ErrorClassSource = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ErrorClassSource = types.StringValue(string(jsonBytes))
            } else {
                data.ErrorClassSource = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ErrorClassSource = types.StringValue(string(jsonBytes))
            } else {
                data.ErrorClassSource = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ErrorClassSource = types.StringValue(string(jsonBytes))
        } else {
            data.ErrorClassSource = types.StringNull()
        }
    } else if arr, ok := dataMap["errorClassSource"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.ErrorClassSource = types.StringValue(string(jsonBytes))
        } else {
            data.ErrorClassSource = types.StringNull()
        }
    } else if val, ok := dataMap["errorClassSource"].(string); ok {
        data.ErrorClassSource = types.StringValue(val)
    } else {
        data.ErrorClassSource = types.StringNull()
    }
    if obj, ok := dataMap["aiFixDeclinedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.AiFixDeclinedAt = NewRFC3339Value(val)
        } else {
            data.AiFixDeclinedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["aiFixDeclinedAt"].(string); ok && val != "" {
        data.AiFixDeclinedAt = NewRFC3339Value(val)
    } else {
        data.AiFixDeclinedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["aiExplanation"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiExplanation = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AiExplanation = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AiExplanation = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AiExplanation = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.AiExplanation = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AiExplanation = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.AiExplanation = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AiExplanation = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.AiExplanation = NewJSONSubsetNull()
        }
    } else if arr, ok := dataMap["aiExplanation"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.AiExplanation = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.AiExplanation = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["aiExplanation"].(string); ok {
        data.AiExplanation = NewJSONSubsetValue(val)
    } else {
        data.AiExplanation = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["aiExplanationGeneratedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.AiExplanationGeneratedAt = NewRFC3339Value(val)
        } else {
            data.AiExplanationGeneratedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["aiExplanationGeneratedAt"].(string); ok && val != "" {
        data.AiExplanationGeneratedAt = NewRFC3339Value(val)
    } else {
        data.AiExplanationGeneratedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["mergedIntoTelemetryExceptionId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MergedIntoTelemetryExceptionId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MergedIntoTelemetryExceptionId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MergedIntoTelemetryExceptionId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MergedIntoTelemetryExceptionId = types.StringValue(string(jsonBytes))
            } else {
                data.MergedIntoTelemetryExceptionId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MergedIntoTelemetryExceptionId = types.StringValue(string(jsonBytes))
            } else {
                data.MergedIntoTelemetryExceptionId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MergedIntoTelemetryExceptionId = types.StringValue(string(jsonBytes))
        } else {
            data.MergedIntoTelemetryExceptionId = types.StringNull()
        }
    } else if arr, ok := dataMap["mergedIntoTelemetryExceptionId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.MergedIntoTelemetryExceptionId = types.StringValue(string(jsonBytes))
        } else {
            data.MergedIntoTelemetryExceptionId = types.StringNull()
        }
    } else if val, ok := dataMap["mergedIntoTelemetryExceptionId"].(string); ok {
        data.MergedIntoTelemetryExceptionId = types.StringValue(val)
    } else {
        data.MergedIntoTelemetryExceptionId = types.StringNull()
    }
    if obj, ok := dataMap["mergedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.MergedAt = NewRFC3339Value(val)
        } else {
            data.MergedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["mergedAt"].(string); ok && val != "" {
        data.MergedAt = NewRFC3339Value(val)
    } else {
        data.MergedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["mergedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MergedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MergedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MergedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MergedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.MergedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MergedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.MergedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MergedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.MergedByUserId = types.StringNull()
        }
    } else if arr, ok := dataMap["mergedByUserId"].([]interface{}); ok {
        // A JSON array column: the array, as JSON
        if jsonBytes, err := json.Marshal(r.normalizeURLWrappers(arr)); err == nil {
            data.MergedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.MergedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["mergedByUserId"].(string); ok {
        data.MergedByUserId = types.StringValue(val)
    } else {
        data.MergedByUserId = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    data.Id = state.Id

    // Unconfigured attributes keep their planned value; see keepPlannedValues.
    r.keepPlannedValues(&data, &plan, &config)
    // The address of a person named by ID (see people.go).
    data.AssignToUserEmail = personEmailAfterWrite(ctx, r.client, data.AssignToUserId, data.AssignToUserEmail)

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExceptionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data ExceptionResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    httpResp, err := r.client.Delete(ctx, "/telemetry-exception/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete exception, got error: %s", err))
        return
    }

    // A failed delete must keep the resource in state — silently dropping it
    // orphans real infrastructure. 404 means it is already gone.
    if httpResp.StatusCode >= 400 && httpResp.StatusCode != http.StatusNotFound {
        err = r.client.ParseResponse(httpResp, nil)
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to delete exception: %s", err))
        return
    }
    if httpResp.Body != nil {
        httpResp.Body.Close()
    }
}


func (r *ExceptionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}


// keepPlannedValues puts back, after a create or an update, the planned value
// of each optional attribute the configuration leaves out. The server keeps
// some of these up to date on its own (when it last checked a heartbeat, the
// status a probe last reported...), so the value read back after the write can
// already differ from the plan, and Terraform would fail the apply with
// "Provider produced inconsistent result after apply". The next refresh reads
// the server's value, which is never a diff for an attribute nobody configured.
func (r *ExceptionResource) keepPlannedValues(data *ExceptionResourceModel, plan *ExceptionResourceModel, config *ExceptionResourceModel) {
    if config.PrimaryEntityType.IsNull() && !plan.PrimaryEntityType.IsUnknown() {
        data.PrimaryEntityType = plan.PrimaryEntityType
    }
    if config.Message.IsNull() && !plan.Message.IsUnknown() {
        data.Message = plan.Message
    }
    if config.StackTrace.IsNull() && !plan.StackTrace.IsUnknown() {
        data.StackTrace = plan.StackTrace
    }
    if config.ExceptionType.IsNull() && !plan.ExceptionType.IsUnknown() {
        data.ExceptionType = plan.ExceptionType
    }
    if config.Fingerprint.IsNull() && !plan.Fingerprint.IsUnknown() {
        data.Fingerprint = plan.Fingerprint
    }
    if config.FirstSeenAt.IsNull() && !plan.FirstSeenAt.IsUnknown() {
        data.FirstSeenAt = plan.FirstSeenAt
    }
    if config.LastSeenAt.IsNull() && !plan.LastSeenAt.IsUnknown() {
        data.LastSeenAt = plan.LastSeenAt
    }
    if config.AssignToUserId.IsNull() && !plan.AssignToUserId.IsUnknown() {
        data.AssignToUserId = plan.AssignToUserId
    }
    if config.AssignToTeamId.IsNull() && !plan.AssignToTeamId.IsUnknown() {
        data.AssignToTeamId = plan.AssignToTeamId
    }
    if config.IsResolved.IsNull() && !plan.IsResolved.IsUnknown() {
        data.IsResolved = plan.IsResolved
    }
    if config.IsArchived.IsNull() && !plan.IsArchived.IsUnknown() {
        data.IsArchived = plan.IsArchived
    }
    if config.OccuranceCount.IsNull() && !plan.OccuranceCount.IsUnknown() {
        data.OccuranceCount = plan.OccuranceCount
    }
    if config.FirstSeenInRelease.IsNull() && !plan.FirstSeenInRelease.IsUnknown() {
        data.FirstSeenInRelease = plan.FirstSeenInRelease
    }
    if config.LastSeenInRelease.IsNull() && !plan.LastSeenInRelease.IsUnknown() {
        data.LastSeenInRelease = plan.LastSeenInRelease
    }
    if config.Environment.IsNull() && !plan.Environment.IsUnknown() {
        data.Environment = plan.Environment
    }
    if config.RegressedAt.IsNull() && !plan.RegressedAt.IsUnknown() {
        data.RegressedAt = plan.RegressedAt
    }
    if config.ResolvedInReleaseAfter.IsNull() && !plan.ResolvedInReleaseAfter.IsUnknown() {
        data.ResolvedInReleaseAfter = plan.ResolvedInReleaseAfter
    }
    if config.AssignToUserEmail.IsNull() && !plan.AssignToUserEmail.IsUnknown() {
        data.AssignToUserEmail = plan.AssignToUserEmail
    }
}

// Helper method to convert Terraform map to Go interface{}
func (r *ExceptionResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *ExceptionResource) convertTerraformListToInterface(terraformList types.List) interface{} {
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
func (r *ExceptionResource) convertTerraformSetToInterface(terraformSet types.Set) interface{} {
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
func (r *ExceptionResource) parseJSONField(terraformString basetypes.StringValuable) interface{} {
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
func (r *ExceptionResource) normalizeURLWrappers(value interface{}) interface{} {
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

func (r *ExceptionResource) normalizeURLString(value string) string {
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
func (r *ExceptionResource) bigFloatToFloat64(bf *big.Float) interface{} {
    if bf == nil {
        return nil
    }
    f, _ := bf.Float64()
    return f
}

// Helper method to check if a type string is a valid OneUptime ObjectType.
// The registry itself lives in objecttypes.go, shared across the package.
func (r *ExceptionResource) isValidOneUptimeObjectType(typeStr string) bool {
    return validOneUptimeObjectTypes[typeStr]
}
