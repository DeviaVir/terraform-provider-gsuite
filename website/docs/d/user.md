---
layout: "gsuite"
page_title: "G Suite: user data source"
sidebar_current: "docs-gsuite-datasource-user"
description: |-
  Retrieves User in G Suite.
---

# gsuite\_user

Reads attributes of a User in G Suite.

## Example Usage

```hcl
data "gsuite_user" "example" {
  primary_email = "example@domain.ext"
}
```

## Argument Reference

The following arguments are supported:

* `primary_email` - (Required) The primary email address of the user.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `org_unit_path`

* `aliases`

* `agreed_to_terms`

* `change_password_next_login`

* `creation_time`

* `customer_id`

* `deletion_time`

* `etag`

* `include_in_global_list`

* `is_ip_whitelisted`

* `is_admin`

* `is_delegated_admin`

* `2s_enforced`

* `2s_enrolled`

* `is_mailbox_setup`

* `last_login_time`

* `name` - contains a set of `family_name`, `full_name`, `given_name`.

* `password`

* `hash_function` - md5, sha-1 and crypt.

* `posix_accounts` - contains a list of sets containing `account_id`, `gecos`,
  `gid`, `home_directory`, `shell`, `system_id`, `primary`, `uid`, `username`.

* `ssh_public_keys` - contains a list of sets containing `expiration_time_usec`,
  `key`, `fingerprint`.

* `is_suspended`

* `suspension_reason`

* `recovery_email`

* `recovery_phone`

* `custom_schema` - contains a schema of `name` and `value`.

* `external_ids` - contains a list of sets containing `custom_type`, `type` and
  `value`.

* `organizations` - contains a list of sets containing `cost_center`,
  `custom_type`, `department`, `description`, `domain`, `full_time_equivalent`,
  `location`, `name`, `primary`, `symbol`, `title` and `type`.
