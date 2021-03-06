/*
 * Copyright 2020 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package commands

import (
	"fmt"
	"github.com/nalej/public-api/internal/app/options"
	"github.com/spf13/cobra"
)

var key string
var value string

var optionsCmd = &cobra.Command{
	Use:     "option",
	Aliases: []string{"options", "opt"},
	Short:   "Manage default options",
	Long:    `Manage default values for the commands parameters`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

func init() {
	optionsCmd.PersistentFlags().StringVar(&key, "key", "", "Specify the key")
	optionsCmd.PersistentFlags().StringVar(&value, "value", "", "Specify the value")
	rootCmd.AddCommand(optionsCmd)
	optionsCmd.AddCommand(setOptionCmd)
	optionsCmd.AddCommand(getOptionCmd)
	optionsCmd.AddCommand(deleteOptionCmd)
	optionsCmd.AddCommand(listOptionsCmd)
}

var setOptionCmd = &cobra.Command{
	Use:   "set",
	Short: "Set the value for a given key",
	Long:  `Set the value for a given key`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		opts := options.NewOptions()
		opts.Set(key, value)
	},
}

var getOptionCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"get"},
	Short:   "Get the value for a given key",
	Long:    `Get the value for a given key`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		opts := options.NewOptions()
		retrieved := opts.Get(key)
		fmt.Printf("Key: %s Value: %s\n", key, retrieved)
	},
}

var listOptionsCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List the stored values",
	Long:    `List the stored values`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		opts := options.NewOptions()
		opts.List()
	},
}

var deleteOptionCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"remove", "del", "rm"},
	Short:   "Delete a given key",
	Long:    `Delete a given key`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		opts := options.NewOptions()
		opts.Delete(key)
		fmt.Printf("Key: %s has been deleted\n", key)
	},
}
