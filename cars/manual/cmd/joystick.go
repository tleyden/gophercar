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
	"fmt"

	"github.com/spf13/cobra"
	"github.com/hybridgroup/gophercar"
)

// joystickCmd represents the joystick command
var joystickCmd = &cobra.Command{
	Use:   "joystick",
	Short: "Manually drive the car via a joystick, with optional camera view",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {

		if enableCamera {
			go gophercar.MpegStream(cameraId, streamListenUrl)
		}

		gophercar.DriveJoystick(enableMpu6050)
	},
}

func init() {
	rootCmd.AddCommand(joystickCmd)

	joystickCmd.PersistentFlags().BoolVarP(
		&enableCamera,
		"enable-camera",
		"c",
		false,
		"Enable camera streaming.  If true, must provide other camera related flags.",
	)

	joystickCmd.PersistentFlags().IntVar(
		&cameraId,
		"camera-id",
		0,
		"The camera id, eg, 0",
	)

	joystickCmd.PersistentFlags().StringVar(
		&streamListenUrl,
		"stream-listen-url",
		"0.0.0.0:8080",
		"The interface and port to listen on to stream the video",
	)

	joystickCmd.PersistentFlags().BoolVarP(
		&enableMpu6050,
		"enable-mpu6050",
		"m",
		false,
		"Enable mpu6050 accelerometer/gyroscope.  ",
	)

}
