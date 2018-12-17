package authhelper

import (
	"context"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-authx-go"
	"google.golang.org/grpc/metadata"
	"strings"
)

const UserIdField = "user_id"
const OrganizationIdField = "organization_id"

type RequestMetadata struct {
	UserID            string
	OrganizationID    string
	OrgPrimitive      bool
	AppsPrimitive     bool
	ResourcePrimitive bool
	ProfilePrimitive  bool
	AppClusterOpsPrimitive bool
}

func GetRequestMetadata(ctx context.Context) (*RequestMetadata, derrors.Error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, derrors.NewInvalidArgumentError("expecting JWT metadata")
	}
	userID, found := md[UserIdField]
	if !found {
		return nil, derrors.NewUnauthenticatedError("userID not found")
	}
	organizationID, found := md[OrganizationIdField]
	if !found {
		return nil, derrors.NewUnauthenticatedError("organizationID not found")
	}
	_, orgPrimitive := md[strings.ToLower(grpc_authx_go.AccessPrimitive_ORG.String())]
	_, appsPrimitive := md[strings.ToLower(grpc_authx_go.AccessPrimitive_APPS.String())]
	_, resourcePrimitive := md[strings.ToLower(grpc_authx_go.AccessPrimitive_RESOURCES.String())]
	_, profilePrimitive := md[strings.ToLower(grpc_authx_go.AccessPrimitive_PROFILE.String())]
	_, appClusterOpsPrimitive := md[strings.ToLower(grpc_authx_go.AccessPrimitive_APPCLUSTEROPS.String())]

	return &RequestMetadata{
		UserID:                 userID[0],
		OrganizationID:         organizationID[0],
		OrgPrimitive:           orgPrimitive,
		AppsPrimitive:          appsPrimitive,
		ResourcePrimitive:      resourcePrimitive,
		ProfilePrimitive:       profilePrimitive,
		AppClusterOpsPrimitive: appClusterOpsPrimitive,
	}, nil
}
