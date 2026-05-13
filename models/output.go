package powerbankModels

import (
	"strconv"
	"strings"
	"time"

	"github.com/techpartners-asia/powerbank/constants"
)

const (
	PopupFailed  = 0x00
	PopupSuccess = 0x01

	ReturnFailed  = 0x00
	ReturnSuccess = 0x01
)

type (
	GetUserResponse struct {
		Listener                         string      `json:"listener"`
		SubscriptionsMax                 string      `json:"subscriptions_max"`
		SendCnt                          int         `json:"send_cnt"`
		MailboxLen                       int         `json:"mailbox_len"`
		Port                             int         `json:"port"`
		InflightCnt                      int         `json:"inflight_cnt"`
		SendMsgDropped                   int         `json:"send_msg.dropped"`
		CleanStart                       bool        `json:"clean_start"`
		MqueueMax                        int         `json:"mqueue_max"`
		RecvCnt                          int         `json:"recv_cnt"`
		ExpiryInterval                   int         `json:"expiry_interval"`
		MqueueLen                        int         `json:"mqueue_len"`
		ProtoVer                         int         `json:"proto_ver"`
		SendMsg                          int         `json:"send_msg"`
		SendPkt                          int         `json:"send_pkt"`
		SendMsgDroppedExpired            int         `json:"send_msg.dropped.expired"`
		IsBridge                         bool        `json:"is_bridge"`
		SendMsgQos2                      int         `json:"send_msg.qos2"`
		MqueueDropped                    int         `json:"mqueue_dropped"`
		AwaitingRelCnt                   int         `json:"awaiting_rel_cnt"`
		AwaitingRelMax                   int         `json:"awaiting_rel_max"`
		EnableAuthn                      bool        `json:"enable_authn"`
		SendMsgQos1                      int         `json:"send_msg.qos1"`
		HeapSize                         int         `json:"heap_size"`
		IsPersistent                     bool        `json:"is_persistent"`
		Keepalive                        int         `json:"keepalive"`
		RecvMsgDroppedAwaitPubrelTimeout int         `json:"recv_msg.dropped.await_pubrel_timeout"`
		SendMsgDroppedTooLarge           int         `json:"send_msg.dropped.too_large"`
		RecvMsgQos1                      int         `json:"recv_msg.qos1"`
		SendMsgQos0                      int         `json:"send_msg.qos0"`
		RecvMsgQos2                      int         `json:"recv_msg.qos2"`
		RecvOct                          int         `json:"recv_oct"`
		SendOct                          int         `json:"send_oct"`
		ConnectedAt                      time.Time   `json:"connected_at"`
		Connected                        bool        `json:"connected"`
		RecvMsgDropped                   int         `json:"recv_msg.dropped"`
		SendMsgDroppedQueueFull          int         `json:"send_msg.dropped.queue_full"`
		IpAddress                        string      `json:"ip_address"`
		Node                             string      `json:"node"`
		Reductions                       int         `json:"reductions"`
		Username                         interface{} `json:"username"`
		Mountpoint                       interface{} `json:"mountpoint"`
		RecvPkt                          int         `json:"recv_pkt"`
		InflightMax                      int         `json:"inflight_max"`
		SubscriptionsCnt                 int         `json:"subscriptions_cnt"`
		ProtoName                        string      `json:"proto_name"`
		ClientID                         string      `json:"clientid"`
		Durable                          bool        `json:"durable"`
		RecvMsgQos0                      int         `json:"recv_msg.qos0"`
		Peerport                         int         `json:"peerport"`
		RecvMsg                          int         `json:"recv_msg"`
		CreatedAt                        time.Time   `json:"created_at"`
	}

	// PowerBankPopupByHoleResponse represents the 0x21 Pop-up By Hole response (9 bytes).
	PowerBankPopupByHoleResponse struct {
		Head         byte // Byte[0] - 0xA8
		Length       int  // Byte[1-2] - packet length
		Cmd          byte // Byte[3] - 0x21
		ControlIndex int  // Byte[4] - Control board address
		HoleIndex    int  // Byte[5] - Position number
		State        int  // Byte[6] - Pop-up state
		Reserved     byte // Byte[7] - 0x00
		Verify       byte // Byte[8] - Check code
	}

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

	// PowerBankReturnResponse represents the 0x40 standard Return response (15 bytes).
	// For the 0x28 self-test variant with charge diagnostics, see PowerBankReturnFixResponse.
	PowerBankReturnResponse struct {
		Head         byte   // Byte[0] - 0xA8
		Length       int    // Byte[1-2] - 0x000E = 14 (excludes header byte; total frame is 15 bytes)
		Cmd          byte   // Byte[3] - 0x40
		ControlIndex int    // Byte[4] - Control board address
		HoleIndex    int    // Byte[5] - Position address
		Area         int    // Byte[6] - Area code
		PowerbankSN  string // Byte[7-10] - Power bank SN
		State        int    // Byte[11] - Return status
		SoftVersion  int    // Byte[12] - Power bank version
		SOC          int    // Byte[13] - Battery percentage (0-100)
		Verify       byte   // Byte[14] - Verification code
	}

	// PowerBankReturnFixResponse represents the byte protocol response for
	// the 0x28 Return-Fix (self-test) command, reported over MQTT regardless
	// of success or failure. 21 bytes total.
	PowerBankReturnFixResponse struct {
		Head         byte    // Byte[0] - 0xA8
		Length       int     // Byte[1-2] - 0x0015 = 21
		Cmd          byte    // Byte[3] - 0x28
		ControlIndex int     // Byte[4] - Movement board address
		HoleIndex    int     // Byte[5] - Position address
		State        int     // Byte[6] - Return status
		Reserved1    byte    // Byte[7] - 0x00
		Reserved2    byte    // Byte[8] - 0x00
		Area         int     // Byte[9] - Area code
		PowerbankSN  string  // Byte[10-13] - Power bank SN
		SOC          int     // Byte[14] - Battery percentage (0-100)
		Temperature  int     // Byte[15] - Temperature (°C)
		ChargeVolt   float64 // Byte[16] - Charging voltage (0.1V units)
		ChargeCurr   float64 // Byte[17] - Charging current (0.1A units)
		SoftVersion  int     // Byte[18] - Software version
		HardVersion  int     // Byte[19] - Hardware version
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
		Area          int     // Byte[4] - Area code
		PowerbankSN   string  // Byte[5-8] - Power bank SN
		SOC           int     // Byte[9] - Battery percentage
		Temperature   int     // Byte[10] - Temperature
		ChargeVolt    float64 // Byte[11] - Charging voltage
		ChargeCurr    float64 // Byte[12] - Charging current
		SoftVersion   int     // Byte[13] - Software version
		Sensor        byte    // Byte[14] - Position detection
	}

	PowerBankHealthCheckResponse struct {
		Head         byte   // Byte[0] - Head code (Default: 0xA8)
		Length       int    // Byte[1-2] - Packet length
		Cmd          byte   // Byte[3] - Command name (Default: 0x7A)
		ControlIndex int    // Byte[4] - Control board address
		Signal       string // Byte[5-n] - Signal value
		Verify       byte   // Byte[n+1] - Check code
	}
)

// PowerBankUploadResponse is an alias for PowerBankCheckResponse — both check
// and upload_all return the same 0x10 cabinet-info layout.
type PowerBankUploadResponse = PowerBankCheckResponse

func (healthCheck *PowerBankHealthCheckResponse) GetSignalStrength() constants.CabinetSignal {
	// Parse the signal string which is in format "CSQ:26;BP:0"
	// Extract the CSQ value (signal strength)

	// If signal is "99", it means no network
	if healthCheck.Signal == "99" {
		return constants.CabinetSignal_NoNetwork
	}

	// Try to parse CSQ value from the signal string
	// The signal format is typically "CSQ:26;BP:0" where 26 is the CSQ value
	if len(healthCheck.Signal) > 0 {
		// Look for "CSQ:" prefix
		if len(healthCheck.Signal) >= 4 && healthCheck.Signal[:4] == "CSQ:" {
			// Extract the number after "CSQ:"
			csqStr := healthCheck.Signal[4:]
			// Find the semicolon or end of string
			if semicolonIndex := strings.Index(csqStr, ";"); semicolonIndex != -1 {
				csqStr = csqStr[:semicolonIndex]
			}

			// Convert to integer
			if csq, err := strconv.Atoi(csqStr); err == nil {
				// Map CSQ value to signal strength based on the table:
				// CSQ signal: 99, 0~12, 13~16, 17~20, 21~25, 26~31
				// description: No network, Very poor, Poor, Average, Good, Very good
				// Signal bars: Offline, 0, 1, 2, 3, 4

				switch {
				case csq == 99:
					return constants.CabinetSignal_NoNetwork
				case csq >= 0 && csq <= 12:
					return constants.CabinetSignal_Weak
				case csq >= 13 && csq <= 16:
					return constants.CabinetSignal_Weak
				case csq >= 17 && csq <= 20:
					return constants.CabinetSignal_Normal
				case csq >= 21 && csq <= 25:
					return constants.CabinetSignal_Better
				case csq >= 26 && csq <= 31:
					return constants.CabinetSignal_Better
				default:
					return constants.CabinetSignal_NoNetwork
				}
			}
		}

		// If the signal doesn't match expected format, try to parse as direct number
		if csq, err := strconv.Atoi(healthCheck.Signal); err == nil {
			switch {
			case csq == 99:
				return constants.CabinetSignal_NoNetwork
			case csq >= 0 && csq <= 12:
				return constants.CabinetSignal_Weak
			case csq >= 13 && csq <= 16:
				return constants.CabinetSignal_Weak
			case csq >= 17 && csq <= 20:
				return constants.CabinetSignal_Normal
			case csq >= 21 && csq <= 25:
				return constants.CabinetSignal_Better
			case csq >= 26 && csq <= 31:
				return constants.CabinetSignal_Better
			default:
				return constants.CabinetSignal_NoNetwork
			}
		}
	}

	// Default fallback
	return constants.CabinetSignal_NoNetwork
}

// GetCSQValue extracts the CSQ (signal strength) value from the signal string
// Returns -1 if parsing fails
func (healthCheck *PowerBankHealthCheckResponse) GetCSQValue() int {
	if len(healthCheck.Signal) == 0 {
		return -1
	}

	// If signal is "99", return 99
	if healthCheck.Signal == "99" {
		return 99
	}

	// Try to parse CSQ value from the signal string format "CSQ:26;BP:0"
	if len(healthCheck.Signal) >= 4 && healthCheck.Signal[:4] == "CSQ:" {
		csqStr := healthCheck.Signal[4:]
		if semicolonIndex := strings.Index(csqStr, ";"); semicolonIndex != -1 {
			csqStr = csqStr[:semicolonIndex]
		}

		if csq, err := strconv.Atoi(csqStr); err == nil {
			return csq
		}
	}

	// Try to parse as direct number
	if csq, err := strconv.Atoi(healthCheck.Signal); err == nil {
		return csq
	}

	return -1
}

// GetBackupPowerStatus extracts the backup power status from the signal string
// Returns -1 if parsing fails
func (healthCheck *PowerBankHealthCheckResponse) GetBackupPowerStatus() int {
	if len(healthCheck.Signal) == 0 {
		return -1
	}

	// Look for "BP:" in the signal string
	bpIndex := strings.Index(healthCheck.Signal, "BP:")
	if bpIndex == -1 {
		return -1
	}

	// Extract the number after "BP:"
	bpStr := healthCheck.Signal[bpIndex+3:]
	if semicolonIndex := strings.Index(bpStr, ";"); semicolonIndex != -1 {
		bpStr = bpStr[:semicolonIndex]
	}

	if bp, err := strconv.Atoi(bpStr); err == nil {
		return bp
	}

	return -1
}

// GetSignalDescription returns a human-readable description of the signal strength
func (healthCheck *PowerBankHealthCheckResponse) GetSignalDescription() string {
	csq := healthCheck.GetCSQValue()

	switch {
	case csq == 99:
		return "No network"
	case csq >= 0 && csq <= 12:
		return "Very poor"
	case csq >= 13 && csq <= 16:
		return "Poor"
	case csq >= 17 && csq <= 20:
		return "Average"
	case csq >= 21 && csq <= 25:
		return "Good"
	case csq >= 26 && csq <= 31:
		return "Very good"
	default:
		return "Unknown"
	}
}

// GetSignalBars returns the number of signal bars (0-4) based on CSQ value
func (healthCheck *PowerBankHealthCheckResponse) GetSignalBars() int {
	csq := healthCheck.GetCSQValue()

	switch {
	case csq == 99:
		return 0 // Offline
	case csq >= 0 && csq <= 12:
		return 0
	case csq >= 13 && csq <= 16:
		return 1
	case csq >= 17 && csq <= 20:
		return 2
	case csq >= 21 && csq <= 25:
		return 3
	case csq >= 26 && csq <= 31:
		return 4
	default:
		return 0
	}
}

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
	case 0xFC:
		return constants.PowerbankStatus_PopupTargetSnNotFound
	case 0xFF:
		return constants.PowerbankStatus_PopupCommandParsingFailed
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
	case 0xFC:
		return "Target SN not found"
	case 0xFF:
		return "Command parsing failed"
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
	case 0x10:
		return "Region code error"
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
	case 0x10:
		return constants.PowerbankStatus_ReturnRegionCodeError
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

func (rt *PowerBankReturnFixResponse) GetDescription() string {
	switch rt.State {
	case 0x00:
		return "Return failed"
	case 0x01:
		return "Return successful"
	case 0x10:
		return "Region code error"
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
		return "Anti‑theft switch detection failed"
	default:
		return "Unknown error"
	}
}

func (rt *PowerBankReturnFixResponse) GetStatus() constants.PowerbankStatus {
	switch rt.State {
	case 0x00:
		return constants.PowerbankStatus_ReturnFailed
	case 0x01:
		return constants.PowerbankStatus_ReturnSuccessful
	case 0x10:
		return constants.PowerbankStatus_ReturnRegionCodeError
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

func (p *PowerBankPopupByHoleResponse) GetDescription() string {
	switch p.State {
	case 0x00:
		return "Pop-up failed"
	case 0x01:
		return "Pop-up successful"
	case 0x11:
		return "Serial timeout"
	case 0x12:
		return "SN readable but power bank not ejected"
	case 0x87:
		return "Timestamp retrieval failed"
	case 0x88:
		return "TTL exceeded"
	case 0xFF:
		return "Command parsing failed"
	default:
		return "Unknown error"
	}
}

func (p *PowerBankPopupByHoleResponse) GetStatus() constants.PowerbankStatus {
	switch p.State {
	case 0x00:
		return constants.PowerbankStatus_PopupFailed
	case 0x01:
		return constants.PowerbankStatus_PopupSuccessful
	case 0x11:
		return constants.PowerbankStatus_FailedToObtainSn
	case 0x12:
		return constants.PowerbankStatus_PopupCompleteMotorHomeSnReadable
	case 0x87:
		return constants.PowerbankStatus_PopupTimestampRetrievalFailed
	case 0x88:
		return constants.PowerbankStatus_PopupTTLExceeded
	case 0xFF:
		return constants.PowerbankStatus_PopupCommandParsingFailed
	}
	return constants.PowerbankStatus_UnknownError
}
