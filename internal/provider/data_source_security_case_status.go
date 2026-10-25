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
var _ datasource.DataSource = &SecurityCaseStatusDataSource{}

func NewSecurityCaseStatusDataSource() datasource.DataSource {
    return &SecurityCaseStatusDataSource{}
}

// SecurityCaseStatusDataSource defines the data source implementation.
type SecurityCaseStatusDataSource struct {
    client *Client
}

// SecurityCaseStatusDataSourceModel describes the data source data model.
type SecurityCaseStatusDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    Name types.String `tfsdk:"name"`
    Slug types.String `tfsdk:"slug"`
    Description types.String `tfsdk:"description"`
    Color types.String `tfsdk:"color"`
    IsNewStatus types.Bool `tfsdk:"is_new_status"`
    IsAcknowledgedStatus types.Bool `tfsdk:"is_acknowledged_status"`
    IsContainedStatus types.Bool `tfsdk:"is_contained_status"`
    IsClosedStatus types.Bool `tfsdk:"is_closed_status"`
    Order types.Number `tfsdk:"order"`
    IsArchived types.Bool `tfsdk:"is_archived"`
    ArchivedAt types.String `tfsdk:"archived_at"`
    ArchivedByUserId types.String `tfsdk:"archived_by_user_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (d *SecurityCaseStatusDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_security_case_status"
}

func (d *SecurityCaseStatusDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "The statuses a security case moves through (New, Investigating, Contained, Closed and your own). Add, rename, reorder or archive them. Look up an existing security case status by `id`, or by any of its other arguments (`name`, `archived_by_user_id`, `created_by_user_id`, ...): each one set must match, and exactly one security case status may match them all.",

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
                MarkdownDescription: "ID of the project this status belongs to. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "What the status is called, like Investigating or Triage.",
                Optional: true,
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for the status.",
                Optional: true,
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "What a case in this status is waiting for, in a sentence.",
                Optional: true,
                Computed: true,
            },
            "color": schema.StringAttribute{
                MarkdownDescription: "The status's colour, in hex (#32a852 for example).",
                Computed: true,
            },
            "is_new_status": schema.BoolAttribute{
                MarkdownDescription: "Is it the status every new case starts in? A project has one, and it can't be deleted or archived.",
                Optional: true,
                Computed: true,
            },
            "is_acknowledged_status": schema.BoolAttribute{
                MarkdownDescription: "Is it the status a reopened case moves to? From it down, a case counts as taken on. A project has one, and it can't be deleted or archived.",
                Optional: true,
                Computed: true,
            },
            "is_contained_status": schema.BoolAttribute{
                MarkdownDescription: "Is it the status the time to contain is measured to? From it down, a case counts as contained. A project has one.",
                Optional: true,
                Computed: true,
            },
            "is_closed_status": schema.BoolAttribute{
                MarkdownDescription: "Is it the status Close Case moves a case to? From it down, a case counts as closed and carries a verdict. A project has one.",
                Optional: true,
                Computed: true,
            },
            "order": schema.NumberAttribute{
                MarkdownDescription: "Where this status sits in the project's list of security case statuses: 1 is the top. A case in a status at or below the acknowledged (contained, closed) status counts as acknowledged (contained, closed), so the new, acknowledged, contained and closed statuses have to stay in that order. A new status without a number goes just above the closed status. Setting a number moves the status to that place, and the ones in between shift by one. In the dashboard, drag the rows to reorder them.",
                Optional: true,
                Computed: true,
            },
            "is_archived": schema.BoolAttribute{
                MarkdownDescription: "Archived statuses stay on the cases in them, but can't be picked again. The built-in statuses can't be archived.",
                Optional: true,
                Computed: true,
            },
            "archived_at": schema.StringAttribute{
                MarkdownDescription: "When this status was archived. Empty while it is not archived.",
                Computed: true,
            },
            "archived_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who archived this status. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "ID of the person who added this status. The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *SecurityCaseStatusDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SecurityCaseStatusDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data SecurityCaseStatusDataSourceModel

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
    if !data.IsNewStatus.IsNull() && !data.IsNewStatus.IsUnknown() {
        filters["isNewStatus"] = data.IsNewStatus.ValueBool()
        filterNames = append(filterNames, "is_new_status = "+fmt.Sprintf("%t", data.IsNewStatus.ValueBool()))
    }
    if !data.IsAcknowledgedStatus.IsNull() && !data.IsAcknowledgedStatus.IsUnknown() {
        filters["isAcknowledgedStatus"] = data.IsAcknowledgedStatus.ValueBool()
        filterNames = append(filterNames, "is_acknowledged_status = "+fmt.Sprintf("%t", data.IsAcknowledgedStatus.ValueBool()))
    }
    if !data.IsContainedStatus.IsNull() && !data.IsContainedStatus.IsUnknown() {
        filters["isContainedStatus"] = data.IsContainedStatus.ValueBool()
        filterNames = append(filterNames, "is_contained_status = "+fmt.Sprintf("%t", data.IsContainedStatus.ValueBool()))
    }
    if !data.IsClosedStatus.IsNull() && !data.IsClosedStatus.IsUnknown() {
        filters["isClosedStatus"] = data.IsClosedStatus.ValueBool()
        filterNames = append(filterNames, "is_closed_status = "+fmt.Sprintf("%t", data.IsClosedStatus.ValueBool()))
    }
    if !data.Order.IsNull() && !data.Order.IsUnknown() {
        filters["order"] = lookupNumber(data.Order)
        filterNames = append(filterNames, "order = "+data.Order.ValueBigFloat().String())
    }
    if !data.IsArchived.IsNull() && !data.IsArchived.IsUnknown() {
        filters["isArchived"] = data.IsArchived.ValueBool()
        filterNames = append(filterNames, "is_archived = "+fmt.Sprintf("%t", data.IsArchived.ValueBool()))
    }
    if !data.ArchivedByUserId.IsNull() && !data.ArchivedByUserId.IsUnknown() {
        filters["archivedByUserId"] = data.ArchivedByUserId.ValueString()
        filterNames = append(filterNames, "archived_by_user_id = "+fmt.Sprintf("%q", data.ArchivedByUserId.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the security case status up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the security case status up by.",
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
        "color": true,
        "isNewStatus": true,
        "isAcknowledgedStatus": true,
        "isContainedStatus": true,
        "isClosedStatus": true,
        "order": true,
        "isArchived": true,
        "archivedAt": true,
        "archivedByUserId": true,
        "createdByUserId": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/security-case-status/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read security_case_status, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No security case status found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read security_case_status: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/security-case-status/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list security_case_status, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list security_case_status: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No security case status matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one security case status matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for security_case_status.")
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
    if obj, ok := item["color"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Color = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Color = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Color = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Color = types.StringValue(string(jsonBytes))
        } else {
            data.Color = types.StringNull()
        }
    } else if val, ok := item["color"].(string); ok {
        data.Color = types.StringValue(val)
    } else {
        data.Color = types.StringNull()
    }
    if val, ok := item["isNewStatus"].(bool); ok {
        data.IsNewStatus = types.BoolValue(val)
    } else {
        data.IsNewStatus = types.BoolNull()
    }
    if val, ok := item["isAcknowledgedStatus"].(bool); ok {
        data.IsAcknowledgedStatus = types.BoolValue(val)
    } else {
        data.IsAcknowledgedStatus = types.BoolNull()
    }
    if val, ok := item["isContainedStatus"].(bool); ok {
        data.IsContainedStatus = types.BoolValue(val)
    } else {
        data.IsContainedStatus = types.BoolNull()
    }
    if val, ok := item["isClosedStatus"].(bool); ok {
        data.IsClosedStatus = types.BoolValue(val)
    } else {
        data.IsClosedStatus = types.BoolNull()
    }
    data.Order = numberFromAPI(data.Order, item["order"])
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
