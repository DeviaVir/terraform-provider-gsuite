---
layout: "gsuite"
page_title: "G Suite: gsuite_user_attributes"
sidebar_current: "docs-gsuite-resource-user-attributes"
description: |-
  Manage a G Suite's user's attributes.
---

# gsuite\_user\_attributes

Provides a resource to create and manage a User's attributes, currently limited
to the Custom Schema.

**Note:** requires the `https://www.googleapis.com/auth/admin.directory.userschema`
oauth scope.

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

resource "gsuite_user" "user" {
  primary_email = "flast@domain.ext"

  name {
    given_name  = "First"
    family_name = "Last"
  }
}

resource "gsuite_user_attributes" "user_attributes" {
  primary_email = gsuite_user.user.primary_email
  custom_schema {
    name  = gsuite_user_schema.details.schema_name
    value = data.gsuite_user_attributes.details.json
  }
}
```

## Argument Reference

The following arguments are supported:

* `primary_email` - (Required; Forces new resource) Email address of the G Suite
  user.

* `custom_schema` - (Required) The `gsuite_user_schema` custom schema for this
  user.

## Import

G Suite User Attributes can be imported using `group-email`, e.g.:

```
terraform import gsuite_user_attributes.user_attributes "flast@domain.ext"
```
