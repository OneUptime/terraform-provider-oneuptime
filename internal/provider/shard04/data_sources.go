package shard04

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard04", []OrderedDataSource{
        {Order: 43, New: NewDataStoreFindingDataSource},
        {Order: 173, New: NewAlertPrivacyRuleDataSource},
        {Order: 237, New: NewPodmanHostLabelRuleDataSource},
        {Order: 263, New: NewStatusPageGroupDataSource},
        {Order: 295, New: NewAutoRemediationDecisionDataSource},
        {Order: 300, New: NewRunnerUserOwnerDataSource},
        {Order: 359, New: NewAiAgentDataSource},
        {Order: 404, New: NewStorageArrayFeedDataSource},
        {Order: 424, New: NewNetworkDeviceOwnerRuleDataSource},
        {Order: 450, New: NewPodmanHostDataSource},
        {Order: 492, New: NewCloudResourceUserOwnerDataSource},
        {Order: 511, New: NewRumApplicationOwnerRuleDataSource},
    }, nil)
}
