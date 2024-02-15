package warp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	apiVersion    = "v0a1922"
	apiURL        = "https://api.cloudflareclient.com"
	regURL        = apiURL + "/" + apiVersion + "/reg"
	_identityFile = "wgcf-identity.json"
	_profileFile  = "wgcf-profile.ini"
)

var (
	identityFile = "wgcf-identity.json"
	profileFile  = "wgcf-profile.ini"
	dnsAddresses = []string{"8.8.8.8", "8.8.4.4"}
	dc           = 0
)

var defaultHeaders = makeDefaultHeaders()
var client = makeClient()

type AccountData struct {
	AccountID   string `json:"account_id"`
	AccessToken string `json:"access_token"`
	PrivateKey  string `json:"private_key"`
	LicenseKey  string `json:"license_key"`
}

type ConfigurationData struct {
	LocalAddressIPv4    string `json:"local_address_ipv4"`
	LocalAddressIPv6    string `json:"local_address_ipv6"`
	EndpointAddressHost string `json:"endpoint_address_host"`
	EndpointAddressIPv4 string `json:"endpoint_address_ipv4"`
	EndpointAddressIPv6 string `json:"endpoint_address_ipv6"`
	EndpointPublicKey   string `json:"endpoint_public_key"`
	WarpEnabled         bool   `json:"warp_enabled"`
	AccountType         string `json:"account_type"`
	WarpPlusEnabled     bool   `json:"warp_plus_enabled"`
	LicenseKeyUpdated   bool   `json:"license_key_updated"`
}

func makeDefaultHeaders() map[string]string {
	return map[string]string{
		"User-Agent":        "okhttp/3.12.1",
		"CF-Client-Version": "a-6.3-1922",
	}
}

func makeClient() *http.Client {
	// Create a custom dialer using the TLS config
	plainDialer := &net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 5 * time.Second,
	}
	tlsDialer := Dialer{}
	// Create a custom HTTP transport
	transport := &http.Transport{
		DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return tlsDialer.TLSDial(plainDialer, network, addr)
		},
	}

	// Create a custom HTTP client using the transport
	return &http.Client{
		Transport: transport,
		// Other client configurations can be added here
	}
}

func MergeMaps(maps ...map[string]string) map[string]string {
	out := make(map[string]string)

	for _, m := range maps {
		for k, v := range m {
			out[k] = v
		}
	}

	return out
}

func getConfigURL(accountID string) string {
	return fmt.Sprintf("%s/%s", regURL, accountID)
}

func getAccountURL(accountID string) string {
	return fmt.Sprintf("%s/account", getConfigURL(accountID))
}

func getDevicesURL(accountID string) string {
	return fmt.Sprintf("%s/devices", getAccountURL(accountID))
}

func getAccountRegURL(accountID, deviceToken string) string {
	return fmt.Sprintf("%s/reg/%s", getAccountURL(accountID), deviceToken)
}

func getTimestamp() string {
	timestamp := time.Now().Format(time.RFC3339Nano)
	return timestamp
}

func genKeyPair() (string, string, error) {
	// Generate private key
	priv, err := GeneratePrivateKey()
	if err != nil {
		fmt.Println("Error generating private key:", err)
		return "", "", err
	}
	privateKey := priv.String()
	publicKey := priv.PublicKey().String()
	return privateKey, publicKey, nil
}

func doRegister() (*AccountData, error) {
	timestamp := getTimestamp()
	privateKey, publicKey, err := genKeyPair()
	if err != nil {
		return nil, err
	}
	data := map[string]interface{}{
		"install_id": "",
		"fcm_token":  "",
		"tos":        timestamp,
		"key":        publicKey,
		"type":       "Android",
		"model":      "PC",
		"locale":     "en_US",
	}

	headers := map[string]string{
		"Content-Type": "application/json; charset=UTF-8",
	}

	jsonBody, _ := json.Marshal(data)

	req, err := http.NewRequest("POST", regURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	// Set headers
	for k, v := range MergeMaps(defaultHeaders, headers) {
		req.Header.Set(k, v)
	}

	// Create HTTP client and execute request
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("sending request to remote server", err)
		return nil, err
	}

	// convert response to byte array
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("reading response body", err)
		return nil, err
	}

	var rspData interface{}

	err = json.Unmarshal(responseData, &rspData)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	m := rspData.(map[string]interface{})

	return &AccountData{
		AccountID:   m["id"].(string),
		AccessToken: m["token"].(string),
		PrivateKey:  privateKey,
		LicenseKey:  m["account"].(map[string]interface{})["license"].(string),
	}, nil
}

func saveIdentity(accountData *AccountData, identityPath string) error {
	file, err := os.Create(identityPath)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(accountData)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	return file.Close()
}

func loadIdentity(identityPath string) (accountData *AccountData, err error) {
	file, err := os.Open(identityPath)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			fmt.Println("Error:", err)
		}
		return
	}(file)

	accountData = &AccountData{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&accountData)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return accountData, nil
}

func enableWarp(accountData *AccountData) error {
	data := map[string]interface{}{
		"warp_enabled": true,
	}

	jsonData, _ := json.Marshal(data)

	url := getConfigURL(accountData.AccountID)

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// Set headers
	headers := map[string]string{
		"Authorization": "Bearer " + accountData.AccessToken,
		"Content-Type":  "application/json; charset=UTF-8",
	}

	for k, v := range MergeMaps(defaultHeaders, headers) {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error enabling WARP, status %d", resp.StatusCode)
	}

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}

	if !response["warp_enabled"].(bool) {
		return errors.New("warp not enabled")
	}

	return nil
}

func getServerConf(accountData *AccountData) (*ConfigurationData, error) {

	req, err := http.NewRequest("GET", getConfigURL(accountData.AccountID), nil)
	if err != nil {
		return nil, err
	}

	// Set headers
	headers := map[string]string{
		"Authorization": "Bearer " + accountData.AccessToken,
		"Content-Type":  "application/json; charset=UTF-8",
	}

	for k, v := range MergeMaps(defaultHeaders, headers) {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting config, status %d", resp.StatusCode)
	}

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	addresses := response["config"].(map[string]interface{})["interface"].(map[string]interface{})["addresses"]
	lv4 := addresses.(map[string]interface{})["v4"].(string)
	lv6 := addresses.(map[string]interface{})["v6"].(string)

	peer := response["config"].(map[string]interface{})["peers"].([]interface{})[0].(map[string]interface{})
	publicKey := peer["public_key"].(string)

	endpoint := peer["endpoint"].(map[string]interface{})
	host := endpoint["host"].(string)
	v4 := endpoint["v4"].(string)
	v6 := endpoint["v6"].(string)

	account, ok := response["account"].(map[string]interface{})
	if !ok {
		account = make(map[string]interface{})
	}

	warpEnabled := response["warp_enabled"].(bool)

	return &ConfigurationData{
		LocalAddressIPv4:    lv4,
		LocalAddressIPv6:    lv6,
		EndpointAddressHost: host,
		EndpointAddressIPv4: v4,
		EndpointAddressIPv6: v6,
		EndpointPublicKey:   publicKey,
		WarpEnabled:         warpEnabled,
		AccountType:         account["account_type"].(string),
		WarpPlusEnabled:     account["warp_plus"].(bool),
		LicenseKeyUpdated:   false, // omit for brevity
	}, nil
}

func updateLicenseKey(accountData *AccountData, confData *ConfigurationData) (bool, error) {

	if confData.AccountType == "free" && accountData.LicenseKey != "" {

		data := map[string]interface{}{
			"license": accountData.LicenseKey,
		}

		jsonData, _ := json.Marshal(data)

		url := getAccountURL(accountData.AccountID)

		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
		if err != nil {
			return false, err
		}

		// Set headers
		headers := map[string]string{
			"Authorization": "Bearer " + accountData.AccessToken,
			"Content-Type":  "application/json; charset=UTF-8",
		}

		for k, v := range MergeMaps(defaultHeaders, headers) {
			req.Header.Set(k, v)
		}

		resp, err := client.Do(req)
		if err != nil {
			return false, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			s, _ := io.ReadAll(resp.Body)
			return false, fmt.Errorf("activation error, status %d %s", resp.StatusCode, string(s))
		}

		var activationResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&activationResp)
		if err != nil {
			return false, err
		}

		return activationResp["warp_plus"].(bool), nil

	} else if confData.AccountType == "unlimited" {
		return true, nil
	}

	return false, nil
}

func getDeviceActive(accountData *AccountData) (bool, error) {

	req, err := http.NewRequest("GET", getDevicesURL(accountData.AccountID), nil)
	if err != nil {
		return false, err
	}

	// Set headers
	headers := map[string]string{
		"Authorization": "Bearer " + accountData.AccessToken,
		"Accept":        "application/json",
	}

	for k, v := range MergeMaps(defaultHeaders, headers) {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("error getting devices, status %d", resp.StatusCode)
	}

	var devices []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&devices)

	for _, d := range devices {
		if d["id"] == accountData.AccountID {
			active := d["active"].(bool)
			return active, nil
		}
	}

	return false, nil
}

func setDeviceActive(accountData *AccountData, status bool) (bool, error) {

	data := map[string]interface{}{
		"active": status,
	}

	jsonData, _ := json.Marshal(data)

	url := getAccountRegURL(accountData.AccountID, accountData.AccountID)

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return false, err
	}

	// Set headers
	headers := map[string]string{
		"Authorization": "Bearer " + accountData.AccessToken,
		"Accept":        "application/json",
	}

	for k, v := range MergeMaps(defaultHeaders, headers) {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("error setting active status, status %d", resp.StatusCode)
	}

	var devices []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&devices)

	for _, d := range devices {
		if d["id"] == accountData.AccountID {
			return d["active"].(bool), nil
		}
	}

	return false, nil
}

func getWireguardConfig(privateKey, address1, address2, publicKey, endpoint string) string {

	var buffer bytes.Buffer

	buffer.WriteString("[Interface]\n")
	buffer.WriteString(fmt.Sprintf("PrivateKey = %s\n", privateKey))
	buffer.WriteString(fmt.Sprintf("DNS = %s\n", dnsAddresses[dc%len(dnsAddresses)]))
	dc++
	buffer.WriteString(fmt.Sprintf("Address = %s\n", address1+"/24"))
	buffer.WriteString(fmt.Sprintf("Address = %s\n", address2+"/128"))

	buffer.WriteString("[Peer]\n")
	buffer.WriteString(fmt.Sprintf("PublicKey = %s\n", publicKey))
	buffer.WriteString("AllowedIPs = 0.0.0.0/0\n")
	buffer.WriteString("AllowedIPs = ::/0\n")
	buffer.WriteString(fmt.Sprintf("Endpoint = %s\n", endpoint))

	return buffer.String()
}

func createConf(accountData *AccountData, confData *ConfigurationData) error {

	config := getWireguardConfig(accountData.PrivateKey, confData.LocalAddressIPv4,
		confData.LocalAddressIPv6, confData.EndpointPublicKey, confData.EndpointAddressHost)

	return os.WriteFile(profileFile, []byte(config), 0600)
}

func LoadOrCreateIdentity(license string) error {
	var accountData *AccountData

	if _, err := os.Stat(identityFile); os.IsNotExist(err) {
		fmt.Println("Creating new identity...")
		accountData, err = doRegister()
		if err != nil {
			return err
		}
		accountData.LicenseKey = license
		saveIdentity(accountData, identityFile)
	} else {
		fmt.Println("Loading existing identity...")
		accountData, err = loadIdentity(identityFile)
		if err != nil {
			return err
		}
	}

	fmt.Println("Getting configuration...")
	confData, err := getServerConf(accountData)
	if err != nil {
		return err
	}

	// updating license key
	fmt.Println("Updating account license key...")
	result, err := updateLicenseKey(accountData, confData)
	if err != nil {
		return err
	}
	if result {
		confData, err = getServerConf(accountData)
		if err != nil {
			return err
		}
	}

	deviceStatus, err := getDeviceActive(accountData)
	if err != nil {
		return err
	}
	if !deviceStatus {
		fmt.Println("This device is not registered to the account!")
	}

	if confData.WarpPlusEnabled && !deviceStatus {
		fmt.Println("Enabling device...")
		deviceStatus, err = setDeviceActive(accountData, true)
	}

	if !confData.WarpEnabled {
		fmt.Println("Enabling Warp...")
		err := enableWarp(accountData)
		if err != nil {
			return err
		}
		confData.WarpEnabled = true
	}

	fmt.Printf("Warp+ enabled: %t\n", confData.WarpPlusEnabled)
	fmt.Printf("Device activated: %t\n", deviceStatus)
	fmt.Printf("Account type: %s\n", confData.AccountType)
	fmt.Printf("Warp+ enabled: %t\n", confData.WarpPlusEnabled)

	fmt.Println("Creating WireGuard configuration...")
	err = createConf(accountData, confData)
	if err != nil {
		return fmt.Errorf("unable to enable write config file, Error: %v", err.Error())
	}

	fmt.Println("All done! Find your files here:")
	fmt.Println(filepath.Abs(identityFile))
	fmt.Println(filepath.Abs(profileFile))
	return nil
}

func fileExist(f string) bool {
	if _, err := os.Stat(f); os.IsNotExist(err) {
		return false
	}
	return true
}
func removeFile(f string) {
	if fileExist(f) {
		e := os.Remove(f)
		if e != nil {
			log.Fatal(e)
		}
	}
}

func UpdatePath(path string) {
	identityFile = path + "/" + _identityFile
	profileFile = path + "/" + _profileFile
}

func CheckProfileExists(license string) bool {
	isOk := true
	if !fileExist(identityFile) || !fileExist(profileFile) {
		isOk = false
	}

	ad := &AccountData{} // Read errors caught by unmarshal
	if isOk {
		fileBytes, _ := os.ReadFile(identityFile)
		err := json.Unmarshal(fileBytes, ad)
		if err != nil {
			isOk = false
		} else if license != "notset" && ad.LicenseKey != license {
			isOk = false
		}
	}
	if !isOk {
		removeFile(profileFile)
		removeFile(identityFile)
	}
	return isOk
}

func RemoveDevice(account AccountData) error {

	headers := map[string]string{
		"Content-Type":      "application/json",
		"User-Agent":        "okhttp/3.12.1",
		"CF-Client-Version": "a-6.30-3596",
		"Authorization":     "Bearer " + account.AccessToken,
	}

	req, err := http.NewRequest("DELETE", "https://api.cloudflareclient.com/v0a3596/reg/"+account.AccountID, nil)
	if err != nil {
		return err
	}

	// Set headers
	for k, v := range MergeMaps(defaultHeaders, headers) {
		req.Header.Set(k, v)
	}

	// Create HTTP client and execute request
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("sending request to remote server", err)
		return err
	}

	if response.StatusCode != 204 {
		return fmt.Errorf("Error in deleting account %d %s", response.StatusCode, response.Status)
	}

	return nil
}
