package powerbankUtils

import (
	"fmt"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/techpartners-asia/powerbank/constants"
	powerbankModels "github.com/techpartners-asia/powerbank/models"
)

func ParsePowerBankUploadResponse(data []byte) (*powerbankModels.PowerBankUploadResponse, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("data too short: expected at least 4 bytes, got %d", len(data))
	}

	resp := &powerbankModels.PowerBankUploadResponse{
		Head:   data[0],
		Length: int(data[1])<<8 | int(data[2]),
		Cmd:    data[3],
	}

	if resp.Head != 0xA8 {
		return nil, fmt.Errorf("invalid header: expected 0xA8, got 0x%02X", resp.Head)
	}
	if resp.Cmd != 0x10 {
		return nil, fmt.Errorf("invalid command: expected 0x10, got 0x%02X", resp.Cmd)
	}
	if len(data) != resp.Length {
		return nil, fmt.Errorf("data length mismatch: expected %d, got %d", resp.Length, len(data))
	}

	currentPos := 4
	for currentPos+6 <= len(data)-1 { // Ensure at least 6 bytes for control board
		board := powerbankModels.ControlBoard{
			ControlIndex: int(data[currentPos]),
			Undefined1:   int(data[currentPos+1]),
			Undefined2:   int(data[currentPos+2]),
			Temperature:  int(data[currentPos+3]),
			SoftVersion:  int(data[currentPos+4]),
			HardVersion:  int(data[currentPos+5]),
		}
		currentPos += 6

		holes := make([]powerbankModels.Hole, 0)

		for currentPos+15 <= len(data)-1 {
			// Check for next control board based on known pattern: address 0x10, 0x20, 0x30
			if data[currentPos] == 0x10 || data[currentPos] == 0x20 || data[currentPos] == 0x30 {
				// Additional sanity check: next 6 bytes should not exceed length
				if currentPos+6 <= len(data)-1 {
					break
				}
			}

			hole := powerbankModels.Hole{
				HoleIndex:     int(data[currentPos]),
				State:         int(data[currentPos+1]),
				PowerbankCurr: float64(data[currentPos+2]) / 10.0,
				PowerbankVolt: float64(data[currentPos+3]) / 10.0,
				Area:          int(data[currentPos+4]),
				PowerbankSN:   strconv.FormatUint(uint64(data[currentPos+5])<<24|uint64(data[currentPos+6])<<16|uint64(data[currentPos+7])<<8|uint64(data[currentPos+8]), 10),
				SOC:           int(data[currentPos+9]),
				Temperature:   int(data[currentPos+10]),
				ChargeVolt:    float64(data[currentPos+11]) / 10.0,
				ChargeCurr:    float64(data[currentPos+12]) / 10.0,
				SoftVersion:   int(data[currentPos+13]),
				Sensor:        data[currentPos+14],
			}
			holes = append(holes, hole)
			currentPos += 15
		}

		board.Holes = holes
		resp.ControlBoards = append(resp.ControlBoards, board)
	}

	if len(data) > 0 {
		resp.Verify = data[len(data)-1]
	}

	return resp, nil
}

func ParseReturnPowerBankResponse(response []byte) (*powerbankModels.PowerBankReturnResponse, error) {
	if len(response) < 16 {
		return nil, fmt.Errorf("invalid data length: expected at least 16 bytes, got %d", len(response))
	}

	packetLength := int(response[1])<<8 | int(response[2])
	expectedTotalLength := packetLength + 1 // packet length excludes head byte

	if len(response) != expectedTotalLength {
		return nil, fmt.Errorf("data length mismatch: expected %d, got %d", expectedTotalLength, len(response))
	}

	return &powerbankModels.PowerBankReturnResponse{
		Head:         response[0],
		Length:       packetLength,
		Cmd:          response[3],
		ControlIndex: int(response[4]),
		HoleIndex:    int(response[5]),
		Area:         int(response[6]),
		PowerbankSN:  strconv.FormatUint(uint64(response[7])<<24|uint64(response[8])<<16|uint64(response[9])<<8|uint64(response[10]), 10),
		State:        int(response[11]),
		SoftVersion:  int(response[12]),
		SOC:          int(response[13]),
		Verify:       response[14],
	}, nil
}
func ParsePopupPowerBankResponse(response []byte) (*powerbankModels.PowerBankPopupResponse, error) {
	if len(response) < 4 {
		return nil, fmt.Errorf("data too short: expected at least 4 bytes, got %d", len(response))
	}

	packetLength := int(response[1])<<8 | int(response[2])
	expectedTotalLength := packetLength + 1 // Length excludes Head

	if len(response) != expectedTotalLength {
		return nil, fmt.Errorf("invalid data length: expected %d, got %d", expectedTotalLength, len(response))
	}

	if packetLength < 9 {
		return nil, fmt.Errorf("invalid packet length: expected at least 9, got %d", packetLength)
	}

	return &powerbankModels.PowerBankPopupResponse{
		Head:          response[0],
		Length:        packetLength,
		Cmd:           response[3],
		ControlIndex:  int(response[4]),
		PowerbankSN:   strconv.FormatUint(uint64(response[5])<<24|uint64(response[6])<<16|uint64(response[7])<<8|uint64(response[8]), 10),
		State:         int(response[9]),
		SolenoidValve: int(response[10]),
		Verify:        response[11],
	}, nil
}
func ParseCheckResponse(response []byte) (*powerbankModels.PowerBankCheckResponse, error) {
	if len(response) < 9 {
		return nil, fmt.Errorf("invalid data length: expected at least 9 bytes, got %d", len(response))
	}

	checkResponse := &powerbankModels.PowerBankCheckResponse{
		Head:   response[0],
		Length: int(response[1])<<8 | int(response[2]),
		Cmd:    response[3],
	}

	offset := 4
	controlBoards := make([]powerbankModels.ControlBoard, 0)

	for offset+6 <= len(response)-1 { // Ensure at least 6 bytes for control board
		controlBoard := powerbankModels.ControlBoard{
			ControlIndex: int(response[offset]),
			Undefined1:   int(response[offset+1]),
			Undefined2:   int(response[offset+2]),
			Temperature:  int(response[offset+3]),
			SoftVersion:  int(response[offset+4]),
			HardVersion:  int(response[offset+5]),
		}
		offset += 6

		holes := make([]powerbankModels.Hole, 0)
		for offset+15 <= len(response)-1 {
			// Check if next 6 bytes look like a control board start
			// Based on your example, valid control board addresses are 0x10, 0x20, 0x30
			if response[offset] == 0x10 || response[offset] == 0x20 || response[offset] == 0x30 {
				// Additional sanity checks possible here
				break
			}

			hole := powerbankModels.Hole{
				HoleIndex:     int(response[offset]),
				State:         int(response[offset+1]),
				PowerbankCurr: float64(response[offset+2]) / 10,
				PowerbankVolt: float64(response[offset+3]) / 10,
				Area:          int(response[offset+4]),
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

func ParseResponse(msg mqtt.Message) (constants.PUBLISH_TYPE, interface{}, error) {
	payload := msg.Payload()
	if len(payload) < 4 {
		fmt.Printf("Invalid response length: %d\n", len(payload))
		return "", nil, fmt.Errorf("invalid response length: expected at least 4 bytes, got %d", len(payload))
	}

	// Check command type from byte[3]
	cmd := payload[3]
	switch cmd {
	case 0x10: // Check command response
		response, err := ParseCheckResponse(payload)
		if err != nil {
			fmt.Printf("Error parsing check response: %v\n", err)
			return "", nil, err
		}
		return constants.PUBLISH_TYPE_CHECK, response, nil
	case 0x31: // Popup command response
		response, err := ParsePopupPowerBankResponse(payload)
		if err != nil {
			fmt.Printf("Error parsing popup response: %v\n", err)
			return "", nil, err
		}
		return constants.PUBLISH_TYPE_POPUP, response, nil
	case 0x40: // Return command response
		response, err := ParseReturnPowerBankResponse(payload)
		if err != nil {
			fmt.Printf("Error parsing return response: %v\n", err)
			return "", nil, err
		}
		return constants.PUBLISH_TYPE_RETURN, response, nil
	default:
		fmt.Printf("Unknown command type: 0x%02X\n", cmd)
		return "", nil, fmt.Errorf("unknown command type: 0x%02X", cmd)
	}

}
