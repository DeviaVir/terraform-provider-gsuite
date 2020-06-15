---
layout: "gsuite"
page_title: "G Suite: gsuite_user_schema"
sidebar_current: "docs-gsuite-resource-user-schema"
description: |-
  Managing a G Suite User Schema.
---

# gsuite\_user\_schema

Provides a resource to create and manage a G Suite User Schema.

**Note:** If you get an error when applying schema changes such as:
```
googleapi: Error 400: Invalid Input: custom_schema, invalid
```

The most likely cause is you are attempting to update a schema with value(s) it
does not support. For example if you apply data.gsuite_user_attributes.test to
the `additional-details` custom_schema below you'll get the error above as a
result.

## Example Usage

```hcl
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

// A data resource that defines full schema that contains all field types mentioned above. Not very
// interesting but can be used for testing changes to the terraform provider. Note,
// the reason we
resource "gsuite_user_schema" "test" {
  schema_name  = "test-schema"
  display_name = "test-schema"

  field {
    field_type = "STRING"
    field_name = "string"
  }

  field {
    field_type   = "STRING"
    field_name   = "strings"
    multi_valued = true
  }

  field {
    field_type = "BOOL"
    field_name = "bool"
  }

  field {
    field_type   = "BOOL"
    field_name   = "bools"
    multi_valued = true
  }

  field {
    field_type = "INT64"
    field_name = "integer"
  }

  field {
    field_type   = "INT64"
    field_name   = "integers"
    multi_valued = true
  }

  field {
    field_type = "DOUBLE"
    field_name = "double"
  }

  field {
    field_type   = "DOUBLE"
    field_name   = "doubles"
    multi_valued = true
  }

  field {
    field_type = "DATE"
    field_name = "date"
  }

  field {
    field_type   = "DATE"
    field_name   = "dates"
    multi_valued = true
  }

  field {
    field_type = "EMAIL"
    field_name = "email"
  }

  field {
    field_type   = "EMAIL"
    field_name   = "emails"
    multi_valued = true
  }

  field {
    field_type = "PHONE"
    field_name = "phone"
  }

  field {
    field_type   = "PHONE"
    field_name   = "phones"
    multi_valued = true
  }
}

// This data resource is used to fill in the `test-schema` defined above. Useful
// for testing changes to the terraform provider. Note, passing in a value that
// has the wrong type (ex. a non-date to a DATE type filed) results in
// the following error:
//   * gsuite_user.user: Error updating user: {"error":{"errors":[{"domain":"global","reason":"invalid","message":"Invalid Input"}],"code":400,"message":"Invalid Input"}}: googleapi: Error 400: Invalid Input, invalid
// The reason a data source is provided is that this allows the GSuite provider
// to properly handle the construction of the json payload to Google's directory
// API. If you choose not to use this data source you'll need to construct it
// yourself according to the SDK documentation:
//   https://developers.google.com/admin-sdk/directory/v1/guides/manage-schemas#set_fields
data "gsuite_user_attributes" "test" {
  string {
    name  = "string"
    value = "some-string"
  }

  strings {
    name  = "strings"
    value = ["string1", "string2", "string3"]
  }

  bool {
    name  = "bool"
    value = false
  }

  bools {
    name  = "bools"
    value = [true, false, true, true, true]
  }

  integer {
    name  = "integer"
    value = 1001
  }

  integers {
    name  = "integers"
    value = [1001, 1002, 1003]
  }

  double {
    name  = "double"
    value = 3.1415926
  }

  doubles {
    name  = "doubles"
    value = [3.1415926, 1.2e-8]
  }

  date {
    name  = "date"
    value = "1970"
  }

  dates {
    name  = "dates"
    value = ["1970", "1970-01-01"]
  }

  email {
    name  = "email"
    value = "test@domain.ext"
  }

  emails {
    name  = "emails"
    value = ["test1@domain.ext", "test2@domain.ext"]
  }

  phone {
    name  = "phone"
    value = "1-555-555-5555"
  }

  phones {
    name  = "phones"
    value = ["1-555-555-5555", "1-555-555-5556"]
  }
}

resource "gsuite_user" "user" {
  // depends_on is required. Otherwise changes to the schema won't be applied
  // before we update this resource.
  depends_on = [gsuite_user_schema.test, gsuite_user_schema.details]

  primary_email = "flast@domain.ext"

  name {
    given_name  = "First"
    family_name = "Last"
  }

  // Set attributes on the user for the `test-schema` custom schema.
  custom_schema {
    name  = "test-schema"
    value = data.gsuite_user_attributes.test.json
  }
  // Set attributes on the user for the `additional-details` custom schema.
  custom_schema {
    name  = "additional-details"
    value = data.gsuite_user_attributes.details.json
  }
}
```

## Argument Reference

The following arguments are supported:

* `schema_name` - (Required) Name of the user schema.

* `field` - (Required) See the examples above.

* `display_name` - (Optional) Human friendly name for this User Schema.

## Attribute Reference

In addition to the above arguments, the following attributes are exported:

* `schema_id`


## Import

A G Suite User Schema can be imported using `schema_id`, e.g.:

```
terraform import gsuite_user_schema.test "test-schema"
```
