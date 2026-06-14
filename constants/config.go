package constants

const ()

type TOPIC string

const (
	TOPIC_SUBSCRIBE    TOPIC = "/powerbank/+/user/update"
	TOPIC_PUBLISH      TOPIC = "/powerbank/%s/user/get"
	TOPIC_HEALTH_CHECK TOPIC = "/powerbank/%s/user/heart"
)

type PUBLISH_TYPE string

// Volinks Powerbank Protocol V1 command tags. Each constant is documented at:
// https://docs.volinks.com/powerbank-protocol-v1/en/guide/
const (
	// Pop-up By SN — https://docs.volinks.com/powerbank-protocol-v1/en/guide/protocol-popupsn.html
	PUBLISH_TYPE_POPUP PUBLISH_TYPE = "popup_sn"
	// Pop-up By Hole — https://docs.volinks.com/powerbank-protocol-v1/en/guide/protocol-popup.html
	PUBLISH_TYPE_POPUP_BY_HOLE PUBLISH_TYPE = "popup"
	// Cabinet check — https://docs.volinks.com/powerbank-protocol-v1/en/guide/protocol-check.html
	PUBLISH_TYPE_CHECK PUBLISH_TYPE = "check"
	// Periodic upload trigger — https://docs.volinks.com/powerbank-protocol-v1/en/guide/protocol-upload.html
	PUBLISH_TYPE_UPLOAD PUBLISH_TYPE = "upload_all"
	// Standard return response tag — https://docs.volinks.com/powerbank-protocol-v1/en/guide/protocol-return.html
	PUBLISH_TYPE_RETURN PUBLISH_TYPE = "return"
	// Return-Fix self-test response tag — https://docs.volinks.com/powerbank-protocol-v1/en/guide/protocol-return-fix.html
	PUBLISH_TYPE_RETURN_FIX PUBLISH_TYPE = "return_fix"
	// Heart response tag — https://docs.volinks.com/powerbank-protocol-v1/en/guide/protocol-heart.html
	PUBLISH_TYPE_HEALTH_CHECK PUBLISH_TYPE = "health_check"
	// Trigger ad refresh — https://docs.volinks.com/powerbank-protocol-v1/en/guide/protocol-advert.html
	PUBLISH_TYPE_LOAD_AD PUBLISH_TYPE = "load_ad"
	// Reboot the cabinet (resets main control board) — https://docs.volinks.com/powerbank-protocol-v1/en/guide/protocol-reboot.html
	PUBLISH_TYPE_REBOOT PUBLISH_TYPE = "reboot"
)

type CabinetSignal string

const (
	CabinetSignal_Weak        CabinetSignal = "weak"
	CabinetSignal_Normal      CabinetSignal = "normal"
	CabinetSignal_Better      CabinetSignal = "better"
	CabinetSignal_NoNetwork   CabinetSignal = "no-network"
	CabinetSignal_NoPower     CabinetSignal = "no-power"
	CabinetSignal_Charging    CabinetSignal = "charging"
	CabinetSignal_Discharging CabinetSignal = "discharging"
)

type PowerbankStatus string

const (
	// Hole
	PowerbankStatus_NoPowerSupply                        PowerbankStatus = "no-power-supply"
	PowerbankStatus_Normal                               PowerbankStatus = "normal"
	PowerbankStatus_ChargingAbnormality                  PowerbankStatus = "charging-abnormality"
	PowerbankStatus_CommunicationException               PowerbankStatus = "communication-exception"
	PowerbankStatus_KabaoDamaged                         PowerbankStatus = "kabao-damaged"
	PowerbankStatus_KeyForceRelease                      PowerbankStatus = "key-force-release"
	PowerbankStatus_SolenoidValveNotReturn               PowerbankStatus = "solenoid-valve-not-return"
	PowerbankStatus_Reserved                             PowerbankStatus = "reserved"
	PowerbankStatus_AntiTheftProtocolCommunicationFailed PowerbankStatus = "anti-theft-protocol-communication-failed"
	PowerbankStatus_TypecShortCircuit                    PowerbankStatus = "typec-short-circuit"
	PowerbankStatus_ReturnFailedBatteryDoesNotPopOut     PowerbankStatus = "return-failed-battery-does-not-pop-out"

	// Popup
	PowerbankStatus_PopupFailed                               PowerbankStatus = "popup-failed"
	PowerbankStatus_PopupSuccessful                           PowerbankStatus = "popup-successful"
	PowerbankStatus_PowerSupplyChargingAbnormally             PowerbankStatus = "power-supply-charging-abnormally"
	PowerbankStatus_CommunicationAbnormalityFirstReturnFailed PowerbankStatus = "communication-abnormality-first-return-failed"
	PowerbankStatus_SlotCannotPopOut                          PowerbankStatus = "slot-cannot-pop-out"
	PowerbankStatus_SlotForciblyReleased                      PowerbankStatus = "slot-forcibly-released"
	PowerbankStatus_SolenoidNotReturned                       PowerbankStatus = "solenoid-not-returned"
	PowerbankStatus_AntiTheftCommFailed                       PowerbankStatus = "anti-theft-comm-failed"
	PowerbankStatus_FailedToObtainSn                          PowerbankStatus = "failed-to-obtain-sn"
	PowerbankStatus_PopupCompleteMotorHomeSnReadable          PowerbankStatus = "popup-complete-motor-home-sn-readable"
	PowerbankStatus_FailedToObtainTraceback                   PowerbankStatus = "failed-to-obtain-traceback"
	PowerbankStatus_BatteryLockCommandFailed                  PowerbankStatus = "battery-lock-command-failed"
	PowerbankStatus_SnAcquisitionAndMotorFailed               PowerbankStatus = "sn-acquisition-and-motor-failed"
	PowerbankStatus_InfoAcquisitionAndMotorFailed             PowerbankStatus = "info-acquisition-and-motor-failed"
	PowerbankStatus_BatteryLockAndMotorFailed                 PowerbankStatus = "battery-lock-and-motor-failed"
	PowerbankStatus_AntiTheftSwitchDetectionFailed            PowerbankStatus = "anti-theft-switch-detection-failed"
	PowerbankStatus_PopupTargetSnNotFound                     PowerbankStatus = "popup-target-sn-not-found"
	PowerbankStatus_PopupCommandParsingFailed                 PowerbankStatus = "popup-command-parsing-failed"

	PowerbankStatus_UnknownError PowerbankStatus = "unknown-error"

	PowerbankStatus_ReturnFailed                                                           PowerbankStatus = "return-failed"
	PowerbankStatus_ReturnSuccessful                                                       PowerbankStatus = "return-successful"
	PowerbankStatus_FailedToObtainSnAndMotorFailed                                         PowerbankStatus = "failed-to-obtain-sn-and-motor-failed"
	PowerbankStatus_FailedToObtainVoltageTemperatureOrOtherInformation                     PowerbankStatus = "failed-to-obtain-voltage-temperature-or-other-information"
	PowerbankStatus_FailedToObtainSoftwareAndHardwareVersionInformation                    PowerbankStatus = "failed-to-obtain-software-and-hardware-version-information"
	PowerbankStatus_FailedToObtainVoltageTemperatureOrOtherInformationAndMotorFailed       PowerbankStatus = "failed-to-obtain-voltage-temperature-or-other-information-and-motor-failed"
	PowerbankStatus_BatteryLockCommandFailedAndMotorFailed                                 PowerbankStatus = "battery-lock-command-failed-and-motor-failed"
	PowerbankStatus_FailedToObtainSnAndMotorActionFailed                                   PowerbankStatus = "failed-to-obtain-sn-and-motor-action-failed"
	PowerbankStatus_FailedToObtainVoltageTemperatureOrOtherInformationAndMotorActionFailed PowerbankStatus = "failed-to-obtain-voltage-temperature-or-other-information-and-motor-action-failed"
	PowerbankStatus_BatteryLockCommandFailedAndMotorActionFailed                           PowerbankStatus = "battery-lock-command-failed-and-motor-action-failed"
	PowerbankStatus_AntiTheftSwitchDetectionFailedAndMotorActionFailed                     PowerbankStatus = "anti-theft-switch-detection-failed-and-motor-action-failed"
	PowerbankStatus_ReturnRegionCodeError                                                  PowerbankStatus = "return-region-code-error"
	PowerbankStatus_PopupTimestampRetrievalFailed                                          PowerbankStatus = "popup-timestamp-retrieval-failed"
	PowerbankStatus_PopupTTLExceeded                                                       PowerbankStatus = "popup-ttl-exceeded"
	PowerbankStatus_PopupSerialTimeout                                                     PowerbankStatus = "popup-serial-timeout"
	PowerbankStatus_PopupBankUnpoppedSnReadable                                            PowerbankStatus = "popup-bank-unpopped-sn-readable"
	PowerbankStatus_PopupNoMatchingBattery                                                 PowerbankStatus = "popup-no-matching-battery"
	PowerbankStatus_PopupAddTaskFailed                                                     PowerbankStatus = "popup-add-task-failed"
	PowerbankStatus_PopupPreviousRentalIncomplete                                          PowerbankStatus = "popup-previous-rental-incomplete"
)
