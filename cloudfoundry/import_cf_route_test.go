package cloudfoundry

import (
	"testing"

	"fmt"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccRoute_importBasic(t *testing.T) {

	_, orgName := defaultTestOrg(t)
	_, spaceName := defaultTestSpace(t)

	resourceName := "cloudfoundry_route.test-app-route"

	resource.Test(t,
		resource.TestCase{
			PreCheck:     func() { testAccPreCheck(t) },
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckRouteDestroyed([]string{"test-app-single", "test-app-multi"}, defaultAppDomain()),
			Steps: []resource.TestStep{

				resource.TestStep{
					Config: fmt.Sprintf(routeResource,
						defaultAppDomain(),
						orgName, spaceName),
				},

				resource.TestStep{
					ResourceName:      resourceName,
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
}
