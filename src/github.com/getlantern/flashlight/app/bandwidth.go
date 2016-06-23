package app

import (
	"github.com/getlantern/bandwidth"
	"github.com/getlantern/i18n"
	"github.com/getlantern/notifier"

	"github.com/getlantern/flashlight/ui"
)

func serveBandwidth(uiaddr string) error {
	helloFn := func(write func(interface{}) error) error {
		log.Debugf("Sending current bandwidth quota to new client")
		return write(bandwidth.GetQuota())
	}

	service, err := ui.Register("bandwidth", helloFn)
	if err == nil {
		go func() {
			n := notify.NewNotifications()
			var notified bool
			for quota := range bandwidth.Updates {
				service.Out <- quota
				if quota.MiBAllowed <= quota.MiBUsed {
					if !notified {
						go notifyUser(n, uiaddr)
						notified = true
					}
				}
			}
		}()
	}

	return err
}

func notifyUser(n notify.Notifier, uiaddr string) {
	// TODO: We need to translate these strings somehow.
	logo := "http://" + uiaddr + "/img/lantern_logo.png"

	msg := &notify.Notification{
		Title:    i18n.T("BACKEND_DATA_TITLE"),
		Message:  i18n.T("BACKEND_DATA_MESSAGE"),
		ClickURL: uiaddr,
		IconURL:  logo,
	}
	err := n.Notify(msg)
	if err != nil {
		log.Errorf("Could not notify? %v", err)
	}
}
