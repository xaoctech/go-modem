// SPDX-License-Identifier: MIT
//
// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.

// modeminfo collects and displays information related to the modem and its
// current configuration.
//
// This serves as an example of how interact with a modem, as well as
// providing information which may be useful for debugging.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/warthog618/modem/at"
	"github.com/warthog618/modem/serial"
	"github.com/warthog618/modem/trace"
)

func main() {
	dev := flag.String("d", "/dev/ttyUSB0", "path to modem device")
	baud := flag.Int("b", 115200, "baud rate")
	timeout := flag.Duration("t", 400*time.Millisecond, "command timeout period")
	verbose := flag.Bool("v", false, "log modem interactions")
	flag.Parse()
	m, err := serial.New(serial.WithPort(*dev), serial.WithBaud(*baud))
	if err != nil {
		log.Println(err)
		return
	}
	defer m.Close()
	var mio io.ReadWriter = m
	if *verbose {
		mio = trace.New(m)
	}
	a := at.New(mio)
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	err = a.Init(ctx)
	cancel()
	if err != nil {
		log.Println(err)
		return
	}
	cmds := []string{
		"I",
		"+GCAP",
		"+CMEE=2",
		"+CGMI",
		"+CGMM",
		"+CGMR",
		"+CGSN",
		"+CSQ",
		"+CIMI",
		"+CREG?",
		"+CNUM",
		"+CPIN?",
		"+CEER",
		"+CSCA?",
		"+CSMS?",
		"+CSMS=?",
		"+CPMS=?",
		"+CNMI?",
		"+CNMI=?",
		"+CNMA=?",
		"+CMGF=?",
	}
	for _, cmd := range cmds {
		ctx, cancel := context.WithTimeout(context.Background(), *timeout)
		info, err := a.Command(ctx, cmd)
		cancel()
		fmt.Println("AT" + cmd)
		if err != nil {
			fmt.Printf(" %s\n", err)
			continue
		}
		for _, l := range info {
			fmt.Printf(" %s\n", l)
		}

	}
}
