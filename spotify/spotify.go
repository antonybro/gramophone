package spotify

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zmb3/spotify"
	"strconv"
	"os"
	"math/rand"
)

const (
	redirectURI = "http://localhost:8080/callback"
	userId = "antonybrro"
	playlistId = "09nqGWwNZht8tKYbsOpnNT"
)

var (
	client *spotify.Client
	auth  = spotify.NewAuthenticator(redirectURI,
		spotify.ScopeUserModifyPlaybackState, spotify.ScopeUserReadCurrentlyPlaying,
		spotify.ScopeUserReadPlaybackState, spotify.ScopePlaylistModifyPublic,
		spotify.ScopePlaylistModifyPrivate)
	ch    = make(chan *spotify.Client)
	state = strconv.Itoa(rand.Int())
)

func GetURL() (string) {
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)

	auth.SetAuthInfo(os.Getenv("spotify_client_id"), os.Getenv("spotify_client_secret"))

	url := auth.AuthURL(state)

	return "Please log in to Spotify by visiting the following page in your browser:\n\n" + url
}

func AuthorizationListener() (string) {
	client = <-ch

	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}

	return "You are logged in as: " + user.ID
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	// use the token to get an authenticated client
	client := auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")
	ch <- &client
}

func checkError(text string, err error) (string) {
	if err != nil {
		return "⛔ " + err.Error()
	}

	return text
}

func Play() (string) {
	err := client.Play()

	return checkError("▶️", err)
}

func Pause() (string) {
	err := client.Pause()

	return checkError("⏸", err)
}

func Next() (string) {
	err := client.Next()

	return checkError("⏭", err)
}

func Previous() (string) {
	err := client.Previous()

	return checkError("⏮", err)
}

func Volume(percent int) (string) {
	err := client.Volume(percent)

	return checkError("🔊 Volume " + strconv.Itoa(percent), err)
}

func Shuffle() (string) {
	state, err := client.PlayerState()
	if err == nil {
		client.Shuffle(!state.ShuffleState)

		if !state.ShuffleState {
			return "🔀 Shuffle enabled"
		} else {
			return "🔁 Shuffle disabled"
		}
	}

	return "⛔️ Error while changing shuffle state"
}

func CurrentTrack() (string) {
	current, err := client.PlayerCurrentlyPlaying()

	return checkError(current.Item.ExternalURLs["spotify"], err)
}

func GetPlaylist() (string) {
	current, err := client.GetPlaylist(userId, playlistId)
	if err == nil {
		text := current.Description + "\n\n"

		tracks := current.Tracks.Tracks

		for index, track := range tracks {
			text += strconv.Itoa(index+1) + ") " + "🎹" + track.Track.Name + " - 🎙"
			for i := 0; i < len(track.Track.Artists); i++ {
				text += track.Track.Artists[i].Name

				if i < len(track.Track.Artists)-1 {
					text += ","
				}

			}

			text += " (" + track.Track.Album.Name + ")" + "\n"
		}

		return "📝" + text
	}

	return "⛔️ Error while getting playlist"
}

func Add(trackId string) (string) {
	_, err := client.AddTracksToPlaylist(userId, playlistId, spotify.ID(trackId))

	if err == nil {
		track, err := client.GetTrack(spotify.ID(trackId))
		return checkError("✅" + track.ExternalURLs["spotify"], err)
	}

	return "⛔️ Error while adding track"
}

func Remove(trackId string) (string) {
	_, err := client.RemoveTracksFromPlaylist(userId, playlistId, spotify.ID(trackId))

	return checkError("❌ Removed!", err)
}