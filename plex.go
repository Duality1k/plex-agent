package main

import (
	"fmt"
	"log"
	superplayer "plex-god/api/player"
	"strings"
)

type Server struct {
	CityId  string
	Players superplayer.Players
	Owners  *superplayer.Owners
	Dev     bool
}

func (server *Server) CalculateOwners() {

	//newOwners = InitOwners()
	admin := server.Players[0]

	newOwners := &superplayer.Owners{
		Airplane:    admin,
		Blimp:       admin,
		Boat:        server.Players.PoorestPlayerId("money"),
		TruckWood:   server.Players.PoorestPlayerId("wood"),
		TruckCement: server.Players.PoorestPlayerId("cement"),
		TruckSteel:  server.Players.PoorestPlayerId("steel"),
	}

	messages := make(map[*superplayer.Player][]string)
	changes := false

	currentOwnersList := superplayer.Players{
		server.Owners.Airplane,
		server.Owners.Blimp,
		server.Owners.Boat,
		server.Owners.TruckWood,
		server.Owners.TruckSteel,
	}

	newOwnersList := superplayer.Players{
		server.Owners.Airplane,
		server.Owners.Blimp,
		server.Owners.Boat,
		server.Owners.TruckWood,
		server.Owners.TruckSteel,
	}

	for i, owner := range currentOwnersList {
		if owner == newOwnersList[i] {
			changes = true
		}
	}

	if changes {
		if newOwners.Airplane != nil {
			messages[newOwners.Airplane] = append(messages[newOwners.Airplane], "✈️")
		}

		if newOwners.Blimp != nil {
			messages[newOwners.Blimp] = append(messages[newOwners.Blimp], "🎈")
		}
		if newOwners.Boat != nil {
			messages[newOwners.Boat] = append(messages[newOwners.Boat], "💰")
		}
		if newOwners.TruckWood != nil {
			messages[newOwners.TruckWood] = append(messages[newOwners.TruckWood], "🪵")
		}
		if newOwners.TruckCement != nil {
			messages[newOwners.TruckCement] = append(messages[newOwners.TruckCement], "🧱")
		}
		if newOwners.TruckSteel != nil {
			messages[newOwners.TruckSteel] = append(messages[newOwners.TruckSteel], "🔩")
		}

		for player, items := range messages {
			if player != nil {
				fmt.Printf("@%s → %s\n💰%d, 🪵%d, 🧱%d, 🔩%d\n",
					player.ResidentId, //player.Stats.Username,
					strings.Join(items, ", "),
					player.WarehouseUsage.Money,
					player.WarehouseUsage.Wood,
					player.WarehouseUsage.Cement,
					player.WarehouseUsage.Steel,
				)
			} else {
				fmt.Println("nil player")
			}
		}

	}

	server.Owners = newOwners
}

func main() {

	_ = &Server{
		Dev:    false,
		Owners: superplayer.InitOwners(),
		Players: superplayer.Players{
			// gabri
			superplayer.NewPlayer("AAC34261-A00A-4A93-880D-7AE34C677520", "5616078598701056", "5319377040179200"),
			// giagia
			superplayer.NewPlayer("722E8EF8-5B8C-44E3-ACD0-8B7C1315E216", "5358963736969216", "5319377040179200"),
			// rebe
			superplayer.NewPlayer("91ABF76F-34EB-46BF-BB48-BB22BDC47F99", "5237005657571328", "5319377040179200"),
		},
	}

	_ = &Server{
		Dev:    false,
		Owners: superplayer.InitOwners(),
		Players: superplayer.Players{
			// gabri
			superplayer.NewPlayer("AAC34261-A00A-4A93-880D-7AE34C677520", "6430000268443648", "6392074448928768"),
		},
	}

	go func() {
		err := StartWebsocket(superplayer.NewPlayer("AAC34261-A00A-4A93-880D-7AE34C677520", "5616078598701056", "5319377040179200"), nil)
		if err != nil {
			log.Fatalf("Dial error: %v", err)
		}
	}()

	go func() {
		err := StartWebsocket(superplayer.NewPlayer("AAC34261-A00A-4A93-880D-7AE34C677520", "6430000268443648", "6392074448928768"), nil)
		if err != nil {
			log.Fatalf("Dial error: %v", err)
		}
	}()

	/*
		go func() {
			err := StartWebsocket(superplayer.NewPlayer("AAC34261-A00A-4A93-880D-7AE34C677520", "5806037016248320", "5009735852490752"), nil)
			if err != nil {
				log.Fatalf("Dial error: %v", err)
			}
		}()

	*/
	fmt.Println("Connected to WebSocket")

	//privateServer.StartBot()
	//server2.StartBot()
	select {}
}
