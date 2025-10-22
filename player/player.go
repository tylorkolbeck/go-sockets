package player

type Vec3 struct{ X, Y, Z float64 }

type Player struct {
	Id  string
	pos Vec3
}
