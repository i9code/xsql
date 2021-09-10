package xsql

import (
	`errors`
	`fmt`
	`strings`
	`time`
)

type DatabaseConfig struct {
	// 数据库类型
	Type string `default:"psql" json:"type" yaml:"type" validate:"required,oneof=mysql psql"`
	// 地址，填写服务器地址
	Host string `default:"127.0.0.1" yaml:"host" validate:"required"`
	// 端口
	Port int `default:"5432" yaml:"port" validate:"required"`
	// 授权，用户名
	Username string `json:"username,omitempty" yaml:"username"`
	// 授权，密码
	Password string `json:"password,omitempty" yaml:"password"`
	// 连接协议
	Protocol string `default:"tcp" json:"protocol" yaml:"protocol" validate:"required,oneof=tcp udp"`
	// 连接池配置
	Connection ConnectionConfig `json:"connection" yaml:"connection"`

	// 表名的前缀
	Suffix string `json:"suffix,omitempty" yaml:"suffix"`
	// 表名后缀
	Prefix string `json:"prefix,omitempty" yaml:"prefix"`
	// 连接的数据库名
	Schema string `json:"schema" yaml:"schema" validate:"required"`

	// 额外参数
	Parameters string `json:"parameters,omitempty" yaml:"parameters"`
	// 是否连接时使用Ping测试数据库连接是否完好
	Ping bool `default:"true" json:"ping" yaml:"ping"`
	// 是否显示SQL执行语句
	Show bool `default:"false" json:"show" yaml:"show"`
}

func (c *DatabaseConfig) dsn() (dsn string, err error) {
	switch strings.ToLower(c.Type) {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@%s(%s:%s)", c.Username, c.Password, c.Protocol, c.Host, c.Port)
		if "" != strings.TrimSpace(c.Schema) {
			dsn = fmt.Sprintf("%s/%s", dsn, strings.TrimSpace(c.Schema))
		}
	case "psql":
		if len(c.Password) == 0 {
			dsn = fmt.Sprintf(
				`host=%s port=%d dbname=%s user=%s sslmode=disable`,
				c.Host, c.Port, c.Schema, c.Username,
			)
		} else {
			dsn = fmt.Sprintf(
				`host=%s port=%d dbname=%s user=%s password=%s sslmode=disable`,

				c.Host, c.Port, c.Schema, c.Username, c.Password,
			)
		}
	default:
		err = errors.New("不支持的数据库类型")
	}
	if nil != err {
		return
	}

	// 增加参数
	if "" != strings.TrimSpace(c.Parameters) {
		dsn = fmt.Sprintf("%s?%s", dsn, strings.TrimSpace(c.Parameters))
	}

	return
}

// ConnectionConfig 连接池配置
type ConnectionConfig struct {
	// 最大打开连接数
	MaxOpen int `default:"150" yaml:"maxOpen" json:"maxOpen"`
	// 最大休眠连接数
	MaxIdle int `default:"30" yaml:"maxIdle" json:"maxIdle"`
	// 每个连接最大存活时间
	MaxLifetime time.Duration `default:"5s" yaml:"maxLifetime" json:"maxLifetime"`
}
