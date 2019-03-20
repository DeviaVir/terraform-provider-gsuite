// Contains functions that don't really belong anywhere else.

package gsuite

import (
	"fmt"
	"log"
	"time"
	"strings"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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

func retry(retryFunc func() error) error {
	return retryTime(retryFunc, 1, false)
}

func retryNotFound(retryFunc func() error) error {
	return retryTime(retryFunc, 1, true)
}

func retryTime(retryFunc func() error, minutes int, retryNotFound bool) error {
	return resource.Retry(time.Duration(minutes)*time.Minute, func() *resource.RetryError {
		err := retryFunc()
		if err == nil {
			return nil
		}
		if gerr, ok := err.(*googleapi.Error); ok && (gerr.Errors[0].Reason == "quotaExceeded" || gerr.Code == 409 || gerr.Code == 429 || gerr.Code == 500 || gerr.Code == 502 || gerr.Code == 503) {
			return resource.RetryableError(gerr)
		}
		if retryNotFound {
			if gerr, ok := err.(*googleapi.Error); ok && (gerr.Code == 404) {
				return resource.RetryableError(gerr)
			}
		}

		// Deal with the broken API
		if strings.Contains(fmt.Sprintf("%s", err), "Invalid Input: Bad request for ") && strings.Contains(fmt.Sprintf("%s", err), "\"code\":400") {
			log.Printf("[DEBUG] Retrying invalid response from API")
			return resource.RetryableError(err)
		}
		if strings.Contains(fmt.Sprintf("%s", err), "Service unavailable. Please try again") {
			log.Printf("[DEBUG] Retrying service unavailable from API")
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
