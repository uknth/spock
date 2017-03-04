/*
* @Author: Ujjwal Kanth
*
* @Email: ujjwal.kanth@unbxd.com
* @Project: spock
* @File: src/engine/ha.go
* @Description:
 */

package engine

// HealthCheck exposes interface for various types of HealthCheck
type HealthCheck interface {
	Status() uint
}
