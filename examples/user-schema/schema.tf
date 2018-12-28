// The below define a custom directory user schema. For additional details on
// possible types of fields, their values, etc see Google's reference documentation:
//     https://developers.google.com/admin-sdk/directory/v1/reference/schemas
resource "gsuite_user_schema" "details" {
  schema_name  = "additional-details" // Required
  display_name = "Additional Details" // Optional (default: value of `schema_name`)

  // Defines the schema of a specific field. The `field` element may be
  // repeated multiple times to define more than a single custom field.
  field {
    field_type       = "PHONE"            // Required (valid values: BOOL, DATA, DOUBLE, EMAIL, INT64, PHONE, STRING)
    field_name       = "internal-phone"   // Required
    display_name     = "Internal Phone"   // Optional (default: value of `field_name`)
    indexed          = true               // Optional (default: true)
    read_access_type = "ALL_DOMAIN_USERS" // Optional (default: ADMINS_AND_SELF)
    multi_valued     = false              // Optional (default: false)
  }
}
