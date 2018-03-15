// Contains functions that don't really belong anywhere else.

package gsuite

import (
	"fmt"
	"log"
	"time"

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
	return retryTime(retryFunc, 1)
}

func retryTime(retryFunc func() error, minutes int) error {
	return resource.Retry(time.Duration(minutes)*time.Minute, func() *resource.RetryError {
		err := retryFunc()
		if err == nil {
			return nil
		}
		if gerr, ok := err.(*googleapi.Error); ok && (gerr.Errors[0].Reason == "quotaExceeded" || gerr.Code == 429 || gerr.Code == 500 || gerr.Code == 502 || gerr.Code == 503) {
			return resource.RetryableError(gerr)
		}
		return resource.NonRetryableError(err)
	})
}
