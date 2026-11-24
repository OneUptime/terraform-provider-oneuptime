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
        {Order: 55, New: NewCustomDataTypeDataSource},
        {Order: 83, New: NewBlockDataSource},
        {Order: 89, New: NewSecurityLogSourceDataSource},
        {Order: 125, New: NewMonitorTemplateDataSource},
        {Order: 135, New: NewIncidentTemplateUserOwnerDataSource},
        {Order: 207, New: NewAlertMeasurementDataSource},
        {Order: 230, New: NewHostOwnerRuleDataSource},
        {Order: 251, New: NewStorageArrayLabelRuleDataSource},
        {Order: 265, New: NewStatusPageCustomFieldDataSource},
        {Order: 268, New: NewStatusPageAnnouncementTemplateDataSource},
        {Order: 320, New: NewCallLogDataSource},
        {Order: 367, New: NewAiInsightDataSource},
        {Order: 417, New: NewKubernetesClusterDataSource},
        {Order: 466, New: NewCephClusterDataSource},
        {Order: 470, New: NewStorageArrayTeamOwnerDataSource},
        {Order: 481, New: NewServerlessFunctionTeamOwnerDataSource},
        {Order: 483, New: NewServerlessFunctionLabelRuleDataSource},
        {Order: 527, New: NewThreatIntelIndicatorDataSource},
    }, nil)
}
