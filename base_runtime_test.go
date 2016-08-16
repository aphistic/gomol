package gomol

import . "gopkg.in/check.v1"

/*
This is in its own file so the line numbers don't change.
These tests are testing calling locations so putting them in their own
file will limit the number of changes to that data.
*/

func (s *GomolSuite) TestIsGomolCaller(c *C) {
	res, file := isGomolCaller("/home/gomoltest/some/sub/dir/that/is/long/filename.go")
	c.Check(res, Equals, false)
	c.Check(file, Equals, "filename.go")
}
