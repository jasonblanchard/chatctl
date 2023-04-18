/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"github.com/tiktoken-go/tokenizer"
)

// tokenizeCmd represents the tokenize command
var tokenizeCmd = &cobra.Command{
	Use:   "tokenize",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		prompt, err := cmd.Flags().GetString("prompt")
		if err != nil {
			return fmt.Errorf("error getting prompt flag: %w", err)
		}

		justCount, err := cmd.Flags().GetBool("count")
		if err != nil {
			return fmt.Errorf("error getting count flag: %w", err)
		}

		messagesFile, err := cmd.Flags().GetString("file")
		if err != nil {
			return fmt.Errorf("error getting messages file: %w", err)
		}

		messages := []openai.ChatCompletionMessage{}

		if messagesFile != "" {
			bytes, err := os.ReadFile(messagesFile)
			if err != nil {
				return fmt.Errorf("error reading file: %w", err)
			}

			err = json.Unmarshal(bytes, &messages)
			if err != nil {
				return fmt.Errorf("error unmarshalling json: %w", err)
			}
		}

		if prompt != "" {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			})
		}

		stringToEncode := ""

		for _, message := range messages {
			stringToEncode = stringToEncode + message.Content
		}

		enc, err := tokenizer.Get(tokenizer.Cl100kBase)
		if err != nil {
			return fmt.Errorf("error creating tokenizer: %w", err)
		}

		ids, tokens, err := enc.Encode(stringToEncode)

		if err != nil {
			return fmt.Errorf("error encoding prompt: %e", err)
		}

		count := len(ids)

		if justCount {
			fmt.Println(count)
			return nil
		}

		colors := []func(a ...interface{}) string{
			color.New(color.FgWhite).SprintFunc(),
			color.New(color.FgCyan).SprintFunc(),
			color.New(color.FgGreen).SprintFunc(),
			color.New(color.FgYellow).SprintFunc(),
			color.New(color.FgMagenta).SprintFunc(),
		}

		fmt.Println()

		for i, token := range tokens {
			colorfn := colors[i%len(colors)]
			fmt.Print(colorfn(token))
		}

		fmt.Println()
		fmt.Printf("\n%v tokens\n", count)
		fmt.Println()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(tokenizeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tokenizeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tokenizeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	tokenizeCmd.Flags().StringP("prompt", "p", "", "Single user prompt for simple instructions")
	tokenizeCmd.Flags().BoolP("count", "c", false, "Output just the count")
	tokenizeCmd.Flags().StringP("file", "f", "", "Path to a file with message input")
}
