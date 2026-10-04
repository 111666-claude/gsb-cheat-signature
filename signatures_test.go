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

func TestLiveSignatureTable(t *testing.T) {
	scanner := New()
	scanner.Scan("p", "s1", []string{"aim"})
	scanner.AddSignature("aim")
	got := scanner.Scan("p", "s2", []string{"aim"})
	if len(got) != 1 || got[0] != "aim" {
		t.Fatalf("首次扫描后登记的特征应该立刻生效，得到 %v", got)
	}
}

func TestDuplicateSampleIdempotent(t *testing.T) {
	scanner := New()
	scanner.AddSignature("aim")
	first := scanner.Scan("p", "s1", []string{"aim"})
	second := scanner.Scan("p", "s1", []string{"aim"})
	if len(second) != 1 || second[0] != "aim" {
		t.Fatalf("重复扫描应该返回同样的命中，首次 %v 重复 %v", first, second)
	}
	if got := scanner.Hits("p"); len(got) != 1 || got[0] != "aim" {
		t.Fatalf("同一个 sampleID 只记一次账，得到 %v", got)
	}
}

func TestExactMatchOnly(t *testing.T) {
	scanner := New()
	scanner.AddSignature("aimbot")
	if got := scanner.Scan("p", "s1", []string{"aim"}); len(got) != 0 {
		t.Fatalf("部分包含不该命中，得到 %v", got)
	}
	scanner.AddSignature("aim")
	if got := scanner.Scan("p", "s2", []string{"aimbotx"}); len(got) != 0 {
		t.Fatalf("上报特征包含登记特征也不该命中，得到 %v", got)
	}
}

func TestHitsLedgerSortedAndDeduped(t *testing.T) {
	scanner := New()
	scanner.AddSignature("wallhack")
	scanner.AddSignature("aim")
	scanner.Scan("p", "s1", []string{"wallhack", "aim"})
	scanner.Scan("p", "s2", []string{"aim"})
	got := scanner.Hits("p")
	if len(got) != 2 || got[0] != "aim" || got[1] != "wallhack" {
		t.Fatalf("命中台账应该是升序去重的不同特征 id，得到 %v", got)
	}
}

func TestScanDoesNotSweepSignatureTable(t *testing.T) {
	scanner := New()
	for index := 0; index < 3000; index++ {
		scanner.AddSignature(string(rune('a' + index%26)) + string(rune(index)))
	}
	for index := 0; index < 3000; index++ {
		scanner.Scan("p", string(rune(index)), []string{"missing"})
	}
	if scanned := scanner.Scanned(); scanned > 6000 {
		t.Fatalf("Scanned 不该随特征数乘扫描次数增长，得到 %d", scanned)
	}
}
