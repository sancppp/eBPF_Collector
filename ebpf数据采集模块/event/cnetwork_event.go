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

type CNetwork_event struct {
	Type          string
	Timestamp     uint64
	Pid           uint32
	Comm          string
	Cid           string
	ContainerName string
	Flag          uint8
	Daddr         [4]byte
	Dport         uint16
	Saddr         [4]byte
	Sport         uint16
}

func (CNetwork_event) GetName() string {
	return "CNetwork_event"
}

func (e CNetwork_event) GetTimestamp() uint64 {
	return e.Timestamp
}

func (e CNetwork_event) GetPid() uint32 {
	return e.Pid
}

func (e CNetwork_event) GetComm() string {
	return e.Comm
}
