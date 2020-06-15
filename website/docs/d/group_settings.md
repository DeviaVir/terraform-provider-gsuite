---
layout: "gsuite"
page_title: "G Suite: group_settings data source"
sidebar_current: "docs-gsuite-datasource-group-settings"
description: |-
  Retrieves the settings of a Group in G Suite.
---

# gsuite\_group\_settings

Reads the Settings of a Group from G Suite

## Example Usage

```hcl
data "gsuite_group" "example" {
  email = "example@domain.ext"
}

data "gsuite_group_settings" "example" {
  email = gsuite_group.example.email
}

output "group-settings" {
  value = data.gsuite_group_settings.example
}
```

## Argument Reference

The following arguments are supported:

* `email` - (Required) The email of group to retrieve the Settings for.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `is_archived`

* `kind`

* `name`

* `description`

* `allow_external_members`

* `allow_google_communication`

* `allow_web_posting`

* `archive_only`

* `custom_footer_text`

* `custom_reply_to`

* `favorite_replies_on_top`

* `include_custom_footer`

* `include_in_global_address_list`

* `max_message_bytes`

* `members_can_post_as_the_group`

* `message_display_font`

* `message_moderation_level`

* `primary_language`

* `reply_to`

* `send_message_deny_notification`

* `show_in_group_directory`

* `spam_moderation_level`

* `who_can_add`

* `who_can_add_references`

* `who_can_approve_members`

* `who_can_approve_messages`

* `who_can_assign_topics`

* `who_can_assist_content`

* `who_can_ban_users`

* `who_can_contact_owner`

* `who_can_delete_any_post`

* `who_can_delete_topics`

* `who_can_discover_group`

* `who_can_enter_free_form_tags`

* `who_can_hide_abuse`

* `who_can_invite`

* `who_can_join`

* `who_can_leave_group`

* `who_can_lock_topics`

* `who_can_make_topics_sticky`

* `who_can_mark_duplicate`

* `who_can_mark_favorite_reply_on_any_topic`

* `who_can_mark_favorite_reply_on_own_topic`

* `who_can_mark_no_response_needed`

* `who_can_moderate_content`

* `who_can_moderate_members`

* `who_can_modify_members`

* `who_can_modify_tags_and_categories`

* `who_can_move_topics_in`

* `who_can_move_topics_out`

* `who_can_post_announcements`

* `who_can_post_message`

* `who_can_take_topics`

* `who_can_unassign_topic`

* `who_can_unmark_favorite_reply_on_any_topic`

* `who_can_view_group`

* `who_can_view_membership`
