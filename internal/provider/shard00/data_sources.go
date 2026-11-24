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
        {Order: 47, New: NewDataStoreUserOwnerDataSource},
        {Order: 66, New: NewDataProductDataSource},
        {Order: 75, New: NewPipelineUserOwnerDataSource},
        {Order: 84, New: NewBlockEnforcementDataSource},
        {Order: 106, New: NewStatusPageDataSource},
        {Order: 114, New: NewUserOverrideDataSource},
        {Order: 158, New: NewAlertUserOwnerDataSource},
        {Order: 159, New: NewAlertSeverityDataSource},
        {Order: 184, New: NewIncidentEpisodeStateTimelineDataSource},
        {Order: 319, New: NewWorkspaceNotificationLogDataSource},
        {Order: 383, New: NewTableViewDataSource},
        {Order: 384, New: NewDashboardDataSource},
        {Order: 408, New: NewVcenterFeedDataSource},
        {Order: 416, New: NewStatusPageScimLogDataSource},
        {Order: 543, New: NewVcenterFeedLegacyDataSource},
    }, []LegacyDataSourceAlias{
        {Order: 408, Name: "vcenter_feed", LegacyName: "v_center_feed", New: NewVcenterFeedDataSource, NewLegacy: NewVcenterFeedLegacyDataSource},
    })
}
