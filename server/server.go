package server

import (
	"container/list"
	"net"
	"strconv"
	"strings"
	"time"

	"../common"
	"./session"
)

const (
	PROTOCOL_TCP = "tcp"
	PROTOCOL_UDP = "udp"
)

// Server is the instance of the server being created.
type Server struct {
	Name     string
	IP       string
	Port     int
	Protocol string
	Sessions *list.List

	lastTick int64

	ClientCBs *session.CallbackSet

	OnNewClient   func(s *Server, c *session.Instance)
	OnClientDeath func(s *Server, c *session.Instance)
	OnHeartbeat   func(s *Server)
}

func (s *Server) onHeartbeatEvent(callback func(s *Server)) {
	s.OnHeartbeat = callback
}

func (s *Server) onNewClientEvent(callback func(s *Server, c *session.Instance)) {
	s.OnNewClient = callback
}

func (s *Server) onClientDeathEvent(callback func(s *Server, c *session.Instance)) {
	s.OnClientDeath = callback
}

// SetClientCallbacks sets the default callbacks for new clients.
func (s *Server) SetClientCallbacks(cbSet *session.CallbackSet) {
	s.ClientCBs = cbSet
}

// Initialize the server Instance.
// 	bindIP: IP to bind to.
//  bindPort: port to bind to.
//  protocol: desider protocol between udp and tcp, the enum usage is encouraged.
//	returns: wheter it fails to initialize or not.
func (s *Server) Initialize(name string, bindIP string, bindPort int, protocol string) bool {
	s.Name = name
	s.IP = bindIP
	s.Port = bindPort

	switch protocol {
	case PROTOCOL_TCP:
		s.Protocol = protocol

		break
	case PROTOCOL_UDP:
	default:

		return false
	}
	s.Sessions = list.New()
	s.Sessions.Init()

	return true
}

func (s *Server) hbLoop() {
	for {
		time.Sleep(1 * time.Second)
		s.cleanDeadClients()
		s.callClientTicks()
		if s.OnHeartbeat != nil {
			s.OnHeartbeat(s)
		}
	}
}

// Start the server
func (s *Server) Start() {
	go s.Listen()
}

// GetName returns the name of the server
func (s *Server) GetName() string {
	return s.Name
}

// Listen start the server's listening
func (s *Server) Listen() {
	go s.hbLoop()
	paramString := s.IP + ":" + strconv.Itoa(s.Port)
	listener, err := net.Listen(s.Protocol, paramString)
	if common.ErrorCheck("Server Listen", err) {
		common.Print("Server listen failed: "+err.Error(), common.MTYPE_ERROR)
	}
	defer listener.Close()

	for {

		connection, _ := listener.Accept()
		newSession := session.Initialize(connection)
		newSession.Callbacks = s.ClientCBs

		s.Sessions.PushFront(newSession)

		if s.OnNewClient != nil {
			s.OnNewClient(s, newSession)
		}

		go newSession.Listener()

	}
}

// OnlineClients returns the amount of active sessions.
func (s *Server) OnlineClients() int {
	return s.Sessions.Len()
}

// GetClientByUID retreives client's instance by its UID.
//  uid: the client's expected uid.
//	returns: the instance.
func (s *Server) GetClientByUID(uid string) *session.Instance {
	for e := s.Sessions.Front(); e != nil; e.Next() {
		var sessionSlave *session.Instance = e.Value.(*session.Instance)
		if strings.Compare(sessionSlave.GetID(), uid) == 0 {
			return sessionSlave
		}
	}

	return nil
}

func (s *Server) callClientTicks() {
	for e := s.Sessions.Front(); e != nil; e.Next() {
		var sessionSlave *session.Instance = e.Value.(*session.Instance)
		sessionSlave.TickClient()
	}
}

func (s *Server) cleanDeadClients() {
	var next *list.Element
	for e := s.Sessions.Front(); e != nil; e = next {
		next = e.Next()
		var sessionSlave *session.Instance = e.Value.(*session.Instance)
		if strings.Compare(sessionSlave.GetID(), "") == 0 {
			s.OnClientDeath(s, sessionSlave)
			s.Sessions.Remove(e)
		}
	}
}

// GetClientByIP retreives client's instance by its IP.
//  ip: the client's expected ip.
//	returns: the instance.
func (s *Server) GetClientByIP(ip string) *session.Instance {
	for e := s.Sessions.Front(); e != nil; e.Next() {
		var sessionSlave *session.Instance = e.Value.(*session.Instance)
		if strings.Compare(sessionSlave.GetIP(), ip) == 0 {
			return sessionSlave
		}
	}

	return nil
}
