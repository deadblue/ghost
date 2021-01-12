package context

func (i *Impl) PutPathVar(name, value string) *Impl {
	if i.pv == nil {
		i.pv = make(map[string]string)
	}
	i.pv[name] = value
	return i
}
