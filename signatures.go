// Package signatures 是作弊特征扫描：用已知特征集合匹配玩家上报的特征，命中记账。
// 精确匹配、特征表实时生效、同一样本幂等、命中台账去重升序、扫描按特征逐个精确查找。
package signatures

import "sort"

// Scanner 是特征扫描台账。
type Scanner struct {
	sigs    map[string]bool
	seen    map[string][]string
	hits    map[string]map[string]bool
	scanned int
}

// New 建空台账。
func New() *Scanner {
	return &Scanner{
		sigs: map[string]bool{},
		seen: map[string][]string{},
		hits: map[string]map[string]bool{},
	}
}

// AddSignature 登记一条作弊特征，返回是否成功；新特征立刻生效。
func (s *Scanner) AddSignature(id string) bool {
	s.sigs[id] = true
	return true
}

// Scan 扫描一次样本，返回本次命中的特征 id（升序、去重）。
// 同一个 sampleID 只记一次，重复扫描返回相同命中但不重复记账。
func (s *Scanner) Scan(player, sampleID string, features []string) []string {
	if ids, ok := s.seen[sampleID]; ok {
		return append([]string{}, ids...)
	}
	matched := map[string]bool{}
	for _, feature := range features {
		s.scanned++
		if s.sigs[feature] {
			matched[feature] = true
		}
	}
	out := make([]string, 0, len(matched))
	for id := range matched {
		out = append(out, id)
	}
	sort.Strings(out)
	s.seen[sampleID] = out
	if s.hits[player] == nil {
		s.hits[player] = map[string]bool{}
	}
	for _, id := range out {
		s.hits[player][id] = true
	}
	return append([]string{}, out...)
}

// Hits 返回该玩家命中过的不同特征 id（升序、去重）。
func (s *Scanner) Hits(player string) []string {
	out := make([]string, 0, len(s.hits[player]))
	for id := range s.hits[player] {
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}

// Scanned 是累计扫描的特征数（规模观测，按上报特征逐个精确查找计一次）。
func (s *Scanner) Scanned() int { return s.scanned }
