package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	superplayer "plex-god/api/player"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

// GenerateRandomBase64 generates a random base64-encoded string
func GenerateRandomBase64(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(randomBytes), nil
}

var messageID int

func StartWebsocket(player *superplayer.Player, owners *superplayer.Owners) (err error) {
	// Define the WebSocket URL
	u := url.URL{
		Scheme:   "wss",
		Host:     "s-usc1a-nss-2059.firebaseio.com",
		Path:     "/.ws",
		RawQuery: "v=5&ns=citygame-dev&ls=Kfdbm4KYRlAdFpB3nPa325FlZHthH3m4",
	}

	// Generate a proper Sec-WebSocket-Key
	webSocketKey, err := GenerateRandomBase64(16)
	if err != nil {
		log.Fatalf("Failed to generate WebSocket key: %v", err)
		return err
	}

	// Set up the WebSocket configuration
	config, err := websocket.NewConfig(u.String(), "https://dummy-origin.com")
	if err != nil {
		log.Fatalf("Failed to create WebSocket config: %v", err)
	}
	config.Header = http.Header{
		"Accept":                   {"*/*"},
		"Accept-Encoding":          {"gzip, deflate"},
		"Sec-WebSocket-Version":    {"13"},
		"Sec-WebSocket-Key":        {webSocketKey},
		"X-Firebase-GMPID":         {"1:959368751113:ios:5ea6429a3af354f5"},
		"User-Agent":               {"Firebase/5/11.4.0_Dec 2 2024/15.8.3/iPhone_com.toastmobile.PixelPlex"},
		"Sec-WebSocket-Extensions": {"permessage-deflate; client_max_window_bits"},
		"Accept-Language":          {"it-IT,it;q=0.9"},
		"Connection":               {"Upgrade"},
		"Upgrade":                  {"websocket"},
	}

	// Connect to the WebSocket server
	ws, err := websocket.DialConfig(config)
	if err != nil {
		log.Printf("[wss:error]: %v\n", err)
		return err
	}

	defer ws.Close()
	once := true

	// Log detailed responses from the server
	for {
		var message string
		err = websocket.Message.Receive(ws, &message)
		if err != nil {
			log.Printf("Read error: %v\n", err)
			continue
		}
		//log.Printf("[C<-S][response]: %s\n", message)

		//wallet := `{"t":"d","d":{"r":72,"a":"q","b":{"p":"\/users\/5576037373575168\/wallet","h":""}}}`
		//resources := `{"t":"d","d":{"r":72,"a":"q","b":{"p":"\/cities/5319377040179200\/resources\/resident\/5616078598701056\/current","h":""}}}`
		collects := fmt.Sprintf(`{"t":"d","d":{"r":13,"a":"q","b":{"p":"\/cities\/%s\/collects","h":""}}}`, player.CityId)
		limits := fmt.Sprintf(`{"t":"d","d":{"r":13,"a":"q","b":{"p":"cities/%s/limits/resident/%s","h":""}}}`, player.CityId, player.ResidentId)

		if strings.Contains(message, "cities/5319377040179200/limits/resident/5616078598701056") {

			var limits WssLimitsReponse
			err := json.Unmarshal([]byte(message), &limits)
			if err != nil {
				log.Printf("[C<-S][wss::limits.response]: %v\n", message)
				continue
			}

			player.UpdateResources(nil, &superplayer.Warehouse{
				Money:  limits.D.B.D.Current.Resources.Money,
				Wood:   limits.D.B.D.Current.Resources.Wood,
				Cement: limits.D.B.D.Current.Resources.Cement,
				Steel:  limits.D.B.D.Current.Resources.Steel,
			})
			//log.Printf("[C<-S][wss::limits.response]: %v\n", limits.D.B.D.Current.Resources.Money)
		} else if strings.Contains(message, fmt.Sprintf("\"cities/%s/collects\"", player.CityId)) {
			var collects WssCollectsResponse
			errcollects := json.Unmarshal([]byte(message), &collects)
			if errcollects != nil {
				fmt.Println("c.error: ", errcollects)
			}
			log.Printf("wss::collects.first")
			if collects.D.B.D.Rover.Available {
				player.CollectRover()
				log.Printf("[%s] wss::collect->rover\n", player.ResidentId)
			}

			if collects.D.B.D.Airplane.Available {
				player.CollectAirplane()
				log.Printf("[%s] wss::collect->airplane\n", player.ResidentId)
			}

			if collects.D.B.D.Blimp.Available {
				player.CollectBlimp()
				log.Printf("[%s] wss::collect->blimp\n", player.ResidentId)
			}

			if collects.D.B.D.Boat.Available {
				player.CollectBoat()
				log.Printf("[%s] wss::collect->boat\n", player.ResidentId)
			}

			if collects.D.B.D.TruckWood.Available {
				player.CollectTruck("wood")
				log.Printf("[%s] wss::collect->truck-wood\n", player.ResidentId)
			}

			if collects.D.B.D.TruckCement.Available {
				player.CollectTruck("cement")
				log.Printf("[%s] wss::collect->truck-cement\n", player.ResidentId)
			}

			if collects.D.B.D.TruckSteel.Available {
				player.CollectTruck("steel")
				log.Printf("[%s] wss::collect->truck-steel\n", player.ResidentId)
			}

		} else if strings.Contains(message, fmt.Sprintf("\"cities/%s/collects/", player.CityId)) {
			var update WssCollectsUpdate

			errupdate := json.Unmarshal([]byte(message), &update)
			if errupdate != nil {
				fmt.Println("u.error: ", errupdate)
			}
			log.Printf("[wss::collects.update]: %v\n", update.D.B.P)
			if strings.Contains(update.D.B.P, "rover") && update.D.B.D.Available {
				player.CollectRover()
				log.Printf("[%s] wss::collect->rover\n", player.ResidentId)
			} else if strings.Contains(update.D.B.P, "blimp") && update.D.B.D.Available {
				player.CollectBlimp()
				log.Printf("[%s] wss::collect->blimp\n", player.ResidentId)
			} else if strings.Contains(update.D.B.P, "airplane") && update.D.B.D.Available {
				player.CollectAirplane()
				log.Printf("[%s] wss::collect->airplane\n", player.ResidentId)
			} else if strings.Contains(update.D.B.P, "boat") && update.D.B.D.Available {
				player.CollectBoat()
				log.Printf("[%s] wss::collect->boat\n", player.ResidentId)
			} else if strings.Contains(update.D.B.P, "truck") && update.D.B.D.Available {
				if strings.Contains(update.D.B.P, "truck-wood") {
					player.CollectTruck("wood")
					log.Printf("[%s] wss::collect->truck-wood\n", player.ResidentId)
				} else if strings.Contains(update.D.B.P, "truck-cement") {
					player.CollectTruck("cement")
					log.Printf("[%s] wss::collect->truck-cement\n", player.ResidentId)
				} else if strings.Contains(update.D.B.P, "truck-steel") {
					player.CollectTruck("steel")
					log.Printf("[%s] wss::collect->truck-steel\n", player.ResidentId)
				}
			}
			//} else if strings.Contains(message,  fmt.Sprintf("\"users/%s", player.CityId)) ||
			//strings.Contains(message, `cities/5319377040179200/resources/resident/5616078598701056/current`)) {
			//log.Printf("[C<-S][wss::resources.response]: %v\n", message)
		} else {
			//fmt.Println("[S->C][wss::not_match]: ", message)
		}

		if once {
			// * subscribe to collects firestore collection
			log.Printf("[wss][C->S][collects.subscribe]: %s\n", collects)
			websocket.Message.Send(ws, collects)

			go func() {
				for {
					websocket.Message.Send(ws, limits)
					player.CollectBlimp()
					player.CollectRover()
					player.CollectCity()
					time.Sleep(10 * time.Second)
				}
			}()
			once = false
		}

		continue
	}
}

// lastonline := "{\"t\":\"d\",\"d\":{\"r\":65,\"a\":\"q\",\"b\":{\"p\":\"\\/lastonline\\/5616078598701056\",\"h\":\"\"}}}",
//jsonString := `{"t":"d","d":{"r":5,"a":"m","b":{"p":"\/users\/5576037373575168\/wallet","d":{"Pixels":100}}}}`
//jsonString := `{"t":"d","d":{"r":6,"a":"p","b":{"p":{"/users/5576037373575168/wallet/Pixels":100}}}}`
//jsonString := `{"t":"d","d":{"r":63,"a":"p","b":{"p":"\/users\/5576037373575168\/wallet\/Pixels","d":100}}}`
//jsonString := `{"t":"d","d":{"r":3,"a":"q","b":{"p":"\/","h":""}}}`
//jsonString := `{"t":"d","d":{"r":3,"a":"qs","b":{"p":"\/","h":"","sh":true}}}`
//jsonString := `{"t":"d","d":{"r":4,"a":"q","b":{"p":"\/","q":{"l":1,"vf":"r"},"h":""}}}`
//jsonString := `{"t":"d","d":{"r":5,"a":"q","b":{"p":"\/","sh":true}}}`
//jsonString := `{"t":"d","d":{"r":6,"a":"q","b":{"p":"\/cities\/5319377040179200\/residents","h":""}}}`
//jsonString := `{"t":"d","d":{"r":6,"a":"q","b":{"p":"\/lastonline","h":""}}}`
//jsonString := `{"t":"d","d":{"r":72,"a":"n","b":{"p":"\/lastonline\/5616078598701056"}}}`
//jsonStringO := `{"t":"d","d":{"r":61,"a":"o","b":{"p":"\/lastonline\/5616078598701056","d":1737384099000}}}`
//jsonStringP := `{"t":"d","d":{"r":64,"a":"p","b":{"p":"\/lastonline\/5616078598701056","d":true}}}`
//jsonStringQ := "{\"t\":\"d\",\"d\":{\"r\":65,\"a\":\"q\",\"b\":{\"p\":\"\\/lastonline\\/5616078598701056\",\"h\":\"\"}}}"

//jsonX := `{"t":"d","d":{"r":1,"a":"n","b":{"p":"\/cities\/5319377040179200\/collects"}}}`
//jsonString := `{"t":"d","d":{"r":62,"a":"o","b":{"p":"cities/5319377040179200/resources/resident/5616078598701056/current","d":{"Cement":1000,"Energy":500,"Food":300,"Money":9000,"Nano":50,"Oxygen":200,"Steel":2200,"Wood":3300}}}}`
//jsonStringDelete := `{"t":"d","d":{"r":71,"a":"d","b":{"p":"cities\/5319377040179200\/resources\/resident\/5616078598701056\/current"}}}`
//jsonString := `{"t":"d","d":{"r":3,"a":"q","b":{"p":"\/","sh":true}}}`
