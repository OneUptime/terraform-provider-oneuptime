package shard05

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard05", []OrderedDataSource{
        {Order: 37, New: NewSecurityCaseNoteTemplateDataSource},
        {Order: 58, New: NewResidencyRuleDataSource},
        {Order: 63, New: NewDataCheckDataSource},
        {Order: 94, New: NewLlmModelPriceDataSource},
        {Order: 105, New: NewSmtpConfigDataSource},
        {Order: 113, New: NewOnCallDutyExecutionLogTimelineDataSource},
        {Order: 128, New: NewIncidentFeedDataSource},
        {Order: 131, New: NewIncidentInternalNoteDataSource},
        {Order: 138, New: NewIncidentRoleDataSource},
        {Order: 146, New: NewFormSubmissionDataSource},
        {Order: 154, New: NewAlertCustomFieldDataSource},
        {Order: 219, New: NewSloOwnerRuleDataSource},
        {Order: 257, New: NewWorkflowLabelRuleDataSource},
        {Order: 307, New: NewMonitorProbeDataSource},
        {Order: 342, New: NewOnCallPolicyLabelRuleDataSource},
        {Order: 349, New: NewIncomingCallPolicyOwnerRuleDataSource},
        {Order: 355, New: NewServiceUserOwnerDataSource},
        {Order: 402, New: NewDockerSwarmClusterFeedDataSource},
        {Order: 407, New: NewOpenstackCloudFeedDataSource},
        {Order: 449, New: NewDockerHostUserOwnerDataSource},
        {Order: 462, New: NewIotFleetDataSource},
        {Order: 472, New: NewVcenterDataSource},
        {Order: 519, New: NewNetworkFlowDataSource},
        {Order: 542, New: NewOpenstackCloudFeedLegacyDataSource},
        {Order: 547, New: NewIotFleetLegacyDataSource},
        {Order: 551, New: NewVcenterLegacyDataSource},
    }, []LegacyDataSourceAlias{
        {Order: 407, Name: "openstack_cloud_feed", LegacyName: "open_stack_cloud_feed", New: NewOpenstackCloudFeedDataSource, NewLegacy: NewOpenstackCloudFeedLegacyDataSource},
        {Order: 462, Name: "iot_fleet", LegacyName: "io_t_fleet", New: NewIotFleetDataSource, NewLegacy: NewIotFleetLegacyDataSource},
        {Order: 472, Name: "vcenter", LegacyName: "v_center", New: NewVcenterDataSource, NewLegacy: NewVcenterLegacyDataSource},
    })
}
