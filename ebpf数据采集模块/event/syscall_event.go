// Copyright (C) 2024 Tianzhenxiong
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package event

type Syscall_event struct {
	Type          string
	Timestamp     uint64
	Flag          uint8 // 0 for enter, 1 for exit, 2 for count
	Pid           uint32
	Comm          string
	Syscall       uint32
	Ret           int64
	Cid           string
	ContainerName string
	Info          string
}

func (Syscall_event) GetName() string {
	return "System_event"
}

func (e Syscall_event) GetTimestamp() uint64 {
	return e.Timestamp
}

func (e Syscall_event) GetPid() uint32 {
	return e.Pid
}

func (e Syscall_event) GetComm() string {
	return e.Comm
}
