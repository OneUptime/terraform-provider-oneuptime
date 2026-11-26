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
        {Order: 54, New: NewDataFlowDataSource},
        {Order: 101, New: NewLogScrubRuleDataSource},
        {Order: 133, New: NewMonitorStatusDataSource},
        {Order: 297, New: NewWorkflowUserOwnerDataSource},
        {Order: 306, New: NewPostmortemRuleDataSource},
        {Order: 415, New: NewPodmanHostFeedDataSource},
        {Order: 495, New: NewCloudAccountDataSource},
        {Order: 513, New: NewQueueTeamOwnerDataSource},
    }, nil)
}
