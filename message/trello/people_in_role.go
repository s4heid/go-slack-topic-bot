package trello

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/VojtechVitek/go-trello"
	"github.com/dpb587/go-slack-topic-bot/message"
)

type PeopleInRole struct {
	Team   string
	Role   string
	People map[string]string
}

var _ message.Messager = &PeopleInRole{}

func getRequiredEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	panic("missing required environment variable " + key)
}

func getOptionalEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getCurrentList() string {
	t := time.Now()
	return fmt.Sprintf("pairing(%02d-%02d-%d)",
		t.Day(), t.Month(), t.Year())
}

func getBoard(user *trello.Member, board string) (trello.Board, error) {
	boards, err := user.Boards()
	if err != nil {
		log.Fatal(err)
	}

	for _, b := range boards {
		if b.Name == board {
			return b, nil
		}
	}
	return trello.Board{}, fmt.Errorf("cannot find the board %v", board)
}

func getList(board *trello.Board, list string) (trello.List, error) {
	lists, err := board.Lists()
	if err != nil {
		log.Fatal(err)
	}

	for _, l := range lists {
		if l.Name == list {
			return l, nil
		}
	}
	return trello.List{}, fmt.Errorf("cannot find the list %v", list)
}

func getCard(list *trello.List, card string) (trello.Card, error) {
	cards, err := list.Cards()
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range cards {
		if c.Name == card {
			return c, nil
		}
	}
	return trello.Card{}, fmt.Errorf("cannot find the card %v", card)
}

func (m PeopleInRole) Message() (string, error) {
	appKey := getRequiredEnv("TRELLO_KEY")
	token := getRequiredEnv("TRELLO_TOKEN")

	trello, err := trello.NewAuthClient(appKey, &token)
	if err != nil {
		return "", err
	}

	user, err := trello.Member(getRequiredEnv("TRELLO_USER"))
	if err != nil {
		return "", err
	}

	board, err := getBoard(user, getRequiredEnv("TRELLO_BOARD"))
	if err != nil {
		return "", err
	}

	listName := getOptionalEnv("TRELLO_LIST", getCurrentList())
	list, err := getList(&board, listName)
	if err != nil {
		return "", err
	}

	cardName := getOptionalEnv("TRELLO_CARD", "interrupt")
	card, err := getCard(&list, cardName)
	if err != nil {
		return "", err
	}

	var handles []string

	for _, memberID := range card.IdMembers {
		member, _ := trello.Member(memberID)
		if handle, ok := m.People[member.Username]; ok {
			handles = append(handles, fmt.Sprintf("<@%s>", handle))
		} else {
			handles = append(handles, member.FullName)
		}
	}

	if len(handles) == 0 {
		return "", nil
	}

	sort.Strings(handles)

	return strings.Join(handles, " "), nil
}
