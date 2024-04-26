package message

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"time"
)

const (
	networkMagic    = 0xDAB5BFFA // Magic value for regtest
	protocolVersion = 70015
)

type Command string

var (
	Version = Command("version")
	Verack  = Command("verack")
)

type Message struct {
	Header  Header
	Payload []byte
}

type Header struct {
	Magic    uint32
	Command  [12]byte
	Length   uint32
	Checksum [4]byte
}

type VersionPayload struct {
	Version     int32
	Services    uint64
	Timestamp   int64
	AddrRecv    [26]byte // Dummy addresses
	AddrFrom    [26]byte // Dummy addresses
	Nonce       uint64
	UserAgent   []byte
	StartHeight int32
	Relay       bool
}

func getDefaultVersionPayload() VersionPayload {
	return VersionPayload{
		Version:     protocolVersion,
		Services:    1,
		Timestamp:   time.Now().Unix(),
		UserAgent:   []byte("/TestNode:0.0.1/"),
		StartHeight: 0,
		Relay:       false,
	}
}

func New(command Command) (*Message, error) {
	payload := new(bytes.Buffer)

	if command == Version {
		msg := getDefaultVersionPayload()

		if err := binary.Write(payload, binary.LittleEndian, msg.Version); err != nil {
			return nil, err
		}
		if err := binary.Write(payload, binary.LittleEndian, msg.Services); err != nil {
			return nil, err
		}
		if err := binary.Write(payload, binary.LittleEndian, msg.Timestamp); err != nil {
			return nil, err
		}
		if _, err := payload.Write(msg.AddrRecv[:]); err != nil {
			return nil, err
		}
		if _, err := payload.Write(msg.AddrFrom[:]); err != nil {
			return nil, err
		}
		if err := binary.Write(payload, binary.LittleEndian, msg.Nonce); err != nil {
			return nil, err
		}
		if err := binary.Write(payload, binary.LittleEndian, uint8(len(msg.UserAgent))); err != nil {
			return nil, err
		}
		if _, err := payload.Write(msg.UserAgent); err != nil {
			return nil, err
		}
		if err := binary.Write(payload, binary.LittleEndian, msg.StartHeight); err != nil {
			return nil, err
		}
		if err := binary.Write(payload, binary.LittleEndian, msg.Relay); err != nil {
			return nil, err
		}
	}

	p := payload.Bytes()

	h := Header{}
	h.Magic = networkMagic
	copy(h.Command[:], Version)

	h.Length = uint32(len(p))
	h.Checksum = calculateChecksum(p)

	return &Message{
		Header:  h,
		Payload: p,
	}, nil
}

func (m *Message) Serialize() ([]byte, error) {
	var buffer bytes.Buffer
	if err := binary.Write(&buffer, binary.LittleEndian, m.Header); err != nil {
		return nil, err
	}
	buffer.Write(m.Payload)

	return buffer.Bytes(), nil
}

func calculateChecksum(payload []byte) [4]byte {
	hasher := sha256.New()
	hasher.Write(payload)
	firstHash := hasher.Sum(nil)

	hasher.Reset()
	hasher.Write(firstHash)
	secondHash := hasher.Sum(nil)

	var checksum [4]byte
	copy(checksum[:], secondHash[:4])
	return checksum
}
