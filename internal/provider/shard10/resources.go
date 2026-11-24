package shard10

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The resources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterResources("shard10", []OrderedResource{
        {Order: 14, New: NewMetricSavedViewResource},
        {Order: 26, New: NewThreatIntelFeedResource},
        {Order: 47, New: NewDataSecuritySettingResource},
        {Order: 53, New: NewDataReconciliationResource},
        {Order: 140, New: NewAlertEpisodeTeamOwnerResource},
        {Order: 146, New: NewAlertIncidentRuleResource},
        {Order: 201, New: NewStatusPageMonitorRuleResource},
        {Order: 209, New: NewPodmanHostLabelRuleResource},
        {Order: 435, New: NewCloudResourceUserOwnerResource},
        {Order: 444, New: NewDatabaseOwnerRuleResource},
    }, nil)
}
