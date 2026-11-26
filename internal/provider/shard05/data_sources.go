package shard05

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard05", []OrderedDataSource{
        {Order: 40, New: NewSecurityCaseNoteTemplateDataSource},
        {Order: 63, New: NewResidencyRuleDataSource},
        {Order: 68, New: NewDataCheckDataSource},
        {Order: 104, New: NewLlmModelPriceDataSource},
        {Order: 115, New: NewSmtpConfigDataSource},
        {Order: 123, New: NewOnCallDutyExecutionLogTimelineDataSource},
        {Order: 138, New: NewIncidentFeedDataSource},
        {Order: 141, New: NewIncidentInternalNoteDataSource},
        {Order: 148, New: NewIncidentRoleDataSource},
        {Order: 156, New: NewFormSubmissionDataSource},
        {Order: 164, New: NewAlertCustomFieldDataSource},
        {Order: 229, New: NewSloOwnerRuleDataSource},
        {Order: 267, New: NewWorkflowLabelRuleDataSource},
        {Order: 317, New: NewMonitorProbeDataSource},
        {Order: 352, New: NewOnCallPolicyLabelRuleDataSource},
        {Order: 359, New: NewIncomingCallPolicyOwnerRuleDataSource},
        {Order: 365, New: NewServiceUserOwnerDataSource},
        {Order: 412, New: NewDockerSwarmClusterFeedDataSource},
        {Order: 417, New: NewOpenstackCloudFeedDataSource},
        {Order: 459, New: NewDockerHostUserOwnerDataSource},
        {Order: 472, New: NewIotFleetDataSource},
        {Order: 482, New: NewVcenterDataSource},
        {Order: 529, New: NewNetworkFlowDataSource},
        {Order: 552, New: NewOpenstackCloudFeedLegacyDataSource},
        {Order: 557, New: NewIotFleetLegacyDataSource},
        {Order: 561, New: NewVcenterLegacyDataSource},
    }, []LegacyDataSourceAlias{
        {Order: 417, Name: "openstack_cloud_feed", LegacyName: "open_stack_cloud_feed", New: NewOpenstackCloudFeedDataSource, NewLegacy: NewOpenstackCloudFeedLegacyDataSource},
        {Order: 472, Name: "iot_fleet", LegacyName: "io_t_fleet", New: NewIotFleetDataSource, NewLegacy: NewIotFleetLegacyDataSource},
        {Order: 482, Name: "vcenter", LegacyName: "v_center", New: NewVcenterDataSource, NewLegacy: NewVcenterLegacyDataSource},
    })
}
