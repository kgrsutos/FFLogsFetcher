package internal

import (
	"fmt"
	"os"
	"strings"
)

func OutputDeathEvent(path string, events []DeathEventOutput) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, ev := range events {
		abn := strings.Replace(ev.AbilityName, "{", "", -1)
		abn = strings.Replace(abn, "}", "", -1)
		_, err := file.WriteString(fmt.Sprintf("%s\t%v\t%v\n", ev.PlayerName, ev.ReportName, abn))
		if err != nil {
			return err
		}
	}
	return nil
}
