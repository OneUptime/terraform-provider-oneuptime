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
    "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DataStoreResource{}
var _ resource.ResourceWithImportState = &DataStoreResource{}
var _ resource.ResourceWithModifyPlan = &DataStoreResource{}

func NewDataStoreResource() resource.Resource {
    return &DataStoreResource{}
}

// DataStoreResource defines the resource implementation.
type DataStoreResource struct {
    client *Client
}

// DataStoreResourceModel describes the resource data model.
type DataStoreResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    Kind types.String `tfsdk:"kind"`
    Environment types.String `tfsdk:"environment"`
    SensitivityOverride types.String `tfsdk:"sensitivity_override"`
    IsPublicByDesign types.Bool `tfsdk:"is_public_by_design"`
    IsAiData types.Bool `tfsdk:"is_ai_data"`
    IsContentScanEnabled types.Bool `tfsdk:"is_content_scan_enabled"`
    ScanProbeId types.String `tfsdk:"scan_probe_id"`
    ProbeConnectionVariable types.String `tfsdk:"probe_connection_variable"`
    ConnectionString types.String `tfsdk:"connection_string"`
    ScanExclusions JSONSubsetValue `tfsdk:"scan_exclusions"`
    Purpose types.String `tfsdk:"purpose"`
    Retention types.String `tfsdk:"retention"`
    Labels types.Set `tfsdk:"labels"`
    CreatedAt RFC3339Value `tfsdk:"created_at"`
    UpdatedAt RFC3339Value `tfsdk:"updated_at"`
    StoreKey types.String `tfsdk:"store_key"`
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
    Sensitivity types.String `tfsdk:"sensitivity"`
    DataTypeIds JSONSubsetValue `tfsdk:"data_type_ids"`
    DataCategories JSONSubsetValue `tfsdk:"data_categories"`
    SubjectCountries JSONSubsetValue `tfsdk:"subject_countries"`
    Exposure types.String `tfsdk:"exposure"`
    Posture JSONSubsetValue `tfsdk:"posture"`
    CloudTags JSONSubsetValue `tfsdk:"cloud_tags"`
    SizeInBytes types.Number `tfsdk:"size_in_bytes"`
    ObjectCount types.Number `tfsdk:"object_count"`
    EstimatedRecords types.Number `tfsdk:"estimated_records"`
    LastActivityAt RFC3339Value `tfsdk:"last_activity_at"`
    ConfigurationReadAt RFC3339Value `tfsdk:"configuration_read_at"`
    ContentScannedAt RFC3339Value `tfsdk:"content_scanned_at"`
    ContentScanStatus types.String `tfsdk:"content_scan_status"`
    NextContentScanAt RFC3339Value `tfsdk:"next_content_scan_at"`
    FindingCount types.Number `tfsdk:"finding_count"`
    OpenIssueCount types.Number `tfsdk:"open_issue_count"`
    HighestIssueSeverity types.String `tfsdk:"highest_issue_severity"`
    RiskScore types.Number `tfsdk:"risk_score"`
    IsArchived types.Bool `tfsdk:"is_archived"`
    FirstSeenAt RFC3339Value `tfsdk:"first_seen_at"`
    LastSeenAt RFC3339Value `tfsdk:"last_seen_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (r *DataStoreResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_data_store"
}

func (r *DataStoreResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = r.schemaDefinition()
}

func (r *DataStoreResource) schemaDefinition() schema.Schema {
    return schema.Schema{
        MarkdownDescription: "Every bucket, database, warehouse and telemetry stream Data Security knows: what sensitive data it holds, who can reach it and how it is configured.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of the project this data store belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "The store's name: a bucket, database, dataset or service name.",
                Required: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "What the store is for, in your team's words.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "kind": schema.StringAttribute{
                MarkdownDescription: "What kind of store it is: an S3 bucket, an RDS instance, a PostgreSQL database, a service's logs.",
                Required: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
                Validators: []validator.String{
                    stringvalidator.OneOf("AWS S3 Bucket", "AWS RDS Instance", "AWS Aurora Cluster", "AWS RDS Snapshot", "AWS DynamoDB Table", "AWS Redshift Cluster", "AWS DocumentDB Cluster", "AWS OpenSearch Domain", "AWS ElastiCache Cluster", "AWS EFS File System", "Azure Storage Account", "Azure SQL Database", "Azure Cosmos DB Account", "Azure PostgreSQL Server", "Azure MySQL Server", "Google Cloud Storage Bucket", "Google Cloud SQL Instance", "Google BigQuery Dataset", "Google Firestore Database", "Google Spanner Instance", "Google Bigtable Instance", "PostgreSQL", "MySQL", "MongoDB", "SQL Server", "Oracle", "Redis", "Elasticsearch", "Logs", "Traces", "Exceptions", "Session Replays", "AI Conversations"),
                },
            },
            "environment": schema.StringAttribute{
                MarkdownDescription: "Which environment the store serves: Production, Staging, Development, Test or Unknown. Read from its tags and name; people may set it.",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("Unknown"),
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    stringvalidator.OneOf("Production", "Staging", "Development", "Test", "Unknown"),
                },
            },
            "sensitivity_override": schema.StringAttribute{
                MarkdownDescription: "A sensitivity people set by hand, which the store is treated as instead of what the scans found.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    stringvalidator.OneOf("Restricted", "Confidential", "Internal", "Public"),
                },
            },
            "is_public_by_design": schema.BoolAttribute{
                MarkdownDescription: "The store is meant to be public, such as a website's assets. Public access checks pass for it.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "is_ai_data": schema.BoolAttribute{
                MarkdownDescription: "The store holds data used to train or ground an AI model.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "is_content_scan_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether Data Security samples the store's content. Configuration checks run either way.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "scan_probe_id": schema.StringAttribute{
                MarkdownDescription: "A probe that scans the store's content from inside your network instead of OneUptime's workers. The ID of a `oneuptime_probe`.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: recordIDValidators("oneuptime_probe"),
            },
            "probe_connection_variable": schema.StringAttribute{
                MarkdownDescription: "For a database scanned from a probe: the environment variable on that probe that holds its connection string. The connection string itself never reaches OneUptime.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "connection_string": schema.StringAttribute{
                MarkdownDescription: "For a database scanned by OneUptime's workers: a read-only connection string. Encrypted at rest and never returned by the API.",
                Optional: true,
                Sensitive: true,
            },
            "scan_exclusions": schema.StringAttribute{
                MarkdownDescription: "Path prefixes, tables or attribute keys the content scan skips. A JSON value: write it with `jsonencode()`.",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    JSONSubsetKeepImported(),
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "purpose": schema.StringAttribute{
                MarkdownDescription: "Why the store holds personal data, for the GDPR record of processing.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "retention": schema.StringAttribute{
                MarkdownDescription: "How long data is kept in the store, for the GDPR record of processing.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "labels": schema.SetAttribute{
                MarkdownDescription: "Labels to group and filter by. IDs of `oneuptime_label` resources.",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
                Validators: recordIDSetValidators("oneuptime_label"),
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
            "store_key": schema.StringAttribute{
                MarkdownDescription: "What identifies the store within the project: its provider, kind and native id. One row per store, however often it is discovered.",
                Computed: true,
            },
            "provider_value": schema.StringAttribute{
                MarkdownDescription: "Where the store lives: aws, azure, gcp, self-managed or oneuptime.",
                Computed: true,
            },
            "family": schema.StringAttribute{
                MarkdownDescription: "The store's family: object storage, database, warehouse, search, cache, file system, snapshot or telemetry.",
                Computed: true,
            },
            "discovered_from": schema.StringAttribute{
                MarkdownDescription: "How OneUptime found the store: cloud inventory, the Database product, telemetry, the service map, or added by hand.",
                Computed: true,
            },
            "cloud_account_id": schema.StringAttribute{
                MarkdownDescription: "The cloud account the store belongs to, for a cloud store. The ID of a `oneuptime_cloud_account`.",
                Computed: true,
            },
            "cloud_resource_id": schema.StringAttribute{
                MarkdownDescription: "The cloud resource the store was discovered from, if any. The ID of a `oneuptime_cloud_resource`.",
                Computed: true,
            },
            "cloud_connection_id": schema.StringAttribute{
                MarkdownDescription: "The cloud connection that reads the store's configuration, if any. The ID of a `oneuptime_cloud_connection`.",
                Computed: true,
            },
            "database_server_id": schema.StringAttribute{
                MarkdownDescription: "The Database product server the store is, for a self-managed database. The ID of a `oneuptime_database`.",
                Computed: true,
            },
            "service_id": schema.StringAttribute{
                MarkdownDescription: "The telemetry service whose logs, traces, exceptions, replays or AI conversations the store is. The ID of a `oneuptime_service`.",
                Computed: true,
            },
            "provider_resource_id": schema.StringAttribute{
                MarkdownDescription: "The store's id in its cloud: an ARN, an Azure resource id or a Google Cloud resource name.",
                Computed: true,
            },
            "cloud_account_number": schema.StringAttribute{
                MarkdownDescription: "The cloud's own id for the account: an AWS account id, an Azure subscription id or a Google Cloud project id.",
                Computed: true,
            },
            "region": schema.StringAttribute{
                MarkdownDescription: "The region the store lives in.",
                Computed: true,
            },
            "sensitivity": schema.StringAttribute{
                MarkdownDescription: "How sensitive the store is: Restricted, Confidential, Internal or Public, from the most sensitive data type found in it.",
                Computed: true,
            },
            "data_type_ids": schema.StringAttribute{
                MarkdownDescription: "The data types found in the store with at least Medium confidence. A JSON value: read it with `jsondecode()`.",
                CustomType: JSONSubsetType{},
                Computed: true,
            },
            "data_categories": schema.StringAttribute{
                MarkdownDescription: "The categories of those data types. A JSON value: read it with `jsondecode()`.",
                CustomType: JSONSubsetType{},
                Computed: true,
            },
            "subject_countries": schema.StringAttribute{
                MarkdownDescription: "Countries whose people's data the store holds, from country-bound data types (ISO 3166-1 alpha-2). A JSON value: read it with `jsondecode()`.",
                CustomType: JSONSubsetType{},
                Computed: true,
            },
            "exposure": schema.StringAttribute{
                MarkdownDescription: "Where the store sits relative to the internet and other organizations: Public, Other accounts, Internet reachable, Private or Unknown.",
                Computed: true,
            },
            "posture": schema.StringAttribute{
                MarkdownDescription: "What the store's configuration says about its safety: public access, encryption, backups, logging and the rest. A JSON value: read it with `jsondecode()`.",
                CustomType: JSONSubsetType{},
                Computed: true,
            },
            "cloud_tags": schema.StringAttribute{
                MarkdownDescription: "The store's tags or labels in its cloud. A JSON value: read it with `jsondecode()`.",
                CustomType: JSONSubsetType{},
                Computed: true,
            },
            "size_in_bytes": schema.NumberAttribute{
                MarkdownDescription: "How much data the store holds, from the metrics already collected.",
                Computed: true,
            },
            "object_count": schema.NumberAttribute{
                MarkdownDescription: "How many objects, rows or documents the store holds, when known.",
                Computed: true,
            },
            "estimated_records": schema.NumberAttribute{
                MarkdownDescription: "An estimate of how many records with sensitive data the store holds, from its samples.",
                Computed: true,
            },
            "last_activity_at": schema.StringAttribute{
                MarkdownDescription: "When the store was last read or written, when its metrics or traces say.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "configuration_read_at": schema.StringAttribute{
                MarkdownDescription: "When the store's configuration and access were last read.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "content_scanned_at": schema.StringAttribute{
                MarkdownDescription: "When the store's content was last sampled.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "content_scan_status": schema.StringAttribute{
                MarkdownDescription: "How the last content scan ended, or why it could not run.",
                Computed: true,
            },
            "next_content_scan_at": schema.StringAttribute{
                MarkdownDescription: "When the store's content is next sampled.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "finding_count": schema.NumberAttribute{
                MarkdownDescription: "How many data types were found at how many locations.",
                Computed: true,
            },
            "open_issue_count": schema.NumberAttribute{
                MarkdownDescription: "How many of the store's issues are open.",
                Computed: true,
            },
            "highest_issue_severity": schema.StringAttribute{
                MarkdownDescription: "The severity of the store's worst open issue.",
                Computed: true,
            },
            "risk_score": schema.NumberAttribute{
                MarkdownDescription: "The risk score of the store's worst open issue, from 0 to 100.",
                Computed: true,
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "The store is gone from its cloud or source; it is kept with its history.",
                Computed: true,
            },
            "first_seen_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime first found the store.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "last_seen_at": schema.StringAttribute{
                MarkdownDescription: "When OneUptime last found the store.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the user who created this data store. The ID of a `oneuptime_user` (see the data source).",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
        },
    }
}

// ModifyPlan keeps the state when nothing a configuration sets changes (see
// importsettle.go).
func (r *DataStoreResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
    // Destroying: nothing to plan.
    if req.Plan.Raw.IsNull() {
        return
    }

    keepStateWhenNothingChanges(ctx, req, resp)
}

func (r *DataStoreResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *DataStoreResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data DataStoreResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // What the configuration sets, and what Terraform planned (see keepPlannedValues).
    var config DataStoreResourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
    if resp.Diagnostics.HasError() {
        return
    }
    plan := data



    // Create API request body. Unset (null/unknown) optional fields are
    // omitted so server-side defaults apply instead of being overwritten
    // with zero values.
    dataStoreRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := dataStoreRequest["data"].(map[string]interface{})

    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        requestDataMap["name"] = data.Name.ValueString()
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        requestDataMap["description"] = data.Description.ValueString()
    }
    if !data.Kind.IsNull() && !data.Kind.IsUnknown() {
        requestDataMap["kind"] = data.Kind.ValueString()
    }
    if !data.Environment.IsNull() && !data.Environment.IsUnknown() {
        requestDataMap["environment"] = data.Environment.ValueString()
    }
    if !data.SensitivityOverride.IsNull() && !data.SensitivityOverride.IsUnknown() {
        requestDataMap["sensitivityOverride"] = data.SensitivityOverride.ValueString()
    }
    if !data.IsPublicByDesign.IsNull() && !data.IsPublicByDesign.IsUnknown() {
        requestDataMap["isPublicByDesign"] = data.IsPublicByDesign.ValueBool()
    }
    if !data.IsAiData.IsNull() && !data.IsAiData.IsUnknown() {
        requestDataMap["isAiData"] = data.IsAiData.ValueBool()
    }
    if !data.IsContentScanEnabled.IsNull() && !data.IsContentScanEnabled.IsUnknown() {
        requestDataMap["isContentScanEnabled"] = data.IsContentScanEnabled.ValueBool()
    }
    if !data.ScanProbeId.IsNull() && !data.ScanProbeId.IsUnknown() {
        requestDataMap["scanProbeId"] = data.ScanProbeId.ValueString()
    }
    if !data.ProbeConnectionVariable.IsNull() && !data.ProbeConnectionVariable.IsUnknown() {
        requestDataMap["probeConnectionVariable"] = data.ProbeConnectionVariable.ValueString()
    }
    if !data.ConnectionString.IsNull() && !data.ConnectionString.IsUnknown() {
        requestDataMap["connectionString"] = data.ConnectionString.ValueString()
    }
    if parsedScanExclusions := r.parseJSONField(data.ScanExclusions); parsedScanExclusions != nil {
        requestDataMap["scanExclusions"] = parsedScanExclusions
    }
    if !data.Purpose.IsNull() && !data.Purpose.IsUnknown() {
        requestDataMap["purpose"] = data.Purpose.ValueString()
    }
    if !data.Retention.IsNull() && !data.Retention.IsUnknown() {
        requestDataMap["retention"] = data.Retention.ValueString()
    }
    if !data.Labels.IsNull() && !data.Labels.IsUnknown() {
        requestDataMap["labels"] = r.convertTerraformSetToInterface(data.Labels)
    }

    // Make API call
    httpResp, err := r.client.Post(ctx, "/data-store", dataStoreRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create data_store, got error: %s", err))
        return
    }

    var dataStoreResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &dataStoreResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to create data_store: %s", err))
        return
    }

    // Extract the new resource id from the create response.
    createdId := ""
    if wrapper, ok := dataStoreResponse["data"].(map[string]interface{}); ok {
        if val, ok := wrapper["_id"].(string); ok {
            createdId = val
        }
    } else if val, ok := dataStoreResponse["_id"].(string); ok {
        createdId = val
    }
    if createdId == "" {
        resp.Diagnostics.AddError("OneUptime API Error", "Create response for data_store did not contain an id. This is a bug in the provider or the API; please report it.")
        return
    }
    data.Id = types.StringValue(createdId)

    /*
     * The server has committed the row. Persist what we know to state BEFORE
     * the read-back: if the read-back fails and we return without setting
     * state, Terraform never learns the resource exists and the created
     * data_store is orphaned server-side — never refreshed, never
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
        "description": true,
        "kind": true,
        "environment": true,
        "sensitivityOverride": true,
        "isPublicByDesign": true,
        "isAiData": true,
        "isContentScanEnabled": true,
        "scanProbeId": true,
        "probeConnectionVariable": true,
        "scanExclusions": true,
        "purpose": true,
        "retention": true,
        "labels": true,
        "createdAt": true,
        "updatedAt": true,
        "storeKey": true,
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
        "sensitivity": true,
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
        "createdByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/data-store/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        /*
         * State already owns the id, so the resource is tracked and the next
         * refresh reconciles the remaining attributes. Warn rather than
         * error: erroring here would strand a real resource.
         */
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created data_store but could not read it back; state is incomplete until the next refresh: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created data_store but could not parse the read-back response; state is incomplete until the next refresh: %s", err))
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
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := dataMap["kind"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Kind = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Kind = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Kind = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Kind = types.StringValue(string(jsonBytes))
            } else {
                data.Kind = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Kind = types.StringValue(string(jsonBytes))
            } else {
                data.Kind = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Kind = types.StringValue(string(jsonBytes))
        } else {
            data.Kind = types.StringNull()
        }
    } else if val, ok := dataMap["kind"].(string); ok {
        data.Kind = types.StringValue(val)
    } else {
        data.Kind = types.StringNull()
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
    } else if val, ok := dataMap["environment"].(string); ok {
        data.Environment = types.StringValue(val)
    } else {
        data.Environment = types.StringNull()
    }
    if obj, ok := dataMap["sensitivityOverride"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SensitivityOverride = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SensitivityOverride = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SensitivityOverride = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SensitivityOverride = types.StringValue(string(jsonBytes))
            } else {
                data.SensitivityOverride = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SensitivityOverride = types.StringValue(string(jsonBytes))
            } else {
                data.SensitivityOverride = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SensitivityOverride = types.StringValue(string(jsonBytes))
        } else {
            data.SensitivityOverride = types.StringNull()
        }
    } else if val, ok := dataMap["sensitivityOverride"].(string); ok {
        data.SensitivityOverride = types.StringValue(val)
    } else {
        data.SensitivityOverride = types.StringNull()
    }
    if val, ok := dataMap["isPublicByDesign"].(bool); ok {
        data.IsPublicByDesign = types.BoolValue(val)
    }
    if val, ok := dataMap["isAiData"].(bool); ok {
        data.IsAiData = types.BoolValue(val)
    }
    if val, ok := dataMap["isContentScanEnabled"].(bool); ok {
        data.IsContentScanEnabled = types.BoolValue(val)
    }
    if obj, ok := dataMap["scanProbeId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScanProbeId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ScanProbeId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ScanProbeId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ScanProbeId = types.StringValue(string(jsonBytes))
            } else {
                data.ScanProbeId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ScanProbeId = types.StringValue(string(jsonBytes))
            } else {
                data.ScanProbeId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ScanProbeId = types.StringValue(string(jsonBytes))
        } else {
            data.ScanProbeId = types.StringNull()
        }
    } else if val, ok := dataMap["scanProbeId"].(string); ok {
        data.ScanProbeId = types.StringValue(val)
    } else {
        data.ScanProbeId = types.StringNull()
    }
    if obj, ok := dataMap["probeConnectionVariable"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProbeConnectionVariable = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ProbeConnectionVariable = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ProbeConnectionVariable = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ProbeConnectionVariable = types.StringValue(string(jsonBytes))
            } else {
                data.ProbeConnectionVariable = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ProbeConnectionVariable = types.StringValue(string(jsonBytes))
            } else {
                data.ProbeConnectionVariable = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ProbeConnectionVariable = types.StringValue(string(jsonBytes))
        } else {
            data.ProbeConnectionVariable = types.StringNull()
        }
    } else if val, ok := dataMap["probeConnectionVariable"].(string); ok {
        data.ProbeConnectionVariable = types.StringValue(val)
    } else {
        data.ProbeConnectionVariable = types.StringNull()
    }
    if obj, ok := dataMap["scanExclusions"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScanExclusions = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ScanExclusions = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ScanExclusions = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ScanExclusions = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.ScanExclusions = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ScanExclusions = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.ScanExclusions = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ScanExclusions = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.ScanExclusions = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["scanExclusions"].(string); ok {
        data.ScanExclusions = NewJSONSubsetValue(val)
    } else {
        data.ScanExclusions = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["purpose"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Purpose = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Purpose = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Purpose = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Purpose = types.StringValue(string(jsonBytes))
            } else {
                data.Purpose = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Purpose = types.StringValue(string(jsonBytes))
            } else {
                data.Purpose = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Purpose = types.StringValue(string(jsonBytes))
        } else {
            data.Purpose = types.StringNull()
        }
    } else if val, ok := dataMap["purpose"].(string); ok {
        data.Purpose = types.StringValue(val)
    } else {
        data.Purpose = types.StringNull()
    }
    if obj, ok := dataMap["retention"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Retention = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Retention = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Retention = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Retention = types.StringValue(string(jsonBytes))
            } else {
                data.Retention = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Retention = types.StringValue(string(jsonBytes))
            } else {
                data.Retention = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Retention = types.StringValue(string(jsonBytes))
        } else {
            data.Retention = types.StringNull()
        }
    } else if val, ok := dataMap["retention"].(string); ok {
        data.Retention = types.StringValue(val)
    } else {
        data.Retention = types.StringNull()
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
    if obj, ok := dataMap["storeKey"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StoreKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StoreKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StoreKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.StoreKey = types.StringValue(string(jsonBytes))
            } else {
                data.StoreKey = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.StoreKey = types.StringValue(string(jsonBytes))
            } else {
                data.StoreKey = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StoreKey = types.StringValue(string(jsonBytes))
        } else {
            data.StoreKey = types.StringNull()
        }
    } else if val, ok := dataMap["storeKey"].(string); ok {
        data.StoreKey = types.StringValue(val)
    } else {
        data.StoreKey = types.StringNull()
    }
    if obj, ok := dataMap["provider"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProviderValue = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ProviderValue = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ProviderValue = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ProviderValue = types.StringValue(string(jsonBytes))
            } else {
                data.ProviderValue = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ProviderValue = types.StringValue(string(jsonBytes))
            } else {
                data.ProviderValue = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ProviderValue = types.StringValue(string(jsonBytes))
        } else {
            data.ProviderValue = types.StringNull()
        }
    } else if val, ok := dataMap["provider"].(string); ok {
        data.ProviderValue = types.StringValue(val)
    } else {
        data.ProviderValue = types.StringNull()
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
    } else if val, ok := dataMap["family"].(string); ok {
        data.Family = types.StringValue(val)
    } else {
        data.Family = types.StringNull()
    }
    if obj, ok := dataMap["discoveredFrom"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DiscoveredFrom = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DiscoveredFrom = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DiscoveredFrom = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DiscoveredFrom = types.StringValue(string(jsonBytes))
            } else {
                data.DiscoveredFrom = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DiscoveredFrom = types.StringValue(string(jsonBytes))
            } else {
                data.DiscoveredFrom = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DiscoveredFrom = types.StringValue(string(jsonBytes))
        } else {
            data.DiscoveredFrom = types.StringNull()
        }
    } else if val, ok := dataMap["discoveredFrom"].(string); ok {
        data.DiscoveredFrom = types.StringValue(val)
    } else {
        data.DiscoveredFrom = types.StringNull()
    }
    if obj, ok := dataMap["cloudAccountId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudAccountId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CloudAccountId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CloudAccountId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CloudAccountId = types.StringValue(string(jsonBytes))
            } else {
                data.CloudAccountId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CloudAccountId = types.StringValue(string(jsonBytes))
            } else {
                data.CloudAccountId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CloudAccountId = types.StringValue(string(jsonBytes))
        } else {
            data.CloudAccountId = types.StringNull()
        }
    } else if val, ok := dataMap["cloudAccountId"].(string); ok {
        data.CloudAccountId = types.StringValue(val)
    } else {
        data.CloudAccountId = types.StringNull()
    }
    if obj, ok := dataMap["cloudResourceId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudResourceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CloudResourceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CloudResourceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CloudResourceId = types.StringValue(string(jsonBytes))
            } else {
                data.CloudResourceId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CloudResourceId = types.StringValue(string(jsonBytes))
            } else {
                data.CloudResourceId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CloudResourceId = types.StringValue(string(jsonBytes))
        } else {
            data.CloudResourceId = types.StringNull()
        }
    } else if val, ok := dataMap["cloudResourceId"].(string); ok {
        data.CloudResourceId = types.StringValue(val)
    } else {
        data.CloudResourceId = types.StringNull()
    }
    if obj, ok := dataMap["cloudConnectionId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudConnectionId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CloudConnectionId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CloudConnectionId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CloudConnectionId = types.StringValue(string(jsonBytes))
            } else {
                data.CloudConnectionId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CloudConnectionId = types.StringValue(string(jsonBytes))
            } else {
                data.CloudConnectionId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CloudConnectionId = types.StringValue(string(jsonBytes))
        } else {
            data.CloudConnectionId = types.StringNull()
        }
    } else if val, ok := dataMap["cloudConnectionId"].(string); ok {
        data.CloudConnectionId = types.StringValue(val)
    } else {
        data.CloudConnectionId = types.StringNull()
    }
    if obj, ok := dataMap["databaseServerId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DatabaseServerId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DatabaseServerId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DatabaseServerId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DatabaseServerId = types.StringValue(string(jsonBytes))
            } else {
                data.DatabaseServerId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DatabaseServerId = types.StringValue(string(jsonBytes))
            } else {
                data.DatabaseServerId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DatabaseServerId = types.StringValue(string(jsonBytes))
        } else {
            data.DatabaseServerId = types.StringNull()
        }
    } else if val, ok := dataMap["databaseServerId"].(string); ok {
        data.DatabaseServerId = types.StringValue(val)
    } else {
        data.DatabaseServerId = types.StringNull()
    }
    if obj, ok := dataMap["serviceId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServiceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ServiceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ServiceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ServiceId = types.StringValue(string(jsonBytes))
            } else {
                data.ServiceId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ServiceId = types.StringValue(string(jsonBytes))
            } else {
                data.ServiceId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ServiceId = types.StringValue(string(jsonBytes))
        } else {
            data.ServiceId = types.StringNull()
        }
    } else if val, ok := dataMap["serviceId"].(string); ok {
        data.ServiceId = types.StringValue(val)
    } else {
        data.ServiceId = types.StringNull()
    }
    if obj, ok := dataMap["providerResourceId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProviderResourceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ProviderResourceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ProviderResourceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ProviderResourceId = types.StringValue(string(jsonBytes))
            } else {
                data.ProviderResourceId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ProviderResourceId = types.StringValue(string(jsonBytes))
            } else {
                data.ProviderResourceId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ProviderResourceId = types.StringValue(string(jsonBytes))
        } else {
            data.ProviderResourceId = types.StringNull()
        }
    } else if val, ok := dataMap["providerResourceId"].(string); ok {
        data.ProviderResourceId = types.StringValue(val)
    } else {
        data.ProviderResourceId = types.StringNull()
    }
    if obj, ok := dataMap["cloudAccountNumber"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudAccountNumber = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CloudAccountNumber = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CloudAccountNumber = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CloudAccountNumber = types.StringValue(string(jsonBytes))
            } else {
                data.CloudAccountNumber = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CloudAccountNumber = types.StringValue(string(jsonBytes))
            } else {
                data.CloudAccountNumber = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CloudAccountNumber = types.StringValue(string(jsonBytes))
        } else {
            data.CloudAccountNumber = types.StringNull()
        }
    } else if val, ok := dataMap["cloudAccountNumber"].(string); ok {
        data.CloudAccountNumber = types.StringValue(val)
    } else {
        data.CloudAccountNumber = types.StringNull()
    }
    if obj, ok := dataMap["region"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Region = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Region = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Region = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Region = types.StringValue(string(jsonBytes))
            } else {
                data.Region = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Region = types.StringValue(string(jsonBytes))
            } else {
                data.Region = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Region = types.StringValue(string(jsonBytes))
        } else {
            data.Region = types.StringNull()
        }
    } else if val, ok := dataMap["region"].(string); ok {
        data.Region = types.StringValue(val)
    } else {
        data.Region = types.StringNull()
    }
    if obj, ok := dataMap["sensitivity"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Sensitivity = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Sensitivity = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Sensitivity = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Sensitivity = types.StringValue(string(jsonBytes))
            } else {
                data.Sensitivity = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Sensitivity = types.StringValue(string(jsonBytes))
            } else {
                data.Sensitivity = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Sensitivity = types.StringValue(string(jsonBytes))
        } else {
            data.Sensitivity = types.StringNull()
        }
    } else if val, ok := dataMap["sensitivity"].(string); ok {
        data.Sensitivity = types.StringValue(val)
    } else {
        data.Sensitivity = types.StringNull()
    }
    if obj, ok := dataMap["dataTypeIds"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataTypeIds = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataTypeIds = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataTypeIds = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataTypeIds = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataTypeIds = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataTypeIds = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataTypeIds = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataTypeIds = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataTypeIds = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["dataTypeIds"].(string); ok {
        data.DataTypeIds = NewJSONSubsetValue(val)
    } else {
        data.DataTypeIds = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["dataCategories"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataCategories = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataCategories = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataCategories = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataCategories = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataCategories = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataCategories = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataCategories = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataCategories = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataCategories = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["dataCategories"].(string); ok {
        data.DataCategories = NewJSONSubsetValue(val)
    } else {
        data.DataCategories = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["subjectCountries"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubjectCountries = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubjectCountries = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubjectCountries = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SubjectCountries = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SubjectCountries = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SubjectCountries = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SubjectCountries = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubjectCountries = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.SubjectCountries = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["subjectCountries"].(string); ok {
        data.SubjectCountries = NewJSONSubsetValue(val)
    } else {
        data.SubjectCountries = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["exposure"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Exposure = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Exposure = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Exposure = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Exposure = types.StringValue(string(jsonBytes))
            } else {
                data.Exposure = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Exposure = types.StringValue(string(jsonBytes))
            } else {
                data.Exposure = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Exposure = types.StringValue(string(jsonBytes))
        } else {
            data.Exposure = types.StringNull()
        }
    } else if val, ok := dataMap["exposure"].(string); ok {
        data.Exposure = types.StringValue(val)
    } else {
        data.Exposure = types.StringNull()
    }
    if obj, ok := dataMap["posture"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Posture = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Posture = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Posture = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Posture = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.Posture = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Posture = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.Posture = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Posture = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.Posture = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["posture"].(string); ok {
        data.Posture = NewJSONSubsetValue(val)
    } else {
        data.Posture = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["cloudTags"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudTags = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CloudTags = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CloudTags = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CloudTags = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CloudTags = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CloudTags = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CloudTags = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CloudTags = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CloudTags = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["cloudTags"].(string); ok {
        data.CloudTags = NewJSONSubsetValue(val)
    } else {
        data.CloudTags = NewJSONSubsetNull()
    }
    data.SizeInBytes = numberFromAPI(data.SizeInBytes, dataMap["sizeInBytes"])
    data.ObjectCount = numberFromAPI(data.ObjectCount, dataMap["objectCount"])
    data.EstimatedRecords = numberFromAPI(data.EstimatedRecords, dataMap["estimatedRecords"])
    if obj, ok := dataMap["lastActivityAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastActivityAt = NewRFC3339Value(val)
        } else {
            data.LastActivityAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastActivityAt"].(string); ok && val != "" {
        data.LastActivityAt = NewRFC3339Value(val)
    } else {
        data.LastActivityAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["configurationReadAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ConfigurationReadAt = NewRFC3339Value(val)
        } else {
            data.ConfigurationReadAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["configurationReadAt"].(string); ok && val != "" {
        data.ConfigurationReadAt = NewRFC3339Value(val)
    } else {
        data.ConfigurationReadAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["contentScannedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ContentScannedAt = NewRFC3339Value(val)
        } else {
            data.ContentScannedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["contentScannedAt"].(string); ok && val != "" {
        data.ContentScannedAt = NewRFC3339Value(val)
    } else {
        data.ContentScannedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["contentScanStatus"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ContentScanStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ContentScanStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ContentScanStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ContentScanStatus = types.StringValue(string(jsonBytes))
            } else {
                data.ContentScanStatus = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ContentScanStatus = types.StringValue(string(jsonBytes))
            } else {
                data.ContentScanStatus = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ContentScanStatus = types.StringValue(string(jsonBytes))
        } else {
            data.ContentScanStatus = types.StringNull()
        }
    } else if val, ok := dataMap["contentScanStatus"].(string); ok {
        data.ContentScanStatus = types.StringValue(val)
    } else {
        data.ContentScanStatus = types.StringNull()
    }
    if obj, ok := dataMap["nextContentScanAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.NextContentScanAt = NewRFC3339Value(val)
        } else {
            data.NextContentScanAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["nextContentScanAt"].(string); ok && val != "" {
        data.NextContentScanAt = NewRFC3339Value(val)
    } else {
        data.NextContentScanAt = NewRFC3339Null()
    }
    data.FindingCount = numberFromAPI(data.FindingCount, dataMap["findingCount"])
    data.OpenIssueCount = numberFromAPI(data.OpenIssueCount, dataMap["openIssueCount"])
    if obj, ok := dataMap["highestIssueSeverity"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HighestIssueSeverity = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.HighestIssueSeverity = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.HighestIssueSeverity = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.HighestIssueSeverity = types.StringValue(string(jsonBytes))
            } else {
                data.HighestIssueSeverity = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.HighestIssueSeverity = types.StringValue(string(jsonBytes))
            } else {
                data.HighestIssueSeverity = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.HighestIssueSeverity = types.StringValue(string(jsonBytes))
        } else {
            data.HighestIssueSeverity = types.StringNull()
        }
    } else if val, ok := dataMap["highestIssueSeverity"].(string); ok {
        data.HighestIssueSeverity = types.StringValue(val)
    } else {
        data.HighestIssueSeverity = types.StringNull()
    }
    data.RiskScore = numberFromAPI(data.RiskScore, dataMap["riskScore"])
    if val, ok := dataMap["isArchived"].(bool); ok {
        data.IsArchived = types.BoolValue(val)
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


    // Nothing a create made was imported (see importsettle.go).
    resp.Diagnostics.Append(forgetImportedFromServer(ctx, resp.Private)...)
    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DataStoreResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data DataStoreResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "description": true,
        "kind": true,
        "environment": true,
        "sensitivityOverride": true,
        "isPublicByDesign": true,
        "isAiData": true,
        "isContentScanEnabled": true,
        "scanProbeId": true,
        "probeConnectionVariable": true,
        "scanExclusions": true,
        "purpose": true,
        "retention": true,
        "labels": true,
        "createdAt": true,
        "updatedAt": true,
        "storeKey": true,
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
        "sensitivity": true,
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
        "createdByUserId": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect(ctx, "/data-store/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read data_store, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var dataStoreResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &dataStoreResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse data_store response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := dataStoreResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = dataStoreResponse
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
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := dataMap["kind"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Kind = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Kind = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Kind = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Kind = types.StringValue(string(jsonBytes))
            } else {
                data.Kind = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Kind = types.StringValue(string(jsonBytes))
            } else {
                data.Kind = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Kind = types.StringValue(string(jsonBytes))
        } else {
            data.Kind = types.StringNull()
        }
    } else if val, ok := dataMap["kind"].(string); ok {
        data.Kind = types.StringValue(val)
    } else {
        data.Kind = types.StringNull()
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
    } else if val, ok := dataMap["environment"].(string); ok {
        data.Environment = types.StringValue(val)
    } else {
        data.Environment = types.StringNull()
    }
    if obj, ok := dataMap["sensitivityOverride"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SensitivityOverride = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SensitivityOverride = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SensitivityOverride = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SensitivityOverride = types.StringValue(string(jsonBytes))
            } else {
                data.SensitivityOverride = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SensitivityOverride = types.StringValue(string(jsonBytes))
            } else {
                data.SensitivityOverride = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SensitivityOverride = types.StringValue(string(jsonBytes))
        } else {
            data.SensitivityOverride = types.StringNull()
        }
    } else if val, ok := dataMap["sensitivityOverride"].(string); ok {
        data.SensitivityOverride = types.StringValue(val)
    } else {
        data.SensitivityOverride = types.StringNull()
    }
    if val, ok := dataMap["isPublicByDesign"].(bool); ok {
        data.IsPublicByDesign = types.BoolValue(val)
    }
    if val, ok := dataMap["isAiData"].(bool); ok {
        data.IsAiData = types.BoolValue(val)
    }
    if val, ok := dataMap["isContentScanEnabled"].(bool); ok {
        data.IsContentScanEnabled = types.BoolValue(val)
    }
    if obj, ok := dataMap["scanProbeId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScanProbeId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ScanProbeId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ScanProbeId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ScanProbeId = types.StringValue(string(jsonBytes))
            } else {
                data.ScanProbeId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ScanProbeId = types.StringValue(string(jsonBytes))
            } else {
                data.ScanProbeId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ScanProbeId = types.StringValue(string(jsonBytes))
        } else {
            data.ScanProbeId = types.StringNull()
        }
    } else if val, ok := dataMap["scanProbeId"].(string); ok {
        data.ScanProbeId = types.StringValue(val)
    } else {
        data.ScanProbeId = types.StringNull()
    }
    if obj, ok := dataMap["probeConnectionVariable"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProbeConnectionVariable = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ProbeConnectionVariable = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ProbeConnectionVariable = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ProbeConnectionVariable = types.StringValue(string(jsonBytes))
            } else {
                data.ProbeConnectionVariable = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ProbeConnectionVariable = types.StringValue(string(jsonBytes))
            } else {
                data.ProbeConnectionVariable = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ProbeConnectionVariable = types.StringValue(string(jsonBytes))
        } else {
            data.ProbeConnectionVariable = types.StringNull()
        }
    } else if val, ok := dataMap["probeConnectionVariable"].(string); ok {
        data.ProbeConnectionVariable = types.StringValue(val)
    } else {
        data.ProbeConnectionVariable = types.StringNull()
    }
    if obj, ok := dataMap["scanExclusions"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScanExclusions = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ScanExclusions = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ScanExclusions = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ScanExclusions = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.ScanExclusions = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ScanExclusions = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.ScanExclusions = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ScanExclusions = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.ScanExclusions = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["scanExclusions"].(string); ok {
        data.ScanExclusions = NewJSONSubsetValue(val)
    } else {
        data.ScanExclusions = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["purpose"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Purpose = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Purpose = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Purpose = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Purpose = types.StringValue(string(jsonBytes))
            } else {
                data.Purpose = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Purpose = types.StringValue(string(jsonBytes))
            } else {
                data.Purpose = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Purpose = types.StringValue(string(jsonBytes))
        } else {
            data.Purpose = types.StringNull()
        }
    } else if val, ok := dataMap["purpose"].(string); ok {
        data.Purpose = types.StringValue(val)
    } else {
        data.Purpose = types.StringNull()
    }
    if obj, ok := dataMap["retention"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Retention = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Retention = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Retention = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Retention = types.StringValue(string(jsonBytes))
            } else {
                data.Retention = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Retention = types.StringValue(string(jsonBytes))
            } else {
                data.Retention = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Retention = types.StringValue(string(jsonBytes))
        } else {
            data.Retention = types.StringNull()
        }
    } else if val, ok := dataMap["retention"].(string); ok {
        data.Retention = types.StringValue(val)
    } else {
        data.Retention = types.StringNull()
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
    if obj, ok := dataMap["storeKey"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StoreKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StoreKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StoreKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.StoreKey = types.StringValue(string(jsonBytes))
            } else {
                data.StoreKey = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.StoreKey = types.StringValue(string(jsonBytes))
            } else {
                data.StoreKey = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StoreKey = types.StringValue(string(jsonBytes))
        } else {
            data.StoreKey = types.StringNull()
        }
    } else if val, ok := dataMap["storeKey"].(string); ok {
        data.StoreKey = types.StringValue(val)
    } else {
        data.StoreKey = types.StringNull()
    }
    if obj, ok := dataMap["provider"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProviderValue = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ProviderValue = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ProviderValue = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ProviderValue = types.StringValue(string(jsonBytes))
            } else {
                data.ProviderValue = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ProviderValue = types.StringValue(string(jsonBytes))
            } else {
                data.ProviderValue = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ProviderValue = types.StringValue(string(jsonBytes))
        } else {
            data.ProviderValue = types.StringNull()
        }
    } else if val, ok := dataMap["provider"].(string); ok {
        data.ProviderValue = types.StringValue(val)
    } else {
        data.ProviderValue = types.StringNull()
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
    } else if val, ok := dataMap["family"].(string); ok {
        data.Family = types.StringValue(val)
    } else {
        data.Family = types.StringNull()
    }
    if obj, ok := dataMap["discoveredFrom"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DiscoveredFrom = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DiscoveredFrom = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DiscoveredFrom = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DiscoveredFrom = types.StringValue(string(jsonBytes))
            } else {
                data.DiscoveredFrom = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DiscoveredFrom = types.StringValue(string(jsonBytes))
            } else {
                data.DiscoveredFrom = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DiscoveredFrom = types.StringValue(string(jsonBytes))
        } else {
            data.DiscoveredFrom = types.StringNull()
        }
    } else if val, ok := dataMap["discoveredFrom"].(string); ok {
        data.DiscoveredFrom = types.StringValue(val)
    } else {
        data.DiscoveredFrom = types.StringNull()
    }
    if obj, ok := dataMap["cloudAccountId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudAccountId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CloudAccountId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CloudAccountId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CloudAccountId = types.StringValue(string(jsonBytes))
            } else {
                data.CloudAccountId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CloudAccountId = types.StringValue(string(jsonBytes))
            } else {
                data.CloudAccountId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CloudAccountId = types.StringValue(string(jsonBytes))
        } else {
            data.CloudAccountId = types.StringNull()
        }
    } else if val, ok := dataMap["cloudAccountId"].(string); ok {
        data.CloudAccountId = types.StringValue(val)
    } else {
        data.CloudAccountId = types.StringNull()
    }
    if obj, ok := dataMap["cloudResourceId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudResourceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CloudResourceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CloudResourceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CloudResourceId = types.StringValue(string(jsonBytes))
            } else {
                data.CloudResourceId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CloudResourceId = types.StringValue(string(jsonBytes))
            } else {
                data.CloudResourceId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CloudResourceId = types.StringValue(string(jsonBytes))
        } else {
            data.CloudResourceId = types.StringNull()
        }
    } else if val, ok := dataMap["cloudResourceId"].(string); ok {
        data.CloudResourceId = types.StringValue(val)
    } else {
        data.CloudResourceId = types.StringNull()
    }
    if obj, ok := dataMap["cloudConnectionId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudConnectionId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CloudConnectionId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CloudConnectionId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CloudConnectionId = types.StringValue(string(jsonBytes))
            } else {
                data.CloudConnectionId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CloudConnectionId = types.StringValue(string(jsonBytes))
            } else {
                data.CloudConnectionId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CloudConnectionId = types.StringValue(string(jsonBytes))
        } else {
            data.CloudConnectionId = types.StringNull()
        }
    } else if val, ok := dataMap["cloudConnectionId"].(string); ok {
        data.CloudConnectionId = types.StringValue(val)
    } else {
        data.CloudConnectionId = types.StringNull()
    }
    if obj, ok := dataMap["databaseServerId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DatabaseServerId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DatabaseServerId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DatabaseServerId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DatabaseServerId = types.StringValue(string(jsonBytes))
            } else {
                data.DatabaseServerId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DatabaseServerId = types.StringValue(string(jsonBytes))
            } else {
                data.DatabaseServerId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DatabaseServerId = types.StringValue(string(jsonBytes))
        } else {
            data.DatabaseServerId = types.StringNull()
        }
    } else if val, ok := dataMap["databaseServerId"].(string); ok {
        data.DatabaseServerId = types.StringValue(val)
    } else {
        data.DatabaseServerId = types.StringNull()
    }
    if obj, ok := dataMap["serviceId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServiceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ServiceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ServiceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ServiceId = types.StringValue(string(jsonBytes))
            } else {
                data.ServiceId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ServiceId = types.StringValue(string(jsonBytes))
            } else {
                data.ServiceId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ServiceId = types.StringValue(string(jsonBytes))
        } else {
            data.ServiceId = types.StringNull()
        }
    } else if val, ok := dataMap["serviceId"].(string); ok {
        data.ServiceId = types.StringValue(val)
    } else {
        data.ServiceId = types.StringNull()
    }
    if obj, ok := dataMap["providerResourceId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProviderResourceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ProviderResourceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ProviderResourceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ProviderResourceId = types.StringValue(string(jsonBytes))
            } else {
                data.ProviderResourceId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ProviderResourceId = types.StringValue(string(jsonBytes))
            } else {
                data.ProviderResourceId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ProviderResourceId = types.StringValue(string(jsonBytes))
        } else {
            data.ProviderResourceId = types.StringNull()
        }
    } else if val, ok := dataMap["providerResourceId"].(string); ok {
        data.ProviderResourceId = types.StringValue(val)
    } else {
        data.ProviderResourceId = types.StringNull()
    }
    if obj, ok := dataMap["cloudAccountNumber"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudAccountNumber = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CloudAccountNumber = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CloudAccountNumber = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CloudAccountNumber = types.StringValue(string(jsonBytes))
            } else {
                data.CloudAccountNumber = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CloudAccountNumber = types.StringValue(string(jsonBytes))
            } else {
                data.CloudAccountNumber = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CloudAccountNumber = types.StringValue(string(jsonBytes))
        } else {
            data.CloudAccountNumber = types.StringNull()
        }
    } else if val, ok := dataMap["cloudAccountNumber"].(string); ok {
        data.CloudAccountNumber = types.StringValue(val)
    } else {
        data.CloudAccountNumber = types.StringNull()
    }
    if obj, ok := dataMap["region"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Region = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Region = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Region = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Region = types.StringValue(string(jsonBytes))
            } else {
                data.Region = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Region = types.StringValue(string(jsonBytes))
            } else {
                data.Region = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Region = types.StringValue(string(jsonBytes))
        } else {
            data.Region = types.StringNull()
        }
    } else if val, ok := dataMap["region"].(string); ok {
        data.Region = types.StringValue(val)
    } else {
        data.Region = types.StringNull()
    }
    if obj, ok := dataMap["sensitivity"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Sensitivity = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Sensitivity = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Sensitivity = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Sensitivity = types.StringValue(string(jsonBytes))
            } else {
                data.Sensitivity = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Sensitivity = types.StringValue(string(jsonBytes))
            } else {
                data.Sensitivity = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Sensitivity = types.StringValue(string(jsonBytes))
        } else {
            data.Sensitivity = types.StringNull()
        }
    } else if val, ok := dataMap["sensitivity"].(string); ok {
        data.Sensitivity = types.StringValue(val)
    } else {
        data.Sensitivity = types.StringNull()
    }
    if obj, ok := dataMap["dataTypeIds"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataTypeIds = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataTypeIds = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataTypeIds = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataTypeIds = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataTypeIds = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataTypeIds = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataTypeIds = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataTypeIds = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataTypeIds = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["dataTypeIds"].(string); ok {
        data.DataTypeIds = NewJSONSubsetValue(val)
    } else {
        data.DataTypeIds = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["dataCategories"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataCategories = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataCategories = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataCategories = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataCategories = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataCategories = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataCategories = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataCategories = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataCategories = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataCategories = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["dataCategories"].(string); ok {
        data.DataCategories = NewJSONSubsetValue(val)
    } else {
        data.DataCategories = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["subjectCountries"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubjectCountries = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubjectCountries = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubjectCountries = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SubjectCountries = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SubjectCountries = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SubjectCountries = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SubjectCountries = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubjectCountries = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.SubjectCountries = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["subjectCountries"].(string); ok {
        data.SubjectCountries = NewJSONSubsetValue(val)
    } else {
        data.SubjectCountries = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["exposure"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Exposure = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Exposure = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Exposure = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Exposure = types.StringValue(string(jsonBytes))
            } else {
                data.Exposure = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Exposure = types.StringValue(string(jsonBytes))
            } else {
                data.Exposure = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Exposure = types.StringValue(string(jsonBytes))
        } else {
            data.Exposure = types.StringNull()
        }
    } else if val, ok := dataMap["exposure"].(string); ok {
        data.Exposure = types.StringValue(val)
    } else {
        data.Exposure = types.StringNull()
    }
    if obj, ok := dataMap["posture"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Posture = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Posture = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Posture = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Posture = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.Posture = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Posture = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.Posture = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Posture = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.Posture = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["posture"].(string); ok {
        data.Posture = NewJSONSubsetValue(val)
    } else {
        data.Posture = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["cloudTags"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudTags = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CloudTags = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CloudTags = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CloudTags = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CloudTags = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CloudTags = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CloudTags = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CloudTags = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CloudTags = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["cloudTags"].(string); ok {
        data.CloudTags = NewJSONSubsetValue(val)
    } else {
        data.CloudTags = NewJSONSubsetNull()
    }
    data.SizeInBytes = numberFromAPI(data.SizeInBytes, dataMap["sizeInBytes"])
    data.ObjectCount = numberFromAPI(data.ObjectCount, dataMap["objectCount"])
    data.EstimatedRecords = numberFromAPI(data.EstimatedRecords, dataMap["estimatedRecords"])
    if obj, ok := dataMap["lastActivityAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastActivityAt = NewRFC3339Value(val)
        } else {
            data.LastActivityAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastActivityAt"].(string); ok && val != "" {
        data.LastActivityAt = NewRFC3339Value(val)
    } else {
        data.LastActivityAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["configurationReadAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ConfigurationReadAt = NewRFC3339Value(val)
        } else {
            data.ConfigurationReadAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["configurationReadAt"].(string); ok && val != "" {
        data.ConfigurationReadAt = NewRFC3339Value(val)
    } else {
        data.ConfigurationReadAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["contentScannedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ContentScannedAt = NewRFC3339Value(val)
        } else {
            data.ContentScannedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["contentScannedAt"].(string); ok && val != "" {
        data.ContentScannedAt = NewRFC3339Value(val)
    } else {
        data.ContentScannedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["contentScanStatus"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ContentScanStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ContentScanStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ContentScanStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ContentScanStatus = types.StringValue(string(jsonBytes))
            } else {
                data.ContentScanStatus = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ContentScanStatus = types.StringValue(string(jsonBytes))
            } else {
                data.ContentScanStatus = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ContentScanStatus = types.StringValue(string(jsonBytes))
        } else {
            data.ContentScanStatus = types.StringNull()
        }
    } else if val, ok := dataMap["contentScanStatus"].(string); ok {
        data.ContentScanStatus = types.StringValue(val)
    } else {
        data.ContentScanStatus = types.StringNull()
    }
    if obj, ok := dataMap["nextContentScanAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.NextContentScanAt = NewRFC3339Value(val)
        } else {
            data.NextContentScanAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["nextContentScanAt"].(string); ok && val != "" {
        data.NextContentScanAt = NewRFC3339Value(val)
    } else {
        data.NextContentScanAt = NewRFC3339Null()
    }
    data.FindingCount = numberFromAPI(data.FindingCount, dataMap["findingCount"])
    data.OpenIssueCount = numberFromAPI(data.OpenIssueCount, dataMap["openIssueCount"])
    if obj, ok := dataMap["highestIssueSeverity"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HighestIssueSeverity = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.HighestIssueSeverity = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.HighestIssueSeverity = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.HighestIssueSeverity = types.StringValue(string(jsonBytes))
            } else {
                data.HighestIssueSeverity = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.HighestIssueSeverity = types.StringValue(string(jsonBytes))
            } else {
                data.HighestIssueSeverity = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.HighestIssueSeverity = types.StringValue(string(jsonBytes))
        } else {
            data.HighestIssueSeverity = types.StringNull()
        }
    } else if val, ok := dataMap["highestIssueSeverity"].(string); ok {
        data.HighestIssueSeverity = types.StringValue(val)
    } else {
        data.HighestIssueSeverity = types.StringNull()
    }
    data.RiskScore = numberFromAPI(data.RiskScore, dataMap["riskScore"])
    if val, ok := dataMap["isArchived"].(bool); ok {
        data.IsArchived = types.BoolValue(val)
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

func (r *DataStoreResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data DataStoreResourceModel
    var state DataStoreResourceModel

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

    // The attributes this update writes from the configuration (see importsettle.go).
    writtenFromConfig := map[string]bool{
        "scan_exclusions": !data.ScanExclusions.Equal(state.ScanExclusions),
    }

    // What the configuration sets, and what Terraform planned (see keepPlannedValues).
    var config DataStoreResourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
    if resp.Diagnostics.HasError() {
        return
    }
    plan := data

    // Create API request body
    dataStoreRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := dataStoreRequest["data"].(map[string]interface{})

    if !data.Name.IsUnknown() && !state.Name.IsUnknown() && !data.Name.Equal(state.Name) {
        requestDataMap["name"] = data.Name.ValueString()
    }
    if !data.Description.IsUnknown() && !state.Description.IsUnknown() && !data.Description.Equal(state.Description) {
        requestDataMap["description"] = data.Description.ValueString()
    }
    if !data.Environment.IsUnknown() && !state.Environment.IsUnknown() && !data.Environment.Equal(state.Environment) {
        requestDataMap["environment"] = data.Environment.ValueString()
    }
    if !data.SensitivityOverride.IsUnknown() && !state.SensitivityOverride.IsUnknown() && !data.SensitivityOverride.Equal(state.SensitivityOverride) {
        requestDataMap["sensitivityOverride"] = data.SensitivityOverride.ValueString()
    }
    if !data.IsPublicByDesign.IsUnknown() && !state.IsPublicByDesign.IsUnknown() && !data.IsPublicByDesign.Equal(state.IsPublicByDesign) {
        requestDataMap["isPublicByDesign"] = data.IsPublicByDesign.ValueBool()
    }
    if !data.IsAiData.IsUnknown() && !state.IsAiData.IsUnknown() && !data.IsAiData.Equal(state.IsAiData) {
        requestDataMap["isAiData"] = data.IsAiData.ValueBool()
    }
    if !data.IsContentScanEnabled.IsUnknown() && !state.IsContentScanEnabled.IsUnknown() && !data.IsContentScanEnabled.Equal(state.IsContentScanEnabled) {
        requestDataMap["isContentScanEnabled"] = data.IsContentScanEnabled.ValueBool()
    }
    if !data.ScanProbeId.IsUnknown() && !state.ScanProbeId.IsUnknown() && !data.ScanProbeId.Equal(state.ScanProbeId) {
        requestDataMap["scanProbeId"] = data.ScanProbeId.ValueString()
    }
    if !data.ProbeConnectionVariable.IsUnknown() && !state.ProbeConnectionVariable.IsUnknown() && !data.ProbeConnectionVariable.Equal(state.ProbeConnectionVariable) {
        requestDataMap["probeConnectionVariable"] = data.ProbeConnectionVariable.ValueString()
    }
    if !data.ConnectionString.IsUnknown() && !state.ConnectionString.IsUnknown() && !data.ConnectionString.Equal(state.ConnectionString) {
        requestDataMap["connectionString"] = data.ConnectionString.ValueString()
    }
    if !data.ScanExclusions.IsUnknown() && !state.ScanExclusions.IsUnknown() && !data.ScanExclusions.Equal(state.ScanExclusions) {
        var scanexclusionsData interface{}
        if err := json.Unmarshal([]byte(data.ScanExclusions.ValueString()), &scanexclusionsData); err == nil {
            requestDataMap["scanExclusions"] = scanexclusionsData
        } else {
            requestDataMap["scanExclusions"] = data.ScanExclusions.ValueString()
        }
    }
    if !data.Purpose.IsUnknown() && !state.Purpose.IsUnknown() && !data.Purpose.Equal(state.Purpose) {
        requestDataMap["purpose"] = data.Purpose.ValueString()
    }
    if !data.Retention.IsUnknown() && !state.Retention.IsUnknown() && !data.Retention.Equal(state.Retention) {
        requestDataMap["retention"] = data.Retention.ValueString()
    }
    if !data.Labels.IsUnknown() && !state.Labels.IsUnknown() && !data.Labels.Equal(state.Labels) {
        requestDataMap["labels"] = r.convertTerraformSetToInterface(data.Labels)
    }

    // Only call the API when there are changed fields to send. An empty
    // update body is rejected by the API; state is still refreshed below so
    // this method never writes unverified plan values into state.
    if len(dataStoreRequest["data"].(map[string]interface{})) > 0 {
        httpResp, err := r.client.Put(ctx, "/data-store/" + data.Id.ValueString() + "", dataStoreRequest)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update data_store, got error: %s", err))
            return
        }

        // Parse the update response
        var dataStoreResponse map[string]interface{}
        err = r.client.ParseResponse(httpResp, &dataStoreResponse)
        if err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to update data_store: %s", err))
            return
        }
        _ = dataStoreResponse
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "name": true,
        "description": true,
        "kind": true,
        "environment": true,
        "sensitivityOverride": true,
        "isPublicByDesign": true,
        "isAiData": true,
        "isContentScanEnabled": true,
        "scanProbeId": true,
        "probeConnectionVariable": true,
        "scanExclusions": true,
        "purpose": true,
        "retention": true,
        "labels": true,
        "createdAt": true,
        "updatedAt": true,
        "storeKey": true,
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
        "sensitivity": true,
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
        "createdByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/data-store/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read data_store after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read data_store after update: %s", err))
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
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["description"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Description = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Description = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Description = types.StringValue(string(jsonBytes))
            } else {
                data.Description = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Description = types.StringValue(string(jsonBytes))
        } else {
            data.Description = types.StringNull()
        }
    } else if val, ok := dataMap["description"].(string); ok {
        data.Description = types.StringValue(val)
    } else {
        data.Description = types.StringNull()
    }
    if obj, ok := dataMap["kind"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Kind = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Kind = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Kind = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Kind = types.StringValue(string(jsonBytes))
            } else {
                data.Kind = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Kind = types.StringValue(string(jsonBytes))
            } else {
                data.Kind = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Kind = types.StringValue(string(jsonBytes))
        } else {
            data.Kind = types.StringNull()
        }
    } else if val, ok := dataMap["kind"].(string); ok {
        data.Kind = types.StringValue(val)
    } else {
        data.Kind = types.StringNull()
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
    } else if val, ok := dataMap["environment"].(string); ok {
        data.Environment = types.StringValue(val)
    } else {
        data.Environment = types.StringNull()
    }
    if obj, ok := dataMap["sensitivityOverride"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SensitivityOverride = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SensitivityOverride = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SensitivityOverride = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SensitivityOverride = types.StringValue(string(jsonBytes))
            } else {
                data.SensitivityOverride = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SensitivityOverride = types.StringValue(string(jsonBytes))
            } else {
                data.SensitivityOverride = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SensitivityOverride = types.StringValue(string(jsonBytes))
        } else {
            data.SensitivityOverride = types.StringNull()
        }
    } else if val, ok := dataMap["sensitivityOverride"].(string); ok {
        data.SensitivityOverride = types.StringValue(val)
    } else {
        data.SensitivityOverride = types.StringNull()
    }
    if val, ok := dataMap["isPublicByDesign"].(bool); ok {
        data.IsPublicByDesign = types.BoolValue(val)
    }
    if val, ok := dataMap["isAiData"].(bool); ok {
        data.IsAiData = types.BoolValue(val)
    }
    if val, ok := dataMap["isContentScanEnabled"].(bool); ok {
        data.IsContentScanEnabled = types.BoolValue(val)
    }
    if obj, ok := dataMap["scanProbeId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScanProbeId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ScanProbeId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ScanProbeId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ScanProbeId = types.StringValue(string(jsonBytes))
            } else {
                data.ScanProbeId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ScanProbeId = types.StringValue(string(jsonBytes))
            } else {
                data.ScanProbeId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ScanProbeId = types.StringValue(string(jsonBytes))
        } else {
            data.ScanProbeId = types.StringNull()
        }
    } else if val, ok := dataMap["scanProbeId"].(string); ok {
        data.ScanProbeId = types.StringValue(val)
    } else {
        data.ScanProbeId = types.StringNull()
    }
    if obj, ok := dataMap["probeConnectionVariable"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProbeConnectionVariable = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ProbeConnectionVariable = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ProbeConnectionVariable = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ProbeConnectionVariable = types.StringValue(string(jsonBytes))
            } else {
                data.ProbeConnectionVariable = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ProbeConnectionVariable = types.StringValue(string(jsonBytes))
            } else {
                data.ProbeConnectionVariable = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ProbeConnectionVariable = types.StringValue(string(jsonBytes))
        } else {
            data.ProbeConnectionVariable = types.StringNull()
        }
    } else if val, ok := dataMap["probeConnectionVariable"].(string); ok {
        data.ProbeConnectionVariable = types.StringValue(val)
    } else {
        data.ProbeConnectionVariable = types.StringNull()
    }
    if obj, ok := dataMap["scanExclusions"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScanExclusions = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ScanExclusions = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ScanExclusions = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ScanExclusions = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.ScanExclusions = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ScanExclusions = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.ScanExclusions = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ScanExclusions = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.ScanExclusions = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["scanExclusions"].(string); ok {
        data.ScanExclusions = NewJSONSubsetValue(val)
    } else {
        data.ScanExclusions = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["purpose"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Purpose = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Purpose = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Purpose = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Purpose = types.StringValue(string(jsonBytes))
            } else {
                data.Purpose = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Purpose = types.StringValue(string(jsonBytes))
            } else {
                data.Purpose = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Purpose = types.StringValue(string(jsonBytes))
        } else {
            data.Purpose = types.StringNull()
        }
    } else if val, ok := dataMap["purpose"].(string); ok {
        data.Purpose = types.StringValue(val)
    } else {
        data.Purpose = types.StringNull()
    }
    if obj, ok := dataMap["retention"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Retention = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Retention = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Retention = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Retention = types.StringValue(string(jsonBytes))
            } else {
                data.Retention = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Retention = types.StringValue(string(jsonBytes))
            } else {
                data.Retention = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Retention = types.StringValue(string(jsonBytes))
        } else {
            data.Retention = types.StringNull()
        }
    } else if val, ok := dataMap["retention"].(string); ok {
        data.Retention = types.StringValue(val)
    } else {
        data.Retention = types.StringNull()
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
    if obj, ok := dataMap["storeKey"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StoreKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StoreKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StoreKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.StoreKey = types.StringValue(string(jsonBytes))
            } else {
                data.StoreKey = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.StoreKey = types.StringValue(string(jsonBytes))
            } else {
                data.StoreKey = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StoreKey = types.StringValue(string(jsonBytes))
        } else {
            data.StoreKey = types.StringNull()
        }
    } else if val, ok := dataMap["storeKey"].(string); ok {
        data.StoreKey = types.StringValue(val)
    } else {
        data.StoreKey = types.StringNull()
    }
    if obj, ok := dataMap["provider"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProviderValue = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ProviderValue = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ProviderValue = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ProviderValue = types.StringValue(string(jsonBytes))
            } else {
                data.ProviderValue = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ProviderValue = types.StringValue(string(jsonBytes))
            } else {
                data.ProviderValue = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ProviderValue = types.StringValue(string(jsonBytes))
        } else {
            data.ProviderValue = types.StringNull()
        }
    } else if val, ok := dataMap["provider"].(string); ok {
        data.ProviderValue = types.StringValue(val)
    } else {
        data.ProviderValue = types.StringNull()
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
    } else if val, ok := dataMap["family"].(string); ok {
        data.Family = types.StringValue(val)
    } else {
        data.Family = types.StringNull()
    }
    if obj, ok := dataMap["discoveredFrom"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DiscoveredFrom = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DiscoveredFrom = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DiscoveredFrom = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DiscoveredFrom = types.StringValue(string(jsonBytes))
            } else {
                data.DiscoveredFrom = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DiscoveredFrom = types.StringValue(string(jsonBytes))
            } else {
                data.DiscoveredFrom = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DiscoveredFrom = types.StringValue(string(jsonBytes))
        } else {
            data.DiscoveredFrom = types.StringNull()
        }
    } else if val, ok := dataMap["discoveredFrom"].(string); ok {
        data.DiscoveredFrom = types.StringValue(val)
    } else {
        data.DiscoveredFrom = types.StringNull()
    }
    if obj, ok := dataMap["cloudAccountId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudAccountId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CloudAccountId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CloudAccountId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CloudAccountId = types.StringValue(string(jsonBytes))
            } else {
                data.CloudAccountId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CloudAccountId = types.StringValue(string(jsonBytes))
            } else {
                data.CloudAccountId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CloudAccountId = types.StringValue(string(jsonBytes))
        } else {
            data.CloudAccountId = types.StringNull()
        }
    } else if val, ok := dataMap["cloudAccountId"].(string); ok {
        data.CloudAccountId = types.StringValue(val)
    } else {
        data.CloudAccountId = types.StringNull()
    }
    if obj, ok := dataMap["cloudResourceId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudResourceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CloudResourceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CloudResourceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CloudResourceId = types.StringValue(string(jsonBytes))
            } else {
                data.CloudResourceId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CloudResourceId = types.StringValue(string(jsonBytes))
            } else {
                data.CloudResourceId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CloudResourceId = types.StringValue(string(jsonBytes))
        } else {
            data.CloudResourceId = types.StringNull()
        }
    } else if val, ok := dataMap["cloudResourceId"].(string); ok {
        data.CloudResourceId = types.StringValue(val)
    } else {
        data.CloudResourceId = types.StringNull()
    }
    if obj, ok := dataMap["cloudConnectionId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudConnectionId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CloudConnectionId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CloudConnectionId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CloudConnectionId = types.StringValue(string(jsonBytes))
            } else {
                data.CloudConnectionId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CloudConnectionId = types.StringValue(string(jsonBytes))
            } else {
                data.CloudConnectionId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CloudConnectionId = types.StringValue(string(jsonBytes))
        } else {
            data.CloudConnectionId = types.StringNull()
        }
    } else if val, ok := dataMap["cloudConnectionId"].(string); ok {
        data.CloudConnectionId = types.StringValue(val)
    } else {
        data.CloudConnectionId = types.StringNull()
    }
    if obj, ok := dataMap["databaseServerId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DatabaseServerId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DatabaseServerId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DatabaseServerId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DatabaseServerId = types.StringValue(string(jsonBytes))
            } else {
                data.DatabaseServerId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DatabaseServerId = types.StringValue(string(jsonBytes))
            } else {
                data.DatabaseServerId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DatabaseServerId = types.StringValue(string(jsonBytes))
        } else {
            data.DatabaseServerId = types.StringNull()
        }
    } else if val, ok := dataMap["databaseServerId"].(string); ok {
        data.DatabaseServerId = types.StringValue(val)
    } else {
        data.DatabaseServerId = types.StringNull()
    }
    if obj, ok := dataMap["serviceId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ServiceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ServiceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ServiceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ServiceId = types.StringValue(string(jsonBytes))
            } else {
                data.ServiceId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ServiceId = types.StringValue(string(jsonBytes))
            } else {
                data.ServiceId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ServiceId = types.StringValue(string(jsonBytes))
        } else {
            data.ServiceId = types.StringNull()
        }
    } else if val, ok := dataMap["serviceId"].(string); ok {
        data.ServiceId = types.StringValue(val)
    } else {
        data.ServiceId = types.StringNull()
    }
    if obj, ok := dataMap["providerResourceId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProviderResourceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ProviderResourceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ProviderResourceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ProviderResourceId = types.StringValue(string(jsonBytes))
            } else {
                data.ProviderResourceId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ProviderResourceId = types.StringValue(string(jsonBytes))
            } else {
                data.ProviderResourceId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ProviderResourceId = types.StringValue(string(jsonBytes))
        } else {
            data.ProviderResourceId = types.StringNull()
        }
    } else if val, ok := dataMap["providerResourceId"].(string); ok {
        data.ProviderResourceId = types.StringValue(val)
    } else {
        data.ProviderResourceId = types.StringNull()
    }
    if obj, ok := dataMap["cloudAccountNumber"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudAccountNumber = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CloudAccountNumber = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CloudAccountNumber = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CloudAccountNumber = types.StringValue(string(jsonBytes))
            } else {
                data.CloudAccountNumber = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CloudAccountNumber = types.StringValue(string(jsonBytes))
            } else {
                data.CloudAccountNumber = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CloudAccountNumber = types.StringValue(string(jsonBytes))
        } else {
            data.CloudAccountNumber = types.StringNull()
        }
    } else if val, ok := dataMap["cloudAccountNumber"].(string); ok {
        data.CloudAccountNumber = types.StringValue(val)
    } else {
        data.CloudAccountNumber = types.StringNull()
    }
    if obj, ok := dataMap["region"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Region = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Region = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Region = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Region = types.StringValue(string(jsonBytes))
            } else {
                data.Region = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Region = types.StringValue(string(jsonBytes))
            } else {
                data.Region = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Region = types.StringValue(string(jsonBytes))
        } else {
            data.Region = types.StringNull()
        }
    } else if val, ok := dataMap["region"].(string); ok {
        data.Region = types.StringValue(val)
    } else {
        data.Region = types.StringNull()
    }
    if obj, ok := dataMap["sensitivity"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Sensitivity = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Sensitivity = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Sensitivity = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Sensitivity = types.StringValue(string(jsonBytes))
            } else {
                data.Sensitivity = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Sensitivity = types.StringValue(string(jsonBytes))
            } else {
                data.Sensitivity = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Sensitivity = types.StringValue(string(jsonBytes))
        } else {
            data.Sensitivity = types.StringNull()
        }
    } else if val, ok := dataMap["sensitivity"].(string); ok {
        data.Sensitivity = types.StringValue(val)
    } else {
        data.Sensitivity = types.StringNull()
    }
    if obj, ok := dataMap["dataTypeIds"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataTypeIds = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataTypeIds = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataTypeIds = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataTypeIds = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataTypeIds = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataTypeIds = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataTypeIds = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataTypeIds = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataTypeIds = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["dataTypeIds"].(string); ok {
        data.DataTypeIds = NewJSONSubsetValue(val)
    } else {
        data.DataTypeIds = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["dataCategories"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataCategories = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataCategories = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataCategories = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataCategories = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataCategories = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataCategories = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DataCategories = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataCategories = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DataCategories = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["dataCategories"].(string); ok {
        data.DataCategories = NewJSONSubsetValue(val)
    } else {
        data.DataCategories = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["subjectCountries"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubjectCountries = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubjectCountries = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubjectCountries = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SubjectCountries = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SubjectCountries = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SubjectCountries = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.SubjectCountries = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubjectCountries = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.SubjectCountries = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["subjectCountries"].(string); ok {
        data.SubjectCountries = NewJSONSubsetValue(val)
    } else {
        data.SubjectCountries = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["exposure"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Exposure = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Exposure = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Exposure = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Exposure = types.StringValue(string(jsonBytes))
            } else {
                data.Exposure = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Exposure = types.StringValue(string(jsonBytes))
            } else {
                data.Exposure = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Exposure = types.StringValue(string(jsonBytes))
        } else {
            data.Exposure = types.StringNull()
        }
    } else if val, ok := dataMap["exposure"].(string); ok {
        data.Exposure = types.StringValue(val)
    } else {
        data.Exposure = types.StringNull()
    }
    if obj, ok := dataMap["posture"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Posture = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Posture = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Posture = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Posture = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.Posture = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Posture = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.Posture = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Posture = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.Posture = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["posture"].(string); ok {
        data.Posture = NewJSONSubsetValue(val)
    } else {
        data.Posture = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["cloudTags"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CloudTags = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CloudTags = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CloudTags = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CloudTags = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CloudTags = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CloudTags = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CloudTags = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CloudTags = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CloudTags = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["cloudTags"].(string); ok {
        data.CloudTags = NewJSONSubsetValue(val)
    } else {
        data.CloudTags = NewJSONSubsetNull()
    }
    data.SizeInBytes = numberFromAPI(data.SizeInBytes, dataMap["sizeInBytes"])
    data.ObjectCount = numberFromAPI(data.ObjectCount, dataMap["objectCount"])
    data.EstimatedRecords = numberFromAPI(data.EstimatedRecords, dataMap["estimatedRecords"])
    if obj, ok := dataMap["lastActivityAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.LastActivityAt = NewRFC3339Value(val)
        } else {
            data.LastActivityAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["lastActivityAt"].(string); ok && val != "" {
        data.LastActivityAt = NewRFC3339Value(val)
    } else {
        data.LastActivityAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["configurationReadAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ConfigurationReadAt = NewRFC3339Value(val)
        } else {
            data.ConfigurationReadAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["configurationReadAt"].(string); ok && val != "" {
        data.ConfigurationReadAt = NewRFC3339Value(val)
    } else {
        data.ConfigurationReadAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["contentScannedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.ContentScannedAt = NewRFC3339Value(val)
        } else {
            data.ContentScannedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["contentScannedAt"].(string); ok && val != "" {
        data.ContentScannedAt = NewRFC3339Value(val)
    } else {
        data.ContentScannedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["contentScanStatus"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ContentScanStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ContentScanStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ContentScanStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ContentScanStatus = types.StringValue(string(jsonBytes))
            } else {
                data.ContentScanStatus = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ContentScanStatus = types.StringValue(string(jsonBytes))
            } else {
                data.ContentScanStatus = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ContentScanStatus = types.StringValue(string(jsonBytes))
        } else {
            data.ContentScanStatus = types.StringNull()
        }
    } else if val, ok := dataMap["contentScanStatus"].(string); ok {
        data.ContentScanStatus = types.StringValue(val)
    } else {
        data.ContentScanStatus = types.StringNull()
    }
    if obj, ok := dataMap["nextContentScanAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.NextContentScanAt = NewRFC3339Value(val)
        } else {
            data.NextContentScanAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["nextContentScanAt"].(string); ok && val != "" {
        data.NextContentScanAt = NewRFC3339Value(val)
    } else {
        data.NextContentScanAt = NewRFC3339Null()
    }
    data.FindingCount = numberFromAPI(data.FindingCount, dataMap["findingCount"])
    data.OpenIssueCount = numberFromAPI(data.OpenIssueCount, dataMap["openIssueCount"])
    if obj, ok := dataMap["highestIssueSeverity"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.HighestIssueSeverity = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.HighestIssueSeverity = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.HighestIssueSeverity = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.HighestIssueSeverity = types.StringValue(string(jsonBytes))
            } else {
                data.HighestIssueSeverity = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.HighestIssueSeverity = types.StringValue(string(jsonBytes))
            } else {
                data.HighestIssueSeverity = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.HighestIssueSeverity = types.StringValue(string(jsonBytes))
        } else {
            data.HighestIssueSeverity = types.StringNull()
        }
    } else if val, ok := dataMap["highestIssueSeverity"].(string); ok {
        data.HighestIssueSeverity = types.StringValue(val)
    } else {
        data.HighestIssueSeverity = types.StringNull()
    }
    data.RiskScore = numberFromAPI(data.RiskScore, dataMap["riskScore"])
    if val, ok := dataMap["isArchived"].(bool); ok {
        data.IsArchived = types.BoolValue(val)
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

    // What this update wrote is what the configuration says now (see importsettle.go).
    resp.Diagnostics.Append(keepImportedFromServer(ctx, resp.Private, writtenFromConfig)...)

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DataStoreResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data DataStoreResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    httpResp, err := r.client.Delete(ctx, "/data-store/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete data_store, got error: %s", err))
        return
    }

    // A failed delete must keep the resource in state — silently dropping it
    // orphans real infrastructure. 404 means it is already gone.
    if httpResp.StatusCode >= 400 && httpResp.StatusCode != http.StatusNotFound {
        err = r.client.ParseResponse(httpResp, nil)
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to delete data_store: %s", err))
        return
    }
    if httpResp.Body != nil {
        httpResp.Body.Close()
    }
}


func (r *DataStoreResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
    // These hold the record as OneUptime keeps it, keys it fills in itself
    // included, until an update writes them (see importsettle.go).
    resp.Diagnostics.Append(markImportedFromServer(ctx, resp.Private, []string{"scan_exclusions"})...)
}


// keepPlannedValues puts back, after a create or an update, the planned value
// of each optional attribute the configuration leaves out. The server keeps
// some of these up to date on its own (when it last checked a heartbeat, the
// status a probe last reported...), so the value read back after the write can
// already differ from the plan, and Terraform would fail the apply with
// "Provider produced inconsistent result after apply". The next refresh reads
// the server's value, which is never a diff for an attribute nobody configured.
func (r *DataStoreResource) keepPlannedValues(data *DataStoreResourceModel, plan *DataStoreResourceModel, config *DataStoreResourceModel) {
    if config.Description.IsNull() && !plan.Description.IsUnknown() {
        data.Description = plan.Description
    }
    if config.Environment.IsNull() && !plan.Environment.IsUnknown() {
        data.Environment = plan.Environment
    }
    if config.SensitivityOverride.IsNull() && !plan.SensitivityOverride.IsUnknown() {
        data.SensitivityOverride = plan.SensitivityOverride
    }
    if config.IsPublicByDesign.IsNull() && !plan.IsPublicByDesign.IsUnknown() {
        data.IsPublicByDesign = plan.IsPublicByDesign
    }
    if config.IsAiData.IsNull() && !plan.IsAiData.IsUnknown() {
        data.IsAiData = plan.IsAiData
    }
    if config.IsContentScanEnabled.IsNull() && !plan.IsContentScanEnabled.IsUnknown() {
        data.IsContentScanEnabled = plan.IsContentScanEnabled
    }
    if config.ScanProbeId.IsNull() && !plan.ScanProbeId.IsUnknown() {
        data.ScanProbeId = plan.ScanProbeId
    }
    if config.ProbeConnectionVariable.IsNull() && !plan.ProbeConnectionVariable.IsUnknown() {
        data.ProbeConnectionVariable = plan.ProbeConnectionVariable
    }
    if config.ScanExclusions.IsNull() && !plan.ScanExclusions.IsUnknown() {
        data.ScanExclusions = plan.ScanExclusions
    }
    if config.Purpose.IsNull() && !plan.Purpose.IsUnknown() {
        data.Purpose = plan.Purpose
    }
    if config.Retention.IsNull() && !plan.Retention.IsUnknown() {
        data.Retention = plan.Retention
    }
    if config.Labels.IsNull() && !plan.Labels.IsUnknown() {
        data.Labels = plan.Labels
    }
}

// Helper method to convert Terraform map to Go interface{}
func (r *DataStoreResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *DataStoreResource) convertTerraformListToInterface(terraformList types.List) interface{} {
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
func (r *DataStoreResource) convertTerraformSetToInterface(terraformSet types.Set) interface{} {
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
func (r *DataStoreResource) parseJSONField(terraformString basetypes.StringValuable) interface{} {
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
func (r *DataStoreResource) normalizeURLWrappers(value interface{}) interface{} {
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

func (r *DataStoreResource) normalizeURLString(value string) string {
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
func (r *DataStoreResource) bigFloatToFloat64(bf *big.Float) interface{} {
    if bf == nil {
        return nil
    }
    f, _ := bf.Float64()
    return f
}

// Helper method to check if a type string is a valid OneUptime ObjectType.
// The registry itself lives in objecttypes.go, shared across the package.
func (r *DataStoreResource) isValidOneUptimeObjectType(typeStr string) bool {
    return validOneUptimeObjectTypes[typeStr]
}
