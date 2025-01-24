package main

type WssCollectsResponse struct {
	T string `json:"t"`
	D struct {
		B struct {
			P string `json:"p"`
			D struct {
				Airplane struct {
					Available     bool   `json:"available"`
					Lastcollector string `json:"lastcollector"`
				} `json:"airplane"`
				Blimp struct {
					Available     bool   `json:"available"`
					Lastcollector string `json:"lastcollector"`
				} `json:"blimp"`
				Boat struct {
					Available     bool   `json:"available"`
					Lastcollector string `json:"lastcollector"`
				} `json:"boat"`
				Rover struct {
					Available     bool   `json:"available"`
					Lastcollector string `json:"lastcollector"`
				} `json:"rover"`
				TruckCement struct {
					Available     bool   `json:"available"`
					Lastcollector string `json:"lastcollector"`
				} `json:"truck-cement"`
				TruckSteel struct {
					Available     bool   `json:"available"`
					Lastcollector string `json:"lastcollector"`
				} `json:"truck-steel"`
				TruckWood struct {
					Available     bool   `json:"available"`
					Lastcollector string `json:"lastcollector"`
				} `json:"truck-wood"`
			} `json:"d"`
		} `json:"b"`
		A string `json:"a"`
	} `json:"d"`
}

type WssCollectsUpdate struct {
	T string `json:"t"`
	D struct {
		B struct {
			P string `json:"p"`
			D struct {
				Available     bool   `json:"available"`
				Lastcollector string `json:"lastcollector"`
			} `json:"d"`
		} `json:"b"`
		A string `json:"a"`
	} `json:"d"`
}

type WssResourcesResponse struct {
	T string `json:"t"`
	D struct {
		B struct {
			P string `json:"p"`
			D struct {
				Cement int `json:"Cement"`
				Energy int `json:"Energy"`
				Food   int `json:"Food"`
				Money  int `json:"Money"`
				Nano   int `json:"Nano"`
				Oxygen int `json:"Oxygen"`
				Steel  int `json:"Steel"`
				Wood   int `json:"Wood"`
			} `json:"d"`
		} `json:"b"`
		A string `json:"a"`
	} `json:"d"`
}

type WssWalletResponse struct {
	T string `json:"t"`
	D struct {
		B struct {
			P string `json:"p"`
			D struct {
				Achievements struct {
					Num0  int `json:"0"`
					Num1  int `json:"1"`
					Num2  int `json:"2"`
					Num3  int `json:"3"`
					Num4  int `json:"4"`
					Num5  int `json:"5"`
					Num6  int `json:"6"`
					Num7  int `json:"7"`
					Num8  int `json:"8"`
					Num9  int `json:"9"`
					Num10 int `json:"10"`
					Num11 int `json:"11"`
					Num12 int `json:"12"`
					Num13 int `json:"13"`
					Num14 int `json:"14"`
					Num15 int `json:"15"`
					Num16 int `json:"16"`
					Num17 int `json:"17"`
					Num18 int `json:"18"`
					Num19 int `json:"19"`
					Num20 int `json:"20"`
					Num21 int `json:"21"`
					Num22 int `json:"22"`
					Num23 int `json:"23"`
					Num24 int `json:"24"`
					Num25 int `json:"25"`
					Num26 int `json:"26"`
					Num27 int `json:"27"`
					Num28 int `json:"28"`
				} `json:"Achievements"`
				Pixels  int `json:"Pixels"`
				Unlocks struct {
					Num0 string `json:"0"`
					Num1 string `json:"1"`
					Num2 string `json:"2"`
				} `json:"Unlocks"`
				UserID int64 `json:"UserId"`
			} `json:"d"`
		} `json:"b"`
		A string `json:"a"`
	} `json:"d"`
}

type WssLimitsReponse struct {
	T string `json:"t"`
	D struct {
		B struct {
			P string `json:"p"`
			D struct {
				Current struct {
					BirthCenterLevel int   `json:"BirthCenterLevel"`
					BuildLevel       int   `json:"BuildLevel"`
					Constructions    int   `json:"Constructions"`
					ContractorLevel  int   `json:"ContractorLevel"`
					HQLevel          int   `json:"HQLevel"`
					MaxShares        int   `json:"MaxShares"`
					ResidentID       int64 `json:"ResidentId"`
					Resources        struct {
						Cement int `json:"Cement"`
						Energy int `json:"Energy"`
						Food   int `json:"Food"`
						Money  int `json:"Money"`
						Nano   int `json:"Nano"`
						Oxygen int `json:"Oxygen"`
						Steel  int `json:"Steel"`
						Wood   int `json:"Wood"`
					} `json:"Resources"`
					TradeAmount int `json:"TradeAmount"`
				} `json:"current"`
			} `json:"d"`
		} `json:"b"`
		A string `json:"a"`
	} `json:"d"`
}
