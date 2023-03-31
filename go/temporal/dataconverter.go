package dataconverter

import (
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/vault/api"
	"go.temporal.io/sdk/converter"

	commonpb "go.temporal.io/api/common/v1"
)

type VaultTransitDataConverter struct {
	vaultClient *api.Client
	keyName     string
	encyrptPath string
	decryptPath string
}

func NewVaultTransitDataConverter(vaultAddr, vaultToken, keyName string) (*VaultTransitDataConverter, error) {
	config := &api.Config{Address: vaultAddr}
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	client.SetToken(vaultToken)

	return &VaultTransitDataConverter{
		vaultClient: client,
		keyName:     keyName,
	}, nil
}

// Converts a single payload
func (c *VaultTransitDataConverter) ToPayload(payload interface{}) (*commonpb.Payload, error) {
	// Convert to bytes
	input, err := converter.GetDefaultDataConverter().ToPayload(payload)
	if err != nil {
		return nil, fmt.Errorf("data conversion failed: %w", err)
	}

	encrypted, err := c.encrypt(input.Data)
	if err != nil {
		return nil, fmt.Errorf("encryption failed: %w", err)
	}

	// Return the encrypted payload
	return &commonpb.Payload{Data: encrypted}, nil
}

func (c *VaultTransitDataConverter) FromPayload(input []byte, valuePtr interface{}) error {
	decrypted, err := c.decrypt(input)
	if err != nil {
		return fmt.Errorf("decryption failed: %w", err)
	}

	err = converter.GetDefaultDataConverter().FromData(decrypted, valuePtr)
	if err != nil {
		return fmt.Errorf("data conversion failed: %w", err)
	}

	return nil
}

func (c *VaultTransitDataConverter) ToPayloads(payload ...interface{}) (*commonpb.Payloads, errors) {
}

func (c *VaultTransitDataConverter) encrypt(plaintext []byte) ([]byte, error) {
	encodedPlaintext := base64.StdEncoding.EncodeToString(plaintext)
	path := fmt.Sprintf("transit/encrypt/%s", c.keyName)

	secret, err := c.vaultClient.Logical().Write(path, map[string]interface{}{
		"plaintext": encodedPlaintext,
	})
	if err != nil {
		return nil, err
	}

	encodedCiphertext := secret.Data["ciphertext"].(string)
	ciphertext, err := base64.StdEncoding.DecodeString(encodedCiphertext)
	if err != nil {
		return nil, err
	}

	return ciphertext, nil
}

func (c *VaultTransitDataConverter) decrypt(ciphertext []byte) ([]byte, error) {
	encodedCiphertext := base64.StdEncoding.EncodeToString(ciphertext)
	path := fmt.Sprintf("transit/decrypt/%s", c.keyName)

	secret, err := c.vaultClient.Logical().Write(path, map[string]interface{}{
		"ciphertext": encodedCiphertext,
	})
	if err != nil {
		return nil, err
	}

	encodedPlaintext := secret.Data["plaintext"].(string)
	plaintext, err := base64.StdEncoding.DecodeString(encodedPlaintext)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
