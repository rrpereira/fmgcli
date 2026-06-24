package fmgcli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"
)

//todo include verbose in the requests that need verbose or just take it out from all of them, decide something here
//todo decide the omit empty if that's done in all or none

type ServiceRespData struct {
	Name string `json:"name"`
}

type ServiceRespResult struct {
	Data   ServiceRespData `json:"data"`
	Status Status          `json:"status"`
	Url    string          `json:"url"`
}

type ServiceResp struct {
	Result []ServiceRespResult `json:"result"`
}

type ServiceReqData struct {
	Name         string                 `json:"name"`
	Protocol     string                 `json:"protocol,omitempty"`
	TCPPortRange []string               `json:"tcp-portrange,omitempty"`
	UDPPortRange []string               `json:"udp-portrange,omitempty"`
	Comment      string                 `json:"comment,omitempty"`
	Metafields   map[string]interface{} `json:"meta fields,omitempty"`
}

type ServiceReqParams struct {
	Data []ServiceReqData `json:"data"`
	Url  string           `json:"url"`
}

type ServiceReq struct {
	Method  string             `json:"method"`
	Params  []ServiceReqParams `json:"params"`
	Session string             `json:"session,omitempty" redact:"true"`
}

type SubnetAddressRespData struct {
	Name string `json:"name"`
}

type SubnetAddressRespResult struct {
	Status Status                `json:"status"`
	Url    string                `json:"url"`
	Data   SubnetAddressRespData `json:"data"`
}

type SubnetAddressResp struct {
	Result []SubnetAddressRespResult `json:"result"`
}

type SubnetAddressReqData struct {
	Name       string                 `json:"name"`
	Type       string                 `json:"type,omitempty"`
	Subnet     string                 `json:"subnet,omitempty"`
	Comment    string                 `json:"comment,omitempty"`
	Metafields map[string]interface{} `json:"meta fields,omitempty"`
}

type SubnetAddressReqParams struct {
	Data []SubnetAddressReqData `json:"data"`
	Url  string                 `json:"url"`
}

type SubnetAddressReq struct {
	Method  string                   `json:"method"`
	Params  []SubnetAddressReqParams `json:"params"`
	Session string                   `json:"session,omitempty" redact:"true"`
}

type GetPolicyData struct {
	Action  string   `json:"action"`
	Srcaddr []string `json:"srcaddr"`
	Dstaddr []string `json:"dstaddr"`
	Objseq  int      `json:"obj seq"`
	//Oid 	 int      `json:"oid"`
	PolicyID   int                    `json:"policyid"`
	Service    []string               `json:"service"`
	Status     string                 `json:"status"`
	Metafields map[string]interface{} `json:"meta fields"`
	//Vpn_dst_node string   `json:"vpn_dst_node"`
	//Vpn_src_node string   `json:"vpn_src_node"`
}

type GetPolicyResult struct {
	Data   GetPolicyData `json:"data"`
	Status Status        `json:"status"`
	Url    string        `json:"url"`
}

type GetPolicyResp struct {
	Result []GetPolicyResult `json:"result"`
}

type GetPoliciesResult struct {
	Data   []GetPolicyData `json:"data"`
	Status Status          `json:"status"`
	Url    string          `json:"url"`
}

type GetPoliciesResp struct {
	Result []GetPoliciesResult `json:"result"`
}

type GetPoliciesParams struct {
	Fields []string `json:"fields"`
	Option []string `json:"option"`
	Url    string   `json:"url"`
}

type GetPoliciesReq struct {
	Method  string              `json:"method"`
	Params  []GetPoliciesParams `json:"params"`
	Session string              `json:"session,omitempty" redact:"true"`
	Verbose int                 `json:"verbose,omitempty"`
}

type CommitResp struct {
	Result []Result `json:"result"`
}

type CommitParams struct {
	Url string `json:"url"`
}

type CommitReq struct {
	Method  string         `json:"method"`
	Params  []CommitParams `json:"params"`
	Session string         `json:"session,omitempty" redact:"true"`
}

type GetServiceData struct {
	Name         string                 `json:"name"`
	ObjSeq       int                    `json:"obj seq"`
	Oid          int                    `json:"oid"`
	Protocol     string                 `json:"protocol,omitempty"`
	TCPPortRange []string               `json:"tcp-portrange"`
	UDPPortRange []string               `json:"udp-portrange"`
	Metafields   map[string]interface{} `json:"meta fields,omitempty"`
}

type GetServiceResult struct {
	Data   GetServiceData `json:"data"`
	Status Status         `json:"status"`
	Url    string         `json:"url"`
}

type GetServiceResp struct {
	Result []GetServiceResult `json:"result"`
}

type GetServicesParams struct {
	Fields []string `json:"fields"`
	Option []string `json:"option,omitempty"`
	Url    string   `json:"url"`
}

type GetServicesReq struct {
	Method  string              `json:"method"`
	Params  []GetServicesParams `json:"params"`
	Session string              `json:"session,omitempty" redact:"true"`
	Verbose int                 `json:"verbose,omitempty"`
}

type GetServicesResult struct {
	Data   []GetServiceData `json:"data,omitempty"`
	Status Status           `json:"status,omitempty"`
	Url    string           `json:"url,omitempty"`
}

type GetServicesResp struct {
	Result []GetServicesResult `json:"result,omitempty"`
}

type GetAddressData struct {
	DynamicMapping []string               `json:"dynamic_mapping"` ////////todo check if this is a list of strings
	List           []string               `json:"list"`            /////////todo check if this is a list of strings
	Metafields     map[string]interface{} `json:"meta fields"`
	Name           string                 `json:"name"`
	Oid            int                    `json:"oid"`
	Subnet         []string               `json:"subnet"`
	Tagging        []string               `json:"tagging"` ////////todo check if this is a list of strings
}

type GetAddressResult struct {
	Data   GetAddressData `json:"data,omitempty"`
	Status Status         `json:"status,omitempty"`
	Url    string         `json:"url,omitempty"`
}

type GetAddressResp struct {
	Result []GetAddressResult `json:"result,omitempty"`
}

type GetAddressesResult struct {
	Data   []GetAddressData `json:"data,omitempty"`
	Status Status           `json:"status,omitempty"`
	Url    string           `json:"url,omitempty"`
}

type GetAddressesResp struct {
	Result []GetAddressesResult `json:"result,omitempty"`
}

type GetAddressesParams struct {
	Fields []string `json:"fields"`
	Option []string `json:"option,omitempty"`
	Url    string   `json:"url"`
}

type GetAddressesReq struct {
	Method  string               `json:"method"`
	Params  []GetAddressesParams `json:"params"`
	Session string               `json:"session,omitempty" redact:"true"`
	Verbose int                  `json:"verbose,omitempty"`
}

type CreatePolicyRespData struct {
	PolicyID int `json:"policyid"`
}

type CreatePolicyResult struct {
	Data   CreatePolicyRespData `json:"data"`
	Status Status               `json:"status"`
	Url    string               `json:"url"`
}

type AddPolicyResp struct {
	Result []CreatePolicyResult `json:"result"`
}

type CreatePolicyReqData struct {
	Action     string                 `json:"action"`
	Srcintf    []string               `json:"srcintf"`
	Dstinft    []string               `json:"dstintf"`
	Srcaddr    []string               `json:"srcaddr"`
	Dstaddr    []string               `json:"dstaddr"`
	Service    []string               `json:"service"`
	Status     string                 `json:"status"`
	Schedule   string                 `json:"schedule"`
	Logtraffic string                 `json:"logtraffic"`
	Metafields map[string]interface{} `json:"meta fields,omitempty"`
	Comments   string                 `json:"comments"`
}

type CreatePolicyParams struct {
	Data CreatePolicyReqData `json:"data"`
	Url  string              `json:"url"`
}

type CreatePolicyReq struct {
	Method  string               `json:"method"`
	Params  []CreatePolicyParams `json:"params"`
	Session string               `json:"session,omitempty" redact:"true"`
}

type DisableRespData struct {
	PolicyID int `json:"policyid"`
}

type DisableResult struct {
	Data   DisableRespData `json:"data"`
	Status Status          `json:"status"`
	Url    string          `json:"url"`
}

type DisableResp struct {
	Result []DisableResult `json:"result"`
}

type DisableReqData struct {
	Status string `json:"status"`
}

type DisableParams struct {
	Data DisableReqData `json:"data"`
	Url  string         `json:"url"`
}

type DisableReq struct {
	Method  string          `json:"method"`
	Params  []DisableParams `json:"params"`
	Session string          `json:"session,omitempty" redact:"true"`
}

type DeleteResp struct {
	Result []Result `json:"result"`
}

type DeleteParams struct {
	Data []string `json:"data"`
	Url  string   `json:"url"`
}

type DeleteReq struct {
	Method  string       `json:"method"`
	Params  DeleteParams `json:"params"`
	Session string       `json:"session,omitempty" redact:"true"`
}

type LoginData struct {
	User     string `json:"user"`
	Password string `json:"password" redact:"true"`
}

type LoginParams struct {
	Data LoginData `json:"data"`
	Url  string    `json:"url"`
}

type LogoutParams struct {
	Url string `json:"url"`
}

type LoginRequest struct {
	Method string        `json:"method"`
	Params []LoginParams `json:"params"`
}

type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Result struct {
	Status Status `json:"status"`
	Url    string `json:"url"`
}

type LoginResp struct {
	Result  []Result `json:"result"`
	Session string   `json:"session,omitempty" redact:"true"`
}

type LogoutResp struct {
	Result  []Result `json:"result"`
	Session string   `json:"session,omitempty" redact:"true"`
}

type LogoutReq struct {
	Method  string         `json:"method"`
	Params  []LogoutParams `json:"params"`
	Session string         `json:"session,omitempty" redact:"true"`
}

type UnlockResp struct {
	Result []Result `json:"result"`
}

type UnlockParams struct {
	Url string `json:"url"`
}

type UnlockReq struct {
	Method  string         `json:"method"`
	Params  []UnlockParams `json:"params"`
	Session string         `json:"session,omitempty" redact:"true"`
}

type LockResp struct {
	Result []Result `json:"result"`
}

type UrlParam struct {
	Url string `json:"url"`
}

type LockReq struct {
	Method  string     `json:"method"`
	Params  []UrlParam `json:"params"`
	Session string     `json:"session,omitempty" redact:"true"`
}

type Client struct {
	Host     string
	Key      string
	User     string
	Password string
	Session  string
	log      *slog.Logger
}

type ClientOptions func(*Client)
type PolicyOptions func(*CreatePolicyReqData)
type AddressOptions func(*SubnetAddressReqData)
type ServiceOptions func(*ServiceReqData)

func WithLog(l *slog.Logger) ClientOptions {
	return func(c *Client) { c.log = l }
}

func WithMetafields(metafields map[string]interface{}) PolicyOptions {
	return func(data *CreatePolicyReqData) {
		data.Metafields = metafields
	}
}

func WithAddressMetafields(metafields map[string]interface{}) AddressOptions {
	return func(data *SubnetAddressReqData) {
		data.Metafields = metafields
	}
}

func WithServiceMetafields(metafields map[string]interface{}) ServiceOptions {
	return func(data *ServiceReqData) {
		data.Metafields = metafields
	}
}

func NewUserClient(host, user, password string, opts ...ClientOptions) *Client {
	client := &Client{
		Host:     host,
		User:     user,
		Password: password,
		Session:  "",
		log:      slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})),
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func NewAPIClient(host, key string, opts ...ClientOptions) *Client {
	client := &Client{
		Host: host,
		Key:  key,
		log:  slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})),
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (fm *Client) CreateService(adom, name, protocol string, minPort, maxPort int, comment string, opts ...ServiceOptions) error {
	request := ServiceReq{
		Method: "add",
		Params: []ServiceReqParams{
			{
				Data: []ServiceReqData{
					{
						Name:    name,
						Comment: comment,
					},
				},
				Url: fmt.Sprintf("/pm/config/adom/%s/obj/firewall/service/custom", adom),
			},
		},
		Session: fm.Session,
	}

	if minPort > maxPort {
		fm.log.Error("Invalid port range for service: minPort can't be greater than maxPort.", "name", name, "minPort", minPort, "maxPort", maxPort)
		return fmt.Errorf("invalid port range for service '%s': minPort (%d) cannot be greater than maxPort (%d)", name, minPort, maxPort)
	}

	if protocol == "tcp" {
		request.Params[0].Data[0].Protocol = "TCP/UDP/SCTP"
		request.Params[0].Data[0].TCPPortRange = []string{fmt.Sprintf("%d-%d", minPort, maxPort)}
	} else if protocol == "udp" {
		request.Params[0].Data[0].Protocol = "TCP/UDP/SCTP"
		request.Params[0].Data[0].UDPPortRange = []string{fmt.Sprintf("%d-%d", minPort, maxPort)}
	} else {
		fm.log.Error("Invalid protocol for service creation.", "protocol", protocol, "name", name)
		return fmt.Errorf("invalid protocol '%s' for service '%s'. Only 'tcp' or 'udp' are allowed.", protocol, name)
	}

	for _, opt := range opts {
		opt(&request.Params[0].Data[0])
	}

	var response GetServiceResp

	err := fm.makeRequest(http.MethodPost, fm.Host, "/jsonrpc", request, &response)
	if err != nil {
		fm.log.Error("Error creating service.", "error", err, "name", name, "protocol", protocol, "minPort", minPort, "maxPort", maxPort)
		return err
	}

	if len(response.Result) == 0 || response.Result[0].Status.Code != 0 || response.Result[0].Status.Message != "OK" {
		fm.log.Error("Failed to create service.", "name", name, "protocol", protocol, "minPort", minPort, "maxPort", maxPort, "error", response.Result[0].Status.Message)
		return fmt.Errorf("failed to create service: %s", response.Result[0].Status.Message)
	}

	fm.log.Debug("Successfully created service.", "adom", adom, "name", name, "protocol", protocol, "minPort", minPort, "maxPort", maxPort)

	return nil
}

func (fm *Client) CreateSubnetAddress(adom, name, subnet, netmask, comment string, opts ...AddressOptions) error {
	request := SubnetAddressReq{
		Method: "add",
		Params: []SubnetAddressReqParams{
			{
				Data: []SubnetAddressReqData{
					{
						Name:    name,
						Type:    "ipmask",
						Subnet:  fmt.Sprintf("%s/%s", subnet, netmask),
						Comment: comment,
					},
				},
				Url: fmt.Sprintf("/pm/config/adom/%s/obj/firewall/address", adom),
			},
		},
		Session: fm.Session,
	}

	for _, opt := range opts {
		opt(&request.Params[0].Data[0])
	}

	var response SubnetAddressResp

	err := fm.makeRequest(http.MethodPost, fm.Host, "/jsonrpc", request, &response)
	if err != nil {
		fm.log.Error("Error creating subnet address.", "error", err, "adom", adom, "name", name, "subnet", subnet, "netmask", netmask)
		return err
	}

	if len(response.Result) == 0 || response.Result[0].Status.Code != 0 || response.Result[0].Status.Message != "OK" {
		fm.log.Error("Failed to create subnet address.", "adom", adom, "name", name, "subnet", subnet, "netmask", netmask, "error", response.Result[0].Status.Message)
		return fmt.Errorf("failed to create subnet address: %s", response.Result[0].Status.Message)
	}

	fm.log.Debug("Successfully created subnet address.", "adom", adom, "name", name, "subnet", subnet, "netmask", netmask)

	return nil
}

func (fm *Client) CreatePolicy(pkg, adom, comment string, srcs, dsts, services []string, opts ...PolicyOptions) error {
	request := CreatePolicyReq{
		Method: "add",
		Params: []CreatePolicyParams{
			{
				Data: CreatePolicyReqData{
					Action:     "accept",
					Srcintf:    []string{"any"},
					Dstinft:    []string{"any"},
					Srcaddr:    srcs,
					Dstaddr:    dsts,
					Service:    services,
					Status:     "enable",
					Schedule:   "always",
					Logtraffic: "all",
					Comments:   comment,
				},
				Url: fmt.Sprintf("/pm/config/adom/%s/pkg/%s/firewall/policy", adom, pkg),
			},
		},
		Session: fm.Session,
	}

	for _, opt := range opts {
		opt(&request.Params[0].Data)
	}

	var response AddPolicyResp

	err := fm.makeRequest(http.MethodPost, fm.Host, "/jsonrpc", request, &response)
	if err != nil {
		fm.log.Error("Error creating policy.", "error", err, "adom", adom, "pkg", pkg, "srcs", srcs, "dsts", dsts, "services", services)
		return err
	}

	if len(response.Result) == 0 || response.Result[0].Status.Code != 0 || response.Result[0].Status.Message != "OK" {
		fm.log.Error("Failed to create policy.", "adom", adom, "pkg", pkg, "srcs", srcs, "dsts", dsts, "services", services, "error", response.Result[0].Status.Message)
		return fmt.Errorf("failed to create policy: %s", response.Result[0].Status.Message)
	}

	fm.log.Debug("Successfully created policy.", "adom", adom, "pkg", pkg, "id", response.Result[0].Data.PolicyID, "sources", srcs, "destinations", dsts, "services", services)

	return nil
}

func (fm *Client) Commit(adom string) error {
	request := CommitReq{
		Method: "exec",
		Params: []CommitParams{
			{
				Url: fmt.Sprintf("/dvmdb/adom/%s/workspace/commit", adom),
			},
		},
		Session: fm.Session,
	}

	var response CommitResp

	err := fm.makeRequest(http.MethodPost, fm.Host, "/jsonrpc", request, &response)
	if err != nil {
		fm.log.Error("Error committing ADOM.", "error", err, "adom", adom)
		return err
	}

	if len(response.Result) == 0 || response.Result[0].Status.Code != 0 || response.Result[0].Status.Message != "OK" {
		fm.log.Error("Failed to commit ADOM.", "error", response.Result[0].Status.Message, "adom", adom)
		return fmt.Errorf("failed to save ADOM: %s", response.Result[0].Status.Message)
	}

	fm.log.Debug("Successfully committed ADOM.", "adom", adom)
	return nil
}

func (fm *Client) DeleteFromPolicy(id int, objects []string, role, vdom, device, adom string) error {
	if role != "srcaddr" && role != "dstaddr" {
		fm.log.Error("Invalid role when deleting objects from policy.", "role", role, "adom", adom, "device", device, "vdom", vdom, "id", id, "objects", objects)
		return errors.New("invalid role")
	}

	request := DeleteReq{
		Method: "delete",
		Params: DeleteParams{
			Data: objects,
			Url:  fmt.Sprintf("/pm/config/adom/%s/pkg/%s/%s/firewall/policy/%d/%s", adom, device, vdom, id, role),
		},
		Session: fm.Session,
	}

	var response DeleteResp

	err := fm.makeRequest(http.MethodPost, fm.Host, "/jsonrpc", request, &response)
	if err != nil {
		fm.log.Error("Error deleting objects from policy.", "error", err, "role", role, "adom", adom, "device", device, "vdom", vdom, "id", id, "objects", objects)
		return err
	}

	if len(response.Result) == 0 || response.Result[0].Status.Code != 0 || response.Result[0].Status.Message != "OK" {
		fm.log.Error("Failed to delete objects from policy.", "error", response.Result[0].Status.Message, "role", role, "adom", adom, "device", device, "vdom", vdom, "id", id, "objects", objects)
		return fmt.Errorf("failed to delete objects: %s", response.Result[0].Status.Message)
	}

	fm.log.Debug("Successfully deleted objects from policy.", "adom", adom, "device", device, "vdom", vdom, "id", id, "role", role, "objects", objects)

	return nil
}

func (fm *Client) DisablePolicies(ids []int, vdom, device, adom string) error {
	var params []DisableParams
	var problematic []int
	var successful []int

	for _, id := range ids {
		params = append(params, DisableParams{
			Data: DisableReqData{
				Status: "disable",
			},
			Url: fmt.Sprintf("/pm/config/adom/%s/pkg/%s/%s/firewall/policy/%d", adom, device, vdom, id),
		})
	}

	request := DisableReq{
		Method:  "set",
		Params:  params,
		Session: fm.Session,
	}

	var response DisableResp

	err := fm.makeRequest(http.MethodPost, fm.Host, "/jsonrpc", request, &response)
	if err != nil {
		fm.log.Error("Error disabling policies.", "error", err, "adom", adom, "device", device, "vdom", vdom, "ids", ids)
		return err
	}

	for i, result := range response.Result {
		if result.Status.Code != 0 || result.Status.Message != "OK" {
			problematic = append(problematic, ids[i])
		} else {
			successful = append(successful, ids[i])
		}
	}

	if len(problematic) != 0 {
		fm.log.Error("Failed to disable some policies.", "adom", adom, "device", device, "vdom", vdom, "problematic", problematic)

		if len(successful) != 0 {
			fm.log.Debug("Successfully disabled some policies.", "adom", adom, "device", device, "vdom", vdom, "successful", successful)
		}

		return fmt.Errorf("failed to disable the following policies: %v", problematic)
	}

	fm.log.Debug("Successfully disabled policies.", "adom", adom, "device", device, "vdom", vdom, "ids", ids)

	return nil
}

func (fm *Client) DisablePolicy(adom, pkg string, id int) error {
	request := DisableReq{
		Method: "set",
		Params: []DisableParams{
			{
				Data: DisableReqData{
					Status: "disable",
				},
				Url: fmt.Sprintf("/pm/config/adom/%s/pkg/%s/firewall/policy/%d", adom, pkg, id),
			},
		},
		Session: fm.Session,
	}

	var response DisableResp

	err := fm.makeRequest(http.MethodPost, fm.Host, "/jsonrpc", request, &response)
	if err != nil {
		fm.log.Error("Error disabling policy.", "error", err, "adom", adom, "pkg", pkg, "id", id)
		return err
	}

	if len(response.Result) == 0 || response.Result[0].Status.Code != 0 || response.Result[0].Status.Message != "OK" {
		fm.log.Error("Failed to disable policy.", "error", response.Result[0].Status.Message, "adom", adom, "pkg", pkg, "id", id)
		return fmt.Errorf("failed to disable policy: %s", response.Result[0].Status.Message)
	}

	fm.log.Debug("Successfully disabled policy.", "adom", adom, "pkg", pkg, "id", id)

	return nil
}

func (fm *Client) GetAddressByName(adom, objectName string) (*GetAddressData, error) {
	request := GetAddressesReq{
		Method: "get",
		Params: []GetAddressesParams{
			{
				Fields: []string{"name", "subnet"},
				Option: []string{"get meta"},
				Url:    fmt.Sprintf("/pm/config/adom/%s/obj/firewall/address/%s", adom, objectName),
			},
		},
		Session: fm.Session,
	}

	var response GetAddressResp

	err := fm.makeRequest(http.MethodPost, fm.Host, "/jsonrpc", request, &response)
	if err != nil {
		fm.log.Error("Error fetching address by name.", "error", err, "adom", adom, "objectName", objectName)
		return nil, err
	}

	if len(response.Result) > 1 {
		fm.log.Error("Multiple addresses found for object name.", "objectName", objectName, "count", len(response.Result))
		return nil, fmt.Errorf("multiple addresses found for '%s': %d results", objectName, len(response.Result))
	}

	if response.Result[0].Status.Code == -3 || response.Result[0].Status.Message == "Object does not exist" {
		fm.log.Debug("Address not found in ADOM.", "adom", adom, "objectName", objectName)
		return nil, fmt.Errorf("address '%s' not found in ADOM '%s'", objectName, adom)
	}

	if response.Result[0].Status.Code != 0 || response.Result[0].Status.Message != "OK" {
		fm.log.Error("Error fetching address by name.", "error", response.Result[0].Status.Message, "objectName", objectName)
		return nil, fmt.Errorf("error fetching address '%s': %s", objectName, response.Result[0].Status.Message)
	}

	fm.log.Debug("Successfully fetched address by name.", "adom", adom, "objectName", objectName, "addressData", response.Result[0].Data)

	return &response.Result[0].Data, nil
}

func (fm *Client) GetAddressByNameIPAndNetmask(adom, objectName, ip, netmask string) (*GetAddressData, error) {
	request := GetAddressesReq{
		Method: "get",
		Params: []GetAddressesParams{
			{
				Fields: []string{"name", "subnet"},
				Option: []string{"get meta"},
				Url:    fmt.Sprintf("/pm/config/adom/%s/obj/firewall/address/%s", adom, objectName),
			},
		},
		Session: fm.Session,
	}

	var response GetAddressResp

	err := fm.makeRequest(http.MethodPost, fm.Host, "/jsonrpc", request, &response)
	if err != nil {
		fm.log.Error("Error fetching address by name, IP and netmask.", "error", err, "adom", adom, "objectName", objectName, "ip", ip, "netmask", netmask)
		return nil, err
	}

	if len(response.Result) > 1 {
		fm.log.Error("Multiple addresses found for object name.", "objectName", objectName, "count", len(response.Result))
		return nil, fmt.Errorf("multiple addresses found for '%s': %d results", objectName, len(response.Result))
	}

	if response.Result[0].Status.Code == -3 || response.Result[0].Status.Message == "Object does not exist" {
		fm.log.Error("Address not found in ADOM.", "adom", adom, "objectName", objectName)
		return nil, fmt.Errorf("address '%s' not found in ADOM '%s'", objectName, adom)
	}

	if response.Result[0].Status.Code != 0 || response.Result[0].Status.Message != "OK" {
		fm.log.Error("Error fetching address by name, IP and netmask.", "error", response.Result[0].Status.Message, "objectName", objectName, "ip", ip, "netmask", netmask)
		return nil, fmt.Errorf("error fetching address '%s': %s", objectName, response.Result[0].Status.Message)
	}

	if response.Result[0].Data.Subnet[0] != ip || response.Result[0].Data.Subnet[1] != netmask {
		fm.log.Error("Address does not match IP or netmask.", "adom", adom, "objectName", objectName, "ip", ip, "netmask", netmask)
		return nil, fmt.Errorf("address '%s' does not match IP '%s' and netmask '%s' in ADOM '%s'", objectName, ip, netmask, adom)
	}

	fm.log.Debug("Successfully fetched address by name, IP and netmask.", "adom", adom, "objectName", objectName, "ip", ip, "netmask", netmask, "addressData", response.Result[0].Data)

	return &response.Result[0].Data, nil
}

func (fm *Client) GetAddressByMetafield(adom, key string, value interface{}) (*GetAddressData, error) {
	found := false
	var address GetAddressData

	request := GetAddressesReq{
		Method: "get",
		Params: []GetAddressesParams{
			{
				Fields: []string{"name", "subnet", "type"},
				Option: []string{"get meta"},
				Url:    fmt.Sprintf("/pm/config/adom/%s/obj/firewall/address", adom),
			},
		},
		Session: fm.Session,
		Verbose: 1, //todo make this dynamic, without verbose, type will be 15 instead of "ipmask"
	}

	var response GetAddressesResp

	err := fm.makeRequest(http.MethodPost, fm.Host, "/jsonrpc", request, &response)
	if err != nil {
		fm.log.Error("Error fetching addresses by metafield.", "error", err, "adom", adom, "key", key, "value", value)
		return nil, err
	}

	if len(response.Result) == 0 {
		fm.log.Error("Error fetching addresses: no result in response.", "adom", adom, "key", key, "value", value, "response", response)
		return nil, fmt.Errorf("error fetching addresses: no result in response for ADOM '%s'", adom)
	}

	if response.Result[0].Status.Code != 0 || response.Result[0].Status.Message != "OK" {
		fm.log.Error("Error fetching addresses.", "error", response.Result[0].Status.Message, "adom", adom, "key", key, "value", value)
		return nil, fmt.Errorf("error fetching addresses: %s", response.Result[0].Status.Message)
	}

	for _, result := range response.Result {
		for _, data := range result.Data {
			if metafieldValue, ok := data.Metafields[key]; ok && reflect.DeepEqual(metafieldValue, value) {
				if found {
					fm.log.Error("Multiple addresses found with metafield.", "key", key, "value", value, "adom", adom)
					return nil, fmt.Errorf("multiple addresses found with metafield '%s'='%s' in ADOM '%s'", key, value, adom)
				}

				found = true
				address = data
			}
		}
	}

	if !found {
		fm.log.Debug("No address found with metafield.", "key", key, "value", value, "adom", adom)
		return nil, fmt.Errorf("no address found with metafield '%s'='%s' in ADOM '%s'", key, value, adom)
	}

	fm.log.Debug("Successfully fetched address by metafield.", "adom", adom, "key", key, "value", value, "addressData", &response.Result[0].Data)

	return &address, nil
}

func (fm *Client) GetAddressesByMetafield(adom, key string, values []interface{}) ([]GetAddressData, error) {
	var addresses []GetAddressData

	request := GetAddressesReq{
		Method: "get",
		Params: []GetAddressesParams{
			{
				Fields: []string{"name", "subnet", "type"},
				Option: []string{"get meta"},
				Url:    fmt.Sprintf("/pm/config/adom/%s/obj/firewall/address", adom),
			},
		},
		Session: fm.Session,
		Verbose: 1, //todo make this dynamic, without verbose, type will be 15 instead of "ipmask"
	}

	var response GetAddressesResp

	err := fm.makeRequest(http.MethodPost, fm.Host, "/jsonrpc", request, &response)
	if err != nil {
		fm.log.Error("Error fetching addresses by metafield.", "error", err, "adom", adom, "key", key, "values", values)
		return nil, err
	}

	if len(response.Result) == 0 {
		fm.log.Error("Error fetching addresses: no result in response.", "adom", adom, "key", key, "values", values, "response", response)
		return nil, fmt.Errorf("error fetching addresses: no result in response for ADOM '%s'", adom)
	}

	if response.Result[0].Status.Code != 0 || response.Result[0].Status.Message != "OK" {
		fm.log.Error("Error fetching addresses.", "error", response.Result[0].Status.Message, "adom", adom, "key", key, "values", values)
		return nil, fmt.Errorf("error fetching addresses: %s", response.Result[0].Status.Message)
	}

	for _, value := range values {
		found := false
		for _, result := range response.Result {
			for _, data := range result.Data {
				if metafieldValue, ok := data.Metafields[key]; ok && reflect.DeepEqual(metafieldValue, value) {
					if found {
						fm.log.Error("Multiple addresses found with metafield.", "adom", adom, "key", key, "value", value)
						return nil, fmt.Errorf("multiple addresses found with metafield '%s'='%v' in ADOM '%s'", key, value, adom)
					}

					addresses = append(addresses, data)
					fm.log.Debug("Found address with metafield.", "adom", adom, "key", key, "value", value, "addressData", data)
					found = true
				}
			}
		}
	}

	if len(addresses) != len(values) {
		fm.log.Warn("Found addresses by metafield, but count does not match expected.", "adom", adom, "key", key, "values", values, "foundCount", len(addresses), "expectedCount", len(values))
		return nil, fmt.Errorf("found %d addresses corresponding the following '%s' metafield values: '%v', but expected %d in ADOM '%s'", len(addresses), key, values, len(values), adom)
	}

	fm.log.Debug("Successfully fetched addresses by metafield.", "adom", adom, "key", key, "values", values, "addressesCount", len(addresses))

	return addresses, nil
}

func (fm *Client) GetAddressesByName(objectNames []string, adom string) ([]GetAddressData, error) {
	var params []GetAddressesParams
	var problematic []string

	for _, objectName := range objectNames {
		params = append(params, GetAddressesParams{
			Fields: []string{"name", "subnet"},
			Url:    fmt.Sprintf("/pm/config/adom/%s/obj/firewall/address/%s", adom, objectName),
		})
	}

	request := GetAddressesReq{
		Method:  "get",
		Params:  params,
		Session: fm.Session,
		Verbose: 1, //todo make this dynamic, without verbose, type will be 15 instead of "ipmask"
	}

	var response GetAddressResp

	err := fm.makeRequest(http.MethodPost, fm.Host, "/jsonrpc", request, &response)
	if err != nil {
		fm.log.Error("Error fetching addresses by names.", "error", err, "adom", adom, "objectNames", objectNames)
		return nil, err
	}

	var addressesData []GetAddressData
	for i, result := range response.Result {
		if result.Status.Code != 0 || result.Status.Message != "OK" {
			problematic = append(problematic, objectNames[i])
		} else {
			addressesData = append(addressesData, result.Data)
		}
	}

	if len(problematic) != 0 {
		fm.log.Error("There were errors when fetching addresses.", "adom", adom, "problematic", problematic)
		return nil, fmt.Errorf("there were errors when fetching the following objects: %v", problematic)
	}

	fm.log.Debug("Successfully fetched addresses by names.", "adom", adom, "addressesCount", len(addressesData), "addressesData", addressesData)

	return addressesData, nil
}

func (fm *Client) GetPolicyByMetafield(adom, pkg, key string, value interface{}) (*GetPolicyData, error) {
	found := false
	var policy GetPolicyData

	request := GetPoliciesReq{
		Method: "get",
		Params: []GetPoliciesParams{
			{
				Fields: []string{"obj seq", "status", "policyid", "srcaddr", "dstaddr", "service", "action", "schedule", "extra info", "_last_hit"},
				Option: []string{"get meta"},
				Url:    fmt.Sprintf("/pm/config/adom/%s/pkg/%s/firewall/policy", adom, pkg),
			},
		},
		Session: fm.Session,
		Verbose: 1, //todo make this dynamic, without verbose, type will be 15 instead of "ipmask"
	}

	var response GetPoliciesResp

	err := fm.makeRequest(http.MethodPost, fm.Host, "/jsonrpc", request, &response)
	if err != nil {
		fm.log.Error("Error fetching policy by metafield.", "error", err, "adom", adom, "pkg", pkg, "key", key, "value", value)
		return nil, err
	}

	if len(response.Result) == 0 {
		fm.log.Error("Error fetching policies: no result in response.", "adom", adom, "pkg", pkg, "key", key, "value", value, "response", response)
		return nil, fmt.Errorf("error fetching policies: no result in response for ADOM '%s'", adom)
	}

	if response.Result[0].Status.Code != 0 || response.Result[0].Status.Message != "OK" {
		fm.log.Error("Error fetching policies.", "error", response.Result[0].Status.Message, "adom", adom, "pkg", pkg, "key", key, "value", value)
		return nil, fmt.Errorf("error fetching policies: %s", response.Result[0].Status.Message)
	}

	for _, result := range response.Result {
		for _, data := range result.Data {
			if metafieldValue, ok := data.Metafields[key]; ok && reflect.DeepEqual(metafieldValue, value) {
				if found {
					fm.log.Error("Multiple policies found with metafield.", "key", key, "value", value, "adom", adom)
					return nil, fmt.Errorf("multiple policies found with metafield '%s'='%v' in ADOM '%s'", key, value, adom)
				}
				found = true
				policy = data
			}
		}
	}

	if !found {
		fm.log.Debug("No policy found with metafield.", "key", key, "value", value, "adom", adom)
		return nil, fmt.Errorf("no policy found with metafield '%s'='%v' in ADOM '%s'", key, value, adom)
	}

	fm.log.Debug("Successfully fetched policy by metafield.", "adom", adom, "pkg", pkg, "key", key, "value", value, "policyData", &policy)

	return &policy, nil
}

func (fm *Client) GetPoliciesByID(pkg, adom string, ids []int) ([]GetPolicyData, error) {
	var params []GetPoliciesParams
	var problematic []int

	for _, id := range ids {
		params = append(params, GetPoliciesParams{
			Fields: []string{"obj seq", "status", "policyid", "srcaddr", "dstaddr", "service", "action", "schedule", "extra info", "_last_hit"},
			Option: []string{"get meta"},
			Url:    fmt.Sprintf("/pm/config/adom/%s/pkg/%s/firewall/policy/%d", adom, pkg, id),
		})
	}

	request := GetPoliciesReq{
		Method:  "get",
		Params:  params,
		Session: fm.Session,
		Verbose: 1,
	}

	var response GetPolicyResp

	err := fm.makeRequest(http.MethodPost, fm.Host, "/jsonrpc", request, &response)
	if err != nil {
		fm.log.Error("Error fetching policies.", "error", err, "adom", adom, "pkg", pkg, "ids", ids)
		return nil, err
	}

	var policiesData []GetPolicyData
	for i, result := range response.Result {
		if result.Status.Code != 0 || result.Status.Message != "OK" {
			problematic = append(problematic, ids[i])
		} else {
			policiesData = append(policiesData, result.Data)
		}
	}

	if len(problematic) != 0 {
		fm.log.Error("There were errors when fetching policies.", "adom", adom, "pkg", pkg, "problematic", problematic)
		return nil, fmt.Errorf("there were errors when fetching the following policies: %v", problematic)
	}

	fm.log.Debug("Successfully fetched policies.", "adom", adom, "pkg", pkg, "policiesCount", len(policiesData), "policiesData", policiesData)

	return policiesData, nil
}

func (fm *Client) GetPoliciesByMetafield(adom, pkg, key string, values []interface{}) ([]GetPolicyData, error) {
	var policies []GetPolicyData

	request := GetPoliciesReq{
		Method: "get",
		Params: []GetPoliciesParams{
			{
				Fields: []string{"obj seq", "status", "policyid", "srcaddr", "dstaddr", "service", "action", "schedule", "extra info", "_last_hit"},
				Option: []string{"get meta"},
				Url:    fmt.Sprintf("/pm/config/adom/%s/pkg/%s/firewall/policy", adom, pkg),
			},
		},
		Session: fm.Session,
		Verbose: 1, //todo make this dynamic, without verbose, type will be 15 instead of "ipmask"
	}

	var response GetPoliciesResp

	err := fm.makeRequest(http.MethodPost, fm.Host, "/jsonrpc", request, &response)
	if err != nil {
		fm.log.Error("Error fetching policies by metafield.", "error", err, "adom", adom, "pkg", pkg, "key", key, "values", values)
		return nil, err
	}

	//todo make sure this is consistent with others requests

	if len(response.Result) == 0 {
		fm.log.Error("Error fetching policies: no result in response.", "adom", adom, "pkg", pkg, "key", key, "values", values, "response", response)
		return nil, fmt.Errorf("error fetching policies: no result in response for ADOM '%s'", adom)
	}

	if response.Result[0].Status.Code != 0 || response.Result[0].Status.Message != "OK" {
		fm.log.Error("Error fetching policies.", "error", response.Result[0].Status.Message, "adom", adom, "pkg", pkg, "key", key, "values", values)
		return nil, fmt.Errorf("error fetching policies: %s", response.Result[0].Status.Message)
	}

	for _, value := range values {
		found := false
		for _, result := range response.Result {
			for _, data := range result.Data {
				if metafieldValue, ok := data.Metafields[key]; ok && reflect.DeepEqual(metafieldValue, value) {
					if found {
						fm.log.Error("Multiple policies found with metafield.", "key", key, "value", value, "adom", adom)
						return nil, fmt.Errorf("multiple policies found with metafield '%s'='%v' in ADOM '%s'", key, value, adom)
					}

					policies = append(policies, data)
					fm.log.Debug("Found policy by metafield.", "adom", adom, "pkg", pkg, "key", key, "value", value, "policyData", data)
					found = true
				}
			}
		}
	}

	if len(policies) != len(values) {
		fm.log.Error("Found policies by metafield, but count does not match expected.", "adom", adom, "pkg", pkg, "key", key, "values", values, "foundCount", len(policies), "expectedCount", len(values))
		return nil, fmt.Errorf("found %d policies corresponding the following '%s' metafield values: '%v', but expected %d in ADOM '%s'", len(policies), key, values, len(values), adom)
	}

	fm.log.Debug("Successfully fetched policies by metafield.", "adom", adom, "pkg", pkg, "key", key, "values", values, "policiesCount", len(policies))
	return policies, nil
}

func (fm *Client) GetServiceByMetafield(adom, key string, value interface{}) (*GetServiceData, error) {
	found := false
	var service GetServiceData

	request := GetServicesReq{
		Method: "get",
		Params: []GetServicesParams{
			{
				Fields: []string{"name", "protocol", "tcp-portrange", "udp-portrange"},
				Option: []string{"get meta"},
				Url:    fmt.Sprintf("/pm/config/adom/%s/obj/firewall/service/custom", adom),
			},
		},
		Session: fm.Session,
		Verbose: 1,
	}

	var response GetServicesResp

	err := fm.makeRequest(http.MethodPost, fm.Host, "/jsonrpc", request, &response)
	if err != nil {
		fm.log.Error("Error fetching services by metafield.", "error", err, "adom", adom, "key", key, "value", value)
		return nil, err
	}

	if len(response.Result) == 0 {
		fm.log.Error("Error fetching services: no result in response.", "adom", adom, "key", key, "value", value, "response", response)
		return nil, fmt.Errorf("error fetching services: no result in response for ADOM '%s'", adom)
	}

	if response.Result[0].Status.Code != 0 || response.Result[0].Status.Message != "OK" {
		fm.log.Error("Error fetching services.", "error", response.Result[0].Status.Message, "adom", adom, "key", key, "value", value)
		return nil, fmt.Errorf("error fetching services: %s", response.Result[0].Status.Message)
	}

	for _, result := range response.Result {
		for _, data := range result.Data {
			if metafieldValue, ok := data.Metafields[key]; ok && reflect.DeepEqual(metafieldValue, value) {
				if found {
					fm.log.Error("Multiple services found with metafield.", "key", key, "value", value, "adom", adom)
					return nil, fmt.Errorf("multiple services found with metafield '%s'='%v' in ADOM '%s'", key, value, adom)
				}
				found = true
				service = data
			}
		}
	}

	if !found {
		fm.log.Debug("No service found with metafield.", "key", key, "value", value, "adom", adom)
		return nil, fmt.Errorf("no service found with metafield '%s'='%v' in ADOM '%s'", key, value, adom)
	}

	fm.log.Debug("Successfully fetched service by metafield.", "adom", adom, "key", key, "value", value, "serviceData", &service)

	return &service, nil
}

func (fm *Client) GetServicesByMetafield(adom, key string, values []interface{}) ([]GetServiceData, error) {
	var services []GetServiceData

	request := GetServicesReq{
		Method: "get",
		Params: []GetServicesParams{
			{
				Fields: []string{"name", "tcp-portrange", "udp-portrange"},
				Option: []string{"get meta"},
				Url:    fmt.Sprintf("/pm/config/adom/%s/obj/firewall/service/custom", adom),
			},
		},
		Session: fm.Session,
		Verbose: 1, //todo make this dynamic, without verbose, type will be 15 instead of "ipmask"
	}

	var response GetServicesResp

	err := fm.makeRequest(http.MethodPost, fm.Host, "/jsonrpc", request, &response)
	if err != nil {
		fm.log.Error("Error fetching services by metafield.", "error", err, "adom", adom, "key", key, "values", values)
		return nil, err
	}

	if len(response.Result) == 0 {
		fm.log.Error("Error fetching services: no result in response.", "adom", adom, "key", key, "values", values, "response", response)
		return nil, fmt.Errorf("error fetching services: no result in response for ADOM '%s'", adom)
	}

	if response.Result[0].Status.Code != 0 || response.Result[0].Status.Message != "OK" {
		fm.log.Error("Error fetching services.", "error", response.Result[0].Status.Message, "adom", adom, "key", key, "values", values)
		return nil, fmt.Errorf("error fetching services: %s", response.Result[0].Status.Message)
	}

	for _, value := range values {
		found := false
		for _, result := range response.Result {
			for _, data := range result.Data {
				if metafieldValue, ok := data.Metafields[key]; ok && reflect.DeepEqual(metafieldValue, value) {
					if found {
						fm.log.Error("Multiple services found with metafield.", "adom", adom, "key", key, "value", value)
						return nil, fmt.Errorf("multiple services found with metafield '%s'='%v' in ADOM '%s'", key, value, adom)
					}

					services = append(services, data)
					fm.log.Debug("Found service by metafield.", "adom", adom, "key", key, "value", value, "serviceData", data)
					found = true
				}
			}
		}
	}

	if len(services) != len(values) {
		fm.log.Warn("Found services by metafield, but count does not match expected.", "adom", adom, "key", key, "values", values, "foundCount", len(services), "expectedCount", len(values))
		return nil, fmt.Errorf("found %d services corresponding the following '%s' metafield values: '%v', but expected %d in ADOM '%s'", len(services), key, values, len(values), adom)
	}

	fm.log.Debug("Successfully fetched services by metafield.", "adom", adom, "key", key, "values", values, "servicesCount", len(services))

	return services, nil
}

func (fm *Client) GetServicesByName(serviceNames []string, adom string) ([]GetServiceData, error) {
	var params []GetServicesParams
	var problematic []string

	for _, serviceName := range serviceNames {
		params = append(params, GetServicesParams{
			Fields: []string{"name", "tcp-portrange", "udp-portrange"},
			Url:    fmt.Sprintf("/pm/config/adom/%s/obj/firewall/service/custom/%s", adom, serviceName),
		})
	}

	request := GetServicesReq{
		Method:  "get",
		Params:  params,
		Session: fm.Session,
		Verbose: 1, //todo make this dynamic, without verbose, type will be 15 instead of "ipmask"
	}

	var response GetServiceResp

	err := fm.makeRequest(http.MethodPost, fm.Host, "/jsonrpc", request, &response)
	if err != nil {
		fm.log.Error("Error fetching services by names.", "error", err, "adom", adom, "serviceNames", serviceNames)
		return nil, err
	}

	var servicesData []GetServiceData
	for i, result := range response.Result {
		if result.Status.Code != 0 || result.Status.Message != "OK" {
			problematic = append(problematic, serviceNames[i])
		} else {
			servicesData = append(servicesData, result.Data)
		}
	}

	if len(problematic) != 0 {
		fm.log.Error("There were errors when fetching services.", "adom", adom, "problematic", problematic)
		return nil, fmt.Errorf("there were errors when fetching the following services: %v", problematic)
	}

	fm.log.Debug("Successfully fetched services by names.", "adom", adom, "servicesCount", len(servicesData), "servicesData", servicesData)

	return servicesData, nil
}

func (fm *Client) GetServiceByNamePortAndProtocol(adom, name, protocol string, minPort, maxPort int) (*GetServiceData, error) {
	request := GetServicesReq{
		Method: "get",
		Params: []GetServicesParams{
			{
				Fields: []string{"name", "protocol", "tcp-portrange", "udp-portrange"},
				Option: []string{"get meta"},
				Url:    fmt.Sprintf("/pm/config/adom/%s/obj/firewall/service/custom/%s", adom, name),
			},
		},
		Session: fm.Session,
		Verbose: 1, //todo make this dynamic, without verbose, type will be 15 instead of "ipmask"
	}

	var response GetServiceResp

	err := fm.makeRequest(http.MethodPost, fm.Host, "/jsonrpc", request, &response)
	if err != nil {
		fm.log.Error("Error fetching service by name, port and protocol.", "error", err, "adom", adom, "name", name, "protocol", protocol, "minPort", minPort, "maxPort", maxPort)
		return nil, err
	}

	if len(response.Result) > 1 {
		fm.log.Error("Multiple services found for name.", "name", name, "count", len(response.Result))
		return nil, fmt.Errorf("multiple services found for '%s': %d results", name, len(response.Result))
	}

	if response.Result[0].Status.Code == -3 || response.Result[0].Status.Message == "Object does not exist" {
		fm.log.Debug("Service not found in ADOM.", "adom", adom, "name", name)
		return nil, fmt.Errorf("service '%s' not found in ADOM '%s'", name, adom)
	}

	if response.Result[0].Status.Code != 0 || response.Result[0].Status.Message != "OK" {
		fm.log.Error("Error fetching service by name, port and protocol.", "error", response.Result[0].Status.Message, "adom", adom, "name", name, "protocol", protocol, "minPort", minPort, "maxPort", maxPort)
		return nil, fmt.Errorf("error fetching service '%s': %s", name, response.Result[0].Status.Message)
	}

	if protocol == "tcp" {
		if response.Result[0].Data.Protocol != "TCP/UDP/SCTP" || len(response.Result[0].Data.TCPPortRange) != 1 || response.Result[0].Data.TCPPortRange[0] != strconv.Itoa(minPort)+"-"+strconv.Itoa(maxPort) {
			fm.log.Error("Service does not match TCP port range.", "adom", adom, "name", name, "minPort", minPort, "maxPort", maxPort)
			return nil, fmt.Errorf("service '%s' does not match TCP port range '%d'-'%d' in ADOM '%s'", name, minPort, maxPort, adom)
		}
	} else if protocol == "udp" {
		if response.Result[0].Data.Protocol != "TCP/UDP/SCTP" || len(response.Result[0].Data.UDPPortRange) != 1 || response.Result[0].Data.UDPPortRange[0] != strconv.Itoa(minPort)+"-"+strconv.Itoa(maxPort) {
			fm.log.Error("Service does not match UDP port range.", "adom", adom, "name", name, "minPort", minPort, "maxPort", maxPort)
			return nil, fmt.Errorf("service '%s' does not match UDP port range '%d'-'%d' in ADOM '%s'", name, minPort, maxPort, adom)
		}
	} else {
		fm.log.Error("Invalid protocol for service.", "adom", adom, "name", name, "protocol", protocol)
		return nil, fmt.Errorf("invalid protocol '%s' for service '%s' in ADOM '%s'", protocol, name, adom)
	}

	return &response.Result[0].Data, nil
}

func (fm *Client) Lock(adom string) error {
	request := LockReq{
		Method: "exec",
		Params: []UrlParam{
			{
				Url: fmt.Sprintf("/dvmdb/adom/%s/workspace/lock", adom),
			},
		},
		Session: fm.Session,
	}

	var response LockResp

	err := fm.makeRequest(http.MethodPost, fm.Host, "/jsonrpc", request, &response)
	if err != nil {
		fm.log.Error("Error locking ADOM.", "error", err, "adom", adom)
		return err
	}

	if len(response.Result) == 0 || response.Result[0].Status.Code != 0 || response.Result[0].Status.Message != "OK" {
		fm.log.Error("Failed to lock ADOM.", "error", response.Result[0].Status.Message, "adom", adom)
		return fmt.Errorf("failed to lock ADOM: %s", response.Result[0].Status.Message)
	}

	fm.log.Debug("Successfully locked ADOM.", "adom", adom)
	return nil
}

func (fm *Client) Login() error {

	if fm.User == "" || fm.Password == "" {
		fm.log.Error("Username or password is empty.")
		return errors.New("username or password is empty")
	}

	request := LoginRequest{
		Method: "exec",
		Params: []LoginParams{
			{
				Data: LoginData{
					User:     fm.User,
					Password: fm.Password,
				},
				Url: "/sys/login/user",
			},
		},
	}

	var response LoginResp

	err := fm.makeRequest(http.MethodPost, fm.Host, "/jsonrpc", request, &response)
	if err != nil {
		fm.log.Error("Error logging in.", "error", err, "host", fm.Host, "user", fm.User)
		return err
	}

	//check if response result exists

	if len(response.Result) == 0 || response.Result[0].Status.Code != 0 || response.Result[0].Status.Message != "OK" {
		fm.log.Error("Login failed.", "error", response.Result[0].Status.Message, "host", fm.Host, "user", fm.User)
		return errors.New("login failed")
	}

	fm.Session = response.Session

	fm.log.Debug("Login successful.", "host", fm.Host, "user", fm.User, "session", fm.Session)

	return nil
}

func (fm *Client) Logout() error {
	request := LogoutReq{
		Method: "exec",
		Params: []LogoutParams{
			{
				Url: "/sys/logout",
			},
		},
		Session: fm.Session,
	}

	var response LoginResp

	err := fm.makeRequest(http.MethodPost, fm.Host, "/jsonrpc", request, &response)
	if err != nil {
		fm.log.Error("Error logging out.", "error", err, "host", fm.Host, "user", fm.User)
		return err
	}

	if len(response.Result) == 0 || response.Result[0].Status.Code != 0 || response.Result[0].Status.Message != "OK" {
		fm.log.Error("Logout failed.", "error", response.Result[0].Status.Message, "host", fm.Host, "user", fm.User)
		return errors.New("logout failed")
	}

	fm.Session = ""

	fm.log.Debug("Logout successful.", "host", fm.Host, "user", fm.User)

	return nil
}

func (fm *Client) Unlock(adom string) error {
	request := UnlockReq{
		Method: "exec",
		Params: []UnlockParams{
			{
				Url: fmt.Sprintf("/dvmdb/adom/%s/workspace/unlock", adom),
			},
		},
		Session: fm.Session,
	}

	var response UnlockResp

	err := fm.makeRequest(http.MethodPost, fm.Host, "/jsonrpc", request, &response)
	if err != nil {
		fm.log.Error("Error unlocking ADOM.", "error", err, "adom", adom)
		return err
	}

	if len(response.Result) == 0 || response.Result[0].Status.Code != 0 || response.Result[0].Status.Message != "OK" {
		fm.log.Error("Failed to unlock ADOM.", "error", response.Result[0].Status.Message, "adom", adom)
		return fmt.Errorf("failed to unlock ADOM: %s", response.Result[0].Status.Message)
	}

	fm.log.Debug("Successfully unlocked ADOM.", "adom", adom)

	return nil
}

func (fm *Client) UpdateServiceWithMetafields(adom, name string, metafields map[string]interface{}) error {
	request := ServiceReq{
		Method: "set",
		Params: []ServiceReqParams{
			{
				Data: []ServiceReqData{
					{
						Name:       name,
						Metafields: metafields,
					},
				},
				Url: fmt.Sprintf("/pm/config/adom/%s/obj/firewall/service/custom", adom),
			},
		},
		Session: fm.Session,
	}

	var response ServiceResp

	err := fm.makeRequest(http.MethodPost, fm.Host, "/jsonrpc", request, &response)
	if err != nil {
		fm.log.Error("Error updating service with metafields.", "error", err, "adom", adom, "name", name)
		return err
	}

	if len(response.Result) == 0 || response.Result[0].Status.Code != 0 || response.Result[0].Status.Message != "OK" {
		fm.log.Error("Failed to update service with metafields.", "error", response.Result[0].Status.Message, "adom", adom, "name", name)
		return fmt.Errorf("failed to update service with metafields: %s", response.Result[0].Status.Message)
	}

	fm.log.Debug("Successfully updated service with metafields.", "adom", adom, "name", name, "metafields", metafields)

	return nil
}

func (fm *Client) UpdateSubnetAddressWithMetafields(adom, name string, metafields map[string]interface{}) error {
	request := SubnetAddressReq{
		Method: "set",
		Params: []SubnetAddressReqParams{
			{
				Data: []SubnetAddressReqData{
					{
						Name:       name,
						Metafields: metafields,
					},
				},
				Url: fmt.Sprintf("/pm/config/adom/%s/obj/firewall/address", adom),
			},
		},
		Session: fm.Session,
	}

	var response SubnetAddressResp

	err := fm.makeRequest(http.MethodPost, fm.Host, "/jsonrpc", request, &response)
	if err != nil {
		fm.log.Error("Error updating address with metafields.", "error", err, "adom", adom, "name", name)
		return err
	}

	if len(response.Result) == 0 || response.Result[0].Status.Code != 0 || response.Result[0].Status.Message != "OK" {
		fm.log.Error("Failed to update address with metafields.", "error", response.Result[0].Status.Message, "adom", adom, "name", name)
		return fmt.Errorf("failed to update address with metafields: %s", response.Result[0].Status.Message)
	}

	fm.log.Debug("Successfully updated address with metafields.", "adom", adom, "name", name, "metafields", metafields)

	return nil
}

// todo this function is quite inefficient - the redacting is not working - do not commit before finishing this
// maybe this has to be an internal function
// doesn't make sense to use fmt in some cases and the arguement in other cases, either this belongs to the client or the caller
//
//	and if it belongs to the client, then we can remove some of the arguments
func (fm *Client) makeRequest(method, host, endpoint string, reqBody interface{}, respBody interface{}) error {
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		fm.log.Error("Error encoding JSON request body.", "error", err, "requestBody", reqBody)
		return err
	}

	req, err := http.NewRequest(method, host+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		fm.log.Error("Error creating HTTP request.", "error", err, "method", method, "url", host+endpoint, "requestBody", reqBody)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	if fm.Key != "" {
		req.Header.Set("Authorization", "Bearer "+fm.Key)
	}

	/*redactedReqBody, err := redactSensitiveData(jsonData, reqBody)
	if err != nil {
		return fmt.Errorf("error redacting sensitive data: %v", err)
	}*
	fm.infoLog.Printf("Request: %s %s %s\n", host+method, endpoint, reqBody)*/
	fm.log.Debug("Making request.", "method", method, "url", host+endpoint, "requestBody", reqBody)

	client := &http.Client{}

	durationStart := time.Now()

	resp, err := client.Do(req)
	if err != nil {
		fm.log.Error("Error making request.", "error", err, "method", method, "url", host+endpoint, "requestBody", jsonData)
		return err
	}
	defer resp.Body.Close()

	durationEnd := time.Since(durationStart).Milliseconds()

	if resp.Body == nil {
		fm.log.Error("Response body is empty.", "method", method, "url", host+endpoint, "requestBody", jsonData)
		return errors.New("response body is empty")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fm.log.Error("Error reading response body.", "error", err, "method", method, "url", host+endpoint, "requestBody", jsonData)
		return fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode >= 400 && resp.StatusCode < 600 {
		fm.log.Error("Error response from server.", "statusCode", resp.StatusCode, "method", method, "url", host+endpoint, "responseBody", string(body))
		return fmt.Errorf("error response: %s", body)
	}

	/*redactedRespBody, err := redactSensitiveData(body, respBody)
	if err != nil {
		return fmt.Errorf("error redacting sensitive data: %v", err)
	}
	fm.infoLog.Printf("Response (%d ms): %s %s %s\n", durationEnd, resp.Status, host+endpoint, redactedRespBody)*/
	fm.log.Debug("Received response.", "status", resp.Status, "url", host+endpoint, "durationMs", durationEnd, "responseBody", string(body))

	err = json.Unmarshal(body, respBody)
	if err != nil {
		fm.log.Error("Error unmarshalling response body.", "error", err, "responseBody", string(body))
		return fmt.Errorf("error unmarshalling response: %v", err)
	}

	return nil
}

// todo restructure this function
func redactSensitiveData(body []byte, v interface{}) ([]byte, error) {
	err := json.Unmarshal(body, v)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	val := reflect.ValueOf(v).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		if tag, ok := fieldType.Tag.Lookup("redact"); ok && tag == "true" {
			if field.Kind() == reflect.String {
				field.SetString("[REDACTED]")
			}
		}
	}

	redactedBody, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("error marshalling JSON: %v", err)
	}

	return redactedBody, nil
}
