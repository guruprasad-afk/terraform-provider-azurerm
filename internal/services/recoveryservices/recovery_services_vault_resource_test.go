package recoveryservices_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type RecoveryServicesVaultResource struct {
}

func TestAccRecoveryServicesVault_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")
	r := RecoveryServicesVaultResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRecoveryServicesVault_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")
	r := RecoveryServicesVaultResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRecoveryServicesVault_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")
	r := RecoveryServicesVaultResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRecoveryServicesVault_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")
	r := RecoveryServicesVaultResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("sku").HasValue("Standard"),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccRecoveryServicesVault_basicWithIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")
	r := RecoveryServicesVaultResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t RecoveryServicesVaultResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.VaultID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.RecoveryServices.VaultsClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading Recovery Service (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (RecoveryServicesVaultResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-recovery-%d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-Vault-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  soft_delete_enabled = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RecoveryServicesVaultResource) basicWithIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-recovery-%d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-Vault-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  identity {
    type = "SystemAssigned"
  }

  soft_delete_enabled = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RecoveryServicesVaultResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-recovery-%d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-Vault-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  soft_delete_enabled = false
  storage_mode_type   = "LocallyRedundant"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r RecoveryServicesVaultResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_recovery_services_vault" "import" {
  name                = azurerm_recovery_services_vault.test.name
  location            = azurerm_recovery_services_vault.test.location
  resource_group_name = azurerm_recovery_services_vault.test.resource_group_name
  sku                 = azurerm_recovery_services_vault.test.sku
}
`, r.basic(data))
}
