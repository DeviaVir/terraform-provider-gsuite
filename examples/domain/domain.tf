# Make sure you have allowed the domain scope before you try to use the
# gsuite_domain resources.
#
# https://www.googleapis.com/auth/admin.directory.domain
# https://github.com/DeviaVir/terraform-provider-gsuite/issues/65#issuecomment-492879619

resource "gsuite_domain" "my_domain" {
  domain_name = "example.com"
}
