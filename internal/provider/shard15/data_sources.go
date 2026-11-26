package shard15

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard15", []OrderedDataSource{
        {Order: 6, New: NewTeamMemberCustomFieldDataSource},
        {Order: 60, New: NewCustomDataTypeDataSource},
        {Order: 88, New: NewBlockDataSource},
        {Order: 99, New: NewSecurityLogSourceDataSource},
        {Order: 135, New: NewMonitorTemplateDataSource},
        {Order: 145, New: NewIncidentTemplateUserOwnerDataSource},
        {Order: 217, New: NewAlertMeasurementDataSource},
        {Order: 240, New: NewHostOwnerRuleDataSource},
        {Order: 261, New: NewStorageArrayLabelRuleDataSource},
        {Order: 275, New: NewStatusPageCustomFieldDataSource},
        {Order: 278, New: NewStatusPageAnnouncementTemplateDataSource},
        {Order: 330, New: NewCallLogDataSource},
        {Order: 377, New: NewAiInsightDataSource},
        {Order: 427, New: NewKubernetesClusterDataSource},
        {Order: 476, New: NewCephClusterDataSource},
        {Order: 480, New: NewStorageArrayTeamOwnerDataSource},
        {Order: 491, New: NewServerlessFunctionTeamOwnerDataSource},
        {Order: 493, New: NewServerlessFunctionLabelRuleDataSource},
        {Order: 537, New: NewThreatIntelIndicatorDataSource},
    }, nil)
}
