package gsuite

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"

	directory "google.golang.org/api/admin/directory/v1"
)

var googleLookup = map[string]string{
	"aliases":                    "Aliases",
	"agreed_to_terms":            "AgreedToTerms",
	"change_password_next_login": "ChangePasswordAtNextLogin",
	"creation_time":              "CreationTime",
	"customer_id":                "CustomerId",
	"deletion_time":              "DeletionTime",
	"etag":                       "Etag",
	"include_in_global_list": "IncludeInGlobalAddressList",
	"is_ip_whitelisted":      "IpWhitelisted",
	"is_admin":               "IsAdmin",
	"is_delegated_admin":     "IsDelegatedAdmin",
	"2s_enforced":            "IsEnforcedIn2Sv",
	"2s_enrolled":            "IsEnrolledIn2Sv",
	"is_mailbox_setup":       "IsMailboxSetup",
	"last_login_time":        "LastLoginTime",
	"password":               "Password",
	"hash_function":          "HashFunction",
	"primary_email":          "PrimaryEmail",
	"is_suspended":           "Suspended",
	"suspension_reason":      "SuspensionReason",
}

func flattenUserName(name *directory.UserName) map[string]interface{} {
	return map[string]interface{}{
		"family_name": name.FamilyName,
		"full_name":   name.FullName,
		"given_name":  name.GivenName,
	}
}

func flattenUserPosixAccounts(posixAccounts []*directory.UserPosixAccount) []map[string]interface{} {
	result := make([]map[string]interface{}, len(posixAccounts))
	for i, posixAccount := range posixAccounts {
		result[i] = map[string]interface{}{
			"account_id":     posixAccount.AccountId,
			"gecos":          posixAccount.Gecos,
			"gid":            posixAccount.Gid,
			"home_directory": posixAccount.HomeDirectory,
			"system_id":      posixAccount.SystemId,
			"primary":        posixAccount.Primary,
			"shell":          posixAccount.Shell,
			"uid":            posixAccount.Uid,
			"username":       posixAccount.Username,
		}
	}
	return result
}

func flattenUserSshPublicKeys(sshPublicKeys []*directory.UserSshPublicKey) []map[string]interface{} {
	result := make([]map[string]interface{}, len(sshPublicKeys))
	for i, sshPublicKey := range sshPublicKeys {
		result[i] = map[string]interface{}{
			"expiration_time_usec": sshPublicKey.ExpirationTimeUsec,
			"key":         sshPublicKey.Key,
			"fingerprint": sshPublicKey.Fingerprint,
		}
	}
	return result
}

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			State: resourceUserImporter,
		},

		Schema: map[string]*schema.Schema{
			"aliases": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"agreed_to_terms": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"change_password_next_login": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"creation_time": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"customer_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"deletion_time": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"etag": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"include_in_global_list": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"is_ip_whitelisted": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"is_admin": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},

			"is_delegated_admin": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},

			"2s_enforced": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},

			"2s_enrolled": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},

			"is_mailbox_setup": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},

			"last_login_time": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"family_name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"full_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"given_name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			// md5, sha-1 and crypt
			"hash_function": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"posix_accounts": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"gecos": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"gid": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"home_directory": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"shell": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"system_id": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"primary": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"uid": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"username": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"primary_email": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"ssh_public_keys": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expiration_time_usec": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"key": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"fingerprint": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"is_suspended": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"suspension_reason": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	user := &directory.User{}

	if v, ok := d.GetOk("deletion_time"); ok {
		log.Printf("[DEBUG] Setting %s: %s", "deletion_time", v.(string))
		user.DeletionTime = v.(string)
	}
	if v, ok := d.GetOk("primary_email"); ok {
		log.Printf("[DEBUG] Setting %s: %s", "primary_email", v.(string))
		user.PrimaryEmail = v.(string)
	}
	if v, ok := d.GetOk("password"); ok {
		log.Printf("[DEBUG] Setting %s: %s", "password", v.(string))
		user.Password = v.(string)
	}
	if v, ok := d.GetOk("hash_function"); ok {
		log.Printf("[DEBUG] Setting %s: %s", "hash_function", v.(string))
		user.HashFunction = v.(string)
	}
	if v, ok := d.GetOk("suspension_reason"); ok {
		log.Printf("[DEBUG] Setting %s: %s", "suspension_reason", v.(string))
		user.SuspensionReason = v.(string)
	}

	if v, ok := d.GetOk("change_password_next_login"); ok {
		log.Printf("[DEBUG] Setting %s: %t", "change_password_next_login", v.(bool))
		user.ChangePasswordAtNextLogin = v.(bool)
	}
	if v, ok := d.GetOk("include_in_global_list"); ok {
		log.Printf("[DEBUG] Setting %s: %t", "include_in_global_list", v.(bool))
		user.IncludeInGlobalAddressList = v.(bool)
	}
	if v, ok := d.GetOk("is_ip_whitelisted"); ok {
		log.Printf("[DEBUG] Setting %s: %t", "is_ip_whitelisted", v.(bool))
		user.IpWhitelisted = v.(bool)
	}
	if v, ok := d.GetOk("is_suspended"); ok {
		log.Printf("[DEBUG] Setting %s: %t", "is_suspended", v.(bool))
		user.Suspended = v.(bool)
	}

	userSshs := []*directory.UserSshPublicKey{}
	sshCount := d.Get("ssh_public_keys.#").(int)
	for i := 0; i < sshCount; i++ {
		sshConfig := d.Get(fmt.Sprintf("ssh_public_keys.%d", i)).(map[string]interface{})
		userSsh := &directory.UserSshPublicKey{}

		if v, ok := sshConfig["expiration_time_usec"]; ok {
			log.Printf("[DEBUG] Setting ssh %d expiration_time_usec: %v", i, int64(v.(int)))
			userSsh.ExpirationTimeUsec = int64(v.(int))
		}
		if v, ok := sshConfig["key"]; ok {
			log.Printf("[DEBUG] Setting ssh %d key: %s", i, v.(string))
			userSsh.Key = v.(string)
		}

		userSshs = append(userSshs, userSsh)
	}
	user.SshPublicKeys = userSshs

	userPosixs := []*directory.UserPosixAccount{}
	posixCount := d.Get("posix_accounts.#").(int)
	for i := 0; i < posixCount; i++ {
		posixConfig := d.Get(fmt.Sprintf("posix_accounts.%d", i)).(map[string]interface{})
		userPosix := &directory.UserPosixAccount{}

		if posixConfig["gecos"] != "" {
			log.Printf("[DEBUG] Setting posix %d gecos: %s", i, posixConfig["gecos"].(string))
			userPosix.Gecos = posixConfig["gecos"].(string)
		}
		if posixConfig["gid"] != 0 {
			log.Printf("[DEBUG] Setting posix %d gid: %d", i, uint64(posixConfig["gid"].(int)))
			userPosix.Gid = uint64(posixConfig["gid"].(int))
		}
		if posixConfig["home_directory"] != "" {
			log.Printf("[DEBUG] Setting posix %d home_directory: %s", i, posixConfig["home_directory"].(string))
			userPosix.HomeDirectory = posixConfig["home_directory"].(string)
		}
		if posixConfig["system_id"] != "" {
			log.Printf("[DEBUG] Setting posix %d system_id: %s", i, posixConfig["system_id"].(string))
			userPosix.SystemId = posixConfig["system_id"].(string)
		}
		if posixConfig["shell"] != "" {
			log.Printf("[DEBUG] Setting posix %d shell: %s", i, posixConfig["shell"].(string))
			userPosix.Shell = posixConfig["shell"].(string)
		}
		if posixConfig["primary"] != "" {
			log.Printf("[DEBUG] Setting posix %d primary: %t", i, posixConfig["primary"].(bool))
			userPosix.Primary = posixConfig["primary"].(bool)
		}
		if posixConfig["uid"] != 0 {
			log.Printf("[DEBUG] Setting posix %d uid: %d", i, uint64(posixConfig["uid"].(int)))
			userPosix.Uid = uint64(posixConfig["uid"].(int))
		}
		if posixConfig["username"] != "" {
			log.Printf("[DEBUG] Setting posix %d username: %s", i, posixConfig["username"].(string))
			userPosix.Username = posixConfig["username"].(string)
		}

		userPosixs = append(userPosixs, userPosix)
	}
	user.PosixAccounts = userPosixs

	userNamePrefix := "name.0"
	userName := &directory.UserName{
		FamilyName: d.Get(userNamePrefix + ".family_name").(string),
		GivenName:  d.Get(userNamePrefix + ".given_name").(string),
	}
	user.Name = userName

	var createdUser *directory.User
	var err error
	err = retry(func() error {
		createdUser, err = config.directory.Users.Insert(user).Do()
		return err
	})

	if err != nil {
		return fmt.Errorf("Error creating user: %s", err)
	}

	d.SetId(createdUser.Id)
	log.Printf("[INFO] Created user: %s", createdUser.PrimaryEmail)
	return resourceUserRead(d, meta)
}

func resourceUserUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	user := &directory.User{}
	nullFields := []string{}

	if d.HasChange("deletion_time") {
		if v, ok := d.GetOk("deletion_time"); ok {
			log.Printf("[DEBUG] Updating user deletion_time: %s", d.Get("deletion_time").(string))
			user.DeletionTime = v.(string)
		} else {
			log.Printf("[DEBUG] Removing user deletion_time")
			user.DeletionTime = ""
			nullFields = append(nullFields, "deletion_time")
		}
	}
	if d.HasChange("primary_email") {
		if v, ok := d.GetOk("primary_email"); ok {
			log.Printf("[DEBUG] Updating user primary_email: %s", d.Get("primary_email").(string))
			user.PrimaryEmail = v.(string)
		} else {
			log.Printf("[DEBUG] Removing user primary_email")
			user.PrimaryEmail = ""
			nullFields = append(nullFields, "primary_email")
		}
	}
	if d.HasChange("password") {
		if v, ok := d.GetOk("password"); ok {
			log.Printf("[DEBUG] Updating user password: %s", d.Get("password").(string))
			user.Password = v.(string)
		} else {
			log.Printf("[DEBUG] Removing user password")
			user.Password = ""
			nullFields = append(nullFields, "password")
		}
	}
	if d.HasChange("hash_function") {
		if v, ok := d.GetOk("hash_function"); ok {
			log.Printf("[DEBUG] Updating user hash_function: %s", d.Get("hash_function").(string))
			user.HashFunction = v.(string)
		} else {
			log.Printf("[DEBUG] Removing user hash_function")
			user.HashFunction = ""
			nullFields = append(nullFields, "hash_function")
		}
	}
	if d.HasChange("suspension_reason") {
		if v, ok := d.GetOk("suspension_reason"); ok {
			log.Printf("[DEBUG] Updating user suspension_reason: %s", d.Get("suspension_reason").(string))
			user.SuspensionReason = v.(string)
		} else {
			log.Printf("[DEBUG] Removing user suspension_reason")
			user.SuspensionReason = ""
			nullFields = append(nullFields, "suspension_reason")
		}
	}

	if d.HasChange("change_password_next_login") {
		if v, ok := d.GetOk("change_password_next_login"); ok {
			log.Printf("[DEBUG] Updating user change_password_next_login: %t", d.Get("change_password_next_login").(bool))
			user.ChangePasswordAtNextLogin = v.(bool)
		} else {
			log.Printf("[DEBUG] Removing user change_password_next_login")
			user.ChangePasswordAtNextLogin = false
			nullFields = append(nullFields, "change_password_next_login")
		}
	}
	if d.HasChange("include_in_global_list") {
		if v, ok := d.GetOk("include_in_global_list"); ok {
			log.Printf("[DEBUG] Updating user include_in_global_list: %t", d.Get("include_in_global_list").(bool))
			user.IncludeInGlobalAddressList = v.(bool)
		} else {
			log.Printf("[DEBUG] Removing user include_in_global_list")
			user.IncludeInGlobalAddressList = true
			nullFields = append(nullFields, "include_in_global_list")
		}
	}
	if d.HasChange("is_ip_whitelisted") {
		if v, ok := d.GetOk("is_ip_whitelisted"); ok {
			log.Printf("[DEBUG] Updating user is_ip_whitelisted: %t", d.Get("is_ip_whitelisted").(bool))
			user.IpWhitelisted = v.(bool)
		} else {
			log.Printf("[DEBUG] Removing user is_ip_whitelisted")
			user.IpWhitelisted = false
			nullFields = append(nullFields, "is_ip_whitelisted")
		}
	}
	if d.HasChange("is_suspended") {
		if v, ok := d.GetOk("is_suspended"); ok {
			log.Printf("[DEBUG] Updating user is_suspended: %t", d.Get("is_suspended").(bool))
			user.Suspended = v.(bool)
		} else {
			log.Printf("[DEBUG] Removing user is_suspended")
			user.Suspended = false
			nullFields = append(nullFields, "is_suspended")
		}
	}

	if d.HasChange("ssh_public_keys") {
		userSshs := []*directory.UserSshPublicKey{}
		sshCount := d.Get("ssh_public_keys.#").(int)
		for i := 0; i < sshCount; i++ {
			sshConfig := d.Get(fmt.Sprintf("ssh_public_keys.%d", i)).(map[string]interface{})
			userSsh := &directory.UserSshPublicKey{}

			if v, ok := sshConfig["expiration_time_usec"]; ok {
				log.Printf("[DEBUG] Setting ssh %d expiration_time_usec: %v", i, int64(v.(int)))
				userSsh.ExpirationTimeUsec = int64(v.(int))
			}
			if v, ok := sshConfig["key"]; ok {
				log.Printf("[DEBUG] Setting ssh %d key: %s", i, v.(string))
				userSsh.Key = v.(string)
			}

			userSshs = append(userSshs, userSsh)
		}
		user.SshPublicKeys = userSshs
	}

	if d.HasChange("posix_accounts") {
		userPosixs := []*directory.UserPosixAccount{}
		posixCount := d.Get("posix_accounts.#").(int)
		for i := 0; i < posixCount; i++ {
			posixConfig := d.Get(fmt.Sprintf("posix_accounts.%d", i)).(map[string]interface{})
			userPosix := &directory.UserPosixAccount{}

			if posixConfig["gecos"] != "" {
				log.Printf("[DEBUG] Setting posix %d gecos: %s", i, posixConfig["gecos"].(string))
				userPosix.Gecos = posixConfig["gecos"].(string)
			}
			if posixConfig["gid"] != 0 {
				log.Printf("[DEBUG] Setting posix %d gid: %d", i, uint64(posixConfig["gid"].(int)))
				userPosix.Gid = uint64(posixConfig["gid"].(int))
			}
			if posixConfig["home_directory"] != "" {
				log.Printf("[DEBUG] Setting posix %d home_directory: %s", i, posixConfig["home_directory"].(string))
				userPosix.HomeDirectory = posixConfig["home_directory"].(string)
			}
			if posixConfig["system_id"] != "" {
				log.Printf("[DEBUG] Setting posix %d system_id: %s", i, posixConfig["system_id"].(string))
				userPosix.SystemId = posixConfig["system_id"].(string)
			}
			if posixConfig["shell"] != "" {
				log.Printf("[DEBUG] Setting posix %d shell: %s", i, posixConfig["shell"].(string))
				userPosix.Shell = posixConfig["shell"].(string)
			}
			if posixConfig["primary"] != "" {
				log.Printf("[DEBUG] Setting posix %d primary: %t", i, posixConfig["primary"].(bool))
				userPosix.Primary = posixConfig["primary"].(bool)
			}
			if posixConfig["uid"] != 0 {
				log.Printf("[DEBUG] Setting posix %d uid: %d", i, uint64(posixConfig["uid"].(int)))
				userPosix.Uid = uint64(posixConfig["uid"].(int))
			}
			if posixConfig["username"] != "" {
				log.Printf("[DEBUG] Setting posix %d username: %s", i, posixConfig["username"].(string))
				userPosix.Username = posixConfig["username"].(string)
			}

			userPosixs = append(userPosixs, userPosix)
		}
		user.PosixAccounts = userPosixs
	}

	userNamePrefix := "name.0"
	userName := &directory.UserName{
		FamilyName: d.Get(userNamePrefix + ".family_name").(string),
		GivenName:  d.Get(userNamePrefix + ".given_name").(string),
	}
	user.Name = userName

	if len(nullFields) > 0 {
		user.NullFields = nullFields
	}

	var updatedUser *directory.User
	var err error
	err = retry(func() error {
		updatedUser, err = config.directory.Users.Update(d.Id(), user).Do()
		return err
	})

	if err != nil {
		return fmt.Errorf("Error updating user: %s", err)
	}

	log.Printf("[INFO] Updated user: %s", updatedUser.PrimaryEmail)
	return resourceUserRead(d, meta)
}

func resourceUserRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	var user *directory.User
	var err error
	err = retry(func() error {
		user, err = config.directory.Users.Get(d.Id()).Do()
		return err
	})

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("User %q", d.Get("name").(string)))
	}

	d.SetId(user.Id)
	d.Set("deletion_time", user.DeletionTime)
	d.Set("primary_email", user.PrimaryEmail)
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

	d.Set("name", flattenUserName(user.Name))
	d.Set("posix_accounts", user.PosixAccounts)
	d.Set("ssh_public_keys", user.SshPublicKeys)

	return nil
}

func resourceUserDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	var err error
	err = retry(func() error {
		err = config.directory.Users.Delete(d.Id()).Do()
		return err
	})
	if err != nil {
		return fmt.Errorf("Error deleting user: %s", err)
	}

	d.SetId("")
	return nil
}

// Allow importing using any key (id, email, alias)
func resourceUserImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	id, err := config.directory.Users.Get(d.Id()).Do()

	if err != nil {
		return nil, fmt.Errorf("Error fetching user. Make sure the user exists: %s ", err)
	}

	d.SetId(id.Id)
	d.Set("deletion_time", id.DeletionTime)
	d.Set("primary_email", id.PrimaryEmail)
	d.Set("password", id.Password)
	d.Set("hash_function", id.HashFunction)
	d.Set("suspension_reason", id.SuspensionReason)
	d.Set("change_password_next_login", id.ChangePasswordAtNextLogin)
	d.Set("include_in_global_list", id.IncludeInGlobalAddressList)
	d.Set("is_ip_whitelisted", id.IpWhitelisted)
	d.Set("is_admin", id.IsAdmin)
	d.Set("is_delegated_admin", id.IsDelegatedAdmin)
	d.Set("is_suspended", id.Suspended)
	d.Set("2s_enrolled", id.IsEnrolledIn2Sv)
	d.Set("2s_enforced", id.IsEnforcedIn2Sv)
	d.Set("aliases", id.Aliases)
	d.Set("agreed_to_terms", id.AgreedToTerms)
	d.Set("creation_time", id.CreationTime)
	d.Set("customer_id", id.CustomerId)
	d.Set("etag", id.Etag)
	d.Set("last_login_time", id.LastLoginTime)
	d.Set("is_mailbox_setup", id.IsMailboxSetup)

	d.Set("name", flattenUserName(id.Name))
	d.Set("posix_accounts", id.PosixAccounts)
	d.Set("ssh_public_keys", id.SshPublicKeys)

	return []*schema.ResourceData{d}, nil
}
