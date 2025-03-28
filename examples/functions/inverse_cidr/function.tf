output "inverse_cidr" {
  value = provider::iactools::inverse_cidr("192.168.0.0/16", "192.168.1.0/24")
}