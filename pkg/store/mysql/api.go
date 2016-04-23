package mysql

import (
	"github.com/lewgun/argyroneta/pkg/types"
)

//Rules get all fetch rules
func (s *store) Rules() ([]types.Rule, error) {

	rules := make([]types.Rule, 0)
	err := s.Engine.Find(&rules)
	return rules, err

}
