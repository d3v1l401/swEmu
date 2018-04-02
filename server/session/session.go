package session

import (
	"net"
	"time"

	"../ssettings"
	// Credits to Maxim Bublis
	"github.com/satori/go.uuid"
)

// Statistics contains general statistics about the client, it can be disabled on initialization.
type Statistics struct {
	Sent            int
	Received        int
	SentPackets     int
	ReceivedPackets int
	Alive           int64
}

// Instance contains callbacks and vital information of the actual connected client.
type Instance struct {
	IClient    net.Conn
	UID        string
	Statistics Statistics
	ClientData map[string]string

	// These can be null, if not needed.
	Callbacks *CallbackSet
}

// CallbackSet exists because we need to organize callbacks better, since there are too many.
type CallbackSet struct {
	BeforePacketSend func(client *Instance, buffer []byte)
	OnPacketReceive  func(client *Instance, buffer []byte) bool
	OnPacketSent     func(client *Instance)
	OnDestroy        func(client *Instance)
	OnClientTick     func(client *Instance)
}

func (s *Instance) onClientTickEvent(callback func(client *Instance)) {
	s.Callbacks.OnClientTick = callback
}

func (s *Instance) onClientDestroyEvent(callback func(client *Instance)) {
	s.Callbacks.OnDestroy = callback
}

func (s *Instance) onPacketSentEvent(callback func(client *Instance)) {
	s.Callbacks.OnPacketSent = callback
}

func (s *Instance) onPacketReceiveEvent(callback func(client *Instance, buffer []byte) bool) {
	s.Callbacks.OnPacketReceive = callback
}

func (s *Instance) onPacketBeforeSendEvent(callback func(client *Instance, buffer []byte)) {
	s.Callbacks.BeforePacketSend = callback
}

// Generates a unique ID for the client, respecting RFC4122 implementation & only randomic.
func (s *Instance) assignID() {
	s.UID = uuid.NewV4().String()
}

// Listener is the client's listening.
func (s *Instance) Listener() {

	for {
		buffer := s.Recv()
		if buffer != nil {
			if s.Callbacks != nil && s.Callbacks.OnPacketReceive != nil {
				s.Callbacks.OnPacketReceive(s, buffer)
			}
		}

		if !ssettings.GServerSettings.KeptAlive() {
			s.Destroy()
			return
		}
	}
}

// GetID returns the UID.
func (s *Instance) GetID() string {
	return s.UID
}

// GetStatistics returns the statistics of the client.
func (s *Instance) GetStatistics() Statistics {
	return s.Statistics
}

// Destroy simply closes the connection with the client and ensures vital parameters are removed for safety.
func (s *Instance) Destroy() {
	if s.Callbacks != nil && s.Callbacks.OnDestroy != nil {
		s.Callbacks.OnDestroy(s)
	}

	s.UID = ""
	s.IClient.Close()
}

// Send the buffer to the client.
//  buffer: actual buffer.
//  returns: failure or success.
func (s *Instance) Send(buffer []byte) bool {

	if s.Callbacks != nil && s.Callbacks.BeforePacketSend != nil {
		s.Callbacks.BeforePacketSend(s, buffer)
	}

	_, err := s.IClient.Write(buffer)
	if err != nil {
		return false
	}

	if s.Callbacks != nil && s.Callbacks.OnPacketSent != nil {
		s.Callbacks.OnPacketSent(s)
	}

	s.Statistics.SentPackets++
	s.Statistics.Sent += len(buffer)
	s.Statistics.Alive = time.Now().UnixNano()

	return true
}

// Recv the buffer to the client.
//  returns: the buffer, or null if nothing.
func (s *Instance) Recv() []byte {
	buffer := make([]byte, ssettings.GServerSettings.GetWindowSize())
	size, err := s.IClient.Read(buffer)
	if err != nil {
		return nil
	}

	if size < ssettings.GServerSettings.GetWindowSize() {
		s.Statistics.ReceivedPackets++
		s.Statistics.Received += size
		s.Statistics.Alive = time.Now().UnixNano()

		nBuffer := make([]byte, size)
		copy(nBuffer, buffer)

		return nBuffer
	}

	return nil
}

// TickClient have to be called every heartbeat for statistics update and checks.
// An heartbeat is just a packet that can be sent every amount of seconds to check wheter the client is still online, very common in online games.
func (s *Instance) TickClient() {
	if (time.Now().UnixNano() - s.Statistics.Alive) > ssettings.GServerSettings.GetClientTimeout() {
		if s.Callbacks != nil && s.Callbacks.OnClientTick != nil {
			s.Callbacks.OnClientTick(s)
		}
		return
	}
}

// GetIP returns the client's IP.
func (s *Instance) GetIP() string {
	return s.IClient.RemoteAddr().String()
}

// Initialize statically a structure for parsing
//  client: client instance of the socket.
//  returns: the Instance structure or null.
func Initialize(client net.Conn) *Instance {
	nInstance := &Instance{}
	nInstance.Initialize(client, nil)
	nInstance.SetData(make(map[string]string))
	return nInstance
}

// SetData sets custom data of a client.
func (s *Instance) SetData(data map[string]string) {
	s.ClientData = data
}

// AccessData returns custom data of a client.
func (s *Instance) AccessData() map[string]string {
	return s.ClientData
}

// Initialize the client instance and calls respective callbacks.
//  client: client instance.
//  returns: wheter creation is successful or not.
func (s *Instance) Initialize(client net.Conn, dataInterface map[string]string) bool {
	s.IClient = client
	s.assignID()
	s.Statistics.Alive = time.Now().UnixNano()
	s.ClientData = dataInterface

	return true
}
