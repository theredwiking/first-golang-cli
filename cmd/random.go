package cmd

import (
	"fmt"
	"net/http"
    "io/ioutil"
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/spf13/cobra"
)

// randomCmd represents the random command
var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Get random dad joke",
	Long: `This commands fetch random dad joke from icanhazdadjoke api`,
	Run: func(cmd *cobra.Command, args []string) {
		jokeTerm, _ := cmd.Flags().GetString("term")

		if jokeTerm != "" {
			getJokeWithTerm(jokeTerm)
		} else {
			getRandomJoke()
		}
	},
}

func init() {
	rootCmd.AddCommand(randomCmd)

	randomCmd.PersistentFlags().String("term", "", "Search term for dad joke")
}

type Joke struct {
	ID string `json:"id"`
	Joke string `json:"joke"`
	Status int `json:"status"`
}

type SearchResult struct {
	Results json.RawMessage `json:"results"`
	SearchTerm string `json:"search_term"`
	Status int `json:"status"`
	TotalJokes int `json:"total_jokes"`
}

func getRandomJoke() {
	url := "https://icanhazdadjoke.com/"
	responseBytes := getJokeData(url)
	joke := Joke{}

	if err := json.Unmarshal(responseBytes, &joke); err != nil {
		fmt.Printf("Error unmarshalling. %v", err)
	}

	fmt.Println(string(joke.Joke))
}

func getJokeWithTerm(jokeTerm string) {
	total, results := getJokeDataTerm(jokeTerm)
	randomiseJokeList(total, results)
}

func randomiseJokeList (length int, jokeList []Joke) {
	rand.Seed(time.Now().Unix())

	min := 0
	max := length - 1

	if length <= 0 {
		err := fmt.Errorf("No jokes found with this term")
		fmt.Println(err.Error())
	} else {
		randomNum := min + rand.Intn(max-min)
		fmt.Println(jokeList[randomNum].Joke)
	}
}

func getJokeData(baseAPI string) []byte {
	request, err := http.NewRequest(
		http.MethodGet,
		baseAPI,
		nil,
	)

	if err != nil {
		log.Printf("Failed request to dadjoke. %v", err)
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("User-agent", "Dadjoke CLI (https://github.com/theredwiking/first-golang-cli")

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		log.Printf("Failed to create requets. %v", err)
	}

	responseBytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Printf("Failed to read body. %v", err)
	}

	return responseBytes
}

func getJokeDataTerm(jokeTerm string) (totalJokes int, jokeList []Joke){
	url := fmt.Sprintf("https://icanhazdadjoke.com/search?term=%s", jokeTerm)
	responseBytes := getJokeData(url)
	jokeListRaw := SearchResult{}

	if err := json.Unmarshal(responseBytes, &jokeListRaw); err != nil {
		log.Printf("Error unmarshalling responseBytes. %v", err)
	}

	jokes := []Joke{}
	if err := json.Unmarshal(jokeListRaw.Results, &jokes); err != nil {
		log.Printf("Error unmarshalling jokeListRaw. %v", err)
	}

	return jokeListRaw.TotalJokes, jokes
}
