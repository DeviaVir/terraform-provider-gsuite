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

* `org_unit_path` - OrgUnit of User.

* `aliases` - List of aliases.

* `agreed_to_terms` - Indicates if user has agreed to terms.

* `change_password_next_login` - Boolean indicating if the user should
  change password in next login

* `creation_time` - User's G Suite account creation time.

* `customer_id` - CustomerId of User.

* `deletion_time` - User's G Suite account deletion time.

* `etag` - ETag of the resource.

* `include_in_global_list` - Boolean indicating if user is included in
  Global Address List.

* `is_ip_whitelisted` - Boolean indicating if ip is whitelisted.

* `is_admin` - Boolean indicating if the user is admin.

* `is_delegated_admin` - Boolean indicating if the user is delegated admin.

* `2s_enforced` - Is 2-step verification enforced.

* `2s_enrolled` - Is enrolled in 2-step verification.

* `is_mailbox_setup` - Is mailbox setup.

* `last_login_time` - User's last login time.

* `name` - User's name.
  contains a set of `family_name`, `full_name`, `given_name`.

* `password` - User's password.

* `hash_function` - Hash function name for password. Supported are MD5,
  SHA-1 and crypt

* `posix_accounts` - contains a list of sets containing `account_id`, `gecos`,
  `gid`, `home_directory`, `shell`, `system_id`, `primary`, `uid`, `username`.

* `ssh_public_keys` - SSH public keys of the user.
  contains a list of sets containing `expiration_time_usec`, `key`, `fingerprint`.

* `is_suspended` - Indicates if user is suspended.

* `suspension_reason` -  Suspension reason if user is suspended.

* `recovery_email` - Recovery email of the user.

* `recovery_phone` - Recovery phone of the user.

* `custom_schema` - Custom fields of the user.

* `external_ids` - A list of external IDs for the user, such as an employee or network ID. 
  contains a list of sets containing `custom_type`, `type` and `value`.

* `organizations` - A list of organizations the user belongs to.
  contains a list of sets containing `cost_center`,
  `custom_type`, `department`, `description`, `domain`, `full_time_equivalent`,
  `location`, `name`, `primary`, `symbol`, `title` and `type`.
