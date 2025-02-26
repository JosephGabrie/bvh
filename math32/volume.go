package math32

type VolumeType[E Number] interface {
	MinBounds(volumes ...VolumeType[E])
	Score() float32
	Equals(VolumeType[E]) bool
	Overlaps(VolumeType[E]) bool
	Contains(VolumeType[E]) bool
	Intersects(VolumeType[E], *Coordinate[E]) float32
	GetPoint() Coordinate[E]
	GetDelta() Coordinate[E]
	String() string
	New() VolumeType[E]
	IsNil() bool
	IsSame(VolumeType[E]) bool
}
