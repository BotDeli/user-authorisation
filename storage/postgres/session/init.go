package session

//go:generate go run github.com/vektra/mockery/v2@v2.32.0 --name=Display
type Display interface {
	NewSession(login string) (string, error)
	GetLoginFromSession(session string) (string, error)
	UpdateSessionLifeTime(login string)
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
