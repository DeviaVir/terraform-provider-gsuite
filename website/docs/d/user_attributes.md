---
layout: "gsuite"
page_title: "G Suite: user attributes data source"
sidebar_current: "docs-gsuite-datasource-user-attributes"
description: |-
  Retrieves Attributes of a User in G Suite.
---

# gsuite\_user\_attributes

Reads attributes of a User in G Suite.

## Example Usage

```hcl
data "gsuite_user_attributes" "example" {
  name = "example"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the user.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `value`
