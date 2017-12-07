package main

import (
	"log"
	"time"

	"./mbserver"
)

func main() {

	serv := mbserver.NewServer()

	// Override ReadDiscreteInputs function.
	serv.RegisterFunctionHandler(2,
		func(s *mbserver.Server, frame mbserver.Framer) ([]byte, *mbserver.Exception) {
			register, numRegs, endRegister := mbserver.RegisterAddressAndNumber(frame)

			log.Printf("%v\n", register)
			log.Printf("%v\n", numRegs)
			log.Printf("%v\n", endRegister)

			// Check the request is within the allocated memory
			if endRegister > 65535 {
				return []byte{}, &mbserver.IllegalDataAddress
			}
			dataSize := numRegs / 8
			if (numRegs % 8) != 0 {
				dataSize++
			}
			data := make([]byte, 1+dataSize)
			data[0] = byte(dataSize)
			for i := range s.DiscreteInputs[register:endRegister] {
				// Return all 1s, regardless of the value in the DiscreteInputs array.
				shift := uint(i) % 8
				data[1+i/8] |= byte(1 << shift)
			}


			return data, &mbserver.Success
		})

	// Start the server.
	err := serv.ListenTCP("0.0.0.0:1502")
	if err != nil {
		log.Printf("%v\n", err)
		return
	}
	defer serv.Close()

	// Wait for the server to start
	time.Sleep(1 * time.Millisecond)

	// Wait forever
	for {
		time.Sleep(1 * time.Second)
	}
}
