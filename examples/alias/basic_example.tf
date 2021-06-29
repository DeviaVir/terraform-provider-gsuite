resource "gsuite_user" "test" {
  name = {
    family_name = "TestAcc_replaceWithUuid"
    given_name  = "Test"
  }

  primary_email = "test_replaceWithUuid@domain.ext"

  lifecycle {
    ignore = [aliases]
  }
}

resource "gsuite_user_alias" "test" {
  user_id = gsuite_user.test.primary_email
  alias   = "test-alias-replaceWithUuid@domain.ext"
}
