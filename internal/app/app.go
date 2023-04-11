package app

func Run() {
	s := NewServer()
	s.ListenAndServe()
}
