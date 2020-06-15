Terraform Provider - G Suite
==================

- Website: https://registry.terraform.io/providers/DeviaVir/gsuite/latest/docs
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Maintainers
-----------

This provider plugin is maintained by Chase Sillevis.

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.11.x
-	[Go](https://golang.org/doc/install) 1.14 (to build the provider plugin)

Installing the Provider
---------------------

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

Building The Provider
---------------------

1. `cd` into `$HOME/.terraform.d/plugins/terraform-provider-gsuite`

1. Run `make vendor` to fetch the go vendor files

1. Make your changes

1. Run `make dev` and in your `terraform` directory, remove the current `.terraform` and re-run `terraform init`

1. Next time you run `terraform plan` it'll use your updated version
