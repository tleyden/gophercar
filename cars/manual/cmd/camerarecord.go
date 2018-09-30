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
	"time"
	"gocv.io/x/gocv"
)


var (
	targetRecordFile string
	recordNumFrames int
)

// camerarecordCmd represents the camerarecord command
var camerarecordCmd = &cobra.Command{
	Use:   "camerarecord",
	Short: "Record from the camera into a file",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {

		// open webcam
		webcam, err := gocv.OpenVideoCapture(cameraId)
		if err != nil {
			fmt.Printf("Error opening capture device: %v\n", cameraId)
			return
		}
		defer webcam.Close()

		img := gocv.NewMat()
		defer img.Close()

		if ok := webcam.Read(&img); !ok {
			fmt.Printf("Cannot read device %v\n", cameraId)
			return
		}

		writer, err := gocv.VideoWriterFile(targetRecordFile, "MJPG", 25, img.Cols(), img.Rows(), true)
		if err != nil {
			fmt.Printf("error opening video writer device: %v\n", targetRecordFile)
			return
		}
		defer writer.Close()

		fmt.Printf("Recording...")

		numFramesRead := 0
		for {

			if recordNumFrames != 0 && numFramesRead >= recordNumFrames {
				break
			}
			numFramesRead += 1

			if ok := webcam.Read(&img); !ok {
				fmt.Printf("Device closed: %v\n", cameraId)
				return
			}
			if img.Empty() {
				continue
			}

			writer.Write(img)
		}

	},
}



func init() {
	rootCmd.AddCommand(camerarecordCmd)

	camerarecordCmd.PersistentFlags().IntVar(
		&cameraId,
		"camera-id",
		0,
		"The camera id, eg, 0",
	)

	camerarecordCmd.PersistentFlags().StringVar(
		&targetRecordFile,
		"target-record-file",
		fmt.Sprintf("camera-%s.mjpg", time.Now()),
		"The target file where to save the mjpg video",
	)

	camerarecordCmd.PersistentFlags().IntVar(
		&recordNumFrames,
		"num-frames",
		0,
		"The number of frames to record, or 0 for unlimited",
	)


}
