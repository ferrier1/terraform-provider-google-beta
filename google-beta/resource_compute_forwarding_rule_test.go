package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccComputeForwardingRule_update(t *testing.T) {
	t.Parallel()

	poolName := fmt.Sprintf("tf-%s", acctest.RandString(10))
	ruleName := fmt.Sprintf("tf-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeForwardingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeForwardingRule_basic(poolName, ruleName),
			},
			{
				ResourceName:      "google_compute_forwarding_rule.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeForwardingRule_update(poolName, ruleName),
			},
			{
				ResourceName:      "google_compute_forwarding_rule.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeForwardingRule_ip(t *testing.T) {
	t.Parallel()

	addrName := fmt.Sprintf("tf-%s", acctest.RandString(10))
	poolName := fmt.Sprintf("tf-%s", acctest.RandString(10))
	ruleName := fmt.Sprintf("tf-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeForwardingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeForwardingRule_ip(addrName, poolName, ruleName),
			},
			{
				ResourceName:      "google_compute_forwarding_rule.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeForwardingRule_networkTier(t *testing.T) {
	t.Parallel()

	poolName := fmt.Sprintf("tf-%s", acctest.RandString(10))
	ruleName := fmt.Sprintf("tf-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeForwardingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeForwardingRule_networkTier(poolName, ruleName),
			},

			{
				ResourceName:      "google_compute_forwarding_rule.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeForwardingRule_mirror(t *testing.T) {
	t.Parallel()

	backendName := fmt.Sprintf("tf-%s", acctest.RandString(10))
	healthchechName := fmt.Sprintf("tf-%s", acctest.RandString(10))
	ruleName := fmt.Sprintf("tf-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeForwardingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeForwardingRule_mirror(backendName, healthchechName, ruleName),
			},

			{
				ResourceName:      "google_compute_forwarding_rule.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeForwardingRule_basic(poolName, ruleName string) string {
	return fmt.Sprintf(`
resource "google_compute_target_pool" "foo-tp" {
  description = "Resource created for Terraform acceptance testing"
  instances   = ["us-central1-a/foo", "us-central1-b/bar"]
  name        = "foo-%s"
}

resource "google_compute_forwarding_rule" "foobar" {
  description = "Resource created for Terraform acceptance testing"
  ip_protocol = "UDP"
  name        = "%s"
  port_range  = "80-81"
  target      = google_compute_target_pool.foo-tp.self_link
  labels = {
    "foo" = "bar"
  }
}
`, poolName, ruleName)
}

func testAccComputeForwardingRule_update(poolName, ruleName string) string {
	return fmt.Sprintf(`
resource "google_compute_target_pool" "foo-tp" {
  description = "Resource created for Terraform acceptance testing"
  instances   = ["us-central1-a/foo", "us-central1-b/bar"]
  name        = "foo-%s"
}

resource "google_compute_target_pool" "bar-tp" {
  description = "Resource created for Terraform acceptance testing"
  instances   = ["us-central1-a/foo", "us-central1-b/bar"]
  name        = "bar-%s"
}

resource "google_compute_forwarding_rule" "foobar" {
  description = "Resource created for Terraform acceptance testing"
  ip_protocol = "UDP"
  name        = "%s"
  port_range  = "80-81"
  target      = google_compute_target_pool.bar-tp.self_link
  labels = {
    "baz" = "qux"
  }
}
`, poolName, poolName, ruleName)
}

func testAccComputeForwardingRule_ip(addrName, poolName, ruleName string) string {
	return fmt.Sprintf(`
resource "google_compute_address" "foo" {
  name = "%s"
}

resource "google_compute_target_pool" "foobar-tp" {
  description = "Resource created for Terraform acceptance testing"
  instances   = ["us-central1-a/foo", "us-central1-b/bar"]
  name        = "%s"
}

resource "google_compute_forwarding_rule" "foobar" {
  description = "Resource created for Terraform acceptance testing"
  ip_address  = google_compute_address.foo.address
  ip_protocol = "TCP"
  name        = "%s"
  port_range  = "80-81"
  target      = google_compute_target_pool.foobar-tp.self_link
}
`, addrName, poolName, ruleName)
}

func testAccComputeForwardingRule_networkTier(poolName, ruleName string) string {
	return fmt.Sprintf(`
resource "google_compute_target_pool" "foobar-tp" {
  description = "Resource created for Terraform acceptance testing"
  instances   = ["us-central1-a/foo", "us-central1-b/bar"]
  name        = "%s"
}

resource "google_compute_forwarding_rule" "foobar" {
  description = "Resource created for Terraform acceptance testing"
  ip_protocol = "UDP"
  name        = "%s"
  port_range  = "80-81"
  target      = google_compute_target_pool.foobar-tp.self_link

  network_tier = "STANDARD"
}
`, poolName, ruleName)
}

func testAccComputeForwardingRule_mirror(backendName, healthchechName, ruleName string) string {
	return fmt.Sprintf(`
resource "google_compute_region_backend_service" "foobar-be" {
	name        	= "%s"
	health_checks 	= [google_compute_health_check.foobar-hc.self_link]
	}
	
	resource "google_compute_health_check" "foobar-hc" {
	name        		= "%s"
	check_interval_sec 	= 1
	timeout_sec        	= 1
	
	tcp_health_check {
		port = "80"
	}
}

resource "google_compute_forwarding_rule" "foobar" {
  description 			= "Resource created for Terraform acceptance testing"
  ip_protocol 			= "TCP"
  name        			= "%s"
  all_ports  			= true
  mirroring_collector   = true
  backend_service       = google_compute_region_backend_service.foobar-be.self_link
  load_balancing_scheme = "INTERNAL"
}
`, backendName, healthchechName, ruleName)
}
