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
var _ datasource.DataSource = &CloudResourceDataSource{}

func NewCloudResourceDataSource() datasource.DataSource {
    return &CloudResourceDataSource{}
}

// CloudResourceDataSource defines the data source implementation.
type CloudResourceDataSource struct {
    client *Client
}

// CloudResourceDataSourceModel describes the data source data model.
type CloudResourceDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Slug types.String `tfsdk:"slug"`
    Description types.String `tfsdk:"description"`
    ResourceIdentifier types.String `tfsdk:"resource_identifier"`
    CloudPlatform types.String `tfsdk:"cloud_platform"`
    CloudProvider types.String `tfsdk:"cloud_provider"`
    CloudRegion types.String `tfsdk:"cloud_region"`
    CloudAccountId types.String `tfsdk:"cloud_account_id"`
    CloudResourceKind types.String `tfsdk:"cloud_resource_kind"`
    CloudResourceType types.String `tfsdk:"cloud_resource_type"`
    ProviderResourceId types.String `tfsdk:"provider_resource_id"`
    CloudResourceGroup types.String `tfsdk:"cloud_resource_group"`
    TelemetryAttributes types.String `tfsdk:"telemetry_attributes"`
    RuntimeName types.String `tfsdk:"runtime_name"`
    RuntimeVersion types.String `tfsdk:"runtime_version"`
    OtelCollectorStatus types.String `tfsdk:"otel_collector_status"`
    AgentVersion types.String `tfsdk:"agent_version"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
    Labels types.Set `tfsdk:"labels"`
    RetainTelemetryDataForDays types.Number `tfsdk:"retain_telemetry_data_for_days"`
    TelemetryRetentionConfig types.String `tfsdk:"telemetry_retention_config"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsArchived types.Bool `tfsdk:"is_archived"`
    ArchivedAt types.String `tfsdk:"archived_at"`
    ArchivedByUserId types.String `tfsdk:"archived_by_user_id"`
    AutoArchivedAt types.String `tfsdk:"auto_archived_at"`
}

func (d *CloudResourceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_cloud_resource"
}

func (d *CloudResourceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Cloud environments - managed compute auto-discovered from OpenTelemetry cloud.platform (e.g. AWS ECS/Fargate, GCP Cloud Run, Azure Container Apps, Elastic Beanstalk, App Runner) - and cloud resources: the IaaS and PaaS resources (virtual machines, load balancers, buckets, managed databases, queues, ...) discovered from the metrics Azure Monitor, Amazon CloudWatch and Google Cloud Monitoring publish about them. Look up an existing cloud resource by `id`, or by any of its other arguments (`name`, `agent_version`, `archived_by_user_id`, ...): each one set must match, and exactly one cloud resource may match them all.",

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
            "name": schema.StringAttribute{
                MarkdownDescription: "Friendly name for this cloud environment. Ingest names a discovered environment after its platform, region and account.",
                Optional: true,
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description that will help you remember.",
                Optional: true,
                Computed: true,
            },
            "resource_identifier": schema.StringAttribute{
                MarkdownDescription: "Environment key: the cloud.platform, cloud.account.id and cloud.region OpenTelemetry resource attributes joined with '|' (e.g. aws_ecs|123456789012|us-east-1; missing parts stay as empty segments). Built by buildCloudEnvironmentKey in Common/Types/Cloud/CloudPlatform. An environment created by hand must carry the same key for ingest to find it instead of creating a duplicate.",
                Optional: true,
                Computed: true,
            },
            "cloud_platform": schema.StringAttribute{
                MarkdownDescription: "Last-seen cloud.platform OpenTelemetry resource attribute, e.g. aws_ecs, gcp_cloud_run, azure_container_apps.",
                Optional: true,
                Computed: true,
            },
            "cloud_provider": schema.StringAttribute{
                MarkdownDescription: "Last-seen cloud.provider OpenTelemetry resource attribute, e.g. aws, gcp, azure.",
                Optional: true,
                Computed: true,
            },
            "cloud_region": schema.StringAttribute{
                MarkdownDescription: "Last-seen cloud.region OpenTelemetry resource attribute, e.g. us-east-1.",
                Optional: true,
                Computed: true,
            },
            "cloud_account_id": schema.StringAttribute{
                MarkdownDescription: "Last-seen cloud.account.id OpenTelemetry resource attribute.",
                Optional: true,
                Computed: true,
            },
            "cloud_resource_kind": schema.StringAttribute{
                MarkdownDescription: "environment: a managed compute environment discovered from the cloud.platform, cloud.account.id and cloud.region resource attributes of a workload's own telemetry. resource: one IaaS or PaaS resource (a virtual machine, a load balancer, a bucket, a managed database, ...) discovered from the metrics Azure Monitor, Amazon CloudWatch or Google Cloud Monitoring publish about it.",
                Optional: true,
                Computed: true,
            },
            "cloud_resource_type": schema.StringAttribute{
                MarkdownDescription: "For a resource: the provider's type for it - the Azure Resource Manager type (Microsoft.Compute/virtualMachines), the AWS CloudFormation type (AWS::EC2::Instance) or the Google Cloud Monitoring resource type (gce_instance).",
                Optional: true,
                Computed: true,
            },
            "provider_resource_id": schema.StringAttribute{
                MarkdownDescription: "For a resource: the provider's id for it - its Azure resource id, its AWS ARN or its Google Cloud full resource name. Where the metrics do not name the resource completely (an AWS resource whose ARN needs an id no metric reports), a readable composite of what they do name.",
                Optional: true,
                Computed: true,
            },
            "cloud_resource_group": schema.StringAttribute{
                MarkdownDescription: "For an Azure resource: the resource group it belongs to.",
                Optional: true,
                Computed: true,
            },
            "telemetry_attributes": schema.StringAttribute{
                MarkdownDescription: "For a resource: the metric attributes, exactly as stored, that select its metrics - for example azuremonitor.resource_id, or the CloudWatch Namespace and identifying Dimensions with the account and region. The resource's pages and the monitors created from them filter on these. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "runtime_name": schema.StringAttribute{
                MarkdownDescription: "Last-seen process.runtime.name OpenTelemetry resource attribute.",
                Optional: true,
                Computed: true,
            },
            "runtime_version": schema.StringAttribute{
                MarkdownDescription: "Last-seen process.runtime.version OpenTelemetry resource attribute.",
                Optional: true,
                Computed: true,
            },
            "otel_collector_status": schema.StringAttribute{
                MarkdownDescription: "Whether telemetry is currently being received (connected) or has gone stale (disconnected).",
                Optional: true,
                Computed: true,
            },
            "agent_version": schema.StringAttribute{
                MarkdownDescription: "Version of the OneUptime agent reporting this resource.",
                Optional: true,
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "When telemetry was last received for this resource.",
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "retain_telemetry_data_for_days": schema.NumberAttribute{
                MarkdownDescription: "Number of days to retain telemetry data for this resource. Leave blank to use the project-wide default.",
                Optional: true,
                Computed: true,
            },
            "telemetry_retention_config": schema.StringAttribute{
                MarkdownDescription: "Per-pillar retention overrides for this resource. Unset fields fall back to the resource default, then the project's retention settings. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Is this cloud resource archived? Archived cloud resources are hidden from lists but keep collecting telemetry.",
                Optional: true,
                Computed: true,
            },
            "archived_at": schema.StringAttribute{
                MarkdownDescription: "When was this cloud resource archived?",
                Computed: true,
            },
            "archived_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "auto_archived_at": schema.StringAttribute{
                MarkdownDescription: "For a resource: when it was archived automatically for sending no metrics for the auto-archive period. Cleared when it reports again.",
                Computed: true,
            },
        },
    }
}

func (d *CloudResourceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CloudResourceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data CloudResourceDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.Slug.IsNull() && !data.Slug.IsUnknown() {
        filters["slug"] = data.Slug.ValueString()
        filterNames = append(filterNames, "slug = "+fmt.Sprintf("%q", data.Slug.ValueString()))
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.ResourceIdentifier.IsNull() && !data.ResourceIdentifier.IsUnknown() {
        filters["resourceIdentifier"] = data.ResourceIdentifier.ValueString()
        filterNames = append(filterNames, "resource_identifier = "+fmt.Sprintf("%q", data.ResourceIdentifier.ValueString()))
    }
    if !data.CloudPlatform.IsNull() && !data.CloudPlatform.IsUnknown() {
        filters["cloudPlatform"] = data.CloudPlatform.ValueString()
        filterNames = append(filterNames, "cloud_platform = "+fmt.Sprintf("%q", data.CloudPlatform.ValueString()))
    }
    if !data.CloudProvider.IsNull() && !data.CloudProvider.IsUnknown() {
        filters["cloudProvider"] = data.CloudProvider.ValueString()
        filterNames = append(filterNames, "cloud_provider = "+fmt.Sprintf("%q", data.CloudProvider.ValueString()))
    }
    if !data.CloudRegion.IsNull() && !data.CloudRegion.IsUnknown() {
        filters["cloudRegion"] = data.CloudRegion.ValueString()
        filterNames = append(filterNames, "cloud_region = "+fmt.Sprintf("%q", data.CloudRegion.ValueString()))
    }
    if !data.CloudAccountId.IsNull() && !data.CloudAccountId.IsUnknown() {
        filters["cloudAccountId"] = data.CloudAccountId.ValueString()
        filterNames = append(filterNames, "cloud_account_id = "+fmt.Sprintf("%q", data.CloudAccountId.ValueString()))
    }
    if !data.CloudResourceKind.IsNull() && !data.CloudResourceKind.IsUnknown() {
        filters["cloudResourceKind"] = data.CloudResourceKind.ValueString()
        filterNames = append(filterNames, "cloud_resource_kind = "+fmt.Sprintf("%q", data.CloudResourceKind.ValueString()))
    }
    if !data.CloudResourceType.IsNull() && !data.CloudResourceType.IsUnknown() {
        filters["cloudResourceType"] = data.CloudResourceType.ValueString()
        filterNames = append(filterNames, "cloud_resource_type = "+fmt.Sprintf("%q", data.CloudResourceType.ValueString()))
    }
    if !data.ProviderResourceId.IsNull() && !data.ProviderResourceId.IsUnknown() {
        filters["providerResourceId"] = data.ProviderResourceId.ValueString()
        filterNames = append(filterNames, "provider_resource_id = "+fmt.Sprintf("%q", data.ProviderResourceId.ValueString()))
    }
    if !data.CloudResourceGroup.IsNull() && !data.CloudResourceGroup.IsUnknown() {
        filters["cloudResourceGroup"] = data.CloudResourceGroup.ValueString()
        filterNames = append(filterNames, "cloud_resource_group = "+fmt.Sprintf("%q", data.CloudResourceGroup.ValueString()))
    }
    if !data.RuntimeName.IsNull() && !data.RuntimeName.IsUnknown() {
        filters["runtimeName"] = data.RuntimeName.ValueString()
        filterNames = append(filterNames, "runtime_name = "+fmt.Sprintf("%q", data.RuntimeName.ValueString()))
    }
    if !data.RuntimeVersion.IsNull() && !data.RuntimeVersion.IsUnknown() {
        filters["runtimeVersion"] = data.RuntimeVersion.ValueString()
        filterNames = append(filterNames, "runtime_version = "+fmt.Sprintf("%q", data.RuntimeVersion.ValueString()))
    }
    if !data.OtelCollectorStatus.IsNull() && !data.OtelCollectorStatus.IsUnknown() {
        filters["otelCollectorStatus"] = data.OtelCollectorStatus.ValueString()
        filterNames = append(filterNames, "otel_collector_status = "+fmt.Sprintf("%q", data.OtelCollectorStatus.ValueString()))
    }
    if !data.AgentVersion.IsNull() && !data.AgentVersion.IsUnknown() {
        filters["agentVersion"] = data.AgentVersion.ValueString()
        filterNames = append(filterNames, "agent_version = "+fmt.Sprintf("%q", data.AgentVersion.ValueString()))
    }
    if !data.RetainTelemetryDataForDays.IsNull() && !data.RetainTelemetryDataForDays.IsUnknown() {
        filters["retainTelemetryDataForDays"] = lookupNumber(data.RetainTelemetryDataForDays)
        filterNames = append(filterNames, "retain_telemetry_data_for_days = "+data.RetainTelemetryDataForDays.ValueBigFloat().String())
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.IsArchived.IsNull() && !data.IsArchived.IsUnknown() {
        filters["isArchived"] = data.IsArchived.ValueBool()
        filterNames = append(filterNames, "is_archived = "+fmt.Sprintf("%t", data.IsArchived.ValueBool()))
    }
    if !data.ArchivedByUserId.IsNull() && !data.ArchivedByUserId.IsUnknown() {
        filters["archivedByUserId"] = data.ArchivedByUserId.ValueString()
        filterNames = append(filterNames, "archived_by_user_id = "+fmt.Sprintf("%q", data.ArchivedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the cloud resource up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the cloud resource up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "name": true,
        "slug": true,
        "description": true,
        "resourceIdentifier": true,
        "cloudPlatform": true,
        "cloudProvider": true,
        "cloudRegion": true,
        "cloudAccountId": true,
        "cloudResourceKind": true,
        "cloudResourceType": true,
        "providerResourceId": true,
        "cloudResourceGroup": true,
        "telemetryAttributes": true,
        "runtimeName": true,
        "runtimeVersion": true,
        "otelCollectorStatus": true,
        "agentVersion": true,
        "lastSeenAt": true,
        "labels": true,
        "retainTelemetryDataForDays": true,
        "telemetryRetentionConfig": true,
        "createdByUserId": true,
        "isArchived": true,
        "archivedAt": true,
        "archivedByUserId": true,
        "autoArchivedAt": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/cloud-resource/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read cloud_resource, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No cloud resource found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read cloud_resource: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/cloud-resource/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list cloud_resource, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list cloud_resource: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No cloud resource matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one cloud resource matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for cloud_resource.")
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
    if obj, ok := item["resourceIdentifier"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResourceIdentifier = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ResourceIdentifier = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ResourceIdentifier = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ResourceIdentifier = types.StringValue(string(jsonBytes))
        } else {
            data.ResourceIdentifier = types.StringNull()
        }
    } else if val, ok := item["resourceIdentifier"].(string); ok {
        data.ResourceIdentifier = types.StringValue(val)
    } else {
        data.ResourceIdentifier = types.StringNull()
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
    if obj, ok := item["cloudResourceKind"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudResourceKind = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CloudResourceKind = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CloudResourceKind = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CloudResourceKind = types.StringValue(string(jsonBytes))
        } else {
            data.CloudResourceKind = types.StringNull()
        }
    } else if val, ok := item["cloudResourceKind"].(string); ok {
        data.CloudResourceKind = types.StringValue(val)
    } else {
        data.CloudResourceKind = types.StringNull()
    }
    if obj, ok := item["cloudResourceType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudResourceType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CloudResourceType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CloudResourceType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CloudResourceType = types.StringValue(string(jsonBytes))
        } else {
            data.CloudResourceType = types.StringNull()
        }
    } else if val, ok := item["cloudResourceType"].(string); ok {
        data.CloudResourceType = types.StringValue(val)
    } else {
        data.CloudResourceType = types.StringNull()
    }
    if obj, ok := item["providerResourceId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProviderResourceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProviderResourceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProviderResourceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProviderResourceId = types.StringValue(string(jsonBytes))
        } else {
            data.ProviderResourceId = types.StringNull()
        }
    } else if val, ok := item["providerResourceId"].(string); ok {
        data.ProviderResourceId = types.StringValue(val)
    } else {
        data.ProviderResourceId = types.StringNull()
    }
    if obj, ok := item["cloudResourceGroup"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudResourceGroup = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CloudResourceGroup = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CloudResourceGroup = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CloudResourceGroup = types.StringValue(string(jsonBytes))
        } else {
            data.CloudResourceGroup = types.StringNull()
        }
    } else if val, ok := item["cloudResourceGroup"].(string); ok {
        data.CloudResourceGroup = types.StringValue(val)
    } else {
        data.CloudResourceGroup = types.StringNull()
    }
    if obj, ok := item["telemetryAttributes"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TelemetryAttributes = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TelemetryAttributes = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TelemetryAttributes = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TelemetryAttributes = types.StringValue(string(jsonBytes))
        } else {
            data.TelemetryAttributes = types.StringNull()
        }
    } else if val, ok := item["telemetryAttributes"].(string); ok {
        data.TelemetryAttributes = types.StringValue(val)
    } else {
        data.TelemetryAttributes = types.StringNull()
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
    if obj, ok := item["otelCollectorStatus"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OtelCollectorStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OtelCollectorStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OtelCollectorStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OtelCollectorStatus = types.StringValue(string(jsonBytes))
        } else {
            data.OtelCollectorStatus = types.StringNull()
        }
    } else if val, ok := item["otelCollectorStatus"].(string); ok {
        data.OtelCollectorStatus = types.StringValue(val)
    } else {
        data.OtelCollectorStatus = types.StringNull()
    }
    if obj, ok := item["agentVersion"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AgentVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AgentVersion = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AgentVersion = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AgentVersion = types.StringValue(string(jsonBytes))
        } else {
            data.AgentVersion = types.StringNull()
        }
    } else if val, ok := item["agentVersion"].(string); ok {
        data.AgentVersion = types.StringValue(val)
    } else {
        data.AgentVersion = types.StringNull()
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
    if obj, ok := item["autoArchivedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AutoArchivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AutoArchivedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AutoArchivedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AutoArchivedAt = types.StringValue(string(jsonBytes))
        } else {
            data.AutoArchivedAt = types.StringNull()
        }
    } else if val, ok := item["autoArchivedAt"].(string); ok {
        data.AutoArchivedAt = types.StringValue(val)
    } else {
        data.AutoArchivedAt = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
