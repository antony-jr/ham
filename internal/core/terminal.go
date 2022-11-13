package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/creack/pty"
)

type Terminal struct {
	cmd   *exec.Cmd
	term  *os.File
	index int
	uid   string
}

func NewTerminal(UniqueID string) (Terminal, error) {
	t := Terminal{}

	t.uid = UniqueID

	t.cmd = exec.Command("bash")
	ptmx, err := pty.Start(t.cmd)
	if err != nil {
		return t, err
	}
	t.term = ptmx

	file, err := os.Create(fmt.Sprintf("/tmp/%s.ham.command.status", UniqueID))
	if err != nil {
		return t, err
	}

	file.Write([]byte("-1 success"))
	file.Close()

	t.term.Write([]byte("set -e \n"))
	t.term.Write([]byte("export HAM_CMD_INDEX=0 \n"))
	// t.term.Write([]byte(fmt.Sprintf("trap \"echo $HAM_CMD_INDEX' failed' > /tmp/%s.ham.command.status\" ERR \n", UniqueID)))
	t.term.Write([]byte(fmt.Sprintf("trap \"env | grep HAM_CMD_INDEX | cut -d'=' -f2 | sed 's/$/ failed/'|cat > /tmp/%s.ham.command.status\" ERR \n", UniqueID)))
	t.term.Write([]byte(fmt.Sprintf("trap \"env | grep HAM_CMD_INDEX | cut -d'=' -f2 | sed 's/$/ failed/'|cat > /tmp/%s.ham.command.status\" EXIT \n", UniqueID)))

	return t, nil
}

func (Term *Terminal) WaitTerminal(Index int) error {
	for {
		time.Sleep(1 * time.Second)
		status, err := ioutil.ReadFile(fmt.Sprintf("/tmp/%s.ham.command.status", Term.uid))
		if err != nil {
			return err
		}

		parts := strings.Split(string(status[:]), " ")
		idx, err := strconv.Atoi(parts[0])
		if err != nil {
			return err
		}
		wStatus := parts[1]

		if idx == Index {
			if !strings.Contains(wStatus, "success") {
				return errors.New("Command Failed")
			}

			return nil
		} else {
			time.Sleep(10 * time.Second)
		}
	}
}

func (Term *Terminal) ExecTerminal(Index int, Command string) error {
	status, err := ioutil.ReadFile(fmt.Sprintf("/tmp/%s.ham.command.status", Term.uid))
	if err != nil {
		return err
	}

	parts := strings.Split(string(status[:]), " ")
	previdx, err := strconv.Atoi(parts[0])
	if err != nil {
		return err
	}
	prevStatus := parts[1]

	if !strings.Contains(prevStatus, "success") {
		return errors.New("Previous Command Failed")
	}

	if previdx > Index {
		return errors.New("Out of Order Execution")
	}

	Term.term.Write([]byte(fmt.Sprintf("export HAM_CMD_INDEX=%d \n", Index)))
	Term.term.Write([]byte(fmt.Sprintf("%s\n", Command)))

	time.Sleep(1 * time.Second)
	status, err = ioutil.ReadFile(fmt.Sprintf("/tmp/%s.ham.command.status", Term.uid))
	if err != nil {
		return err
	}

	parts = strings.Split(string(status[:]), " ")
	idx, err := strconv.Atoi(parts[0])
	if err != nil {
		return err
	}
	wStatus := parts[1]

	if idx == Index {
		if strings.Contains(wStatus, "failed") {
			return errors.New("Command Failed")
		}
	} 

	Term.term.Write([]byte(fmt.Sprintf("echo $HAM_CMD_INDEX' success' > /tmp/%s.ham.command.status\n", Term.uid)))

	return nil
}

func (Term *Terminal) CloseTerminal() error {
	err := os.Remove(fmt.Sprintf("/tmp/%s.ham.command.status", Term.uid))
	if err != nil {
		return err
	}
	Term.term.Close()
	return nil
}
