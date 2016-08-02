package methodman

import (
	"testing"
)

var (
	Method1 = func(p1 int, p2 string) (r1 string) {
		return "I'm original Method1"
	}

	Method2 = func(p1 int) (r1 string) {
		return "I'm original Method2"
	}

	Method3 = Method1
)

func TestTypeTool(t *testing.T) {

	idA := GetMethodUniqueID(Method1)
	idB := GetMethodUniqueID(Method2)
	idC := GetMethodUniqueID(Method3)

	idAPointer := GetMethodUniqueID(&Method1)

	pointerMethod1 := &Method1
	idAPointerPointer := GetMethodUniqueID(&pointerMethod1)

	if idA != GetMethodUniqueID(Method1) {
		t.Fatal("MethodUniqueID should not change.")
	}
	if idA == idB {
		t.Fatal("MethodUniqueID should not repeat.")
	}
	if idA != idC {
		t.Fatal("MethodUniqueID should not change when being copied.")
	}
	if idA == "" {
		t.Fatal("MethodUniqueID should not be empty.")
	}
	if idA != idAPointer {
		t.Fatal("MethodUniqueID should not be affected by indirecting pointer.")
	}
	if idAPointer != idAPointerPointer {
		t.Fatal("MethodUniqueID should not be affected by indirecting multiple pointer.")
	}
	if GetMethodType(Method1) != GetMethodType(pointerMethod1) {
		t.Fatal("GetMethodType should not be affected by indirecting pointer.")
	}
	if !IsMethod(Method1) {
		t.Fatal("IsMethod failed to identify.")
	}
	if IsMethod(pointerMethod1) {
		t.Fatal("IsMethod failed to identify.")
	}
	if !IsMethodPointer(pointerMethod1) {
		t.Fatal("IsMethodPointer failed to identify.")
	}
	if IsMethodPointer(Method1) {
		t.Fatal("IsMethodPointer failed to identify.")
	}
}
