Что выведет программа? Объяснить вывод программы.

```go
package main

import (
    "fmt"
)

func main() {
    a := [5]int{76, 77, 78, 79, 80}
    var b []int = a[1:4]
    fmt.Println(b)
}
```

Ответ:
```
Ответ будет - 77, 78, 79. Индексы начинаются с 0.
При создании слайса таким образом, последний индекс включен не будет.
```