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
        {Order: 42, New: NewDataStoreDataSource},
        {Order: 92, New: NewMetricPipelineRuleDataSource},
        {Order: 145, New: NewFormDataSource},
        {Order: 152, New: NewAlertDataSource},
        {Order: 160, New: NewAlertNoteTemplateDataSource},
        {Order: 164, New: NewAlertEpisodeUserOwnerDataSource},
        {Order: 204, New: NewIncidentSlaRuleDataSource},
        {Order: 297, New: NewExternalMcpServerDataSource},
        {Order: 439, New: NewNetworkEndpointDataSource},
        {Order: 476, New: NewHostTeamOwnerDataSource},
        {Order: 493, New: NewCloudResourceLabelRuleDataSource},
        {Order: 506, New: NewQueueOwnerRuleDataSource},
    }, nil)
}
