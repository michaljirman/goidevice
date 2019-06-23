package main

import (
	"fmt"

	"github.com/michaljirman/goidevice/usbmuxd"
)

func main() {
	client, _ := usbmuxd.NewDefaultClient()
	client.ListDevices()
	//client.ListDevices()
	//client.ListDevices()
	fmt.Println("finished")
}

//	payload := []byte(`<?xml version="1.0" encoding="UTF-8"?>
//<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
//<plist version="1.0">
//<dict>
//	<key>BundleID</key>
//	<string>com.apple.configurator.xpc.DeviceService</string>
//	<key>ClientVersionString</key>
//	<string>usbmuxd-423.258.2</string>
//	<key>MessageType</key>
//	<string>ListDevices</string>
//	<key>ProgName</key>
//	<string>com.apple.configurator.xpc.DeviceService</string>
//</dict>
//</plist>`)
//	payload := []byte(`<?xml version="1.0" encoding="UTF-8"?>
//<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
//<plist version="1.0">
//<dict>
//	<key>BundleID</key>
//	<string>com.apple.configurator.xpc.DeviceService</string>
//	<key>ClientVersionString</key>
//	<string>usbmuxd-423.258.2</string>
//	<key>DeviceID</key>
//	<integer>5</integer>
//	<key>MessageType</key>
//	<string>Connect</string>
//	<key>PortNumber</key>
//	<integer>32498</integer>
//	<key>ProgName</key>
//	<string>com.apple.configurator.xpc.DeviceService</string>
//</dict>
//</plist>`)

//data := make([]byte, 16+len(payload))
//binary.LittleEndian.PutUint32(data[0:4], 1)
//binary.LittleEndian.PutUint32(data[4:8], 2) // 8 packet type request, 2 packet type connect
//binary.LittleEndian.PutUint32(data[8:12], 16)
//copy(data[12:], payload)
