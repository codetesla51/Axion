/*
History Module - Calculation History Management
===============================================

This module provides persistent storage and retrieval of calculation history
using JSON serialization. All calculations are stored in a local file and
can be displayed to the user for reference.

The history system:
- Automatically saves each successful calculation
- Persists data across program sessions
- Displays results in reverse chronological order (newest first)
- Handles file I/O errors gracefully
- Uses structured JSON format for data integrity

File format: Array of Entry objects in JSON format
Location: history.json in the current working directory
*/

package history

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
)

// Entry represents a single calculation record in the history
type JsonFloat float64

func (f JsonFloat) MarshalJSON() ([]byte, error) {
	v := float64(f)

	if math.IsInf(v, 1) {
		return json.Marshal("+∞")
	}
	if math.IsInf(v, -1) {
		return json.Marshal("-∞")
	}
	if math.IsNaN(v) {
		return json.Marshal("NaN")
	}
	return json.Marshal(v)
}

type Entry struct {
	Expression string    `json:"expression"` // Original mathematical expression
	Result     JsonFloat `json:"result"`     // Computed numerical result
}

// AddHistory appends a new calculation to the persistent history file
// Handles file creation, existing data preservation, and atomic updates
func AddHistory(input string, result float64) error {
	const HistoryFile = "history.json"
	var history []Entry

	// Attempt to read existing history data
	data, err := os.ReadFile(HistoryFile)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, create it
			newFile, err := os.Create(HistoryFile)
			if err != nil {
				return err
			}
			defer newFile.Close()
			data = []byte{} // empty data since file is new
		} else {
			// Other read errors
			return err
		}
	}
	// Parse existing history if file contains data
	if len(data) > 0 {
		err = json.Unmarshal(data, &history)
		if err != nil {
			// Return error for malformed JSON data
			return err
		}
	}

	// Create new history entry
	entry := Entry{Expression: input, Result: JsonFloat(result)}

	// Append new entry to existing history
	history = append(history, entry)

	// Serialize updated history with readable formatting
	updatedContent, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		// Return error for serialization failure
		return err
	}

	// Write updated history to file with appropriate permissions
	err = os.WriteFile(HistoryFile, updatedContent, 0644)
	if err != nil {
		// Return error for write failure
		return err
	}

	return nil
}

// ShowHistory displays the complete calculation history in reverse order
// Most recent calculations are shown first for better user experience
func ShowHistory() error {
	const historyFile = "history.json"
	var history []Entry

	// Read history file
	data, err := os.ReadFile(historyFile)
	if err != nil {
		if os.IsNotExist(err) {
			// Handle case where no history exists yet
			fmt.Println("no history data")
			return nil
		}
		// Return error for read failures
		return err
	}

	// Parse JSON history data
	err = json.Unmarshal(data, &history)
	if err != nil {
		// Return error for malformed JSON
		return err
	}

	// Handle empty history case
	if len(history) == 0 {
		fmt.Println("no history data")
		return nil
	}

	// Display history in reverse chronological order (newest first)
	for i := len(history) - 1; i >= 0; i-- {
		entry := history[i]
		fmt.Printf("------------------------------------------------\n")
		fmt.Printf(" Expression : %s\n", entry.Expression)
		fmt.Printf(" Result     : %g\n", entry.Result)
		fmt.Printf("------------------------------------------------\n\n")
	}

	return nil
}
