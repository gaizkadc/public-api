/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package nodes

import "github.com/nalej/grpc-infrastructure-go"

// Manager structure with the required clients for node operations.
type Manager struct {

}

func (m *Manager) ClusterNodes(clusterId *grpc_infrastructure_go.ClusterId) (*grpc_infrastructure_go.NodeList, error) {
	panic("implement me")
}