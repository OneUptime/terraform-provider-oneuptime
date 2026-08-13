package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "math/big"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ServiceDataSource{}

func NewServiceDataSource() datasource.DataSource {
    return &ServiceDataSource{}
}

// ServiceDataSource defines the data source implementation.
type ServiceDataSource struct {
    client *Client
}

// ServiceDataSourceModel describes the data source data model.
type ServiceDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Slug types.String `tfsdk:"slug"`
    Description types.String `tfsdk:"description"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsArchived types.Bool `tfsdk:"is_archived"`
    ArchivedAt types.String `tfsdk:"archived_at"`
    ArchivedByUserId types.String `tfsdk:"archived_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
    Labels types.Set `tfsdk:"labels"`
    ServiceColor types.String `tfsdk:"service_color"`
    ServiceLanguage types.String `tfsdk:"service_language"`
    TechStack types.String `tfsdk:"tech_stack"`
    RetainTelemetryDataForDays types.Number `tfsdk:"retain_telemetry_data_for_days"`
    MetricCardinalityBudget types.Number `tfsdk:"metric_cardinality_budget"`
    MetricDownsamplingRetentionDays types.String `tfsdk:"metric_downsampling_retention_days"`
    TelemetryRetentionConfig types.String `tfsdk:"telemetry_retention_config"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
    ServiceVersion types.String `tfsdk:"service_version"`
    DeploymentEnvironment types.String `tfsdk:"deployment_environment"`
    ServiceNamespace types.String `tfsdk:"service_namespace"`
    RuntimeName types.String `tfsdk:"runtime_name"`
    RuntimeVersion types.String `tfsdk:"runtime_version"`
    TelemetrySdkLanguage types.String `tfsdk:"telemetry_sdk_language"`
    CloudProvider types.String `tfsdk:"cloud_provider"`
    CloudPlatform types.String `tfsdk:"cloud_platform"`
    CloudRegion types.String `tfsdk:"cloud_region"`
    CloudAccountId types.String `tfsdk:"cloud_account_id"`
}

func (d *ServiceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_service"
}

func (d *ServiceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Services is a collection of services that you have in your organization. It can be a collection of services that you are monitoring or services that you are providing to your customers. It can be anything that you want to keep track of. Look up an existing service by `id` or by `name`.",

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
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description that will help you remember.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Is this service archived? Archived services are hidden from lists but keep collecting telemetry..",
                Computed: true,
            },
            "archived_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "archived_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "deleted_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in..",
                Computed: true,
                ElementType: types.StringType,
            },
            "service_color": schema.StringAttribute{
                MarkdownDescription: "Color object",
                Computed: true,
            },
            "service_language": schema.StringAttribute{
                MarkdownDescription: "Language in which this service is written",
                Computed: true,
            },
            "tech_stack": schema.StringAttribute{
                MarkdownDescription: "Tech stack used in the service. This will help other developers understand the service better..",
                Computed: true,
            },
            "retain_telemetry_data_for_days": schema.NumberAttribute{
                MarkdownDescription: "Number of days to retain telemetry data for this service. Leave blank to use the project-wide default..",
                Computed: true,
            },
            "metric_cardinality_budget": schema.NumberAttribute{
                MarkdownDescription: "Max number of distinct metric series this service may emit per metric. When exceeded, the highest-cardinality attribute is auto-bucketed. Null inherits the project default..",
                Computed: true,
            },
            "metric_downsampling_retention_days": schema.StringAttribute{
                MarkdownDescription: "Per-tier retention override (raw, 1m, 5m, 1h, 1d) in days. Null fields inherit the project default..",
                Computed: true,
            },
            "telemetry_retention_config": schema.StringAttribute{
                MarkdownDescription: "Per-pillar retention overrides for this service (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the service default, then the project's retention settings..",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "service_version": schema.StringAttribute{
                MarkdownDescription: "Last-seen value of the service.version OpenTelemetry resource attribute..",
                Computed: true,
            },
            "deployment_environment": schema.StringAttribute{
                MarkdownDescription: "Last-seen value of the deployment.environment.name (or deployment.environment) OpenTelemetry resource attribute, e.g. production, staging..",
                Computed: true,
            },
            "service_namespace": schema.StringAttribute{
                MarkdownDescription: "Last-seen value of the service.namespace OpenTelemetry resource attribute..",
                Computed: true,
            },
            "runtime_name": schema.StringAttribute{
                MarkdownDescription: "Last-seen value of the process.runtime.name OpenTelemetry resource attribute, e.g. nodejs, go, OpenJDK Runtime Environment..",
                Computed: true,
            },
            "runtime_version": schema.StringAttribute{
                MarkdownDescription: "Last-seen value of the process.runtime.version OpenTelemetry resource attribute..",
                Computed: true,
            },
            "telemetry_sdk_language": schema.StringAttribute{
                MarkdownDescription: "Last-seen value of the telemetry.sdk.language OpenTelemetry resource attribute, e.g. java, dotnet, nodejs, python, go. Drives technology-specific golden metrics on the service overview..",
                Computed: true,
            },
            "cloud_provider": schema.StringAttribute{
                MarkdownDescription: "Last-seen value of the cloud.provider OpenTelemetry resource attribute, e.g. aws, gcp, azure..",
                Computed: true,
            },
            "cloud_platform": schema.StringAttribute{
                MarkdownDescription: "Last-seen value of the cloud.platform OpenTelemetry resource attribute, e.g. aws_ecs, gcp_cloud_run, aws_lambda..",
                Computed: true,
            },
            "cloud_region": schema.StringAttribute{
                MarkdownDescription: "Last-seen value of the cloud.region OpenTelemetry resource attribute, e.g. us-east-1..",
                Computed: true,
            },
            "cloud_account_id": schema.StringAttribute{
                MarkdownDescription: "Last-seen value of the cloud.account.id OpenTelemetry resource attribute..",
                Computed: true,
            },
        },
    }
}

func (d *ServiceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ServiceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ServiceDataSourceModel

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
            "Exactly one of `id` or `name` must be set to look up a service.",
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
        "slug": true,
        "description": true,
        "createdByUserId": true,
        "isArchived": true,
        "archivedAt": true,
        "archivedByUserId": true,
        "deletedByUserId": true,
        "labels": true,
        "serviceColor": true,
        "serviceLanguage": true,
        "techStack": true,
        "retainTelemetryDataForDays": true,
        "metricCardinalityBudget": true,
        "metricDownsamplingRetentionDays": true,
        "telemetryRetentionConfig": true,
        "lastSeenAt": true,
        "serviceVersion": true,
        "deploymentEnvironment": true,
        "serviceNamespace": true,
        "runtimeName": true,
        "runtimeVersion": true,
        "telemetrySdkLanguage": true,
        "cloudProvider": true,
        "cloudPlatform": true,
        "cloudRegion": true,
        "cloudAccountId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/service/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read service, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No service found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read service: %s", err))
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
        httpResp, err := d.client.Post(ctx, "/service/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list service, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list service: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No service found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one service matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for service.")
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
    if obj, ok := item["slug"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := item["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if obj, ok := item["description"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := item["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
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
    if val, ok := item["isArchived"].(bool); ok {
        data.IsArchived = types.BoolValue(val)
    } else {
        data.IsArchived = types.BoolNull()
    }
    if obj, ok := item["archivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ArchivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ArchivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ArchivedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ArchivedAt = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedAt = types.StringNull()
        }
    } else if val, ok := item["archivedAt"].(string); ok {
        data.ArchivedAt = types.StringValue(val)
    } else {
        data.ArchivedAt = types.StringNull()
    }
    if obj, ok := item["archivedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ArchivedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ArchivedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ArchivedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.ArchivedByUserId = types.StringNull()
        }
    } else if val, ok := item["archivedByUserId"].(string); ok {
        data.ArchivedByUserId = types.StringValue(val)
    } else {
        data.ArchivedByUserId = types.StringNull()
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
    if val, ok := item["labels"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.Labels = types.SetValueMust(types.StringType, setItems)
    } else {
        data.Labels = types.SetNull(types.StringType)
    }
    if obj, ok := item["serviceColor"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServiceColor = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ServiceColor = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ServiceColor = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ServiceColor = types.StringValue(string(jsonBytes))
        } else {
            data.ServiceColor = types.StringNull()
        }
    } else if val, ok := item["serviceColor"].(string); ok {
        data.ServiceColor = types.StringValue(val)
    } else {
        data.ServiceColor = types.StringNull()
    }
    if obj, ok := item["serviceLanguage"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServiceLanguage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ServiceLanguage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ServiceLanguage = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ServiceLanguage = types.StringValue(string(jsonBytes))
        } else {
            data.ServiceLanguage = types.StringNull()
        }
    } else if val, ok := item["serviceLanguage"].(string); ok {
        data.ServiceLanguage = types.StringValue(val)
    } else {
        data.ServiceLanguage = types.StringNull()
    }
    if obj, ok := item["techStack"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TechStack = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TechStack = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TechStack = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TechStack = types.StringValue(string(jsonBytes))
        } else {
            data.TechStack = types.StringNull()
        }
    } else if val, ok := item["techStack"].(string); ok {
        data.TechStack = types.StringValue(val)
    } else {
        data.TechStack = types.StringNull()
    }
    if val, ok := item["retainTelemetryDataForDays"].(float64); ok {
        data.RetainTelemetryDataForDays = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["retainTelemetryDataForDays"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.RetainTelemetryDataForDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.RetainTelemetryDataForDays = types.NumberNull()
        }
    } else {
        data.RetainTelemetryDataForDays = types.NumberNull()
    }
    if val, ok := item["metricCardinalityBudget"].(float64); ok {
        data.MetricCardinalityBudget = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["metricCardinalityBudget"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.MetricCardinalityBudget = types.NumberValue(big.NewFloat(val))
        } else {
            data.MetricCardinalityBudget = types.NumberNull()
        }
    } else {
        data.MetricCardinalityBudget = types.NumberNull()
    }
    if obj, ok := item["metricDownsamplingRetentionDays"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MetricDownsamplingRetentionDays = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MetricDownsamplingRetentionDays = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MetricDownsamplingRetentionDays = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MetricDownsamplingRetentionDays = types.StringValue(string(jsonBytes))
        } else {
            data.MetricDownsamplingRetentionDays = types.StringNull()
        }
    } else if val, ok := item["metricDownsamplingRetentionDays"].(string); ok {
        data.MetricDownsamplingRetentionDays = types.StringValue(val)
    } else {
        data.MetricDownsamplingRetentionDays = types.StringNull()
    }
    if obj, ok := item["telemetryRetentionConfig"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TelemetryRetentionConfig = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TelemetryRetentionConfig = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TelemetryRetentionConfig = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TelemetryRetentionConfig = types.StringValue(string(jsonBytes))
        } else {
            data.TelemetryRetentionConfig = types.StringNull()
        }
    } else if val, ok := item["telemetryRetentionConfig"].(string); ok {
        data.TelemetryRetentionConfig = types.StringValue(val)
    } else {
        data.TelemetryRetentionConfig = types.StringNull()
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
    if obj, ok := item["serviceVersion"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServiceVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ServiceVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ServiceVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ServiceVersion = types.StringValue(string(jsonBytes))
        } else {
            data.ServiceVersion = types.StringNull()
        }
    } else if val, ok := item["serviceVersion"].(string); ok {
        data.ServiceVersion = types.StringValue(val)
    } else {
        data.ServiceVersion = types.StringNull()
    }
    if obj, ok := item["deploymentEnvironment"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeploymentEnvironment = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeploymentEnvironment = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeploymentEnvironment = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeploymentEnvironment = types.StringValue(string(jsonBytes))
        } else {
            data.DeploymentEnvironment = types.StringNull()
        }
    } else if val, ok := item["deploymentEnvironment"].(string); ok {
        data.DeploymentEnvironment = types.StringValue(val)
    } else {
        data.DeploymentEnvironment = types.StringNull()
    }
    if obj, ok := item["serviceNamespace"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServiceNamespace = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ServiceNamespace = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ServiceNamespace = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ServiceNamespace = types.StringValue(string(jsonBytes))
        } else {
            data.ServiceNamespace = types.StringNull()
        }
    } else if val, ok := item["serviceNamespace"].(string); ok {
        data.ServiceNamespace = types.StringValue(val)
    } else {
        data.ServiceNamespace = types.StringNull()
    }
    if obj, ok := item["runtimeName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RuntimeName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RuntimeName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RuntimeName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RuntimeName = types.StringValue(string(jsonBytes))
        } else {
            data.RuntimeName = types.StringNull()
        }
    } else if val, ok := item["runtimeName"].(string); ok {
        data.RuntimeName = types.StringValue(val)
    } else {
        data.RuntimeName = types.StringNull()
    }
    if obj, ok := item["runtimeVersion"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RuntimeVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RuntimeVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RuntimeVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RuntimeVersion = types.StringValue(string(jsonBytes))
        } else {
            data.RuntimeVersion = types.StringNull()
        }
    } else if val, ok := item["runtimeVersion"].(string); ok {
        data.RuntimeVersion = types.StringValue(val)
    } else {
        data.RuntimeVersion = types.StringNull()
    }
    if obj, ok := item["telemetrySdkLanguage"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TelemetrySdkLanguage = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TelemetrySdkLanguage = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TelemetrySdkLanguage = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TelemetrySdkLanguage = types.StringValue(string(jsonBytes))
        } else {
            data.TelemetrySdkLanguage = types.StringNull()
        }
    } else if val, ok := item["telemetrySdkLanguage"].(string); ok {
        data.TelemetrySdkLanguage = types.StringValue(val)
    } else {
        data.TelemetrySdkLanguage = types.StringNull()
    }
    if obj, ok := item["cloudProvider"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudProvider = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CloudProvider = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CloudProvider = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CloudProvider = types.StringValue(string(jsonBytes))
        } else {
            data.CloudProvider = types.StringNull()
        }
    } else if val, ok := item["cloudProvider"].(string); ok {
        data.CloudProvider = types.StringValue(val)
    } else {
        data.CloudProvider = types.StringNull()
    }
    if obj, ok := item["cloudPlatform"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudPlatform = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CloudPlatform = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CloudPlatform = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CloudPlatform = types.StringValue(string(jsonBytes))
        } else {
            data.CloudPlatform = types.StringNull()
        }
    } else if val, ok := item["cloudPlatform"].(string); ok {
        data.CloudPlatform = types.StringValue(val)
    } else {
        data.CloudPlatform = types.StringNull()
    }
    if obj, ok := item["cloudRegion"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudRegion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CloudRegion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CloudRegion = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CloudRegion = types.StringValue(string(jsonBytes))
        } else {
            data.CloudRegion = types.StringNull()
        }
    } else if val, ok := item["cloudRegion"].(string); ok {
        data.CloudRegion = types.StringValue(val)
    } else {
        data.CloudRegion = types.StringNull()
    }
    if obj, ok := item["cloudAccountId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudAccountId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CloudAccountId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CloudAccountId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CloudAccountId = types.StringValue(string(jsonBytes))
        } else {
            data.CloudAccountId = types.StringNull()
        }
    } else if val, ok := item["cloudAccountId"].(string); ok {
        data.CloudAccountId = types.StringValue(val)
    } else {
        data.CloudAccountId = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
