package constants

const ()

type TOPIC string

const (
	TOPIC_SUBSCRIBE    TOPIC = "/powerbank/+/user/update"
	TOPIC_PUBLISH      TOPIC = "/powerbank/%s/user/get"
	TOPIC_HEALTH_CHECK TOPIC = "/powerbank/%s/user/heart"
)

type PUBLISH_TYPE string

const (
	PUBLISH_TYPE_POPUP  PUBLISH_TYPE = "popup_sn"
	PUBLISH_TYPE_CHECK  PUBLISH_TYPE = "check"
	PUBLISH_TYPE_UPLOAD PUBLISH_TYPE = "upload_all"
	PUBLISH_TYPE_RETURN PUBLISH_TYPE = "return"
)
