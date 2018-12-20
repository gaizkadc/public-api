package cli

import (
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-user-go"
	"github.com/nalej/grpc-user-manager-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type Users struct {
	Connection
	Credentials
}

func NewUsers(address string, port int, insecure bool, caCertPath string) *Users {
	return &Users{
		Connection:  *NewConnection(address, port, insecure, caCertPath),
		Credentials: *NewEmptyCredentials(DefaultPath),
	}
}

func (u *Users) load() {
	err := u.LoadCredentials()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot load credentials, try login first")
	}
}

func (u *Users) getClient() (grpc_public_api_go.UsersClient, *grpc.ClientConn) {
	conn, err := u.GetConnection()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot create the connection with the Nalej platform")
	}
	client := grpc_public_api_go.NewUsersClient(conn)
	return client, conn
}

// Add a new user to the organization.
func (u *Users) Add(organizationID string, email string, password string, name string, roleName string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	if email == "" {
		log.Fatal().Msg("email cannot be empty")
	}

	u.load()
	ctx, cancel := u.GetContext()
	client, conn := u.getClient()
	defer conn.Close()
	defer cancel()

	addRequest := &grpc_public_api_go.AddUserRequest{
		OrganizationId: organizationID,
		Email:          email,
		Password:       password,
		Name:           name,
		RoleName:       roleName,
	}
	added, err := client.Add(ctx, addRequest)
	u.PrintResultOrError(added, err, "cannot add user")
}

// Info retrieves the information of a user.
func (u *Users) Info(organizationID string, email string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}

	u.load()
	ctx, cancel := u.GetContext()
	client, conn := u.getClient()
	defer conn.Close()
	defer cancel()

	userID := &grpc_user_go.UserId{
		OrganizationId: organizationID,
		Email:          email,
	}
	info, err := client.Info(ctx, userID)
	u.PrintResultOrError(info, err, "cannot obtain user info")
}

// List the users of an organization.
func (u *Users) List(organizationID string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}

	u.load()
	ctx, cancel := u.GetContext()
	client, conn := u.getClient()
	defer conn.Close()
	defer cancel()

	userID := &grpc_organization_go.OrganizationId{
		OrganizationId: organizationID,
	}
	users, err := client.List(ctx, userID)
	u.PrintResultOrError(users, err, "cannot obtain user list")
}

// Delete a user from an organization.
func (u *Users) Delete(organizationID string, email string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}

	u.load()
	ctx, cancel := u.GetContext()
	client, conn := u.getClient()
	defer conn.Close()
	defer cancel()

	userID := &grpc_user_go.UserId{
		OrganizationId: organizationID,
		Email:          email,
	}
	done, err := client.Delete(ctx, userID)
	u.PrintResultOrError(done, err, "cannot delete user")
}

// Reset the password of a user.
func (u *Users) ChangePassword(organizationID string, email string, password string, newPassword string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	u.load()
	ctx, cancel := u.GetContext()
	client, conn := u.getClient()
	defer conn.Close()
	defer cancel()

	passwordRequest := &grpc_user_manager_go.ChangePasswordRequest{
		OrganizationId: organizationID,
		Email:          email,
		Password:       password,
		NewPassword:    newPassword,
	}
	done, err := client.ChangePassword(ctx, passwordRequest)
	u.PrintResultOrError(done, err, "cannot change password")
}

// Update the user information.
func (u *Users) Update(organizationID string, email string, newName string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	u.load()
	ctx, cancel := u.GetContext()
	client, conn := u.getClient()
	defer conn.Close()
	defer cancel()

	updateRequest := &grpc_user_go.UpdateUserRequest{
		OrganizationId: organizationID,
		Email:          email,
	}
	if newName != "" {
		updateRequest.Name = newName
	}
	log.Debug().Interface("updateRequest", updateRequest).Msg("sending update request")
	done, err := client.Update(ctx, updateRequest)
	u.PrintResultOrError(done, err, "cannot update user")
}
