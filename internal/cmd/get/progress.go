package get

/*
import (
    "bufio"
    "errors"
    "io"

    "encoding/json"
)

type ProgressT struct {
    error  bool `json:"error"`
    status string `json:"status"`
    progress string `json:"progress"`
    percentage int `json:"percentage"`
    message string `json:"message"`
}

func GetProgress(stdin *io.WriteCloser, stdout *io.Reader, stderr *io.Reader) (ProgressT, error) {
   status_command := []byte("ham build-status | cat |  grep -a Status | cut -c 10-")
   err_prog := ProgressT {
      error: true,
      status: "Failed",
      progress: "Gathering error report",
      percentage: 100,
      message: "Cannot get progress over SSH.",
   }

   _, err := *stdin.Write(status_command)
   if err != nil {
      return err_prog, err 
   }

   scanner := bufio.NewScanner(stdout)
   if tkn := scanner.Scan(); tkn {
      rcv := scanner.Bytes()

      raw := make([]byte, len(rcv))
      copy(raw, rcv)

      var statusJson ProgressT
      err = json.Unmarshal(raw, &statusJson)
      if err != nil {
	 return err_prog, err
      }

      return statusJson, nil
   } else {
      if scanner.Err() != nil {
	 return err_prog, scanner.Err()
      } else {
	 return err_prog, errors.New("EOF")
      }
   }
}

*/
