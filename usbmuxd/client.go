package usbmuxd

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"

	"howett.net/plist"
)

type tag uint32
type PacketHandlerFunc func(client *usbMuxdClient, data []byte, err error)

type usbMuxdClient struct {
	Conn            net.Conn
	protocolVersion uint32
	tag             tag
	mux             sync.Mutex
	replyHandlers   map[tag]PacketHandlerFunc
}

func NewDefaultClient() (*usbMuxdClient, error) {
	return &usbMuxdClient{
		protocolVersion: 1,
		replyHandlers:   make(map[tag]PacketHandlerFunc),
	}, nil
}

func (client *usbMuxdClient) CloseWithErr(err error) error {
	return client.Conn.Close()
}

func (client *usbMuxdClient) HandleAll() error {
	for {
		if client.Conn == nil {
			return io.EOF
		}
		var size uint32
		err := binary.Read(client.Conn, binary.LittleEndian, &size)
		if err != nil {
			client.CloseWithErr(err)
			return err
		}
		packetData := make([]byte, size+4)
		binary.LittleEndian.PutUint32(packetData[:4], size)
		_, err = client.Conn.Read(packetData[4:])
		if err != nil {
			client.CloseWithErr(err)
			return err
		}

		tag := tag(binary.LittleEndian.Uint32(packetData[12:16]))
		client.mux.Lock()
		defer client.mux.Unlock()

		if replyHandler := client.replyHandlers[tag]; replyHandler != nil {
			delete(client.replyHandlers, tag)
			replyHandler(client, packetData, nil)
			return nil
		}
		return nil
	}
}

func (client *usbMuxdClient) ListDevices() {
	if err := client.connect(); err != nil {
		fmt.Println(err)
	}
	err := client.SendMessage(context.Background(), NewListDevicesMessage())
	if err != nil {
		fmt.Println(err)
	}
}

func (client *usbMuxdClient) connect() error {
	conn, err := net.Dial("unix", "/var/run/usbmuxd")
	if err != nil {
		return fmt.Errorf("failed to connect to usbmuxd: %v", err)
	}
	client.Conn = conn
	go func() {
		defer func() {
			x := recover()
			if x != nil {
				conn.Close()
			}
		}()
		err := client.HandleAll()
		if err != nil {
			client.CloseWithErr(err)
		}
	}()
	return nil
}

func (client *usbMuxdClient) NextTag() tag {
	return tag(atomic.AddUint32((*uint32)(&client.tag), 1))
}

func (client *usbMuxdClient) setReplyHandler(tag tag, handler PacketHandlerFunc) {
	client.mux.Lock()
	defer client.mux.Unlock()
	client.replyHandlers[tag] = handler
}

func (client *usbMuxdClient) sendPacket(ctx context.Context, packet Packet) error {
	replyChan := make(chan struct{}, 1)
	client.setReplyHandler(packet.GetTag(), func(client *usbMuxdClient, data []byte, err error) {
		fmt.Printf("%x\n", data)
		packet, _ := RawPacketFromData(data)
		fmt.Println(packet)
		replyChan <- struct{}{}
	})
	data, err := packet.MarshalPacket()
	if err != nil {
		return err
	}

	err = binary.Write(client.Conn, binary.LittleEndian, uint32(len(data)))
	if err != nil {
		return err
	}
	_, err = client.Conn.Write(data)
	if err != nil {
		return err
	}
	<-replyChan
	return nil
}

func (client *usbMuxdClient) SendMessage(ctx context.Context, message Message) error {
	msgData, err := plist.Marshal(message, plist.XMLFormat)
	if err != nil {
		return err
	}
	packet := &RawPacket{
		Version: client.protocolVersion,
		Type:    uint32(PacketTypeRequest),
		Tag:     client.NextTag(),
		Payload: msgData,
	}
	return client.sendPacket(ctx, packet)
}
