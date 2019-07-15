
resource "gsuite_group" "testing_team" {
  email       = "testing-team@${local.domain}"
  name        = "testing-team@${local.domain}"
  description = "Testing team group"
}

resource "gsuite_group_member_list" "testing_team_members" {
  group_email = "${gsuite_group.testing_team.email}"
  member_list = [
    {
      email = "karthikprabu@${local.domain}"
      role  = "MANAGER"
    },
    {
      email = "bobjoseph@${local.domain}"
      role  = "MEMBER"
    },
    {
      email = "ravi.shastri@${local.domain}"
      role  = "MEMBER"
    }
  ]
}
