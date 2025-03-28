package provider

import (
	"fmt"
	"log"
	"net"
	"strings"
)

// InverseCIDR determines the address type and calls the appropriate function.
func InverseCIDR(parentCIDR, childCIDR string) ([]string, error) {
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

	// Find the inverse CIDRs leading to the child CIDR
	inverseCIDRs, found := findInverseCIDRs(parentNet, childNet)
	if !found {
		return nil, fmt.Errorf("child CIDR not found within parent CIDR")
	}

	return convertToStringSlice(inverseCIDRs), nil
}

// Helper functions

// splitCIDR splits a CIDR into two smaller CIDRs.
func splitCIDR(ipnet *net.IPNet) ([]*net.IPNet, error) {
	ones, bits := ipnet.Mask.Size()
	if ones >= bits {
		return nil, fmt.Errorf("cannot split CIDR %s: mask size is too large", ipnet.String())
	}

	// Calculate the new subnet mask
	newOnes := ones + 1
	newMask := net.CIDRMask(newOnes, bits)

	// Create two new subnets
	firstIP := ipnet.IP.Mask(newMask)
	secondIP := make(net.IP, len(firstIP))
	copy(secondIP, firstIP)
	secondIP[ones/8] |= 1 << (7 - uint(ones%8))

	firstSubnet := &net.IPNet{IP: firstIP, Mask: newMask}
	secondSubnet := &net.IPNet{IP: secondIP, Mask: newMask}

	return []*net.IPNet{firstSubnet, secondSubnet}, nil
}

// Helper function to get the sibling subnet
func getSiblingSubnet(subnets []*net.IPNet, target *net.IPNet) *net.IPNet {
	if subnets[0].String() == target.String() {
		return subnets[1]
	}
	return subnets[0]
}

// Recursive function to find the path to the child CIDR and the inverse CIDR ranges
func findInverseCIDRs(parentCIDR, childCIDR *net.IPNet) ([]*net.IPNet, bool) {
	subnets, err := splitCIDR(parentCIDR)
	if err != nil {
		log.Printf("Error splitting CIDR %s: %v", parentCIDR, err)
		return nil, false
	}

	for _, subnet := range subnets {
		if subnet.Contains(childCIDR.IP) {
			if subnet.String() == childCIDR.String() {
				// Add the sibling CIDR when the child CIDR is found
				return []*net.IPNet{getSiblingSubnet(subnets, subnet)}, true
			}
			inverseCIDRs, found := findInverseCIDRs(subnet, childCIDR)
			if found {
				return append([]*net.IPNet{getSiblingSubnet(subnets, subnet)}, inverseCIDRs...), true
			}
		}
	}

	return nil, false
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
