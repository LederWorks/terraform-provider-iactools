// Copyright (c) LederWorks
// SPDX-FileCopyrightText: The terraform-provider-iactools Authors
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"net"
	"strings"
)

// ReverseDNSIPv4 generates the reverse DNS PTR record for an IPv4 address.
func ReverseDNSIPv4(ipAddress string) string {
	splitted := strings.Split(ipAddress, ".")
	parts := make([]string, 0)
	for i := len(splitted) - 1; i >= 0; i-- {
		parts = append(parts, splitted[i])
	}
	joined := strings.Join(parts, ".")

	return fmt.Sprintf("%v.in-addr.arpa.", joined)
}

// ReverseDNSIPv6 generates the reverse DNS PTR record for an IPv6 address.
func ReverseDNSIPv6(ipAddress net.IP) string {
	expandedAddress := expandIPv6Address(ipAddress)
	splitted := strings.Split(strings.ReplaceAll(expandedAddress, ":", ""), "")
	parts := make([]string, 0)
	for i := len(splitted) - 1; i >= 0; i-- {
		parts = append(parts, splitted[i])
	}
	joined := strings.Join(parts, ".")

	return fmt.Sprintf("%v.ip6.arpa.", joined)
}

func expandIPv6Address(ip net.IP) string {
	b := make([]byte, 0, len(ip))

	for i := 0; i < len(ip); i += 2 {
		if i > 0 {
			b = append(b, ':')
		}
		s := (uint32(ip[i]) << 8) | uint32(ip[i+1])
		bHex := fmt.Sprintf("%04x", s)
		b = append(b, bHex...)
	}
	return string(b)
}
