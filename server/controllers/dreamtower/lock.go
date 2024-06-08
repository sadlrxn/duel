package dreamtower

func (c *Controller) lockUser(userID uint) {
	c.lockedUsers.Store(userID, true)
}

func (c *Controller) releaseUser(userID uint) {
	if _, prs := c.lockedUsers.Load(userID); prs {
		c.lockedUsers.Delete(userID)
	}
}

func (c *Controller) checkUserLocked(userID uint) (prs bool) {
	_, prs = c.lockedUsers.Load(userID)
	return
}
