package bird

import (
	"fmt"

	"github.com/satori/go.uuid"
)

type Egg struct {
	Id         string            `json:"id" valid:"-"`
	Headers    map[string]string `json:"headers,omitempty" valid:"-"`
	Body       string            `json:"body" valid:"-"`
	StatusCode int               `json:"status_code" valid:"required"`
	TTL        int               `json:"ttl,omitempty" valid:"-"`
}

func (self *Egg) initialize() {
	if self.Id == "" {
		e, err := uuid.NewV4()
		if err != nil {
			panic(err)
		}
		self.Id = fmt.Sprintf("%s", e)
	}
}
