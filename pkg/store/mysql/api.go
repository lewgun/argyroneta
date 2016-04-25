package mysql

import (
	"github.com/lewgun/argyroneta/pkg/errutil"
	"github.com/lewgun/argyroneta/pkg/types"
)

//Rules get all fetch rules
func (s *store) Rules() ([]types.Rule, error) {

	rules := make([]types.Rule, 0)
	err := s.Engine.Find(&rules)
	return rules, err

}

func (s *store) AddEntry(e *types.Entry) error {

	_, err := s.Engine.Insert(e)
	return err

}

func (s *store) EntryTopN(req *types.TopReq) ([]types.Entry, error) {

	if req == nil {
		return nil, errutil.ErrInvalidParameter
	}

	entries := make([]types.Entry, 0)
	err := s.Engine.Limit(req.N, 0).Find(&entries)
	return entries, err

}
