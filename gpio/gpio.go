package gpio

import (
	lcd1602 "github.com/pimvanhespen/go-pi-lcd1602"
	"github.com/pimvanhespen/go-pi-lcd1602/synchronized"
	"time"
)

// Set the values of GPIO pins for the LCD
var (
	LCD_RS = 26 // Register Select
	LCD_E  = 19 // Enable
	LCD_D4 = 13 // Data 4
	LCD_D5 = 6  // Data 5
	LCD_D6 = 5  // Data 6
	LCD_D7 = 11 // Data 7
	LED_ON = 15 // LED on/off (not used)
)

func LCD() {
	lcdi := lcd1602.New(
		LCD_RS,                                //rs
		LCD_E,                                 //enable
		[]int{LCD_D4, LCD_D5, LCD_D6, LCD_D7}, //datapins
		16,                                    //lineSize
	)
	lcd := synchronized.NewSynchronizedLCD(lcdi)
	lcd.Initialize()
	lcd.WriteLines("This is working", "I hope! - Emily")
	time.Sleep(1 * time.Second)
	lcd.WriteLines("I guess", "we'll see...")
	lcd.Clear()
	lcd.Close()
}

func Rotary() {

}
