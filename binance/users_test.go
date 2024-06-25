package binance

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/spf13/viper"
)

func GetEnv() (string, string) {
	viper.SetConfigFile("../.env") // .env 파일 설정
	viper.AutomaticEnv()           // 환경 변수를 자동으로 읽도록 설정
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	accKey := viper.GetString("AccessKey")
	seckey := viper.GetString("SecretKey")
	return accKey, seckey
}

func GetUsrs() *BinanceUser {
	AccessKey, SecritKey := GetEnv()
	binance := new(BinanceUser)
	binance.Init(AccessKey, SecritKey)
	return binance
}

func GetUsrsData(AccessKey, SecritKey string) *BinanceUser {
	binance := new(BinanceUser)
	binance.Init(AccessKey, SecritKey)
	return binance
}

func TestGetStartUserStreamService(t *testing.T) {
	res, err := GetUsrs().GetStartUserStreamService()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func TestGetLeverageBracket(t *testing.T) {
	res, err := GetUsrs().GetLeverageBracket("ETHUSDT")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

// positionModel
func TestGetChangePositionModeService(t *testing.T) {
	err := GetUsrs().GetChangePositionModeService(true)
	if err != nil {
		fmt.Println(err)
	}
}

func TestGetBalanceService(t *testing.T) {
	res, err := GetUsrs().GetBalanceService()
	if err != nil {
		fmt.Println(err)
	}

	for _, v := range res {
		// fmt.Println(v.AccountAlias, v.Asset, v.Balance, v.CrossWalletBalance, v.CrossUnPnl, v.AvailableBalance, v.MaxWithdrawAmount)
		if v.Asset == "USDT" {
			fmt.Println("v.AvailableBalance ", v.AvailableBalance)
		}
	}
}

func TestGetPositionRiskService(t *testing.T) {
	res, err := GetUsrs().GetPositionRiskService("ETHUSDT")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func TestGetListOpenOrdersService(t *testing.T) {
	res, err := GetUsrs().GetListOpenOrdersService("")
	if err != nil {
		fmt.Println(err)
	}

	for _, v := range res {
		fmt.Printf("%+v\n", v)
	}
}

func TestCreateMuiOrder(t *testing.T) {
	var send []*futures.CreateOrderService
	var order OpenOrder
	order.Symbol = "BIGTIMEUSDT"
	order.Side = futures.SideTypeBuy                  // SideTypeBuy SideTypeSell
	order.PositionSide = futures.PositionSideTypeLong // PositionSideTypeLong PositionSideTypeShort
	order.Type = futures.OrderTypeLimit               // OrderTypeLimit OrderTypeMarket
	order.Quantity = "50"
	order.Price = "0.1000"
	order.TimeInForce = futures.TimeInForceTypeGTC

	createOrderService := GetUsrs().CreateOrderLimitMarket(order)
	send = append(send, createOrderService)
	res, err := GetUsrs().CreateMuiOrder(send)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func TestCreateMuiOrderMarket(t *testing.T) {
	var send []*futures.CreateOrderService
	var order OpenOrder
	order.Symbol = "BIGTIMEUSDT"
	order.Side = futures.SideTypeSell                  // SideTypeBuy SideTypeSell
	order.PositionSide = futures.PositionSideTypeShort // PositionSideTypeLong PositionSideTypeShort
	order.Type = futures.OrderTypeMarket               // OrderTypeLimit OrderTypeMarket
	order.Quantity = "500"
	// order.Price = "0.2177"
	// order.TimeInForce = futures.TimeInForceTypeGTC

	createOrderService := GetUsrs().CreateOrderLimitMarket(order)
	send = append(send, createOrderService)
	res, err := GetUsrs().CreateMuiOrder(send)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

// 시장가 청산
func TestCreateClosePosition(t *testing.T) {
	var send []*futures.CreateOrderService
	var order OpenOrder
	order.Symbol = "BIGTIMEUSDT"
	order.Side = futures.SideTypeBuy                   // SideTypeBuy SideTypeSell
	order.PositionSide = futures.PositionSideTypeShort // PositionSideTypeLong PositionSideTypeShort
	order.Type = futures.OrderTypeMarket               // OrderTypeLimit OrderTypeMarket
	order.Quantity = "500"

	createOrderService := GetUsrs().CreateOrderLimitMarket(order)
	send = append(send, createOrderService)
	res, err := GetUsrs().CreateMuiOrder(send)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func TestCreateMuiOrderTakeProfit(t *testing.T) {
	// var send []*futures.CreateOrderService
	var order OrderType
	order.Symbol = "BIGTIMEUSDT"
	order.Side = futures.SideTypeBuy                    // SideTypeBuy SideTypeSell
	order.PositionSide = futures.PositionSideTypeShort  // PositionSideTypeLong PositionSideTypeShort
	order.OrderType = futures.OrderTypeTakeProfitMarket // OrderTypeLimit OrderTypeMarket OrderTypeTakeProfitMarket
	order.StopPrice = "0.210"
	order.ClosePosition = true
	order.WorkingType = futures.WorkingTypeContractPrice // WorkingTypeContractPrice
	// createOrderService := CreateOrderLimitMarket(order)
	// send = append(send, createOrderService)
	// res, err := GetUsrs().CreateMuiOrder(send)
	res, err := GetUsrs().CreateOrderService(order)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func TestCreateMuiOrderStopMarket(t *testing.T) {
	// var send []*futures.CreateOrderService
	var order OrderType
	order.Symbol = "BIGTIMEUSDT"
	order.Side = futures.SideTypeBuy                  // SideTypeBuy SideTypeSell
	order.PositionSide = futures.PositionSideTypeLong // PositionSideTypeLong PositionSideTypeShort
	order.OrderType = futures.OrderTypeStopMarket     // OrderTypeLimit OrderTypeMarket OrderTypeTakeProfitMarket
	order.StopPrice = "0.2150"
	order.Quantity = "500"
	order.ClosePosition = false
	order.WorkingType = futures.WorkingTypeContractPrice // WorkingTypeContractPrice
	// createOrderService := CreateOrderLimitMarket(order)
	// send = append(send, createOrderService)
	// res, err := GetUsrs().CreateMuiOrder(send)
	res, err := GetUsrs().CreateOrderService(order)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

/////////////// 개별테스트

///// 데이터 가져오기 - 우선 포지션이 있는지 확인

type BinanceKey struct {
	ID        string
	AccessKey string
	SecretKey string
}

func TestHavePosition(t *testing.T) {

	list := []BinanceKey{
		{ID: "4039964145338795009", AccessKey: "YAzzGlOZCfc1ztnwG9YSzZkbA4SS24bynIPCgiaqRVGpu1CWgboXFwO9C1R2lngb", SecretKey: "QM24BeRle4Q99YDDXshCekJu0JPzRLV48zkuPFacsumd66OyLln3P5w4U8IksRdf"},
		{ID: "3954806416859529985", AccessKey: "s0xQSJ0aZH5LUmLh9kZbqRqTICCc5CozZ1yqwyyJrSD2kPbfhkaHjaadILIFXD1B", SecretKey: "eAbeZy2HZAHZYrPeKSKkV7FHI06VSKrd7E9Ogik3Nr9kTrVNPxipJ4LOx3ph1vVz"},
		{ID: "4039157981181872128", AccessKey: "Yg0aACxN64bVByllQPEa6lLkG5YornlNExhA3NPkJ68SFBijuDii1OSx9JyGSNSm", SecretKey: "MdVhdUt9BMVMbDAzRjKKNCu3ozhzv7W2xunWKUoLLbISeh3HfSlkYUofuPvVdWJ7"},
		{ID: "4039038463775381248", AccessKey: "xmtLw3IBUfuiHQwBgxab59ftgnQqw5X0GOiC5gGVknkBwgFOe96vNgLdMeyqYBY2", SecretKey: "lyxjSvcXFO7lsS9sPm2J7t5FXF2DGMPLtV1cCTTQxjkJmWPDNXMqzt9NM2TXfYTt"},
		{ID: "4039156260863751424", AccessKey: "XCH510F361bm47EgxYjZ4s9KfHqFv7XsdCKNtqlGLToXKnesJAmpNpur5fpXT4Uu", SecretKey: "gKeF9gDNnC5jSaXDvU2XjoJCjLuyPbj1qiHRDHpveHKoJNvNAKpX4zUSMIm1KLid"},
		{ID: "4015682822927577856", AccessKey: "aAOUtZanr4GhgtX3lh53f1TJYrlfrvCcyhsaBpXRrT9SKFVEhCblGEea6AYw2sT8", SecretKey: "rgoOCQH2xGNO3QCTe0HHDyYkrbj1PAz3OOMiyLjLi9k4fdgLjXLLCFjMGfgx0QK8"},
		{ID: "4018735578681644800", AccessKey: "fZTUxku1U8m3HL8KoMTiSVxUfBD7Nx8uXgOSJEsiXYg3sSEYcN2fGnESW6UouqF7", SecretKey: "giZfOozF6xwvM1iihAJlATmT3Vlg84IxvT3yD5faPJFQfKxu98PVNbV9eSFadtqI"},
		{ID: "4021226350652327680", AccessKey: "mEoMpXGUfaPbSgDmFjGS98LVl1lvg6StqxtuRs07r1W0y0nabIcX4TWilftM1TXl", SecretKey: "PZu2bhM1UWro9qi6ub2eat8ZLBBT1b7x3kvS3x2aUlQpkyd67j8bjTbVAaSoRYkn"},
		{ID: "4040174027775910913", AccessKey: "PMvVlmaR6bWnwMf3QO3KRHC51mq1aQJC37MAHtdoGgrAwg5Kxz5pm9N05KupY4cD", SecretKey: "3HctK0Wjhto3GDXOCloFEslmUImKv9AockOkOFeLGVpFkChqLqlxEY6tzEfXbWpW"},
		{ID: "4011046519534295041", AccessKey: "Y9L2Ldf6SaLu8UpKjRxQKI7OZUHe0PMd5IGm5dkQo0kX0WCljilgsZarKFiWcNld", SecretKey: "RbVdztn4hyblFNbP6vV0dF2GDfiBkBA59Vte6acPn51FvShBSkjA2kScodRKy4OF"},
		{ID: "4011033153694187520", AccessKey: "t4SrsHHeI3uej25VPB6q9HfEOBq1Ai4amfue4oCfR9MnJdppQEfRcf9Exgacn8bP", SecretKey: "LNQwBZsLwmMdOiK2LD2UqIaOoK3O21fqlw234SHcYiPaQWeINw1Dk39wvbGQ2Yhc"},
		{ID: "4039080765758989313", AccessKey: "Am9phZvTOEqaUEigb5bf4zAQNhagdXg5C3xnJegQFPCHtLi21k0zO5DsKPb0a8A9", SecretKey: "zsGWXmTMy7juGCBqEqTO7BZJZhlsip4xAPqmmHt9QUTcRtrhHnZsOdIJu10esmFU"},
		{ID: "4039116610783141376", AccessKey: "gOKMFWnaqKkj9YykddYfHH64iP5OeP6UgfsTzj7DQWl6wWrHvlUXR709zexvH3AS", SecretKey: "7tW1lqvRGUkvXFwcHips3Db7miCOSzhwxTTSl76y8lGgVOHYCFRz511yxYx4EWNc"},
		{ID: "3997012228324541952", AccessKey: "YoRVAlqAuXUxNUxFk9UHLhxeDtpKsNjTkt0HkNDqcVV53P33rwW9ncU1BkKgOXXn", SecretKey: "JSpFymLENLWCT9P5vSNX5tEKNzvMA0BFGsjOuXyLTeiHJ6gjbh5ExGyXO1qDtWgd"},
		{ID: "4018747643979672576", AccessKey: "2Ow7uvlz2Au8c4xittqjTFu0XMvg7T7Tx1XDGxDAPEzuMCUpq9eLocmXyF9CEtmb", SecretKey: "8NSQGTq4j2PfbK2LavrOcCbzHM84MrhnNNValVrNlN3nAmHnijUJJ6NiECiYJIFY"},
		{ID: "4018745341040561153", AccessKey: "zNXG7P2iBhooHPREEE3aqgzSEhGooOZ5k2PaZjPG4fiCHAViX1EdRjZnit6wdE1i", SecretKey: "5moGTngqRrYO1xXEyEH6KEaqnAJdDgefRqhWFRATY0VKDcoSa21eOTrHLPT3McwP"},
		{ID: "4025500238509219328", AccessKey: "NCVL2svrUQwSbW8W7v7449LnV2dVFLvIek9F1Op6yizNnDCJnhJ30iTCMWI94EtC", SecretKey: "XNRSnFAVqyWLqNucJCZFsYT1rdIB6bemyWtnj94FrkRWhGq3ZnHOFa8kVfVWT5wo"},
		{ID: "4018750548821782017", AccessKey: "PBRPFw4tYAVTAC8Pya1eyj0K2WLqIPLLgHd8lGikEVVVHuS90QEjstSd6tp1V6cL", SecretKey: "iEL3kK5nJhCXGAb3bO8drzPjsS5uTp0ATRB323YXFb3AjldCruByINuMK3voGaQP"},
		{ID: "4018731341273424385", AccessKey: "NK3owPQnbFpGfUhMcEvOtcYke1wp7LkNuSar4UouTHvB5h0dZrCaYOInjw8ZzuJ9", SecretKey: "VcY2C9TphEg9fRrrxXsoheQ4kl1XFNw5MsIQGNp9oodgIVcOqteStjZ0Nc4exQv8"},
		{ID: "4020187803795927553", AccessKey: "S0gI8aLagAuOu5E4kV2n6UEUyAFjQcwzRw6vKcGDUVTatt7R6gW2F6xGtvHJzaJh", SecretKey: "nBxZ1blpY05wVXIQZID6hpMPpdB8VBEjMkDCYJzGz6iYUEcql4ohWcprYzQiC4mx"},
		{ID: "4039076213614265345", AccessKey: "f4Wx20IDrMqlbgXznMpONq764f09XTdi31jZd6Wq0olukUYerU2HhPgBZMqCuqle", SecretKey: "mD8IoQfMIlbMyBUiyM88xysdBL9QikMvvqcchbZAw97dFZbK9TABg1SFiDOtXKLT"},
		{ID: "4039123057902068737", AccessKey: "Cp6XT9V8Uqgg7xb6jksr5OgfSHIxjDPRlDE7TSE3iF0cPFm6rnqAdAmVsHyvfY8J", SecretKey: "6WCu5mHvYcvuLKHQ5M4e0RbP87iuBeHc2fMXgU3YEGyJY0wo0udDJ7h1CanW4yCT"},
		{ID: "4039053759711716097", AccessKey: "V1SDxgNL50bWb9WYxh24g5Ty5dOVf5TAevWIXo2bTHwrmAPUUDBSyMql5YIN4rNN", SecretKey: "EEo1Qcj1wycP8viIIlf4n0zKD85ZHnMmdIP3QXyDi9H7dTobgajDM0GRzcuc7XbB"},
		{ID: "4039114540087713025", AccessKey: "59TUevAXeg1Qkv7HWKCr6YWbXzzmxEyuaFf8n77sYSc1p3uqTORxo28DUh37AAH9", SecretKey: "rkXsxFBBrrOantydD3oauWVUkNm2Nv4gBTX4lDn65QHbMVNvuHpKzZnPz9dWjr7K"},
		{ID: "4018734070798159105", AccessKey: "Hjg90Uz87jdxi5wgSsc4feTp2HNZJia7qbTYfgAXF5Q5uaRCEPPRNjgh0Ewb5I0e", SecretKey: "0zuOpzI7XqjMx6RB4mBWQRMwFYzlKoNEjxnjbJmGYSUHPRxRhsb3bFTlkHC2Gibn"},
		{ID: "4009873035671330816", AccessKey: "BwOe3fW58vJPOBir31qMwZYvNPB3f3I4o7vSqDlQ9OSCDW9DEtQWKFcSEjNvTUQz", SecretKey: "aiWpEpNAHKEhAF62mKfhNS7rEivGPBDDtg1ACU5vKmIpKQyVVzDCyIMr5HySiLYq"},
		{ID: "4038948155727837441", AccessKey: "vZephxoRcUcfFD2OpGpdgWvV7OwRsTPubazfchKIfD9ZPIch0CUAAnuxh8MmPiAT", SecretKey: "HZsfdoSY3BJjhRwSySwuzlItphfpe9l8T3xzT0ujsHaR12ab5mGsvOx74XZURaIv"},
		{ID: "4018743075963202561", AccessKey: "UMZygwbkWhJDfFktVAeJnPTVHBQkfeSvBtb3BNYf9o7c63UyYB3DOrETeMcImOkv", SecretKey: "ivjWtfKkrO6m2Sc522ZlasKVyEur3jnp3TDIl0XTX1gbz6peLHggpMKHR11CPurM"},
		{ID: "4011041169401465856", AccessKey: "of4JsNjSFbZQJP6ABBqLyIhaiX5yOkJSOEIGzwC0V6R4YZlG8pXoBEz2PQPzd25p", SecretKey: "tgJA9XCwq5nd8xuxglVN4apf7lPGyoxWcJjEvYmsJfmjvOD4YGTsqHmZbZExvw1X"},
		{ID: "4011038316181309440", AccessKey: "uZomHDIzDQezQY0VBOgfv30kBaq7VhuR5s3sBu02cjZLqqX0qYF5qyX5NUIIhJ5r", SecretKey: "LZcGEZkQnIuzP6HINAdcisRNzzl8nqMzwBwgbcKzSypUCWHeoA2PvgliBmPZqUkp"},
		{ID: "3971951508369543937", AccessKey: "BC4OhmWJOr9J9dEXyOtZsNJewLoV6s2DEdA9oLePBKjwTA2KnlGZq2AFs3iwft8s", SecretKey: "mdCF9iiA3qSaTeSckMplXQ1M7UP1JRaibETik1vdepApnqGCOjtbEJD6tEvQvAi8"},
		{ID: "4039045090432982272", AccessKey: "o0Ixp0Cl8H9uQKh82LPTDTvxx9RkuNkedIGgwH37sWJ0cFv6zzwUfIkKrZ0bAtxA", SecretKey: "8DkcPzsyJLuFI7jx8fQ9o3Qc8Rw3InZ0LINjJNvJ3xjkR2sFsm9mPJq7E2KUHL9H"},
		{ID: "4031120570965721344", AccessKey: "0vAs0MrT7XaRMdZFlQH9uW4LuUEOdFFKpPZpUdXMFIvqHfvvhfK38nHWYcxP8WJJ", SecretKey: "d02Shxrg7yqMTU0cT1SM9PTfDBhbs9UQOTLqpUmueMaMzIkF3Je1WQLBvoRRK8P3"},
		{ID: "4024637971593983233", AccessKey: "kYDtRgERgMXBcgCYx7vPQn2VLKEW15cSh9akmxgqqqBpNvavYwwJULb71c1VRa8M", SecretKey: "VvF2R9pIvbn65pyev9qFW5kL0pz3abm6yS98HSZpyt0pnwJY57PsxS3lfHgx1Kem"},
	}

	for _, v := range list {
		client := GetUsrsData(v.AccessKey, v.SecretKey)

		res, err := client.GetPositionRiskService("")
		if err != nil {
			fmt.Println(err)
		}

		zeroValues := map[string]bool{
			"0":     true,
			"0.0":   true,
			"0.00":  true,
			"0.000": true,
		}

		for _, k := range res {
			/*
				amount, err := utils.Float64(k.PositionAmt)
				if err != nil {
					fmt.Println(err)
				}

				if amount > 0 {
					fmt.Println(v.ID, " || ", k.Symbol, " : ", k.PositionSide, " : ", k.PositionAmt)
				}
			*/

			if !zeroValues[k.PositionAmt] {
				fmt.Println(v.ID, " || ", k.Symbol, " : ", k.PositionSide, " : ", k.PositionAmt, " || ", k.UnRealizedProfit)
			}

		}
		time.Sleep(time.Second)

		/*
			// 벨런스
			res, err := client.GetBalanceService()
			if err != nil {
				fmt.Println(err)
			}

			for _, k := range res {

				if k.Asset == "USDT" {
					fmt.Println(v.ID, " : ", k)
				}
			}
		*/

	}

}
