---
layout: "gsuite"
page_title: "G Suite: group data source"
sidebar_current: "docs-gsuite-datasource-group"
description: |-
  Retrieves a Group in G Suite.
---

# gsuite\_group

Reads a Group from G Suite

## Example Usage

```hcl
data "gsuite_group" "example" {
  email = "example@domain.ext"
}

output "group" {
  value = data.gsuite_group.example
}
```

## Argument Reference

The following arguments are supported:

* `email` - (Required) The email of the group.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `id`

* `aliases`

* `name`

* `description`

* `direct_members_count`

* `admin_created`

* `non_editable_aliases`

* `member` - Lists the set of members in this group.
