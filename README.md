# simple_json_compare
A very rudimentary JSON compare package

Ignore path that end with "->*" will ignore last element
Example:
Ignore Spec ==> "root->target->end_time->*"
Ignored     ==> "root->target->end_time"
Ignored     ==> "root->target->end_time->08:13:13"
Not Ignored ==> "root->target->end_time->date->08-13-1988"

Ignore path that ends with "->*R[regex]" will ignore last element if it matches the supplied regular expression
Ignore Spec ==> "root->address->*->connectivity->*->ips->*R[->\\d+->\\d+\\.\\d+\\.\\d+\\.\\d+;\\d+]"
Ignored     ==> "root->address->*->connectivity->*->ips->0->169.254.169.254;8"
Not Ignored ==> "root->address->*->connectivity->*->ips->0->169.254.169.254;8a"

Sample code

```
package main

import (
	"fmt"
	"github.com/CalypsoSys/simple_json_compare"
)

func main() {
	ignorePaths := []string{"root->target->end_time->*", "root->data->*->pct_used->*"}
	diff, differences, err := simple_json_compare.CompareJSONFiles("test1.json", "test2.json", ignorePaths)
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



