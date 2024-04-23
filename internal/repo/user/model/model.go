package user

type User struct {
	UUID string `bson:"_id"`
	Name string `bson:"name"`
}
