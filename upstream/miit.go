package upstream

import (
	"github.com/fghwett/icp/abbreviateinfo"
	"icp-search/entity"
	"log"
)
type MiitUpstream struct{}
var Miit = &MiitUpstream{}
func (u *MiitUpstream) Search(domain string) (*entity.Icp, error) {

	icp2 := &abbreviateinfo.Icp{}

	domainInfo, err := icp2.Query(domain)
	if err == abbreviateinfo.IcpNotForRecord {
		return nil, Norecord
	} else if err != nil {
		return nil, SearchErr
	} else {
		icp := &entity.Icp{
			Domain:   domainInfo.Domain,
			Unit:     domainInfo.UnitName,
			Type:     domainInfo.NatureName,
			IcpCode:  domainInfo.ServiceLicence,
			Name:     domainInfo.ServiceName,
			PassTime: domainInfo.UpdateRecordTime,
		}
		log.Printf("%v", domainInfo)
		return icp, nil
	}
}