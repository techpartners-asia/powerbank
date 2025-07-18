package powerbankUtils

import (
	"fmt"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/techpartners-asia/powerbank/constants"
	powerbankModels "github.com/techpartners-asia/powerbank/models"
)

func ParsePowerBankUploadResponse(data []byte) (*powerbankModels.PowerBankUploadResponse, error) {
	if len(data) < 5 {
		return nil, fmt.Errorf("invalid data length: expected at least 5 bytes, got %d", len(data))
	}

	resp := &powerbankModels.PowerBankUploadResponse{
		Head:   data[0],
		Length: int(data[1])<<8 | int(data[2]),
		Cmd:    data[3],
	}

	// Verify header and command
	if resp.Head != 0xA8 {
		return nil, fmt.Errorf("invalid header: expected 0xA8, got 0x%02X", resp.Head)
	}
	if resp.Cmd != 0x10 {
		return nil, fmt.Errorf("invalid command: expected 0x10, got 0x%02X", resp.Cmd)
	}

	pos := 4 // Start after header, length, cmd

	// Parse until we reach the verification byte (last byte)
	for pos < len(data)-1 {
		// Check if we have enough bytes for a control board (6 bytes)
		if pos+6 > len(data)-1 {
			break
		}

		// Parse Control Board (6 bytes)
		cb := powerbankModels.ControlBoard{
			ControlIndex: int(data[pos]),
			Undefined1:   int(data[pos+1]),
			Undefined2:   int(data[pos+2]),
			Temperature:  int(data[pos+3]),
			SoftVersion:  int(data[pos+4]),
			HardVersion:  int(data[pos+5]),
		}
		pos += 6

		holes := make([]powerbankModels.Hole, 0)

		// Parse holes for this control board (4 holes per control board, 15 bytes each)
		for i := 0; i < 4 && pos+15 <= len(data)-1; i++ {
			h := powerbankModels.Hole{
				HoleIndex:     int(data[pos]),
				State:         int(data[pos+1]),
				PowerbankCurr: float64(data[pos+2]) / 10,
				PowerbankVolt: float64(data[pos+3]) / 10,
				Area:          int(data[pos+4]),
				PowerbankSN: fmt.Sprintf("%d",
					uint32(data[pos+5])<<24|
						uint32(data[pos+6])<<16|
						uint32(data[pos+7])<<8|
						uint32(data[pos+8])),
				SOC:         int(data[pos+9]),
				Temperature: int(data[pos+10]),
				ChargeVolt:  float64(data[pos+11]) / 10,
				ChargeCurr:  float64(data[pos+12]) / 10,
				SoftVersion: int(data[pos+13]),
				Sensor:      data[pos+14],
			}
			holes = append(holes, h)
			pos += 15
		}

		cb.Holes = holes
		resp.ControlBoards = append(resp.ControlBoards, cb)
	}

	// Final byte is verification
	if len(data) > 0 {
		resp.Verify = data[len(data)-1]
	}

	return resp, nil
}

func ParseReturnPowerBankResponse(response []byte) (*powerbankModels.PowerBankReturnResponse, error) {
	if len(response) < 15 {
		return nil, fmt.Errorf("invalid data length: expected at least 15 bytes, got %d", len(response))
	}

	return &powerbankModels.PowerBankReturnResponse{
		Head:         response[0],
		Length:       int(response[1]<<8 | response[2]),
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
	if len(response) < 9 {
		return nil, fmt.Errorf("invalid data length: expected at least 9 bytes, got %d", len(response))
	}

	return &powerbankModels.PowerBankPopupResponse{
		Head:          response[0],
		Length:        int(response[1]<<8 | response[2]),
		Cmd:           response[3],
		ControlIndex:  int(response[4]),
		PowerbankSN:   strconv.FormatUint(uint64(response[5])<<24|uint64(response[6])<<16|uint64(response[7])<<8|uint64(response[8]), 10),
		State:         int(response[9]),
		SolenoidValve: int(response[10]),
		Verify:        response[11],
	}, nil
}
func ParseCheckResponse(response []byte) (*powerbankModels.PowerBankCheckResponse, error) {
	if len(response) < 5 {
		return nil, fmt.Errorf("invalid data length: expected at least 5 bytes, got %d", len(response))
	}

	resp := &powerbankModels.PowerBankCheckResponse{
		Head:   response[0],
		Length: int(response[1])<<8 | int(response[2]),
		Cmd:    response[3],
	}

	pos := 4 // Start after header, length, cmd

	// Parse until we reach the verification byte (last byte)
	for pos < len(response)-1 {
		// Check if we have enough bytes for a control board (6 bytes)
		if pos+6 > len(response)-1 {
			break
		}

		// Parse Control Board (6 bytes)
		cb := powerbankModels.ControlBoard{
			ControlIndex: int(response[pos]),
			Undefined1:   int(response[pos+1]),
			Undefined2:   int(response[pos+2]),
			Temperature:  int(response[pos+3]),
			SoftVersion:  int(response[pos+4]),
			HardVersion:  int(response[pos+5]),
		}
		pos += 6

		holes := make([]powerbankModels.Hole, 0)

		// Parse holes for this control board (4 holes per control board, 15 bytes each)
		for i := 0; i < 4 && pos+15 <= len(response)-1; i++ {
			h := powerbankModels.Hole{
				HoleIndex:     int(response[pos]),
				State:         int(response[pos+1]),
				PowerbankCurr: float64(response[pos+2]) / 10,
				PowerbankVolt: float64(response[pos+3]) / 10,
				Area:          int(response[pos+4]),
				PowerbankSN: fmt.Sprintf("%d",
					uint32(response[pos+5])<<24|
						uint32(response[pos+6])<<16|
						uint32(response[pos+7])<<8|
						uint32(response[pos+8])),
				SOC:         int(response[pos+9]),
				Temperature: int(response[pos+10]),
				ChargeVolt:  float64(response[pos+11]) / 10,
				ChargeCurr:  float64(response[pos+12]) / 10,
				SoftVersion: int(response[pos+13]),
				Sensor:      response[pos+14],
			}
			holes = append(holes, h)
			pos += 15
		}

		cb.Holes = holes
		resp.ControlBoards = append(resp.ControlBoards, cb)
	}

	// Final byte is verification
	if len(response) > 0 {
		resp.Verify = response[len(response)-1]
	}

	return resp, nil
}

func ParseHealthCheckResponse(msg mqtt.Message) (*powerbankModels.PowerBankHealthCheckResponse, error) {
	response := msg.Payload()
	if len(response) < 9 {
		return nil, fmt.Errorf("invalid data length: expected at least 9 bytes, got %d", len(response))
	}

	return &powerbankModels.PowerBankHealthCheckResponse{
		Head:         response[0],
		Length:       int(response[1]<<8 | response[2]),
		Cmd:          response[3],
		ControlIndex: int(response[4]),
		Signal:       string(response[5:]),
		Verify:       response[len(response)-1],
	}, nil
}

func ParseResponse(msg mqtt.Message) (constants.PUBLISH_TYPE, interface{}, error) {
	payload := msg.Payload()
	if len(payload) < 4 {
		fmt.Printf("Invalid response length: %d\n", len(payload))
		return "", nil, fmt.Errorf("invalid response length: expected at least 4 bytes, got %d", len(payload))
	}

	fmt.Printf("Payload: % X\n", payload)

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
	// case 0x7A: // Health check command response
	// 	response, err := ParseHealthCheckResponse(payload)
	// 	if err != nil {
	// 		fmt.Printf("Error parsing health check response: %v\n", err)
	// 		return "", nil, err
	// 	}
	// 	return constants.PUBLISH_TYPE_HEALTH_CHECK, response, nil
	default:
		fmt.Printf("Unknown command type: 0x%02X\n", cmd)
		return "", nil, fmt.Errorf("unknown command type: 0x%02X", cmd)
	}

}
