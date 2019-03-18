package server

import (
	"github.com/nalej/authx/pkg/interceptor"
	"github.com/nalej/derrors"
	"github.com/nalej/public-api/version"
	"github.com/rs/zerolog/log"
	"strings"
)

type Config struct {
	// Port where the gRPC API service will listen requests.
	Port int
	// HTTPPort where the HTTP gRPC gateway will be listening.
	HTTPPort int
	// SystemModelAddress with the host:port to connect to System Model
	SystemModelAddress string
	// InfrastructureManagerAddress with the host:port to connect to the Infrastructure Manager.
	InfrastructureManagerAddress string
	// ApplicationsManagerAddress with the host:port to connect to the Applications manager.
	ApplicationsManagerAddress string
	// UserManagerAddress with the host:port to connect to the Access manager.
	UserManagerAddress string
	// DeviceManagerAddress with the host:port to connect to the Device Manager component.
	DeviceManagerAddress string
	// UnifiedLoggingAddress with the host:port to connect to the Unified Logging Coordinator component.
	UnifiedLoggingAddress string
	// InfrastructureMonitorAddress with the host:port to connect to the Infrastructure Monitor Coordinator component.
	InfrastructureMonitorAddress string
	// AuthSecret contains the shared authx secret.
	AuthSecret string
	// AuthHeader contains the name of the target header.
	AuthHeader string
	// AuthConfigPath contains the path of the file with the authentication configuration.
	AuthConfigPath string
}

func (conf *Config) Validate() derrors.Error {

	if conf.Port <= 0 || conf.HTTPPort <= 0 {
		return derrors.NewInvalidArgumentError("ports must be valid")
	}

	if conf.SystemModelAddress == "" {
		return derrors.NewInvalidArgumentError("systemModelAddress must be set")
	}

	if conf.InfrastructureManagerAddress == "" {
		return derrors.NewInvalidArgumentError("infrastructureManagerAddress must be set")
	}

	if conf.ApplicationsManagerAddress == "" {
		return derrors.NewInvalidArgumentError("applicationsManagerAddress must be set")
	}

	if conf.InfrastructureManagerAddress == "" {
		return derrors.NewInvalidArgumentError("infrastructureManagerAddress must be set")
	}

	if conf.UserManagerAddress == "" {
		return derrors.NewInvalidArgumentError("userManagerAddress must be set")
	}

	if conf.DeviceManagerAddress == "" {
		return derrors.NewInvalidArgumentError("deviceManagerAddress must be set")
	}

	if conf.UnifiedLoggingAddress == "" {
		return derrors.NewInvalidArgumentError("unifiedLoggingAddress must be set")
	}

	if conf.InfrastructureMonitorAddress == "" {
		return derrors.NewInvalidArgumentError("infrastructureMonitorAddress must be set")
	}

	if conf.AuthHeader == "" || conf.AuthSecret == "" {
		return derrors.NewInvalidArgumentError("Authorization header and secret must be set")
	}

	if conf.AuthConfigPath == "" {
		return derrors.NewInvalidArgumentError("authConfigPath must be set")
	}

	return nil
}

// LoadAuthConfig loads the security configuration.
func (conf *Config) LoadAuthConfig() (*interceptor.AuthorizationConfig, derrors.Error) {
	return interceptor.LoadAuthorizationConfig(conf.AuthConfigPath)
}

func (conf *Config) Print() {
	log.Info().Str("app", version.AppVersion).Str("commit", version.Commit).Msg("Version")
	log.Info().Int("port", conf.Port).Msg("gRPC port")
	log.Info().Int("port", conf.HTTPPort).Msg("HTTP port")
	log.Info().Str("URL", conf.SystemModelAddress).Msg("System Model")
	log.Info().Str("URL", conf.InfrastructureManagerAddress).Msg("Infrastructure Manager")
	log.Info().Str("URL", conf.ApplicationsManagerAddress).Msg("Applications Manager")
	log.Info().Str("URL", conf.UserManagerAddress).Msg("User Manager")
	log.Info().Str("URL", conf.UnifiedLoggingAddress).Msg("Unified Logging Coordinator Service")
	log.Info().Str("URL", conf.InfrastructureMonitorAddress).Msg("Infrastructure Monitor Coordinator Service")
	log.Info().Str("URL", conf.DeviceManagerAddress).Msg("Device Manager Service")
	log.Info().Str("header", conf.AuthHeader).Str("secret", strings.Repeat("*", len(conf.AuthSecret))).Msg("Authorization")
	log.Info().Str("path", conf.AuthConfigPath).Msg("Permissions file")
}
