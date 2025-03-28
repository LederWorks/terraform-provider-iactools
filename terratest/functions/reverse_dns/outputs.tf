# Copyright (c) LederWorks
# SPDX-FileCopyrightText: The terraform-provider-iactools Authors
# SPDX-License-Identifier: MPL-2.0

output "reverse_dns" {
  value = provider::iactools::reverse_dns(var.ip_address)
}
