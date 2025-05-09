package powerbankSdk

const (
	CMD_POPUP  = 0x31
	CMD_RETURN = 0x40
	CMD_CHECK  = "check"

	TOPIC_SUBSCRIBE    = "/powerbank/+/user/update"
	TOPIC_PUBLISH      = "/powerbank/%s/user/get"
	TOPIC_HEALTH_CHECK = "/powerbank/%s/user/heart"

	PUBLISH_TYPE_POPUP = "popup_sn"
	PUBLISH_TYPE_CHECK = "check"
)

type (
	PowerBankResponse struct {
		Head         byte // Byte[0] - Header code - ( head )
		Length       int  // Byte[1-2] - Packet length - ( length )
		Cmd          byte // Byte[3] - Command name - ( cmd ) - 0x31 Popup | 0x40 Return
		ControlIndex int  // Byte[4] - Movement board address - ( controlIndex )
	}
)
