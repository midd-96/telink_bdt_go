package utils

import (
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/gousb"
)

var lastMessage []byte

func NewSend(wValue uint16, msg []byte, dev *gousb.Device) error {
	ctx := gousb.NewContext()
	defer ctx.Close()

	// Open the USB device with the specific VID and PID.
	// dev, err := ctx.OpenDeviceWithVIDPID(targetVID, targetPID)
	// if err != nil {
	// 	return fmt.Errorf("failed to open USB device: %w", err)
	// }
	// defer dev.Close()

	// Perform the control transfer.
	_, err := dev.Control(
		0x41,
		0x02, // bRequest
		wValue,
		0, // wIndex
		msg,
	)
	if err != nil {
		return fmt.Errorf("error during control transfer: %w", err)
	}

	// Add a delay if needed.
	// time.Sleep(3 * time.Millisecond)

	return nil
}

func receive(wValue uint16, dataOrWLength, timeOut int, reset bool, testIndex int, testData byte, testLast bool, dev *gousb.Device) ([]byte, error) {
	timeoutStart := time.Now().Add(time.Duration(timeOut) * time.Second)

	// Initialize the USB context.
	// ctx := gousb.NewContext()
	// defer ctx.Close()

	// Open the default USB device.
	// dev, err := ctx.OpenDeviceWithVIDPID(targetVID, targetPID)
	// if err != nil {
	// 	fmt.Println("Exiting receive 1")
	// 	return nil, err
	// }
	// defer dev.Close()

	// if reset {
	// 	time.Sleep(3 * time.Millisecond)
	// 	// Perform the reset control transfer.
	// 	res, _ := dev.Control(
	// 		// gousb.ControlOut|gousb.ControlClass|gousb.ControlInterface,
	// 		0xC1,
	// 		0x02, // bRequest
	// 		wValue,
	// 		0, // wIndex
	// 		make([]byte, dataOrWLength),
	// 	)
	// 	if res == 0 {
	// 		fmt.Println("Exiting receive 2")
	// 		return nil, err
	// 	}

	// }

	var data = make([]byte, dataOrWLength)

	for {
		time.Sleep(3 * time.Millisecond)
		// Prepare a buffer to read the data.

		// Perform the control transfer for bulk read (IN endpoint 0x81).
		res, err := dev.Control(
			// gousb.ControlIn|gousb.ControlClass|gousb.ControlInterface,
			0xC1,
			0x02, // bRequest
			wValue,
			0, // wIndex
			data,
		)
		if res == 0 {
			fmt.Println("Exiting receive 3")
			return nil, errors.New("result 0")
		}

		if err != nil {
			fmt.Println("Error in dev.Control ", err)
		}

		if time.Now().After(timeoutStart) {
			break
		}

		if testLast {
			if !equalSlices(data, lastMessage) {
				break
			}
		}

		if testIndex >= 0 && testIndex < len(data) && data[testIndex] == testData {
			break
		}
	}

	if time.Now().After(timeoutStart) && timeOut != 0 {
		return nil, fmt.Errorf("Timeout")
	}

	lastMessage = data
	time.Sleep(3 * time.Millisecond)
	return data, nil
}

// Utility function to check if two slices are equal.
func equalSlices(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func EraseInit(dev *gousb.Device) bool {
	time.Sleep(10 * time.Millisecond)

	msg1, _ := hex.DecodeString("029fff505050505050")
	msg2, _ := hex.DecodeString("029ff800000000000000000001000007")

	// ################1

	NewSend(0x9fff, msg1, dev)
	time.Sleep(30 * time.Millisecond)
	NewSend(0x9ff8, msg2, dev)
	time.Sleep(3 * time.Millisecond)
	receive(0x9ff0, 24, 5, true, 0, 0x01, false, dev)

	// ################2

	msg1, _ = hex.DecodeString("029ff800000000000000000001000007")

	NewSend(0x9ff8, msg1, dev)
	time.Sleep(30 * time.Millisecond)
	receive(0x9ff0, 24, 5, true, 0, 0x01, false, dev)

	// ################3

	msg1, _ = hex.DecodeString("029ff8b2b2b2b2b2b200000001000081")

	NewSend(0x9ff8, msg1, dev)
	time.Sleep(30 * time.Millisecond)
	receive(0x9ff0, 24, 5, true, 0, 0x01, false, dev)

	// ################4

	msg1, _ = hex.DecodeString("02a000abababababab")
	msg2, _ = hex.DecodeString("029ff804040404040400040001000041")

	NewSend(0xa000, msg1, dev)
	time.Sleep(30 * time.Millisecond)
	NewSend(0x9ff8, msg2, dev)
	time.Sleep(3 * time.Millisecond)
	receive(0x9ff0, 24, 5, true, 0, 0x01, false, dev)

	// ################5

	msg1, _ = hex.DecodeString("029ff8040404040404000400010000c1") //

	NewSend(0x9ff8, msg1, dev)
	time.Sleep(30 * time.Millisecond)
	receive(0x9ff0, 24, 5, true, 0, 0x01, false, dev)

	receive(0xa000, 9, 5, true, 0, 0xab, false, dev)

	// ################6

	msg1, _ = hex.DecodeString("02a000050505050505")
	msg2, _ = hex.DecodeString("029ff802020202020206000001000041")

	NewSend(0xa000, msg1, dev)
	// time.Sleep(30 * time.Millisecond)
	NewSend(0x9ff8, msg2, dev)
	// time.Sleep(3 * time.Millisecond)
	receive(0x9ff0, 24, 5, true, 0, 0x01, false, dev)

	// ################7

	msg1, _ = hex.DecodeString("029ff8020202020202060000010000c1")

	NewSend(0x9ff8, msg1, dev)
	time.Sleep(30 * time.Millisecond)
	receive(0x9ff0, 24, 5, true, 0, 0x01, false, dev)

	receive(0xa000, 9, 5, true, 0, 0x05, false, dev)

	// ################8

	msg1, _ = hex.DecodeString("02a000000000000000")
	msg2, _ = hex.DecodeString("029ff822222222222206000001000041")

	NewSend(0xa000, msg1, dev)
	// time.Sleep(30 * time.Millisecond)
	NewSend(0x9ff8, msg2, dev)
	// time.Sleep(3 * time.Millisecond)
	receive(0x9ff0, 24, 5, true, 0, 0x01, false, dev)

	// ################9

	msg1, _ = hex.DecodeString("02a000000000000000")
	msg2, _ = hex.DecodeString("029ff843434343434306000001000041")

	NewSend(0xa000, msg1, dev)
	// time.Sleep(30 * time.Millisecond)
	NewSend(0x9ff8, msg2, dev)
	// time.Sleep(3 * time.Millisecond)
	receive(0x9ff0, 24, 5, true, 0, 0x01, false, dev)

	// ################10
	// string1 := fmt.Sprintln(strings.ReplaceAll("02 a0 00 28 28 28 28 28 28 80 05 00 00 00 00 00 4b 4e 4c 54 c1 01 88 00 a6 80 00 00 58 82 00 00 04 1c 00 00 00 00 00 00 0c 64 81 a2 0a 0b 1a 40 c0 06 c0 06 c0 06 c0 06 c0 06 c0 06 c0 06 c0 06 c0 06 c0 06 c0 06 c0 06 c0 06 c0 06 c0 06 c0 06 0c 6c 70 07 c0 46 c0 46 6f 00 80 00 26 08 c0 6b 27 08 85 06 25 08 c0 6b 26 08 85 06 00 a0 26 09 26 0a 91 02 02 ca 08 50 04 b1 fa 87 24 09 27 08 00 fa 08 40 48 40 2f 08 2f 09 01 50 09 fe 01 41 41 41 2f 08 00 a1 41 40 ab a1 01 40 00 a2 06 a3 01 b2 9a 02 fc cd 01 a1 41 40 24 08 24 09 02 ec 0a 40 40 a2 8a 40 8a 48 d2 f7 d2 ff 01 aa fa c0 4a 48 00 a0 82 02 09 c1 16 09 17 0a 17 0b 9a 02 04 ca 08 58 10 50 04 b1 04 b2 f8 87 14 09 15 0a 15 0b 9a 02 04 ca 08 58 10 50 04 b1 04 b2 f8 87 01 90 10 9c fe 87 c0 46 12 00 00 00 13 00 00 00 d0 21 84 00 00 fe 84 00 50 20 84 00 dd 28 84 00 0c 06 80 00 00 09 84 00 00 0a 84 00 00 f0 00 00 00 09 00 00 bc 1b 00 00 00 20 84 00 00 20 84 00 bc 1b 00 00 00 20 84 00 48 20 84 00 7e 00 00 00 b8 00 80 00 60 00 80 00 00 00 00 ff 00 00 80 00 0c 00 80 00 c4 00 00 00 ff ff ff ff 00 20 84 00 00 00 85 00 09 00 00 00 00 65 ff 64 d8 6b 41 06 4a 06 53 06 5c 06 65 06 3f 64 00 90 0d 98 3f 6c 88 06 91 06 9a 06 a3 06 ac 06 d0 6b ff 6c 00 69 c0 46 c0 46 c0 46 c0 46 70 07 c0 46 10 65 00 f6 00 fe 0a 0b 1c 48 00 a2 1a 40 09 0b 18 40 09 09 40 a3 0b 40 01 a2 0b 48 13 00 fc c1 06 0a 10 48 01 b2 13 40 01 0b 1c 40 10 6d c0 46 43 06 80 00 b8 00 80 00 ba 00 80 00 b9 00 80 00 10 65 00 f6 00 fe 09 f6 09 fe 0a 0b 1c 48 00 a2 1a 40 09 0b 18 40 01 b3 19 40 08 09 60 a3 0b 40 01 a2 0b 48 13 00 fc c1 04 0a 13 40 01 0b 1c 40 10 6d c0 46 43 06 80 00 b8 00 80 00 ba 00 80 00 03 0a 11 58 00 f1 13 58 5b ea 98 02 fb c2 70 07 40 07 80 00 00 65 c8 a0 80 a1 ff 97 d1 9f 30 a0 01 a1 07 a2 07 a3 01 90 ad 9b c7 a0 0e a1 ff 97 c7 9f c7 a0 0f a1 ff 97 c3 9f cf a0 ff 97 a0 9f 03 f6 fa c5 cb a0 ff 97 9b 9f 01 ec 33 a0 ff 97 b7 9f 30 a0 00 a1 07 a2 07 a3 01 90 93 9b c7 a0 0e a1 ff 97 ad 9f 00 6d 10 65 30 a0 60 a1 ff 97 a7 9f c6 a0 f6 a1 ff 97 a3 9f c6 a0 f7 a1 ff 97 9f 9f 40 a4 cf a0 ff 97 7b 9f 04 02 fa c0 c9 a0 ff 97 76 9f 01 ec 32 a0 ff 97 92 9f ca a0 ff 97 6f 9f 01 ec 31 a0 ff 97 8b 9f c6 a0 f6 a1 ff 97 87 9f 30 a0 20 a1 ff 97 83 9f 10 6d 00 65 2d a0 ff 97 5e 9f 7f a1 01 00 2d a0 ff 97 79 9f 09 0b 19 48 02 a2 0a 03 12 f6 12 fe 1a 40 06 0b 1a 48 02 a1 8a 03 1a 40 22 b3 1a 48 0c a1 8a 03 1a 40 00 6d c0 46 73 00 80 00 86 05 80 00 02 f2 12 fe 0c 0b 1a 40 0c 09 10 a2 0b 48 1a 02 fc c1 02 f4 12 fe 08 0b 1a 40 08 09 10 a2 0b 48 1a 02 fc c1 00 f6 00 fe 03 0b 18 40 03 09 10 a2 0b 48 1a 02 fc c1 70 07 0c 00 80 00 0d 00 80 00 30 65 05 ec 07 0c 01 a3 23 40 01 a0 ff 97 5a 9f 00 a3 23 40 04 0b 1d 40 10 a2 23 48 1a 02 fc c1 30 6d c0 46 0d 00 80 00 0c 00 80 00 70 65 64 a0 ff 97 48 9f 05 a0 ff 97 e3 9f 0a 08 0a 0c 00 a6 0a 09 10 a2 01 a5 26 40 0b 48 1a 02 fc c1 23 48 1d 02 02 c0 01 b8 00 a8 f5 c1 01 a2 03 0b 1a 40 70 6d c0 46 80 96 98 00 0c 00 80 00 0d 00 80 00 70 65 06 ec 09 0c 25 48 00 a3 23 40 06 a0 ff 97 bf 9f 20 a0 ff 97 bc 9f 30 ec ff 97 99 9f 01 a2 03 0b 1a 40 ff 97 ca 9f 25 40 70 6d 43 06 80 00 0d 00 80 00 f0 65 06 ec 0c ec 15 ec 11 0b 1f 48", " ", ""))

	msg1, _ = hex.DecodeString("02a000282828282828800500000000004b4e4c54c1018800a680000058820000041c0000000000000c6481a20a0b1a40c006c006c006c006c006c006c006c006c006c006c006c006c006c006c006c0060c6c7007c046c0466f0080002608c06b270885062508c06b2608850600a02609260a910202ca085004b1fa872409270800fa084048402f082f09015009fe014141412f0800a14140aba1014000a206a301b29a02fccd01a141402408240902ec0a4040a28a408a48d2f7d2ff01aafac04a4800a0820209c11609170a170b9a0204ca0858105004b104b2f8871409150a150b9a0204ca0858105004b104b2f8870190109cfe87c0461200000013000000d021840000fe840050208400dd2884000c06800000098400000a840000f0000000090000bc1b00000020840000208400bc1b000000208400482084007e000000b800800060008000000000ff000080000c008000c4000000ffffffff0020840000008500090000000065ff64d86b41064a0653065c0665063f6400900d983f6c880691069a06a306ac06d06bff6c0069c046c046c046c0467007c046106500f600fe0a0b1c4800a21a40090b1840090940a30b4001a20b481300fcc1060a104801b21340010b1c40106dc04643068000b8008000ba008000b9008000106500f600fe09f609fe0a0b1c4800a21a40090b184001b31940080960a30b4001a20b481300fcc1040a1340010b1c40106dc04643068000b8008000ba008000030a115800f113585bea9802fbc27007400780000065c8a080a1ff97d19f30a001a107a207a30190ad9bc7a00ea1ff97c79fc7a00fa1ff97c39fcfa0ff97a09f03f6fac5cba0ff979b9f01ec33a0ff97b79f30a000a107a207a30190939bc7a00ea1ff97ad9f006d106530a060a1ff97a79fc6a0f6a1ff97a39fc6a0f7a1ff979f9f40a4cfa0ff977b9f0402fac0c9a0ff97769f01ec32a0ff97929fcaa0ff976f9f01ec31a0ff978b9fc6a0f6a1ff97879f30a020a1ff97839f106d00652da0ff975e9f7fa101002da0ff97799f090b194802a20a0312f612fe1a40060b1a4802a18a031a4022b31a480ca18a031a40006dc046730080008605800002f212fe0c0b1a400c0910a20b481a02fcc102f412fe080b1a40080910a20b481a02fcc100f600fe030b1840030910a20b481a02fcc170070c0080000d008000306505ec070c01a3234001a0ff975a9f00a32340040b1d4010a223481a02fcc1306dc0460d0080000c008000706564a0ff97489f05a0ff97e39f0a080a0c00a60a0910a201a526400b481a02fcc123481d0202c001b800a8f5c101a2030b1a40706dc046809698000c0080000d008000706506ec090c254800a3234006a0ff97bf9f20a0ff97bc9f30ecff97999f01a2030b1a40ff97ca9f2540706d430680000d008000f06506ec0cec15ec110b1f48")
	msg2, _ = hex.DecodeString("029ff800000000000000040000040041")

	NewSend(0xa000, msg1, dev)
	// time.Sleep(30 * time.Millisecond)
	NewSend(0x9ff8, msg2, dev)
	// time.Sleep(3 * time.Millisecond)
	receive(0x9ff0, 24, 5, true, 1, 0x04, false, dev)

	// ###11

	msg1, _ = hex.DecodeString("02a000000000000000a21a4006a0ff97a39f02a0ff97a09f30ecff977d9f00ac0bc000a00a0e0b0910a22b1c33400b481a02fcc101b08402f7c801a2050b1a40ff97a09f010b1f40f06dc046430680000c0080000d008000f0654706806407ec0cec15ec170b1a48900600a61e4003a0ff97729f38ecff974f9f130b1e40130810a1120a03481902fbc10aa3134010a113481902fcc100ac0bc000a00a0e0b0910a233482b140b481a02fcc101b0a002f7c101a2050b1a40020b42061a40046c9006f06d430680000c0080000d00800030650b0b1d4800a41c4005a0ff973c9f080b1c40080910a20b481a02fcc1050b184801a201b31a40010b1d40306dc046430680000c0080000d008000706504f624fe150b1e4800a21a4006a0ff971c9f01a0ff97199f110b1c40110a10a3100c15481d00fbc101a32340ff97239f64a0ff976c9e05a0ff97079f080b1d4010a223481a02fcc1050b184801a201b31a40010b1e40706dc046430680000c0080000d008000706506ec090c254800a3234006a0ff97e99e52a0ff97e69e30ecff97c39e01a2030b1a40ff97f49e2540706d430680000d008000706506ec090c254800a3234006a0ff97cf9ed8a0ff97cc9e30ecff97a99e01a2030b1a40ff97da9e2540706d430680000d0080003065070c254800a32340b9a0ff97b69e01a2040b1a4001a0ff97129e2540306d430680000d008000f065070c274800a32340aba0ff97a29e040d01a62e40ff97b39e2e402740f06d430680000d008000706506ec090c254800a3234006a0ff978d9e81a0ff978a9e30ecff97679e01a2030b1a40ff97989e2540706d430680000d0080003065080c254800a3234006a0ff97749e60a0ff97719e01a2030b1a40ff97829e2540306d430680000d008000706504ec130b1e4800a51d409fa0ff975d9e110b1d40110810a1100a03481902fbc10aa3134010a113481902fcc1e5ec09080a0910a20348234001b40b481a02fcc1ac02f7c101a2040b1a40010b1e40706dc046430680000c0080000d008000f0650cec00f605fe210b1f4800a61e4028ecff972b9e4bad23c05aad2cc000a21c0b1a401c0810a11b0a03481902fbc10aa3134010a113481902fcc125ec10b51408150910a20348234001b40b481a02fcc1ac02f7c101a20f0b1a400c0b1f40f06d00a0ff97e29d0a0b1e400a0910a20b481a02fcc1d28780a0ff97d79d050b1e40050910a20b481a02fcc1c787c046430680000c0080000d008000f06547068064816034a087a1ff97229d00a3260a1340b9a101ba11400033003b01ab05c8003b01b30033003b01abf9c91e0f01a33b401e0e00a434400ba038a1ff97089d82a00ca1ff97049d190b1d482df2190bede82b589806180b2b50ff971f9c43062b5082a064a1ff97f39c0ba03ba1ff97ef9c0fa333403c40aba2100b")
	msg2, _ = hex.DecodeString("029ff800000000000004040000040041")

	NewSend(0xa000, msg1, dev)
	// time.Sleep(30 * time.Millisecond)
	NewSend(0x9ff8, msg2, dev)
	// time.Sleep(3 * time.Millisecond)
	receive(0x9ff0, 24, 5, true, 1, 0x04, false, dev)

	// ###12

	msg1, _ = hex.DecodeString("02a0001a1a1a1a1a1a400034003b01ab05c8003b01b30033003b01abf9c901a2040b1a4034a080a1ff97d89c0160046c9006f06d0d008000a10580000d06800058008400c006c0060c008000f06528a2150b1a4000a400a5140a08a108a7140e1740c006c006c006c006c006c006c006c006c006c006c006c006c006c006c006c00613481902fcc1305800ac06c043eb01ab07c92bec430001ab02c001b405ecde8728ecf06dc0464c0780004f078000540780000065080817a101903198faa35bf1060a132001a2050b1a4000a3050a134001b21340006d601b0000500780004f078000200c800000000000000000000000000000000000000000000000000000000000000000001065810208c200aa05c000a3c41ccc1401b39302fac1106d00aafcc053eec0e8c9e8540200a3c21cca1401bba302fac1f187c04600aa06c009f609fe00a3c11401b39302fbc17007106504ec08ec21ecff97d69f20ec106d826003eccb00013320a359ea88000030013b0038180302607007c046f0655f0656064d064406f064b860073008788006187a023203ec084800f64a4812f482e8c84812e8884800f212e804d304b1023c9c02f0c1287d0035307e0136073908b130ec20a2ff97c49f303d9908820600a7307c028001b710af60c0e1598906265930ec06a1ff97bc9f043030ec0ba1ff97b79f830630ec19a1ff97b29fa2599406530604b39a0604bb04db430604b3980604bb02db8ae84a045906043b5900480010e8625961064a00160063067300c6e828ec02a1ff97949f830628ec0da1ff978f9f810628ec16a1ff978a9f63581aec2a03a1580a002b0015ec1d034b065a065300580028e8e3589be9e35086e9e65104bc21ec20b1255a255007a33b0007aba8c1013b003a61da61d361da61d360da60d3255a0cec01b710af9ec16408810640a08104087e0a800fa33b000fab00c18f8001b740af00c197802558043eb15b880608ec11a1ff974b9f0330400613a1ff97469f0530725a930633ed04337058820607a1ff973c9f0630500612a1ff97379f33589b04053a03394a00430699fa4a0013ec5b04063948005106caf850001be89a063354e2590332265930ec06a1ff971e9f830630ec0ba1ff97199f800630ec19a1ff97149fa3599c06490604b1890604b904d941065b0659004800033940e882e86159630659000e007300d6e8560428ec02a1ff97fb9e820628ec0da1ff97f69e800628ec16a1ff97f19e63581aec2a03a1580a002b0015ec1d03430651064b00580028e8e3589be9e35086e9e65104bc21ec20b1235a235007a33b0007ab00c07887013b003a70da70d331da31d341da41d30cec0fa33b000fab00c06f870878187940a2ff97ba9e043940b9043101b740af00c06787287c00a1073d40dc88ec80f0471d2be85a4812f23a039f483ff43a03df48")
	msg2, _ = hex.DecodeString("029ff800000000000008040000040041")

	NewSend(0xa000, msg1, dev)
	// time.Sleep(30 * time.Millisecond)
	NewSend(0x9ff8, msg2, dev)
	// time.Sleep(3 * time.Millisecond)
	receive(0x9ff0, 24, 5, true, 1, 0x04, false, dev)

	// ###13

	msg1, _ = hex.DecodeString("02a0003f3f3f3f3f3ff63a0392e92a1410fa584010fc984012feda4001b108a9e6c138603c6c90069906a206ab06f06d181a0000006500a168a2ff977f9e006df0654706806405ec0cec16ee35c0024843481bf21303824812f41303c24812f613033fa01800f3e82b401afa6a401afcaa401bfeeb4033e898063fa3980520c900a810c040a73fea2eec28b630e83aecff975e9e28ec31ecff97709ee4e9460640be3fae07c928ec21ecff97679e40b440be3faef7c800ae02c1046c9006f06d00a028b528e821ec32ecff97419ef487f06547068064906004ec8806034847483ff21f0383481bf41f03c3481bf61f03fef06d06680600a140a2ff971f9e80a32b403fa2170037af32c838a2d2eb20ec6906ff97959f00a3003333fe2b4133fc6b4133faab41ee4120ec690608a2ff97879f430600a188f802b080f0051d20e8424812f22a0385482df42a03c04800f6020310fe184010fc584010fa9840da4004b104b320a9e6c11060046c9006f06d78a2d2ebcb87c046f065470680649a6006ec880617ec6806ff97549f00a3003301330a0a027b23da23d323da23d303da03d3680631ec4206ff974a9f680639ecff978e9f1a60046c9006f06d181b0000080b1bf31bfb080a9be8080a1350080a125819f609fe11401bf41bfe534001a07007c046d02184000080fcffd42884004420840070651a0e33581a49584900f210039a4912f41003da4912f610031a4a5c4a24f2994adb4a1403120d21ec2aecff971a9b00ac0dc00f082a4803489a0214c100a30380e91cc21c910208c101b39c02f8c8335824f624fe1c4001a0706d1bf61bfe3258534000a0f88700a3f98744208400d0238400d0218400f065230b1a581349504900f2180393491bf4d14918031f0b1850134953499349d54900ad2ec000a31b0c164a514a09f23103964a36f43103d64a36f63103def0f100e11401b39d02efc1140e29ec32ecff97cc9a12096de80bec00a7f2e8110952e811481940e0e80e0940e811480248910200c001a701b3ab02efc178ee47eeb801f06d29ec050aff97b09a00a7f58744208400d0288400d0218400d023840008000401f8fffbfe7065180b1b581a49584900f210039a4912f4d9491003140a10501a495a499a49dd49120c29ec22ecff978c9a00a600ad14cd0f0968e80bec028001b383020dc0e2e80c0952e8114819401248ffaaf4c101b636f636fe01b38302f1c1a8eb45eea801706d44208400d0288400d023840008000401f8fffbfe70658260150b1b5819495a4912f20a03994909f40a03d84900f601ec11031a4a5c4a24f2984adb4a140363ee00a007ab10c802ac10c90a0e6d0608ec21ec6a06ff97449a00a39ae9e91c114001b3a302f9c101a00260706d020eed87442084000800040104000401f06557064e064506e064270b1b581a49584900f2")
	msg2, _ = hex.DecodeString("029ff80000000000000c040000040041")

	NewSend(0xa000, msg1, dev)
	// time.Sleep(30 * time.Millisecond)
	NewSend(0x9ff8, msg2, dev)
	// time.Sleep(3 * time.Millisecond)
	receive(0x9ff0, 24, 5, true, 1, 0x04, false, dev)

	// ###14

	msg1, _ = hex.DecodeString("02a000101010101010039a4912f4d9491003230e3050194a5a4a12f20a03994a09f40a03db4a1bf613039bf81d0f3b5001a200ab1bc000a399061b0d6cecffb480a35bf19806ff97bf9900a39a063058500480a149f02aecff97fe992bec028001b3a30209c01a48ffaaf9c000a210ec1c6c90069906a206f06d80a35bf09a04c205e4c1305880a35bf1c0e8305001a399043b584b05d6c801a2e88744208400d0288400d8288400d0238400f0655f0656064d064406f06498609b0e335819495a4912f20a03994909f40a03d94909f611038806194a5a4a12f20a03994a09f40a03db4a1bf613039906117d8f0b2aec13db13d21b5813508d0b4006180382068c0c02a122ecff97a799234866ab1bc0634866ab18c01778ff97b89a173b1bf41bfc173384081ae800aa36c083095ae800aa32c0820a93022fc0820a1320335806a21a4000a02180400608a122ecff9783997d0b00a57d08028001b383020dc0e2e87b0952e8114819401248ffaaf4c101b52df62dfe01b38302f1c108ad42c0335899a21a4088a25a4001a018603c6c90069906a206ab06f06d0d7f4ba039ecff97a49a00a300a0028001b310ab09c0f91cea1c9102f8c101b000f600fe01b310abf5c110a821c0420613fb1bf39b0618ecff97f99880a35bf118ec5804003065ecffb5580680a149f022ecff97349923ec028001b3ab0225c01a48ffaaf9c0335803a29d87335801a29a87400608a122ecff9721994c0b00a54c08028001b383020bc0e2e84a0952e8114819401248ffaaf4c101b52df62dfef08708ad53c0335866a29e8780a149f08b04003a9305c8c1017d38ec10a12aec01a3ff978c9d00a0097f840663063b1403ec00a2e91c4a0008b31fabfac93a1401b008a8f2c13b483fab02c87a488baa44c0480602f2157d2b407b486b40bb48ab40e84013fc2b4112fe6a41fb48ab413b49eb41400608a12aecff97a298400608a122ecff97cb9800a300a020095ae8e11c1140e91ce21c910200c001a001b308abf3c100a805c0335804a22c87335802a2298766a32b40500601a12aecff978098500601a122ecff97a9982a4823489a0200c13c87335805a21587ba489eaab7c17b49ba497a40fa49ba40b18744208400381b0000fe0f0000d0238400389fffffafbfffff85600000088000010800040110000401f8fffbfef06557064606c064250d2b587fa4250e00a2920677ecffb7b80602805a4801aa1fc0da48220049aa35c11a4801aaf5c11a49584900f210039a4912f41003da4912f610031a4a594a09f29f4adb4a110332ecff9729982b5852061a405a4801aadfc11a49584900f210039a4912f41003da4912f610031a4a594a09f29f4adb4a11034206ff9710982b5852065a40da48220049aac9c001a00c6c90069a06f06d44208400d0218400f06547068064220f3b581949")
	msg2, _ = hex.DecodeString("029ff800000000000010040000040041")

	NewSend(0xa000, msg1, dev)
	// time.Sleep(30 * time.Millisecond)
	NewSend(0x9ff8, msg2, dev)
	// time.Sleep(3 * time.Millisecond)
	receive(0x9ff0, 24, 5, true, 1, 0x04, false, dev)

	// ###15

	msg1, _ = hex.DecodeString("02a0005a5a5a5a5a5a4912f20a03994909f40a03d94909f6110388061a4a5c4a24f2994adb4a1403190e400621ec32ecfe97e49f170d400621ec2aecff970c9800ac0cc02a4833489a021bc100a30380e91cf21c91020cc101b39c02f8c83b5822f612fe1a4024fa5c4001a0046c9006f06d19f609fe1bf41bfe3a581140534000a0f38700a300a1f787c04644208400d0218400d0238400f065280b1a581349504900f2180393491bf4d1491803240e3050134953499349d54900ad33c000a3200c174a514a09f23903974a3ff43903d74a3ff63903dff0f900e11401b39d02efc129ec22ecfe97899f3058160e29ec32ecfe97b19f15096de80bec00a7f2e8130952e811481940e0e8110940e811480248910200c001a701b3ab02efc178ee47eeb801f06d29ec060afe97679f305829ec050afe97909f00a7f08744208400d0288400d0218400d023840008000401f8fffbfef06557064e064506e064410f3a581349504900f2180393491bf4d14918033d098a060850134953499349d34900ab98065ac000a3380d144a514a09f22103944a24f42103d44a24f62103dcf0e100e91401b39805efc1310c410622ecfe97529f2f0a16ec460413ec00a1890699a2940688a0e2e82b0952e8114819401248ffaa05c03a5861061140504001a2910601b3b302eec100a001a3990520c05106085841062aecfe97009f52061058410622ecfe97289f1a0b00a0e2e81a0952e811481940efe817097fe811483a48910200c001a001b3b302efc101b843ee98011c6c90069906a206f06d0c0c410622ecfe97099f520610584106070afe97d59e53061858410622ecfe97fd9e00a0e487c04644208400d0288400d0218400d023840008000401f8fffbfe30650e0d2b5819495a4912f20a03994909f40a03db491bf6130300a401ab0ac020ecfe97359f2b58184044001fa0200044026001306d28a4f287c04644208400f06557064606c064900600a2900531c01a0d2b4800ab1ec1190b9a061aec00a701a5180e80a35bf09c063bec00a404805bf8730001b408ac05c01d02f8c15bf801b408acf9c108d201b76705edc101a30a0d2b4001800a0a920600a3ffa4ca1c4200220092f05506aa1800fa500001b39805f4c80c6c90069a06f06ddc288400d02484002083b8edf0655f0656064d064406f06481603e088206035819495a4912f20a03994909f40a03d94909f611030031194a5a4a12f20a03994a09f40a03db4a1bf6130399065bfa98064906cbf5dbfd9b0600a2900555c0003d00a401a676022c0f28ec80a189f03aecfe97569e30ec39ec80a292f0ff97829f06ec01b480a39bf0ede8a005ecc8430600a083050dc05bf2003958e81e0c590622ecfe973d9e30ec21ec5a06ff976a9f06ecf60352061358ffa211ec3100184a194231fa1100584a594231fc1100984a")
	msg2, _ = hex.DecodeString("029ff800000000000014040000040041")

	NewSend(0xa000, msg1, dev)
	// time.Sleep(30 * time.Millisecond)
	NewSend(0x9ff8, msg2, dev)
	// time.Sleep(3 * time.Millisecond)
	receive(0x9ff0, 24, 5, true, 1, 0x04, false, dev)

	// ###16

	msg1, _ = hex.DecodeString("02a0009999999999994236fed94ade424906110018491941480601fa110058495941480601fc0a0099499a4102fed949da4101a001603c6c90069906a206ab06f06d00a301a67602be87c04644208400d02184000065ff972d9802a00090c09804a2010b1a40006d61008000f0654706806401a04002230901a2ff971f9f220b98061958210d7fa400a6cb48a303fcc02aec00a3c84840f640fe1748bc06840507c001b305b20dabf4c1cb482300cb40eb870e404e408e409af0d3e8ebe859489a4812f20a03d94809f40a031b491bf6130300901d9800f600fe01a80bc047063958ca4852f652feffa353008b40cb482300cb40c98743061958cb485bf65bfe8b40f487d021840044208400002084001807c0460065ff979d9fff97a79ffe87f06557064606c0640eee29cd04ec00a540a780a213f498063fa29206078008aa23c007aa27c001b504b4ae0218cd2348614809f208ec1803a148e3481f02f2c042060203940652061a0003aae7c180a212f413ec6304194001b504b4ae02e6cc30ec0c6c90069a06f06d630618f600fefe97249cd78700f240e8fe973f9cd287f0650ff63ffe16f636fe1cf624fe05f62dfe28ecfe97f29b01b636eb01a3b3001eec01be33eca30001ec99033700a700390309f609fe28ecfe97009cf06dc04600f600fe05a808c900a20b0b1a400ab31a4801a18a031a40700780f0070b1b189f0620a2f18744a2ef8743a2ed8742a2eb8760a2e987c04666008000481b0000982f8a4291443771cffbc0b5a5dbb5e95bc25639f111f159a4823f92d55e1cab98aa07d8015b8312be853124c37d0c55745dbe72feb1de80a706dc9b74f19bc1c1699be48647beefc69dc10fcca10c246f2ce92daa84744adca9b05cda88f97652513e986dc631a8c82703b0c77f59bff30be0c64791a7d55163ca0667292914850ab72738211b2efc6d2c4d130d385354730a65bb0a6a762ec9c281852c7292a1e8bfa24b661aa8708b4bc2a3516cc719e892d1240699d685350ef470a06a1016c1a419086c371e4c774827b5bcb034b30c1c394aaad84e4fca9c5bf36f2e68ee828f746f63a5781478c8840802c78cfaffbe90eb6c50a4f7a3f9bef27871c667e6096a85ae67bb72f36e3c3af54fa57f520e518c68059babd9831f19cde05b51015101510151015101510151015101fe190000021a0000061a00000a1a0000fa190000e0190000600000c3610000c3620000c36300ffc36400ffc36500ffc3820064c8340080c80b003bc88c0002c8270000c8280000c8290000c82a0000c8400c04c3410c04c3420c04c3430c04c3440c04c3450c04c3460c04c3470c04c3480c04c340b90d000041f513000042ed0d0000439114000044650e0000454d150000460d0f00004ded0f00004e751600004f850f0000494d130000503d17000056a5100000000000")
	msg2, _ = hex.DecodeString("029ff800000000000018040000040041")

	NewSend(0xa000, msg1, dev)
	// time.Sleep(30 * time.Millisecond)
	NewSend(0x9ff8, msg2, dev)
	// time.Sleep(3 * time.Millisecond)
	receive(0x9ff0, 24, 5, true, 1, 0x04, false, dev)

	// ###17

	msg1, _ = hex.DecodeString("02a000040404040404008400")
	msg2, _ = hex.DecodeString("029ff80000000000001c040004000041")

	NewSend(0xa000, msg1, dev)
	NewSend(0x9ff8, msg2, dev)
	receive(0x9ff0, 24, 5, true, 0, 0x04, false, dev)

	// ###18

	msg1, _ = hex.DecodeString("02a000888888888888")
	msg2, _ = hex.DecodeString("029ff802020202020206000001000041")

	NewSend(0xa000, msg1, dev)
	// time.Sleep(30 * time.Millisecond)
	NewSend(0x9ff8, msg2, dev)
	// time.Sleep(30 * time.Millisecond)
	receive(0x9ff0, 24, 5, true, 0, 0x01, false, dev)

	// ###
	return true
}

func EraseAdr(adr int, dev *gousb.Device) bool {
	msg1, _ := hex.DecodeString("02a0004d4d4d4d4d4d0000000004000000")

	adrBin := make([]byte, 2)
	adrBin[0] = byte(adr & 0xFF)
	adrBin[1] = byte((adr >> 8) & 0xFF)

	copy(msg1[10:12], adrBin)

	msg2, _ := hex.DecodeString("029ff807070707070700040009000041")

	NewSend(0xa000, msg1, dev)
	time.Sleep(3 * time.Millisecond)
	NewSend(0x9ff8, msg2, dev)
	receive(0x9ff0, 24, 5, true, 0, 0x09, false, dev)
	// fmt.Printf("%x", result)

	// ######################

	msg1, _ = hex.DecodeString("02a000cdcdcdcdcdcd")
	msg2, _ = hex.DecodeString("029ff807070707070700040001000041")

	NewSend(0xa000, msg1, dev)
	time.Sleep(30 * time.Millisecond)
	NewSend(0x9ff8, msg2, dev)
	receive(0x9ff0, 24, 5, true, 0, 0x01, false, dev)
	// fmt.Printf("%x\n", result)

	msg1, _ = hex.DecodeString("029ff8070707070707000400010000c1")

	timeoutStart := time.Now().Add(5 * time.Second)

	for {
		NewSend(0x9ff8, msg1, dev)

		result, _ := receive(0x9ff0, 24, 5, true, 0, 0x01, false, dev)

		result, _ = receive(0xa000, 9, 0, false, 0, 0x01, false, dev)

		if time.Now().After(timeoutStart) {
			break
		}

		if result[0] == 0x4d {
			break
		}
	}

	if time.Now().After(timeoutStart) {
		return false
	}

	msg1, _ = hex.DecodeString("029ff80404040404040004000c0000c1")

	NewSend(0xa000, msg1, dev)
	receive(0x9ff0, 24, 5, true, 0, 0x0c, false, dev)

	receive(0xa000, 20, 0, false, 0, 0, false, dev)
	return true
}

func Reset(dev *gousb.Device) bool {
	msg1, _ := hex.DecodeString("029fff505050505050")
	err := NewSend(0x9fff, msg1, dev)
	if err != nil {
		fmt.Println("returns from here")
		return false
	}

	msg1, _ = hex.DecodeString("02a000202020202020")
	msg2, _ := hex.DecodeString("029ff86f6f6f6f6f6f00000001000041")

	err = NewSend(0xa000, msg1, dev)
	if err != nil {
		return false
	}
	err = NewSend(0x9ff8, msg2, dev)
	if err != nil {
		return false
	}
	receive(0x9ff0, 24, 5, true, 0, 0x01, false, dev)

	return true

}

func Activate(dev *gousb.Device) bool {
	msg1, _ := hex.DecodeString("029fff505050505050")
	msg2, _ := hex.DecodeString("029fff4d4d4d4d4d4d")
	msg3, _ := hex.DecodeString("029ff87f7f7f7f7f7f000000010000c1")

	err := NewSend(0x9fff, msg1, dev)
	if err != nil {
		return false
	}
	receive(0x9ff4, 12, 0, true, 0, 0, false, dev)

	err = NewSend(0x9fff, msg2, dev)
	if err != nil {
		return false
	}
	receive(0x9ff4, 12, 0, true, 0, 0, false, dev)

	err = NewSend(0x9fff, msg1, dev)
	if err != nil {
		return false
	}
	receive(0x9ff4, 12, 0, true, 0, 0, false, dev)

	err = NewSend(0x9fff, msg2, dev)
	if err != nil {
		return false
	}
	receive(0x9ff4, 12, 0, true, 0, 0, false, dev)

	err = NewSend(0x9ff8, msg3, dev)
	if err != nil {
		return false
	}
	time.Sleep(3 * time.Millisecond)

	receive(0xa000, 9, 5, true, 0, 0x55, false, dev)

	return true

}

func DownloadInit(dev *gousb.Device) bool {

	EraseInit(dev)

	msg1, _ := hex.DecodeString("02a0004040404040400000000000000000")
	msg2, _ := hex.DecodeString("029ff807070707070700040009000041")
	err := NewSend(0xa000, msg1, dev)
	if err != nil {
		return false
	}
	err = NewSend(0x9ff8, msg2, dev)
	if err != nil {
		return false
	}

	receive(0x9ff0, 24, 5, true, 0, 0x09, false, dev)

	msg1, _ = hex.DecodeString("02a000c0c0c0c0c0c0")
	msg2, _ = hex.DecodeString("029ff807070707070700040001000041")
	err = NewSend(0xa000, msg1, dev)
	if err != nil {
		return false
	}
	err = NewSend(0x9ff8, msg2, dev)
	if err != nil {
		return false
	}

	receive(0x9ff0, 24, 5, true, 0, 0x01, false, dev)

	msg1, _ = hex.DecodeString("029ff8070707070707000400010000c1")
	err = NewSend(0x9ff8, msg1, dev)
	if err != nil {
		return false
	}

	receive(0x9ff8, 24, 5, true, 0, 0x01, false, dev)
	receive(0xa000, 9, 5, true, 0, 0x40, false, dev)

	msg1, _ = hex.DecodeString("029ff80404040404040004000c0000c1")
	err = NewSend(0x9ff8, msg1, dev)
	if err != nil {
		return false
	}

	receive(0x9ff0, 24, 5, true, 0, 0x0c, false, dev)
	receive(0xa000, 20, 5, true, 0, 0, false, dev)

	return true
}

func Download_Block_Init(adr int, dev *gousb.Device) bool {

	msg1, _ := hex.DecodeString("02a0004d4d4d4d4d4d0000000004000000")

	adrBin := make([]byte, 2)
	adrBin[0] = byte(adr & 0xFF)
	adrBin[1] = byte((adr >> 8) & 0xFF)

	copy(msg1[10:12], adrBin)
	msg2, _ := hex.DecodeString("029ff807070707070700040009000041")

	err := NewSend(0xa000, msg1, dev)
	if err != nil {
		return false
	}
	err = NewSend(0x9ff8, msg2, dev)
	if err != nil {
		return false
	}
	receive(0x9ff0, 24, 5, true, 0, 0x09, false, dev)

	msg1, _ = hex.DecodeString("02a000cdcdcdcdcdcd")
	msg2, _ = hex.DecodeString("029ff807070707070700040001000041")
	err = NewSend(0xa000, msg1, dev)
	if err != nil {
		return false
	}
	err = NewSend(0x9ff8, msg2, dev)
	if err != nil {
		return false
	}

	receive(0x9ff0, 24, 5, true, 0, 0x01, false, dev)

	msg1, _ = hex.DecodeString("029ff8070707070707000400010000c1")

	timeoutStart := time.Now().Add(5 * time.Minute)

	for {
		err = NewSend(0x9ff8, msg1, dev)
		if err != nil {
			return false
		}
		result, _ := receive(0x9ff0, 24, 5, true, 0, 0x01, false, dev)
		receive(0xa000, 9, 0, false, 0, 0, false, dev)

		if time.Now().After(timeoutStart) {
			break
		}
		byteSlice := []byte{0x4d}

		if equalSlices(result, byteSlice) {
			break
		}

	}
	if time.Now().After(timeoutStart) || timeoutStart == time.Now() {
		fmt.Println("Timeout")
		return false
	}

	msg1, _ = hex.DecodeString("029ff80404040404040004000c0000c1")
	err = NewSend(0x9ff8, msg1, dev)
	if err != nil {
		return false
	}
	receive(0x9ff0, 24, 5, true, 0, 0x0c, false, dev)
	receive(0xa000, 20, 0, false, 0, 0, false, dev)
	return true
}

func Download_Adr(adr int, dev *gousb.Device) bool {
	msg1, _ := hex.DecodeString("02a000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")

	copy(msg1[8:], []byte{byte(adr)})
	msg1[3] = byte(adr)
	msg1[4] = byte(adr)
	msg1[5] = byte(adr)
	msg1[6] = byte(adr)
	msg1[7] = byte(adr)

	msg2, _ := hex.DecodeString("029ff8d0d0d0d0d0d021040000010041")
	err := NewSend(0xa000, msg1, dev)
	if err != nil {
		return false
	}
	err = NewSend(0x9ff8, msg2, dev)
	if err != nil {
		return false
	}
	receive(0x9ff0, 24, 5, true, 8, 0xd0, false, dev)
	receive(0x9ff0, 24, 5, true, 1, 0x01, false, dev)

	msg1, _ = hex.DecodeString("02a0004141414141410000000000010000")

	adrBin := make([]byte, 2)
	copy(msg1[10:12], adrBin)

	msg2, _ = hex.DecodeString("029ff807070707070700040009000041")

	err = NewSend(0xa000, msg1, dev)
	if err != nil {
		return false
	}
	err = NewSend(0x9ff8, msg2, dev)
	if err != nil {
		return false
	}
	receive(0x9ff0, 24, 5, true, 5, 0x09, false, dev)

	msg1, _ = hex.DecodeString("02a000c1c1c1c1c1c1")
	msg2, _ = hex.DecodeString("029ff807070707070700040001000041")

	err = NewSend(0xa000, msg1, dev)
	if err != nil {
		return false
	}
	err = NewSend(0x9ff8, msg2, dev)
	if err != nil {
		return false
	}
	receive(0x9ff0, 24, 5, true, 5, 0x01, false, dev)

	msg1, _ = hex.DecodeString("029ff8070707070707000400010000c1")
	timeStart := time.Now().Add(5 * time.Minute)
	for {
		err = NewSend(0x9ff8, msg1, dev)
		if err != nil {
			return false
		}
		result, _ := receive(0x9ff0, 24, 5, true, 0, 0x01, false, dev)
		receive(0xa000, 9, 0, false, 0, 0, false, dev)

		if time.Now().After(timeStart) {
			break
		}
		byteSlice := []byte{0x41}

		if equalSlices(result, byteSlice) {
			break
		}
	}
	if time.Now().After(timeStart) || timeStart == time.Now() {
		fmt.Println("Timeout")
		return false
	}
	msg1, _ = hex.DecodeString("029ff80404040404040004000c0000c1")
	err = NewSend(0x9ff8, msg1, dev)
	if err != nil {
		return false
	}
	receive(0x9ff0, 24, 5, true, 0, 0x0c, false, dev)
	receive(0xa000, 20, 0, false, 0, 0, false, dev)
	return true
}

func Download_End(dev *gousb.Device) bool {
	msg1, _ := hex.DecodeString("02a0005050505050500000000000fc0300")
	msg2, _ := hex.DecodeString("029ff807070707070700040009000041")
	err := NewSend(0xa000, msg1, dev)
	if err != nil {
		return false
	}
	err = NewSend(0x9ff8, msg2, dev)
	if err != nil {
		return false
	}
	receive(0x9ff0, 24, 5, true, 0, 0x09, false, dev)

	msg1, _ = hex.DecodeString("02a000d0d0d0d0d0d0")
	msg2, _ = hex.DecodeString("029ff807070707070700040001000041")
	err = NewSend(0xa000, msg1, dev)
	if err != nil {
		return false
	}
	err = NewSend(0x9ff8, msg2, dev)
	if err != nil {
		return false
	}
	receive(0x9ff0, 24, 5, true, 0, 0x01, false, dev)

	msg1, _ = hex.DecodeString("029ff8070707070707000400010000c1")

	timeoutStart := time.Now().Add(5 * time.Minute)
	for {
		err = NewSend(0x9ff8, msg1, dev)
		if err != nil {
			return false
		}
		result, _ := receive(0x9ff0, 24, 5, true, 0, 0x01, false, dev)
		receive(0xa000, 9, 0, false, 0, 0, false, dev)

		if time.Now().After(timeoutStart) {
			break
		}
		byteSlice := []byte{0x41}

		if equalSlices(result, byteSlice) {
			break
		}
	}
	if time.Now().After(timeoutStart) || timeoutStart == time.Now() {
		fmt.Println("Timeout")
		return false
	}
	msg1, _ = hex.DecodeString("029ff80404040404040004000c0000c1")
	err = NewSend(0x9ff8, msg1, dev)
	if err != nil {
		return false
	}
	receive(0x9ff0, 24, 5, true, 0, 0x0c, false, dev)
	receive(0xa000, 20, 0, false, 0, 0, false, dev)

	return true
}
