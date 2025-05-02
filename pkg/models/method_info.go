package models

// MethodInfo представляет информацию о методе
type MethodInfo struct {
	Name        string   // Имя метода
	Signature   string   // Полная сигнатура метода
	Body        string   // Тело метода
	Params      []string // Параметры метода (для внутреннего использования)
	Returns     []string // Возвращаемые значения (для внутреннего использования)
	Parameters  []string // Параметры метода (для совместимости с оркестратором)
	ReturnType  []string // Типы возвращаемых значений (для совместимости с оркестратором)
	Description string   // Описание метода (может быть заполнено с помощью ЛЛМ)
}
