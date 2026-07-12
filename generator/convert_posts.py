#!/usr/bin/env python3
"""One-time conversion of existing posts/*.html into content/posts/*.md sources."""
import glob
import json
import os
import re
import sys

import yaml
from bs4 import BeautifulSoup, NavigableString, Tag

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
POSTS = os.path.join(ROOT, 'posts')
OUT = os.path.join(ROOT, 'content', 'posts')


def collapse(text):
    return re.sub(r'\s+', ' ', text)


def inline(node):
    """Convert an inline-level node tree to markdown text."""
    parts = []
    for child in node.children:
        if isinstance(child, NavigableString):
            parts.append(collapse(str(child)))
        elif isinstance(child, Tag):
            if child.name == 'code':
                parts.append('`%s`' % child.get_text())
            elif child.name == 'a':
                parts.append('[%s](%s)' % (inline(child).strip(), child.get('href', '')))
            elif child.name == 'strong':
                parts.append('**%s**' % inline(child).strip())
            elif child.name == 'em':
                parts.append('*%s*' % inline(child).strip())
            elif child.name == 'br':
                parts.append('  \n')
            elif child.name == 'img':
                parts.append(img_html(child))
            else:
                sys.exit('unhandled inline tag: %s' % child.name)
    return ''.join(parts)


def img_html(tag):
    attrs = ['src="%s"' % tag['src'], 'alt="%s"' % tag.get('alt', '')]
    if tag.get('loading'):
        attrs.append('loading="%s"' % tag['loading'])
    if tag.get('width'):
        attrs.append('width="%s"' % tag['width'])
    return '<img %s>' % ' '.join(attrs)


def convert_list(ul, depth=0):
    lines = []
    for li in ul.find_all('li', recursive=False):
        # inline content of the li, excluding nested lists
        chunks = []
        nested = []
        for child in li.children:
            if isinstance(child, Tag) and child.name in ('ul', 'ol'):
                nested.append(child)
            elif isinstance(child, NavigableString):
                chunks.append(collapse(str(child)))
            else:
                tmp = BeautifulSoup('<span></span>', 'html.parser').span
                tmp.append(child.__copy__())
                chunks.append(inline(tmp))
        text = ''.join(chunks).strip()
        lines.append('%s- %s' % ('    ' * depth, text))
        for n in nested:
            lines.extend(convert_list(n, depth + 1))
    return lines


def convert_article(article):
    out = []
    for child in article.children:
        if isinstance(child, NavigableString):
            if child.strip():
                sys.exit('stray text in article: %r' % child)
            continue
        name = child.name
        if name == 'h1':
            continue
        if name == 'p':
            klass = child.get('class') or []
            if 'text-sm' in klass:  # date line
                continue
            # paragraph that only wraps an image
            tags = [c for c in child.children if isinstance(c, Tag)]
            if len(tags) == 1 and tags[0].name == 'img' and not child.get_text(strip=True):
                out.append(img_html(tags[0]))
            else:
                out.append(inline(child).strip())
        elif name in ('h2', 'h3', 'h4'):
            out.append('%s %s' % ('#' * int(name[1]), inline(child).strip()))
        elif name == 'pre':
            code = child.find('code')
            text = (code or child).get_text()
            out.append('```\n%s\n```' % text.rstrip('\n'))
        elif name in ('ul', 'ol'):
            out.append('\n'.join(convert_list(child)))
        else:
            sys.exit('unhandled block tag: %s' % name)
    return '\n\n'.join(out) + '\n'


def summaries():
    soup = BeautifulSoup(open(os.path.join(POSTS, 'index.html')).read(), 'html.parser')
    result = {}
    for art in soup.find_all('article'):
        a = art.find('a')
        p = art.find_all('p')[-1]
        result[a['href']] = collapse(p.get_text()).strip()
    return result


def meta_content(soup, **attrs):
    tag = soup.find('meta', attrs=attrs)
    return tag['content'] if tag else None


def main():
    cards = summaries()
    for path in sorted(glob.glob(os.path.join(POSTS, '*.html'))):
        base = os.path.basename(path)
        if base == 'index.html':
            continue
        slug = base[:-5]
        soup = BeautifulSoup(open(path).read(), 'html.parser')
        ld = json.loads(soup.find('script', type='application/ld+json').string)

        fm = {
            'title': ld['headline'],
            'date': ld['datePublished'],
            'lang': ld.get('inLanguage', 'en'),
            'description': meta_content(soup, name='description'),
        }
        if ld.get('dateModified'):
            fm['updated'] = ld['dateModified']
        og_desc = meta_content(soup, property='og:description')
        if og_desc and og_desc != fm['description']:
            fm['ogDescription'] = og_desc
        if ld.get('image'):
            fm['image'] = ld['image']
        if ld.get('keywords'):
            fm['keywords'] = ld['keywords']
        summary = cards.get(base)
        if summary and summary != fm['description']:
            fm['summary'] = summary

        body = convert_article(soup.find('article'))
        front = yaml.safe_dump(fm, allow_unicode=True, sort_keys=False,
                               default_flow_style=False, width=1000)
        with open(os.path.join(OUT, slug + '.md'), 'w') as f:
            f.write('---\n%s---\n\n%s' % (front, body))
        print('wrote content/posts/%s.md' % slug)


if __name__ == '__main__':
    main()
