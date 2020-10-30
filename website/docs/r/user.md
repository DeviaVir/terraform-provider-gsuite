---
layout: "gsuite"
page_title: "G Suite: gsuite_user_schema"
sidebar_current: "docs-gsuite-resource-user-schema"
description: |-
  Managing a G Suite User Schema.
---

# gsuite\_user\_schema

Provides a resource to create and manage a G Suite User Schema.

**Note** the following behaviors regarding passwords:

- When running `terraform import` on a user resource:
  - The `password` and `hash_function` fields are ignored.
- When running `terraform apply` with a new user resource in your terraform state:
  - If the user does not exist in GSuite the following applies:
  - The `password` field should be set or a secured password will be automatically generated.
  - The `hash_function` field must be set only if the `password` field contains a hashed value.
  - The GSuite account will be configured to require password change on next login.
- If the user exists in GSuite the following applies:
  - The `password` and `hash_function` fields will be ignored.
- When running `terraform apply` with an existing user resource:
  - Empty `password` and `hash_function` fields will be ignored.

**Warn:** it is possible on-creation of a new account that the POSIX data is
found to not be unique, and a 503 backend error is returned indefinitely.
In that case, the account is created, but without the POSIX data. Simply
update the POSIX data and terraform apply to update until it works.

## Example Usage

```hcl
resource "gsuite_user" "developer" {

  aliases = [
    "chase@domain.ext"
  ]

  name {
    family_name = "Sillevis"
    given_name  = "Chase"
  }

  password = "testtest123!"

  primary_email = "developer@domain.ext"

  ssh_public_keys {
    key                  = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDUYJKI2gGdZr5Brd1IaT8OQSSt81mBBXQnAfmmjw5hOK9VaJ1MmDB5qY7V1nuXftmLBLvaA7L6k21FDJeWxwD8vKuYwbuJyh1DKB6PMXAQxnX7uLSSi9a/ZOzh3gIHXdil0fSWFpFBmImznqbzaEb7nya+tnK4RONoEjJcRe8Tl+8hET/29XBd3oxlfwwjQA9A84iKhAMLdJIQ28z2GA/2mRJ8RkHLrkQL8kMCj4GJYxy3PR9JU0aFAtWh2mXGfOzaBTh/IhpMW53d8puxihBbIN87MoGngYLt4eBEdE0SiHb0Zdqp5ZDCkwNmAKiWOOrDQxtWvUOThHV5eLMMObqA06XFiwNlojl9ZTH0Y2w/LZmvgb98T/1lBY6mb1iRERGKqYNBeSNwh1Afvu1miDau2f5AYqcf7yxvuD8d0O4cb1xfl7WJwWPJraYaN1X+WmCGTIA+Vve+Kp9TaGoE5n5EGz2a7RNzWj0L0hkf8923iEEtTrsfWewnTnq7XzFoaW53xjWcN7jQplisjWr6AWYApyinw0qGD3dzKgPLyOOcdC3YLhYFpGJcMbegrNdmhbxqIXCB3vBpEFV6o4GqdEy2OVFOM6kSydEQUsMHl5WU8l4gYW28ekZZtbrE52v1dMNzKwfrpVPpUfwn4jbeaqYoIWEwFNVnvbJaFu1vjfrshw== chase"
    expiration_time_usec = "1549735670773"
  }

  posix_accounts {
    home_directory = "/home/chase"
    primary        = true
    gid            = 1001
    uid            = 1001
    shell          = "/bin/bash"
    system_id      = "uid"
    username       = "chase"
  }

  external_ids {
    type  = "organization"
    value = "1234"
  }

  # If omitted or `true` existing GSuite users defined as Terraform resources will be imported by `terraform apply`.
  update_existing = true
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the user. Schema of `name` contains `family_name`
  and `given_name`.

* `primary_email` - (Required) Email of the user.

* `password` - (Optional) See the note on passwords above.

* `aliases` - (Optional) Alternative names for this user, expects a list of
  email addresses.

* `include_in_global_list` - (Optional) Boolean switch to show or hide this user
  in the global list. Defaults to true.

* `is_ip_whitelisted` - (Optional) Boolean switch to enforce whitelisting of the
  user's IP.

* `hash_function` - (Optional) `md5`, `sha-1` or `crypt`

* `posix_accounts` - (Optional) List with the following schema:
  * `account_id` - A POSIX account field identifier.
  * `gecos` - The GECOS (user information) for this account.
  * `gid` - The default group ID.
  * `home_directory` - The path to the home directory for this account.
  * `shell` - The path to the login shell for this account.
  * `system_id` - System identifier for which account Username or Uid apply to.
  * `primary` - If this is user's primary account within the SystemId.
  * `uid` - The POSIX compliant user ID.
  * `username` - The username of the account.

* `recovery_email` - (Optional) Recovery email of the user. Does not have to be
  in the domain.

* `recovery_phone` - (Optional) Recovery phone number of the user.

* `org_unit_path` - (Optional) Organizational unit path, defaults to `/`.

* `ssh_public_keys` - (Optional) SSH public keys of the user. Schema contains
  the following items:
  * `expiration_time_usec` - An expiration time in microseconds since epoch.
  * `key` - An SSH public key.

* `is_suspended` - (Optional) Suspend the user, defaults to false.

* `suspension_reason` - (Optional) Why is the user suspended?

* `custom_schema` - (Optional) See `user_custom_schema` for more details.

* `external_ids` - (Optional) List of `external_ids`. Schema contains:
  * `custom_type` - Custom type.
  * `type` - The type of the Id.
  * `value` - The value of the id.

* `update_existing` - (Optional) Boolean, defaults to false. Allows overwriting
  existing values instead of erroring out when a user already exists.

* `organizations` - (Optional) List of organizations. Schema of organization
  contains:
  * `cost_center` - The cost center of the users department.
  * `custom_type` - Custom type.
  * `department` - Department within the organization.
  * `description` - Description of the organization.
  * `domain` - The domain to which the organization belongs to.
  * `full_time_equivalent` - The full-time equivalent millipercent within the organization (100000 = 100%).
  * `location` - Location of the organization. This need not be fully qualified address.
  * `name` - Name of the organization.
  * `primary` - If it user's primary organization.
  * `symbol` - Symbol of the organization.
  * `title` - Title (designation) of the user in the organization.
  * `type` - Each entry can have a type which indicates standard types of
    that entry. For example organization could be of school, work etc. In
    addition to the standard type, an entry can have a custom type and
    can give it any name. Such types should have the CUSTOM value as type
    and also have a CustomType value.

## Attribute Reference

In addition to the above arguments, the following attributes are exported:

* `deletion_time` - User's G Suite account deletion time.

* `agreed_to_terms` - Indicates if user has agreed to terms.

* `creation_time` - User's G Suite account creation time.

* `customer_id` - CustomerId of User. 

* `etag` - ETag of the resource.

* `is_admin` - Boolean indicating if the user is admin.

* `is_delegated_admin` - Boolean indicating if the user is delegated admin.

* `2s_enforced` - Is 2-step verification enforced.

* `2s_enrolled` - Is enrolled in 2-step verification.

* `is_mailbox_setup` - Is mailbox setup.

* `last_login_time` - User's last login time.

## Import

A G Suite User can be imported using any key (`id`, `email`, `alias`), e.g.:

```
terraform import gsuite_user.developer "developer@domain.ext"
```
