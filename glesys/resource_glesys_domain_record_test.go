package glesys

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDNSDomainRecord_ARecord(t *testing.T) {
	time.Sleep(2 * time.Second)
	domain := "tf-" + acctest.RandString(6) + ".com"
	ip := "192.0.2.123"

	resourceName := "glesys_dnsdomain_record.test"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testGlesysProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSDomainRecordConfig(domain, "www", ip, "A", 3600, "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "host", "www"),
					resource.TestCheckResourceAttr(resourceName, "data", ip),
					resource.TestCheckResourceAttr(resourceName, "type", "A"),
				),
			},
		},
	})
}

func TestAccDNSDomainRecord_TXT_DKIM(t *testing.T) {
	time.Sleep(2 * time.Second)
	domain := "tf-" + acctest.RandString(6) + ".com"
	dkim := `v=DKIM1;g=*;k=rsa;p=MIGf...trimmed...AQAB`

	resourceName := "glesys_dnsdomain_record.dkim"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testGlesysProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSDomainRecordConfig(domain, "selector._domainkey", dkim, "TXT", 3600, "dkim"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "host", "selector._domainkey"),
					resource.TestCheckResourceAttr(resourceName, "data", dkim),
					resource.TestCheckResourceAttr(resourceName, "type", "TXT"),
				),
			},
		},
	})
}

func TestAccDNSDomainRecord_MX(t *testing.T) {
	time.Sleep(2 * time.Second)
	domain := "tf-" + acctest.RandString(6) + ".com"
	resourceName := "glesys_dnsdomain_record.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testGlesysProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "glesys_dnsdomain" "test" {
					  name = "%s"
					}

					resource "glesys_dnsdomain_record" "test" {
					  domain = glesys_dnsdomain.test.name
					  host   = "@"
					  type   = "MX"
					  data   = "10 mail.%s"
					  ttl    = 3600
					}
				`, domain, domain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "MX"),
					resource.TestCheckResourceAttr(resourceName, "data", fmt.Sprintf("10 mail.%s", domain)),
				),
			},
		},
	})
}

func TestAccDNSDomainRecord_SRV(t *testing.T) {
	time.Sleep(2 * time.Second)
	domain := "tf-" + acctest.RandString(6) + ".com"
	resourceName := "glesys_dnsdomain_record.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testGlesysProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "glesys_dnsdomain" "test" {
					  name = "%s"
					}

					resource "glesys_dnsdomain_record" "test" {
					  domain = glesys_dnsdomain.test.name
					  host   = "_sip._tcp"
					  type   = "SRV"
					  data   = "10 20 5060 sip.%s"
					  ttl    = 3600
					}
				`, domain, domain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "SRV"),
					resource.TestCheckResourceAttr(resourceName, "data", fmt.Sprintf("10 20 5060 sip.%s", domain)),
				),
			},
		},
	})
}

func TestAccDNSDomainRecord_conflict(t *testing.T) {
	time.Sleep(2 * time.Second)
	domain := "tf-" + acctest.RandString(6) + ".com"
	ip := "198.51.100.42"

	resourceName := "glesys_dnsdomain_record.test"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testGlesysProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSDomainRecordConfig(domain, "api", ip, "A", 3600, "test"),
				Check:  resource.TestCheckResourceAttr(resourceName, "data", ip),
			},
			{
				Config: testAccDNSDomainRecordConfig(domain, "api", ip, "A", 3600, "test"), // identical apply
				Check:  resource.TestCheckResourceAttr(resourceName, "data", ip),
			},
		},
	})
}

func testAccDNSDomainRecordConfig(domain, host, data, recordType string, ttl int, name string) string {
	return fmt.Sprintf(`
resource "glesys_dnsdomain" "test" {
  name = "%s"
}

resource "glesys_dnsdomain_record" "%s" {
  domain = glesys_dnsdomain.test.name
  host   = "%s"
  data   = "%s"
  type   = "%s"
  ttl    = %d
}
`, domain, name, host, data, recordType, ttl)
}

func TestAccDNSDomainRecord_DuplicateRecord(t *testing.T) {
	domain := "tf-duplicate-" + acctest.RandString(6) + ".com"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testGlesysProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "glesys_dnsdomain" "test" {
						name = "%s"
					}

					resource "glesys_dnsdomain_record" "first" {
						domain = glesys_dnsdomain.test.name
						host   = "www"
						type   = "A"
						data   = "1.2.3.4"
						ttl    = 3600
					}
				`, domain),
			},
			{
				Config: fmt.Sprintf(`
					resource "glesys_dnsdomain" "test" {
						name = "%s"
					}

					resource "glesys_dnsdomain_record" "first" {
						domain = glesys_dnsdomain.test.name
						host   = "www"
						type   = "A"
						data   = "1.2.3.4"
						ttl    = 3600
					}

					resource "glesys_dnsdomain_record" "second" {
						domain = glesys_dnsdomain.test.name
						host   = "www"
						type   = "A"
						data   = "1.2.3.4"
						ttl    = 3600
					}
				`, domain),
			},
		},
	})
}

func TestAccDNSDomainRecord_InvalidMXRecordFails(t *testing.T) {
	domain := "tf-invalidmx-" + acctest.RandString(6) + ".com"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testGlesysProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "glesys_dnsdomain" "test" {
						name = "%s"
					}

					resource "glesys_dnsdomain_record" "test" {
						domain = glesys_dnsdomain.test.name
						host   = "@"
						type   = "MX"
						data   = "mail.%s"  // Missing priority causes 422
						ttl    = 3600
					}
				`, domain, domain),
				ExpectError: regexp.MustCompile(`422.*Invalid data for this record type`),
			},
		},
	})
}

func TestAccDNSDomainRecord_AdoptsDuplicateRecord(t *testing.T) {
	domain := "tf-adopt-" + acctest.RandString(5) + ".com"
	host := "www"
	recType := "A"
	data := "192.0.2.42"
	ttl := 3600

	config := fmt.Sprintf(`
resource "glesys_dnsdomain" "test" {
  name = "%s"
}

resource "glesys_dnsdomain_record" "test" {
  domain = glesys_dnsdomain.test.name
  host   = "%s"
  type   = "%s"
  data   = "%s"
  ttl    = %d
}
`, domain, host, recType, data, ttl)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testGlesysProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				// Reapply same config, should be adopted not recreated
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("glesys_dnsdomain_record.test", "data", data),
					resource.TestCheckResourceAttr("glesys_dnsdomain_record.test", "ttl", fmt.Sprintf("%d", ttl)),
				),
			},
		},
	})
}

func TestAccDNSDomainRecord_MultipleARecordsWithDifferentData(t *testing.T) {
	domain := "tf-conflict-" + acctest.RandString(5) + ".com"
	host := "mail"
	recType := "A"
	data1 := "192.0.2.10"
	data2 := "192.0.2.20"

	initial := fmt.Sprintf(`
resource "glesys_dnsdomain" "test" {
  name = "%s"
}

resource "glesys_dnsdomain_record" "test" {
  domain = glesys_dnsdomain.test.name
  host   = "%s"
  type   = "%s"
  data   = "%s"
}
`, domain, host, recType, data1)

	addSecondRecord := fmt.Sprintf(`
resource "glesys_dnsdomain" "test" {
  name = "%s"
}

resource "glesys_dnsdomain_record" "test" {
  domain = glesys_dnsdomain.test.name
  host   = "%s"
  type   = "%s"
  data   = "%s"
}

resource "glesys_dnsdomain_record" "test2" {
  domain = glesys_dnsdomain.test.name
  host   = "%s"
  type   = "%s"
  data   = "%s"
}
`, domain, host, recType, data1, host, recType, data2)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testGlesysProviders,
		Steps: []resource.TestStep{
			{Config: initial},
			{Config: addSecondRecord}, // no ExpectError here, both records should coexist
		},
	})
}

func TestAccDNSDomainRecord_AdoptsWithSameTTL(t *testing.T) {
	domain := "tf-ttl-match-" + acctest.RandString(5) + ".com"
	host := "static"
	recType := "A"
	data := "192.0.2.50"
	ttl := 1800

	config := fmt.Sprintf(`
resource "glesys_dnsdomain" "test" {
  name = "%s"
}

resource "glesys_dnsdomain_record" "test" {
  domain = glesys_dnsdomain.test.name
  host   = "%s"
  type   = "%s"
  data   = "%s"
  ttl    = %d
}
`, domain, host, recType, data, ttl)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testGlesysProviders,
		Steps: []resource.TestStep{
			{Config: config},
			{
				// Reapply same TTL — should adopt, not recreate
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("glesys_dnsdomain_record.test", "ttl", fmt.Sprintf("%d", ttl)),
				),
			},
		},
	})
}

func TestAccDNSDomainRecord_UpdatesTTL(t *testing.T) {
	domain := "tf-ttl-change-" + acctest.RandString(5) + ".com"
	host := "app"
	recType := "A"
	data := "192.0.2.60"
	ttl1 := 300
	ttl2 := 900

	config1 := fmt.Sprintf(`
resource "glesys_dnsdomain" "test" {
  name = "%s"
}

resource "glesys_dnsdomain_record" "test" {
  domain = glesys_dnsdomain.test.name
  host   = "%s"
  type   = "%s"
  data   = "%s"
  ttl    = %d
}
`, domain, host, recType, data, ttl1)

	config2 := strings.Replace(config1, fmt.Sprintf("ttl    = %d", ttl1), fmt.Sprintf("ttl    = %d", ttl2), 1)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testGlesysProviders,
		Steps: []resource.TestStep{
			{
				Config: config1,
				Check:  resource.TestCheckResourceAttr("glesys_dnsdomain_record.test", "ttl", fmt.Sprintf("%d", ttl1)),
			},
			{
				Config: config2,
				Check:  resource.TestCheckResourceAttr("glesys_dnsdomain_record.test", "ttl", fmt.Sprintf("%d", ttl2)),
			},
		},
	})
}

func TestAccDNSDomainRecord_RejectsDuplicateRecordWithDifferentTTL(t *testing.T) {
	domain := "tf-ttl-change-" + acctest.RandString(5) + ".com"

	initialConfig := fmt.Sprintf(`
resource "glesys_dnsdomain" "test" {
  name = "%s"
}

resource "glesys_dnsdomain_record" "ttltest600" {
  domain = glesys_dnsdomain.test.name
  host   = "ttltest"
  type   = "A"
  data   = "203.0.113.10"
  ttl    = 600
}
`, domain)

	duplicateWithDifferentTTL := fmt.Sprintf(`
resource "glesys_dnsdomain" "test" {
  name = "%s"
}

resource "glesys_dnsdomain_record" "ttltest3600" {
  domain = glesys_dnsdomain.test.name
  host   = "ttltest"
  type   = "A"
  data   = "203.0.113.10"
  ttl    = 3600
}
`, domain)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testGlesysProviders,
		Steps: []resource.TestStep{
			{
				Config: initialConfig,
			},
			{
				Config:      duplicateWithDifferentTTL,
				ExpectError: regexp.MustCompile(`already exists with different TTL`),
			},
		},
	})
}

func TestAccDNSDomainRecord_AdoptsIdenticalRecord(t *testing.T) {
	domain := "tf-ttl-change-" + acctest.RandString(5) + ".com"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testGlesysProviders,
		Steps: []resource.TestStep{

			{
				Config: fmt.Sprintf(`
resource "glesys_dnsdomain" "test" {
  name = "%s"
}

resource "glesys_dnsdomain_record" "dup" {
  domain = glesys_dnsdomain.test.name
  host   = "www"
  type   = "A"
  data   = "198.51.100.42"
  ttl    = 3600
}
`, domain),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("glesys_dnsdomain_record.dup", "data", "198.51.100.42"),
				),
			},
			{
				// Re-run exact config — should adopt, not recreate or fail
				Config: fmt.Sprintf(`
resource "glesys_dnsdomain" "test" {
  name = "%s"
}

resource "glesys_dnsdomain_record" "dup" {
  domain = glesys_dnsdomain.test.name
  host   = "www"
  type   = "A"
  data   = "198.51.100.42"
  ttl    = 3600
}
`, domain),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("glesys_dnsdomain_record.dup", "type", "A"),
				),
			},
		},
	})
}

func TestAccDNSDomainRecord_MultipleTXTRecords(t *testing.T) {
	domain := "tf-ttl-change-" + acctest.RandString(5) + ".com"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testGlesysProviders,
		Steps: []resource.TestStep{

			{
				Config: fmt.Sprintf(`
resource "glesys_dnsdomain" "test" {
  name = "%s"
}

resource "glesys_dnsdomain_record" "txt1" {
  domain = glesys_dnsdomain.test.name
  host   = "spf"
  type   = "TXT"
  data   = "v=spf1 include:example.com ~all"
  ttl    = 600
}

resource "glesys_dnsdomain_record" "txt2" {
  domain = glesys_dnsdomain.test.name
  host   = "spf"
  type   = "TXT"
  data   = "v=spf1 ip4:192.0.2.0/24 -all"
  ttl    = 600
}
`, domain),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("glesys_dnsdomain_record.txt1", "type", "TXT"),
					resource.TestCheckResourceAttr("glesys_dnsdomain_record.txt2", "type", "TXT"),
				),
			},
		},
	})
}

func TestAccDNSDomainRecord_MultipleARecords(t *testing.T) {
	domain := "tf-ttl-change-" + acctest.RandString(5) + ".com"

	resourceName1 := "glesys_dnsdomain_record.record1"
	resourceName2 := "glesys_dnsdomain_record.record2"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testGlesysProviders,
		Steps: []resource.TestStep{

			{
				Config: fmt.Sprintf(`
resource "glesys_dnsdomain" "test" {
  name = "%s"
}

resource "glesys_dnsdomain_record" "record1" {
  domain = glesys_dnsdomain.test.name
  host   = "multi"
  type   = "A"
  data   = "192.0.2.1"
  ttl    = 3600
}

resource "glesys_dnsdomain_record" "record2" {
  domain = glesys_dnsdomain.test.name
  host   = "multi"
  type   = "A"
  data   = "192.0.2.2"
  ttl    = 3600
}
`, domain),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName1, "data", "192.0.2.1"),
					resource.TestCheckResourceAttr(resourceName2, "data", "192.0.2.2"),
				),
			},
		},
	})
}
