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
	CreateUserResponse struct {
		UserID string `json:"user_id"`
	}
)
