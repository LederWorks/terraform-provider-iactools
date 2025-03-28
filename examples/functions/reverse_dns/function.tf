# Copyright (c) LederWorks
# SPDX-FileCopyrightText: The terraform-provider-iactools Authors
# SPDX-License-Identifier: MPL-2.0

output "reverse_dns_ipv4" {
  value = provider::iactools::reverse_dns("1.2.3.4")
}

output "reverse_dns_ipv6" {
  value = provider::iactools::reverse_dns("2001:db8::567:89ab")
}
