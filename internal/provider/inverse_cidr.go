// Copyright (c) LederWorks
// SPDX-FileCopyrightText: The terraform-provider-iactools Authors
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"net"
	"strings"
)

// InverseCIDR determines the address type and calls the appropriate function (IPv4 or IPv6).
func InverseCIDR(parentCIDR, childCIDR string) ([]string, error) {
	_, parentNet, err := net.ParseCIDR(parentCIDR)
	if err != nil {
		return nil, fmt.Errorf("invalid parent CIDR: %v", err)
	}

	_, childNet, err := net.ParseCIDR(childCIDR)
	if err != nil {
		return nil, fmt.Errorf("invalid child CIDR: %v", err)
	}

	if isIPv4(parentNet.IP) && isIPv4(childNet.IP) {
		return InverseCIDRIPv4(parentCIDR, childCIDR)
	} else {
		return InverseCIDRIPv6(parentCIDR, childCIDR)
	}
}

// InverseCIDRIPv4 calculates the inverse CIDR ranges for a given childCIDR within a parentCIDR for IPv4 addresses.
func InverseCIDRIPv4(parentCIDR, childCIDR string) ([]string, error) {
	_, parentNet, err := net.ParseCIDR(parentCIDR)
	if err != nil {
		return nil, fmt.Errorf("invalid parent CIDR: %v", err)
	}

	_, childNet, err := net.ParseCIDR(childCIDR)
	if err != nil {
		return nil, fmt.Errorf("invalid child CIDR: %v", err)
	}

	if !parentNet.Contains(childNet.IP) {
		return nil, fmt.Errorf("child CIDR %s is not within parent CIDR %s", childCIDR, parentCIDR)
	}

	inverseCIDRs := excludeSubnets(parentNet, []*net.IPNet{childNet})
	return convertToStringSlice(inverseCIDRs), nil
}

// InverseCIDRIPv6 calculates the inverse CIDR ranges for a given childCIDR within a parentCIDR for IPv4 and IPv6 addresses.
func InverseCIDRIPv6(parentCIDR, childCIDR string) ([]string, error) {
	_, childNet, err := net.ParseCIDR(childCIDR)
	if err != nil {
		return nil, fmt.Errorf("invalid child CIDR: %v", err)
	}

	// Check if the childCIDR is an IPv4 address
	if isIPv4(childNet.IP) {
		return InverseCIDRIPv4(parentCIDR, childCIDR)
	}

	_, parentNet, err := net.ParseCIDR(parentCIDR)
	if err != nil {
		return nil, fmt.Errorf("invalid parent CIDR: %v", err)
	}

	if !parentNet.Contains(childNet.IP) {
		return nil, fmt.Errorf("child CIDR %s is not within parent CIDR %s", childCIDR, parentCIDR)
	}

	inverseCIDRs := excludeSubnets(parentNet, []*net.IPNet{childNet})
	return convertToStringSlice(inverseCIDRs), nil
}

// Helper functions

// excludeSubnets calculates the inverse CIDRs excluding the given subnets from the parent network.
func excludeSubnets(parentNet *net.IPNet, subnets []*net.IPNet) []*net.IPNet {
	var result []*net.IPNet
	exclude := make(map[string]struct{})
	for _, subnet := range subnets {
		exclude[subnet.String()] = struct{}{}
	}

	parentCIDRs := []*net.IPNet{parentNet}
	for len(parentCIDRs) > 0 {
		parent := parentCIDRs[0]
		parentCIDRs = parentCIDRs[1:]

		var excluded bool
		for subnet := range exclude {
			_, subnetNet, _ := net.ParseCIDR(subnet)
			if subnetNet.Contains(parent.IP) && subnetNet.String() == parent.String() {
				excluded = true
				break
			}
		}

		if !excluded {
			/* if len(parentCIDRs) == 0 || parent.String() != parentCIDRs[0].String() {
				result = append(result, parent)
			} */
			result = append(result, parent)
		} else {
			sub1, sub2 := splitCIDR(parent)
			parentCIDRs = append(parentCIDRs, sub1, sub2)
		}
	}
	return result
}

// splitCIDR splits a CIDR into two smaller CIDRs.
func splitCIDR(ipnet *net.IPNet) (*net.IPNet, *net.IPNet) {
	prefixLen, _ := ipnet.Mask.Size()
	newPrefixLen := prefixLen + 1

	first := &net.IPNet{
		IP:   ipnet.IP,
		Mask: net.CIDRMask(newPrefixLen, 8*len(ipnet.IP)),
	}
	second := &net.IPNet{
		IP:   make(net.IP, len(ipnet.IP)),
		Mask: net.CIDRMask(newPrefixLen, 8*len(ipnet.IP)),
	}
	copy(second.IP, ipnet.IP)
	// second.IP[newPrefixLen/8-1] |= 1 << (8 - newPrefixLen%8)
	second.IP[len(second.IP)-1] |= 1 << (7 - (newPrefixLen-1)%8)

	return first, second
}

// convertToStringSlice converts a slice of *net.IPNet to a slice of strings.
func convertToStringSlice(ipnets []*net.IPNet) []string {
	var result []string
	for _, ipnet := range ipnets {
		result = append(result, ipnet.String())
	}
	return result
}

// isIPv4 checks if an IP address is IPv4.
func isIPv4(ip net.IP) bool {
	return strings.Contains(ip.String(), ".")
}
