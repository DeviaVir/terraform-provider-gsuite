---
layout: "gsuite"
page_title: "G Suite: gsuite_group_members"
sidebar_current: "docs-gsuite-resource-group-members"
description: |-
  Managing all Group Members of a G Suite Group
---

# gsuite\_group\_members

Provides a resource to create and manage all Members of a G Suite Group.

**Note:** do not use this resource in conjunction with `gsuite_group_member`!

## Example Usage

```hcl
resource "gsuite_group" "example" {
  email       = "example@domain.ext"
  name        = "example@domain.ext"
  description = "Example group"
}

resource "gsuite_group_members" "members" {
  group_email = gsuite_group.example.email

  member {
    email = "member@domain.ext"
    role  = "MEMBER"
  }

  member {
    email = "owner@domain.ext"
    role  = "OWNER"
  }
}
```

## Argument Reference

The following arguments are supported:

* `group_email` - (Required; Forces new resource) Email address of the G Suite
  group.


## Attribute Reference

In addition to the above arguments, the following attributes are exported:

* `members` - contains the set of members, with the following schema:
  * `email` - Email of member.
  * `etag` - ETag of the resource.
  * `kind` - Kind of resource this is.
  * `status` - Status of member.
  * `type` - Type of member.
  * `role` - Role of member.

## Import

G Suite Group Members can be imported using `group-email`, e.g.:

```
terraform import gsuite_group_members.members "example@domain.ext"
```
