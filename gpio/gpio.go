package gpio

import (
	lcd1602 "github.com/pimvanhespen/go-pi-lcd1602"
	"github.com/pimvanhespen/go-pi-lcd1602/synchronized"
	"time"
)

var (
	LCD_RS = 26
	LCD_E  = 19
	LCD_D4 = 13
	LCD_D5 = 6
	LCD_D6 = 5
	LCD_D7 = 11
	LED_ON = 15
)

func main() {
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
