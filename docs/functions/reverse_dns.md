---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "reverse_dns function - iactools"
subcategory: ""
description: |-
  Calculate the reverse DNS name of an IP address
---

# function: reverse_dns

Accepts both IPv4 and IPv6 addresses and outputs their reverse DNS entry.

## Example Usage

```terraform
# Copyright (c) LederWorks
# SPDX-FileCopyrightText: The terraform-provider-iactools Authors
# SPDX-License-Identifier: MPL-2.0

output "reverse_dns_ipv4" {
  value = provider::iactools::reverse_dns("1.2.3.4")
}

output "reverse_dns_ipv6" {
  value = provider::iactools::reverse_dns("2001:db8::567:89ab")
}
```

## Signature

<!-- signature generated by tfplugindocs -->
```text
reverse_dns(ip_address string) string
```

## Arguments

<!-- arguments generated by tfplugindocs -->
1. `ip_address` (String) The IPv4 or IPv6 address itself

