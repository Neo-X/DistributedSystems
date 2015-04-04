package s3dm

import "math"

type Sphere struct {
	Xform
	Radius float64
}

func NewSphere(pos Position, radius float64) *Sphere {
	s := new(Sphere)
	s.Xform = XformIdentity
	s.Position = pos
	s.Radius = radius
	return s
}

// Returns the normal vector for a point 'p' on sphere 's'
func (s *Sphere) Normal(p Position) V3 {
	delta := p.Sub(s.Position).V3()
	return delta.Unit()
}

/*
	Returns the point of intersection and the normal at that point if 
	there is an intersection
	Returns two empty vectors otherwise????
	
	I don't think this works perfectly
*/
func (s *Sphere) Intersect(r *Ray) (bool, Position, V3) {
	pos := s.Position
	ro, rd := r.Origin, r.Dir
	rp := ro.Sub(pos).V3()
	A := rd.Dot(rd)
	B := float64(2) * (rd.X*rp.X +
		rd.Y*rp.Y +
		rd.Z*rp.Z)
	C := (rp.X*rp.X +
		rp.Y*rp.Y +
		rp.Z*rp.Z) -
		s.Radius*s.Radius

	delta := B*B - 4*A*C
	if delta > 0 {
		t0 := (-B - math.Sqrt(delta)) / 2
		t1 := (-B + math.Sqrt(delta)) / 2

		t := float64(0)

		// t0 must be smaller than t1
		if t0 > t1 {
			t0, t1 = t1, t0
		}

		// Sphere behind ray
		if t1 < 0 {
			return false, Position{}, V3{}
		}

		if t0 < 0 {
			t = t1
		} else {
			t = t0
		}

		intersection := ro.Addf(rd.Muls(t))
		normal := s.Normal(intersection)
		return true, intersection, normal
	}
	return false, Position{}, V3{}
}
