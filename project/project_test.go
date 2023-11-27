package project

import (
	"bytes"
	OTBasic "cryptographic-computing/project/OTBasic"
	OTExt "cryptographic-computing/project/OTExtension"
	"cryptographic-computing/project/elgamal"
	utils "cryptographic-computing/project/utils"
	"math"
	"testing"
)

func TestOTBasicProtocol(t *testing.T) {
	l := 1
	m := int(math.Pow(2, float64(8)))

	// create cryptoalgorithm, messages and selection bits for algorithms.
	elGamal := elgamal.ElGamal{}
	elGamal.Init()
	selectionBits := utils.RandomSelectionBits(m)
	var messages []*utils.MessagePair
	for i := 0; i < m; i++ {
		msg := utils.MessagePair{
			Message0: utils.RandomBytes(l),
			Message1: utils.RandomBytes(l),
		}
		messages = append(messages, &msg)
	}

	plaintext := OTBasic.OTBasicProtocol(l, m, selectionBits, messages, elGamal)

	// Check if the plaintext is correct
	for i := 0; i < m; i++ {
		if selectionBits[i] == 0 {
			if !bytes.Equal(plaintext[i], messages[i].Message0) {
				t.Errorf("Plaintext is not correct")
			}
		} else {
			if !bytes.Equal(plaintext[i], messages[i].Message1) {
				t.Errorf("Plaintext is not correct")
			}
		}
	}

}

func TestOTExtensionProtocol(t *testing.T) {
	k := 128
	l := 1

	for iters := 1; iters < 12; iters++ {
		m := int(math.Pow(2, float64(iters)))

		// create cryptoalgorithm, messages and selection bits for algorithms.
		elGamal := elgamal.ElGamal{}
		elGamal.Init()
		selectionBits := utils.RandomSelectionBits(m)
		var messages []*utils.MessagePair
		for i := 0; i < m; i++ {
			msg := utils.MessagePair{
				Message0: utils.RandomBytes(l),
				Message1: utils.RandomBytes(l),
			}
			messages = append(messages, &msg)
		}

		plaintext := OTExt.OTExtensionProtocol(k, l, m, selectionBits, messages, elGamal)

		// Check if the plaintext is correct
		for i := 0; i < m; i++ {
			if selectionBits[i] == 0 {
				if !bytes.Equal(plaintext[i], messages[i].Message0) {
					t.Errorf("Plaintext is not correct")
				}
			} else {
				if !bytes.Equal(plaintext[i], messages[i].Message1) {
					t.Errorf("Plaintext is not correct")
				}
			}
		}
	}

}

func TestOTExtensionProtocolTranspose(t *testing.T) {
	k := 128
	l := 1

	for iters := 1; iters < 12; iters++ {
		m := int(math.Pow(2, float64(iters)))

		// create cryptoalgorithm, messages and selection bits for algorithms.
		elGamal := elgamal.ElGamal{}
		elGamal.Init()
		selectionBits := utils.RandomSelectionBits(m)
		var messages []*utils.MessagePair
		for i := 0; i < m; i++ {
			msg := utils.MessagePair{
				Message0: utils.RandomBytes(l),
				Message1: utils.RandomBytes(l),
			}
			messages = append(messages, &msg)
		}

		plaintext := OTExt.OTExtensionProtocolTranspose(k, l, m, selectionBits, messages, elGamal)

		// Check if the plaintext is correct
		for i := 0; i < m; i++ {
			if selectionBits[i] == 0 {
				if !bytes.Equal(plaintext[i], messages[i].Message0) {
					t.Errorf("Plaintext is not correct")
				}
			} else {
				if !bytes.Equal(plaintext[i], messages[i].Message1) {
					t.Errorf("Plaintext is not correct")
				}
			}
		}
	}

}

func TestOTExtensionProtocolEklundh(t *testing.T) {
	k := 128
	l := 1

	for iters := 1; iters < 12; iters++ {
		m := int(math.Pow(2, float64(iters)))

		// create cryptoalgorithm, messages and selection bits for algorithms.
		elGamal := elgamal.ElGamal{}
		elGamal.Init()
		selectionBits := utils.RandomSelectionBits(m)
		var messages []*utils.MessagePair
		for i := 0; i < m; i++ {
			msg := utils.MessagePair{
				Message0: utils.RandomBytes(l),
				Message1: utils.RandomBytes(l),
			}
			messages = append(messages, &msg)
		}

		plaintext := OTExt.OTExtensionProtocolEklundh(k, l, m, selectionBits, messages, elGamal)

		// Check if the plaintext is correct
		for i := 0; i < m; i++ {
			if selectionBits[i] == 0 {
				if !bytes.Equal(plaintext[i], messages[i].Message0) {
					t.Errorf("Plaintext is not correct")
				}
			} else {
				if !bytes.Equal(plaintext[i], messages[i].Message1) {
					t.Errorf("Plaintext is not correct")
				}
			}
		}
	}

}
