package pkg

func ChanModeLevel(mode string) (size int){
	switch mode {
	case "slow":
		return 100
	case "hight":
		return 500
	case "fast":
		return 1000
	case "quick":
		return 1500
	default:
		return 50
	}
}



