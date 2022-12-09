package models

// Value 值类型
type Value struct {
	hasData bool        // 是否有数据
	value   interface{} // 值
}

// NewValue 构建一个值类型
func NewValue(value interface{}) *Value {
	return &Value{hasData: true, value: value}
}

// Val 获取实际值
func (value *Value) Val() interface{} {
	if value == nil {
		return nil
	}
	if value.hasData {
		return value.value
	}
	return nil
}

// String 数据转字符串
func (value *Value) String() string {
	data := value.Val()
	if data != nil {
		return data.(string)
	}
	return ""
}
