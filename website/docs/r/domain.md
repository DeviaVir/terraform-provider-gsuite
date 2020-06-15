---
layout: "gsuite"
page_title: "G Suite: gsuite_domain"
sidebar_current: "docs-gsuite-resource-domain"
description: |-
  Managing domains in G Suite
---

# gsuite\_domain

Provides a resource to create and manage domains in a G Suite account.

**Note:** Requires the `https://www.googleapis.com/auth/admin.directory.domain`
oauth scope.

## Example Usage

```hcl
resource "gsuite_domain" "example" {
    domain_name = "example"
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required; Forces new resource) Name of the domain.

## Attribute Reference

In addition to the above arguments, the following attributes are exported:

* `creation_time`

* `etag`

## Import

Domains can currently not be imported.
