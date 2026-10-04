package signatures

import "testing"

func TestNoSignaturesNoHits(t *testing.T) {
	if got := New().Scan("p", "s1", []string{"aim"}); len(got) != 0 {
		t.Fatalf("没有特征时不该命中，得到 %v", got)
	}
}

func TestExactFeatureHits(t *testing.T) {
	scanner := New()
	scanner.AddSignature("aim")
	if got := scanner.Scan("p", "s1", []string{"aim"}); len(got) != 1 || got[0] != "aim" {
		t.Fatalf("登记过的特征应该命中，得到 %v", got)
	}
}

func TestHitsStartEmpty(t *testing.T) {
	if got := New().Hits("p"); len(got) != 0 {
		t.Fatalf("还没扫描过时命中台账应该是空的，得到 %v", got)
	}
}

func TestScannedStartsAtZero(t *testing.T) {
	if New().Scanned() != 0 {
		t.Fatal("Scanned 初始应该是 0")
	}
}
