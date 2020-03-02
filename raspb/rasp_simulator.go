package raspb

import (
	"errors"
	"fmt"
	"github.com/oloose/ard"
	"log"
)

var usedRaspbIds uint8 = 0

type Raspberry struct {
	id         uint8
	usedArdIds uint8
	ards       map[uint8]*ard.Arduino
}

//Initializes and returns a new Raspberry instance
//mArds = number of Arduinos to initialize at Raspberry creation
func NewRaspberry(mArds int) *Raspberry {
	rasp := Raspberry{usedRaspbIds, 0, make(map[uint8]*ard.Arduino)}
	usedRaspbIds++
	log.Printf("## CREATED RASPERRY: R%d/...\n", rasp.id)

	if mArds != 0 {
		for i := 0; i < mArds; i++ {
			rasp.AddArd()
		}
	}

	return &rasp
}

//Adds a new Arduino to a Raspberry
func (rRaspb *Raspberry) AddArd() ard.Arduino {
	nArd := ard.NewArduino(rRaspb.usedArdIds) //create new ard
	log.Printf("## CREATED ARDUINO: R%d/Rig%d/...\n", rRaspb.id, nArd.Id())
	rRaspb.ards[rRaspb.usedArdIds] = nArd // add nArd to map with rRaspb.usedArdIds as key
	rRaspb.usedArdIds++                   // increase id counter (always to prevent double used ids)
	return *nArd
}

//Removes the ard with the given mId value from a Raspberry
//mId = id of the Arduino to remove
//Return = Error if key is not present
func (rRaspb *Raspberry) RemoveArd(mId uint8) error {
	if _, b := rRaspb.ards[mId]; !b {
		return errors.New(fmt.Sprintf("Key {%d} not present", mId))
	}
	delete(rRaspb.ards, mId)
	return nil
}

//Returns a reference to all Arduinos a Raspberry has
func (rRaspb *Raspberry) Ards() map[uint8]ard.Arduino {
	var ards = make(map[uint8]ard.Arduino)
	for ai, a := range rRaspb.ards {
		ards[ai] = *a
	}
	return ards
}

//Returns the last used Arduino id value
func (rRaspb *Raspberry) UsedArdIds() uint8 {
	return rRaspb.usedArdIds
}

func (rRaspb *Raspberry) Id() uint8 {
	return rRaspb.id
}

//Returns the last used Raspberry id
func UsedRaspIds() uint8 {
	return usedRaspbIds
}
