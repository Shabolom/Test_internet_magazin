package migrate

//import (
//	"Arkadiy_2Service/config"
//	"Arkadiy_2Service/iternal/domain"
//	"github.com/gofrs/uuid"
//	"github.com/jinzhu/gorm"
//	log "github.com/sirupsen/logrus"
//	"gopkg.in/gormigrate.v1"
//
//	_ "github.com/jinzhu/gorm/dialects/postgres"
//)
//
//// Migrate запустите миграцию для всех объектов и добавьте для них ограничения
//// создаем таблицы и закидываем в бд тут
//func Migrate() {
//	db := config.DB
//
//	orderID, _ := uuid.NewV4()
//	productID, _ := uuid.NewV4()
//	paletteID, _ := uuid.NewV4()
//	ordersProducts, _ := uuid.NewV4()
//	palettesProducts, _ := uuid.NewV4()
//
//	// создаем объект миграции данная строка всегда статична (всегда такая)
//	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
//		{
//			// id всех миграций кторые были проведены
//			ID: orderID.String(),
//			Migrate: func(tx *gorm.DB) error {
//				err := tx.AutoMigrate(&domain.Order{}).Error
//				if err != nil {
//					return err
//				}
//				return nil
//			},
//			Rollback: func(tx *gorm.DB) error {
//				err := tx.DropTable("Orders").Error
//				if err != nil {
//					return err
//				}
//				return nil
//			},
//		}, {
//			// id всех миграций кторые были проведены
//			ID: productID.String(),
//			Migrate: func(tx *gorm.DB) error {
//				err := tx.AutoMigrate(&domain.Product{}).Error
//				if err != nil {
//					return err
//				}
//				return nil
//			},
//			Rollback: func(tx *gorm.DB) error {
//				err := tx.DropTable("Products").Error
//				if err != nil {
//					return err
//				}
//				return nil
//			},
//		}, {
//			// id всех миграций кторые были проведены
//			ID: paletteID.String(),
//			Migrate: func(tx *gorm.DB) error {
//				err := tx.AutoMigrate(&domain.Palette{}).Error
//				if err != nil {
//					return err
//				}
//				return nil
//			},
//			Rollback: func(tx *gorm.DB) error {
//				err := tx.DropTable("palette").Error
//				if err != nil {
//					return err
//				}
//				return nil
//			},
//		}, {
//			// id всех миграций кторые были проведены
//			ID: ordersProducts.String(),
//			Migrate: func(tx *gorm.DB) error {
//				err := tx.AutoMigrate(&domain.OrdersProducts{}).Error
//				if err != nil {
//					return err
//				}
//				return nil
//			},
//			Rollback: func(tx *gorm.DB) error {
//				err := tx.DropTable("orders_products").Error
//				if err != nil {
//					return err
//				}
//				return nil
//			},
//		}, {
//			// id всех миграций кторые были проведены
//			ID: palettesProducts.String(),
//			Migrate: func(tx *gorm.DB) error {
//				err := tx.AutoMigrate(&domain.PalettesProducts{}).Error
//				if err != nil {
//					return err
//				}
//				return nil
//			},
//			Rollback: func(tx *gorm.DB) error {
//				err := tx.DropTable("palettes_products").Error
//				if err != nil {
//					return err
//				}
//				return nil
//			},
//		},
//	})
//
//	err := m.Migrate()
//	if err != nil {
//		log.WithField("component", "migration").Panic(err)
//	}
//
//	if err == nil {
//		log.WithField("component", "migration").Info("Migration did run successfully")
//	} else {
//		log.WithField("component", "migration").Infof("Could not migrate: %v", err)
//	}
//}
