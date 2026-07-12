#!/usr/bin/env python3
"""One-time conversion of index.html portfolio cards into content/projects/*.md
and of index.html itself into templates/home.html."""
import os
import re
import unicodedata

import yaml
from bs4 import BeautifulSoup

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
OUT = os.path.join(ROOT, 'content', 'projects')

SECTIONS = [
    ('Professional Portfolio', 'professional'),
    ('Open Source Projects', 'opensource'),
    ('Unmaintained Open Source Applications', 'unmaintained'),
]


def slugify(title):
    s = unicodedata.normalize('NFKD', title).encode('ascii', 'ignore').decode()
    s = re.sub(r'[^a-zA-Z0-9]+', '-', s).strip('-').lower()
    return s


def main():
    src = open(os.path.join(ROOT, 'index.html')).read()
    soup = BeautifulSoup(src, 'html.parser')

    heading_to_key = dict(SECTIONS)
    weight = 0
    for h3 in soup.find_all('h3'):
        key = heading_to_key.get(h3.get_text(strip=True))
        if not key:
            continue
        grid = h3.find_next_sibling('div')
        for card in grid.find_all('div', recursive=False):
            weight += 10
            img = card.find('img')
            title = card.find('h4').get_text(strip=True)
            desc = re.sub(r'\s+', ' ', card.find('p').get_text()).strip()
            fm = {
                'title': title,
                'section': key,
                'weight': weight,
                'image': img['src'],
                'alt': img.get('alt', ''),
            }
            stack = card.find('div', class_='font-mono')
            if stack:
                fm['stack'] = re.sub(r'\s+', ' ', stack.get_text()).strip().removeprefix('Stack: ')
            link = card.find('a')
            if link is not None:
                fm['link'] = link.get('href', '')
                fm['linkText'] = link.get_text(strip=True)
            front = yaml.safe_dump(fm, allow_unicode=True, sort_keys=False,
                                   default_flow_style=False, width=1000)
            name = slugify(title)
            with open(os.path.join(OUT, name + '.md'), 'w') as f:
                f.write('---\n%s---\n\n%s\n' % (front, desc))
            print('wrote content/projects/%s.md' % name)

    # --- templates/home.html: same page with the portfolio grids templated ---
    start = src.index('      <!-- Professional Portfolio Section -->')
    end = src.index('    </div>\n  </section>\n\n  <!-- Contact Section -->')
    portfolio = '''      {{range .Sections}}<h3 class="text-2xl font-semibold mb-6{{if .ExtraHeadingClass}} {{.ExtraHeadingClass}}{{end}}">{{.Title}}</h3>
      <div class="grid md:grid-cols-2 gap-10{{if .ExtraGridClass}} {{.ExtraGridClass}}{{end}}">
        {{- range .Projects}}
        <div class="bg-[var(--bg-light)] p-6 rounded-lg shadow-lg">
          <div
            class="w-full h-48 mb-4 rounded-lg overflow-hidden bg-white/70 border border-[var(--bg-light)] flex items-center justify-center p-2">
            <img src="{{.Image}}" alt="{{.Alt}}" class="w-full h-full object-contain">
          </div>
          <h4 class="text-2xl font-semibold mb-2">{{.Title}}</h4>
          {{.Body}}
          {{- if .Stack}}
          <div class="text-sm text-[var(--accent)] font-mono">Stack: {{.Stack}}</div>
          {{- end}}
          {{- if .LinkText}}
          <a href="{{.Link}}" class="text-[var(--accent)] hover:underline">{{.LinkText}}</a>
          {{- end}}
        </div>
        {{- end}}
      </div>

      {{end}}'''
    home = src[:start] + portfolio + src[end:]
    with open(os.path.join(ROOT, 'templates', 'home.html'), 'w') as f:
        f.write(home)
    print('wrote templates/home.html')


if __name__ == '__main__':
    main()
