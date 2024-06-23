package static

import "fmt"

func Asset(name string) string {
	return fmt.Sprintf("/public/assets/%s", name)
}
