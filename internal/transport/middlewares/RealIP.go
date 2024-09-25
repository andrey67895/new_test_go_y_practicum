package middlewares

import (
	"fmt"
	"net"
	"net/http"

	"github.com/andrey67895/new_test_go_y_practicum/internal/config"
)

func RealIP(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if rip := config.TrustedSubnet; rip != "" {
			ip := net.ParseIP(r.Header.Get("X-Real-IP"))
			ones, _ := ip.DefaultMask().Size()
			_, i, _ := net.ParseCIDR(fmt.Sprintf("%s/%d", ip.To4(), ones))
			mask := i.String()
			if mask != rip {
				http.Error(w, "Ошибка доступа ip", http.StatusForbidden)
				return
			}
			r.RemoteAddr = r.Header.Get("X-Real-IP")
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
