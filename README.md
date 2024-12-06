# 内存数据库项目

## 简介  
这是一个个人的学习项目，一个高性能、轻量级的内存数据库，专为低延迟、高频率数据操作场景设计。现在功能不多，仅供学习交流使用，数据库通过将数据存储在内存中实现极快的读写速度，非常适合缓存、实时分析和其他对性能要求较高的应用场景。

---

## 特性  
- **高性能**：所有数据操作均在内存中完成，提供极低的延迟。  
- **灵活的数据模型**：支持自定义表结构和简单的关系映射。  
- **轻量级实现**：无外部依赖，适合嵌入式和微服务应用。  
- **【计划中】事务支持**：提供基本的事务机制，确保数据一致性。  
- **【计划中】易用的接口**：提供简单的 API，方便开发者快速集成。  

---

## 功能列表  
1. **基本数据操作**  
   - 插入数据  
   - 查询数据（支持条件过滤）  
   - 更新数据  
   - 【计划中】删除数据  

2. **事务支持**  
   - 提供事务的开始、提交和回滚功能，确保复杂操作的一致性。

3. **索引机制**  
   - 支持简单的索引功能以提升查询性能。  

4. **多表支持**  
   - 可以创建多个表，并支持表之间的简单关联操作。

5. **持久化选项（可选）**  
   - 提供数据快照功能，将内存中的数据保存到磁盘。  
   - 支持从快照中恢复数据。

---

## 使用指南  

### 1. 安装  
将源码克隆到本地并编译即可，无需额外依赖。  
```bash
git clone https://github.com/cloakscn/go-redis.git
go build -o gordb
```

### 2. 示例代码  
以下是一个简单的使用示例：  

```go
package main

import (
	"memorydb"
	"fmt"
)

func main() {
	// 初始化数据库
	db := memorydb.New()

	// 创建表
	db.CreateTable("users", []string{"id", "name", "age"})

	// 插入数据
	db.Insert("users", map[string]interface{}{"id": 1, "name": "Alice", "age": 25})
	db.Insert("users", map[string]interface{}{"id": 2, "name": "Bob", "age": 30})

	// 查询数据
	results := db.Query("users", func(row map[string]interface{}) bool {
		return row["age"].(int) > 25
	})
	fmt.Println("查询结果:", results)

	// 更新数据
	db.Update("users", func(row map[string]interface{}) bool {
		return row["name"].(string) == "Alice"
	}, map[string]interface{}{"age": 26})

	// 删除数据
	db.Delete("users", func(row map[string]interface{}) bool {
		return row["id"].(int) == 2
	})

	// 显示表内容
	fmt.Println("表内容:", db.Dump("users"))
}
```

### 3. API 文档  
| 方法           | 功能说明                                     | 示例                                         |  
|----------------|--------------------------------------------|---------------------------------------------|  
| `CreateTable`  | 创建一个新表                                 | `db.CreateTable("users", []string{"id"})`   |  
| `Insert`       | 插入一行数据                                 | `db.Insert("users", {"id": 1, "name": "A"})`|  
| `Query`        | 查询数据，支持条件过滤                       | `db.Query("users", func(row){ ... })`       |  
| `Update`       | 更新数据，支持条件过滤                       | `db.Update("users", func(row){ ... }, {...})` |  
| `Delete`       | 删除数据，支持条件过滤                       | `db.Delete("users", func(row){ ... })`      |  
| `Dump`         | 导出表中所有数据                             | `db.Dump("users")`                          |  

---

## 性能优化建议  
- **索引**：针对频繁查询的字段建立索引以提高查询效率。  
- **事务分区**：将不同的事务隔离到不同的线程中执行，减少锁冲突。  
- **快照机制**：对持久化场景，利用定期快照减少内存占用。  

---

## 适用场景  
- **缓存服务**：作为高性能的临时数据存储。  
- **实时分析**：支持快速统计和筛选操作。  
- **嵌入式数据库**：适用于轻量级、单机应用。  

---

## 开发计划  
- **支持更多查询语法**：例如范围查询、模糊匹配等。  
- **支持分布式**：提供多节点数据同步和分片存储功能。  
- **更复杂的事务支持**：实现 MVCC 或类似的隔离级别。  
- **GUI 工具**：提供可视化操作界面，简化开发和调试。  

---

## 贡献  
欢迎任何形式的贡献，包括但不限于：  
- 报告问题  
- 提交代码  
- 提出新特性需求  

如有问题，请通过 [Issues](https://github.com/cloakscn/go-redis/issues) 联系我们。

---

## 许可证  
本项目基于 MIT 协议开源，详细信息请参考 [LICENSE](LICENSE)。

---

如果有更多功能需要加入，或实际实现中有细节特点需要体现，可以进一步调整 README。