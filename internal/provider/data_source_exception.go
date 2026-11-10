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
var _ datasource.DataSource = &ExceptionDataSource{}

func NewExceptionDataSource() datasource.DataSource {
    return &ExceptionDataSource{}
}

// ExceptionDataSource defines the data source implementation.
type ExceptionDataSource struct {
    client *Client
}

// ExceptionDataSourceModel describes the data source data model.
type ExceptionDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    PrimaryEntityId types.String `tfsdk:"primary_entity_id"`
    PrimaryEntityType types.String `tfsdk:"primary_entity_type"`
    Message types.String `tfsdk:"message"`
    StackTrace types.String `tfsdk:"stack_trace"`
    ExceptionType types.String `tfsdk:"exception_type"`
    Fingerprint types.String `tfsdk:"fingerprint"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    MarkedAsResolvedAt types.String `tfsdk:"marked_as_resolved_at"`
    MarkedAsArchivedAt types.String `tfsdk:"marked_as_archived_at"`
    FirstSeenAt types.String `tfsdk:"first_seen_at"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
    AssignToUserId types.String `tfsdk:"assign_to_user_id"`
    AssignToTeamId types.String `tfsdk:"assign_to_team_id"`
    MarkedAsResolvedByUserId types.String `tfsdk:"marked_as_resolved_by_user_id"`
    MarkedAsArchivedByUserId types.String `tfsdk:"marked_as_archived_by_user_id"`
    IsResolved types.Bool `tfsdk:"is_resolved"`
    IsArchived types.Bool `tfsdk:"is_archived"`
    OccuranceCount types.Number `tfsdk:"occurance_count"`
    FirstSeenInRelease types.String `tfsdk:"first_seen_in_release"`
    LastSeenInRelease types.String `tfsdk:"last_seen_in_release"`
    Environment types.String `tfsdk:"environment"`
    Unhandled types.Bool `tfsdk:"unhandled"`
    AiClassification types.String `tfsdk:"ai_classification"`
    ErrorClass types.String `tfsdk:"error_class"`
    ErrorClassSource types.String `tfsdk:"error_class_source"`
    AiFixDeclinedAt types.String `tfsdk:"ai_fix_declined_at"`
    RegressedAt types.String `tfsdk:"regressed_at"`
    ResolvedInReleaseAfter types.String `tfsdk:"resolved_in_release_after"`
    AiExplanation types.String `tfsdk:"ai_explanation"`
    AiExplanationGeneratedAt types.String `tfsdk:"ai_explanation_generated_at"`
    MergedIntoTelemetryExceptionId types.String `tfsdk:"merged_into_telemetry_exception_id"`
    MergedAt types.String `tfsdk:"merged_at"`
    MergedByUserId types.String `tfsdk:"merged_by_user_id"`
    AssignToUserEmail types.String `tfsdk:"assign_to_user_email"`
}

func (d *ExceptionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_exception"
}

func (d *ExceptionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "List of all Telemetry Exceptions created for the telemetry service for this OneUptime project and it's status. Look up an existing exception by `id`, or by any of its other arguments (`ai_classification`, `assign_to_team_id`, `assign_to_user_email`, ...): each one set must match, and exactly one exception may match them all.",

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
                MarkdownDescription: "ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "primary_entity_id": schema.StringAttribute{
                MarkdownDescription: "ID of the resource this exception belongs to (Service / Host / DockerHost / KubernetesCluster, or the projectId for unattributed telemetry — disambiguated by primaryEntityType).",
                Optional: true,
                Computed: true,
            },
            "primary_entity_type": schema.StringAttribute{
                MarkdownDescription: "Resource type that produced this exception (e.g. OpenTelemetry service, Host, DockerHost, KubernetesCluster, or Unknown for unattributed telemetry).",
                Optional: true,
                Computed: true,
            },
            "message": schema.StringAttribute{
                MarkdownDescription: "Exception message that was thrown by the telemetry service.",
                Optional: true,
                Computed: true,
            },
            "stack_trace": schema.StringAttribute{
                MarkdownDescription: "Stack trace of the exception that was thrown by the telemetry service.",
                Optional: true,
                Computed: true,
            },
            "exception_type": schema.StringAttribute{
                MarkdownDescription: "Type of the exception that was thrown by the telemetry service.",
                Optional: true,
                Computed: true,
            },
            "fingerprint": schema.StringAttribute{
                MarkdownDescription: "Finger print of the exception that was thrown by the telemetry service.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: dataSourceRecordIDValidators("oneuptime_user"),
            },
            "marked_as_resolved_at": schema.StringAttribute{
                MarkdownDescription: "When this exception was last resolved. Empty while it is unresolved.",
                Computed: true,
            },
            "marked_as_archived_at": schema.StringAttribute{
                MarkdownDescription: "When this exception was last archived. Empty while it is not archived.",
                Computed: true,
            },
            "first_seen_at": schema.StringAttribute{
                MarkdownDescription: "When the earliest occurrence of this exception happened.",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "When the latest occurrence of this exception happened.",
                Computed: true,
            },
            "assign_to_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who this exception is assigned to. The ID of a `oneuptime_user` (see the data source). Or look them up by email address, in `assign_to_user_email`.",
                Optional: true,
                Computed: true,
                Validators: personIDValidators("assign_to_user_email", false),
            },
            "assign_to_team_id": schema.StringAttribute{
                MarkdownDescription: "Team ID who this exception is assigned to. The ID of a `oneuptime_team`.",
                Optional: true,
                Computed: true,
                Validators: recordIDValidators("oneuptime_team"),
            },
            "marked_as_resolved_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who marked this exception as resolved. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: dataSourceRecordIDValidators("oneuptime_user"),
            },
            "marked_as_archived_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who marked this exception as archived. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: dataSourceRecordIDValidators("oneuptime_user"),
            },
            "is_resolved": schema.BoolAttribute{
                MarkdownDescription: "Is this exception resolved?",
                Optional: true,
                Computed: true,
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Is this exception archived?",
                Optional: true,
                Computed: true,
            },
            "occurance_count": schema.NumberAttribute{
                MarkdownDescription: "Number of times this exception has occurred.",
                Optional: true,
                Computed: true,
            },
            "first_seen_in_release": schema.StringAttribute{
                MarkdownDescription: "The service version / release in which this exception was first observed.",
                Optional: true,
                Computed: true,
            },
            "last_seen_in_release": schema.StringAttribute{
                MarkdownDescription: "The most recent service version / release in which this exception was observed.",
                Optional: true,
                Computed: true,
            },
            "environment": schema.StringAttribute{
                MarkdownDescription: "Deployment environment from deployment.environment resource attribute.",
                Optional: true,
                Computed: true,
            },
            "unhandled": schema.BoolAttribute{
                MarkdownDescription: "True when at least one occurrence of this exception escaped its span scope (was unhandled, per OTel exception.escaped).",
                Optional: true,
                Computed: true,
            },
            "ai_classification": schema.StringAttribute{
                MarkdownDescription: "AI triage verdict for this exception group (code-fault, user-error, expected-denial, infrastructure).",
                Optional: true,
                Computed: true,
            },
            "error_class": schema.StringAttribute{
                MarkdownDescription: "Fault domain of this exception group (code-fault, user-error, expected-denial, infrastructure, unknown). Non-actionable classes are excluded from the Issues list.",
                Optional: true,
                Computed: true,
            },
            "error_class_source": schema.StringAttribute{
                MarkdownDescription: "Where the error class came from: default (unclassified), declared (by the emitting code), ai (triage verdict) or manual (a human).",
                Optional: true,
                Computed: true,
            },
            "ai_fix_declined_at": schema.StringAttribute{
                MarkdownDescription: "Set when an AI-authored fix pull request for this exception was closed without merging; suppresses further automatic fix attempts.",
                Computed: true,
            },
            "regressed_at": schema.StringAttribute{
                MarkdownDescription: "When this exception last came back after it was resolved. Empty when it has not come back since someone last resolved or reopened it.",
                Computed: true,
            },
            "resolved_in_release_after": schema.StringAttribute{
                MarkdownDescription: "Set when the exception was resolved in the next release: the release it was resolved after. Occurrences from this release do not reopen the exception; an occurrence from any other release does.",
                Optional: true,
                Computed: true,
            },
            "ai_explanation": schema.StringAttribute{
                MarkdownDescription: "OneUptime AI's explanation of this exception: a summary, the likely root cause, a suggested fix, how confident it is, the suspect stack frames and next steps. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "ai_explanation_generated_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime AI last wrote its explanation of this exception.",
                Computed: true,
            },
            "merged_into_telemetry_exception_id": schema.StringAttribute{
                MarkdownDescription: "ID of the exception this one was merged into as the same bug. Empty for an exception that was not merged. The ID of a `oneuptime_exception`.",
                Optional: true,
                Computed: true,
                Validators: recordIDValidators("oneuptime_exception"),
            },
            "merged_at": schema.StringAttribute{
                MarkdownDescription: "When this exception was merged into another.",
                Computed: true,
            },
            "merged_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the user who merged this exception into another. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: dataSourceRecordIDValidators("oneuptime_user"),
            },
            "assign_to_user_email": schema.StringAttribute{
                MarkdownDescription: "Look it up by the email address of the person in `assign_to_user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.",
                Optional: true,
                Computed: true,
                Validators: personEmailValidators(),
            },
        },
    }
}

func (d *ExceptionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ExceptionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ExceptionDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.PrimaryEntityId.IsNull() && !data.PrimaryEntityId.IsUnknown() {
        filters["primaryEntityId"] = data.PrimaryEntityId.ValueString()
        filterNames = append(filterNames, "primary_entity_id = "+fmt.Sprintf("%q", data.PrimaryEntityId.ValueString()))
    }
    if !data.PrimaryEntityType.IsNull() && !data.PrimaryEntityType.IsUnknown() {
        filters["primaryEntityType"] = data.PrimaryEntityType.ValueString()
        filterNames = append(filterNames, "primary_entity_type = "+fmt.Sprintf("%q", data.PrimaryEntityType.ValueString()))
    }
    if !data.Message.IsNull() && !data.Message.IsUnknown() {
        filters["message"] = data.Message.ValueString()
        filterNames = append(filterNames, "message = "+fmt.Sprintf("%q", data.Message.ValueString()))
    }
    if !data.StackTrace.IsNull() && !data.StackTrace.IsUnknown() {
        filters["stackTrace"] = data.StackTrace.ValueString()
        filterNames = append(filterNames, "stack_trace = "+fmt.Sprintf("%q", data.StackTrace.ValueString()))
    }
    if !data.ExceptionType.IsNull() && !data.ExceptionType.IsUnknown() {
        filters["exceptionType"] = data.ExceptionType.ValueString()
        filterNames = append(filterNames, "exception_type = "+fmt.Sprintf("%q", data.ExceptionType.ValueString()))
    }
    if !data.Fingerprint.IsNull() && !data.Fingerprint.IsUnknown() {
        filters["fingerprint"] = data.Fingerprint.ValueString()
        filterNames = append(filterNames, "fingerprint = "+fmt.Sprintf("%q", data.Fingerprint.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.AssignToUserId.IsNull() && !data.AssignToUserId.IsUnknown() {
        filters["assignToUserId"] = data.AssignToUserId.ValueString()
        filterNames = append(filterNames, "assign_to_user_id = "+fmt.Sprintf("%q", data.AssignToUserId.ValueString()))
    }
    if !data.AssignToTeamId.IsNull() && !data.AssignToTeamId.IsUnknown() {
        filters["assignToTeamId"] = data.AssignToTeamId.ValueString()
        filterNames = append(filterNames, "assign_to_team_id = "+fmt.Sprintf("%q", data.AssignToTeamId.ValueString()))
    }
    if !data.MarkedAsResolvedByUserId.IsNull() && !data.MarkedAsResolvedByUserId.IsUnknown() {
        filters["markedAsResolvedByUserId"] = data.MarkedAsResolvedByUserId.ValueString()
        filterNames = append(filterNames, "marked_as_resolved_by_user_id = "+fmt.Sprintf("%q", data.MarkedAsResolvedByUserId.ValueString()))
    }
    if !data.MarkedAsArchivedByUserId.IsNull() && !data.MarkedAsArchivedByUserId.IsUnknown() {
        filters["markedAsArchivedByUserId"] = data.MarkedAsArchivedByUserId.ValueString()
        filterNames = append(filterNames, "marked_as_archived_by_user_id = "+fmt.Sprintf("%q", data.MarkedAsArchivedByUserId.ValueString()))
    }
    if !data.IsResolved.IsNull() && !data.IsResolved.IsUnknown() {
        filters["isResolved"] = data.IsResolved.ValueBool()
        filterNames = append(filterNames, "is_resolved = "+fmt.Sprintf("%t", data.IsResolved.ValueBool()))
    }
    if !data.IsArchived.IsNull() && !data.IsArchived.IsUnknown() {
        filters["isArchived"] = data.IsArchived.ValueBool()
        filterNames = append(filterNames, "is_archived = "+fmt.Sprintf("%t", data.IsArchived.ValueBool()))
    }
    if !data.OccuranceCount.IsNull() && !data.OccuranceCount.IsUnknown() {
        filters["occuranceCount"] = lookupNumber(data.OccuranceCount)
        filterNames = append(filterNames, "occurance_count = "+data.OccuranceCount.ValueBigFloat().String())
    }
    if !data.FirstSeenInRelease.IsNull() && !data.FirstSeenInRelease.IsUnknown() {
        filters["firstSeenInRelease"] = data.FirstSeenInRelease.ValueString()
        filterNames = append(filterNames, "first_seen_in_release = "+fmt.Sprintf("%q", data.FirstSeenInRelease.ValueString()))
    }
    if !data.LastSeenInRelease.IsNull() && !data.LastSeenInRelease.IsUnknown() {
        filters["lastSeenInRelease"] = data.LastSeenInRelease.ValueString()
        filterNames = append(filterNames, "last_seen_in_release = "+fmt.Sprintf("%q", data.LastSeenInRelease.ValueString()))
    }
    if !data.Environment.IsNull() && !data.Environment.IsUnknown() {
        filters["environment"] = data.Environment.ValueString()
        filterNames = append(filterNames, "environment = "+fmt.Sprintf("%q", data.Environment.ValueString()))
    }
    if !data.Unhandled.IsNull() && !data.Unhandled.IsUnknown() {
        filters["unhandled"] = data.Unhandled.ValueBool()
        filterNames = append(filterNames, "unhandled = "+fmt.Sprintf("%t", data.Unhandled.ValueBool()))
    }
    if !data.AiClassification.IsNull() && !data.AiClassification.IsUnknown() {
        filters["aiClassification"] = data.AiClassification.ValueString()
        filterNames = append(filterNames, "ai_classification = "+fmt.Sprintf("%q", data.AiClassification.ValueString()))
    }
    if !data.ErrorClass.IsNull() && !data.ErrorClass.IsUnknown() {
        filters["errorClass"] = data.ErrorClass.ValueString()
        filterNames = append(filterNames, "error_class = "+fmt.Sprintf("%q", data.ErrorClass.ValueString()))
    }
    if !data.ErrorClassSource.IsNull() && !data.ErrorClassSource.IsUnknown() {
        filters["errorClassSource"] = data.ErrorClassSource.ValueString()
        filterNames = append(filterNames, "error_class_source = "+fmt.Sprintf("%q", data.ErrorClassSource.ValueString()))
    }
    if !data.ResolvedInReleaseAfter.IsNull() && !data.ResolvedInReleaseAfter.IsUnknown() {
        filters["resolvedInReleaseAfter"] = data.ResolvedInReleaseAfter.ValueString()
        filterNames = append(filterNames, "resolved_in_release_after = "+fmt.Sprintf("%q", data.ResolvedInReleaseAfter.ValueString()))
    }
    if !data.MergedIntoTelemetryExceptionId.IsNull() && !data.MergedIntoTelemetryExceptionId.IsUnknown() {
        filters["mergedIntoTelemetryExceptionId"] = data.MergedIntoTelemetryExceptionId.ValueString()
        filterNames = append(filterNames, "merged_into_telemetry_exception_id = "+fmt.Sprintf("%q", data.MergedIntoTelemetryExceptionId.ValueString()))
    }
    if !data.MergedByUserId.IsNull() && !data.MergedByUserId.IsUnknown() {
        filters["mergedByUserId"] = data.MergedByUserId.ValueString()
        filterNames = append(filterNames, "merged_by_user_id = "+fmt.Sprintf("%q", data.MergedByUserId.ValueString()))
    }
    if !data.AssignToUserEmail.IsNull() && !data.AssignToUserEmail.IsUnknown() {
        personID, ok := personIDForLookup(ctx, d.client, data.AssignToUserEmail.ValueString(), "assign_to_user_email", &resp.Diagnostics)
        if !ok {
            return
        }
        filters["assignToUserId"] = personID
        filterNames = append(filterNames, "assign_to_user_email = "+fmt.Sprintf("%q", data.AssignToUserEmail.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the exception up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the exception up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "primaryEntityId": true,
        "primaryEntityType": true,
        "message": true,
        "stackTrace": true,
        "exceptionType": true,
        "fingerprint": true,
        "createdByUserId": true,
        "markedAsResolvedAt": true,
        "markedAsArchivedAt": true,
        "firstSeenAt": true,
        "lastSeenAt": true,
        "assignToUserId": true,
        "assignToTeamId": true,
        "markedAsResolvedByUserId": true,
        "markedAsArchivedByUserId": true,
        "isResolved": true,
        "isArchived": true,
        "occuranceCount": true,
        "firstSeenInRelease": true,
        "lastSeenInRelease": true,
        "environment": true,
        "unhandled": true,
        "aiClassification": true,
        "errorClass": true,
        "errorClassSource": true,
        "aiFixDeclinedAt": true,
        "regressedAt": true,
        "resolvedInReleaseAfter": true,
        "aiExplanation": true,
        "aiExplanationGeneratedAt": true,
        "mergedIntoTelemetryExceptionId": true,
        "mergedAt": true,
        "mergedByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/telemetry-exception/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read exception, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No exception found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read exception: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/telemetry-exception/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list exception, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list exception: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No exception matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one exception matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for exception.")
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
    } else if arr, ok := item["_id"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["createdAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["updatedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["projectId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ProjectId = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := item["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := item["primaryEntityId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PrimaryEntityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PrimaryEntityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PrimaryEntityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PrimaryEntityId = types.StringValue(string(jsonBytes))
        } else {
            data.PrimaryEntityId = types.StringNull()
        }
    } else if arr, ok := item["primaryEntityId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.PrimaryEntityId = types.StringValue(string(jsonBytes))
        } else {
            data.PrimaryEntityId = types.StringNull()
        }
    } else if val, ok := item["primaryEntityId"].(string); ok {
        data.PrimaryEntityId = types.StringValue(val)
    } else {
        data.PrimaryEntityId = types.StringNull()
    }
    if obj, ok := item["primaryEntityType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PrimaryEntityType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PrimaryEntityType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PrimaryEntityType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PrimaryEntityType = types.StringValue(string(jsonBytes))
        } else {
            data.PrimaryEntityType = types.StringNull()
        }
    } else if arr, ok := item["primaryEntityType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.PrimaryEntityType = types.StringValue(string(jsonBytes))
        } else {
            data.PrimaryEntityType = types.StringNull()
        }
    } else if val, ok := item["primaryEntityType"].(string); ok {
        data.PrimaryEntityType = types.StringValue(val)
    } else {
        data.PrimaryEntityType = types.StringNull()
    }
    if obj, ok := item["message"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Message = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Message = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Message = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Message = types.StringValue(string(jsonBytes))
        } else {
            data.Message = types.StringNull()
        }
    } else if arr, ok := item["message"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Message = types.StringValue(string(jsonBytes))
        } else {
            data.Message = types.StringNull()
        }
    } else if val, ok := item["message"].(string); ok {
        data.Message = types.StringValue(val)
    } else {
        data.Message = types.StringNull()
    }
    if obj, ok := item["stackTrace"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StackTrace = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StackTrace = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StackTrace = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StackTrace = types.StringValue(string(jsonBytes))
        } else {
            data.StackTrace = types.StringNull()
        }
    } else if arr, ok := item["stackTrace"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.StackTrace = types.StringValue(string(jsonBytes))
        } else {
            data.StackTrace = types.StringNull()
        }
    } else if val, ok := item["stackTrace"].(string); ok {
        data.StackTrace = types.StringValue(val)
    } else {
        data.StackTrace = types.StringNull()
    }
    if obj, ok := item["exceptionType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ExceptionType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ExceptionType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ExceptionType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ExceptionType = types.StringValue(string(jsonBytes))
        } else {
            data.ExceptionType = types.StringNull()
        }
    } else if arr, ok := item["exceptionType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ExceptionType = types.StringValue(string(jsonBytes))
        } else {
            data.ExceptionType = types.StringNull()
        }
    } else if val, ok := item["exceptionType"].(string); ok {
        data.ExceptionType = types.StringValue(val)
    } else {
        data.ExceptionType = types.StringNull()
    }
    if obj, ok := item["fingerprint"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Fingerprint = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Fingerprint = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Fingerprint = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Fingerprint = types.StringValue(string(jsonBytes))
        } else {
            data.Fingerprint = types.StringNull()
        }
    } else if arr, ok := item["fingerprint"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Fingerprint = types.StringValue(string(jsonBytes))
        } else {
            data.Fingerprint = types.StringNull()
        }
    } else if val, ok := item["fingerprint"].(string); ok {
        data.Fingerprint = types.StringValue(val)
    } else {
        data.Fingerprint = types.StringNull()
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
    } else if arr, ok := item["createdByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := item["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if obj, ok := item["markedAsResolvedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MarkedAsResolvedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MarkedAsResolvedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MarkedAsResolvedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MarkedAsResolvedAt = types.StringValue(string(jsonBytes))
        } else {
            data.MarkedAsResolvedAt = types.StringNull()
        }
    } else if arr, ok := item["markedAsResolvedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.MarkedAsResolvedAt = types.StringValue(string(jsonBytes))
        } else {
            data.MarkedAsResolvedAt = types.StringNull()
        }
    } else if val, ok := item["markedAsResolvedAt"].(string); ok {
        data.MarkedAsResolvedAt = types.StringValue(val)
    } else {
        data.MarkedAsResolvedAt = types.StringNull()
    }
    if obj, ok := item["markedAsArchivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MarkedAsArchivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MarkedAsArchivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MarkedAsArchivedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MarkedAsArchivedAt = types.StringValue(string(jsonBytes))
        } else {
            data.MarkedAsArchivedAt = types.StringNull()
        }
    } else if arr, ok := item["markedAsArchivedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.MarkedAsArchivedAt = types.StringValue(string(jsonBytes))
        } else {
            data.MarkedAsArchivedAt = types.StringNull()
        }
    } else if val, ok := item["markedAsArchivedAt"].(string); ok {
        data.MarkedAsArchivedAt = types.StringValue(val)
    } else {
        data.MarkedAsArchivedAt = types.StringNull()
    }
    if obj, ok := item["firstSeenAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FirstSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FirstSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FirstSeenAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FirstSeenAt = types.StringValue(string(jsonBytes))
        } else {
            data.FirstSeenAt = types.StringNull()
        }
    } else if arr, ok := item["firstSeenAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.FirstSeenAt = types.StringValue(string(jsonBytes))
        } else {
            data.FirstSeenAt = types.StringNull()
        }
    } else if val, ok := item["firstSeenAt"].(string); ok {
        data.FirstSeenAt = types.StringValue(val)
    } else {
        data.FirstSeenAt = types.StringNull()
    }
    if obj, ok := item["lastSeenAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastSeenAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastSeenAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastSeenAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastSeenAt = types.StringNull()
        }
    } else if arr, ok := item["lastSeenAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastSeenAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastSeenAt = types.StringNull()
        }
    } else if val, ok := item["lastSeenAt"].(string); ok {
        data.LastSeenAt = types.StringValue(val)
    } else {
        data.LastSeenAt = types.StringNull()
    }
    if obj, ok := item["assignToUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AssignToUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AssignToUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AssignToUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AssignToUserId = types.StringValue(string(jsonBytes))
        } else {
            data.AssignToUserId = types.StringNull()
        }
    } else if arr, ok := item["assignToUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AssignToUserId = types.StringValue(string(jsonBytes))
        } else {
            data.AssignToUserId = types.StringNull()
        }
    } else if val, ok := item["assignToUserId"].(string); ok {
        data.AssignToUserId = types.StringValue(val)
    } else {
        data.AssignToUserId = types.StringNull()
    }
    if obj, ok := item["assignToTeamId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AssignToTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AssignToTeamId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AssignToTeamId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AssignToTeamId = types.StringValue(string(jsonBytes))
        } else {
            data.AssignToTeamId = types.StringNull()
        }
    } else if arr, ok := item["assignToTeamId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AssignToTeamId = types.StringValue(string(jsonBytes))
        } else {
            data.AssignToTeamId = types.StringNull()
        }
    } else if val, ok := item["assignToTeamId"].(string); ok {
        data.AssignToTeamId = types.StringValue(val)
    } else {
        data.AssignToTeamId = types.StringNull()
    }
    if obj, ok := item["markedAsResolvedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MarkedAsResolvedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MarkedAsResolvedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MarkedAsResolvedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MarkedAsResolvedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.MarkedAsResolvedByUserId = types.StringNull()
        }
    } else if arr, ok := item["markedAsResolvedByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.MarkedAsResolvedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.MarkedAsResolvedByUserId = types.StringNull()
        }
    } else if val, ok := item["markedAsResolvedByUserId"].(string); ok {
        data.MarkedAsResolvedByUserId = types.StringValue(val)
    } else {
        data.MarkedAsResolvedByUserId = types.StringNull()
    }
    if obj, ok := item["markedAsArchivedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MarkedAsArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MarkedAsArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MarkedAsArchivedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MarkedAsArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.MarkedAsArchivedByUserId = types.StringNull()
        }
    } else if arr, ok := item["markedAsArchivedByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.MarkedAsArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.MarkedAsArchivedByUserId = types.StringNull()
        }
    } else if val, ok := item["markedAsArchivedByUserId"].(string); ok {
        data.MarkedAsArchivedByUserId = types.StringValue(val)
    } else {
        data.MarkedAsArchivedByUserId = types.StringNull()
    }
    if val, ok := item["isResolved"].(bool); ok {
        data.IsResolved = types.BoolValue(val)
    } else {
        data.IsResolved = types.BoolNull()
    }
    if val, ok := item["isArchived"].(bool); ok {
        data.IsArchived = types.BoolValue(val)
    } else {
        data.IsArchived = types.BoolNull()
    }
    data.OccuranceCount = numberFromAPI(data.OccuranceCount, item["occuranceCount"])
    if obj, ok := item["firstSeenInRelease"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FirstSeenInRelease = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FirstSeenInRelease = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FirstSeenInRelease = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FirstSeenInRelease = types.StringValue(string(jsonBytes))
        } else {
            data.FirstSeenInRelease = types.StringNull()
        }
    } else if arr, ok := item["firstSeenInRelease"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.FirstSeenInRelease = types.StringValue(string(jsonBytes))
        } else {
            data.FirstSeenInRelease = types.StringNull()
        }
    } else if val, ok := item["firstSeenInRelease"].(string); ok {
        data.FirstSeenInRelease = types.StringValue(val)
    } else {
        data.FirstSeenInRelease = types.StringNull()
    }
    if obj, ok := item["lastSeenInRelease"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastSeenInRelease = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastSeenInRelease = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastSeenInRelease = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastSeenInRelease = types.StringValue(string(jsonBytes))
        } else {
            data.LastSeenInRelease = types.StringNull()
        }
    } else if arr, ok := item["lastSeenInRelease"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastSeenInRelease = types.StringValue(string(jsonBytes))
        } else {
            data.LastSeenInRelease = types.StringNull()
        }
    } else if val, ok := item["lastSeenInRelease"].(string); ok {
        data.LastSeenInRelease = types.StringValue(val)
    } else {
        data.LastSeenInRelease = types.StringNull()
    }
    if obj, ok := item["environment"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Environment = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Environment = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Environment = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Environment = types.StringValue(string(jsonBytes))
        } else {
            data.Environment = types.StringNull()
        }
    } else if arr, ok := item["environment"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Environment = types.StringValue(string(jsonBytes))
        } else {
            data.Environment = types.StringNull()
        }
    } else if val, ok := item["environment"].(string); ok {
        data.Environment = types.StringValue(val)
    } else {
        data.Environment = types.StringNull()
    }
    if val, ok := item["unhandled"].(bool); ok {
        data.Unhandled = types.BoolValue(val)
    } else {
        data.Unhandled = types.BoolNull()
    }
    if obj, ok := item["aiClassification"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiClassification = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AiClassification = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AiClassification = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AiClassification = types.StringValue(string(jsonBytes))
        } else {
            data.AiClassification = types.StringNull()
        }
    } else if arr, ok := item["aiClassification"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AiClassification = types.StringValue(string(jsonBytes))
        } else {
            data.AiClassification = types.StringNull()
        }
    } else if val, ok := item["aiClassification"].(string); ok {
        data.AiClassification = types.StringValue(val)
    } else {
        data.AiClassification = types.StringNull()
    }
    if obj, ok := item["errorClass"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ErrorClass = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ErrorClass = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ErrorClass = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ErrorClass = types.StringValue(string(jsonBytes))
        } else {
            data.ErrorClass = types.StringNull()
        }
    } else if arr, ok := item["errorClass"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ErrorClass = types.StringValue(string(jsonBytes))
        } else {
            data.ErrorClass = types.StringNull()
        }
    } else if val, ok := item["errorClass"].(string); ok {
        data.ErrorClass = types.StringValue(val)
    } else {
        data.ErrorClass = types.StringNull()
    }
    if obj, ok := item["errorClassSource"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ErrorClassSource = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ErrorClassSource = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ErrorClassSource = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ErrorClassSource = types.StringValue(string(jsonBytes))
        } else {
            data.ErrorClassSource = types.StringNull()
        }
    } else if arr, ok := item["errorClassSource"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ErrorClassSource = types.StringValue(string(jsonBytes))
        } else {
            data.ErrorClassSource = types.StringNull()
        }
    } else if val, ok := item["errorClassSource"].(string); ok {
        data.ErrorClassSource = types.StringValue(val)
    } else {
        data.ErrorClassSource = types.StringNull()
    }
    if obj, ok := item["aiFixDeclinedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiFixDeclinedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AiFixDeclinedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AiFixDeclinedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AiFixDeclinedAt = types.StringValue(string(jsonBytes))
        } else {
            data.AiFixDeclinedAt = types.StringNull()
        }
    } else if arr, ok := item["aiFixDeclinedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AiFixDeclinedAt = types.StringValue(string(jsonBytes))
        } else {
            data.AiFixDeclinedAt = types.StringNull()
        }
    } else if val, ok := item["aiFixDeclinedAt"].(string); ok {
        data.AiFixDeclinedAt = types.StringValue(val)
    } else {
        data.AiFixDeclinedAt = types.StringNull()
    }
    if obj, ok := item["regressedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RegressedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RegressedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RegressedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RegressedAt = types.StringValue(string(jsonBytes))
        } else {
            data.RegressedAt = types.StringNull()
        }
    } else if arr, ok := item["regressedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.RegressedAt = types.StringValue(string(jsonBytes))
        } else {
            data.RegressedAt = types.StringNull()
        }
    } else if val, ok := item["regressedAt"].(string); ok {
        data.RegressedAt = types.StringValue(val)
    } else {
        data.RegressedAt = types.StringNull()
    }
    if obj, ok := item["resolvedInReleaseAfter"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResolvedInReleaseAfter = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ResolvedInReleaseAfter = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ResolvedInReleaseAfter = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ResolvedInReleaseAfter = types.StringValue(string(jsonBytes))
        } else {
            data.ResolvedInReleaseAfter = types.StringNull()
        }
    } else if arr, ok := item["resolvedInReleaseAfter"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ResolvedInReleaseAfter = types.StringValue(string(jsonBytes))
        } else {
            data.ResolvedInReleaseAfter = types.StringNull()
        }
    } else if val, ok := item["resolvedInReleaseAfter"].(string); ok {
        data.ResolvedInReleaseAfter = types.StringValue(val)
    } else {
        data.ResolvedInReleaseAfter = types.StringNull()
    }
    if obj, ok := item["aiExplanation"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiExplanation = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AiExplanation = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AiExplanation = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AiExplanation = types.StringValue(string(jsonBytes))
        } else {
            data.AiExplanation = types.StringNull()
        }
    } else if arr, ok := item["aiExplanation"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AiExplanation = types.StringValue(string(jsonBytes))
        } else {
            data.AiExplanation = types.StringNull()
        }
    } else if val, ok := item["aiExplanation"].(string); ok {
        data.AiExplanation = types.StringValue(val)
    } else {
        data.AiExplanation = types.StringNull()
    }
    if obj, ok := item["aiExplanationGeneratedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AiExplanationGeneratedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AiExplanationGeneratedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AiExplanationGeneratedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AiExplanationGeneratedAt = types.StringValue(string(jsonBytes))
        } else {
            data.AiExplanationGeneratedAt = types.StringNull()
        }
    } else if arr, ok := item["aiExplanationGeneratedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AiExplanationGeneratedAt = types.StringValue(string(jsonBytes))
        } else {
            data.AiExplanationGeneratedAt = types.StringNull()
        }
    } else if val, ok := item["aiExplanationGeneratedAt"].(string); ok {
        data.AiExplanationGeneratedAt = types.StringValue(val)
    } else {
        data.AiExplanationGeneratedAt = types.StringNull()
    }
    if obj, ok := item["mergedIntoTelemetryExceptionId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MergedIntoTelemetryExceptionId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MergedIntoTelemetryExceptionId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MergedIntoTelemetryExceptionId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MergedIntoTelemetryExceptionId = types.StringValue(string(jsonBytes))
        } else {
            data.MergedIntoTelemetryExceptionId = types.StringNull()
        }
    } else if arr, ok := item["mergedIntoTelemetryExceptionId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.MergedIntoTelemetryExceptionId = types.StringValue(string(jsonBytes))
        } else {
            data.MergedIntoTelemetryExceptionId = types.StringNull()
        }
    } else if val, ok := item["mergedIntoTelemetryExceptionId"].(string); ok {
        data.MergedIntoTelemetryExceptionId = types.StringValue(val)
    } else {
        data.MergedIntoTelemetryExceptionId = types.StringNull()
    }
    if obj, ok := item["mergedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MergedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MergedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MergedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MergedAt = types.StringValue(string(jsonBytes))
        } else {
            data.MergedAt = types.StringNull()
        }
    } else if arr, ok := item["mergedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.MergedAt = types.StringValue(string(jsonBytes))
        } else {
            data.MergedAt = types.StringNull()
        }
    } else if val, ok := item["mergedAt"].(string); ok {
        data.MergedAt = types.StringValue(val)
    } else {
        data.MergedAt = types.StringNull()
    }
    if obj, ok := item["mergedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MergedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MergedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MergedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MergedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.MergedByUserId = types.StringNull()
        }
    } else if arr, ok := item["mergedByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.MergedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.MergedByUserId = types.StringNull()
        }
    } else if val, ok := item["mergedByUserId"].(string); ok {
        data.MergedByUserId = types.StringValue(val)
    } else {
        data.MergedByUserId = types.StringNull()
    }
    data.AssignToUserEmail = personEmailAfterLookup(ctx, d.client, data.AssignToUserId, data.AssignToUserEmail)

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
