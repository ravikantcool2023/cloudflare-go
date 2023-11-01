package cloudflare

import (
	"context"
	"fmt"
	"net/http"

	"github.com/goccy/go-json"
)

type Enabled struct {
	Enabled bool `json:"enabled"`
}

// DeviceClientCertificatesZone identifies if the zero trust zone is configured for an account.
type DeviceClientCertificatesZone struct {
	Response
	Result Enabled
}

type ServiceMode string

const (
	oneDotOne                                 ServiceMode = "1dot1"
	warp                                      ServiceMode = "warp"
	proxy                                     ServiceMode = "proxy"
	postureOnly                               ServiceMode = "posture_only"
	warpTunnelOnly                            ServiceMode = "warp_tunnel_only"
	listDeviceSettingsPoliciesDefaultPageSize             = 20
)

type ServiceModeV2 struct {
	Mode ServiceMode `json:"mode,omitempty"`
	Port int         `json:"port,omitempty"`
}

type DeviceSettingsPolicy struct {
	ServiceModeV2       *ServiceModeV2    `json:"service_mode_v2"`
	DisableAutoFallback *bool             `json:"disable_auto_fallback"`
	FallbackDomains     *[]FallbackDomain `json:"fallback_domains"`
	Include             *[]SplitTunnel    `json:"include"`
	Exclude             *[]SplitTunnel    `json:"exclude"`
	GatewayUniqueID     *string           `json:"gateway_unique_id"`
	SupportURL          *string           `json:"support_url"`
	CaptivePortal       *int              `json:"captive_portal"`
	AllowModeSwitch     *bool             `json:"allow_mode_switch"`
	SwitchLocked        *bool             `json:"switch_locked"`
	AllowUpdates        *bool             `json:"allow_updates"`
	AutoConnect         *int              `json:"auto_connect"`
	AllowedToLeave      *bool             `json:"allowed_to_leave"`
	PolicyID            *string           `json:"policy_id"`
	Enabled             *bool             `json:"enabled"`
	Name                *string           `json:"name"`
	Match               *string           `json:"match"`
	Precedence          *int              `json:"precedence"`
	Default             bool              `json:"default"`
	ExcludeOfficeIps    *bool             `json:"exclude_office_ips"`
	Description         *string           `json:"description"`
}

type DeviceSettingsPolicyResponse struct {
	Response
	Result DeviceSettingsPolicy
}

type DeleteDeviceSettingsPolicyResponse struct {
	Response
	Result []DeviceSettingsPolicy
}

type DeviceSettingsPolicyRequest struct {
	DisableAutoFallback *bool          `json:"disable_auto_fallback,omitempty"`
	CaptivePortal       *int           `json:"captive_portal,omitempty"`
	AllowModeSwitch     *bool          `json:"allow_mode_switch,omitempty"`
	SwitchLocked        *bool          `json:"switch_locked,omitempty"`
	AllowUpdates        *bool          `json:"allow_updates,omitempty"`
	AutoConnect         *int           `json:"auto_connect,omitempty"`
	AllowedToLeave      *bool          `json:"allowed_to_leave,omitempty"`
	SupportURL          *string        `json:"support_url,omitempty"`
	ServiceModeV2       *ServiceModeV2 `json:"service_mode_v2,omitempty"`
	Precedence          *int           `json:"precedence,omitempty"`
	Name                *string        `json:"name,omitempty"`
	Match               *string        `json:"match,omitempty"`
	Enabled             *bool          `json:"enabled,omitempty"`
	ExcludeOfficeIps    *bool          `json:"exclude_office_ips"`
	Description         *string        `json:"description,omitempty"`
}

type ListDeviceSettingsPoliciesResponse struct {
	Response
	ResultInfo ResultInfo             `json:"result_info"`
	Result     []DeviceSettingsPolicy `json:"result"`
}

// UpdateDeviceClientCertificates controls the zero trust zone used to provision client certificates.
//
// API reference: https://api.cloudflare.com/#device-client-certificates
func (api *API) UpdateDeviceClientCertificatesZone(ctx context.Context, zoneID string, enable bool) (DeviceClientCertificatesZone, error) {
	uri := fmt.Sprintf("/%s/%s/devices/policy/certificates", ZoneRouteRoot, zoneID)

	result := DeviceClientCertificatesZone{}
	res, err := api.makeRequestContext(ctx, http.MethodPatch, uri, Enabled{enable})
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal(res, &result); err != nil {
		return result, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}

	return result, err
}

// GetDeviceClientCertificatesZone controls the zero trust zone used to provision client certificates.
//
// API reference: https://api.cloudflare.com/#device-client-certificates
func (api *API) GetDeviceClientCertificatesZone(ctx context.Context, zoneID string) (DeviceClientCertificatesZone, error) {
	uri := fmt.Sprintf("/%s/%s/devices/policy/certificates", ZoneRouteRoot, zoneID)

	result := DeviceClientCertificatesZone{}
	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal(res, &result); err != nil {
		return result, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}

	return result, err
}

// CreateDeviceSettingsPolicy creates a settings policy against devices that match the policy
//
// API reference: https://api.cloudflare.com/#devices-create-device-settings-policy
func (api *API) CreateDeviceSettingsPolicy(ctx context.Context, accountID string, req DeviceSettingsPolicyRequest) (DeviceSettingsPolicyResponse, error) {
	uri := fmt.Sprintf("/%s/%s/devices/policy", AccountRouteRoot, accountID)

	result := DeviceSettingsPolicyResponse{}
	res, err := api.makeRequestContext(ctx, http.MethodPost, uri, req)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal(res, &result); err != nil {
		return result, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}

	return result, err
}

// UpdateDefaultDeviceSettingsPolicy updates the default settings policy for an account
//
// API reference: https://api.cloudflare.com/#devices-update-default-device-settings-policy
func (api *API) UpdateDefaultDeviceSettingsPolicy(ctx context.Context, accountID string, req DeviceSettingsPolicyRequest) (DeviceSettingsPolicyResponse, error) {
	result := DeviceSettingsPolicyResponse{}
	uri := fmt.Sprintf("/%s/%s/devices/policy", AccountRouteRoot, accountID)
	res, err := api.makeRequestContext(ctx, http.MethodPatch, uri, req)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal(res, &result); err != nil {
		return result, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}

	return result, err
}

// UpdateDeviceSettingsPolicy updates a settings policy
//
// API reference: https://api.cloudflare.com/#devices-update-device-settings-policy
func (api *API) UpdateDeviceSettingsPolicy(ctx context.Context, accountID, policyID string, req DeviceSettingsPolicyRequest) (DeviceSettingsPolicyResponse, error) {
	uri := fmt.Sprintf("/%s/%s/devices/policy/%s", AccountRouteRoot, accountID, policyID)

	result := DeviceSettingsPolicyResponse{}
	res, err := api.makeRequestContext(ctx, http.MethodPatch, uri, req)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal(res, &result); err != nil {
		return result, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}

	return result, err
}

// DeleteDeviceSettingsPolicy deletes a settings policy and returns a list
// of all of the other policies in the account
//
// API reference: https://api.cloudflare.com/#devices-delete-device-settings-policy
func (api *API) DeleteDeviceSettingsPolicy(ctx context.Context, accountID, policyID string) (DeleteDeviceSettingsPolicyResponse, error) {
	uri := fmt.Sprintf("/%s/%s/devices/policy/%s", AccountRouteRoot, accountID, policyID)

	result := DeleteDeviceSettingsPolicyResponse{}
	res, err := api.makeRequestContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal(res, &result); err != nil {
		return result, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}

	return result, err
}

// GetDefaultDeviceSettings gets the default device settings policy
//
// API reference: https://api.cloudflare.com/#devices-get-default-device-settings-policy
func (api *API) GetDefaultDeviceSettingsPolicy(ctx context.Context, accountID string) (DeviceSettingsPolicyResponse, error) {
	uri := fmt.Sprintf("/%s/%s/devices/policy", AccountRouteRoot, accountID)

	result := DeviceSettingsPolicyResponse{}
	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal(res, &result); err != nil {
		return result, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}

	return result, err
}

// GetDefaultDeviceSettings gets the device settings policy by its policyID
//
// API reference: https://api.cloudflare.com/#devices-get-device-settings-policy-by-id
func (api *API) GetDeviceSettingsPolicy(ctx context.Context, accountID, policyID string) (DeviceSettingsPolicyResponse, error) {
	uri := fmt.Sprintf("/%s/%s/devices/policy/%s", AccountRouteRoot, accountID, policyID)

	result := DeviceSettingsPolicyResponse{}
	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal(res, &result); err != nil {
		return result, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}

	return result, err
}

type ListDeviceSettingsPoliciesParams struct {
	ResultInfo
}

// ListDeviceSettingsPolicies returns all device settings policies for an account
//
// API reference: https://api.cloudflare.com/#devices-list-device-settings-policies
func (api *API) ListDeviceSettingsPolicies(ctx context.Context, accountID string, params ListDeviceSettingsPoliciesParams) ([]DeviceSettingsPolicy, *ResultInfo, error) {

	autoPaginate := true
	if params.PerPage >= 1 || params.Page >= 1 {
		autoPaginate = false
	}

	if params.PerPage < 1 {
		params.PerPage = listDeviceSettingsPoliciesDefaultPageSize
	}

	var policies []DeviceSettingsPolicy
	var lastResultInfo ResultInfo
	for {
		uri := buildURI(fmt.Sprintf("/%s/%s/devices/policies", AccountRouteRoot, accountID), params)
		res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
		if err != nil {
			return nil, nil, err
		}
		var r ListDeviceSettingsPoliciesResponse
		err = json.Unmarshal(res, &r)
		if err != nil {
			return nil, nil, fmt.Errorf("%s: %w", errUnmarshalError, err)
		}
		policies = append(policies, r.Result...)
		lastResultInfo = r.ResultInfo
		params.ResultInfo = r.ResultInfo.Next()
		if params.ResultInfo.Done() || !autoPaginate {
			break
		}
	}
	return policies, &lastResultInfo, nil
}
