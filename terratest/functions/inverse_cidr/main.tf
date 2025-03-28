# Copyright (c) LederWorks
# SPDX-FileCopyrightText: The terraform-provider-iactools Authors
# SPDX-License-Identifier: MPL-2.0

terraform {
  required_providers {
    iactools = {
      source  = "localhost/lederworks/iactools"
      version = "9999.99.99"
    }
  }
}

provider "iactools" {}
