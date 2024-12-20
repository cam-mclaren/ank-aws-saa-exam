package main

import "fmt"
import "os"
import "log/slog"
import "bytes"
import "regexp"
import "strconv"

var logger *slog.Logger

type MultipleChoiceCard struct {
	Question string
	A        string
	B        string
	C        string
	D        string
	Correct  string
	Tags     string
}

func (m *MultipleChoiceCard) GetOption(option string, buff *bytes.Buffer) error {

	var update *string
	switch option {
	case `a\)`:
		update = &(m.A)

	case `b\)`:
		update = &(m.B)

	case `c\)`:
		update = &(m.C)

	case `d\)`:
		update = &(m.D)

	case `Question\:`:
		update = &(m.Question)

	case `Correct\:`:
		update = &(m.Correct)
	case `Tags\:`:
		update = &(m.Tags)
	default:
		return fmt.Errorf("invalid option: %s\n", option)
	}

	re_pattern, err := regexp.Compile(`(?m)^\s*` + option + `(.*)\n`)
	if err != nil {
		return err
	}
	result := re_pattern.FindSubmatch(buff.Bytes())
	*update = string(result[1])
	logger.Info(string(result[1]))
	return nil
}

func main() {
	log_file, err := os.OpenFile("./Card_composer.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		fmt.Printf("Failed to open log file\n")
		os.Exit(1)
	}

	log_options := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	logger = slog.New(slog.NewTextHandler(log_file, log_options))

	slog.SetLogLoggerLevel(slog.LevelInfo)

	deck_file_path := "/home/cam/Projects/My-SA-Deck/My-SA-Deck.csv"
	filepath := os.Args[1]
	card_file_fd, err := os.OpenFile(filepath, os.O_RDONLY, 0777)
	if err != nil {
		logger.Error("Failed to open card file", "error", err.Error())
		os.Exit(1)
	}
	defer card_file_fd.Close()

	card_buf := new(bytes.Buffer)
	num, err := card_buf.ReadFrom(card_file_fd)
	slog.Debug("Read data from card file", "bytes read", num)
	if err != nil {
		logger.Error("Failed to read card file data", "error", err.Error())
		os.Exit(1)
	}
	card := MultipleChoiceCard{}
	options := []string{"Question\\:", `a\)`, `b\)`, `c\)`, `d\)`, `Correct\:`, `Tags\:`}
	for _, option := range options {
		err := card.GetOption(option, card_buf)
		if err != nil {
			logger.Error("GetOption returned an error", "error", err)
			logger.Error("Failed to find option", "option", option)
			os.Exit(1)
		}
	}
	// Open file for reading
	deck_buffer := new(bytes.Buffer)
	deck_file_fd, err := os.OpenFile(deck_file_path, os.O_RDONLY, 0777)
	if err != nil {
		logger.Error("Could not open file", "file", deck_file_path)
		os.Exit(1)
	}

	// Read entire file content
	num, err = deck_buffer.ReadFrom(deck_file_fd)
	slog.Debug("Read data from deck file", "bytes read", num)
	if err != nil {
		logger.Error("Failed to read deck file data", "error", err.Error())
		os.Exit(1)
	}

	// Close file for now
	deck_file_fd.Close()

	fileContent := deck_buffer.String()

	pat, err := regexp.Compile(`(?m)^#number_of_cards\:([0-9]*)`)
	if err != nil {
		logger.Error("Could not compile regex", "error", err.Error())
		os.Exit(1)
	}

	match := pat.FindStringSubmatch(fileContent)
	number_string := match[1]
	number, err := strconv.Atoi(number_string)
	if err != nil {
		logger.Error("Could na", "error", err.Error())
		os.Exit(1)
	}
	number++
	logger.Info("Incrementing number of cards", "number_of_cards", fmt.Sprintf("%d", number))
	ReplacePattern, _ := regexp.Compile(`(?m)#number_of_cards:.*\n`)
	fileContent = ReplacePattern.ReplaceAllString(fileContent, fmt.Sprintf("#number_of_cards:%d\n", number))
	newline := fmt.Sprintf("Card%03d;Cams cloze+;%s<br>a)%s<br>b)%s<br>c)%s<br>d)%s<br>Correct:{{c1::%s}}<br>;;%s\n", number, card.Question, card.A, card.B, card.C, card.D, card.Correct, card.Tags)
	fileContent = fileContent + newline
	fmt.Println(newline)
	deck_file_fd, err = os.OpenFile(deck_file_path, os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		logger.Error("Could not open file", "file", deck_file_path)
		os.Exit(1)
	}
	num_written, err := deck_file_fd.WriteString(fileContent)
	if err != nil {
		logger.Error("Could not write to  file", "file", deck_file_path, "error", err.Error())
		os.Exit(1)
	}
	logger.Info("Wrote data to file", "bytes written", fmt.Sprintf("%d", num_written))
}
