package connect
import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func Banana() {

	s, err := discordgo.New("qwe011235@gmail.com", "QWERTY011235")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = s.Open()
	if err != nil {
		fmt.Println(err)
		return
	}

	s.ChannelMessageSend("835209409616412704", "test")

}