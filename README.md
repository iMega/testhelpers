# Test helpers

-   DB

### Usage

```go
import (
    "github.com/imega/txwrapper"
    "github.com/imega/testhelpers/db"
)

func Test_MyTest(t *testing.T) {
    tx := func(ctx context.Context, tx *sql.Tx) error {
        // your code

        return nil
    }

    db, closedb, err := db.Create("dbname", tx)
    if err != nil {
        t.Fatalf("failed to create db, %s", err)
    }
    defer closedb()

    // ...
```
