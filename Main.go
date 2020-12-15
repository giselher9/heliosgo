package main

import (
	"fmt"
	"strings"

	"github.com/goburrow/modbus"
)

func main() {
	// Modbus TCP
	handler := modbus.NewTCPClientHandler("helios.fritz.box:502")
	//handler.Timeout = 3 * time.Second
	handler.SlaveId = 180
	//handler.Logger = log.New(os.Stdout, "Debug: ", log.LstdFlags)
	// Connect manually so that multiple requests are handled in one connection session
	err := handler.Connect()
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to Helios device: %v", err))
	}
	defer handler.Close()

	client := modbus.NewClient(handler)

	fmt.Println("Außenluft:", getValueFromHelios(client, außenluft, 8))
	fmt.Println("Zuluft:", getValueFromHelios(client, zuluft, 8))
	fmt.Println("Fortluft:", getValueFromHelios(client, fortluft, 8))
	fmt.Println("Abluft:", getValueFromHelios(client, abluft, 8))
	fmt.Println()
	fmt.Println("Zuluft RPM:", getValueFromHelios(client, zuluftRpm, 8))
	fmt.Println("Abluft RPM:", getValueFromHelios(client, abluftRpm, 8))
	fmt.Println("Lüfterstufe:", getValueFromHelios(client, lüfterstufe, 8))
}

const (
	lüfterstufe = "v00102" // 0x76, 0x30, 0x30, 0x31, 0x30, 0x32, 0x00, 0x00
	außenluft   = "v00104"
	zuluft      = "v00105"
	fortluft    = "v00106"
	abluft      = "v00107"
	zuluftRpm   = "v00348"
	abluftRpm   = "v00349" // 0x76, 0x30, 0x30, 0x33, 0x34, 0x39, 0x00, 0x00
)

func getValueFromHelios(client modbus.Client, heliosVariable string, countBytesToRead byte) string {
	registers := []byte(heliosVariable)
	registers = append(registers, 0x0, 0x0)

	results, err := client.WriteMultipleRegisters(0, 4, registers)
	if err != nil {
		panic(fmt.Sprintf("Unable to send value to Helios device: %v", err))
	}
	results, err = client.ReadHoldingRegisters(0, 6)
	if err != nil {
		panic(fmt.Sprintf("Unable to receive value from Helios device: %v", err))
	}

	return getValueFromHeliosResponse(results)
}

func getValueFromHeliosResponse(heliosResult []byte) (result string) {
	return strings.SplitAfter(modbusResponseToString(heliosResult), "=")[1]
}

func modbusResponseToString(data []byte) (result string) {
	for _, v := range data {
		if v > 0 {
			result += string(v)
		}
	}
	return result
}
