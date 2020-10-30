---
layout: "gsuite"
page_title: "G Suite: gsuite_group"
sidebar_current: "docs-gsuite-resource-group"
description: |-
  Managing a G Suite Group.
---

# gsuite\_group

Provides a resource to create and manage a G Suite group.

## Example Usage

```hcl
resource "gsuite_group" "example" {
  email       = "example@domain.ext"
  name        = "example@domain.ext"
  description = "Example group"
}
```

## Argument Reference

The following arguments are supported:

* `email` - (Required; Forces new resource) Email address of the G Suite
  group.

* `aliases` - (Optional) Provide a list of aliases for this Group.

* `name` - (Optional) Group name.

* `description` - (Optional) Description of the group.

## Attribute Reference

In addition to the above arguments, the following attributes are exported:

* `direct_members_count` - Group direct members count.

* `admin_created` - Is the group created by admin.

* `non_editable_aliases` - List of non editable aliases.

## Import

A G Suite Group can be imported using `group-email`, e.g.:

```
terraform import gsuite_group.example "example@domain.ext"
```
