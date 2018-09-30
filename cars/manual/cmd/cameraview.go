// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/hybridgroup/gophercar"
	"github.com/spf13/cobra"
)

// cameraviewCmd represents the cameraview command
var cameraviewCmd = &cobra.Command{
	Use:   "cameraview",
	Short: "Connect to the onboard camera and stream video to the browser.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		gophercar.MpegStream(cameraId, streamListenUrl)

	},
}

func init() {

	rootCmd.AddCommand(cameraviewCmd)

	cameraviewCmd.PersistentFlags().IntVar(
		&cameraId,
		"camera-id",
		0,
		"The camera id, eg, 0",
	)

	cameraviewCmd.PersistentFlags().StringVar(
		&streamListenUrl,
		"stream-listen-url",
		"0.0.0.0:8080",
		"The interface and port to listen on to stream the video",
	)

}
