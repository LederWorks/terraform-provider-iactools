# Copyright (c) LederWorks
# SPDX-FileCopyrightText: The terraform-provider-iactools Authors
# SPDX-License-Identifier: MPL-2.0

output "inverse_cidrs" {
  value = provider::iactools::inverse_cidrs(var.parent_cidr, var.child_cidr)
}
