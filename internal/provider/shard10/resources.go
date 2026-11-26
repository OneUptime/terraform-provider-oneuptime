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
        {Order: 48, New: NewDataSecuritySettingResource},
        {Order: 54, New: NewDataReconciliationResource},
        {Order: 144, New: NewAlertEpisodeTeamOwnerResource},
        {Order: 150, New: NewAlertIncidentRuleResource},
        {Order: 205, New: NewStatusPageMonitorRuleResource},
        {Order: 213, New: NewPodmanHostLabelRuleResource},
        {Order: 439, New: NewCloudResourceUserOwnerResource},
        {Order: 448, New: NewDatabaseOwnerRuleResource},
    }, nil)
}
