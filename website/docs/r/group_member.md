---
layout: "gsuite"
page_title: "G Suite: gsuite_group_member"
sidebar_current: "docs-gsuite-resource-group-member"
description: |-
  Managing a single Group Member in a Group
---

# gsuite\_group\_member

Provides a resource to create and manage a single group member.

**Note:** do not use this resource in conjunction with `gsuite_group_members`!

## Example Usage

```hcl
resource "gsuite_group" "example" {
  email       = "example@domain.ext"
  name        = "example@domain.ext"
  description = "Example group"
}

resource "gsuite_group_member" "owner" {
  group = gsuite_group.example.email
  email = "owner@domain.ext"
  role  = "OWNER"
}
```

## Argument Reference

The following arguments are supported:

* `email` - (Required; Forces new resource) Email address of the member.

* `role` - (Optional) Defaults to `MEMBER`. Other groups cannot be `OWNER`.


## Attribute Reference

In addition to the above arguments, the following attributes are exported:

* `etag` - ETag of the resource.

* `kind` - Kind of resource this is.

* `status` - Status of member.

* `type`- Type of member.

## Import

A G Suite Group Member can be imported using `group-email/user-email`, e.g.:

```
terraform import gsuite_group_member.owner "example@domain.ext/owner@domain.ext"
```
