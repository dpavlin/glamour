package glamour

import (
	"strings"
	"testing"

	styles "charm.land/glamour/v2/styles"
)

func TestMermaidRendering(t *testing.T) {
	md := `
# Title

` + "```mermaid" + `
sequenceDiagram
    participant Alice
    participant Bob
    Alice->>Bob: Hello Bob
    Bob-->>Alice: Hi Alice
` + "```" + `

Some trailing text.
`

	r, err := NewTermRenderer(
		WithStyles(styles.DarkStyleConfig),
		WithWordWrap(80),
	)
	if err != nil {
		t.Fatalf("failed to create term renderer: %v", err)
	}

	out, err := r.Render(md)
	if err != nil {
		t.Fatalf("failed to render markdown with mermaid: %v", err)
	}

	t.Logf("Rendered output:\n%s", out)

	if !strings.Contains(out, "Alice") || !strings.Contains(out, "Bob") {
		t.Errorf("rendered output should contain participants Alice and Bob, got:\n%s", out)
	}

	if !strings.Contains(out, "Hello Bob") {
		t.Errorf("rendered output should contain 'Hello Bob', got:\n%s", out)
	}
}

func TestMermaidFlowchart(t *testing.T) {
	md := `
` + "```mermaid" + `
flowchart LR
    A[Start] --> B[Process]
    B --> C[End]
` + "```" + `
`

	r, err := NewTermRenderer(
		WithStyles(styles.DarkStyleConfig),
		WithWordWrap(80),
	)
	if err != nil {
		t.Fatalf("failed to create term renderer: %v", err)
	}

	out, err := r.Render(md)
	if err != nil {
		t.Fatalf("failed to render markdown with mermaid flowchart: %v", err)
	}

	t.Logf("Rendered flowchart output:\n%s", out)

	if !strings.Contains(out, "Start") || !strings.Contains(out, "Process") || !strings.Contains(out, "End") {
		t.Errorf("rendered output should contain node labels, got:\n%s", out)
	}
}

func TestMermaidKohaSequenceDiagram(t *testing.T) {
	md := `
# AAIEduHR LDAP Expiration

` + "```mermaid" + `
sequenceDiagram
    autonumber
    participant Cron as Cron Daemon (04:00)
    participant Run as run.sh
    participant Dedup as koha-remove-duplicate-borrowers.pl
    participant Sync as sync-dateexpiry-from-ldap.pl
    participant MySQL as MariaDB (koha)
    participant LDAP as LDAPS (ldap.ffzg.hr:636)

    Cron->>Run: Execute run.sh
    Run->>Dedup: Run borrower deduplication
    Dedup->>MySQL: Query duplicate OIB/JMBAG attributes
    Note over Dedup,MySQL: Purges duplicate attribute rows for same borrower.<br/>Merges borrower records if KOHA_UPDATE=1.
    Run->>Sync: Execute sync-dateexpiry-from-ldap.pl
    Sync->>MySQL: Query borrowers WHERE userid LIKE '%@ffzg.hr'
    Sync->>LDAP: Bind using service account (cn=admin,dc=ffzg,dc=hr)
    loop For each borrower (7,300+ patrons)
        Sync->>LDAP: Search HrEduPersonUniqueID = userid
        alt Found in LDAP
            Sync->>MySQL: Insert / update OIB in borrower_attributes
            alt Expiry date changed
                Sync->>MySQL: UPDATE borrowers SET dateexpiry = ldap_date
            else Expiry is NONE
                Sync->>MySQL: UPDATE borrowers SET dateexpiry = '2035-12-31'
            end
        else Not found in LDAP
            Sync->>MySQL: Rename userid & cardnumber to <user>@expired
            Sync->>MySQL: UPDATE borrowers SET dateexpiry = CURRENT_DATE
        end
    end
    Sync-->>Run: Pipe stdout to ldap-dateexpiry.log.YYYY-MM-DD
` + "```" + `
`

	r, err := NewTermRenderer(
		WithStyles(styles.DarkStyleConfig),
		WithWordWrap(80),
	)
	if err != nil {
		t.Fatalf("failed to create term renderer: %v", err)
	}

	out, err := r.Render(md)
	if err != nil {
		t.Fatalf("failed to render koha sequence diagram: %v", err)
	}

	t.Logf("Rendered Koha sequence diagram:\n%s", out)

	if !strings.Contains(out, "Cron Daemon") || !strings.Contains(out, "sync-dateexpiry-from-ldap.pl") {
		t.Errorf("rendered output missing participants: %s", out)
	}
}
