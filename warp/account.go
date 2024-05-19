package warp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	apiVersion = "v0a3596"
	apiURL     = "https://api.cloudflareclient.com"
	regURL     = apiURL + "/" + apiVersion + "/reg"
)

var (
	identityFile = "wgcf-identity.json"
	profileFile  = "wgcf-profile.ini"
)

var (
	defaultHeaders = makeDefaultHeaders()
	client         = makeClient()
)

type IdentityAccount struct {
	Created                  string `json:"created"`
	Updated                  string `json:"updated"`
	License                  string `json:"license"`
	PremiumData              int64  `json:"premium_data"`
	WarpPlus                 bool   `json:"warp_plus"`
	AccountType              string `json:"account_type"`
	ReferralRenewalCountdown int64  `json:"referral_renewal_countdown"`
	Role                     string `json:"role"`
	ID                       string `json:"id"`
	Quota                    int64  `json:"quota"`
	Usage                    int64  `json:"usage"`
	ReferralCount            int64  `json:"referral_count"`
	TTL                      string `json:"ttl"`
}

type IdentityConfigPeerEndpoint struct {
	V4    string   `json:"v4"`
	V6    string   `json:"v6"`
	Host  string   `json:"host"`
	Ports []uint16 `json:"ports"`
}

type IdentityConfigPeer struct {
	PublicKey string                     `json:"public_key"`
	Endpoint  IdentityConfigPeerEndpoint `json:"endpoint"`
}

type IdentityConfigInterfaceAddresses struct {
	V4 string `json:"v4"`
	V6 string `json:"v6"`
}

type IdentityConfigInterface struct {
	Addresses IdentityConfigInterfaceAddresses `json:"addresses"`
}
type IdentityConfigServices struct {
	HTTPProxy string `json:"http_proxy"`
}

type IdentityConfig struct {
	Peers     []IdentityConfigPeer    `json:"peers"`
	Interface IdentityConfigInterface `json:"interface"`
	Services  IdentityConfigServices  `json:"services"`
	ClientID  string                  `json:"client_id"`
}

type Identity struct {
	PrivateKey      string          `json:"private_key"`
	Key             string          `json:"key"`
	Account         IdentityAccount `json:"account"`
	Place           int64           `json:"place"`
	FCMToken        string          `json:"fcm_token"`
	Name            string          `json:"name"`
	TOS             string          `json:"tos"`
	Locale          string          `json:"locale"`
	InstallID       string          `json:"install_id"`
	WarpEnabled     bool            `json:"warp_enabled"`
	Type            string          `json:"type"`
	Model           string          `json:"model"`
	Config          IdentityConfig  `json:"config"`
	Token           string          `json:"token"`
	Enabled         bool            `json:"enabled"`
	ID              string          `json:"id"`
	Created         string          `json:"created"`
	Updated         string          `json:"updated"`
	WaitlistEnabled bool            `json:"waitlist_enabled"`
}

func makeDefaultHeaders() map[string]string {
	return map[string]string{
		"Content-Type":      "application/json; charset=UTF-8",
		"User-Agent":        "okhttp/3.12.1",
		"CF-Client-Version": "a-6.30-3596",
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

func doRegister(publicKey string) (Identity, error) {
	data := map[string]interface{}{
		"install_id":   "",
		"fcm_token":    "",
		"tos":          time.Now().Format(time.RFC3339Nano),
		"key":          publicKey,
		"type":         "Android",
		"model":        "PC",
		"locale":       "en_US",
		"warp_enabled": true,
	}

	jsonBody, err := json.Marshal(data)
	if err != nil {
		return Identity{}, err
	}

	req, err := http.NewRequest("POST", regURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return Identity{}, err
	}

	// Set headers
	for k, v := range defaultHeaders {
		req.Header.Set(k, v)
	}

	// Create HTTP client and execute request
	resp, err := client.Do(req)
	if err != nil {
		return Identity{}, err
	}
	defer resp.Body.Close()

	// convert response to byte array
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return Identity{}, err
	}

	var rspData = Identity{}
	err = json.Unmarshal(responseData, &rspData)
	if err != nil {
		return Identity{}, err
	}

	return rspData, nil
}

func saveIdentity(a Identity, path string) error {
	file, err := os.Create(filepath.Join(path, identityFile))
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(a)
	if err != nil {
		return err
	}

	return file.Close()
}

func updateLicenseKey(accountID, accessToken, license string) (IdentityAccount, error) {
	jsonData, err := json.Marshal(map[string]string{"license": license})
	if err != nil {
		return IdentityAccount{}, err
	}

	url := fmt.Sprintf("%s/%s/account", regURL, accountID)

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return IdentityAccount{}, err
	}

	headers := defaultHeaders
	headers["Authorization"] = "Bearer " + accessToken
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return IdentityAccount{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s, err := io.ReadAll(resp.Body)
		if err != nil {
			return IdentityAccount{}, err
		}

		return IdentityAccount{}, fmt.Errorf("activation error, status %d %s", resp.StatusCode, string(s))
	}

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return IdentityAccount{}, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp1, err := client.Do(req)
	if err != nil {
		return IdentityAccount{}, err
	}
	defer resp1.Body.Close()

	if resp1.StatusCode != http.StatusOK {
		s, err := io.ReadAll(resp1.Body)
		if err != nil {
			return IdentityAccount{}, err
		}

		return IdentityAccount{}, fmt.Errorf("activation error, status %d %s", resp1.StatusCode, string(s))
	}

	var activationResp1 = IdentityAccount{}
	err = json.NewDecoder(resp1.Body).Decode(&activationResp1)
	if err != nil {
		return IdentityAccount{}, err
	}

	return activationResp1, nil
}

func createConf(i Identity, path string) error {
	var buffer bytes.Buffer

	buffer.WriteString("[Interface]\n")
	buffer.WriteString(fmt.Sprintf("PrivateKey = %s\n", i.PrivateKey))

	buffer.WriteString(fmt.Sprintf("Address = %s/24\n", i.Config.Interface.Addresses.V4))
	buffer.WriteString(fmt.Sprintf("Address = %s/128\n", i.Config.Interface.Addresses.V6))

	buffer.WriteString("[Peer]\n")
	buffer.WriteString(fmt.Sprintf("PublicKey = %s\n", i.Config.Peers[0].PublicKey))
	buffer.WriteString("AllowedIPs = 0.0.0.0/0\n")
	buffer.WriteString("AllowedIPs = ::/0\n")
	buffer.WriteString(fmt.Sprintf("Endpoint = %s\n", i.Config.Peers[0].Endpoint.Host))

	return os.WriteFile(filepath.Join(path, profileFile), buffer.Bytes(), 0o600)
}

func LoadOrCreateIdentity(l *slog.Logger, path, license string) error {
	i, err := LoadIdentity(path)
	if err != nil {
		l.Info("failed to load identity", "path", path, "error", err)
		if err := os.RemoveAll(path); err != nil {
			return err
		}
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
		i, err = CreateIdentity(l, path, license)
		if err != nil {
			return err
		}
	}

	if license != "" && i.Account.License != license {
		l.Info("license recreating identity with new license")
		if err := os.RemoveAll(path); err != nil {
			return err
		}
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
		i, err = CreateIdentity(l, path, license)
		if err != nil {
			return err
		}
	}

	err = createConf(i, path)
	if err != nil {
		return fmt.Errorf("unable to enable write config file: %w", err)
	}

	l.Info("successfully generated wireguard configuration")
	return nil
}

func LoadIdentity(path string) (Identity, error) {
	// If either of the identity or profile files doesn't exist.
	identityPath := filepath.Join(path, identityFile)
	_, err := os.Stat(identityPath)
	if err != nil {
		return Identity{}, err
	}

	profilePath := filepath.Join(path, profileFile)
	_, err = os.Stat(profilePath)
	if err != nil {
		return Identity{}, err
	}

	i := &Identity{}

	fileBytes, err := os.ReadFile(identityPath)
	if err != nil {
		return Identity{}, err
	}

	err = json.Unmarshal(fileBytes, i)
	if err != nil {
		return Identity{}, err
	}

	if len(i.Config.Peers) < 1 {
		return Identity{}, errors.New("identity contains 0 peers")
	}

	return *i, nil
}

func CreateIdentity(l *slog.Logger, path, license string) (Identity, error) {
	priv, err := GeneratePrivateKey()
	if err != nil {
		return Identity{}, err
	}

	privateKey, publicKey := priv.String(), priv.PublicKey().String()

	l.Info("creating new identity")
	i, err := doRegister(publicKey)
	if err != nil {
		return Identity{}, err
	}

	if license != "" {
		l.Info("updating account license key")
		ac, err := updateLicenseKey(i.ID, i.Token, license)
		if err != nil {
			return Identity{}, err
		}
		i.Account = ac
	}

	i.PrivateKey = privateKey

	err = saveIdentity(i, path)
	if err != nil {
		return Identity{}, err
	}

	return i, nil
}

func RemoveDevice(l *slog.Logger, accountID, accessToken string) error {
	url := fmt.Sprintf("%s/%s", regURL, accountID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	headers := defaultHeaders
	headers["Authorization"] = "Bearer " + accessToken
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Create HTTP client and execute request
	resp, err := client.Do(req)
	if err != nil {
		l.Info("sending request to remote server", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		return fmt.Errorf("error in deleting account %d %s", resp.StatusCode, resp.Status)
	}

	return nil
}
