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

* `id` - Unique identifier of Group.

* `aliases` - List of aliases.

* `name` - Group name.

* `description` - Description of the group.

* `direct_members_count` - Group direct members count.

* `admin_created` - Is the group created by admin.

* `non_editable_aliases` - List of non editable aliases.

* `member` - Lists the set of members in this group.
