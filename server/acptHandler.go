package server

import (
	"encoding/hex"
	"fmt"

	s "./serializer"
	"./session"
)

var counter int = 0

func OpSwitch(r *s.PacketReader) []byte {
	r.Dump()
	fmt.Println(r.Class, r.Type)

	switch r.Class {
	case CLASS_UTILS:

		switch r.Type {
		case TYPE_PING:
			fmt.Println("Ping packet")
			break
		}

		break
	case CLASS_LOGINCLASS:
		switch r.Type {
		case TYPE_AUTHCLASSICREQ:
			r.Skip(2)

			username := r.ReadStringUTF16()
			password := r.ReadStringUTF16()
			mac := r.ReadStringUTF16()
			constant_ver := r.ReadInt()
			fmt.Println(username, password, mac, constant_ver)

			pW := s.NewWriter()
			pW.Write("\x02\x02\x00\x00\x00\x00\x00\x00\x01\x00\x01\x00\x00\x00\x00\x00\x00\x01\x00\x01\x00\x00\x00\x00\x00\x00\x01\x00\x01\x00\x00\x00\x00\x00\x00\x01")
			return pW.Finalize()

			break
		}
		break
	case CLASS_AUTH:
		switch r.Type {
		// w/ token
		case TYPE_AUTHREQ:

			r.Skip(2)
			authToken := r.ReadStringUTF16()
			moboMac := r.ReadString()
			fmt.Println("Token:", authToken, "\nMAC:", moboMac)

			pW := s.NewWriter()
			pW.Write("\x02\xd6\x05\x00\x00\x01")
			pW.Write(moboMac)
			pW.Write("\x00\x01\x00\x00\x00\x00\x00\x00\x01")
			pW.WriteStringUTF16("divinerain0373")
			// 4 bytes unknown, seems random
			pW.Write("\x4a\x22\x28\x00")
			pW.Write("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00")
			return pW.Finalize()
		}
	case CLASS_SERVERS:
		switch r.Type {
		case TYPE_SERVREQ:
			fmt.Println("Server list request")
		}
	}
	return nil
}

func OnNewClient(s *Server, c *session.Instance) {
	fmt.Printf("New client %s (%s) @ %s\n", c.GetID(), c.GetIP(), s.GetName())
}

func OnClientDeath(s *Server, c *session.Instance) {
	//fmt.Printf("Client %s (%s) @ %s being removed\n", c.GetID(), c.GetIP(), "server#1")
}

func BeforePacketSend(client *session.Instance, buffer []byte) {
	fmt.Printf("Packet being sent to %s...\n", client.GetID())
}

func OnTick(c *session.Instance) {
	//fmt.Printf("Client %s (%s) @ %s is dead\n", c.GetID(), c.GetIP(), "server#1")
}

func OnHeartbeat(s *Server) {
	//fmt.Printf("Online clients: %d.\n", s.OnlineClients())
}

func OnPacketRecv(client *session.Instance, buffer []byte) bool {
	if buffer != nil {
		fmt.Printf("Request received from %s.\n", client.GetID())

		if ack := OpSwitch(s.NewReader(buffer)); ack != nil {
			client.Send(ack)
			fmt.Printf("SENT:\n%s\n", hex.Dump(ack))
		}

		return true
	}
	return false
}

func OnPacketSent(client *session.Instance) {
	fmt.Printf("Packet sent to %s\n", client.GetID())
}

func OnDestroy(client *session.Instance) {
	fmt.Printf("Client %s destroyed\n", client.GetID())
}
