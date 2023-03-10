package gungi

import "fmt"

func (piece *Piece) GetMove() {
	switch int(*piece) {
	case BPA, WPA:
		fmt.Println("hello")
	case BLG, WLG:
		fmt.Println("hello")
	case BMG, WMG:
		fmt.Println("hello")
	case BGE, WGE:
		fmt.Println("hello")
	case BFO, WFO:
		fmt.Println("hello")
	case BKN, WKN:
		fmt.Println("hello")
	case BAR, WAR:
		fmt.Println("hello")
	case BMU, WMU:
		fmt.Println("hello")
	case BSA, WSA:
		fmt.Println("hello")
	case BCN, WCN:
		fmt.Println("hello")
	case BSP, WSP:
		fmt.Println("hello")
	case BCP, WCP:
		fmt.Println("hello")
	case BMA, WMA:
		fmt.Println("hello")
	default:
		fmt.Println("false")
	}
}
