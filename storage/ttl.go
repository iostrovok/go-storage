package storage

import (
	"log"
	"time"
)

func SetShieldTTL(ttl time.Duration) {
	Singleton.SetShieldTTL(ttl)
}

func (s *Storage) SetShieldTTL(ttl time.Duration) {

	if s.IsDebug {
		log.Printf("SetShieldTTL. ttl: %s\n", ttl)
	}

	s.ShieldTTL = ttl
	s.TimerCh <- s.ShieldTTL
}

func (s *Storage) _check_ttl() {
	ticker := time.NewTicker(s.ShieldTTL)

	for {
		select {
		case dur, ok := <-s.TimerCh:
			if !ok {
				return
			}

			ticker.Stop()
			ticker = time.NewTicker(dur)

		case <-ticker.C:
			mes := newMessage(Clean, "", "", "")
			s.OneAct(mes)
		}
	}
}
