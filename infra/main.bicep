// Copyright 2024 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

targetScope = 'subscription'

@minLength(1)
@maxLength(64)
@description('Name of the environment that can be used as part of naming resource convention')
param environmentName string

@minLength(1)
@description('Primary location for all resources')
param location string

@description('User principal ID')
param principalId string

@description('The vault name; default is a unique string based on the resource group ID')
param vaultName string = ''

var projectName = 'go-dotazure'
var tags = {
  'azd-env-name': environmentName
}

resource rg 'Microsoft.Resources/resourceGroups@2024-11-01' = {
  name: 'rg-${projectName}-${environmentName}'
  location: location
  tags: tags
}

module resources 'resources.bicep' = {
  name: 'resources'
  scope: rg
  params: {
    environmentName: environmentName
    location: location
    principalId: principalId
    vaultName: vaultName
  }
}

output AZURE_PRINCIPAL_ID string = resources.outputs.AZURE_PRINCIPAL_ID
output AZURE_RESOURCE_GROUP string = rg.name
output AZURE_KEYVAULT_NAME string = resources.outputs.AZURE_KEYVAULT_NAME
output AZURE_KEYVAULT_URL string = resources.outputs.AZURE_KEYVAULT_URL
