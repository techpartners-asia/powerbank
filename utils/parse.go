package powerbankUtils

import (
	"fmt"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	powerbankModels "github.com/techpartners-asia/powerbank/models"
)

func ParseReturnPowerBankResponse(response []byte) (*powerbankModels.PowerBankReturnResponse, error) {
	if len(response) < 9 {
		return nil, fmt.Errorf("invalid data length: expected at least 9 bytes, got %d", len(response))
	}

	return &powerbankModels.PowerBankReturnResponse{
		Head:         response[0],
		Length:       int(response[1]<<8 | response[2]),
		Cmd:          response[3],
		ControlIndex: int(response[4]),
		HoleIndex:    int(response[5]),
		State:        int(response[6]),
		Undefined1:   int(response[7]),
		Undefined2:   int(response[8]),
		Area:         string(response[9]),
		PowerbankSN:  strconv.FormatUint(uint64(response[10])<<24|uint64(response[11])<<16|uint64(response[12])<<8|uint64(response[13]), 10),
		SOC:          int(response[14]),
		Temperature:  int(response[15]),
		ChargeVolt:   float64(response[16]) / 10,
		ChargeCurr:   float64(response[17]) / 10,
		SoftVersion:  int(response[18]),
		HardVersion:  int(response[19]),
		Verify:       response[20],
	}, nil
}

func ParsePopupPowerBankResponse(response []byte) (*powerbankModels.PowerBankPopupResponse, error) {
	if len(response) < 21 {
		return nil, fmt.Errorf("invalid data length: expected at least 9 bytes, got %d", len(response))
	}

	return &powerbankModels.PowerBankPopupResponse{
		Head:          response[0],
		Length:        int(response[1]<<8 | response[2]),
		Cmd:           response[3],
		ControlIndex:  int(response[4]),
		PowerbankSN:   strconv.FormatUint(uint64(response[5])<<24|uint64(response[6])<<16|uint64(response[7])<<8|uint64(response[8]), 10),
		State:         int(response[6]),
		SolenoidValve: int(response[7]),
		Verify:        response[8],
	}, nil
}

func ParseCheckResponse(response []byte) (*powerbankModels.PowerBankCheckResponse, error) {
	if len(response) < 9 {
		return nil, fmt.Errorf("invalid data length: expected at least 9 bytes, got %d", len(response))
	}

	// Parse header information
	checkResponse := &powerbankModels.PowerBankCheckResponse{
		Head:   response[0],
		Length: int(response[1]<<8 | response[2]),
		Cmd:    response[3],
	}

	// Start parsing from byte 4 (after header)
	offset := 4
	controlBoards := make([]powerbankModels.ControlBoard, 0)

	// Parse control boards and their holes
	for offset < len(response)-1 { // -1 to leave room for verify byte
		// Parse control board (6 bytes)
		if offset+6 > len(response) {
			return nil, fmt.Errorf("invalid data length: expected control board data at offset %d", offset)
		}

		controlBoard := powerbankModels.ControlBoard{
			ControlIndex: int(response[offset]),
			Undefined1:   int(response[offset+1]),
			Undefined2:   int(response[offset+2]),
			Temperature:  int(response[offset+3]),
			SoftVersion:  int(response[offset+4]),
			HardVersion:  int(response[offset+5]),
		}
		offset += 6

		// Parse holes (15 bytes each)
		holes := make([]powerbankModels.Hole, 0)
		for offset+15 <= len(response)-1 { // -1 to leave room for verify byte
			// Check if next bytes are control board data
			if offset+6 <= len(response)-1 && response[offset] >= 0x10 && response[offset] <= 0x30 {
				break
			}

			hole := powerbankModels.Hole{
				HoleIndex:     int(response[offset]),
				State:         int(response[offset+1]),
				PowerbankCurr: float64(response[offset+2]) / 10,
				PowerbankVolt: float64(response[offset+3]) / 10,
				Area:          string(response[offset+4]),
				PowerbankSN:   strconv.FormatUint(uint64(response[offset+5])<<24|uint64(response[offset+6])<<16|uint64(response[offset+7])<<8|uint64(response[offset+8]), 10),
				SOC:           int(response[offset+9]),
				Temperature:   int(response[offset+10]),
				ChargeVolt:    float64(response[offset+11]) / 10,
				ChargeCurr:    float64(response[offset+12]) / 10,
				SoftVersion:   int(response[offset+13]),
				Sensor:        response[offset+14],
			}
			holes = append(holes, hole)
			offset += 15
		}

		controlBoard.Holes = holes
		controlBoards = append(controlBoards, controlBoard)
	}

	checkResponse.ControlBoards = controlBoards
	checkResponse.Verify = response[len(response)-1]

	return checkResponse, nil
}

func ParseResponse(msg mqtt.Message) (interface{}, error) {
	payload := msg.Payload()
	if len(payload) < 4 {
		fmt.Printf("Invalid response length: %d\n", len(payload))
		return nil, fmt.Errorf("invalid response length: expected at least 4 bytes, got %d", len(payload))
	}

	// Check command type from byte[3]
	cmd := payload[3]
	switch cmd {
	case 0x10: // Check command response
		response, err := ParseCheckResponse(payload)
		if err != nil {
			fmt.Printf("Error parsing check response: %v\n", err)
			return nil, err
		}
		return response, nil
	case 0x31: // Popup command response
		response, err := ParsePopupPowerBankResponse(payload)
		if err != nil {
			fmt.Printf("Error parsing popup response: %v\n", err)
			return nil, err
		}
		return response, nil
	case 0x28: // Return command response
		response, err := ParseReturnPowerBankResponse(payload)
		if err != nil {
			fmt.Printf("Error parsing return response: %v\n", err)
			return nil, err
		}
		return response, nil
	default:
		fmt.Printf("Unknown command type: 0x%02X\n", cmd)
		return nil, fmt.Errorf("unknown command type: 0x%02X", cmd)
	}

}
