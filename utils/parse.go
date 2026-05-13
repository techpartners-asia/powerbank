package powerbankUtils

import (
	"fmt"
	"strconv"

	"github.com/techpartners-asia/powerbank/constants"
	powerbankModels "github.com/techpartners-asia/powerbank/models"
)

// ParsePowerBankUploadResponse parses the upload_all (0x10) cabinet info frame.
// The layout is identical to the check response, so it delegates to ParseCheckResponse.
func ParsePowerBankUploadResponse(data []byte) (*powerbankModels.PowerBankUploadResponse, error) {
	return ParseCheckResponse(data)
}

func ParseReturnPowerBankResponse(response []byte) (*powerbankModels.PowerBankReturnResponse, error) {
	if len(response) < 15 {
		return nil, fmt.Errorf("invalid data length: expected at least 15 bytes, got %d", len(response))
	}

	return &powerbankModels.PowerBankReturnResponse{
		Head:         response[0],
		Length:       int(response[1])<<8 | int(response[2]),
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
func ParseReturnFixPowerBankResponse(response []byte) (*powerbankModels.PowerBankReturnFixResponse, error) {
	if len(response) < 21 {
		return nil, fmt.Errorf("invalid data length: expected at least 21 bytes, got %d", len(response))
	}

	return &powerbankModels.PowerBankReturnFixResponse{
		Head:         response[0],
		Length:       int(response[1])<<8 | int(response[2]),
		Cmd:          response[3],
		ControlIndex: int(response[4]),
		HoleIndex:    int(response[5]),
		State:        int(response[6]),
		Reserved1:    response[7],
		Reserved2:    response[8],
		Area:         int(response[9]),
		PowerbankSN:  strconv.FormatUint(uint64(response[10])<<24|uint64(response[11])<<16|uint64(response[12])<<8|uint64(response[13]), 10),
		SOC:          int(response[14]),
		Temperature:  int(response[15]),
		ChargeVolt:   float64(response[16]) / 10.0,
		ChargeCurr:   float64(response[17]) / 10.0,
		SoftVersion:  int(response[18]),
		HardVersion:  int(response[19]),
		Verify:       response[20],
	}, nil
}

func ParsePopupByHolePowerBankResponse(response []byte) (*powerbankModels.PowerBankPopupByHoleResponse, error) {
	if len(response) < 9 {
		return nil, fmt.Errorf("invalid data length: expected at least 9 bytes, got %d", len(response))
	}

	return &powerbankModels.PowerBankPopupByHoleResponse{
		Head:         response[0],
		Length:       int(response[1])<<8 | int(response[2]),
		Cmd:          response[3],
		ControlIndex: int(response[4]),
		HoleIndex:    int(response[5]),
		State:        int(response[6]),
		Reserved:     response[7],
		Verify:       response[8],
	}, nil
}

func ParsePopupPowerBankResponse(response []byte) (*powerbankModels.PowerBankPopupResponse, error) {
	if len(response) < 9 {
		return nil, fmt.Errorf("invalid data length: expected at least 9 bytes, got %d", len(response))
	}

	return &powerbankModels.PowerBankPopupResponse{
		Head:          response[0],
		Length:        int(response[1])<<8 | int(response[2]),
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

func ParseHealthCheckResponse(response []byte) (*powerbankModels.PowerBankHealthCheckResponse, error) {
	if len(response) < 9 {
		return nil, fmt.Errorf("invalid data length: expected at least 9 bytes, got %d", len(response))
	}

	return &powerbankModels.PowerBankHealthCheckResponse{
		Head:         response[0],
		Length:       int(response[1])<<8 | int(response[2]),
		Cmd:          response[3],
		ControlIndex: int(response[4]),
		Signal:       string(response[5 : len(response)-1]),
		Verify:       response[len(response)-1],
	}, nil
}

// Debug enables verbose stdout traces in ParseResponse. Off by default.
// Set from the host application (e.g., when constructing the server) if you
// want raw payload hex and per-cmd parse error messages.
var Debug bool

func debugf(format string, args ...interface{}) {
	if Debug {
		fmt.Printf(format, args...)
	}
}

func ParseResponse(payload []byte) (constants.PUBLISH_TYPE, interface{}, error) {
	if len(payload) < 4 {
		return "", nil, fmt.Errorf("invalid response length: expected at least 4 bytes, got %d", len(payload))
	}

	debugf("Payload: % X\n", payload)

	cmd := payload[3]
	switch cmd {
	case 0x10:
		response, err := ParseCheckResponse(payload)
		if err != nil {
			return "", nil, fmt.Errorf("parse check (0x10): %w", err)
		}
		return constants.PUBLISH_TYPE_CHECK, response, nil
	case 0x31:
		response, err := ParsePopupPowerBankResponse(payload)
		if err != nil {
			return "", nil, fmt.Errorf("parse popup (0x31): %w", err)
		}
		return constants.PUBLISH_TYPE_POPUP, response, nil
	case 0x21:
		response, err := ParsePopupByHolePowerBankResponse(payload)
		if err != nil {
			return "", nil, fmt.Errorf("parse popup-by-hole (0x21): %w", err)
		}
		return constants.PUBLISH_TYPE_POPUP_BY_HOLE, response, nil
	case 0x40:
		response, err := ParseReturnPowerBankResponse(payload)
		if err != nil {
			return "", nil, fmt.Errorf("parse return (0x40): %w", err)
		}
		return constants.PUBLISH_TYPE_RETURN, response, nil
	case 0x28:
		response, err := ParseReturnFixPowerBankResponse(payload)
		if err != nil {
			return "", nil, fmt.Errorf("parse return-fix (0x28): %w", err)
		}
		return constants.PUBLISH_TYPE_RETURN_FIX, response, nil
	default:
		return "", nil, fmt.Errorf("unknown command type: 0x%02X", cmd)
	}
}
