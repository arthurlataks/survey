package survey

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/AlecAivazis/survey/v2/core"
	"github.com/AlecAivazis/survey/v2/terminal"
	expect "github.com/Netflix/go-expect"
	"github.com/stretchr/testify/assert"
)

func init() {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true
}

func TestConfirmRender(t *testing.T) {

	tests := []struct {
		title    string
		prompt   Confirm
		data     ConfirmTemplateData
		expected string
	}{
		{
			"Test Confirm question output with default true",
			Confirm{Message: "Is pizza your favorite food?", Default: true},
			ConfirmTemplateData{},
			fmt.Sprintf("%s Is pizza your favorite food? (Y/n) ", defaultAskOptions().PromptConfig.Icons.Question),
		},
		{
			"Test Confirm question output with default false",
			Confirm{Message: "Is pizza your favorite food?", Default: false},
			ConfirmTemplateData{},
			fmt.Sprintf("%s Is pizza your favorite food? (y/N) ", defaultAskOptions().PromptConfig.Icons.Question),
		},
		{
			"Test Confirm answer output",
			Confirm{Message: "Is pizza your favorite food?"},
			ConfirmTemplateData{Answer: "Yes"},
			fmt.Sprintf("%s Is pizza your favorite food? Yes\n", defaultAskOptions().PromptConfig.Icons.Question),
		},
		{
			"Test Confirm with help but help message is hidden",
			Confirm{Message: "Is pizza your favorite food?", Help: "This is helpful"},
			ConfirmTemplateData{},
			fmt.Sprintf("%s Is pizza your favorite food? [%s for help] (y/N) ", defaultAskOptions().PromptConfig.Icons.Question, string(defaultAskOptions().PromptConfig.HelpInput)),
		},
		{
			"Test Confirm help output with help message shown",
			Confirm{Message: "Is pizza your favorite food?", Help: "This is helpful"},
			ConfirmTemplateData{ShowHelp: true},
			fmt.Sprintf("%s This is helpful\n%s Is pizza your favorite food? (y/N) ", defaultAskOptions().PromptConfig.Icons.Help, defaultAskOptions().PromptConfig.Icons.Question),
		},
	}

	for _, test := range tests {
		r, w, err := os.Pipe()
		assert.Nil(t, err, test.title)

		test.prompt.WithStdio(terminal.Stdio{Out: w})
		test.data.Confirm = test.prompt

		// set the runtime config
		test.data.Config = &defaultAskOptions().PromptConfig

		err = test.prompt.Render(
			ConfirmQuestionTemplate,
			test.data,
		)
		assert.Nil(t, err, test.title)

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)

		assert.Contains(t, buf.String(), test.expected, test.title)
	}
}

func TestConfirmPrompt(t *testing.T) {
	tests := []PromptTest{
		{
			"Test Confirm prompt interaction",
			&Confirm{
				Message: "Is pizza your favorite food?",
			},
			func(c *expect.Console) {
				c.ExpectString("Is pizza your favorite food? (y/N)")
				c.SendLine("n")
				c.ExpectEOF()
			},
			false,
		},
		{
			"Test Confirm prompt interaction with default",
			&Confirm{
				Message: "Is pizza your favorite food?",
				Default: true,
			},
			func(c *expect.Console) {
				c.ExpectString("Is pizza your favorite food? (Y/n)")
				c.SendLine("")
				c.ExpectEOF()
			},
			true,
		},
		{
			"Test Confirm prompt interaction overriding default",
			&Confirm{
				Message: "Is pizza your favorite food?",
				Default: true,
			},
			func(c *expect.Console) {
				c.ExpectString("Is pizza your favorite food? (Y/n)")
				c.SendLine("n")
				c.ExpectEOF()
			},
			false,
		},
		{
			"Test Confirm prompt interaction and prompt for help",
			&Confirm{
				Message: "Is pizza your favorite food?",
				Help:    "It probably is",
			},
			func(c *expect.Console) {
				c.ExpectString(
					fmt.Sprintf(
						"Is pizza your favorite food? [%s for help] (y/N)",
						string(defaultAskOptions().PromptConfig.HelpInput),
					),
				)
				c.SendLine("?")
				c.ExpectString("It probably is")
				c.SendLine("Y")
				c.ExpectEOF()
			},
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunPromptTest(t, test)
		})
	}
}
