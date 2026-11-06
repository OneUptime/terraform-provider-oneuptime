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
var _ datasource.DataSource = &OpenstackCloudOwnerRuleDataSource{}

func NewOpenstackCloudOwnerRuleDataSource() datasource.DataSource {
    return &OpenstackCloudOwnerRuleDataSource{}
}

// NewOpenstackCloudOwnerRuleLegacyDataSource registers this data source under the
// name it had before type names kept mixed-case words whole,
// oneuptime_open_stack_cloud_owner_rule. Deprecated.
func NewOpenstackCloudOwnerRuleLegacyDataSource() datasource.DataSource {
    return &OpenstackCloudOwnerRuleDataSource{isLegacyAlias: true}
}

// OpenstackCloudOwnerRuleDataSource defines the data source implementation.
type OpenstackCloudOwnerRuleDataSource struct {
    client *Client
    // Registered under the name this data source had before.
    isLegacyAlias bool
}

// OpenstackCloudOwnerRuleDataSourceModel describes the data source data model.
type OpenstackCloudOwnerRuleDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    Criteria types.String `tfsdk:"criteria"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    NotifyOwners types.Bool `tfsdk:"notify_owners"`
    OpenStackCloudLabels types.Set `tfsdk:"open_stack_cloud_labels"`
    OpenStackCloudNamePattern types.String `tfsdk:"open_stack_cloud_name_pattern"`
    OpenStackCloudDescriptionPattern types.String `tfsdk:"open_stack_cloud_description_pattern"`
    OpenStackProjectId types.String `tfsdk:"open_stack_project_id"`
    OpenStackDomainId types.String `tfsdk:"open_stack_domain_id"`
    OwnerUsers types.Set `tfsdk:"owner_users"`
    OwnerTeams types.Set `tfsdk:"owner_teams"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *OpenstackCloudOwnerRuleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    if d.isLegacyAlias {
        resp.TypeName = req.ProviderTypeName + "_open_stack_cloud_owner_rule"
        return
    }
    resp.TypeName = req.ProviderTypeName + "_openstack_cloud_owner_rule"
}

func (d *OpenstackCloudOwnerRuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Configure rules for automatically assigning owner users and teams when matching OpenStack clouds are created Look up an existing openstack cloud owner rule by `id`, or by any of its other arguments (`name`, `created_by_user_id`, `description`, ...): each one set must match, and exactly one openstack cloud owner rule may match them all.",

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
                MarkdownDescription: "Versioned conditions that determine whether this rule matches a resource. A JSON value: read it with `jsondecode()`.",
                Computed: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name of this OpenStack cloud owner rule.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of this OpenStack cloud owner rule.",
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
            "open_stack_cloud_labels": schema.SetAttribute{
                MarkdownDescription: "Only trigger for OpenStack clouds that have at least one of these labels. Leave empty to match regardless of labels. IDs of `oneuptime_label` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "open_stack_cloud_name_pattern": schema.StringAttribute{
                MarkdownDescription: "Regex (case-insensitive) matched against the OpenStack cloud name. Leave empty to match any name.",
                Optional: true,
                Computed: true,
            },
            "open_stack_cloud_description_pattern": schema.StringAttribute{
                MarkdownDescription: "Regex (case-insensitive) matched against the OpenStack cloud description. Leave empty to match any description.",
                Optional: true,
                Computed: true,
            },
            "open_stack_project_id": schema.StringAttribute{
                MarkdownDescription: "Keystone's ID for one project. Set, this rule adds its owners to every alert and incident that reaches that project instead of to the cloud.",
                Optional: true,
                Computed: true,
            },
            "open_stack_domain_id": schema.StringAttribute{
                MarkdownDescription: "Keystone's ID for one domain. Set, this rule adds its owners to every alert and incident that reaches a project in that domain instead of to the cloud.",
                Optional: true,
                Computed: true,
            },
            "owner_users": schema.SetAttribute{
                MarkdownDescription: "Users to add as owners on the OpenStack cloud when this rule matches. IDs of `oneuptime_user` records.",
                Computed: true,
                ElementType: types.StringType,
            },
            "owner_teams": schema.SetAttribute{
                MarkdownDescription: "Teams to add as owners on the OpenStack cloud when this rule matches. IDs of `oneuptime_team` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
                Validators: dataSourceRecordIDValidators("oneuptime_user"),
            },
        },
    }
    if d.isLegacyAlias {
        resp.Schema.DeprecationMessage = "oneuptime_open_stack_cloud_owner_rule has been renamed to oneuptime_openstack_cloud_owner_rule. Use the new name; the old one keeps working until then."
    }
}

func (d *OpenstackCloudOwnerRuleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OpenstackCloudOwnerRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data OpenstackCloudOwnerRuleDataSourceModel

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
    if !data.OpenStackCloudNamePattern.IsNull() && !data.OpenStackCloudNamePattern.IsUnknown() {
        filters["openStackCloudNamePattern"] = data.OpenStackCloudNamePattern.ValueString()
        filterNames = append(filterNames, "open_stack_cloud_name_pattern = "+fmt.Sprintf("%q", data.OpenStackCloudNamePattern.ValueString()))
    }
    if !data.OpenStackCloudDescriptionPattern.IsNull() && !data.OpenStackCloudDescriptionPattern.IsUnknown() {
        filters["openStackCloudDescriptionPattern"] = data.OpenStackCloudDescriptionPattern.ValueString()
        filterNames = append(filterNames, "open_stack_cloud_description_pattern = "+fmt.Sprintf("%q", data.OpenStackCloudDescriptionPattern.ValueString()))
    }
    if !data.OpenStackProjectId.IsNull() && !data.OpenStackProjectId.IsUnknown() {
        filters["openStackProjectId"] = data.OpenStackProjectId.ValueString()
        filterNames = append(filterNames, "open_stack_project_id = "+fmt.Sprintf("%q", data.OpenStackProjectId.ValueString()))
    }
    if !data.OpenStackDomainId.IsNull() && !data.OpenStackDomainId.IsUnknown() {
        filters["openStackDomainId"] = data.OpenStackDomainId.ValueString()
        filterNames = append(filterNames, "open_stack_domain_id = "+fmt.Sprintf("%q", data.OpenStackDomainId.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the openstack cloud owner rule up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the openstack cloud owner rule up by.",
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
        "openStackCloudLabels": true,
        "openStackCloudNamePattern": true,
        "openStackCloudDescriptionPattern": true,
        "openStackProjectId": true,
        "openStackDomainId": true,
        "ownerUsers": true,
        "ownerTeams": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/openstack-cloud-owner-rule/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read openstack_cloud_owner_rule, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No openstack cloud owner rule found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read openstack_cloud_owner_rule: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/openstack-cloud-owner-rule/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list openstack_cloud_owner_rule, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list openstack_cloud_owner_rule: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No openstack cloud owner rule matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one openstack cloud owner rule matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for openstack_cloud_owner_rule.")
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
    if val, ok := item["openStackCloudLabels"].([]interface{}); ok {
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
        data.OpenStackCloudLabels = types.SetValueMust(types.StringType, setItems)
    } else {
        data.OpenStackCloudLabels = types.SetNull(types.StringType)
    }
    if obj, ok := item["openStackCloudNamePattern"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OpenStackCloudNamePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OpenStackCloudNamePattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OpenStackCloudNamePattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OpenStackCloudNamePattern = types.StringValue(string(jsonBytes))
        } else {
            data.OpenStackCloudNamePattern = types.StringNull()
        }
    } else if val, ok := item["openStackCloudNamePattern"].(string); ok {
        data.OpenStackCloudNamePattern = types.StringValue(val)
    } else {
        data.OpenStackCloudNamePattern = types.StringNull()
    }
    if obj, ok := item["openStackCloudDescriptionPattern"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OpenStackCloudDescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OpenStackCloudDescriptionPattern = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OpenStackCloudDescriptionPattern = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OpenStackCloudDescriptionPattern = types.StringValue(string(jsonBytes))
        } else {
            data.OpenStackCloudDescriptionPattern = types.StringNull()
        }
    } else if val, ok := item["openStackCloudDescriptionPattern"].(string); ok {
        data.OpenStackCloudDescriptionPattern = types.StringValue(val)
    } else {
        data.OpenStackCloudDescriptionPattern = types.StringNull()
    }
    if obj, ok := item["openStackProjectId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OpenStackProjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OpenStackProjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OpenStackProjectId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OpenStackProjectId = types.StringValue(string(jsonBytes))
        } else {
            data.OpenStackProjectId = types.StringNull()
        }
    } else if val, ok := item["openStackProjectId"].(string); ok {
        data.OpenStackProjectId = types.StringValue(val)
    } else {
        data.OpenStackProjectId = types.StringNull()
    }
    if obj, ok := item["openStackDomainId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.OpenStackDomainId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.OpenStackDomainId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.OpenStackDomainId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.OpenStackDomainId = types.StringValue(string(jsonBytes))
        } else {
            data.OpenStackDomainId = types.StringNull()
        }
    } else if val, ok := item["openStackDomainId"].(string); ok {
        data.OpenStackDomainId = types.StringValue(val)
    } else {
        data.OpenStackDomainId = types.StringNull()
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
