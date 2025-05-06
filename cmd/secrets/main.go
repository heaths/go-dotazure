// Copyright 2025 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/heaths/go-dotazure"
)

func main() {
	if err := dotazure.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		panic(err)
	}

	vaultURL, _ := os.LookupEnv("AZURE_KEYVAULT_URL")
	if vaultURL == "" {
		panic("AZURE_KEYVAULT_URL not set")
	}
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(fmt.Errorf("create credential: %w", err))
	}
	client, err := azsecrets.NewClient(vaultURL, credential, nil)
	if err != nil {
		panic(fmt.Errorf("create secret client: %w", err))
	}
	secret, err := client.GetSecret(context.Background(), "my-secret", "", nil)
	if err != nil {
		panic(fmt.Errorf("get secret: %w", err))
	}
	if secret.Value != nil {
		fmt.Println(*secret.Value)
	}
}
