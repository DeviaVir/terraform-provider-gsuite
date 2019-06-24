// Make sure to authorize the following scope:
// https://www.googleapis.com/auth/apps.groups.settings
// oauth_scopes and possibly your service account scopes need this to function.

resource "gsuite_group" "devteam" {
  email       = "devteam2@sillevis.net"
  name        = "devteam2@sillevis.net"
  description = "Developer team2"
}

resource "gsuite_group_settings" "devteam" {
  email = "${gsuite_group.devteam.email}"

  allow_external_members  = true
  show_in_group_directory = true
  who_can_discover_group  = "ALL_IN_DOMAIN_CAN_DISCOVER"
}
