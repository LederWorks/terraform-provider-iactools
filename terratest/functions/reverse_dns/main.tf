# Copyright (c) LederWorks
# SPDX-FileCopyrightText: The terraform-provider-iactools Authors
# SPDX-License-Identifier: MPL-2.0

terraform {
  required_providers {
    iactools = {
      source  = "lederworks/iactools"
    }
  }
}

provider "iactools" {}
