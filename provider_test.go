package main

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]terraform.ResourceProvider{
		"windowsdns": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("Error: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("WINRM_USERNAME"); v == "" {
		t.Fatal("WINRM_USERNAME must be set for tests")
	}

	if v := os.Getenv("WINRM_PASSWORD"); v == "" {
		t.Fatal("WINRM_PASSWORD must be set for tests")
	}

	if v := os.Getenv("WINRM_SERVER"); v == "" {
		t.Fatal("WINRM_SERVER must be set for tests")
	}

	if v := os.Getenv("WINRM_DOMAIN"); v == "" {
		t.Fatal("WINRM_DOMAIN must be set for tests")
	}
}

func testPorts(t *testing.T) {
	tables := []struct {
		p int
		https bool
		e int
	}{
		{1234, false, 1234},
		{1234, true, 1234},
		{0, false, 5985},
		{0, true, 5986},
	}
	
for _, table := range tables {
	derivedPort := derivePort(table.p, table.https)
	if (derivedPort != table.e) {
		t.Errorf("Port test for port %d with HTTPS %t returned %d instead of %d", table.p, table.https, derivedPort, table.e)
	}
}

}
