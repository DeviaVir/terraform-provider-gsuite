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

* `is_archived` - Allows the Group contents to be archived.

* `kind` - The type of the resource.

* `name` - Name of the group, which has a maximum size of 75 characters.

* `description` - Description of the group.

* `allow_external_members` - Identifies whether members external
  to your organization can join the group.

* `allow_google_communication` - Deprecated. Allows Google to contact
  administrator of the group.

* `allow_web_posting` - Allows posting from web.

* `archive_only` - Allows the group to be archived only.

* `custom_footer_text` - Set the content of custom footer text.

* `custom_reply_to` - An email address used when replying to a message
  if the replyTo property is set to REPLY_TO_CUSTOM. This address is defined
  by an account administrator.

* `favorite_replies_on_top` - Indicates if favorite replies should be
  displayed above other replies.

* `include_custom_footer` - Whether to include custom footer.

* `include_in_global_address_list` - Enables the group to be
  included in the Global Address List.

* `max_message_bytes` - Deprecated. The maximum size of a message is 25Mb.

* `members_can_post_as_the_group` - Enables members to post messages as the group.

* `message_display_font` - Deprecated. The default message display font
  always has a value of "DEFAULT_FONT".

* `message_moderation_level` - Moderation level of incoming messages.

* `primary_language` - The primary language for group. For a group's primary language use the language tags from
  the G Suite languages found at G Suite Email Settings API Email Language Tags.

* `reply_to` - Specifies who should the default reply go to.

* `send_message_deny_notification` - Allows a member to be notified if the
  member's message to the group is denied by the group owner.

* `show_in_group_directory` - Deprecated. This is merged into the new
  whoCanDiscoverGroup setting.

* `spam_moderation_level` - Specifies moderation levels for messages detected as spam.

* `who_can_add` - This is merged into the new whoCanModerateMembers setting.

* `who_can_add_references` - Deprecated. This functionality is no longer
   supported in the Google Groups UI.

* `who_can_approve_members` - Specifies who can approve members who ask to join groups.

* `who_can_approve_messages` - Deprecated. This is merged into the new
  whoCanModerateContent setting.

* `who_can_assign_topics` - Deprecated. This is merged into the new
  whoCanAssistContent setting.

* `who_can_assist_content` - Specifies who can moderate metadata.

* `who_can_ban_users` - Specifies who can deny membership to users.

* `who_can_contact_owner` - Permission to contact owner of the group via web UI.

* `who_can_delete_any_post` - Deprecated. This is merged into the new
  whoCanModerateContent setting.

* `who_can_delete_topics` - Deprecated. This is merged into the new
  whoCanModerateContent setting.

* `who_can_discover_group` - Specifies the set of users for whom this group
  is discoverable.

* `who_can_enter_free_form_tags` - Deprecated. This is merged into the new
  whoCanAssistContent setting.

* `who_can_hide_abuse` - Deprecated. This is merged into the new
  whoCanModerateContent setting.

* `who_can_invite` - Deprecated. This is merged into the new
  whoCanModerateMembers setting.

* `who_can_join` - Permission to join group.

* `who_can_leave_group` - Permission to leave the group.

* `who_can_lock_topics` - Deprecated. This is merged into the new
  whoCanModerateContent setting.

* `who_can_make_topics_sticky`  - Deprecated. This is merged into the new
  whoCanModerateContent setting.

* `who_can_mark_duplicate` - Deprecated. This is merged into the new
  whoCanAssistContent setting.

* `who_can_mark_favorite_reply_on_any_topic` - Deprecated. This is merged into
  the new whoCanAssistContent setting.

* `who_can_mark_favorite_reply_on_own_topic` - Deprecated. This is merged into
  the new whoCanAssistContent setting.

* `who_can_mark_no_response_needed` - Deprecated. This is merged into the new
  whoCanAssistContent setting.

* `who_can_moderate_content` - Specifies who can moderate content.

* `who_can_moderate_members` - Specifies who can manage members.

* `who_can_modify_members` - Deprecated. This is merged into the new
  whoCanModerateMembers setting.

* `who_can_modify_tags_and_categories` - Deprecated. This is merged into the
  new whoCanAssistContent setting.

* `who_can_move_topics_in` - Deprecated. This is merged into the new
  whoCanModerateContent setting.

* `who_can_move_topics_out` - Deprecated. This is merged into the new
  whoCanModerateContent setting.

* `who_can_post_announcements` - Deprecated. This is merged into the new
  whoCanModerateContent setting.

* `who_can_post_message` - Permissions to post messages.

* `who_can_take_topics` - Deprecated. This is merged into the new
  whoCanAssistContent setting.

* `who_can_unassign_topic` - Deprecated. This is merged into the new
  whoCanAssistContent setting.

* `who_can_unmark_favorite_reply_on_any_topic` - Deprecated. This is merged into
  the new whoCanAssistContent setting.

* `who_can_view_group` - Permissions to view group messages.

* `who_can_view_membership` - Permissions to view membership.
