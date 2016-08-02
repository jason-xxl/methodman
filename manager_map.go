package methodman

import "github.com/myteksi/go/commons/util/type/methodtool"

// ManagerMap stores mapping relation from pointer of manager method to its state struct.
// Note that,
// 1) the pointer of original method is in the manager state.
// 2) assuming the map is formed from beginning of test and remain no changed,
// so no need to add sync.RWMutex protection.
type ManagerMap map[methodtool.MethodUniqueID]*Manager

// ManagerMapNew ...
func ManagerMapNew() (o ManagerMap) {
	o = make(ManagerMap)
	return
}
