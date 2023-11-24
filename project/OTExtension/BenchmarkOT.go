// BenchmarkOT.go
package OTExtension

import (
	"cryptographic-computing/project/elgamal"
	"encoding/csv"
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
	_ = csvwriter.Write([]string{"m_size", "time_OT_extension"})

	time_OT_extension := 0
	k := 256
	l := 1

	for i := 1; i < iterations; i++ {

		m := int(math.Pow(2, float64(i)))

		// create cryptoalgorithm, messages and selection bits for algorithms.
		elGamal := elgamal.ElGamal{}
		selectionBits := RandomSelectionBits(m)
		var messages []*MessagePair
		for i := 0; i < m; i++ {
			msg := MessagePair{
				Message0: RandomBytes(l),
				Message1: RandomBytes(l),
			}
			messages = append(messages, &msg)
		}

		time_start := time.Now()
		OTExtensionProtocol(k, l, m, selectionBits, messages, elGamal)
		time_end := int(time.Since(time_start))
		time_OT_extension = time_end

		_ = csvwriter.Write([]string{strconv.Itoa(m), strconv.Itoa(time_OT_extension)})
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
