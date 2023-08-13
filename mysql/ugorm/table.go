// Package ugorm 包描述
// Author: wanlizhan
// Date: 2023/7/2
package ugorm

type FieldInfo struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default string
	Extra   string
}

type TableInfo struct {
	TableName  string
	FieldInfos []FieldInfo
}

// GetTables 获取所有表
func (g *GormMysql) GetTables() ([]string, error) {
	tables := make([]string, 0)
	tx := g.Db.Raw("SHOW TABLES").Scan(&tables)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tables, nil
}

func (g *GormMysql) GetTableInfo(table string) (TableInfo, error) {
	var tableInfo TableInfo

	tableInfo.TableName = table
	sql := "Desc" + " " + table

	tx := g.Db.Raw(sql).Scan(&tableInfo.FieldInfos)
	if tx.Error != nil {
		return tableInfo, tx.Error
	}

	return tableInfo, nil
}

func (g *GormMysql) HasTable(table string) bool {
	return g.Db.Migrator().HasTable(table)
}
