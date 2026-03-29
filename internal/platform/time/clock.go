package time

import stdtime "time"

type Clock struct{}

func (Clock) Now() stdtime.Time { return stdtime.Now() }
