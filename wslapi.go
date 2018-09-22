// go-wsl, a Golang interface to Windows Services for Linux
// Copyright (C) 2018  Hilko Bengen
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// Package wsl implements wrapper functions the Windwos Subsystem For
// Linux API as documented in
// https://docs.microsoft.com/en-us/previous-versions/windows/desktop/api/_wsl/
package wsl

import (
	"golang.org/x/sys/windows"
	"reflect"
	"unsafe"
)

type DistributionFlags uint32

const (
	DISTRIBUTION_FLAGS_NONE DistributionFlags = iota
	DISTRIBUTION_FLAGS_ENABLE_INTEROP
	DISTRIBUTION_FLAGS_APPEND_NT_PATH
	DISTRIBUTION_FLAGS_ENABLE_DRIVE_MOUNTING
)

//sys	coTaskMemFree(p unsafe.Pointer) (err error) = Ole32.CoTaskMemFree

//sys	configureDistribution(distributionName *uint16, defaultUID uint32, wslDistributionFlags uint32) (err error) = wslapi.WslConfigureDistribution

// ConfigureDistribution Modifies the behavior of a distribution
// registered with the Windows Subsystem for Linux.
//
// See https://docs.microsoft.com/en-us/previous-versions/windows/desktop/api/wslapi/nf-wslapi-wslconfiguredistribution
func ConfigureDistribution(name string, defaultUID uint32, flags DistributionFlags) error {
	n, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return err
	}
	return configureDistribution(n, defaultUID, uint32(flags))
}

//sys	getDistributionConfiguration(distributionName *uint16, distributionVersion *uint32, defaultUID *uint32,  wslDistributionFlags *uint32, defaultEnvironmentVariables ***uint16, defaultEnvironmentVariableCount *uint32) (err error) = wslapi.WslGetDistributionConfiguration

// GetDistributionConfiguration retrieves the current configuration of
// a distribution registered with the Windows Subsystem for Linux.
//
// See https://docs.microsoft.com/en-us/previous-versions/windows/desktop/api/wslapi/nf-wslapi-wslgetdistributionconfiguration
func GetDistributionConfiguration(name string) (version uint32, defaultUID uint32, flags DistributionFlags, environment []string, err error) {
	var tmpEnv **uint16
	var envCount uint32
	var tmpName *uint16
	if tmpName, err = windows.UTF16PtrFromString(name); err != nil {
		return
	}
	if err = getDistributionConfiguration(tmpName, &version, &defaultUID, (*uint32)(&flags), &tmpEnv, &envCount); err != nil {
		return
	}
	for e, i := uintptr(unsafe.Pointer(tmpEnv)), uintptr(0); i < uintptr(envCount); i++ {
		var tmpUTF16 []uint16
		p := e + i*unsafe.Sizeof(tmpEnv)
		hdr := (*reflect.SliceHeader)(unsafe.Pointer(&tmpUTF16))
		// Assume that individual environment strings will not be
		// larger than 4096.
		hdr.Data, hdr.Len, hdr.Cap = p, 4096, 4096
		environment = append(environment, windows.UTF16ToString(tmpUTF16))
		coTaskMemFree(unsafe.Pointer(p))
	}
	return
}

//sys	isDistributionRegistered(distributionName *uint16) (rv bool) = wslapi.WslIsDistributionRegistered

// IsDistributionRegistered determines if a distribution is registered with the Windows Subsystem for Linux.
//
// See https://docs.microsoft.com/en-us/previous-versions/windows/desktop/api/wslapi/nf-wslapi-wslisdistributionregistered
func IsDistributionRegistered(name string) bool {
	n, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return false
	}
	return isDistributionRegistered(n)
}

//sys	launch(distributionName *uint16, command *uint16, useCurrentWorkingDirectory bool, stdIn windows.Handle, stdOut windows.Handle, stdErr windows.Handle, process *windows.Handle) (err error) = wslapi.WslLaunch

// Launch launches a Windows Subsystem for Linux (WSL) process in the context of a particular distribution.
//
// See https://docs.microsoft.com/en-us/previous-versions/windows/desktop/api/wslapi/nf-wslapi-wsllaunch
func Launch(name string, command string, useCwd bool, stdin, stdout, stderr windows.Handle) (process windows.Handle, err error) {
	var n, c *uint16
	if n, err = windows.UTF16PtrFromString(name); err != nil {
		return
	}
	if c, err = windows.UTF16PtrFromString(command); err != nil {
		return
	}
	launch(n, c, useCwd, stdin, stdout, stderr, &process)
	return
}

//sys	launchInteractive(distributionName *uint16, command *uint16, useCurrentWorkingDirectory bool, exitCode *uint32) (err error) = wslapi.WslLaunchInteractive

// LaunchInteractive Launches an interactive Windows Subsystem for
// Linux (WSL) process in the context of a particular distribution.
// This differs from Launch in that the end user will be able to
// interact with the newly-created process.
//
// See https://docs.microsoft.com/en-us/previous-versions/windows/desktop/api/wslapi/nf-wslapi-wsllaunchinteractive
func LaunchInteractive(name string, command string, useCwd bool) (exitCode uint32, err error) {
	var n, c *uint16
	if n, err = windows.UTF16PtrFromString(name); err != nil {
		return
	}
	if c, err = windows.UTF16PtrFromString(command); err != nil {
		return
	}
	launchInteractive(n, c, useCwd, &exitCode)
	return
}

//sys	registerDistribution(distributionName *uint16, tarGzFilename *uint16) (err error) = wslapi.WslRegisterDistribution

// RegisterDistribution registers a new distribution with the Windows
// Subsystem for Linux.
//
// See https://docs.microsoft.com/en-us/previous-versions/windows/desktop/api/wslapi/nf-wslapi-wslregisterdistribution
func RegisterDistribution(name string, tarball string) (err error) {
	var n, t *uint16
	if n, err = windows.UTF16PtrFromString(name); err != nil {
		return
	}
	if t, err = windows.UTF16PtrFromString(tarball); err != nil {
		return
	}
	return registerDistribution(n, t)
}

//sys	unregisterDistribution(distributionName *uint16) (err error) = wslapi.WslUnregisterDistribution

// UnregisterDistribution unregisters a distribution from the Windows
// Subsystem for Linux.
//
// See https://docs.microsoft.com/en-us/previous-versions/windows/desktop/api/wslapi/nf-wslapi-wslunregisterdistribution
func UnregisterDistribution(name string, tarball string) (err error) {
	var n *uint16
	if n, err = windows.UTF16PtrFromString(name); err != nil {
		return
	}
	return unregisterDistribution(n)
}
