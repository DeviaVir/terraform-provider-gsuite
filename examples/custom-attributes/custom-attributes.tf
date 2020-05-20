// The below define a custom directory user schema. For additional details on
// possible types of fields, their values, etc see Google's reference documentation:
//     https://developers.google.com/admin-sdk/directory/v1/reference/schemas
resource "gsuite_user_schema" "details" {
  schema_name  = "additional-details" // Required
  display_name = "Additional Details" // Optional (default: value of `schema_name`)

  // Defines the schema of a specific field. The `field` element may be
  // repeated multiple times to define more than a single custom field.
  field {
    field_type       = "PHONE"            // Required (valid values: BOOL, DATE, DOUBLE, EMAIL, INT64, PHONE, STRING)
    field_name       = "internal-phone"   // Required
    display_name     = "Internal Phone"   // Optional (default: value of `field_name`)
    indexed          = true               // Optional (default: true)
    read_access_type = "ALL_DOMAIN_USERS" // Optional (default: ADMINS_AND_SELF)
    multi_valued     = false              // Optional (default: false)
  }
}

data "gsuite_user_attributes" "details" {
  phone {
    name  = "internal-phone"
    value = "555-555-5555"
  }
}

resource "gsuite_user" "user" {
  primary_email = "flast@example.com"

  name {
    given_name  = "First"
    family_name = "Last"
  }
}

resource "gsuite_user_attributes" "user_attributes" {
  primary_email = gsuite_user.user.primary_email
  custom_schema =  data.gsuite_user_attributes.details.json
}