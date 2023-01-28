package twitch // mudar depois para o package message

import (
	"strings"
)

type Message struct {
	Source     Source  `json:"source"`
	Command    Command `json:"command"`
	Parameters string  `json:"parameters"`
	Message    string  `json:"message"`
	Author     string  `json:"author"`
	Color      string  `json:"color"`
	Emotes     []Emote `json:"emotes"`
	Badges     []Badge `json:"badges"`
	EmoteSets  []Emote `json:"emotes-sets"`
	Tags       []Tag   `json:"tags"`
}

type Command struct {
	Command             string `json:"command"`
	Channel             string `json:"channel"`
	IsCapRequestEnabled bool
}

type Badge struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Emote struct {
	Id        string `json:"id"`
	Positions []TextPosition
}

type TextPosition struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type Tag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Source struct {
	Nick string `json:"nick"`
	Host string `json:"host"`
}

// func main() {
// Parse("@badges=staff/1,broadcaster/1,turbo/1;color=#FF0000;display-name=PetsgomOO;emote-only=1;emotes=33:0-7;flags=0-7:A.6/P.6,25-36:A.1/I.2;id=c285c9ed-8b1b-4702-ae1c-c64d76cc74ef;mod=0;room-id=81046256;subscriber=0;turbo=0;tmi-sent-ts=1550868292494;user-id=81046256;user-type=staff :petsgomoo!petsgomoo@petsgomoo.tmi.twitch.tv PRIVMSG #petsgomoo :DansGame")
// Parse("PING :tmi.twitch.tv")
// Parse(":lovingt3s!lovingt3s@lovingt3s.tmi.twitch.tv PRIVMSG #lovingt3s :!dilly")
// }

func parse(message string) Message {
	index := 0
	endIndex := 0

	parsed := Message{}
	var rawTags string
	var rawSource string
	var rawCommand string
	var rawParameters string

	if message[index] == '@' {
		endIndex = strings.Index(message, " ")
		rawTags = message[1:endIndex]
		message = message[endIndex+1:]
	}

	if message[index] == ':' {
		endIndex = strings.Index(message, " ")
		rawSource = message[1:endIndex]
		message = message[endIndex+1:]
	}

	endIndex = strings.Index(message, ":")

	if endIndex == -1 {
		endIndex = len(message)
	}

	rawCommand = strings.Trim(message[index:endIndex], "")

	if endIndex != len(message) {
		endIndex++
		rawParameters = message[endIndex:]
	}
	parsed.Message = message[endIndex:]

	command := parseCommand(rawCommand)

	if command.Command == "" {
		return Message{}
	}

	parsed.Command = command

	if rawTags != "" {
		badges := parseBadges(rawTags)
		parsed.Badges = append(parsed.Badges, badges...)

		emotes := parseEmotes(rawTags)
		parsed.Emotes = append(parsed.Emotes, emotes...)

		emoteSets := parseEmoteSets(rawTags)
		parsed.EmoteSets = append(parsed.EmoteSets, emoteSets...)

		tags, author, color := parseTags(rawTags)
		parsed.Tags = append(parsed.Tags, tags...)
		parsed.Author = author
		parsed.Color = color
	}

	parsed.Source = parseSource(rawSource)

	parsed.Parameters = strings.TrimSpace(rawParameters)

	return parsed
}

func parseCommand(rawCommand string) Command {
	commandParts := strings.Split(rawCommand, " ")
	command := commandParts[0]
	secondCommand := commandParts[1]

	switch command {
	case "JOIN", "PART", "NOTICE", "CLEARCHAT", "HOSTTARGET", "PRIVMSG", "USERSTATE", "ROOMSTATE", "001":
		return Command{Command: command, Channel: secondCommand}
	case "PING", "GLOBALUSERSTATE", "RECONNECT":
		return Command{Command: command}
	case "CAP":
		return Command{Command: command, IsCapRequestEnabled: secondCommand == "ACK"}
	}

	return Command{}
}

func parseBadges(rawTags string) []Badge {
	parsedTags := strings.Split(rawTags, ";")

	badges := []Badge{}

	for _, t := range parsedTags {
		tag := strings.Split(t, "=")
		tagName := tag[0]
		tagValue := tag[1]

		if tagValue == "" {
			continue
		}

		if tagName != "badges" && tagName != "badge-info" {
			continue
		}

		rawBadges := strings.Split(tagValue, ",")
		for _, b := range rawBadges {
			badge := strings.Split(b, "/")
			badges = append(badges, Badge{Name: badge[0], Value: badge[1]})
		}
	}

	return badges
}

func parseEmotes(rawTags string) []Emote {
	parsedTags := strings.Split(rawTags, ";")

	emotes := []Emote{}

	for _, t := range parsedTags {
		tag := strings.Split(t, "=")
		tagName := tag[0]
		tagValue := tag[1]

		if tagValue == "" {
			continue
		}

		if tagName != "emotes" {
			continue
		}

		rawEmotes := strings.Split(tagValue, "/")
		for _, e := range rawEmotes {
			positions := []TextPosition{}

			emoteParts := strings.Split(e, ":")
			rawPositions := strings.Split(emoteParts[1], ",")
			for _, p := range rawPositions {
				positionParts := strings.Split(p, "-")
				positions = append(positions, TextPosition{
					Start: positionParts[0],
					End:   positionParts[1],
				})
			}
			emotes = append(emotes, Emote{
				Id:        emoteParts[0],
				Positions: positions})
		}
	}

	return emotes
}

func parseEmoteSets(rawTags string) []Emote {
	parsedTags := strings.Split(rawTags, ";")

	emotes := []Emote{}

	for _, t := range parsedTags {
		tag := strings.Split(t, "=")
		tagName := tag[0]
		tagValue := tag[1]

		if tagValue == "" {
			continue
		}

		if tagName != "emote-sets" {
			continue
		}

		rawSet := strings.Split(tagValue, ",")
		for _, s := range rawSet {
			emotes = append(emotes, Emote{Id: s})
		}
	}

	return emotes
}

func parseTags(rawTags string) (tags []Tag, displayName string, color string) {
	parsedTags := strings.Split(rawTags, ";")

	tags = []Tag{}

	tagsToIgnore := []string{"client-nonce", "flags", "badges", "badges-info", "emotes", "emote-sets"}

	for _, t := range parsedTags {
		tag := strings.Split(t, "=")
		tagName := tag[0]
		tagValue := tag[1]

		if tagValue == "" {
			continue
		}

		if strings.Contains(strings.Join(tagsToIgnore, " "), tagName) {
			continue
		}

		if tagName == "display-name" {
			displayName = tagValue
		}

		if tagName == "color" {
			color = tagValue
		}

		tags = append(tags, Tag{Name: tagName, Value: tagValue})
	}

	return tags, displayName, color
}

func parseSource(rawSource string) Source {
	if rawSource == "" {
		return Source{}
	}

	sourceParts := strings.Split(rawSource, "!")

	if len(sourceParts) != 2 {
		return Source{Host: sourceParts[0]}
	}

	return Source{Nick: sourceParts[0], Host: sourceParts[1]}
}
