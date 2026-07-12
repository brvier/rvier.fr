// Command generator builds the rvier.fr website into public/:
//
//	static/*              -> public/* (copied verbatim: images, css, robots.txt, redirects)
//	content/posts/*.md    -> public/posts/<slug>.html, public/posts/index.html, public/sitemap.xml
//	content/projects/*.md -> public/index.html (portfolio sections of templates/home.html)
//
// Run from the repository root: go run ./generator
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v3"
)

const (
	staticDir = "static"
	outDir    = "public"
)

var markdown = goldmark.New(
	goldmark.WithRendererOptions(html.WithUnsafe()), // posts embed raw <img> tags
)

type Post struct {
	Title         string `yaml:"title"`
	Date          string `yaml:"date"`
	Updated       string `yaml:"updated"`
	Lang          string `yaml:"lang"`
	Description   string `yaml:"description"`
	OgDescription string `yaml:"ogDescription"`
	Image         string `yaml:"image"`
	Keywords      string `yaml:"keywords"`
	Summary       string `yaml:"summary"`

	Slug string        `yaml:"-"`
	Body template.HTML `yaml:"-"`
}

func (p Post) URL() string            { return "https://rvier.fr/posts/" + p.Slug + ".html" }
func (p Post) LangTag() string        { return strings.ToUpper(p.Lang) }
func (p Post) DisplayDate() string    { return displayDate(p.Date, p.Lang) }
func (p Post) DisplayUpdated() string { return displayDate(p.Updated, p.Lang) }
func (p Post) LastMod() string {
	if p.Updated != "" {
		return p.Updated
	}
	return p.Date
}

func (p Post) JSONLD() (template.JS, error) {
	type author struct {
		Type string `json:"@type"`
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	ld := struct {
		Context       string `json:"@context"`
		Type          string `json:"@type"`
		Headline      string `json:"headline"`
		DatePublished string `json:"datePublished"`
		DateModified  string `json:"dateModified,omitempty"`
		InLanguage    string `json:"inLanguage"`
		Keywords      string `json:"keywords,omitempty"`
		Image         string `json:"image,omitempty"`
		Description   string `json:"description"`
		MainEntity    string `json:"mainEntityOfPage"`
		Author        author `json:"author"`
	}{
		Context:       "https://schema.org",
		Type:          "BlogPosting",
		Headline:      p.Title,
		DatePublished: p.Date,
		DateModified:  p.Updated,
		InLanguage:    p.Lang,
		Keywords:      p.Keywords,
		Image:         p.Image,
		Description:   p.Description,
		MainEntity:    p.URL(),
		Author:        author{Type: "Person", Name: "Benoît Hervier", URL: "https://rvier.fr/"},
	}
	b, err := json.MarshalIndent(ld, "  ", "  ")
	return template.JS(b), err
}

type Project struct {
	Title    string `yaml:"title"`
	Section  string `yaml:"section"`
	Weight   int    `yaml:"weight"`
	Image    string `yaml:"image"`
	Alt      string `yaml:"alt"`
	Stack    string `yaml:"stack"`
	Link     string `yaml:"link"`
	LinkText string `yaml:"linkText"`

	Body template.HTML `yaml:"-"`
}

type Section struct {
	Key               string
	Title             string
	ExtraHeadingClass string
	ExtraGridClass    string
	Projects          []Project
}

var sections = []Section{
	{Key: "professional", Title: "Professional Portfolio", ExtraGridClass: "mb-16"},
	{Key: "opensource", Title: "Open Source Projects", ExtraGridClass: "mb-16"},
	{Key: "unmaintained", Title: "Unmaintained Open Source Applications", ExtraHeadingClass: "mt-12"},
}

// splitFrontMatter separates the YAML front matter from the markdown body.
func splitFrontMatter(raw []byte) (fm, body []byte, err error) {
	const sep = "---\n"
	s := string(raw)
	if !strings.HasPrefix(s, sep) {
		return nil, nil, fmt.Errorf("missing front matter")
	}
	rest := s[len(sep):]
	i := strings.Index(rest, "\n"+sep)
	if i < 0 {
		return nil, nil, fmt.Errorf("unterminated front matter")
	}
	return []byte(rest[:i+1]), []byte(rest[i+1+len(sep):]), nil
}

func render(md []byte) (template.HTML, error) {
	var buf bytes.Buffer
	if err := markdown.Convert(md, &buf); err != nil {
		return "", err
	}
	return template.HTML(strings.TrimRight(buf.String(), "\n")), nil
}

var months = map[string][]string{
	"en": {"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December"},
	"fr": {"janvier", "février", "mars", "avril", "mai", "juin",
		"juillet", "août", "septembre", "octobre", "novembre", "décembre"},
}

// displayDate formats an ISO date per language: "July 12, 2026" / "12 juillet 2026".
func displayDate(iso, lang string) string {
	var y, m, d int
	if _, err := fmt.Sscanf(iso, "%d-%d-%d", &y, &m, &d); err != nil {
		log.Fatalf("invalid date %q: %v", iso, err)
	}
	names, ok := months[lang]
	if !ok {
		names = months["en"]
	}
	if lang == "fr" {
		return fmt.Sprintf("%d %s %d", d, names[m-1], y)
	}
	return fmt.Sprintf("%s %d, %d", names[m-1], d, y)
}

func loadPosts() ([]Post, error) {
	files, err := filepath.Glob("content/posts/*.md")
	if err != nil {
		return nil, err
	}
	var posts []Post
	for _, f := range files {
		raw, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		fm, body, err := splitFrontMatter(raw)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", f, err)
		}
		var p Post
		if err := yaml.Unmarshal(fm, &p); err != nil {
			return nil, fmt.Errorf("%s: %w", f, err)
		}
		if p.Title == "" || p.Date == "" || p.Description == "" {
			return nil, fmt.Errorf("%s: title, date and description are required", f)
		}
		if p.Lang == "" {
			p.Lang = "en"
		}
		if p.OgDescription == "" {
			p.OgDescription = p.Description
		}
		if p.Summary == "" {
			p.Summary = p.Description
		}
		p.Slug = strings.TrimSuffix(filepath.Base(f), ".md")
		if p.Body, err = render(body); err != nil {
			return nil, fmt.Errorf("%s: %w", f, err)
		}
		posts = append(posts, p)
	}
	sort.Slice(posts, func(i, j int) bool {
		if posts[i].Date != posts[j].Date {
			return posts[i].Date > posts[j].Date
		}
		return posts[i].Slug < posts[j].Slug
	})
	return posts, nil
}

func loadSections() ([]Section, error) {
	files, err := filepath.Glob("content/projects/*.md")
	if err != nil {
		return nil, err
	}
	byKey := map[string][]Project{}
	for _, f := range files {
		raw, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		fm, body, err := splitFrontMatter(raw)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", f, err)
		}
		var p Project
		if err := yaml.Unmarshal(fm, &p); err != nil {
			return nil, fmt.Errorf("%s: %w", f, err)
		}
		html, err := render(body)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", f, err)
		}
		p.Body = template.HTML(strings.ReplaceAll(string(html), "<p>", `<p class="mb-4">`))
		if _, ok := byKey[p.Section]; !ok {
			known := false
			for _, s := range sections {
				known = known || s.Key == p.Section
			}
			if !known {
				return nil, fmt.Errorf("%s: unknown section %q", f, p.Section)
			}
		}
		byKey[p.Section] = append(byKey[p.Section], p)
	}
	out := make([]Section, 0, len(sections))
	for _, s := range sections {
		s.Projects = byKey[s.Key]
		sort.Slice(s.Projects, func(i, j int) bool { return s.Projects[i].Weight < s.Projects[j].Weight })
		if len(s.Projects) > 0 {
			out = append(out, s)
		}
	}
	return out, nil
}

func writeSitemap(posts []Post) error {
	latest := ""
	for _, p := range posts {
		if p.LastMod() > latest {
			latest = p.LastMod()
		}
	}
	var b strings.Builder
	b.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	b.WriteString("<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n")
	entry := func(loc, lastmod string) {
		fmt.Fprintf(&b, "  <url>\n    <loc>%s</loc>\n    <lastmod>%s</lastmod>\n  </url>\n", loc, lastmod)
	}
	entry("https://rvier.fr/", latest)
	entry("https://rvier.fr/posts/index.html", latest)
	for _, p := range posts {
		entry(p.URL(), p.LastMod())
	}
	b.WriteString("</urlset>\n")
	return os.WriteFile(filepath.Join(outDir, "sitemap.xml"), []byte(b.String()), 0o644)
}

// copyStatic mirrors static/ into public/.
func copyStatic() error {
	return filepath.WalkDir(staticDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(staticDir, path)
		if err != nil {
			return err
		}
		dst := filepath.Join(outDir, rel)
		if d.IsDir() {
			return os.MkdirAll(dst, 0o755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(dst, data, 0o644)
	})
}

func renderToFile(t *template.Template, path string, data any) error {
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	return os.WriteFile(path, buf.Bytes(), 0o644)
}

func main() {
	posts, err := loadPosts()
	if err != nil {
		log.Fatal(err)
	}
	secs, err := loadSections()
	if err != nil {
		log.Fatal(err)
	}

	postTpl := template.Must(template.ParseFiles("templates/post.html", "templates/partials.html"))
	indexTpl := template.Must(template.ParseFiles("templates/blogindex.html", "templates/partials.html"))
	homeTpl := template.Must(template.ParseFiles("templates/home.html"))

	if err := os.RemoveAll(outDir); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(outDir, "posts"), 0o755); err != nil {
		log.Fatal(err)
	}
	if err := copyStatic(); err != nil {
		log.Fatal(err)
	}

	for _, p := range posts {
		if err := renderToFile(postTpl, filepath.Join(outDir, "posts", p.Slug+".html"), p); err != nil {
			log.Fatal(err)
		}
	}
	if err := renderToFile(indexTpl, filepath.Join(outDir, "posts", "index.html"), map[string]any{"Posts": posts}); err != nil {
		log.Fatal(err)
	}
	if err := renderToFile(homeTpl, filepath.Join(outDir, "index.html"), map[string]any{"Sections": secs}); err != nil {
		log.Fatal(err)
	}
	if err := writeSitemap(posts); err != nil {
		log.Fatal(err)
	}
	log.Printf("generated %d posts, blog index, homepage (%d projects), sitemap",
		len(posts), func() (n int) {
			for _, s := range secs {
				n += len(s.Projects)
			}
			return
		}())
}
