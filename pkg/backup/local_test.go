package backup

import (
	"strconv"
	"testing"
	"time"

	"github.com/stefanprodan/mgob/pkg/config"
	"github.com/stretchr/testify/assert"
)

func Test_runDump(t *testing.T) {
	cmd := `mongodump --archive=test.gz --gzip --uri "mongodb://user:password@localhost:27017" --db test --authenticationDatabase admin`
	retryPlan := config.Retry{
		Attempts:      3,
		BackoffFactor: 0,
	}
	archive := "test.gz"
	retryAttempt := 0.0
	timeout := time.Duration(1) * time.Second
	_, retryCount, err := runDump(cmd, retryPlan, archive, retryAttempt, timeout)
	assert.Error(t, err)
	assert.Equal(t, retryPlan.Attempts, int(retryCount))
}

func Test_getDumpedDocMap(t *testing.T) {
	dumpOutput := []byte(`2022-09-12T23:00:00.102+0000	writing DBCollection.Contents_Published to archive '/tmp/content-1663023600.gz'
	2022-09-12T23:00:00.102+0000	writing DBCollection.Contents_Count to archive '/tmp/content-1663023600.gz'
	2022-09-12T23:00:00.102+0000	writing DBCollection.Contents_Published_Count to archive '/tmp/content-1663023600.gz'
	2022-09-12T23:00:00.116+0000	writing DBCollection.Contents_All to archive '/tmp/content-1663023600.gz'
	2022-09-12T23:00:00.128+0000	done dumping DBCollection.Contents_Published_Count (1 document)
	2022-09-12T23:00:00.140+0000	done dumping DBCollection.Contents_Count (3 documents)
	2022-09-12T23:00:00.565+0000	done dumping DBCollection.Contents_Published (7415 documents)
	2022-09-12T23:00:03.010+0000	[####################....]  DBCollection.Contents_All  84357/99427  (84.8%)
	2022-09-12T23:00:03.171+0000	[########################]  DBCollection.Contents_All  99427/99427  (100.0%)
	2022-09-12T23:00:03.294+0000	done dumping DBCollection.Contents_All (99427 documents)
	`)
	result := getDumpedDocMap(string(dumpOutput))
	assert.Len(t, result, 4)
	assert.Equal(t, strconv.Itoa(7415), result["Contents_Published"])
}
