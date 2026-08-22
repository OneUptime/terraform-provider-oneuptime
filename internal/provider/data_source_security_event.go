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
var _ datasource.DataSource = &SecurityEventDataSource{}

func NewSecurityEventDataSource() datasource.DataSource {
    return &SecurityEventDataSource{}
}

// SecurityEventDataSource defines the data source implementation.
type SecurityEventDataSource struct {
    client *Client
}

// SecurityEventDataSourceModel describes the data source data model.
type SecurityEventDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    ProjectId types.String `tfsdk:"project_id"`
    PrimaryEntityId types.String `tfsdk:"primary_entity_id"`
    PrimaryEntityType types.String `tfsdk:"primary_entity_type"`
    Time types.String `tfsdk:"time"`
    EventUid types.String `tfsdk:"event_uid"`
    CategoryUid types.Number `tfsdk:"category_uid"`
    CategoryName types.String `tfsdk:"category_name"`
    ClassUid types.Number `tfsdk:"class_uid"`
    ClassName types.String `tfsdk:"class_name"`
    ActivityName types.String `tfsdk:"activity_name"`
    SeverityId types.Number `tfsdk:"severity_id"`
    SeverityName types.String `tfsdk:"severity_name"`
    StatusName types.String `tfsdk:"status_name"`
    Message types.String `tfsdk:"message"`
    VendorName types.String `tfsdk:"vendor_name"`
    ProductName types.String `tfsdk:"product_name"`
    RuleId types.String `tfsdk:"rule_id"`
    RuleName types.String `tfsdk:"rule_name"`
    MitreTactics types.Set `tfsdk:"mitre_tactics"`
    MitreTechniques types.Set `tfsdk:"mitre_techniques"`
    PrincipalUser types.String `tfsdk:"principal_user"`
    PrincipalHost types.String `tfsdk:"principal_host"`
    PrincipalIp types.String `tfsdk:"principal_ip"`
    PrincipalProcess types.String `tfsdk:"principal_process"`
    TargetUser types.String `tfsdk:"target_user"`
    TargetHost types.String `tfsdk:"target_host"`
    TargetIp types.String `tfsdk:"target_ip"`
    TargetPort types.Number `tfsdk:"target_port"`
    TargetResource types.String `tfsdk:"target_resource"`
    Observables types.Set `tfsdk:"observables"`
    Attributes types.String `tfsdk:"attributes"`
    AttributeKeys types.Set `tfsdk:"attribute_keys"`
}

func (d *SecurityEventDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_security_event"
}

func (d *SecurityEventDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "API endpoints for Security Event Look up an existing security_event by `id` or by `name`.",

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
            "project_id": schema.StringAttribute{
                MarkdownDescription: "Project ID",
                Computed: true,
            },
            "primary_entity_id": schema.StringAttribute{
                MarkdownDescription: "Source ID",
                Computed: true,
            },
            "primary_entity_type": schema.StringAttribute{
                MarkdownDescription: "Source Type",
                Computed: true,
            },
            "time": schema.StringAttribute{
                MarkdownDescription: "Time",
                Computed: true,
            },
            "event_uid": schema.StringAttribute{
                MarkdownDescription: "Event UID",
                Computed: true,
            },
            "category_uid": schema.NumberAttribute{
                MarkdownDescription: "Category UID",
                Computed: true,
            },
            "category_name": schema.StringAttribute{
                MarkdownDescription: "Category",
                Computed: true,
            },
            "class_uid": schema.NumberAttribute{
                MarkdownDescription: "Class UID",
                Computed: true,
            },
            "class_name": schema.StringAttribute{
                MarkdownDescription: "Event Class",
                Computed: true,
            },
            "activity_name": schema.StringAttribute{
                MarkdownDescription: "Activity",
                Computed: true,
            },
            "severity_id": schema.NumberAttribute{
                MarkdownDescription: "Severity ID",
                Computed: true,
            },
            "severity_name": schema.StringAttribute{
                MarkdownDescription: "Severity",
                Computed: true,
            },
            "status_name": schema.StringAttribute{
                MarkdownDescription: "Status",
                Computed: true,
            },
            "message": schema.StringAttribute{
                MarkdownDescription: "Message",
                Computed: true,
            },
            "vendor_name": schema.StringAttribute{
                MarkdownDescription: "Vendor",
                Computed: true,
            },
            "product_name": schema.StringAttribute{
                MarkdownDescription: "Product",
                Computed: true,
            },
            "rule_id": schema.StringAttribute{
                MarkdownDescription: "Rule ID",
                Computed: true,
            },
            "rule_name": schema.StringAttribute{
                MarkdownDescription: "Rule Name",
                Computed: true,
            },
            "mitre_tactics": schema.SetAttribute{
                MarkdownDescription: "MITRE Tactics",
                Computed: true,
                ElementType: types.StringType,
            },
            "mitre_techniques": schema.SetAttribute{
                MarkdownDescription: "MITRE Techniques",
                Computed: true,
                ElementType: types.StringType,
            },
            "principal_user": schema.StringAttribute{
                MarkdownDescription: "Principal User",
                Computed: true,
            },
            "principal_host": schema.StringAttribute{
                MarkdownDescription: "Principal Host",
                Computed: true,
            },
            "principal_ip": schema.StringAttribute{
                MarkdownDescription: "Principal IP",
                Computed: true,
            },
            "principal_process": schema.StringAttribute{
                MarkdownDescription: "Principal Process",
                Computed: true,
            },
            "target_user": schema.StringAttribute{
                MarkdownDescription: "Target User",
                Computed: true,
            },
            "target_host": schema.StringAttribute{
                MarkdownDescription: "Target Host",
                Computed: true,
            },
            "target_ip": schema.StringAttribute{
                MarkdownDescription: "Target IP",
                Computed: true,
            },
            "target_port": schema.NumberAttribute{
                MarkdownDescription: "Target Port",
                Computed: true,
            },
            "target_resource": schema.StringAttribute{
                MarkdownDescription: "Target Resource",
                Computed: true,
            },
            "observables": schema.SetAttribute{
                MarkdownDescription: "Observables",
                Computed: true,
                ElementType: types.StringType,
            },
            "attributes": schema.StringAttribute{
                MarkdownDescription: "Attributes",
                Computed: true,
            },
            "attribute_keys": schema.SetAttribute{
                MarkdownDescription: "Attribute Keys",
                Computed: true,
                ElementType: types.StringType,
            },
        },
    }
}

func (d *SecurityEventDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SecurityEventDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data SecurityEventDataSourceModel

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
            "Exactly one of `id` or `name` must be set to look up a security_event.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "name": true,
        "projectId": true,
        "primaryEntityId": true,
        "primaryEntityType": true,
        "time": true,
        "eventUid": true,
        "categoryUid": true,
        "categoryName": true,
        "classUid": true,
        "className": true,
        "activityName": true,
        "severityId": true,
        "severityName": true,
        "statusName": true,
        "message": true,
        "vendorName": true,
        "productName": true,
        "ruleId": true,
        "ruleName": true,
        "mitreTactics": true,
        "mitreTechniques": true,
        "principalUser": true,
        "principalHost": true,
        "principalIp": true,
        "principalProcess": true,
        "targetUser": true,
        "targetHost": true,
        "targetIp": true,
        "targetPort": true,
        "targetResource": true,
        "observables": true,
        "attributes": true,
        "attributeKeys": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/security-events/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read security_event, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No security_event found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read security_event: %s", err))
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
        httpResp, err := d.client.Post(ctx, "/security-events/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list security_event, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list security_event: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No security_event found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one security_event matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for security_event.")
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
    if obj, ok := item["time"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Time = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Time = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Time = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Time = types.StringValue(string(jsonBytes))
        } else {
            data.Time = types.StringNull()
        }
    } else if val, ok := item["time"].(string); ok {
        data.Time = types.StringValue(val)
    } else {
        data.Time = types.StringNull()
    }
    if obj, ok := item["eventUid"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.EventUid = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.EventUid = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.EventUid = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.EventUid = types.StringValue(string(jsonBytes))
        } else {
            data.EventUid = types.StringNull()
        }
    } else if val, ok := item["eventUid"].(string); ok {
        data.EventUid = types.StringValue(val)
    } else {
        data.EventUid = types.StringNull()
    }
    if val, ok := item["categoryUid"].(float64); ok {
        data.CategoryUid = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["categoryUid"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.CategoryUid = types.NumberValue(big.NewFloat(val))
        } else {
            data.CategoryUid = types.NumberNull()
        }
    } else {
        data.CategoryUid = types.NumberNull()
    }
    if obj, ok := item["categoryName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CategoryName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CategoryName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CategoryName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CategoryName = types.StringValue(string(jsonBytes))
        } else {
            data.CategoryName = types.StringNull()
        }
    } else if val, ok := item["categoryName"].(string); ok {
        data.CategoryName = types.StringValue(val)
    } else {
        data.CategoryName = types.StringNull()
    }
    if val, ok := item["classUid"].(float64); ok {
        data.ClassUid = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["classUid"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.ClassUid = types.NumberValue(big.NewFloat(val))
        } else {
            data.ClassUid = types.NumberNull()
        }
    } else {
        data.ClassUid = types.NumberNull()
    }
    if obj, ok := item["className"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ClassName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ClassName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ClassName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ClassName = types.StringValue(string(jsonBytes))
        } else {
            data.ClassName = types.StringNull()
        }
    } else if val, ok := item["className"].(string); ok {
        data.ClassName = types.StringValue(val)
    } else {
        data.ClassName = types.StringNull()
    }
    if obj, ok := item["activityName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ActivityName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ActivityName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ActivityName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ActivityName = types.StringValue(string(jsonBytes))
        } else {
            data.ActivityName = types.StringNull()
        }
    } else if val, ok := item["activityName"].(string); ok {
        data.ActivityName = types.StringValue(val)
    } else {
        data.ActivityName = types.StringNull()
    }
    if val, ok := item["severityId"].(float64); ok {
        data.SeverityId = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["severityId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.SeverityId = types.NumberValue(big.NewFloat(val))
        } else {
            data.SeverityId = types.NumberNull()
        }
    } else {
        data.SeverityId = types.NumberNull()
    }
    if obj, ok := item["severityName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SeverityName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SeverityName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SeverityName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SeverityName = types.StringValue(string(jsonBytes))
        } else {
            data.SeverityName = types.StringNull()
        }
    } else if val, ok := item["severityName"].(string); ok {
        data.SeverityName = types.StringValue(val)
    } else {
        data.SeverityName = types.StringNull()
    }
    if obj, ok := item["statusName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StatusName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StatusName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StatusName = types.StringValue(string(jsonBytes))
        } else {
            data.StatusName = types.StringNull()
        }
    } else if val, ok := item["statusName"].(string); ok {
        data.StatusName = types.StringValue(val)
    } else {
        data.StatusName = types.StringNull()
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
    if obj, ok := item["vendorName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.VendorName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.VendorName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.VendorName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.VendorName = types.StringValue(string(jsonBytes))
        } else {
            data.VendorName = types.StringNull()
        }
    } else if val, ok := item["vendorName"].(string); ok {
        data.VendorName = types.StringValue(val)
    } else {
        data.VendorName = types.StringNull()
    }
    if obj, ok := item["productName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProductName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProductName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProductName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProductName = types.StringValue(string(jsonBytes))
        } else {
            data.ProductName = types.StringNull()
        }
    } else if val, ok := item["productName"].(string); ok {
        data.ProductName = types.StringValue(val)
    } else {
        data.ProductName = types.StringNull()
    }
    if obj, ok := item["ruleId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RuleId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RuleId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RuleId = types.StringValue(string(jsonBytes))
        } else {
            data.RuleId = types.StringNull()
        }
    } else if val, ok := item["ruleId"].(string); ok {
        data.RuleId = types.StringValue(val)
    } else {
        data.RuleId = types.StringNull()
    }
    if obj, ok := item["ruleName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RuleName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RuleName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RuleName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RuleName = types.StringValue(string(jsonBytes))
        } else {
            data.RuleName = types.StringNull()
        }
    } else if val, ok := item["ruleName"].(string); ok {
        data.RuleName = types.StringValue(val)
    } else {
        data.RuleName = types.StringNull()
    }
    if val, ok := item["mitreTactics"].([]interface{}); ok {
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
        data.MitreTactics = types.SetValueMust(types.StringType, setItems)
    } else {
        data.MitreTactics = types.SetNull(types.StringType)
    }
    if val, ok := item["mitreTechniques"].([]interface{}); ok {
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
        data.MitreTechniques = types.SetValueMust(types.StringType, setItems)
    } else {
        data.MitreTechniques = types.SetNull(types.StringType)
    }
    if obj, ok := item["principalUser"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PrincipalUser = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PrincipalUser = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PrincipalUser = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PrincipalUser = types.StringValue(string(jsonBytes))
        } else {
            data.PrincipalUser = types.StringNull()
        }
    } else if val, ok := item["principalUser"].(string); ok {
        data.PrincipalUser = types.StringValue(val)
    } else {
        data.PrincipalUser = types.StringNull()
    }
    if obj, ok := item["principalHost"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PrincipalHost = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PrincipalHost = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PrincipalHost = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PrincipalHost = types.StringValue(string(jsonBytes))
        } else {
            data.PrincipalHost = types.StringNull()
        }
    } else if val, ok := item["principalHost"].(string); ok {
        data.PrincipalHost = types.StringValue(val)
    } else {
        data.PrincipalHost = types.StringNull()
    }
    if obj, ok := item["principalIp"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PrincipalIp = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PrincipalIp = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PrincipalIp = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PrincipalIp = types.StringValue(string(jsonBytes))
        } else {
            data.PrincipalIp = types.StringNull()
        }
    } else if val, ok := item["principalIp"].(string); ok {
        data.PrincipalIp = types.StringValue(val)
    } else {
        data.PrincipalIp = types.StringNull()
    }
    if obj, ok := item["principalProcess"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PrincipalProcess = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PrincipalProcess = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PrincipalProcess = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PrincipalProcess = types.StringValue(string(jsonBytes))
        } else {
            data.PrincipalProcess = types.StringNull()
        }
    } else if val, ok := item["principalProcess"].(string); ok {
        data.PrincipalProcess = types.StringValue(val)
    } else {
        data.PrincipalProcess = types.StringNull()
    }
    if obj, ok := item["targetUser"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TargetUser = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TargetUser = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TargetUser = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TargetUser = types.StringValue(string(jsonBytes))
        } else {
            data.TargetUser = types.StringNull()
        }
    } else if val, ok := item["targetUser"].(string); ok {
        data.TargetUser = types.StringValue(val)
    } else {
        data.TargetUser = types.StringNull()
    }
    if obj, ok := item["targetHost"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TargetHost = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TargetHost = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TargetHost = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TargetHost = types.StringValue(string(jsonBytes))
        } else {
            data.TargetHost = types.StringNull()
        }
    } else if val, ok := item["targetHost"].(string); ok {
        data.TargetHost = types.StringValue(val)
    } else {
        data.TargetHost = types.StringNull()
    }
    if obj, ok := item["targetIp"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TargetIp = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TargetIp = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TargetIp = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TargetIp = types.StringValue(string(jsonBytes))
        } else {
            data.TargetIp = types.StringNull()
        }
    } else if val, ok := item["targetIp"].(string); ok {
        data.TargetIp = types.StringValue(val)
    } else {
        data.TargetIp = types.StringNull()
    }
    if val, ok := item["targetPort"].(float64); ok {
        data.TargetPort = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["targetPort"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.TargetPort = types.NumberValue(big.NewFloat(val))
        } else {
            data.TargetPort = types.NumberNull()
        }
    } else {
        data.TargetPort = types.NumberNull()
    }
    if obj, ok := item["targetResource"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TargetResource = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TargetResource = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TargetResource = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TargetResource = types.StringValue(string(jsonBytes))
        } else {
            data.TargetResource = types.StringNull()
        }
    } else if val, ok := item["targetResource"].(string); ok {
        data.TargetResource = types.StringValue(val)
    } else {
        data.TargetResource = types.StringNull()
    }
    if val, ok := item["observables"].([]interface{}); ok {
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
        data.Observables = types.SetValueMust(types.StringType, setItems)
    } else {
        data.Observables = types.SetNull(types.StringType)
    }
    if obj, ok := item["attributes"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Attributes = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Attributes = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Attributes = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Attributes = types.StringValue(string(jsonBytes))
        } else {
            data.Attributes = types.StringNull()
        }
    } else if val, ok := item["attributes"].(string); ok {
        data.Attributes = types.StringValue(val)
    } else {
        data.Attributes = types.StringNull()
    }
    if val, ok := item["attributeKeys"].([]interface{}); ok {
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
        data.AttributeKeys = types.SetValueMust(types.StringType, setItems)
    } else {
        data.AttributeKeys = types.SetNull(types.StringType)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
