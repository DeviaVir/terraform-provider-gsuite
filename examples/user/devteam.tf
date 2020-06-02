resource "gsuite_user" "developer" {

  aliases = [
    "chase@sillevis.net"
  ]

  name {
    family_name = "Sillevis"
    given_name  = "Chase"
  }

  # Note the following behaviors regarding passwords:
  #
  #   - When running `terraform import` on a user resource:
  #     - The `password` and `hash_function` fields are ignored.
  #   - When running `terraform apply` with a new user resource in your terraform state:
  #     - If the user does not exist in GSuite the following applies:
  #       - The `password` field should be set or a secured password will be automatically generated.
  #       - The `hash_function` field must be set only if the `password` field contains a hashed value.
  #       - The GSuite account will be configured to require password change on next login.
  #     - If the user exists in GSuite the following applies:
  #       - The `password` and `hash_function` fields will be ignored.
  #   - When running `terraform apply` with an existing user resource:
  #     - Empty `password` and `hash_function` fields will be ignored.
  password = "testtest123!"

  primary_email = "developer@sillevis.net"

  ssh_public_keys {
    key                  = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDUYJKI2gGdZr5Brd1IaT8OQSSt81mBBXQnAfmmjw5hOK9VaJ1MmDB5qY7V1nuXftmLBLvaA7L6k21FDJeWxwD8vKuYwbuJyh1DKB6PMXAQxnX7uLSSi9a/ZOzh3gIHXdil0fSWFpFBmImznqbzaEb7nya+tnK4RONoEjJcRe8Tl+8hET/29XBd3oxlfwwjQA9A84iKhAMLdJIQ28z2GA/2mRJ8RkHLrkQL8kMCj4GJYxy3PR9JU0aFAtWh2mXGfOzaBTh/IhpMW53d8puxihBbIN87MoGngYLt4eBEdE0SiHb0Zdqp5ZDCkwNmAKiWOOrDQxtWvUOThHV5eLMMObqA06XFiwNlojl9ZTH0Y2w/LZmvgb98T/1lBY6mb1iRERGKqYNBeSNwh1Afvu1miDau2f5AYqcf7yxvuD8d0O4cb1xfl7WJwWPJraYaN1X+WmCGTIA+Vve+Kp9TaGoE5n5EGz2a7RNzWj0L0hkf8923iEEtTrsfWewnTnq7XzFoaW53xjWcN7jQplisjWr6AWYApyinw0qGD3dzKgPLyOOcdC3YLhYFpGJcMbegrNdmhbxqIXCB3vBpEFV6o4GqdEy2OVFOM6kSydEQUsMHl5WU8l4gYW28ekZZtbrE52v1dMNzKwfrpVPpUfwn4jbeaqYoIWEwFNVnvbJaFu1vjfrshw== chase"
    expiration_time_usec = "1549735670773"
  }

  #
  # WARN: it is possible on-creation of a new account that the POSIX data is
  # found to not be unique, and a 503 backend error is returned indefinitely.
  # In that case, the account is created, but without the POSIX data. Simply
  # update the POSIX data and terraform apply to update until it works.
  #
  # best to fill out all of these fields or face the consequences
  # might get 503 backend errors if you try to change this too often/fast
  posix_accounts {
    home_directory = "/home/chase"
    primary        = true
    gid            = 1001
    uid            = 1001
    shell          = "/bin/bash"
    system_id      = "uid"
    username       = "chase"
  }

  external_ids {
    type  = "organization"
    value = "1234"
  }

  # If omitted or `true` existing GSuite users defined as Terraform resources will be imported by `terraform apply`.
  update_existing = true
}
