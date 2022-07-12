package dict

type Dictionary map[string]string

func (d Dictionary) Search(word string) string {
	return d[word]
}

func (d Dictionary) Add(word, definition string) {
	d[word] = definition
}
