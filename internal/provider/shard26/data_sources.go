package shard26

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard26", []OrderedDataSource{
        {Order: 4, New: NewTeamCustomFieldDataSource},
        {Order: 41, New: NewSecurityPlaybookRunDataSource},
        {Order: 45, New: NewDataSecurityIdentityDataSource},
        {Order: 68, New: NewBiDashboardDataSource},
        {Order: 90, New: NewSecurityTelemetrySourceDataSource},
        {Order: 120, New: NewMonitorDataSource},
        {Order: 267, New: NewStatusPageAnnouncementDataSource},
        {Order: 294, New: NewAutoRemediationSuggestionDataSource},
        {Order: 299, New: NewRunnerTeamOwnerDataSource},
        {Order: 375, New: NewTeamOrLabelBudgetDataSource},
        {Order: 447, New: NewAddressNameDataSource},
        {Order: 453, New: NewProxmoxClusterDataSource},
        {Order: 478, New: NewCollectorConfigDataSource},
        {Order: 491, New: NewCloudResourceTeamOwnerDataSource},
        {Order: 497, New: NewDatabaseFeedDataSource},
    }, nil)
}
