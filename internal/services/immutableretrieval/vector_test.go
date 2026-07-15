package immutableretrieval

import "testing"

func TestCollapseDocumentsAndFuseRRF(t *testing.T) {
	bm25 := []ChunkHit{{Rank: 1, ChunkID: "a-1", DocumentRevisionID: "a", Score: 9}, {Rank: 2, ChunkID: "a-2", DocumentRevisionID: "a", Score: 8}, {Rank: 3, ChunkID: "b-1", DocumentRevisionID: "b", Score: 7}}
	vector := []ChunkHit{{Rank: 1, ChunkID: "b-1", DocumentRevisionID: "b", Score: 0.9}, {Rank: 2, ChunkID: "a-2", DocumentRevisionID: "a", Score: 0.8}}
	collapsed := CollapseDocuments(bm25)
	if len(collapsed) != 2 || collapsed[0].ChunkID != "a-1" || collapsed[1].Rank != 2 {
		t.Fatalf("collapsed = %#v", collapsed)
	}
	fused := FuseRRF(map[string][]ChunkHit{"bm25": bm25, "vector": vector}, 60, 10)
	if len(fused) != 2 {
		t.Fatalf("fused count=%d", len(fused))
	}
	if fused[0].DocumentRevisionID != "a" {
		t.Fatalf("first fused document=%s, want a", fused[0].DocumentRevisionID)
	}
	if fused[0].Components["bm25"].WinningChunkID != "a-1" {
		t.Fatalf("winning evidence=%#v", fused[0].Components)
	}
}

func TestFuseWeightedRRFUsesConfiguredChannelWeight(t *testing.T) {
	channels := map[string][]ChunkHit{
		"bm25":   {{Rank: 1, ChunkID: "a-1", DocumentRevisionID: "a", Score: 4}},
		"vector": {{Rank: 1, ChunkID: "b-1", DocumentRevisionID: "b", Score: 0.9}},
	}
	fused := FuseWeightedRRF(channels, 60, map[string]float64{"vector": 2}, 10)
	if len(fused) != 2 || fused[0].DocumentRevisionID != "b" {
		t.Fatalf("weighted fusion = %#v", fused)
	}
	if fused[0].Components["vector"].Contribution != 2.0/61.0 {
		t.Fatalf("weighted contribution = %#v", fused[0].Components)
	}
}
