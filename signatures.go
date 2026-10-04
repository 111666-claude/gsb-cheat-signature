// Package signatures 是作弊特征扫描：用已知特征集合匹配玩家上报的特征，命中记账。
// 缺陷：第一次扫描后特征表被冻结（新特征不生效）、同一样本重复扫描重复记账、
// 部分包含也算命中、每次扫描全扫特征表。
package signatures

import (
	"sort"
	"strings"
)

// Scanner 是特征扫描台账。
type Scanner struct {
	sigs    []string
	frozen  map[string]bool
	hits    map[string][]string
	scanned int
}

// New 建空台账。
func New() *Scanner {
	return &Scanner{hits: map[string][]string{}}
}

// AddSignature 登记一条作弊特征，返回是否成功。
func (s *Scanner) AddSignature(id string) bool {
	s.sigs = append(s.sigs, id)
	return true
}

// Scan 扫描一次样本，返回本次命中的特征 id（升序）。
// 缺陷：用冻结的特征表、重复样本重复记账、用子串包含判定、全扫特征表。
func (s *Scanner) Scan(player, sampleID string, features []string) []string {
	s.scanned += len(s.sigs)
	if s.frozen == nil {
		s.frozen = map[string]bool{}
		for _, id := range s.sigs {
			s.frozen[id] = true
		}
	}
	out := []string{}
	for _, feature := range features {
		for id := range s.frozen {
			if strings.Contains(id, feature) || strings.Contains(feature, id) {
				out = append(out, id)
			}
		}
	}
	sort.Strings(out)
	for _, id := range out {
		s.hits[player] = append(s.hits[player], id)
	}
	return out
}

// Hits 返回该玩家命中过的特征 id（升序）。缺陷：全扫命中台账。
func (s *Scanner) Hits(player string) []string {
	s.scanned += len(s.hits[player])
	out := append([]string{}, s.hits[player]...)
	sort.Strings(out)
	return out
}

// Scanned 是累计扫描的特征数（规模观测）。
func (s *Scanner) Scanned() int { return s.scanned }
