/*
* @Author: Ujjwal Kanth
*
* @Email: ujjwal.kanth@unbxd.com
* @Project: spock
* @File: src/config/viper.go
* @Description: use spf13/viper to parse the config file
 */

package config

// Check folders for config file
var folders = []string{
	".",
	"./conf",
	"/etc/spock/",
}

// TODO: struct that implements Conf encapsulating Viper

// NewVIPERConf returns a new Conf object using spf13/viper config parser
func NewVIPERConf(filename string) (Conf, error) {
	// Add Config lookup folders
	// for _, folder := range folders {
	// 	viper.AddConfigPath(folder)
	// }

	// Config Filename (Without ext)
	// filename = filename[0 : len(filename)-
	// 	len(filepath.Ext(filename))]

	// Set config name
	// viper.SetConfigName(filename)

	// err := viper.ReadInConfig()
	// if err != nil {
	// 	return nil, errors.Wrap(err, "Viper: Error in reading config")
	// }
	return nil, nil
}
