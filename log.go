package codefixture

import (
	"log"
	"os"
)

func init() {
	log.SetOutput(os.Stderr)
}
