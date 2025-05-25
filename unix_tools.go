//go:build (darwin || freebsd || linux) && cgo
// +build darwin freebsd linux
// +build cgo

// code fromm : https://github.com/wfd3/go-groups/tree/master
//				Apache Licence 2.0
//				auth: wdf3
//
// topklean :: adding cache for optimistion
//

package main

// package group

import (
	"fmt"
	"strconv"
	"syscall"
	"unsafe"
)

/*
#include <unistd.h>
#include <sys/types.h>
#include <grp.h>
#include <stdlib.h>
*/
import "C"

// topklean:: adding cache
type lookupGroupCache map[int]Group

var groupNameCache lookupGroupCache

func init() {
	groupNameCache = make(lookupGroupCache)
}

// topklean::

// UnknownGroupIdError is returned by LookupId when a group ID cannot be found.
type UnknownGroupIdError int

func (e UnknownGroupIdError) Error() string {
	return "group: unknown gid " + strconv.Itoa(int(e))
}

// UnknownGroupError is returned by Lookup when a group name cannot be found.
type UnknownGroupError string

func (e UnknownGroupError) Error() string {
	return "group: unknown group " + string(e)
}

// Convert (**char)clist to []string
func convert(clist **C.char) []string {
	var members []string

	p := (*[1 << 30]*C.char)(unsafe.Pointer(clist))
	for i := 0; p[i] != nil; i++ {
		members = append(members, C.GoString(p[i]))
	}

	return members
}

type Group struct {
	Name    string
	Gid     int
	Members []string
}

// Current returns the curreng  group information (from getgid())
func Current() (*Group, error) {
	return lookup(syscall.Getgid(), "", false)
}

// LookupId returns the group information for a specific GID.  If the group cannot be found,
// the error UnknownGroupIdError is returned.
// cache system

// topklean:: adding cache function
func LookupId(gid int) (*Group, error) {
	// cache group name
	group, ok := groupNameCache[gid]
	if ok {
		return &group, nil
	} else {
		group, err := lookup(gid, "", false)
		if err == nil {
			groupNameCache[gid] = *group
			return group, nil
		} else {
			return nil, err
		}
	}
}

// Lookup returns the group information for a specific group name.  If the group cannot be
// found, the error UnknownGroupError is returned.
func Lookup(groupname string) (*Group, error) {
	return lookup(-1, groupname, true)
}

func lookup(gid int, groupname string, lookupByName bool) (*Group, error) {
	var grp C.struct_group
	var result *C.struct_group
	var bufsize C.long

	bufsize = C.sysconf(C._SC_GETGR_R_SIZE_MAX)
	if bufsize == -1 {
		bufsize = 1024
	}
	buf := C.malloc(C.size_t(bufsize))
	defer C.free(buf)

	var rv C.int
	if lookupByName {
		CGroup := C.CString(groupname)
		defer C.free(unsafe.Pointer(CGroup))
		rv = C.getgrnam_r(CGroup, &grp, (*C.char)(buf),
			C.size_t(bufsize), &result)
		if rv != 0 {
			return nil,
				fmt.Errorf("group: lookup group name %s: %s",
					groupname, syscall.Errno(rv))
		}
		if result == nil {
			return nil, UnknownGroupError(groupname)
		}
	} else {
		rv = C.getgrgid_r(C.gid_t(gid), &grp, (*C.char)(buf),
			C.size_t(bufsize), &result)
		if rv != 0 {
			return nil, fmt.Errorf("group: lookup gid %d: %s",
				gid, syscall.Errno(rv))
		}
		if result == nil {
			return nil, UnknownGroupIdError(gid)
		}
	}

	g := &Group{
		Gid:     int(grp.gr_gid),
		Name:    C.GoString(grp.gr_name),
		Members: convert(grp.gr_mem),
	}

	return g, nil
}
