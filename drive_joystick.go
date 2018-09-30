package gophercar

import (
	"sync/atomic"
	"time"

	"github.com/fogleman/gg"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/joystick"
	"gobot.io/x/gobot/platforms/raspi"
)

type JoystickValuePair struct {
	x float64
	y float64
}

func DriveJoystick(enableMpu6050 bool) {

	// joystick
	var leftX, leftY, rightX, rightY atomic.Value

	throttleZero := 350
	steering := 0.0

	r := raspi.NewAdaptor()
	pca9685 := i2c.NewPCA9685Driver(r)
	oled := i2c.NewSSD1306Driver(r)
	var mpu6050 *i2c.MPU6050Driver
	if enableMpu6050 {
		mpu6050 = i2c.NewMPU6050Driver(r)
	}

	joystickAdaptor := joystick.NewAdaptor()
	stick := joystick.NewDriver(joystickAdaptor, "dualshock3")

	ggCtx := gg.NewContext(oled.Buffer.Width, oled.Buffer.Height)

	work := func() {
		leftX.Store(float64(0.0))
		leftY.Store(float64(0.0))
		rightX.Store(float64(0.0))
		rightY.Store(float64(0.0))

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

		stick.On(joystick.LeftX, func(data interface{}) {
			val := float64(data.(int16))
			leftX.Store(val)
		})

		stick.On(joystick.LeftY, func(data interface{}) {
			val := float64(data.(int16))
			leftY.Store(val)
		})

		stick.On(joystick.RightX, func(data interface{}) {
			val := float64(data.(int16))
			rightX.Store(val)
		})

		stick.On(joystick.RightY, func(data interface{}) {
			val := float64(data.(int16))
			rightY.Store(val)
		})

		gobot.Every(10*time.Millisecond, func() {
			// right stick is steering
			rightStick := getRightStick(rightX, rightY)

			switch {
			case rightStick.x > 10:
				setSteering(pca9685, gobot.Rescale(rightStick.x, -32767.0, 32767.0, -1.0, 1.0))
			case rightStick.x < -10:
				setSteering(pca9685, gobot.Rescale(rightStick.x, -32767.0, 32767.0, -1.0, 1.0))
			default:
				setSteering(pca9685, 0)
			}
		})

		gobot.Every(10*time.Millisecond, func() {
			leftStick := getLeftStick(leftX, leftY)
			// left stick is throttle

			switch {
			case leftStick.y < -10:
				setThrottle(pca9685, gobot.Rescale(leftStick.y, -32767.0, 32767.0, -1.0, 1.0))
			case leftStick.y > 10:
				setThrottle(pca9685, gobot.Rescale(leftStick.y, -32767.0, 32767.0, -1.0, 1.0))
			default:
				setThrottle(pca9685, 0)
			}
		})
	}

	robot := gobot.NewRobot("gophercar",
		[]gobot.Connection{r, joystickAdaptor},
		[]gobot.Device{pca9685, oled, mpu6050, stick},
		work,
	)

	robot.Start()

}

func getLeftStick(leftX, leftY atomic.Value) JoystickValuePair {
	s := JoystickValuePair{x: 0, y: 0}
	s.x = leftX.Load().(float64)
	s.y = leftY.Load().(float64)
	return s
}

func getRightStick(rightX, rightY atomic.Value) JoystickValuePair {
	s := JoystickValuePair{x: 0, y: 0}
	s.x = rightX.Load().(float64)
	s.y = rightY.Load().(float64)
	return s
}
