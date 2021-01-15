package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	Len() int
	Front() *listItem
	Back() *listItem
	PushFront(v interface{})
	PushBack(v interface{})
	Remove(i *listItem)      // удалить элемент
	MoveToFront(i *listItem) // переместить элемент в начало
}

type listItem struct {
	Value interface{} // значение
	Next  *listItem   // следующий элемент
	Prev  *listItem   // предыдущий элемент
}

type list struct {
	first *listItem
	last  *listItem
	len   int
}

// Если список пуст добавляем первое значение и ставим указатели первого и последнего элементов на него же.
func (l *list) addFirst(lI *listItem) {
	l.first = lI
	l.last = lI
	l.len++
}

func (l *list) Front() *listItem {
	return l.first
}

func (l *list) Back() *listItem {
	return l.last
}

func (l *list) Len() int {
	return l.len
}

func (l *list) PushFront(v interface{}) {
	lI := &listItem{v, nil, nil}

	if l.Len() == 0 {
		l.addFirst(lI)
		return
	}

	lI.Next = l.first
	l.first, l.first.Prev = lI, lI
	l.len++
}

func (l *list) PushBack(v interface{}) {
	lI := &listItem{v, nil, nil}

	if l.Len() == 0 {
		l.addFirst(lI)
		return
	}

	lI.Prev = l.last
	l.last, l.last.Next = lI, lI
	l.len++
}

func (l *list) Remove(i *listItem) {
	// Если элемент один или список пуст, то просто обнуляем список
	if l.len <= 1 {
		l.first = nil
		l.last = nil
		l.len = 0
		return
	}

	switch i {
	case l.first:
		l.first.Next.Prev = nil
		l.first = l.first.Next
	case l.last:
		l.last.Prev.Next = nil
		l.last = l.last.Prev
	default:
		// Нам известно куда указывает каждый элемент
		i.Next.Prev = i.Prev
		i.Prev.Next = i.Next
	}

	i = nil
	l.len--
}

func (l *list) MoveToFront(i *listItem) {
	switch i {
	case l.first:
		return
	case l.last:
		// Делаем предпоследний элемент последним
		l.last = i.Prev
		l.last.Next = nil
	default:
		// Меняем указатели для соседних элементов
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}

	// Меняем указатель первого элемента на текущий и делаем первым текущий
	i.Next = l.first
	l.first, l.first.Prev = i, i
	i.Prev = nil
}

func NewList() List {
	return &list{}
}
