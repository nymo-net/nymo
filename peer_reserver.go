package nymo

type serverReserver struct {
	u      *user
	cohort uint32
	id     *[hashTruncate]byte
}

func (s *serverReserver) reserveId(id *[hashTruncate]byte) bool {
	if s.id != nil {
		panic("multiple reservation")
	}

	if !s.u.reserveId(*id) {
		return false
	}

	s.id = id
	return true
}

func (s *serverReserver) rollback() {
	u := s.u
	if u == nil {
		return
	}

	u.peerLock.Lock()
	defer u.peerLock.Unlock()

	if s.id != nil {
		delete(u.peers, *s.id)
	}

	u.total--
	if sameCohort(u.cohort, s.cohort) {
		u.numIn--
	}
	s.u = nil
}

func (s *serverReserver) commit(p *peer) {
	u := s.u
	u.peerLock.Lock()
	defer u.peerLock.Unlock()

	if !sameCohort(p.cohort, s.cohort) && sameCohort(u.cohort, s.cohort) {
		// bad path, the peer's cohort changed, and
		// originally we thought it was in-cohort
		u.numIn--
	}

	u.peers[*s.id] = p
	s.u = nil
}

type clientReserver struct {
	u      *user
	cohort *uint32
	id     [hashTruncate]byte
}

func (c *clientReserver) reserveCohort(cohort uint32) bool {
	if c.cohort != nil {
		panic("multiple reservation")
	}

	if !c.u.reserveCohort(cohort) {
		return false
	}
	c.cohort = &cohort
	return true
}

func (c *clientReserver) rollback() {
	u := c.u
	if u == nil {
		return
	}

	u.peerLock.Lock()
	defer u.peerLock.Unlock()

	if c.cohort != nil {
		u.total--
		if sameCohort(u.cohort, *c.cohort) {
			u.numIn--
		}
	}

	delete(u.peers, c.id)
	c.u = nil
}

func (c *clientReserver) commit(p *peer) {
	u := c.u

	u.peerLock.Lock()
	defer u.peerLock.Unlock()

	u.peers[c.id] = p
	c.u = nil
}

func (u *user) shouldConnectPeers() bool {
	u.peerLock.RLock()
	defer u.peerLock.RUnlock()

	return u.total < u.cfg.MaxConcurrentConn
}

func (u *user) peerCleanup() {
	u.peerLock.Lock()
	defer u.peerLock.Unlock()

	for k, p := range u.peers {
		if p == nil {
			continue
		}
		if p.ctx.Err() != nil {
			delete(u.peers, k)
			u.total--
			if sameCohort(u.cohort, p.cohort) {
				u.numIn--
			}
		}
	}
}

func (u *user) reserveCohort(cohort uint32) bool {
	u.peerLock.Lock()
	defer u.peerLock.Unlock()

	maxIn := uint(float64(u.cfg.MaxConcurrentConn) * (1 - epsilon))
	if sameCohort(u.cohort, cohort) {
		if u.numIn >= maxIn {
			return false
		}
		u.numIn++
	} else {
		if u.total-u.numIn >= u.cfg.MaxConcurrentConn-maxIn {
			return false
		}
	}
	u.total++
	return true
}

func (u *user) reserveId(id [hashTruncate]byte) bool {
	u.peerLock.Lock()
	defer u.peerLock.Unlock()

	if _, ok := u.peers[id]; ok {
		return false
	}

	u.peers[id] = nil
	return true
}

func (u *user) reserveServer(cohort uint32) *serverReserver {
	if !u.reserveCohort(cohort) {
		return nil
	}
	return &serverReserver{
		u:      u,
		cohort: cohort,
	}
}

func (u *user) reserveClient(id [hashTruncate]byte) *clientReserver {
	if !u.reserveId(id) {
		return nil
	}
	return &clientReserver{
		u:  u,
		id: id,
	}
}
