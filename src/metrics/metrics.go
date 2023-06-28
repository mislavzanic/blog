package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ViewIndex = promauto.NewCounter(prometheus.CounterOpts{
		Name: "view_index",
		Help: "Views on the main page",
	})

	BlogHits = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "blog_hits",
		Help: "Views on a blog site",
	}, []string{"name"})
)
