package models

// MethodInfo представляет информацию о методе
type MethodInfo struct {
	Name      string   // Имя метода
	Signature string   // Полная сигнатура метода
	Body      string   // Тело метода
	Params    []string // Параметры метода
	Returns   []string // Возвращаемые значения
}
