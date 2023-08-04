package main

import (
	utils "flash_telink/Utils"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/gousb"
)

var dev *gousb.Device
var lastMessage []byte

func driverFind() []*gousb.Device {
	ctx := gousb.NewContext()
	devices, _ := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		return desc.Vendor == gousb.ID(0x248a) && (desc.Product == gousb.ID(0x5320) || desc.Product == gousb.ID(0x8266))
	})
	return devices
}

func driverInit(device int) {
	devices := driverFind()
	if len(devices) == 0 {
		fmt.Println("Device not found")
	}

	if device >= len(devices) {
		fmt.Println("Device not found")
	}

	dev = devices[device]

	// Detach any kernel driver that might be bound to the interface.
	err := dev.SetAutoDetach(true)
	if err != nil {
		fmt.Println("Error setting auto detach:", err)
		return
	}

	cfg, err := dev.Config(1) // changed to dev.Config(1) from dev.Config(0)
	if err != nil {
		// Check if the error message contains the available configurations
		if strings.Contains(err.Error(), "Available config ids:") {
			// Print the error message along with the available configurations
			fmt.Println("Error in fetching config:", err)
		} else {
			// Handle other errors that might occur during configuration retrieval
			fmt.Println("Error:", err)
		}
	}

	intf, err := cfg.Interface(0, 0)
	if err != nil {
		fmt.Println("Error in fetching interface --- : ", err)
		return
	}

	outEndpoint, err := intf.OutEndpoint(0x05) // Replacd 0x01 with the endpoint address with 0x88
	if err != nil {
		fmt.Println("Error --- : ", err)
	}

	fmt.Println("outEndpoint", outEndpoint)
	err = dev.SetAutoDetach(true)
	if err != nil {
		fmt.Println("Error --- : ", err)
	}
	lastMessage = make([]byte, 1)
}

func initDevice(device int) {
	devices := driverFind()
	if len(devices) == 0 {
		fmt.Println("Device not found")
	}

	if device >= len(devices) {
		fmt.Println("Device not found")
	}

	i := devices[device]
	fmt.Printf("Using Device %d (Bus: %d Address: %d IdVendor: %d IdProduct: %d)\n",
		device, i.Desc.Bus, i.Desc.Address, i.Desc.Vendor, i.Desc.Product)
	driverInit(device)
}

func main() {
	args := os.Args
	deviceIndex := 0

	switch args[1] {
	case "-r":
		fmt.Println("Reseting:")
		initDevice(deviceIndex)

		if !utils.Reset() {
			fmt.Println("Reset Error")
			return
		}
		fmt.Println("Reset OK!")

	case "-e":
		fmt.Println("Erasing:")
		if utils.EraseInit() {
			fmt.Println("TC32 EVK : Swire OK")
		}
		var test bool = false
		if !test {
			firmwareSize := 2048 // 524288/16
			barLen := 50

			for i := 0; i < firmwareSize; i += 16 {
				// Placeholder for the eraseAdr function in Go under td package
				utils.EraseAdr(i)

				// hexValue := fmt.Sprintf("%x", i*0x100)
				// fmt.Println("hexValue  : ", hexValue)
				firmwareAddr := i

				percent := (firmwareAddr * 100) / (firmwareSize - 16)

				barProgress := strings.Repeat("#", percent*barLen/100)
				barRemaining := strings.Repeat("=", barLen-(percent*barLen/100))

				fmt.Printf("\r%d%% [\033[3;91m%s\033[0m%s]0x%05x", percent, barProgress, barRemaining, firmwareAddr*256)
				time.Sleep(50 * time.Millisecond) // Simulate the delay as in Python code
			}
		}
	case "-h":
		fmt.Println("____________Help____________\n-h for Help\n-r for Reset\n-e for Erase")

	default:
		fmt.Println("Invalid option \n use 'go run main.go -h' for help")

	}
}
