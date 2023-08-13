// Package ugorm 包描述
// Author: wanlizhan
// Date: 2023/7/2
package ugorm

import (
	"fmt"
	"github.com/xxzhwl/wdk/dict"
	"github.com/xxzhwl/wdk/uconfig"
	"github.com/xxzhwl/wdk/ulog"
	"github.com/xxzhwl/wdk/ustr"
	"gorm.io/gen"
)

// Generator gorm-model生成器
type Generator struct {
	*gen.Generator
	Db *GormMysql
}

// NewSchemaGenerator 根据schema获取生成器
func NewSchemaGenerator(schema string) (*Generator, error) {
	ormGenConf, err := uconfig.ME("OrmGen")
	if err != nil {
		err = fmt.Errorf("NewSchemaGenerator Err:%s", err.Error())
		ulog.Error("GormGen", err.Error())
		return nil, err
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:      dict.S(ormGenConf, "OutPath"),
		ModelPkgPath: dict.S(ormGenConf, "ModelPkgPath"),
		WithUnitTest: false,
		Mode:         gen.WithDefaultQuery | gen.WithoutContext,
	})

	db, err := NewMysqlBySchema(schema)
	if err != nil {
		return nil, err
	}
	g.UseDB(db.Db)
	return &Generator{g, db}, nil
}

// NewGenConfig 获取一个生成器配置
func NewGenConfig(confName string) (gen.Config, error) {
	ormGenConf, err := uconfig.ME(confName)
	if err != nil {
		err = fmt.Errorf("NewDefaultGenerator Err:%s", err.Error())
		ulog.Error("GormGen", err.Error())
		return gen.Config{}, err
	}
	return gen.Config{
		OutPath:      dict.S(ormGenConf, "OutPath"),
		ModelPkgPath: dict.S(ormGenConf, "ModelPkgPath"),
		WithUnitTest: false,
		Mode:         gen.WithDefaultQuery | gen.WithoutContext | gen.WithQueryInterface,
	}, nil
}

// NewDefaultGenerator 获取一个生成器
func NewDefaultGenerator() (*Generator, error) {
	config, err := NewGenConfig("OrmGen")
	if err != nil {
		return nil, err
	}
	g := gen.NewGenerator(config)
	db, err := NewMysqlDefault()
	if err != nil {
		return nil, err
	}
	g.UseDB(db.Db)

	return &Generator{g, db}, nil
}

// NewGenerator 获取一个生成器
func NewGenerator(config gen.Config) (*Generator, error) {
	g := gen.NewGenerator(config)
	db, err := NewMysqlDefault()
	if err != nil {
		return nil, err
	}
	g.UseDB(db.Db) // reuse y

	return &Generator{g, db}, nil
}

// GenModel 生成对象
func (g *Generator) GenModel() error {
	tables, err := g.Db.GetTables()
	if err != nil {
		ulog.ErrorF("GormGen", "GenModel Err:%s", err.Error())
		return err
	}
	g.GenerateBasicModel(tables)
	return nil
}

// GenerateBasicModel 生成基础DAO对象
func (g *Generator) GenerateBasicModel(tables []string) {
	softDeleteField := gen.FieldType("deleted", "gorm.DeletedAt")
	ctTimeType := gen.FieldType("create_time", "ugorm.LocalTime")
	utTimeType := gen.FieldType("update_time", "ugorm.LocalTime")

	for _, table := range tables {
		g.ApplyBasic(g.GenerateModelAs(table, ustr.SnakeToBigCamel(table)+"Dao", softDeleteField, ctTimeType, utTimeType))
	}
	g.Execute()
}
