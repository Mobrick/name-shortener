package handler

import (
	"encoding/json"
	"net"
	"net/http"

	"github.com/Mobrick/name-shortener/internal/logger"
)

// StatsHandle показывает статистику по серверу.
func (env Env) StatsHandle(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	trustedSubnet := env.ConfigStruct.FlagTrustedSubnet
	if len(trustedSubnet) == 0 {
		logger.Log.Debug("subnet from config is empty")
		res.WriteHeader(http.StatusForbidden)
		return
	}

	subnetIP := net.ParseIP(trustedSubnet)
	if len(subnetIP) == 0 {
		logger.Log.Debug("could not parse subnet ip from config")
		res.WriteHeader(http.StatusForbidden)
		return
	}

	ipFromReq := req.RemoteAddr
	ip := net.ParseIP(ipFromReq)
	if len(ip) == 0 {
		logger.Log.Debug("could not parse ip from request")
		res.WriteHeader(http.StatusForbidden)
		return
	}

	subnet := net.IPNet{
		IP:   subnetIP,
		Mask: net.CIDRMask(24, 32),
	}

	if !subnet.Contains(ip) {
		logger.Log.Debug("ip is not in the subnet")
		res.WriteHeader(http.StatusForbidden)
	}

	// TODO: тут заменить на получение количества сокращенных URL
	stats, err := env.Storage.GetStats(ctx)
	if err != nil {
		logger.Log.Debug("could not get stats")
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(stats)
	if err != nil {
		logger.Log.Debug("could not marshal response")
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(resp))
}
