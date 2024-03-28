package tools

import (
	"Arkadiy_2Service/config"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"os"

	log "github.com/sirupsen/logrus"
)

func InitLogger() error {

	// можно выбрать только один режим среди O_RDONLY, O_WRONLY и O_RDWR
	//O_RDONLY int = syscall.O_RDONLY // открыть файл в режиме чтения
	//O_WRONLY int = syscall.O_WRONLY // открыть файл в режиме записи
	//O_RDWR   int = syscall.O_RDWR   // открыть файл в режиме чтения и записи

	// значения для управления поведением файла
	//O_APPEND int = syscall.O_APPEND // добавлять новые данные в файл при записи
	//O_CREATE int = syscall.O_CREAT  // создать новый файл, если файла не существует
	//O_EXCL   int = syscall.O_EXCL   // используется вместе с O_CREATE и возвращает
	// ошибку, если файл уже существует
	//O_SYNC  int = syscall.O_SYNC  // открыть в режиме синхронного ввода/вывода
	//O_TRUNC int = syscall.O_TRUNC // очистить файл при открытии

	logsFile, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	defer logsFile.Close()

	log.SetReportCaller(true)

	if config.Env.Production {
		log.SetLevel(log.WarnLevel)
		log.SetOutput(logsFile)
		log.SetFormatter(&nested.Formatter{
			ShowFullLevel: true,
			HideKeys:      true,
			FieldsOrder:   []string{"component", "category"},
		})
	} else {
		log.SetLevel(log.DebugLevel)
		log.SetOutput(os.Stdout)
		log.SetFormatter(&nested.Formatter{
			ShowFullLevel: true,
			HideKeys:      true,
			FieldsOrder:   []string{"component", "category"},
		})
	}
	return nil
}
