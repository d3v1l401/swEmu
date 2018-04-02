package server

import (
	"fmt"
	"strings"

	"./session"
)

func Adm_OnNewClient(s *Server, c *session.Instance) {
	fmt.Printf("New client %s (%s) @ %s\n", c.GetID(), c.GetIP(), s.GetName())
	c.AccessData()["Auth"] = "no"
}

func Adm_OnClientDeath(s *Server, c *session.Instance) {
	//fmt.Printf("Client %s (%s) @ %s being removed\n", c.GetID(), c.GetIP(), "server#1")
}

func Adm_BeforePacketSend(client *session.Instance, buffer []byte) {
	fmt.Printf("Packet being sent to %s...\n", client.GetID())
}

func Adm_OnTick(c *session.Instance) {
	//fmt.Printf("Client %s (%s) @ %s is dead\n", c.GetID(), c.GetIP(), "server#1")
}

func Adm_OnHeartbeat(s *Server) {
	//fmt.Printf("Online clients: %d.\n", s.OnlineClients())
}

func Adm_OnPacketRecv(client *session.Instance, buffer []byte) bool {
	if buffer != nil {
		fmt.Printf("Request received from %s: %s\n", client.GetID(), string(buffer))

		if strings.Compare(string(buffer), "stop") == 0 {
			fmt.Printf("Shutdown everything.\n")
		}

		return true
	}
	return false
}

func Adm_OnPacketSent(client *session.Instance) {
	fmt.Printf("Packet sent to %s\n", client.GetID())
}

func Adm_OnDestroy(client *session.Instance) {
	fmt.Printf("Client %s destroyed\n", client.GetID())
}
