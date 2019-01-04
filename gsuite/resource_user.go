package gsuite

import (
	"fmt"
	"google.golang.org/api/googleapi"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	directory "google.golang.org/api/admin/directory/v1"
)

func flattenUserName(name *directory.UserName) map[string]interface{} {
	return map[string]interface{}{
		"family_name": name.FamilyName,
		"full_name":   name.FullName,
		"given_name":  name.GivenName,
	}
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
				Optional: true,
				Default:  false,
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
				Optional: true,
			},

			"etag": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"include_in_global_list": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"is_ip_whitelisted": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"family_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"full_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"given_name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// md5, sha-1 and crypt
			"hash_function": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"posix_accounts": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gecos": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"gid": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"home_directory": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"shell": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"system_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"primary": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"uid": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"username": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"primary_email": {
				Type:     schema.TypeString,
				Required: true,
				StateFunc: func(val interface{}) string {
					return strings.ToLower(val.(string))
				},
			},

			"ssh_public_keys": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expiration_time_usec": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"key": {
							Type:     schema.TypeString,
							Required: true,
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
				Optional: true,
				Default:  false,
			},

			"suspension_reason": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"custom_schema": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
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
		user.PrimaryEmail = strings.ToLower(v.(string))
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

	userSSHs := []*directory.UserSshPublicKey{}
	sshCount := d.Get("ssh_public_keys.#").(int)
	for i := 0; i < sshCount; i++ {
		sshConfig := d.Get(fmt.Sprintf("ssh_public_keys.%d", i)).(map[string]interface{})
		userSSH := &directory.UserSshPublicKey{}

		if v, ok := sshConfig["expiration_time_usec"]; ok {
			log.Printf("[DEBUG] Setting ssh %d expiration_time_usec: %v", i, int64(v.(int)))
			userSSH.ExpirationTimeUsec = int64(v.(int))
		}
		if v, ok := sshConfig["key"]; ok {
			log.Printf("[DEBUG] Setting ssh %d key: %s", i, v.(string))
			userSSH.Key = v.(string)
		}

		userSSHs = append(userSSHs, userSSH)
	}
	user.SshPublicKeys = userSSHs

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

	customSchemas := map[string]googleapi.RawMessage{}
	for i := 0; i < d.Get("custom_schema.#").(int); i++ {
		entry := d.Get(fmt.Sprintf("custom_schema.%d", i)).(map[string]interface{})
		customSchemas[entry["name"].(string)] = []byte(entry["value"].(string))
	}
	if len(customSchemas) > 0 {
		user.CustomSchemas = customSchemas
	}

	user.SshPublicKeys = userSSHs

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
		userSSHs := []*directory.UserSshPublicKey{}
		sshCount := d.Get("ssh_public_keys.#").(int)
		for i := 0; i < sshCount; i++ {
			sshConfig := d.Get(fmt.Sprintf("ssh_public_keys.%d", i)).(map[string]interface{})
			userSSH := &directory.UserSshPublicKey{}

			if v, ok := sshConfig["expiration_time_usec"]; ok {
				log.Printf("[DEBUG] Setting ssh %d expiration_time_usec: %v", i, int64(v.(int)))
				userSSH.ExpirationTimeUsec = int64(v.(int))
			}
			if v, ok := sshConfig["key"]; ok {
				log.Printf("[DEBUG] Setting ssh %d key: %s", i, v.(string))
				userSSH.Key = v.(string)
			}

			userSSHs = append(userSSHs, userSSH)
		}
		user.SshPublicKeys = userSSHs
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

	if d.HasChange("custom_schema") {
		customSchemas := map[string]googleapi.RawMessage{}
		for i := 0; i < d.Get("custom_schema.#").(int); i++ {
			entry := d.Get(fmt.Sprintf("custom_schema.%d", i)).(map[string]interface{})
			customSchemas[entry["name"].(string)] = []byte(entry["value"].(string))
		}
		user.CustomSchemas = customSchemas
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
		if e, ok := err.(*googleapi.Error); ok {
			panic(e.Body)
		}
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
	d.Set("custom_schemas", user.CustomSchemas)

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
	d.Set("custom_schemas", id.CustomSchemas)

	return []*schema.ResourceData{d}, nil
}
