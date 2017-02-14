/*
* @Author: Ujjwal Kanth
* @Email: ujjwal.kanth@unbxd.com
* @File:
* @Description:
 */

package network

// Defines the state of the server
const (
	UP = iota
	DOWN
)

// Server is a wrapper around echo.Echo, with addtional functionalities
type Server struct {
	name  string
	port  string
	state int
}
