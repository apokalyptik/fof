package main

import "net/url"

type slackMsg struct {
	Where string
}

func (s *slackMsg) to(where string) *slackMsg {
	s.Where = where
	return s
}

func (s *slackMsg) send(text string) error {
	if s.Where[0] == '@' {
		if channel, err := udb.getChannelForIM(s.Where[1:]); err != nil {
			return err
		} else {
			s.Where = channel
		}
	} else {
		s.Where = "#slack-tools-testing"
	}

	slackMsgQueue <- url.Values{
		"channel": []string{s.Where},
		"text":    []string{text},
	}
	return nil
}
