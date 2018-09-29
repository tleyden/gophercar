package gophercar

import (
	"time"

	"fmt"
	"math"

	"github.com/fogleman/gg"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/keyboard"
	"gobot.io/x/gobot/platforms/raspi"
)

// Drive the car via keyboard control
func DriveKeyboard(enableMpu6050 bool) {

	r := raspi.NewAdaptor()
	pca9685 := i2c.NewPCA9685Driver(r)
	oled := i2c.NewSSD1306Driver(r)

	var mpu6050 *i2c.MPU6050Driver
	if enableMpu6050 {
		mpu6050 = i2c.NewMPU6050Driver(r)

	}
	keys := keyboard.NewDriver()

	steering := 0.0
	throttleZero := 350
	throttlePower := 0.25

	ggCtx := gg.NewContext(oled.Buffer.Width, oled.Buffer.Height)

	work := func() {
		gobot.Every(1*time.Second, func() {
			handleOLED(ggCtx, steering, oled)
		})

		gobot.Every(100*time.Millisecond, func() {
			if enableMpu6050 {
				handleAccel(mpu6050)
			}
		})

		// init the PWM controller
		pca9685.SetPWMFreq(60)

		// init the ESC controller for throttle zero
		pca9685.SetPWM(0, 0, uint16(throttleZero))

		keys.On(keyboard.Key, func(data interface{}) {
			key := data.(keyboard.KeyEvent)

			switch key.Key {
			case keyboard.ArrowUp:
				setThrottle(pca9685, throttlePower)

				gobot.After(1*time.Second, func() {
					setThrottle(pca9685, 0)
				})
			case keyboard.ArrowDown:
				setThrottle(pca9685, -1*throttlePower)

				gobot.After(1*time.Second, func() {
					setThrottle(pca9685, 0)
				})
			case keyboard.ArrowRight:
				if steering < 1.0 {
					steering = round(steering+0.1, 0.05)
				}

				setSteering(pca9685, steering)
			case keyboard.ArrowLeft:
				if round(steering, 0.05) > -1.0 {
					steering = round(steering-0.1, 0.05)
				}

				setSteering(pca9685, steering)
			}
		})
	}

	robot := gobot.NewRobot("gophercar",
		[]gobot.Connection{r},
		[]gobot.Device{pca9685, oled, keys},
		work,
	)

	robot.Start()

}

func handleOLED(ggCtx *gg.Context, steering float64, oled *i2c.SSD1306Driver) {
	ggCtx.SetRGB(0, 0, 0)
	ggCtx.Clear()
	ggCtx.SetRGB(1, 1, 1)
	ggCtx.DrawStringAnchored(time.Now().Format("15:04:05"), 0, 0, 0, 1)

	ggCtx.DrawStringAnchored(fmt.Sprint("Steering: ", steering), 0, 32, 0, 1)
	oled.ShowImage(ggCtx.Image())
}

func handleAccel(mpu6050 *i2c.MPU6050Driver) {

	mpu6050.GetData()

	fmt.Println("Accelerometer", mpu6050.Accelerometer)
	fmt.Println("Gyroscope", mpu6050.Gyroscope)
	fmt.Println("Temperature", mpu6050.Temperature)
}

func setSteering(pca9685 *i2c.PCA9685Driver, steering float64) {
	steeringVal := getSteeringPulse(steering)
	pca9685.SetPWM(1, 0, uint16(steeringVal))
}

func setThrottle(pca9685 *i2c.PCA9685Driver, throttle float64) {
	throttleVal := getThrottlePulse(throttle)
	pca9685.SetPWM(0, 0, uint16(throttleVal))
}

// adjusts the steering from -1.0 (hard left) <-> 1.0 (hardright) to the correct
// pwm pulse values.
func getSteeringPulse(val float64) float64 {
	return gobot.Rescale(val, -1, 1, 290, 490)
}

// adjusts the throttle from -1.0 (hard back) <-> 1.0 (hard forward) to the correct
// pwm pulse values.
func getThrottlePulse(val float64) int {
	if val > 0 {
		return int(gobot.Rescale(val, 0, 1, 350, 300))
	}
	return int(gobot.Rescale(val, -1, 0, 490, 350))
}

func round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}
