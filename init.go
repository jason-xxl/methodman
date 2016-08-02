package methodman

import "github.com/myteksi/go/commons/util/type/methodtool"

var (
	// queueLength is max element number in a queue. you can change it if you really
	// need larger queue size.
	queueLength = 200

	// managerMap stores mapping of method -> manager
	managerMap = ManagerMapNew()

	// queueKeyPrefix is used to identify queueMap in goroutine local storage
	queueKeyPrefix = "_methodman_queue_132435_"
)

func getFullKey(method interface{}) (fullKey string) {
	fullKey = queueKeyPrefix + string(methodtool.GetMethodUniqueID(method))
	return
}

// SetQueueLength ...
func SetQueueLength(newLength int) {
	if newLength < queueLength {
		panic("newLength<queueLength")
	}
	queueLength = newLength
}
