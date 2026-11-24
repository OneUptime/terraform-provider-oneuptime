package shard27

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard27", []OrderedDataSource{
        {Order: 2, New: NewProjectDataSource},
        {Order: 33, New: NewSecurityCaseTaskDataSource},
        {Order: 51, New: NewDataFlowDataSource},
        {Order: 91, New: NewLogScrubRuleDataSource},
        {Order: 123, New: NewMonitorStatusDataSource},
        {Order: 287, New: NewWorkflowUserOwnerDataSource},
        {Order: 296, New: NewPostmortemRuleDataSource},
        {Order: 405, New: NewPodmanHostFeedDataSource},
        {Order: 485, New: NewCloudAccountDataSource},
        {Order: 503, New: NewQueueTeamOwnerDataSource},
    }, nil)
}
