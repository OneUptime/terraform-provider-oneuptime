package shard06

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
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
var _ datasource.DataSource = &TelemetryIngestionKeyDataSource{}

func NewTelemetryIngestionKeyDataSource() datasource.DataSource {
    return &TelemetryIngestionKeyDataSource{}
}

// TelemetryIngestionKeyDataSource defines the data source implementation.
type TelemetryIngestionKeyDataSource struct {
    client *Client
}

// TelemetryIngestionKeyDataSourceModel describes the data source data model.
type TelemetryIngestionKeyDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    CreatedByActorType types.String `tfsdk:"created_by_actor_type"`
    CreatedByActorId types.String `tfsdk:"created_by_actor_id"`
    LastModifiedAt types.String `tfsdk:"last_modified_at"`
    LastModifiedByActorType types.String `tfsdk:"last_modified_by_actor_type"`
    LastModifiedByActorId types.String `tfsdk:"last_modified_by_actor_id"`
    LastModifiedByUserId types.String `tfsdk:"last_modified_by_user_id"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    SecretKey types.String `tfsdk:"secret_key"`
    KeyType types.String `tfsdk:"key_type"`
    AllowedOrigins types.String `tfsdk:"allowed_origins"`
    PinnedServiceName types.String `tfsdk:"pinned_service_name"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    ExpiresAt types.String `tfsdk:"expires_at"`
    LastUsedAt types.String `tfsdk:"last_used_at"`
    RequestsPerMinuteLimit types.Number `tfsdk:"requests_per_minute_limit"`
    LogsDailyIngestionBudgetInGb types.Number `tfsdk:"logs_daily_ingestion_budget_in_gb"`
    TracesDailyIngestionBudgetInGb types.Number `tfsdk:"traces_daily_ingestion_budget_in_gb"`
    MetricsDailyIngestionBudgetInGb types.Number `tfsdk:"metrics_daily_ingestion_budget_in_gb"`
    ProfilesDailyIngestionBudgetInGb types.Number `tfsdk:"profiles_daily_ingestion_budget_in_gb"`
    SecurityEventsDailyIngestionBudgetInGb types.Number `tfsdk:"security_events_daily_ingestion_budget_in_gb"`
    LogsIngestionRateLimitInMbPerSecond types.Number `tfsdk:"logs_ingestion_rate_limit_in_mb_per_second"`
    TracesIngestionRateLimitInMbPerSecond types.Number `tfsdk:"traces_ingestion_rate_limit_in_mb_per_second"`
    MetricsIngestionRateLimitInMbPerSecond types.Number `tfsdk:"metrics_ingestion_rate_limit_in_mb_per_second"`
    ProfilesIngestionRateLimitInMbPerSecond types.Number `tfsdk:"profiles_ingestion_rate_limit_in_mb_per_second"`
    SecurityEventsIngestionRateLimitInMbPerSecond types.Number `tfsdk:"security_events_ingestion_rate_limit_in_mb_per_second"`
}

func (d *TelemetryIngestionKeyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_telemetry_ingestion_key"
}

func (d *TelemetryIngestionKeyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Manage ingestion keys for your project Look up an existing telemetry ingestion key by `id`, or by any of its other arguments (`name`, `created_by_actor_id`, `created_by_actor_type`, ...): each one set must match, and exactly one telemetry ingestion key may match them all.",

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
            "created_by_actor_type": schema.StringAttribute{
                MarkdownDescription: "What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.",
                Optional: true,
                Computed: true,
            },
            "created_by_actor_id": schema.StringAttribute{
                MarkdownDescription: "ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.",
                Optional: true,
                Computed: true,
            },
            "last_modified_at": schema.StringAttribute{
                MarkdownDescription: "When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.",
                Computed: true,
            },
            "last_modified_by_actor_type": schema.StringAttribute{
                MarkdownDescription: "What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.",
                Optional: true,
                Computed: true,
            },
            "last_modified_by_actor_id": schema.StringAttribute{
                MarkdownDescription: "ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.",
                Optional: true,
                Computed: true,
            },
            "last_modified_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Any friendly name of this object.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Friendly description that will help you remember.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
            "secret_key": schema.StringAttribute{
                MarkdownDescription: "Secret ingestion key.",
                Computed: true,
                Sensitive: true,
            },
            "key_type": schema.StringAttribute{
                MarkdownDescription: "Server keys are for backend services and OpenTelemetry collectors: full ingest, no origin checks. Browser keys are write-only client keys: listed web origins may send traces, logs, metrics, session replays and Sentry events, while listed app:// identities currently authorize React Native session replay only. Mobile keys are write-only keys for one mobile app: pinned to its service name, they send only that app's OTLP traces, logs and metrics over HTTP and its session replays, never from a web page. This cannot be changed after the key is created - create a new key instead.",
                Optional: true,
                Computed: true,
            },
            "allowed_origins": schema.StringAttribute{
                MarkdownDescription: "Web origins (for example https://app.example.com or https://*.example.com) and exact React Native identities (for example app://com.example.mobile) that may use this key. Required on a Browser key. Web requests need a listed Origin; mobile replay requests without Origin need a listed app identity. app:// entries cannot use wildcards and are self-asserted identifiers, not platform attestation. Ignored on a Server key and on a Mobile key, because neither a server nor a mobile app sends an Origin header. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "pinned_service_name": schema.StringAttribute{
                MarkdownDescription: "When set, every OpenTelemetry resource ingested with this key has its service.name REPLACED with this value. This is what stops data written with a scraped key from masquerading as another service: forged spans land in one service you can see and mute, instead of poisoning your backend services' dashboards and alerts. Required on a Mobile key, where it is the service.name of the app the key ships in, and the only RUM application it may send session replays for.",
                Optional: true,
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Turn this off to immediately stop accepting telemetry written with this key, without deleting it. Turn it back on to resume.",
                Optional: true,
                Computed: true,
            },
            "expires_at": schema.StringAttribute{
                MarkdownDescription: "Date and time after which this key stops being accepted. Empty means it never expires, which is how every key behaved before this column existed. Setting one on a Browser key bounds how long a scraped copy stays useful.",
                Computed: true,
            },
            "last_used_at": schema.StringAttribute{
                MarkdownDescription: "The last time telemetry was accepted with this key. Empty means it has never been used since this was recorded. Use it to find keys that are safe to rotate or delete.",
                Computed: true,
            },
            "requests_per_minute_limit": schema.NumberAttribute{
                MarkdownDescription: "Maximum ingest requests per minute accepted with this key. Leave empty to use the shipped default for a Browser or Mobile key, and to leave a Server key unlimited. The limit is per key, across every client using it, so it has to clear your whole fleet - see DEFAULT_BROWSER_KEY_REQUESTS_PER_MINUTE for the default and the reasoning behind its size.",
                Optional: true,
                Computed: true,
            },
            "logs_daily_ingestion_budget_in_gb": schema.NumberAttribute{
                MarkdownDescription: "The most logs accepted with this key in a budget day, in GB, counted decompressed as they arrive; once used up they are refused until the budget resets. At most the project's own budget. Empty: no budget of the key's own (the project's still applies).",
                Optional: true,
                Computed: true,
            },
            "traces_daily_ingestion_budget_in_gb": schema.NumberAttribute{
                MarkdownDescription: "The most traces accepted with this key in a budget day, in GB, counted decompressed as they arrive; once used up they are refused until the budget resets. At most the project's own budget. Empty: no budget of the key's own (the project's still applies).",
                Optional: true,
                Computed: true,
            },
            "metrics_daily_ingestion_budget_in_gb": schema.NumberAttribute{
                MarkdownDescription: "The most metrics accepted with this key in a budget day, in GB, counted decompressed as they arrive; once used up they are refused until the budget resets. At most the project's own budget. Empty: no budget of the key's own (the project's still applies).",
                Optional: true,
                Computed: true,
            },
            "profiles_daily_ingestion_budget_in_gb": schema.NumberAttribute{
                MarkdownDescription: "The most profiles accepted with this key in a budget day, in GB, counted decompressed as they arrive; once used up they are refused until the budget resets. At most the project's own budget. Empty: no budget of the key's own (the project's still applies).",
                Optional: true,
                Computed: true,
            },
            "security_events_daily_ingestion_budget_in_gb": schema.NumberAttribute{
                MarkdownDescription: "The most security events accepted with this key in a budget day, in GB, counted decompressed as they arrive; once used up they are refused until the budget resets. At most the project's own budget. Empty: no budget of the key's own (the project's still applies).",
                Optional: true,
                Computed: true,
            },
            "logs_ingestion_rate_limit_in_mb_per_second": schema.NumberAttribute{
                MarkdownDescription: "The fastest logs may be sent with this key, in MB per second; faster requests are answered 429 with Retry-After. At most the project's own limit. Empty: no limit of the key's own.",
                Optional: true,
                Computed: true,
            },
            "traces_ingestion_rate_limit_in_mb_per_second": schema.NumberAttribute{
                MarkdownDescription: "The fastest traces may be sent with this key, in MB per second; faster requests are answered 429 with Retry-After. At most the project's own limit. Empty: no limit of the key's own.",
                Optional: true,
                Computed: true,
            },
            "metrics_ingestion_rate_limit_in_mb_per_second": schema.NumberAttribute{
                MarkdownDescription: "The fastest metrics may be sent with this key, in MB per second; faster requests are answered 429 with Retry-After. At most the project's own limit. Empty: no limit of the key's own.",
                Optional: true,
                Computed: true,
            },
            "profiles_ingestion_rate_limit_in_mb_per_second": schema.NumberAttribute{
                MarkdownDescription: "The fastest profiles may be sent with this key, in MB per second; faster requests are answered 429 with Retry-After. At most the project's own limit. Empty: no limit of the key's own.",
                Optional: true,
                Computed: true,
            },
            "security_events_ingestion_rate_limit_in_mb_per_second": schema.NumberAttribute{
                MarkdownDescription: "The fastest security events may be sent with this key, in MB per second; faster requests are answered 429 with Retry-After. At most the project's own limit. Empty: no limit of the key's own.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *TelemetryIngestionKeyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TelemetryIngestionKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data TelemetryIngestionKeyDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.CreatedByActorType.IsNull() && !data.CreatedByActorType.IsUnknown() {
        filters["createdByActorType"] = data.CreatedByActorType.ValueString()
        filterNames = append(filterNames, "created_by_actor_type = "+fmt.Sprintf("%q", data.CreatedByActorType.ValueString()))
    }
    if !data.CreatedByActorId.IsNull() && !data.CreatedByActorId.IsUnknown() {
        filters["createdByActorId"] = data.CreatedByActorId.ValueString()
        filterNames = append(filterNames, "created_by_actor_id = "+fmt.Sprintf("%q", data.CreatedByActorId.ValueString()))
    }
    if !data.LastModifiedByActorType.IsNull() && !data.LastModifiedByActorType.IsUnknown() {
        filters["lastModifiedByActorType"] = data.LastModifiedByActorType.ValueString()
        filterNames = append(filterNames, "last_modified_by_actor_type = "+fmt.Sprintf("%q", data.LastModifiedByActorType.ValueString()))
    }
    if !data.LastModifiedByActorId.IsNull() && !data.LastModifiedByActorId.IsUnknown() {
        filters["lastModifiedByActorId"] = data.LastModifiedByActorId.ValueString()
        filterNames = append(filterNames, "last_modified_by_actor_id = "+fmt.Sprintf("%q", data.LastModifiedByActorId.ValueString()))
    }
    if !data.LastModifiedByUserId.IsNull() && !data.LastModifiedByUserId.IsUnknown() {
        filters["lastModifiedByUserId"] = data.LastModifiedByUserId.ValueString()
        filterNames = append(filterNames, "last_modified_by_user_id = "+fmt.Sprintf("%q", data.LastModifiedByUserId.ValueString()))
    }
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.Description.IsNull() && !data.Description.IsUnknown() {
        filters["description"] = data.Description.ValueString()
        filterNames = append(filterNames, "description = "+fmt.Sprintf("%q", data.Description.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.KeyType.IsNull() && !data.KeyType.IsUnknown() {
        filters["keyType"] = data.KeyType.ValueString()
        filterNames = append(filterNames, "key_type = "+fmt.Sprintf("%q", data.KeyType.ValueString()))
    }
    if !data.PinnedServiceName.IsNull() && !data.PinnedServiceName.IsUnknown() {
        filters["pinnedServiceName"] = data.PinnedServiceName.ValueString()
        filterNames = append(filterNames, "pinned_service_name = "+fmt.Sprintf("%q", data.PinnedServiceName.ValueString()))
    }
    if !data.IsEnabled.IsNull() && !data.IsEnabled.IsUnknown() {
        filters["isEnabled"] = data.IsEnabled.ValueBool()
        filterNames = append(filterNames, "is_enabled = "+fmt.Sprintf("%t", data.IsEnabled.ValueBool()))
    }
    if !data.RequestsPerMinuteLimit.IsNull() && !data.RequestsPerMinuteLimit.IsUnknown() {
        filters["requestsPerMinuteLimit"] = LookupNumber(data.RequestsPerMinuteLimit)
        filterNames = append(filterNames, "requests_per_minute_limit = "+data.RequestsPerMinuteLimit.ValueBigFloat().String())
    }
    if !data.LogsDailyIngestionBudgetInGb.IsNull() && !data.LogsDailyIngestionBudgetInGb.IsUnknown() {
        filters["logsDailyIngestionBudgetInGB"] = LookupNumber(data.LogsDailyIngestionBudgetInGb)
        filterNames = append(filterNames, "logs_daily_ingestion_budget_in_gb = "+data.LogsDailyIngestionBudgetInGb.ValueBigFloat().String())
    }
    if !data.TracesDailyIngestionBudgetInGb.IsNull() && !data.TracesDailyIngestionBudgetInGb.IsUnknown() {
        filters["tracesDailyIngestionBudgetInGB"] = LookupNumber(data.TracesDailyIngestionBudgetInGb)
        filterNames = append(filterNames, "traces_daily_ingestion_budget_in_gb = "+data.TracesDailyIngestionBudgetInGb.ValueBigFloat().String())
    }
    if !data.MetricsDailyIngestionBudgetInGb.IsNull() && !data.MetricsDailyIngestionBudgetInGb.IsUnknown() {
        filters["metricsDailyIngestionBudgetInGB"] = LookupNumber(data.MetricsDailyIngestionBudgetInGb)
        filterNames = append(filterNames, "metrics_daily_ingestion_budget_in_gb = "+data.MetricsDailyIngestionBudgetInGb.ValueBigFloat().String())
    }
    if !data.ProfilesDailyIngestionBudgetInGb.IsNull() && !data.ProfilesDailyIngestionBudgetInGb.IsUnknown() {
        filters["profilesDailyIngestionBudgetInGB"] = LookupNumber(data.ProfilesDailyIngestionBudgetInGb)
        filterNames = append(filterNames, "profiles_daily_ingestion_budget_in_gb = "+data.ProfilesDailyIngestionBudgetInGb.ValueBigFloat().String())
    }
    if !data.SecurityEventsDailyIngestionBudgetInGb.IsNull() && !data.SecurityEventsDailyIngestionBudgetInGb.IsUnknown() {
        filters["securityEventsDailyIngestionBudgetInGB"] = LookupNumber(data.SecurityEventsDailyIngestionBudgetInGb)
        filterNames = append(filterNames, "security_events_daily_ingestion_budget_in_gb = "+data.SecurityEventsDailyIngestionBudgetInGb.ValueBigFloat().String())
    }
    if !data.LogsIngestionRateLimitInMbPerSecond.IsNull() && !data.LogsIngestionRateLimitInMbPerSecond.IsUnknown() {
        filters["logsIngestionRateLimitInMBPerSecond"] = LookupNumber(data.LogsIngestionRateLimitInMbPerSecond)
        filterNames = append(filterNames, "logs_ingestion_rate_limit_in_mb_per_second = "+data.LogsIngestionRateLimitInMbPerSecond.ValueBigFloat().String())
    }
    if !data.TracesIngestionRateLimitInMbPerSecond.IsNull() && !data.TracesIngestionRateLimitInMbPerSecond.IsUnknown() {
        filters["tracesIngestionRateLimitInMBPerSecond"] = LookupNumber(data.TracesIngestionRateLimitInMbPerSecond)
        filterNames = append(filterNames, "traces_ingestion_rate_limit_in_mb_per_second = "+data.TracesIngestionRateLimitInMbPerSecond.ValueBigFloat().String())
    }
    if !data.MetricsIngestionRateLimitInMbPerSecond.IsNull() && !data.MetricsIngestionRateLimitInMbPerSecond.IsUnknown() {
        filters["metricsIngestionRateLimitInMBPerSecond"] = LookupNumber(data.MetricsIngestionRateLimitInMbPerSecond)
        filterNames = append(filterNames, "metrics_ingestion_rate_limit_in_mb_per_second = "+data.MetricsIngestionRateLimitInMbPerSecond.ValueBigFloat().String())
    }
    if !data.ProfilesIngestionRateLimitInMbPerSecond.IsNull() && !data.ProfilesIngestionRateLimitInMbPerSecond.IsUnknown() {
        filters["profilesIngestionRateLimitInMBPerSecond"] = LookupNumber(data.ProfilesIngestionRateLimitInMbPerSecond)
        filterNames = append(filterNames, "profiles_ingestion_rate_limit_in_mb_per_second = "+data.ProfilesIngestionRateLimitInMbPerSecond.ValueBigFloat().String())
    }
    if !data.SecurityEventsIngestionRateLimitInMbPerSecond.IsNull() && !data.SecurityEventsIngestionRateLimitInMbPerSecond.IsUnknown() {
        filters["securityEventsIngestionRateLimitInMBPerSecond"] = LookupNumber(data.SecurityEventsIngestionRateLimitInMbPerSecond)
        filterNames = append(filterNames, "security_events_ingestion_rate_limit_in_mb_per_second = "+data.SecurityEventsIngestionRateLimitInMbPerSecond.ValueBigFloat().String())
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the telemetry ingestion key up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the telemetry ingestion key up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "createdByActorType": true,
        "createdByActorId": true,
        "lastModifiedAt": true,
        "lastModifiedByActorType": true,
        "lastModifiedByActorId": true,
        "lastModifiedByUserId": true,
        "projectId": true,
        "name": true,
        "description": true,
        "createdByUserId": true,
        "secretKey": true,
        "keyType": true,
        "allowedOrigins": true,
        "pinnedServiceName": true,
        "isEnabled": true,
        "expiresAt": true,
        "lastUsedAt": true,
        "requestsPerMinuteLimit": true,
        "logsDailyIngestionBudgetInGB": true,
        "tracesDailyIngestionBudgetInGB": true,
        "metricsDailyIngestionBudgetInGB": true,
        "profilesDailyIngestionBudgetInGB": true,
        "securityEventsDailyIngestionBudgetInGB": true,
        "logsIngestionRateLimitInMBPerSecond": true,
        "tracesIngestionRateLimitInMBPerSecond": true,
        "metricsIngestionRateLimitInMBPerSecond": true,
        "profilesIngestionRateLimitInMBPerSecond": true,
        "securityEventsIngestionRateLimitInMBPerSecond": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/telemetry-ingestion-key/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read telemetry_ingestion_key, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No telemetry ingestion key found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read telemetry_ingestion_key: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/telemetry-ingestion-key/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list telemetry_ingestion_key, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list telemetry_ingestion_key: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No telemetry ingestion key matches %s.", DescribeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one telemetry ingestion key matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", DescribeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for telemetry_ingestion_key.")
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
    if obj, ok := item["createdByActorType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByActorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorType = types.StringNull()
        }
    } else if arr, ok := item["createdByActorType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CreatedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorType = types.StringNull()
        }
    } else if val, ok := item["createdByActorType"].(string); ok {
        data.CreatedByActorType = types.StringValue(val)
    } else {
        data.CreatedByActorType = types.StringNull()
    }
    if obj, ok := item["createdByActorId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByActorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorId = types.StringNull()
        }
    } else if arr, ok := item["createdByActorId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.CreatedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByActorId = types.StringNull()
        }
    } else if val, ok := item["createdByActorId"].(string); ok {
        data.CreatedByActorId = types.StringValue(val)
    } else {
        data.CreatedByActorId = types.StringNull()
    }
    if obj, ok := item["lastModifiedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedAt = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedAt = types.StringNull()
        }
    } else if val, ok := item["lastModifiedAt"].(string); ok {
        data.LastModifiedAt = types.StringValue(val)
    } else {
        data.LastModifiedAt = types.StringNull()
    }
    if obj, ok := item["lastModifiedByActorType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedByActorType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedByActorType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorType = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedByActorType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedByActorType = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorType = types.StringNull()
        }
    } else if val, ok := item["lastModifiedByActorType"].(string); ok {
        data.LastModifiedByActorType = types.StringValue(val)
    } else {
        data.LastModifiedByActorType = types.StringNull()
    }
    if obj, ok := item["lastModifiedByActorId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedByActorId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedByActorId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorId = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedByActorId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedByActorId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByActorId = types.StringNull()
        }
    } else if val, ok := item["lastModifiedByActorId"].(string); ok {
        data.LastModifiedByActorId = types.StringValue(val)
    } else {
        data.LastModifiedByActorId = types.StringNull()
    }
    if obj, ok := item["lastModifiedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastModifiedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastModifiedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastModifiedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByUserId = types.StringNull()
        }
    } else if arr, ok := item["lastModifiedByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastModifiedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.LastModifiedByUserId = types.StringNull()
        }
    } else if val, ok := item["lastModifiedByUserId"].(string); ok {
        data.LastModifiedByUserId = types.StringValue(val)
    } else {
        data.LastModifiedByUserId = types.StringNull()
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
    } else if arr, ok := item["name"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    } else if arr, ok := item["description"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
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
    if obj, ok := item["secretKey"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SecretKey = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SecretKey = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SecretKey = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SecretKey = types.StringValue(string(jsonBytes))
        } else {
            data.SecretKey = types.StringNull()
        }
    } else if arr, ok := item["secretKey"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.SecretKey = types.StringValue(string(jsonBytes))
        } else {
            data.SecretKey = types.StringNull()
        }
    } else if val, ok := item["secretKey"].(string); ok {
        data.SecretKey = types.StringValue(val)
    } else {
        data.SecretKey = types.StringNull()
    }
    if obj, ok := item["keyType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.KeyType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.KeyType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.KeyType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.KeyType = types.StringValue(string(jsonBytes))
        } else {
            data.KeyType = types.StringNull()
        }
    } else if arr, ok := item["keyType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.KeyType = types.StringValue(string(jsonBytes))
        } else {
            data.KeyType = types.StringNull()
        }
    } else if val, ok := item["keyType"].(string); ok {
        data.KeyType = types.StringValue(val)
    } else {
        data.KeyType = types.StringNull()
    }
    if obj, ok := item["allowedOrigins"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AllowedOrigins = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AllowedOrigins = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AllowedOrigins = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AllowedOrigins = types.StringValue(string(jsonBytes))
        } else {
            data.AllowedOrigins = types.StringNull()
        }
    } else if arr, ok := item["allowedOrigins"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.AllowedOrigins = types.StringValue(string(jsonBytes))
        } else {
            data.AllowedOrigins = types.StringNull()
        }
    } else if val, ok := item["allowedOrigins"].(string); ok {
        data.AllowedOrigins = types.StringValue(val)
    } else {
        data.AllowedOrigins = types.StringNull()
    }
    if obj, ok := item["pinnedServiceName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PinnedServiceName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PinnedServiceName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PinnedServiceName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PinnedServiceName = types.StringValue(string(jsonBytes))
        } else {
            data.PinnedServiceName = types.StringNull()
        }
    } else if arr, ok := item["pinnedServiceName"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.PinnedServiceName = types.StringValue(string(jsonBytes))
        } else {
            data.PinnedServiceName = types.StringNull()
        }
    } else if val, ok := item["pinnedServiceName"].(string); ok {
        data.PinnedServiceName = types.StringValue(val)
    } else {
        data.PinnedServiceName = types.StringNull()
    }
    if val, ok := item["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    } else {
        data.IsEnabled = types.BoolNull()
    }
    if obj, ok := item["expiresAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ExpiresAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ExpiresAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ExpiresAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ExpiresAt = types.StringValue(string(jsonBytes))
        } else {
            data.ExpiresAt = types.StringNull()
        }
    } else if arr, ok := item["expiresAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ExpiresAt = types.StringValue(string(jsonBytes))
        } else {
            data.ExpiresAt = types.StringNull()
        }
    } else if val, ok := item["expiresAt"].(string); ok {
        data.ExpiresAt = types.StringValue(val)
    } else {
        data.ExpiresAt = types.StringNull()
    }
    if obj, ok := item["lastUsedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.LastUsedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.LastUsedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.LastUsedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.LastUsedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastUsedAt = types.StringNull()
        }
    } else if arr, ok := item["lastUsedAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.LastUsedAt = types.StringValue(string(jsonBytes))
        } else {
            data.LastUsedAt = types.StringNull()
        }
    } else if val, ok := item["lastUsedAt"].(string); ok {
        data.LastUsedAt = types.StringValue(val)
    } else {
        data.LastUsedAt = types.StringNull()
    }
    data.RequestsPerMinuteLimit = NumberFromAPI(data.RequestsPerMinuteLimit, item["requestsPerMinuteLimit"])
    data.LogsDailyIngestionBudgetInGb = NumberFromAPI(data.LogsDailyIngestionBudgetInGb, item["logsDailyIngestionBudgetInGB"])
    data.TracesDailyIngestionBudgetInGb = NumberFromAPI(data.TracesDailyIngestionBudgetInGb, item["tracesDailyIngestionBudgetInGB"])
    data.MetricsDailyIngestionBudgetInGb = NumberFromAPI(data.MetricsDailyIngestionBudgetInGb, item["metricsDailyIngestionBudgetInGB"])
    data.ProfilesDailyIngestionBudgetInGb = NumberFromAPI(data.ProfilesDailyIngestionBudgetInGb, item["profilesDailyIngestionBudgetInGB"])
    data.SecurityEventsDailyIngestionBudgetInGb = NumberFromAPI(data.SecurityEventsDailyIngestionBudgetInGb, item["securityEventsDailyIngestionBudgetInGB"])
    data.LogsIngestionRateLimitInMbPerSecond = NumberFromAPI(data.LogsIngestionRateLimitInMbPerSecond, item["logsIngestionRateLimitInMBPerSecond"])
    data.TracesIngestionRateLimitInMbPerSecond = NumberFromAPI(data.TracesIngestionRateLimitInMbPerSecond, item["tracesIngestionRateLimitInMBPerSecond"])
    data.MetricsIngestionRateLimitInMbPerSecond = NumberFromAPI(data.MetricsIngestionRateLimitInMbPerSecond, item["metricsIngestionRateLimitInMBPerSecond"])
    data.ProfilesIngestionRateLimitInMbPerSecond = NumberFromAPI(data.ProfilesIngestionRateLimitInMbPerSecond, item["profilesIngestionRateLimitInMBPerSecond"])
    data.SecurityEventsIngestionRateLimitInMbPerSecond = NumberFromAPI(data.SecurityEventsIngestionRateLimitInMbPerSecond, item["securityEventsIngestionRateLimitInMBPerSecond"])

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
