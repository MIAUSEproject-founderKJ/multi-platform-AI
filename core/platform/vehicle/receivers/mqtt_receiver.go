//core/platform/vehicle/receivers/mqtt_receiver.go

type MQTTReceiver struct {
    client mqtt.Client
    output chan<- []byte
}

func (r *MQTTReceiver) Start() error {
    return r.client.Subscribe("#", 0, func(client mqtt.Client, msg mqtt.Message) {
        r.output <- msg.Payload()
    })
}