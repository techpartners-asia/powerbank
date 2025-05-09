package powerbankModels

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
)
