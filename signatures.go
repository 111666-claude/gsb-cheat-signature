// Package signatures 是作弊特征扫描：用已知特征集合匹配玩家上报的特征，命中记账。
package signatures

import (
	"sort"
)

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

// AddSignature 登记一条作弊特征，返回是否成功。
func (s *Scanner) AddSignature(id string) bool {
	if s.sigs[id] {
		return false
	}
	s.sigs[id] = true
	return true
}

// Scan 扫描一次样本，返回本次命中的特征 id（升序）。
// 精确匹配当前特征表；同一个 sampleID 只记一次账，重复扫描直接返回首次结果。
func (s *Scanner) Scan(player, sampleID string, features []string) []string {
	if out, ok := s.seen[sampleID]; ok {
		return append([]string{}, out...)
	}
	s.scanned += len(features)
	out := []string{}
	for _, feature := range features {
		if s.sigs[feature] {
			out = append(out, feature)
		}
	}
	sort.Strings(out)
	if s.hits[player] == nil {
		s.hits[player] = map[string]bool{}
	}
	for _, id := range out {
		s.hits[player][id] = true
	}
	s.seen[sampleID] = out
	return out
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

// Scanned 是累计扫描的特征数（规模观测）。
func (s *Scanner) Scanned() int { return s.scanned }
