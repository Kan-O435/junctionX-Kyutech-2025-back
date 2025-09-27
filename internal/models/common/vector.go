package common

import "math"

type Vector3D struct {
    X float64 `json:"x"`
    Y float64 `json:"y"`
    Z float64 `json:"z"`
}

func (v Vector3D) Add(other Vector3D) Vector3D {
    return Vector3D{
        X: v.X + other.X,
        Y: v.Y + other.Y,
        Z: v.Z + other.Z,
    }
}

func (v Vector3D) Subtract(other Vector3D) Vector3D {
    return Vector3D{
        X: v.X - other.X,
        Y: v.Y - other.Y,
        Z: v.Z - other.Z,
    }
}

func (v Vector3D) Scale(scalar float64) Vector3D {
    return Vector3D{
        X: v.X * scalar,
        Y: v.Y * scalar,
        Z: v.Z * scalar,
    }
}

func (v Vector3D) Magnitude() float64 {
    return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vector3D) Normalize() Vector3D {
    mag := v.Magnitude()
    if mag == 0 {
        return Vector3D{}
    }
    return v.Scale(1.0 / mag)
}

func (v Vector3D) Dot(other Vector3D) float64 {
    return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

func (v Vector3D) Cross(other Vector3D) Vector3D {
    return Vector3D{
        X: v.Y*other.Z - v.Z*other.Y,
        Y: v.Z*other.X - v.X*other.Z,
        Z: v.X*other.Y - v.Y*other.X,
    }
}