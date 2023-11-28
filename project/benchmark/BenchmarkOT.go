// BenchmarkOT.go
package OTExtension

import (
	OTExt "cryptographic-computing/project/OTExtension"
	"cryptographic-computing/project/elgamal"
	"cryptographic-computing/project/utils"
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"
)

func TestMakeDataFixL(iterations int) {
	csvFile, err := os.Create("./fixed_m_data.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvwriter := csv.NewWriter(csvFile)
	_ = csvwriter.Write([]string{"m_size", "time_OT_Basic", "time_OT_Extension", "time_OT_Extension_Transpose", "time_OT_Extension_Eklundh", "time_OT_Extension_Eklundh_Multithreaded"})

	k := 128
	l := 1

	for i := 2; i < iterations; i++ {

		m := int(math.Pow(2, float64(i)))

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

		// time_start := time.Now()
		// OTBasic.OTBasicProtocol(l, m, selectionBits, messages, elGamal)
		// time_end := time.Since(time_start).Seconds()
		// time_OT_Basic := fmt.Sprintf("%.2f", time_end)
		time_OT_Basic := "0"

		time_start := time.Now()
		OTExt.OTExtensionProtocol(k, l, m, selectionBits, messages, elGamal)
		time_end := time.Since(time_start).Seconds()
		time_OT_Extension := fmt.Sprintf("%.2f", time_end)

		time_start = time.Now()
		OTExt.OTExtensionProtocolTranspose(k, l, m, selectionBits, messages, elGamal)
		time_end = time.Since(time_start).Seconds()
		time_OT_Extension_Transpose := fmt.Sprintf("%.2f", time_end)

		time_start = time.Now()
		OTExt.OTExtensionProtocolEklundh(k, l, m, selectionBits, messages, elGamal, false)
		time_end = time.Since(time_start).Seconds()
		time_OT_Extension_Eklundh := fmt.Sprintf("%.2f", time_end)

		time_start = time.Now()
		OTExt.OTExtensionProtocolEklundh(k, l, m, selectionBits, messages, elGamal, true)
		time_end = time.Since(time_start).Seconds()
		time_OT_Extension_Eklundh_Multithreaded := fmt.Sprintf("%.2f", time_end)

		_ = csvwriter.Write([]string{strconv.Itoa(m), time_OT_Basic, time_OT_Extension, time_OT_Extension_Transpose, time_OT_Extension_Eklundh, time_OT_Extension_Eklundh_Multithreaded})
		csvwriter.Flush()
	}

}

/*
func TestMakeDataFixM(iterations int) {
	csvFile, err := os.Create("./testdata/fixed_m_data.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvwriter := csv.NewWriter(csvFile)
	_ = csvwriter.Write([]string{"m_size", "time_ot_basic", "time_OT_extension"})

	time_OT_basic := 0
	time_OT_extension := 0

	for i := 1; i < iterations; i++ {

		num_of_l := math.Pow(2, float64(i))

		time_start := time.Now()
		// insert OT Basic
		time_end := int(time.Since(time_start))
		time_OT_basic = time_end

		time_start = time.Now()
		// insert OT Extension
		time_end = int(time.Since(time_start))
		time_OT_extension = time_end

		_ = csvwriter.Write([]string{strconv.Itoa(int(num_of_l)), strconv.Itoa(time_OT_basic), strconv.Itoa(time_OT_extension)})
		csvwriter.Flush()

	}

} */
