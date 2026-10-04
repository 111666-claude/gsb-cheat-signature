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

func TestSignatureTakesEffectImmediately(t *testing.T) {
	scanner := New()
	scanner.Scan("p", "s1", []string{"aim"})
	if got := scanner.Scan("p", "s2", []string{"aim"}); len(got) != 0 {
		t.Fatalf("登记前不应命中，得到 %v", got)
	}
	scanner.AddSignature("aim")
	if got := scanner.Scan("p", "s3", []string{"aim"}); len(got) != 1 || got[0] != "aim" {
		t.Fatalf("登记后新扫描应立即命中，得到 %v", got)
	}
}

func TestDuplicateSampleIDIsIdempotent(t *testing.T) {
	scanner := New()
	scanner.AddSignature("aim")
	first := scanner.Scan("p", "s1", []string{"aim"})
	again := scanner.Scan("p", "s1", []string{"aim"})
	if len(again) != 1 || again[0] != "aim" {
		t.Fatalf("重复扫描应仍返回命中，得到 %v", again)
	}
	if len(first) != len(again) || first[0] != again[0] {
		t.Fatalf("重复扫描返回应与首次一致：%v vs %v", first, again)
	}
	if got := scanner.Hits("p"); len(got) != 1 || got[0] != "aim" {
		t.Fatalf("同一 sampleID 只记一次，台账应只有一条，得到 %v", got)
	}
	if scanner.Scanned() != 1 {
		t.Fatalf("重复扫描不应重复计数，得到 %d", scanner.Scanned())
	}
}

func TestPartialMatchDoesNotHit(t *testing.T) {
	scanner := New()
	scanner.AddSignature("aimbot")
	if got := scanner.Scan("p", "s1", []string{"aim"}); len(got) != 0 {
		t.Fatalf("部分包含不算命中，得到 %v", got)
	}
	if got := scanner.Scan("p", "s2", []string{"aimbot-v2", "xaimbot"}); len(got) != 0 {
		t.Fatalf("超串也不算命中，得到 %v", got)
	}
	if got := scanner.Scan("p", "s3", []string{"aimbot", "aimbot"}); len(got) != 1 || got[0] != "aimbot" {
		t.Fatalf("完全相同才算命中且样本内去重，得到 %v", got)
	}
}

func TestHitsLedgerSortedAndDeduped(t *testing.T) {
	scanner := New()
	scanner.AddSignature("wallhack")
	scanner.AddSignature("aimbot")
	scanner.AddSignature("speed")
	scanner.Scan("p", "s1", []string{"speed", "wallhack"})
	scanner.Scan("p", "s2", []string{"wallhack", "aimbot"})
	got := scanner.Hits("p")
	want := []string{"aimbot", "speed", "wallhack"}
	if len(got) != len(want) {
		t.Fatalf("台账去重后期望 %v，得到 %v", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("台账应升序去重，期望 %v，得到 %v", want, got)
		}
	}
	if got := scanner.Hits("other"); len(got) != 0 {
		t.Fatalf("其他玩家台账应为空，得到 %v", got)
	}
}

func TestScannedDoesNotScaleWithSignatureCount(t *testing.T) {
	scanner := New()
	for i := 0; i < 3000; i++ {
		scanner.AddSignature("sig-" + string(rune('a'+i%26)))
	}
	for i := 0; i < 3000; i++ {
		scanner.Scan("p", "s-"+string(rune('a'+i)), []string{"sig-a"})
	}
	if got := scanner.Scanned(); got != 3000 {
		t.Fatalf("Scanned 只应按上报特征计数，期望 3000，得到 %d", got)
	}
}

func TestRepeatedScanSequenceIsDeterministic(t *testing.T) {
	build := func() ([]string, []string) {
		scanner := New()
		scanner.AddSignature("aim")
		first := scanner.Scan("p", "s1", []string{"aim"})
		second := scanner.Scan("p", "s1", []string{"aim"})
		return first, second
	}
	a1, a2 := build()
	b1, b2 := build()
	if a1[0] != b1[0] || a2[0] != b2[0] {
		t.Fatalf("同一串扫描重复执行结果应相同：%v/%v vs %v/%v", a1, a2, b1, b2)
	}
}
