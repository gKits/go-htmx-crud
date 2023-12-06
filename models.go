package crud

type User struct {
	ID     uint    `json:"id"`
	Name   string  `json:"name"`
	Age    int     `json:"age"`
	Height float32 `json:"height"`
}

func NewUser(id uint, name string, age int, height float32) User {
	return User{id, name, age, height}
}

type Averages struct {
	Age    int
	Height float32
}
