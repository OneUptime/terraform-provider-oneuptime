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
        {Order: 44, New: NewSecurityPlaybookRunDataSource},
        {Order: 48, New: NewDataSecurityIdentityDataSource},
        {Order: 73, New: NewBiDashboardDataSource},
        {Order: 91, New: NewBlocklistFeedDataSource},
        {Order: 100, New: NewSecurityTelemetrySourceDataSource},
        {Order: 130, New: NewMonitorDataSource},
        {Order: 277, New: NewStatusPageAnnouncementDataSource},
        {Order: 304, New: NewAutoRemediationSuggestionDataSource},
        {Order: 309, New: NewRunnerTeamOwnerDataSource},
        {Order: 385, New: NewTeamOrLabelBudgetDataSource},
        {Order: 457, New: NewAddressNameDataSource},
        {Order: 463, New: NewProxmoxClusterDataSource},
        {Order: 488, New: NewCollectorConfigDataSource},
        {Order: 501, New: NewCloudResourceTeamOwnerDataSource},
        {Order: 507, New: NewDatabaseFeedDataSource},
    }, nil)
}
