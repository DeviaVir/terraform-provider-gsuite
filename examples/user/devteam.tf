resource "gsuite_user" "developer" {
  # advise to set this field to true on creation, then false afterwards
  change_password_next_login = true

  name {
    family_name = "Sillevis"
    given_name  = "Chase"
  }

  # on creation this field is required, later during updates it is ignored;
  # it is expected that the user and Google will handle passwords from there on
  # out
  password = "testtest123!"

  primary_email = "developer@sillevis.net"

  ssh_public_keys {
    key                  = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDUYJKI2gGdZr5Brd1IaT8OQSSt81mBBXQnAfmmjw5hOK9VaJ1MmDB5qY7V1nuXftmLBLvaA7L6k21FDJeWxwD8vKuYwbuJyh1DKB6PMXAQxnX7uLSSi9a/ZOzh3gIHXdil0fSWFpFBmImznqbzaEb7nya+tnK4RONoEjJcRe8Tl+8hET/29XBd3oxlfwwjQA9A84iKhAMLdJIQ28z2GA/2mRJ8RkHLrkQL8kMCj4GJYxy3PR9JU0aFAtWh2mXGfOzaBTh/IhpMW53d8puxihBbIN87MoGngYLt4eBEdE0SiHb0Zdqp5ZDCkwNmAKiWOOrDQxtWvUOThHV5eLMMObqA06XFiwNlojl9ZTH0Y2w/LZmvgb98T/1lBY6mb1iRERGKqYNBeSNwh1Afvu1miDau2f5AYqcf7yxvuD8d0O4cb1xfl7WJwWPJraYaN1X+WmCGTIA+Vve+Kp9TaGoE5n5EGz2a7RNzWj0L0hkf8923iEEtTrsfWewnTnq7XzFoaW53xjWcN7jQplisjWr6AWYApyinw0qGD3dzKgPLyOOcdC3YLhYFpGJcMbegrNdmhbxqIXCB3vBpEFV6o4GqdEy2OVFOM6kSydEQUsMHl5WU8l4gYW28ekZZtbrE52v1dMNzKwfrpVPpUfwn4jbeaqYoIWEwFNVnvbJaFu1vjfrshw== chase"
    expiration_time_usec = "1549735670773"
  }

  #
  # WARN: on-create the posix account is ignored!
  # trigger another apply when the account has been logged into and verified.
  # posix accounts are not trivial to set up automatically, triggering this
  # too many times may lead to all terraform created accounts to automatically
  # be suspended for some time. Take care!
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
}
