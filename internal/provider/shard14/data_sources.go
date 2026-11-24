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
        {Order: 48, New: NewDataStoreTeamOwnerDataSource},
        {Order: 76, New: NewPipelineTeamOwnerDataSource},
        {Order: 130, New: NewIncidentStateTimelineDataSource},
        {Order: 157, New: NewAlertTeamOwnerDataSource},
        {Order: 210, New: NewScheduledMaintenanceMeasurementValueDataSource},
        {Order: 232, New: NewServiceOwnerRuleDataSource},
        {Order: 234, New: NewDockerHostOwnerRuleDataSource},
        {Order: 275, New: NewStatusPageHistoryChartBarColorDataSource},
        {Order: 283, New: NewWorkflowVariableDataSource},
        {Order: 298, New: NewRunnerDataSource},
        {Order: 361, New: NewAiAgentUserOwnerDataSource},
        {Order: 452, New: NewPodmanHostUserOwnerDataSource},
        {Order: 490, New: NewCloudResourceDataSource},
    }, nil)
}
