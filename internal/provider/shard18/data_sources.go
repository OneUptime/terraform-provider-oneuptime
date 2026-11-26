package shard18

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard18", []OrderedDataSource{
        {Order: 13, New: NewLabelDataSource},
        {Order: 45, New: NewDataStoreDataSource},
        {Order: 102, New: NewMetricPipelineRuleDataSource},
        {Order: 155, New: NewFormDataSource},
        {Order: 162, New: NewAlertDataSource},
        {Order: 170, New: NewAlertNoteTemplateDataSource},
        {Order: 174, New: NewAlertEpisodeUserOwnerDataSource},
        {Order: 214, New: NewIncidentSlaRuleDataSource},
        {Order: 307, New: NewExternalMcpServerDataSource},
        {Order: 449, New: NewNetworkEndpointDataSource},
        {Order: 486, New: NewHostTeamOwnerDataSource},
        {Order: 503, New: NewCloudResourceLabelRuleDataSource},
        {Order: 516, New: NewQueueOwnerRuleDataSource},
    }, nil)
}
