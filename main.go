package main

import (
	"github.com/antonybro/gramophone/telegram"
	"github.com/antonybro/gramophone/spotify"

	"time"
	"strings"
	"regexp"
	"strconv"
)

func main() {
	telegram.Authorization()
	telegram.Listen(handler)
}

func handler(command string, text string) {
	switch command {
	case "start":
		text := `
			In order to work with the bot, you need to log into your Spotify account. 
			Bot uses the OAuth 2.0 protocol for authorization.
			In order to authorize in Spotify you need to execute the command /login`

		telegram.Send(text)

	case "login":
		url := spotify.GetURL()
		print(url)
		telegram.Send(url)

		text := spotify.AuthorizationListener()
		telegram.Send(text)

	case "play":
		telegram.Send(spotify.Play())
		updateState()

	case "pause":
		telegram.Send(spotify.Pause())

	case "next":
		telegram.Send(spotify.Next())
		updateState()

	case "previous":
		telegram.Send(spotify.Previous())
		updateState()

	case "volume":
		str := strings.Trim(text, "/volume ")
		percent, err := strconv.Atoi(str)
		if err == nil {
			telegram.Send(spotify.Volume(percent))
		}

	case "shuffle":
		telegram.Send(spotify.Shuffle())

	case "playlist":
		telegram.Send(spotify.GetPlaylist())

	case "add":
		id := managePlaylist(text)
		telegram.Send(spotify.Add(id))

	case "remove":
		id := managePlaylist(text)
		telegram.Send(spotify.Remove(id))

	case "help":
		text := `
			Gramophone commands:

			/start - Activate bot
			/login - Login Spotify (OAuth 2.0)
			/play - Play or resume current playback
			/pause - Pause Playback
			/next - Next track
			/previous - Previous track
			/volume - Percent is must be a value from 0 to 100 inclusive.
			/playlist - Show current playlist
			/add - Add track to current playlist
			/remove - remove track from playlist if exist
			/help - List of commands`

		telegram.Send(text)

	default:
		id := managePlaylist(text)
		if id != "" {
			telegram.Send(spotify.Add(id))
		}
	}

	easterEgg(text)
}

func updateState() {
	rateLimit := time.Tick(500 * time.Millisecond)
	<-rateLimit
	telegram.Send(spotify.CurrentTrack())
}

func managePlaylist(text string) (id string) {
	url, _ := regexp.MatchString("open.spotify.com/track/.*", text)
	uri, _ := regexp.MatchString("spotify:track:.*", text)

	if url {
		id := strings.SplitAfter(text, "track/")[1]

		return strings.Split(id, "?")[0]
	}

	if uri {
		id := strings.SplitAfter(text, "track:")[1]
		return  id
	}

	return ""
}

func easterEgg(text string) {
	if strings.Contains(text, "баг") {
		telegram.Send("http://memesmix.net/media/created/kzi6fa.jpg")
	}

	if strings.Contains(text, "android") {
		telegram.Send("https://2ch.hk/v/arch/2017-02-15/src/1756924/14868514913860.jpg")
	}

	if strings.Contains(text, "ios") {
		telegram.Send("http://m.memegen.com/mdt0c7.jpg")
	}
}

