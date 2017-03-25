package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/elliottsam/winrm-dns-client/dns"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestWinDNSRecord_Basic(t *testing.T) {
	var record dns.Record
	domain := os.Getenv("WINRM_DOMAIN")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWinDNSRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckWinDNSRecordConfig_basic, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWinDNSRecordExists("windows-dns_record.foobar", &record),
					testAccCheckWinDNSRecordAttributes(&record),
					resource.TestCheckResourceAttr(
						"windows-dns_record.foobar", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"windows-dns_record.foobar", "domain", domain),
					resource.TestCheckResourceAttr(
						"windows-dns_record.foobar", "value", "10.99.0.10"),
				),
			},
		},
	})
}

func TestAccWinDNSRecord_Updated(t *testing.T) {
	var record dns.Record
	domain := os.Getenv("WINRM_DOMAIN")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWinDNSRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckWinDNSRecordConfig_basic, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWinDNSRecordExists("windows-dns_record.foobar", &record),
					testAccCheckWinDNSRecordAttributes(&record),
					resource.TestCheckResourceAttr(
						"windows-dns_record.foobar", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"windows-dns_record.foobar", "domain", domain),
					resource.TestCheckResourceAttr(
						"windows-dns_record.foobar", "value", "10.99.0.10"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckWinDNSRecordConfig_new_value, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWinDNSRecordExists("windows-dns_record.foobar", &record),
					testAccCheckWinDNSRecordAttributesUpdated(&record),
					resource.TestCheckResourceAttr(
						"windows-dns_record.foobar", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"windows-dns_record.foobar", "domain", domain),
					resource.TestCheckResourceAttr(
						"windows-dns_record.foobar", "value", "10.99.99.99"),
				),
			},
		},
	})
}

func TestAccWinDNSRecord_Multiple(t *testing.T) {
	var record dns.Record
	domain := os.Getenv("WINRM_DOMAIN")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWinDNSRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckWinDNSRecordConfig_multiple, domain, domain, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWinDNSRecordExists("windows-dns_record.foobar1", &record),
					testAccCheckWinDNSRecordAttributes(&record),
					resource.TestCheckResourceAttr(
						"windows-dns_record.foobar1", "name", "terraform1"),
					resource.TestCheckResourceAttr(
						"windows-dns_record.foobar1", "domain", domain),
					resource.TestCheckResourceAttr(
						"windows-dns_record.foobar1", "value", "10.99.0.10"),
					resource.TestCheckResourceAttr(
						"windows-dns_record.foobar2", "name", "terraform2"),
					resource.TestCheckResourceAttr(
						"windows-dns_record.foobar2", "domain", domain),
					resource.TestCheckResourceAttr(
						"windows-dns_record.foobar2", "value", "10.99.1.10"),
					resource.TestCheckResourceAttr(
						"windows-dns_record.foobar3", "name", "terraform3"),
					resource.TestCheckResourceAttr(
						"windows-dns_record.foobar3", "domain", domain),
					resource.TestCheckResourceAttr(
						"windows-dns_record.foobar3", "value", "10.99.2.10"),
				),
			},
		},
	})
}

func testAccCheckWinDNSRecordDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*dns.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "windows-dns_record" {
			continue
		}

		foundRecord := dns.Record{
			Dnszone: rs.Primary.Attributes["doamin"],
			Name:    rs.Primary.Attributes["name"],
			ID:      rs.Primary.ID,
			Type:    rs.Primary.Attributes["type"],
		}

		if client.RecordExist(foundRecord) {
			return fmt.Errorf("Record still exists")
		}

	}

	return nil
}

func testAccCheckWinDNSRecordAttributes(record *dns.Record) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if record.Value != "10.99.0.10" {
			return fmt.Errorf("Bad value: %s", record.Value)
		}

		return nil
	}
}

func testAccCheckWinDNSRecordAttributesUpdated(record *dns.Record) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if record.Value != "10.99.99.99" {
			return fmt.Errorf("Bad value: %s", record.Value)
		}

		return nil
	}
}

func testAccCheckWinDNSRecordExists(n string, record *dns.Record) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*dns.Client)

		foundRecord := dns.Record{
			Dnszone: rs.Primary.Attributes["domain"],
			Name:    rs.Primary.Attributes["name"],
			Value:   rs.Primary.Attributes["value"],
			ID:      rs.Primary.ID,
			Type:    rs.Primary.Attributes["type"],
		}

		if !client.RecordExist(foundRecord) {
			return fmt.Errorf("Record not found")
		}

		*record = foundRecord

		return nil
	}
}

const testAccCheckWinDNSRecordConfig_basic = `
resource "windows-dns_record" "foobar" {
	domain = "%s"
	name = "terraform"
	value = "10.99.0.10"
	type = "A"
	ttl = "1h0m0s"
}`

const testAccCheckWinDNSRecordConfig_new_value = `
resource "windows-dns_record" "foobar" {
	domain = "%s"
	name = "terraform"
	value = "10.99.99.99"
	type = "A"
	ttl = "5m0s"
}`

const testAccCheckWinDNSRecordConfig_multiple = `
resource "windows-dns_record" "foobar1" {
	domain = "%s"
	name = "terraform1"
	value = "10.99.0.10"
	type = "A"
	ttl = "1h0m0s"
}
resource "windows-dns_record" "foobar2" {
	domain = "%s"
	name = "terraform2"
	value = "10.99.1.10"
	type = "A"
	ttl = "1h0m0s"
}
resource "windows-dns_record" "foobar3" {
	domain = "%s"
	name = "terraform3"
	value = "10.99.2.10"
	type = "A"
	ttl = "1h0m0s"
}`
