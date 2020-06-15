package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"go.bug.st/serial"
)

func main() {
	fmt.Println("Terminal...")
	baudPtr := flag.Int("baud", 38400, "Baud Rate de comunicação")
	portPtr := flag.String("port", "/dev/tnt1", "Porta serial")
	hexPtr := flag.Bool("hex", false, "Habilitar saída de dados hexadecimal")
	flag.Parse()

	options := &serial.Mode{
		BaudRate: *baudPtr,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}
	dev, err := serial.Open(*portPtr, options)
	if err != nil {
		fmt.Println("ERRO OPENPORT", err.Error())
	}

	exitCh := make(chan []byte)
	entryCh := make(chan byte)

	go receiveFromKeyboard(exitCh)
	go receiveFromSerial(dev, entryCh)
	for {
		time.Sleep(time.Millisecond)
		select {
		case msg := <-exitCh:
			_, err := dev.Write(msg)
			if err != nil {
				fmt.Println("Erro ao escrever na serial:", err.Error())
			}
		case b := <-entryCh:
			if *hexPtr {
				fmt.Printf("%02X ", b)
			} else {
				fmt.Printf("%v ", b)
			}

		}
	}
}

func receiveFromKeyboard(ch chan []byte) {
	teclado := bufio.NewReader(os.Stdin)

	for {
		time.Sleep(time.Millisecond)
		linha, err := teclado.ReadString('\n')
		if err != nil {
			fmt.Println("Erro ao ler uma linha:", err.Error())
		}
		concat := strings.Join(strings.Fields(linha), "")
		arr, err := hex.DecodeString(concat)
		if err != nil {
			fmt.Println("Erro ao decodificar hex:", err.Error())
		}
		fmt.Println(arr)
		ch <- arr
	}
}

func receiveFromSerial(dev io.Reader, ch chan byte) {
	serial := bufio.NewReader(dev)
	for {
		time.Sleep(time.Millisecond)
		b, err := serial.ReadByte()
		if err != nil {
			fmt.Println("Erro ao ler byte:", err.Error())
		}
		ch <- b
	}

}
