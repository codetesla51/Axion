package history

import (
	"encoding/json"
	"fmt"
	"os"
)

type Entry struct {
	Expression string  `json:"expression"`
	Result     float64 `json:"result"`
}

func AddHistory(input string, result float64) error {
	HistoryFile := "history.json"
	var history []Entry

	data, err := os.ReadFile(HistoryFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if len(data) > 0 {
		err = json.Unmarshal(data, &history)
		if err != nil {
			return err

		}
	}
	entry := Entry{Expression: input, Result: result}
	history = append(history, entry)
	updatedContent, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err

	}

	err = os.WriteFile(HistoryFile, updatedContent, 0644)
	if err != nil {
		return err

	}
	return nil
}
func ShowHistory() error {
	historyFile := "history.json"
	var history []Entry

	data, err := os.ReadFile(historyFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("no history data")
			return nil
		}
		return err
	}

	err = json.Unmarshal(data, &history)
	if err != nil {
		return err
	}

	if len(history) == 0 {
		fmt.Println("no history data")
		return nil
	}
	for i := len(history) - 1; i >= 0; i-- {
		entry := history[i]
		fmt.Printf("------------------------------------------------\n")
		fmt.Printf(" Expression : %s\n", entry.Expression)
		fmt.Printf(" Result     : %g\n", entry.Result)
		fmt.Printf("------------------------------------------------\n\n")
	}

	return nil
}
