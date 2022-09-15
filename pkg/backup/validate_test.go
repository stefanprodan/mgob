package backup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_checkIfAnyFailure_No_Failure(t *testing.T) {
	output := `2022-09-15T19:18:00.068+0000	preparing collections to restore from
	2022-09-15T19:18:00.078+0000	reading metadata for test_restore.test from archive '/tmp/mongo-test-1663269480.gz'
	2022-09-15T19:18:00.085+0000	restoring test_restore.test from archive '/tmp/mongo-test-1663269480.gz'
	2022-09-15T19:18:00.096+0000	finished restoring test_restore.test (1 document, 0 failures)
	2022-09-15T19:18:00.096+0000	no indexes to restore for collection test_restore.test
	2022-09-15T19:18:00.096+0000	1 document(s) restored successfully. 0 document(s) failed to restore.
	`
	assert.NoError(t, checkIfAnyFailure(output))
}

func Test_checkIfAnyFailure_No_Success(t *testing.T) {
	output := `2022-09-15T17:04:00.096+0000\	The --db and --collection flags are deprecated for this use-case; please use --nsInclude instead, i.e. with --nsInclude=${DATABASE}.${COLLECTION}
	2022-09-15T17:04:00.115+0000\	preparing collections to restore from
	2022-09-15T17:04:00.135+0000\	0 document(s) restored successfully. 0 document(s) failed to restore.
	`
	assert.Error(t, checkIfAnyFailure(output))
}

func Test_checkIfAnyFailure_Failure(t *testing.T) {
	output := `2022-09-15T16:46:00.065+0000	The --db and --collection flags are deprecated for this use-case; please use --nsInclude instead, i.e. with --nsInclude=${DATABASE}.${COLLECTION}
	2022-09-15T16:46:00.074+0000	preparing collections to restore from
	2022-09-15T16:46:00.086+0000	reading metadata for test.test from archive '/tmp/mongo-test-1663260360.gz'
	2022-09-15T16:46:00.087+0000	restoring to existing collection test.test without dropping
	2022-09-15T16:46:00.087+0000	restoring test.test from archive '/tmp/mongo-test-1663260360.gz'
	2022-09-15T16:46:00.089+0000	continuing through error: E11000 duplicate key error collection: test.test index: _id_ dup key: { _id: ObjectId('6323560453723833469b3091') }
	2022-09-15T16:46:00.098+0000	finished restoring test.test (0 documents, 1 failure)
	2022-09-15T16:46:00.098+0000	no indexes to restore for collection test.test
	2022-09-15T16:46:00.098+0000	0 document(s) restored successfully. 1 document(s) failed to restore.`
	assert.Error(t, checkIfAnyFailure(output))
}

func Test_checkRetoreDatabase(t *testing.T) {
	backupResult := map[string]string{
		"test":  "1",
		"test2": "1",
	}
	collectionNames := []string{"test", "test2"}

	assert.NoError(t, checkRetoreDatabase(backupResult, collectionNames))
}

func Test_checkRetoreDatabase_Missing_In_Restore(t *testing.T) {
	backupResult := map[string]string{
		"test":  "1",
		"test2": "1",
	}
	collectionNames := []string{"test"}

	assert.Error(t, checkRetoreDatabase(backupResult, collectionNames))
}

func Test_checkRetoreDatabase_Mismatch_In_Restore(t *testing.T) {
	backupResult := map[string]string{
		"test":  "1",
		"test2": "1",
	}
	collectionNames := []string{"test2"}

	assert.Error(t, checkRetoreDatabase(backupResult, collectionNames))
}
