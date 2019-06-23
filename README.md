# go-idevice library

## List attached devices
```go
client, _ := usbmuxd.NewDefaultClient()
client.ListDevices()
```