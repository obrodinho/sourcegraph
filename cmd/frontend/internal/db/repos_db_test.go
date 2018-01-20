package db

import (
	"reflect"
	"sort"
	"strings"
	"testing"

	"context"

	"sourcegraph.com/sourcegraph/sourcegraph/cmd/frontend/internal/pkg/types"
	"sourcegraph.com/sourcegraph/sourcegraph/pkg/actor"
)

/*
 * Helpers
 */

func sortedRepoURIs(repos []*types.Repo) []string {
	uris := repoURIs(repos)
	sort.Strings(uris)
	return uris
}

func repoURIs(repos []*types.Repo) []string {
	var uris []string
	for _, repo := range repos {
		uris = append(uris, repo.URI)
	}
	return uris
}

func createRepo(ctx context.Context, t *testing.T, repo *types.Repo) {
	if err := Repos.TryInsertNew(ctx, repo.URI, repo.Description, repo.Fork, true); err != nil {
		t.Fatal(err)
	}
}

func mustCreate(ctx context.Context, t *testing.T, repos ...*types.Repo) []*types.Repo {
	var createdRepos []*types.Repo
	for _, repo := range repos {
		createRepo(ctx, t, repo)
		repo, err := Repos.GetByURI(ctx, repo.URI)
		if err != nil {
			t.Fatal(err)
		}
		createdRepos = append(createdRepos, repo)
	}
	return createdRepos
}

/*
 * Tests
 */

func TestRepos_Get(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := testContext()

	want := mustCreate(ctx, t, &types.Repo{URI: "r"})

	repo, err := Repos.Get(ctx, want[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	if !jsonEqual(t, repo, want[0]) {
		t.Errorf("got %v, want %v", repo, want[0])
	}
}

func TestRepos_List(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := testContext()

	ctx = actor.WithActor(ctx, &actor.Actor{})

	want := mustCreate(ctx, t, &types.Repo{URI: "r"})

	repos, err := Repos.List(ctx, ReposListOptions{Enabled: true})
	if err != nil {
		t.Fatal(err)
	}
	if !jsonEqual(t, repos, want) {
		t.Errorf("got %v, want %v", repos, want)
	}
}

func TestRepos_List_pagination(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := testContext()

	ctx = actor.WithActor(ctx, &actor.Actor{})

	createdRepos := []*types.Repo{
		{URI: "r1"},
		{URI: "r2"},
		{URI: "r3"},
	}
	for _, repo := range createdRepos {
		mustCreate(ctx, t, repo)
	}

	type testcase struct {
		limit  int
		offset int
		exp    []string
	}
	tests := []testcase{
		{limit: 1, offset: 0, exp: []string{"r1"}},
		{limit: 1, offset: 1, exp: []string{"r2"}},
		{limit: 1, offset: 2, exp: []string{"r3"}},
		{limit: 2, offset: 0, exp: []string{"r1", "r2"}},
		{limit: 2, offset: 2, exp: []string{"r3"}},
		{limit: 3, offset: 0, exp: []string{"r1", "r2", "r3"}},
		{limit: 3, offset: 3, exp: nil},
		{limit: 4, offset: 0, exp: []string{"r1", "r2", "r3"}},
		{limit: 4, offset: 4, exp: nil},
	}
	for _, test := range tests {
		repos, err := Repos.List(ctx, ReposListOptions{Enabled: true, LimitOffset: &LimitOffset{Limit: test.limit, Offset: test.offset}})
		if err != nil {
			t.Fatal(err)
		}
		if got := sortedRepoURIs(repos); !reflect.DeepEqual(got, test.exp) {
			t.Errorf("for test case %v, got %v (want %v)", test, got, test.exp)
		}
	}
}

// TestRepos_List_query tests the behavior of Repos.List when called with
// a query.
// Test batch 1 (correct filtering)
func TestRepos_List_query1(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := testContext()

	ctx = actor.WithActor(ctx, &actor.Actor{})

	createdRepos := []*types.Repo{
		{URI: "abc/def"},
		{URI: "def/ghi"},
		{URI: "jkl/mno/pqr"},
		{URI: "github.com/abc/xyz"},
	}
	for _, repo := range createdRepos {
		createRepo(ctx, t, repo)
	}
	tests := []struct {
		query string
		want  []string
	}{
		{"def", []string{"abc/def", "def/ghi"}},
		{"ABC/DEF", []string{"abc/def"}},
		{"xyz", []string{"github.com/abc/xyz"}},
		{"mno/p", []string{"jkl/mno/pqr"}},
		{"jkl mno pqr", []string{"jkl/mno/pqr"}},
	}
	for _, test := range tests {
		repos, err := Repos.List(ctx, ReposListOptions{Query: test.query, Enabled: true})
		if err != nil {
			t.Fatal(err)
		}
		if got := repoURIs(repos); !reflect.DeepEqual(got, test.want) {
			t.Errorf("%q: got repos %q, want %q", test.query, got, test.want)
		}
	}
}

// Test batch 2 (correct ranking)
func TestRepos_List_query2(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := testContext()

	ctx = actor.WithActor(ctx, &actor.Actor{})

	createdRepos := []*types.Repo{
		{URI: "a/def"},
		{URI: "b/def"},
		{URI: "c/def"},
		{URI: "def/ghi"},
		{URI: "def/jkl"},
		{URI: "def/mno"},
		{URI: "abc/m"},
	}
	for _, repo := range createdRepos {
		createRepo(ctx, t, repo)
	}
	tests := []struct {
		query string
		want  []string
	}{
		{"def", []string{"a/def", "b/def", "c/def", "def/ghi", "def/jkl", "def/mno"}},
		{"b/def", []string{"b/def"}},
		{"def/", []string{"def/ghi", "def/jkl", "def/mno"}},
		{"def/m", []string{"def/mno"}},
	}
	for _, test := range tests {
		repos, err := Repos.List(ctx, ReposListOptions{Query: test.query, Enabled: true})
		if err != nil {
			t.Fatal(err)
		}
		if got := repoURIs(repos); !reflect.DeepEqual(got, test.want) {
			t.Errorf("Unexpected repo result for query %q:\ngot:  %q\nwant: %q", test.query, got, test.want)
		}
	}
}

// TestRepos_List_patterns tests the behavior of Repos.List when called with
// IncludePatterns and ExcludePattern.
func TestRepos_List_patterns(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := testContext()

	ctx = actor.WithActor(ctx, &actor.Actor{})

	createdRepos := []*types.Repo{
		{URI: "a/b"},
		{URI: "c/d"},
		{URI: "e/f"},
		{URI: "g/h"},
	}
	for _, repo := range createdRepos {
		createRepo(ctx, t, repo)
	}
	tests := []struct {
		includePatterns []string
		excludePattern  string
		want            []string
	}{
		{
			includePatterns: []string{"(a|c)"},
			want:            []string{"a/b", "c/d"},
		},
		{
			includePatterns: []string{"(a|c)", "b"},
			want:            []string{"a/b"},
		},
		{
			includePatterns: []string{"(a|c)"},
			excludePattern:  "d",
			want:            []string{"a/b"},
		},
		{
			excludePattern: "(d|e)",
			want:           []string{"a/b", "g/h"},
		},
	}
	for _, test := range tests {
		repos, err := Repos.List(ctx, ReposListOptions{
			IncludePatterns: test.includePatterns,
			ExcludePattern:  test.excludePattern,
			Enabled:         true,
		})
		if err != nil {
			t.Fatal(err)
		}
		if got := repoURIs(repos); !reflect.DeepEqual(got, test.want) {
			t.Errorf("include %q exclude %q: got repos %q, want %q", test.includePatterns, test.excludePattern, got, test.want)
		}
	}
}

func TestRepos_List_queryAndPatternsMutuallyExclusive(t *testing.T) {
	ctx := context.Background()
	wantErr := "Query and IncludePatterns/ExcludePattern options are mutually exclusive"

	t.Run("Query and IncludePatterns", func(t *testing.T) {
		_, err := Repos.List(ctx, ReposListOptions{Query: "x", IncludePatterns: []string{"y"}, Enabled: true})
		if err == nil || !strings.Contains(err.Error(), wantErr) {
			t.Fatalf("got error %v, want it to contain %q", err, wantErr)
		}
	})

	t.Run("Query and ExcludePattern", func(t *testing.T) {
		_, err := Repos.List(ctx, ReposListOptions{Query: "x", ExcludePattern: "y", Enabled: true})
		if err == nil || !strings.Contains(err.Error(), wantErr) {
			t.Fatalf("got error %v, want it to contain %q", err, wantErr)
		}
	})
}

func TestRepos_Create(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := testContext()

	// Add a repo.
	createRepo(ctx, t, &types.Repo{URI: "a/b"})

	repo, err := Repos.GetByURI(ctx, "a/b")
	if err != nil {
		t.Fatal(err)
	}
	if repo.CreatedAt.IsZero() {
		t.Fatal("got CreatedAt.IsZero()")
	}
}

func TestRepos_Create_dupe(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := testContext()

	// Add a repo.
	createRepo(ctx, t, &types.Repo{URI: "a/b"})

	// Add another repo with the same name.
	createRepo(ctx, t, &types.Repo{URI: "a/b"})
}

func TestMakeFuzzyLikeRepoQuery(t *testing.T) {
	cases := map[string]string{
		"":           "%",
		"/":          "%/%",
		"foo":        "%foo%",
		"/foo":       "%/%foo%",
		"foo/":       "%foo%/%",
		"//foo":      "%/%/%foo%",
		"foo//":      "%foo%/%/%",
		"/foo/":      "%/%foo%/%",
		"foo/bar":    "%foo%/%bar%",
		"/foo/bar/":  "%/%foo%/%bar%/%",
		"/foo%/bar/": "%/%foo%\\%%/%bar%/%",
	}
	for query, want := range cases {
		got := makeFuzzyLikeRepoQuery(query)
		if want != got {
			t.Errorf("makeFuzzyLikeQuery(%q) == %q != %q", query, got, want)
		}
	}
}
