AUTHOR = "Benoît HERVIER"
SITENAME = "Benoît Rvier.fr"
SITEURL = ""

PATH = "content"

TIMEZONE = "Europe/Paris"

DEFAULT_LANG = "en"
LOCALE = "en_US.utf8"
DATE_FORMATS = {
    "en": "%a, %d %b %Y",
}


# Feed generation is usually not desired when developing
FEED_ALL_ATOM = None
CATEGORY_FEED_ATOM = None
TRANSLATION_FEED_ATOM = None
AUTHOR_FEED_ATOM = None
AUTHOR_FEED_RSS = None

# Blogroll
LINKS = (
    ("Pelican", "https://getpelican.com/"),
    ("Python.org", "https://www.python.org/"),
    ("Jinja2", "https://palletsprojects.com/p/jinja/"),
    ("You can modify those links in your config file", "#"),
)

# Social widget
SOCIAL = (
    ("You can add links", "#"),
    ("Another social link", "#"),
)

DEFAULT_PAGINATION = 3

# Uncomment following line if you want document-relative URLs when developing
# RELATIVE_URLS = True


AUTHOR = "Benoît HERVIER"
THEME = "./theme"

# Menu
MENUITEMS = (
    ("posts", "/posts"),
    ("morg", "/pages/morg.html"),
    ("portfolio", "/pages/portfolio.html"),
)
ARCHIVES_SAVE_AS = "posts.html"

DISPLAY_CATEGORIES_ON_MENU = False
DISPLAY_PAGES_ON_MENU = False

ARTICLE_ORDER_BY = "reversed-date"

# Theme customizations
MINIMALXY_CUSTOM_CSS = "static/custom.css"
MINIMALXY_FAVICON = "/theme/images/favicon.ico"
MINIMALXY_START_YEAR = 2022
MINIMALXY_CURRENT_YEAR = 2022

MARKDOWN = {
    'extension_configs': {
        'markdown.extensions.codehilite': {'css_class': 'highlight'},
        'markdown.extensions.extra': {},
        'markdown.extensions.meta': {},
    },
    'output_format': 'html5',
    'tab_length': 2
}

