/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cli

import (
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-organization-manager-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/rs/zerolog/log"
)

type Organizations struct {
	Connection
	Credentials
}

func NewOrganizations(address string, port int, insecure bool, useTLS bool, caCertPath string, output string, labelLength int) *Organizations {
	return &Organizations{
		Connection:  *NewConnection(address, port, insecure, useTLS, caCertPath, output, labelLength),
		Credentials: *NewEmptyCredentials(DefaultPath),
	}
}

func (o *Organizations) Info(organizationID string) *grpc_organization_manager_go.Organization {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}

	err := o.LoadCredentials()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot load credentials, try login first")
	}

	c, err := o.GetConnection()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot create the connection with the Nalej platform")
	}
	defer c.Close()
	ctx, cancel := o.GetContext()
	defer cancel()

	orgClient := grpc_public_api_go.NewOrganizationsClient(c)
	orgID := &grpc_organization_go.OrganizationId{
		OrganizationId: organizationID,
	}
	info, iErr := orgClient.Info(ctx, orgID)

	o.PrintResultOrError(info, iErr, "cannot obtain organization info")
	return info
}

func (o *Organizations) Update(organizationID string, updateName bool, name string, updateFullAddress bool, fullAddress string,
	updateCity bool, city string, updateState bool, state string, updateCountry bool, country string,
	updateZipCode bool, zipCode string, updatePhoto bool, photoPath string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}

	err := o.LoadCredentials()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot load credentials, try login first")
	}

	c, err := o.GetConnection()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot create the connection with the Nalej platform")
	}
	defer c.Close()

	photo64, err := PhotoToBase64(photoPath)
	if err != nil {
		o.PrintResultOrError(photoPath, err, "cannot open photo file")
	}

	ctx, cancel := o.GetContext()
	defer cancel()

	orgClient := grpc_public_api_go.NewOrganizationsClient(c)
	updateRequest := &grpc_organization_go.UpdateOrganizationRequest{
		OrganizationId:    organizationID,
		UpdateName:        updateName,
		Name:              name,
		UpdateFullAddress: updateFullAddress,
		FullAddress:       fullAddress,
		UpdateCity:        updateCity,
		City:              city,
		UpdateState:       updateState,
		State:             state,
		UpdateCountry:     updateCountry,
		Country:           country,
		UpdateZipCode:     updateZipCode,
		ZipCode:           zipCode,
		UpdatePhoto:       updatePhoto,
		PhotoBase64:       photo64,
	}
	success, iErr := orgClient.Update(ctx, updateRequest)
	o.PrintResultOrError(success, iErr, "cannot update organization info")
	return
}
