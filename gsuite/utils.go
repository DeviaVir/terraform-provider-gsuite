// Contains functions that don't really belong anywhere else.

package gsuite

import (
	"fmt"
	"log"
	"math/rand"
	"net/mail"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/api/googleapi"
)

func handleNotFoundError(err error, d *schema.ResourceData, resource string) error {
	if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
		log.Printf("[WARN] Removing %s because it's gone", resource)
		// The resource doesn't exist anymore
		d.SetId("")

		return nil
	}

	return fmt.Errorf("Error reading %s: %s", resource, err)
}

func retry(retryFunc func() error, minutes int) error {
	return retryTime(retryFunc, minutes, false, false, false)
}

func retryNotFound(retryFunc func() error, minutes int) error {
	return retryTime(retryFunc, minutes, true, false, false)
}

func retryInvalid(retryFunc func() error, minutes int) error {
	return retryTime(retryFunc, minutes, false, false, true)
}

func retryPassDuplicate(retryFunc func() error, minutes int) error {
	return retryTime(retryFunc, minutes, true, true, false)
}

func retryTime(retryFunc func() error, minutes int, retryNotFound bool, retryPassDuplicate bool, retryInvalid bool) error {
	wait := 1
	return resource.Retry(time.Duration(minutes)*time.Minute, func() *resource.RetryError {
		err := retryFunc()
		if err == nil {
			return nil
		}

		rand.Seed(time.Now().UnixNano())
		randomNumberMiliseconds := rand.Intn(1001)

		if gerr, ok := err.(*googleapi.Error); ok {
			code := gerr.Code
			var reason string
			if len(gerr.Errors) > 0 {
				reason = gerr.Errors[0].Reason
			}

			if code == 500 || code == 502 || code == 503 {
				log.Printf("[DEBUG] Retrying server error code...")
				time.Sleep(time.Duration(wait)*time.Second + time.Duration(randomNumberMiliseconds))
				wait = wait * 2
				return resource.RetryableError(gerr)
			}
			if reason == "quotaExceeded" || code == 401 || code == 429 || (!retryPassDuplicate && code == 409) {
					log.Printf("[DEBUG] Retrying quota/server error code...")
					time.Sleep(time.Duration(wait)*time.Second + time.Duration(randomNumberMiliseconds))
					wait = wait * 2
					return resource.RetryableError(gerr)
			}
			if retryNotFound && code == 404 {
				log.Printf("[DEBUG] Retrying for eventual consistency...")
				time.Sleep(time.Duration(wait)*time.Second + time.Duration(randomNumberMiliseconds))
				wait = wait * 2
				return resource.RetryableError(gerr)
			}
			if retryInvalid && reason == "invalid" || code == 400 {
				log.Printf("[DEBUG] Retrying invalid error code...")
				time.Sleep(time.Duration(wait)*time.Second + time.Duration(randomNumberMiliseconds))
				wait = wait * 2
				return resource.RetryableError(gerr)
			}
		}

		// Deal with the broken API
		if strings.Contains(fmt.Sprintf("%s", err), "Invalid Input: Bad request for \"") && strings.Contains(fmt.Sprintf("%s", err), "\"code\":400") {
			log.Printf("[DEBUG] Retrying invalid response from API")
			return resource.RetryableError(err)
		}
		if strings.Contains(fmt.Sprintf("%s", err), "Service unavailable. Please try again") {
			log.Printf("[DEBUG] Retrying service unavailable from API")
			return resource.RetryableError(err)
		}
		if strings.Contains(fmt.Sprintf("%s", err), "Eventual consistency. Please try again") {
			log.Printf("[DEBUG] Retrying due to eventual consistency")
			return resource.RetryableError(err)
		}

		return resource.NonRetryableError(err)
	})
}

func mergeSchemas(a, b map[string]*schema.Schema) map[string]*schema.Schema {
	merged := make(map[string]*schema.Schema)

	for k, v := range a {
		merged[k] = v
	}

	for k, v := range b {
		merged[k] = v
	}

	return merged
}

func convertStringSet(set *schema.Set) []string {
	s := make([]string, 0, set.Len())
	for _, v := range set.List() {
		s = append(s, v.(string))
	}
	return s
}

func stringSliceDifference(left []string, right []string) []string {
	var d []string
	for _, l := range left {
		f := false
		for _, r := range right {
			if r == l {
				f = true
				break
			}
		}
		if !f {
			d = append(d, l)
		}
	}
	return d
}

func validateEmail(v interface{}, k string) (warnings []string, errors []error) {
	if v == nil || v.(string) == "" {
		return
	}
	email := v.(string)

	e, err := mail.ParseAddress(email)
	if err != nil {
		errors = append(errors,
			fmt.Errorf("unable to parse email address %s", email))
	}

	if e.Name != "" {
		errors = append(errors,
			fmt.Errorf("unexpected email format for %s expected an email format of myemail@domain.com", email))
	}

	parts := strings.Split(e.Address, "@")
	local := strings.Join(parts[0:len(parts)-1], "@")
	if len(local) > 63 {
		errors = append(errors,
			fmt.Errorf("local portion of email %s exceeds 63 characters", email))
	}

	return
}
