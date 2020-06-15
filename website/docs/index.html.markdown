---
layout: "gsuite"
page_title: "Provider: G Suite"
sidebar_current: "docs-gsuite-index"
description: |-
  The G Suite provider allows Terraform to read from, write to, and configure Google G Suite
---

# G Suite Provider

The G Suite provider allows Terraform to read from, write to, and configure
[Google G Suite](https://gsuite.google.com/).

There are two possible authentication mechanisms for using this provider. Using
a service account, or a personal admin account. The latter requires user
interaction, whereas a service account could be used in an automated workflow.

## Best practices

Using a personal admin account should ideally only be used for testing.

Follow these [instructions](https://developers.google.com/admin-sdk/directory/v1/guides/delegation)
to set up a service account for use in this provider.

We advise to at a minimum set the following oauth scopes:

* `https://www.googleapis.com/auth/admin.directory.group`
* `https://www.googleapis.com/auth/admin.directory.user`

When setting oauth scopes, the scopes need to be set in both the G Suite
service account settings, and in this provider's `oauth_scopes` parameter.

### Relevant Google Admin SDK Documentation

#### General

* http://google.golang.org/api/admin/directory/v1
* https://developers.google.com/admin-sdk/directory/v1/reference/

#### Schema Types

* https://developers.google.com/admin-sdk/directory/v1/reference/users
* https://developers.google.com/admin-sdk/directory/v1/reference/groups
* https://developers.google.com/admin-sdk/directory/v1/reference/schemas

## Provider Arguments

The provider configuration block accepts the following arguments.
In most cases it is recommended to set them via the indicated environment
variables in order to keep credential information out of the configuration.

* `credentials` - (Optional) Path to or string content of your credentials. If
  you have authenticated using `gcloud auth login` and want to test using your
  personal account you may leave this empty. May be set via the
  `GOOGLE_CREDENTIALS`, `GOOGLE_CLOUD_KEYFILE_JSON`, `GOOGLE_KEYFILE_JSON`,
  `GOOGLE_APPLICATION_CREDENTIALS` environment variables.

* `impersonated_user_email` - (Optional) Service accounts cannot be granted
  access to the Admin API SDK, therefore the service account needs to
  impersonate one of the users to access the Admin SDK. May be set via the
  `IMPERSONATED_USER_EMAIL` environment variable. No default impersonated user
  email is set.

* `oauth_scopes` - (Optional) When granting the service account oauth scopes,
  you need to let this provider know it can use them. For a list of oauth scopes
  see this [link](https://developers.google.com/admin-sdk/directory/v1/guides/authorizing).
  No default oauth scopes are set.

* `customer_id` - (Optional) By default we use my_customer as customer ID, which
  means the API will use the G Suite customer ID associated with the
  impersonating account. Override this setting when you know what you are doing.
  Default value of `my_customer`.

* `timeout_minutes` - (Optional) G Suite API's are eventually consistent. This
  means that we sometimes need to wait before resources become available. See
  [implementing exponential backoff](https://developers.google.com/admin-sdk/directory/v1/limits#backoff)
  for more information on why this value is `1 minute` by default. You can
  increase this value if you persistently run into backoffs and timeouts.

* `update_existing` - (Optional) Many terraform providers are not authoritative
  by default and do not allow the provider to be set as such. By setting this to
  `true` (default `false`) you tell the provider it is okay to overwrite
  existing values (import on create).

## Example Usage

```hcl
provider "gsuite" {
  # It is strongly recommended to configure this provider through the
  # environment variables described above for "credentials" and
  # "impersonated_user_email", so that each user can have separate credentials
  # set in the environment.
  oauth_scopes = [
    "https://www.googleapis.com/auth/admin.directory.group",
    "https://www.googleapis.com/auth/apps.groups.settings",
    "https://www.googleapis.com/auth/admin.directory.user",
    "https://www.googleapis.com/auth/admin.directory.userschema",
  ]
  # Oauth scopes do not need to be set when using a personal admin account.
}

resource "gsuite_group" "example" {
  email       = "example@domain.ext"
  name        = "example"
  description = "Example group"
}

resource "gsuite_group_member" "owner" {
  group = gsuite_group.example.email
  email = "owner@domain.ext"
  role  = "OWNER"
}

resource "gsuite_group_settings" "example" {
  email = gsuite_group.example.email

  allow_external_members     = true
  allow_google_communication = false
  show_in_group_directory    = "true"
  who_can_discover_group     = "ALL_IN_DOMAIN_CAN_DISCOVER"
}
```
