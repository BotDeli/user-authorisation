package session

type Display interface {
	GetSession(login string) string
	IsSession(login string) bool
}

//type Postgres struct {
//	db *sql.DB
//}
//
//func initDisplay() Display {
//	return &Postgres{}
//}
//
//func (p *Postgres) GetSession(login string) string {
//	return ""
//}
//
//func (p *Postgres) IsSession(login string) bool {
//	return false
//}
