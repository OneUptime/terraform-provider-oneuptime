package shard00

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard00", []OrderedDataSource{
        {Order: 5, New: NewTeamMemberDataSource},
        {Order: 50, New: NewDataStoreUserOwnerDataSource},
        {Order: 71, New: NewDataProductDataSource},
        {Order: 80, New: NewPipelineUserOwnerDataSource},
        {Order: 89, New: NewBlockEnforcementDataSource},
        {Order: 116, New: NewStatusPageDataSource},
        {Order: 124, New: NewUserOverrideDataSource},
        {Order: 168, New: NewAlertUserOwnerDataSource},
        {Order: 169, New: NewAlertSeverityDataSource},
        {Order: 194, New: NewIncidentEpisodeStateTimelineDataSource},
        {Order: 329, New: NewWorkspaceNotificationLogDataSource},
        {Order: 393, New: NewTableViewDataSource},
        {Order: 394, New: NewDashboardDataSource},
        {Order: 418, New: NewVcenterFeedDataSource},
        {Order: 426, New: NewStatusPageScimLogDataSource},
        {Order: 553, New: NewVcenterFeedLegacyDataSource},
    }, []LegacyDataSourceAlias{
        {Order: 418, Name: "vcenter_feed", LegacyName: "v_center_feed", New: NewVcenterFeedDataSource, NewLegacy: NewVcenterFeedLegacyDataSource},
    })
}
