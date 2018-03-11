resource "gsuite_group" "testing_team" {
  email       = "testing-team@xxx.com"
  name        = "testing-teamg@xxx.com"
  description = "Testing team group"
}

resource "gsuite_group_members" "testing_team_members" {
  group = "${gsuite_group.testing_team.id}"

  member {
    email = "a@xxx.com",
    role  = "MEMBER"
  }

  member {
    email = "b@xxx.com",
    role  = "OWNER"
  }
}
