package main

func main() {
	discordMain(Token)
}

var (
	Token string
)

// func init() {
// 	flag.StringVar(&Token, "t", "", "-t <discord_bot_token>")
// 	flag.Parse()

// 	if Token == "" {
// 		flag.Usage()
// 		os.Exit(1)
// 	}
// }