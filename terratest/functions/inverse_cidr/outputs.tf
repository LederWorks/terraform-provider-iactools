# Copyright (c) LederWorks
# SPDX-FileCopyrightText: The terraform-provider-iactools Authors
# SPDX-License-Identifier: MPL-2.0

output "inverse_cidr" {
  value = provider::iactools::inverse_cidr(var.parent_cidr, var.child_cidr)
}
