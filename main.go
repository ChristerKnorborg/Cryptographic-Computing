package main

import (
	"cryptographic-computing/project/OTBasic"
	// "cryptographic-computing/project/OTExtension"
)

// k int, l int, m int
func main() {
	// k := 5000  // Security parameter
	// l := 1     // Byte length of each message
	// m := 10000 // Number of messages to be sent

	// OTExtension.OTExtensionProtocol(k, l, m)

	l := 4 // Byte length of each message
	m := 2 // Number of messages to be sent

	OTBasic.OTBasicProtocol(l, m)

}
