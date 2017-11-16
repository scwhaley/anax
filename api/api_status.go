package api

import (
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/open-horizon/anax/policy"
)

func (a *API) status(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":

		info := NewInfo(a.Config)

		if err := WriteConnectionStatus(info); err != nil {
			glog.Errorf(apiLogString(fmt.Sprintf("Unable to get connectivity status: %v", err)))
		}

		a.bcStateLock.Lock()
		defer a.bcStateLock.Unlock()

		for _, bc := range a.bcState[policy.Ethereum_bc] {
			geth := NewGeth()

			gethURL := fmt.Sprintf("http://%v:%v", bc.service, bc.servicePort)
			if err := WriteGethStatus(gethURL, geth); err != nil {
				glog.Errorf(apiLogString(fmt.Sprintf("Unable to determine geth service facts: %v", err)))
			}

			info.AddGeth(geth)
		}

		writeResponse(w, info, http.StatusOK)
	case "OPTIONS":
		w.Header().Set("Allow", "GET, OPTIONS")
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
