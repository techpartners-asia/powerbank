package main

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	powerbankSdk "github.com/techpartners-asia/powerbank/api"
	"github.com/techpartners-asia/powerbank/constants"
	powerbankModels "github.com/techpartners-asia/powerbank/models"
	powerbankUtils "github.com/techpartners-asia/powerbank/utils"
)

func testParseCheckResponse() {
	fmt.Println("=== Testing ParseCheckResponse ===")
	// Example data from the documentation
	// A8 01 25 10 10 00 00 00 6E 6E 01 01 00 F1 04 05 8E 17 08 5F 18 00 00 CA 00 02 01 00 E1 04 05 8E 15 F0 63 18 00 00 CA 00 03 00 00 00 00 00 00 00 00 FF FF 00 00 00 80 04 01 00 EB 04 05 8D 99 F7 5F 18 00 00 CA 00 05 01 00 EA 04 05 8D FC 12 5F 18 00 00 CA 00 06 01 00 B6 04 05 8D 93 E7 1A 1C 32 14 CA 00 20 00 00 00 6E 6E 07 01 00 F1 04 05 8E 17 08 5F 18 00 00 CA 00 08 01 00 E1 04 05 8E 15 F0 63 18 00 00 CA 00 09 00 00 00 00 00 00 00 00 FF FF 00 00 00 80 0A 01 00 EB 04 05 8D 99 F7 5F 18 00 00 CA 00 0B 01 00 EA 04 05 8D FC 12 5F 18 00 00 CA 00 0C 01 00 B6 04 05 8D 93 E7 1A 1C 32 14 CA 00 30 00 00 00 6E 6E 0D 01 00 F1 04 05 8E 17 08 5F 18 00 00 CA 00 0E 01 00 E1 04 05 8E 15 F0 63 18 00 00 CA 00 0F 00 00 00 00 00 00 00 00 FF FF 00 00 00 80 10 01 00 EB 04 05 8D 99 F7 5F 18 00 00 CA 00 11 01 00 EA 04 05 8D FC 12 5F 18 00 00 CA 00 12 01 00 B6 04 05 8D 93 E7 1A 1C 32 14 CA 00 88
	testData := []byte{
		0xA8, 0x01, 0x25, 0x10, // Header: A8, Length: 01 25 (293), Cmd: 10
		// Control Board 1: 10 00 00 00 6E 6E
		0x10, 0x00, 0x00, 0x00, 0x6E, 0x6E,
		// Holes 1-4 (15 bytes each)
		0x01, 0x01, 0x00, 0xF1, 0x04, 0x05, 0x8E, 0x17, 0x08, 0x5F, 0x18, 0x00, 0x00, 0xCA, 0x00,
		0x02, 0x01, 0x00, 0xE1, 0x04, 0x05, 0x8E, 0x15, 0xF0, 0x63, 0x18, 0x00, 0x00, 0xCA, 0x00,
		0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0x80,
		0x04, 0x01, 0x00, 0xEB, 0x04, 0x05, 0x8D, 0x99, 0xF7, 0x5F, 0x18, 0x00, 0x00, 0xCA, 0x00,
		// Control Board 2: 20 00 00 00 6E 6E
		0x20, 0x00, 0x00, 0x00, 0x6E, 0x6E,
		// Holes 5-8 (15 bytes each)
		0x05, 0x01, 0x00, 0xEA, 0x04, 0x05, 0x8D, 0xFC, 0x12, 0x5F, 0x18, 0x00, 0x00, 0xCA, 0x00,
		0x06, 0x01, 0x00, 0xB6, 0x04, 0x05, 0x8D, 0x93, 0xE7, 0x1A, 0x1C, 0x32, 0x14, 0xCA, 0x00,
		0x07, 0x01, 0x00, 0xF1, 0x04, 0x05, 0x8E, 0x17, 0x08, 0x5F, 0x18, 0x00, 0x00, 0xCA, 0x00,
		0x08, 0x01, 0x00, 0xE1, 0x04, 0x05, 0x8E, 0x15, 0xF0, 0x63, 0x18, 0x00, 0x00, 0xCA, 0x00,
		// Control Board 3: 30 00 00 00 6E 6E
		0x30, 0x00, 0x00, 0x00, 0x6E, 0x6E,
		// Holes 9-12 (15 bytes each)
		0x09, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0x80,
		0x0A, 0x01, 0x00, 0xEB, 0x04, 0x05, 0x8D, 0x99, 0xF7, 0x5F, 0x18, 0x00, 0x00, 0xCA, 0x00,
		0x0B, 0x01, 0x00, 0xEA, 0x04, 0x05, 0x8D, 0xFC, 0x12, 0x5F, 0x18, 0x00, 0x00, 0xCA, 0x00,
		0x0C, 0x01, 0x00, 0xB6, 0x04, 0x05, 0x8D, 0x93, 0xE7, 0x1A, 0x1C, 0x32, 0x14, 0xCA, 0x00,
		// Verification byte
		0x88,
	}

	fmt.Printf("Test data length: %d bytes\n", len(testData))
	fmt.Printf("Test data: % X\n", testData)

	response, err := powerbankUtils.ParseCheckResponse(testData)
	if err != nil {
		fmt.Printf("Error parsing check response: %v\n", err)
		return
	}

	fmt.Printf("Parsed response:\n")
	fmt.Printf("Head: 0x%02X\n", response.Head)
	fmt.Printf("Length: %d\n", response.Length)
	fmt.Printf("Cmd: 0x%02X\n", response.Cmd)
	fmt.Printf("Verify: 0x%02X\n", response.Verify)
	fmt.Printf("Number of control boards: %d\n", len(response.ControlBoards))

	for i, cb := range response.ControlBoards {
		fmt.Printf("\nControl Board %d:\n", i+1)
		fmt.Printf("  ControlIndex: 0x%02X (%d)\n", cb.ControlIndex, cb.ControlIndex)
		fmt.Printf("  Undefined1: %d\n", cb.Undefined1)
		fmt.Printf("  Undefined2: %d\n", cb.Undefined2)
		fmt.Printf("  Temperature: %d\n", cb.Temperature)
		fmt.Printf("  SoftVersion: %d\n", cb.SoftVersion)
		fmt.Printf("  HardVersion: %d\n", cb.HardVersion)
		fmt.Printf("  Number of holes: %d\n", len(cb.Holes))

		for j, hole := range cb.Holes {
			fmt.Printf("    Hole %d:\n", j+1)
			fmt.Printf("      HoleIndex: %d\n", hole.HoleIndex)
			fmt.Printf("      State: %d (%s)\n", hole.State, hole.GetStateDescription())
			fmt.Printf("      PowerbankCurr: %.1fA\n", hole.PowerbankCurr)
			fmt.Printf("      PowerbankVolt: %.1fV\n", hole.PowerbankVolt)
			fmt.Printf("      Area: %d\n", hole.Area)
			fmt.Printf("      PowerbankSN: %s\n", hole.PowerbankSN)
			fmt.Printf("      SOC: %d%%\n", hole.SOC)
			fmt.Printf("      Temperature: %d°C\n", hole.Temperature)
			fmt.Printf("      ChargeVolt: %.1fV\n", hole.ChargeVolt)
			fmt.Printf("      ChargeCurr: %.1fA\n", hole.ChargeCurr)
			fmt.Printf("      SoftVersion: %d\n", hole.SoftVersion)
			fmt.Printf("      Sensor: 0x%02X\n", hole.Sensor)
		}
	}
}

func testParsePowerBankUploadResponse() {
	fmt.Println("\n=== Testing ParsePowerBankUploadResponse ===")
	// Same test data as ParseCheckResponse since they parse the same format
	testData := []byte{
		0xA8, 0x01, 0x25, 0x10, // Header: A8, Length: 01 25 (293), Cmd: 10
		// Control Board 1: 10 00 00 00 6E 6E
		0x10, 0x00, 0x00, 0x00, 0x6E, 0x6E,
		// Holes 1-4 (15 bytes each)
		0x01, 0x01, 0x00, 0xF1, 0x04, 0x05, 0x8E, 0x17, 0x08, 0x5F, 0x18, 0x00, 0x00, 0xCA, 0x00,
		0x02, 0x01, 0x00, 0xE1, 0x04, 0x05, 0x8E, 0x15, 0xF0, 0x63, 0x18, 0x00, 0x00, 0xCA, 0x00,
		0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0x80,
		0x04, 0x01, 0x00, 0xEB, 0x04, 0x05, 0x8D, 0x99, 0xF7, 0x5F, 0x18, 0x00, 0x00, 0xCA, 0x00,
		// Control Board 2: 20 00 00 00 6E 6E
		0x20, 0x00, 0x00, 0x00, 0x6E, 0x6E,
		// Holes 5-8 (15 bytes each)
		0x05, 0x01, 0x00, 0xEA, 0x04, 0x05, 0x8D, 0xFC, 0x12, 0x5F, 0x18, 0x00, 0x00, 0xCA, 0x00,
		0x06, 0x01, 0x00, 0xB6, 0x04, 0x05, 0x8D, 0x93, 0xE7, 0x1A, 0x1C, 0x32, 0x14, 0xCA, 0x00,
		0x07, 0x01, 0x00, 0xF1, 0x04, 0x05, 0x8E, 0x17, 0x08, 0x5F, 0x18, 0x00, 0x00, 0xCA, 0x00,
		0x08, 0x01, 0x00, 0xE1, 0x04, 0x05, 0x8E, 0x15, 0xF0, 0x63, 0x18, 0x00, 0x00, 0xCA, 0x00,
		// Control Board 3: 30 00 00 00 6E 6E
		0x30, 0x00, 0x00, 0x00, 0x6E, 0x6E,
		// Holes 9-12 (15 bytes each)
		0x09, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0x80,
		0x0A, 0x01, 0x00, 0xEB, 0x04, 0x05, 0x8D, 0x99, 0xF7, 0x5F, 0x18, 0x00, 0x00, 0xCA, 0x00,
		0x0B, 0x01, 0x00, 0xEA, 0x04, 0x05, 0x8D, 0xFC, 0x12, 0x5F, 0x18, 0x00, 0x00, 0xCA, 0x00,
		0x0C, 0x01, 0x00, 0xB6, 0x04, 0x05, 0x8D, 0x93, 0xE7, 0x1A, 0x1C, 0x32, 0x14, 0xCA, 0x00,
		// Verification byte
		0x88,
	}

	fmt.Printf("Test data length: %d bytes\n", len(testData))
	fmt.Printf("Test data: % X\n", testData)

	response, err := powerbankUtils.ParsePowerBankUploadResponse(testData)
	if err != nil {
		fmt.Printf("Error parsing upload response: %v\n", err)
		return
	}

	fmt.Printf("Parsed upload response:\n")
	fmt.Printf("Head: 0x%02X\n", response.Head)
	fmt.Printf("Length: %d\n", response.Length)
	fmt.Printf("Cmd: 0x%02X\n", response.Cmd)
	fmt.Printf("Verify: 0x%02X\n", response.Verify)
	fmt.Printf("Number of control boards: %d\n", len(response.ControlBoards))

	for i, cb := range response.ControlBoards {
		fmt.Printf("\nControl Board %d:\n", i+1)
		fmt.Printf("  ControlIndex: 0x%02X (%d)\n", cb.ControlIndex, cb.ControlIndex)
		fmt.Printf("  Undefined1: %d\n", cb.Undefined1)
		fmt.Printf("  Undefined2: %d\n", cb.Undefined2)
		fmt.Printf("  Temperature: %d\n", cb.Temperature)
		fmt.Printf("  SoftVersion: %d\n", cb.SoftVersion)
		fmt.Printf("  HardVersion: %d\n", cb.HardVersion)
		fmt.Printf("  Number of holes: %d\n", len(cb.Holes))

		for j, hole := range cb.Holes {
			fmt.Printf("    Hole %d:\n", j+1)
			fmt.Printf("      HoleIndex: %d\n", hole.HoleIndex)
			fmt.Printf("      State: %d (%s)\n", hole.State, hole.GetStateDescription())
			fmt.Printf("      PowerbankCurr: %.1fA\n", hole.PowerbankCurr)
			fmt.Printf("      PowerbankVolt: %.1fV\n", hole.PowerbankVolt)
			fmt.Printf("      Area: %d\n", hole.Area)
			fmt.Printf("      PowerbankSN: %s\n", hole.PowerbankSN)
			fmt.Printf("      SOC: %d%%\n", hole.SOC)
			fmt.Printf("      Temperature: %d°C\n", hole.Temperature)
			fmt.Printf("      ChargeVolt: %.1fV\n", hole.ChargeVolt)
			fmt.Printf("      ChargeCurr: %.1fA\n", hole.ChargeCurr)
			fmt.Printf("      SoftVersion: %d\n", hole.SoftVersion)
			fmt.Printf("      Sensor: 0x%02X\n", hole.Sensor)
		}
	}
}

func main() {
	// Test the ParseCheckResponse function
	testParseCheckResponse()

	// Test the ParsePowerBankUploadResponse function
	testParsePowerBankUploadResponse()

	// Original main function code
	service := powerbankSdk.NewServer(powerbankModels.ServerInput{
		Host:     "103.50.205.106",
		Port:     "1883",
		Username: "backend",
		Password: "Mongol123@",
		CallbackSubscribe: func(typ constants.PUBLISH_TYPE, clientID string, msg interface{}) {
			fmt.Println(typ, clientID, msg)
		},

		CallbackPublish: func(msg mqtt.Message) {
			fmt.Println(string(msg.Payload()))
		},
	})

	// fmt.Println(service)

	// service.Publish(powerbankModels.PublishInput{
	// 	ClientID:    "864601068412899",
	// 	PublishType: constants.PUBLISH_TYPE_POPUP,
	// 	Data:        "85021618",
	// })

	service.Publish(powerbankModels.PublishInput{
		ClientID:    "864601069972081",
		PublishType: constants.PUBLISH_TYPE_CHECK,
	})

	// // Keep the program running indefinitely
	select {}
}
