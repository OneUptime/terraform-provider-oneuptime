package shard17

import (
    . "github.com/oneuptime/terraform-provider-oneuptime/internal/provider"
)

// The data sources this package holds, handed to internal/provider with their
// places among the provider's, and the deprecated aliases among them (its
// shards.go).
func init() {
    RegisterDataSources("shard17", []OrderedDataSource{
        {Order: 17, New: NewExceptionSavedViewDataSource},
        {Order: 60, New: NewDataAssetDataSource},
        {Order: 64, New: NewDataPipelineDataSource},
        {Order: 98, New: NewTracePipelineProcessorDataSource},
        {Order: 116, New: NewIncomingCallPolicyEscalationRuleDataSource},
        {Order: 121, New: NewMonitorSecretDataSource},
        {Order: 122, New: NewIncomingRequestSecretDataSource},
        {Order: 177, New: NewAlertEpisodeOwnerRuleDataSource},
        {Order: 323, New: NewUserOnCallLogTimelineDataSource},
        {Order: 332, New: NewMonitorGroupResourceDataSource},
        {Order: 353, New: NewServiceDataSource},
        {Order: 409, New: NewHostFeedDataSource},
        {Order: 418, New: NewKubernetesClusterTeamOwnerDataSource},
        {Order: 420, New: NewDockerHostDataSource},
        {Order: 464, New: NewIotFleetUserOwnerDataSource},
        {Order: 467, New: NewCephClusterTeamOwnerDataSource},
        {Order: 474, New: NewVcenterUserOwnerDataSource},
        {Order: 496, New: NewDatabaseEndpointDataSource},
        {Order: 520, New: NewKubernetesCostAllocationDataSource},
        {Order: 549, New: NewIotFleetUserOwnerLegacyDataSource},
        {Order: 553, New: NewVcenterUserOwnerLegacyDataSource},
    }, []LegacyDataSourceAlias{
        {Order: 464, Name: "iot_fleet_user_owner", LegacyName: "io_t_fleet_user_owner", New: NewIotFleetUserOwnerDataSource, NewLegacy: NewIotFleetUserOwnerLegacyDataSource},
        {Order: 474, Name: "vcenter_user_owner", LegacyName: "v_center_user_owner", New: NewVcenterUserOwnerDataSource, NewLegacy: NewVcenterUserOwnerLegacyDataSource},
    })
}
