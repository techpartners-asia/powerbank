package powerbankUtils

import (
	"fmt"
	"strconv"

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

// func ParseCheckResponse(response []byte) (*models.PowerBankCheckResponse, error) {
// 	if len(response) < 9 {
// 		return nil, fmt.Errorf("invalid data length: expected at least 9 bytes, got %d", len(response))
// 	}

// 	return &models.PowerBankCheckResponse{
// 		Head:         response[0],
// 		Length:       int(response[1]<<8 | response[2]),
// 		Cmd:          response[3],
// 		ControlIndex: int(response[4]),
// 	}, nil
// }
