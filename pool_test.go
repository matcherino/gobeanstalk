package gobeanstalk

import "testing"

func newPool(t *testing.T) *Pool {
	p, err := NewPool(address, 10)
	if err != nil {
		t.Fatal("NewPool failed.err = :", err.Error())
	}
	return p
}

func assertPoolLength(t *testing.T, fName string, poolLength int, expectedLength int) {
	if poolLength != expectedLength {
		t.Fatalf("%s failed.err = pool size of %d expected got %d", fName, poolLength, expectedLength)
	}
}

func TestNewPool(t *testing.T) {
	p, err := NewPool(address, 10)
	if err != nil {
		t.Fatal("NewPool failed.err = :", err.Error())
	}
	assertPoolLength(t, "NewPool", len(p.pool), 10)

	// Test getting and using a pool member
	conn, err := p.Get(testtube)
	if err != nil {
		t.Fatal("Pool.Get failed.Err = ", err.Error())
	}
	assertPoolLength(t, "PoolGet", len(p.pool), 9)
	err = conn.Use(testtube)
	if err != nil {
		t.Fatal("Get (use) failed.Err = ", err.Error())
	}

	//Test releasing the connection back into the pool
	p.Release(conn)
	assertPoolLength(t, "Pool.Release", len(p.pool), 10)

	//  Test emptying the pool
	p.Empty()
	assertPoolLength(t, "Pool.Release", len(p.pool), 0)
}
