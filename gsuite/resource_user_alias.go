package gsuite

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	admin "google.golang.org/api/admin/directory/v1"
)

func resourceUserAlias() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserAliasCreate,
		Read:   resourceUserAliasRead,
		Update: nil,
		Delete: resourceUserAliasDelete,
		Importer: &schema.ResourceImporter{
			State: resourceUserAliasImport,
		},
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Description: "ID (userKey) of the user the alias should be applied to.",
				Required:    true,
				ForceNew:    true,
			},
			"alias": {
				Type:         schema.TypeString,
				Description:  "Email alias which should be applied to the user.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateEmail,
			},
		},
	}
}

func resourceUserAliasCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	userId := d.Get("user_id").(string)
	setAlias := d.Get("alias").(string)

	alias := &admin.Alias{
		Alias: setAlias,
	}
	resp, err := config.directory.Users.Aliases.Insert(userId, alias).Do()
	if err != nil {
		return fmt.Errorf("[ERROR] failed to add alias for user (%s): %v", userId, err)
	}

	bOff := backoff.NewExponentialBackOff()
	bOff.MaxElapsedTime = time.Duration(config.TimeoutMinutes) * time.Minute
	bOff.InitialInterval = time.Second

	err = backoff.Retry(func() error {
		resp, err := config.directory.Users.Aliases.List(userId).Do()
		if err != nil {
			return backoff.Permanent(fmt.Errorf("[ERROR] could not retrieve aliases for user (%s): %v", userId, err))
		}

		_, ok := doesAliasExist(resp, setAlias)
		if ok {
			return nil
		}
		return fmt.Errorf(fmt.Sprintf("[WARN] no matching alias (%s) found for user (%s).", setAlias, userId))

	}, bOff)

	d.SetId(fmt.Sprintf("%s/%s", resp.PrimaryEmail, resp.Alias))
	return resourceUserAliasRead(d, meta)
}

func resourceUserAliasRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	userId := d.Get("user_id").(string)
	expectedAlias := d.Get("alias").(string)

	resp, err := config.directory.Users.Aliases.List(userId).Do()
	if err != nil {
		return fmt.Errorf("[ERROR] could not retrieve aliases for user (%s): %v", userId, err)
	}

	alias, ok := doesAliasExist(resp, expectedAlias)
	if !ok {
		log.Println(fmt.Sprintf("[WARN] no matching alias (%s) found for user (%s).", expectedAlias, userId))
		d.SetId("")
		return nil
	}
	d.SetId(fmt.Sprintf("%s/%s", userId, alias))
	d.Set("user_id", userId)
	d.Set("alias", alias)
	return nil
}

func resourceUserAliasDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	userId := d.Get("user_id").(string)
	alias := d.Get("alias").(string)

	err := config.directory.Users.Aliases.Delete(userId, alias).Do()
	if err != nil {
		return fmt.Errorf("[ERROR] unable to remove alias (%s) from user (%s): %v", alias, userId, err)
	}

	d.SetId("")
	return nil
}

func resourceUserAliasImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	userId := strings.Split(d.Id(), "/")[0]
	expectedAlias := strings.Split(d.Id(), "/")[1]

	resp, err := config.directory.Users.Aliases.List(userId).Do()
	if err != nil {
		return nil, fmt.Errorf("[ERROR] could not retrieve aliases for user (%s): %v", userId, err)
	}

	alias, ok := doesAliasExist(resp, expectedAlias)
	if !ok {
		return nil, fmt.Errorf("[ERROR] no matching alias (%s) found for user (%s).", expectedAlias, userId)
	}
	d.SetId(fmt.Sprintf("%s/%s", userId, alias))
	d.Set("user_id", userId)
	d.Set("alias", alias)

	return []*schema.ResourceData{d}, nil
}

func doesAliasExist(aliasesResp *admin.Aliases, expectedAlias string) (string, bool) {
	for _, aliasInt := range aliasesResp.Aliases {
		alias, ok := aliasInt.(map[string]interface{})
		if ok {
			value := alias["alias"].(string)
			if expectedAlias == alias["alias"].(string) {
				return value, true
			}
		}
		if !ok {
			log.Println(fmt.Sprintf("[ERROR] alias format in response did not match sdk struct, this indicates a probelm with provider or sdk handling: %v", reflect.TypeOf(aliasInt)))
			return "", false
		}
	}
	return "", false
}
