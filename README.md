# Terraform G Suite Provider

This is a terraform provider for managing G Suite (Admin SDK) resources on Google

## Authentication

There are two possible authentication mechanisms for using this provider.
Using a service account, or a personal admin account. The latter requires
user interaction, whereas a service account could be used in an automated
workflow.

See the necessary oauth scopes both for service accounts and users below:
- https://www.googleapis.com/auth/admin.directory.group
- https://www.googleapis.com/auth/admin.directory.user

You could also provide the minimal set of scopes using the
`oauth_scopes` variable in the provider configuration.

```
provider "gsuite" {
  oauth_scopes = [
    "https://www.googleapis.com/auth/admin.directory.group",
    "https://www.googleapis.com/auth/admin.directory.user"
  ]
}
```

**NOTE** If you are creating or modifying schemas and custom user attributes
you will need the following additional scope:

    https://www.googleapis.com/auth/admin.directory.userschema


### Using a service account

Service accounts are great for automated workflows.

Only users with access to the Admin APIs can access the Admin SDK Directory API,
therefore the service account needs to impersonate one of those users
to access the Admin SDK Directory API.

Follow the instruction at
https://developers.google.com/admin-sdk/directory/v1/guides/delegation.

Add `credentials` and `impersonated_user_email` when initializing the provider.
```
provider "gsuite" {
  credentials = "/full/path/service-account.json"
  impersonated_user_email = "admin@xxx.com"
}
```

Credentials can also be provided via the following environment variables:
- GOOGLE_CREDENTIALS
- GOOGLE_CLOUD_KEYFILE_JSON
- GCLOUD_KEYFILE_JSON
- IMPERSONATED_USER_EMAIL

### Using a personal administrator account

In order to use the Admin SDK with a project, we will first need to create
credentials for that project, you can do so here:

https://console.cloud.google.com/apis/credentials?project=[project_ID]

Please make sure to create an OAuth 2.0 client, and download the file to your
local directory.

You can now use that credential to authenticate:

```
$ gcloud auth application-default login \
  --client-id-file=client_id.json \
  --scopes \
https://www.googleapis.com/auth/admin.directory.group,\
https://www.googleapis.com/auth/admin.directory.user,
```

Now that you have a credential that is allowed to the Admin SDK, you can use the
G Suite provider.

## Configuration

Your G Suite Customer ID is required for some Admins SDK Directory APIs,
therefore add also `customer_id` when initializing the provider.

```
provider "gsuite" {
  customer_id = "xxxxxxxx"
}
```

## Installation

1. Download the latest compiled binary from [GitHub releases](https://github.com/DeviaVir/terraform-provider-gsuite/releases).

1. Unzip/untar the archive.

1. Move it into `$HOME/.terraform.d/plugins`:

    ```sh
    $ mkdir -p $HOME/.terraform.d/plugins
    $ mv terraform-provider-gsuite $HOME/.terraform.d/plugins/terraform-provider-gsuite
    ```

1. Create your Terraform configurations as normal, and run `terraform init`:

    ```sh
    $ terraform init
    ```

    This will find the plugin locally.

## Development

1. `cd` into `$HOME/.terraform.d/plugins/terraform-provider-gsuite`

1. Run `dep ensure` to fetch the go vendor files

1. Make your changes

1. Run `make dev` and in your `terraform` directory, remove the current `.terraform` and re-run `terraform init`

1. Next time you run `terraform plan` it'll use your updated version

### Relevant Google Admin SDK Documentation
#### General
* http://google.golang.org/api/admin/directory/v1
* https://developers.google.com/admin-sdk/directory/v1/reference/

#### Schema Types
* https://developers.google.com/admin-sdk/directory/v1/reference/users
* https://developers.google.com/admin-sdk/directory/v1/reference/groups
* https://developers.google.com/admin-sdk/directory/v1/reference/schemas

When using a service account, make sure to add:
`https://www.googleapis.com/auth/admin.directory.userschema`
to the `oauth_scopes` list, otherwise you will be missing permissions to manage
user schemas.

## Notes

- Asking too many permissions right now, but rather start out with too much and tone down later on
- Quite limited, as it is a huge API, I have only added the parts I plan on using
  - Open for PR's to extend functionality
- Documentation is still to be written, you can refer to the `examples` directory for now
