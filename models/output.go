package powerbankModels

import "github.com/techpartners-asia/powerbank/constants"

const (
	PopupFailed  = 0x00
	PopupSuccess = 0x01

	ReturnFailed  = 0x00
	ReturnSuccess = 0x01
)

type (

	// PowerBankPopupResponse represents the byte protocol response for power bank pop-up
	PowerBankPopupResponse struct {
		Head          byte   // Byte[0] - Head code (Default: 0xA8)
		Length        int    // Byte[1-2] - Packet length
		Cmd           byte   // Byte[3] - Command name (Default: 0x31)
		ControlIndex  int    // Byte[4] - Control board address
		PowerbankSN   string // Byte[5-8] - Power bank SN
		State         int    // Byte[6] - Popup state
		SolenoidValve int    // Byte[7] - Solenoid valve status
		Verify        byte   // Byte[8] - Check code
	}

	// PowerBankReturnResponse represents the byte protocol response for power bank return
	PowerBankReturnResponse struct {
		Head         byte    // Byte[0] - Header code (Default: 0xA8)
		Length       int     // Byte[1-2] - Packet length (Default: 0x0015 = 21)
		Cmd          byte    // Byte[3] - Command name (Default: 0x28)
		ControlIndex int     // Byte[4] - Movement board address (Default: 0x10)
		HoleIndex    int     // Byte[5] - Position address (Default: 0x01)
		State        int     // Byte[6] - Return status
		Undefined1   int     // Byte[7] - Reserved 1 (Default: 0x00)
		Undefined2   int     // Byte[8] - Reserved 2 (Default: 0x00)
		Area         string  // Byte[9] - Area code
		PowerbankSN  string  // Byte[10-13] - Power bank SN
		SOC          int     // Byte[14] - Power (0-255%)
		Temperature  int     // Byte[15] - Temperature (0-100℃)
		ChargeVolt   float64 // Byte[16] - Charging voltage (1 decimal place)
		ChargeCurr   float64 // Byte[17] - Charging current (1 decimal place)
		SoftVersion  int     // Byte[18] - Software version number
		HardVersion  int     // Byte[19] - Hardware version number
		Verify       byte    // Byte[20] - Verification code
	}

	CreateUserResponse struct {
		UserID string `json:"user_id"`
	}

	// PowerBankCheckResponse represents the byte protocol response for cabinet check
	PowerBankCheckResponse struct {
		Head          byte           // Byte[0] - Head code (Default: 0xA8)
		Length        int            // Byte[1-2] - Packet length
		Cmd           byte           // Byte[3] - Command name (Default: 0x10)
		ControlBoards []ControlBoard // Control board information
		Verify        byte           // Last byte - Check code
	}

	// ControlBoard represents a single control board's information
	ControlBoard struct {
		ControlIndex int    // Byte[0] - Control board address
		Undefined1   int    // Byte[1] - Reserved 1
		Undefined2   int    // Byte[2] - Reserved 2
		Temperature  int    // Byte[3] - Temperature
		SoftVersion  int    // Byte[4] - Software version
		HardVersion  int    // Byte[5] - Hardware version
		Holes        []Hole // Position information
	}

	// Hole represents a single position's information
	Hole struct {
		HoleIndex     int     // Byte[0] - Position address
		State         int     // Byte[1] - State information
		PowerbankCurr float64 // Byte[2] - Power bank current
		PowerbankVolt float64 // Byte[3] - Power bank voltage
		Area          string  // Byte[4] - Area code
		PowerbankSN   string  // Byte[5-8] - Power bank SN
		SOC           int     // Byte[9] - Battery percentage
		Temperature   int     // Byte[10] - Temperature
		ChargeVolt    float64 // Byte[11] - Charging voltage
		ChargeCurr    float64 // Byte[12] - Charging current
		SoftVersion   int     // Byte[13] - Software version
		Sensor        byte    // Byte[14] - Position detection
	}

	PowerBankUploadResponse struct {
		Head          byte           // Byte[0] - Head code (Default: 0xA8)
		Length        int            // Byte[1-2] - Packet length
		Cmd           byte           // Byte[3] - Command name (Default: 0x10)
		ControlBoards []ControlBoard // Byte[4~n] - Control board information
		Verify        byte           // Byte[n+1] - Check code
	}
)

func (hole *Hole) GetStateDescription() string {
	switch hole.State {
	case 0x00:
		return "No mobile power supply"
	case 0x01:
		return "Power bank is normal"
	case 0x02:
		return "Charging abnormality"
	case 0x03:
		return "Communication exception"
	case 0x04:
		return "KaBao/Damaged"
	case 0x05:
		return "The key is forcibly released"
	case 0x06:
		return "The solenoid valve did not return to the position when returned"
	case 0x07:
		return "Reserved"
	case 0x08:
		return "Anti-theft protocol communication failed"
	case 0x09:
		return "Typec short circuit"
	case 0x0A:
		return "Return failed, battery does not pop out"
	default:
		return "Reserved"
	}
}

func (hole *Hole) GetStatus() constants.PowerbankStatus {

	switch hole.State {
	case 0x00:
		return constants.PowerbankStatus_NoPowerSupply
	case 0x01:
		return constants.PowerbankStatus_Normal
	case 0x02:
		return constants.PowerbankStatus_ChargingAbnormality
	case 0x03:
		return constants.PowerbankStatus_CommunicationException
	case 0x04:
		return constants.PowerbankStatus_KabaoDamaged
	case 0x05:
		return constants.PowerbankStatus_KeyForceRelease
	case 0x06:
		return constants.PowerbankStatus_SolenoidValveNotReturn
	case 0x07:
		return constants.PowerbankStatus_Reserved
	case 0x08:
		return constants.PowerbankStatus_AntiTheftProtocolCommunicationFailed
	case 0x09:
		return constants.PowerbankStatus_TypecShortCircuit
	case 0x0A:
		return constants.PowerbankStatus_ReturnFailedBatteryDoesNotPopOut
	default:
		return constants.PowerbankStatus_Reserved
	}
}

func (popup *PowerBankPopupResponse) GetStatus() constants.PowerbankStatus {

	switch popup.State {
	case 0x00:
		return constants.PowerbankStatus_PopupFailed
	case 0x01:
		return constants.PowerbankStatus_PopupSuccessful
	case 0x02:
		return constants.PowerbankStatus_PowerSupplyChargingAbnormally
	case 0x03:
		return constants.PowerbankStatus_CommunicationAbnormalityFirstReturnFailed
	case 0x04:
		return constants.PowerbankStatus_SlotCannotPopOut
	case 0x05:
		return constants.PowerbankStatus_SlotForciblyReleased
	case 0x06:
		return constants.PowerbankStatus_SolenoidNotReturned
	case 0x08:
		return constants.PowerbankStatus_AntiTheftCommFailed
	case 0x11:
		return constants.PowerbankStatus_FailedToObtainSn
	case 0x12:
		return constants.PowerbankStatus_PopupCompleteMotorHomeSnReadable
	case 0x13:
		return constants.PowerbankStatus_FailedToObtainTraceback
	case 0x14:
		return constants.PowerbankStatus_BatteryLockCommandFailed
	case 0x21:
		return constants.PowerbankStatus_SnAcquisitionAndMotorFailed
	case 0x22:
		return constants.PowerbankStatus_InfoAcquisitionAndMotorFailed
	case 0x23:
		return constants.PowerbankStatus_BatteryLockAndMotorFailed
	case 0x24:
		return constants.PowerbankStatus_AntiTheftSwitchDetectionFailed
	default:
		return constants.PowerbankStatus_UnknownError
	}
}

// GetDescription converts the byte held in popup.State into a
// human‑readable explanation of what went wrong (or right).
func (popup *PowerBankPopupResponse) GetDescription() string {
	switch popup.State {
	case 0x00:
		return "Pop‑up failed"
	case 0x01:
		return "Pop‑up successful"
	case 0x02:
		return "Power‑supply charging abnormally"
	case 0x03:
		return "Communication abnormality (first return failed)"
	case 0x04:
		return "This slot cannot pop out a mobile power supply normally"
	case 0x05:
		return "Slot forcibly released"
	case 0x06:
		return "Solenoid valve did not return to home position"
	case 0x08:
		return "Anti‑theft protocol communication failed"
	case 0x11:
		return "Failed to obtain SN"
	case 0x12:
		return "Pop‑up completed; motor is home and SN can be read"
	case 0x13:
		return "Failed to obtain traceback information"
	case 0x14:
		return "Battery‑lock command failed"
	case 0x21:
		return "Failed to obtain SN and motor action failed"
	case 0x22:
		return "Failed to obtain all information and motor operation failed"
	case 0x23:
		return "Battery‑lock command failed and motor action failed"
	case 0x24:
		return "Anti‑theft‑switch detection failed"
	default:
		return "Unknown error"
	}
}

func (rt *PowerBankReturnResponse) GetDescription() string {
	switch rt.State {
	case 0x00:
		return "Return failed"
	case 0x01:
		return "Return successful"
	case 0x11:
		return "Failed to obtain SN"
	case 0x12:
		return "Failed to obtain voltage, temperature, or other information"
	case 0x13:
		return "Failed to obtain software and hardware version information"
	case 0x14:
		return "Battery‑lock command failed"
	case 0x21:
		return "Failed to obtain SN, and motor action failed"
	case 0x22:
		return "Failed to obtain voltage, temperature, or other information, and motor action failed"
	case 0x23:
		return "Battery‑lock command failed, and motor action failed"
	case 0x24:
		return "Anti‑theft switch detection failed (within 5 minutes, an 0x28 self‑test command will be reported)"
	default:
		return "Unknown error"
	}
}

func (rt *PowerBankReturnResponse) GetStatus() constants.PowerbankStatus {
	switch rt.State {
	case 0x00:
		return constants.PowerbankStatus_ReturnFailed
	case 0x01:
		return constants.PowerbankStatus_ReturnSuccessful
	case 0x11:
		return constants.PowerbankStatus_FailedToObtainSn
	case 0x12:
		return constants.PowerbankStatus_FailedToObtainVoltageTemperatureOrOtherInformation
	case 0x13:
		return constants.PowerbankStatus_FailedToObtainSoftwareAndHardwareVersionInformation
	case 0x14:
		return constants.PowerbankStatus_BatteryLockCommandFailed
	case 0x21:
		return constants.PowerbankStatus_FailedToObtainSnAndMotorActionFailed
	case 0x22:
		return constants.PowerbankStatus_FailedToObtainVoltageTemperatureOrOtherInformationAndMotorActionFailed
	case 0x23:
		return constants.PowerbankStatus_BatteryLockCommandFailedAndMotorActionFailed
	case 0x24:
		return constants.PowerbankStatus_AntiTheftSwitchDetectionFailedAndMotorActionFailed
	}
	return constants.PowerbankStatus_UnknownError
}
