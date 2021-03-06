package nymo

type serverReserver struct {
	u      *User
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
	if u.peerSameCohort(s.cohort) {
		u.numIn--
	}
	s.u = nil
}

func (s *serverReserver) commit(p *peer) {
	u := s.u
	u.peerLock.Lock()
	defer u.peerLock.Unlock()

	if !sameCohort(p.cohort, s.cohort) && u.peerSameCohort(s.cohort) {
		// bad path, the peer's cohort changed, and
		// originally we thought it was in-cohort
		u.numIn--
	}

	u.peers[*s.id] = p
	s.u = nil
}

type clientReserver struct {
	u      *User
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
		if u.peerSameCohort(*c.cohort) {
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

func (u *User) shouldConnectPeers() bool {
	u.peerLock.RLock()
	defer u.peerLock.RUnlock()

	return u.total < u.cfg.MaxInCohortConn+u.cfg.MaxOutCohortConn
}

func (u *User) peerCleanup() {
	u.peerLock.Lock()
	defer u.peerLock.Unlock()

	for k, p := range u.peers {
		if p == nil {
			continue
		}
		if p.ctx.Err() != nil {
			delete(u.peers, k)
			u.total--
			if u.peerSameCohort(p.cohort) {
				u.numIn--
			}
		}
	}
}

func (u *User) reserveCohort(cohort uint32) bool {
	u.peerLock.Lock()
	defer u.peerLock.Unlock()

	if u.peerSameCohort(cohort) {
		if u.numIn >= u.cfg.MaxInCohortConn {
			return false
		}
		u.numIn++
	} else {
		if u.total-u.numIn >= u.cfg.MaxOutCohortConn {
			return false
		}
	}
	u.total++
	return true
}

func (u *User) reserveId(id [hashTruncate]byte) bool {
	u.peerLock.Lock()
	defer u.peerLock.Unlock()

	if _, ok := u.peers[id]; ok {
		return false
	}

	u.peers[id] = nil
	return true
}

func (u *User) reserveServer(cohort uint32) *serverReserver {
	if !u.reserveCohort(cohort) {
		return nil
	}
	return &serverReserver{
		u:      u,
		cohort: cohort,
	}
}

func (u *User) reserveClient(id [hashTruncate]byte) *clientReserver {
	if !u.reserveId(id) {
		return nil
	}
	return &clientReserver{
		u:  u,
		id: id,
	}
}
