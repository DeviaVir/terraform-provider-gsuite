# Terraform GSuite Provider

This is a terraform provider for managing GSuite (Admin SDK) resources on Google

## Setup

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
  https://www.googleapis.com/auth/admin.directory.customer,\
  https://www.googleapis.com/auth/admin.directory.group,\
  https://www.googleapis.com/auth/admin.directory.orgunit,\
  https://www.googleapis.com/auth/admin.directory.user,\
  https://www.googleapis.com/auth/admin.directory.userschema,
```

Now that you have a credential that is allowed to the Admin SDK, you can use the
GSuite provider.

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

Some useful resources:

* http://google.golang.org/api/admin/directory/v1
* https://developers.google.com/admin-sdk/directory/v1/reference/

## Notes

- Asking too many permissions right now, but rather start out with too much and tone down later on
- Quite limited, as it is a huge API, I have only added the parts I plan on using
  - Open for PR's to extend functionality
- Documentation is still to be written, you can refer to the `examples` directory for now
