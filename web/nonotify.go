package web

import "todo-app/dto"

type NoNotify struct{}

func (d NoNotify) Add(s *dto.SavedToDo)    {}
func (d NoNotify) Change(s *dto.SavedToDo) {}
func (d NoNotify) Delete(s *dto.SavedToDo) {}
