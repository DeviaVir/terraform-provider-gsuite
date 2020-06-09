package gsuite

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	directory "google.golang.org/api/admin/directory/v1"
)

func dataUser() *schema.Resource {
	return &schema.Resource{
		Read: dataUserRead,
		Schema: map[string]*schema.Schema{
			"primary_email": {
				Type:     schema.TypeString,
				Required: true,
				StateFunc: func(val interface{}) string {
					return strings.ToLower(val.(string))
				},
			},

			"org_unit_path": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"aliases": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"agreed_to_terms": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"change_password_next_login": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"customer_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"deletion_time": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"etag": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"include_in_global_list": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"is_ip_whitelisted": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"is_admin": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"is_delegated_admin": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"2s_enforced": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"2s_enrolled": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"is_mailbox_setup": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"last_login_time": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"family_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"full_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"given_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"password": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// md5, sha-1 and crypt
			"hash_function": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"posix_accounts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gecos": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"home_directory": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"shell": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"system_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"primary": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"uid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"ssh_public_keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expiration_time_usec": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fingerprint": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"is_suspended": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"suspension_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"recovery_email": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"recovery_phone": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"custom_schema": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"external_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"custom_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"organizations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cost_center": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"custom_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"department": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"full_time_equivalent": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"location": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"primary": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"symbol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"title": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataUserRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	var user *directory.User
	var err error
	err = retry(func() error {
		user, err = config.directory.Users.Get(d.Get("primary_email").(string)).Do()
		return err
	}, config.TimeoutMinutes)

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("User %q", d.Id()))
	}

	d.SetId(user.Id)
	d.Set("deletion_time", user.DeletionTime)
	d.Set("primary_email", user.PrimaryEmail)
	d.Set("org_unit_path", user.OrgUnitPath)
	d.Set("password", user.Password)
	d.Set("hash_function", user.HashFunction)
	d.Set("suspension_reason", user.SuspensionReason)
	d.Set("change_password_next_login", user.ChangePasswordAtNextLogin)
	d.Set("include_in_global_list", user.IncludeInGlobalAddressList)
	d.Set("is_ip_whitelisted", user.IpWhitelisted)
	d.Set("is_admin", user.IsAdmin)
	d.Set("is_delegated_admin", user.IsDelegatedAdmin)
	d.Set("is_suspended", user.Suspended)
	d.Set("2s_enrolled", user.IsEnrolledIn2Sv)
	d.Set("2s_enforced", user.IsEnforcedIn2Sv)
	d.Set("aliases", user.Aliases)
	d.Set("agreed_to_terms", user.AgreedToTerms)
	d.Set("creation_time", user.CreationTime)
	d.Set("customer_id", user.CustomerId)
	d.Set("etag", user.Etag)
	d.Set("last_login_time", user.LastLoginTime)
	d.Set("is_mailbox_setup", user.IsMailboxSetup)
	d.Set("recovery_email", user.RecoveryEmail)
	d.Set("recovery_phone", user.RecoveryPhone)
	d.Set("name", flattenUserName(user.Name))
	d.Set("posix_accounts", user.PosixAccounts)
	d.Set("ssh_public_keys", user.SshPublicKeys)
	d.Set("custom_schema", user.CustomSchemas)
	d.Set("external_ids", user.ExternalIds)
	d.Set("organizations", user.Organizations)

	return nil
}
