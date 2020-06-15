---
layout: "gsuite"
page_title: "G Suite: gsuite_group_settings"
sidebar_current: "docs-gsuite-resource-group-settings"
description: |-
  Managing Settings of a G Suite Group.
---

# gsuite\_group\_settings

Provides a resource to create and manage Settings of a G Suite group.

**Note:** requires the `https://www.googleapis.com/auth/apps.groups.settings`
oauth scope.

## Example Usage

```hcl
resource "gsuite_group" "example" {
  email       = "example@domain.ext"
  name        = "example@domain.ext"
  description = "Example group"
}

resource "gsuite_group_settings" "example" {
  email = gsuite_group.example.email

  allow_external_members     = true
  allow_google_communication = false
  show_in_group_directory    = "true"
  who_can_discover_group     = "ALL_IN_DOMAIN_CAN_DISCOVER"
}
```

## Argument Reference

The following arguments are supported:

* `email` - (Required; Forces new resource) Email address of the G Suite
  group.

* `allow_external_members`

* `allow_web_posting`

* `archive_only`

* `custom_footer_text`

* `custom_reply_to`

* `description`

* `favorite_replies_on_top`

* `include_custom_footer`

* `include_in_global_address_list`

* `members_can_post_as_the_group`

* `message_moderation_level`

* `primary_language`

* `reply_to`

* `send_message_deny_notification`

* `spam_moderation_level`

* `who_can_approve_members`

* `who_can_assist_content`

* `who_can_contact_owner`

* `who_can_discover_group`

* `who_can_join`

* `who_can_leave_group`

* `who_can_moderate_content`

* `who_can_moderate_members`

* `who_can_post_message`

* `who_can_view_group`

* `who_can_view_membership`


## Attribute Reference

In addition to the above arguments, the following attributes are exported:

* `kind`

* `is_archived`

* `name`

* `description`

## Import

G Suite Group Settings can be imported using `group-email`, e.g.:

```
terraform import gsuite_group_settings.example "example@domain.ext"
```
