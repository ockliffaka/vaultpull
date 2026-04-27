package dotenv

import (
	"testing"
)

func TestChain_LaterSourceWins(t *testing.T) {
	opts := ChainOptions{
		Overwrite: true,
		Sources: []ChainSource{
			{Name: "base", Secrets: map[string]string{"KEY": "base_val", "ONLY_BASE": "1"}},
			{Name: "override", Secrets: map[string]string{"KEY": "override_val"}},
		},
	}
	res, err := Chain(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Merged["KEY"] != "override_val" {
		t.Errorf("expected override_val, got %s", res.Merged["KEY"])
	}
	if res.Merged["ONLY_BASE"] != "1" {
		t.Errorf("expected ONLY_BASE to be preserved")
	}
	if res.Provenance["KEY"] != "override" {
		t.Errorf("expected provenance 'override', got %s", res.Provenance["KEY"])
	}
}

func TestChain_FirstSourceWins_NoOverwrite(t *testing.T) {
	opts := ChainOptions{
		Overwrite: false,
		Sources: []ChainSource{
			{Name: "first", Secrets: map[string]string{"KEY": "first_val"}},
			{Name: "second", Secrets: map[string]string{"KEY": "second_val"}},
		},
	}
	res, err := Chain(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Merged["KEY"] != "first_val" {
		t.Errorf("expected first_val, got %s", res.Merged["KEY"])
	}
	if res.Provenance["KEY"] != "first" {
		t.Errorf("expected provenance 'first', got %s", res.Provenance["KEY"])
	}
}

func TestChain_NoSources_ReturnsError(t *testing.T) {
	_, err := Chain(ChainOptions{})
	if err == nil {
		t.Fatal("expected error for empty sources")
	}
}

func TestChain_UnnamedSource_ReturnsError(t *testing.T) {
	opts := ChainOptions{
		Overwrite: true,
		Sources: []ChainSource{
			{Name: "", Secrets: map[string]string{"K": "v"}},
		},
	}
	_, err := Chain(opts)
	if err == nil {
		t.Fatal("expected error for unnamed source")
	}
}

func TestChain_ProvenanceSummary_Sorted(t *testing.T) {
	opts := DefaultChainOptions()
	opts.Sources = []ChainSource{
		{Name: "env", Secrets: map[string]string{"ZEBRA": "z", "ALPHA": "a"}},
	}
	res, err := Chain(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	summary := res.ProvenanceSummary()
	if len(summary) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(summary))
	}
	if summary[0] != "ALPHA <- env" {
		t.Errorf("unexpected first line: %s", summary[0])
	}
	if summary[1] != "ZEBRA <- env" {
		t.Errorf("unexpected second line: %s", summary[1])
	}
}
