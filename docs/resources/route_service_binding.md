---
layout: "cloudfoundry"
page_title: "Cloud Foundry: cloudfoundry_route_service_binding"
sidebar_current: "docs-cf-resource-route-service-binding"
description: |-
  Provides a Cloud Foundry resource to bind a service instance to a route.
---

# cloudfoundry\_route\_service\_binding

Provides a Cloud Foundry resource for [binding](https://docs.cloudfoundry.org/devguide/services/route-binding.html#bind) of service instances to routes.

## Example Usage

The following example binds a specific route to the given service instance.

```hcl
resource "cloudfoundry_route_service_binding" "route-bind" {
  service_instance = cloudfoundry_service_instance.myservice.id
  route            = cloudfoundry_route.myroute.id
}
```

## Argument Reference

The following arguments are supported:

* `service_instance` - (Required, String) The ID the service instance to bind to the route.
* `route` - (Required, String) The ID of the route to bind the service instance to.
* `json_params` - (Optional, String) Arbitrary parameters in the form of stringified JSON object to pass to the service bind handler. Defaults to empty map.

## Import

Existing Route Service Binding can be imported using the composite `id` formed
with service instance's GUID and route's GUID.

Import does not support `json_params` attribute. Specifying non-empty `json_params` in
terraform files after import will lead to the recreation of the resource.

E.g.

```bash
$ terraform import cloudfoundry_route_service_binding.mybind <service-guid>/<route-guid>
```
