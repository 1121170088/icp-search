package upstream

import "icp-search/entity"

type Upstream interface {
	Search(domain string)  (*entity.Icp, error)
}