package dataconverter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	TestVaultAddr      = "<VAULT_ADDR>"
	TestVaultToken     = "<VAULT_TOKEN>"
	TestTransitKeyName = "test-transit-key"
)

func TestDataConverter(t *testing.T) {
	// Replace <VAULT_ADDR> and <VAULT_TOKEN> with your Vault address and token
	dataConverter, err := NewVaultTransitDataConverter(TestVaultAddr, TestVaultToken, TestTransitKeyName)
	require.NoError(t, err)

	// Test data
	type TestData struct {
		Field1 string
		Field2 int
	}
	testData := TestData{
		Field1: "Test data",
		Field2: 42,
	}

	// Test ToData (encryption)
	encryptedData, err := dataConverter.ToData(testData)
	require.NoError(t, err)
	assert.NotEmpty(t, encryptedData)

	// Test FromData (decryption)
	var decryptedData TestData
	err = dataConverter.FromData(encryptedData, &decryptedData)
	require.NoError(t, err)
	assert.Equal(t, testData, decryptedData)
}
