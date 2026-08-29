package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "math/big"

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
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    PrimaryEntityId types.String `tfsdk:"primary_entity_id"`
    PrimaryEntityType types.String `tfsdk:"primary_entity_type"`
    Message types.String `tfsdk:"message"`
    StackTrace types.String `tfsdk:"stack_trace"`
    ExceptionType types.String `tfsdk:"exception_type"`
    Fingerprint types.String `tfsdk:"fingerprint"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
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
    AiFixDeclinedAt types.String `tfsdk:"ai_fix_declined_at"`
}

func (d *ExceptionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_exception"
}

func (d *ExceptionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "List of all Telemetry Exceptions created for the telemetry service for this OneUptime project and it's status. Look up an existing exception by `id` or by `name`.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Exactly one of `id` or `name` must be set.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.",
                Optional: true,
                Computed: true,
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "deleted_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "version": schema.NumberAttribute{
                MarkdownDescription: "Object version",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "primary_entity_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "primary_entity_type": schema.StringAttribute{
                MarkdownDescription: "Resource type that produced this exception (e.g. OpenTelemetry service, Host, DockerHost, KubernetesCluster, or Unknown for unattributed telemetry)..",
                Computed: true,
            },
            "message": schema.StringAttribute{
                MarkdownDescription: "Exception message that was thrown by the telemetry service.",
                Computed: true,
            },
            "stack_trace": schema.StringAttribute{
                MarkdownDescription: "Stack trace of the exception that was thrown by the telemetry service.",
                Computed: true,
            },
            "exception_type": schema.StringAttribute{
                MarkdownDescription: "Type of the exception that was thrown by the telemetry service.",
                Computed: true,
            },
            "fingerprint": schema.StringAttribute{
                MarkdownDescription: "Finger print of the exception that was thrown by the telemetry service.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "deleted_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "marked_as_resolved_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "marked_as_archived_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "first_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "assign_to_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "assign_to_team_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "marked_as_resolved_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "marked_as_archived_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_resolved": schema.BoolAttribute{
                MarkdownDescription: "Is this exception resolved?.",
                Computed: true,
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Is this exception archived?.",
                Computed: true,
            },
            "occurance_count": schema.NumberAttribute{
                MarkdownDescription: "Number of times this exception has occurred.",
                Computed: true,
            },
            "first_seen_in_release": schema.StringAttribute{
                MarkdownDescription: "The service version / release in which this exception was first observed.",
                Computed: true,
            },
            "last_seen_in_release": schema.StringAttribute{
                MarkdownDescription: "The most recent service version / release in which this exception was observed.",
                Computed: true,
            },
            "environment": schema.StringAttribute{
                MarkdownDescription: "Deployment environment from deployment.environment resource attribute.",
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
            "ai_fix_declined_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
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

    hasId := !data.Id.IsNull() && data.Id.ValueString() != ""
    hasName := !data.Name.IsNull() && data.Name.ValueString() != ""
    if hasId == hasName {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Exactly one of `id` or `name` must be set to look up a exception.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "name": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "projectId": true,
        "primaryEntityId": true,
        "primaryEntityType": true,
        "message": true,
        "stackTrace": true,
        "exceptionType": true,
        "fingerprint": true,
        "createdByUserId": true,
        "deletedByUserId": true,
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
        "aiFixDeclinedAt": true,
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
    } else {
        listBody := map[string]interface{}{
            "query": map[string]interface{}{
                "name": data.Name.ValueString(),
            },
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No exception found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one exception matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
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
    } else if val, ok := item["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
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
    if obj, ok := item["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
        }
    } else if val, ok := item["deletedAt"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    } else {
        data.DeletedAt = types.StringNull()
    }
    if val, ok := item["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["version"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        data.Version = types.NumberNull()
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
    } else if val, ok := item["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if obj, ok := item["deletedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeletedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedByUserId = types.StringNull()
        }
    } else if val, ok := item["deletedByUserId"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    } else {
        data.DeletedByUserId = types.StringNull()
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
    if val, ok := item["occuranceCount"].(float64); ok {
        data.OccuranceCount = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["occuranceCount"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.OccuranceCount = types.NumberValue(big.NewFloat(val))
        } else {
            data.OccuranceCount = types.NumberNull()
        }
    } else {
        data.OccuranceCount = types.NumberNull()
    }
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
    } else if val, ok := item["aiClassification"].(string); ok {
        data.AiClassification = types.StringValue(val)
    } else {
        data.AiClassification = types.StringNull()
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
    } else if val, ok := item["aiFixDeclinedAt"].(string); ok {
        data.AiFixDeclinedAt = types.StringValue(val)
    } else {
        data.AiFixDeclinedAt = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
