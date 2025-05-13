package main

func main() {

	// service := powerbankSdk.NewServer(powerbankModels.ServerInput{
	// 	Host:     "103.50.205.106",
	// 	Port:     "1883",
	// 	Username: "backend",
	// 	Password: "Mongol123@",
	// 	CallbackSubscribe: func(typ constants.PUBLISH_TYPE, clientID string, msg interface{}) {
	// 		fmt.Println(typ, clientID, msg)
	// 	},

	// 	CallbackPublish: func(msg mqtt.Message) {
	// 		fmt.Println(string(msg.Payload()))
	// 	},
	// })

	// fmt.Println(service)

	// // service.Publish(powerbankModels.PublishInput{
	// // 	ClientID:    "864601068412899",
	// // 	PublishType: constants.PUBLISH_TYPE_POPUP,
	// // 	Data:        "85021618",
	// // })

	// service.Publish(powerbankModels.PublishInput{
	// 	ClientID:    "864601068412899",
	// 	PublishType: constants.PUBLISH_TYPE_CHECK,
	// })

	// // Keep the program running indefinitely
	// select {}
}
