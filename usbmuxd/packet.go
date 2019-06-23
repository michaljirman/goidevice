package usbmuxd

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type PacketType uint32

const (
	PacketTypeConnect PacketType = 2
	PacketTypeRequest PacketType = 8
)

type Packet interface {
	MarshalPacket() ([]byte, error)
	UnmarshalPacket(data []byte) error
	String() string
	GetTag() tag
}

type RawPacket struct {
	Version uint32
	Type    uint32
	Tag     tag
	Payload []byte
}

func RawPacketFromData(data []byte) (*RawPacket, error) {
	packet := &RawPacket{}
	err := packet.UnmarshalPacket(data)
	if err != nil {
		return nil, err
	}
	return packet, nil
}

func (packet *RawPacket) UnmarshalPacket(data []byte) error {
	packet.Version = binary.LittleEndian.Uint32(data[4:8])
	packet.Type = binary.LittleEndian.Uint32(data[8:12])
	packet.Tag = tag(binary.LittleEndian.Uint32(data[12:16]))
	packet.Payload = data[16:]
	return nil
}

func (packet *RawPacket) MarshalPacket() ([]byte, error) {
	data := make([]byte, 16+len(packet.Payload))
	binary.LittleEndian.PutUint32(data[0:4], packet.Version)
	binary.LittleEndian.PutUint32(data[4:8], packet.Type)
	binary.LittleEndian.PutUint32(data[8:12], uint32(packet.Tag))
	copy(data[12:], packet.Payload)
	return data, nil
}

func (packet RawPacket) String() string {
	payloadString := strings.Replace(string(packet.Payload), "\n", "\n    ", -1)
	return fmt.Sprintf("Sent packet\n  Type: %v\n  Tag: %v\n  Payload:\n    %v\n", packet.Type, packet.Tag, payloadString)
}

func (packet *RawPacket) GetTag() tag {
	return packet.Tag
}
