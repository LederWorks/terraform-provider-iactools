# Copyright (c) LederWorks
# SPDX-FileCopyrightText: The terraform-provider-iactools Authors
# SPDX-License-Identifier: MPL-2.0

output "inverse_cidr_ipv4" {
  value = provider::iactools::inverse_cidr("192.168.0.0/16", "192.168.1.0/24")
}

output "inverse_cidr_ipv6" {
  value = provider::iactools::inverse_cidr("2001:db8::/32", "2001:db8:1::/48")
}