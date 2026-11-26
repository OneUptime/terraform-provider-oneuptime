package shard04

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard04", []OrderedDataSource{
        {Order: 46, New: NewDataStoreFindingDataSource},
        {Order: 183, New: NewAlertPrivacyRuleDataSource},
        {Order: 247, New: NewPodmanHostLabelRuleDataSource},
        {Order: 273, New: NewStatusPageGroupDataSource},
        {Order: 305, New: NewAutoRemediationDecisionDataSource},
        {Order: 310, New: NewRunnerUserOwnerDataSource},
        {Order: 369, New: NewAiAgentDataSource},
        {Order: 414, New: NewStorageArrayFeedDataSource},
        {Order: 434, New: NewNetworkDeviceOwnerRuleDataSource},
        {Order: 460, New: NewPodmanHostDataSource},
        {Order: 502, New: NewCloudResourceUserOwnerDataSource},
        {Order: 521, New: NewRumApplicationOwnerRuleDataSource},
    }, nil)
}
