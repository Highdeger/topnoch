package xpath

import (
	"fmt"
)

func GetInternalDatabaseFilepath() string {
	return fmt.Sprintf("data/internal.db")
}

func GetConfigDatabaseFilepath() string {
	return fmt.Sprintf("data/config.db")
}

func GetNodeDatabaseFilepath(nodeKey string) string {
	return fmt.Sprintf("data/%s.db", nodeKey)
}
