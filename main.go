package main

import (
	"./common"
	"./crypto"
	"./server"
	"./server/session"
	"./server/ssettings"
)

var (
	gConf *common.Configuration
)

func main() {
	common.Print("Loading configuration...", common.MTYPE_NORMAL)

	if gConf = common.ImportConfiguration(".\\data\\config\\config.json"); gConf == nil {
		common.Print("Cannot load configuration file!", common.MTYPE_ERROR)
		return
	}

	common.Print("Loading protocol keyTable...", common.MTYPE_NORMAL)

	if err := crypto.ImportKeyTable(".\\data\\crypto\\keyTable"); err != nil {
		common.Print("Key table loading failed: "+err.Error(), common.MTYPE_ERROR)
		return
	}

	common.Print("Starting server...", common.MTYPE_NORMAL)
	ssettings.Initialize(1024, 10000000, true)

	authServer := &server.Server{}
	if authServer.Initialize("AUTH1", "127.0.0.1", 10000, server.PROTOCOL_TCP) {
		authServer.OnNewClient = server.OnNewClient
		authServer.OnClientDeath = server.OnClientDeath
		authServer.OnHeartbeat = server.OnHeartbeat
		authServer.SetClientCallbacks(
			&session.CallbackSet{
				BeforePacketSend: server.BeforePacketSend,
				OnPacketReceive:  server.OnPacketRecv,
				OnPacketSent:     server.OnPacketSent,
				OnDestroy:        server.OnDestroy,
				OnClientTick:     server.OnTick,
			})
		common.Print("Authorization Server started.", common.MTYPE_NORMAL)
		authServer.Start()

	}

	admServer := &server.Server{}
	if admServer.Initialize("ADM1", "127.0.0.1", 10001, server.PROTOCOL_TCP) {
		admServer.OnNewClient = server.Adm_OnNewClient
		admServer.OnClientDeath = server.Adm_OnClientDeath
		admServer.OnHeartbeat = server.Adm_OnHeartbeat
		admServer.SetClientCallbacks(
			&session.CallbackSet{
				BeforePacketSend: server.Adm_BeforePacketSend,
				OnPacketReceive:  server.Adm_OnPacketRecv,
				OnPacketSent:     server.Adm_OnPacketSent,
				OnDestroy:        server.Adm_OnDestroy,
				OnClientTick:     server.Adm_OnTick,
			})
		common.Print("Administration Server started.", common.MTYPE_NORMAL)
		admServer.Start()
	}

	for {

	}

	common.Print("Server stopped.", common.MTYPE_WARNING)
}
