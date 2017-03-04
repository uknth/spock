/*
* @Author: Ujjwal Kanth
*
* @Email: ujjwal.kanth@unbxd.com
* @Project: spock
* @File: src/ctx/logrus.go
* @Description:
 */

package ctx

import (
	"common/logrus/hook"
	"common/writer"
	"config"

	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	logrusSection = "logrus"

	// Formatter
	formatterTypeOpt     = "formatter-type"
	defaultFormatterType = "text"

	textFormatterType = "text"
	jsonFormatterType = "json"

	// Console
	consoleOpt = "console"

	// File Writer
	filePathOpt     = "file-path"
	defaultFilePath = "logrus.log"

	// Kafka Keys
	kafkaHookEnabledOpt   = "hook-kafka-enabled"
	kafkaHookIdOpt        = "hook-kafka-id"
	kafkaDefaultHookId    = "kafka-hook"
	kafkaDefaultTopicsOpt = "hook-kafka-default-topics"
	kafkaBrokersOpt       = "hook-kafka-brokers"
)

// Logrus Ctx Loader
// This sets logrus's global configuration. If a package needs a custom
// logger, use `logger.New()` method in logrus package
type logrusCtx struct {
	name string
}

// Initialize logrus
func (lc *logrusCtx) Init(cf config.Conf) error {
	// Load Properties
	// Formatter
	lc.setFormatter(cf)
	// Output
	lc.setOutput(cf)
	return nil
}
func (lc *logrusCtx) setFormatter(cf config.Conf) {
	switch cf.Section(logrusSection).
		Key(formatterTypeOpt).
		MustString(defaultFormatterType) {
	case textFormatterType:
		log.SetFormatter(&log.TextFormatter{})
	case jsonFormatterType:
		log.SetFormatter(&log.TextFormatter{})
	default:
		log.SetFormatter(&log.TextFormatter{})
	}
}

func (lc *logrusCtx) setHooks(cf config.Conf) error {
	section := cf.Section(logrusSection)
	if section.Key(kafkaHookEnabledOpt).
		MustBool(false) {
		brokers := section.Key(kafkaBrokersOpt).StringS(",")
		topics := section.Key(kafkaDefaultTopicsOpt).StringS(",")
		// Add Kafka hook
		kafkaHook, err := hook.NewKafkaHook(
			section.Key(kafkaHookIdOpt).MustString(kafkaDefaultHookId),
			[]log.Level{
				log.ErrorLevel,
				log.InfoLevel,
			},
			&log.JSONFormatter{},
			brokers,
			topics,
		)
		if err != nil {
			return errors.Wrap(err, "Error Initializing Kafka HOOK")
		}
		// Add to logrus
		log.AddHook(kafkaHook)
	}
	return nil
}

func (lc *logrusCtx) setOutput(cf config.Conf) error {
	// config - console = false by default
	if cf.Section(logrusSection).Key(consoleOpt).MustBool(false) {
		// Write it in console
		log.SetOutput(os.Stdout)
		return nil
	}
	// Open a file by reading file path from config
	filePath := cf.Section(logrusSection).Key(filePathOpt).
		MustString(defaultFilePath)

	// Rotate Writer
	wr, err := writer.NewWriter(filePath)
	if err != nil {
		return errors.Wrap(err, "Error opening file:"+filePath)
	}
	log.SetOutput(wr)
	return nil
}

func (lc *logrusCtx) Name() string {
	return lc.name
}
