package gsuite

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataGroupSettings() *schema.Resource {
	return &schema.Resource{
		Read: dataGroupSettingsRead,
		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},

			"is_archived": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kind": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"allow_external_members": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"allow_google_communication": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"allow_web_posting": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"archive_only": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_footer_text": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_reply_to": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"favorite_replies_on_top": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"include_custom_footer": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"include_in_global_address_list": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"max_message_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"members_can_post_as_the_group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"message_display_font": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"message_moderation_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"primary_language": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reply_to": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"send_message_deny_notification": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"show_in_group_directory": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"spam_moderation_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"who_can_add": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"who_can_add_references": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"who_can_approve_members": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"who_can_approve_messages": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"who_can_assign_topics": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"who_can_assist_content": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"who_can_ban_users": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"who_can_contact_owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"who_can_delete_any_post": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"who_can_delete_topics": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"who_can_discover_group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"who_can_enter_free_form_tags": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"who_can_hide_abuse": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"who_can_invite": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"who_can_join": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"who_can_leave_group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"who_can_lock_topics": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"who_can_make_topics_sticky": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"who_can_mark_duplicate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"who_can_mark_favorite_reply_on_any_topic": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"who_can_mark_favorite_reply_on_own_topic": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"who_can_mark_no_response_needed": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"who_can_moderate_content": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"who_can_moderate_members": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"who_can_modify_members": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"who_can_modify_tags_and_categories": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"who_can_move_topics_in": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"who_can_move_topics_out": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"who_can_post_announcements": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"who_can_post_message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"who_can_take_topics": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"who_can_unassign_topic": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"who_can_unmark_favorite_reply_on_any_topic": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"who_can_view_group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"who_can_view_membership": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataGroupSettingsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	id, err := config.groupSettings.Groups.Get(d.Get("email").(string)).Do()
	if err != nil {
		return fmt.Errorf("[ERROR] Error fetching group settings. Make sure the group exists: %s ", err)
	}

	d.SetId(d.Get("email").(string))
	d.Set("allow_external_members", id.AllowExternalMembers)
	d.Set("allow_google_communication", id.AllowGoogleCommunication)
	d.Set("allow_web_posting", id.AllowWebPosting)
	d.Set("archive_only", id.ArchiveOnly)
	d.Set("custom_footer_text", id.CustomFooterText)
	d.Set("custom_reply_to", id.CustomReplyTo)
	d.Set("description", id.Description)
	d.Set("favorite_replies_on_top", id.FavoriteRepliesOnTop)
	d.Set("include_custom_footer", id.IncludeCustomFooter)
	d.Set("include_in_global_address_list", id.IncludeInGlobalAddressList)
	d.Set("max_message_bytes", id.MaxMessageBytes)
	d.Set("members_can_post_as_the_group", id.MembersCanPostAsTheGroup)
	d.Set("message_display_font", id.MessageDisplayFont)
	d.Set("message_moderation_level", id.MessageModerationLevel)
	d.Set("primary_language", id.PrimaryLanguage)
	d.Set("reply_to", id.ReplyTo)
	d.Set("send_message_deny_notification", id.SendMessageDenyNotification)
	d.Set("show_in_group_directory", id.ShowInGroupDirectory)
	d.Set("spam_moderation_level", id.SpamModerationLevel)
	d.Set("who_can_add", id.WhoCanAdd)
	d.Set("who_can_add_references", id.WhoCanAddReferences)
	d.Set("who_can_approve_members", id.WhoCanApproveMembers)
	d.Set("who_can_approve_messages", id.WhoCanApproveMessages)
	d.Set("who_can_assign_topics", id.WhoCanAssignTopics)
	d.Set("who_can_assist_content", id.WhoCanAssistContent)
	d.Set("who_can_ban_users", id.WhoCanBanUsers)
	d.Set("who_can_contact_owner", id.WhoCanContactOwner)
	d.Set("who_can_delete_any_post", id.WhoCanDeleteAnyPost)
	d.Set("who_can_delete_topics", id.WhoCanDeleteTopics)
	d.Set("who_can_discover_group", id.WhoCanDiscoverGroup)
	d.Set("who_can_enter_free_form_tags", id.WhoCanEnterFreeFormTags)
	d.Set("who_can_hide_abuse", id.WhoCanHideAbuse)
	d.Set("who_can_invite", id.WhoCanInvite)
	d.Set("who_can_join", id.WhoCanJoin)
	d.Set("who_can_leave_group", id.WhoCanLeaveGroup)
	d.Set("who_can_lock_topics", id.WhoCanLockTopics)
	d.Set("who_can_make_topics_sticky", id.WhoCanMakeTopicsSticky)
	d.Set("who_can_mark_duplicate", id.WhoCanMarkDuplicate)
	d.Set("who_can_mark_favorite_reply_on_any_topic", id.WhoCanMarkFavoriteReplyOnAnyTopic)
	d.Set("who_can_mark_favorite_reply_on_own_topic", id.WhoCanMarkFavoriteReplyOnOwnTopic)
	d.Set("who_can_mark_no_response_needed", id.WhoCanMarkNoResponseNeeded)
	d.Set("who_can_moderate_content", id.WhoCanModerateContent)
	d.Set("who_can_moderate_members", id.WhoCanModerateMembers)
	d.Set("who_can_modify_members", id.WhoCanModifyMembers)
	d.Set("who_can_modify_tags_and_categories", id.WhoCanModifyTagsAndCategories)
	d.Set("who_can_move_topics_in", id.WhoCanMoveTopicsIn)
	d.Set("who_can_move_topics_out", id.WhoCanMoveTopicsOut)
	d.Set("who_can_post_announcements", id.WhoCanPostAnnouncements)
	d.Set("who_can_post_message", id.WhoCanPostMessage)
	d.Set("who_can_take_topics", id.WhoCanTakeTopics)
	d.Set("who_can_unassign_topic", id.WhoCanUnassignTopic)
	d.Set("who_can_unmark_favorite_reply_on_any_topic", id.WhoCanUnmarkFavoriteReplyOnAnyTopic)
	d.Set("who_can_view_group", id.WhoCanViewGroup)
	d.Set("who_can_view_membership", id.WhoCanViewMembership)

	return nil
}
