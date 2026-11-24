package shard20

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
var _ datasource.DataSource = &ComplianceSettingDataSource{}

func NewComplianceSettingDataSource() datasource.DataSource {
    return &ComplianceSettingDataSource{}
}

// ComplianceSettingDataSource defines the data source implementation.
type ComplianceSettingDataSource struct {
    client *Client
}

// ComplianceSettingDataSourceModel describes the data source data model.
type ComplianceSettingDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    ResourceType types.String `tfsdk:"resource_type"`
    ResourceId types.String `tfsdk:"resource_id"`
    Framework types.String `tfsdk:"framework"`
    Requirement types.String `tfsdk:"requirement"`
    Namespace types.String `tfsdk:"namespace"`
    ObjectKind types.String `tfsdk:"object_kind"`
    ObjectName types.String `tfsdk:"object_name"`
    ExpiresAt types.String `tfsdk:"expires_at"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    Note types.String `tfsdk:"note"`
    TurnedOffByUserId types.String `tfsdk:"turned_off_by_user_id"`
    TurnedOffAt types.String `tfsdk:"turned_off_at"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *ComplianceSettingDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_compliance_setting"
}

func (d *ComplianceSettingDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "A compliance framework, or one of its requirements, a resource turned off for its own data stores. Look up an existing compliance setting by `id`, or by any of its other arguments (`created_by_user_id`, `framework`, `is_enabled`, ...): each one set must match, and exactly one compliance setting may match them all.",

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
                MarkdownDescription: "ID of the project this compliance setting belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "resource_type": schema.StringAttribute{
                MarkdownDescription: "The kind of resource that turned the framework off: KubernetesCluster, DockerHost, PodmanHost, DatabaseServer, CloudResource or Service.",
                Optional: true,
                Computed: true,
            },
            "resource_id": schema.StringAttribute{
                MarkdownDescription: "The ID of the resource that turned the framework off.",
                Optional: true,
                Computed: true,
            },
            "framework": schema.StringAttribute{
                MarkdownDescription: "The compliance framework this setting changes for the resource: GDPR, CCPA, HIPAA, PCI DSS, SOC 2, ISO 27001, NIST CSF, LGPD or DPDP, or for a Kubernetes cluster PSS Baseline, PSS Restricted or CIS Kubernetes.",
                Optional: true,
                Computed: true,
            },
            "requirement": schema.StringAttribute{
                MarkdownDescription: "One requirement of the framework, as Data Security names it (for example Art. 32 or CC6.1). Empty for the whole framework.",
                Optional: true,
                Computed: true,
            },
            "namespace": schema.StringAttribute{
                MarkdownDescription: "For a Kubernetes cluster's PSS Baseline, PSS Restricted or CIS Kubernetes setting, the one namespace it is for. Empty for the whole resource.",
                Optional: true,
                Computed: true,
            },
            "object_kind": schema.StringAttribute{
                MarkdownDescription: "For an accepted Kubernetes finding, the kind of the object it is accepted for: Deployment, DaemonSet, CronJob, ClusterRoleBinding and the others. Empty for a framework, a control or a namespace.",
                Optional: true,
                Computed: true,
            },
            "object_name": schema.StringAttribute{
                MarkdownDescription: "For an accepted Kubernetes finding, the name of the object it is accepted for, as Kubernetes names it. Empty for a framework, a control or a namespace.",
                Optional: true,
                Computed: true,
            },
            "expires_at": schema.StringAttribute{
                MarkdownDescription: "For an accepted Kubernetes finding, when it fails again. Empty: it stays accepted until revoked.",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether the framework or requirement applies to the resource. Off takes the resource's data stores out of its scope, and for an accepted Kubernetes finding, accepts it.",
                Optional: true,
                Computed: true,
            },
            "note": schema.StringAttribute{
                MarkdownDescription: "Why it does not apply to the resource, or why the Kubernetes finding is accepted.",
                Optional: true,
                Computed: true,
            },
            "turned_off_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the user who last turned the framework or requirement off for the resource, or accepted the finding. Empty while it applies, or when an API key turned it off. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
            "turned_off_at": schema.StringAttribute{
                MarkdownDescription: "When the framework or requirement was last turned off for the resource, or the finding accepted. Empty while it applies.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the user who created this compliance setting. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: DataSourceRecordIDValidators("oneuptime_user"),
            },
        },
    }
}

func (d *ComplianceSettingDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ComplianceSettingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ComplianceSettingDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.ResourceType.IsNull() && !data.ResourceType.IsUnknown() {
        filters["resourceType"] = data.ResourceType.ValueString()
        filterNames = append(filterNames, "resource_type = "+fmt.Sprintf("%q", data.ResourceType.ValueString()))
    }
    if !data.ResourceId.IsNull() && !data.ResourceId.IsUnknown() {
        filters["resourceId"] = data.ResourceId.ValueString()
        filterNames = append(filterNames, "resource_id = "+fmt.Sprintf("%q", data.ResourceId.ValueString()))
    }
    if !data.Framework.IsNull() && !data.Framework.IsUnknown() {
        filters["framework"] = data.Framework.ValueString()
        filterNames = append(filterNames, "framework = "+fmt.Sprintf("%q", data.Framework.ValueString()))
    }
    if !data.Requirement.IsNull() && !data.Requirement.IsUnknown() {
        filters["requirement"] = data.Requirement.ValueString()
        filterNames = append(filterNames, "requirement = "+fmt.Sprintf("%q", data.Requirement.ValueString()))
    }
    if !data.Namespace.IsNull() && !data.Namespace.IsUnknown() {
        filters["namespace"] = data.Namespace.ValueString()
        filterNames = append(filterNames, "namespace = "+fmt.Sprintf("%q", data.Namespace.ValueString()))
    }
    if !data.ObjectKind.IsNull() && !data.ObjectKind.IsUnknown() {
        filters["objectKind"] = data.ObjectKind.ValueString()
        filterNames = append(filterNames, "object_kind = "+fmt.Sprintf("%q", data.ObjectKind.ValueString()))
    }
    if !data.ObjectName.IsNull() && !data.ObjectName.IsUnknown() {
        filters["objectName"] = data.ObjectName.ValueString()
        filterNames = append(filterNames, "object_name = "+fmt.Sprintf("%q", data.ObjectName.ValueString()))
    }
    if !data.IsEnabled.IsNull() && !data.IsEnabled.IsUnknown() {
        filters["isEnabled"] = data.IsEnabled.ValueBool()
        filterNames = append(filterNames, "is_enabled = "+fmt.Sprintf("%t", data.IsEnabled.ValueBool()))
    }
    if !data.Note.IsNull() && !data.Note.IsUnknown() {
        filters["note"] = data.Note.ValueString()
        filterNames = append(filterNames, "note = "+fmt.Sprintf("%q", data.Note.ValueString()))
    }
    if !data.TurnedOffByUserId.IsNull() && !data.TurnedOffByUserId.IsUnknown() {
        filters["turnedOffByUserId"] = data.TurnedOffByUserId.ValueString()
        filterNames = append(filterNames, "turned_off_by_user_id = "+fmt.Sprintf("%q", data.TurnedOffByUserId.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the compliance setting up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the compliance setting up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "resourceType": true,
        "resourceId": true,
        "framework": true,
        "requirement": true,
        "namespace": true,
        "objectKind": true,
        "objectName": true,
        "expiresAt": true,
        "isEnabled": true,
        "note": true,
        "turnedOffByUserId": true,
        "turnedOffAt": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/data-security-compliance-setting/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read compliance_setting, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No compliance setting found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read compliance_setting: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/data-security-compliance-setting/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list compliance_setting, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list compliance_setting: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No compliance setting matches %s.", DescribeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one compliance setting matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", DescribeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for compliance_setting.")
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
    if obj, ok := item["resourceType"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResourceType = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ResourceType = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ResourceType = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ResourceType = types.StringValue(string(jsonBytes))
        } else {
            data.ResourceType = types.StringNull()
        }
    } else if arr, ok := item["resourceType"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ResourceType = types.StringValue(string(jsonBytes))
        } else {
            data.ResourceType = types.StringNull()
        }
    } else if val, ok := item["resourceType"].(string); ok {
        data.ResourceType = types.StringValue(val)
    } else {
        data.ResourceType = types.StringNull()
    }
    if obj, ok := item["resourceId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResourceId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ResourceId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ResourceId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ResourceId = types.StringValue(string(jsonBytes))
        } else {
            data.ResourceId = types.StringNull()
        }
    } else if arr, ok := item["resourceId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ResourceId = types.StringValue(string(jsonBytes))
        } else {
            data.ResourceId = types.StringNull()
        }
    } else if val, ok := item["resourceId"].(string); ok {
        data.ResourceId = types.StringValue(val)
    } else {
        data.ResourceId = types.StringNull()
    }
    if obj, ok := item["framework"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Framework = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Framework = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Framework = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Framework = types.StringValue(string(jsonBytes))
        } else {
            data.Framework = types.StringNull()
        }
    } else if arr, ok := item["framework"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Framework = types.StringValue(string(jsonBytes))
        } else {
            data.Framework = types.StringNull()
        }
    } else if val, ok := item["framework"].(string); ok {
        data.Framework = types.StringValue(val)
    } else {
        data.Framework = types.StringNull()
    }
    if obj, ok := item["requirement"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Requirement = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Requirement = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Requirement = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Requirement = types.StringValue(string(jsonBytes))
        } else {
            data.Requirement = types.StringNull()
        }
    } else if arr, ok := item["requirement"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Requirement = types.StringValue(string(jsonBytes))
        } else {
            data.Requirement = types.StringNull()
        }
    } else if val, ok := item["requirement"].(string); ok {
        data.Requirement = types.StringValue(val)
    } else {
        data.Requirement = types.StringNull()
    }
    if obj, ok := item["namespace"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Namespace = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Namespace = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Namespace = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Namespace = types.StringValue(string(jsonBytes))
        } else {
            data.Namespace = types.StringNull()
        }
    } else if arr, ok := item["namespace"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Namespace = types.StringValue(string(jsonBytes))
        } else {
            data.Namespace = types.StringNull()
        }
    } else if val, ok := item["namespace"].(string); ok {
        data.Namespace = types.StringValue(val)
    } else {
        data.Namespace = types.StringNull()
    }
    if obj, ok := item["objectKind"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ObjectKind = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ObjectKind = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ObjectKind = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ObjectKind = types.StringValue(string(jsonBytes))
        } else {
            data.ObjectKind = types.StringNull()
        }
    } else if arr, ok := item["objectKind"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ObjectKind = types.StringValue(string(jsonBytes))
        } else {
            data.ObjectKind = types.StringNull()
        }
    } else if val, ok := item["objectKind"].(string); ok {
        data.ObjectKind = types.StringValue(val)
    } else {
        data.ObjectKind = types.StringNull()
    }
    if obj, ok := item["objectName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ObjectName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ObjectName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ObjectName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ObjectName = types.StringValue(string(jsonBytes))
        } else {
            data.ObjectName = types.StringNull()
        }
    } else if arr, ok := item["objectName"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.ObjectName = types.StringValue(string(jsonBytes))
        } else {
            data.ObjectName = types.StringNull()
        }
    } else if val, ok := item["objectName"].(string); ok {
        data.ObjectName = types.StringValue(val)
    } else {
        data.ObjectName = types.StringNull()
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
    if val, ok := item["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    } else {
        data.IsEnabled = types.BoolNull()
    }
    if obj, ok := item["note"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Note = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Note = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Note = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Note = types.StringValue(string(jsonBytes))
        } else {
            data.Note = types.StringNull()
        }
    } else if arr, ok := item["note"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.Note = types.StringValue(string(jsonBytes))
        } else {
            data.Note = types.StringNull()
        }
    } else if val, ok := item["note"].(string); ok {
        data.Note = types.StringValue(val)
    } else {
        data.Note = types.StringNull()
    }
    if obj, ok := item["turnedOffByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TurnedOffByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TurnedOffByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TurnedOffByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TurnedOffByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.TurnedOffByUserId = types.StringNull()
        }
    } else if arr, ok := item["turnedOffByUserId"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.TurnedOffByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.TurnedOffByUserId = types.StringNull()
        }
    } else if val, ok := item["turnedOffByUserId"].(string); ok {
        data.TurnedOffByUserId = types.StringValue(val)
    } else {
        data.TurnedOffByUserId = types.StringNull()
    }
    if obj, ok := item["turnedOffAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TurnedOffAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TurnedOffAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TurnedOffAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TurnedOffAt = types.StringValue(string(jsonBytes))
        } else {
            data.TurnedOffAt = types.StringNull()
        }
    } else if arr, ok := item["turnedOffAt"].([]interface{}); ok {
        // A JSON array column, such as a runbook's steps: the array, as JSON
        if jsonBytes, err := json.Marshal(arr); err == nil {
            data.TurnedOffAt = types.StringValue(string(jsonBytes))
        } else {
            data.TurnedOffAt = types.StringNull()
        }
    } else if val, ok := item["turnedOffAt"].(string); ok {
        data.TurnedOffAt = types.StringValue(val)
    } else {
        data.TurnedOffAt = types.StringNull()
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

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
