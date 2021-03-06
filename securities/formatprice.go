package securities

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"math"
	"strings"
	//"github.com/RileyR387/sc-dtc-client/dtcproto"
)

func (s *Security) FormatPrice(p float64) string {
	dfmt := s.Definition.PriceDisplayFormat
	dfmtString := dfmt.String()
	dfmtVal := int(dfmt)
	var fmtStr string
	//log.Debugf("Format value is %v with string val: %v", dfmtVal, dfmtString)
	if dfmtString == "PRICE_DISPLAY_FORMAT_UNSET" || dfmtString == "PRICE_DISPLAY_FORMAT_DECIMAL_0" {
		fmtStr = "%.f"
	} else if strings.HasPrefix(dfmtString, "PRICE_DISPLAY_FORMAT_DECIMAL") {
		fmtStr = "%." + fmt.Sprint(dfmtVal) + "f"
	} else if strings.HasPrefix(dfmtString, "PRICE_DISPLAY_FORMAT_DENOMINATOR") {
		//log.Debugf("Format value for %v is %v with string val: %v", s.Definition.Symbol, dfmtVal, dfmtString)
		intVal, fraction := math.Modf(float64(p))
		if intVal == 0 {
			return fmt.Sprintf("'%05.2f", fraction*32)
		} else {
			return fmt.Sprintf("%v'%05.2f", intVal, fraction*32)
		}
	} else {
		log.Warnf("Unknown price display format: %v", dfmtString)
		fmtStr = "%.6f"
	}
	return fmt.Sprintf(fmtStr, p)
}

func (s *Security) String() string {
	bs := s.BidString()
	as := s.AskString()
	s.AddingDataMutex.Lock()
	defer s.AddingDataMutex.Unlock()
	return fmt.Sprintf("%15v %10v %10v %v", s.Definition.Symbol, bs, as, s.Definition)
}
func (s *Security) BidString() string {
	return s.FormatPrice(s.BidPrice())
}
func (s *Security) AskString() string {
	return s.FormatPrice(s.AskPrice())
}
func (s *Security) LastString() string {
	return s.FormatPrice(s.GetLastPrice())
}
func (s *Security) DchgString() string {
	return s.FormatPrice(s.GetLastPrice() - s.GetSettlementPrice())
}

func (s *Security) SettlementString() string {
	return s.FormatPrice(s.GetSettlementPrice())
}
