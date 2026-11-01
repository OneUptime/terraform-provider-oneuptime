package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DataStoreDataSource{}

func NewDataStoreDataSource() datasource.DataSource {
    return &DataStoreDataSource{}
}

// DataStoreDataSource defines the data source implementation.
type DataStoreDataSource struct {
    client *Client
}

// DataStoreDataSourceModel describes the data source data model.
type DataStoreDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    StoreKey types.String `tfsdk:"store_key"`
    Kind types.String `tfsdk:"kind"`
    ProviderValue types.String `tfsdk:"provider_value"`
    Family types.String `tfsdk:"family"`
    DiscoveredFrom types.String `tfsdk:"discovered_from"`
    CloudAccountId types.String `tfsdk:"cloud_account_id"`
    CloudResourceId types.String `tfsdk:"cloud_resource_id"`
    CloudConnectionId types.String `tfsdk:"cloud_connection_id"`
    DatabaseServerId types.String `tfsdk:"database_server_id"`
    ServiceId types.String `tfsdk:"service_id"`
    ProviderResourceId types.String `tfsdk:"provider_resource_id"`
    CloudAccountNumber types.String `tfsdk:"cloud_account_number"`
    Region types.String `tfsdk:"region"`
    Environment types.String `tfsdk:"environment"`
    Sensitivity types.String `tfsdk:"sensitivity"`
    SensitivityOverride types.String `tfsdk:"sensitivity_override"`
    DataTypeIds types.String `tfsdk:"data_type_ids"`
    DataCategories types.String `tfsdk:"data_categories"`
    SubjectCountries types.String `tfsdk:"subject_countries"`
    Exposure types.String `tfsdk:"exposure"`
    Posture types.String `tfsdk:"posture"`
    CloudTags types.String `tfsdk:"cloud_tags"`
    SizeInBytes types.Number `tfsdk:"size_in_bytes"`
    ObjectCount types.Number `tfsdk:"object_count"`
    EstimatedRecords types.Number `tfsdk:"estimated_records"`
    LastActivityAt types.String `tfsdk:"last_activity_at"`
    IsPublicByDesign types.Bool `tfsdk:"is_public_by_design"`
    IsAiData types.Bool `tfsdk:"is_ai_data"`
    IsContentScanEnabled types.Bool `tfsdk:"is_content_scan_enabled"`
    ScanProbeId types.String `tfsdk:"scan_probe_id"`
    ProbeConnectionVariable types.String `tfsdk:"probe_connection_variable"`
    ScanExclusions types.String `tfsdk:"scan_exclusions"`
    Purpose types.String `tfsdk:"purpose"`
    Retention types.String `tfsdk:"retention"`
    ConfigurationReadAt types.String `tfsdk:"configuration_read_at"`
    ContentScannedAt types.String `tfsdk:"content_scanned_at"`
    ContentScanStatus types.String `tfsdk:"content_scan_status"`
    NextContentScanAt types.String `tfsdk:"next_content_scan_at"`
    FindingCount types.Number `tfsdk:"finding_count"`
    OpenIssueCount types.Number `tfsdk:"open_issue_count"`
    HighestIssueSeverity types.String `tfsdk:"highest_issue_severity"`
    RiskScore types.Number `tfsdk:"risk_score"`
    IsArchived types.Bool `tfsdk:"is_archived"`
    FirstSeenAt types.String `tfsdk:"first_seen_at"`
    LastSeenAt types.String `tfsdk:"last_seen_at"`
    Labels types.Set `tfsdk:"labels"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *DataStoreDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_data_store"
}

func (d *DataStoreDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Every bucket, database, warehouse and telemetry stream Data Security knows: what sensitive data it holds, who can reach it and how it is configured. Look up an existing data store by `id`, or by any of its other arguments (`name`, `cloud_account_id`, `cloud_account_number`, ...): each one set must match, and exactly one data store may match them all.",

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
                MarkdownDescription: "ID of the project this data store belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "The store's name: a bucket, database, dataset or service name.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "What the store is for, in your team's words.",
                Optional: true,
                Computed: true,
            },
            "store_key": schema.StringAttribute{
                MarkdownDescription: "What identifies the store within the project: its provider, kind and native id. One row per store, however often it is discovered.",
                Optional: true,
                Computed: true,
            },
            "kind": schema.StringAttribute{
                MarkdownDescription: "What kind of store it is: an S3 bucket, an RDS instance, a PostgreSQL database, a service's logs.",
                Optional: true,
                Computed: true,
            },
            "provider_value": schema.StringAttribute{
                MarkdownDescription: "Where the store lives: aws, azure, gcp, self-managed or oneuptime.",
                Optional: true,
                Computed: true,
            },
            "family": schema.StringAttribute{
                MarkdownDescription: "The store's family: object storage, database, warehouse, search, cache, file system, snapshot or telemetry.",
                Optional: true,
                Computed: true,
            },
            "discovered_from": schema.StringAttribute{
                MarkdownDescription: "How OneUptime found the store: cloud inventory, the Database product, telemetry, the service map, or added by hand.",
                Optional: true,
                Computed: true,
            },
            "cloud_account_id": schema.StringAttribute{
                MarkdownDescription: "The cloud account the store belongs to, for a cloud store. The ID of a `oneuptime_cloud_account`.",
                Optional: true,
                Computed: true,
            },
            "cloud_resource_id": schema.StringAttribute{
                MarkdownDescription: "The cloud resource the store was discovered from, if any. The ID of a `oneuptime_cloud_resource`.",
                Optional: true,
                Computed: true,
            },
            "cloud_connection_id": schema.StringAttribute{
                MarkdownDescription: "The cloud connection that reads the store's configuration, if any. The ID of a `oneuptime_cloud_connection`.",
                Optional: true,
                Computed: true,
            },
            "database_server_id": schema.StringAttribute{
                MarkdownDescription: "The Database product server the store is, for a self-managed database. The ID of a `oneuptime_database`.",
                Optional: true,
                Computed: true,
            },
            "service_id": schema.StringAttribute{
                MarkdownDescription: "The telemetry service whose logs, traces, exceptions, replays or AI conversations the store is. The ID of a `oneuptime_service`.",
                Optional: true,
                Computed: true,
            },
            "provider_resource_id": schema.StringAttribute{
                MarkdownDescription: "The store's id in its cloud: an ARN, an Azure resource id or a Google Cloud resource name.",
                Optional: true,
                Computed: true,
            },
            "cloud_account_number": schema.StringAttribute{
                MarkdownDescription: "The cloud's own id for the account: an AWS account id, an Azure subscription id or a Google Cloud project id.",
                Optional: true,
                Computed: true,
            },
            "region": schema.StringAttribute{
                MarkdownDescription: "The region the store lives in.",
                Optional: true,
                Computed: true,
            },
            "environment": schema.StringAttribute{
                MarkdownDescription: "Which environment the store serves: Production, Staging, Development, Test or Unknown. Read from its tags and name; people may set it.",
                Optional: true,
                Computed: true,
            },
            "sensitivity": schema.StringAttribute{
                MarkdownDescription: "How sensitive the store is: Restricted, Confidential, Internal or Public, from the most sensitive data type found in it.",
                Optional: true,
                Computed: true,
            },
            "sensitivity_override": schema.StringAttribute{
                MarkdownDescription: "A sensitivity people set by hand, which the store is treated as instead of what the scans found.",
                Optional: true,
                Computed: true,
            },
            "data_type_ids": schema.StringAttribute{
                MarkdownDescription: "The data types found in the store with at least Medium confidence. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "data_categories": schema.StringAttribute{
                MarkdownDescription: "The categories of those data types. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "subject_countries": schema.StringAttribute{
                MarkdownDescription: "Countries whose people's data the store holds, from country-bound data types (ISO 3166-1 alpha-2). A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "exposure": schema.StringAttribute{
                MarkdownDescription: "Where the store sits relative to the internet and other organizations: Public, Other accounts, Internet reachable, Private or Unknown.",
                Optional: true,
                Computed: true,
            },
            "posture": schema.StringAttribute{
                MarkdownDescription: "What the store's configuration says about its safety: public access, encryption, backups, logging and the rest. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "cloud_tags": schema.StringAttribute{
                MarkdownDescription: "The store's tags or labels in its cloud. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "size_in_bytes": schema.NumberAttribute{
                MarkdownDescription: "How much data the store holds, from the metrics already collected.",
                Optional: true,
                Computed: true,
            },
            "object_count": schema.NumberAttribute{
                MarkdownDescription: "How many objects, rows or documents the store holds, when known.",
                Optional: true,
                Computed: true,
            },
            "estimated_records": schema.NumberAttribute{
                MarkdownDescription: "An estimate of how many records with sensitive data the store holds, from its samples.",
                Optional: true,
                Computed: true,
            },
            "last_activity_at": schema.StringAttribute{
                MarkdownDescription: "When the store was last read or written, when its metrics or traces say.",
                Computed: true,
            },
            "is_public_by_design": schema.BoolAttribute{
                MarkdownDescription: "The store is meant to be public, such as a website's assets. Public access checks pass for it.",
                Optional: true,
                Computed: true,
            },
            "is_ai_data": schema.BoolAttribute{
                MarkdownDescription: "The store holds data used to train or ground an AI model.",
                Optional: true,
                Computed: true,
            },
            "is_content_scan_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether Data Security samples the store's content. Configuration checks run either way.",
                Optional: true,
                Computed: true,
            },
            "scan_probe_id": schema.StringAttribute{
                MarkdownDescription: "A probe that scans the store's content from inside your network instead of OneUptime's workers. The ID of a `oneuptime_probe`.",
                Optional: true,
                Computed: true,
            },
            "probe_connection_variable": schema.StringAttribute{
                MarkdownDescription: "For a database scanned from a probe: the environment variable on that probe that holds its connection string. The connection string itself never reaches OneUptime.",
                Optional: true,
                Computed: true,
            },
            "scan_exclusions": schema.StringAttribute{
                MarkdownDescription: "Path prefixes, tables or attribute keys the content scan skips. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "purpose": schema.StringAttribute{
                MarkdownDescription: "Why the store holds personal data, for the GDPR record of processing.",
                Optional: true,
                Computed: true,
            },
            "retention": schema.StringAttribute{
                MarkdownDescription: "How long data is kept in the store, for the GDPR record of processing.",
                Optional: true,
                Computed: true,
            },
            "configuration_read_at": schema.StringAttribute{
                MarkdownDescription: "When the store's configuration and access were last read.",
                Computed: true,
            },
            "content_scanned_at": schema.StringAttribute{
                MarkdownDescription: "When the store's content was last sampled.",
                Computed: true,
            },
            "content_scan_status": schema.StringAttribute{
                MarkdownDescription: "How the last content scan ended, or why it could not run.",
                Optional: true,
                Computed: true,
            },
            "next_content_scan_at": schema.StringAttribute{
                MarkdownDescription: "When the store's content is next sampled.",
                Computed: true,
            },
            "finding_count": schema.NumberAttribute{
                MarkdownDescription: "How many data types were found at how many locations.",
                Optional: true,
                Computed: true,
            },
            "open_issue_count": schema.NumberAttribute{
                MarkdownDescription: "How many of the store's issues are open.",
                Optional: true,
                Computed: true,
            },
            "highest_issue_severity": schema.StringAttribute{
                MarkdownDescription: "The severity of the store's worst open issue.",
                Optional: true,
                Computed: true,
            },
            "risk_score": schema.NumberAttribute{
                MarkdownDescription: "The risk score of the store's worst open issue, from 0 to 100.",
                Optional: true,
                Computed: true,
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "The store is gone from its cloud or source; it is kept with its history.",
                Optional: true,
                Computed: true,
            },
            "first_seen_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime first found the store.",
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime last found the store.",
                Computed: true,
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the user who created this data store. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *DataStoreDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataStoreDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data DataStoreDataSourceModel

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
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.StoreKey.IsNull() && !data.StoreKey.IsUnknown() {
        filters["storeKey"] = data.StoreKey.ValueString()
        filterNames = append(filterNames, "store_key = "+fmt.Sprintf("%q", data.StoreKey.ValueString()))
    }
    if !data.Kind.IsNull() && !data.Kind.IsUnknown() {
        filters["kind"] = data.Kind.ValueString()
        filterNames = append(filterNames, "kind = "+fmt.Sprintf("%q", data.Kind.ValueString()))
    }
    if !data.ProviderValue.IsNull() && !data.ProviderValue.IsUnknown() {
        filters["provider"] = data.ProviderValue.ValueString()
        filterNames = append(filterNames, "provider_value = "+fmt.Sprintf("%q", data.ProviderValue.ValueString()))
    }
    if !data.Family.IsNull() && !data.Family.IsUnknown() {
        filters["family"] = data.Family.ValueString()
        filterNames = append(filterNames, "family = "+fmt.Sprintf("%q", data.Family.ValueString()))
    }
    if !data.DiscoveredFrom.IsNull() && !data.DiscoveredFrom.IsUnknown() {
        filters["discoveredFrom"] = data.DiscoveredFrom.ValueString()
        filterNames = append(filterNames, "discovered_from = "+fmt.Sprintf("%q", data.DiscoveredFrom.ValueString()))
    }
    if !data.CloudAccountId.IsNull() && !data.CloudAccountId.IsUnknown() {
        filters["cloudAccountId"] = data.CloudAccountId.ValueString()
        filterNames = append(filterNames, "cloud_account_id = "+fmt.Sprintf("%q", data.CloudAccountId.ValueString()))
    }
    if !data.CloudResourceId.IsNull() && !data.CloudResourceId.IsUnknown() {
        filters["cloudResourceId"] = data.CloudResourceId.ValueString()
        filterNames = append(filterNames, "cloud_resource_id = "+fmt.Sprintf("%q", data.CloudResourceId.ValueString()))
    }
    if !data.CloudConnectionId.IsNull() && !data.CloudConnectionId.IsUnknown() {
        filters["cloudConnectionId"] = data.CloudConnectionId.ValueString()
        filterNames = append(filterNames, "cloud_connection_id = "+fmt.Sprintf("%q", data.CloudConnectionId.ValueString()))
    }
    if !data.DatabaseServerId.IsNull() && !data.DatabaseServerId.IsUnknown() {
        filters["databaseServerId"] = data.DatabaseServerId.ValueString()
        filterNames = append(filterNames, "database_server_id = "+fmt.Sprintf("%q", data.DatabaseServerId.ValueString()))
    }
    if !data.ServiceId.IsNull() && !data.ServiceId.IsUnknown() {
        filters["serviceId"] = data.ServiceId.ValueString()
        filterNames = append(filterNames, "service_id = "+fmt.Sprintf("%q", data.ServiceId.ValueString()))
    }
    if !data.ProviderResourceId.IsNull() && !data.ProviderResourceId.IsUnknown() {
        filters["providerResourceId"] = data.ProviderResourceId.ValueString()
        filterNames = append(filterNames, "provider_resource_id = "+fmt.Sprintf("%q", data.ProviderResourceId.ValueString()))
    }
    if !data.CloudAccountNumber.IsNull() && !data.CloudAccountNumber.IsUnknown() {
        filters["cloudAccountNumber"] = data.CloudAccountNumber.ValueString()
        filterNames = append(filterNames, "cloud_account_number = "+fmt.Sprintf("%q", data.CloudAccountNumber.ValueString()))
    }
    if !data.Region.IsNull() && !data.Region.IsUnknown() {
        filters["region"] = data.Region.ValueString()
        filterNames = append(filterNames, "region = "+fmt.Sprintf("%q", data.Region.ValueString()))
    }
    if !data.Environment.IsNull() && !data.Environment.IsUnknown() {
        filters["environment"] = data.Environment.ValueString()
        filterNames = append(filterNames, "environment = "+fmt.Sprintf("%q", data.Environment.ValueString()))
    }
    if !data.Sensitivity.IsNull() && !data.Sensitivity.IsUnknown() {
        filters["sensitivity"] = data.Sensitivity.ValueString()
        filterNames = append(filterNames, "sensitivity = "+fmt.Sprintf("%q", data.Sensitivity.ValueString()))
    }
    if !data.SensitivityOverride.IsNull() && !data.SensitivityOverride.IsUnknown() {
        filters["sensitivityOverride"] = data.SensitivityOverride.ValueString()
        filterNames = append(filterNames, "sensitivity_override = "+fmt.Sprintf("%q", data.SensitivityOverride.ValueString()))
    }
    if !data.Exposure.IsNull() && !data.Exposure.IsUnknown() {
        filters["exposure"] = data.Exposure.ValueString()
        filterNames = append(filterNames, "exposure = "+fmt.Sprintf("%q", data.Exposure.ValueString()))
    }
    if !data.SizeInBytes.IsNull() && !data.SizeInBytes.IsUnknown() {
        filters["sizeInBytes"] = lookupNumber(data.SizeInBytes)
        filterNames = append(filterNames, "size_in_bytes = "+data.SizeInBytes.ValueBigFloat().String())
    }
    if !data.ObjectCount.IsNull() && !data.ObjectCount.IsUnknown() {
        filters["objectCount"] = lookupNumber(data.ObjectCount)
        filterNames = append(filterNames, "object_count = "+data.ObjectCount.ValueBigFloat().String())
    }
    if !data.EstimatedRecords.IsNull() && !data.EstimatedRecords.IsUnknown() {
        filters["estimatedRecords"] = lookupNumber(data.EstimatedRecords)
        filterNames = append(filterNames, "estimated_records = "+data.EstimatedRecords.ValueBigFloat().String())
    }
    if !data.IsPublicByDesign.IsNull() && !data.IsPublicByDesign.IsUnknown() {
        filters["isPublicByDesign"] = data.IsPublicByDesign.ValueBool()
        filterNames = append(filterNames, "is_public_by_design = "+fmt.Sprintf("%t", data.IsPublicByDesign.ValueBool()))
    }
    if !data.IsAiData.IsNull() && !data.IsAiData.IsUnknown() {
        filters["isAiData"] = data.IsAiData.ValueBool()
        filterNames = append(filterNames, "is_ai_data = "+fmt.Sprintf("%t", data.IsAiData.ValueBool()))
    }
    if !data.IsContentScanEnabled.IsNull() && !data.IsContentScanEnabled.IsUnknown() {
        filters["isContentScanEnabled"] = data.IsContentScanEnabled.ValueBool()
        filterNames = append(filterNames, "is_content_scan_enabled = "+fmt.Sprintf("%t", data.IsContentScanEnabled.ValueBool()))
    }
    if !data.ScanProbeId.IsNull() && !data.ScanProbeId.IsUnknown() {
        filters["scanProbeId"] = data.ScanProbeId.ValueString()
        filterNames = append(filterNames, "scan_probe_id = "+fmt.Sprintf("%q", data.ScanProbeId.ValueString()))
    }
    if !data.ProbeConnectionVariable.IsNull() && !data.ProbeConnectionVariable.IsUnknown() {
        filters["probeConnectionVariable"] = data.ProbeConnectionVariable.ValueString()
        filterNames = append(filterNames, "probe_connection_variable = "+fmt.Sprintf("%q", data.ProbeConnectionVariable.ValueString()))
    }
    if !data.Purpose.IsNull() && !data.Purpose.IsUnknown() {
        filters["purpose"] = data.Purpose.ValueString()
        filterNames = append(filterNames, "purpose = "+fmt.Sprintf("%q", data.Purpose.ValueString()))
    }
    if !data.Retention.IsNull() && !data.Retention.IsUnknown() {
        filters["retention"] = data.Retention.ValueString()
        filterNames = append(filterNames, "retention = "+fmt.Sprintf("%q", data.Retention.ValueString()))
    }
    if !data.ContentScanStatus.IsNull() && !data.ContentScanStatus.IsUnknown() {
        filters["contentScanStatus"] = data.ContentScanStatus.ValueString()
        filterNames = append(filterNames, "content_scan_status = "+fmt.Sprintf("%q", data.ContentScanStatus.ValueString()))
    }
    if !data.FindingCount.IsNull() && !data.FindingCount.IsUnknown() {
        filters["findingCount"] = lookupNumber(data.FindingCount)
        filterNames = append(filterNames, "finding_count = "+data.FindingCount.ValueBigFloat().String())
    }
    if !data.OpenIssueCount.IsNull() && !data.OpenIssueCount.IsUnknown() {
        filters["openIssueCount"] = lookupNumber(data.OpenIssueCount)
        filterNames = append(filterNames, "open_issue_count = "+data.OpenIssueCount.ValueBigFloat().String())
    }
    if !data.HighestIssueSeverity.IsNull() && !data.HighestIssueSeverity.IsUnknown() {
        filters["highestIssueSeverity"] = data.HighestIssueSeverity.ValueString()
        filterNames = append(filterNames, "highest_issue_severity = "+fmt.Sprintf("%q", data.HighestIssueSeverity.ValueString()))
    }
    if !data.RiskScore.IsNull() && !data.RiskScore.IsUnknown() {
        filters["riskScore"] = lookupNumber(data.RiskScore)
        filterNames = append(filterNames, "risk_score = "+data.RiskScore.ValueBigFloat().String())
    }
    if !data.IsArchived.IsNull() && !data.IsArchived.IsUnknown() {
        filters["isArchived"] = data.IsArchived.ValueBool()
        filterNames = append(filterNames, "is_archived = "+fmt.Sprintf("%t", data.IsArchived.ValueBool()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the data store up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the data store up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "name": true,
        "description": true,
        "storeKey": true,
        "kind": true,
        "provider": true,
        "family": true,
        "discoveredFrom": true,
        "cloudAccountId": true,
        "cloudResourceId": true,
        "cloudConnectionId": true,
        "databaseServerId": true,
        "serviceId": true,
        "providerResourceId": true,
        "cloudAccountNumber": true,
        "region": true,
        "environment": true,
        "sensitivity": true,
        "sensitivityOverride": true,
        "dataTypeIds": true,
        "dataCategories": true,
        "subjectCountries": true,
        "exposure": true,
        "posture": true,
        "cloudTags": true,
        "sizeInBytes": true,
        "objectCount": true,
        "estimatedRecords": true,
        "lastActivityAt": true,
        "isPublicByDesign": true,
        "isAiData": true,
        "isContentScanEnabled": true,
        "scanProbeId": true,
        "probeConnectionVariable": true,
        "scanExclusions": true,
        "purpose": true,
        "retention": true,
        "configurationReadAt": true,
        "contentScannedAt": true,
        "contentScanStatus": true,
        "nextContentScanAt": true,
        "findingCount": true,
        "openIssueCount": true,
        "highestIssueSeverity": true,
        "riskScore": true,
        "isArchived": true,
        "firstSeenAt": true,
        "lastSeenAt": true,
        "labels": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/data-store/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read data_store, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data store found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read data_store: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/data-store/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list data_store, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list data_store: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No data store matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one data store matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for data_store.")
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
    if obj, ok := item["storeKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StoreKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StoreKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StoreKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StoreKey = types.StringValue(string(jsonBytes))
        } else {
            data.StoreKey = types.StringNull()
        }
    } else if val, ok := item["storeKey"].(string); ok {
        data.StoreKey = types.StringValue(val)
    } else {
        data.StoreKey = types.StringNull()
    }
    if obj, ok := item["kind"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Kind = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Kind = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Kind = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Kind = types.StringValue(string(jsonBytes))
        } else {
            data.Kind = types.StringNull()
        }
    } else if val, ok := item["kind"].(string); ok {
        data.Kind = types.StringValue(val)
    } else {
        data.Kind = types.StringNull()
    }
    if obj, ok := item["provider"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProviderValue = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProviderValue = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProviderValue = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProviderValue = types.StringValue(string(jsonBytes))
        } else {
            data.ProviderValue = types.StringNull()
        }
    } else if val, ok := item["provider"].(string); ok {
        data.ProviderValue = types.StringValue(val)
    } else {
        data.ProviderValue = types.StringNull()
    }
    if obj, ok := item["family"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Family = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Family = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Family = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Family = types.StringValue(string(jsonBytes))
        } else {
            data.Family = types.StringNull()
        }
    } else if val, ok := item["family"].(string); ok {
        data.Family = types.StringValue(val)
    } else {
        data.Family = types.StringNull()
    }
    if obj, ok := item["discoveredFrom"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DiscoveredFrom = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DiscoveredFrom = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DiscoveredFrom = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DiscoveredFrom = types.StringValue(string(jsonBytes))
        } else {
            data.DiscoveredFrom = types.StringNull()
        }
    } else if val, ok := item["discoveredFrom"].(string); ok {
        data.DiscoveredFrom = types.StringValue(val)
    } else {
        data.DiscoveredFrom = types.StringNull()
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
    if obj, ok := item["cloudResourceId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudResourceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CloudResourceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CloudResourceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CloudResourceId = types.StringValue(string(jsonBytes))
        } else {
            data.CloudResourceId = types.StringNull()
        }
    } else if val, ok := item["cloudResourceId"].(string); ok {
        data.CloudResourceId = types.StringValue(val)
    } else {
        data.CloudResourceId = types.StringNull()
    }
    if obj, ok := item["cloudConnectionId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudConnectionId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CloudConnectionId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CloudConnectionId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CloudConnectionId = types.StringValue(string(jsonBytes))
        } else {
            data.CloudConnectionId = types.StringNull()
        }
    } else if val, ok := item["cloudConnectionId"].(string); ok {
        data.CloudConnectionId = types.StringValue(val)
    } else {
        data.CloudConnectionId = types.StringNull()
    }
    if obj, ok := item["databaseServerId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DatabaseServerId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DatabaseServerId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DatabaseServerId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DatabaseServerId = types.StringValue(string(jsonBytes))
        } else {
            data.DatabaseServerId = types.StringNull()
        }
    } else if val, ok := item["databaseServerId"].(string); ok {
        data.DatabaseServerId = types.StringValue(val)
    } else {
        data.DatabaseServerId = types.StringNull()
    }
    if obj, ok := item["serviceId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServiceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ServiceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ServiceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ServiceId = types.StringValue(string(jsonBytes))
        } else {
            data.ServiceId = types.StringNull()
        }
    } else if val, ok := item["serviceId"].(string); ok {
        data.ServiceId = types.StringValue(val)
    } else {
        data.ServiceId = types.StringNull()
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
    if obj, ok := item["cloudAccountNumber"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudAccountNumber = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CloudAccountNumber = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CloudAccountNumber = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CloudAccountNumber = types.StringValue(string(jsonBytes))
        } else {
            data.CloudAccountNumber = types.StringNull()
        }
    } else if val, ok := item["cloudAccountNumber"].(string); ok {
        data.CloudAccountNumber = types.StringValue(val)
    } else {
        data.CloudAccountNumber = types.StringNull()
    }
    if obj, ok := item["region"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Region = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Region = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Region = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Region = types.StringValue(string(jsonBytes))
        } else {
            data.Region = types.StringNull()
        }
    } else if val, ok := item["region"].(string); ok {
        data.Region = types.StringValue(val)
    } else {
        data.Region = types.StringNull()
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
    if obj, ok := item["sensitivity"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Sensitivity = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Sensitivity = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Sensitivity = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Sensitivity = types.StringValue(string(jsonBytes))
        } else {
            data.Sensitivity = types.StringNull()
        }
    } else if val, ok := item["sensitivity"].(string); ok {
        data.Sensitivity = types.StringValue(val)
    } else {
        data.Sensitivity = types.StringNull()
    }
    if obj, ok := item["sensitivityOverride"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SensitivityOverride = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SensitivityOverride = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SensitivityOverride = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SensitivityOverride = types.StringValue(string(jsonBytes))
        } else {
            data.SensitivityOverride = types.StringNull()
        }
    } else if val, ok := item["sensitivityOverride"].(string); ok {
        data.SensitivityOverride = types.StringValue(val)
    } else {
        data.SensitivityOverride = types.StringNull()
    }
    if obj, ok := item["dataTypeIds"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataTypeIds = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DataTypeIds = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DataTypeIds = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DataTypeIds = types.StringValue(string(jsonBytes))
        } else {
            data.DataTypeIds = types.StringNull()
        }
    } else if val, ok := item["dataTypeIds"].(string); ok {
        data.DataTypeIds = types.StringValue(val)
    } else {
        data.DataTypeIds = types.StringNull()
    }
    if obj, ok := item["dataCategories"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataCategories = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DataCategories = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DataCategories = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DataCategories = types.StringValue(string(jsonBytes))
        } else {
            data.DataCategories = types.StringNull()
        }
    } else if val, ok := item["dataCategories"].(string); ok {
        data.DataCategories = types.StringValue(val)
    } else {
        data.DataCategories = types.StringNull()
    }
    if obj, ok := item["subjectCountries"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubjectCountries = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SubjectCountries = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SubjectCountries = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SubjectCountries = types.StringValue(string(jsonBytes))
        } else {
            data.SubjectCountries = types.StringNull()
        }
    } else if val, ok := item["subjectCountries"].(string); ok {
        data.SubjectCountries = types.StringValue(val)
    } else {
        data.SubjectCountries = types.StringNull()
    }
    if obj, ok := item["exposure"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Exposure = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Exposure = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Exposure = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Exposure = types.StringValue(string(jsonBytes))
        } else {
            data.Exposure = types.StringNull()
        }
    } else if val, ok := item["exposure"].(string); ok {
        data.Exposure = types.StringValue(val)
    } else {
        data.Exposure = types.StringNull()
    }
    if obj, ok := item["posture"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Posture = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Posture = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Posture = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Posture = types.StringValue(string(jsonBytes))
        } else {
            data.Posture = types.StringNull()
        }
    } else if val, ok := item["posture"].(string); ok {
        data.Posture = types.StringValue(val)
    } else {
        data.Posture = types.StringNull()
    }
    if obj, ok := item["cloudTags"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudTags = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CloudTags = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CloudTags = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CloudTags = types.StringValue(string(jsonBytes))
        } else {
            data.CloudTags = types.StringNull()
        }
    } else if val, ok := item["cloudTags"].(string); ok {
        data.CloudTags = types.StringValue(val)
    } else {
        data.CloudTags = types.StringNull()
    }
    data.SizeInBytes = numberFromAPI(data.SizeInBytes, item["sizeInBytes"])
    data.ObjectCount = numberFromAPI(data.ObjectCount, item["objectCount"])
    data.EstimatedRecords = numberFromAPI(data.EstimatedRecords, item["estimatedRecords"])
    if obj, ok := item["lastActivityAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastActivityAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastActivityAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastActivityAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastActivityAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastActivityAt = types.StringNull()
        }
    } else if val, ok := item["lastActivityAt"].(string); ok {
        data.LastActivityAt = types.StringValue(val)
    } else {
        data.LastActivityAt = types.StringNull()
    }
    if val, ok := item["isPublicByDesign"].(bool); ok {
        data.IsPublicByDesign = types.BoolValue(val)
    } else {
        data.IsPublicByDesign = types.BoolNull()
    }
    if val, ok := item["isAiData"].(bool); ok {
        data.IsAiData = types.BoolValue(val)
    } else {
        data.IsAiData = types.BoolNull()
    }
    if val, ok := item["isContentScanEnabled"].(bool); ok {
        data.IsContentScanEnabled = types.BoolValue(val)
    } else {
        data.IsContentScanEnabled = types.BoolNull()
    }
    if obj, ok := item["scanProbeId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScanProbeId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ScanProbeId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ScanProbeId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ScanProbeId = types.StringValue(string(jsonBytes))
        } else {
            data.ScanProbeId = types.StringNull()
        }
    } else if val, ok := item["scanProbeId"].(string); ok {
        data.ScanProbeId = types.StringValue(val)
    } else {
        data.ScanProbeId = types.StringNull()
    }
    if obj, ok := item["probeConnectionVariable"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProbeConnectionVariable = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProbeConnectionVariable = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProbeConnectionVariable = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProbeConnectionVariable = types.StringValue(string(jsonBytes))
        } else {
            data.ProbeConnectionVariable = types.StringNull()
        }
    } else if val, ok := item["probeConnectionVariable"].(string); ok {
        data.ProbeConnectionVariable = types.StringValue(val)
    } else {
        data.ProbeConnectionVariable = types.StringNull()
    }
    if obj, ok := item["scanExclusions"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScanExclusions = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ScanExclusions = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ScanExclusions = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ScanExclusions = types.StringValue(string(jsonBytes))
        } else {
            data.ScanExclusions = types.StringNull()
        }
    } else if val, ok := item["scanExclusions"].(string); ok {
        data.ScanExclusions = types.StringValue(val)
    } else {
        data.ScanExclusions = types.StringNull()
    }
    if obj, ok := item["purpose"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Purpose = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Purpose = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Purpose = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Purpose = types.StringValue(string(jsonBytes))
        } else {
            data.Purpose = types.StringNull()
        }
    } else if val, ok := item["purpose"].(string); ok {
        data.Purpose = types.StringValue(val)
    } else {
        data.Purpose = types.StringNull()
    }
    if obj, ok := item["retention"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Retention = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Retention = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Retention = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Retention = types.StringValue(string(jsonBytes))
        } else {
            data.Retention = types.StringNull()
        }
    } else if val, ok := item["retention"].(string); ok {
        data.Retention = types.StringValue(val)
    } else {
        data.Retention = types.StringNull()
    }
    if obj, ok := item["configurationReadAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ConfigurationReadAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ConfigurationReadAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ConfigurationReadAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ConfigurationReadAt = types.StringValue(string(jsonBytes))
        } else {
            data.ConfigurationReadAt = types.StringNull()
        }
    } else if val, ok := item["configurationReadAt"].(string); ok {
        data.ConfigurationReadAt = types.StringValue(val)
    } else {
        data.ConfigurationReadAt = types.StringNull()
    }
    if obj, ok := item["contentScannedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ContentScannedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ContentScannedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ContentScannedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ContentScannedAt = types.StringValue(string(jsonBytes))
        } else {
            data.ContentScannedAt = types.StringNull()
        }
    } else if val, ok := item["contentScannedAt"].(string); ok {
        data.ContentScannedAt = types.StringValue(val)
    } else {
        data.ContentScannedAt = types.StringNull()
    }
    if obj, ok := item["contentScanStatus"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ContentScanStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ContentScanStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ContentScanStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ContentScanStatus = types.StringValue(string(jsonBytes))
        } else {
            data.ContentScanStatus = types.StringNull()
        }
    } else if val, ok := item["contentScanStatus"].(string); ok {
        data.ContentScanStatus = types.StringValue(val)
    } else {
        data.ContentScanStatus = types.StringNull()
    }
    if obj, ok := item["nextContentScanAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.NextContentScanAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.NextContentScanAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.NextContentScanAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.NextContentScanAt = types.StringValue(string(jsonBytes))
        } else {
            data.NextContentScanAt = types.StringNull()
        }
    } else if val, ok := item["nextContentScanAt"].(string); ok {
        data.NextContentScanAt = types.StringValue(val)
    } else {
        data.NextContentScanAt = types.StringNull()
    }
    data.FindingCount = numberFromAPI(data.FindingCount, item["findingCount"])
    data.OpenIssueCount = numberFromAPI(data.OpenIssueCount, item["openIssueCount"])
    if obj, ok := item["highestIssueSeverity"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HighestIssueSeverity = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.HighestIssueSeverity = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.HighestIssueSeverity = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.HighestIssueSeverity = types.StringValue(string(jsonBytes))
        } else {
            data.HighestIssueSeverity = types.StringNull()
        }
    } else if val, ok := item["highestIssueSeverity"].(string); ok {
        data.HighestIssueSeverity = types.StringValue(val)
    } else {
        data.HighestIssueSeverity = types.StringNull()
    }
    data.RiskScore = numberFromAPI(data.RiskScore, item["riskScore"])
    if val, ok := item["isArchived"].(bool); ok {
        data.IsArchived = types.BoolValue(val)
    } else {
        data.IsArchived = types.BoolNull()
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

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
