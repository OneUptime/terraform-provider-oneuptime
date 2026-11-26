package shard14

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard14", []OrderedDataSource{
        {Order: 14, New: NewLogSavedViewDataSource},
        {Order: 21, New: NewDetectionRuleDataSource},
        {Order: 51, New: NewDataStoreTeamOwnerDataSource},
        {Order: 81, New: NewPipelineTeamOwnerDataSource},
        {Order: 93, New: NewAutomaticResponsePolicyDataSource},
        {Order: 140, New: NewIncidentStateTimelineDataSource},
        {Order: 167, New: NewAlertTeamOwnerDataSource},
        {Order: 220, New: NewScheduledMaintenanceMeasurementValueDataSource},
        {Order: 242, New: NewServiceOwnerRuleDataSource},
        {Order: 244, New: NewDockerHostOwnerRuleDataSource},
        {Order: 285, New: NewStatusPageHistoryChartBarColorDataSource},
        {Order: 293, New: NewWorkflowVariableDataSource},
        {Order: 308, New: NewRunnerDataSource},
        {Order: 371, New: NewAiAgentUserOwnerDataSource},
        {Order: 462, New: NewPodmanHostUserOwnerDataSource},
        {Order: 500, New: NewCloudResourceDataSource},
    }, nil)
}
