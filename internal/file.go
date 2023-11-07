package internal

import (
	"fmt"
	"os"
	"strings"
)

func OutputDeathEvent(path string, events []DeathEvent) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, ev := range events {
		tmpStr := strings.Replace(ev.Ability.Name, "{", "", -1)
		tmpStr = strings.Replace(tmpStr, "}", "", -1)
		_, err := file.WriteString(fmt.Sprintf("%s\t%v\n", ev.Name, tmpStr))
		if err != nil {
			return err
		}
	}
	return nil
}
