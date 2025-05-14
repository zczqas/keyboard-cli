package game

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Challenge struct {
	Text            string
	CurrentPosition int
	StartTime       time.Time
	EndTime         time.Time
	Mistakes        int
	Completed       bool
}

// PracticeTexts contains sample texts for typing practice
// I might do something different in the future.
var PracticeTexts = []string{
	"The quick brown fox jumps over the lazy dog.",
	"Programming is the art of telling another human what one wants the computer to do.",
	"Simplicity is prerequisite for reliability.",
	"The best error message is the one that never shows up.",
	"First, solve the problem. Then, write the code.",
	"Any fool can write code that a computer can understand. Good programmers write code that humans can understand.",
	"Experience is the name everyone gives to their mistakes.",
	"Make it work, make it right, make it fast.",
	"Weeks of coding can save you hours of planning.",
	"Perfection is achieved not when there is nothing more to add, but rather when there is nothing more to take away.",
}

func NewChallenge(difficulty string) *Challenge {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	var text string
	switch strings.ToLower(difficulty) {
	case "easy":
		text = PracticeTexts[r.Intn(3)]
	case "hard":
		text = PracticeTexts[len(PracticeTexts)-3+r.Intn(3)]
	default:
		text = PracticeTexts[3+r.Intn(4)]
	}

	return &Challenge{
		Text:            text,
		CurrentPosition: 0,
		Mistakes:        0,
		Completed:       false,
	}
}

func (c *Challenge) Start() {
	c.StartTime = time.Now()
}

func (c *Challenge) ProcessKey(key string) bool {
	if c.Completed {
		return false
	}

	if c.CurrentPosition == 0 && c.StartTime.IsZero() {
		c.Start()
	}

	if c.CurrentPosition < len(c.Text) {
		expectedChar := string(c.Text[c.CurrentPosition])
		if key == expectedChar {
			c.CurrentPosition++

			if c.CurrentPosition >= len(c.Text) {
				c.EndTime = time.Now()
				c.Completed = true
			}
			return true
		} else {
			c.Mistakes++
			return false
		}
	}

	return false
}

func (c *Challenge) GetStats() map[string]interface{} {
	stats := make(map[string]interface{})
	stats["total_chars"] = len(c.Text)
	stats["typed_chars"] = c.CurrentPosition
	stats["mistakes"] = c.Mistakes
	stats["completed"] = c.Completed

	if c.CurrentPosition > 0 {
		accuracy := 100.0 * float64(c.CurrentPosition-c.Mistakes) / float64(c.CurrentPosition)
		stats["accuracy"] = accuracy
	} else {
		stats["accuracy"] = 100.0
	}

	if !c.StartTime.IsZero() {
		var duration time.Duration
		if c.Completed {
			duration = c.EndTime.Sub(c.StartTime)
		} else {
			duration = time.Since(c.StartTime)
		}

		minutes := duration.Minutes()
		if minutes > 0 {
			wpm := float64(c.CurrentPosition) / 5.0 / minutes
			stats["wpm"] = wpm
		} else {
			stats["wpm"] = 0.0
		}
	} else {
		stats["wpm"] = 0.0
	}

	return stats
}

func (c *Challenge) GetFormattedText() string {
	var result strings.Builder

	if c.CurrentPosition > 0 {
		result.WriteString(fmt.Sprintf("\x1b[32m%s\x1b[0m", c.Text[:c.CurrentPosition]))
	}

	if c.CurrentPosition < len(c.Text) {
		result.WriteString(fmt.Sprintf("\x1b[47m\x1b[30m%s\x1b[0m", string(c.Text[c.CurrentPosition])))
	}

	if c.CurrentPosition+1 < len(c.Text) {
		result.WriteString(c.Text[c.CurrentPosition+1:])
	}

	return result.String()
}

func (c *Challenge) GetProgressText() string {
	stats := c.GetStats()

	if c.Completed {
		return fmt.Sprintf("Completed! WPM: %.1f, Accuracy: %.1f%%",
			stats["wpm"].(float64),
			stats["accuracy"].(float64))
	}

	return fmt.Sprintf("Progress: %d/%d chars, WPM: %.1f, Accuracy: %.1f%%",
		c.CurrentPosition,
		len(c.Text),
		stats["wpm"].(float64),
		stats["accuracy"].(float64))
}
