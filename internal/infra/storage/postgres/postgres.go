package postgres

import (
	"fmt"
	"github.com/vaberof/go-chat/pkg/xtime/xlocation"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

const novosibirskTimeZone = "Asia/Novosibirsk"

func New(config *PostgresDatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprint(
		"postgres://" + config.User +
			":" + config.Password +
			"@" + config.Host +
			":" + fmt.Sprintf("%d", config.Port) +
			"/" + config.Name)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{NowFunc: currentTimeWithTimezone})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func currentTimeWithTimezone() time.Time {
	novosibirsk := xlocation.Must(novosibirskTimeZone)
	return time.Now().In(novosibirsk)
}
