package main

import (
	"bufio"
	"log"
	"os"

	"github.com/nmi/gokvm/kvm"

	// change to own library.
	"golang.org/x/term"
)

func main() {
	g, err := kvm.NewLinuxGuest("./bzImage", "./initrd")
	if err != nil {
		panic(err)
	}

	go func() {
		if err = g.RunInfiniteLoop(); err != nil {
			panic(err)
		}
	}()

	// fd 0 is stdin
	state, err := term.MakeRaw(0)
	if err != nil {
		log.Fatalln("setting stdin to raw:", err)
	}

	defer func() {
		if err := term.Restore(0, state); err != nil {
			log.Println("warning, failed to restore terminal:", err)
		}
	}()

	var before byte = 0

	in := bufio.NewReader(os.Stdin)

	for {
		b, err := in.ReadByte()
		if err != nil {
			log.Println("stdin:", err)

			break
		}
		g.GetInputChan() <- b

		if len(g.GetInputChan()) > 0 {
			g.InjectSerialIRQ()
		}

		if before == 0x1 && b == 'x' {
			break
		}

		before = b
	}
}