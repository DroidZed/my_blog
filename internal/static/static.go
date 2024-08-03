package static

import "fmt"

func Asset(name, folder string) string {
	return fmt.Sprintf("/static/%s/%s", folder, name)
}
