PACKAGE DOCUMENTATION

package driver
    import "./driver"



CONSTANTS

const (
    N_FLOORS  = 4
    N_BUTTONS = 3

    BTN_UP   = 0
    BTN_DOWN = 1
    BTN_CMD  = 2

    DIR_DOWN = -1
    DIR_STOP = 0
    DIR_UP   = 1
)


FUNCTIONS

func GetFloorSignal() int

func Init() int

func ReadButton(button int, floor int) bool

func RunBottomFloor()

func RunDown()

func RunStop()

func RunTopFloor()

func RunUp()

func SetButtonLight(button int, floor int, value bool)

func SetFloorIndicator(floor int) bool


SUBDIRECTORIES

	channels
	io
	simulator

