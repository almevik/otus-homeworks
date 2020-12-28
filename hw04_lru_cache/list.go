package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) bool
	PushBack(v interface{}) bool
	Remove(i *ListItem) // удалить элемент
	MoveToFront(i *ListItem) // переместить элемент в начало
}

type ListItem struct {
	Value interface{} // значение
	Next  *ListItem   // следующий элемент
	Prev  *ListItem   // предыдущий элемент
}

type list struct {
	first *ListItem
	last  *ListItem
	len   int
}

// Если список пуст добавляем первое значение и ставим указатели первого и последнего элементов на него же
func (l *list) addFirst(lI *ListItem) {
	l.first = lI
	l.last = lI
	l.len++
}

func (l *list) Front() *ListItem {
	return l.first
}

func (l *list) Back() *ListItem {
	return l.last
}

func (l *list) Len() int {
	return l.len
}

func (l *list) PushFront(v interface{}) bool {
	var lI ListItem
	lI.Value = v

	if l.Len() == 0 {
		l.addFirst(&lI)
		return true
	}

	lI.Next, l.first, l.first.Prev = l.first, &lI, &lI
	l.len++
	return true
}

func (l *list) PushBack(v interface{}) bool {
	var lI ListItem
	lI.Value = v

	if l.Len() == 0 {
		l.addFirst(&lI)
		return true
	}

	lI.Prev, l.last, l.last.Next = l.last, &lI, &lI
	l.len++
	return true
}

func (l *list) Remove(i *ListItem) {
	// Если элемент один, то просто обнуляем список
	if l.len == 1 {
		l = nil
		return
	}

	switch i {
	case l.first:
		l.first, l.first.Next.Prev = l.first.Next, nil
	case l.last:
		l.last, l.last.Prev.Next = l.last.Prev, nil
	default:
		// Нам известно куда указывает элемент
		i.Next.Prev = i.Prev
		i.Prev.Next = i.Next
	}

	i = nil

	l.len--
	return
}

func (l *list) MoveToFront(i *ListItem) {
	switch i {
	case l.first:
		// Ничего не делаем
		break
	case l.last:
		// Делаем предпоследний элемент последним
		l.last = i.Prev
		l.last.Next = nil
		// Меняем указатель первого элемента на текущий и делаем первым текущий
		i.Next, l.first, l.first.Prev = l.first, i, i
		i.Prev = nil
	default:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
		// Меняем указатель первого элемента на текущий и делаем первым текущий
		i.Next, l.first, l.first.Prev  = l.first, i, i
		i.Prev = nil
	}
}

func NewList() List {
	return &list{}
}
