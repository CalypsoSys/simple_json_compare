# simple_json_compare
A very rudimentary JSON compare package

Sample code

```
package main

import (
	"fmt"
	"github.com/CalypsoSys/simple_json_compare"
)

func main() {
	ignorePaths := []string{"root->target->end_time->*", "root->data->*->pct_used->*"}
	diff, differences, err := CompareJSONFiles("test1.json", "test2.json", ignorePaths)
	if err != nil {
		fmt.Printf("Errors: %v+\n", err)
	} else if diff == false {
		fmt.Println("JSON is idential")
	} else {
		fmt.Println("JSON Differences:")
		for _, difference := range differences {
			fmt.Printf("\t%s\n", difference)
		}
	}
}
```



