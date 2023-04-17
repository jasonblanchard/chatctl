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
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		key := viper.GetString("key")
		client := openai.NewClient(key)

		prompt, err := cmd.Flags().GetString("prompt")
		if err != nil {
			return fmt.Errorf("error getting prompt: %w", err)
		}

		useJson, err := cmd.Flags().GetBool("json")
		if err != nil {
			return fmt.Errorf("error getting json: %w", err)
		}

		messages := []openai.ChatCompletionMessage{}

		messagesFile, err := cmd.Flags().GetString("messages")
		if err != nil {
			return fmt.Errorf("error getting messages file: %w", err)
		}

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

		s := spinner.New(spinner.CharSets[26], 250*time.Millisecond)
		s.Start()

		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:       openai.GPT3Dot5Turbo,
				Temperature: 0.8,
				Messages:    messages,
			},
		)

		s.Stop()

		if err != nil {
			return err
		}

		if useJson {
			bytes, err := json.Marshal(resp)
			if err != nil {
				return fmt.Errorf("error marshalling json: %w", err)
			}
			fmt.Println(string(bytes))
		} else {
			fmt.Println(resp.Choices[0].Message.Content)
		}

		// TODO: Allow writing back to file to make it work like a chat

		return nil
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chatCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	chatCmd.Flags().BoolP("json", "j", false, "Output result as JSON")
	chatCmd.Flags().StringP("prompt", "p", "", "Single user prompt for simple instructions")
	chatCmd.Flags().StringP("messages", "m", "", "Path to a file with message input")
}
