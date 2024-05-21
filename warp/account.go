package warp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

var (
	identityFile = "wgcf-identity.json"
	profileFile  = "wgcf-profile.ini"
)

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
	i, err := Register(publicKey)
	if err != nil {
		return Identity{}, err
	}

	if license != "" {
		l.Info("updating account license key")
		ac, err := UpdateAccount(i.Token, i.ID, license)
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
