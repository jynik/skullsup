// SPDX License Identifier: MIT
package network

type User struct {
	Name       string `json:"name"`        // User name pulled from cert's CN
	CertSerial string `json:"cert_serial"` // User ID pulled from cert's serial #

	ReadQueues  []string `json:"read_queues"`  // Queues user is permitted to read
	WriteQueues []string `json:"write_queues"` // Queues user is permitted to write
}

func (u *User) canAccess(target string, queues []string) bool {
	for _, queue := range queues {
		if queue == target {
			return true
		}
	}
	return false
}

func (u *User) CanRead(target string) bool {
	return u.canAccess(target, u.ReadQueues)
}

func (u *User) CanWrite(target string) bool {
	return u.canAccess(target, u.WriteQueues)
}

type UserList []User
