/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package resources

import (
	"context"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
)

// Manager structure with the required clients for resources operations.
type Manager struct {
	clustClient grpc_infrastructure_go.ClustersClient
	nodeClient grpc_infrastructure_go.NodesClient
}

// NewManager creates a Manager using a set of clients.
func NewManager(clustClient grpc_infrastructure_go.ClustersClient,
	nodeClient grpc_infrastructure_go.NodesClient) Manager {
	return Manager{
		clustClient: clustClient, nodeClient:nodeClient,
	}
}

func (m * Manager) getNumNodes(organizationID string, clusterID string) (int, derrors.Error){
	// Return number of nodes in a cluster
	cID := &grpc_infrastructure_go.ClusterId{
		OrganizationId:       organizationID,
		ClusterId:            clusterID,
	}
	clusterNodes, err := m.nodeClient.ListNodes(context.Background(), cID)
	if err != nil{
		return 0, conversions.ToDerror(err)
	}
	return len(clusterNodes.Nodes), nil
}

func (m * Manager) getSummary(organizationID *grpc_organization_go.OrganizationId) (int, int, derrors.Error){
	// Obtain list of clusters
	totalNodes := 0
	list, err := m.clustClient.ListClusters(context.Background(), organizationID)
	if err != nil{
		return 0, 0, conversions.ToDerror(err)
	}
	for _, c := range list.Clusters{
		n, err := m.getNumNodes(c.OrganizationId, c.ClusterId)
		if err != nil{
			return 0, 0, err
		}
		totalNodes += n
	}
	return len(list.Clusters), totalNodes, nil
}

func (m * Manager) Summary(organizationID *grpc_organization_go.OrganizationId) (*grpc_public_api_go.ResourceSummary, error) {
	totalClusters, totalNodes, err := m.getSummary(organizationID)
	if err != nil{
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_public_api_go.ResourceSummary{
		OrganizationId:       organizationID.OrganizationId,
		TotalClusters:        int64(totalClusters),
		TotalNodes:           int64(totalNodes),
	}, nil
}