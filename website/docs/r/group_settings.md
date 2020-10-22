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

* `allow_external_members` - (Optional) Identifies whether members external
  to your organization can join the group.
  Valid values are `true` or `false`. Defaults to `false`.

* `allow_web_posting` - (Optional) Allows posting from web.
  Valid values are `true` or `false`. Defaults to `true`.

* `archive_only` - (Optional) Allows the group to be archived only.
  Valid values are `true` or `false`. Defaults to `false`.
  If true, the `whoCanPostMessage` property is set to `NONE_CAN_POST`.
  When false, updating `whoCanPostMessage` to `NONE_CAN_POST`, results in an error.

* `custom_footer_text` - (Optional) Set the content of custom footer text.
  The maximum number of characters is 1,000.

* `custom_reply_to` - (Optional) An email address used when replying to a message
  if the replyTo property is set to REPLY_TO_CUSTOM. This address is defined
  by an account administrator.

* `description` - (Optional) A longer, human-readable description for the group.

* `favorite_replies_on_top` - (Optional) Indicates if favorite replies should be
  displayed above other replies.
  Valid values are `true` or `false`. Defaults to `true`.

* `include_custom_footer` - (Optional) Whether to include custom footer. 
  Valid values are `true` or `false`. Defaults to `false`.

* `include_in_global_address_list` - (Optional) Enables the group to be
  included in the Global Address List. For more information, see the help center.
  Valid values are `true` or `false`. Defaults to `true`.

* `members_can_post_as_the_group` - (Optional) Enables members to post messages as the group.
  Valid values are `true` or `false`. Defaults to `false`.

* `message_moderation_level` - (Optional) Moderation level of incoming messages.
  The valid values are `MODERATE_ALL_MESSAGES`, `MODERATE_NON_MEMBERS`, `MODERATE_NEW_MEMBERS` and `MODERATE_NONE`. Defaults to `MODERATE_NONE`.

* `primary_language` - (Optional) The primary language for group. For a group's primary language use the language tags from
  the G Suite languages found at G Suite Email Settings API Email Language Tags.

* `reply_to` - (Optional) Specifies who should the default reply go to.
  The valid values are `REPLY_TO_CUSTOM`, `REPLY_TO_SENDER`, `REPLY_TO_LIST`, `REPLY_TO_OWNER`, `REPLY_TO_IGNORE` and `REPLY_TO_MANAGERS`. Defaults to `REPLY_TO_IGNORE`.

* `send_message_deny_notification` - (Optional) Allows a member to be notified if the
  member's message to the group is denied by the group owner.
  Valid values are `true` or `false`. Defaults to `false`.

* `spam_moderation_level` - (Optional) Specifies moderation levels for messages detected as spam.
  The valid values are `ALLOW`, `MODERATE`, `SILENTLY_MODERATE` and `REJECT`. Defaults to `MODERATE`.

* `who_can_approve_members` - (Optional) Specifies who can approve members who ask to
  join groups. This permission will be deprecated once it is merged
  into the new whoCanModerateMembers setting.
  The valid values are `ALL_OWNERS_CAN_APPROVE`, `ALL_MANAGERS_CAN_APPROVE`, `ALL_MEMBERS_CAN_APPROVE` and `NONE_CAN_APPROVE`. Defaults to `ALL_MANAGERS_CAN_APPROVE`.

* `who_can_assist_content` - (Optional) Specifies who can moderate metadata.
  The valid values are `NONE`, `OWNERS_ONLY`, `MANAGERS_ONLY`, `OWNERS_AND_MANAGERS` and `ALL_MEMBERS`. Defaults to `NONE`.

* `who_can_contact_owner` - (Optional) Permission to contact owner of the group via web UI.
  The valid values are `ANYONE_CAN_CONTACT`, `ALL_IN_DOMAIN_CAN_CONTACT`, `ALL_MEMBERS_CAN_CONTACT` and `ALL_MANAGERS_CAN_CONTACT`. Defaults to `ANYONE_CAN_CONTACT`.

* `who_can_discover_group` - (Optional) Specifies the set of users for whom this group
  is discoverable.
  The valid values are `ALL_MEMBERS_CAN_DISCOVER`, `ALL_IN_DOMAIN_CAN_DISCOVER` and `ANYONE_CAN_DISCOVER`. Defaults to `ALL_MEMBERS_CAN_DISCOVER`.

* `who_can_join` - (Optional) Permission to join group. 
  The valid values are `ANYONE_CAN_JOIN`, `ALL_IN_DOMAIN_CAN_JOIN`, `INVITED_CAN_JOIN` and `CAN_REQUEST_TO_JOIN`. Defaults to `CAN_REQUEST_TO_JOIN`.

* `who_can_leave_group` - (Optional) Permission to leave the group.
  The valid values are `ALL_MANAGERS_CAN_LEAVE`, `ALL_OWNERS_CAN_LEAVE`, `ALL_MEMBERS_CAN_LEAVE` and `NONE_CAN_LEAVE`. Defaults to `ALL_MEMBERS_CAN_LEAVE`.

* `who_can_moderate_content` - (Optional) Specifies who can moderate content.
  The valid values are `NONE`, `OWNERS_ONLY`, `OWNERS_AND_MANAGERS` and `ALL_MEMBERS`. Defaults to `OWNERS_AND_MANAGERS`.

* `who_can_moderate_members` - (Optional) Specifies who can manage members.
  The valid values are `NONE`, `OWNERS_ONLY`, `OWNERS_AND_MANAGERS` and `ALL_MEMBERS`. Defaults to `OWNERS_AND_MANAGERS`.

* `who_can_post_message` - (Optional) Permissions to post messages.
  The valid values are `NONE_CAN_POST`, `ALL_MANAGERS_CAN_POST`, `ALL_MEMBERS_CAN_POST`, `ALL_OWNERS_CAN_POST`, `ALL_IN_DOMAIN_CAN_POST` and `ANYONE_CAN_POST`.Defaults to `ANYONE_CAN_POST`.

* `who_can_view_group` - (Optional) Permissions to view group messages.
  The valid values are `ANYONE_CAN_VIEW`, `ALL_IN_DOMAIN_CAN_VIEW`, `ALL_MEMBERS_CAN_VIEW`, `ALL_MANAGERS_CAN_VIEW` and `ALL_OWNERS_CAN_VIEW`. Defaults to `ALL_MEMBERS_CAN_VIEW`.

* `who_can_view_membership` - (Optional) Permissions to view membership.
  The valid values are `ALL_IN_DOMAIN_CAN_VIEW`, `ALL_MEMBERS_CAN_VIEW`, `ALL_MANAGERS_CAN_VIEW` and `ALL_OWNERS_CAN_VIEW`. Defaults to `ALL_MEMBERS_CAN_VIEW`.


## Attribute Reference

In addition to the above arguments, the following attributes are exported:

* `kind` - The type of the resource. It is always groupsSettings#groups.

* `is_archived` - Allows the Group contents to be archived.
  Valid values are `true` or `false`.

* `name` - Name of the group, which has a maximum size of 75 characters.

* `description` - Description of the group. This property value may be an empty
  string if no group description has been entered. If entered, the maximum group
  description is no more than 300 characters. 

## Import

G Suite Group Settings can be imported using `group-email`, e.g.:

```
terraform import gsuite_group_settings.example "example@domain.ext"
```
