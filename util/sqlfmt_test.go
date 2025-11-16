package util_test

import (
	"fmt"
	"testing"

	"github.com/neura-flow/common/util"
)

func TestFormatSql(t *testing.T) {
	sql := `
		select
                xxx
                , xxx
                , xxx,
                xxx

        from xxx a
        join xxx b
        on xxx = xxx
        join xxx xxx
        on xxx = xxx
        left outer join xxx
        on xxx = xxx
        where xxx = xxx
        and xxx = true
        and xxx is null
	`
	v, err := util.FmtSQL(util.DefaultFmtSqlConfig(), []string{sql})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(v)
}
