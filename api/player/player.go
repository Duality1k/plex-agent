package superplayer

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

type Warehouse struct {
	Pixels int
	Money  int
	Wood   int
	Cement int
	Steel  int
}

type Player struct {
	UUID       string
	ResidentId string
	HttpClient *http.Client
	AuthCookie string
	CityId     string
	// Dynamic Set
	WarehouseUsage Warehouse
	WarehouseLimit Warehouse
	DomainControl  float64 // % of domain
}

type Players []*Player

func NewPlayer(uuid string, residentId string, cityId string) *Player {
	player := Player{
		UUID:       uuid,
		ResidentId: residentId,
		CityId:     cityId,
		HttpClient: &http.Client{Timeout: 10 * time.Second},
	}

	player.AuthCookie = player.LoginByUUID(uuid)

	//player.Me()
	//player.Basedata()
	return &player
}

func (player *Player) Login(endpoint, method, payload string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", "https://citygame-dev.appspot.com", endpoint)

	req, err := http.NewRequest(method, url, bytes.NewBufferString(payload))
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %w", err)
	}

	req.Header.Add("Host", "citygame-dev.appspot.com")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "PixelPlex/1 CFNetwork/1335.0.3.4 Darwin/21.6.0")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Language", "it-IT,it;q=0.9")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("App-Version", "4.4.1")

	res, err := player.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error making request: %w", err)
	}
	defer res.Body.Close()

	cookies := res.Header["Set-Cookie"]
	return []byte(cookies[0]), nil
}

func (player *Player) LoginByUUID(uuid string) string {
	endpoint := "/api/login"
	payload := fmt.Sprintf("uuid=%s&", uuid)
	method := "POST"

	body, _ := player.Login(endpoint, method, payload)

	return string(body)
}

func (player *Player) Me() {
	endpoint := "/api/me"
	method := "GET"

	body, _, err := player.HttpRequest(endpoint, method, "")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
}

func (player *Player) UpdateResources(usage *Warehouse, limit *Warehouse) {
	if limit != nil {
		player.WarehouseLimit = Warehouse{
			Money:  limit.Money,
			Wood:   limit.Wood,
			Cement: limit.Cement,
			Steel:  limit.Steel,
		}
	}

	if usage != nil {
		player.WarehouseUsage = Warehouse{
			Pixels: usage.Pixels,
			Money:  usage.Money,
			Wood:   usage.Wood,
			Cement: usage.Cement,
			Steel:  usage.Steel,
		}
	}
}

// money | wood | cement | steel
func (players Players) PoorestPlayerId(material string) (result *Player) {
	var minValue int
	var found bool

	for _, player := range players {
		var value int
		switch material {
		case "money":
			value = player.WarehouseUsage.Money
		case "wood":
			value = player.WarehouseUsage.Wood
		case "cement":
			value = player.WarehouseUsage.Cement
		case "steel":
			value = player.WarehouseUsage.Steel
		}

		if !found || value < minValue {
			minValue = value
			result = player
			found = true
		}
	}

	return result
}

func (player *Player) CollectGeneric(collectType string) (body []byte, code int, err error) {
	endpoint := fmt.Sprintf("/api/cities/%s/%s/collect", player.CityId, collectType)
	payload := fmt.Sprintf("resident=%s&", player.ResidentId)
	return player.HttpRequest(endpoint, "POST", payload)
}

func (player *Player) CollectBlimp() (body []byte, code int, err error) {
	return player.CollectGeneric("blimp")
}

func (player *Player) CollectRover() (body []byte, code int, err error) {
	return player.CollectGeneric("rover")
}

func (player *Player) CollectAirplane() (body []byte, code int, err error) {
	return player.CollectGeneric("airplane")
}

func (player *Player) CollectBoat() (body []byte, code int, err error) {
	return player.CollectGeneric("boats")
}

func (player *Player) CollectTruck(material string) (body []byte, code int, err error) {
	endpoint := fmt.Sprintf("/api/cities/%s/trucks/collect", player.CityId)
	payload := fmt.Sprintf("resident=%s&type=%s&", player.ResidentId, material)
	return player.HttpRequest(endpoint, "POST", payload)
}

func (player *Player) CollectCity() (body []byte, code int, err error) {
	endpoint := fmt.Sprintf("/api/cities/%s/collect", player.CityId)
	payload := fmt.Sprintf("resident=%s&isMars=false&", player.ResidentId)
	return player.HttpRequest(endpoint, "POST", payload)
}

// ! use wss

type Owners struct {
	Airplane    *Player
	Blimp       *Player
	Boat        *Player
	TruckWood   *Player
	TruckCement *Player
	TruckSteel  *Player
}

func InitOwners() *Owners {
	return &Owners{
		Airplane:    nil,
		Blimp:       nil,
		Boat:        nil,
		TruckWood:   nil,
		TruckCement: nil,
		TruckSteel:  nil,
	}
}
