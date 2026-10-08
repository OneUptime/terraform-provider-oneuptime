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
var _ datasource.DataSource = &IncidentOwnerRuleDataSource{}

func NewIncidentOwnerRuleDataSource() datasource.DataSource {
    return &IncidentOwnerRuleDataSource{}
}

// IncidentOwnerRuleDataSource defines the data source implementation.
type IncidentOwnerRuleDataSource struct {
    client *Client
}

// IncidentOwnerRuleDataSourceModel describes the data source data model.
type IncidentOwnerRuleDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    Criteria types.String `tfsdk:"criteria"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    NotifyOwners types.Bool `tfsdk:"notify_owners"`
    Monitors types.Set `tfsdk:"monitors"`
    IncidentSeverities types.Set `tfsdk:"incident_severities"`
    IncidentLabels types.Set `tfsdk:"incident_labels"`
    MonitorLabels types.Set `tfsdk:"monitor_labels"`
    IncidentTitlePattern types.String `tfsdk:"incident_title_pattern"`
    IncidentDescriptionPattern types.String `tfsdk:"incident_description_pattern"`
    MonitorNamePattern types.String `tfsdk:"monitor_name_pattern"`
    MonitorDescriptionPattern types.String `tfsdk:"monitor_description_pattern"`
    OwnerUsers types.Set `tfsdk:"owner_users"`
    OwnerTeams types.Set `tfsdk:"owner_teams"`
    InheritOwnersFromMonitors types.Bool `tfsdk:"inherit_owners_from_monitors"`
    InheritOwnersFromHosts types.Bool `tfsdk:"inherit_owners_from_hosts"`
    InheritOwnersFromKubernetesClusters types.Bool `tfsdk:"inherit_owners_from_kubernetes_clusters"`
    InheritOwnersFromDockerHosts types.Bool `tfsdk:"inherit_owners_from_docker_hosts"`
    InheritOwnersFromPodmanHosts types.Bool `tfsdk:"inherit_owners_from_podman_hosts"`
    InheritOwnersFromServices types.Bool `tfsdk:"inherit_owners_from_services"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *IncidentOwnerRuleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_incident_owner_rule"
}

func (d *IncidentOwnerRuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Configure rules for automatically assigning owner users and teams when matching incidents are created Look up an existing incident owner rule by `id`, or by any of its other arguments (`name`, `created_by_user_id`, `description`, ...): each one set must match, and exactly one incident owner rule may match them all.",

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
            "criteria": schema.StringAttribute{
                MarkdownDescription: "Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name of this incident owner rule.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of this incident owner rule.",
                Optional: true,
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Whether this rule is enabled.",
                Optional: true,
                Computed: true,
            },
            "notify_owners": schema.BoolAttribute{
                MarkdownDescription: "Send notifications to owner users and teams when they are added by this rule.",
                Optional: true,
                Computed: true,
            },
            "monitors": schema.SetAttribute{
                MarkdownDescription: "Only trigger for incidents from these monitors. Leave empty to match incidents from any monitor. IDs of `oneuptime_monitor` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "incident_severities": schema.SetAttribute{
                MarkdownDescription: "Only trigger for incidents with these severities. Leave empty to match incidents of any severity. IDs of `oneuptime_incident_severity` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "incident_labels": schema.SetAttribute{
                MarkdownDescription: "Only trigger for incidents that have at least one of these labels. Leave empty to match regardless of incident labels. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "monitor_labels": schema.SetAttribute{
                MarkdownDescription: "Only trigger for incidents from monitors that have at least one of these labels. Leave empty to match regardless of monitor labels. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "incident_title_pattern": schema.StringAttribute{
                MarkdownDescription: "Regex (case-insensitive) matched against the incident title. Leave empty to match any title.",
                Optional: true,
                Computed: true,
            },
            "incident_description_pattern": schema.StringAttribute{
                MarkdownDescription: "Regex (case-insensitive) matched against the incident description. Leave empty to match any description.",
                Optional: true,
                Computed: true,
            },
            "monitor_name_pattern": schema.StringAttribute{
                MarkdownDescription: "Regex (case-insensitive) matched against any of the incident's monitor names. Leave empty to match any monitor.",
                Optional: true,
                Computed: true,
            },
            "monitor_description_pattern": schema.StringAttribute{
                MarkdownDescription: "Regex (case-insensitive) matched against any of the incident's monitor descriptions. Leave empty to match any description.",
                Optional: true,
                Computed: true,
            },
            "owner_users": schema.SetAttribute{
                MarkdownDescription: "Users to add as owners on the incident when this rule matches. IDs of `oneuptime_user` records.",
                Computed: true,
                ElementType: types.StringType,
            },
            "owner_teams": schema.SetAttribute{
                MarkdownDescription: "Teams to add as owners on the incident when this rule matches. IDs of `oneuptime_team` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "inherit_owners_from_monitors": schema.BoolAttribute{
                MarkdownDescription: "When this rule matches, also assign every owner of the incident's monitors to the incident.",
                Optional: true,
                Computed: true,
            },
            "inherit_owners_from_hosts": schema.BoolAttribute{
                MarkdownDescription: "When this rule matches, also assign every owner of the incident's affected hosts to the incident.",
                Optional: true,
                Computed: true,
            },
            "inherit_owners_from_kubernetes_clusters": schema.BoolAttribute{
                MarkdownDescription: "When this rule matches, also assign every owner of the incident's affected Kubernetes clusters to the incident.",
                Optional: true,
                Computed: true,
            },
            "inherit_owners_from_docker_hosts": schema.BoolAttribute{
                MarkdownDescription: "When this rule matches, also assign every owner of the incident's affected Docker hosts to the incident.",
                Optional: true,
                Computed: true,
            },
            "inherit_owners_from_podman_hosts": schema.BoolAttribute{
                MarkdownDescription: "When this rule matches, also assign every owner of the incident's affected Podman hosts to the incident.",
                Optional: true,
                Computed: true,
            },
            "inherit_owners_from_services": schema.BoolAttribute{
                MarkdownDescription: "When this rule matches, also assign every owner of the incident's affected services to the incident.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *IncidentOwnerRuleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IncidentOwnerRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data IncidentOwnerRuleDataSourceModel

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
    if !data.IsEnabled.IsNull() && !data.IsEnabled.IsUnknown() {
        filters["isEnabled"] = data.IsEnabled.ValueBool()
        filterNames = append(filterNames, "is_enabled = "+fmt.Sprintf("%t", data.IsEnabled.ValueBool()))
    }
    if !data.NotifyOwners.IsNull() && !data.NotifyOwners.IsUnknown() {
        filters["notifyOwners"] = data.NotifyOwners.ValueBool()
        filterNames = append(filterNames, "notify_owners = "+fmt.Sprintf("%t", data.NotifyOwners.ValueBool()))
    }
    if !data.IncidentTitlePattern.IsNull() && !data.IncidentTitlePattern.IsUnknown() {
        filters["incidentTitlePattern"] = data.IncidentTitlePattern.ValueString()
        filterNames = append(filterNames, "incident_title_pattern = "+fmt.Sprintf("%q", data.IncidentTitlePattern.ValueString()))
    }
    if !data.IncidentDescriptionPattern.IsNull() && !data.IncidentDescriptionPattern.IsUnknown() {
        filters["incidentDescriptionPattern"] = data.IncidentDescriptionPattern.ValueString()
        filterNames = append(filterNames, "incident_description_pattern = "+fmt.Sprintf("%q", data.IncidentDescriptionPattern.ValueString()))
    }
    if !data.MonitorNamePattern.IsNull() && !data.MonitorNamePattern.IsUnknown() {
        filters["monitorNamePattern"] = data.MonitorNamePattern.ValueString()
        filterNames = append(filterNames, "monitor_name_pattern = "+fmt.Sprintf("%q", data.MonitorNamePattern.ValueString()))
    }
    if !data.MonitorDescriptionPattern.IsNull() && !data.MonitorDescriptionPattern.IsUnknown() {
        filters["monitorDescriptionPattern"] = data.MonitorDescriptionPattern.ValueString()
        filterNames = append(filterNames, "monitor_description_pattern = "+fmt.Sprintf("%q", data.MonitorDescriptionPattern.ValueString()))
    }
    if !data.InheritOwnersFromMonitors.IsNull() && !data.InheritOwnersFromMonitors.IsUnknown() {
        filters["inheritOwnersFromMonitors"] = data.InheritOwnersFromMonitors.ValueBool()
        filterNames = append(filterNames, "inherit_owners_from_monitors = "+fmt.Sprintf("%t", data.InheritOwnersFromMonitors.ValueBool()))
    }
    if !data.InheritOwnersFromHosts.IsNull() && !data.InheritOwnersFromHosts.IsUnknown() {
        filters["inheritOwnersFromHosts"] = data.InheritOwnersFromHosts.ValueBool()
        filterNames = append(filterNames, "inherit_owners_from_hosts = "+fmt.Sprintf("%t", data.InheritOwnersFromHosts.ValueBool()))
    }
    if !data.InheritOwnersFromKubernetesClusters.IsNull() && !data.InheritOwnersFromKubernetesClusters.IsUnknown() {
        filters["inheritOwnersFromKubernetesClusters"] = data.InheritOwnersFromKubernetesClusters.ValueBool()
        filterNames = append(filterNames, "inherit_owners_from_kubernetes_clusters = "+fmt.Sprintf("%t", data.InheritOwnersFromKubernetesClusters.ValueBool()))
    }
    if !data.InheritOwnersFromDockerHosts.IsNull() && !data.InheritOwnersFromDockerHosts.IsUnknown() {
        filters["inheritOwnersFromDockerHosts"] = data.InheritOwnersFromDockerHosts.ValueBool()
        filterNames = append(filterNames, "inherit_owners_from_docker_hosts = "+fmt.Sprintf("%t", data.InheritOwnersFromDockerHosts.ValueBool()))
    }
    if !data.InheritOwnersFromPodmanHosts.IsNull() && !data.InheritOwnersFromPodmanHosts.IsUnknown() {
        filters["inheritOwnersFromPodmanHosts"] = data.InheritOwnersFromPodmanHosts.ValueBool()
        filterNames = append(filterNames, "inherit_owners_from_podman_hosts = "+fmt.Sprintf("%t", data.InheritOwnersFromPodmanHosts.ValueBool()))
    }
    if !data.InheritOwnersFromServices.IsNull() && !data.InheritOwnersFromServices.IsUnknown() {
        filters["inheritOwnersFromServices"] = data.InheritOwnersFromServices.ValueBool()
        filterNames = append(filterNames, "inherit_owners_from_services = "+fmt.Sprintf("%t", data.InheritOwnersFromServices.ValueBool()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the incident owner rule up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the incident owner rule up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "criteria": true,
        "projectId": true,
        "name": true,
        "description": true,
        "isEnabled": true,
        "notifyOwners": true,
        "monitors": true,
        "incidentSeverities": true,
        "incidentLabels": true,
        "monitorLabels": true,
        "incidentTitlePattern": true,
        "incidentDescriptionPattern": true,
        "monitorNamePattern": true,
        "monitorDescriptionPattern": true,
        "ownerUsers": true,
        "ownerTeams": true,
        "inheritOwnersFromMonitors": true,
        "inheritOwnersFromHosts": true,
        "inheritOwnersFromKubernetesClusters": true,
        "inheritOwnersFromDockerHosts": true,
        "inheritOwnersFromPodmanHosts": true,
        "inheritOwnersFromServices": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/incident-owner-rule/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read incident_owner_rule, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No incident owner rule found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read incident_owner_rule: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/incident-owner-rule/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list incident_owner_rule, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list incident_owner_rule: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No incident owner rule matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one incident owner rule matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for incident_owner_rule.")
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
    if obj, ok := item["criteria"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Criteria = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Criteria = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Criteria = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Criteria = types.StringValue(string(jsonBytes))
        } else {
            data.Criteria = types.StringNull()
        }
    } else if val, ok := item["criteria"].(string); ok {
        data.Criteria = types.StringValue(val)
    } else {
        data.Criteria = types.StringNull()
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
    if val, ok := item["isEnabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    } else {
        data.IsEnabled = types.BoolNull()
    }
    if val, ok := item["notifyOwners"].(bool); ok {
        data.NotifyOwners = types.BoolValue(val)
    } else {
        data.NotifyOwners = types.BoolNull()
    }
    if val, ok := item["monitors"].([]interface{}); ok {
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
        data.Monitors = types.SetValueMust(types.StringType, setItems)
    } else {
        data.Monitors = types.SetNull(types.StringType)
    }
    if val, ok := item["incidentSeverities"].([]interface{}); ok {
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
        data.IncidentSeverities = types.SetValueMust(types.StringType, setItems)
    } else {
        data.IncidentSeverities = types.SetNull(types.StringType)
    }
    if val, ok := item["incidentLabels"].([]interface{}); ok {
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
        data.IncidentLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        data.IncidentLabels = types.SetNull(types.StringType)
    }
    if val, ok := item["monitorLabels"].([]interface{}); ok {
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
        data.MonitorLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        data.MonitorLabels = types.SetNull(types.StringType)
    }
    if obj, ok := item["incidentTitlePattern"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentTitlePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncidentTitlePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncidentTitlePattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncidentTitlePattern = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentTitlePattern = types.StringNull()
        }
    } else if val, ok := item["incidentTitlePattern"].(string); ok {
        data.IncidentTitlePattern = types.StringValue(val)
    } else {
        data.IncidentTitlePattern = types.StringNull()
    }
    if obj, ok := item["incidentDescriptionPattern"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentDescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncidentDescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncidentDescriptionPattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncidentDescriptionPattern = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentDescriptionPattern = types.StringNull()
        }
    } else if val, ok := item["incidentDescriptionPattern"].(string); ok {
        data.IncidentDescriptionPattern = types.StringValue(val)
    } else {
        data.IncidentDescriptionPattern = types.StringNull()
    }
    if obj, ok := item["monitorNamePattern"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorNamePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MonitorNamePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MonitorNamePattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MonitorNamePattern = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorNamePattern = types.StringNull()
        }
    } else if val, ok := item["monitorNamePattern"].(string); ok {
        data.MonitorNamePattern = types.StringValue(val)
    } else {
        data.MonitorNamePattern = types.StringNull()
    }
    if obj, ok := item["monitorDescriptionPattern"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MonitorDescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MonitorDescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MonitorDescriptionPattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MonitorDescriptionPattern = types.StringValue(string(jsonBytes))
        } else {
            data.MonitorDescriptionPattern = types.StringNull()
        }
    } else if val, ok := item["monitorDescriptionPattern"].(string); ok {
        data.MonitorDescriptionPattern = types.StringValue(val)
    } else {
        data.MonitorDescriptionPattern = types.StringNull()
    }
    if val, ok := item["ownerUsers"].([]interface{}); ok {
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
        data.OwnerUsers = types.SetValueMust(types.StringType, setItems)
    } else {
        data.OwnerUsers = types.SetNull(types.StringType)
    }
    if val, ok := item["ownerTeams"].([]interface{}); ok {
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
        data.OwnerTeams = types.SetValueMust(types.StringType, setItems)
    } else {
        data.OwnerTeams = types.SetNull(types.StringType)
    }
    if val, ok := item["inheritOwnersFromMonitors"].(bool); ok {
        data.InheritOwnersFromMonitors = types.BoolValue(val)
    } else {
        data.InheritOwnersFromMonitors = types.BoolNull()
    }
    if val, ok := item["inheritOwnersFromHosts"].(bool); ok {
        data.InheritOwnersFromHosts = types.BoolValue(val)
    } else {
        data.InheritOwnersFromHosts = types.BoolNull()
    }
    if val, ok := item["inheritOwnersFromKubernetesClusters"].(bool); ok {
        data.InheritOwnersFromKubernetesClusters = types.BoolValue(val)
    } else {
        data.InheritOwnersFromKubernetesClusters = types.BoolNull()
    }
    if val, ok := item["inheritOwnersFromDockerHosts"].(bool); ok {
        data.InheritOwnersFromDockerHosts = types.BoolValue(val)
    } else {
        data.InheritOwnersFromDockerHosts = types.BoolNull()
    }
    if val, ok := item["inheritOwnersFromPodmanHosts"].(bool); ok {
        data.InheritOwnersFromPodmanHosts = types.BoolValue(val)
    } else {
        data.InheritOwnersFromPodmanHosts = types.BoolNull()
    }
    if val, ok := item["inheritOwnersFromServices"].(bool); ok {
        data.InheritOwnersFromServices = types.BoolValue(val)
    } else {
        data.InheritOwnersFromServices = types.BoolNull()
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
