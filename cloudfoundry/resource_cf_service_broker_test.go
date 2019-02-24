package cloudfoundry

import (
	"fmt"
	"testing"

	"code.cloudfoundry.org/cli/cf/errors"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-cloudfoundry/cloudfoundry/cfapi"
)

const sbResource = `

resource "cloudfoundry_service_broker" "test" {
	name = "test"
	url = "%s"
	username = "%s"
	password = "%s"
}
`

const sbResourceUpdate = `

resource "cloudfoundry_service_broker" "test" {
	name = "test-renamed"
	url = "%s"
	username = "%s"
	password = "%s"
}
`

func TestAccServiceBroker_normal(t *testing.T) {

	serviceBrokerURL, serviceBrokerUser, serviceBrokerPassword, serviceBrokerPlanPath := getTestBrokerCredentials(t)

	// Ensure any test artifacts from a
	// failed run are deleted if the exist
	deleteServiceBroker("test")
	deleteServiceBroker("test-renamed")

	ref := "cloudfoundry_service_broker.test"

	resource.Test(t,
		resource.TestCase{
			PreCheck:     func() { testAccPreCheck(t) },
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckServiceBrokerDestroyed("test"),
			Steps: []resource.TestStep{

				resource.TestStep{
					Config: fmt.Sprintf(sbResource,
						serviceBrokerURL, serviceBrokerUser, serviceBrokerPassword),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckServiceBrokerExists(ref),
						resource.TestCheckResourceAttr(
							ref, "name", "test"),
						resource.TestCheckResourceAttr(
							ref, "url", serviceBrokerURL),
						resource.TestCheckResourceAttr(
							ref, "username", serviceBrokerUser),
						resource.TestCheckResourceAttrSet(
							ref, "service_plans."+serviceBrokerPlanPath),
					),
				},

				resource.TestStep{
					Config: fmt.Sprintf(sbResourceUpdate,
						serviceBrokerURL, serviceBrokerUser, serviceBrokerPassword),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckServiceBrokerExists(ref),
						resource.TestCheckResourceAttr(
							ref, "name", "test-renamed"),
					),
				},
			},
		})
}

func testAccCheckServiceBrokerExists(resource string) resource.TestCheckFunc {

	return func(s *terraform.State) (err error) {

		session := testAccProvider.Meta().(*cfapi.Session)

		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("service broker '%s' not found in terraform state", resource)
		}

		session.Log.DebugMessage(
			"terraform state for resource '%s': %# v",
			resource, rs)

		id := rs.Primary.ID
		attributes := rs.Primary.Attributes

		var (
			serviceBroker cfapi.CCServiceBroker
		)

		sm := session.ServiceManager()
		if serviceBroker, err = sm.ReadServiceBroker(id); err != nil {
			return
		}

		if err := assertEquals(attributes, "name", serviceBroker.Name); err != nil {
			return err
		}
		if err := assertEquals(attributes, "url", serviceBroker.BrokerURL); err != nil {
			return err
		}
		if err := assertEquals(attributes, "username", serviceBroker.AuthUserName); err != nil {
			return err
		}

		return
	}
}

func testAccCheckServiceBrokerDestroyed(name string) resource.TestCheckFunc {

	return func(s *terraform.State) error {

		session := testAccProvider.Meta().(*cfapi.Session)
		if _, err := session.ServiceManager().GetServiceBrokerID(name); err != nil {
			switch err.(type) {
			case *errors.ModelNotFoundError:
				return nil
			default:
				return err
			}
		}

		return fmt.Errorf("service broker with name '%s' still exists in cloud foundry", name)
	}
}
