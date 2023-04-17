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
	"text/tabwriter"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// moderationsCmd represents the moderations command
var moderationsCmd = &cobra.Command{
	Use:   "moderations",
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

		outputJson, err := cmd.Flags().GetBool("json")
		if err != nil {
			return fmt.Errorf("error getting json: %w", err)
		}

		s := spinner.New(spinner.CharSets[26], 250*time.Millisecond)
		s.Start()

		resp, err := client.Moderations(
			context.Background(),
			openai.ModerationRequest{
				Input: prompt,
			},
		)

		s.Stop()

		if err != nil {
			return fmt.Errorf("error getting moderation: %w", err)
		}

		if outputJson {
			bytes, err := json.Marshal(resp)
			if err != nil {
				return fmt.Errorf("error marshalling json: %w", err)
			}
			fmt.Println(string(bytes))
		}

		result := resp.Results[0]

		type ResultMap struct {
			Name      string
			IsFlagged bool
			Score     float32
		}

		displayResults := []ResultMap{
			{
				Name:      "hate",
				IsFlagged: result.Categories.Hate,
				Score:     result.CategoryScores.Hate,
			},
			{
				Name:      "hate/threatening",
				IsFlagged: result.Categories.HateThreatening,
				Score:     result.CategoryScores.HateThreatening,
			},
			{
				Name:      "self-harm",
				IsFlagged: result.Categories.SelfHarm,
				Score:     result.CategoryScores.SelfHarm,
			},
			{
				Name:      "sexual",
				IsFlagged: result.Categories.Sexual,
				Score:     result.CategoryScores.Sexual,
			},
			{
				Name:      "sexual/minors",
				IsFlagged: result.Categories.SexualMinors,
				Score:     result.CategoryScores.SexualMinors,
			},
			{
				Name:      "violence",
				IsFlagged: result.Categories.Violence,
				Score:     result.CategoryScores.Violence,
			},
			{
				Name:      "violence/graphic",
				IsFlagged: result.Categories.ViolenceGraphic,
				Score:     result.CategoryScores.ViolenceGraphic,
			},
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

		for _, displayResult := range displayResults {
			gray := color.New(color.FgHiBlack).SprintFunc()
			green := color.New(color.FgGreen).SprintFunc()
			red := color.New(color.FgRed).SprintFunc()
			var isFlagged string
			var score string

			if displayResult.IsFlagged {
				isFlagged = red("true")
				score = red(displayResult.Score)
			} else {
				isFlagged = green("false")
				score = gray(displayResult.Score)
			}

			fmt.Fprintf(w, "%v\t%v\t%v\n", displayResult.Name, isFlagged, score)
		}

		fmt.Println()
		w.Flush()
		fmt.Println()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(moderationsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// moderationsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// moderationsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	moderationsCmd.Flags().StringP("prompt", "p", "", "Single user prompt for simple instructions")
	moderationsCmd.Flags().BoolP("json", "j", false, "Output result as JSON")
}
