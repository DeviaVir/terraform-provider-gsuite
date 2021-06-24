---
layout: "gsuite"
page_title: "G Suite: gsuite_alias"
sidebar_current: "docs-gsuite-resource-user-alias"
description: |-
  Managing a G Suite User Alias.
---

# gsuite\_user\_alias

Provides a resource for creating and managing an email alias for a GSuite user account.

## Example Usage

```hcl
resource "gsuite_user_alias" "test" {
  user_id = "test-user-replaceWithUuid@domain.ext"
  alias   = "test-alias-replaceWithUuid@domain.ext"
}
```

## Argument Reference

* `user_id` (Required) Primary email (userKey) of the user who will have the alias applied to them.
* `alias` (Required) Email alias to be applied to the user.


## Attribute Reference

N/A apart from the included arguments

## Import 

An alias can be imported by passing the ID format of "<user_id>/<alias>"

For example:
```
terraform import gsuite_user_alias.test "test-user@domain.ext/test-alias@domain.ext"
```
