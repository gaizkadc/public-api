/*
 * Copyright 2020 Nalej
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

package cli2

import (
	"context"
	"github.com/nalej/derrors"
	"github.com/nalej/public-api/internal/app/options"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/metadata"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// TokenFileName with the name of the file we use to store the token.
const TokenFileName = "token"

// RefreshTokenFileName with the name of the file that contains the refresh token
const RefreshTokenFileName = "refresh_token"

type Credentials struct {
	BasePath     string
	Token        string
	RefreshToken string
}

func NewEmptyCredentials(basePath string) *Credentials {
	return &Credentials{BasePath: basePath}
}

// NewCredentials creates a new Credentials structure.
func NewCredentials(basePath string, token string, refreshToken string) *Credentials {
	return &Credentials{basePath, token, refreshToken}
}

// NewCredentials from disk reads the stored credentials from the Login operation.
func NewCredentialsFromDisk(basePath string) (*Credentials, derrors.Error) {
	tokenPath := filepath.Join(options.GetPath(basePath), TokenFileName)
	refreshTokenPath := filepath.Join(options.GetPath(basePath), RefreshTokenFileName)
	token, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		return nil, derrors.AsError(err, "cannot read token file")
	}
	refreshToken, err := ioutil.ReadFile(refreshTokenPath)
	if err != nil {
		return nil, derrors.AsError(err, "cannot read refresh token file")
	}
	return NewCredentials(basePath, string(token), string(refreshToken)), nil
}

func (c *Credentials) LoadCredentials() derrors.Error {
	tokenPath := filepath.Join(options.GetPath(c.BasePath), TokenFileName)
	refreshTokenPath := filepath.Join(options.GetPath(c.BasePath), RefreshTokenFileName)
	token, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		return derrors.AsError(err, "cannot read token file")
	}
	refreshToken, err := ioutil.ReadFile(refreshTokenPath)
	if err != nil {
		return derrors.AsError(err, "cannot read refresh token file")
	}
	c.Token = string(token)
	c.RefreshToken = string(refreshToken)
	return nil
}

// Store the credentials in disk
func (c *Credentials) Store() derrors.Error {
	rPath := options.GetPath(c.BasePath)
	_ = os.MkdirAll(rPath, 0700)
	tokenPath := filepath.Join(options.GetPath(c.BasePath), TokenFileName)
	refreshTokenPath := filepath.Join(options.GetPath(c.BasePath), RefreshTokenFileName)
	err := ioutil.WriteFile(tokenPath, []byte(c.Token), 0600)
	if err != nil {
		return derrors.AsError(err, "cannot write token file")
	}
	err = ioutil.WriteFile(refreshTokenPath, []byte(c.RefreshToken), 0600)
	if err != nil {
		return derrors.AsError(err, "cannot write refresh token file")
	}
	return nil
}

func (c *Credentials) GetContext(timeout ...time.Duration) (context.Context, context.CancelFunc) {
	md := metadata.New(map[string]string{AuthHeader: c.Token})
	log.Debug().Interface("md", md).Msg("metadata has been created")
	if len(timeout) == 0 {
		baseContext, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
		return metadata.NewOutgoingContext(baseContext, md), cancel
	}
	baseContext, cancel := context.WithTimeout(context.Background(), timeout[0])
	return metadata.NewOutgoingContext(baseContext, md), cancel
}
