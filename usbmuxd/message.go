package usbmuxd

const (
	MessageTypeListDevices = "ListDevices"
	MessageTypeResult      = "Result"
)

type Message interface {
	GetMessageType() string
}

var DefaultMessage = BaseMessage{
	BundleID:            "com.apple.configurator.xpc.DeviceService",
	ClientVersionString: "usbmuxd-423.258.2",
	ProgName:            "com.apple.configurator.xpc.DeviceService",
}

type BaseMessage struct {
	BundleID            string `plist:"BundleID"`
	MessageType         string `plist:"MessageType"`
	ClientVersionString string `plist:"ClientVersionString"`
	ProgName            string `plist:"ProgName"`
}

type ResultMessage struct {
	Number int `plist:"Number"`
}

type ListDevicesMessage struct {
	BaseMessage
}

func NewListDevicesMessage() ListDevicesMessage {
	msg := ListDevicesMessage{BaseMessage: DefaultMessage}
	msg.MessageType = MessageTypeListDevices
	return msg
}

func (msg ListDevicesMessage) GetMessageType() string {
	return msg.MessageType
}
